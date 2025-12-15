package mcp

import (
	"context"
	"fmt"
	"log"

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

// registerToolHandlers registers all MCP tool handlers with the server
func (s *Server) registerToolHandlers() {
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

	// Source tools
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

	// Audio tools
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

	// Status tool
	mcpsdk.AddTool(s.mcpServer,
		&mcpsdk.Tool{
			Name:        "get_obs_status",
			Description: "Get overall OBS status including version, connection state, and active scene",
		},
		s.handleGetOBSStatus,
	)

	// Scene preset tools
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

	mcpsdk.AddTool(s.mcpServer,
		&mcpsdk.Tool{
			Name:        "get_input_volume",
			Description: "Get the current volume level of an audio input (returns dB and multiplier values)",
		},
		s.handleGetInputVolume,
	)

	log.Println("Tool handlers registered successfully")
}

// Tool handler implementations

func (s *Server) handleSetCurrentScene(ctx context.Context, request *mcpsdk.CallToolRequest, input SceneNameInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Setting current scene to: %s", input.SceneName)

	if err := s.obsClient.SetCurrentScene(input.SceneName); err != nil {
		return nil, nil, fmt.Errorf("failed to set current scene: %w", err)
	}

	return nil, SimpleResult{
		Message: fmt.Sprintf("Successfully switched to scene: %s", input.SceneName),
	}, nil
}

func (s *Server) handleCreateScene(ctx context.Context, request *mcpsdk.CallToolRequest, input SceneNameInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Creating scene: %s", input.SceneName)

	if err := s.obsClient.CreateScene(input.SceneName); err != nil {
		return nil, nil, fmt.Errorf("failed to create scene: %w", err)
	}

	return nil, SimpleResult{
		Message: fmt.Sprintf("Successfully created scene: %s", input.SceneName),
	}, nil
}

func (s *Server) handleRemoveScene(ctx context.Context, request *mcpsdk.CallToolRequest, input SceneNameInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Removing scene: %s", input.SceneName)

	if err := s.obsClient.RemoveScene(input.SceneName); err != nil {
		return nil, nil, fmt.Errorf("failed to remove scene: %w", err)
	}

	return nil, SimpleResult{
		Message: fmt.Sprintf("Successfully removed scene: %s", input.SceneName),
	}, nil
}

func (s *Server) handleStartRecording(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	log.Println("Starting recording")

	if err := s.obsClient.StartRecording(); err != nil {
		return nil, nil, fmt.Errorf("failed to start recording: %w", err)
	}

	return nil, SimpleResult{
		Message: "Successfully started recording",
	}, nil
}

func (s *Server) handleStopRecording(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	log.Println("Stopping recording")

	outputPath, err := s.obsClient.StopRecording()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to stop recording: %w", err)
	}

	return nil, SimpleResult{
		Message: fmt.Sprintf("Successfully stopped recording. Output saved to: %s", outputPath),
	}, nil
}

func (s *Server) handleGetRecordingStatus(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	log.Println("Getting recording status")

	status, err := s.obsClient.GetRecordingStatus()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get recording status: %w", err)
	}

	// Return the status as JSON
	return nil, status, nil
}

func (s *Server) handleStartStreaming(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	log.Println("Starting streaming")

	if err := s.obsClient.StartStreaming(); err != nil {
		return nil, nil, fmt.Errorf("failed to start streaming: %w", err)
	}

	return nil, SimpleResult{
		Message: "Successfully started streaming",
	}, nil
}

func (s *Server) handleStopStreaming(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	log.Println("Stopping streaming")

	if err := s.obsClient.StopStreaming(); err != nil {
		return nil, nil, fmt.Errorf("failed to stop streaming: %w", err)
	}

	return nil, SimpleResult{
		Message: "Successfully stopped streaming",
	}, nil
}

func (s *Server) handleGetStreamingStatus(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	log.Println("Getting streaming status")

	status, err := s.obsClient.GetStreamingStatus()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get streaming status: %w", err)
	}

	// Return the status as JSON
	return nil, status, nil
}

func (s *Server) handleGetOBSStatus(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	log.Println("Getting OBS status")

	status, err := s.obsClient.GetOBSStatus()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get OBS status: %w", err)
	}

	// Return the status as JSON
	return nil, status, nil
}

// New P1 tool handlers

func (s *Server) handleListScenes(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	log.Println("Listing all scenes")

	scenes, currentScene, err := s.obsClient.GetSceneList()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list scenes: %w", err)
	}

	return nil, map[string]interface{}{
		"scenes":        scenes,
		"current_scene": currentScene,
	}, nil
}

