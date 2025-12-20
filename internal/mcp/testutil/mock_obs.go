package testutil

import (
	"fmt"
	"sync"

	"github.com/andreykaipov/goobs/api/typedefs"
	"github.com/ironystock/agentic-obs/internal/obs"
)

// Ensure MockOBSClient implements the OBSClient interface
// This is verified at compile time in the mcp package

// MockOBSClient implements a mock OBS client for testing.
// It simulates OBS behavior without requiring a real OBS instance.
type MockOBSClient struct {
	mu sync.RWMutex

	// Connection state
	connected bool

	// Mock data
	scenes         []string
	currentScene   string
	sources        []*typedefs.Input
	sceneItems     map[string][]obs.SceneSource
	sourceSettings map[string]map[string]interface{}
	inputMutes     map[string]bool
	inputVolumes   map[string]float64

	// Recording/Streaming state
	recording bool
	paused    bool
	streaming bool

	// Error injection for testing error paths
	ErrorOnConnect             error
	ErrorOnGetSceneList        error
	ErrorOnSetCurrentScene     error
	ErrorOnCreateScene         error
	ErrorOnRemoveScene         error
	ErrorOnStartRecording      error
	ErrorOnStopRecording       error
	ErrorOnPauseRecording      error
	ErrorOnResumeRecording     error
	ErrorOnStartStreaming      error
	ErrorOnStopStreaming       error
	ErrorOnListSources         error
	ErrorOnGetSourceSettings   error
	ErrorOnToggleVisibility    error
	ErrorOnGetInputMute        error
	ErrorOnToggleInputMute     error
	ErrorOnSetInputVolume      error
	ErrorOnGetInputVolume      error
	ErrorOnGetOBSStatus        error
	ErrorOnCaptureSceneState   error
	ErrorOnApplyScenePreset    error
	ErrorOnTakeScreenshot      error
	ErrorOnCreateBrowserSource error

	// Screenshot mock data
	mockScreenshotData string // Base64 PNG data to return

	// Design tool mock data
	sceneItemTransforms map[string]map[int]*obs.SceneItemTransform // scene -> itemID -> transform
	sceneItemLocked     map[string]map[int]bool                    // scene -> itemID -> locked
	inputKinds          []string                                   // available input kinds
	nextSceneItemID     int                                        // counter for new scene items

	// Error injection for design tools
	ErrorOnCreateInput           error
	ErrorOnGetSceneItemTransform error
	ErrorOnSetSceneItemTransform error
	ErrorOnSetSceneItemIndex     error
	ErrorOnSetSceneItemLocked    error
	ErrorOnGetSceneItemLocked    error
	ErrorOnDuplicateSceneItem    error
	ErrorOnRemoveSceneItem       error
	ErrorOnGetInputKindList      error

	// Filter mock data
	sourceFilters map[string][]obs.FilterInfo              // source -> filters
	filterDetails map[string]map[string]*obs.FilterDetails // source -> filter name -> details
	filterKinds   []string                                 // available filter types

	// Transition mock data
	transitions       []obs.TransitionInfo
	currentTransition *obs.TransitionDetails
	studioModeEnabled bool

	// Error injection for filters
	ErrorOnGetSourceFilterList     error
	ErrorOnGetSourceFilter         error
	ErrorOnCreateSourceFilter      error
	ErrorOnRemoveSourceFilter      error
	ErrorOnSetSourceFilterEnabled  error
	ErrorOnSetSourceFilterSettings error
	ErrorOnGetSourceFilterKindList error

	// Error injection for transitions
	ErrorOnGetSceneTransitionList            error
	ErrorOnGetCurrentSceneTransition         error
	ErrorOnSetCurrentSceneTransition         error
	ErrorOnSetCurrentSceneTransitionDuration error
	ErrorOnTriggerStudioModeTransition       error
}

