package obs

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/andreykaipov/goobs/api/typedefs"
)

// Scene represents an OBS scene with its sources.
type Scene struct {
	Name    string        `json:"name"`
	Index   int           `json:"index"`
	Sources []SceneSource `json:"sources"`
}

// SceneSource represents a source within a scene.
type SceneSource struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Enabled  bool    `json:"enabled"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	ScaleX   float64 `json:"scale_x"`
	ScaleY   float64 `json:"scale_y"`
	Rotation float64 `json:"rotation"`
	Visible  bool    `json:"visible"`
	Locked   bool    `json:"locked"`
}

// RecordingStatus represents the current recording state.
type RecordingStatus struct {
	Active      bool   `json:"active"`
	Paused      bool   `json:"paused"`
	Timecode    string `json:"timecode,omitempty"`
	OutputPath  string `json:"output_path,omitempty"`
	OutputBytes int64  `json:"output_bytes,omitempty"`
}

// StreamingStatus represents the current streaming state.
type StreamingStatus struct {
	Active       bool   `json:"active"`
	Reconnecting bool   `json:"reconnecting"`
	Timecode     string `json:"timecode,omitempty"`
	TotalBytes   int64  `json:"total_bytes,omitempty"`
	TotalFrames  int    `json:"total_frames,omitempty"`
}

// GetSceneList retrieves all scenes from OBS.
// Returns a list of scene names and the currently active scene.
func (c *Client) GetSceneList() ([]string, string, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, "", err
	}

	resp, err := client.Scenes.GetSceneList()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get scene list from OBS: %w", err)
	}

	sceneNames := make([]string, len(resp.Scenes))
	for i, scene := range resp.Scenes {
		sceneNames[i] = scene.SceneName
	}

	return sceneNames, resp.CurrentProgramSceneName, nil
}

// GetSceneByName retrieves detailed information about a specific scene, including its sources.
func (c *Client) GetSceneByName(name string) (*Scene, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	// Get scene item list
	resp, err := client.SceneItems.GetSceneItemList(&sceneitems.GetSceneItemListParams{
		SceneName: &name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get scene '%s' from OBS: %w", name, err)
	}

	scene := &Scene{
		Name:    name,
		Sources: make([]SceneSource, 0, len(resp.SceneItems)),
	}

	// Convert scene items to our SceneSource format
	for _, item := range resp.SceneItems {
		source := SceneSource{
			ID:      int(item.SceneItemID),
			Name:    item.SourceName,
			Type:    item.SourceType,
			Enabled: item.SceneItemEnabled,
			Locked:  item.SceneItemLocked,
		}

		// Extract transform information (always available as a struct)
		transform := item.SceneItemTransform
		source.X = transform.PositionX
		source.Y = transform.PositionY
		source.Width = transform.Width
		source.Height = transform.Height
		source.ScaleX = transform.ScaleX
		source.ScaleY = transform.ScaleY
		source.Rotation = transform.Rotation

		scene.Sources = append(scene.Sources, source)
	}

	return scene, nil
}

// SetCurrentScene switches the active scene in OBS.
func (c *Client) SetCurrentScene(name string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneName: &name,
	})
	if err != nil {
		return fmt.Errorf("failed to set current scene to '%s': %w. Scene may not exist", name, err)
	}

	return nil
}

// CreateScene creates a new scene in OBS.
func (c *Client) CreateScene(name string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	sceneName := name
	_, err = client.Scenes.CreateScene(&scenes.CreateSceneParams{
		SceneName: &sceneName,
	})
	if err != nil {
		return fmt.Errorf("failed to create scene '%s': %w. Scene may already exist", name, err)
	}

	return nil
}

// RemoveScene deletes a scene from OBS.
func (c *Client) RemoveScene(name string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	sceneName := name
	_, err = client.Scenes.RemoveScene(&scenes.RemoveSceneParams{
		SceneName: &sceneName,
	})
	if err != nil {
		return fmt.Errorf("failed to remove scene '%s': %w. Scene may not exist or may be the only scene", name, err)
	}

	return nil
}

// StartRecording begins recording in OBS.
func (c *Client) StartRecording() error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Record.StartRecord()
	if err != nil {
		return fmt.Errorf("failed to start recording: %w. Check OBS recording settings and output path", err)
	}

	return nil
}

// StopRecording stops the current recording in OBS.
func (c *Client) StopRecording() (string, error) {
	client, err := c.getClient()
	if err != nil {
		return "", err
	}

	resp, err := client.Record.StopRecord()
	if err != nil {
		return "", fmt.Errorf("failed to stop recording: %w. Recording may not be active", err)
	}

	return resp.OutputPath, nil
}

// GetRecordingStatus retrieves the current recording status from OBS.
func (c *Client) GetRecordingStatus() (*RecordingStatus, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Record.GetRecordStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get recording status: %w", err)
	}

	status := &RecordingStatus{
		Active:      resp.OutputActive,
		Paused:      resp.OutputPaused,
		Timecode:    resp.OutputTimecode,
		OutputBytes: int64(resp.OutputBytes),
	}

	return status, nil
}

// PauseRecording pauses the current recording.
func (c *Client) PauseRecording() error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Record.PauseRecord()
	if err != nil {
		return fmt.Errorf("failed to pause recording: %w. Recording may not be active", err)
	}

	return nil
}

// ResumeRecording resumes a paused recording.
func (c *Client) ResumeRecording() error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Record.ResumeRecord()
	if err != nil {
		return fmt.Errorf("failed to resume recording: %w. Recording may not be paused", err)
	}

	return nil
}

// StartStreaming begins streaming in OBS.
func (c *Client) StartStreaming() error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Stream.StartStream()
	if err != nil {
		return fmt.Errorf("failed to start streaming: %w. Check OBS stream settings and credentials", err)
	}

	return nil
}

// StopStreaming stops the current stream in OBS.
func (c *Client) StopStreaming() error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Stream.StopStream()
	if err != nil {
		return fmt.Errorf("failed to stop streaming: %w. Stream may not be active", err)
	}

	return nil
}

// GetStreamingStatus retrieves the current streaming status from OBS.
func (c *Client) GetStreamingStatus() (*StreamingStatus, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Stream.GetStreamStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get streaming status: %w", err)
	}

	status := &StreamingStatus{
		Active:       resp.OutputActive,
		Reconnecting: resp.OutputReconnecting,
		Timecode:     resp.OutputTimecode,
		TotalBytes:   int64(resp.OutputBytes),
	}

	return status, nil
}

// ListSources retrieves all input sources available in OBS.
func (c *Client) ListSources() ([]*typedefs.Input, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Inputs.GetInputList(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list sources from OBS: %w", err)
	}

	return resp.Inputs, nil
}

// GetSourceSettings retrieves the settings for a specific source.
func (c *Client) GetSourceSettings(sourceName string) (map[string]interface{}, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	inputName := sourceName
	resp, err := client.Inputs.GetInputSettings(&inputs.GetInputSettingsParams{
		InputName: &inputName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get settings for source '%s': %w. Source may not exist", sourceName, err)
	}

	return resp.InputSettings, nil
}

// ToggleSourceVisibility toggles the visibility of a source in a specific scene.
func (c *Client) ToggleSourceVisibility(sceneName string, sourceID int) (bool, error) {
	client, err := c.getClient()
	if err != nil {
		return false, err
	}

	// First get the current state
	sceneItemIDInt := sourceID
	getResp, err := client.SceneItems.GetSceneItemEnabled(&sceneitems.GetSceneItemEnabledParams{
		SceneName:   &sceneName,
		SceneItemId: &sceneItemIDInt,
	})
	if err != nil {
		return false, fmt.Errorf("failed to get visibility state for source %d in scene '%s': %w", sourceID, sceneName, err)
	}

	// Toggle the state
	newState := !getResp.SceneItemEnabled
	newStateBool := newState
	_, err = client.SceneItems.SetSceneItemEnabled(&sceneitems.SetSceneItemEnabledParams{
		SceneName:        &sceneName,
		SceneItemId:      &sceneItemIDInt,
		SceneItemEnabled: &newStateBool,
	})
	if err != nil {
		return false, fmt.Errorf("failed to toggle visibility for source %d in scene '%s': %w", sourceID, sceneName, err)
	}

	return newState, nil
}

// GetInputMute retrieves the mute state of an audio input.
func (c *Client) GetInputMute(inputName string) (bool, error) {
	client, err := c.getClient()
	if err != nil {
		return false, err
	}

	name := inputName
	resp, err := client.Inputs.GetInputMute(&inputs.GetInputMuteParams{
		InputName: &name,
	})
	if err != nil {
		return false, fmt.Errorf("failed to get mute state for input '%s': %w. Input may not exist", inputName, err)
	}

	return resp.InputMuted, nil
}

// ToggleInputMute toggles the mute state of an audio input.
func (c *Client) ToggleInputMute(inputName string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	name := inputName
	_, err = client.Inputs.ToggleInputMute(&inputs.ToggleInputMuteParams{
		InputName: &name,
	})
	if err != nil {
		return fmt.Errorf("failed to toggle mute for input '%s': %w. Input may not exist", inputName, err)
	}

	return nil
}

// SetInputVolume sets the volume level of an audio input.
// volumeDb is in decibels (-100.0 to 26.0), or use volumeMul (0.0 to 20.0) for multiplier.
func (c *Client) SetInputVolume(inputName string, volumeDb *float64, volumeMul *float64) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	name := inputName
	params := &inputs.SetInputVolumeParams{
		InputName: &name,
	}

	if volumeDb != nil {
		params.InputVolumeDb = volumeDb
	}
	if volumeMul != nil {
		params.InputVolumeMul = volumeMul
	}

	_, err = client.Inputs.SetInputVolume(params)
	if err != nil {
		return fmt.Errorf("failed to set volume for input '%s': %w. Input may not exist", inputName, err)
	}

	return nil
}

// GetInputVolume retrieves the volume level of an audio input.
func (c *Client) GetInputVolume(inputName string) (float64, float64, error) {
	client, err := c.getClient()
	if err != nil {
		return 0, 0, err
	}

	name := inputName
	resp, err := client.Inputs.GetInputVolume(&inputs.GetInputVolumeParams{
		InputName: &name,
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get volume for input '%s': %w. Input may not exist", inputName, err)
	}

	return resp.InputVolumeDb, resp.InputVolumeMul, nil
}

// GetOBSStatus retrieves overall OBS status information.
func (c *Client) GetOBSStatus() (*OBSStatus, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	// Get version info
	versionResp, err := client.General.GetVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get OBS version: %w", err)
	}

	// Get stats
	statsResp, err := client.General.GetStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get OBS stats: %w", err)
	}

	// Get recording status
	recordStatus, err := c.GetRecordingStatus()
	if err != nil {
		// Non-fatal, continue with empty status
		recordStatus = &RecordingStatus{}
	}

	// Get streaming status
	streamStatus, err := c.GetStreamingStatus()
	if err != nil {
		// Non-fatal, continue with empty status
		streamStatus = &StreamingStatus{}
	}

	// Get current scene
	_, currentScene, err := c.GetSceneList()
	if err != nil {
		// Non-fatal, continue with empty scene name
		currentScene = ""
	}

	status := &OBSStatus{
		Version:          versionResp.ObsVersion,
		WebSocketVersion: versionResp.ObsWebSocketVersion,
		Platform:         versionResp.Platform,
		CurrentScene:     currentScene,
		Recording:        recordStatus.Active,
		Streaming:        streamStatus.Active,
		FPS:              statsResp.ActiveFps,
		FrameTime:        statsResp.AverageFrameRenderTime,
		Frames:           int(statsResp.OutputTotalFrames),
		DroppedFrames:    int(statsResp.OutputSkippedFrames),
	}

	return status, nil
}

// OBSStatus represents overall OBS status information.
type OBSStatus struct {
	Version          string  `json:"version"`
	WebSocketVersion string  `json:"websocket_version"`
	Platform         string  `json:"platform"`
	CurrentScene     string  `json:"current_scene"`
	Recording        bool    `json:"recording"`
	Streaming        bool    `json:"streaming"`
	FPS              float64 `json:"fps"`
	FrameTime        float64 `json:"frame_time_ms"`
	Frames           int     `json:"frames"`
	DroppedFrames    int     `json:"dropped_frames"`
}

// SourceState represents the visibility state of a source for preset capture/apply.
type SourceState struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// CaptureSceneState captures the current state of all sources in a scene.
// Returns source IDs, names, and their enabled (visible) states.
func (c *Client) CaptureSceneState(sceneName string) ([]SourceState, error) {
	scene, err := c.GetSceneByName(sceneName)
	if err != nil {
		return nil, fmt.Errorf("failed to capture scene state: %w", err)
	}

	states := make([]SourceState, len(scene.Sources))
	for i, src := range scene.Sources {
		states[i] = SourceState{
			ID:      src.ID,
			Name:    src.Name,
			Enabled: src.Enabled,
		}
	}

	return states, nil
}

// ApplyScenePreset applies source visibility states to a scene.
// This sets each source's enabled state according to the provided states.
func (c *Client) ApplyScenePreset(sceneName string, sources []SourceState) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	for _, src := range sources {
		sceneItemID := src.ID
		enabled := src.Enabled
		_, err := client.SceneItems.SetSceneItemEnabled(&sceneitems.SetSceneItemEnabledParams{
			SceneName:        &sceneName,
			SceneItemId:      &sceneItemID,
			SceneItemEnabled: &enabled,
		})
		if err != nil {
			return fmt.Errorf("failed to set visibility for source '%s' (ID %d): %w", src.Name, src.ID, err)
		}
	}

	return nil
}
