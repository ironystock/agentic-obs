package http

import (
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ironystock/agentic-obs/internal/storage"
)

//go:embed static/*
var staticFiles embed.FS

// Config holds HTTP server configuration options.
//
// Security Note: This server is designed for local use (localhost binding by default).
// If exposed externally, consider adding rate limiting and authentication at the
// reverse proxy or load balancer level.
type Config struct {
	// Host to bind to (default: "localhost")
	Host string
	// Port to listen on (default: 8765)
	Port int
	// ThumbnailCacheSec is the Cache-Control max-age for thumbnails (0 to disable)
	ThumbnailCacheSec int
}

// isValidSourceName validates that a source name is safe to use.
// Returns false if the name contains path traversal patterns, path separators,
// null bytes, or other potentially dangerous characters.
func isValidSourceName(name string) bool {
	// Reject empty names
	if name == "" {
		return false
	}

	// Reject names containing null bytes
	if strings.ContainsRune(name, '\x00') {
		return false
	}

	// Reject path separators (both Unix and Windows style)
	if strings.ContainsAny(name, "/\\") {
		return false
	}

	// Reject parent directory patterns
	if strings.Contains(name, "..") {
		return false
	}

	// Reject names that are just dots
	if name == "." {
		return false
	}

	return true
}

// DefaultConfig returns the default HTTP server configuration.
func DefaultConfig() Config {
	return Config{
		Host:              "localhost",
		Port:              8765,
		ThumbnailCacheSec: 5, // 5 seconds default, 0 for development
	}
}

// Server provides HTTP endpoints for serving screenshot images and the web dashboard.
type Server struct {
	storage        *storage.DB
	httpServer     *http.Server
	addr           string
	cfg            Config
	startTime      time.Time
	statusProvider StatusProvider
	uiHandlers     *UIHandlers

	mu       sync.RWMutex
	running  bool
	listener net.Listener
}

// NewServer creates a new HTTP server for serving screenshots.
// Use SetStatusProvider to enable MCP-UI endpoints.
func NewServer(db *storage.DB, cfg Config) *Server {
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == 0 {
		cfg.Port = 8765
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	s := &Server{
		storage:   db,
		cfg:       cfg,
		addr:      addr,
		startTime: time.Now(),
	}

	return s
}

// SetStatusProvider configures the status provider for MCP-UI endpoints.
// This must be called before Start() to enable UI routes.
// Returns an error if the server is already running.
func (s *Server) SetStatusProvider(provider StatusProvider) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("cannot set status provider: server already running")
	}

	s.statusProvider = provider
	return nil
}

// setupRoutes configures all HTTP routes.
func (s *Server) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Screenshot and health endpoints
	mux.HandleFunc("/screenshot/", s.handleScreenshot)
	mux.HandleFunc("/health", s.handleHealth)

	// API endpoints
	mux.HandleFunc("/api/status", s.handleAPIStatus)
	mux.HandleFunc("/api/history", s.handleAPIHistory)
	mux.HandleFunc("/api/history/stats", s.handleAPIHistoryStats)
	mux.HandleFunc("/api/screenshots", s.handleAPIScreenshots)
	mux.HandleFunc("/api/config", s.handleAPIConfig)

	// MCP-UI endpoints (only if status provider is configured)
	if s.statusProvider != nil {
		baseURL := fmt.Sprintf("http://%s", s.addr)
		s.uiHandlers = NewUIHandlers(s.statusProvider, baseURL, s.cfg.ThumbnailCacheSec)

		// Set action executor if provider implements it
		if executor, ok := s.statusProvider.(ActionExecutor); ok {
			s.uiHandlers.SetActionExecutor(executor)
		}

		mux.HandleFunc("/ui/status", s.uiHandlers.HandleUIStatus)
		mux.HandleFunc("/ui/scenes", s.uiHandlers.HandleUIScenes)
		mux.HandleFunc("/ui/audio", s.uiHandlers.HandleUIAudio)
		mux.HandleFunc("/ui/screenshots", s.uiHandlers.HandleUIScreenshots)
		mux.HandleFunc("/ui/scene-thumbnail/", s.uiHandlers.HandleSceneThumbnail)
		mux.HandleFunc("/ui/action", s.uiHandlers.HandleUIAction)
	}

	// Documentation endpoints
	mux.HandleFunc("/docs", s.handleDocsIndex)
	mux.HandleFunc("/docs/", s.handleDocView)
	mux.HandleFunc("/api/docs", s.handleDocsAPI)

	// Serve static files for the web dashboard
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Printf("Warning: failed to setup static file serving: %v", err)
	} else {
		mux.Handle("/", http.FileServer(http.FS(staticFS)))
	}

	return mux
}

// Start begins serving HTTP requests in a background goroutine.
// Returns an error if the server is already running or fails to bind.
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("HTTP server already running")
	}

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to bind to %s: %w", s.addr, err)
	}
	s.listener = listener

	// Update addr with actual bound address (useful if port was 0)
	s.addr = listener.Addr().String()

	// Setup routes (must happen after addr is finalized for UI handlers)
	mux := s.setupRoutes()

	s.httpServer = &http.Server{
		Addr:         s.addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	s.running = true

	go func() {
		if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
		}
	}()

	return nil
}

// Stop gracefully shuts down the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = false
	s.mu.Unlock()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	return nil
}

// GetAddr returns the base URL of the HTTP server.
func (s *Server) GetAddr() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return fmt.Sprintf("http://%s", s.addr)
}

// GetScreenshotURL returns the full URL for a screenshot source.
func (s *Server) GetScreenshotURL(sourceName string) string {
	return fmt.Sprintf("%s/screenshot/%s", s.GetAddr(), sourceName)
}

// IsRunning returns whether the server is currently running.
func (s *Server) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// handleScreenshot serves the latest screenshot for a source.
// URL pattern: GET /screenshot/{source_name}
func (s *Server) handleScreenshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract source name from path: /screenshot/{source_name}
	path := strings.TrimPrefix(r.URL.Path, "/screenshot/")
	if path == "" || path == "/" {
		http.Error(w, "Source name required", http.StatusBadRequest)
		return
	}
	sourceName := path

	// Validate source name to prevent path traversal attacks
	if !isValidSourceName(sourceName) {
		log.Printf("Invalid source name rejected: %q", sourceName)
		http.Error(w, "Invalid source name", http.StatusBadRequest)
		return
	}

	// Look up the screenshot source by name
	source, err := s.storage.GetScreenshotSourceByName(r.Context(), sourceName)
	if err != nil {
		log.Printf("Screenshot source lookup failed for %q: %v", sourceName, err)
		http.Error(w, "Source not found", http.StatusNotFound)
		return
	}

	// Get the latest screenshot for this source
	screenshot, err := s.storage.GetLatestScreenshot(r.Context(), source.ID)
	if err != nil {
		log.Printf("No screenshots available for source %q (ID: %d): %v", sourceName, source.ID, err)
		http.Error(w, "No screenshots available", http.StatusNotFound)
		return
	}

	// Decode base64 image data
	imageData, err := base64.StdEncoding.DecodeString(screenshot.ImageData)
	if err != nil {
		log.Printf("Failed to decode screenshot data for source %q: %v", sourceName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", screenshot.MimeType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(imageData)))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("X-Screenshot-Source", sourceName)
	w.Header().Set("X-Screenshot-Captured", screenshot.CapturedAt.Format(time.RFC3339))

	// Write the image data
	w.WriteHeader(http.StatusOK)
	w.Write(imageData)
}

// handleHealth provides a simple health check endpoint.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
