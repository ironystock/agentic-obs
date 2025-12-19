package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ironystock/agentic-obs/internal/storage"
)

// StatusResponse represents the server status API response.
type StatusResponse struct {
	Status      string                 `json:"status"`
	ServerName  string                 `json:"server_name"`
	Version     string                 `json:"version"`
	Uptime      string                 `json:"uptime"`
	HTTPAddress string                 `json:"http_address"`
	OBS         map[string]interface{} `json:"obs,omitempty"`
}

// handleAPIStatus returns server status as JSON.
func (s *Server) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := StatusResponse{
		Status:      "ok",
		ServerName:  "agentic-obs",
		Version:     "0.1.0",
		Uptime:      time.Since(s.startTime).Round(time.Second).String(),
		HTTPAddress: s.addr,
	}

	writeJSON(w, http.StatusOK, response)
}

// handleAPIHistory returns action history as JSON.
func (s *Server) handleAPIHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 500 {
			limit = parsed
		}
	}

	toolFilter := r.URL.Query().Get("tool")

	var records []storage.ActionRecord
	var err error

	if toolFilter != "" {
		records, err = s.storage.GetActionsByTool(r.Context(), toolFilter, limit)
	} else {
		records, err = s.storage.GetRecentActions(r.Context(), limit)
	}

	if err != nil {
		log.Printf("Failed to get action history: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve history"})
		return
	}

	response := map[string]interface{}{
		"count":   len(records),
		"limit":   limit,
		"actions": records,
	}

	writeJSON(w, http.StatusOK, response)
}

// handleAPIHistoryStats returns action history statistics.
func (s *Server) handleAPIHistoryStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := s.storage.GetActionStats(r.Context())
	if err != nil {
		log.Printf("Failed to get action stats: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve stats"})
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

// handleAPIScreenshots returns list of screenshot sources with URLs.
func (s *Server) handleAPIScreenshots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sources, err := s.storage.ListScreenshotSources(r.Context())
	if err != nil {
		log.Printf("Failed to list screenshot sources: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve screenshots"})
		return
	}

	// Enrich with URLs
	enriched := make([]map[string]interface{}, len(sources))
	for i, src := range sources {
		enriched[i] = map[string]interface{}{
			"id":           src.ID,
			"name":         src.Name,
			"source_name":  src.SourceName,
			"cadence_ms":   src.CadenceMs,
			"image_format": src.ImageFormat,
			"enabled":      src.Enabled,
			"url":          s.GetScreenshotURL(src.Name),
			"created_at":   src.CreatedAt.Format(time.RFC3339),
		}
	}

	response := map[string]interface{}{
		"count":   len(sources),
		"sources": enriched,
	}

	writeJSON(w, http.StatusOK, response)
}

// handleAPIConfig returns current configuration (GET) or updates it (POST).
func (s *Server) handleAPIConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetConfig(w, r)
	case http.MethodPost:
		s.handleUpdateConfig(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetConfig returns current configuration.
func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	obsConfig, err := s.storage.LoadOBSConfig(r.Context())
	if err != nil {
		log.Printf("Failed to load OBS config: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to load config"})
		return
	}

	toolGroups, err := s.storage.LoadToolGroupConfig(r.Context())
	if err != nil {
		log.Printf("Failed to load tool group config: %v", err)
		toolGroups = storage.ToolGroupConfig{Core: true, Visual: true, Layout: true, Audio: true, Sources: true}
	}

	webServer, err := s.storage.LoadWebServerConfig(r.Context())
	if err != nil {
		log.Printf("Failed to load web server config: %v", err)
		webServer = storage.WebServerConfig{Enabled: true, Host: "localhost", Port: 8765}
	}

	response := map[string]interface{}{
		"obs": map[string]interface{}{
			"host": obsConfig.Host,
			"port": obsConfig.Port,
			// Password intentionally omitted for security
		},
		"tool_groups": map[string]bool{
			"core":    toolGroups.Core,
			"visual":  toolGroups.Visual,
			"layout":  toolGroups.Layout,
			"audio":   toolGroups.Audio,
			"sources": toolGroups.Sources,
		},
		"web_server": map[string]interface{}{
			"enabled": webServer.Enabled,
			"host":    webServer.Host,
			"port":    webServer.Port,
		},
	}

	writeJSON(w, http.StatusOK, response)
}

// handleUpdateConfig updates configuration from POST body.
func (s *Server) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	// Limit request body size to prevent memory exhaustion attacks (64KB is plenty for config)
	const maxBodySize = 64 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON body"})
		return
	}

	// Process tool_groups updates
	if tg, ok := updates["tool_groups"].(map[string]interface{}); ok {
		config := storage.ToolGroupConfig{
			Core:    getBool(tg, "core", true),
			Visual:  getBool(tg, "visual", true),
			Layout:  getBool(tg, "layout", true),
			Audio:   getBool(tg, "audio", true),
			Sources: getBool(tg, "sources", true),
		}
		if err := s.storage.SaveToolGroupConfig(r.Context(), config); err != nil {
			log.Printf("Failed to save tool group config: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save tool groups"})
			return
		}
	}

	// Process web_server updates
	if ws, ok := updates["web_server"].(map[string]interface{}); ok {
		config := storage.WebServerConfig{
			Enabled: getBool(ws, "enabled", true),
			Host:    getString(ws, "host", "localhost"),
			Port:    getInt(ws, "port", 8765),
		}
		if err := s.storage.SaveWebServerConfig(r.Context(), config); err != nil {
			log.Printf("Failed to save web server config: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save web server config"})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "Configuration updated. Restart server for changes to take effect.",
	})
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

// Helper functions for safely extracting values from map[string]interface{}
func getBool(m map[string]interface{}, key string, defaultVal bool) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return defaultVal
}

func getString(m map[string]interface{}, key string, defaultVal string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return defaultVal
}

func getInt(m map[string]interface{}, key string, defaultVal int) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return defaultVal
}
