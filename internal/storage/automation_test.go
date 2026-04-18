package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAutomationRule(t *testing.T) {
	db, cleanup := testDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("creates rule successfully", func(t *testing.T) {
		rule := AutomationRule{
			Name:        "test-rule",
			Description: "A test automation rule",
			Enabled:     true,
			TriggerType: TriggerTypeEvent,
			TriggerConfig: map[string]interface{}{
				"event_type": "scene_changed",
			},
			Actions: []RuleAction{
				{
					Type: "set_mute",
					Parameters: map[string]interface{}{
						"input_name": "Microphone",
						"muted":      true,
					},
				},
			},
			CooldownMs: 5000,
			Priority:   10,
		}

		id, err := db.CreateAutomationRule(ctx, rule)
		require.NoError(t, err)
		assert.Greater(t, id, int64(0))

		// Verify rule was created
		saved, err := db.GetAutomationRule(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "test-rule", saved.Name)
		assert.Equal(t, "A test automation rule", saved.Description)
		assert.True(t, saved.Enabled)
		assert.Equal(t, TriggerTypeEvent, saved.TriggerType)
		assert.Equal(t, "scene_changed", saved.TriggerConfig["event_type"])
		assert.Len(t, saved.Actions, 1)
		assert.Equal(t, "set_mute", saved.Actions[0].Type)
		assert.Equal(t, 5000, saved.CooldownMs)
		assert.Equal(t, 10, saved.Priority)
	})

	t.Run("fails on duplicate name", func(t *testing.T) {
		rule := AutomationRule{
			Name:        "duplicate-test",
			TriggerType: TriggerTypeEvent,
			TriggerConfig: map[string]interface{}{
				"event_type": "recording_started",
			},
			Actions: []RuleAction{},
		}

		_, err := db.CreateAutomationRule(ctx, rule)
		require.NoError(t, err)

		// Try to create another with same name
		_, err = db.CreateAutomationRule(ctx, rule)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("creates schedule rule", func(t *testing.T) {
		rule := AutomationRule{
			Name:        "hourly-check",
			TriggerType: TriggerTypeSchedule,
			TriggerConfig: map[string]interface{}{
				"schedule": "0 * * * *",
			},
			Actions: []RuleAction{
				{Type: "set_scene", Parameters: map[string]interface{}{"scene_name": "Break"}},
			},
		}

		id, err := db.CreateAutomationRule(ctx, rule)
		require.NoError(t, err)

		saved, err := db.GetAutomationRule(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, TriggerTypeSchedule, saved.TriggerType)
		assert.Equal(t, "0 * * * *", saved.TriggerConfig["schedule"])
	})
}

func TestGetAutomationRuleByName(t *testing.T) {
	db, cleanup := testDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("retrieves existing rule", func(t *testing.T) {
		rule := AutomationRule{
			Name:        "find-me",
			Description: "Find this rule",
			Enabled:     true,
			TriggerType: TriggerTypeEvent,
			TriggerConfig: map[string]interface{}{
				"event_type": "streaming_started",
			},
			Actions: []RuleAction{
				{Type: "set_scene", Parameters: map[string]interface{}{"scene_name": "Live"}},
			},
		}

		_, err := db.CreateAutomationRule(ctx, rule)
		require.NoError(t, err)

		found, err := db.GetAutomationRuleByName(ctx, "find-me")
		require.NoError(t, err)
		assert.Equal(t, "find-me", found.Name)
		assert.Equal(t, "streaming_started", found.TriggerConfig["event_type"])
	})

	t.Run("returns error for non-existent rule", func(t *testing.T) {
		_, err := db.GetAutomationRuleByName(ctx, "non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestListAutomationRules(t *testing.T) {
	db, cleanup := testDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create test rules
	rules := []AutomationRule{
		{Name: "rule-1", Enabled: true, TriggerType: TriggerTypeEvent, TriggerConfig: map[string]interface{}{"event_type": "e1"}, Actions: []RuleAction{}, Priority: 5},
		{Name: "rule-2", Enabled: false, TriggerType: TriggerTypeEvent, TriggerConfig: map[string]interface{}{"event_type": "e2"}, Actions: []RuleAction{}, Priority: 10},
		{Name: "rule-3", Enabled: true, TriggerType: TriggerTypeSchedule, TriggerConfig: map[string]interface{}{"schedule": "* * * * *"}, Actions: []RuleAction{}, Priority: 1},
	}

	for _, rule := range rules {
		_, err := db.CreateAutomationRule(ctx, rule)
		require.NoError(t, err)
	}

	t.Run("lists all rules", func(t *testing.T) {
		listed, err := db.ListAutomationRules(ctx, false)
		require.NoError(t, err)
		assert.Len(t, listed, 3)

		// Should be sorted by priority DESC
		assert.Equal(t, "rule-2", listed[0].Name) // Priority 10
		assert.Equal(t, "rule-1", listed[1].Name) // Priority 5
		assert.Equal(t, "rule-3", listed[2].Name) // Priority 1
	})

	t.Run("lists only enabled rules", func(t *testing.T) {
		listed, err := db.ListAutomationRules(ctx, true)
		require.NoError(t, err)
		assert.Len(t, listed, 2)

		// All should be enabled
		for _, r := range listed {
			assert.True(t, r.Enabled)
		}
	})
}

func TestUpdateAutomationRule(t *testing.T) {
	db, cleanup := testDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("updates rule successfully", func(t *testing.T) {
		rule := AutomationRule{
			Name:        "update-me",
			Description: "Original description",
			Enabled:     true,
			TriggerType: TriggerTypeEvent,
			TriggerConfig: map[string]interface{}{
				"event_type": "scene_changed",
			},
			Actions: []RuleAction{
				{Type: "set_mute", Parameters: map[string]interface{}{"input_name": "Mic", "muted": true}},
			},
		}

		id, err := db.CreateAutomationRule(ctx, rule)
		require.NoError(t, err)

		// Update the rule
		rule.ID = id
		rule.Description = "Updated description"
		rule.Enabled = false
		rule.Actions = []RuleAction{
			{Type: "set_scene", Parameters: map[string]interface{}{"scene_name": "New"}},
			{Type: "delay", Parameters: map[string]interface{}{"delay_ms": float64(1000)}},
		}

		err = db.UpdateAutomationRule(ctx, rule)
		require.NoError(t, err)

		// Verify update
		updated, err := db.GetAutomationRule(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated description", updated.Description)
		assert.False(t, updated.Enabled)
		assert.Len(t, updated.Actions, 2)
	})

	t.Run("fails for non-existent rule", func(t *testing.T) {
		rule := AutomationRule{
			ID:            99999,
			Name:          "ghost",
			TriggerType:   TriggerTypeEvent,
			TriggerConfig: map[string]interface{}{},
			Actions:       []RuleAction{},
		}

		err := db.UpdateAutomationRule(ctx, rule)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestDeleteAutomationRule(t *testing.T) {
	db, cleanup := testDB(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("deletes rule by ID", func(t *testing.T) {
		rule := AutomationRule{
			Name:          "delete-by-id",
			TriggerType:   TriggerTypeEvent,
			TriggerConfig: map[string]interface{}{"event_type": "test"},
			Actions:       []RuleAction{},
		}

		id, err := db.CreateAutomationRule(ctx, rule)
		require.NoError(t, err)

		err = db.DeleteAutomationRule(ctx, id)
		require.NoError(t, err)

		// Verify deletion
		_, err = db.GetAutomationRule(ctx, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("deletes rule by name", func(t *testing.T) {
		rule := AutomationRule{
			Name:          "delete-by-name",
			TriggerType:   TriggerTypeEvent,
			TriggerConfig: map[string]interface{}{"event_type": "test"},
			Actions:       []RuleAction{},
		}

		_, err := db.CreateAutomationRule(ctx, rule)
		require.NoError(t, err)

		err = db.DeleteAutomationRuleByName(ctx, "delete-by-name")
		require.NoError(t, err)

		// Verify deletion
		_, err = db.GetAutomationRuleByName(ctx, "delete-by-name")
		assert.Error(t, err)
	})

	t.Run("fails for non-existent rule", func(t *testing.T) {
		err := db.DeleteAutomationRule(ctx, 99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestSetAutomationRuleEnabled(t *testing.T) {
	db, cleanup := testDB(t)
	defer cleanup()

	ctx := context.Background()

	rule := AutomationRule{
		Name:          "toggle-enabled",
		Enabled:       true,
		TriggerType:   TriggerTypeEvent,
		TriggerConfig: map[string]interface{}{"event_type": "test"},
		Actions:       []RuleAction{},
	}

	id, err := db.CreateAutomationRule(ctx, rule)
	require.NoError(t, err)

	t.Run("disables rule", func(t *testing.T) {
		err := db.SetAutomationRuleEnabled(ctx, id, false)
		require.NoError(t, err)

		updated, err := db.GetAutomationRule(ctx, id)
		require.NoError(t, err)
		assert.False(t, updated.Enabled)
	})

	t.Run("enables rule", func(t *testing.T) {
		err := db.SetAutomationRuleEnabled(ctx, id, true)
		require.NoError(t, err)

		updated, err := db.GetAutomationRule(ctx, id)
		require.NoError(t, err)
		assert.True(t, updated.Enabled)
	})
}

func TestUpdateRuleRunStats(t *testing.T) {
	db, cleanup := testDB(t)
	defer cleanup()

	ctx := context.Background()

	rule := AutomationRule{
		Name:          "track-runs",
		TriggerType:   TriggerTypeEvent,
		TriggerConfig: map[string]interface{}{"event_type": "test"},
		Actions:       []RuleAction{},
	}

	id, err := db.CreateAutomationRule(ctx, rule)
	require.NoError(t, err)

	t.Run("updates run stats", func(t *testing.T) {
		runTime := time.Now()
		err := db.UpdateRuleRunStats(ctx, id, runTime)
		require.NoError(t, err)

		updated, err := db.GetAutomationRule(ctx, id)
		require.NoError(t, err)
		assert.NotNil(t, updated.LastRun)
		assert.Equal(t, int64(1), updated.RunCount)

		// Run again
		err = db.UpdateRuleRunStats(ctx, id, time.Now())
		require.NoError(t, err)

		updated, err = db.GetAutomationRule(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, int64(2), updated.RunCount)
	})
}

func TestRuleExecution(t *testing.T) {
	db, cleanup := testDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create a rule first
	rule := AutomationRule{
		Name:          "execution-test",
		TriggerType:   TriggerTypeEvent,
		TriggerConfig: map[string]interface{}{"event_type": "scene_changed"},
		Actions:       []RuleAction{},
	}

	ruleID, err := db.CreateAutomationRule(ctx, rule)
	require.NoError(t, err)

	t.Run("creates execution record", func(t *testing.T) {
		exec := RuleExecution{
			RuleID:      ruleID,
			RuleName:    "execution-test",
			TriggerType: TriggerTypeEvent,
			TriggerData: map[string]interface{}{
				"scene_name": "Gaming",
			},
			StartedAt: time.Now(),
			Status:    ExecutionStatusRunning,
		}

		execID, err := db.CreateRuleExecution(ctx, exec)
		require.NoError(t, err)
		assert.Greater(t, execID, int64(0))
	})

	t.Run("updates execution after completion", func(t *testing.T) {
		exec := RuleExecution{
			RuleID:      ruleID,
			RuleName:    "execution-test",
			TriggerType: TriggerTypeEvent,
			StartedAt:   time.Now(),
			Status:      ExecutionStatusRunning,
		}

		execID, err := db.CreateRuleExecution(ctx, exec)
		require.NoError(t, err)

		// Complete the execution
		completedAt := time.Now()
		exec.ID = execID
		exec.CompletedAt = &completedAt
		exec.Status = ExecutionStatusCompleted
		exec.DurationMs = 150
		exec.ActionResults = []ActionResult{
			{ActionType: "set_scene", Index: 0, Success: true, DurationMs: 50},
			{ActionType: "delay", Index: 1, Success: true, DurationMs: 100},
		}

		err = db.UpdateRuleExecution(ctx, exec)
		require.NoError(t, err)
	})

	t.Run("retrieves recent executions", func(t *testing.T) {
		executions, err := db.GetRecentRuleExecutions(ctx, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(executions), 2)
	})

	t.Run("retrieves executions for specific rule", func(t *testing.T) {
		executions, err := db.GetRuleExecutions(ctx, ruleID, 10)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(executions), 2)

		for _, exec := range executions {
			assert.Equal(t, ruleID, exec.RuleID)
		}
	})

	t.Run("clears old executions", func(t *testing.T) {
		// Create an old execution
		oldExec := RuleExecution{
			RuleID:      ruleID,
			RuleName:    "execution-test",
			TriggerType: TriggerTypeEvent,
			StartedAt:   time.Now().Add(-48 * time.Hour),
			Status:      ExecutionStatusCompleted,
		}

		_, err := db.CreateRuleExecution(ctx, oldExec)
		require.NoError(t, err)

		// Clear executions older than 24 hours
		deleted, err := db.ClearOldRuleExecutions(ctx, 24*time.Hour)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, deleted, int64(1))
	})
}

func TestCountAutomationRules(t *testing.T) {
	db, cleanup := testDB(t)
	defer cleanup()

	ctx := context.Background()

	// Initial count should be 0
	count, err := db.CountAutomationRules(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Create some rules
	for i := 0; i < 3; i++ {
		rule := AutomationRule{
			Name:          "count-test-" + string(rune('a'+i)),
			TriggerType:   TriggerTypeEvent,
			TriggerConfig: map[string]interface{}{"event_type": "test"},
			Actions:       []RuleAction{},
		}
		_, err := db.CreateAutomationRule(ctx, rule)
		require.NoError(t, err)
	}

	count, err = db.CountAutomationRules(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestExecutionCascadeDelete(t *testing.T) {
	db, cleanup := testDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create a rule
	rule := AutomationRule{
		Name:          "cascade-test",
		TriggerType:   TriggerTypeEvent,
		TriggerConfig: map[string]interface{}{"event_type": "test"},
		Actions:       []RuleAction{},
	}

	ruleID, err := db.CreateAutomationRule(ctx, rule)
	require.NoError(t, err)

	// Create execution records
	for i := 0; i < 3; i++ {
		exec := RuleExecution{
			RuleID:      ruleID,
			RuleName:    "cascade-test",
			TriggerType: TriggerTypeEvent,
			StartedAt:   time.Now(),
			Status:      ExecutionStatusCompleted,
		}
		_, err := db.CreateRuleExecution(ctx, exec)
		require.NoError(t, err)
	}

	// Verify executions exist
	executions, err := db.GetRuleExecutions(ctx, ruleID, 10)
	require.NoError(t, err)
	assert.Len(t, executions, 3)

	// Delete the rule
	err = db.DeleteAutomationRule(ctx, ruleID)
	require.NoError(t, err)

	// Executions should be cascaded deleted
	executions, err = db.GetRuleExecutions(ctx, ruleID, 10)
	require.NoError(t, err)
	assert.Len(t, executions, 0)
}
