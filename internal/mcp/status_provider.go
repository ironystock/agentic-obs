package mcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"net/url"
	"strings"

	agenthttp "github.com/ironystock/agentic-obs/internal/http"
	"github.com/ironystock/agentic-obs/internal/obs"
)

// Ensure Server implements StatusProvider and ActionExecutor
var _ agenthttp.StatusProvider = (*Server)(nil)
var _ agenthttp.ActionExecutor = (*Server)(nil)

// GetStatus returns the current OBS status.
func (s *Server) GetStatus() (any, error) {
	if !s.obsClient.IsConnected() {
		return map[string]any{
			"connected":    false,
			"recording":    false,
			"streaming":    false,
			"currentScene": "",
		}, nil
	}

	status, err := s.obsClient.GetOBSStatus()
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"connected":        true,
		"recording":        status.Recording,
		"streaming":        status.Streaming,
		"currentScene":     status.CurrentScene,
		"obsVersion":       status.Version,
		"websocketVersion": status.WebSocketVersion,
	}, nil
}

// GetScenes returns a list of scenes for UI display.
func (s *Server) GetScenes() ([]agenthttp.SceneInfo, error) {
	if !s.obsClient.IsConnected() {
		return []agenthttp.SceneInfo{}, nil
	}

	sceneNames, currentScene, err := s.obsClient.GetSceneList()
	if err != nil {
		return nil, err
	}

	// Build base URL for thumbnails
	baseURL := ""
	if s.httpServer != nil {
		baseURL = s.httpServer.GetAddr()
	}

	result := make([]agenthttp.SceneInfo, len(sceneNames))
	for i, name := range sceneNames {
		// Get source count for each scene
		sourceCount := 0
		scene, err := s.obsClient.GetSceneByName(name)
		if err == nil && scene != nil {
			sourceCount = len(scene.Sources)
		}

		// Build thumbnail URL (encode scene name for URL safety)
		thumbnailURL := ""
		if baseURL != "" {
			thumbnailURL = baseURL + "/ui/scene-thumbnail/" + url.PathEscape(name)
		}

		result[i] = agenthttp.SceneInfo{
			Name:         name,
			Index:        i,
			IsCurrent:    name == currentScene,
			SourceCount:  sourceCount,
			ThumbnailURL: thumbnailURL,
		}
	}

	return result, nil
}

// GetAudioInputs returns a list of audio inputs with their current state.
func (s *Server) GetAudioInputs() ([]agenthttp.AudioInputInfo, error) {
	if !s.obsClient.IsConnected() {
		return []agenthttp.AudioInputInfo{}, nil
	}

	inputs, err := s.obsClient.ListSources()
	if err != nil {
		return nil, err
	}

	result := []agenthttp.AudioInputInfo{}
	for _, input := range inputs {
		// Only include audio inputs (those that support volume control)
		volumeMul, volumeDB, err := s.obsClient.GetInputVolume(input.InputName)
		if err != nil {
			continue // Skip non-audio inputs
		}

		muted, _ := s.obsClient.GetInputMute(input.InputName)

		// Convert dB to slider percentage (0-100) using logarithmic curve
		// Formula: sliderPercent = 10^(dB/30) * 100
		// This gives: 0dB=100%, -10dB≈46%, -20dB≈22%, -30dB=10%, -60dB≈1%
		volumePercent := 0.0
		if volumeMul > 0 && volumeDB > -60 {
			volumePercent = math.Pow(10, volumeDB/30) * 100
			if volumePercent > 100 {
				volumePercent = 100 // Clamp for any gain above 0dB
			}
		}

		log.Printf("[Audio] Input %q: volumeMul=%.4f, volumeDB=%.2f, volumePercent=%.1f, muted=%v",
			input.InputName, volumeMul, volumeDB, volumePercent, muted)

		result = append(result, agenthttp.AudioInputInfo{
			Name:          input.InputName,
			Volume:        volumeMul,
			VolumePercent: volumePercent,
			VolumeDB:      volumeDB,
			IsMuted:       muted,
			InputKind:     input.InputKind,
		})
	}

	return result, nil
}

