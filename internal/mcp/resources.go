package mcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ironystock/agentic-obs/internal/obs"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Resource URI prefixes for MCP resource identification
const (
	SceneURIPrefix         = "obs://scene/"
	ScreenshotURIPrefix    = "obs://screenshot/"
	ScreenshotURLURIPrefix = "obs://screenshot-url/"
	PresetURIPrefix        = "obs://preset/"
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
	resourceCount := 0

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
	resourceCount++

	// Register screenshot sources as resources with a URI template
	// This allows accessing screenshots at obs://screenshot/{sourceName}
	s.mcpServer.AddResourceTemplate(
		&mcpsdk.ResourceTemplate{
			URITemplate: "obs://screenshot/{sourceName}",
			Name:        "Screenshot Source",
			Description: "Latest screenshot from a configured screenshot source",
			MIMEType:    "image/png",
		},
		s.handleScreenshotResourceRead,
	)
	resourceCount++

	// Register screenshot URL resource (only if HTTP server is enabled)
	// This provides a lightweight JSON alternative to binary screenshots
	if s.httpServer != nil {
		s.mcpServer.AddResourceTemplate(
			&mcpsdk.ResourceTemplate{
				URITemplate: "obs://screenshot-url/{sourceName}",
				Name:        "Screenshot URL",
				Description: "Get HTTP URL for accessing screenshot (lightweight alternative to binary)",
				MIMEType:    "application/json",
			},
			s.handleScreenshotURLResourceRead,
		)
		resourceCount++
		log.Println("Screenshot URL resource registered (HTTP server enabled)")
	}

	// Register scene presets as resources with a URI template
	// This allows accessing presets at obs://preset/{presetName}
	s.mcpServer.AddResourceTemplate(
		&mcpsdk.ResourceTemplate{
			URITemplate: "obs://preset/{presetName}",
			Name:        "Scene Preset",
			Description: "Saved source visibility configuration for a scene",
			MIMEType:    "application/json",
		},
		s.handlePresetResourceRead,
	)
	resourceCount++

	// Register MCP-UI resources (only if HTTP server is enabled)
	if s.httpServer != nil {
		s.registerUIResources(&resourceCount)
	}

	log.Printf("Resource handlers registered successfully (%d resources)", resourceCount)
}

// registerUIResources registers UI resources for MCP-UI protocol support
func (s *Server) registerUIResources(resourceCount *int) {
	// Status Dashboard UI
	s.mcpServer.AddResourceTemplate(
		&mcpsdk.ResourceTemplate{
			URITemplate: UIStatusDashboardURI,
			Name:        "OBS Status Dashboard",
			Description: "Interactive status dashboard with connection info, recording/streaming state, and scene overview",
			MIMEType:    "text/html",
		},
		s.handleUIStatusResource,
	)
	*resourceCount++

	// Scene Preview UI
	s.mcpServer.AddResourceTemplate(
		&mcpsdk.ResourceTemplate{
			URITemplate: UIScenePreviewURI,
			Name:        "Scene Preview",
			Description: "Visual scene grid with thumbnails for interactive scene switching",
			MIMEType:    "text/html",
		},
		s.handleUIScenePreviewResource,
	)
	*resourceCount++

	// Audio Mixer UI
	s.mcpServer.AddResourceTemplate(
		&mcpsdk.ResourceTemplate{
			URITemplate: UIAudioMixerURI,
			Name:        "Audio Mixer",
			Description: "Audio input controls with volume sliders and mute toggles",
			MIMEType:    "text/html",
		},
		s.handleUIAudioMixerResource,
	)
	*resourceCount++

	// Screenshot Gallery UI
	s.mcpServer.AddResourceTemplate(
		&mcpsdk.ResourceTemplate{
			URITemplate: UIScreenshotGalleryURI,
			Name:        "Screenshot Gallery",
			Description: "Live screenshot gallery from configured screenshot sources",
			MIMEType:    "text/html",
		},
		s.handleUIScreenshotGalleryResource,
	)
	*resourceCount++

	log.Println("MCP-UI resources registered (4 UI resources)")
}

