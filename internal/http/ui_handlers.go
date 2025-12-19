package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// StatusProvider defines the interface for retrieving OBS status data.
//
// This interface decouples the HTTP server from direct OBS dependencies,
// allowing the UI layer to fetch current state without knowing about the
// underlying MCP server or OBS WebSocket connection.
//
// The MCP server implements this interface to provide real-time OBS data
// to the web UI endpoints (/ui/status, /ui/scenes, /ui/audio, /ui/screenshots).
//
// Flow: HTTP Request → StatusProvider → MCP Server → OBS WebSocket → OBS Studio
type StatusProvider interface {
	// GetStatus returns the current OBS connection and streaming/recording state.
	// Returns a map containing: connected, recording, streaming, currentScene, obsVersion.
	GetStatus() (any, error)

	// GetScenes returns all available OBS scenes with their current state.
	// Each scene includes name, index, source count, and whether it's currently active.
	GetScenes() ([]SceneInfo, error)

	// GetAudioInputs returns all audio inputs with volume and mute state.
	// Only inputs that support volume control are included.
	GetAudioInputs() ([]AudioInputInfo, error)

	// GetScreenshotSources returns configured screenshot capture sources.
	// These are agentic monitoring sources, not OBS sources directly.
	GetScreenshotSources() ([]ScreenshotSourceInfo, error)
}

// ActionExecutor defines the interface for executing UI-triggered actions.
//
// This interface allows the web UI to trigger OBS operations through the
// MCP server without direct access to the OBS WebSocket connection.
// Actions are executed synchronously and return errors on failure.
//
// The MCP server implements this interface to handle POST requests to
// /ui/action, which are sent by JavaScript in the embedded UI templates.
//
// Flow: UI Click → POST /ui/action → ActionExecutor → MCP Server → OBS WebSocket → OBS Studio
type ActionExecutor interface {
	// SetCurrentScene switches OBS to the specified scene.
	// Returns an error if the scene doesn't exist or OBS is disconnected.
	SetCurrentScene(sceneName string) error

	// ToggleInputMute toggles the mute state for an audio input.
	// Returns an error if the input doesn't exist or doesn't support muting.
	ToggleInputMute(inputName string) error

	// SetInputVolume sets the volume for an audio input in decibels.
	// volumeDb should be in the range -100 (mute) to 0 (full volume).
	// Returns an error if the input doesn't exist or doesn't support volume.
	SetInputVolume(inputName string, volumeDb float64) error

	// TakeSceneThumbnail captures a thumbnail image of the specified scene.
	// Returns the image data, MIME type (e.g., "image/jpeg"), and any error.
	// Used by /ui/scene-thumbnail/{sceneName} endpoint.
	TakeSceneThumbnail(sceneName string) ([]byte, string, error)
}

// SceneInfo represents an OBS scene for UI display.
// Used by the scene preview grid to show available scenes with thumbnails.
type SceneInfo struct {
	Name         string `json:"name"`                   // Scene name as shown in OBS
	Index        int    `json:"index"`                  // Zero-based index in scene list
	IsCurrent    bool   `json:"isCurrent"`              // True if this is the active program scene
	SourceCount  int    `json:"sourceCount"`            // Number of sources in the scene
	ThumbnailURL string `json:"thumbnailUrl,omitempty"` // URL to fetch scene thumbnail image
}

// AudioInputInfo represents an audio input for the mixer UI.
// Includes both linear multiplier and dB values for flexible display.
type AudioInputInfo struct {
	Name          string  `json:"name"`          // Input name as shown in OBS
	Volume        float64 `json:"volume"`        // Volume multiplier (0.0-1.0, can exceed 1.0 for gain)
	VolumePercent float64 `json:"volumePercent"` // Logarithmic slider position (0-100)
	VolumeDB      float64 `json:"volumeDb"`      // Volume in decibels (-inf to 0, typically -60 to 0)
	IsMuted       bool    `json:"isMuted"`       // True if input is muted
	InputKind     string  `json:"inputKind"`     // OBS input type (e.g., "wasapi_input_capture")
}