// GetScreenshotSources returns a list of screenshot sources for UI display.
func (s *Server) GetScreenshotSources() ([]agenthttp.ScreenshotSourceInfo, error) {
	ctx := context.Background()

	sources, err := s.storage.ListScreenshotSources(ctx)
	if err != nil {
		return nil, err
	}

	baseURL := ""
	if s.httpServer != nil {
		baseURL = s.httpServer.GetAddr()
	}

	result := make([]agenthttp.ScreenshotSourceInfo, len(sources))
	for i, source := range sources {
		// Build image URL (encode source name for URL safety)
		imageURL := ""
		if baseURL != "" {
			imageURL = baseURL + "/screenshot/" + url.PathEscape(source.Name)
		}

		lastCapture := "Never"
		screenshot, err := s.storage.GetLatestScreenshot(ctx, source.ID)
		if err == nil && screenshot != nil {
			lastCapture = screenshot.CapturedAt.Format("2006-01-02 15:04:05")
		}

		result[i] = agenthttp.ScreenshotSourceInfo{
			Name:        source.Name,
			SourceName:  source.SourceName,
			CadenceMs:   source.CadenceMs,
			ImageURL:    imageURL,
			LastCapture: lastCapture,
		}
	}

	return result, nil
}

// ActionExecutor implementation

// SetCurrentScene switches to the specified scene.
func (s *Server) SetCurrentScene(sceneName string) error {
	if !s.obsClient.IsConnected() {
		return fmt.Errorf("OBS not connected")
	}
	return s.obsClient.SetCurrentScene(sceneName)
}

// ToggleInputMute toggles mute state for an audio input.
func (s *Server) ToggleInputMute(inputName string) error {
	if !s.obsClient.IsConnected() {
		return fmt.Errorf("OBS not connected")
	}
	return s.obsClient.ToggleInputMute(inputName)
}

// SetInputVolume sets the volume for an audio input.
func (s *Server) SetInputVolume(inputName string, volumeDb float64) error {
	if !s.obsClient.IsConnected() {
		return fmt.Errorf("OBS not connected")
	}
	return s.obsClient.SetInputVolume(inputName, &volumeDb, nil)
}

// TakeSceneThumbnail captures a thumbnail of the specified scene.
// Uses caching to reduce OBS API calls for frequently accessed scenes.
func (s *Server) TakeSceneThumbnail(sceneName string) ([]byte, string, error) {
	if !s.obsClient.IsConnected() {
		return nil, "", fmt.Errorf("OBS not connected")
	}

	// Check cache first
	if s.thumbnailCache != nil {
		if imageData, mimeType, ok := s.thumbnailCache.get(sceneName); ok {
			log.Printf("[Thumbnail] Cache hit for scene: %s", sceneName)
			return imageData, mimeType, nil
		}
	}

	// Cache miss - capture new thumbnail
	log.Printf("[Thumbnail] Cache miss for scene: %s, capturing...", sceneName)

	// Take screenshot of the scene
	opts := obs.ScreenshotOptions{
		SourceName: sceneName,
		Format:     "jpg",
		Width:      320,
		Height:     180,
		Quality:    75,
	}

	dataURI, err := s.obsClient.TakeSourceScreenshot(opts)
	if err != nil {
		return nil, "", err
	}

	// Parse data URI: data:image/jpeg;base64,<data>
	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) != 2 {
		return nil, "", fmt.Errorf("invalid data URI format")
	}

	// Extract MIME type
	mimeType := "image/jpeg"
	if strings.Contains(parts[0], "image/png") {
		mimeType = "image/png"
	}

	// Decode base64
	imageData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Store in cache
	if s.thumbnailCache != nil {
		s.thumbnailCache.set(sceneName, imageData, mimeType)
	}

	return imageData, mimeType, nil
}