func (s *Server) handlePauseRecording(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	log.Println("Pausing recording")

	if err := s.obsClient.PauseRecording(); err != nil {
		return nil, nil, fmt.Errorf("failed to pause recording: %w", err)
	}

	return nil, SimpleResult{
		Message: "Successfully paused recording",
	}, nil
}

func (s *Server) handleResumeRecording(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	log.Println("Resuming recording")

	if err := s.obsClient.ResumeRecording(); err != nil {
		return nil, nil, fmt.Errorf("failed to resume recording: %w", err)
	}

	return nil, SimpleResult{
		Message: "Successfully resumed recording",
	}, nil
}

func (s *Server) handleListSources(ctx context.Context, request *mcpsdk.CallToolRequest, input struct{}) (*mcpsdk.CallToolResult, any, error) {
	log.Println("Listing all sources")

	sources, err := s.obsClient.ListSources()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list sources: %w", err)
	}

	return nil, sources, nil
}

func (s *Server) handleToggleSourceVisibility(ctx context.Context, request *mcpsdk.CallToolRequest, input SourceVisibilityInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Toggling visibility for source %d in scene: %s", input.SourceID, input.SceneName)

	newState, err := s.obsClient.ToggleSourceVisibility(input.SceneName, int(input.SourceID))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to toggle source visibility: %w", err)
	}

	return nil, map[string]interface{}{
		"scene_name": input.SceneName,
		"source_id":  input.SourceID,
		"visible":    newState,
	}, nil
}

func (s *Server) handleGetSourceSettings(ctx context.Context, request *mcpsdk.CallToolRequest, input SourceNameInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Getting settings for source: %s", input.SourceName)

	settings, err := s.obsClient.GetSourceSettings(input.SourceName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get source settings: %w", err)
	}

	return nil, settings, nil
}

func (s *Server) handleGetInputMute(ctx context.Context, request *mcpsdk.CallToolRequest, input InputNameInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Getting mute status for input: %s", input.InputName)

	isMuted, err := s.obsClient.GetInputMute(input.InputName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get input mute status: %w", err)
	}

	return nil, map[string]interface{}{
		"input_name": input.InputName,
		"is_muted":   isMuted,
	}, nil
}

func (s *Server) handleToggleInputMute(ctx context.Context, request *mcpsdk.CallToolRequest, input InputNameInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Toggling mute for input: %s", input.InputName)

	if err := s.obsClient.ToggleInputMute(input.InputName); err != nil {
		return nil, nil, fmt.Errorf("failed to toggle input mute: %w", err)
	}

	return nil, SimpleResult{
		Message: fmt.Sprintf("Successfully toggled mute for input: %s", input.InputName),
	}, nil
}

func (s *Server) handleSetInputVolume(ctx context.Context, request *mcpsdk.CallToolRequest, input SetVolumeInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Setting volume for input: %s", input.InputName)

	if err := s.obsClient.SetInputVolume(input.InputName, input.VolumeDb, input.VolumeMul); err != nil {
		return nil, nil, fmt.Errorf("failed to set input volume: %w", err)
	}

	return nil, SimpleResult{
		Message: fmt.Sprintf("Successfully set volume for input: %s", input.InputName),
	}, nil
}

// Scene preset tool handlers

