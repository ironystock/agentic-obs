package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ironystock/agentic-obs/internal/obs"
)

// SceneDetails contains detailed information about a scene
type SceneDetails struct {
	Name        string                   `json:"name"`
	IsActive    bool                     `json:"isActive"`
	SceneIndex  int                      `json:"sceneIndex,omitempty"`
	Sources     []map[string]interface{} `json:"sources,omitempty"`
	Description string                   `json:"description,omitempty"`
}

// registerResourceHandlers registers all resource handlers with the MCP server
func (s *Server) registerResourceHandlers() {
	// Register scenes as resources with a URI template
	// This allows accessing scenes at obs://scene/{sceneName}
	s.mcpServer.AddResourceTemplate(
		&mcpsdk.ResourceTemplate{
			URITemplate: "obs://scene/{sceneName}",
			Name:        "OBS Scene",
			Description: "Access OBS scene configuration and source details",
			MIMEType:    "application/json",
		},
		s.handleResourceRead,
	)

	log.Println("Resource handlers registered successfully")
}

// handleResourceRead returns detailed information about a specific scene
func (s *Server) handleResourceRead(ctx context.Context, request *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	uri := request.Params.URI
	log.Printf("Handling resource read request for URI: %s", uri)

	// Extract scene name from URI (format: obs://scene/{scene_name})
	sceneName, err := extractSceneNameFromURI(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid resource URI: %w", err)
	}

	// Get scene details from OBS
	scene, err := s.obsClient.GetSceneByName(sceneName)
	if err != nil {
		return nil, fmt.Errorf("failed to get scene details: %w", err)
	}

	// Get current scene to determine if this one is active
	_, currentScene, err := s.obsClient.GetSceneList()
	if err != nil {
		log.Printf("Warning: failed to get current scene: %v", err)
		currentScene = ""
	}

	// Build scene details
	details := SceneDetails{
		Name:        scene.Name,
		IsActive:    sceneName == currentScene,
		SceneIndex:  scene.Index,
		Sources:     convertSourcesToMap(scene.Sources),
		Description: fmt.Sprintf("Scene with %d sources", len(scene.Sources)),
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(details, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal scene details: %w", err)
	}

	result := &mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{
			{
				URI:      uri,
				MIMEType: "application/json",
				Text:     string(jsonData),
			},
		},
	}

	log.Printf("Returning scene details for: %s", sceneName)
	return result, nil
}

// extractSceneNameFromURI extracts the scene name from a resource URI
// Expected format: obs://scene/{scene_name}
func extractSceneNameFromURI(uri string) (string, error) {
	const prefix = "obs://scene/"
	if len(uri) <= len(prefix) {
		return "", fmt.Errorf("URI too short")
	}
	if uri[:len(prefix)] != prefix {
		return "", fmt.Errorf("URI must start with %s", prefix)
	}
	return uri[len(prefix):], nil
}

// convertSourcesToMap converts OBS SceneSource structs to generic map format for JSON serialization
func convertSourcesToMap(sources []obs.SceneSource) []map[string]interface{} {
	result := make([]map[string]interface{}, len(sources))
	for i, src := range sources {
		result[i] = map[string]interface{}{
			"id":       src.ID,
			"name":     src.Name,
			"type":     src.Type,
			"enabled":  src.Enabled,
			"visible":  src.Visible,
			"locked":   src.Locked,
			"x":        src.X,
			"y":        src.Y,
			"width":    src.Width,
			"height":   src.Height,
			"scale_x":  src.ScaleX,
			"scale_y":  src.ScaleY,
			"rotation": src.Rotation,
		}
	}
	return result
}
