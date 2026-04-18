package automation

import (
	"fmt"
	"log"
	"time"

	"github.com/ironystock/agentic-obs/internal/obs"
)

// OBSClient defines the interface for OBS operations used by the executor.
// This matches the mcp.OBSClient interface.
type OBSClient interface {
	// Scene operations
	SetCurrentScene(name string) error

	// Recording operations
	StartRecording() error
	StopRecording() (string, error)
	PauseRecording() error
	ResumeRecording() error

	// Streaming operations
	StartStreaming() error
	StopStreaming() error

	// Audio operations
	GetInputMute(inputName string) (bool, error)
	ToggleInputMute(inputName string) error
	SetInputVolume(inputName string, volumeDb *float64, volumeMul *float64) error

	// Source visibility
	ToggleSourceVisibility(sceneName string, sourceID int) (bool, error)

	// Virtual camera operations
	ToggleVirtualCam() (bool, error)
	StartVirtualCam() error
	StopVirtualCam() error

	// Replay buffer operations
	ToggleReplayBuffer() (bool, error)
	SaveReplayBuffer() error

	// Studio mode operations
	SetCurrentPreviewScene(sceneName string) error
	TriggerStudioModeTransition() error

	// Hotkey operations
	TriggerHotkeyByName(hotkeyName string) error

	// Event handling (needed for automation bridge)
	SetEventCallback(callback obs.EventCallback)
}

// Executor handles action execution against OBS.
type Executor struct {
	obsClient OBSClient
}

// NewExecutor creates a new action executor.
func NewExecutor(client OBSClient) *Executor {
	return &Executor{
		obsClient: client,
	}
}

// ExecuteAction runs a single action and returns the result.
func (e *Executor) ExecuteAction(action Action, index int) ActionResult {
	start := time.Now()
	result := ActionResult{
		ActionType: action.Type,
		Index:      index,
	}

	err := e.runAction(action)

	result.DurationMs = time.Since(start).Milliseconds()
	result.Success = err == nil
	if err != nil {
		result.Error = err.Error()
		log.Printf("[Automation] Action %d (%s) failed: %v", index, action.Type, err)
	} else {
		log.Printf("[Automation] Action %d (%s) completed in %dms", index, action.Type, result.DurationMs)
	}

	return result
}

