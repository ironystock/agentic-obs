package mcp

import (
	"github.com/andreykaipov/goobs/api/typedefs"
	"github.com/ironystock/agentic-obs/internal/obs"
)

// OBSClient defines the interface for OBS client operations.
// This interface allows the Server to use either a real OBS client or a mock for testing.
type OBSClient interface {
	// Connection management
	Connect() error
	Disconnect() error
	Close() error
	IsConnected() bool
	GetConnectionStatus() (obs.ConnectionStatus, error)
	HealthCheck() error

	// Scene operations
	GetSceneList() ([]string, string, error)
	GetSceneByName(name string) (*obs.Scene, error)
	SetCurrentScene(name string) error
	CreateScene(name string) error
	RemoveScene(name string) error

	// Recording operations
	StartRecording() error
	StopRecording() (string, error)
	GetRecordingStatus() (*obs.RecordingStatus, error)
	PauseRecording() error
	ResumeRecording() error

	// Streaming operations
	StartStreaming() error
	StopStreaming() error
	GetStreamingStatus() (*obs.StreamingStatus, error)

	// Source operations
	ListSources() ([]*typedefs.Input, error)
	GetSourceSettings(sourceName string) (map[string]interface{}, error)
	ToggleSourceVisibility(sceneName string, sourceID int) (bool, error)

	// Audio operations
	GetInputMute(inputName string) (bool, error)
	ToggleInputMute(inputName string) error
	SetInputVolume(inputName string, volumeDb *float64, volumeMul *float64) error
	GetInputVolume(inputName string) (float64, float64, error)

	// Status
	GetOBSStatus() (*obs.OBSStatus, error)

	// Scene preset operations
	CaptureSceneState(sceneName string) ([]obs.SourceState, error)
	ApplyScenePreset(sceneName string, sources []obs.SourceState) error

	// Screenshot operations
	TakeSourceScreenshot(opts obs.ScreenshotOptions) (string, error)
	CreateBrowserSource(sceneName, sourceName string, settings obs.BrowserSourceSettings) (int, error)

	// Scene design - source creation
	CreateInput(sceneName, sourceName, inputKind string, settings map[string]interface{}) (int, error)

	// Scene design - transform operations
	GetSceneItemTransform(sceneName string, sceneItemID int) (*obs.SceneItemTransform, error)
	SetSceneItemTransform(sceneName string, sceneItemID int, transform *obs.SceneItemTransform) error
	SetSceneItemIndex(sceneName string, sceneItemID int, index int) error

	// Scene design - item management
	SetSceneItemLocked(sceneName string, sceneItemID int, locked bool) error
	GetSceneItemLocked(sceneName string, sceneItemID int) (bool, error)
	DuplicateSceneItem(sceneName string, sceneItemID int, destScene string) (int, error)
	RemoveSceneItem(sceneName string, sceneItemID int) error

	// Scene design - input kinds
	GetInputKindList() ([]string, error)

	// Event handling
	SetEventCallback(callback obs.EventCallback)
}

// Verify that obs.Client implements OBSClient at compile time
var _ OBSClient = (*obs.Client)(nil)
