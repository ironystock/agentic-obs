package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// TriggerType defines the category of automation trigger
const (
	TriggerTypeEvent    = "event"
	TriggerTypeSchedule = "schedule"
	TriggerTypeManual   = "manual"
)

// ExecutionStatus defines the state of a rule execution
const (
	ExecutionStatusRunning   = "running"
	ExecutionStatusCompleted = "completed"
	ExecutionStatusFailed    = "failed"
	ExecutionStatusSkipped   = "skipped"
)

// AutomationRule represents an automation rule with trigger and actions.
type AutomationRule struct {
	ID            int64                  `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description,omitempty"`
	Enabled       bool                   `json:"enabled"`
	TriggerType   string                 `json:"trigger_type"`
	TriggerConfig map[string]interface{} `json:"trigger_config"`
	Actions       []RuleAction           `json:"actions"`
	CooldownMs    int                    `json:"cooldown_ms,omitempty"`
	Priority      int                    `json:"priority,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	LastRun       *time.Time             `json:"last_run,omitempty"`
	RunCount      int64                  `json:"run_count"`
}

// RuleAction represents a single action in an automation rule.
type RuleAction struct {
	Type       string                 `json:"type"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	OnError    string                 `json:"on_error,omitempty"` // "continue" or "stop"
}

// RuleExecution represents a single execution of an automation rule.
type RuleExecution struct {
	ID            int64                  `json:"id"`
	RuleID        int64                  `json:"rule_id"`
	RuleName      string                 `json:"rule_name"`
	TriggerType   string                 `json:"trigger_type"`
	TriggerData   map[string]interface{} `json:"trigger_data,omitempty"`
	StartedAt     time.Time              `json:"started_at"`
	CompletedAt   *time.Time             `json:"completed_at,omitempty"`
	Status        string                 `json:"status"`
	ActionResults []ActionResult         `json:"action_results,omitempty"`
	Error         string                 `json:"error,omitempty"`
	DurationMs    int64                  `json:"duration_ms,omitempty"`
}

// ActionResult represents the result of a single action execution.
type ActionResult struct {
	ActionType string `json:"action_type"`
	Index      int    `json:"index"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
	DurationMs int64  `json:"duration_ms"`
}

// CreateAutomationRule creates a new automation rule in the database.
// Returns the ID of the newly created rule.
func (db *DB) CreateAutomationRule(ctx context.Context, rule AutomationRule) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Serialize trigger_config to JSON
	triggerJSON, err := json.Marshal(rule.TriggerConfig)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize trigger_config to JSON: %w", err)
	}

	// Serialize actions to JSON
	actionsJSON, err := json.Marshal(rule.Actions)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize actions to JSON: %w", err)
	}

	// Default enabled to true if not set
	enabled := 1
	if !rule.Enabled {
		enabled = 0
	}

	result, err := db.conn.ExecContext(ctx, `
		INSERT INTO automation_rules (name, description, enabled, trigger_type, trigger_config, actions, cooldown_ms, priority, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, rule.Name, rule.Description, enabled, rule.TriggerType, string(triggerJSON), string(actionsJSON), rule.CooldownMs, rule.Priority)

	if err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "UNIQUE constraint failed: automation_rules.name") {
			return 0, fmt.Errorf("automation rule with name '%s' already exists", rule.Name)
		}
		return 0, fmt.Errorf("failed to create automation rule: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get inserted rule ID: %w", err)
	}

	return id, nil
}

