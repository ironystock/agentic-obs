package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ScreenshotSource represents a configured screenshot capture source.
// Sources define which OBS scene/source to capture and at what interval.
type ScreenshotSource struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`         // User-friendly unique name
	SourceName  string    `json:"source_name"`  // OBS scene or source name to capture
	CadenceMs   int       `json:"cadence_ms"`   // Capture interval in milliseconds
	ImageFormat string    `json:"image_format"` // "png" or "jpg"
	ImageWidth  int       `json:"image_width"`  // Optional resize width (0 = original)
	ImageHeight int       `json:"image_height"` // Optional resize height (0 = original)
	Quality     int       `json:"quality"`      // Compression quality 0-100
	Enabled     bool      `json:"enabled"`      // Whether capture is active
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Screenshot represents a captured screenshot image.
type Screenshot struct {
	ID         int64     `json:"id"`
	SourceID   int64     `json:"source_id"`  // FK to ScreenshotSource
	ImageData  string    `json:"image_data"` // Base64-encoded image data
	MimeType   string    `json:"mime_type"`  // "image/png" or "image/jpeg"
	CapturedAt time.Time `json:"captured_at"`
	SizeBytes  int       `json:"size_bytes"`
}

// CreateScreenshotSource creates a new screenshot source in the database.
// Returns the ID of the newly created source.
func (db *DB) CreateScreenshotSource(ctx context.Context, source ScreenshotSource) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Set defaults
	if source.CadenceMs <= 0 {
		source.CadenceMs = 5000
	}
	if source.ImageFormat == "" {
		source.ImageFormat = "png"
	}
	if source.Quality <= 0 {
		source.Quality = 80
	}

	result, err := db.conn.ExecContext(ctx, `
		INSERT INTO screenshot_sources (name, source_name, cadence_ms, image_format, image_width, image_height, quality, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, source.Name, source.SourceName, source.CadenceMs, source.ImageFormat, source.ImageWidth, source.ImageHeight, source.Quality, source.Enabled)

	if err != nil {
		// Check for unique constraint violation
		if err.Error() == "UNIQUE constraint failed: screenshot_sources.name" {
			return 0, fmt.Errorf("screenshot source with name '%s' already exists", source.Name)
		}
		return 0, fmt.Errorf("failed to create screenshot source: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get inserted source ID: %w", err)
	}

	return id, nil
}

// GetScreenshotSource retrieves a screenshot source by ID.
func (db *DB) GetScreenshotSource(ctx context.Context, id int64) (*ScreenshotSource, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var source ScreenshotSource
	var createdAt, updatedAt string
	var imageWidth, imageHeight sql.NullInt64

	err := db.conn.QueryRowContext(ctx, `
		SELECT id, name, source_name, cadence_ms, image_format, image_width, image_height, quality, enabled, created_at, updated_at
		FROM screenshot_sources
		WHERE id = ?
	`, id).Scan(&source.ID, &source.Name, &source.SourceName, &source.CadenceMs, &source.ImageFormat,
		&imageWidth, &imageHeight, &source.Quality, &source.Enabled, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("screenshot source with ID %d not found", id)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get screenshot source with ID %d: %w", id, err)
	}

	// Handle nullable width/height
	if imageWidth.Valid {
		source.ImageWidth = int(imageWidth.Int64)
	}
	if imageHeight.Valid {
		source.ImageHeight = int(imageHeight.Int64)
	}

	// Parse timestamps
	source.CreatedAt, _ = parseTimestamp(createdAt)
	source.UpdatedAt, _ = parseTimestamp(updatedAt)

	return &source, nil
}

// GetScreenshotSourceByName retrieves a screenshot source by name.
func (db *DB) GetScreenshotSourceByName(ctx context.Context, name string) (*ScreenshotSource, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var source ScreenshotSource
	var createdAt, updatedAt string
	var imageWidth, imageHeight sql.NullInt64

	err := db.conn.QueryRowContext(ctx, `
		SELECT id, name, source_name, cadence_ms, image_format, image_width, image_height, quality, enabled, created_at, updated_at
		FROM screenshot_sources
		WHERE name = ?
	`, name).Scan(&source.ID, &source.Name, &source.SourceName, &source.CadenceMs, &source.ImageFormat,
		&imageWidth, &imageHeight, &source.Quality, &source.Enabled, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("screenshot source '%s' not found", name)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get screenshot source '%s': %w", name, err)
	}

	// Handle nullable width/height
	if imageWidth.Valid {
		source.ImageWidth = int(imageWidth.Int64)
	}
	if imageHeight.Valid {
		source.ImageHeight = int(imageHeight.Int64)
	}

	// Parse timestamps
	source.CreatedAt, _ = parseTimestamp(createdAt)
	source.UpdatedAt, _ = parseTimestamp(updatedAt)

	return &source, nil
}

