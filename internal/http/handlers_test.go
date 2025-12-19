package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ironystock/agentic-obs/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testServer creates a Server with a test database for testing handlers.
func testServer(t *testing.T) (*Server, func()) {
	t.Helper()

	// Create temp directory for test database
	tempDir, err := os.MkdirTemp("", "agentic-obs-http-test-*")
	require.NoError(t, err)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := storage.New(context.Background(), storage.Config{Path: dbPath})
	require.NoError(t, err)

	cfg := DefaultConfig()
	s := &Server{
		storage:   db,
		cfg:       cfg,
		addr:      "localhost:8765",
		startTime: time.Now(),
	}

	cleanup := func() {
		db.Close()
		os.RemoveAll(tempDir)
	}

	return s, cleanup
}

func TestHandleAPIStatus(t *testing.T) {
	t.Run("returns status on GET", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
		w := httptest.NewRecorder()

		s.handleAPIStatus(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response StatusResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response.Status)
		assert.Equal(t, "agentic-obs", response.ServerName)
		assert.Equal(t, "0.1.0", response.Version)
		assert.NotEmpty(t, response.Uptime)
	})

	t.Run("rejects non-GET methods", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
			req := httptest.NewRequest(method, "/api/status", nil)
			w := httptest.NewRecorder()

			s.handleAPIStatus(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code, "method %s should be rejected", method)
		}
	})
}

