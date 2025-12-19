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
	SceneName string `json:"scene_name,omitempty" jsonschema:"description=Optional scene name to filter presets by"`
}

// PresetNameInput is the input for preset operations by name
type PresetNameInput struct {
	PresetName string `json:"preset_name" jsonschema:"required,description=Name of the preset to operate on"`
}

// RenamePresetInput is the input for renaming a preset
type RenamePresetInput struct {
	OldName string `json:"old_name" jsonschema:"required,description=Current name of the preset to rename"`
	NewName string `json:"new_name" jsonschema:"required,description=New name for the preset"`
}

// SavePresetInput is the input for saving a scene preset
type SavePresetInput struct {
	PresetName string `json:"preset_name" jsonschema:"required,description=Name to give the new preset"`
	SceneName  string `json:"scene_name" jsonschema:"required,description=Name of the OBS scene to capture state from"`
}

// CreateScreenshotSourceInput is the input for creating a screenshot source
type CreateScreenshotSourceInput struct {
	Name        string `json:"name" jsonschema:"required,description=Unique name for this screenshot source"`
	SourceName  string `json:"source_name" jsonschema:"required,description=OBS scene or source name to capture"`
	CadenceMs   int    `json:"cadence_ms,omitempty" jsonschema:"description=Capture interval in milliseconds (default: 5000)"`
	ImageFormat string `json:"image_format,omitempty" jsonschema:"description=Image format: png or jpg (default: png)"`
	ImageWidth  int    `json:"image_width,omitempty" jsonschema:"description=Optional resize width (0 = original)"`
	ImageHeight int    `json:"image_height,omitempty" jsonschema:"description=Optional resize height (0 = original)"`
	Quality     int    `json:"quality,omitempty" jsonschema:"description=Compression quality 0-100 (default: 80)"`
}

// ScreenshotSourceNameInput is the input for screenshot source operations by name
type ScreenshotSourceNameInput struct {
	Name string `json:"name" jsonschema:"required,description=Name of the screenshot source"`
}

// ConfigureScreenshotCadenceInput is the input for updating screenshot cadence
type ConfigureScreenshotCadenceInput struct {
	Name      string `json:"name" jsonschema:"required,description=Name of the screenshot source"`
	CadenceMs int    `json:"cadence_ms" jsonschema:"required,description=New capture interval in milliseconds"`
}

// Design tool input types

// CreateTextSourceInput is the input for creating a text source
type CreateTextSourceInput struct {
	SceneName  string `json:"scene_name" jsonschema:"required,description=Name of the scene to add the source to"`
	SourceName string `json:"source_name" jsonschema:"required,description=Name for the new text source"`
	Text       string `json:"text" jsonschema:"required,description=Text content to display"`
	FontName   string `json:"font_name,omitempty" jsonschema:"description=Font face name (default: Arial)"`
	FontSize   int    `json:"font_size,omitempty" jsonschema:"description=Font size in points (default: 36)"`
	Color      int64  `json:"color,omitempty" jsonschema:"description=Text color as ABGR integer (default: white)"`
}

// CreateImageSourceInput is the input for creating an image source
type CreateImageSourceInput struct {
	SceneName  string `json:"scene_name" jsonschema:"required,description=Name of the scene to add the source to"`
	SourceName string `json:"source_name" jsonschema:"required,description=Name for the new image source"`
	FilePath   string `json:"file_path" jsonschema:"required,description=Path to the image file"`
}

// CreateColorSourceInput is the input for creating a color source
type CreateColorSourceInput struct {
	SceneName  string `json:"scene_name" jsonschema:"required,description=Name of the scene to add the source to"`
	SourceName string `json:"source_name" jsonschema:"required,description=Name for the new color source"`
	Color      int64  `json:"color" jsonschema:"required,description=Color as ABGR integer (e.g., 0xFF0000FF for red)"`
	Width      int    `json:"width,omitempty" jsonschema:"description=Width in pixels (default: 1920)"`
	Height     int    `json:"height,omitempty" jsonschema:"description=Height in pixels (default: 1080)"`
}

// CreateBrowserSourceInput is the input for creating a browser source
type CreateBrowserSourceInput struct {
	SceneName  string `json:"scene_name" jsonschema:"required,description=Name of the scene to add the source to"`
	SourceName string `json:"source_name" jsonschema:"required,description=Name for the new browser source"`
	URL        string `json:"url" jsonschema:"required,description=URL to load in the browser source"`
	Width      int    `json:"width,omitempty" jsonschema:"description=Browser width in pixels (default: 800)"`
	Height     int    `json:"height,omitempty" jsonschema:"description=Browser height in pixels (default: 600)"`
	FPS        int    `json:"fps,omitempty" jsonschema:"description=Frame rate (default: 30)"`
}

// CreateMediaSourceInput is the input for creating a media/video source
type CreateMediaSourceInput struct {
	SceneName  string `json:"scene_name" jsonschema:"required,description=Name of the scene to add the source to"`
	SourceName string `json:"source_name" jsonschema:"required,description=Name for the new media source"`
	FilePath   string `json:"file_path" jsonschema:"required,description=Path to the media file"`
	Loop       bool   `json:"loop,omitempty" jsonschema:"description=Whether to loop the media (default: false)"`
}

