package obs

import (
	"fmt"
	"log"
)

// EventHandler implements the EventCallback interface and handles OBS events,
// dispatching them to the MCP server for resource notifications.
type EventHandler struct {
	// notificationFunc is called when an event occurs that should trigger
	// an MCP resource notification
	notificationFunc NotificationFunc
}

// NotificationFunc is the signature for functions that handle event notifications.
// The function receives the event type and relevant data (e.g., scene name).
type NotificationFunc func(eventType EventType, data map[string]interface{})

// EventType represents the type of OBS event that occurred.
type EventType string

const (
	// Scene events
	EventTypeSceneCreated EventType = "scene_created"
	EventTypeSceneRemoved EventType = "scene_removed"
	EventTypeSceneChanged EventType = "scene_changed"

	// Recording events
	EventTypeRecordingStarted     EventType = "recording_started"
	EventTypeRecordingStopped     EventType = "recording_stopped"
	EventTypeRecordingPaused      EventType = "recording_paused"
	EventTypeRecordingResumed     EventType = "recording_resumed"
	EventTypeRecordingFileChanged EventType = "recording_file_changed"

	// Streaming events
	EventTypeStreamingStarted EventType = "streaming_started"
	EventTypeStreamingStopped EventType = "streaming_stopped"

	// Virtual camera events
	EventTypeVirtualCamStarted EventType = "virtual_cam_started"
	EventTypeVirtualCamStopped EventType = "virtual_cam_stopped"

	// Replay buffer events
	EventTypeReplayBufferSaved EventType = "replay_buffer_saved"

	// Input events
	EventTypeInputMuteChanged EventType = "input_mute_changed"

	// Scene item events
	EventTypeSourceVisibilityChanged EventType = "source_visibility_changed"

	// Transition events
	EventTypeTransitionStarted EventType = "transition_started"

	// Studio mode events
	EventTypeStudioModeChanged EventType = "studio_mode_changed"

	// Canvas events (FB-42, OBS 30+ multi-canvas via obs-websocket 5.7+)
	EventTypeCanvasCreated     EventType = "canvas_created"
	EventTypeCanvasNameChanged EventType = "canvas_name_changed"
	EventTypeCanvasRemoved     EventType = "canvas_removed"
)

// NewEventHandler creates a new event handler with the specified notification function.
// The notification function will be called whenever an OBS event occurs that requires
// an MCP resource notification.
func NewEventHandler(notificationFunc NotificationFunc) *EventHandler {
	return &EventHandler{
		notificationFunc: notificationFunc,
	}
}

// OnSceneCreated is called when a new scene is created in OBS.
// This triggers a "resources/list_changed" notification to MCP clients.
func (h *EventHandler) OnSceneCreated(sceneName string) {
	log.Printf("[OBS Event] Scene created: %s", sceneName)

	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeSceneCreated, map[string]interface{}{
			"scene_name": sceneName,
			"action":     "created",
		})
	}
}

// OnSceneRemoved is called when a scene is removed/deleted from OBS.
// This triggers a "resources/list_changed" notification to MCP clients.
func (h *EventHandler) OnSceneRemoved(sceneName string) {
	log.Printf("[OBS Event] Scene removed: %s", sceneName)

	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeSceneRemoved, map[string]interface{}{
			"scene_name": sceneName,
			"action":     "removed",
		})
	}
}

// OnCurrentProgramSceneChanged is called when the active scene changes in OBS.
// This triggers a "resources/updated" notification for the specific scene URI.
func (h *EventHandler) OnCurrentProgramSceneChanged(sceneName string) {
	log.Printf("[OBS Event] Current program scene changed to: %s", sceneName)

	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeSceneChanged, map[string]interface{}{
			"scene_name": sceneName,
			"action":     "changed",
		})
	}
}

// OnRecordingStarted is called when recording begins.
func (h *EventHandler) OnRecordingStarted() {
	log.Printf("[OBS Event] Recording started")
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeRecordingStarted, map[string]interface{}{})
	}
}

// OnRecordingStopped is called when recording stops.
func (h *EventHandler) OnRecordingStopped(outputPath string) {
	log.Printf("[OBS Event] Recording stopped: %s", outputPath)
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeRecordingStopped, map[string]interface{}{
			"output_path": outputPath,
		})
	}
}

