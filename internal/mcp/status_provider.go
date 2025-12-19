package mcp

import (
	"context"

	agenthttp "github.com/ironystock/agentic-obs/internal/http"
)

// Ensure Server implements StatusProvider
var _ agenthttp.StatusProvider = (*Server)(nil)

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

	result := make([]agenthttp.SceneInfo, len(sceneNames))
	for i, name := range sceneNames {
		result[i] = agenthttp.SceneInfo{
			Name:      name,
			Index:     i,
			IsCurrent: name == currentScene,
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

		result = append(result, agenthttp.AudioInputInfo{
			Name:          input.InputName,
			Volume:        volumeMul,
			VolumePercent: volumeMul * 100,
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
		imageURL := ""
		if baseURL != "" {
			imageURL = baseURL + "/screenshot/" + source.Name
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
