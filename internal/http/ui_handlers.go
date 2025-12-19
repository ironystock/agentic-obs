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

// SceneInfo represents a scene for UI display.
type SceneInfo struct {
	Name      string `json:"name"`
	Index     int    `json:"index"`
	IsCurrent bool   `json:"isCurrent"`
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
	baseURL        string
}

// NewUIHandlers creates a new UIHandlers instance.
func NewUIHandlers(provider StatusProvider, baseURL string) *UIHandlers {
	return &UIHandlers{
		statusProvider: provider,
		baseURL:        baseURL,
	}
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

	// TODO: Route action to appropriate handler based on type
	// For now, return acknowledgment
	response := map[string]any{
		"type":      "ui-message-received",
		"messageId": action.MessageID,
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
    <style>` + sharedCSS + `</style>
</head>
<body>
    <div class="container">
        <h1>Scene <span class="accent">Preview</span></h1>

        <div class="card">
            <h2>Available Scenes</h2>
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
        document.querySelectorAll('.scene-card').forEach(card => {
            card.addEventListener('click', () => {
                const sceneName = card.dataset.scene;
                if (!sceneName) return;

                fetch('{{.BaseURL}}/ui/action', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        type: 'tool',
                        messageId: 'scene-' + Date.now(),
                        payload: { toolName: 'set_current_scene', params: { scene_name: sceneName } }
                    })
                }).then(() => location.reload());
            });
        });
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
    <style>` + sharedCSS + `</style>
</head>
<body>
    <div class="container">
        <h1>Audio <span class="accent">Mixer</span></h1>

        <div class="card">
            <h2>Audio Inputs</h2>
            {{range .Inputs}}
            <div class="slider-container" data-input="{{.Name}}">
                <div class="slider-header">
                    <span class="slider-name">{{.Name}}</span>
                    <span class="slider-value">{{printf "%.1f" .VolumeDB}} dB</span>
                    <button class="btn mute-btn {{if .IsMuted}}muted{{end}}" onclick="toggleMute('{{.Name}}')">
                        {{if .IsMuted}}Unmute{{else}}Mute{{end}}
                    </button>
                </div>
                <input type="range" min="0" max="100" value="{{printf "%.0f" .VolumePercent}}"
                       onchange="setVolume('{{.Name}}', this.value)">
            </div>
            {{else}}
            <p style="color: var(--text-secondary);">No audio inputs available</p>
            {{end}}
        </div>
    </div>

    <script>
        function toggleMute(inputName) {
            fetch('{{.BaseURL}}/ui/action', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    type: 'tool',
                    messageId: 'mute-' + Date.now(),
                    payload: { toolName: 'toggle_input_mute', params: { input_name: inputName } }
                })
            }).then(() => location.reload());
        }

        function setVolume(inputName, value) {
            const db = (value / 100) * 26 - 26; // Map 0-100 to -26 to 0 dB
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