// OnRecordingPaused is called when recording is paused.
func (h *EventHandler) OnRecordingPaused() {
	log.Printf("[OBS Event] Recording paused")
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeRecordingPaused, map[string]interface{}{})
	}
}

// OnRecordingResumed is called when recording resumes after being paused.
func (h *EventHandler) OnRecordingResumed() {
	log.Printf("[OBS Event] Recording resumed")
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeRecordingResumed, map[string]interface{}{})
	}
}

// OnRecordingFileChanged is called when the record output rotates to a new file
// (e.g. OBS 30+ file splits).
func (h *EventHandler) OnRecordingFileChanged(newOutputPath string) {
	log.Printf("[OBS Event] Recording file changed: %s", newOutputPath)
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeRecordingFileChanged, map[string]interface{}{
			"new_output_path": newOutputPath,
		})
	}
}

// OnStreamingStarted is called when streaming begins.
func (h *EventHandler) OnStreamingStarted() {
	log.Printf("[OBS Event] Streaming started")
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeStreamingStarted, map[string]interface{}{})
	}
}

// OnStreamingStopped is called when streaming stops.
func (h *EventHandler) OnStreamingStopped() {
	log.Printf("[OBS Event] Streaming stopped")
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeStreamingStopped, map[string]interface{}{})
	}
}

// OnVirtualCamStarted is called when the virtual camera is started.
func (h *EventHandler) OnVirtualCamStarted() {
	log.Printf("[OBS Event] Virtual camera started")
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeVirtualCamStarted, map[string]interface{}{})
	}
}

// OnVirtualCamStopped is called when the virtual camera is stopped.
func (h *EventHandler) OnVirtualCamStopped() {
	log.Printf("[OBS Event] Virtual camera stopped")
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeVirtualCamStopped, map[string]interface{}{})
	}
}

// OnReplayBufferSaved is called when a replay buffer is saved.
func (h *EventHandler) OnReplayBufferSaved(savedPath string) {
	log.Printf("[OBS Event] Replay buffer saved: %s", savedPath)
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeReplayBufferSaved, map[string]interface{}{
			"saved_path": savedPath,
		})
	}
}

// OnInputMuteChanged is called when an input's mute state changes.
func (h *EventHandler) OnInputMuteChanged(inputName string, muted bool) {
	log.Printf("[OBS Event] Input mute changed: %s = %v", inputName, muted)
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeInputMuteChanged, map[string]interface{}{
			"input_name": inputName,
			"muted":      muted,
		})
	}
}

// OnSceneItemVisibilityChanged is called when a scene item's visibility changes.
func (h *EventHandler) OnSceneItemVisibilityChanged(sceneName string, sceneItemId int, visible bool) {
	log.Printf("[OBS Event] Scene item visibility changed: %s item %d = %v", sceneName, sceneItemId, visible)
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeSourceVisibilityChanged, map[string]interface{}{
			"scene_name":    sceneName,
			"scene_item_id": sceneItemId,
			"visible":       visible,
		})
	}
}

// OnTransitionStarted is called when a scene transition starts.
func (h *EventHandler) OnTransitionStarted(transitionName string) {
	log.Printf("[OBS Event] Transition started: %s", transitionName)
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeTransitionStarted, map[string]interface{}{
			"transition_name": transitionName,
		})
	}
}

// OnStudioModeChanged is called when studio mode is enabled or disabled.
func (h *EventHandler) OnStudioModeChanged(enabled bool) {
	log.Printf("[OBS Event] Studio mode changed: %v", enabled)
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeStudioModeChanged, map[string]interface{}{
			"enabled": enabled,
		})
	}
}

// OnCanvasCreated is called when a new canvas is created in OBS (FB-42).
// Triggers a "resources/list_changed" notification so clients re-list canvases.
func (h *EventHandler) OnCanvasCreated(canvasName, canvasUUID string) {
	log.Printf("[OBS Event] Canvas created: %s (uuid=%s)", canvasName, canvasUUID)
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeCanvasCreated, map[string]interface{}{
			"canvas_name": canvasName,
			"canvas_uuid": canvasUUID,
			"action":      "created",
		})
	}
}

// OnCanvasNameChanged is called when a canvas is renamed in OBS (FB-42).
// Triggers a "resources/updated" notification for the renamed canvas URI.
func (h *EventHandler) OnCanvasNameChanged(oldName, newName, canvasUUID string) {
	log.Printf("[OBS Event] Canvas renamed: %s -> %s (uuid=%s)", oldName, newName, canvasUUID)
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeCanvasNameChanged, map[string]interface{}{
			"old_canvas_name": oldName,
			"canvas_name":     newName,
			"canvas_uuid":     canvasUUID,
			"action":          "renamed",
		})
	}
}

