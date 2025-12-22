package mcp

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ironystock/agentic-obs/internal/obs"
	"github.com/ironystock/agentic-obs/internal/storage"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool input/output types

// SceneNameInput is the input for scene operations
type SceneNameInput struct {
	SceneName string `json:"scene_name"`
}

// SimpleResult is a simple text result
type SimpleResult struct {
	Message string `json:"message"`
}

// SourceNameInput is the input for source-related operations
type SourceNameInput struct {
	SourceName string `json:"source_name"`
}

// SourceVisibilityInput is the input for toggling source visibility
type SourceVisibilityInput struct {
	SceneName string `json:"scene_name"`
	SourceID  int64  `json:"source_id"`
}

// InputNameInput is the input for audio input operations
type InputNameInput struct {
	InputName string `json:"input_name"`
}

// SetVolumeInput is the input for setting audio input volume
type SetVolumeInput struct {
	InputName string   `json:"input_name"`
	VolumeDb  *float64 `json:"volume_db,omitempty"`
	VolumeMul *float64 `json:"volume_mul,omitempty"`
}

// ListPresetsInput is the input for listing scene presets
type ListPresetsInput struct {
	SceneName string `json:"scene_name,omitempty" jsonschema:"Optional scene name to filter presets by"`
}

// PresetNameInput is the input for preset operations by name
type PresetNameInput struct {
	PresetName string `json:"preset_name" jsonschema:"Name of the preset to operate on"`
}

// RenamePresetInput is the input for renaming a preset
type RenamePresetInput struct {
	OldName string `json:"old_name" jsonschema:"Current name of the preset to rename"`
	NewName string `json:"new_name" jsonschema:"New name for the preset"`
}

// SavePresetInput is the input for saving a scene preset
type SavePresetInput struct {
	PresetName string `json:"preset_name" jsonschema:"Name to give the new preset"`
	SceneName  string `json:"scene_name" jsonschema:"Name of the OBS scene to capture state from"`
}

// CreateScreenshotSourceInput is the input for creating a screenshot source
type CreateScreenshotSourceInput struct {
	Name        string `json:"name" jsonschema:"Unique name for this screenshot source"`
	SourceName  string `json:"source_name" jsonschema:"OBS scene or source name to capture"`
	CadenceMs   int    `json:"cadence_ms,omitempty" jsonschema:"Capture interval in milliseconds (default: 5000)"`
	ImageFormat string `json:"image_format,omitempty" jsonschema:"Image format: png or jpg (default: png)"`
	ImageWidth  int    `json:"image_width,omitempty" jsonschema:"Optional resize width (0 = original)"`
	ImageHeight int    `json:"image_height,omitempty" jsonschema:"Optional resize height (0 = original)"`
	Quality     int    `json:"quality,omitempty" jsonschema:"Compression quality 0-100 (default: 80)"`
}

// ScreenshotSourceNameInput is the input for screenshot source operations by name
type ScreenshotSourceNameInput struct {
	Name string `json:"name" jsonschema:"Name of the screenshot source"`
}

// ConfigureScreenshotCadenceInput is the input for updating screenshot cadence
type ConfigureScreenshotCadenceInput struct {
	Name      string `json:"name" jsonschema:"Name of the screenshot source"`
	CadenceMs int    `json:"cadence_ms" jsonschema:"New capture interval in milliseconds"`
}

// Design tool input types

// CreateTextSourceInput is the input for creating a text source
type CreateTextSourceInput struct {
	SceneName  string `json:"scene_name" jsonschema:"Name of the scene to add the source to"`
	SourceName string `json:"source_name" jsonschema:"Name for the new text source"`
	Text       string `json:"text" jsonschema:"Text content to display"`
	FontName   string `json:"font_name,omitempty" jsonschema:"Font face name (default: Arial)"`
	FontSize   int    `json:"font_size,omitempty" jsonschema:"Font size in points (default: 36)"`
	Color      int64  `json:"color,omitempty" jsonschema:"Text color as ABGR integer (default: white)"`
}

// CreateImageSourceInput is the input for creating an image source
type CreateImageSourceInput struct {
	SceneName  string `json:"scene_name" jsonschema:"Name of the scene to add the source to"`
	SourceName string `json:"source_name" jsonschema:"Name for the new image source"`
	FilePath   string `json:"file_path" jsonschema:"Path to the image file"`
}

// CreateColorSourceInput is the input for creating a color source
type CreateColorSourceInput struct {
	SceneName  string `json:"scene_name" jsonschema:"Name of the scene to add the source to"`
	SourceName string `json:"source_name" jsonschema:"Name for the new color source"`
	Color      int64  `json:"color" jsonschema:"Color as ABGR integer (e.g., 0xFF0000FF for red)"`
	Width      int    `json:"width,omitempty" jsonschema:"Width in pixels (default: 1920)"`
	Height     int    `json:"height,omitempty" jsonschema:"Height in pixels (default: 1080)"`
}

// CreateBrowserSourceInput is the input for creating a browser source
type CreateBrowserSourceInput struct {
	SceneName  string `json:"scene_name" jsonschema:"Name of the scene to add the source to"`
	SourceName string `json:"source_name" jsonschema:"Name for the new browser source"`
	URL        string `json:"url" jsonschema:"URL to load in the browser source"`
	Width      int    `json:"width,omitempty" jsonschema:"Browser width in pixels (default: 800)"`
	Height     int    `json:"height,omitempty" jsonschema:"Browser height in pixels (default: 600)"`
	FPS        int    `json:"fps,omitempty" jsonschema:"Frame rate (default: 30)"`
}

// CreateMediaSourceInput is the input for creating a media/video source
type CreateMediaSourceInput struct {
	SceneName  string `json:"scene_name" jsonschema:"Name of the scene to add the source to"`
	SourceName string `json:"source_name" jsonschema:"Name for the new media source"`
	FilePath   string `json:"file_path" jsonschema:"Path to the media file"`
	Loop       bool   `json:"loop,omitempty" jsonschema:"Whether to loop the media (default: false)"`
}

// SetSourceTransformInput is the input for setting source transform properties
type SetSourceTransformInput struct {
	SceneName   string   `json:"scene_name" jsonschema:"Name of the scene containing the source"`
	SceneItemID int      `json:"scene_item_id" jsonschema:"Scene item ID of the source"`
	X           *float64 `json:"x,omitempty" jsonschema:"X position in pixels"`
	Y           *float64 `json:"y,omitempty" jsonschema:"Y position in pixels"`
	ScaleX      *float64 `json:"scale_x,omitempty" jsonschema:"X scale factor (1.0 = 100%)"`
	ScaleY      *float64 `json:"scale_y,omitempty" jsonschema:"Y scale factor (1.0 = 100%)"`
	Rotation    *float64 `json:"rotation,omitempty" jsonschema:"Rotation in degrees"`
}

// GetSourceTransformInput is the input for getting source transform properties
type GetSourceTransformInput struct {
	SceneName   string `json:"scene_name" jsonschema:"Name of the scene containing the source"`
	SceneItemID int    `json:"scene_item_id" jsonschema:"Scene item ID of the source"`
}

// SetSourceCropInput is the input for setting source crop
type SetSourceCropInput struct {
	SceneName   string `json:"scene_name" jsonschema:"Name of the scene containing the source"`
	SceneItemID int    `json:"scene_item_id" jsonschema:"Scene item ID of the source"`
	CropTop     int    `json:"crop_top,omitempty" jsonschema:"Pixels to crop from top"`
	CropBottom  int    `json:"crop_bottom,omitempty" jsonschema:"Pixels to crop from bottom"`
	CropLeft    int    `json:"crop_left,omitempty" jsonschema:"Pixels to crop from left"`
	CropRight   int    `json:"crop_right,omitempty" jsonschema:"Pixels to crop from right"`
}

// SetSourceBoundsInput is the input for setting source bounds
type SetSourceBoundsInput struct {
	SceneName    string  `json:"scene_name" jsonschema:"Name of the scene containing the source"`
	SceneItemID  int     `json:"scene_item_id" jsonschema:"Scene item ID of the source"`
	BoundsType   string  `json:"bounds_type" jsonschema:"Bounds type: OBS_BOUNDS_NONE, OBS_BOUNDS_STRETCH, OBS_BOUNDS_SCALE_INNER, OBS_BOUNDS_SCALE_OUTER, OBS_BOUNDS_SCALE_TO_WIDTH, OBS_BOUNDS_SCALE_TO_HEIGHT, OBS_BOUNDS_MAX_ONLY"`
	BoundsWidth  float64 `json:"bounds_width,omitempty" jsonschema:"Bounds width in pixels"`
	BoundsHeight float64 `json:"bounds_height,omitempty" jsonschema:"Bounds height in pixels"`
}

// SetSourceOrderInput is the input for setting source z-order
type SetSourceOrderInput struct {
	SceneName   string `json:"scene_name" jsonschema:"Name of the scene containing the source"`
	SceneItemID int    `json:"scene_item_id" jsonschema:"Scene item ID of the source"`
	Index       int    `json:"index" jsonschema:"New index position (0 = bottom, higher = front)"`
}

// SetSourceLockedInput is the input for locking/unlocking a source
type SetSourceLockedInput struct {
	SceneName   string `json:"scene_name" jsonschema:"Name of the scene containing the source"`
	SceneItemID int    `json:"scene_item_id" jsonschema:"Scene item ID of the source"`
	Locked      bool   `json:"locked" jsonschema:"Whether the source should be locked"`
}

// DuplicateSourceInput is the input for duplicating a source
type DuplicateSourceInput struct {
	SceneName     string `json:"scene_name" jsonschema:"Name of the scene containing the source"`
	SceneItemID   int    `json:"scene_item_id" jsonschema:"Scene item ID of the source to duplicate"`
	DestSceneName string `json:"dest_scene_name,omitempty" jsonschema:"Destination scene name (default: same scene)"`
}

// RemoveSourceInput is the input for removing a source from a scene
type RemoveSourceInput struct {
	SceneName   string `json:"scene_name" jsonschema:"Name of the scene containing the source"`
	SceneItemID int    `json:"scene_item_id" jsonschema:"Scene item ID of the source to remove"`
}

// Filter tool input types (FB-23)

// ListSourceFiltersInput is the input for listing filters on a source
type ListSourceFiltersInput struct {
	SourceName string `json:"source_name" jsonschema:"Name of the source to list filters for"`
}