// SetSourceTransformInput is the input for setting source transform properties
type SetSourceTransformInput struct {
	SceneName   string   `json:"scene_name" jsonschema:"required,description=Name of the scene containing the source"`
	SceneItemID int      `json:"scene_item_id" jsonschema:"required,description=Scene item ID of the source"`
	X           *float64 `json:"x,omitempty" jsonschema:"description=X position in pixels"`
	Y           *float64 `json:"y,omitempty" jsonschema:"description=Y position in pixels"`
	ScaleX      *float64 `json:"scale_x,omitempty" jsonschema:"description=X scale factor (1.0 = 100%)"`
	ScaleY      *float64 `json:"scale_y,omitempty" jsonschema:"description=Y scale factor (1.0 = 100%)"`
	Rotation    *float64 `json:"rotation,omitempty" jsonschema:"description=Rotation in degrees"`
}

// GetSourceTransformInput is the input for getting source transform properties
type GetSourceTransformInput struct {
	SceneName   string `json:"scene_name" jsonschema:"required,description=Name of the scene containing the source"`
	SceneItemID int    `json:"scene_item_id" jsonschema:"required,description=Scene item ID of the source"`
}

// SetSourceCropInput is the input for setting source crop
type SetSourceCropInput struct {
	SceneName   string `json:"scene_name" jsonschema:"required,description=Name of the scene containing the source"`
	SceneItemID int    `json:"scene_item_id" jsonschema:"required,description=Scene item ID of the source"`
	CropTop     int    `json:"crop_top,omitempty" jsonschema:"description=Pixels to crop from top"`
	CropBottom  int    `json:"crop_bottom,omitempty" jsonschema:"description=Pixels to crop from bottom"`
	CropLeft    int    `json:"crop_left,omitempty" jsonschema:"description=Pixels to crop from left"`
	CropRight   int    `json:"crop_right,omitempty" jsonschema:"description=Pixels to crop from right"`
}

// SetSourceBoundsInput is the input for setting source bounds
type SetSourceBoundsInput struct {
	SceneName    string  `json:"scene_name" jsonschema:"required,description=Name of the scene containing the source"`
	SceneItemID  int     `json:"scene_item_id" jsonschema:"required,description=Scene item ID of the source"`
	BoundsType   string  `json:"bounds_type" jsonschema:"required,description=Bounds type: OBS_BOUNDS_NONE, OBS_BOUNDS_STRETCH, OBS_BOUNDS_SCALE_INNER, OBS_BOUNDS_SCALE_OUTER, OBS_BOUNDS_SCALE_TO_WIDTH, OBS_BOUNDS_SCALE_TO_HEIGHT, OBS_BOUNDS_MAX_ONLY"`
	BoundsWidth  float64 `json:"bounds_width,omitempty" jsonschema:"description=Bounds width in pixels"`
	BoundsHeight float64 `json:"bounds_height,omitempty" jsonschema:"description=Bounds height in pixels"`
}

// SetSourceOrderInput is the input for setting source z-order
type SetSourceOrderInput struct {
	SceneName   string `json:"scene_name" jsonschema:"required,description=Name of the scene containing the source"`
	SceneItemID int    `json:"scene_item_id" jsonschema:"required,description=Scene item ID of the source"`
	Index       int    `json:"index" jsonschema:"required,description=New index position (0 = bottom, higher = front)"`
}

// SetSourceLockedInput is the input for locking/unlocking a source
type SetSourceLockedInput struct {
	SceneName   string `json:"scene_name" jsonschema:"required,description=Name of the scene containing the source"`
	SceneItemID int    `json:"scene_item_id" jsonschema:"required,description=Scene item ID of the source"`
	Locked      bool   `json:"locked" jsonschema:"required,description=Whether the source should be locked"`
}

// DuplicateSourceInput is the input for duplicating a source
type DuplicateSourceInput struct {
	SceneName     string `json:"scene_name" jsonschema:"required,description=Name of the scene containing the source"`
	SceneItemID   int    `json:"scene_item_id" jsonschema:"required,description=Scene item ID of the source to duplicate"`
	DestSceneName string `json:"dest_scene_name,omitempty" jsonschema:"description=Destination scene name (default: same scene)"`
}

// RemoveSourceInput is the input for removing a source from a scene
type RemoveSourceInput struct {
	SceneName   string `json:"scene_name" jsonschema:"required,description=Name of the scene containing the source"`
	SceneItemID int    `json:"scene_item_id" jsonschema:"required,description=Scene item ID of the source to remove"`
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

		toolCount += 13
		log.Println("Core tools registered (13 tools)")
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
	log.Printf("Removing scene: %s", input.SceneName)

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
	log.Println("Starting streaming")

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
	log.Println("Stopping streaming")

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
	log.Printf("Deleting scene preset: %s", input.PresetName)

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