// OnCanvasRemoved is called when a canvas is removed from OBS (FB-42).
// Triggers a "resources/list_changed" notification so clients re-list canvases.
func (h *EventHandler) OnCanvasRemoved(canvasName, canvasUUID string) {
	log.Printf("[OBS Event] Canvas removed: %s (uuid=%s)", canvasName, canvasUUID)
	if h.notificationFunc != nil {
		h.notificationFunc(EventTypeCanvasRemoved, map[string]interface{}{
			"canvas_name": canvasName,
			"canvas_uuid": canvasUUID,
			"action":      "removed",
		})
	}
}

// EventLogger is a simple event callback implementation that just logs events
// without triggering MCP notifications. Useful for testing and debugging.
type EventLogger struct{}

// NewEventLogger creates a new event logger.
func NewEventLogger() *EventLogger {
	return &EventLogger{}
}

// OnSceneCreated logs scene creation events.
func (l *EventLogger) OnSceneCreated(sceneName string) {
	log.Printf("[OBS Event Logger] Scene created: %s", sceneName)
}

// OnSceneRemoved logs scene removal events.
func (l *EventLogger) OnSceneRemoved(sceneName string) {
	log.Printf("[OBS Event Logger] Scene removed: %s", sceneName)
}

// OnCurrentProgramSceneChanged logs scene change events.
func (l *EventLogger) OnCurrentProgramSceneChanged(sceneName string) {
	log.Printf("[OBS Event Logger] Current program scene changed to: %s", sceneName)
}

// OnRecordingStarted logs recording start events.
func (l *EventLogger) OnRecordingStarted() {
	log.Printf("[OBS Event Logger] Recording started")
}

// OnRecordingStopped logs recording stop events.
func (l *EventLogger) OnRecordingStopped(outputPath string) {
	log.Printf("[OBS Event Logger] Recording stopped: %s", outputPath)
}

// OnRecordingPaused logs recording pause events.
func (l *EventLogger) OnRecordingPaused() {
	log.Printf("[OBS Event Logger] Recording paused")
}

// OnRecordingResumed logs recording resume events.
func (l *EventLogger) OnRecordingResumed() {
	log.Printf("[OBS Event Logger] Recording resumed")
}

// OnRecordingFileChanged logs recording file-rotation events.
func (l *EventLogger) OnRecordingFileChanged(newOutputPath string) {
	log.Printf("[OBS Event Logger] Recording file changed: %s", newOutputPath)
}

// OnStreamingStarted logs streaming start events.
func (l *EventLogger) OnStreamingStarted() {
	log.Printf("[OBS Event Logger] Streaming started")
}

// OnStreamingStopped logs streaming stop events.
func (l *EventLogger) OnStreamingStopped() {
	log.Printf("[OBS Event Logger] Streaming stopped")
}

// OnVirtualCamStarted logs virtual camera start events.
func (l *EventLogger) OnVirtualCamStarted() {
	log.Printf("[OBS Event Logger] Virtual camera started")
}

// OnVirtualCamStopped logs virtual camera stop events.
func (l *EventLogger) OnVirtualCamStopped() {
	log.Printf("[OBS Event Logger] Virtual camera stopped")
}

// OnReplayBufferSaved logs replay buffer saved events.
func (l *EventLogger) OnReplayBufferSaved(savedPath string) {
	log.Printf("[OBS Event Logger] Replay buffer saved: %s", savedPath)
}

// OnInputMuteChanged logs input mute change events.
func (l *EventLogger) OnInputMuteChanged(inputName string, muted bool) {
	log.Printf("[OBS Event Logger] Input mute changed: %s = %v", inputName, muted)
}

// OnSceneItemVisibilityChanged logs scene item visibility change events.
func (l *EventLogger) OnSceneItemVisibilityChanged(sceneName string, sceneItemId int, visible bool) {
	log.Printf("[OBS Event Logger] Scene item visibility changed: %s item %d = %v", sceneName, sceneItemId, visible)
}

// OnTransitionStarted logs transition start events.
func (l *EventLogger) OnTransitionStarted(transitionName string) {
	log.Printf("[OBS Event Logger] Transition started: %s", transitionName)
}

