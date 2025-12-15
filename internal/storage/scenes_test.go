package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateScenePreset(t *testing.T) {
	t.Run("creates preset successfully", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		preset := ScenePreset{
			Name:      "Gaming Setup",
			SceneName: "Gaming",
			Sources: []SourceState{
				{Name: "Webcam", Visible: true},
				{Name: "Game Capture", Visible: true},
			},
		}

		id, err := db.CreateScenePreset(context.Background(), preset)
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})

	t.Run("creates preset with empty sources", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		preset := ScenePreset{
			Name:      "Empty Preset",
			SceneName: "Idle",
			Sources:   []SourceState{},
		}

		id, err := db.CreateScenePreset(context.Background(), preset)
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})

	t.Run("creates preset with source settings", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		preset := ScenePreset{
			Name:      "Complex Preset",
			SceneName: "Stream",
			Sources: []SourceState{
				{
					Name:    "Webcam",
					Visible: true,
					Settings: map[string]interface{}{
						"device_id":  "usb-camera-1",
						"resolution": "1920x1080",
						"fps":        30,
					},
				},
			},
		}

		id, err := db.CreateScenePreset(context.Background(), preset)
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})

	t.Run("fails on duplicate name", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		preset := ScenePreset{
			Name:      "Duplicate",
			SceneName: "Test",
		}

		_, err := db.CreateScenePreset(context.Background(), preset)
		require.NoError(t, err)

		_, err = db.CreateScenePreset(context.Background(), preset)
		assert.Error(t, err)
		// Check for SQLite unique constraint error
		assert.Contains(t, err.Error(), "UNIQUE constraint failed")
	})
}

func TestGetScenePreset(t *testing.T) {
	t.Run("retrieves preset by name", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		original := ScenePreset{
			Name:      "Test Preset",
			SceneName: "Gaming",
			Sources: []SourceState{
				{Name: "Webcam", Visible: true},
			},
		}

		id, err := db.CreateScenePreset(context.Background(), original)
		require.NoError(t, err)

		preset, err := db.GetScenePreset(context.Background(), "Test Preset")
		assert.NoError(t, err)
		assert.Equal(t, id, preset.ID)
		assert.Equal(t, "Test Preset", preset.Name)
		assert.Equal(t, "Gaming", preset.SceneName)
		assert.Len(t, preset.Sources, 1)
		assert.Equal(t, "Webcam", preset.Sources[0].Name)
		assert.True(t, preset.Sources[0].Visible)
	})

	t.Run("returns error for non-existent preset", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		_, err := db.GetScenePreset(context.Background(), "NonExistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("handles empty sources", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		original := ScenePreset{
			Name:      "Empty",
			SceneName: "Test",
		}

		_, err := db.CreateScenePreset(context.Background(), original)
		require.NoError(t, err)

		preset, err := db.GetScenePreset(context.Background(), "Empty")
		assert.NoError(t, err)
		assert.Empty(t, preset.Sources)
	})
}