// GetAutomationRule retrieves an automation rule by ID.
func (db *DB) GetAutomationRule(ctx context.Context, id int64) (*AutomationRule, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var rule AutomationRule
	var triggerJSON, actionsJSON string
	var enabled int
	var createdAt, updatedAt string
	var lastRun sql.NullString

	err := db.conn.QueryRowContext(ctx, `
		SELECT id, name, description, enabled, trigger_type, trigger_config, actions, cooldown_ms, priority, created_at, updated_at, last_run, run_count
		FROM automation_rules
		WHERE id = ?
	`, id).Scan(&rule.ID, &rule.Name, &rule.Description, &enabled, &rule.TriggerType, &triggerJSON, &actionsJSON, &rule.CooldownMs, &rule.Priority, &createdAt, &updatedAt, &lastRun, &rule.RunCount)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("automation rule with ID %d not found", id)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get automation rule with ID %d: %w", id, err)
	}

	rule.Enabled = enabled == 1

	// Parse trigger_config JSON
	if triggerJSON != "" {
		if err := json.Unmarshal([]byte(triggerJSON), &rule.TriggerConfig); err != nil {
			return nil, fmt.Errorf("failed to parse trigger_config JSON for rule ID %d: %w", id, err)
		}
	}

	// Parse actions JSON
	if actionsJSON != "" {
		if err := json.Unmarshal([]byte(actionsJSON), &rule.Actions); err != nil {
			return nil, fmt.Errorf("failed to parse actions JSON for rule ID %d: %w", id, err)
		}
	}

	// Parse timestamps
	rule.CreatedAt, err = parseTimestamp(createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at timestamp for rule ID %d: %w", id, err)
	}
	rule.UpdatedAt, err = parseTimestamp(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at timestamp for rule ID %d: %w", id, err)
	}
	if lastRun.Valid {
		t, err := parseTimestamp(lastRun.String)
		if err == nil {
			rule.LastRun = &t
		}
	}

	return &rule, nil
}

// GetAutomationRuleByName retrieves an automation rule by name.
func (db *DB) GetAutomationRuleByName(ctx context.Context, name string) (*AutomationRule, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var rule AutomationRule
	var triggerJSON, actionsJSON string
	var enabled int
	var createdAt, updatedAt string
	var lastRun sql.NullString

	err := db.conn.QueryRowContext(ctx, `
		SELECT id, name, description, enabled, trigger_type, trigger_config, actions, cooldown_ms, priority, created_at, updated_at, last_run, run_count
		FROM automation_rules
		WHERE name = ?
	`, name).Scan(&rule.ID, &rule.Name, &rule.Description, &enabled, &rule.TriggerType, &triggerJSON, &actionsJSON, &rule.CooldownMs, &rule.Priority, &createdAt, &updatedAt, &lastRun, &rule.RunCount)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("automation rule '%s' not found", name)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get automation rule '%s': %w", name, err)
	}

	rule.Enabled = enabled == 1

	// Parse trigger_config JSON
	if triggerJSON != "" {
		if err := json.Unmarshal([]byte(triggerJSON), &rule.TriggerConfig); err != nil {
			return nil, fmt.Errorf("failed to parse trigger_config JSON for rule '%s': %w", name, err)
		}
	}

	// Parse actions JSON
	if actionsJSON != "" {
		if err := json.Unmarshal([]byte(actionsJSON), &rule.Actions); err != nil {
			return nil, fmt.Errorf("failed to parse actions JSON for rule '%s': %w", name, err)
		}
	}

	// Parse timestamps
	rule.CreatedAt, err = parseTimestamp(createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at timestamp for rule '%s': %w", name, err)
	}
	rule.UpdatedAt, err = parseTimestamp(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at timestamp for rule '%s': %w", name, err)
	}
	if lastRun.Valid {
		t, err := parseTimestamp(lastRun.String)
		if err == nil {
			rule.LastRun = &t
		}
	}

	return &rule, nil
}

