package http

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ironystock/agentic-obs/internal/storage"
)

// Config holds HTTP server configuration options.
type Config struct {
	// Host to bind to (default: "localhost")
	Host string
	// Port to listen on (default: 8765)
	Port int
}

// DefaultConfig returns the default HTTP server configuration.
func DefaultConfig() Config {
	return Config{
		Host: "localhost",
		Port: 8765,
	}
}

// Server provides HTTP endpoints for serving screenshot images.
type Server struct {
	storage    *storage.DB
	httpServer *http.Server
	addr       string
	cfg        Config

	mu       sync.RWMutex
	running  bool
	listener net.Listener
}

// NewServer creates a new HTTP server for serving screenshots.
func NewServer(db *storage.DB, cfg Config) *Server {
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	if cfg.Port == 0 {
		cfg.Port = 8765
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	s := &Server{
		storage: db,
		cfg:     cfg,
		addr:    addr,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/screenshot/", s.handleScreenshot)
	mux.HandleFunc("/health", s.handleHealth)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
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
	s.httpServer.Addr = s.addr

	s.running = true

	go func() {
		if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			// Log error but don't crash - the server will be marked as not running
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

	// Look up the screenshot source by name
	source, err := s.storage.GetScreenshotSourceByName(r.Context(), sourceName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Screenshot source not found: %s", sourceName), http.StatusNotFound)
		return
	}

	// Get the latest screenshot for this source
	screenshot, err := s.storage.GetLatestScreenshot(r.Context(), source.ID)
	if err != nil {
		http.Error(w, "No screenshots available", http.StatusNotFound)
		return
	}

	// Decode base64 image data
	imageData, err := base64.StdEncoding.DecodeString(screenshot.ImageData)
	if err != nil {
		http.Error(w, "Failed to decode image data", http.StatusInternalServerError)
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