func TestGetScenePresetByID(t *testing.T) {
	t.Run("retrieves preset by ID", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		original := ScenePreset{
			Name:      "By ID Test",
			SceneName: "Stream",
		}

		id, err := db.CreateScenePreset(context.Background(), original)
		require.NoError(t, err)

		preset, err := db.GetScenePresetByID(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, id, preset.ID)
		assert.Equal(t, "By ID Test", preset.Name)
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		_, err := db.GetScenePresetByID(context.Background(), 99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestListScenePresets(t *testing.T) {
	t.Run("lists all presets", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Create multiple presets
		for i, name := range []string{"Preset A", "Preset B", "Preset C"} {
			_, err := db.CreateScenePreset(context.Background(), ScenePreset{
				Name:      name,
				SceneName: "Scene " + string(rune('A'+i)),
			})
			require.NoError(t, err)
		}

		presets, err := db.ListScenePresets(context.Background(), "")
		assert.NoError(t, err)
		assert.Len(t, presets, 3)
	})

	t.Run("filters by scene name", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Create presets for different scenes
		_, err := db.CreateScenePreset(context.Background(), ScenePreset{
			Name:      "Gaming 1",
			SceneName: "Gaming",
		})
		require.NoError(t, err)

		_, err = db.CreateScenePreset(context.Background(), ScenePreset{
			Name:      "Gaming 2",
			SceneName: "Gaming",
		})
		require.NoError(t, err)

		_, err = db.CreateScenePreset(context.Background(), ScenePreset{
			Name:      "Chat 1",
			SceneName: "Chat",
		})
		require.NoError(t, err)

		// Filter by Gaming scene
		presets, err := db.ListScenePresets(context.Background(), "Gaming")
		assert.NoError(t, err)
		assert.Len(t, presets, 2)

		for _, p := range presets {
			assert.Equal(t, "Gaming", p.SceneName)
		}
	})

	t.Run("returns empty list when no presets", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		presets, err := db.ListScenePresets(context.Background(), "")
		assert.NoError(t, err)
		assert.Empty(t, presets)
	})

	t.Run("returns all presets regardless of order", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Create presets
		names := []string{"First", "Second", "Third"}
		for _, name := range names {
			_, err := db.CreateScenePreset(context.Background(), ScenePreset{
				Name:      name,
				SceneName: "Test",
			})
			require.NoError(t, err)
		}

		presets, err := db.ListScenePresets(context.Background(), "")
		assert.NoError(t, err)
		assert.Len(t, presets, 3)

		// All presets should be present (order depends on timing)
		presetNames := make([]string, len(presets))
		for i, p := range presets {
			presetNames[i] = p.Name
		}
		assert.ElementsMatch(t, names, presetNames)
	})
}

