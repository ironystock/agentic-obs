package http

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

// StatusProvider interface for getting OBS status data.
// This allows the HTTP server to fetch status without direct OBS dependency.
type StatusProvider interface {
	// GetStatus returns the current OBS status as a JSON-serializable struct.
	GetStatus() (any, error)
	// GetScenes returns a list of scenes.
	GetScenes() ([]SceneInfo, error)
	// GetAudioInputs returns a list of audio inputs with their current state.
	GetAudioInputs() ([]AudioInputInfo, error)
	// GetScreenshotSources returns a list of screenshot sources.
	GetScreenshotSources() ([]ScreenshotSourceInfo, error)
}

// ActionExecutor interface for executing UI actions.
// This allows the UI to trigger tool calls on the MCP server.
type ActionExecutor interface {
	// SetCurrentScene switches to the specified scene.
	SetCurrentScene(sceneName string) error
	// ToggleInputMute toggles mute state for an audio input.
	ToggleInputMute(inputName string) error
	// SetInputVolume sets the volume for an audio input.
	SetInputVolume(inputName string, volumeDb float64) error
	// TakeSceneThumbnail captures a thumbnail of the specified scene.
	TakeSceneThumbnail(sceneName string) ([]byte, string, error)
}

// SceneInfo represents a scene for UI display.
type SceneInfo struct {
	Name         string `json:"name"`
	Index        int    `json:"index"`
	IsCurrent    bool   `json:"isCurrent"`
	SourceCount  int    `json:"sourceCount"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
}

// AudioInputInfo represents an audio input for UI display.
type AudioInputInfo struct {
	Name          string  `json:"name"`
	Volume        float64 `json:"volume"`        // 0.0-1.0
	VolumePercent float64 `json:"volumePercent"` // 0-100 for slider
	VolumeDB      float64 `json:"volumeDb"`      // in dB
	IsMuted       bool    `json:"isMuted"`
	InputKind     string  `json:"inputKind"`
}

// ScreenshotSourceInfo represents a screenshot source for UI display.
type ScreenshotSourceInfo struct {
	Name        string `json:"name"`
	SourceName  string `json:"sourceName"`
	CadenceMs   int    `json:"cadenceMs"`
	ImageURL    string `json:"imageUrl"`
	LastCapture string `json:"lastCapture"`
}

// UIHandlers provides HTTP handlers for MCP-UI resources.
type UIHandlers struct {
	statusProvider StatusProvider
	actionExecutor ActionExecutor
	baseURL        string
}

// NewUIHandlers creates a new UIHandlers instance.
func NewUIHandlers(provider StatusProvider, baseURL string) *UIHandlers {
	return &UIHandlers{
		statusProvider: provider,
		baseURL:        baseURL,
	}
}

// SetActionExecutor sets the action executor for handling UI actions.
func (h *UIHandlers) SetActionExecutor(executor ActionExecutor) {
	h.actionExecutor = executor
}

// HandleUIStatus serves the status dashboard UI.
func (h *UIHandlers) HandleUIStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current status
	status, err := h.statusProvider.GetStatus()
	if err != nil {
		h.renderError(w, "Failed to get OBS status", err)
		return
	}

	// Get scenes
	scenes, err := h.statusProvider.GetScenes()
	if err != nil {
		scenes = []SceneInfo{} // Empty on error, show what we can
	}

	data := map[string]any{
		"Status":  status,
		"Scenes":  scenes,
		"BaseURL": h.baseURL,
	}

	h.renderTemplate(w, statusDashboardTemplate, data)
}

// HandleUIScenes serves the scene preview UI.
func (h *UIHandlers) HandleUIScenes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	scenes, err := h.statusProvider.GetScenes()
	if err != nil {
		h.renderError(w, "Failed to get scenes", err)
		return
	}

	data := map[string]any{
		"Scenes":  scenes,
		"BaseURL": h.baseURL,
	}

	h.renderTemplate(w, scenePreviewTemplate, data)
}

// HandleUIAudio serves the audio mixer UI.
func (h *UIHandlers) HandleUIAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	inputs, err := h.statusProvider.GetAudioInputs()
	if err != nil {
		h.renderError(w, "Failed to get audio inputs", err)
		return
	}

	data := map[string]any{
		"Inputs":  inputs,
		"BaseURL": h.baseURL,
	}

	h.renderTemplate(w, audioMixerTemplate, data)
}

// HandleUIScreenshots serves the screenshot gallery UI.
func (h *UIHandlers) HandleUIScreenshots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sources, err := h.statusProvider.GetScreenshotSources()
	if err != nil {
		h.renderError(w, "Failed to get screenshot sources", err)
		return
	}

	data := map[string]any{
		"Sources": sources,
		"BaseURL": h.baseURL,
	}

	h.renderTemplate(w, screenshotGalleryTemplate, data)
}

// HandleSceneThumbnail serves a thumbnail image for a scene.
// GET /ui/scene-thumbnail/{sceneName}
func (h *UIHandlers) HandleSceneThumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract scene name from path
	path := r.URL.Path
	prefix := "/ui/scene-thumbnail/"
	if len(path) <= len(prefix) {
		http.Error(w, "Scene name required", http.StatusBadRequest)
		return
	}
	sceneName := path[len(prefix):]

	if h.actionExecutor == nil {
		// Return a placeholder SVG if no executor
		h.servePlaceholderThumbnail(w, sceneName)
		return
	}

	// Take a screenshot of the scene
	imageData, mimeType, err := h.actionExecutor.TakeSceneThumbnail(sceneName)
	if err != nil {
		// Return placeholder on error
		h.servePlaceholderThumbnail(w, sceneName)
		return
	}

	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Cache-Control", "max-age=5") // Cache for 5 seconds
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(imageData)))
	w.Write(imageData)
}

func (h *UIHandlers) servePlaceholderThumbnail(w http.ResponseWriter, sceneName string) {
	// Generate a simple SVG placeholder
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="320" height="180" viewBox="0 0 320 180">
		<rect fill="#1a1a2e" width="320" height="180"/>
		<rect fill="#16213e" x="10" y="10" width="300" height="160" rx="8"/>
		<text fill="#e94560" x="160" y="85" text-anchor="middle" font-family="sans-serif" font-size="14" font-weight="bold">%s</text>
		<text fill="#666" x="160" y="105" text-anchor="middle" font-family="sans-serif" font-size="11">Scene Preview</text>
	</svg>`, sceneName)

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "max-age=60")
	w.Write([]byte(svg))
}

