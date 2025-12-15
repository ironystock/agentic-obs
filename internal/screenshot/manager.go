package screenshot

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ironystock/agentic-obs/internal/obs"
	"github.com/ironystock/agentic-obs/internal/storage"
)

// OBSScreenshotter defines the interface for taking screenshots.
// This allows the manager to work with both real and mock OBS clients.
type OBSScreenshotter interface {
	TakeSourceScreenshot(opts obs.ScreenshotOptions) (string, error)
	IsConnected() bool
}

// Config holds screenshot manager configuration options.
//
// Storage Considerations: Each screenshot is stored as base64-encoded image data in SQLite.
// A typical 1080p PNG screenshot is ~1-2MB, resulting in ~1.3-2.6MB of base64 data per image.
// With the default of 10 screenshots per source, expect ~13-26MB of storage per source.
// Adjust MaxScreenshotsPerSource based on available disk space and number of sources.
type Config struct {
	// MaxScreenshotsPerSource is the maximum number of screenshots to keep per source.
	// Older screenshots are automatically cleaned up. Default: 10
	MaxScreenshotsPerSource int

	// CleanupInterval is how often to run the cleanup routine. Default: 1 minute
	CleanupInterval time.Duration
}

// DefaultConfig returns the default manager configuration.
func DefaultConfig() Config {
	return Config{
		MaxScreenshotsPerSource: 10,
		CleanupInterval:         time.Minute,
	}
}

// Manager coordinates periodic screenshot capture from OBS sources.
type Manager struct {
	obsClient OBSScreenshotter
	storage   *storage.DB
	cfg       Config

	mu      sync.RWMutex
	workers map[int64]*worker
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	running bool
}

// NewManager creates a new screenshot manager.
func NewManager(obsClient OBSScreenshotter, db *storage.DB, cfg Config) *Manager {
	if cfg.MaxScreenshotsPerSource <= 0 {
		cfg.MaxScreenshotsPerSource = 10
	}
	if cfg.CleanupInterval <= 0 {
		cfg.CleanupInterval = time.Minute
	}

	return &Manager{
		obsClient: obsClient,
		storage:   db,
		cfg:       cfg,
		workers:   make(map[int64]*worker),
	}
}

// Start initializes the manager, loads existing sources, and starts capture workers.
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("screenshot manager already running")
	}

	m.ctx, m.cancel = context.WithCancel(ctx)
	m.running = true

	// Load all enabled screenshot sources and start workers
	sources, err := m.storage.ListScreenshotSources(m.ctx)
	if err != nil {
		m.running = false
		m.cancel()
		return fmt.Errorf("failed to load screenshot sources: %w", err)
	}

	for _, source := range sources {
		if source.Enabled {
			if err := m.startWorkerLocked(source); err != nil {
				// Log but continue - don't fail entirely for one bad source
				continue
			}
		}
	}

	// Start cleanup goroutine
	m.wg.Add(1)
	go m.cleanupLoop()

	return nil
}

// Stop gracefully shuts down all capture workers.
func (m *Manager) Stop() {
	m.mu.Lock()
	if !m.running {
		m.mu.Unlock()
		return
	}
	m.running = false
	m.cancel()
	m.mu.Unlock()

	// Wait for all workers and cleanup goroutine to finish
	m.wg.Wait()

	// Clear workers map
	m.mu.Lock()
	m.workers = make(map[int64]*worker)
	m.mu.Unlock()
}

// AddSource adds a new screenshot source and starts capturing if enabled.
func (m *Manager) AddSource(source *storage.ScreenshotSource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return fmt.Errorf("screenshot manager not running")
	}

	if _, exists := m.workers[source.ID]; exists {
		return fmt.Errorf("source %d already has a running worker", source.ID)
	}

	if source.Enabled {
		return m.startWorkerLocked(source)
	}

	return nil
}

// RemoveSource stops and removes a screenshot source.
func (m *Manager) RemoveSource(id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if w, exists := m.workers[id]; exists {
		w.stop()
		delete(m.workers, id)
	}

	return nil
}

// UpdateSource updates a source's configuration and restarts its worker if needed.
func (m *Manager) UpdateSource(source *storage.ScreenshotSource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return fmt.Errorf("screenshot manager not running")
	}

	// Stop existing worker if any
	if w, exists := m.workers[source.ID]; exists {
		w.stop()
		delete(m.workers, source.ID)
	}

	// Start new worker if source is enabled
	if source.Enabled {
		return m.startWorkerLocked(source)
	}

	return nil
}