// OnStudioModeChanged logs studio mode change events.
func (l *EventLogger) OnStudioModeChanged(enabled bool) {
	log.Printf("[OBS Event Logger] Studio mode changed: %v", enabled)
}

// OnCanvasCreated logs canvas creation events.
func (l *EventLogger) OnCanvasCreated(canvasName, canvasUUID string) {
	log.Printf("[OBS Event Logger] Canvas created: %s (uuid=%s)", canvasName, canvasUUID)
}

// OnCanvasNameChanged logs canvas rename events.
func (l *EventLogger) OnCanvasNameChanged(oldName, newName, canvasUUID string) {
	log.Printf("[OBS Event Logger] Canvas renamed: %s -> %s (uuid=%s)", oldName, newName, canvasUUID)
}

// OnCanvasRemoved logs canvas removal events.
func (l *EventLogger) OnCanvasRemoved(canvasName, canvasUUID string) {
	log.Printf("[OBS Event Logger] Canvas removed: %s (uuid=%s)", canvasName, canvasUUID)
}

// FormatEventNotification formats an event into a structured notification message
// suitable for MCP resource notifications.
func FormatEventNotification(eventType EventType, data map[string]interface{}) (string, error) {
	sceneName, ok := data["scene_name"].(string)
	if !ok {
		return "", fmt.Errorf("scene_name not found in event data")
	}

	switch eventType {
	case EventTypeSceneCreated:
		return fmt.Sprintf("Scene '%s' was created in OBS", sceneName), nil

	case EventTypeSceneRemoved:
		return fmt.Sprintf("Scene '%s' was removed from OBS", sceneName), nil

	case EventTypeSceneChanged:
		return fmt.Sprintf("OBS switched to scene '%s'", sceneName), nil

	default:
		return "", fmt.Errorf("unknown event type: %s", eventType)
	}
}

// GetResourceURIForScene returns the MCP resource URI for a given scene name.
// This follows the pattern: obs://scene/{scene_name}
func GetResourceURIForScene(sceneName string) string {
	return fmt.Sprintf("obs://scene/%s", sceneName)
}

// GetResourceURIForCanvas returns the MCP resource URI for a given canvas name.
// This follows the pattern: obs://canvas/{canvas_name}
func GetResourceURIForCanvas(canvasName string) string {
	return fmt.Sprintf("obs://canvas/%s", canvasName)
}

// ShouldTriggerListChanged returns true if the event type should trigger
// a "resources/list_changed" notification (scene/canvas creation or removal).
func ShouldTriggerListChanged(eventType EventType) bool {
	switch eventType {
	case EventTypeSceneCreated, EventTypeSceneRemoved,
		EventTypeCanvasCreated, EventTypeCanvasRemoved:
		return true
	}
	return false
}

// ShouldTriggerResourceUpdated returns true if the event type should trigger
// a "resources/updated" notification for a specific resource (scene change or canvas rename).
func ShouldTriggerResourceUpdated(eventType EventType) bool {
	return eventType == EventTypeSceneChanged || eventType == EventTypeCanvasNameChanged
}

// EventMetrics tracks statistics about OBS events for monitoring and debugging.
type EventMetrics struct {
	SceneCreatedCount            int
	SceneRemovedCount            int
	SceneChangedCount            int
	RecordingStartedCount        int
	RecordingStoppedCount        int
	RecordingPausedCount         int
	RecordingResumedCount        int
	RecordingFileChangedCount    int
	StreamingStartedCount        int
	StreamingStoppedCount        int
	VirtualCamStartedCount       int
	VirtualCamStoppedCount       int
	ReplayBufferSavedCount       int
	InputMuteChangedCount        int
	SourceVisibilityChangedCount int
	TransitionStartedCount       int
	StudioModeChangedCount       int
	CanvasCreatedCount           int
	CanvasNameChangedCount       int
	CanvasRemovedCount           int
}

// EventMetricsTracker is an event callback that tracks event counts.
type EventMetricsTracker struct {
	metrics EventMetrics
}

// NewEventMetricsTracker creates a new metrics tracker.
func NewEventMetricsTracker() *EventMetricsTracker {
	return &EventMetricsTracker{
		metrics: EventMetrics{},
	}
}

// OnSceneCreated increments the scene created counter.
func (t *EventMetricsTracker) OnSceneCreated(sceneName string) {
	t.metrics.SceneCreatedCount++
	log.Printf("[Metrics] Scene created: %s (total: %d)", sceneName, t.metrics.SceneCreatedCount)
}