// HandleUIAction handles UIAction requests from embedded UIs.
// POST /ui/action - receives UIAction JSON, executes, returns UIResponse.
func (h *UIHandlers) HandleUIAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var action struct {
		Type      string          `json:"type"`
		MessageID string          `json:"messageId"`
		Payload   json.RawMessage `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		h.jsonError(w, "Invalid action format", http.StatusBadRequest)
		return
	}

	// Handle tool actions
	if action.Type == "tool" && h.actionExecutor != nil {
		var toolPayload struct {
			ToolName string         `json:"toolName"`
			Params   map[string]any `json:"params"`
		}
		if err := json.Unmarshal(action.Payload, &toolPayload); err == nil {
			err := h.executeToolAction(toolPayload.ToolName, toolPayload.Params)
			if err != nil {
				h.sendActionResponse(w, action.MessageID, nil, err)
				return
			}
			h.sendActionResponse(w, action.MessageID, map[string]string{"status": "success"}, nil)
			return
		}
	}

	// Default: return acknowledgment
	response := map[string]any{
		"type":      "ui-message-received",
		"messageId": action.MessageID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UIHandlers) executeToolAction(toolName string, params map[string]any) error {
	switch toolName {
	case "set_current_scene":
		if sceneName, ok := params["scene_name"].(string); ok {
			return h.actionExecutor.SetCurrentScene(sceneName)
		}
		return fmt.Errorf("missing scene_name parameter")

	case "toggle_input_mute":
		if inputName, ok := params["input_name"].(string); ok {
			return h.actionExecutor.ToggleInputMute(inputName)
		}
		return fmt.Errorf("missing input_name parameter")

	case "set_input_volume":
		inputName, ok1 := params["input_name"].(string)
		volumeDb, ok2 := params["volume_db"].(float64)
		if ok1 && ok2 {
			return h.actionExecutor.SetInputVolume(inputName, volumeDb)
		}
		return fmt.Errorf("missing input_name or volume_db parameter")

	default:
		return fmt.Errorf("unsupported tool: %s", toolName)
	}
}

func (h *UIHandlers) sendActionResponse(w http.ResponseWriter, messageID string, result any, err error) {
	response := map[string]any{
		"type":      "ui-message-response",
		"messageId": messageID,
	}

	if err != nil {
		response["payload"] = map[string]any{
			"error": map[string]string{"message": err.Error()},
		}
	} else {
		response["payload"] = map[string]any{
			"response": result,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UIHandlers) renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	t, err := template.New("ui").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		// Headers already sent, just log
		fmt.Printf("Template execution error: %v\n", err)
	}
}

func (h *UIHandlers) renderError(w http.ResponseWriter, message string, err error) {
	errMsg := message
	if err != nil {
		errMsg = fmt.Sprintf("%s: %v", message, err)
	}

	data := map[string]any{
		"Error": errMsg,
	}
	h.renderTemplate(w, errorTemplate, data)
}

func (h *UIHandlers) jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// CSS shared across all UI templates - matches existing dashboard theme
const sharedCSS = `
:root {
    --bg-primary: #1a1a2e;
    --bg-secondary: #16213e;
    --bg-card: #0f3460;
    --text-primary: #eaeaea;
    --text-secondary: #a0a0a0;
    --accent: #e94560;
    --accent-hover: #ff6b6b;
    --success: #4ecca3;
    --warning: #ffc107;
    --error: #ff6b6b;
    --border: #2a2a4a;
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: var(--bg-primary);
    color: var(--text-primary);
    line-height: 1.6;
    padding: 20px;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
}

h1 {
    font-size: 1.5rem;
    margin-bottom: 20px;
    color: var(--text-primary);
}

h1 .accent {
    color: var(--accent);
}

.card {
    background: var(--bg-secondary);
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 20px;
    border: 1px solid var(--border);
}

.card h2 {
    font-size: 1rem;
    color: var(--text-secondary);
    margin-bottom: 15px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 20px;
}

.status-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 0;
    border-bottom: 1px solid var(--border);
}