// ListScreenshotSources retrieves all screenshot sources.
func (db *DB) ListScreenshotSources(ctx context.Context) ([]*ScreenshotSource, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	rows, err := db.conn.QueryContext(ctx, `
		SELECT id, name, source_name, cadence_ms, image_format, image_width, image_height, quality, enabled, created_at, updated_at
		FROM screenshot_sources
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list screenshot sources: %w", err)
	}
	defer rows.Close()

	var sources []*ScreenshotSource
	for rows.Next() {
		var source ScreenshotSource
		var createdAt, updatedAt string
		var imageWidth, imageHeight sql.NullInt64

		if err := rows.Scan(&source.ID, &source.Name, &source.SourceName, &source.CadenceMs, &source.ImageFormat,
			&imageWidth, &imageHeight, &source.Quality, &source.Enabled, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan screenshot source row: %w", err)
		}

		// Handle nullable width/height
		if imageWidth.Valid {
			source.ImageWidth = int(imageWidth.Int64)
		}
		if imageHeight.Valid {
			source.ImageHeight = int(imageHeight.Int64)
		}

		// Parse timestamps
		source.CreatedAt, _ = parseTimestamp(createdAt)
		source.UpdatedAt, _ = parseTimestamp(updatedAt)

		sources = append(sources, &source)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating screenshot sources: %w", err)
	}

	return sources, nil
}

// UpdateScreenshotSource updates an existing screenshot source.
func (db *DB) UpdateScreenshotSource(ctx context.Context, source ScreenshotSource) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result, err := db.conn.ExecContext(ctx, `
		UPDATE screenshot_sources
		SET source_name = ?, cadence_ms = ?, image_format = ?, image_width = ?, image_height = ?, quality = ?, enabled = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, source.SourceName, source.CadenceMs, source.ImageFormat, source.ImageWidth, source.ImageHeight, source.Quality, source.Enabled, source.ID)

	if err != nil {
		return fmt.Errorf("failed to update screenshot source: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("screenshot source with ID %d not found", source.ID)
	}

	return nil
}

// DeleteScreenshotSource deletes a screenshot source and all its screenshots (via CASCADE).
func (db *DB) DeleteScreenshotSource(ctx context.Context, id int64) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result, err := db.conn.ExecContext(ctx, `DELETE FROM screenshot_sources WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete screenshot source: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("screenshot source with ID %d not found", id)
	}

	return nil
}

// SaveScreenshot saves a captured screenshot to the database.
// Returns the ID of the newly created screenshot.
func (db *DB) SaveScreenshot(ctx context.Context, screenshot Screenshot) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result, err := db.conn.ExecContext(ctx, `
		INSERT INTO screenshots (source_id, image_data, mime_type, captured_at, size_bytes)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, ?)
	`, screenshot.SourceID, screenshot.ImageData, screenshot.MimeType, screenshot.SizeBytes)

	if err != nil {
		return 0, fmt.Errorf("failed to save screenshot: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get inserted screenshot ID: %w", err)
	}

	return id, nil
}

// GetLatestScreenshot retrieves the most recent screenshot for a source.
func (db *DB) GetLatestScreenshot(ctx context.Context, sourceID int64) (*Screenshot, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var screenshot Screenshot
	var capturedAt string

	err := db.conn.QueryRowContext(ctx, `
		SELECT id, source_id, image_data, mime_type, captured_at, size_bytes
		FROM screenshots
		WHERE source_id = ?
		ORDER BY captured_at DESC
		LIMIT 1
	`, sourceID).Scan(&screenshot.ID, &screenshot.SourceID, &screenshot.ImageData, &screenshot.MimeType, &capturedAt, &screenshot.SizeBytes)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no screenshots found for source ID %d", sourceID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get latest screenshot for source ID %d: %w", sourceID, err)
	}

	screenshot.CapturedAt, _ = parseTimestamp(capturedAt)

	return &screenshot, nil
}

// DeleteOldScreenshots deletes old screenshots, keeping only the most recent keepCount.
// Returns the number of deleted screenshots.
func (db *DB) DeleteOldScreenshots(ctx context.Context, sourceID int64, keepCount int) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Delete screenshots older than the Nth most recent
	result, err := db.conn.ExecContext(ctx, `
		DELETE FROM screenshots
		WHERE source_id = ? AND id NOT IN (
			SELECT id FROM screenshots
			WHERE source_id = ?
			ORDER BY captured_at DESC
			LIMIT ?
		)
	`, sourceID, sourceID, keepCount)

	if err != nil {
		return 0, fmt.Errorf("failed to delete old screenshots: %w", err)
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows deleted: %w", err)
	}

	return rowsDeleted, nil
}

// CountScreenshots returns the number of screenshots for a source.
func (db *DB) CountScreenshots(ctx context.Context, sourceID int64) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var count int64
	err := db.conn.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM screenshots WHERE source_id = ?
	`, sourceID).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to count screenshots: %w", err)
	}

	return count, nil
}
