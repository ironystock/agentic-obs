package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordAction(t *testing.T) {
	t.Run("records action successfully", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		record := ActionRecord{
			Action:     "Test action",
			ToolName:   "test_tool",
			Input:      `{"key": "value"}`,
			Output:     `{"result": "success"}`,
			Success:    true,
			DurationMs: 150,
		}

		id, err := db.RecordAction(context.Background(), record)
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})

	t.Run("records failed action", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		record := ActionRecord{
			Action:     "Failed action",
			ToolName:   "failing_tool",
			Input:      `{"bad": "input"}`,
			Output:     `{"error": "something went wrong"}`,
			Success:    false,
			DurationMs: 50,
		}

		id, err := db.RecordAction(context.Background(), record)
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})

	t.Run("records action without optional fields", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		record := ActionRecord{
			Action:  "Simple action",
			Success: true,
		}

		id, err := db.RecordAction(context.Background(), record)
		assert.NoError(t, err)
		assert.Greater(t, id, int64(0))
	})

	t.Run("records multiple actions with sequential IDs", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		id1, err := db.RecordAction(context.Background(), ActionRecord{Action: "First", Success: true})
		require.NoError(t, err)

		id2, err := db.RecordAction(context.Background(), ActionRecord{Action: "Second", Success: true})
		require.NoError(t, err)

		assert.Equal(t, id1+1, id2)
	})
}

func TestGetRecentActions(t *testing.T) {
	t.Run("returns empty list when no actions", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		actions, err := db.GetRecentActions(context.Background(), 10)
		assert.NoError(t, err)
		assert.Empty(t, actions)
	})

	t.Run("returns actions in descending order", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Record actions
		_, err := db.RecordAction(context.Background(), ActionRecord{Action: "First", Success: true})
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		_, err = db.RecordAction(context.Background(), ActionRecord{Action: "Second", Success: true})
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		_, err = db.RecordAction(context.Background(), ActionRecord{Action: "Third", Success: true})
		require.NoError(t, err)

		actions, err := db.GetRecentActions(context.Background(), 10)
		assert.NoError(t, err)
		assert.Len(t, actions, 3)
		assert.Equal(t, "Third", actions[0].Action)
		assert.Equal(t, "Second", actions[1].Action)
		assert.Equal(t, "First", actions[2].Action)
	})

	t.Run("respects limit parameter", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Record 5 actions
		for i := 0; i < 5; i++ {
			_, err := db.RecordAction(context.Background(), ActionRecord{Action: "Action", Success: true})
			require.NoError(t, err)
		}

		actions, err := db.GetRecentActions(context.Background(), 3)
		assert.NoError(t, err)
		assert.Len(t, actions, 3)
	})

	t.Run("uses default limit when zero", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		_, err := db.RecordAction(context.Background(), ActionRecord{Action: "Test", Success: true})
		require.NoError(t, err)

		actions, err := db.GetRecentActions(context.Background(), 0)
		assert.NoError(t, err)
		assert.Len(t, actions, 1)
	})

	t.Run("returns all action fields", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		record := ActionRecord{
			Action:     "Full action",
			ToolName:   "my_tool",
			Input:      `{"input": true}`,
			Output:     `{"output": true}`,
			Success:    true,
			DurationMs: 200,
		}
		_, err := db.RecordAction(context.Background(), record)
		require.NoError(t, err)

		actions, err := db.GetRecentActions(context.Background(), 1)
		assert.NoError(t, err)
		require.Len(t, actions, 1)

		retrieved := actions[0]
		assert.Equal(t, "Full action", retrieved.Action)
		assert.Equal(t, "my_tool", retrieved.ToolName)
		assert.Equal(t, `{"input": true}`, retrieved.Input)
		assert.Equal(t, `{"output": true}`, retrieved.Output)
		assert.True(t, retrieved.Success)
		assert.Equal(t, int64(200), retrieved.DurationMs)
		assert.False(t, retrieved.CreatedAt.IsZero())
	})
}

func TestGetActionsByTool(t *testing.T) {
	t.Run("filters by tool name", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Record actions with different tools
		_, err := db.RecordAction(context.Background(), ActionRecord{Action: "A1", ToolName: "tool_a", Success: true})
		require.NoError(t, err)
		_, err = db.RecordAction(context.Background(), ActionRecord{Action: "B1", ToolName: "tool_b", Success: true})
		require.NoError(t, err)
		_, err = db.RecordAction(context.Background(), ActionRecord{Action: "A2", ToolName: "tool_a", Success: true})
		require.NoError(t, err)

		actions, err := db.GetActionsByTool(context.Background(), "tool_a", 10)
		assert.NoError(t, err)
		assert.Len(t, actions, 2)
		for _, a := range actions {
			assert.Equal(t, "tool_a", a.ToolName)
		}
	})

	t.Run("returns empty for non-existent tool", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		_, err := db.RecordAction(context.Background(), ActionRecord{Action: "Test", ToolName: "existing", Success: true})
		require.NoError(t, err)

		actions, err := db.GetActionsByTool(context.Background(), "nonexistent", 10)
		assert.NoError(t, err)
		assert.Empty(t, actions)
	})

	t.Run("respects limit parameter", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Record 5 actions with same tool
		for i := 0; i < 5; i++ {
			_, err := db.RecordAction(context.Background(), ActionRecord{Action: "Test", ToolName: "my_tool", Success: true})
			require.NoError(t, err)
		}

		actions, err := db.GetActionsByTool(context.Background(), "my_tool", 2)
		assert.NoError(t, err)
		assert.Len(t, actions, 2)
	})
}