// UpdateCadence updates the capture interval for a source.
func (m *Manager) UpdateCadence(id int64, cadenceMs int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return fmt.Errorf("screenshot manager not running")
	}

	w, exists := m.workers[id]
	if !exists {
		return fmt.Errorf("no worker found for source %d", id)
	}

	w.updateCadence(time.Duration(cadenceMs) * time.Millisecond)
	return nil
}

// GetWorkerCount returns the number of active workers.
func (m *Manager) GetWorkerCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.workers)
}

// IsRunning returns whether the manager is currently running.
func (m *Manager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// startWorkerLocked starts a worker for the given source.
// Must be called with m.mu held.
func (m *Manager) startWorkerLocked(source *storage.ScreenshotSource) error {
	w := newWorker(m.ctx, m.obsClient, m.storage, source)
	m.workers[source.ID] = w
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		w.run()
	}()
	return nil
}

// cleanupLoop periodically removes old screenshots to keep storage bounded.
func (m *Manager) cleanupLoop() {
	defer m.wg.Done()

	ticker := time.NewTicker(m.cfg.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.runCleanup()
		}
	}
}

// runCleanup removes old screenshots for all sources.
func (m *Manager) runCleanup() {
	sources, err := m.storage.ListScreenshotSources(m.ctx)
	if err != nil {
		return
	}

	for _, source := range sources {
		m.storage.DeleteOldScreenshots(m.ctx, source.ID, m.cfg.MaxScreenshotsPerSource)
	}
}

// worker handles periodic capture for a single screenshot source.
type worker struct {
	ctx       context.Context
	cancel    context.CancelFunc
	obsClient OBSScreenshotter
	storage   *storage.DB
	source    *storage.ScreenshotSource

	mu      sync.RWMutex
	cadence time.Duration
}

// newWorker creates a new capture worker for a source.
func newWorker(parentCtx context.Context, obs OBSScreenshotter, db *storage.DB, source *storage.ScreenshotSource) *worker {
	ctx, cancel := context.WithCancel(parentCtx)
	return &worker{
		ctx:       ctx,
		cancel:    cancel,
		obsClient: obs,
		storage:   db,
		source:    source,
		cadence:   time.Duration(source.CadenceMs) * time.Millisecond,
	}
}

// run starts the capture loop.
func (w *worker) run() {
	// Capture immediately on start
	w.capture()

	ticker := time.NewTicker(w.getCadence())
	defer ticker.Stop()

	lastCadence := w.getCadence()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.capture()
			// Check if cadence changed and reset ticker if needed
			currentCadence := w.getCadence()
			if currentCadence != lastCadence {
				ticker.Reset(currentCadence)
				lastCadence = currentCadence
			}
		}
	}
}

// stop cancels the worker's context.
func (w *worker) stop() {
	w.cancel()
}

// updateCadence changes the capture interval.
func (w *worker) updateCadence(cadence time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.cadence = cadence
}

// getCadence returns the current capture interval.
func (w *worker) getCadence() time.Duration {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.cadence
}

// capture takes a screenshot and saves it to storage.
func (w *worker) capture() {
	if !w.obsClient.IsConnected() {
		return
	}

	opts := obs.ScreenshotOptions{
		SourceName: w.source.SourceName,
		Format:     w.source.ImageFormat,
		Width:      w.source.ImageWidth,
		Height:     w.source.ImageHeight,
		Quality:    w.source.Quality,
	}

	imageData, err := w.obsClient.TakeSourceScreenshot(opts)
	if err != nil {
		log.Printf("Screenshot capture failed for source %q: %v", w.source.Name, err)
		return
	}

	// Determine MIME type
	mimeType := "image/png"
	if w.source.ImageFormat == "jpg" || w.source.ImageFormat == "jpeg" {
		mimeType = "image/jpeg"
	}

	// Calculate approximate size (base64 is ~4/3 of original)
	sizeBytes := len(imageData) * 3 / 4

	screenshot := storage.Screenshot{
		SourceID:  w.source.ID,
		ImageData: imageData,
		MimeType:  mimeType,
		SizeBytes: sizeBytes,
	}

	if _, err := w.storage.SaveScreenshot(w.ctx, screenshot); err != nil {
		log.Printf("Failed to save screenshot for source %q: %v", w.source.Name, err)
	}
}
