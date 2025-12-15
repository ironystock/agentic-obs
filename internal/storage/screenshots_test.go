package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupScreenshotTestDB(t *testing.T) (*DB, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "screenshot_test_*")
	require.NoError(t, err)

	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := New(context.Background(), Config{Path: dbPath})
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
		os.RemoveAll(tmpDir)
	}

	return db, cleanup
}

func TestCreateScreenshotSource(t *testing.T) {
	db, cleanup := setupScreenshotTestDB(t)
	defer cleanup()
	ctx := context.Background()

	t.Run("creates source with all fields", func(t *testing.T) {
		source := ScreenshotSource{
			Name:        "test-source",
			SourceName:  "Scene 1",
			CadenceMs:   3000,
			ImageFormat: "jpg",
			ImageWidth:  1920,
			ImageHeight: 1080,
			Quality:     90,
			Enabled:     true,
		}

		id, err := db.CreateScreenshotSource(ctx, source)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		// Verify the source was created
		got, err := db.GetScreenshotSource(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "test-source", got.Name)
		assert.Equal(t, "Scene 1", got.SourceName)
		assert.Equal(t, 3000, got.CadenceMs)
		assert.Equal(t, "jpg", got.ImageFormat)
		assert.Equal(t, 1920, got.ImageWidth)
		assert.Equal(t, 1080, got.ImageHeight)
		assert.Equal(t, 90, got.Quality)
		assert.True(t, got.Enabled)
	})

	t.Run("applies defaults for missing values", func(t *testing.T) {
		source := ScreenshotSource{
			Name:       "default-source",
			SourceName: "Scene 2",
			Enabled:    true,
		}

		id, err := db.CreateScreenshotSource(ctx, source)
		require.NoError(t, err)

		got, err := db.GetScreenshotSource(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, 5000, got.CadenceMs)    // Default
		assert.Equal(t, "png", got.ImageFormat) // Default
		assert.Equal(t, 80, got.Quality)        // Default
	})

	t.Run("rejects duplicate names", func(t *testing.T) {
		source := ScreenshotSource{
			Name:       "unique-name",
			SourceName: "Scene 3",
			Enabled:    true,
		}

		_, err := db.CreateScreenshotSource(ctx, source)
		require.NoError(t, err)

		// Try to create another with the same name
		_, err = db.CreateScreenshotSource(ctx, source)
		assert.Error(t, err)
		// Error message may vary by SQLite driver
		assert.True(t, err != nil, "expected an error for duplicate name")
	})
}