// GetSourceFilterInput is the input for getting a specific filter's details
type GetSourceFilterInput struct {
	SourceName string `json:"source_name" jsonschema:"Name of the source containing the filter"`
	FilterName string `json:"filter_name" jsonschema:"Name of the filter to get details for"`
}

// CreateSourceFilterInput is the input for creating a new filter on a source
type CreateSourceFilterInput struct {
	SourceName     string                 `json:"source_name" jsonschema:"Name of the source to add the filter to"`
	FilterName     string                 `json:"filter_name" jsonschema:"Name for the new filter"`
	FilterKind     string                 `json:"filter_kind" jsonschema:"Type of filter (use list_filter_kinds to see available types)"`
	FilterSettings map[string]interface{} `json:"filter_settings,omitempty" jsonschema:"Optional initial settings for the filter"`
}

// RemoveSourceFilterInput is the input for removing a filter from a source
type RemoveSourceFilterInput struct {
	SourceName string `json:"source_name" jsonschema:"Name of the source containing the filter"`
	FilterName string `json:"filter_name" jsonschema:"Name of the filter to remove"`
}

// ToggleSourceFilterInput is the input for enabling/disabling a filter
type ToggleSourceFilterInput struct {
	SourceName    string `json:"source_name" jsonschema:"Name of the source containing the filter"`
	FilterName    string `json:"filter_name" jsonschema:"Name of the filter to toggle"`
	FilterEnabled *bool  `json:"filter_enabled,omitempty" jsonschema:"Set to true/false to enable/disable; omit to toggle"`
}

// SetSourceFilterSettingsInput is the input for updating filter settings
type SetSourceFilterSettingsInput struct {
	SourceName     string                 `json:"source_name" jsonschema:"Name of the source containing the filter"`
	FilterName     string                 `json:"filter_name" jsonschema:"Name of the filter to update"`
	FilterSettings map[string]interface{} `json:"filter_settings" jsonschema:"Settings to apply to the filter"`
	Overlay        bool                   `json:"overlay,omitempty" jsonschema:"If true, merge with existing settings; if false, replace entirely (default: true)"`
}

// Transition tool input types (FB-24)

// SetCurrentTransitionInput is the input for setting the current scene transition
type SetCurrentTransitionInput struct {
	TransitionName string `json:"transition_name" jsonschema:"Name of the transition to set as current"`
}

// SetTransitionDurationInput is the input for setting the transition duration
type SetTransitionDurationInput struct {
	TransitionDuration int `json:"transition_duration" jsonschema:"Transition duration in milliseconds"`
}

// Virtual Camera and Replay Buffer input types (FB-25)

// (no input needed for get_virtual_cam_status, toggle_virtual_cam, get_replay_buffer_status, toggle_replay_buffer, save_replay_buffer, get_last_replay)

// Studio Mode input types (FB-26)

// SetStudioModeInput is the input for enabling/disabling studio mode
type SetStudioModeInput struct {
	StudioModeEnabled bool `json:"studio_mode_enabled" jsonschema:"True to enable studio mode, false to disable"`
}

// SetPreviewSceneInput is the input for setting the preview scene in studio mode
type SetPreviewSceneInput struct {
	SceneName string `json:"scene_name" jsonschema:"Name of the scene to set as preview"`
}

// Hotkey input types (FB-26)

// TriggerHotkeyInput is the input for triggering a hotkey by name
type TriggerHotkeyInput struct {
	HotkeyName string `json:"hotkey_name" jsonschema:"Name of the hotkey to trigger (use list_hotkeys to see available hotkeys)"`
}

// registerToolHandlers registers MCP tool handlers based on enabled tool groups
func (s *Server) registerToolHandlers() {
	toolCount := 0

	// Core tools: Scene management, Recording, Streaming, Status
	if s.toolGroups.Core {
		// Scene management tools
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "set_current_scene",
				Description: "Switch to a different scene in OBS",
			},
			s.handleSetCurrentScene,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "create_scene",
				Description: "Create a new scene in OBS",
			},
			s.handleCreateScene,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "remove_scene",
				Description: "Remove a scene from OBS",
			},
			s.handleRemoveScene,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "list_scenes",
				Description: "List all available scenes in OBS and identify the current scene",
			},
			s.handleListScenes,
		)

		// Recording tools
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "start_recording",
				Description: "Start recording in OBS",
			},
			s.handleStartRecording,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "stop_recording",
				Description: "Stop the current recording in OBS",
			},
			s.handleStopRecording,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_recording_status",
				Description: "Get the current recording status from OBS",
			},
			s.handleGetRecordingStatus,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "pause_recording",
				Description: "Pause the current recording in OBS (recording must be active)",
			},
			s.handlePauseRecording,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "resume_recording",
				Description: "Resume a paused recording in OBS (recording must be paused)",
			},
			s.handleResumeRecording,
		)

		// Streaming tools
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "start_streaming",
				Description: "Start streaming in OBS",
			},
			s.handleStartStreaming,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "stop_streaming",
				Description: "Stop the current stream in OBS",
			},
			s.handleStopStreaming,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_streaming_status",
				Description: "Get the current streaming status from OBS",
			},
			s.handleGetStreamingStatus,
		)

		// Status tool
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_obs_status",
				Description: "Get overall OBS status including version, connection state, and active scene",
			},
			s.handleGetOBSStatus,
		)

		// Virtual camera tools (FB-25)
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_virtual_cam_status",
				Description: "Check if the virtual camera is currently active",
			},
			s.handleGetVirtualCamStatus,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "toggle_virtual_cam",
				Description: "Start or stop the virtual camera",
			},
			s.handleToggleVirtualCam,
		)

		// Replay buffer tools (FB-25)
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_replay_buffer_status",
				Description: "Check if the replay buffer is currently active",
			},
			s.handleGetReplayBufferStatus,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "toggle_replay_buffer",
				Description: "Start or stop the replay buffer",
			},
			s.handleToggleReplayBuffer,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "save_replay_buffer",
				Description: "Save the current replay buffer to disk (replay buffer must be active)",
			},
			s.handleSaveReplayBuffer,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_last_replay",
				Description: "Get the file path of the last saved replay buffer",
			},
			s.handleGetLastReplay,
		)

		// Studio mode tools (FB-26)
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_studio_mode_enabled",
				Description: "Check if studio mode is currently enabled in OBS",
			},
			s.handleGetStudioModeEnabled,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "toggle_studio_mode",
				Description: "Enable or disable studio mode in OBS",
			},
			s.handleToggleStudioMode,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_preview_scene",
				Description: "Get the current preview scene in studio mode",
			},
			s.handleGetPreviewScene,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "set_preview_scene",
				Description: "Set the preview scene in studio mode",
			},
			s.handleSetPreviewScene,
		)

		// Hotkey tools (FB-26)
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "trigger_hotkey_by_name",
				Description: "Trigger an OBS hotkey by its name (use list_hotkeys to see available hotkeys)",
			},
			s.handleTriggerHotkeyByName,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "list_hotkeys",
				Description: "List all available OBS hotkey names",
			},
			s.handleListHotkeys,
		)

		toolCount += 25
		log.Println("Core tools registered (25 tools)")
	}

	// Source tools
	if s.toolGroups.Sources {
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "list_sources",
				Description: "List all input sources (audio and video) available in OBS",
			},
			s.handleListSources,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "toggle_source_visibility",
				Description: "Toggle the visibility of a source in a specific scene",
			},
			s.handleToggleSourceVisibility,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_source_settings",
				Description: "Retrieve configuration settings for a specific source",
			},
			s.handleGetSourceSettings,
		)

		toolCount += 3
		log.Println("Source tools registered (3 tools)")
	}

	// Audio tools
	if s.toolGroups.Audio {
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_input_mute",
				Description: "Check whether an audio input is currently muted",
			},
			s.handleGetInputMute,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "toggle_input_mute",
				Description: "Toggle the mute state of an audio input (muted <-> unmuted)",
			},
			s.handleToggleInputMute,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "set_input_volume",
				Description: "Set the volume level of an audio input (supports dB or multiplier format)",
			},
			s.handleSetInputVolume,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_input_volume",
				Description: "Get the current volume level of an audio input (returns dB and multiplier values)",
			},
			s.handleGetInputVolume,
		)

		toolCount += 4
		log.Println("Audio tools registered (4 tools)")
	}

	// Layout tools: Scene presets
	if s.toolGroups.Layout {
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "list_scene_presets",
				Description: "List all saved scene presets, optionally filtered by scene name",
			},
			s.handleListScenePresets,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_preset_details",
				Description: "Get detailed information about a specific scene preset including source states",
			},
			s.handleGetPresetDetails,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "delete_scene_preset",
				Description: "Delete a saved scene preset by name",
			},
			s.handleDeleteScenePreset,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "rename_scene_preset",
				Description: "Rename an existing scene preset",
			},
			s.handleRenameScenePreset,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "save_scene_preset",
				Description: "Save the current state of a scene as a named preset",
			},
			s.handleSaveScenePreset,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "apply_scene_preset",
				Description: "Apply a saved preset to restore source visibility states",
			},
			s.handleApplyScenePreset,
		)

		toolCount += 6
		log.Println("Layout tools registered (6 tools)")
	}

	// Visual tools: Screenshot sources
	if s.toolGroups.Visual {
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "create_screenshot_source",
				Description: "Create a periodic screenshot capture source for visual monitoring of OBS scenes",
			},
			s.handleCreateScreenshotSource,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "remove_screenshot_source",
				Description: "Stop and remove a screenshot capture source",
			},
			s.handleRemoveScreenshotSource,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "list_screenshot_sources",
				Description: "List all configured screenshot sources with their status and HTTP URLs",
			},
			s.handleListScreenshotSources,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "configure_screenshot_cadence",
				Description: "Update the capture interval for a screenshot source",
			},
			s.handleConfigureScreenshotCadence,
		)

		toolCount += 4
		log.Println("Visual tools registered (4 tools)")
	}

	// Design tools: Source creation and manipulation
	if s.toolGroups.Design {
		// Source creation tools
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "create_text_source",
				Description: "Create a text/label source in a scene with customizable font and color",
			},
			s.handleCreateTextSource,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "create_image_source",
				Description: "Create an image source in a scene from a file path",
			},
			s.handleCreateImageSource,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "create_color_source",
				Description: "Create a solid color source in a scene",
			},
			s.handleCreateColorSource,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "create_browser_source",
				Description: "Create a browser source in a scene to display web content",
			},
			s.handleCreateBrowserSource,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "create_media_source",
				Description: "Create a media/video source in a scene from a file path",
			},
			s.handleCreateMediaSource,
		)

		// Layout control tools
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "set_source_transform",
				Description: "Set position, scale, and rotation of a source in a scene",
			},
			s.handleSetSourceTransform,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_source_transform",
				Description: "Get the current transform properties of a source in a scene",
			},
			s.handleGetSourceTransform,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "set_source_crop",
				Description: "Set crop values for a source in a scene",
			},
			s.handleSetSourceCrop,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "set_source_bounds",
				Description: "Set bounds type and size for a source in a scene",
			},
			s.handleSetSourceBounds,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "set_source_order",
				Description: "Set the z-order index of a source in a scene (0 = back, higher = front)",
			},
			s.handleSetSourceOrder,
		)

		// Advanced tools
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "set_source_locked",
				Description: "Lock or unlock a source to prevent accidental changes",
			},
			s.handleSetSourceLocked,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "duplicate_source",
				Description: "Duplicate a source within the same scene or to another scene",
			},
			s.handleDuplicateSource,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "remove_source",
				Description: "Remove a source from a scene",
			},
			s.handleRemoveSource,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "list_input_kinds",
				Description: "List all available input source types in OBS",
			},
			s.handleListInputKinds,
		)

		toolCount += 14
		log.Println("Design tools registered (14 tools)")
	}

	// Filter tools (FB-23)
	if s.toolGroups.Filters {
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "list_source_filters",
				Description: "List all filters applied to a source",
			},
			s.handleListSourceFilters,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_source_filter",
				Description: "Get detailed information about a specific filter on a source",
			},
			s.handleGetSourceFilter,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "create_source_filter",
				Description: "Add a new filter to a source (e.g., color correction, noise suppression)",
			},
			s.handleCreateSourceFilter,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "remove_source_filter",
				Description: "Remove a filter from a source",
			},
			s.handleRemoveSourceFilter,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "toggle_source_filter",
				Description: "Enable or disable a filter on a source",
			},
			s.handleToggleSourceFilter,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "set_source_filter_settings",
				Description: "Modify the configuration settings of a filter",
			},
			s.handleSetSourceFilterSettings,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "list_filter_kinds",
				Description: "List all available filter types in OBS",
			},
			s.handleListFilterKinds,
		)

		toolCount += 7
		log.Println("Filter tools registered (7 tools)")
	}

	// Transition tools (FB-24)
	if s.toolGroups.Transitions {
		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "list_transitions",
				Description: "List all available scene transitions and identify the current one",
			},
			s.handleListTransitions,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "get_current_transition",
				Description: "Get details about the current scene transition including duration and settings",
			},
			s.handleGetCurrentTransition,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "set_current_transition",
				Description: "Change the active scene transition (e.g., Cut, Fade, Swipe)",
			},
			s.handleSetCurrentTransition,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "set_transition_duration",
				Description: "Set the duration of the current scene transition in milliseconds",
			},
			s.handleSetTransitionDuration,
		)

		mcpsdk.AddTool(s.mcpServer,
			&mcpsdk.Tool{
				Name:        "trigger_transition",
				Description: "Trigger the current transition in studio mode (swaps preview and program scenes)",
			},
			s.handleTriggerTransition,
		)

		toolCount += 5
		log.Println("Transition tools registered (5 tools)")
	}

	// Meta tools - always enabled, cannot be disabled
	// These provide help and runtime tool configuration

	mcpsdk.AddTool(s.mcpServer,
		&mcpsdk.Tool{
			Name:        "help",
			Description: "Get detailed help on agentic-obs features, tools, resources, prompts, and workflows. Use topic='overview' for a quick start guide, topic='tools' for all available tools, or topic='<tool_name>' for specific tool help.",
		},
		s.handleHelp,
	)

	mcpsdk.AddTool(s.mcpServer,
		&mcpsdk.Tool{
			Name:        "get_tool_config",
			Description: "Get current tool group configuration showing which tool groups are enabled/disabled. Use group parameter to filter by specific group, verbose=true to include tool names.",
		},
		s.handleGetToolConfig,
	)

	mcpsdk.AddTool(s.mcpServer,
		&mcpsdk.Tool{
			Name:        "set_tool_config",
			Description: "Enable or disable a tool group at runtime. Changes are session-only by default; use persist=true to save to database for future sessions.",
		},
		s.handleSetToolConfig,
	)

	mcpsdk.AddTool(s.mcpServer,
		&mcpsdk.Tool{
			Name:        "list_tool_groups",
			Description: "List all available tool groups with their descriptions and enabled status. Use include_disabled=false to only show enabled groups.",
		},
		s.handleListToolGroups,
	)

	toolCount += 4 // 4 meta-tools (help + 3 config tools)
	log.Println("Meta tools registered (help, get_tool_config, set_tool_config, list_tool_groups)")

	log.Printf("Tool handlers registered successfully (%d tools total)", toolCount)
}

