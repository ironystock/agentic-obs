package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ironystock/agentic-obs/internal/obs"
)

func TestExtractSceneNameFromURI(t *testing.T) {
	t.Run("extracts scene name from valid URI", func(t *testing.T) {
		name, err := extractSceneNameFromURI("obs://scene/Gaming")
		assert.NoError(t, err)
		assert.Equal(t, "Gaming", name)
	})

	t.Run("extracts scene name with spaces", func(t *testing.T) {
		name, err := extractSceneNameFromURI("obs://scene/Starting Soon")
		assert.NoError(t, err)
		assert.Equal(t, "Starting Soon", name)
	})

	t.Run("returns error for URI too short", func(t *testing.T) {
		_, err := extractSceneNameFromURI("obs://scene/")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too short")
	})

	t.Run("returns error for invalid prefix", func(t *testing.T) {
		_, err := extractSceneNameFromURI("http://scene/Gaming")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must start with")
	})

	t.Run("returns error for empty URI", func(t *testing.T) {
		_, err := extractSceneNameFromURI("")
		assert.Error(t, err)
	})
}

func TestConvertSourcesToMap(t *testing.T) {
	t.Run("converts empty sources list", func(t *testing.T) {
		result := convertSourcesToMap([]obs.SceneSource{})
		assert.Empty(t, result)
	})

	t.Run("converts single source", func(t *testing.T) {
		sources := []obs.SceneSource{
			{
				ID:       1,
				Name:     "Webcam",
				Type:     "dshow_input",
				Enabled:  true,
				Visible:  true,
				Locked:   false,
				X:        100,
				Y:        200,
				Width:    1920,
				Height:   1080,
				ScaleX:   1.0,
				ScaleY:   1.0,
				Rotation: 0,
			},
		}

		result := convertSourcesToMap(sources)

		assert.Len(t, result, 1)
		assert.Equal(t, 1, result[0]["id"])
		assert.Equal(t, "Webcam", result[0]["name"])
		assert.Equal(t, "dshow_input", result[0]["type"])
		assert.Equal(t, true, result[0]["enabled"])
		assert.Equal(t, true, result[0]["visible"])
		assert.Equal(t, false, result[0]["locked"])
		assert.Equal(t, float64(100), result[0]["x"])
		assert.Equal(t, float64(200), result[0]["y"])
		assert.Equal(t, float64(1920), result[0]["width"])
		assert.Equal(t, float64(1080), result[0]["height"])
	})

	t.Run("converts multiple sources", func(t *testing.T) {
		sources := []obs.SceneSource{
			{ID: 1, Name: "Source1"},
			{ID: 2, Name: "Source2"},
			{ID: 3, Name: "Source3"},
		}

		result := convertSourcesToMap(sources)

		assert.Len(t, result, 3)
		assert.Equal(t, "Source1", result[0]["name"])
		assert.Equal(t, "Source2", result[1]["name"])
		assert.Equal(t, "Source3", result[2]["name"])
	})
}

func TestSceneDetails(t *testing.T) {
	t.Run("SceneDetails struct fields", func(t *testing.T) {
		details := SceneDetails{
			Name:        "Test Scene",
			IsActive:    true,
			SceneIndex:  0,
			Description: "A test scene",
			Sources: []map[string]interface{}{
				{"id": 1, "name": "Source1"},
			},
		}

		assert.Equal(t, "Test Scene", details.Name)
		assert.True(t, details.IsActive)
		assert.Equal(t, 0, details.SceneIndex)
		assert.Equal(t, "A test scene", details.Description)
		assert.Len(t, details.Sources, 1)
	})
}

func TestExtractScreenshotNameFromURI(t *testing.T) {
	t.Run("extracts screenshot name from valid URI", func(t *testing.T) {
		name, err := extractScreenshotNameFromURI("obs://screenshot/gameplay")
		assert.NoError(t, err)
		assert.Equal(t, "gameplay", name)
	})

	t.Run("extracts screenshot name with spaces", func(t *testing.T) {
		name, err := extractScreenshotNameFromURI("obs://screenshot/my monitor")
		assert.NoError(t, err)
		assert.Equal(t, "my monitor", name)
	})

	t.Run("extracts screenshot name with special characters", func(t *testing.T) {
		name, err := extractScreenshotNameFromURI("obs://screenshot/webcam_1080p")
		assert.NoError(t, err)
		assert.Equal(t, "webcam_1080p", name)
	})

	t.Run("returns error for URI too short", func(t *testing.T) {
		_, err := extractScreenshotNameFromURI("obs://screenshot/")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too short")
	})

	t.Run("returns error for invalid prefix", func(t *testing.T) {
		_, err := extractScreenshotNameFromURI("obs://scene/longenoughtestname")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must start with")
	})

	t.Run("returns error for wrong scheme", func(t *testing.T) {
		_, err := extractScreenshotNameFromURI("http://screenshot/gameplay")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must start with")
	})

	t.Run("returns error for empty URI", func(t *testing.T) {
		_, err := extractScreenshotNameFromURI("")
		assert.Error(t, err)
	})
}

func TestExtractPresetNameFromURI(t *testing.T) {
	t.Run("extracts preset name from valid URI", func(t *testing.T) {
		name, err := extractPresetNameFromURI("obs://preset/streaming")
		assert.NoError(t, err)
		assert.Equal(t, "streaming", name)
	})

	t.Run("extracts preset name with spaces", func(t *testing.T) {
		name, err := extractPresetNameFromURI("obs://preset/my layout")
		assert.NoError(t, err)
		assert.Equal(t, "my layout", name)
	})

	t.Run("extracts preset name with special characters", func(t *testing.T) {
		name, err := extractPresetNameFromURI("obs://preset/recording_setup_v2")
		assert.NoError(t, err)
		assert.Equal(t, "recording_setup_v2", name)
	})

	t.Run("returns error for URI too short", func(t *testing.T) {
		_, err := extractPresetNameFromURI("obs://preset/")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too short")
	})

	t.Run("returns error for invalid prefix", func(t *testing.T) {
		_, err := extractPresetNameFromURI("obs://scene/test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must start with")
	})

	t.Run("returns error for wrong scheme", func(t *testing.T) {
		_, err := extractPresetNameFromURI("http://preset/streaming")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must start with")
	})

	t.Run("returns error for empty URI", func(t *testing.T) {
		_, err := extractPresetNameFromURI("")
		assert.Error(t, err)
	})
}