.status-item:last-child {
    border-bottom: none;
}

.status-label {
    color: var(--text-secondary);
}

.status-value {
    font-weight: 500;
}

.status-value.success {
    color: var(--success);
}

.status-value.error {
    color: var(--error);
}

.status-value.warning {
    color: var(--warning);
}

.badge {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 4px 12px;
    border-radius: 12px;
    font-size: 0.875rem;
}

.badge.online {
    background: rgba(78, 204, 163, 0.15);
    color: var(--success);
}

.badge.offline {
    background: rgba(255, 107, 107, 0.15);
    color: var(--error);
}

.badge .dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: currentColor;
}

.btn {
    background: var(--accent);
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.875rem;
    transition: background 0.2s;
}

.btn:hover {
    background: var(--accent-hover);
}

.btn.secondary {
    background: var(--bg-card);
    border: 1px solid var(--border);
}

.btn.secondary:hover {
    background: var(--border);
}

.scene-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 15px;
}

.scene-card {
    background: var(--bg-card);
    border-radius: 8px;
    padding: 15px;
    cursor: pointer;
    transition: all 0.2s;
    border: 2px solid transparent;
}

.scene-card:hover {
    border-color: var(--accent);
}

.scene-card.active {
    border-color: var(--success);
    background: rgba(78, 204, 163, 0.1);
}

.scene-card .name {
    font-weight: 500;
    margin-bottom: 5px;
}

.scene-card .index {
    font-size: 0.75rem;
    color: var(--text-secondary);
}

.slider-container {
    padding: 15px 0;
    border-bottom: 1px solid var(--border);
}

.slider-container:last-child {
    border-bottom: none;
}

.slider-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 10px;
}

.slider-name {
    font-weight: 500;
}

.slider-value {
    font-size: 0.875rem;
    color: var(--text-secondary);
}

input[type="range"] {
    width: 100%;
    height: 6px;
    border-radius: 3px;
    background: var(--bg-card);
    -webkit-appearance: none;
}

input[type="range"]::-webkit-slider-thumb {
    -webkit-appearance: none;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: var(--accent);
    cursor: pointer;
}

.mute-btn {
    padding: 6px 12px;
    font-size: 0.75rem;
}

.mute-btn.muted {
    background: var(--error);
}

.screenshot-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 20px;
}

.screenshot-card {
    background: var(--bg-card);
    border-radius: 8px;
    overflow: hidden;
}

.screenshot-card img {
    width: 100%;
    height: 180px;
    object-fit: cover;
    background: var(--bg-primary);
}

.screenshot-card .info {
    padding: 15px;
}

