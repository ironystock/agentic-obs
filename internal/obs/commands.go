package obs

import (
	"fmt"
	"strings"

	"github.com/andreykaipov/goobs/api/requests/filters"
	"github.com/andreykaipov/goobs/api/requests/general"
	"github.com/andreykaipov/goobs/api/requests/inputs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/andreykaipov/goobs/api/requests/sources"
	"github.com/andreykaipov/goobs/api/requests/transitions"
	"github.com/andreykaipov/goobs/api/requests/ui"
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
// Returns (volumeMultiplier, volumeDb, error).
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

	return resp.InputVolumeMul, resp.InputVolumeDb, nil
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
// TODO: For scenes with >25 sources, consider implementing batch operations or
// concurrent processing to improve performance.
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
// TODO: For scenes with >25 sources, consider implementing batch operations to improve
// performance and provide better error handling (collect partial failures rather than
// failing on first error).
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

// ScreenshotOptions configures screenshot capture settings.
type ScreenshotOptions struct {
	SourceName string // Name of the source or scene to capture
	Format     string // Image format: "png" or "jpg" (default: "png")
	Width      int    // Optional resize width (0 = original)
	Height     int    // Optional resize height (0 = original)
	Quality    int    // Compression quality 0-100 (default: -1 = library default)
}

// TakeSourceScreenshot captures a screenshot of the specified source or scene.
// Returns the Base64-encoded image data with data URI prefix (e.g., "data:image/png;base64,...").
func (c *Client) TakeSourceScreenshot(opts ScreenshotOptions) (string, error) {
	client, err := c.getClient()
	if err != nil {
		return "", err
	}

	// Set defaults
	imageFormat := opts.Format
	if imageFormat == "" {
		imageFormat = "png"
	}

	// Build parameters
	params := &sources.GetSourceScreenshotParams{
		SourceName:  &opts.SourceName,
		ImageFormat: &imageFormat,
	}

	// Optional width/height
	if opts.Width > 0 {
		width := float64(opts.Width)
		params.ImageWidth = &width
	}
	if opts.Height > 0 {
		height := float64(opts.Height)
		params.ImageHeight = &height
	}

	// Optional quality (for JPEG)
	if opts.Quality >= 0 && opts.Quality <= 100 {
		quality := float64(opts.Quality)
		params.ImageCompressionQuality = &quality
	}

	resp, err := client.Sources.GetSourceScreenshot(params)
	if err != nil {
		return "", fmt.Errorf("failed to capture screenshot of '%s': %w", opts.SourceName, err)
	}

	// OBS returns data as a data URI (e.g., "data:image/png;base64,iVBORw0...")
	// Strip the prefix to return clean base64 data
	imageData := resp.ImageData
	if idx := strings.Index(imageData, ","); idx != -1 {
		imageData = imageData[idx+1:]
	}

	return imageData, nil
}

// BrowserSourceSettings holds configuration for creating a browser source.
type BrowserSourceSettings struct {
	URL         string // URL to display
	Width       int    // Width in pixels
	Height      int    // Height in pixels
	RefreshRate int    // Page refresh rate in seconds (0 = no refresh)
	CSS         string // Custom CSS to inject
}

// SceneItemTransform represents the transform properties of a scene item.
type SceneItemTransform struct {
	PositionX    float64 `json:"position_x"`
	PositionY    float64 `json:"position_y"`
	ScaleX       float64 `json:"scale_x"`
	ScaleY       float64 `json:"scale_y"`
	Rotation     float64 `json:"rotation"`
	Width        float64 `json:"width"`
	Height       float64 `json:"height"`
	SourceWidth  float64 `json:"source_width"`
	SourceHeight float64 `json:"source_height"`
	BoundsType   string  `json:"bounds_type"`
	BoundsWidth  float64 `json:"bounds_width"`
	BoundsHeight float64 `json:"bounds_height"`
	CropTop      int     `json:"crop_top"`
	CropBottom   int     `json:"crop_bottom"`
	CropLeft     int     `json:"crop_left"`
	CropRight    int     `json:"crop_right"`
}

// CreateBrowserSource creates a new browser source in the specified scene.
// Returns the scene item ID of the created source.
func (c *Client) CreateBrowserSource(sceneName, sourceName string, settings BrowserSourceSettings) (int, error) {
	client, err := c.getClient()
	if err != nil {
		return 0, err
	}

	// Build input settings
	inputSettings := map[string]interface{}{
		"url":    settings.URL,
		"width":  settings.Width,
		"height": settings.Height,
	}

	if settings.RefreshRate > 0 {
		inputSettings["refresh_rate"] = settings.RefreshRate
	}
	if settings.CSS != "" {
		inputSettings["css"] = settings.CSS
	}

	// Create the browser source
	enabled := true
	inputKind := "browser_source"
	resp, err := client.Inputs.CreateInput(&inputs.CreateInputParams{
		SceneName:        &sceneName,
		InputName:        &sourceName,
		InputKind:        &inputKind,
		InputSettings:    inputSettings,
		SceneItemEnabled: &enabled,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create browser source '%s' in scene '%s': %w", sourceName, sceneName, err)
	}

	return int(resp.SceneItemId), nil
}

// CreateInput creates a new input source in the specified scene.
// Returns the scene item ID of the created source.
func (c *Client) CreateInput(sceneName, sourceName, inputKind string, settings map[string]interface{}) (int, error) {
	client, err := c.getClient()
	if err != nil {
		return 0, err
	}

	enabled := true
	resp, err := client.Inputs.CreateInput(&inputs.CreateInputParams{
		SceneName:        &sceneName,
		InputName:        &sourceName,
		InputKind:        &inputKind,
		InputSettings:    settings,
		SceneItemEnabled: &enabled,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create input '%s' of type '%s' in scene '%s': %w", sourceName, inputKind, sceneName, err)
	}

	return int(resp.SceneItemId), nil
}

// GetSceneItemTransform retrieves the transform properties of a scene item.
func (c *Client) GetSceneItemTransform(sceneName string, sceneItemID int) (*SceneItemTransform, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.SceneItems.GetSceneItemTransform(&sceneitems.GetSceneItemTransformParams{
		SceneName:   &sceneName,
		SceneItemId: &sceneItemID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get transform for item %d in scene '%s': %w", sceneItemID, sceneName, err)
	}

	t := resp.SceneItemTransform
	return &SceneItemTransform{
		PositionX:    t.PositionX,
		PositionY:    t.PositionY,
		ScaleX:       t.ScaleX,
		ScaleY:       t.ScaleY,
		Rotation:     t.Rotation,
		Width:        t.Width,
		Height:       t.Height,
		SourceWidth:  t.SourceWidth,
		SourceHeight: t.SourceHeight,
		BoundsType:   t.BoundsType,
		BoundsWidth:  t.BoundsWidth,
		BoundsHeight: t.BoundsHeight,
		CropTop:      int(t.CropTop),
		CropBottom:   int(t.CropBottom),
		CropLeft:     int(t.CropLeft),
		CropRight:    int(t.CropRight),
	}, nil
}

// SetSceneItemTransform sets the transform properties of a scene item.
func (c *Client) SetSceneItemTransform(sceneName string, sceneItemID int, transform *SceneItemTransform) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	// Build transform struct
	t := &typedefs.SceneItemTransform{
		PositionX:    transform.PositionX,
		PositionY:    transform.PositionY,
		ScaleX:       transform.ScaleX,
		ScaleY:       transform.ScaleY,
		Rotation:     transform.Rotation,
		BoundsType:   transform.BoundsType,
		BoundsWidth:  transform.BoundsWidth,
		BoundsHeight: transform.BoundsHeight,
		CropTop:      float64(transform.CropTop),
		CropBottom:   float64(transform.CropBottom),
		CropLeft:     float64(transform.CropLeft),
		CropRight:    float64(transform.CropRight),
	}

	_, err = client.SceneItems.SetSceneItemTransform(&sceneitems.SetSceneItemTransformParams{
		SceneName:          &sceneName,
		SceneItemId:        &sceneItemID,
		SceneItemTransform: t,
	})
	if err != nil {
		return fmt.Errorf("failed to set transform for item %d in scene '%s': %w", sceneItemID, sceneName, err)
	}

	return nil
}

// SetSceneItemIndex sets the z-order index of a scene item.
func (c *Client) SetSceneItemIndex(sceneName string, sceneItemID int, index int) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.SceneItems.SetSceneItemIndex(&sceneitems.SetSceneItemIndexParams{
		SceneName:      &sceneName,
		SceneItemId:    &sceneItemID,
		SceneItemIndex: &index,
	})
	if err != nil {
		return fmt.Errorf("failed to set index for item %d in scene '%s': %w", sceneItemID, sceneName, err)
	}

	return nil
}

// SetSceneItemLocked sets the locked state of a scene item.
func (c *Client) SetSceneItemLocked(sceneName string, sceneItemID int, locked bool) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.SceneItems.SetSceneItemLocked(&sceneitems.SetSceneItemLockedParams{
		SceneName:       &sceneName,
		SceneItemId:     &sceneItemID,
		SceneItemLocked: &locked,
	})
	if err != nil {
		return fmt.Errorf("failed to set locked state for item %d in scene '%s': %w", sceneItemID, sceneName, err)
	}

	return nil
}

// GetSceneItemLocked gets the locked state of a scene item.
func (c *Client) GetSceneItemLocked(sceneName string, sceneItemID int) (bool, error) {
	client, err := c.getClient()
	if err != nil {
		return false, err
	}

	resp, err := client.SceneItems.GetSceneItemLocked(&sceneitems.GetSceneItemLockedParams{
		SceneName:   &sceneName,
		SceneItemId: &sceneItemID,
	})
	if err != nil {
		return false, fmt.Errorf("failed to get locked state for item %d in scene '%s': %w", sceneItemID, sceneName, err)
	}

	return resp.SceneItemLocked, nil
}

// DuplicateSceneItem duplicates a scene item within the same scene or to another scene.
// Returns the scene item ID of the duplicated item.
func (c *Client) DuplicateSceneItem(sceneName string, sceneItemID int, destScene string) (int, error) {
	client, err := c.getClient()
	if err != nil {
		return 0, err
	}

	params := &sceneitems.DuplicateSceneItemParams{
		SceneName:   &sceneName,
		SceneItemId: &sceneItemID,
	}

	// If destination scene is specified and different, set it
	if destScene != "" && destScene != sceneName {
		params.DestinationSceneName = &destScene
	}

	resp, err := client.SceneItems.DuplicateSceneItem(params)
	if err != nil {
		return 0, fmt.Errorf("failed to duplicate item %d from scene '%s': %w", sceneItemID, sceneName, err)
	}

	return int(resp.SceneItemId), nil
}

// RemoveSceneItem removes a scene item from a scene.
func (c *Client) RemoveSceneItem(sceneName string, sceneItemID int) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.SceneItems.RemoveSceneItem(&sceneitems.RemoveSceneItemParams{
		SceneName:   &sceneName,
		SceneItemId: &sceneItemID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove item %d from scene '%s': %w", sceneItemID, sceneName, err)
	}

	return nil
}

// GetInputKindList returns a list of available input kinds (source types).
func (c *Client) GetInputKindList() ([]string, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Inputs.GetInputKindList(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get input kind list: %w", err)
	}

	return resp.InputKinds, nil
}

// =============================================================================
// Filter Types and Methods (FB-23)
// =============================================================================

// FilterInfo represents basic information about a source filter.
type FilterInfo struct {
	Name    string `json:"name"`
	Kind    string `json:"kind"`
	Index   int    `json:"index"`
	Enabled bool   `json:"enabled"`
}

// FilterDetails represents detailed information about a source filter.
type FilterDetails struct {
	Name     string                 `json:"name"`
	Kind     string                 `json:"kind"`
	Index    int                    `json:"index"`
	Enabled  bool                   `json:"enabled"`
	Settings map[string]interface{} `json:"settings"`
}

// GetSourceFilterList retrieves all filters on a source.
func (c *Client) GetSourceFilterList(sourceName string) ([]FilterInfo, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Filters.GetSourceFilterList(&filters.GetSourceFilterListParams{
		SourceName: &sourceName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get filters for source '%s': %w", sourceName, err)
	}

	result := make([]FilterInfo, len(resp.Filters))
	for i, f := range resp.Filters {
		result[i] = FilterInfo{
			Name:    f.FilterName,
			Kind:    f.FilterKind,
			Index:   int(f.FilterIndex),
			Enabled: f.FilterEnabled,
		}
	}

	return result, nil
}

// GetSourceFilter retrieves details about a specific filter on a source.
func (c *Client) GetSourceFilter(sourceName, filterName string) (*FilterDetails, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Filters.GetSourceFilter(&filters.GetSourceFilterParams{
		SourceName: &sourceName,
		FilterName: &filterName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get filter '%s' on source '%s': %w", filterName, sourceName, err)
	}

	return &FilterDetails{
		Name:     filterName,
		Kind:     resp.FilterKind,
		Index:    int(resp.FilterIndex),
		Enabled:  resp.FilterEnabled,
		Settings: resp.FilterSettings,
	}, nil
}

// CreateSourceFilter creates a new filter on a source.
func (c *Client) CreateSourceFilter(sourceName, filterName, filterKind string, settings map[string]interface{}) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	params := &filters.CreateSourceFilterParams{
		SourceName: &sourceName,
		FilterName: &filterName,
		FilterKind: &filterKind,
	}

	if settings != nil {
		params.FilterSettings = settings
	}

	_, err = client.Filters.CreateSourceFilter(params)
	if err != nil {
		return fmt.Errorf("failed to create filter '%s' on source '%s': %w", filterName, sourceName, err)
	}

	return nil
}

// RemoveSourceFilter removes a filter from a source.
func (c *Client) RemoveSourceFilter(sourceName, filterName string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Filters.RemoveSourceFilter(&filters.RemoveSourceFilterParams{
		SourceName: &sourceName,
		FilterName: &filterName,
	})
	if err != nil {
		return fmt.Errorf("failed to remove filter '%s' from source '%s': %w", filterName, sourceName, err)
	}

	return nil
}

// SetSourceFilterEnabled enables or disables a filter on a source.
func (c *Client) SetSourceFilterEnabled(sourceName, filterName string, enabled bool) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Filters.SetSourceFilterEnabled(&filters.SetSourceFilterEnabledParams{
		SourceName:    &sourceName,
		FilterName:    &filterName,
		FilterEnabled: &enabled,
	})
	if err != nil {
		return fmt.Errorf("failed to set filter '%s' enabled=%v on source '%s': %w", filterName, enabled, sourceName, err)
	}

	return nil
}

// SetSourceFilterSettings updates the settings of a filter.
// If overlay is true, settings are merged with existing; otherwise they replace entirely.
func (c *Client) SetSourceFilterSettings(sourceName, filterName string, settings map[string]interface{}, overlay bool) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Filters.SetSourceFilterSettings(&filters.SetSourceFilterSettingsParams{
		SourceName:     &sourceName,
		FilterName:     &filterName,
		FilterSettings: settings,
		Overlay:        &overlay,
	})
	if err != nil {
		return fmt.Errorf("failed to set settings for filter '%s' on source '%s': %w", filterName, sourceName, err)
	}

	return nil
}