// handleUIStatusResource returns the status dashboard UI
func (s *Server) handleUIStatusResource(ctx context.Context, request *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	log.Printf("Handling UI resource read for: %s", request.Params.URI)

	// Return URL reference to the HTTP-served UI
	urlContent := fmt.Sprintf(`{"url": "%s/ui/status", "type": "url"}`, s.httpServer.GetAddr())

	return &mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{
			{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     urlContent,
			},
		},
	}, nil
}

// handleUIScenePreviewResource returns the scene preview UI
func (s *Server) handleUIScenePreviewResource(ctx context.Context, request *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	log.Printf("Handling UI resource read for: %s", request.Params.URI)

	urlContent := fmt.Sprintf(`{"url": "%s/ui/scenes", "type": "url"}`, s.httpServer.GetAddr())

	return &mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{
			{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     urlContent,
			},
		},
	}, nil
}

// handleUIAudioMixerResource returns the audio mixer UI
func (s *Server) handleUIAudioMixerResource(ctx context.Context, request *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	log.Printf("Handling UI resource read for: %s", request.Params.URI)

	urlContent := fmt.Sprintf(`{"url": "%s/ui/audio", "type": "url"}`, s.httpServer.GetAddr())

	return &mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{
			{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     urlContent,
			},
		},
	}, nil
}

// handleUIScreenshotGalleryResource returns the screenshot gallery UI
func (s *Server) handleUIScreenshotGalleryResource(ctx context.Context, request *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	log.Printf("Handling UI resource read for: %s", request.Params.URI)

	urlContent := fmt.Sprintf(`{"url": "%s/ui/screenshots", "type": "url"}`, s.httpServer.GetAddr())

	return &mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{
			{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     urlContent,
			},
		},
	}, nil
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
	if len(uri) <= len(SceneURIPrefix) {
		return "", fmt.Errorf("URI too short")
	}
	if uri[:len(SceneURIPrefix)] != SceneURIPrefix {
		return "", fmt.Errorf("URI must start with %s", SceneURIPrefix)
	}
	return uri[len(SceneURIPrefix):], nil
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

// handleScreenshotResourceRead returns the latest screenshot binary data for a screenshot source
func (s *Server) handleScreenshotResourceRead(ctx context.Context, request *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	uri := request.Params.URI
	log.Printf("Handling screenshot resource read request for URI: %s", uri)

	// Extract screenshot source name from URI (format: obs://screenshot/{sourceName})
	sourceName, err := extractScreenshotNameFromURI(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid screenshot resource URI: %w", err)
	}

	// Get screenshot source from database
	source, err := s.storage.GetScreenshotSourceByName(ctx, sourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get screenshot source: %w", err)
	}

	// Get latest screenshot for this source
	screenshot, err := s.storage.GetLatestScreenshot(ctx, source.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest screenshot: %w", err)
	}

	// Decode base64 image data to binary
	imageData, err := base64.StdEncoding.DecodeString(screenshot.ImageData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode screenshot image data (base64 length: %d): %w", len(screenshot.ImageData), err)
	}

	result := &mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{
			{
				URI:      uri,
				MIMEType: screenshot.MimeType,
				Blob:     imageData,
			},
		},
	}

	log.Printf("Returning screenshot data for source: %s (%d bytes)", sourceName, len(imageData))
	return result, nil
}