// ListAutomationRules returns all automation rules, optionally filtered to enabled only.
func (db *DB) ListAutomationRules(ctx context.Context, enabledOnly bool) ([]*AutomationRule, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var rows *sql.Rows
	var err error

	if enabledOnly {
		rows, err = db.conn.QueryContext(ctx, `
			SELECT id, name, description, enabled, trigger_type, trigger_config, actions, cooldown_ms, priority, created_at, updated_at, last_run, run_count
			FROM automation_rules
			WHERE enabled = 1
			ORDER BY priority DESC, created_at ASC
		`)
	} else {
		rows, err = db.conn.QueryContext(ctx, `
			SELECT id, name, description, enabled, trigger_type, trigger_config, actions, cooldown_ms, priority, created_at, updated_at, last_run, run_count
			FROM automation_rules
			ORDER BY priority DESC, created_at ASC
		`)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list automation rules: %w", err)
	}
	defer rows.Close()

	var rules []*AutomationRule
	for rows.Next() {
		var rule AutomationRule
		var triggerJSON, actionsJSON string
		var enabled int
		var createdAt, updatedAt string
		var lastRun sql.NullString

		if err := rows.Scan(&rule.ID, &rule.Name, &rule.Description, &enabled, &rule.TriggerType, &triggerJSON, &actionsJSON, &rule.CooldownMs, &rule.Priority, &createdAt, &updatedAt, &lastRun, &rule.RunCount); err != nil {
			return nil, fmt.Errorf("failed to scan automation rule row: %w", err)
		}

		rule.Enabled = enabled == 1

		// Parse trigger_config JSON
		if triggerJSON != "" {
			if err := json.Unmarshal([]byte(triggerJSON), &rule.TriggerConfig); err != nil {
				return nil, fmt.Errorf("failed to parse trigger_config JSON for rule '%s': %w", rule.Name, err)
			}
		}

		// Parse actions JSON
		if actionsJSON != "" {
			if err := json.Unmarshal([]byte(actionsJSON), &rule.Actions); err != nil {
				return nil, fmt.Errorf("failed to parse actions JSON for rule '%s': %w", rule.Name, err)
			}
		}

		// Parse timestamps
		rule.CreatedAt, _ = parseTimestamp(createdAt)
		rule.UpdatedAt, _ = parseTimestamp(updatedAt)
		if lastRun.Valid {
			t, err := parseTimestamp(lastRun.String)
			if err == nil {
				rule.LastRun = &t
			}
		}

		rules = append(rules, &rule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating automation rule rows: %w", err)
	}

	return rules, nil
}

// UpdateAutomationRule updates an existing automation rule.
// The rule is identified by its ID.
func (db *DB) UpdateAutomationRule(ctx context.Context, rule AutomationRule) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Serialize trigger_config to JSON
	triggerJSON, err := json.Marshal(rule.TriggerConfig)
	if err != nil {
		return fmt.Errorf("failed to serialize trigger_config to JSON: %w", err)
	}

	// Serialize actions to JSON
	actionsJSON, err := json.Marshal(rule.Actions)
	if err != nil {
		return fmt.Errorf("failed to serialize actions to JSON: %w", err)
	}

	enabled := 0
	if rule.Enabled {
		enabled = 1
	}

	result, err := db.conn.ExecContext(ctx, `
		UPDATE automation_rules
		SET name = ?, description = ?, enabled = ?, trigger_type = ?, trigger_config = ?, actions = ?, cooldown_ms = ?, priority = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, rule.Name, rule.Description, enabled, rule.TriggerType, string(triggerJSON), string(actionsJSON), rule.CooldownMs, rule.Priority, rule.ID)

	if err != nil {
		// Check for unique constraint violation
		if strings.Contains(err.Error(), "UNIQUE constraint failed: automation_rules.name") {
			return fmt.Errorf("automation rule with name '%s' already exists", rule.Name)
		}
		return fmt.Errorf("failed to update automation rule '%s': %w", rule.Name, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result for rule '%s': %w", rule.Name, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("automation rule with ID %d not found", rule.ID)
	}

	return nil
}

// DeleteAutomationRule deletes an automation rule by ID.
func (db *DB) DeleteAutomationRule(ctx context.Context, id int64) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result, err := db.conn.ExecContext(ctx,
		"DELETE FROM automation_rules WHERE id = ?",
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete automation rule with ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result for rule ID %d: %w", id, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("automation rule with ID %d not found", id)
	}

	return nil
}

// DeleteAutomationRuleByName deletes an automation rule by name.
func (db *DB) DeleteAutomationRuleByName(ctx context.Context, name string) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result, err := db.conn.ExecContext(ctx,
		"DELETE FROM automation_rules WHERE name = ?",
		name,
	)

	if err != nil {
		return fmt.Errorf("failed to delete automation rule '%s': %w", name, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result for rule '%s': %w", name, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("automation rule '%s' not found", name)
	}

	return nil
}

// SetAutomationRuleEnabled enables or disables an automation rule.
func (db *DB) SetAutomationRuleEnabled(ctx context.Context, id int64, enabled bool) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	enabledInt := 0
	if enabled {
		enabledInt = 1
	}

	result, err := db.conn.ExecContext(ctx,
		"UPDATE automation_rules SET enabled = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		enabledInt, id,
	)

	if err != nil {
		return fmt.Errorf("failed to update enabled state for rule ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result for rule ID %d: %w", id, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("automation rule with ID %d not found", id)
	}

	return nil
}

// UpdateRuleRunStats updates the last_run and run_count for a rule.
func (db *DB) UpdateRuleRunStats(ctx context.Context, id int64, runAt time.Time) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result, err := db.conn.ExecContext(ctx,
		"UPDATE automation_rules SET last_run = ?, run_count = run_count + 1 WHERE id = ?",
		runAt.Format(time.RFC3339), id,
	)

	if err != nil {
		return fmt.Errorf("failed to update run stats for rule ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result for rule ID %d: %w", id, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("automation rule with ID %d not found", id)
	}

	return nil
}

// CountAutomationRules returns the total number of automation rules.
func (db *DB) CountAutomationRules(ctx context.Context) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var count int64
	err := db.conn.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM automation_rules",
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to count automation rules: %w", err)
	}

	return count, nil
}

// CreateRuleExecution creates a new rule execution record.
func (db *DB) CreateRuleExecution(ctx context.Context, exec RuleExecution) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Serialize trigger_data to JSON
	var triggerDataJSON string
	if exec.TriggerData != nil {
		data, err := json.Marshal(exec.TriggerData)
		if err != nil {
			return 0, fmt.Errorf("failed to serialize trigger_data to JSON: %w", err)
		}
		triggerDataJSON = string(data)
	}

	// Serialize action_results to JSON
	var actionResultsJSON string
	if exec.ActionResults != nil {
		data, err := json.Marshal(exec.ActionResults)
		if err != nil {
			return 0, fmt.Errorf("failed to serialize action_results to JSON: %w", err)
		}
		actionResultsJSON = string(data)
	}

	result, err := db.conn.ExecContext(ctx, `
		INSERT INTO rule_executions (rule_id, rule_name, trigger_type, trigger_data, started_at, status, action_results, error, duration_ms)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, exec.RuleID, exec.RuleName, exec.TriggerType, triggerDataJSON, exec.StartedAt.Format(time.RFC3339), exec.Status, actionResultsJSON, exec.Error, exec.DurationMs)

	if err != nil {
		return 0, fmt.Errorf("failed to create rule execution record: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get inserted execution ID: %w", err)
	}

	return id, nil
}