// runAction dispatches to the appropriate handler based on action type.
func (e *Executor) runAction(action Action) error {
	switch action.Type {
	case ActionTypeSetScene:
		return e.setScene(action.Parameters)

	case ActionTypeToggleMute:
		return e.toggleMute(action.Parameters)

	case ActionTypeSetMute:
		return e.setMute(action.Parameters)

	case ActionTypeSetVolume:
		return e.setVolume(action.Parameters)

	case ActionTypeToggleVisibility:
		return e.toggleVisibility(action.Parameters)

	case ActionTypeSetVisibility:
		return e.setVisibility(action.Parameters)

	case ActionTypeStartRecording:
		return e.obsClient.StartRecording()

	case ActionTypeStopRecording:
		_, err := e.obsClient.StopRecording()
		return err

	case ActionTypePauseRecording:
		return e.obsClient.PauseRecording()

	case ActionTypeResumeRecording:
		return e.obsClient.ResumeRecording()

	case ActionTypeStartStreaming:
		return e.obsClient.StartStreaming()

	case ActionTypeStopStreaming:
		return e.obsClient.StopStreaming()

	case ActionTypeToggleVirtualCam:
		_, err := e.obsClient.ToggleVirtualCam()
		return err

	case ActionTypeStartVirtualCam:
		return e.obsClient.StartVirtualCam()

	case ActionTypeStopVirtualCam:
		return e.obsClient.StopVirtualCam()

	case ActionTypeToggleReplayBuffer:
		_, err := e.obsClient.ToggleReplayBuffer()
		return err

	case ActionTypeSaveReplay:
		return e.obsClient.SaveReplayBuffer()

	case ActionTypeTriggerHotkey:
		return e.triggerHotkey(action.Parameters)

	case ActionTypeTriggerTransition:
		return e.obsClient.TriggerStudioModeTransition()

	case ActionTypeSetPreviewScene:
		return e.setPreviewScene(action.Parameters)

	case ActionTypeDelay:
		return e.delay(action.Parameters)

	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// setScene switches to a scene.
func (e *Executor) setScene(params map[string]interface{}) error {
	sceneName, ok := getStringParam(params, "scene_name")
	if !ok {
		return fmt.Errorf("set_scene requires 'scene_name' parameter")
	}
	return e.obsClient.SetCurrentScene(sceneName)
}

// toggleMute toggles mute for an input.
func (e *Executor) toggleMute(params map[string]interface{}) error {
	inputName, ok := getStringParam(params, "input_name")
	if !ok {
		return fmt.Errorf("toggle_mute requires 'input_name' parameter")
	}
	return e.obsClient.ToggleInputMute(inputName)
}

// setMute sets the mute state for an input.
func (e *Executor) setMute(params map[string]interface{}) error {
	inputName, ok := getStringParam(params, "input_name")
	if !ok {
		return fmt.Errorf("set_mute requires 'input_name' parameter")
	}

	muted, ok := getBoolParam(params, "muted")
	if !ok {
		return fmt.Errorf("set_mute requires 'muted' parameter")
	}

	// Get current state
	currentMuted, err := e.obsClient.GetInputMute(inputName)
	if err != nil {
		return fmt.Errorf("failed to get current mute state: %w", err)
	}

	// Only toggle if state needs to change
	if currentMuted != muted {
		return e.obsClient.ToggleInputMute(inputName)
	}

	return nil
}

// setVolume sets volume for an input.
func (e *Executor) setVolume(params map[string]interface{}) error {
	inputName, ok := getStringParam(params, "input_name")
	if !ok {
		return fmt.Errorf("set_volume requires 'input_name' parameter")
	}

	volumeDb, hasDb := getFloat64Param(params, "volume_db")
	volumeMul, hasMul := getFloat64Param(params, "volume_mul")

	if !hasDb && !hasMul {
		return fmt.Errorf("set_volume requires either 'volume_db' or 'volume_mul' parameter")
	}

	var dbPtr, mulPtr *float64
	if hasDb {
		dbPtr = &volumeDb
	}
	if hasMul {
		mulPtr = &volumeMul
	}

	return e.obsClient.SetInputVolume(inputName, dbPtr, mulPtr)
}

// toggleVisibility toggles source visibility.
func (e *Executor) toggleVisibility(params map[string]interface{}) error {
	sceneName, ok := getStringParam(params, "scene_name")
	if !ok {
		return fmt.Errorf("toggle_visibility requires 'scene_name' parameter")
	}

	sourceID, ok := getIntParam(params, "source_id")
	if !ok {
		return fmt.Errorf("toggle_visibility requires 'source_id' parameter")
	}

	_, err := e.obsClient.ToggleSourceVisibility(sceneName, sourceID)
	return err
}

// setVisibility sets source visibility to a specific state.
func (e *Executor) setVisibility(params map[string]interface{}) error {
	// Note: OBS WebSocket doesn't have a direct "set visibility" - only toggle
	// For now, this is the same as toggle. A future enhancement could
	// get current state and only toggle if needed.
	return e.toggleVisibility(params)
}

// triggerHotkey triggers a hotkey by name.
func (e *Executor) triggerHotkey(params map[string]interface{}) error {
	hotkeyName, ok := getStringParam(params, "hotkey_name")
	if !ok {
		return fmt.Errorf("trigger_hotkey requires 'hotkey_name' parameter")
	}
	return e.obsClient.TriggerHotkeyByName(hotkeyName)
}

// setPreviewScene sets the preview scene in studio mode.
func (e *Executor) setPreviewScene(params map[string]interface{}) error {
	sceneName, ok := getStringParam(params, "scene_name")
	if !ok {
		return fmt.Errorf("set_preview_scene requires 'scene_name' parameter")
	}
	return e.obsClient.SetCurrentPreviewScene(sceneName)
}

// delay pauses execution for the specified duration.
func (e *Executor) delay(params map[string]interface{}) error {
	delayMs, ok := getIntParam(params, "delay_ms")
	if !ok {
		return fmt.Errorf("delay requires 'delay_ms' parameter")
	}

	if delayMs < 0 {
		return fmt.Errorf("delay_ms must be non-negative")
	}

	// Cap delay at 5 minutes to prevent excessive waits
	const maxDelayMs = 5 * 60 * 1000
	if delayMs > maxDelayMs {
		log.Printf("[Automation] Capping delay from %dms to %dms", delayMs, maxDelayMs)
		delayMs = maxDelayMs
	}

	time.Sleep(time.Duration(delayMs) * time.Millisecond)
	return nil
}

// Parameter extraction helpers

func getStringParam(params map[string]interface{}, key string) (string, bool) {
	if params == nil {
		return "", false
	}
	if v, ok := params[key].(string); ok {
		return v, true
	}
	return "", false
}

func getBoolParam(params map[string]interface{}, key string) (bool, bool) {
	if params == nil {
		return false, false
	}
	if v, ok := params[key].(bool); ok {
		return v, true
	}
	return false, false
}

func getFloat64Param(params map[string]interface{}, key string) (float64, bool) {
	if params == nil {
		return 0, false
	}
	// JSON numbers are float64
	if v, ok := params[key].(float64); ok {
		return v, true
	}
	// Also handle int for convenience
	if v, ok := params[key].(int); ok {
		return float64(v), true
	}
	return 0, false
}

func getIntParam(params map[string]interface{}, key string) (int, bool) {
	if params == nil {
		return 0, false
	}
	// JSON numbers are float64
	if v, ok := params[key].(float64); ok {
		return int(v), true
	}
	// Also handle int
	if v, ok := params[key].(int); ok {
		return v, true
	}
	return 0, false
}