// GetSourceFilterKindList retrieves all available filter types.
func (c *Client) GetSourceFilterKindList() ([]string, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Filters.GetSourceFilterKindList(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get filter kind list: %w", err)
	}

	return resp.SourceFilterKinds, nil
}

// =============================================================================
// Transition Types and Methods (FB-24)
// =============================================================================

// TransitionInfo represents basic information about a scene transition.
type TransitionInfo struct {
	Name         string `json:"name"`
	Kind         string `json:"kind"`
	Fixed        bool   `json:"fixed"`
	Configurable bool   `json:"configurable"`
}

// TransitionDetails represents the current scene transition with settings.
type TransitionDetails struct {
	Name         string                 `json:"name"`
	Kind         string                 `json:"kind"`
	Duration     int                    `json:"duration_ms"`
	Configurable bool                   `json:"configurable"`
	Settings     map[string]interface{} `json:"settings,omitempty"`
}

// GetSceneTransitionList retrieves all available scene transitions.
func (c *Client) GetSceneTransitionList() ([]TransitionInfo, string, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, "", err
	}

	resp, err := client.Transitions.GetSceneTransitionList()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get transition list: %w", err)
	}

	result := make([]TransitionInfo, len(resp.Transitions))
	for i, t := range resp.Transitions {
		result[i] = TransitionInfo{
			Name:         t.TransitionName,
			Kind:         t.TransitionKind,
			Fixed:        t.TransitionFixed,
			Configurable: t.TransitionConfigurable,
		}
	}

	return result, resp.CurrentSceneTransitionName, nil
}