// UpdateRuleExecution updates a rule execution record (typically after completion).
func (db *DB) UpdateRuleExecution(ctx context.Context, exec RuleExecution) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Serialize action_results to JSON
	var actionResultsJSON string
	if exec.ActionResults != nil {
		data, err := json.Marshal(exec.ActionResults)
		if err != nil {
			return fmt.Errorf("failed to serialize action_results to JSON: %w", err)
		}
		actionResultsJSON = string(data)
	}

	var completedAtStr *string
	if exec.CompletedAt != nil {
		s := exec.CompletedAt.Format(time.RFC3339)
		completedAtStr = &s
	}

	result, err := db.conn.ExecContext(ctx, `
		UPDATE rule_executions
		SET completed_at = ?, status = ?, action_results = ?, error = ?, duration_ms = ?
		WHERE id = ?
	`, completedAtStr, exec.Status, actionResultsJSON, exec.Error, exec.DurationMs, exec.ID)

	if err != nil {
		return fmt.Errorf("failed to update rule execution record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result for execution ID %d: %w", exec.ID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("rule execution with ID %d not found", exec.ID)
	}

	return nil
}

// GetRecentRuleExecutions retrieves recent rule executions.
func (db *DB) GetRecentRuleExecutions(ctx context.Context, limit int) ([]RuleExecution, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if limit <= 0 {
		limit = 20
	}

	rows, err := db.conn.QueryContext(ctx, `
		SELECT id, rule_id, rule_name, trigger_type, trigger_data, started_at, completed_at, status, action_results, error, duration_ms
		FROM rule_executions
		ORDER BY started_at DESC
		LIMIT ?
	`, limit)

	if err != nil {
		return nil, fmt.Errorf("failed to get recent rule executions: %w", err)
	}
	defer rows.Close()

	return scanRuleExecutions(rows)
}