// Tool handler implementations

func (s *Server) handleSetCurrentScene(ctx context.Context, request *mcpsdk.CallToolRequest, input SceneNameInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Setting current scene to: %s", input.SceneName)

	if err := s.obsClient.SetCurrentScene(input.SceneName); err != nil {
		s.recordAction("set_current_scene", "Set current scene", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to set current scene: %w", err)
	}

	result := SimpleResult{Message: fmt.Sprintf("Successfully switched to scene: %s", input.SceneName)}
	s.recordAction("set_current_scene", "Set current scene", input, result, true, time.Since(start))
	return nil, result, nil
}

func (s *Server) handleCreateScene(ctx context.Context, request *mcpsdk.CallToolRequest, input SceneNameInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Creating scene: %s", input.SceneName)

	if err := s.obsClient.CreateScene(input.SceneName); err != nil {
		s.recordAction("create_scene", "Create scene", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to create scene: %w", err)
	}

	result := SimpleResult{Message: fmt.Sprintf("Successfully created scene: %s", input.SceneName)}
	s.recordAction("create_scene", "Create scene", input, result, true, time.Since(start))
	return nil, result, nil
}

func (s *Server) handleRemoveScene(ctx context.Context, request *mcpsdk.CallToolRequest, input SceneNameInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Removing scene: %s - requesting confirmation", input.SceneName)

	// Request user confirmation before deleting scene
	confirmed, err := ElicitDeleteConfirmation(ctx, getSession(request), "scene", input.SceneName)
	if err != nil {
		log.Printf("Elicitation error: %v", err)
		// Continue without confirmation if elicitation fails
	} else if !confirmed {
		result := CancelledResult("Scene removal")
		s.recordAction("remove_scene", "Remove scene (cancelled)", input, result, false, time.Since(start))
		return nil, result, nil
	}

	if err := s.obsClient.RemoveScene(input.SceneName); err != nil {
		s.recordAction("remove_scene", "Remove scene", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to remove scene: %w", err)
	}

	result := SimpleResult{Message: fmt.Sprintf("Successfully removed scene: %s", input.SceneName)}
	s.recordAction("remove_scene", "Remove scene", input, result, true, time.Since(start))
	return nil, result, nil
}

func (s *Server) handleStartRecording(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Starting recording")

	if err := s.obsClient.StartRecording(); err != nil {
		s.recordAction("start_recording", "Start recording", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to start recording: %w", err)
	}

	result := SimpleResult{Message: "Successfully started recording"}
	s.recordAction("start_recording", "Start recording", nil, result, true, time.Since(start))
	return nil, result, nil
}

func (s *Server) handleStopRecording(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Stopping recording")

	outputPath, err := s.obsClient.StopRecording()
	if err != nil {
		s.recordAction("stop_recording", "Stop recording", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to stop recording: %w", err)
	}

	result := SimpleResult{Message: fmt.Sprintf("Successfully stopped recording. Output saved to: %s", outputPath)}
	s.recordAction("stop_recording", "Stop recording", nil, result, true, time.Since(start))
	return nil, result, nil
}

func (s *Server) handleGetRecordingStatus(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Getting recording status")

	status, err := s.obsClient.GetRecordingStatus()
	if err != nil {
		s.recordAction("get_recording_status", "Get recording status", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get recording status: %w", err)
	}

	s.recordAction("get_recording_status", "Get recording status", nil, status, true, time.Since(start))
	return nil, status, nil
}

func (s *Server) handleStartStreaming(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Starting streaming - requesting confirmation")

	// Request user confirmation before starting stream
	confirmed, err := ElicitStreamingConfirmation(ctx, getSession(request))
	if err != nil {
		log.Printf("Elicitation error: %v", err)
		// Continue without confirmation if elicitation fails
	} else if !confirmed {
		result := CancelledResult("Streaming start")
		s.recordAction("start_streaming", "Start streaming (cancelled)", nil, result, false, time.Since(start))
		return nil, result, nil
	}

	if err := s.obsClient.StartStreaming(); err != nil {
		s.recordAction("start_streaming", "Start streaming", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to start streaming: %w", err)
	}

	result := SimpleResult{Message: "Successfully started streaming"}
	s.recordAction("start_streaming", "Start streaming", nil, result, true, time.Since(start))
	return nil, result, nil
}

func (s *Server) handleStopStreaming(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Stopping streaming - requesting confirmation")

	// Request user confirmation before stopping stream
	confirmed, err := ElicitStopStreamingConfirmation(ctx, getSession(request))
	if err != nil {
		log.Printf("Elicitation error: %v", err)
		// Continue without confirmation if elicitation fails
	} else if !confirmed {
		result := CancelledResult("Streaming stop")
		s.recordAction("stop_streaming", "Stop streaming (cancelled)", nil, result, false, time.Since(start))
		return nil, result, nil
	}

	if err := s.obsClient.StopStreaming(); err != nil {
		s.recordAction("stop_streaming", "Stop streaming", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to stop streaming: %w", err)
	}

	result := SimpleResult{Message: "Successfully stopped streaming"}
	s.recordAction("stop_streaming", "Stop streaming", nil, result, true, time.Since(start))
	return nil, result, nil
}

func (s *Server) handleGetStreamingStatus(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Getting streaming status")

	status, err := s.obsClient.GetStreamingStatus()
	if err != nil {
		s.recordAction("get_streaming_status", "Get streaming status", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get streaming status: %w", err)
	}

	s.recordAction("get_streaming_status", "Get streaming status", nil, status, true, time.Since(start))
	return nil, status, nil
}

func (s *Server) handleGetOBSStatus(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Getting OBS status")

	status, err := s.obsClient.GetOBSStatus()
	if err != nil {
		s.recordAction("get_obs_status", "Get OBS status", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get OBS status: %w", err)
	}

	s.recordAction("get_obs_status", "Get OBS status", nil, status, true, time.Since(start))
	return nil, status, nil
}

// New P1 tool handlers

func (s *Server) handleListScenes(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Listing all scenes")

	scenes, currentScene, err := s.obsClient.GetSceneList()
	if err != nil {
		s.recordAction("list_scenes", "List scenes", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to list scenes: %w", err)
	}

	result := map[string]interface{}{
		"scenes":        scenes,
		"current_scene": currentScene,
	}
	s.recordAction("list_scenes", "List scenes", nil, result, true, time.Since(start))
	return nil, result, nil
}

func (s *Server) handlePauseRecording(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Pausing recording")

	if err := s.obsClient.PauseRecording(); err != nil {
		s.recordAction("pause_recording", "Pause recording", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to pause recording: %w", err)
	}

	result := SimpleResult{Message: "Successfully paused recording"}
	s.recordAction("pause_recording", "Pause recording", nil, result, true, time.Since(start))
	return nil, result, nil
}

func (s *Server) handleResumeRecording(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Resuming recording")

	if err := s.obsClient.ResumeRecording(); err != nil {
		s.recordAction("resume_recording", "Resume recording", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to resume recording: %w", err)
	}

	result := SimpleResult{Message: "Successfully resumed recording"}
	s.recordAction("resume_recording", "Resume recording", nil, result, true, time.Since(start))
	return nil, result, nil
}

func (s *Server) handleListSources(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Listing all sources")

	sources, err := s.obsClient.ListSources()
	if err != nil {
		s.recordAction("list_sources", "List sources", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to list sources: %w", err)
	}

	s.recordAction("list_sources", "List sources", nil, sources, true, time.Since(start))
	return nil, sources, nil
}

func (s *Server) handleToggleSourceVisibility(ctx context.Context, request *mcpsdk.CallToolRequest, input SourceVisibilityInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Toggling visibility for source %d in scene: %s", input.SourceID, input.SceneName)

	newState, err := s.obsClient.ToggleSourceVisibility(input.SceneName, int(input.SourceID))
	if err != nil {
		s.recordAction("toggle_source_visibility", "Toggle source visibility", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to toggle source visibility: %w", err)
	}

	result := map[string]interface{}{
		"scene_name": input.SceneName,
		"source_id":  input.SourceID,
		"visible":    newState,
	}
	s.recordAction("toggle_source_visibility", "Toggle source visibility", input, result, true, time.Since(start))
	return nil, result, nil
}

func (s *Server) handleGetSourceSettings(ctx context.Context, request *mcpsdk.CallToolRequest, input SourceNameInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Getting settings for source: %s", input.SourceName)

	settings, err := s.obsClient.GetSourceSettings(input.SourceName)
	if err != nil {
		s.recordAction("get_source_settings", "Get source settings", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get source settings: %w", err)
	}

	s.recordAction("get_source_settings", "Get source settings", input, settings, true, time.Since(start))
	return nil, settings, nil
}

func (s *Server) handleGetInputMute(ctx context.Context, request *mcpsdk.CallToolRequest, input InputNameInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Getting mute status for input: %s", input.InputName)

	isMuted, err := s.obsClient.GetInputMute(input.InputName)
	if err != nil {
		s.recordAction("get_input_mute", "Get input mute", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get input mute status: %w", err)
	}

	result := map[string]interface{}{
		"input_name": input.InputName,
		"is_muted":   isMuted,
	}
	s.recordAction("get_input_mute", "Get input mute", input, result, true, time.Since(start))
	return nil, result, nil
}

func (s *Server) handleToggleInputMute(ctx context.Context, request *mcpsdk.CallToolRequest, input InputNameInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Toggling mute for input: %s", input.InputName)

	if err := s.obsClient.ToggleInputMute(input.InputName); err != nil {
		s.recordAction("toggle_input_mute", "Toggle input mute", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to toggle input mute: %w", err)
	}

	result := SimpleResult{Message: fmt.Sprintf("Successfully toggled mute for input: %s", input.InputName)}
	s.recordAction("toggle_input_mute", "Toggle input mute", input, result, true, time.Since(start))
	return nil, result, nil
}

func (s *Server) handleSetInputVolume(ctx context.Context, request *mcpsdk.CallToolRequest, input SetVolumeInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Setting volume for input: %s", input.InputName)

	if err := s.obsClient.SetInputVolume(input.InputName, input.VolumeDb, input.VolumeMul); err != nil {
		s.recordAction("set_input_volume", "Set input volume", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to set input volume: %w", err)
	}

	result := SimpleResult{Message: fmt.Sprintf("Successfully set volume for input: %s", input.InputName)}
	s.recordAction("set_input_volume", "Set input volume", input, result, true, time.Since(start))
	return nil, result, nil
}

// Scene preset tool handlers

// handleListScenePresets returns all saved scene presets, optionally filtered by scene name.
// Returns a list of preset summaries (id, name, scene_name, created_at) and total count.
func (s *Server) handleListScenePresets(ctx context.Context, request *mcpsdk.CallToolRequest, input ListPresetsInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Listing scene presets (filter: %s)", input.SceneName)

	presets, err := s.storage.ListScenePresets(ctx, input.SceneName)
	if err != nil {
		s.recordAction("list_scene_presets", "List scene presets", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to list scene presets: %w", err)
	}

	// Convert to simpler response format (without full source details)
	presetList := make([]map[string]interface{}, len(presets))
	for i, p := range presets {
		presetList[i] = map[string]interface{}{
			"id":         p.ID,
			"name":       p.Name,
			"scene_name": p.SceneName,
			"created_at": p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	result := map[string]interface{}{
		"presets": presetList,
		"count":   len(presets),
	}
	s.recordAction("list_scene_presets", "List scene presets", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleGetPresetDetails retrieves full details of a scene preset including source states.
// Returns the preset's id, name, scene_name, sources array, and created_at timestamp.
// Returns an error if the preset does not exist.
func (s *Server) handleGetPresetDetails(ctx context.Context, request *mcpsdk.CallToolRequest, input PresetNameInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Getting preset details for: %s", input.PresetName)

	preset, err := s.storage.GetScenePreset(ctx, input.PresetName)
	if err != nil {
		s.recordAction("get_preset_details", "Get preset details", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get preset details: %w", err)
	}

	result := map[string]interface{}{
		"id":         preset.ID,
		"name":       preset.Name,
		"scene_name": preset.SceneName,
		"sources":    preset.Sources,
		"created_at": preset.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	s.recordAction("get_preset_details", "Get preset details", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleDeleteScenePreset permanently removes a scene preset from storage.
// Returns a success message or an error if the preset does not exist.
func (s *Server) handleDeleteScenePreset(ctx context.Context, request *mcpsdk.CallToolRequest, input PresetNameInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Deleting scene preset: %s - requesting confirmation", input.PresetName)

	// Request user confirmation before deleting preset
	confirmed, err := ElicitDeleteConfirmation(ctx, getSession(request), "preset", input.PresetName)
	if err != nil {
		log.Printf("Elicitation error: %v", err)
		// Continue without confirmation if elicitation fails
	} else if !confirmed {
		result := CancelledResult("Preset deletion")
		s.recordAction("delete_scene_preset", "Delete scene preset (cancelled)", input, result, false, time.Since(start))
		return nil, result, nil
	}

	if err := s.storage.DeleteScenePreset(ctx, input.PresetName); err != nil {
		s.recordAction("delete_scene_preset", "Delete scene preset", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to delete preset: %w", err)
	}

	result := SimpleResult{Message: fmt.Sprintf("Successfully deleted preset: %s", input.PresetName)}
	s.recordAction("delete_scene_preset", "Delete scene preset", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleRenameScenePreset changes the name of an existing scene preset.
// Returns a success message or an error if the preset does not exist or the new name conflicts.
func (s *Server) handleRenameScenePreset(ctx context.Context, request *mcpsdk.CallToolRequest, input RenamePresetInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Renaming preset from '%s' to '%s'", input.OldName, input.NewName)

	if err := s.storage.RenameScenePreset(ctx, input.OldName, input.NewName); err != nil {
		s.recordAction("rename_scene_preset", "Rename scene preset", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to rename preset: %w", err)
	}

	result := SimpleResult{Message: fmt.Sprintf("Successfully renamed preset from '%s' to '%s'", input.OldName, input.NewName)}
	s.recordAction("rename_scene_preset", "Rename scene preset", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleGetInputVolume retrieves the current volume level of an audio input.
// Returns volume_db (decibels) and volume_mul (linear multiplier) values.
// Returns an error if the input does not exist or OBS is not connected.
func (s *Server) handleGetInputVolume(ctx context.Context, request *mcpsdk.CallToolRequest, input InputNameInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Getting volume for input: %s", input.InputName)

	volumeDb, volumeMul, err := s.obsClient.GetInputVolume(input.InputName)
	if err != nil {
		s.recordAction("get_input_volume", "Get input volume", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get input volume: %w", err)
	}

	result := map[string]interface{}{
		"input_name": input.InputName,
		"volume_db":  volumeDb,
		"volume_mul": volumeMul,
	}
	s.recordAction("get_input_volume", "Get input volume", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleSaveScenePreset captures the current source visibility states from an OBS scene
// and saves them as a named preset in storage. Returns the preset id, name, scene_name,
// source_count, and a success message. Returns an error if the scene does not exist,
// OBS is not connected, or a preset with the same name already exists.
func (s *Server) handleSaveScenePreset(ctx context.Context, request *mcpsdk.CallToolRequest, input SavePresetInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Saving scene preset '%s' for scene '%s'", input.PresetName, input.SceneName)

	// Capture current scene state from OBS
	states, err := s.obsClient.CaptureSceneState(input.SceneName)
	if err != nil {
		s.recordAction("save_scene_preset", "Save scene preset", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to capture scene state: %w", err)
	}

	// Convert OBS source states to storage format
	sources := make([]storage.SourceState, len(states))
	for i, state := range states {
		sources[i] = storage.SourceState{
			Name:    state.Name,
			Visible: state.Enabled,
		}
	}

	// Create preset in storage
	preset := storage.ScenePreset{
		Name:      input.PresetName,
		SceneName: input.SceneName,
		Sources:   sources,
	}

	id, err := s.storage.CreateScenePreset(ctx, preset)
	if err != nil {
		s.recordAction("save_scene_preset", "Save scene preset", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to save preset: %w", err)
	}

	result := map[string]interface{}{
		"id":           id,
		"preset_name":  input.PresetName,
		"scene_name":   input.SceneName,
		"source_count": len(sources),
		"message":      fmt.Sprintf("Successfully saved preset '%s' with %d sources", input.PresetName, len(sources)),
	}
	s.recordAction("save_scene_preset", "Save scene preset", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleApplyScenePreset loads a saved preset and applies its source visibility states
// to the target OBS scene. Sources that no longer exist in the scene are skipped.
// Returns the preset_name, scene_name, applied_count, and a success message.
// Returns an error if the preset does not exist, the scene no longer exists, or OBS is not connected.
func (s *Server) handleApplyScenePreset(ctx context.Context, request *mcpsdk.CallToolRequest, input PresetNameInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Applying scene preset: %s", input.PresetName)

	// Load preset from storage
	preset, err := s.storage.GetScenePreset(ctx, input.PresetName)
	if err != nil {
		s.recordAction("apply_scene_preset", "Apply scene preset", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to load preset: %w", err)
	}

	// Get current scene items to map names to IDs
	scene, err := s.obsClient.GetSceneByName(preset.SceneName)
	if err != nil {
		s.recordAction("apply_scene_preset", "Apply scene preset", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get scene '%s': %w", preset.SceneName, err)
	}

	// Build name-to-ID map
	nameToID := make(map[string]int)
	for _, src := range scene.Sources {
		nameToID[src.Name] = src.ID
	}

	// Convert storage format to OBS source states
	obsStates := make([]obs.SourceState, 0, len(preset.Sources))
	for _, src := range preset.Sources {
		id, exists := nameToID[src.Name]
		if !exists {
			log.Printf("Warning: source '%s' not found in scene, skipping", src.Name)
			continue
		}
		obsStates = append(obsStates, obs.SourceState{
			ID:      id,
			Name:    src.Name,
			Enabled: src.Visible,
		})
	}

	// Apply preset to OBS
	if err := s.obsClient.ApplyScenePreset(preset.SceneName, obsStates); err != nil {
		s.recordAction("apply_scene_preset", "Apply scene preset", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to apply preset: %w", err)
	}

	result := map[string]interface{}{
		"preset_name":   input.PresetName,
		"scene_name":    preset.SceneName,
		"applied_count": len(obsStates),
		"message":       fmt.Sprintf("Successfully applied preset '%s' to scene '%s'", input.PresetName, preset.SceneName),
	}
	s.recordAction("apply_scene_preset", "Apply scene preset", input, result, true, time.Since(start))
	return nil, result, nil
}

// Screenshot source tool handlers

// handleCreateScreenshotSource creates a new periodic screenshot capture source.
// Returns the source details including the HTTP URL where screenshots can be accessed.
func (s *Server) handleCreateScreenshotSource(ctx context.Context, request *mcpsdk.CallToolRequest, input CreateScreenshotSourceInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Creating screenshot source '%s' for OBS source '%s'", input.Name, input.SourceName)

	// Set defaults
	if input.CadenceMs <= 0 {
		input.CadenceMs = 5000
	}
	if input.ImageFormat == "" {
		input.ImageFormat = "png"
	}
	if input.Quality <= 0 {
		input.Quality = 80
	}

	// Create the source in storage
	source := storage.ScreenshotSource{
		Name:        input.Name,
		SourceName:  input.SourceName,
		CadenceMs:   input.CadenceMs,
		ImageFormat: input.ImageFormat,
		ImageWidth:  input.ImageWidth,
		ImageHeight: input.ImageHeight,
		Quality:     input.Quality,
		Enabled:     true,
	}

	id, err := s.storage.CreateScreenshotSource(ctx, source)
	if err != nil {
		s.recordAction("create_screenshot_source", "Create screenshot source", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to create screenshot source: %w", err)
	}

	// Retrieve the full source with defaults applied
	createdSource, err := s.storage.GetScreenshotSource(ctx, id)
	if err != nil {
		s.recordAction("create_screenshot_source", "Create screenshot source", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to retrieve created source: %w", err)
	}

	// Start the capture worker
	if err := s.screenshotMgr.AddSource(createdSource); err != nil {
		// Log but don't fail - source is created, worker can be started later
		log.Printf("Warning: failed to start capture worker: %v", err)
	}

	// Get the HTTP URL for this source
	screenshotURL := s.httpServer.GetScreenshotURL(input.Name)

	result := map[string]interface{}{
		"id":           id,
		"name":         input.Name,
		"source_name":  input.SourceName,
		"cadence_ms":   input.CadenceMs,
		"image_format": input.ImageFormat,
		"quality":      input.Quality,
		"url":          screenshotURL,
		"message":      fmt.Sprintf("Successfully created screenshot source '%s'. Access at: %s", input.Name, screenshotURL),
	}
	s.recordAction("create_screenshot_source", "Create screenshot source", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleRemoveScreenshotSource stops and removes a screenshot capture source.
func (s *Server) handleRemoveScreenshotSource(ctx context.Context, request *mcpsdk.CallToolRequest, input ScreenshotSourceNameInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Removing screenshot source: %s", input.Name)

	// Get the source to find its ID
	source, err := s.storage.GetScreenshotSourceByName(ctx, input.Name)
	if err != nil {
		s.recordAction("remove_screenshot_source", "Remove screenshot source", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to find screenshot source: %w", err)
	}

	// Stop the capture worker
	if err := s.screenshotMgr.RemoveSource(source.ID); err != nil {
		log.Printf("Warning: failed to stop capture worker: %v", err)
	}

	// Delete from storage (cascades to delete screenshots)
	if err := s.storage.DeleteScreenshotSource(ctx, source.ID); err != nil {
		s.recordAction("remove_screenshot_source", "Remove screenshot source", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to delete screenshot source: %w", err)
	}

	result := SimpleResult{Message: fmt.Sprintf("Successfully removed screenshot source '%s'", input.Name)}
	s.recordAction("remove_screenshot_source", "Remove screenshot source", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleListScreenshotSources lists all configured screenshot sources with their status and URLs.
func (s *Server) handleListScreenshotSources(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Listing screenshot sources")

	sources, err := s.storage.ListScreenshotSources(ctx)
	if err != nil {
		s.recordAction("list_screenshot_sources", "List screenshot sources", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to list screenshot sources: %w", err)
	}

	sourceList := make([]map[string]interface{}, len(sources))
	for i, src := range sources {
		// Get screenshot count for this source
		count, _ := s.storage.CountScreenshots(ctx, src.ID)

		sourceList[i] = map[string]interface{}{
			"id":               src.ID,
			"name":             src.Name,
			"source_name":      src.SourceName,
			"cadence_ms":       src.CadenceMs,
			"image_format":     src.ImageFormat,
			"quality":          src.Quality,
			"enabled":          src.Enabled,
			"url":              s.httpServer.GetScreenshotURL(src.Name),
			"screenshot_count": count,
			"created_at":       src.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	result := map[string]interface{}{
		"sources": sourceList,
		"count":   len(sources),
	}
	s.recordAction("list_screenshot_sources", "List screenshot sources", nil, result, true, time.Since(start))
	return nil, result, nil
}

// handleConfigureScreenshotCadence updates the capture interval for a screenshot source.
func (s *Server) handleConfigureScreenshotCadence(ctx context.Context, request *mcpsdk.CallToolRequest, input ConfigureScreenshotCadenceInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Updating cadence for screenshot source '%s' to %dms", input.Name, input.CadenceMs)

	if input.CadenceMs <= 0 {
		s.recordAction("configure_screenshot_cadence", "Configure screenshot cadence", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("cadence_ms must be greater than 0")
	}

	// Get the source
	source, err := s.storage.GetScreenshotSourceByName(ctx, input.Name)
	if err != nil {
		s.recordAction("configure_screenshot_cadence", "Configure screenshot cadence", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to find screenshot source: %w", err)
	}

	// Update in storage
	source.CadenceMs = input.CadenceMs
	if err := s.storage.UpdateScreenshotSource(ctx, *source); err != nil {
		s.recordAction("configure_screenshot_cadence", "Configure screenshot cadence", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to update screenshot source: %w", err)
	}

	// Update the running worker's cadence
	if err := s.screenshotMgr.UpdateCadence(source.ID, input.CadenceMs); err != nil {
		// Try to restart the source with updated settings
		if err := s.screenshotMgr.UpdateSource(source); err != nil {
			log.Printf("Warning: failed to update capture worker cadence: %v", err)
		}
	}

	result := map[string]interface{}{
		"name":       input.Name,
		"cadence_ms": input.CadenceMs,
		"message":    fmt.Sprintf("Successfully updated cadence for '%s' to %dms", input.Name, input.CadenceMs),
	}
	s.recordAction("configure_screenshot_cadence", "Configure screenshot cadence", input, result, true, time.Since(start))
	return nil, result, nil
}

// Design tool handlers

// handleCreateTextSource creates a text source in a scene
func (s *Server) handleCreateTextSource(ctx context.Context, request *mcpsdk.CallToolRequest, input CreateTextSourceInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Creating text source '%s' in scene '%s'", input.SourceName, input.SceneName)

	// Build settings map
	settings := map[string]interface{}{
		"text": input.Text,
	}

	// Apply optional font settings
	if input.FontName != "" || input.FontSize > 0 {
		font := map[string]interface{}{}
		if input.FontName != "" {
			font["face"] = input.FontName
		}
		if input.FontSize > 0 {
			font["size"] = input.FontSize
		}
		settings["font"] = font
	}

	if input.Color != 0 {
		settings["color"] = input.Color
	}

	// Create the input using the generic method
	sceneItemID, err := s.obsClient.CreateInput(input.SceneName, input.SourceName, "text_gdiplus_v3", settings)
	if err != nil {
		s.recordAction("create_text_source", "Create text source", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to create text source: %w", err)
	}

	result := map[string]interface{}{
		"scene_name":    input.SceneName,
		"source_name":   input.SourceName,
		"scene_item_id": sceneItemID,
		"message":       fmt.Sprintf("Successfully created text source '%s' in scene '%s'", input.SourceName, input.SceneName),
	}
	s.recordAction("create_text_source", "Create text source", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleCreateImageSource creates an image source in a scene
func (s *Server) handleCreateImageSource(ctx context.Context, request *mcpsdk.CallToolRequest, input CreateImageSourceInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Creating image source '%s' in scene '%s'", input.SourceName, input.SceneName)

	settings := map[string]interface{}{
		"file": input.FilePath,
	}

	sceneItemID, err := s.obsClient.CreateInput(input.SceneName, input.SourceName, "image_source", settings)
	if err != nil {
		s.recordAction("create_image_source", "Create image source", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to create image source: %w", err)
	}

	result := map[string]interface{}{
		"scene_name":    input.SceneName,
		"source_name":   input.SourceName,
		"scene_item_id": sceneItemID,
		"file_path":     input.FilePath,
		"message":       fmt.Sprintf("Successfully created image source '%s' in scene '%s'", input.SourceName, input.SceneName),
	}
	s.recordAction("create_image_source", "Create image source", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleCreateColorSource creates a solid color source in a scene
func (s *Server) handleCreateColorSource(ctx context.Context, request *mcpsdk.CallToolRequest, input CreateColorSourceInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Creating color source '%s' in scene '%s'", input.SourceName, input.SceneName)

	// Set defaults
	width := input.Width
	if width <= 0 {
		width = 1920
	}
	height := input.Height
	if height <= 0 {
		height = 1080
	}

	settings := map[string]interface{}{
		"color":  input.Color,
		"width":  width,
		"height": height,
	}

	sceneItemID, err := s.obsClient.CreateInput(input.SceneName, input.SourceName, "color_source_v3", settings)
	if err != nil {
		s.recordAction("create_color_source", "Create color source", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to create color source: %w", err)
	}

	result := map[string]interface{}{
		"scene_name":    input.SceneName,
		"source_name":   input.SourceName,
		"scene_item_id": sceneItemID,
		"width":         width,
		"height":        height,
		"message":       fmt.Sprintf("Successfully created color source '%s' in scene '%s'", input.SourceName, input.SceneName),
	}
	s.recordAction("create_color_source", "Create color source", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleCreateBrowserSource creates a browser source in a scene
func (s *Server) handleCreateBrowserSource(ctx context.Context, request *mcpsdk.CallToolRequest, input CreateBrowserSourceInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Creating browser source '%s' in scene '%s'", input.SourceName, input.SceneName)

	// Set defaults
	width := input.Width
	if width <= 0 {
		width = 800
	}
	height := input.Height
	if height <= 0 {
		height = 600
	}
	fps := input.FPS
	if fps <= 0 {
		fps = 30
	}

	settings := map[string]interface{}{
		"url":    input.URL,
		"width":  width,
		"height": height,
		"fps":    fps,
	}

	sceneItemID, err := s.obsClient.CreateInput(input.SceneName, input.SourceName, "browser_source", settings)
	if err != nil {
		s.recordAction("create_browser_source", "Create browser source", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to create browser source: %w", err)
	}

	result := map[string]interface{}{
		"scene_name":    input.SceneName,
		"source_name":   input.SourceName,
		"scene_item_id": sceneItemID,
		"url":           input.URL,
		"width":         width,
		"height":        height,
		"message":       fmt.Sprintf("Successfully created browser source '%s' in scene '%s'", input.SourceName, input.SceneName),
	}
	s.recordAction("create_browser_source", "Create browser source", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleCreateMediaSource creates a media/video source in a scene
func (s *Server) handleCreateMediaSource(ctx context.Context, request *mcpsdk.CallToolRequest, input CreateMediaSourceInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Creating media source '%s' in scene '%s'", input.SourceName, input.SceneName)

	settings := map[string]interface{}{
		"local_file":   input.FilePath,
		"looping":      input.Loop,
		"hw_decode":    true,
		"clear_on_end": false,
	}

	sceneItemID, err := s.obsClient.CreateInput(input.SceneName, input.SourceName, "ffmpeg_source", settings)
	if err != nil {
		s.recordAction("create_media_source", "Create media source", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to create media source: %w", err)
	}

	result := map[string]interface{}{
		"scene_name":    input.SceneName,
		"source_name":   input.SourceName,
		"scene_item_id": sceneItemID,
		"file_path":     input.FilePath,
		"loop":          input.Loop,
		"message":       fmt.Sprintf("Successfully created media source '%s' in scene '%s'", input.SourceName, input.SceneName),
	}
	s.recordAction("create_media_source", "Create media source", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleSetSourceTransform sets the position, scale, and rotation of a source
func (s *Server) handleSetSourceTransform(ctx context.Context, request *mcpsdk.CallToolRequest, input SetSourceTransformInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Setting transform for scene item %d in scene '%s'", input.SceneItemID, input.SceneName)

	// Get current transform first
	current, err := s.obsClient.GetSceneItemTransform(input.SceneName, input.SceneItemID)
	if err != nil {
		s.recordAction("set_source_transform", "Set source transform", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get current transform: %w", err)
	}

	// Apply changes only for provided values
	if input.X != nil {
		current.PositionX = *input.X
	}
	if input.Y != nil {
		current.PositionY = *input.Y
	}
	if input.ScaleX != nil {
		current.ScaleX = *input.ScaleX
	}
	if input.ScaleY != nil {
		current.ScaleY = *input.ScaleY
	}
	if input.Rotation != nil {
		current.Rotation = *input.Rotation
	}

	if err := s.obsClient.SetSceneItemTransform(input.SceneName, input.SceneItemID, current); err != nil {
		s.recordAction("set_source_transform", "Set source transform", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to set transform: %w", err)
	}

	result := map[string]interface{}{
		"scene_name":    input.SceneName,
		"scene_item_id": input.SceneItemID,
		"x":             current.PositionX,
		"y":             current.PositionY,
		"scale_x":       current.ScaleX,
		"scale_y":       current.ScaleY,
		"rotation":      current.Rotation,
		"message":       "Successfully updated source transform",
	}
	s.recordAction("set_source_transform", "Set source transform", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleGetSourceTransform gets the current transform properties of a source
func (s *Server) handleGetSourceTransform(ctx context.Context, request *mcpsdk.CallToolRequest, input GetSourceTransformInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Getting transform for scene item %d in scene '%s'", input.SceneItemID, input.SceneName)

	transform, err := s.obsClient.GetSceneItemTransform(input.SceneName, input.SceneItemID)
	if err != nil {
		s.recordAction("get_source_transform", "Get source transform", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get transform: %w", err)
	}

	result := map[string]interface{}{
		"scene_name":    input.SceneName,
		"scene_item_id": input.SceneItemID,
		"x":             transform.PositionX,
		"y":             transform.PositionY,
		"scale_x":       transform.ScaleX,
		"scale_y":       transform.ScaleY,
		"rotation":      transform.Rotation,
		"width":         transform.Width,
		"height":        transform.Height,
		"source_width":  transform.SourceWidth,
		"source_height": transform.SourceHeight,
		"bounds_type":   transform.BoundsType,
		"bounds_width":  transform.BoundsWidth,
		"bounds_height": transform.BoundsHeight,
		"crop_top":      transform.CropTop,
		"crop_bottom":   transform.CropBottom,
		"crop_left":     transform.CropLeft,
		"crop_right":    transform.CropRight,
	}
	s.recordAction("get_source_transform", "Get source transform", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleSetSourceCrop sets the crop values for a source
func (s *Server) handleSetSourceCrop(ctx context.Context, request *mcpsdk.CallToolRequest, input SetSourceCropInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Setting crop for scene item %d in scene '%s'", input.SceneItemID, input.SceneName)

	// Get current transform
	current, err := s.obsClient.GetSceneItemTransform(input.SceneName, input.SceneItemID)
	if err != nil {
		s.recordAction("set_source_crop", "Set source crop", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get current transform: %w", err)
	}

	// Apply crop values
	current.CropTop = input.CropTop
	current.CropBottom = input.CropBottom
	current.CropLeft = input.CropLeft
	current.CropRight = input.CropRight

	if err := s.obsClient.SetSceneItemTransform(input.SceneName, input.SceneItemID, current); err != nil {
		s.recordAction("set_source_crop", "Set source crop", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to set crop: %w", err)
	}

	result := map[string]interface{}{
		"scene_name":    input.SceneName,
		"scene_item_id": input.SceneItemID,
		"crop_top":      input.CropTop,
		"crop_bottom":   input.CropBottom,
		"crop_left":     input.CropLeft,
		"crop_right":    input.CropRight,
		"message":       "Successfully updated source crop",
	}
	s.recordAction("set_source_crop", "Set source crop", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleSetSourceBounds sets the bounds type and size for a source
func (s *Server) handleSetSourceBounds(ctx context.Context, request *mcpsdk.CallToolRequest, input SetSourceBoundsInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Setting bounds for scene item %d in scene '%s'", input.SceneItemID, input.SceneName)

	// Get current transform
	current, err := s.obsClient.GetSceneItemTransform(input.SceneName, input.SceneItemID)
	if err != nil {
		s.recordAction("set_source_bounds", "Set source bounds", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get current transform: %w", err)
	}

	// Apply bounds values
	current.BoundsType = input.BoundsType
	current.BoundsWidth = input.BoundsWidth
	current.BoundsHeight = input.BoundsHeight

	if err := s.obsClient.SetSceneItemTransform(input.SceneName, input.SceneItemID, current); err != nil {
		s.recordAction("set_source_bounds", "Set source bounds", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to set bounds: %w", err)
	}

	result := map[string]interface{}{
		"scene_name":    input.SceneName,
		"scene_item_id": input.SceneItemID,
		"bounds_type":   input.BoundsType,
		"bounds_width":  input.BoundsWidth,
		"bounds_height": input.BoundsHeight,
		"message":       "Successfully updated source bounds",
	}
	s.recordAction("set_source_bounds", "Set source bounds", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleSetSourceOrder sets the z-order index of a source
func (s *Server) handleSetSourceOrder(ctx context.Context, request *mcpsdk.CallToolRequest, input SetSourceOrderInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Setting order for scene item %d in scene '%s' to index %d", input.SceneItemID, input.SceneName, input.Index)

	if err := s.obsClient.SetSceneItemIndex(input.SceneName, input.SceneItemID, input.Index); err != nil {
		s.recordAction("set_source_order", "Set source order", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to set order: %w", err)
	}

	result := map[string]interface{}{
		"scene_name":    input.SceneName,
		"scene_item_id": input.SceneItemID,
		"index":         input.Index,
		"message":       fmt.Sprintf("Successfully set source order to index %d", input.Index),
	}
	s.recordAction("set_source_order", "Set source order", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleSetSourceLocked locks or unlocks a source
func (s *Server) handleSetSourceLocked(ctx context.Context, request *mcpsdk.CallToolRequest, input SetSourceLockedInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Setting locked=%v for scene item %d in scene '%s'", input.Locked, input.SceneItemID, input.SceneName)

	if err := s.obsClient.SetSceneItemLocked(input.SceneName, input.SceneItemID, input.Locked); err != nil {
		s.recordAction("set_source_locked", "Set source locked", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to set locked state: %w", err)
	}

	status := "unlocked"
	if input.Locked {
		status = "locked"
	}

	result := map[string]interface{}{
		"scene_name":    input.SceneName,
		"scene_item_id": input.SceneItemID,
		"locked":        input.Locked,
		"message":       fmt.Sprintf("Successfully %s source", status),
	}
	s.recordAction("set_source_locked", "Set source locked", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleDuplicateSource duplicates a source within the same scene or to another scene
func (s *Server) handleDuplicateSource(ctx context.Context, request *mcpsdk.CallToolRequest, input DuplicateSourceInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	destScene := input.DestSceneName
	if destScene == "" {
		destScene = input.SceneName
	}
	log.Printf("Duplicating scene item %d from scene '%s' to '%s'", input.SceneItemID, input.SceneName, destScene)

	newItemID, err := s.obsClient.DuplicateSceneItem(input.SceneName, input.SceneItemID, destScene)
	if err != nil {
		s.recordAction("duplicate_source", "Duplicate source", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to duplicate source: %w", err)
	}

	result := map[string]interface{}{
		"source_scene":      input.SceneName,
		"source_item_id":    input.SceneItemID,
		"dest_scene":        destScene,
		"new_scene_item_id": newItemID,
		"message":           fmt.Sprintf("Successfully duplicated source to scene '%s' with item ID %d", destScene, newItemID),
	}
	s.recordAction("duplicate_source", "Duplicate source", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleRemoveSource removes a source from a scene
func (s *Server) handleRemoveSource(ctx context.Context, request *mcpsdk.CallToolRequest, input RemoveSourceInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Removing scene item %d from scene '%s'", input.SceneItemID, input.SceneName)

	if err := s.obsClient.RemoveSceneItem(input.SceneName, input.SceneItemID); err != nil {
		s.recordAction("remove_source", "Remove source", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to remove source: %w", err)
	}

	result := map[string]interface{}{
		"scene_name":    input.SceneName,
		"scene_item_id": input.SceneItemID,
		"message":       "Successfully removed source from scene",
	}
	s.recordAction("remove_source", "Remove source", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleListInputKinds lists all available input source types in OBS
func (s *Server) handleListInputKinds(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Listing available input kinds")

	kinds, err := s.obsClient.GetInputKindList()
	if err != nil {
		s.recordAction("list_input_kinds", "List input kinds", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to list input kinds: %w", err)
	}

	result := map[string]interface{}{
		"input_kinds": kinds,
		"count":       len(kinds),
	}
	s.recordAction("list_input_kinds", "List input kinds", nil, result, true, time.Since(start))
	return nil, result, nil
}

// =============================================================================
// Filter tool handlers (FB-23)
// =============================================================================

// handleListSourceFilters lists all filters on a source
func (s *Server) handleListSourceFilters(ctx context.Context, request *mcpsdk.CallToolRequest, input ListSourceFiltersInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Listing filters for source: %s", input.SourceName)

	filters, err := s.obsClient.GetSourceFilterList(input.SourceName)
	if err != nil {
		s.recordAction("list_source_filters", "List source filters", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to list filters: %w", err)
	}

	result := map[string]interface{}{
		"source_name": input.SourceName,
		"filters":     filters,
		"count":       len(filters),
	}
	s.recordAction("list_source_filters", "List source filters", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleGetSourceFilter gets details about a specific filter
func (s *Server) handleGetSourceFilter(ctx context.Context, request *mcpsdk.CallToolRequest, input GetSourceFilterInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Getting filter '%s' on source '%s'", input.FilterName, input.SourceName)

	filter, err := s.obsClient.GetSourceFilter(input.SourceName, input.FilterName)
	if err != nil {
		s.recordAction("get_source_filter", "Get source filter", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get filter: %w", err)
	}

	result := map[string]interface{}{
		"source_name": input.SourceName,
		"filter":      filter,
	}
	s.recordAction("get_source_filter", "Get source filter", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleCreateSourceFilter creates a new filter on a source
func (s *Server) handleCreateSourceFilter(ctx context.Context, request *mcpsdk.CallToolRequest, input CreateSourceFilterInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Creating filter '%s' of type '%s' on source '%s'", input.FilterName, input.FilterKind, input.SourceName)

	if err := s.obsClient.CreateSourceFilter(input.SourceName, input.FilterName, input.FilterKind, input.FilterSettings); err != nil {
		s.recordAction("create_source_filter", "Create source filter", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to create filter: %w", err)
	}

	result := map[string]interface{}{
		"source_name": input.SourceName,
		"filter_name": input.FilterName,
		"filter_kind": input.FilterKind,
		"message":     fmt.Sprintf("Successfully created filter '%s' on source '%s'", input.FilterName, input.SourceName),
	}
	s.recordAction("create_source_filter", "Create source filter", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleRemoveSourceFilter removes a filter from a source
func (s *Server) handleRemoveSourceFilter(ctx context.Context, request *mcpsdk.CallToolRequest, input RemoveSourceFilterInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Removing filter '%s' from source '%s' - requesting confirmation", input.FilterName, input.SourceName)

	// Request user confirmation before removing filter
	confirmed, err := ElicitFilterRemovalConfirmation(ctx, getSession(request), input.SourceName, input.FilterName)
	if err != nil {
		log.Printf("Elicitation error: %v", err)
		// Continue without confirmation if elicitation fails
	} else if !confirmed {
		result := CancelledResult("Filter removal")
		s.recordAction("remove_source_filter", "Remove source filter (cancelled)", input, result, false, time.Since(start))
		return nil, result, nil
	}

	if err := s.obsClient.RemoveSourceFilter(input.SourceName, input.FilterName); err != nil {
		s.recordAction("remove_source_filter", "Remove source filter", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to remove filter: %w", err)
	}

	result := map[string]interface{}{
		"source_name": input.SourceName,
		"filter_name": input.FilterName,
		"message":     fmt.Sprintf("Successfully removed filter '%s' from source '%s'", input.FilterName, input.SourceName),
	}
	s.recordAction("remove_source_filter", "Remove source filter", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleToggleSourceFilter enables or disables a filter
func (s *Server) handleToggleSourceFilter(ctx context.Context, request *mcpsdk.CallToolRequest, input ToggleSourceFilterInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Toggling filter '%s' on source '%s'", input.FilterName, input.SourceName)

	var enabled bool
	if input.FilterEnabled != nil {
		// Explicit enable/disable
		enabled = *input.FilterEnabled
	} else {
		// Toggle: get current state and flip it
		filter, err := s.obsClient.GetSourceFilter(input.SourceName, input.FilterName)
		if err != nil {
			s.recordAction("toggle_source_filter", "Toggle source filter", input, nil, false, time.Since(start))
			return nil, nil, fmt.Errorf("failed to get filter state: %w", err)
		}
		enabled = !filter.Enabled
	}

	if err := s.obsClient.SetSourceFilterEnabled(input.SourceName, input.FilterName, enabled); err != nil {
		s.recordAction("toggle_source_filter", "Toggle source filter", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to toggle filter: %w", err)
	}

	status := "disabled"
	if enabled {
		status = "enabled"
	}

	result := map[string]interface{}{
		"source_name":    input.SourceName,
		"filter_name":    input.FilterName,
		"filter_enabled": enabled,
		"message":        fmt.Sprintf("Filter '%s' is now %s", input.FilterName, status),
	}
	s.recordAction("toggle_source_filter", "Toggle source filter", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleSetSourceFilterSettings updates filter settings
func (s *Server) handleSetSourceFilterSettings(ctx context.Context, request *mcpsdk.CallToolRequest, input SetSourceFilterSettingsInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Setting filter settings for '%s' on source '%s'", input.FilterName, input.SourceName)

	// Default to overlay mode (merge settings)
	overlay := true
	if !input.Overlay {
		overlay = false
	}

	if err := s.obsClient.SetSourceFilterSettings(input.SourceName, input.FilterName, input.FilterSettings, overlay); err != nil {
		s.recordAction("set_source_filter_settings", "Set source filter settings", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to set filter settings: %w", err)
	}

	mode := "merged with"
	if !overlay {
		mode = "replaced"
	}

	result := map[string]interface{}{
		"source_name": input.SourceName,
		"filter_name": input.FilterName,
		"overlay":     overlay,
		"message":     fmt.Sprintf("Settings %s existing settings for filter '%s'", mode, input.FilterName),
	}
	s.recordAction("set_source_filter_settings", "Set source filter settings", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleListFilterKinds lists all available filter types
func (s *Server) handleListFilterKinds(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Listing available filter kinds")

	kinds, err := s.obsClient.GetSourceFilterKindList()
	if err != nil {
		s.recordAction("list_filter_kinds", "List filter kinds", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to list filter kinds: %w", err)
	}

	result := map[string]interface{}{
		"filter_kinds": kinds,
		"count":        len(kinds),
	}
	s.recordAction("list_filter_kinds", "List filter kinds", nil, result, true, time.Since(start))
	return nil, result, nil
}

// =============================================================================
// Transition tool handlers (FB-24)
// =============================================================================

// handleListTransitions lists all available scene transitions
func (s *Server) handleListTransitions(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Listing scene transitions")

	transitions, currentName, err := s.obsClient.GetSceneTransitionList()
	if err != nil {
		s.recordAction("list_transitions", "List transitions", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to list transitions: %w", err)
	}

	result := map[string]interface{}{
		"transitions":        transitions,
		"current_transition": currentName,
		"count":              len(transitions),
	}
	s.recordAction("list_transitions", "List transitions", nil, result, true, time.Since(start))
	return nil, result, nil
}

// handleGetCurrentTransition gets the current scene transition details
func (s *Server) handleGetCurrentTransition(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Getting current scene transition")

	transition, err := s.obsClient.GetCurrentSceneTransition()
	if err != nil {
		s.recordAction("get_current_transition", "Get current transition", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get current transition: %w", err)
	}

	result := map[string]interface{}{
		"name":         transition.Name,
		"kind":         transition.Kind,
		"duration_ms":  transition.Duration,
		"configurable": transition.Configurable,
		"settings":     transition.Settings,
	}
	s.recordAction("get_current_transition", "Get current transition", nil, result, true, time.Since(start))
	return nil, result, nil
}

// handleSetCurrentTransition sets the current scene transition
func (s *Server) handleSetCurrentTransition(ctx context.Context, request *mcpsdk.CallToolRequest, input SetCurrentTransitionInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Setting current transition to: %s", input.TransitionName)

	if err := s.obsClient.SetCurrentSceneTransition(input.TransitionName); err != nil {
		s.recordAction("set_current_transition", "Set current transition", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to set current transition: %w", err)
	}

	result := map[string]interface{}{
		"transition_name": input.TransitionName,
		"message":         fmt.Sprintf("Successfully set transition to '%s'", input.TransitionName),
	}
	s.recordAction("set_current_transition", "Set current transition", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleSetTransitionDuration sets the duration of the current scene transition
func (s *Server) handleSetTransitionDuration(ctx context.Context, request *mcpsdk.CallToolRequest, input SetTransitionDurationInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Setting transition duration to: %dms", input.TransitionDuration)

	if input.TransitionDuration <= 0 {
		s.recordAction("set_transition_duration", "Set transition duration", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("transition_duration must be greater than 0")
	}

	if err := s.obsClient.SetCurrentSceneTransitionDuration(input.TransitionDuration); err != nil {
		s.recordAction("set_transition_duration", "Set transition duration", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to set transition duration: %w", err)
	}

	result := map[string]interface{}{
		"duration_ms": input.TransitionDuration,
		"message":     fmt.Sprintf("Successfully set transition duration to %dms", input.TransitionDuration),
	}
	s.recordAction("set_transition_duration", "Set transition duration", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleTriggerTransition triggers the current transition in studio mode
func (s *Server) handleTriggerTransition(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Triggering studio mode transition")

	if err := s.obsClient.TriggerStudioModeTransition(); err != nil {
		s.recordAction("trigger_transition", "Trigger transition", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to trigger transition: %w", err)
	}

	result := SimpleResult{Message: "Successfully triggered studio mode transition"}
	s.recordAction("trigger_transition", "Trigger transition", nil, result, true, time.Since(start))
	return nil, result, nil
}

// =============================================================================
// Virtual Camera and Replay Buffer handlers (FB-25)
// =============================================================================

// handleGetVirtualCamStatus returns the virtual camera status
func (s *Server) handleGetVirtualCamStatus(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Getting virtual camera status")

	status, err := s.obsClient.GetVirtualCamStatus()
	if err != nil {
		s.recordAction("get_virtual_cam_status", "Get virtual camera status", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get virtual camera status: %w", err)
	}

	result := map[string]interface{}{
		"active":  status.Active,
		"message": fmt.Sprintf("Virtual camera is %s", map[bool]string{true: "active", false: "inactive"}[status.Active]),
	}
	s.recordAction("get_virtual_cam_status", "Get virtual camera status", nil, result, true, time.Since(start))
	return nil, result, nil
}

// handleToggleVirtualCam toggles the virtual camera on/off
func (s *Server) handleToggleVirtualCam(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Toggling virtual camera")

	active, err := s.obsClient.ToggleVirtualCam()
	if err != nil {
		s.recordAction("toggle_virtual_cam", "Toggle virtual camera", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to toggle virtual camera: %w", err)
	}

	result := map[string]interface{}{
		"active":  active,
		"message": fmt.Sprintf("Virtual camera is now %s", map[bool]string{true: "active", false: "inactive"}[active]),
	}
	s.recordAction("toggle_virtual_cam", "Toggle virtual camera", nil, result, true, time.Since(start))
	return nil, result, nil
}

// handleGetReplayBufferStatus returns the replay buffer status
func (s *Server) handleGetReplayBufferStatus(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Getting replay buffer status")

	status, err := s.obsClient.GetReplayBufferStatus()
	if err != nil {
		s.recordAction("get_replay_buffer_status", "Get replay buffer status", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get replay buffer status: %w", err)
	}

	result := map[string]interface{}{
		"active":  status.Active,
		"message": fmt.Sprintf("Replay buffer is %s", map[bool]string{true: "active", false: "inactive"}[status.Active]),
	}
	s.recordAction("get_replay_buffer_status", "Get replay buffer status", nil, result, true, time.Since(start))
	return nil, result, nil
}

// handleToggleReplayBuffer toggles the replay buffer on/off
func (s *Server) handleToggleReplayBuffer(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Toggling replay buffer")

	active, err := s.obsClient.ToggleReplayBuffer()
	if err != nil {
		s.recordAction("toggle_replay_buffer", "Toggle replay buffer", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to toggle replay buffer: %w", err)
	}

	result := map[string]interface{}{
		"active":  active,
		"message": fmt.Sprintf("Replay buffer is now %s", map[bool]string{true: "active", false: "inactive"}[active]),
	}
	s.recordAction("toggle_replay_buffer", "Toggle replay buffer", nil, result, true, time.Since(start))
	return nil, result, nil
}

// handleSaveReplayBuffer saves the current replay buffer
func (s *Server) handleSaveReplayBuffer(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Saving replay buffer")

	if err := s.obsClient.SaveReplayBuffer(); err != nil {
		s.recordAction("save_replay_buffer", "Save replay buffer", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to save replay buffer: %w", err)
	}

	result := SimpleResult{Message: "Successfully saved replay buffer"}
	s.recordAction("save_replay_buffer", "Save replay buffer", nil, result, true, time.Since(start))
	return nil, result, nil
}

// handleGetLastReplay returns the path to the last saved replay
func (s *Server) handleGetLastReplay(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Getting last replay path")

	path, err := s.obsClient.GetLastReplayBufferReplay()
	if err != nil {
		s.recordAction("get_last_replay", "Get last replay", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get last replay path: %w", err)
	}

	result := map[string]interface{}{
		"saved_replay_path": path,
		"message":           fmt.Sprintf("Last replay saved to: %s", path),
	}
	s.recordAction("get_last_replay", "Get last replay", nil, result, true, time.Since(start))
	return nil, result, nil
}

// =============================================================================
// Studio Mode handlers (FB-26)
// =============================================================================

// handleGetStudioModeEnabled returns whether studio mode is enabled
func (s *Server) handleGetStudioModeEnabled(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Getting studio mode status")

	enabled, err := s.obsClient.GetStudioModeEnabled()
	if err != nil {
		s.recordAction("get_studio_mode_enabled", "Get studio mode status", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get studio mode status: %w", err)
	}

	result := map[string]interface{}{
		"studio_mode_enabled": enabled,
		"message":             fmt.Sprintf("Studio mode is %s", map[bool]string{true: "enabled", false: "disabled"}[enabled]),
	}
	s.recordAction("get_studio_mode_enabled", "Get studio mode status", nil, result, true, time.Since(start))
	return nil, result, nil
}

// handleToggleStudioMode enables or disables studio mode
func (s *Server) handleToggleStudioMode(ctx context.Context, request *mcpsdk.CallToolRequest, input SetStudioModeInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Setting studio mode enabled=%v", input.StudioModeEnabled)

	if err := s.obsClient.SetStudioModeEnabled(input.StudioModeEnabled); err != nil {
		s.recordAction("toggle_studio_mode", "Toggle studio mode", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to set studio mode: %w", err)
	}

	result := map[string]interface{}{
		"studio_mode_enabled": input.StudioModeEnabled,
		"message":             fmt.Sprintf("Studio mode is now %s", map[bool]string{true: "enabled", false: "disabled"}[input.StudioModeEnabled]),
	}
	s.recordAction("toggle_studio_mode", "Toggle studio mode", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleGetPreviewScene returns the current preview scene in studio mode
func (s *Server) handleGetPreviewScene(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Getting preview scene")

	sceneName, err := s.obsClient.GetCurrentPreviewScene()
	if err != nil {
		s.recordAction("get_preview_scene", "Get preview scene", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to get preview scene: %w", err)
	}

	result := map[string]interface{}{
		"preview_scene": sceneName,
		"message":       fmt.Sprintf("Current preview scene: %s", sceneName),
	}
	s.recordAction("get_preview_scene", "Get preview scene", nil, result, true, time.Since(start))
	return nil, result, nil
}

// handleSetPreviewScene sets the preview scene in studio mode
func (s *Server) handleSetPreviewScene(ctx context.Context, request *mcpsdk.CallToolRequest, input SetPreviewSceneInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Setting preview scene to: %s", input.SceneName)

	if err := s.obsClient.SetCurrentPreviewScene(input.SceneName); err != nil {
		s.recordAction("set_preview_scene", "Set preview scene", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to set preview scene: %w", err)
	}

	result := map[string]interface{}{
		"preview_scene": input.SceneName,
		"message":       fmt.Sprintf("Preview scene set to: %s", input.SceneName),
	}
	s.recordAction("set_preview_scene", "Set preview scene", input, result, true, time.Since(start))
	return nil, result, nil
}

// =============================================================================
// Hotkey handlers (FB-26)
// =============================================================================

// handleTriggerHotkeyByName triggers an OBS hotkey by its name
func (s *Server) handleTriggerHotkeyByName(ctx context.Context, request *mcpsdk.CallToolRequest, input TriggerHotkeyInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Triggering hotkey: %s", input.HotkeyName)

	if err := s.obsClient.TriggerHotkeyByName(input.HotkeyName); err != nil {
		s.recordAction("trigger_hotkey_by_name", "Trigger hotkey", input, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to trigger hotkey: %w", err)
	}

	result := map[string]interface{}{
		"hotkey_name": input.HotkeyName,
		"message":     fmt.Sprintf("Successfully triggered hotkey: %s", input.HotkeyName),
	}
	s.recordAction("trigger_hotkey_by_name", "Trigger hotkey", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleListHotkeys lists all available OBS hotkey names
func (s *Server) handleListHotkeys(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Println("Listing hotkeys")

	hotkeys, err := s.obsClient.GetHotkeyList()
	if err != nil {
		s.recordAction("list_hotkeys", "List hotkeys", nil, nil, false, time.Since(start))
		return nil, nil, fmt.Errorf("failed to list hotkeys: %w", err)
	}

	result := map[string]interface{}{
		"hotkeys": hotkeys,
		"count":   len(hotkeys),
		"message": fmt.Sprintf("Found %d available hotkeys", len(hotkeys)),
	}
	s.recordAction("list_hotkeys", "List hotkeys", nil, result, true, time.Since(start))
	return nil, result, nil
}
