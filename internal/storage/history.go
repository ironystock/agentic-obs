package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ActionRecord represents a single action in the history log.
type ActionRecord struct {
	ID         int64     `json:"id"`
	Action     string    `json:"action"`
	ToolName   string    `json:"tool_name,omitempty"`
	Input      string    `json:"input,omitempty"`
	Output     string    `json:"output,omitempty"`
	Success    bool      `json:"success"`
	DurationMs int64     `json:"duration_ms,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// RecordAction adds a new action to the history log.
func (db *DB) RecordAction(ctx context.Context, record ActionRecord) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	successInt := 0
	if record.Success {
		successInt = 1
	}

	result, err := db.conn.ExecContext(ctx,
		`INSERT INTO action_history (action, tool_name, input, output, success, duration_ms, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		record.Action,
		record.ToolName,
		record.Input,
		record.Output,
		successInt,
		record.DurationMs,
		time.Now(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to record action: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get action ID: %w", err)
	}

	return id, nil
}

// GetRecentActions retrieves the most recent actions from history.
// limit specifies maximum number of records to return (0 = use default of 100).
func (db *DB) GetRecentActions(ctx context.Context, limit int) ([]ActionRecord, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if limit <= 0 {
		limit = 100
	}

	rows, err := db.conn.QueryContext(ctx,
		`SELECT id, action, tool_name, input, output, success, duration_ms, created_at
		 FROM action_history
		 ORDER BY created_at DESC
		 LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query action history: %w", err)
	}
	defer rows.Close()

	return scanActionRecords(rows)
}

// GetActionsByTool retrieves actions filtered by tool name.
func (db *DB) GetActionsByTool(ctx context.Context, toolName string, limit int) ([]ActionRecord, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if limit <= 0 {
		limit = 100
	}

	rows, err := db.conn.QueryContext(ctx,
		`SELECT id, action, tool_name, input, output, success, duration_ms, created_at
		 FROM action_history
		 WHERE tool_name = ?
		 ORDER BY created_at DESC
		 LIMIT ?`,
		toolName,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query actions by tool: %w", err)
	}
	defer rows.Close()

	return scanActionRecords(rows)
}

// GetActionsSince retrieves actions since a given time.
func (db *DB) GetActionsSince(ctx context.Context, since time.Time, limit int) ([]ActionRecord, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if limit <= 0 {
		limit = 100
	}

	rows, err := db.conn.QueryContext(ctx,
		`SELECT id, action, tool_name, input, output, success, duration_ms, created_at
		 FROM action_history
		 WHERE created_at >= ?
		 ORDER BY created_at DESC
		 LIMIT ?`,
		since,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query actions since time: %w", err)
	}
	defer rows.Close()

	return scanActionRecords(rows)
}

// GetActionStats returns statistics about recorded actions.
func (db *DB) GetActionStats(ctx context.Context) (map[string]interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	stats := make(map[string]interface{})

	// Total count
	var total int64
	err := db.conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM action_history").Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count actions: %w", err)
	}
	stats["total_actions"] = total

	// Success count
	var successCount int64
	err = db.conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM action_history WHERE success = 1").Scan(&successCount)
	if err != nil {
		return nil, fmt.Errorf("failed to count successful actions: %w", err)
	}
	stats["successful_actions"] = successCount
	stats["failed_actions"] = total - successCount

	// Average duration
	var avgDuration sql.NullFloat64
	err = db.conn.QueryRowContext(ctx, "SELECT AVG(duration_ms) FROM action_history WHERE duration_ms > 0").Scan(&avgDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate average duration: %w", err)
	}
	if avgDuration.Valid {
		stats["avg_duration_ms"] = avgDuration.Float64
	} else {
		stats["avg_duration_ms"] = 0
	}

	// Most used tools (top 5)
	rows, err := db.conn.QueryContext(ctx,
		`SELECT tool_name, COUNT(*) as count
		 FROM action_history
		 WHERE tool_name IS NOT NULL AND tool_name != ''
		 GROUP BY tool_name
		 ORDER BY count DESC
		 LIMIT 5`)
	if err != nil {
		return nil, fmt.Errorf("failed to query top tools: %w", err)
	}
	defer rows.Close()

	topTools := make([]map[string]interface{}, 0)
	for rows.Next() {
		var toolName string
		var count int64
		if err := rows.Scan(&toolName, &count); err != nil {
			continue
		}
		topTools = append(topTools, map[string]interface{}{
			"tool_name": toolName,
			"count":     count,
		})
	}
	stats["top_tools"] = topTools

	return stats, nil
}

// ClearOldActions removes actions older than the specified duration.
// Returns the number of deleted records.
func (db *DB) ClearOldActions(ctx context.Context, olderThan time.Duration) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	cutoff := time.Now().Add(-olderThan)

	result, err := db.conn.ExecContext(ctx,
		"DELETE FROM action_history WHERE created_at < ?",
		cutoff,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old actions: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get deleted count: %w", err)
	}

	return count, nil
}

// scanActionRecords helper to scan rows into ActionRecord slice.
func scanActionRecords(rows *sql.Rows) ([]ActionRecord, error) {
	var records []ActionRecord

	for rows.Next() {
		var r ActionRecord
		var toolName, input, output sql.NullString
		var durationMs sql.NullInt64
		var success int
		var createdAt string

		err := rows.Scan(
			&r.ID,
			&r.Action,
			&toolName,
			&input,
			&output,
			&success,
			&durationMs,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan action record: %w", err)
		}

		r.ToolName = toolName.String
		r.Input = input.String
		r.Output = output.String
		r.Success = success == 1
		if durationMs.Valid {
			r.DurationMs = durationMs.Int64
		}

		// Parse timestamp - handle multiple SQLite formats
		parsedTime, err := parseTimestamp(createdAt)
		if err != nil {
			r.CreatedAt = time.Now() // Fallback
		} else {
			r.CreatedAt = parsedTime
		}

		records = append(records, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating action records: %w", err)
	}

	return records, nil
}