// OnSceneRemoved increments the scene removed counter.
func (t *EventMetricsTracker) OnSceneRemoved(sceneName string) {
	t.metrics.SceneRemovedCount++
	log.Printf("[Metrics] Scene removed: %s (total: %d)", sceneName, t.metrics.SceneRemovedCount)
}

// OnCurrentProgramSceneChanged increments the scene changed counter.
func (t *EventMetricsTracker) OnCurrentProgramSceneChanged(sceneName string) {
	t.metrics.SceneChangedCount++
	log.Printf("[Metrics] Scene changed: %s (total: %d)", sceneName, t.metrics.SceneChangedCount)
}

// OnRecordingStarted increments the recording started counter.
func (t *EventMetricsTracker) OnRecordingStarted() {
	t.metrics.RecordingStartedCount++
}

// OnRecordingStopped increments the recording stopped counter.
func (t *EventMetricsTracker) OnRecordingStopped(outputPath string) {
	t.metrics.RecordingStoppedCount++
}

// OnRecordingPaused increments the recording paused counter.
func (t *EventMetricsTracker) OnRecordingPaused() {
	t.metrics.RecordingPausedCount++
}

// OnRecordingResumed increments the recording resumed counter.
func (t *EventMetricsTracker) OnRecordingResumed() {
	t.metrics.RecordingResumedCount++
}

// OnRecordingFileChanged increments the recording file-changed counter.
func (t *EventMetricsTracker) OnRecordingFileChanged(newOutputPath string) {
	t.metrics.RecordingFileChangedCount++
}

// OnStreamingStarted increments the streaming started counter.
func (t *EventMetricsTracker) OnStreamingStarted() {
	t.metrics.StreamingStartedCount++
}

// OnStreamingStopped increments the streaming stopped counter.
func (t *EventMetricsTracker) OnStreamingStopped() {
	t.metrics.StreamingStoppedCount++
}

// OnVirtualCamStarted increments the virtual cam started counter.
func (t *EventMetricsTracker) OnVirtualCamStarted() {
	t.metrics.VirtualCamStartedCount++
}

// OnVirtualCamStopped increments the virtual cam stopped counter.
func (t *EventMetricsTracker) OnVirtualCamStopped() {
	t.metrics.VirtualCamStoppedCount++
}

// OnReplayBufferSaved increments the replay buffer saved counter.
func (t *EventMetricsTracker) OnReplayBufferSaved(savedPath string) {
	t.metrics.ReplayBufferSavedCount++
}

// OnInputMuteChanged increments the input mute changed counter.
func (t *EventMetricsTracker) OnInputMuteChanged(inputName string, muted bool) {
	t.metrics.InputMuteChangedCount++
}

// OnSceneItemVisibilityChanged increments the source visibility changed counter.
func (t *EventMetricsTracker) OnSceneItemVisibilityChanged(sceneName string, sceneItemId int, visible bool) {
	t.metrics.SourceVisibilityChangedCount++
}

// OnTransitionStarted increments the transition started counter.
func (t *EventMetricsTracker) OnTransitionStarted(transitionName string) {
	t.metrics.TransitionStartedCount++
}

// OnStudioModeChanged increments the studio mode changed counter.
func (t *EventMetricsTracker) OnStudioModeChanged(enabled bool) {
	t.metrics.StudioModeChangedCount++
}

// OnCanvasCreated increments the canvas created counter.
func (t *EventMetricsTracker) OnCanvasCreated(canvasName, canvasUUID string) {
	t.metrics.CanvasCreatedCount++
}

// OnCanvasNameChanged increments the canvas renamed counter.
func (t *EventMetricsTracker) OnCanvasNameChanged(oldName, newName, canvasUUID string) {
	t.metrics.CanvasNameChangedCount++
}

// OnCanvasRemoved increments the canvas removed counter.
func (t *EventMetricsTracker) OnCanvasRemoved(canvasName, canvasUUID string) {
	t.metrics.CanvasRemovedCount++
}

// GetMetrics returns the current event metrics.
func (t *EventMetricsTracker) GetMetrics() EventMetrics {
	return t.metrics
}

// ResetMetrics resets all event counters to zero.
func (t *EventMetricsTracker) ResetMetrics() {
	t.metrics = EventMetrics{}
}