// handleListScenePresets returns all saved scene presets, optionally filtered by scene name.
// Returns a list of preset summaries (id, name, scene_name, created_at) and total count.
func (s *Server) handleListScenePresets(ctx context.Context, request *mcpsdk.CallToolRequest, input ListPresetsInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Listing scene presets (filter: %s)", input.SceneName)

	presets, err := s.storage.ListScenePresets(ctx, input.SceneName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list scene presets: %w", err)
	}

	// Convert to simpler response format (without full source details)
	result := make([]map[string]interface{}, len(presets))
	for i, p := range presets {
		result[i] = map[string]interface{}{
			"id":         p.ID,
			"name":       p.Name,
			"scene_name": p.SceneName,
			"created_at": p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return nil, map[string]interface{}{
		"presets": result,
		"count":   len(presets),
	}, nil
}

// handleGetPresetDetails retrieves full details of a scene preset including source states.
// Returns the preset's id, name, scene_name, sources array, and created_at timestamp.
// Returns an error if the preset does not exist.
func (s *Server) handleGetPresetDetails(ctx context.Context, request *mcpsdk.CallToolRequest, input PresetNameInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Getting preset details for: %s", input.PresetName)

	preset, err := s.storage.GetScenePreset(ctx, input.PresetName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get preset details: %w", err)
	}

	return nil, map[string]interface{}{
		"id":         preset.ID,
		"name":       preset.Name,
		"scene_name": preset.SceneName,
		"sources":    preset.Sources,
		"created_at": preset.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// handleDeleteScenePreset permanently removes a scene preset from storage.
// Returns a success message or an error if the preset does not exist.
func (s *Server) handleDeleteScenePreset(ctx context.Context, request *mcpsdk.CallToolRequest, input PresetNameInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Deleting scene preset: %s", input.PresetName)

	if err := s.storage.DeleteScenePreset(ctx, input.PresetName); err != nil {
		return nil, nil, fmt.Errorf("failed to delete preset: %w", err)
	}

	return nil, SimpleResult{
		Message: fmt.Sprintf("Successfully deleted preset: %s", input.PresetName),
	}, nil
}

// handleRenameScenePreset changes the name of an existing scene preset.
// Returns a success message or an error if the preset does not exist or the new name conflicts.
func (s *Server) handleRenameScenePreset(ctx context.Context, request *mcpsdk.CallToolRequest, input RenamePresetInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Renaming preset from '%s' to '%s'", input.OldName, input.NewName)

	if err := s.storage.RenameScenePreset(ctx, input.OldName, input.NewName); err != nil {
		return nil, nil, fmt.Errorf("failed to rename preset: %w", err)
	}

	return nil, SimpleResult{
		Message: fmt.Sprintf("Successfully renamed preset from '%s' to '%s'", input.OldName, input.NewName),
	}, nil
}

// handleGetInputVolume retrieves the current volume level of an audio input.
// Returns volume_db (decibels) and volume_mul (linear multiplier) values.
// Returns an error if the input does not exist or OBS is not connected.
func (s *Server) handleGetInputVolume(ctx context.Context, request *mcpsdk.CallToolRequest, input InputNameInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Getting volume for input: %s", input.InputName)

	volumeDb, volumeMul, err := s.obsClient.GetInputVolume(input.InputName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get input volume: %w", err)
	}

	return nil, map[string]interface{}{
		"input_name": input.InputName,
		"volume_db":  volumeDb,
		"volume_mul": volumeMul,
	}, nil
}

// handleSaveScenePreset captures the current source visibility states from an OBS scene
// and saves them as a named preset in storage. Returns the preset id, name, scene_name,
// source_count, and a success message. Returns an error if the scene does not exist,
// OBS is not connected, or a preset with the same name already exists.
func (s *Server) handleSaveScenePreset(ctx context.Context, request *mcpsdk.CallToolRequest, input SavePresetInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Saving scene preset '%s' for scene '%s'", input.PresetName, input.SceneName)

	// Capture current scene state from OBS
	states, err := s.obsClient.CaptureSceneState(input.SceneName)
	if err != nil {
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
		return nil, nil, fmt.Errorf("failed to save preset: %w", err)
	}

	return nil, map[string]interface{}{
		"id":           id,
		"preset_name":  input.PresetName,
		"scene_name":   input.SceneName,
		"source_count": len(sources),
		"message":      fmt.Sprintf("Successfully saved preset '%s' with %d sources", input.PresetName, len(sources)),
	}, nil
}

// handleApplyScenePreset loads a saved preset and applies its source visibility states
// to the target OBS scene. Sources that no longer exist in the scene are skipped.
// Returns the preset_name, scene_name, applied_count, and a success message.
// Returns an error if the preset does not exist, the scene no longer exists, or OBS is not connected.
func (s *Server) handleApplyScenePreset(ctx context.Context, request *mcpsdk.CallToolRequest, input PresetNameInput) (*mcpsdk.CallToolResult, any, error) {
	log.Printf("Applying scene preset: %s", input.PresetName)

	// Load preset from storage
	preset, err := s.storage.GetScenePreset(ctx, input.PresetName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load preset: %w", err)
	}

	// Get current scene items to map names to IDs
	scene, err := s.obsClient.GetSceneByName(preset.SceneName)
	if err != nil {
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
		return nil, nil, fmt.Errorf("failed to apply preset: %w", err)
	}

	return nil, map[string]interface{}{
		"preset_name":   input.PresetName,
		"scene_name":    preset.SceneName,
		"applied_count": len(obsStates),
		"message":       fmt.Sprintf("Successfully applied preset '%s' to scene '%s'", input.PresetName, preset.SceneName),
	}, nil
}