// GetCurrentSceneTransition retrieves the current scene transition details.
func (c *Client) GetCurrentSceneTransition() (*TransitionDetails, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Transitions.GetCurrentSceneTransition()
	if err != nil {
		return nil, fmt.Errorf("failed to get current transition: %w", err)
	}

	return &TransitionDetails{
		Name:         resp.TransitionName,
		Kind:         resp.TransitionKind,
		Duration:     int(resp.TransitionDuration),
		Configurable: resp.TransitionConfigurable,
		Settings:     resp.TransitionSettings,
	}, nil
}

// SetCurrentSceneTransition sets the current scene transition.
func (c *Client) SetCurrentSceneTransition(transitionName string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Transitions.SetCurrentSceneTransition(&transitions.SetCurrentSceneTransitionParams{
		TransitionName: &transitionName,
	})
	if err != nil {
		return fmt.Errorf("failed to set current transition to '%s': %w", transitionName, err)
	}

	return nil
}

// SetCurrentSceneTransitionDuration sets the duration of the current scene transition.
func (c *Client) SetCurrentSceneTransitionDuration(durationMs int) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	duration := float64(durationMs)
	_, err = client.Transitions.SetCurrentSceneTransitionDuration(&transitions.SetCurrentSceneTransitionDurationParams{
		TransitionDuration: &duration,
	})
	if err != nil {
		return fmt.Errorf("failed to set transition duration to %dms: %w", durationMs, err)
	}

	return nil
}