// GetRuleExecutions retrieves executions for a specific rule.
func (db *DB) GetRuleExecutions(ctx context.Context, ruleID int64, limit int) ([]RuleExecution, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if limit <= 0 {
		limit = 20
	}

	rows, err := db.conn.QueryContext(ctx, `
		SELECT id, rule_id, rule_name, trigger_type, trigger_data, started_at, completed_at, status, action_results, error, duration_ms
		FROM rule_executions
		WHERE rule_id = ?
		ORDER BY started_at DESC
		LIMIT ?
	`, ruleID, limit)

	if err != nil {
		return nil, fmt.Errorf("failed to get executions for rule ID %d: %w", ruleID, err)
	}
	defer rows.Close()

	return scanRuleExecutions(rows)
}

// scanRuleExecutions scans rows into RuleExecution structs.
func scanRuleExecutions(rows *sql.Rows) ([]RuleExecution, error) {
	var executions []RuleExecution

	for rows.Next() {
		var exec RuleExecution
		var triggerDataJSON, actionResultsJSON sql.NullString
		var startedAt string
		var completedAt sql.NullString

		if err := rows.Scan(&exec.ID, &exec.RuleID, &exec.RuleName, &exec.TriggerType, &triggerDataJSON, &startedAt, &completedAt, &exec.Status, &actionResultsJSON, &exec.Error, &exec.DurationMs); err != nil {
			return nil, fmt.Errorf("failed to scan rule execution row: %w", err)
		}

		// Parse trigger_data JSON
		if triggerDataJSON.Valid && triggerDataJSON.String != "" {
			if err := json.Unmarshal([]byte(triggerDataJSON.String), &exec.TriggerData); err != nil {
				// Non-fatal, continue without trigger data
				exec.TriggerData = nil
			}
		}

		// Parse action_results JSON
		if actionResultsJSON.Valid && actionResultsJSON.String != "" {
			if err := json.Unmarshal([]byte(actionResultsJSON.String), &exec.ActionResults); err != nil {
				// Non-fatal, continue without action results
				exec.ActionResults = nil
			}
		}

		// Parse timestamps
		exec.StartedAt, _ = parseTimestamp(startedAt)
		if completedAt.Valid {
			t, err := parseTimestamp(completedAt.String)
			if err == nil {
				exec.CompletedAt = &t
			}
		}

		executions = append(executions, exec)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rule execution rows: %w", err)
	}

	return executions, nil
}

// ClearOldRuleExecutions removes executions older than the specified duration.
// Returns the number of deleted records.
func (db *DB) ClearOldRuleExecutions(ctx context.Context, olderThan time.Duration) (int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	cutoff := time.Now().Add(-olderThan).Format(time.RFC3339)

	result, err := db.conn.ExecContext(ctx,
		"DELETE FROM rule_executions WHERE started_at < ?",
		cutoff,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to clear old rule executions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to check delete result: %w", err)
	}

	return rowsAffected, nil
}
