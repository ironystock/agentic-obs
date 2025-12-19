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

// mockActionExecutor implements ActionExecutor for testing
type mockActionExecutor struct {
	setSceneErr      error
	setSceneCalled   string
	toggleMuteErr    error
	toggleMuteCalled string
	setVolumeErr     error
	setVolumeCalled  struct {
		name     string
		volumeDb float64
	}
	thumbnailData     []byte
	thumbnailMimeType string
	thumbnailErr      error
}

func (m *mockActionExecutor) SetCurrentScene(sceneName string) error {
	m.setSceneCalled = sceneName
	return m.setSceneErr
}

func (m *mockActionExecutor) ToggleInputMute(inputName string) error {
	m.toggleMuteCalled = inputName
	return m.toggleMuteErr
}

func (m *mockActionExecutor) SetInputVolume(inputName string, volumeDb float64) error {
	m.setVolumeCalled.name = inputName
	m.setVolumeCalled.volumeDb = volumeDb
	return m.setVolumeErr
}

func (m *mockActionExecutor) TakeSceneThumbnail(sceneName string) ([]byte, string, error) {
	if m.thumbnailErr != nil {
		return nil, "", m.thumbnailErr
	}
	if m.thumbnailData != nil {
		return m.thumbnailData, m.thumbnailMimeType, nil
	}
	// Return simple placeholder
	return []byte("<svg></svg>"), "image/svg+xml", nil
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

func TestSetActionExecutor(t *testing.T) {
	handlers := NewUIHandlers(&mockStatusProvider{}, "http://localhost:8765")
	executor := &mockActionExecutor{}

	handlers.SetActionExecutor(executor)

	assert.Equal(t, executor, handlers.actionExecutor)
}

func TestHandleSceneThumbnail(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		executor       *mockActionExecutor
		wantStatus     int
		wantType       string
		wantContains   string
	}{
		{
			name:       "returns thumbnail image",
			method:     http.MethodGet,
			path:       "/ui/scene-thumbnail/Gaming",
			executor:   &mockActionExecutor{thumbnailData: []byte{0x89, 0x50, 0x4E, 0x47}, thumbnailMimeType: "image/png"},
			wantStatus: http.StatusOK,
			wantType:   "image/png",
		},
		{
			name:       "returns SVG placeholder when no executor",
			method:     http.MethodGet,
			path:       "/ui/scene-thumbnail/Gaming",
			executor:   nil,
			wantStatus: http.StatusOK,
			wantType:   "image/svg+xml",
			wantContains: "Gaming",
		},
		{
			name:       "returns SVG placeholder on error",
			method:     http.MethodGet,
			path:       "/ui/scene-thumbnail/Gaming",
			executor:   &mockActionExecutor{thumbnailErr: errors.New("screenshot failed")},
			wantStatus: http.StatusOK,
			wantType:   "image/svg+xml",
		},
		{
			name:       "rejects POST method",
			method:     http.MethodPost,
			path:       "/ui/scene-thumbnail/Gaming",
			executor:   nil,
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "rejects missing scene name",
			method:     http.MethodGet,
			path:       "/ui/scene-thumbnail/",
			executor:   nil,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewUIHandlers(&mockStatusProvider{}, "http://localhost:8765")
			if tt.executor != nil {
				handlers.SetActionExecutor(tt.executor)
			}

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			handlers.HandleSceneThumbnail(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if tt.wantType != "" {
				assert.Equal(t, tt.wantType, rec.Header().Get("Content-Type"))
			}
			if tt.wantContains != "" {
				assert.Contains(t, rec.Body.String(), tt.wantContains)
			}
		})
	}
}

func TestHandleUIActionWithExecutor(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		executor     *mockActionExecutor
		wantStatus   int
		wantType     string
		wantSuccess  bool
		wantError    bool
	}{
		{
			name: "executes set_current_scene successfully",
			body: `{"type":"tool","messageId":"msg-1","payload":{"toolName":"set_current_scene","params":{"scene_name":"Gaming"}}}`,
			executor: &mockActionExecutor{},
			wantStatus: http.StatusOK,
			wantType:   "ui-message-response",
			wantSuccess: true,
		},
		{
			name: "returns error when set_current_scene fails",
			body: `{"type":"tool","messageId":"msg-2","payload":{"toolName":"set_current_scene","params":{"scene_name":"Invalid"}}}`,
			executor: &mockActionExecutor{setSceneErr: errors.New("scene not found")},
			wantStatus: http.StatusOK,
			wantType:   "ui-message-response",
			wantError: true,
		},
		{
			name: "executes toggle_input_mute successfully",
			body: `{"type":"tool","messageId":"msg-3","payload":{"toolName":"toggle_input_mute","params":{"input_name":"Mic"}}}`,
			executor: &mockActionExecutor{},
			wantStatus: http.StatusOK,
			wantType:   "ui-message-response",
			wantSuccess: true,
		},
		{
			name: "executes set_input_volume successfully",
			body: `{"type":"tool","messageId":"msg-4","payload":{"toolName":"set_input_volume","params":{"input_name":"Mic","volume_db":-10.0}}}`,
			executor: &mockActionExecutor{},
			wantStatus: http.StatusOK,
			wantType:   "ui-message-response",
			wantSuccess: true,
		},
		{
			name: "returns error for unsupported tool",
			body: `{"type":"tool","messageId":"msg-5","payload":{"toolName":"unknown_tool","params":{}}}`,
			executor: &mockActionExecutor{},
			wantStatus: http.StatusOK,
			wantType:   "ui-message-response",
			wantError: true,
		},
		{
			name: "returns error for missing parameters",
			body: `{"type":"tool","messageId":"msg-6","payload":{"toolName":"set_current_scene","params":{}}}`,
			executor: &mockActionExecutor{},
			wantStatus: http.StatusOK,
			wantType:   "ui-message-response",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers := NewUIHandlers(&mockStatusProvider{}, "http://localhost:8765")
			handlers.SetActionExecutor(tt.executor)

			req := httptest.NewRequest(http.MethodPost, "/ui/action", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handlers.HandleUIAction(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)

			var response map[string]any
			err := json.NewDecoder(rec.Body).Decode(&response)
			require.NoError(t, err)

			assert.Equal(t, tt.wantType, response["type"])

			if tt.wantSuccess {
				payload := response["payload"].(map[string]any)
				assert.NotNil(t, payload["response"])
			}
			if tt.wantError {
				payload := response["payload"].(map[string]any)
				assert.NotNil(t, payload["error"])
			}
		})
	}
}

func TestSceneInfoEnhancements(t *testing.T) {
	// Test that SceneInfo includes new fields
	scene := SceneInfo{
		Name:         "Main",
		Index:        0,
		IsCurrent:    true,
		SourceCount:  5,
		ThumbnailURL: "http://localhost:8765/ui/scene-thumbnail/Main",
	}

	assert.Equal(t, "Main", scene.Name)
	assert.Equal(t, 0, scene.Index)
	assert.True(t, scene.IsCurrent)
	assert.Equal(t, 5, scene.SourceCount)
	assert.Equal(t, "http://localhost:8765/ui/scene-thumbnail/Main", scene.ThumbnailURL)
}
