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