// TriggerStudioModeTransition triggers the current scene transition in studio mode.
// This executes the transition from preview to program.
// Returns an error if studio mode is not enabled.
func (c *Client) TriggerStudioModeTransition() error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Transitions.TriggerStudioModeTransition()
	if err != nil {
		// Check if it's because studio mode isn't enabled
		if strings.Contains(err.Error(), "studio mode") {
			return fmt.Errorf("studio mode is not enabled. Enable studio mode in OBS to use this feature")
		}
		return fmt.Errorf("failed to trigger studio mode transition: %w", err)
	}

	return nil
}

// =============================================================================
// Virtual Camera Types and Methods (FB-25)
// =============================================================================

// VirtualCamStatus represents the current virtual camera state.
type VirtualCamStatus struct {
	Active bool `json:"active"`
}

// GetVirtualCamStatus retrieves the current virtual camera status.
func (c *Client) GetVirtualCamStatus() (*VirtualCamStatus, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Outputs.GetVirtualCamStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual camera status: %w", err)
	}

	return &VirtualCamStatus{
		Active: resp.OutputActive,
	}, nil
}

// ToggleVirtualCam toggles the virtual camera on or off.
// Returns the new active state.
func (c *Client) ToggleVirtualCam() (bool, error) {
	client, err := c.getClient()
	if err != nil {
		return false, err
	}

	resp, err := client.Outputs.ToggleVirtualCam()
	if err != nil {
		return false, fmt.Errorf("failed to toggle virtual camera: %w", err)
	}

	return resp.OutputActive, nil
}