.screenshot-card .name {
    font-weight: 500;
    margin-bottom: 5px;
}

.screenshot-card .meta {
    font-size: 0.75rem;
    color: var(--text-secondary);
}

.error-message {
    background: rgba(255, 107, 107, 0.1);
    border: 1px solid var(--error);
    color: var(--error);
    padding: 20px;
    border-radius: 8px;
    text-align: center;
}
`

// Status Dashboard Template
var statusDashboardTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>OBS Status Dashboard</title>
    <style>` + sharedCSS + `</style>
</head>
<body>
    <div class="container">
        <h1><span class="accent">agentic-obs</span> Status</h1>

        <div class="grid">
            <div class="card">
                <h2>Connection</h2>
                <div class="status-item">
                    <span class="status-label">OBS WebSocket</span>
                    <span class="badge online"><span class="dot"></span> Connected</span>
                </div>
            </div>

            <div class="card">
                <h2>Recording</h2>
                <div class="status-item">
                    <span class="status-label">Status</span>
                    <span class="status-value">Idle</span>
                </div>
            </div>

            <div class="card">
                <h2>Streaming</h2>
                <div class="status-item">
                    <span class="status-label">Status</span>
                    <span class="status-value">Offline</span>
                </div>
            </div>
        </div>

        <div class="card">
            <h2>Scenes ({{len .Scenes}})</h2>
            <div class="scene-grid">
                {{range .Scenes}}
                <div class="scene-card {{if .IsCurrent}}active{{end}}" data-scene="{{.Name}}">
                    <div class="name">{{.Name}}</div>
                    <div class="index">Scene #{{.Index}}</div>
                </div>
                {{else}}
                <p style="color: var(--text-secondary);">No scenes available</p>
                {{end}}
            </div>
        </div>
    </div>

    <script>
        // UIAction helper for scene switching
        document.querySelectorAll('.scene-card').forEach(card => {
            card.addEventListener('click', () => {
                const sceneName = card.dataset.scene;
                if (!sceneName) return;

                // Send UIAction to parent
                const action = {
                    type: 'tool',
                    messageId: 'scene-' + Date.now(),
                    payload: {
                        toolName: 'set_current_scene',
                        params: { scene_name: sceneName }
                    }
                };

                // Post to action endpoint
                fetch('{{.BaseURL}}/ui/action', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(action)
                }).then(() => {
                    // Refresh to show updated state
                    location.reload();
                });
            });
        });
    </script>
</body>
</html>`