// NewMockOBSClient creates a new mock OBS client with default test data.
func NewMockOBSClient() *MockOBSClient {
	return &MockOBSClient{
		connected:    false,
		scenes:       []string{"Scene 1", "Scene 2", "Gaming", "Starting Soon"},
		currentScene: "Scene 1",
		sources: []*typedefs.Input{
			{InputName: "Microphone", InputKind: "wasapi_input_capture"},
			{InputName: "Desktop Audio", InputKind: "wasapi_output_capture"},
			{InputName: "Webcam", InputKind: "dshow_input"},
		},
		sceneItems: map[string][]obs.SceneSource{
			"Scene 1": {
				{ID: 1, Name: "Webcam", Type: "dshow_input", Enabled: true, Visible: true},
				{ID: 2, Name: "Text", Type: "text_gdiplus_v2", Enabled: true, Visible: true},
			},
			"Gaming": {
				{ID: 3, Name: "Game Capture", Type: "game_capture", Enabled: true, Visible: true},
				{ID: 4, Name: "Webcam", Type: "dshow_input", Enabled: true, Visible: true},
			},
		},
		sourceSettings: map[string]map[string]interface{}{
			"Microphone":    {"device_id": "default", "use_device_timing": true},
			"Desktop Audio": {"device_id": "default"},
			"Webcam":        {"video_device_id": "default", "resolution": "1920x1080"},
		},
		inputMutes: map[string]bool{
			"Microphone":    false,
			"Desktop Audio": false,
		},
		inputVolumes: map[string]float64{
			"Microphone":    0.0,
			"Desktop Audio": 0.0,
		},
		recording: false,
		paused:    false,
		streaming: false,
		// Design tool mock data
		sceneItemTransforms: map[string]map[int]*obs.SceneItemTransform{
			"Scene 1": {
				1: {PositionX: 0, PositionY: 0, ScaleX: 1.0, ScaleY: 1.0, Rotation: 0, Width: 1920, Height: 1080},
				2: {PositionX: 100, PositionY: 50, ScaleX: 1.0, ScaleY: 1.0, Rotation: 0, Width: 400, Height: 100},
			},
			"Gaming": {
				3: {PositionX: 0, PositionY: 0, ScaleX: 1.0, ScaleY: 1.0, Rotation: 0, Width: 1920, Height: 1080},
				4: {PositionX: 1600, PositionY: 800, ScaleX: 0.25, ScaleY: 0.25, Rotation: 0, Width: 320, Height: 180},
			},
		},
		sceneItemLocked: map[string]map[int]bool{
			"Scene 1": {1: false, 2: false},
			"Gaming":  {3: false, 4: false},
		},
		inputKinds: []string{
			"text_gdiplus_v3", "image_source", "color_source_v3", "browser_source",
			"ffmpeg_source", "wasapi_input_capture", "wasapi_output_capture",
			"dshow_input", "game_capture", "window_capture", "monitor_capture",
		},
		nextSceneItemID: 100,
		// Filter mock data
		sourceFilters: map[string][]obs.FilterInfo{
			"Webcam": {
				{Name: "Color Correction", Kind: "color_filter_v2", Index: 0, Enabled: true},
				{Name: "Sharpen", Kind: "sharpness_filter_v2", Index: 1, Enabled: true},
			},
			"Microphone": {
				{Name: "Noise Suppression", Kind: "noise_suppress_filter_v2", Index: 0, Enabled: true},
				{Name: "Compressor", Kind: "compressor_filter", Index: 1, Enabled: false},
			},
		},
		filterDetails: map[string]map[string]*obs.FilterDetails{
			"Webcam": {
				"Color Correction": {
					Name:     "Color Correction",
					Kind:     "color_filter_v2",
					Index:    0,
					Enabled:  true,
					Settings: map[string]interface{}{"brightness": 0.0, "contrast": 0.0, "saturation": 0.0},
				},
				"Sharpen": {
					Name:     "Sharpen",
					Kind:     "sharpness_filter_v2",
					Index:    1,
					Enabled:  true,
					Settings: map[string]interface{}{"sharpness": 0.08},
				},
			},
			"Microphone": {
				"Noise Suppression": {
					Name:     "Noise Suppression",
					Kind:     "noise_suppress_filter_v2",
					Index:    0,
					Enabled:  true,
					Settings: map[string]interface{}{"suppress_level": -30, "method": "rnnoise"},
				},
				"Compressor": {
					Name:     "Compressor",
					Kind:     "compressor_filter",
					Index:    1,
					Enabled:  false,
					Settings: map[string]interface{}{"ratio": 10.0, "threshold": -18.0},
				},
			},
		},
		filterKinds: []string{
			"color_filter_v2", "sharpness_filter_v2", "noise_suppress_filter_v2",
			"compressor_filter", "limiter_filter", "gain_filter", "chroma_key_filter_v2",
			"luma_key_filter", "mask_filter_v2", "scroll_filter", "crop_filter",
		},
		// Transition mock data
		transitions: []obs.TransitionInfo{
			{Name: "Cut", Kind: "cut_transition", Fixed: true, Configurable: false},
			{Name: "Fade", Kind: "fade_transition", Fixed: false, Configurable: true},
			{Name: "Swipe", Kind: "swipe_transition", Fixed: false, Configurable: true},
			{Name: "Slide", Kind: "slide_transition", Fixed: false, Configurable: true},
			{Name: "Stinger", Kind: "obs_stinger_transition", Fixed: false, Configurable: true},
		},
		currentTransition: &obs.TransitionDetails{
			Name:         "Fade",
			Kind:         "fade_transition",
			Duration:     300,
			Configurable: true,
			Settings:     map[string]interface{}{},
		},
		studioModeEnabled: false,
	}
}

// Connect simulates connecting to OBS.
func (m *MockOBSClient) Connect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnConnect != nil {
		return m.ErrorOnConnect
	}

	m.connected = true
	return nil
}

// Disconnect simulates disconnecting from OBS.
func (m *MockOBSClient) Disconnect() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connected = false
	return nil
}

// Close simulates closing the client.
func (m *MockOBSClient) Close() error {
	return m.Disconnect()
}

// IsConnected returns the connection state.
func (m *MockOBSClient) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connected
}

// GetConnectionStatus returns mock connection status.
func (m *MockOBSClient) GetConnectionStatus() (obs.ConnectionStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return obs.ConnectionStatus{
		Connected:        m.connected,
		Host:             "localhost",
		Port:             "4455",
		OBSVersion:       "30.0.0",
		WebSocketVersion: "5.4.0",
		Platform:         "windows",
	}, nil
}

// HealthCheck simulates a health check.
func (m *MockOBSClient) HealthCheck() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}
	return nil
}

// GetSceneList returns mock scene list.
func (m *MockOBSClient) GetSceneList() ([]string, string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnGetSceneList != nil {
		return nil, "", m.ErrorOnGetSceneList
	}

	if !m.connected {
		return nil, "", fmt.Errorf("not connected to OBS")
	}

	return m.scenes, m.currentScene, nil
}