func TestGetActionsSince(t *testing.T) {
	t.Run("returns actions since timestamp", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Record an action
		_, err := db.RecordAction(context.Background(), ActionRecord{Action: "Recent", Success: true})
		require.NoError(t, err)

		// Get actions since 1 minute ago
		since := time.Now().Add(-time.Minute)
		actions, err := db.GetActionsSince(context.Background(), since, 10)
		assert.NoError(t, err)
		assert.Len(t, actions, 1)
	})

	t.Run("excludes older actions", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		_, err := db.RecordAction(context.Background(), ActionRecord{Action: "Test", Success: true})
		require.NoError(t, err)

		// Get actions since the future - should return nothing
		since := time.Now().Add(time.Hour)
		actions, err := db.GetActionsSince(context.Background(), since, 10)
		assert.NoError(t, err)
		assert.Empty(t, actions)
	})
}

func TestGetActionStats(t *testing.T) {
	t.Run("returns stats for empty database", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		stats, err := db.GetActionStats(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, int64(0), stats["total_actions"])
		assert.Equal(t, int64(0), stats["successful_actions"])
		assert.Equal(t, int64(0), stats["failed_actions"])
	})

	t.Run("counts total and success/failure", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Record successful actions
		_, err := db.RecordAction(context.Background(), ActionRecord{Action: "S1", Success: true})
		require.NoError(t, err)
		_, err = db.RecordAction(context.Background(), ActionRecord{Action: "S2", Success: true})
		require.NoError(t, err)

		// Record failed action
		_, err = db.RecordAction(context.Background(), ActionRecord{Action: "F1", Success: false})
		require.NoError(t, err)

		stats, err := db.GetActionStats(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, int64(3), stats["total_actions"])
		assert.Equal(t, int64(2), stats["successful_actions"])
		assert.Equal(t, int64(1), stats["failed_actions"])
	})

	t.Run("calculates average duration", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		_, err := db.RecordAction(context.Background(), ActionRecord{Action: "A", DurationMs: 100, Success: true})
		require.NoError(t, err)
		_, err = db.RecordAction(context.Background(), ActionRecord{Action: "B", DurationMs: 200, Success: true})
		require.NoError(t, err)

		stats, err := db.GetActionStats(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, float64(150), stats["avg_duration_ms"])
	})

	t.Run("returns top tools", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Record actions with different tools
		for i := 0; i < 3; i++ {
			_, err := db.RecordAction(context.Background(), ActionRecord{Action: "A", ToolName: "tool_a", Success: true})
			require.NoError(t, err)
		}
		for i := 0; i < 2; i++ {
			_, err := db.RecordAction(context.Background(), ActionRecord{Action: "B", ToolName: "tool_b", Success: true})
			require.NoError(t, err)
		}
		_, err := db.RecordAction(context.Background(), ActionRecord{Action: "C", ToolName: "tool_c", Success: true})
		require.NoError(t, err)

		stats, err := db.GetActionStats(context.Background())
		assert.NoError(t, err)

		topTools := stats["top_tools"].([]map[string]interface{})
		assert.Len(t, topTools, 3)
		assert.Equal(t, "tool_a", topTools[0]["tool_name"])
		assert.Equal(t, int64(3), topTools[0]["count"])
	})
}

func TestClearOldActions(t *testing.T) {
	t.Run("deletes old actions", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Record an action
		_, err := db.RecordAction(context.Background(), ActionRecord{Action: "Old", Success: true})
		require.NoError(t, err)

		// Clear actions older than 0 (all of them)
		deleted, err := db.ClearOldActions(context.Background(), 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), deleted)

		// Verify deletion
		actions, err := db.GetRecentActions(context.Background(), 10)
		assert.NoError(t, err)
		assert.Empty(t, actions)
	})

	t.Run("keeps recent actions", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		// Record an action
		_, err := db.RecordAction(context.Background(), ActionRecord{Action: "Recent", Success: true})
		require.NoError(t, err)

		// Clear actions older than 1 hour (should keep all)
		deleted, err := db.ClearOldActions(context.Background(), time.Hour)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), deleted)

		// Verify no deletion
		actions, err := db.GetRecentActions(context.Background(), 10)
		assert.NoError(t, err)
		assert.Len(t, actions, 1)
	})

	t.Run("returns zero when no actions to delete", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		deleted, err := db.ClearOldActions(context.Background(), time.Hour)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), deleted)
	})
}

func TestActionRecordFields(t *testing.T) {
	t.Run("handles nil/empty optional fields", func(t *testing.T) {
		db, cleanup := testDB(t)
		defer cleanup()

		record := ActionRecord{
			Action:  "Minimal",
			Success: true,
		}
		_, err := db.RecordAction(context.Background(), record)
		require.NoError(t, err)

		actions, err := db.GetRecentActions(context.Background(), 1)
		assert.NoError(t, err)
		require.Len(t, actions, 1)

		retrieved := actions[0]
		assert.Equal(t, "Minimal", retrieved.Action)
		assert.Empty(t, retrieved.ToolName)
		assert.Empty(t, retrieved.Input)
		assert.Empty(t, retrieved.Output)
		assert.True(t, retrieved.Success)
		assert.Equal(t, int64(0), retrieved.DurationMs)
	})
}