// CompositeEventCallback allows multiple event callbacks to be registered.
// When an event occurs, all registered callbacks are invoked in sequence.
type CompositeEventCallback struct {
	callbacks []EventCallback
}

// NewCompositeEventCallback creates a new composite callback with the given callbacks.
func NewCompositeEventCallback(callbacks ...EventCallback) *CompositeEventCallback {
	return &CompositeEventCallback{
		callbacks: callbacks,
	}
}

// AddCallback adds a new callback to the composite.
func (c *CompositeEventCallback) AddCallback(callback EventCallback) {
	c.callbacks = append(c.callbacks, callback)
}

// OnSceneCreated dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnSceneCreated(sceneName string) {
	for _, callback := range c.callbacks {
		callback.OnSceneCreated(sceneName)
	}
}

// OnSceneRemoved dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnSceneRemoved(sceneName string) {
	for _, callback := range c.callbacks {
		callback.OnSceneRemoved(sceneName)
	}
}

// OnCurrentProgramSceneChanged dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnCurrentProgramSceneChanged(sceneName string) {
	for _, callback := range c.callbacks {
		callback.OnCurrentProgramSceneChanged(sceneName)
	}
}

// OnRecordingStarted dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnRecordingStarted() {
	for _, callback := range c.callbacks {
		callback.OnRecordingStarted()
	}
}

// OnRecordingStopped dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnRecordingStopped(outputPath string) {
	for _, callback := range c.callbacks {
		callback.OnRecordingStopped(outputPath)
	}
}

// OnRecordingPaused dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnRecordingPaused() {
	for _, callback := range c.callbacks {
		callback.OnRecordingPaused()
	}
}

// OnRecordingResumed dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnRecordingResumed() {
	for _, callback := range c.callbacks {
		callback.OnRecordingResumed()
	}
}

// OnRecordingFileChanged dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnRecordingFileChanged(newOutputPath string) {
	for _, callback := range c.callbacks {
		callback.OnRecordingFileChanged(newOutputPath)
	}
}

// OnStreamingStarted dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnStreamingStarted() {
	for _, callback := range c.callbacks {
		callback.OnStreamingStarted()
	}
}

// OnStreamingStopped dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnStreamingStopped() {
	for _, callback := range c.callbacks {
		callback.OnStreamingStopped()
	}
}

// OnVirtualCamStarted dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnVirtualCamStarted() {
	for _, callback := range c.callbacks {
		callback.OnVirtualCamStarted()
	}
}

// OnVirtualCamStopped dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnVirtualCamStopped() {
	for _, callback := range c.callbacks {
		callback.OnVirtualCamStopped()
	}
}

// OnReplayBufferSaved dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnReplayBufferSaved(savedPath string) {
	for _, callback := range c.callbacks {
		callback.OnReplayBufferSaved(savedPath)
	}
}

// OnInputMuteChanged dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnInputMuteChanged(inputName string, muted bool) {
	for _, callback := range c.callbacks {
		callback.OnInputMuteChanged(inputName, muted)
	}
}

// OnSceneItemVisibilityChanged dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnSceneItemVisibilityChanged(sceneName string, sceneItemId int, visible bool) {
	for _, callback := range c.callbacks {
		callback.OnSceneItemVisibilityChanged(sceneName, sceneItemId, visible)
	}
}

// OnTransitionStarted dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnTransitionStarted(transitionName string) {
	for _, callback := range c.callbacks {
		callback.OnTransitionStarted(transitionName)
	}
}

// OnStudioModeChanged dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnStudioModeChanged(enabled bool) {
	for _, callback := range c.callbacks {
		callback.OnStudioModeChanged(enabled)
	}
}

// OnCanvasCreated dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnCanvasCreated(canvasName, canvasUUID string) {
	for _, callback := range c.callbacks {
		callback.OnCanvasCreated(canvasName, canvasUUID)
	}
}

// OnCanvasNameChanged dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnCanvasNameChanged(oldName, newName, canvasUUID string) {
	for _, callback := range c.callbacks {
		callback.OnCanvasNameChanged(oldName, newName, canvasUUID)
	}
}

// OnCanvasRemoved dispatches to all registered callbacks.
func (c *CompositeEventCallback) OnCanvasRemoved(canvasName, canvasUUID string) {
	for _, callback := range c.callbacks {
		callback.OnCanvasRemoved(canvasName, canvasUUID)
	}
}
