package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockStatusProvider implements StatusProvider for testing
type mockStatusProvider struct {
	status            any
	statusErr         error
	scenes            []SceneInfo
	scenesErr         error
	audioInputs       []AudioInputInfo
	audioInputsErr    error
	screenshotSources []ScreenshotSourceInfo
	screenshotErr     error
}

func (m *mockStatusProvider) GetStatus() (any, error) {
	return m.status, m.statusErr
}

func (m *mockStatusProvider) GetScenes() ([]SceneInfo, error) {
	return m.scenes, m.scenesErr
}

func (m *mockStatusProvider) GetAudioInputs() ([]AudioInputInfo, error) {
	return m.audioInputs, m.audioInputsErr
}

func (m *mockStatusProvider) GetScreenshotSources() ([]ScreenshotSourceInfo, error) {
	return m.screenshotSources, m.screenshotErr
}

func TestHandleUIStatus(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		provider       *mockStatusProvider
		wantStatus     int
		wantContains   []string
		wantNotContain []string
	}{
		{
			name:   "returns status dashboard on GET",
			method: http.MethodGet,
			provider: &mockStatusProvider{
				status: map[string]any{
					"connected":    true,
					"recording":    false,
					"streaming":    false,
					"currentScene": "Main",
				},
				scenes: []SceneInfo{
					{Name: "Main", Index: 0, IsCurrent: true},
					{Name: "Gaming", Index: 1, IsCurrent: false},
				},
			},
			wantStatus:   http.StatusOK,
			wantContains: []string{"agentic-obs", "Status", "Main", "Gaming"},
		},
		{
			name:   "shows error when status fails",
			method: http.MethodGet,
			provider: &mockStatusProvider{
				statusErr: errors.New("connection failed"),
			},
			wantStatus:   http.StatusOK, // Still returns 200 with error template
			wantContains: []string{"Error", "connection failed"},
		},
		{
			name:   "shows empty scenes when scenes fail",
			method: http.MethodGet,
			provider: &mockStatusProvider{
				status: map[string]any{"connected": true},
				scenesErr: errors.New("scenes error"),
			},
			wantStatus:   http.StatusOK,
			wantContains: []string{"Scenes (0)"},
		},
		{
			name:       "rejects POST method",
			method:     http.MethodPost,
			provider:   &mockStatusProvider{},
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewUIHandlers(tt.provider, "http://localhost:8765")

			req := httptest.NewRequest(tt.method, "/ui/status", nil)
			rec := httptest.NewRecorder()

			handlers.HandleUIStatus(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			body := rec.Body.String()
			for _, want := range tt.wantContains {
				assert.Contains(t, body, want, "response should contain %q", want)
			}
			for _, notWant := range tt.wantNotContain {
				assert.NotContains(t, body, notWant, "response should not contain %q", notWant)
			}
		})
	}
}

func TestHandleUIScenes(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		provider     *mockStatusProvider
		wantStatus   int
		wantContains []string
	}{
		{
			name:   "returns scene grid on GET",
			method: http.MethodGet,
			provider: &mockStatusProvider{
				scenes: []SceneInfo{
					{Name: "Main", Index: 0, IsCurrent: true},
					{Name: "BRB", Index: 1, IsCurrent: false},
					{Name: "Ending", Index: 2, IsCurrent: false},
				},
			},
			wantStatus:   http.StatusOK,
			wantContains: []string{"Scene", "Preview", "Main", "BRB", "Ending", "active"},
		},
		{
			name:   "shows empty message when no scenes",
			method: http.MethodGet,
			provider: &mockStatusProvider{
				scenes: []SceneInfo{},
			},
			wantStatus:   http.StatusOK,
			wantContains: []string{"No scenes available"},
		},
		{
			name:   "shows error when scenes fail",
			method: http.MethodGet,
			provider: &mockStatusProvider{
				scenesErr: errors.New("OBS not connected"),
			},
			wantStatus:   http.StatusOK,
			wantContains: []string{"Error", "OBS not connected"},
		},
		{
			name:       "rejects PUT method",
			method:     http.MethodPut,
			provider:   &mockStatusProvider{},
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewUIHandlers(tt.provider, "http://localhost:8765")

			req := httptest.NewRequest(tt.method, "/ui/scenes", nil)
			rec := httptest.NewRecorder()

			handlers.HandleUIScenes(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			body := rec.Body.String()
			for _, want := range tt.wantContains {
				assert.Contains(t, body, want)
			}
		})
	}
}

func TestHandleUIAudio(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		provider     *mockStatusProvider
		wantStatus   int
		wantContains []string
	}{
		{
			name:   "returns audio mixer on GET",
			method: http.MethodGet,
			provider: &mockStatusProvider{
				audioInputs: []AudioInputInfo{
					{Name: "Mic/Aux", Volume: 1.0, VolumePercent: 100, VolumeDB: 0, IsMuted: false, InputKind: "wasapi_input_capture"},
					{Name: "Desktop Audio", Volume: 0.5, VolumePercent: 50, VolumeDB: -6, IsMuted: true, InputKind: "wasapi_output_capture"},
				},
			},
			wantStatus:   http.StatusOK,
			wantContains: []string{"Audio", "Mixer", "Mic/Aux", "Desktop Audio", "Unmute"},
		},
		{
			name:   "shows empty message when no inputs",
			method: http.MethodGet,
			provider: &mockStatusProvider{
				audioInputs: []AudioInputInfo{},
			},
			wantStatus:   http.StatusOK,
			wantContains: []string{"No audio inputs available"},
		},
		{
			name:   "shows error when audio fails",
			method: http.MethodGet,
			provider: &mockStatusProvider{
				audioInputsErr: errors.New("audio subsystem error"),
			},
			wantStatus:   http.StatusOK,
			wantContains: []string{"Error", "audio subsystem error"},
		},
		{
			name:       "rejects DELETE method",
			method:     http.MethodDelete,
			provider:   &mockStatusProvider{},
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewUIHandlers(tt.provider, "http://localhost:8765")

			req := httptest.NewRequest(tt.method, "/ui/audio", nil)
			rec := httptest.NewRecorder()

			handlers.HandleUIAudio(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			body := rec.Body.String()
			for _, want := range tt.wantContains {
				assert.Contains(t, body, want)
			}
		})
	}
}