// ScreenshotSourceInfo represents a configured screenshot monitoring source.
// These are agentic-obs specific sources for AI visual monitoring, not OBS sources.
type ScreenshotSourceInfo struct {
	Name        string `json:"name"`        // User-defined name for this monitoring source
	SourceName  string `json:"sourceName"`  // OBS source being captured
	CadenceMs   int    `json:"cadenceMs"`   // Capture interval in milliseconds
	ImageURL    string `json:"imageUrl"`    // URL to fetch latest screenshot
	LastCapture string `json:"lastCapture"` // Timestamp of most recent capture
}

// UIHandlers provides HTTP handlers for the MCP-UI web interface.
//
// UIHandlers serves embedded HTML templates that display OBS status,
// scene previews, audio mixer controls, and screenshot galleries.
// Each handler fetches data through StatusProvider and executes
// user actions through ActionExecutor.
//
// Endpoints:
//   - GET /ui/status - Status dashboard with connection info
//   - GET /ui/scenes - Scene grid with thumbnails and switching
//   - GET /ui/audio - Audio mixer with volume sliders and mute buttons
//   - GET /ui/screenshots - Screenshot gallery for monitoring sources
//   - GET /ui/scene-thumbnail/{name} - Scene thumbnail image
//   - POST /ui/action - Execute UI-triggered actions (scene switch, mute, volume)
type UIHandlers struct {
	statusProvider StatusProvider
	actionExecutor ActionExecutor
	baseURL        string // Base URL for constructing absolute URLs in templates
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

	h.renderTemplate(w, templateStatusDashboard, data)
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

	log.Printf("[UI Scenes] Serving scene preview with BaseURL: %s, %d scenes", h.baseURL, len(scenes))

	data := map[string]any{
		"Scenes":  scenes,
		"BaseURL": h.baseURL,
	}

	h.renderTemplate(w, templateScenePreview, data)
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

	h.renderTemplate(w, templateAudioMixer, data)
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

	h.renderTemplate(w, templateScreenshotGallery, data)
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
	log.Printf("[UI Action] Received request: %s %s", r.Method, r.URL.Path)

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
	if action.Type == "tool" {
		if h.actionExecutor == nil {
			log.Printf("[UI Action] ERROR: actionExecutor is nil, cannot execute tool")
			h.sendActionResponse(w, action.MessageID, nil, fmt.Errorf("action executor not configured"))
			return
		}

		var toolPayload struct {
			ToolName string         `json:"toolName"`
			Params   map[string]any `json:"params"`
		}
		if err := json.Unmarshal(action.Payload, &toolPayload); err != nil {
			log.Printf("[UI Action] ERROR: failed to parse payload: %v", err)
			h.sendActionResponse(w, action.MessageID, nil, fmt.Errorf("invalid payload: %w", err))
			return
		}

		log.Printf("[UI Action] Executing tool: %s with params: %v", toolPayload.ToolName, toolPayload.Params)
		err := h.executeToolAction(toolPayload.ToolName, toolPayload.Params)
		if err != nil {
			log.Printf("[UI Action] ERROR: tool execution failed: %v", err)
			h.sendActionResponse(w, action.MessageID, nil, err)
			return
		}
		log.Printf("[UI Action] SUCCESS: tool %s executed", toolPayload.ToolName)
		h.sendActionResponse(w, action.MessageID, map[string]string{"status": "success"}, nil)
		return
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

func (h *UIHandlers) renderTemplate(w http.ResponseWriter, templateName string, data map[string]any) {
	tmpl, err := getTemplate(templateName)
	if err != nil {
		log.Printf("Template load error for %s: %v", templateName, err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Inject shared CSS into template data
	data["SharedCSS"] = getSharedCSS()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		// Headers already sent, just log
		log.Printf("Template execution error for %s: %v", templateName, err)
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
	h.renderTemplate(w, templateError, data)
}

func (h *UIHandlers) jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