func TestHandleAPIHistory(t *testing.T) {
	t.Run("returns empty list when no history", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodGet, "/api/history", nil)
		w := httptest.NewRecorder()

		s.handleAPIHistory(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["count"])
	})

	t.Run("returns recorded actions", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		// Record some actions
		_, err := s.storage.RecordAction(context.Background(), storage.ActionRecord{
			Action:   "Test action",
			ToolName: "test_tool",
			Success:  true,
		})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/history", nil)
		w := httptest.NewRecorder()

		s.handleAPIHistory(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(1), response["count"])
	})

	t.Run("respects limit parameter", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		// Record multiple actions
		for i := 0; i < 10; i++ {
			_, err := s.storage.RecordAction(context.Background(), storage.ActionRecord{
				Action:  "Action",
				Success: true,
			})
			require.NoError(t, err)
		}

		req := httptest.NewRequest(http.MethodGet, "/api/history?limit=3", nil)
		w := httptest.NewRecorder()

		s.handleAPIHistory(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(3), response["count"])
		assert.Equal(t, float64(3), response["limit"])
	})

	t.Run("filters by tool parameter", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		// Record actions with different tools
		_, err := s.storage.RecordAction(context.Background(), storage.ActionRecord{
			Action: "A", ToolName: "tool_a", Success: true,
		})
		require.NoError(t, err)
		_, err = s.storage.RecordAction(context.Background(), storage.ActionRecord{
			Action: "B", ToolName: "tool_b", Success: true,
		})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/history?tool=tool_a", nil)
		w := httptest.NewRecorder()

		s.handleAPIHistory(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(1), response["count"])
	})

	t.Run("rejects non-GET methods", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodPost, "/api/history", nil)
		w := httptest.NewRecorder()

		s.handleAPIHistory(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func TestHandleAPIHistoryStats(t *testing.T) {
	t.Run("returns stats for empty database", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodGet, "/api/history/stats", nil)
		w := httptest.NewRecorder()

		s.handleAPIHistoryStats(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["total_actions"])
	})

	t.Run("returns correct stats with data", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		// Record actions
		_, err := s.storage.RecordAction(context.Background(), storage.ActionRecord{
			Action: "S1", Success: true, DurationMs: 100,
		})
		require.NoError(t, err)
		_, err = s.storage.RecordAction(context.Background(), storage.ActionRecord{
			Action: "F1", Success: false, DurationMs: 50,
		})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/history/stats", nil)
		w := httptest.NewRecorder()

		s.handleAPIHistoryStats(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(2), response["total_actions"])
		assert.Equal(t, float64(1), response["successful_actions"])
		assert.Equal(t, float64(1), response["failed_actions"])
	})

	t.Run("rejects non-GET methods", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodPost, "/api/history/stats", nil)
		w := httptest.NewRecorder()

		s.handleAPIHistoryStats(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func TestHandleAPIScreenshots(t *testing.T) {
	t.Run("returns empty list when no sources", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodGet, "/api/screenshots", nil)
		w := httptest.NewRecorder()

		s.handleAPIScreenshots(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), response["count"])
	})

	t.Run("returns screenshot sources with URLs", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		// Create a screenshot source
		_, err := s.storage.CreateScreenshotSource(context.Background(), storage.ScreenshotSource{
			Name:        "test_source",
			SourceName:  "OBS Source",
			CadenceMs:   5000,
			ImageFormat: "png",
			Enabled:     true,
		})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/screenshots", nil)
		w := httptest.NewRecorder()

		s.handleAPIScreenshots(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(1), response["count"])

		sources := response["sources"].([]interface{})
		source := sources[0].(map[string]interface{})
		assert.Equal(t, "test_source", source["name"])
		assert.Contains(t, source["url"], "/screenshot/test_source")
	})

	t.Run("rejects non-GET methods", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodPost, "/api/screenshots", nil)
		w := httptest.NewRecorder()

		s.handleAPIScreenshots(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func TestHandleAPIConfig(t *testing.T) {
	t.Run("GET returns configuration", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		// Save some config first
		err := s.storage.SaveOBSConfig(context.Background(), storage.OBSConfig{
			Host: "localhost", Port: 4455, Password: "secret",
		})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
		w := httptest.NewRecorder()

		s.handleAPIConfig(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		obs := response["obs"].(map[string]interface{})
		assert.Equal(t, "localhost", obs["host"])
		assert.Equal(t, float64(4455), obs["port"])
		// Password should NOT be included
		_, hasPassword := obs["password"]
		assert.False(t, hasPassword, "password should not be returned")
	})

	t.Run("POST updates tool groups", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		body := `{"tool_groups": {"core": true, "visual": false, "layout": true, "audio": false, "sources": true}}`
		req := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		s.handleAPIConfig(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify config was saved
		toolGroups, err := s.storage.LoadToolGroupConfig(context.Background())
		require.NoError(t, err)
		assert.True(t, toolGroups.Core)
		assert.False(t, toolGroups.Visual)
		assert.True(t, toolGroups.Layout)
		assert.False(t, toolGroups.Audio)
		assert.True(t, toolGroups.Sources)
	})

	t.Run("POST updates web server config", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		body := `{"web_server": {"enabled": false, "host": "0.0.0.0", "port": 9000}}`
		req := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		s.handleAPIConfig(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify config was saved
		webServer, err := s.storage.LoadWebServerConfig(context.Background())
		require.NoError(t, err)
		assert.False(t, webServer.Enabled)
		assert.Equal(t, "0.0.0.0", webServer.Host)
		assert.Equal(t, 9000, webServer.Port)
	})

	t.Run("POST rejects invalid JSON", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		body := `{invalid json}`
		req := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		s.handleAPIConfig(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST rejects oversized body", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		// Create a body larger than 64KB
		largeBody := bytes.Repeat([]byte("x"), 65*1024)
		req := httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewReader(largeBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		s.handleAPIConfig(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST rejects invalid host", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		// Test various invalid hosts
		invalidHosts := []string{"192.168.1.1", "10.0.0.1", "example.com", "0.0.0.1"}
		for _, host := range invalidHosts {
			body := `{"web_server": {"host": "` + host + `", "port": 8765}}`
			req := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			s.handleAPIConfig(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code, "host %s should be rejected", host)

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response["error"], "Invalid host")
		}
	})

	t.Run("POST rejects invalid port", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		// Test ports outside valid range
		invalidPorts := []int{0, 80, 443, 1023, 65536, 70000}
		for _, port := range invalidPorts {
			body := `{"web_server": {"host": "localhost", "port": ` + strconv.Itoa(port) + `}}`
			req := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			s.handleAPIConfig(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code, "port %d should be rejected", port)

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response["error"], "Invalid port")
		}
	})

	t.Run("POST accepts valid host and port", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		// Test valid combinations
		validConfigs := []struct {
			host string
			port int
		}{
			{"localhost", 8765},
			{"127.0.0.1", 1024},
			{"0.0.0.0", 65535},
			{"localhost", 9000},
		}
		for _, cfg := range validConfigs {
			body := `{"web_server": {"host": "` + cfg.host + `", "port": ` + strconv.Itoa(cfg.port) + `}}`
			req := httptest.NewRequest(http.MethodPost, "/api/config", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			s.handleAPIConfig(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "host=%s port=%d should be accepted", cfg.host, cfg.port)
		}
	})

	t.Run("rejects unsupported methods", func(t *testing.T) {
		s, cleanup := testServer(t)
		defer cleanup()

		for _, method := range []string{http.MethodPut, http.MethodDelete, http.MethodPatch} {
			req := httptest.NewRequest(method, "/api/config", nil)
			w := httptest.NewRecorder()

			s.handleAPIConfig(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code, "method %s should be rejected", method)
		}
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("getBool returns value when present", func(t *testing.T) {
		m := map[string]interface{}{"key": true}
		assert.True(t, getBool(m, "key", false))
	})

	t.Run("getBool returns default when missing", func(t *testing.T) {
		m := map[string]interface{}{}
		assert.True(t, getBool(m, "key", true))
	})

	t.Run("getBool returns default for wrong type", func(t *testing.T) {
		m := map[string]interface{}{"key": "not a bool"}
		assert.True(t, getBool(m, "key", true))
	})

	t.Run("getString returns value when present", func(t *testing.T) {
		m := map[string]interface{}{"key": "value"}
		assert.Equal(t, "value", getString(m, "key", "default"))
	})

	t.Run("getString returns default when missing", func(t *testing.T) {
		m := map[string]interface{}{}
		assert.Equal(t, "default", getString(m, "key", "default"))
	})

	t.Run("getInt returns value when present", func(t *testing.T) {
		m := map[string]interface{}{"key": float64(42)} // JSON numbers are float64
		assert.Equal(t, 42, getInt(m, "key", 0))
	})

	t.Run("getInt returns default when missing", func(t *testing.T) {
		m := map[string]interface{}{}
		assert.Equal(t, 99, getInt(m, "key", 99))
	})
}

func TestWriteJSON(t *testing.T) {
	t.Run("sets correct content type", func(t *testing.T) {
		w := httptest.NewRecorder()
		writeJSON(w, http.StatusOK, map[string]string{"test": "value"})

		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})

	t.Run("sets correct status code", func(t *testing.T) {
		w := httptest.NewRecorder()
		writeJSON(w, http.StatusCreated, map[string]string{})

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("encodes data as JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		writeJSON(w, http.StatusOK, map[string]string{"key": "value"})

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "value", response["key"])
	})
}

func TestIsValidHost(t *testing.T) {
	t.Run("accepts valid hosts", func(t *testing.T) {
		validHosts := []string{"localhost", "127.0.0.1", "0.0.0.0"}
		for _, host := range validHosts {
			assert.True(t, isValidHost(host), "host %s should be valid", host)
		}
	})

	t.Run("rejects invalid hosts", func(t *testing.T) {
		invalidHosts := []string{
			"192.168.1.1",  // Private IP
			"10.0.0.1",     // Private IP
			"8.8.8.8",      // Public IP
			"example.com",  // Domain
			"0.0.0.1",      // Invalid loopback
			"",             // Empty
			"LOCALHOST",    // Case sensitive
			"127.0.0.2",    // Not standard loopback
		}
		for _, host := range invalidHosts {
			assert.False(t, isValidHost(host), "host %s should be invalid", host)
		}
	})
}