// =============================================================================
// Replay Buffer Types and Methods (FB-25)
// =============================================================================

// ReplayBufferStatus represents the current replay buffer state.
type ReplayBufferStatus struct {
	Active bool `json:"active"`
}

// GetReplayBufferStatus retrieves the current replay buffer status.
func (c *Client) GetReplayBufferStatus() (*ReplayBufferStatus, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Outputs.GetReplayBufferStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get replay buffer status: %w", err)
	}

	return &ReplayBufferStatus{
		Active: resp.OutputActive,
	}, nil
}

// ToggleReplayBuffer toggles the replay buffer on or off.
// Returns the new active state.
func (c *Client) ToggleReplayBuffer() (bool, error) {
	client, err := c.getClient()
	if err != nil {
		return false, err
	}

	resp, err := client.Outputs.ToggleReplayBuffer()
	if err != nil {
		return false, fmt.Errorf("failed to toggle replay buffer: %w", err)
	}

	return resp.OutputActive, nil
}

// SaveReplayBuffer saves the current replay buffer to disk.
// The replay buffer must be active for this to succeed.
func (c *Client) SaveReplayBuffer() error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Outputs.SaveReplayBuffer()
	if err != nil {
		return fmt.Errorf("failed to save replay buffer: %w. Replay buffer may not be active", err)
	}

	return nil
}