// Scene Preview Template
var scenePreviewTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Scene Preview</title>
    <style>` + sharedCSS + `
        .scene-card {
            background: var(--bg-card);
            border-radius: 8px;
            overflow: hidden;
            cursor: pointer;
            transition: all 0.2s;
            border: 2px solid transparent;
        }

        .scene-card:hover {
            border-color: var(--accent);
            transform: translateY(-2px);
        }

        .scene-card.active {
            border-color: var(--success);
        }

        .scene-card.active .current-badge {
            display: inline-block;
        }

        .scene-thumbnail {
            width: 100%;
            height: 120px;
            object-fit: cover;
            background: var(--bg-primary);
        }

        .scene-info {
            padding: 12px;
        }

        .scene-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 4px;
        }

        .scene-name {
            font-weight: 600;
            font-size: 0.95rem;
        }

        .current-badge {
            display: none;
            background: var(--success);
            color: var(--bg-primary);
            font-size: 0.65rem;
            padding: 2px 6px;
            border-radius: 4px;
            text-transform: uppercase;
            font-weight: 600;
        }

        .scene-meta {
            display: flex;
            gap: 12px;
            font-size: 0.75rem;
            color: var(--text-secondary);
        }

        .scene-meta span {
            display: flex;
            align-items: center;
            gap: 4px;
        }

        .loading {
            opacity: 0.6;
            pointer-events: none;
        }

        .keyboard-hint {
            text-align: center;
            color: var(--text-secondary);
            font-size: 0.8rem;
            margin-top: 20px;
        }

        .keyboard-hint kbd {
            background: var(--bg-card);
            padding: 2px 8px;
            border-radius: 4px;
            border: 1px solid var(--border);
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Scene <span class="accent">Preview</span></h1>

        <div class="card">
            <h2>Available Scenes ({{len .Scenes}})</h2>
            <div class="scene-grid">
                {{range .Scenes}}
                <div class="scene-card {{if .IsCurrent}}active{{end}}" data-scene="{{.Name}}" tabindex="0">
                    <img class="scene-thumbnail"
                         src="{{.ThumbnailURL}}"
                         alt="{{.Name}}"
                         loading="lazy"
                         onerror="this.style.display='none'">
                    <div class="scene-info">
                        <div class="scene-header">
                            <span class="scene-name">{{.Name}}</span>
                            <span class="current-badge">Live</span>
                        </div>
                        <div class="scene-meta">
                            <span>Scene #{{.Index}}</span>
                            <span>{{.SourceCount}} sources</span>
                        </div>
                    </div>
                </div>
                {{else}}
                <p style="color: var(--text-secondary);">No scenes available</p>
                {{end}}
            </div>
            <div class="keyboard-hint">
                Press <kbd>1</kbd>-<kbd>9</kbd> to switch scenes, <kbd>Enter</kbd> to select focused
            </div>
        </div>
    </div>

    <script>
        let isLoading = false;

        function switchScene(sceneName) {
            if (isLoading) return;
            isLoading = true;

            // Add loading state to all cards
            document.querySelectorAll('.scene-card').forEach(c => c.classList.add('loading'));

            fetch('{{.BaseURL}}/ui/action', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    type: 'tool',
                    messageId: 'scene-' + Date.now(),
                    payload: { toolName: 'set_current_scene', params: { scene_name: sceneName } }
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.payload && data.payload.error) {
                    alert('Error: ' + data.payload.error.message);
                    isLoading = false;
                    document.querySelectorAll('.scene-card').forEach(c => c.classList.remove('loading'));
                } else {
                    location.reload();
                }
            })
            .catch(() => {
                location.reload();
            });
        }

        // Click handling
        document.querySelectorAll('.scene-card').forEach(card => {
            card.addEventListener('click', () => {
                const sceneName = card.dataset.scene;
                if (sceneName) switchScene(sceneName);
            });
        });

        // Keyboard handling
        document.addEventListener('keydown', (e) => {
            // Number keys 1-9 for quick switch
            if (e.key >= '1' && e.key <= '9') {
                const index = parseInt(e.key) - 1;
                const cards = document.querySelectorAll('.scene-card');
                if (cards[index]) {
                    const sceneName = cards[index].dataset.scene;
                    if (sceneName) switchScene(sceneName);
                }
            }

            // Enter on focused card
            if (e.key === 'Enter' && document.activeElement.classList.contains('scene-card')) {
                const sceneName = document.activeElement.dataset.scene;
                if (sceneName) switchScene(sceneName);
            }

            // Arrow key navigation
            if (e.key === 'ArrowRight' || e.key === 'ArrowLeft') {
                const cards = Array.from(document.querySelectorAll('.scene-card'));
                const current = cards.indexOf(document.activeElement);
                if (current >= 0) {
                    const next = e.key === 'ArrowRight' ? current + 1 : current - 1;
                    if (cards[next]) cards[next].focus();
                } else if (cards.length > 0) {
                    cards[0].focus();
                }
            }
        });

        // Auto-refresh thumbnails every 10 seconds
        setInterval(() => {
            document.querySelectorAll('.scene-thumbnail').forEach(img => {
                const url = new URL(img.src, location.href);
                url.searchParams.set('t', Date.now());
                img.src = url.toString();
            });
        }, 10000);
    </script>
