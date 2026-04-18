// Package automation provides event-triggered actions and scheduling for OBS control.
package automation

import (
	"time"
)

// TriggerType defines the category of automation trigger.
const (
	TriggerTypeEvent    = "event"
	TriggerTypeSchedule = "schedule"
	TriggerTypeManual   = "manual"
)

// ActionType constants for all supported actions.
const (
	ActionTypeSetScene           = "set_scene"
	ActionTypeToggleMute         = "toggle_mute"
	ActionTypeSetMute            = "set_mute"
	ActionTypeSetVolume          = "set_volume"
	ActionTypeToggleVisibility   = "toggle_visibility"
	ActionTypeSetVisibility      = "set_visibility"
	ActionTypeStartRecording     = "start_recording"
	ActionTypeStopRecording      = "stop_recording"
	ActionTypePauseRecording     = "pause_recording"
	ActionTypeResumeRecording    = "resume_recording"
	ActionTypeStartStreaming     = "start_streaming"
	ActionTypeStopStreaming      = "stop_streaming"
	ActionTypeToggleVirtualCam   = "toggle_virtual_cam"
	ActionTypeStartVirtualCam    = "start_virtual_cam"
	ActionTypeStopVirtualCam     = "stop_virtual_cam"
	ActionTypeToggleReplayBuffer = "toggle_replay_buffer"
	ActionTypeSaveReplay         = "save_replay"
	ActionTypeTriggerHotkey      = "trigger_hotkey"
	ActionTypeTriggerTransition  = "trigger_transition"
	ActionTypeSetPreviewScene    = "set_preview_scene"
	ActionTypeDelay              = "delay"
)

// ActionErrorPolicy defines what to do when an action fails.
const (
	ActionErrorContinue = "continue" // Continue to next action (default)
	ActionErrorStop     = "stop"     // Stop rule execution
)

// EventType constants for OBS events that can trigger rules.
const (
	EventSceneChanged            = "scene_changed"
	EventSceneCreated            = "scene_created"
	EventSceneRemoved            = "scene_removed"
	EventRecordingStarted        = "recording_started"
	EventRecordingStopped        = "recording_stopped"
	EventRecordingPaused         = "recording_paused"
	EventRecordingResumed        = "recording_resumed"
	EventStreamingStarted        = "streaming_started"
	EventStreamingStopped        = "streaming_stopped"
	EventVirtualCamStarted       = "virtual_cam_started"
	EventVirtualCamStopped       = "virtual_cam_stopped"
	EventReplayBufferSaved       = "replay_buffer_saved"
	EventInputMuteChanged        = "input_mute_changed"
	EventSourceVisibilityChanged = "source_visibility_changed"
	EventTransitionStarted       = "transition_started"
	EventStudioModeChanged       = "studio_mode_changed"
)

// Rule represents an automation rule with trigger and actions.
type Rule struct {
	ID            int64                  `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description,omitempty"`
	Enabled       bool                   `json:"enabled"`
	TriggerType   string                 `json:"trigger_type"`
	TriggerConfig map[string]interface{} `json:"trigger_config"`
	Actions       []Action               `json:"actions"`
	CooldownMs    int                    `json:"cooldown_ms,omitempty"`
	Priority      int                    `json:"priority,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	LastRun       *time.Time             `json:"last_run,omitempty"`
	RunCount      int64                  `json:"run_count"`
}

// GetEventType returns the event_type from trigger config, or empty string.
func (r *Rule) GetEventType() string {
	if r.TriggerType != TriggerTypeEvent {
		return ""
	}
	if et, ok := r.TriggerConfig["event_type"].(string); ok {
		return et
	}
	return ""
}

// GetEventFilter returns the event_filter from trigger config, or nil.
func (r *Rule) GetEventFilter() map[string]interface{} {
	if filter, ok := r.TriggerConfig["event_filter"].(map[string]interface{}); ok {
		return filter
	}
	return nil
}

// GetSchedule returns the cron schedule from trigger config, or empty string.
func (r *Rule) GetSchedule() string {
	if r.TriggerType != TriggerTypeSchedule {
		return ""
	}
	if schedule, ok := r.TriggerConfig["schedule"].(string); ok {
		return schedule
	}
	return ""
}

// Action represents a single step in an automation sequence.
type Action struct {
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	OnError    string                 `json:"on_error,omitempty"` // "continue" or "stop"
}

// GetOnError returns the error policy, defaulting to "continue".
func (a *Action) GetOnError() string {
	if a.OnError == ActionErrorStop {
		return ActionErrorStop
	}
	return ActionErrorContinue
}

// EventPayload represents an OBS event that may trigger rules.
type EventPayload struct {
	EventType string                 `json:"event_type"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// ActionResult represents the result of executing a single action.
type ActionResult struct {
	ActionType string `json:"action_type"`
	Index      int    `json:"index"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	DurationMs int64  `json:"duration_ms"`
}

// ExecutionResult represents the complete result of rule execution.
type ExecutionResult struct {
	RuleID        int64                  `json:"rule_id"`
	RuleName      string                 `json:"rule_name"`
	TriggerType   string                 `json:"trigger_type"`
	TriggerData   map[string]interface{} `json:"trigger_data,omitempty"`
	StartedAt     time.Time              `json:"started_at"`
	CompletedAt   time.Time              `json:"completed_at"`
	Status        string                 `json:"status"` // "completed", "failed", "skipped"
	ActionResults []ActionResult         `json:"action_results,omitempty"`
	Error         string                 `json:"error,omitempty"`
	DurationMs    int64                  `json:"duration_ms"`
}

// SupportedEventTypes returns all event types that can be used as triggers.
func SupportedEventTypes() []string {
	return []string{
		EventSceneChanged,
		EventSceneCreated,
		EventSceneRemoved,
		EventRecordingStarted,
		EventRecordingStopped,
		EventRecordingPaused,
		EventRecordingResumed,
		EventStreamingStarted,
		EventStreamingStopped,
		EventVirtualCamStarted,
		EventVirtualCamStopped,
		EventReplayBufferSaved,
		EventInputMuteChanged,
		EventSourceVisibilityChanged,
		EventTransitionStarted,
		EventStudioModeChanged,
	}
}

// SupportedActionTypes returns all action types that can be used in rules.
func SupportedActionTypes() []string {
	return []string{
		ActionTypeSetScene,
		ActionTypeToggleMute,
		ActionTypeSetMute,
		ActionTypeSetVolume,
		ActionTypeToggleVisibility,
		ActionTypeSetVisibility,
		ActionTypeStartRecording,
		ActionTypeStopRecording,
		ActionTypePauseRecording,
		ActionTypeResumeRecording,
		ActionTypeStartStreaming,
		ActionTypeStopStreaming,
		ActionTypeToggleVirtualCam,
		ActionTypeStartVirtualCam,
		ActionTypeStopVirtualCam,
		ActionTypeToggleReplayBuffer,
		ActionTypeSaveReplay,
		ActionTypeTriggerHotkey,
		ActionTypeTriggerTransition,
		ActionTypeSetPreviewScene,
		ActionTypeDelay,
	}
}