func TestGetScreenshotSource(t *testing.T) {
	db, cleanup := setupScreenshotTestDB(t)
	defer cleanup()
	ctx := context.Background()

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		_, err := db.GetScreenshotSource(ctx, 99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestGetScreenshotSourceByName(t *testing.T) {
	db, cleanup := setupScreenshotTestDB(t)
	defer cleanup()
	ctx := context.Background()

	source := ScreenshotSource{
		Name:       "named-source",
		SourceName: "Gaming Scene",
		Enabled:    true,
	}

	id, err := db.CreateScreenshotSource(ctx, source)
	require.NoError(t, err)

	t.Run("finds source by name", func(t *testing.T) {
		got, err := db.GetScreenshotSourceByName(ctx, "named-source")
		require.NoError(t, err)
		assert.Equal(t, id, got.ID)
		assert.Equal(t, "Gaming Scene", got.SourceName)
	})

	t.Run("returns error for non-existent name", func(t *testing.T) {
		_, err := db.GetScreenshotSourceByName(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestListScreenshotSources(t *testing.T) {
	db, cleanup := setupScreenshotTestDB(t)
	defer cleanup()
	ctx := context.Background()

	t.Run("returns empty list initially", func(t *testing.T) {
		sources, err := db.ListScreenshotSources(ctx)
		require.NoError(t, err)
		assert.Empty(t, sources)
	})

	t.Run("returns all created sources", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			source := ScreenshotSource{
				Name:       "source-" + string(rune('a'+i)),
				SourceName: "Scene " + string(rune('1'+i)),
				Enabled:    true,
			}
			_, err := db.CreateScreenshotSource(ctx, source)
			require.NoError(t, err)
		}

		sources, err := db.ListScreenshotSources(ctx)
		require.NoError(t, err)
		assert.Len(t, sources, 3)
	})
}

func TestUpdateScreenshotSource(t *testing.T) {
	db, cleanup := setupScreenshotTestDB(t)
	defer cleanup()
	ctx := context.Background()

	source := ScreenshotSource{
		Name:       "update-test",
		SourceName: "Original Scene",
		CadenceMs:  5000,
		Quality:    80,
		Enabled:    true,
	}

	id, err := db.CreateScreenshotSource(ctx, source)
	require.NoError(t, err)

	t.Run("updates source fields", func(t *testing.T) {
		updated := ScreenshotSource{
			ID:          id,
			SourceName:  "Updated Scene",
			CadenceMs:   10000,
			ImageFormat: "jpg",
			Quality:     50,
			Enabled:     false,
		}

		err := db.UpdateScreenshotSource(ctx, updated)
		require.NoError(t, err)

		got, err := db.GetScreenshotSource(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated Scene", got.SourceName)
		assert.Equal(t, 10000, got.CadenceMs)
		assert.Equal(t, "jpg", got.ImageFormat)
		assert.Equal(t, 50, got.Quality)
		assert.False(t, got.Enabled)
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		updated := ScreenshotSource{
			ID:         99999,
			SourceName: "Test",
		}

		err := db.UpdateScreenshotSource(ctx, updated)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestDeleteScreenshotSource(t *testing.T) {
	db, cleanup := setupScreenshotTestDB(t)
	defer cleanup()
	ctx := context.Background()

	source := ScreenshotSource{
		Name:       "delete-test",
		SourceName: "Scene to Delete",
		Enabled:    true,
	}

	id, err := db.CreateScreenshotSource(ctx, source)
	require.NoError(t, err)

	t.Run("deletes existing source", func(t *testing.T) {
		err := db.DeleteScreenshotSource(ctx, id)
		require.NoError(t, err)

		// Verify it's gone
		_, err = db.GetScreenshotSource(ctx, id)
		assert.Error(t, err)
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		err := db.DeleteScreenshotSource(ctx, 99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestSaveScreenshot(t *testing.T) {
	db, cleanup := setupScreenshotTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Create a source first
	source := ScreenshotSource{
		Name:       "screenshot-source",
		SourceName: "Test Scene",
		Enabled:    true,
	}
	sourceID, err := db.CreateScreenshotSource(ctx, source)
	require.NoError(t, err)

	t.Run("saves screenshot successfully", func(t *testing.T) {
		screenshot := Screenshot{
			SourceID:  sourceID,
			ImageData: "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
			MimeType:  "image/png",
			SizeBytes: 100,
		}

		id, err := db.SaveScreenshot(ctx, screenshot)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})
}

func TestGetLatestScreenshot(t *testing.T) {
	db, cleanup := setupScreenshotTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Create a source
	source := ScreenshotSource{
		Name:       "latest-source",
		SourceName: "Test Scene",
		Enabled:    true,
	}
	sourceID, err := db.CreateScreenshotSource(ctx, source)
	require.NoError(t, err)

	t.Run("returns error when no screenshots exist", func(t *testing.T) {
		_, err := db.GetLatestScreenshot(ctx, sourceID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no screenshots found")
	})

	t.Run("returns a screenshot after saving", func(t *testing.T) {
		// Save a screenshot
		screenshot := Screenshot{
			SourceID:  sourceID,
			ImageData: "test-image-data",
			MimeType:  "image/png",
			SizeBytes: 150,
		}
		_, err := db.SaveScreenshot(ctx, screenshot)
		require.NoError(t, err)

		got, err := db.GetLatestScreenshot(ctx, sourceID)
		require.NoError(t, err)
		assert.Equal(t, "test-image-data", got.ImageData)
		assert.Equal(t, "image/png", got.MimeType)
		assert.Equal(t, 150, got.SizeBytes)
		assert.Equal(t, sourceID, got.SourceID)
	})
}

func TestDeleteOldScreenshots(t *testing.T) {
	db, cleanup := setupScreenshotTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Create a source
	source := ScreenshotSource{
		Name:       "cleanup-source",
		SourceName: "Test Scene",
		Enabled:    true,
	}
	sourceID, err := db.CreateScreenshotSource(ctx, source)
	require.NoError(t, err)

	// Save 10 screenshots
	for i := 0; i < 10; i++ {
		screenshot := Screenshot{
			SourceID:  sourceID,
			ImageData: "data" + string(rune('0'+i)),
			MimeType:  "image/png",
			SizeBytes: 100,
		}
		_, err := db.SaveScreenshot(ctx, screenshot)
		require.NoError(t, err)
	}

	t.Run("keeps only specified count", func(t *testing.T) {
		count, err := db.CountScreenshots(ctx, sourceID)
		require.NoError(t, err)
		assert.Equal(t, int64(10), count)

		// Keep only 3
		deleted, err := db.DeleteOldScreenshots(ctx, sourceID, 3)
		require.NoError(t, err)
		assert.Equal(t, int64(7), deleted)

		count, err = db.CountScreenshots(ctx, sourceID)
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

func TestCountScreenshots(t *testing.T) {
	db, cleanup := setupScreenshotTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Create a source
	source := ScreenshotSource{
		Name:       "count-source",
		SourceName: "Test Scene",
		Enabled:    true,
	}
	sourceID, err := db.CreateScreenshotSource(ctx, source)
	require.NoError(t, err)

	t.Run("returns zero for empty source", func(t *testing.T) {
		count, err := db.CountScreenshots(ctx, sourceID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("counts screenshots correctly", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			screenshot := Screenshot{
				SourceID:  sourceID,
				ImageData: "test",
				MimeType:  "image/png",
				SizeBytes: 100,
			}
			_, err := db.SaveScreenshot(ctx, screenshot)
			require.NoError(t, err)
		}

		count, err := db.CountScreenshots(ctx, sourceID)
		require.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}

func TestCascadeDeleteScreenshots(t *testing.T) {
	db, cleanup := setupScreenshotTestDB(t)
	defer cleanup()
	ctx := context.Background()

	// Create a source
	source := ScreenshotSource{
		Name:       "cascade-source",
		SourceName: "Test Scene",
		Enabled:    true,
	}
	sourceID, err := db.CreateScreenshotSource(ctx, source)
	require.NoError(t, err)

	// Save screenshots
	for i := 0; i < 5; i++ {
		screenshot := Screenshot{
			SourceID:  sourceID,
			ImageData: "test",
			MimeType:  "image/png",
			SizeBytes: 100,
		}
		_, err := db.SaveScreenshot(ctx, screenshot)
		require.NoError(t, err)
	}

	// Verify screenshots exist
	count, err := db.CountScreenshots(ctx, sourceID)
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)

	// Delete the source
	err = db.DeleteScreenshotSource(ctx, sourceID)
	require.NoError(t, err)

	// Screenshots should be deleted via CASCADE
	count, err = db.CountScreenshots(ctx, sourceID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