</body>
</html>`

// Audio Mixer Template
var audioMixerTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Audio Mixer</title>
    <style>` + sharedCSS + `
        .audio-channel {
            background: var(--bg-card);
            border-radius: 8px;
            padding: 16px;
            margin-bottom: 12px;
            transition: all 0.2s;
            border: 2px solid transparent;
        }

        .audio-channel:hover {
            border-color: var(--border);
        }

        .audio-channel.focused {
            border-color: var(--accent);
        }

        .audio-channel.muted {
            opacity: 0.6;
        }

        .channel-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 12px;
        }

        .channel-info {
            display: flex;
            align-items: center;
            gap: 12px;
        }

        .channel-name {
            font-weight: 600;
            font-size: 1rem;
        }

        .channel-type {
            background: var(--bg-secondary);
            color: var(--text-secondary);
            font-size: 0.7rem;
            padding: 2px 8px;
            border-radius: 4px;
            text-transform: uppercase;
        }

        .channel-controls {
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .volume-display {
            font-family: monospace;
            font-size: 0.9rem;
            color: var(--text-secondary);
            min-width: 70px;
            text-align: right;
        }

        .mute-toggle {
            width: 40px;
            height: 40px;
            border-radius: 8px;
            border: none;
            cursor: pointer;
            display: flex;
            align-items: center;
            justify-content: center;
            transition: all 0.2s;
            background: var(--bg-secondary);
            color: var(--text-primary);
        }

        .mute-toggle:hover {
            background: var(--border);
        }

        .mute-toggle.muted {
            background: var(--error);
            color: white;
        }

        .mute-toggle svg {
            width: 20px;
            height: 20px;
        }

        .slider-track {
            position: relative;
            height: 8px;
            background: var(--bg-secondary);
            border-radius: 4px;
            overflow: hidden;
        }

        .slider-fill {
            position: absolute;
            left: 0;
            top: 0;
            height: 100%;
            background: linear-gradient(90deg, var(--success) 0%, var(--warning) 70%, var(--error) 100%);
            border-radius: 4px;
            transition: width 0.1s;
        }

        .audio-channel.muted .slider-fill {
            background: var(--text-secondary);
        }

        .volume-slider {
            width: 100%;
            height: 8px;
            -webkit-appearance: none;
            background: transparent;
            position: relative;
            z-index: 1;
            margin-top: -8px;
        }

        .volume-slider::-webkit-slider-thumb {
            -webkit-appearance: none;
            width: 18px;
            height: 18px;
            border-radius: 50%;
            background: var(--text-primary);
            cursor: pointer;
            border: 2px solid var(--bg-card);
            box-shadow: 0 2px 4px rgba(0,0,0,0.3);
            transition: transform 0.1s;
        }

        .volume-slider::-webkit-slider-thumb:hover {
            transform: scale(1.2);
        }

        .volume-slider::-moz-range-thumb {
            width: 18px;
            height: 18px;
            border-radius: 50%;
            background: var(--text-primary);
            cursor: pointer;
            border: 2px solid var(--bg-card);
        }

        .db-markers {
            display: flex;
            justify-content: space-between;
            margin-top: 4px;
            font-size: 0.65rem;
            color: var(--text-secondary);
        }

        .keyboard-hint {
            text-align: center;
            color: var(--text-secondary);
            font-size: 0.8rem;
            margin-top: 20px;
        }

        .keyboard-hint kbd {
            background: var(--bg-card);
            padding: 2px 8px;
            border-radius: 4px;
            border: 1px solid var(--border);
        }

        .refresh-indicator {
            position: fixed;
            top: 20px;
            right: 20px;
            background: var(--bg-secondary);
            color: var(--text-secondary);
            padding: 8px 12px;
            border-radius: 6px;
            font-size: 0.75rem;
            opacity: 0;
            transition: opacity 0.3s;
        }

        .refresh-indicator.visible {
            opacity: 1;
        }

        .no-inputs {
            text-align: center;
            padding: 40px;
            color: var(--text-secondary);
        }

        .no-inputs svg {
            width: 48px;
            height: 48px;
            margin-bottom: 12px;
            opacity: 0.5;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Audio <span class="accent">Mixer</span></h1>

        <div class="card">
            <h2>Audio Inputs ({{len .Inputs}})</h2>
            {{if .Inputs}}
            {{range $i, $input := .Inputs}}
            <div class="audio-channel {{if .IsMuted}}muted{{end}}"
                 data-input="{{.Name}}"
                 data-index="{{$i}}"
                 tabindex="0">
                <div class="channel-header">
                    <div class="channel-info">
                        <span class="channel-name">{{.Name}}</span>
                        <span class="channel-type">{{.InputKind}}</span>
                    </div>
                    <div class="channel-controls">
                        <span class="volume-display">{{printf "%.1f" .VolumeDB}} dB</span>
                        <button class="mute-toggle {{if .IsMuted}}muted{{end}}"
                                onclick="toggleMute('{{.Name}}')"
                                title="{{if .IsMuted}}Unmute{{else}}Mute{{end}} (M)">
                            {{if .IsMuted}}
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <path d="M11 5L6 9H2v6h4l5 4V5z"/>
                                <line x1="23" y1="9" x2="17" y2="15"/>
                                <line x1="17" y1="9" x2="23" y2="15"/>
                            </svg>
                            {{else}}
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <path d="M11 5L6 9H2v6h4l5 4V5z"/>
                                <path d="M19.07 4.93a10 10 0 0 1 0 14.14M15.54 8.46a5 5 0 0 1 0 7.07"/>
                            </svg>
                            {{end}}
                        </button>
                    </div>
                </div>
                <div class="slider-track">
                    <div class="slider-fill" style="width: {{printf "%.0f" .VolumePercent}}%"></div>
                </div>
                <input type="range"
                       class="volume-slider"
                       min="0" max="100"
                       value="{{printf "%.0f" .VolumePercent}}"
                       data-input="{{.Name}}"
                       oninput="updateSlider(this)"
                       onchange="setVolume('{{.Name}}', this.value)">
                <div class="db-markers">
                    <span>-∞</span>
                    <span>-20</span>
                    <span>-10</span>
                    <span>0 dB</span>
                </div>
            </div>
            {{end}}
            {{else}}
            <div class="no-inputs">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 5L6 9H2v6h4l5 4V5z"/>
                    <line x1="23" y1="9" x2="17" y2="15"/>
                    <line x1="17" y1="9" x2="23" y2="15"/>
                </svg>
                <p>No audio inputs available</p>
                <p style="font-size: 0.8rem; margin-top: 8px;">Connect to OBS to see audio sources</p>
            </div>
            {{end}}
        </div>

        <div class="keyboard-hint">
            <kbd>↑</kbd><kbd>↓</kbd> Navigate | <kbd>←</kbd><kbd>→</kbd> Adjust Volume | <kbd>M</kbd> Mute/Unmute | <kbd>0</kbd> Reset
        </div>
    </div>

    <div class="refresh-indicator" id="refreshIndicator">Updating...</div>

    <script>
        let focusedIndex = 0;
        const channels = document.querySelectorAll('.audio-channel');

        function updateSlider(slider) {
            const channel = slider.closest('.audio-channel');
            const fill = channel.querySelector('.slider-fill');
            const display = channel.querySelector('.volume-display');
            const value = slider.value;
            fill.style.width = value + '%';
            // Map 0-100 to -inf to 0 dB (using -60 as practical minimum)
            const db = value === '0' ? '-∞' : ((value / 100) * 26 - 26).toFixed(1);
            display.textContent = db + ' dB';
        }

        function toggleMute(inputName) {
            showRefreshIndicator();
            fetch('{{.BaseURL}}/ui/action', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    type: 'tool',
                    messageId: 'mute-' + Date.now(),
                    payload: { toolName: 'toggle_input_mute', params: { input_name: inputName } }
                })
            })
            .then(response => response.json())
            .then(data => {
                if (data.payload && data.payload.error) {
                    alert('Error: ' + data.payload.error.message);
                    hideRefreshIndicator();
                } else {
                    refreshAudioState();
                }
            })
            .catch(() => hideRefreshIndicator());
        }

        function setVolume(inputName, value) {
            const db = value == 0 ? -100 : (value / 100) * 26 - 26; // Map 0-100 to -100 (mute) to 0 dB
            fetch('{{.BaseURL}}/ui/action', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    type: 'tool',
                    messageId: 'vol-' + Date.now(),
                    payload: { toolName: 'set_input_volume', params: { input_name: inputName, volume_db: db } }
                })
            });
        }

        function showRefreshIndicator() {
            document.getElementById('refreshIndicator').classList.add('visible');
        }

        function hideRefreshIndicator() {
            document.getElementById('refreshIndicator').classList.remove('visible');
        }

        function refreshAudioState() {
            location.reload();
        }

        function setFocus(index) {
            if (index < 0) index = channels.length - 1;
            if (index >= channels.length) index = 0;
            focusedIndex = index;

            channels.forEach((ch, i) => {
                ch.classList.toggle('focused', i === index);
            });
            channels[index].focus();
        }

        // Keyboard navigation
        document.addEventListener('keydown', (e) => {
            if (channels.length === 0) return;

            const focusedChannel = channels[focusedIndex];
            const slider = focusedChannel?.querySelector('.volume-slider');
            const inputName = focusedChannel?.dataset.input;

            switch (e.key) {
                case 'ArrowUp':
                    e.preventDefault();
                    setFocus(focusedIndex - 1);
                    break;
                case 'ArrowDown':
                    e.preventDefault();
                    setFocus(focusedIndex + 1);
                    break;
                case 'ArrowLeft':
                    e.preventDefault();
                    if (slider) {
                        slider.value = Math.max(0, parseInt(slider.value) - 5);
                        updateSlider(slider);
                        setVolume(inputName, slider.value);
                    }
                    break;
                case 'ArrowRight':
                    e.preventDefault();
                    if (slider) {
                        slider.value = Math.min(100, parseInt(slider.value) + 5);
                        updateSlider(slider);
                        setVolume(inputName, slider.value);
                    }
                    break;
                case 'm':
                case 'M':
                    e.preventDefault();
                    if (inputName) toggleMute(inputName);
                    break;
                case '0':
                    e.preventDefault();
                    if (slider) {
                        slider.value = 77; // ~0 dB
                        updateSlider(slider);
                        setVolume(inputName, slider.value);
                    }
                    break;
            }
        });

        // Click to focus channel
        channels.forEach((channel, index) => {
            channel.addEventListener('click', (e) => {
                if (!e.target.closest('.mute-toggle')) {
                    setFocus(index);
                }
            });
            channel.addEventListener('focus', () => {
                focusedIndex = index;
                channels.forEach((ch, i) => ch.classList.toggle('focused', i === index));
            });
        });

        // Auto-refresh audio state every 3 seconds
        setInterval(() => {
            fetch('{{.BaseURL}}/ui/audio')
                .then(response => response.text())
                .then(html => {
                    // Only refresh if no slider is being dragged
                    if (!document.querySelector('.volume-slider:active')) {
                        const parser = new DOMParser();
                        const doc = parser.parseFromString(html, 'text/html');
                        const newChannels = doc.querySelectorAll('.audio-channel');

                        newChannels.forEach(newChannel => {
                            const inputName = newChannel.dataset.input;
                            const oldChannel = document.querySelector('[data-input="' + inputName + '"]');
                            if (oldChannel) {
                                // Update mute state
                                const wasMuted = oldChannel.classList.contains('muted');
                                const isMuted = newChannel.classList.contains('muted');
                                if (wasMuted !== isMuted) {
                                    oldChannel.classList.toggle('muted', isMuted);
                                    const btn = oldChannel.querySelector('.mute-toggle');
                                    btn.classList.toggle('muted', isMuted);
                                    btn.innerHTML = newChannel.querySelector('.mute-toggle').innerHTML;
                                }

                                // Update volume display (but not slider if user isn't touching it)
                                const newDisplay = newChannel.querySelector('.volume-display').textContent;
                                oldChannel.querySelector('.volume-display').textContent = newDisplay;
                            }
                        });
                    }
                });
        }, 3000);

        // Initial focus
        if (channels.length > 0) setFocus(0);
    </script>
</body>
</html>`

