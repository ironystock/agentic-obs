package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ironystock/agentic-obs/internal/mcp/testutil"
	"github.com/ironystock/agentic-obs/internal/obs"
)

// TestOBSClientInterface verifies that both the real client and mock implement OBSClient
func TestOBSClientInterface(t *testing.T) {
	t.Run("mock client implements OBSClient", func(t *testing.T) {
		// This is a compile-time check - if it compiles, the mock implements the interface
		var _ OBSClient = (*testutil.MockOBSClient)(nil)
	})

	t.Run("real client implements OBSClient", func(t *testing.T) {
		// This is a compile-time check - if it compiles, the real client implements the interface
		var _ OBSClient = (*obs.Client)(nil)
	})
}

// TestMockClientBehavior tests basic mock client behavior
func TestMockClientBehavior(t *testing.T) {
	t.Run("mock starts disconnected", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()
		assert.False(t, mock.IsConnected())
	})

	t.Run("mock connects and disconnects", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()

		err := mock.Connect()
		assert.NoError(t, err)
		assert.True(t, mock.IsConnected())

		err = mock.Disconnect()
		assert.NoError(t, err)
		assert.False(t, mock.IsConnected())
	})

	t.Run("mock returns connection status", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()
		mock.Connect()

		status, err := mock.GetConnectionStatus()
		assert.NoError(t, err)
		assert.True(t, status.Connected)
		assert.Equal(t, "localhost", status.Host)
		assert.Equal(t, "4455", status.Port)
	})

	t.Run("mock health check requires connection", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()

		err := mock.HealthCheck()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not connected")

		mock.Connect()
		err = mock.HealthCheck()
		assert.NoError(t, err)
	})
}

// TestMockSceneOperations tests mock scene operations
func TestMockSceneOperations(t *testing.T) {
	t.Run("mock returns default scenes", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()
		mock.Connect()

		scenes, current, err := mock.GetSceneList()
		assert.NoError(t, err)
		assert.Contains(t, scenes, "Scene 1")
		assert.Contains(t, scenes, "Gaming")
		assert.Equal(t, "Scene 1", current)
	})

	t.Run("mock gets scene by name", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()
		mock.Connect()

		scene, err := mock.GetSceneByName("Gaming")
		assert.NoError(t, err)
		assert.Equal(t, "Gaming", scene.Name)
	})

	t.Run("mock creates and removes scenes", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()
		mock.Connect()

		err := mock.CreateScene("New Test Scene")
		assert.NoError(t, err)

		scenes, _, _ := mock.GetSceneList()
		assert.Contains(t, scenes, "New Test Scene")

		err = mock.RemoveScene("New Test Scene")
		assert.NoError(t, err)

		scenes, _, _ = mock.GetSceneList()
		assert.NotContains(t, scenes, "New Test Scene")
	})
}

// TestMockRecordingOperations tests mock recording operations
func TestMockRecordingOperations(t *testing.T) {
	t.Run("mock recording lifecycle", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()
		mock.Connect()

		// Initial state - not recording
		status, _ := mock.GetRecordingStatus()
		assert.False(t, status.Active)

		// Start recording
		err := mock.StartRecording()
		assert.NoError(t, err)
		status, _ = mock.GetRecordingStatus()
		assert.True(t, status.Active)
		assert.False(t, status.Paused)

		// Pause recording
		err = mock.PauseRecording()
		assert.NoError(t, err)
		status, _ = mock.GetRecordingStatus()
		assert.True(t, status.Active)
		assert.True(t, status.Paused)

		// Resume recording
		err = mock.ResumeRecording()
		assert.NoError(t, err)
		status, _ = mock.GetRecordingStatus()
		assert.True(t, status.Active)
		assert.False(t, status.Paused)

		// Stop recording
		path, err := mock.StopRecording()
		assert.NoError(t, err)
		assert.NotEmpty(t, path)
		status, _ = mock.GetRecordingStatus()
		assert.False(t, status.Active)
	})
}

// TestMockStreamingOperations tests mock streaming operations
func TestMockStreamingOperations(t *testing.T) {
	t.Run("mock streaming lifecycle", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()
		mock.Connect()

		// Initial state - not streaming
		status, _ := mock.GetStreamingStatus()
		assert.False(t, status.Active)

		// Start streaming
		err := mock.StartStreaming()
		assert.NoError(t, err)
		status, _ = mock.GetStreamingStatus()
		assert.True(t, status.Active)

		// Stop streaming
		err = mock.StopStreaming()
		assert.NoError(t, err)
		status, _ = mock.GetStreamingStatus()
		assert.False(t, status.Active)
	})
}

// TestMockAudioOperations tests mock audio operations
func TestMockAudioOperations(t *testing.T) {
	t.Run("mock audio mute operations", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()
		mock.Connect()

		// Check initial mute state
		muted, err := mock.GetInputMute("Microphone")
		assert.NoError(t, err)
		assert.False(t, muted)

		// Toggle mute
		err = mock.ToggleInputMute("Microphone")
		assert.NoError(t, err)

		muted, _ = mock.GetInputMute("Microphone")
		assert.True(t, muted)

		// Toggle again
		err = mock.ToggleInputMute("Microphone")
		assert.NoError(t, err)

		muted, _ = mock.GetInputMute("Microphone")
		assert.False(t, muted)
	})

	t.Run("mock audio volume operations", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()
		mock.Connect()

		// Get initial volume
		volDb, volMul, err := mock.GetInputVolume("Microphone")
		assert.NoError(t, err)
		_ = volDb
		_ = volMul

		// Set volume
		newVol := -6.0
		err = mock.SetInputVolume("Microphone", &newVol, nil)
		assert.NoError(t, err)
	})
}

// TestMockErrorInjection tests error injection capabilities
func TestMockErrorInjection(t *testing.T) {
	t.Run("inject connection error", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()
		mock.ErrorOnConnect = assert.AnError

		err := mock.Connect()
		assert.Error(t, err)
	})

	t.Run("inject scene list error", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()
		mock.Connect()
		mock.ErrorOnGetSceneList = assert.AnError

		_, _, err := mock.GetSceneList()
		assert.Error(t, err)
	})

	t.Run("inject start recording error", func(t *testing.T) {
		mock := testutil.NewMockOBSClient()
		mock.Connect()
		mock.ErrorOnStartRecording = assert.AnError

		err := mock.StartRecording()
		assert.Error(t, err)
	})
}