func TestHandleUIScreenshots(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		provider     *mockStatusProvider
		wantStatus   int
		wantContains []string
	}{
		{
			name:   "returns screenshot gallery on GET",
			method: http.MethodGet,
			provider: &mockStatusProvider{
				screenshotSources: []ScreenshotSourceInfo{
					{Name: "main-monitor", SourceName: "Display Capture", CadenceMs: 5000, ImageURL: "http://localhost:8765/screenshot/main-monitor", LastCapture: "2025-01-01 12:00:00"},
					{Name: "game-capture", SourceName: "Game Capture", CadenceMs: 1000, ImageURL: "http://localhost:8765/screenshot/game-capture", LastCapture: "2025-01-01 12:00:01"},
				},
			},
			wantStatus:   http.StatusOK,
			wantContains: []string{"Screenshot", "Gallery", "main-monitor", "game-capture", "Display Capture", "Game Capture"},
		},
		{
			name:   "shows empty message when no sources",
			method: http.MethodGet,
			provider: &mockStatusProvider{
				screenshotSources: []ScreenshotSourceInfo{},
			},
			wantStatus:   http.StatusOK,
			wantContains: []string{"No screenshot sources configured"},
		},
		{
			name:   "shows error when screenshot sources fail",
			method: http.MethodGet,
			provider: &mockStatusProvider{
				screenshotErr: errors.New("database error"),
			},
			wantStatus:   http.StatusOK,
			wantContains: []string{"Error", "database error"},
		},
		{
			name:       "rejects PATCH method",
			method:     http.MethodPatch,
			provider:   &mockStatusProvider{},
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewUIHandlers(tt.provider, "http://localhost:8765")

			req := httptest.NewRequest(tt.method, "/ui/screenshots", nil)
			rec := httptest.NewRecorder()

			handlers.HandleUIScreenshots(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			body := rec.Body.String()
			for _, want := range tt.wantContains {
				assert.Contains(t, body, want)
			}
		})
	}
}

func TestHandleUIAction(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         string
		wantStatus   int
		wantResponse map[string]any
	}{
		{
			name:   "acknowledges tool action",
			method: http.MethodPost,
			body: `{
				"type": "tool",
				"messageId": "msg-123",
				"payload": {"toolName": "set_current_scene", "params": {"scene_name": "Gaming"}}
			}`,
			wantStatus: http.StatusOK,
			wantResponse: map[string]any{
				"type":      "ui-message-received",
				"messageId": "msg-123",
			},
		},
		{
			name:   "acknowledges intent action",
			method: http.MethodPost,
			body: `{
				"type": "intent",
				"messageId": "msg-456",
				"payload": {"intent": "switch_scene", "params": {}}
			}`,
			wantStatus: http.StatusOK,
			wantResponse: map[string]any{
				"type":      "ui-message-received",
				"messageId": "msg-456",
			},
		},
		{
			name:       "rejects GET method",
			method:     http.MethodGet,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "rejects invalid JSON",
			method:     http.MethodPost,
			body:       `{invalid json`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewUIHandlers(&mockStatusProvider{}, "http://localhost:8765")

			var body *strings.Reader
			if tt.body != "" {
				body = strings.NewReader(tt.body)
			} else {
				body = strings.NewReader("")
			}

			req := httptest.NewRequest(tt.method, "/ui/action", body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handlers.HandleUIAction(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			if tt.wantResponse != nil {
				var got map[string]any
				err := json.NewDecoder(rec.Body).Decode(&got)
				require.NoError(t, err)
				assert.Equal(t, tt.wantResponse["type"], got["type"])
				assert.Equal(t, tt.wantResponse["messageId"], got["messageId"])
			}
		})
	}
}

func TestUIHandlersContentType(t *testing.T) {
	provider := &mockStatusProvider{
		status: map[string]any{"connected": true},
		scenes: []SceneInfo{{Name: "Test", Index: 0, IsCurrent: true}},
	}
	handlers := NewUIHandlers(provider, "http://localhost:8765")

	tests := []struct {
		name        string
		handler     func(http.ResponseWriter, *http.Request)
		wantType    string
	}{
		{
			name:     "status returns HTML",
			handler:  handlers.HandleUIStatus,
			wantType: "text/html; charset=utf-8",
		},
		{
			name:     "scenes returns HTML",
			handler:  handlers.HandleUIScenes,
			wantType: "text/html; charset=utf-8",
		},
		{
			name:     "audio returns HTML",
			handler:  handlers.HandleUIAudio,
			wantType: "text/html; charset=utf-8",
		},
		{
			name:     "screenshots returns HTML",
			handler:  handlers.HandleUIScreenshots,
			wantType: "text/html; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/ui/test", nil)
			rec := httptest.NewRecorder()

			tt.handler(rec, req)

			assert.Equal(t, tt.wantType, rec.Header().Get("Content-Type"))
		})
	}
}

func TestNewUIHandlers(t *testing.T) {
	provider := &mockStatusProvider{}
	handlers := NewUIHandlers(provider, "http://example.com:9000")

	assert.NotNil(t, handlers)
	assert.Equal(t, "http://example.com:9000", handlers.baseURL)
	assert.Equal(t, provider, handlers.statusProvider)
}