// GetSceneByName returns mock scene data.
func (m *MockOBSClient) GetSceneByName(name string) (*obs.Scene, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected {
		return nil, fmt.Errorf("not connected to OBS")
	}

	// Check if scene exists
	found := false
	for _, s := range m.scenes {
		if s == name {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("scene '%s' not found", name)
	}

	sources := m.sceneItems[name]
	if sources == nil {
		sources = []obs.SceneSource{}
	}

	return &obs.Scene{
		Name:    name,
		Index:   0,
		Sources: sources,
	}, nil
}

// SetCurrentScene simulates switching scenes.
func (m *MockOBSClient) SetCurrentScene(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnSetCurrentScene != nil {
		return m.ErrorOnSetCurrentScene
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	// Check if scene exists
	found := false
	for _, s := range m.scenes {
		if s == name {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("scene '%s' not found", name)
	}

	m.currentScene = name
	return nil
}

// CreateScene simulates creating a scene.
func (m *MockOBSClient) CreateScene(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnCreateScene != nil {
		return m.ErrorOnCreateScene
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	// Check if scene already exists
	for _, s := range m.scenes {
		if s == name {
			return fmt.Errorf("scene '%s' already exists", name)
		}
	}

	m.scenes = append(m.scenes, name)
	m.sceneItems[name] = []obs.SceneSource{}
	return nil
}

// RemoveScene simulates removing a scene.
func (m *MockOBSClient) RemoveScene(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnRemoveScene != nil {
		return m.ErrorOnRemoveScene
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	// Find and remove scene
	idx := -1
	for i, s := range m.scenes {
		if s == name {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("scene '%s' not found", name)
	}

	m.scenes = append(m.scenes[:idx], m.scenes[idx+1:]...)
	delete(m.sceneItems, name)
	return nil
}

// StartRecording simulates starting recording.
func (m *MockOBSClient) StartRecording() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnStartRecording != nil {
		return m.ErrorOnStartRecording
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	if m.recording {
		return fmt.Errorf("recording already active")
	}

	m.recording = true
	m.paused = false
	return nil
}

// StopRecording simulates stopping recording.
func (m *MockOBSClient) StopRecording() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnStopRecording != nil {
		return "", m.ErrorOnStopRecording
	}

	if !m.connected {
		return "", fmt.Errorf("not connected to OBS")
	}

	if !m.recording {
		return "", fmt.Errorf("recording not active")
	}

	m.recording = false
	m.paused = false
	return "/recordings/test-recording.mkv", nil
}

// GetRecordingStatus returns mock recording status.
func (m *MockOBSClient) GetRecordingStatus() (*obs.RecordingStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected {
		return nil, fmt.Errorf("not connected to OBS")
	}

	return &obs.RecordingStatus{
		Active:      m.recording,
		Paused:      m.paused,
		Timecode:    "00:01:30.000",
		OutputPath:  "/recordings/",
		OutputBytes: 1024000,
	}, nil
}

// PauseRecording simulates pausing recording.
func (m *MockOBSClient) PauseRecording() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnPauseRecording != nil {
		return m.ErrorOnPauseRecording
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	if !m.recording {
		return fmt.Errorf("recording not active")
	}

	if m.paused {
		return fmt.Errorf("recording already paused")
	}

	m.paused = true
	return nil
}

// ResumeRecording simulates resuming recording.
func (m *MockOBSClient) ResumeRecording() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnResumeRecording != nil {
		return m.ErrorOnResumeRecording
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	if !m.recording {
		return fmt.Errorf("recording not active")
	}

	if !m.paused {
		return fmt.Errorf("recording not paused")
	}

	m.paused = false
	return nil
}

// StartStreaming simulates starting streaming.
func (m *MockOBSClient) StartStreaming() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnStartStreaming != nil {
		return m.ErrorOnStartStreaming
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	if m.streaming {
		return fmt.Errorf("streaming already active")
	}

	m.streaming = true
	return nil
}

// StopStreaming simulates stopping streaming.
func (m *MockOBSClient) StopStreaming() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnStopStreaming != nil {
		return m.ErrorOnStopStreaming
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	if !m.streaming {
		return fmt.Errorf("streaming not active")
	}

	m.streaming = false
	return nil
}

// GetStreamingStatus returns mock streaming status.
func (m *MockOBSClient) GetStreamingStatus() (*obs.StreamingStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.connected {
		return nil, fmt.Errorf("not connected to OBS")
	}

	return &obs.StreamingStatus{
		Active:       m.streaming,
		Reconnecting: false,
		Timecode:     "00:05:00.000",
		TotalBytes:   5120000,
		TotalFrames:  9000,
	}, nil
}

// ListSources returns mock source list.
func (m *MockOBSClient) ListSources() ([]*typedefs.Input, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnListSources != nil {
		return nil, m.ErrorOnListSources
	}

	if !m.connected {
		return nil, fmt.Errorf("not connected to OBS")
	}

	return m.sources, nil
}

// GetSourceSettings returns mock source settings.
func (m *MockOBSClient) GetSourceSettings(sourceName string) (map[string]interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnGetSourceSettings != nil {
		return nil, m.ErrorOnGetSourceSettings
	}

	if !m.connected {
		return nil, fmt.Errorf("not connected to OBS")
	}

	settings, exists := m.sourceSettings[sourceName]
	if !exists {
		return nil, fmt.Errorf("source '%s' not found", sourceName)
	}

	return settings, nil
}

// ToggleSourceVisibility simulates toggling source visibility.
func (m *MockOBSClient) ToggleSourceVisibility(sceneName string, sourceID int) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnToggleVisibility != nil {
		return false, m.ErrorOnToggleVisibility
	}

	if !m.connected {
		return false, fmt.Errorf("not connected to OBS")
	}

	items, exists := m.sceneItems[sceneName]
	if !exists {
		return false, fmt.Errorf("scene '%s' not found", sceneName)
	}

	for i, item := range items {
		if item.ID == sourceID {
			m.sceneItems[sceneName][i].Enabled = !item.Enabled
			return m.sceneItems[sceneName][i].Enabled, nil
		}
	}

	return false, fmt.Errorf("source ID %d not found in scene '%s'", sourceID, sceneName)
}

// GetInputMute returns mock input mute state.
func (m *MockOBSClient) GetInputMute(inputName string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnGetInputMute != nil {
		return false, m.ErrorOnGetInputMute
	}

	if !m.connected {
		return false, fmt.Errorf("not connected to OBS")
	}

	muted, exists := m.inputMutes[inputName]
	if !exists {
		return false, fmt.Errorf("input '%s' not found", inputName)
	}

	return muted, nil
}

// ToggleInputMute simulates toggling input mute.
func (m *MockOBSClient) ToggleInputMute(inputName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnToggleInputMute != nil {
		return m.ErrorOnToggleInputMute
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	_, exists := m.inputMutes[inputName]
	if !exists {
		return fmt.Errorf("input '%s' not found", inputName)
	}

	m.inputMutes[inputName] = !m.inputMutes[inputName]
	return nil
}

// SetInputVolume simulates setting input volume.
func (m *MockOBSClient) SetInputVolume(inputName string, volumeDb *float64, volumeMul *float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnSetInputVolume != nil {
		return m.ErrorOnSetInputVolume
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	_, exists := m.inputVolumes[inputName]
	if !exists {
		return fmt.Errorf("input '%s' not found", inputName)
	}

	if volumeDb != nil {
		m.inputVolumes[inputName] = *volumeDb
	}

	return nil
}

// GetInputVolume returns mock input volume.
func (m *MockOBSClient) GetInputVolume(inputName string) (float64, float64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnGetInputVolume != nil {
		return 0, 0, m.ErrorOnGetInputVolume
	}

	if !m.connected {
		return 0, 0, fmt.Errorf("not connected to OBS")
	}

	vol, exists := m.inputVolumes[inputName]
	if !exists {
		return 0, 0, fmt.Errorf("input '%s' not found", inputName)
	}

	return vol, 1.0, nil
}

// GetOBSStatus returns mock OBS status.
func (m *MockOBSClient) GetOBSStatus() (*obs.OBSStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnGetOBSStatus != nil {
		return nil, m.ErrorOnGetOBSStatus
	}

	if !m.connected {
		return nil, fmt.Errorf("not connected to OBS")
	}

	return &obs.OBSStatus{
		Version:          "30.0.0",
		WebSocketVersion: "5.4.0",
		Platform:         "windows",
		CurrentScene:     m.currentScene,
		Recording:        m.recording,
		Streaming:        m.streaming,
		FPS:              60.0,
		FrameTime:        16.67,
		Frames:           10000,
		DroppedFrames:    5,
	}, nil
}

// Helper methods for test setup

// SetScenes sets the available scenes for testing.
func (m *MockOBSClient) SetScenes(scenes []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.scenes = scenes
}

// SetCurrentSceneDirect sets the current scene without validation.
func (m *MockOBSClient) SetCurrentSceneDirect(scene string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentScene = scene
}

// SetRecordingState sets the recording state for testing.
func (m *MockOBSClient) SetRecordingState(recording, paused bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recording = recording
	m.paused = paused
}

// SetStreamingState sets the streaming state for testing.
func (m *MockOBSClient) SetStreamingState(streaming bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.streaming = streaming
}

// AddSource adds a source to the mock data.
func (m *MockOBSClient) AddSource(input *typedefs.Input) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sources = append(m.sources, input)
}

// SetSourceSettings sets the settings for a source.
func (m *MockOBSClient) SetSourceSettings(sourceName string, settings map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sourceSettings[sourceName] = settings
}

// SetInputMuteState sets the mute state for an input.
func (m *MockOBSClient) SetInputMuteState(inputName string, muted bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inputMutes[inputName] = muted
}

// SetInputVolumeState sets the volume for an input.
func (m *MockOBSClient) SetInputVolumeState(inputName string, volume float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inputVolumes[inputName] = volume
}

// SetEventCallback is a no-op for the mock client since we don't need event handling in tests.
func (m *MockOBSClient) SetEventCallback(callback obs.EventCallback) {
	// No-op for mock - events are not simulated
}

// CaptureSceneState returns the current source states for a scene.
func (m *MockOBSClient) CaptureSceneState(sceneName string) ([]obs.SourceState, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnCaptureSceneState != nil {
		return nil, m.ErrorOnCaptureSceneState
	}

	if !m.connected {
		return nil, fmt.Errorf("not connected to OBS")
	}

	items, exists := m.sceneItems[sceneName]
	if !exists {
		return nil, fmt.Errorf("scene '%s' not found", sceneName)
	}

	states := make([]obs.SourceState, len(items))
	for i, item := range items {
		states[i] = obs.SourceState{
			ID:      item.ID,
			Name:    item.Name,
			Enabled: item.Enabled,
		}
	}

	return states, nil
}

// ApplyScenePreset applies source visibility states to a scene.
func (m *MockOBSClient) ApplyScenePreset(sceneName string, sources []obs.SourceState) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnApplyScenePreset != nil {
		return m.ErrorOnApplyScenePreset
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	items, exists := m.sceneItems[sceneName]
	if !exists {
		return fmt.Errorf("scene '%s' not found", sceneName)
	}

	// Apply each source state
	for _, src := range sources {
		found := false
		for i, item := range items {
			if item.ID == src.ID {
				m.sceneItems[sceneName][i].Enabled = src.Enabled
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("source ID %d not found in scene '%s'", src.ID, sceneName)
		}
	}

	return nil
}

// TakeSourceScreenshot simulates taking a screenshot of a source.
func (m *MockOBSClient) TakeSourceScreenshot(opts obs.ScreenshotOptions) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnTakeScreenshot != nil {
		return "", m.ErrorOnTakeScreenshot
	}

	if !m.connected {
		return "", fmt.Errorf("not connected to OBS")
	}

	// Return mock screenshot data if set, otherwise return a minimal valid base64 PNG
	if m.mockScreenshotData != "" {
		return m.mockScreenshotData, nil
	}

	// Return a minimal 1x1 transparent PNG as base64
	// This is a valid PNG that can be used in tests
	return "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==", nil
}

// CreateBrowserSource simulates creating a browser source in a scene.
func (m *MockOBSClient) CreateBrowserSource(sceneName, sourceName string, settings obs.BrowserSourceSettings) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnCreateBrowserSource != nil {
		return 0, m.ErrorOnCreateBrowserSource
	}

	if !m.connected {
		return 0, fmt.Errorf("not connected to OBS")
	}

	// Check if scene exists
	found := false
	for _, s := range m.scenes {
		if s == sceneName {
			found = true
			break
		}
	}
	if !found {
		return 0, fmt.Errorf("scene '%s' not found", sceneName)
	}

	// Generate a mock scene item ID
	newID := 100 // Start from 100 to avoid conflicts with existing mock IDs
	if items, exists := m.sceneItems[sceneName]; exists {
		for _, item := range items {
			if item.ID >= newID {
				newID = item.ID + 1
			}
		}
	}

	// Add the browser source to the scene
	newSource := obs.SceneSource{
		ID:      newID,
		Name:    sourceName,
		Type:    "browser_source",
		Enabled: true,
		Visible: true,
	}

	if m.sceneItems == nil {
		m.sceneItems = make(map[string][]obs.SceneSource)
	}
	m.sceneItems[sceneName] = append(m.sceneItems[sceneName], newSource)

	// Also add to sourceSettings for consistency
	if m.sourceSettings == nil {
		m.sourceSettings = make(map[string]map[string]interface{})
	}
	m.sourceSettings[sourceName] = map[string]interface{}{
		"url":    settings.URL,
		"width":  settings.Width,
		"height": settings.Height,
		"css":    settings.CSS,
	}

	return newID, nil
}

// SetMockScreenshotData sets the screenshot data to return from TakeSourceScreenshot.
func (m *MockOBSClient) SetMockScreenshotData(data string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mockScreenshotData = data
}

// Design tool mock implementations

// CreateInput simulates creating an input source in a scene.
func (m *MockOBSClient) CreateInput(sceneName, sourceName, inputKind string, settings map[string]interface{}) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnCreateInput != nil {
		return 0, m.ErrorOnCreateInput
	}

	if !m.connected {
		return 0, fmt.Errorf("not connected to OBS")
	}

	// Check if scene exists
	found := false
	for _, s := range m.scenes {
		if s == sceneName {
			found = true
			break
		}
	}
	if !found {
		return 0, fmt.Errorf("scene '%s' not found", sceneName)
	}

	// Generate new scene item ID
	m.nextSceneItemID++
	newID := m.nextSceneItemID

	// Add the source to the scene
	newSource := obs.SceneSource{
		ID:      newID,
		Name:    sourceName,
		Type:    inputKind,
		Enabled: true,
		Visible: true,
	}

	if m.sceneItems == nil {
		m.sceneItems = make(map[string][]obs.SceneSource)
	}
	m.sceneItems[sceneName] = append(m.sceneItems[sceneName], newSource)

	// Initialize transform for the new item
	if m.sceneItemTransforms == nil {
		m.sceneItemTransforms = make(map[string]map[int]*obs.SceneItemTransform)
	}
	if m.sceneItemTransforms[sceneName] == nil {
		m.sceneItemTransforms[sceneName] = make(map[int]*obs.SceneItemTransform)
	}
	m.sceneItemTransforms[sceneName][newID] = &obs.SceneItemTransform{
		PositionX: 0, PositionY: 0,
		ScaleX: 1.0, ScaleY: 1.0,
		Rotation: 0,
		Width:    1920, Height: 1080,
	}

	// Initialize locked state
	if m.sceneItemLocked == nil {
		m.sceneItemLocked = make(map[string]map[int]bool)
	}
	if m.sceneItemLocked[sceneName] == nil {
		m.sceneItemLocked[sceneName] = make(map[int]bool)
	}
	m.sceneItemLocked[sceneName][newID] = false

	// Store source settings
	if m.sourceSettings == nil {
		m.sourceSettings = make(map[string]map[string]interface{})
	}
	m.sourceSettings[sourceName] = settings

	return newID, nil
}

// GetSceneItemTransform returns the transform for a scene item.
func (m *MockOBSClient) GetSceneItemTransform(sceneName string, sceneItemID int) (*obs.SceneItemTransform, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnGetSceneItemTransform != nil {
		return nil, m.ErrorOnGetSceneItemTransform
	}

	if !m.connected {
		return nil, fmt.Errorf("not connected to OBS")
	}

	sceneTransforms, exists := m.sceneItemTransforms[sceneName]
	if !exists {
		return nil, fmt.Errorf("scene '%s' not found", sceneName)
	}

	transform, exists := sceneTransforms[sceneItemID]
	if !exists {
		return nil, fmt.Errorf("scene item %d not found in scene '%s'", sceneItemID, sceneName)
	}

	// Return a copy to prevent modification
	result := *transform
	return &result, nil
}

// SetSceneItemTransform sets the transform for a scene item.
func (m *MockOBSClient) SetSceneItemTransform(sceneName string, sceneItemID int, transform *obs.SceneItemTransform) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnSetSceneItemTransform != nil {
		return m.ErrorOnSetSceneItemTransform
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	sceneTransforms, exists := m.sceneItemTransforms[sceneName]
	if !exists {
		return fmt.Errorf("scene '%s' not found", sceneName)
	}

	if _, exists := sceneTransforms[sceneItemID]; !exists {
		return fmt.Errorf("scene item %d not found in scene '%s'", sceneItemID, sceneName)
	}

	// Store a copy
	newTransform := *transform
	m.sceneItemTransforms[sceneName][sceneItemID] = &newTransform
	return nil
}

// SetSceneItemIndex sets the z-order index of a scene item.
func (m *MockOBSClient) SetSceneItemIndex(sceneName string, sceneItemID int, index int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnSetSceneItemIndex != nil {
		return m.ErrorOnSetSceneItemIndex
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	items, exists := m.sceneItems[sceneName]
	if !exists {
		return fmt.Errorf("scene '%s' not found", sceneName)
	}

	// Find the item
	found := false
	for _, item := range items {
		if item.ID == sceneItemID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("scene item %d not found in scene '%s'", sceneItemID, sceneName)
	}

	// In a real implementation, we'd reorder the items
	// For mock, just validate and succeed
	if index < 0 || index >= len(items) {
		return fmt.Errorf("invalid index %d for scene with %d items", index, len(items))
	}

	return nil
}

// SetSceneItemLocked sets the locked state of a scene item.
func (m *MockOBSClient) SetSceneItemLocked(sceneName string, sceneItemID int, locked bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnSetSceneItemLocked != nil {
		return m.ErrorOnSetSceneItemLocked
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	sceneLocked, exists := m.sceneItemLocked[sceneName]
	if !exists {
		return fmt.Errorf("scene '%s' not found", sceneName)
	}

	if _, exists := sceneLocked[sceneItemID]; !exists {
		return fmt.Errorf("scene item %d not found in scene '%s'", sceneItemID, sceneName)
	}

	m.sceneItemLocked[sceneName][sceneItemID] = locked
	return nil
}

// GetSceneItemLocked returns the locked state of a scene item.
func (m *MockOBSClient) GetSceneItemLocked(sceneName string, sceneItemID int) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnGetSceneItemLocked != nil {
		return false, m.ErrorOnGetSceneItemLocked
	}

	if !m.connected {
		return false, fmt.Errorf("not connected to OBS")
	}

	sceneLocked, exists := m.sceneItemLocked[sceneName]
	if !exists {
		return false, fmt.Errorf("scene '%s' not found", sceneName)
	}

	locked, exists := sceneLocked[sceneItemID]
	if !exists {
		return false, fmt.Errorf("scene item %d not found in scene '%s'", sceneItemID, sceneName)
	}

	return locked, nil
}

// DuplicateSceneItem duplicates a scene item to the same or another scene.
func (m *MockOBSClient) DuplicateSceneItem(sceneName string, sceneItemID int, destScene string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnDuplicateSceneItem != nil {
		return 0, m.ErrorOnDuplicateSceneItem
	}

	if !m.connected {
		return 0, fmt.Errorf("not connected to OBS")
	}

	// Find source item
	items, exists := m.sceneItems[sceneName]
	if !exists {
		return 0, fmt.Errorf("scene '%s' not found", sceneName)
	}

	var sourceItem *obs.SceneSource
	for _, item := range items {
		if item.ID == sceneItemID {
			sourceItem = &item
			break
		}
	}
	if sourceItem == nil {
		return 0, fmt.Errorf("scene item %d not found in scene '%s'", sceneItemID, sceneName)
	}

	// Check destination scene exists
	found := false
	for _, s := range m.scenes {
		if s == destScene {
			found = true
			break
		}
	}
	if !found {
		return 0, fmt.Errorf("destination scene '%s' not found", destScene)
	}

	// Generate new ID
	m.nextSceneItemID++
	newID := m.nextSceneItemID

	// Create duplicate
	newSource := obs.SceneSource{
		ID:      newID,
		Name:    sourceItem.Name,
		Type:    sourceItem.Type,
		Enabled: sourceItem.Enabled,
		Visible: sourceItem.Visible,
	}
	m.sceneItems[destScene] = append(m.sceneItems[destScene], newSource)

	// Copy transform
	if m.sceneItemTransforms[destScene] == nil {
		m.sceneItemTransforms[destScene] = make(map[int]*obs.SceneItemTransform)
	}
	if srcTransform, ok := m.sceneItemTransforms[sceneName][sceneItemID]; ok {
		transformCopy := *srcTransform
		m.sceneItemTransforms[destScene][newID] = &transformCopy
	}

	// Copy locked state
	if m.sceneItemLocked[destScene] == nil {
		m.sceneItemLocked[destScene] = make(map[int]bool)
	}
	if srcLocked, ok := m.sceneItemLocked[sceneName][sceneItemID]; ok {
		m.sceneItemLocked[destScene][newID] = srcLocked
	}

	return newID, nil
}

// RemoveSceneItem removes a scene item from a scene.
func (m *MockOBSClient) RemoveSceneItem(sceneName string, sceneItemID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnRemoveSceneItem != nil {
		return m.ErrorOnRemoveSceneItem
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	items, exists := m.sceneItems[sceneName]
	if !exists {
		return fmt.Errorf("scene '%s' not found", sceneName)
	}

	// Find and remove the item
	idx := -1
	for i, item := range items {
		if item.ID == sceneItemID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("scene item %d not found in scene '%s'", sceneItemID, sceneName)
	}

	m.sceneItems[sceneName] = append(items[:idx], items[idx+1:]...)

	// Clean up transform and locked state
	if sceneTransforms, ok := m.sceneItemTransforms[sceneName]; ok {
		delete(sceneTransforms, sceneItemID)
	}
	if sceneLocked, ok := m.sceneItemLocked[sceneName]; ok {
		delete(sceneLocked, sceneItemID)
	}

	return nil
}

// GetInputKindList returns the list of available input kinds.
func (m *MockOBSClient) GetInputKindList() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnGetInputKindList != nil {
		return nil, m.ErrorOnGetInputKindList
	}

	if !m.connected {
		return nil, fmt.Errorf("not connected to OBS")
	}

	// Return a copy
	result := make([]string, len(m.inputKinds))
	copy(result, m.inputKinds)
	return result, nil
}

// =============================================================================
// Filter mock implementations (FB-23)
// =============================================================================

// GetSourceFilterList returns filters for a source.
func (m *MockOBSClient) GetSourceFilterList(sourceName string) ([]obs.FilterInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnGetSourceFilterList != nil {
		return nil, m.ErrorOnGetSourceFilterList
	}

	if !m.connected {
		return nil, fmt.Errorf("not connected to OBS")
	}

	filters, exists := m.sourceFilters[sourceName]
	if !exists {
		// Return empty list for sources without filters
		return []obs.FilterInfo{}, nil
	}

	// Return a copy
	result := make([]obs.FilterInfo, len(filters))
	copy(result, filters)
	return result, nil
}

// GetSourceFilter returns details for a specific filter.
func (m *MockOBSClient) GetSourceFilter(sourceName, filterName string) (*obs.FilterDetails, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnGetSourceFilter != nil {
		return nil, m.ErrorOnGetSourceFilter
	}

	if !m.connected {
		return nil, fmt.Errorf("not connected to OBS")
	}

	sourceFilters, exists := m.filterDetails[sourceName]
	if !exists {
		return nil, fmt.Errorf("source '%s' not found", sourceName)
	}

	filter, exists := sourceFilters[filterName]
	if !exists {
		return nil, fmt.Errorf("filter '%s' not found on source '%s'", filterName, sourceName)
	}

	// Return a copy
	result := *filter
	return &result, nil
}

// CreateSourceFilter creates a filter on a source.
func (m *MockOBSClient) CreateSourceFilter(sourceName, filterName, filterKind string, settings map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnCreateSourceFilter != nil {
		return m.ErrorOnCreateSourceFilter
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	// Check if filter already exists
	if sourceFilters, exists := m.filterDetails[sourceName]; exists {
		if _, exists := sourceFilters[filterName]; exists {
			return fmt.Errorf("filter '%s' already exists on source '%s'", filterName, sourceName)
		}
	}

	// Initialize maps if needed
	if m.sourceFilters == nil {
		m.sourceFilters = make(map[string][]obs.FilterInfo)
	}
	if m.filterDetails == nil {
		m.filterDetails = make(map[string]map[string]*obs.FilterDetails)
	}
	if m.filterDetails[sourceName] == nil {
		m.filterDetails[sourceName] = make(map[string]*obs.FilterDetails)
	}

	// Determine index
	index := len(m.sourceFilters[sourceName])

	// Add to filter list
	m.sourceFilters[sourceName] = append(m.sourceFilters[sourceName], obs.FilterInfo{
		Name:    filterName,
		Kind:    filterKind,
		Index:   index,
		Enabled: true,
	})

	// Add to filter details
	m.filterDetails[sourceName][filterName] = &obs.FilterDetails{
		Name:     filterName,
		Kind:     filterKind,
		Index:    index,
		Enabled:  true,
		Settings: settings,
	}

	return nil
}

// RemoveSourceFilter removes a filter from a source.
func (m *MockOBSClient) RemoveSourceFilter(sourceName, filterName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnRemoveSourceFilter != nil {
		return m.ErrorOnRemoveSourceFilter
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	// Check if filter exists
	sourceFilters, exists := m.filterDetails[sourceName]
	if !exists {
		return fmt.Errorf("source '%s' not found", sourceName)
	}

	if _, exists := sourceFilters[filterName]; !exists {
		return fmt.Errorf("filter '%s' not found on source '%s'", filterName, sourceName)
	}

	// Remove from filter details
	delete(m.filterDetails[sourceName], filterName)

	// Remove from filter list
	filters := m.sourceFilters[sourceName]
	for i, f := range filters {
		if f.Name == filterName {
			m.sourceFilters[sourceName] = append(filters[:i], filters[i+1:]...)
			break
		}
	}

	// Update indices
	for i := range m.sourceFilters[sourceName] {
		m.sourceFilters[sourceName][i].Index = i
		if details, ok := m.filterDetails[sourceName][m.sourceFilters[sourceName][i].Name]; ok {
			details.Index = i
		}
	}

	return nil
}

// SetSourceFilterEnabled enables or disables a filter.
func (m *MockOBSClient) SetSourceFilterEnabled(sourceName, filterName string, enabled bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnSetSourceFilterEnabled != nil {
		return m.ErrorOnSetSourceFilterEnabled
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	sourceFilters, exists := m.filterDetails[sourceName]
	if !exists {
		return fmt.Errorf("source '%s' not found", sourceName)
	}

	filter, exists := sourceFilters[filterName]
	if !exists {
		return fmt.Errorf("filter '%s' not found on source '%s'", filterName, sourceName)
	}

	filter.Enabled = enabled

	// Update in filter list too
	for i, f := range m.sourceFilters[sourceName] {
		if f.Name == filterName {
			m.sourceFilters[sourceName][i].Enabled = enabled
			break
		}
	}

	return nil
}

// SetSourceFilterSettings updates filter settings.
func (m *MockOBSClient) SetSourceFilterSettings(sourceName, filterName string, settings map[string]interface{}, overlay bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnSetSourceFilterSettings != nil {
		return m.ErrorOnSetSourceFilterSettings
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	sourceFilters, exists := m.filterDetails[sourceName]
	if !exists {
		return fmt.Errorf("source '%s' not found", sourceName)
	}

	filter, exists := sourceFilters[filterName]
	if !exists {
		return fmt.Errorf("filter '%s' not found on source '%s'", filterName, sourceName)
	}

	if overlay {
		// Merge settings
		for k, v := range settings {
			filter.Settings[k] = v
		}
	} else {
		// Replace settings
		filter.Settings = settings
	}

	return nil
}

// GetSourceFilterKindList returns available filter types.
func (m *MockOBSClient) GetSourceFilterKindList() ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnGetSourceFilterKindList != nil {
		return nil, m.ErrorOnGetSourceFilterKindList
	}

	if !m.connected {
		return nil, fmt.Errorf("not connected to OBS")
	}

	// Return a copy
	result := make([]string, len(m.filterKinds))
	copy(result, m.filterKinds)
	return result, nil
}

// =============================================================================
// Transition mock implementations (FB-24)
// =============================================================================

// GetSceneTransitionList returns available transitions and current transition name.
func (m *MockOBSClient) GetSceneTransitionList() ([]obs.TransitionInfo, string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnGetSceneTransitionList != nil {
		return nil, "", m.ErrorOnGetSceneTransitionList
	}

	if !m.connected {
		return nil, "", fmt.Errorf("not connected to OBS")
	}

	// Return a copy
	result := make([]obs.TransitionInfo, len(m.transitions))
	copy(result, m.transitions)

	currentName := ""
	if m.currentTransition != nil {
		currentName = m.currentTransition.Name
	}

	return result, currentName, nil
}

// GetCurrentSceneTransition returns the current transition details.
func (m *MockOBSClient) GetCurrentSceneTransition() (*obs.TransitionDetails, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOnGetCurrentSceneTransition != nil {
		return nil, m.ErrorOnGetCurrentSceneTransition
	}

	if !m.connected {
		return nil, fmt.Errorf("not connected to OBS")
	}

	if m.currentTransition == nil {
		return nil, fmt.Errorf("no current transition set")
	}

	// Return a copy
	result := *m.currentTransition
	return &result, nil
}

// SetCurrentSceneTransition sets the current transition.
func (m *MockOBSClient) SetCurrentSceneTransition(transitionName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnSetCurrentSceneTransition != nil {
		return m.ErrorOnSetCurrentSceneTransition
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	// Find the transition
	var found *obs.TransitionInfo
	for _, t := range m.transitions {
		if t.Name == transitionName {
			found = &t
			break
		}
	}

	if found == nil {
		return fmt.Errorf("transition '%s' not found", transitionName)
	}

	// Update current transition
	duration := 300
	if m.currentTransition != nil {
		duration = m.currentTransition.Duration
	}

	m.currentTransition = &obs.TransitionDetails{
		Name:         found.Name,
		Kind:         found.Kind,
		Duration:     duration,
		Configurable: found.Configurable,
		Settings:     map[string]interface{}{},
	}

	return nil
}

// SetCurrentSceneTransitionDuration sets the transition duration.
func (m *MockOBSClient) SetCurrentSceneTransitionDuration(durationMs int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnSetCurrentSceneTransitionDuration != nil {
		return m.ErrorOnSetCurrentSceneTransitionDuration
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	if m.currentTransition == nil {
		return fmt.Errorf("no current transition set")
	}

	m.currentTransition.Duration = durationMs
	return nil
}

// TriggerStudioModeTransition triggers the studio mode transition.
func (m *MockOBSClient) TriggerStudioModeTransition() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOnTriggerStudioModeTransition != nil {
		return m.ErrorOnTriggerStudioModeTransition
	}

	if !m.connected {
		return fmt.Errorf("not connected to OBS")
	}

	if !m.studioModeEnabled {
		return fmt.Errorf("studio mode is not enabled")
	}

	// In a real implementation, this would swap preview and program scenes
	return nil
}

// Helper methods for test setup

// SetStudioModeEnabled sets the studio mode state for testing.
func (m *MockOBSClient) SetStudioModeEnabled(enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.studioModeEnabled = enabled
}

// AddFilter adds a filter to a source for testing.
func (m *MockOBSClient) AddFilter(sourceName string, filter obs.FilterDetails) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sourceFilters == nil {
		m.sourceFilters = make(map[string][]obs.FilterInfo)
	}
	if m.filterDetails == nil {
		m.filterDetails = make(map[string]map[string]*obs.FilterDetails)
	}
	if m.filterDetails[sourceName] == nil {
		m.filterDetails[sourceName] = make(map[string]*obs.FilterDetails)
	}

	m.sourceFilters[sourceName] = append(m.sourceFilters[sourceName], obs.FilterInfo{
		Name:    filter.Name,
		Kind:    filter.Kind,
		Index:   filter.Index,
		Enabled: filter.Enabled,
	})
	m.filterDetails[sourceName][filter.Name] = &filter
}

// AddTransition adds a transition for testing.
func (m *MockOBSClient) AddTransition(transition obs.TransitionInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transitions = append(m.transitions, transition)
}

// SetCurrentTransitionDirect sets the current transition directly for testing.
func (m *MockOBSClient) SetCurrentTransitionDirect(transition *obs.TransitionDetails) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentTransition = transition
}