// GetLastReplayBufferReplay retrieves the path to the last saved replay.
func (c *Client) GetLastReplayBufferReplay() (string, error) {
	client, err := c.getClient()
	if err != nil {
		return "", err
	}

	resp, err := client.Outputs.GetLastReplayBufferReplay()
	if err != nil {
		return "", fmt.Errorf("failed to get last replay path: %w", err)
	}

	return resp.SavedReplayPath, nil
}

// =============================================================================
// Studio Mode Types and Methods (FB-26)
// =============================================================================

// GetStudioModeEnabled retrieves whether studio mode is enabled.
func (c *Client) GetStudioModeEnabled() (bool, error) {
	client, err := c.getClient()
	if err != nil {
		return false, err
	}

	resp, err := client.Ui.GetStudioModeEnabled()
	if err != nil {
		return false, fmt.Errorf("failed to get studio mode status: %w", err)
	}

	return resp.StudioModeEnabled, nil
}

// SetStudioModeEnabled enables or disables studio mode.
func (c *Client) SetStudioModeEnabled(enabled bool) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Ui.SetStudioModeEnabled(&ui.SetStudioModeEnabledParams{
		StudioModeEnabled: &enabled,
	})
	if err != nil {
		return fmt.Errorf("failed to set studio mode enabled=%v: %w", enabled, err)
	}

	return nil
}

// GetCurrentPreviewScene retrieves the current preview scene in studio mode.
// Returns an error if studio mode is not enabled.
func (c *Client) GetCurrentPreviewScene() (string, error) {
	client, err := c.getClient()
	if err != nil {
		return "", err
	}

	resp, err := client.Scenes.GetCurrentPreviewScene()
	if err != nil {
		if strings.Contains(err.Error(), "studio mode") {
			return "", fmt.Errorf("studio mode is not enabled. Enable studio mode in OBS to use preview scenes")
		}
		return "", fmt.Errorf("failed to get current preview scene: %w", err)
	}

	return resp.CurrentPreviewSceneName, nil
}

// SetCurrentPreviewScene sets the current preview scene in studio mode.
// Returns an error if studio mode is not enabled.
func (c *Client) SetCurrentPreviewScene(sceneName string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Scenes.SetCurrentPreviewScene(&scenes.SetCurrentPreviewSceneParams{
		SceneName: &sceneName,
	})
	if err != nil {
		if strings.Contains(err.Error(), "studio mode") {
			return fmt.Errorf("studio mode is not enabled. Enable studio mode in OBS to use preview scenes")
		}
		return fmt.Errorf("failed to set preview scene to '%s': %w", sceneName, err)
	}

	return nil
}

// =============================================================================
// Hotkey Methods (FB-26)
// =============================================================================

// TriggerHotkeyByName triggers a hotkey by its name.
// Use GetHotkeyList() to discover available hotkey names.
func (c *Client) TriggerHotkeyByName(hotkeyName string) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.General.TriggerHotkeyByName(&general.TriggerHotkeyByNameParams{
		HotkeyName: &hotkeyName,
	})
	if err != nil {
		return fmt.Errorf("failed to trigger hotkey '%s': %w", hotkeyName, err)
	}

	return nil
}

// GetHotkeyList retrieves all available hotkey names.
func (c *Client) GetHotkeyList() ([]string, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.General.GetHotkeyList()
	if err != nil {
		return nil, fmt.Errorf("failed to get hotkey list: %w", err)
	}

	return resp.Hotkeys, nil
}

// StartVirtualCam starts the virtual camera.
func (c *Client) StartVirtualCam() error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Outputs.StartVirtualCam()
	if err != nil {
		return fmt.Errorf("failed to start virtual camera: %w", err)
	}

	return nil
}

// StopVirtualCam stops the virtual camera.
func (c *Client) StopVirtualCam() error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Outputs.StopVirtualCam()
	if err != nil {
		return fmt.Errorf("failed to stop virtual camera: %w", err)
	}

	return nil
}

// StartReplayBuffer starts the replay buffer.
func (c *Client) StartReplayBuffer() error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Outputs.StartReplayBuffer()
	if err != nil {
		return fmt.Errorf("failed to start replay buffer: %w", err)
	}

	return nil
}

// StopReplayBuffer stops the replay buffer.
func (c *Client) StopReplayBuffer() error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	_, err = client.Outputs.StopReplayBuffer()
	if err != nil {
		return fmt.Errorf("failed to stop replay buffer: %w", err)
	}

	return nil
}
