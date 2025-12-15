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