func TestUpdateScenePreset(t *testing.T) {
	t.Run("updates preset successfully", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Create preset
		_, err := db.CreateScenePreset(context.Background(), ScenePreset{
			Name:      "Updatable",
			SceneName: "Original",
			Sources:   []SourceState{{Name: "Source1", Visible: true}},
		})
		require.NoError(t, err)

		// Update preset
		err = db.UpdateScenePreset(context.Background(), ScenePreset{
			Name:      "Updatable",
			SceneName: "Updated",
			Sources:   []SourceState{{Name: "Source2", Visible: false}},
		})
		assert.NoError(t, err)

		// Verify update
		preset, err := db.GetScenePreset(context.Background(), "Updatable")
		assert.NoError(t, err)
		assert.Equal(t, "Updated", preset.SceneName)
		assert.Len(t, preset.Sources, 1)
		assert.Equal(t, "Source2", preset.Sources[0].Name)
		assert.False(t, preset.Sources[0].Visible)
	})

	t.Run("returns error for non-existent preset", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.UpdateScenePreset(context.Background(), ScenePreset{
			Name:      "NonExistent",
			SceneName: "Test",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestDeleteScenePreset(t *testing.T) {
	t.Run("deletes preset by name", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		_, err := db.CreateScenePreset(context.Background(), ScenePreset{
			Name:      "ToDelete",
			SceneName: "Test",
		})
		require.NoError(t, err)

		err = db.DeleteScenePreset(context.Background(), "ToDelete")
		assert.NoError(t, err)

		// Verify deleted
		_, err = db.GetScenePreset(context.Background(), "ToDelete")
		assert.Error(t, err)
	})

	t.Run("returns error for non-existent preset", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.DeleteScenePreset(context.Background(), "NonExistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestDeleteScenePresetByID(t *testing.T) {
	t.Run("deletes preset by ID", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		id, err := db.CreateScenePreset(context.Background(), ScenePreset{
			Name:      "ToDeleteByID",
			SceneName: "Test",
		})
		require.NoError(t, err)

		err = db.DeleteScenePresetByID(context.Background(), id)
		assert.NoError(t, err)

		// Verify deleted
		_, err = db.GetScenePresetByID(context.Background(), id)
		assert.Error(t, err)
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.DeleteScenePresetByID(context.Background(), 99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestDeleteScenePresetsByScene(t *testing.T) {
	t.Run("deletes all presets for a scene", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Create presets for different scenes
		for i := 1; i <= 3; i++ {
			_, err := db.CreateScenePreset(context.Background(), ScenePreset{
				Name:      "Gaming " + string(rune('0'+i)),
				SceneName: "Gaming",
			})
			require.NoError(t, err)
		}

		_, err := db.CreateScenePreset(context.Background(), ScenePreset{
			Name:      "Chat 1",
			SceneName: "Chat",
		})
		require.NoError(t, err)

		// Delete all Gaming presets
		count, err := db.DeleteScenePresetsByScene(context.Background(), "Gaming")
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// Verify Gaming presets deleted
		presets, err := db.ListScenePresets(context.Background(), "Gaming")
		assert.NoError(t, err)
		assert.Empty(t, presets)

		// Verify Chat presets still exist
		presets, err = db.ListScenePresets(context.Background(), "Chat")
		assert.NoError(t, err)
		assert.Len(t, presets, 1)
	})

	t.Run("returns zero when no presets match", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		count, err := db.DeleteScenePresetsByScene(context.Background(), "NonExistent")
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestCountScenePresets(t *testing.T) {
	t.Run("counts all presets", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Initially empty
		count, err := db.CountScenePresets(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)

		// Add some presets
		for i := 1; i <= 5; i++ {
			_, err := db.CreateScenePreset(context.Background(), ScenePreset{
				Name:      "Preset " + string(rune('0'+i)),
				SceneName: "Test",
			})
			require.NoError(t, err)
		}

		count, err = db.CountScenePresets(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}

func TestRenameScenePreset(t *testing.T) {
	t.Run("renames preset successfully", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		_, err := db.CreateScenePreset(context.Background(), ScenePreset{
			Name:      "OldName",
			SceneName: "Test",
		})
		require.NoError(t, err)

		err = db.RenameScenePreset(context.Background(), "OldName", "NewName")
		assert.NoError(t, err)

		// Old name should not exist
		_, err = db.GetScenePreset(context.Background(), "OldName")
		assert.Error(t, err)

		// New name should exist
		preset, err := db.GetScenePreset(context.Background(), "NewName")
		assert.NoError(t, err)
		assert.Equal(t, "NewName", preset.Name)
	})

	t.Run("returns error for non-existent preset", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		err := db.RenameScenePreset(context.Background(), "NonExistent", "NewName")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("returns error when new name already exists", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		_, err := db.CreateScenePreset(context.Background(), ScenePreset{
			Name:      "First",
			SceneName: "Test",
		})
		require.NoError(t, err)

		_, err = db.CreateScenePreset(context.Background(), ScenePreset{
			Name:      "Second",
			SceneName: "Test",
		})
		require.NoError(t, err)

		err = db.RenameScenePreset(context.Background(), "First", "Second")
		assert.Error(t, err)
		// Check for SQLite unique constraint error
		assert.Contains(t, err.Error(), "UNIQUE constraint failed")
	})
}

func TestScenePreset_SourceStateJSON(t *testing.T) {
	t.Run("preserves source settings through roundtrip", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		original := ScenePreset{
			Name:      "Complex",
			SceneName: "Stream",
			Sources: []SourceState{
				{
					Name:    "Camera",
					Visible: true,
					Settings: map[string]interface{}{
						"device_id":  "cam-001",
						"resolution": "1920x1080",
						"fps":        60.0,
						"enabled":    true,
					},
				},
				{
					Name:    "Overlay",
					Visible: false,
					Settings: map[string]interface{}{
						"file": "/path/to/overlay.png",
					},
				},
			},
		}

		_, err := db.CreateScenePreset(context.Background(), original)
		require.NoError(t, err)

		preset, err := db.GetScenePreset(context.Background(), "Complex")
		assert.NoError(t, err)
		assert.Len(t, preset.Sources, 2)

		// Check first source
		assert.Equal(t, "Camera", preset.Sources[0].Name)
		assert.True(t, preset.Sources[0].Visible)
		assert.Equal(t, "cam-001", preset.Sources[0].Settings["device_id"])
		assert.Equal(t, "1920x1080", preset.Sources[0].Settings["resolution"])

		// Check second source
		assert.Equal(t, "Overlay", preset.Sources[1].Name)
		assert.False(t, preset.Sources[1].Visible)
	})
}