// handlePresetResourceRead returns detailed information about a specific scene preset
func (s *Server) handlePresetResourceRead(ctx context.Context, request *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	uri := request.Params.URI
	log.Printf("Handling preset resource read request for URI: %s", uri)

	// Extract preset name from URI (format: obs://preset/{presetName})
	presetName, err := extractPresetNameFromURI(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid preset resource URI: %w", err)
	}

	// Get preset from database
	preset, err := s.storage.GetScenePreset(ctx, presetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get scene preset: %w", err)
	}

	// Format preset as JSON
	presetData := map[string]interface{}{
		"name":        preset.Name,
		"scene_name":  preset.SceneName,
		"sources":     preset.Sources,
		"created_at":  preset.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"description": fmt.Sprintf("Preset for scene '%s' with %d sources", preset.SceneName, len(preset.Sources)),
	}

	jsonData, err := json.MarshalIndent(presetData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal preset details: %w", err)
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

	log.Printf("Returning preset details for: %s", presetName)
	return result, nil
}

// extractScreenshotNameFromURI extracts the screenshot source name from a resource URI
// Expected format: obs://screenshot/{sourceName}
func extractScreenshotNameFromURI(uri string) (string, error) {
	if len(uri) <= len(ScreenshotURIPrefix) {
		return "", fmt.Errorf("URI too short")
	}
	if uri[:len(ScreenshotURIPrefix)] != ScreenshotURIPrefix {
		return "", fmt.Errorf("URI must start with %s", ScreenshotURIPrefix)
	}
	return uri[len(ScreenshotURIPrefix):], nil
}

// extractPresetNameFromURI extracts the preset name from a resource URI
// Expected format: obs://preset/{presetName}
func extractPresetNameFromURI(uri string) (string, error) {
	if len(uri) <= len(PresetURIPrefix) {
		return "", fmt.Errorf("URI too short")
	}
	if uri[:len(PresetURIPrefix)] != PresetURIPrefix {
		return "", fmt.Errorf("URI must start with %s", PresetURIPrefix)
	}
	return uri[len(PresetURIPrefix):], nil
}

// handleScreenshotURLResourceRead returns JSON with the HTTP URL for a screenshot source
// This is a lightweight alternative to returning the binary screenshot data directly
func (s *Server) handleScreenshotURLResourceRead(ctx context.Context, request *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	uri := request.Params.URI
	log.Printf("Handling screenshot URL resource read request for URI: %s", uri)

	// Extract screenshot source name from URI (format: obs://screenshot-url/{sourceName})
	sourceName, err := extractScreenshotURLNameFromURI(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid screenshot-url resource URI: %w", err)
	}

	// Verify the screenshot source exists
	source, err := s.storage.GetScreenshotSourceByName(ctx, sourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get screenshot source: %w", err)
	}

	// Get latest screenshot for capture timestamp
	screenshot, err := s.storage.GetLatestScreenshot(ctx, source.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest screenshot: %w", err)
	}

	// Defensive nil check - resource should only be registered when httpServer is enabled
	if s.httpServer == nil {
		return nil, fmt.Errorf("HTTP server not available for screenshot URLs")
	}

	// Build JSON response with URL and metadata
	urlData := map[string]interface{}{
		"url":         s.httpServer.GetScreenshotURL(sourceName),
		"source":      sourceName,
		"captured_at": screenshot.CapturedAt.Format(time.RFC3339),
		"format":      source.ImageFormat,
		"mime_type":   screenshot.MimeType,
	}

	jsonData, err := json.MarshalIndent(urlData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal screenshot URL data: %w", err)
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

	log.Printf("Returning screenshot URL for source: %s", sourceName)
	return result, nil
}

// extractScreenshotURLNameFromURI extracts the screenshot source name from a screenshot-url resource URI
// Expected format: obs://screenshot-url/{sourceName}
func extractScreenshotURLNameFromURI(uri string) (string, error) {
	if len(uri) <= len(ScreenshotURLURIPrefix) {
		return "", fmt.Errorf("URI too short")
	}
	if uri[:len(ScreenshotURLURIPrefix)] != ScreenshotURLURIPrefix {
		return "", fmt.Errorf("URI must start with %s", ScreenshotURLURIPrefix)
	}
	return uri[len(ScreenshotURLURIPrefix):], nil
}
