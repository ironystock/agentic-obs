package mcp

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ironystock/agentic-obs/internal/mcp/testutil"
	"github.com/ironystock/agentic-obs/internal/storage"
)

// testServer creates a minimal test server with a mock OBS client.
func testServer(t *testing.T) (*Server, *testutil.MockOBSClient) {
	t.Helper()

	mock := testutil.NewMockOBSClient()
	mock.Connect() // Start connected

	// Create a minimal server for testing tool handlers
	server := &Server{
		obsClient: mock,
		ctx:       context.Background(),
	}

	return server, mock
}

// Test scene management tools

func TestHandleSetCurrentScene(t *testing.T) {
	t.Run("switches scene successfully", func(t *testing.T) {
		server, mock := testServer(t)

		input := SceneNameInput{SceneName: "Gaming"}
		_, result, err := server.handleSetCurrentScene(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		simpleResult, ok := result.(SimpleResult)
		require.True(t, ok)
		assert.Contains(t, simpleResult.Message, "Gaming")

		// Verify scene was changed
		scenes, current, _ := mock.GetSceneList()
		_ = scenes
		assert.Equal(t, "Gaming", current)
	})

	t.Run("returns error for non-existent scene", func(t *testing.T) {
		server, _ := testServer(t)

		input := SceneNameInput{SceneName: "NonExistent"}
		_, _, err := server.handleSetCurrentScene(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("returns error when not connected", func(t *testing.T) {
		server, mock := testServer(t)
		mock.Disconnect()

		input := SceneNameInput{SceneName: "Gaming"}
		_, _, err := server.handleSetCurrentScene(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not connected")
	})
}

func TestHandleCreateScene(t *testing.T) {
	t.Run("creates scene successfully", func(t *testing.T) {
		server, mock := testServer(t)

		input := SceneNameInput{SceneName: "New Scene"}
		_, result, err := server.handleCreateScene(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		simpleResult, ok := result.(SimpleResult)
		require.True(t, ok)
		assert.Contains(t, simpleResult.Message, "New Scene")

		// Verify scene was created
		scenes, _, _ := mock.GetSceneList()
		assert.Contains(t, scenes, "New Scene")
	})

	t.Run("returns error for duplicate scene", func(t *testing.T) {
		server, _ := testServer(t)

		input := SceneNameInput{SceneName: "Scene 1"} // Already exists in mock
		_, _, err := server.handleCreateScene(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestHandleRemoveScene(t *testing.T) {
	t.Run("removes scene successfully", func(t *testing.T) {
		server, mock := testServer(t)

		input := SceneNameInput{SceneName: "Gaming"}
		_, result, err := server.handleRemoveScene(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		simpleResult, ok := result.(SimpleResult)
		require.True(t, ok)
		assert.Contains(t, simpleResult.Message, "Gaming")

		// Verify scene was removed
		scenes, _, _ := mock.GetSceneList()
		assert.NotContains(t, scenes, "Gaming")
	})

	t.Run("returns error for non-existent scene", func(t *testing.T) {
		server, _ := testServer(t)

		input := SceneNameInput{SceneName: "NonExistent"}
		_, _, err := server.handleRemoveScene(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestHandleListScenes(t *testing.T) {
	t.Run("lists all scenes", func(t *testing.T) {
		server, _ := testServer(t)

		_, result, err := server.handleListScenes(context.Background(), nil, struct{}{})

		assert.NoError(t, err)
		assert.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Contains(t, resultMap, "scenes")
		assert.Contains(t, resultMap, "current_scene")

		scenes, ok := resultMap["scenes"].([]string)
		require.True(t, ok)
		assert.Len(t, scenes, 4) // Default mock has 4 scenes
	})
}

// Test recording tools

func TestHandleStartRecording(t *testing.T) {
	t.Run("starts recording successfully", func(t *testing.T) {
		server, mock := testServer(t)

		_, result, err := server.handleStartRecording(context.Background(), nil, struct{}{})

		assert.NoError(t, err)
		assert.NotNil(t, result)

		simpleResult, ok := result.(SimpleResult)
		require.True(t, ok)
		assert.Contains(t, simpleResult.Message, "started recording")

		// Verify recording state
		status, _ := mock.GetRecordingStatus()
		assert.True(t, status.Active)
	})

	t.Run("returns error when already recording", func(t *testing.T) {
		server, mock := testServer(t)
		mock.SetRecordingState(true, false)

		_, _, err := server.handleStartRecording(context.Background(), nil, struct{}{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already active")
	})
}

func TestHandleStopRecording(t *testing.T) {
	t.Run("stops recording successfully", func(t *testing.T) {
		server, mock := testServer(t)
		mock.SetRecordingState(true, false)

		_, result, err := server.handleStopRecording(context.Background(), nil, struct{}{})

		assert.NoError(t, err)
		assert.NotNil(t, result)

		simpleResult, ok := result.(SimpleResult)
		require.True(t, ok)
		assert.Contains(t, simpleResult.Message, "stopped recording")

		// Verify recording state
		status, _ := mock.GetRecordingStatus()
		assert.False(t, status.Active)
	})

	t.Run("returns error when not recording", func(t *testing.T) {
		server, _ := testServer(t)

		_, _, err := server.handleStopRecording(context.Background(), nil, struct{}{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not active")
	})
}

func TestHandleGetRecordingStatus(t *testing.T) {
	t.Run("returns recording status", func(t *testing.T) {
		server, mock := testServer(t)
		mock.SetRecordingState(true, false)

		_, result, err := server.handleGetRecordingStatus(context.Background(), nil, struct{}{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestHandlePauseRecording(t *testing.T) {
	t.Run("pauses recording successfully", func(t *testing.T) {
		server, mock := testServer(t)
		mock.SetRecordingState(true, false)

		_, result, err := server.handlePauseRecording(context.Background(), nil, struct{}{})

		assert.NoError(t, err)
		assert.NotNil(t, result)

		simpleResult, ok := result.(SimpleResult)
		require.True(t, ok)
		assert.Contains(t, simpleResult.Message, "paused recording")

		// Verify paused state
		status, _ := mock.GetRecordingStatus()
		assert.True(t, status.Paused)
	})

	t.Run("returns error when not recording", func(t *testing.T) {
		server, _ := testServer(t)

		_, _, err := server.handlePauseRecording(context.Background(), nil, struct{}{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not active")
	})
}

func TestHandleResumeRecording(t *testing.T) {
	t.Run("resumes recording successfully", func(t *testing.T) {
		server, mock := testServer(t)
		mock.SetRecordingState(true, true) // Recording and paused

		_, result, err := server.handleResumeRecording(context.Background(), nil, struct{}{})

		assert.NoError(t, err)
		assert.NotNil(t, result)

		simpleResult, ok := result.(SimpleResult)
		require.True(t, ok)
		assert.Contains(t, simpleResult.Message, "resumed recording")

		// Verify resumed state
		status, _ := mock.GetRecordingStatus()
		assert.False(t, status.Paused)
	})

	t.Run("returns error when not paused", func(t *testing.T) {
		server, mock := testServer(t)
		mock.SetRecordingState(true, false) // Recording but not paused

		_, _, err := server.handleResumeRecording(context.Background(), nil, struct{}{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not paused")
	})
}

// Test streaming tools

func TestHandleStartStreaming(t *testing.T) {
	t.Run("starts streaming successfully", func(t *testing.T) {
		server, mock := testServer(t)

		_, result, err := server.handleStartStreaming(context.Background(), nil, struct{}{})

		assert.NoError(t, err)
		assert.NotNil(t, result)

		simpleResult, ok := result.(SimpleResult)
		require.True(t, ok)
		assert.Contains(t, simpleResult.Message, "started streaming")

		// Verify streaming state
		status, _ := mock.GetStreamingStatus()
		assert.True(t, status.Active)
	})

	t.Run("returns error when already streaming", func(t *testing.T) {
		server, mock := testServer(t)
		mock.SetStreamingState(true)

		_, _, err := server.handleStartStreaming(context.Background(), nil, struct{}{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already active")
	})
}

func TestHandleStopStreaming(t *testing.T) {
	t.Run("stops streaming successfully", func(t *testing.T) {
		server, mock := testServer(t)
		mock.SetStreamingState(true)

		_, result, err := server.handleStopStreaming(context.Background(), nil, struct{}{})

		assert.NoError(t, err)
		assert.NotNil(t, result)

		simpleResult, ok := result.(SimpleResult)
		require.True(t, ok)
		assert.Contains(t, simpleResult.Message, "stopped streaming")

		// Verify streaming state
		status, _ := mock.GetStreamingStatus()
		assert.False(t, status.Active)
	})

	t.Run("returns error when not streaming", func(t *testing.T) {
		server, _ := testServer(t)

		_, _, err := server.handleStopStreaming(context.Background(), nil, struct{}{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not active")
	})
}

func TestHandleGetStreamingStatus(t *testing.T) {
	t.Run("returns streaming status", func(t *testing.T) {
		server, mock := testServer(t)
		mock.SetStreamingState(true)

		_, result, err := server.handleGetStreamingStatus(context.Background(), nil, struct{}{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// Test source tools

func TestHandleListSources(t *testing.T) {
	t.Run("lists all sources", func(t *testing.T) {
		server, _ := testServer(t)

		_, result, err := server.handleListSources(context.Background(), nil, struct{}{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestHandleToggleSourceVisibility(t *testing.T) {
	t.Run("toggles visibility successfully", func(t *testing.T) {
		server, _ := testServer(t)

		input := SourceVisibilityInput{
			SceneName: "Scene 1",
			SourceID:  1,
		}
		_, result, err := server.handleToggleSourceVisibility(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Contains(t, resultMap, "visible")
	})

	t.Run("returns error for non-existent scene", func(t *testing.T) {
		server, _ := testServer(t)

		input := SourceVisibilityInput{
			SceneName: "NonExistent",
			SourceID:  1,
		}
		_, _, err := server.handleToggleSourceVisibility(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestHandleGetSourceSettings(t *testing.T) {
	t.Run("returns source settings", func(t *testing.T) {
		server, _ := testServer(t)

		input := SourceNameInput{SourceName: "Microphone"}
		_, result, err := server.handleGetSourceSettings(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		settings, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Contains(t, settings, "device_id")
	})

	t.Run("returns error for non-existent source", func(t *testing.T) {
		server, _ := testServer(t)

		input := SourceNameInput{SourceName: "NonExistent"}
		_, _, err := server.handleGetSourceSettings(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

// Test audio tools

func TestHandleGetInputMute(t *testing.T) {
	t.Run("returns mute status", func(t *testing.T) {
		server, mock := testServer(t)
		mock.SetInputMuteState("Microphone", true)

		input := InputNameInput{InputName: "Microphone"}
		_, result, err := server.handleGetInputMute(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, true, resultMap["is_muted"])
	})

	t.Run("returns error for non-existent input", func(t *testing.T) {
		server, _ := testServer(t)

		input := InputNameInput{InputName: "NonExistent"}
		_, _, err := server.handleGetInputMute(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestHandleToggleInputMute(t *testing.T) {
	t.Run("toggles mute successfully", func(t *testing.T) {
		server, _ := testServer(t)

		input := InputNameInput{InputName: "Microphone"}
		_, result, err := server.handleToggleInputMute(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		simpleResult, ok := result.(SimpleResult)
		require.True(t, ok)
		assert.Contains(t, simpleResult.Message, "Microphone")
	})

	t.Run("returns error for non-existent input", func(t *testing.T) {
		server, _ := testServer(t)

		input := InputNameInput{InputName: "NonExistent"}
		_, _, err := server.handleToggleInputMute(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestHandleSetInputVolume(t *testing.T) {
	t.Run("sets volume successfully with dB", func(t *testing.T) {
		server, _ := testServer(t)

		volumeDb := -6.0
		input := SetVolumeInput{
			InputName: "Microphone",
			VolumeDb:  &volumeDb,
		}
		_, result, err := server.handleSetInputVolume(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		simpleResult, ok := result.(SimpleResult)
		require.True(t, ok)
		assert.Contains(t, simpleResult.Message, "Microphone")
	})

	t.Run("returns error for non-existent input", func(t *testing.T) {
		server, _ := testServer(t)

		volumeDb := -6.0
		input := SetVolumeInput{
			InputName: "NonExistent",
			VolumeDb:  &volumeDb,
		}
		_, _, err := server.handleSetInputVolume(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

// Test status tool

func TestHandleGetOBSStatus(t *testing.T) {
	t.Run("returns OBS status", func(t *testing.T) {
		server, _ := testServer(t)

		_, result, err := server.handleGetOBSStatus(context.Background(), nil, struct{}{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("returns error when not connected", func(t *testing.T) {
		server, mock := testServer(t)
		mock.Disconnect()

		_, _, err := server.handleGetOBSStatus(context.Background(), nil, struct{}{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not connected")
	})
}

// Test input type validation

func TestInputTypes(t *testing.T) {
	t.Run("SceneNameInput has correct JSON tag", func(t *testing.T) {
		input := SceneNameInput{SceneName: "Test"}
		assert.Equal(t, "Test", input.SceneName)
	})

	t.Run("SourceNameInput has correct JSON tag", func(t *testing.T) {
		input := SourceNameInput{SourceName: "Test"}
		assert.Equal(t, "Test", input.SourceName)
	})

	t.Run("InputNameInput has correct JSON tag", func(t *testing.T) {
		input := InputNameInput{InputName: "Test"}
		assert.Equal(t, "Test", input.InputName)
	})

	t.Run("SetVolumeInput handles optional fields", func(t *testing.T) {
		// Test with VolumeDb set
		db := -6.0
		inputDb := SetVolumeInput{InputName: "Test", VolumeDb: &db}
		assert.Equal(t, &db, inputDb.VolumeDb)
		assert.Nil(t, inputDb.VolumeMul)

		// Test with VolumeMul set
		mul := 0.5
		inputMul := SetVolumeInput{InputName: "Test", VolumeMul: &mul}
		assert.Nil(t, inputMul.VolumeDb)
		assert.Equal(t, &mul, inputMul.VolumeMul)
	})
}

// Test SimpleResult type

func TestSimpleResult(t *testing.T) {
	result := SimpleResult{Message: "Test message"}
	assert.Equal(t, "Test message", result.Message)
}

// Integration-style test for tool request flow

func TestToolRequestFlow(t *testing.T) {
	t.Run("complete scene workflow", func(t *testing.T) {
		server, mock := testServer(t)

		// List initial scenes
		_, listResult, err := server.handleListScenes(context.Background(), nil, struct{}{})
		require.NoError(t, err)
		resultMap := listResult.(map[string]interface{})
		initialScenes := resultMap["scenes"].([]string)
		initialCount := len(initialScenes)

		// Create a new scene
		_, _, err = server.handleCreateScene(context.Background(), nil, SceneNameInput{SceneName: "Workflow Test"})
		require.NoError(t, err)

		// Verify scene was added
		scenes, _, _ := mock.GetSceneList()
		assert.Len(t, scenes, initialCount+1)

		// Switch to the new scene
		_, _, err = server.handleSetCurrentScene(context.Background(), nil, SceneNameInput{SceneName: "Workflow Test"})
		require.NoError(t, err)

		// Verify current scene changed
		_, current, _ := mock.GetSceneList()
		assert.Equal(t, "Workflow Test", current)

		// Remove the scene
		_, _, err = server.handleRemoveScene(context.Background(), nil, SceneNameInput{SceneName: "Workflow Test"})
		require.NoError(t, err)

		// Verify scene was removed
		scenes, _, _ = mock.GetSceneList()
		assert.Len(t, scenes, initialCount)
	})

	t.Run("complete recording workflow", func(t *testing.T) {
		server, mock := testServer(t)

		// Start recording
		_, _, err := server.handleStartRecording(context.Background(), nil, struct{}{})
		require.NoError(t, err)

		// Get status - should be recording
		status, _ := mock.GetRecordingStatus()
		assert.True(t, status.Active)
		assert.False(t, status.Paused)

		// Pause recording
		_, _, err = server.handlePauseRecording(context.Background(), nil, struct{}{})
		require.NoError(t, err)

		// Get status - should be paused
		status, _ = mock.GetRecordingStatus()
		assert.True(t, status.Active)
		assert.True(t, status.Paused)

		// Resume recording
		_, _, err = server.handleResumeRecording(context.Background(), nil, struct{}{})
		require.NoError(t, err)

		// Get status - should be recording again
		status, _ = mock.GetRecordingStatus()
		assert.True(t, status.Active)
		assert.False(t, status.Paused)

		// Stop recording
		_, _, err = server.handleStopRecording(context.Background(), nil, struct{}{})
		require.NoError(t, err)

		// Get status - should be stopped
		status, _ = mock.GetRecordingStatus()
		assert.False(t, status.Active)
	})
}

// Verify mock implements OBSClient interface
func TestMockImplementsInterface(t *testing.T) {
	// This test verifies at compile time that MockOBSClient implements OBSClient
	var _ OBSClient = (*testutil.MockOBSClient)(nil)
}

// Ensure we don't have nil pointer issues with the mock request
func TestNilRequest(t *testing.T) {
	t.Run("handlers work with nil CallToolRequest", func(t *testing.T) {
		server, _ := testServer(t)

		// Test that handlers don't panic with nil request
		_, _, err := server.handleListScenes(context.Background(), nil, struct{}{})
		assert.NoError(t, err)

		_, _, err = server.handleGetOBSStatus(context.Background(), nil, struct{}{})
		assert.NoError(t, err)
	})
}

// testServerWithStorage creates a test server with mock OBS client and real storage.
func testServerWithStorage(t *testing.T) (*Server, *testutil.MockOBSClient, *storage.DB) {
	t.Helper()

	mock := testutil.NewMockOBSClient()
	mock.Connect()

	// Create temp database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := storage.New(context.Background(), storage.Config{Path: dbPath})
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
		os.Remove(dbPath)
	})

	server := &Server{
		obsClient: mock,
		storage:   db,
		ctx:       context.Background(),
	}

	return server, mock, db
}

// Test scene preset tools

func TestHandleListScenePresets(t *testing.T) {
	t.Run("returns empty list initially", func(t *testing.T) {
		server, _, _ := testServerWithStorage(t)

		input := ListPresetsInput{}
		_, result, err := server.handleListScenePresets(context.Background(), nil, input)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		resultMap, ok := result.(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, 0, resultMap["count"])
	})

	t.Run("lists created presets", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)

		// Create a preset directly in storage
		_, err := db.CreateScenePreset(context.Background(), storage.ScenePreset{
			Name:      "Test Preset",
			SceneName: "Scene 1",
			Sources:   []storage.SourceState{{Name: "Webcam", Visible: true}},
		})
		require.NoError(t, err)

		input := ListPresetsInput{}
		_, result, err := server.handleListScenePresets(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		assert.Equal(t, 1, resultMap["count"])

		presets := resultMap["presets"].([]map[string]interface{})
		assert.Equal(t, "Test Preset", presets[0]["name"])
	})

	t.Run("filters by scene name", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)

		// Create presets for different scenes
		db.CreateScenePreset(context.Background(), storage.ScenePreset{
			Name: "Preset 1", SceneName: "Scene 1", Sources: []storage.SourceState{},
		})
		db.CreateScenePreset(context.Background(), storage.ScenePreset{
			Name: "Preset 2", SceneName: "Scene 2", Sources: []storage.SourceState{},
		})

		input := ListPresetsInput{SceneName: "Scene 1"}
		_, result, err := server.handleListScenePresets(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		assert.Equal(t, 1, resultMap["count"])
	})
}

func TestHandleGetPresetDetails(t *testing.T) {
	t.Run("returns preset details", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)

		// Create a preset
		db.CreateScenePreset(context.Background(), storage.ScenePreset{
			Name:      "My Preset",
			SceneName: "Gaming",
			Sources:   []storage.SourceState{{Name: "Webcam", Visible: true}},
		})

		input := PresetNameInput{PresetName: "My Preset"}
		_, result, err := server.handleGetPresetDetails(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		assert.Equal(t, "My Preset", resultMap["name"])
		assert.Equal(t, "Gaming", resultMap["scene_name"])
		assert.NotNil(t, resultMap["sources"])
	})

	t.Run("returns error for non-existent preset", func(t *testing.T) {
		server, _, _ := testServerWithStorage(t)

		input := PresetNameInput{PresetName: "NonExistent"}
		_, _, err := server.handleGetPresetDetails(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestHandleDeleteScenePreset(t *testing.T) {
	t.Run("deletes preset successfully", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)

		// Create a preset
		db.CreateScenePreset(context.Background(), storage.ScenePreset{
			Name: "To Delete", SceneName: "Scene 1", Sources: []storage.SourceState{},
		})

		input := PresetNameInput{PresetName: "To Delete"}
		_, result, err := server.handleDeleteScenePreset(context.Background(), nil, input)

		assert.NoError(t, err)
		simpleResult := result.(SimpleResult)
		assert.Contains(t, simpleResult.Message, "To Delete")

		// Verify preset was deleted
		_, err = db.GetScenePreset(context.Background(), "To Delete")
		assert.Error(t, err)
	})

	t.Run("returns error for non-existent preset", func(t *testing.T) {
		server, _, _ := testServerWithStorage(t)

		input := PresetNameInput{PresetName: "NonExistent"}
		_, _, err := server.handleDeleteScenePreset(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestHandleRenameScenePreset(t *testing.T) {
	t.Run("renames preset successfully", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)

		// Create a preset
		db.CreateScenePreset(context.Background(), storage.ScenePreset{
			Name: "Old Name", SceneName: "Scene 1", Sources: []storage.SourceState{},
		})

		input := RenamePresetInput{OldName: "Old Name", NewName: "New Name"}
		_, result, err := server.handleRenameScenePreset(context.Background(), nil, input)

		assert.NoError(t, err)
		simpleResult := result.(SimpleResult)
		assert.Contains(t, simpleResult.Message, "New Name")

		// Verify preset was renamed
		preset, err := db.GetScenePreset(context.Background(), "New Name")
		assert.NoError(t, err)
		assert.Equal(t, "New Name", preset.Name)
	})

	t.Run("returns error for non-existent preset", func(t *testing.T) {
		server, _, _ := testServerWithStorage(t)

		input := RenamePresetInput{OldName: "NonExistent", NewName: "New Name"}
		_, _, err := server.handleRenameScenePreset(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestHandleGetInputVolume(t *testing.T) {
	t.Run("returns input volume", func(t *testing.T) {
		server, mock := testServer(t)
		mock.SetInputVolumeState("Microphone", -6.0)

		input := InputNameInput{InputName: "Microphone"}
		_, result, err := server.handleGetInputVolume(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		assert.Equal(t, "Microphone", resultMap["input_name"])
		assert.Equal(t, -6.0, resultMap["volume_db"])
		assert.NotNil(t, resultMap["volume_mul"])
	})

	t.Run("returns error for non-existent input", func(t *testing.T) {
		server, _ := testServer(t)

		input := InputNameInput{InputName: "NonExistent"}
		_, _, err := server.handleGetInputVolume(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestHandleSaveScenePreset(t *testing.T) {
	t.Run("saves scene preset successfully", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)

		input := SavePresetInput{PresetName: "Gaming Setup", SceneName: "Scene 1"}
		_, result, err := server.handleSaveScenePreset(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		assert.Equal(t, "Gaming Setup", resultMap["preset_name"])
		assert.Equal(t, "Scene 1", resultMap["scene_name"])
		assert.NotNil(t, resultMap["id"])

		// Verify preset was saved
		preset, err := db.GetScenePreset(context.Background(), "Gaming Setup")
		assert.NoError(t, err)
		assert.Equal(t, "Scene 1", preset.SceneName)
		assert.Len(t, preset.Sources, 2) // Scene 1 has 2 sources in mock
	})

	t.Run("returns error for non-existent scene", func(t *testing.T) {
		server, _, _ := testServerWithStorage(t)

		input := SavePresetInput{PresetName: "Test", SceneName: "NonExistent"}
		_, _, err := server.handleSaveScenePreset(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("returns error for duplicate preset name", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)

		// Create a preset first
		db.CreateScenePreset(context.Background(), storage.ScenePreset{
			Name: "Duplicate", SceneName: "Scene 1", Sources: []storage.SourceState{},
		})

		input := SavePresetInput{PresetName: "Duplicate", SceneName: "Scene 1"}
		_, _, err := server.handleSaveScenePreset(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "UNIQUE constraint failed")
	})
}

func TestHandleApplyScenePreset(t *testing.T) {
	t.Run("applies scene preset successfully", func(t *testing.T) {
		server, mock, db := testServerWithStorage(t)

		// Create a preset with specific source states
		db.CreateScenePreset(context.Background(), storage.ScenePreset{
			Name:      "Test Preset",
			SceneName: "Scene 1",
			Sources: []storage.SourceState{
				{Name: "Webcam", Visible: false}, // Will be disabled
				{Name: "Text", Visible: true},    // Will stay enabled
			},
		})

		input := PresetNameInput{PresetName: "Test Preset"}
		_, result, err := server.handleApplyScenePreset(context.Background(), nil, input)

		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		assert.Equal(t, "Test Preset", resultMap["preset_name"])
		assert.Equal(t, "Scene 1", resultMap["scene_name"])

		// Verify source states were applied
		scene, _ := mock.GetSceneByName("Scene 1")
		for _, src := range scene.Sources {
			if src.Name == "Webcam" {
				assert.False(t, src.Enabled)
			}
			if src.Name == "Text" {
				assert.True(t, src.Enabled)
			}
		}
	})

	t.Run("returns error for non-existent preset", func(t *testing.T) {
		server, _, _ := testServerWithStorage(t)

		input := PresetNameInput{PresetName: "NonExistent"}
		_, _, err := server.handleApplyScenePreset(context.Background(), nil, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("skips sources not found in scene", func(t *testing.T) {
		server, _, db := testServerWithStorage(t)

		// Create a preset with a source that doesn't exist in the scene
		db.CreateScenePreset(context.Background(), storage.ScenePreset{
			Name:      "Preset With Missing Source",
			SceneName: "Scene 1",
			Sources: []storage.SourceState{
				{Name: "Webcam", Visible: true},
				{Name: "NonExistentSource", Visible: false}, // This source doesn't exist
			},
		})

		input := PresetNameInput{PresetName: "Preset With Missing Source"}
		_, result, err := server.handleApplyScenePreset(context.Background(), nil, input)

		// Should succeed but only apply 1 source
		assert.NoError(t, err)
		resultMap := result.(map[string]interface{})
		assert.Equal(t, 1, resultMap["applied_count"])
	})
}

// Test complete preset workflow
func TestPresetWorkflow(t *testing.T) {
	t.Run("complete save and apply workflow", func(t *testing.T) {
		server, mock, _ := testServerWithStorage(t)

		// Save the current state of Scene 1
		saveInput := SavePresetInput{PresetName: "Scene1 State", SceneName: "Scene 1"}
		_, _, err := server.handleSaveScenePreset(context.Background(), nil, saveInput)
		require.NoError(t, err)

		// Modify the scene by toggling a source
		mock.ToggleSourceVisibility("Scene 1", 1) // Toggle Webcam

		// Apply the saved preset to restore original state
		applyInput := PresetNameInput{PresetName: "Scene1 State"}
		_, _, err = server.handleApplyScenePreset(context.Background(), nil, applyInput)
		require.NoError(t, err)

		// Verify the scene was restored
		scene, _ := mock.GetSceneByName("Scene 1")
		for _, src := range scene.Sources {
			if src.ID == 1 {
				assert.True(t, src.Enabled) // Should be restored to original state
			}
		}
	})
}
