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
	// EventTypeSceneCreated is fired when a new scene is created
	EventTypeSceneCreated EventType = "scene_created"

	// EventTypeSceneRemoved is fired when a scene is removed/deleted
	EventTypeSceneRemoved EventType = "scene_removed"

	// EventTypeSceneChanged is fired when the current program scene changes
	EventTypeSceneChanged EventType = "scene_changed"
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

// ShouldTriggerListChanged returns true if the event type should trigger
// a "resources/list_changed" notification (scene creation or removal).
func ShouldTriggerListChanged(eventType EventType) bool {
	return eventType == EventTypeSceneCreated || eventType == EventTypeSceneRemoved
}

// ShouldTriggerResourceUpdated returns true if the event type should trigger
// a "resources/updated" notification for a specific resource (scene change).
func ShouldTriggerResourceUpdated(eventType EventType) bool {
	return eventType == EventTypeSceneChanged
}

// EventMetrics tracks statistics about OBS events for monitoring and debugging.
type EventMetrics struct {
	SceneCreatedCount int
	SceneRemovedCount int
	SceneChangedCount int
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