// Screenshot Gallery Template
var screenshotGalleryTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Screenshot Gallery</title>
    <style>` + sharedCSS + `</style>
</head>
<body>
    <div class="container">
        <h1>Screenshot <span class="accent">Gallery</span></h1>

        <div class="screenshot-grid">
            {{range .Sources}}
            <div class="screenshot-card">
                <img src="{{.ImageURL}}" alt="{{.Name}}" onerror="this.src='data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 width=%22300%22 height=%22180%22><rect fill=%22%231a1a2e%22 width=%22300%22 height=%22180%22/><text fill=%22%23666%22 x=%2250%%22 y=%2250%%22 text-anchor=%22middle%22>No Image</text></svg>'">
                <div class="info">
                    <div class="name">{{.Name}}</div>
                    <div class="meta">Source: {{.SourceName}} | Cadence: {{.CadenceMs}}ms</div>
                    <div class="meta">Last: {{.LastCapture}}</div>
                </div>
            </div>
            {{else}}
            <p style="color: var(--text-secondary);">No screenshot sources configured</p>
            {{end}}
        </div>
    </div>

    <script>
        // Auto-refresh images every 5 seconds
        setInterval(() => {
            document.querySelectorAll('.screenshot-card img').forEach(img => {
                const url = new URL(img.src, location.href);
                url.searchParams.set('t', Date.now());
                img.src = url.toString();
            });
        }, 5000);
    </script>
</body>
</html>`

// Error Template
var errorTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Error</title>
    <style>` + sharedCSS + `</style>
</head>
<body>
    <div class="container">
        <h1><span class="accent">Error</span></h1>
        <div class="error-message">
            {{.Error}}
        </div>
    </div>
</body>
</html>`
