package mcp

import (
	"context"
	"fmt"
	"log"

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
