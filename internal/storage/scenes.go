package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// parseTimestamp parses a timestamp string from SQLite.
// SQLite can return timestamps in various formats depending on the driver.
func parseTimestamp(s string) (time.Time, error) {
	// Try common timestamp formats
	formats := []string{
		time.RFC3339,           // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05Z", // ISO 8601 with Z
		"2006-01-02 15:04:05",  // SQLite default
		time.DateTime,          // "2006-01-02 15:04:05"
	}

	// Normalize the timestamp - remove trailing Z if present for some formats
	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", s)
}

// ScenePreset represents a user-defined OBS scene preset.
// Presets allow users to save and restore specific scene configurations,
// including which sources are visible and their settings.
type ScenePreset struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	SceneName string        `json:"scene_name"`
	Sources   []SourceState `json:"sources,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
}

// SourceState represents the state of a source within a scene preset.
type SourceState struct {
	Name     string                 `json:"name"`               // Source name
	Visible  bool                   `json:"visible"`            // Whether the source is visible
	Settings map[string]interface{} `json:"settings,omitempty"` // Source-specific settings (optional)
}

// CreateScenePreset creates a new scene preset in the database.
// Returns the ID of the newly created preset.
func (db *DB) CreateScenePreset(ctx context.Context, preset ScenePreset) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Serialize sources to JSON
	sourcesJSON, err := json.Marshal(preset.Sources)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize sources to JSON: %w", err)
	}

	// Insert preset
	result, err := db.conn.ExecContext(ctx, `
		INSERT INTO scene_presets (name, scene_name, sources, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
	`, preset.Name, preset.SceneName, string(sourcesJSON))

	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "UNIQUE constraint failed: scene_presets.name" {
			return 0, fmt.Errorf("scene preset with name '%s' already exists", preset.Name)
		}
		return 0, fmt.Errorf("failed to create scene preset: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get inserted preset ID: %w", err)
	}

	return id, nil
}

// GetScenePreset retrieves a scene preset by name.
// Returns sql.ErrNoRows if the preset doesn't exist.
func (db *DB) GetScenePreset(ctx context.Context, name string) (*ScenePreset, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var preset ScenePreset
	var sourcesJSON string
	var createdAt string

	err := db.conn.QueryRowContext(ctx, `
		SELECT id, name, scene_name, sources, created_at
		FROM scene_presets
		WHERE name = ?
	`, name).Scan(&preset.ID, &preset.Name, &preset.SceneName, &sourcesJSON, &createdAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("scene preset '%s' not found", name)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get scene preset '%s': %w", name, err)
	}

	// Parse sources JSON
	if sourcesJSON != "" {
		if err := json.Unmarshal([]byte(sourcesJSON), &preset.Sources); err != nil {
			return nil, fmt.Errorf("failed to parse sources JSON for preset '%s': %w", name, err)
		}
	}

	// Parse created_at timestamp
	preset.CreatedAt, err = parseTimestamp(createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at timestamp for preset '%s': %w", name, err)
	}

	return &preset, nil
}

// GetScenePresetByID retrieves a scene preset by its ID.
// Returns sql.ErrNoRows if the preset doesn't exist.
func (db *DB) GetScenePresetByID(ctx context.Context, id int64) (*ScenePreset, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var preset ScenePreset
	var sourcesJSON string
	var createdAt string

	err := db.conn.QueryRowContext(ctx, `
		SELECT id, name, scene_name, sources, created_at
		FROM scene_presets
		WHERE id = ?
	`, id).Scan(&preset.ID, &preset.Name, &preset.SceneName, &sourcesJSON, &createdAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("scene preset with ID %d not found", id)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get scene preset with ID %d: %w", id, err)
	}

	// Parse sources JSON
	if sourcesJSON != "" {
		if err := json.Unmarshal([]byte(sourcesJSON), &preset.Sources); err != nil {
			return nil, fmt.Errorf("failed to parse sources JSON for preset ID %d: %w", id, err)
		}
	}

	// Parse created_at timestamp
	preset.CreatedAt, err = parseTimestamp(createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at timestamp for preset ID %d: %w", id, err)
	}

	return &preset, nil
}

// ListScenePresets returns all scene presets, optionally filtered by scene name.
// If sceneName is empty, returns all presets.
func (db *DB) ListScenePresets(ctx context.Context, sceneName string) ([]*ScenePreset, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var rows *sql.Rows
	var err error

	if sceneName != "" {
		// Filter by scene name
		rows, err = db.conn.QueryContext(ctx, `
			SELECT id, name, scene_name, sources, created_at
			FROM scene_presets
			WHERE scene_name = ?
			ORDER BY created_at DESC
		`, sceneName)
	} else {
		// Return all presets
		rows, err = db.conn.QueryContext(ctx, `
			SELECT id, name, scene_name, sources, created_at
			FROM scene_presets
			ORDER BY created_at DESC
		`)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list scene presets: %w", err)
	}
	defer rows.Close()

	var presets []*ScenePreset
	for rows.Next() {
		var preset ScenePreset
		var sourcesJSON string
		var createdAt string

		if err := rows.Scan(&preset.ID, &preset.Name, &preset.SceneName, &sourcesJSON, &createdAt); err != nil {
			return nil, fmt.Errorf("failed to scan scene preset row: %w", err)
		}

		// Parse sources JSON
		if sourcesJSON != "" {
			if err := json.Unmarshal([]byte(sourcesJSON), &preset.Sources); err != nil {
				return nil, fmt.Errorf("failed to parse sources JSON for preset '%s': %w", preset.Name, err)
			}
		}

		// Parse created_at timestamp
		preset.CreatedAt, err = parseTimestamp(createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at timestamp for preset '%s': %w", preset.Name, err)
		}

		presets = append(presets, &preset)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating scene preset rows: %w", err)
	}

	return presets, nil
}

// UpdateScenePreset updates an existing scene preset.
// The preset is identified by its name.
func (db *DB) UpdateScenePreset(ctx context.Context, preset ScenePreset) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Serialize sources to JSON
	sourcesJSON, err := json.Marshal(preset.Sources)
	if err != nil {
		return fmt.Errorf("failed to serialize sources to JSON: %w", err)
	}

	// Update preset
	result, err := db.conn.ExecContext(ctx, `
		UPDATE scene_presets
		SET scene_name = ?, sources = ?
		WHERE name = ?
	`, preset.SceneName, string(sourcesJSON), preset.Name)

	if err != nil {
		return fmt.Errorf("failed to update scene preset '%s': %w", preset.Name, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result for preset '%s': %w", preset.Name, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("scene preset '%s' not found", preset.Name)
	}

	return nil
}

// DeleteScenePreset deletes a scene preset by name.
func (db *DB) DeleteScenePreset(ctx context.Context, name string) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result, err := db.conn.ExecContext(ctx,
		"DELETE FROM scene_presets WHERE name = ?",
		name,
	)

	if err != nil {
		return fmt.Errorf("failed to delete scene preset '%s': %w", name, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result for preset '%s': %w", name, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("scene preset '%s' not found", name)
	}

	return nil
}

// DeleteScenePresetByID deletes a scene preset by ID.
func (db *DB) DeleteScenePresetByID(ctx context.Context, id int64) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result, err := db.conn.ExecContext(ctx,
		"DELETE FROM scene_presets WHERE id = ?",
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete scene preset with ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result for preset ID %d: %w", id, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("scene preset with ID %d not found", id)
	}

	return nil
}

// DeleteScenePresetsByScene deletes all presets associated with a specific scene.
// This is useful when a scene is deleted from OBS.
func (db *DB) DeleteScenePresetsByScene(ctx context.Context, sceneName string) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result, err := db.conn.ExecContext(ctx,
		"DELETE FROM scene_presets WHERE scene_name = ?",
		sceneName,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to delete presets for scene '%s': %w", sceneName, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to check delete result for scene '%s': %w", sceneName, err)
	}

	return rowsAffected, nil
}

// CountScenePresets returns the total number of scene presets in the database.
func (db *DB) CountScenePresets(ctx context.Context) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var count int64
	err := db.conn.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM scene_presets",
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to count scene presets: %w", err)
	}

	return count, nil
}

// RenameScenePreset renames a scene preset.
func (db *DB) RenameScenePreset(ctx context.Context, oldName, newName string) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result, err := db.conn.ExecContext(ctx,
		"UPDATE scene_presets SET name = ? WHERE name = ?",
		newName, oldName,
	)

	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "UNIQUE constraint failed: scene_presets.name" {
			return fmt.Errorf("scene preset with name '%s' already exists", newName)
		}
		return fmt.Errorf("failed to rename scene preset from '%s' to '%s': %w", oldName, newName, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rename result: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("scene preset '%s' not found", oldName)
	}

	return nil
}
