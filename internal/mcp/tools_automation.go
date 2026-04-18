package mcp

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ironystock/agentic-obs/internal/automation"
	"github.com/ironystock/agentic-obs/internal/storage"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Automation tool input types

// ListAutomationRulesInput is the input for listing automation rules.
type ListAutomationRulesInput struct {
	EnabledOnly bool `json:"enabled_only,omitempty" jsonschema:"Only return enabled rules (default: false)"`
}

// GetAutomationRuleInput is the input for getting a specific automation rule.
type GetAutomationRuleInput struct {
	Name string `json:"name" jsonschema:"Name of the automation rule to retrieve"`
}

// CreateAutomationRuleInput is the input for creating an automation rule.
type CreateAutomationRuleInput struct {
	Name          string                   `json:"name" jsonschema:"Unique name for the rule"`
	Description   string                   `json:"description,omitempty" jsonschema:"Description of what the rule does"`
	TriggerType   string                   `json:"trigger_type" jsonschema:"Trigger type: 'event', 'schedule', or 'manual'"`
	TriggerConfig map[string]interface{}   `json:"trigger_config" jsonschema:"Trigger configuration (event_type+event_filter for event, schedule for schedule)"`
	Actions       []map[string]interface{} `json:"actions" jsonschema:"List of actions to execute (type, parameters, on_error)"`
	CooldownMs    int                      `json:"cooldown_ms,omitempty" jsonschema:"Minimum time between rule executions in milliseconds (default: 0)"`
	Priority      int                      `json:"priority,omitempty" jsonschema:"Higher priority rules execute first (default: 0)"`
	Enabled       *bool                    `json:"enabled,omitempty" jsonschema:"Whether the rule is enabled (default: true)"`
}

// UpdateAutomationRuleInput is the input for updating an automation rule.
type UpdateAutomationRuleInput struct {
	Name          string                   `json:"name" jsonschema:"Name of the rule to update"`
	NewName       string                   `json:"new_name,omitempty" jsonschema:"New name for the rule"`
	Description   *string                  `json:"description,omitempty" jsonschema:"New description"`
	TriggerType   string                   `json:"trigger_type,omitempty" jsonschema:"New trigger type"`
	TriggerConfig map[string]interface{}   `json:"trigger_config,omitempty" jsonschema:"New trigger configuration"`
	Actions       []map[string]interface{} `json:"actions,omitempty" jsonschema:"New list of actions"`
	CooldownMs    *int                     `json:"cooldown_ms,omitempty" jsonschema:"New cooldown in milliseconds"`
	Priority      *int                     `json:"priority,omitempty" jsonschema:"New priority value"`
}

// DeleteAutomationRuleInput is the input for deleting an automation rule.
type DeleteAutomationRuleInput struct {
	Name string `json:"name" jsonschema:"Name of the rule to delete"`
}

// EnableAutomationRuleInput is the input for enabling/disabling an automation rule.
type EnableAutomationRuleInput struct {
	Name string `json:"name" jsonschema:"Name of the rule to enable"`
}

// TriggerAutomationRuleInput is the input for manually triggering an automation rule.
type TriggerAutomationRuleInput struct {
	Name string `json:"name" jsonschema:"Name of the rule to trigger"`
}

// ListRuleExecutionsInput is the input for listing rule execution history.
type ListRuleExecutionsInput struct {
	RuleName string `json:"rule_name,omitempty" jsonschema:"Filter by rule name (optional)"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Maximum number of executions to return (default: 20, max: 100)"`
}

// handleListAutomationRules lists all automation rules.
func (s *Server) handleListAutomationRules(ctx context.Context, request *mcpsdk.CallToolRequest, input ListAutomationRulesInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Listing automation rules (enabled_only=%v)", input.EnabledOnly)

	rules, err := s.storage.ListAutomationRules(ctx, input.EnabledOnly)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list automation rules: %w", err)
	}

	// Convert to response format
	ruleList := make([]map[string]interface{}, len(rules))
	for i, rule := range rules {
		ruleList[i] = map[string]interface{}{
			"id":           rule.ID,
			"name":         rule.Name,
			"description":  rule.Description,
			"enabled":      rule.Enabled,
			"trigger_type": rule.TriggerType,
			"priority":     rule.Priority,
			"run_count":    rule.RunCount,
			"created_at":   rule.CreatedAt.Format(time.RFC3339),
			"updated_at":   rule.UpdatedAt.Format(time.RFC3339),
		}
		if rule.LastRun != nil {
			ruleList[i]["last_run"] = rule.LastRun.Format(time.RFC3339)
		}
	}

	result := map[string]interface{}{
		"rules":   ruleList,
		"count":   len(rules),
		"message": fmt.Sprintf("Found %d automation rules", len(rules)),
	}

	s.recordAction("list_automation_rules", "List automation rules", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleGetAutomationRule retrieves a specific automation rule by name.
func (s *Server) handleGetAutomationRule(ctx context.Context, request *mcpsdk.CallToolRequest, input GetAutomationRuleInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Getting automation rule: %s", input.Name)

	rule, err := s.storage.GetAutomationRuleByName(ctx, input.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get automation rule: %w", err)
	}

	// Convert actions to response format
	actions := make([]map[string]interface{}, len(rule.Actions))
	for i, action := range rule.Actions {
		actions[i] = map[string]interface{}{
			"type":       action.Type,
			"parameters": action.Parameters,
			"on_error":   action.OnError,
		}
	}

	result := map[string]interface{}{
		"id":             rule.ID,
		"name":           rule.Name,
		"description":    rule.Description,
		"enabled":        rule.Enabled,
		"trigger_type":   rule.TriggerType,
		"trigger_config": rule.TriggerConfig,
		"actions":        actions,
		"cooldown_ms":    rule.CooldownMs,
		"priority":       rule.Priority,
		"run_count":      rule.RunCount,
		"created_at":     rule.CreatedAt.Format(time.RFC3339),
		"updated_at":     rule.UpdatedAt.Format(time.RFC3339),
	}
	if rule.LastRun != nil {
		result["last_run"] = rule.LastRun.Format(time.RFC3339)
	}

	s.recordAction("get_automation_rule", "Get automation rule", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleCreateAutomationRule creates a new automation rule.
func (s *Server) handleCreateAutomationRule(ctx context.Context, request *mcpsdk.CallToolRequest, input CreateAutomationRuleInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Creating automation rule: %s", input.Name)

	// Validate trigger type
	if input.TriggerType != automation.TriggerTypeEvent &&
		input.TriggerType != automation.TriggerTypeSchedule &&
		input.TriggerType != automation.TriggerTypeManual {
		return nil, nil, fmt.Errorf("invalid trigger_type '%s'. Must be 'event', 'schedule', or 'manual'", input.TriggerType)
	}

	// Validate schedule if trigger type is schedule
	if input.TriggerType == automation.TriggerTypeSchedule {
		schedule, ok := input.TriggerConfig["schedule"].(string)
		if !ok || schedule == "" {
			return nil, nil, fmt.Errorf("schedule trigger requires 'schedule' in trigger_config")
		}
		if err := automation.ValidateCronExpression(schedule); err != nil {
			return nil, nil, fmt.Errorf("invalid cron schedule: %w", err)
		}
	}

	// Validate event if trigger type is event
	if input.TriggerType == automation.TriggerTypeEvent {
		eventType, ok := input.TriggerConfig["event_type"].(string)
		if !ok || eventType == "" {
			return nil, nil, fmt.Errorf("event trigger requires 'event_type' in trigger_config")
		}
		// Validate event type is known
		validEvent := false
		for _, et := range automation.SupportedEventTypes() {
			if et == eventType {
				validEvent = true
				break
			}
		}
		if !validEvent {
			return nil, nil, fmt.Errorf("unknown event_type '%s'. Valid types: %v", eventType, automation.SupportedEventTypes())
		}
	}

	// Validate actions
	if len(input.Actions) == 0 {
		return nil, nil, fmt.Errorf("at least one action is required")
	}

	// Convert actions to storage format
	actions := make([]storage.RuleAction, len(input.Actions))
	for i, actionMap := range input.Actions {
		actionType, ok := actionMap["type"].(string)
		if !ok || actionType == "" {
			return nil, nil, fmt.Errorf("action %d missing 'type'", i)
		}

		// Validate action type is known
		validAction := false
		for _, at := range automation.SupportedActionTypes() {
			if at == actionType {
				validAction = true
				break
			}
		}
		if !validAction {
			return nil, nil, fmt.Errorf("unknown action type '%s'. Valid types: %v", actionType, automation.SupportedActionTypes())
		}

		params, _ := actionMap["parameters"].(map[string]interface{})
		onError, _ := actionMap["on_error"].(string)
		if onError == "" {
			onError = automation.ActionErrorContinue
		}

		actions[i] = storage.RuleAction{
			Type:       actionType,
			Parameters: params,
			OnError:    onError,
		}
	}

	// Create the rule
	enabled := true
	if input.Enabled != nil {
		enabled = *input.Enabled
	}

	rule := storage.AutomationRule{
		Name:          input.Name,
		Description:   input.Description,
		Enabled:       enabled,
		TriggerType:   input.TriggerType,
		TriggerConfig: input.TriggerConfig,
		Actions:       actions,
		CooldownMs:    input.CooldownMs,
		Priority:      input.Priority,
	}

	id, err := s.storage.CreateAutomationRule(ctx, rule)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create automation rule: %w", err)
	}

	// Notify automation engine if running
	if s.automationEngine != nil && s.automationEngine.IsRunning() {
		s.automationEngine.NotifyRuleChange(id, false)
	}

	result := map[string]interface{}{
		"id":      id,
		"name":    input.Name,
		"enabled": enabled,
		"message": fmt.Sprintf("Automation rule '%s' created successfully", input.Name),
	}

	s.recordAction("create_automation_rule", "Create automation rule", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleUpdateAutomationRule updates an existing automation rule.
func (s *Server) handleUpdateAutomationRule(ctx context.Context, request *mcpsdk.CallToolRequest, input UpdateAutomationRuleInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Updating automation rule: %s", input.Name)

	// Get existing rule
	existing, err := s.storage.GetAutomationRuleByName(ctx, input.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get automation rule: %w", err)
	}

	// Apply updates
	updated := *existing

	if input.NewName != "" {
		updated.Name = input.NewName
	}
	if input.Description != nil {
		updated.Description = *input.Description
	}
	if input.TriggerType != "" {
		updated.TriggerType = input.TriggerType
	}
	if input.TriggerConfig != nil {
		updated.TriggerConfig = input.TriggerConfig
	}
	if input.Actions != nil {
		actions := make([]storage.RuleAction, len(input.Actions))
		for i, actionMap := range input.Actions {
			actionType, ok := actionMap["type"].(string)
			if !ok || actionType == "" {
				return nil, nil, fmt.Errorf("action %d missing 'type'", i)
			}
			params, _ := actionMap["parameters"].(map[string]interface{})
			onError, _ := actionMap["on_error"].(string)
			if onError == "" {
				onError = automation.ActionErrorContinue
			}
			actions[i] = storage.RuleAction{
				Type:       actionType,
				Parameters: params,
				OnError:    onError,
			}
		}
		updated.Actions = actions
	}
	if input.CooldownMs != nil {
		updated.CooldownMs = *input.CooldownMs
	}
	if input.Priority != nil {
		updated.Priority = *input.Priority
	}

	// Validate if trigger type changed
	if input.TriggerType == automation.TriggerTypeSchedule {
		schedule, ok := updated.TriggerConfig["schedule"].(string)
		if !ok || schedule == "" {
			return nil, nil, fmt.Errorf("schedule trigger requires 'schedule' in trigger_config")
		}
		if err := automation.ValidateCronExpression(schedule); err != nil {
			return nil, nil, fmt.Errorf("invalid cron schedule: %w", err)
		}
	}

	if err := s.storage.UpdateAutomationRule(ctx, updated); err != nil {
		return nil, nil, fmt.Errorf("failed to update automation rule: %w", err)
	}

	// Notify automation engine if running
	if s.automationEngine != nil && s.automationEngine.IsRunning() {
		s.automationEngine.NotifyRuleChange(existing.ID, false)
	}

	result := map[string]interface{}{
		"id":      existing.ID,
		"name":    updated.Name,
		"message": fmt.Sprintf("Automation rule '%s' updated successfully", updated.Name),
	}

	s.recordAction("update_automation_rule", "Update automation rule", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleDeleteAutomationRule deletes an automation rule.
func (s *Server) handleDeleteAutomationRule(ctx context.Context, request *mcpsdk.CallToolRequest, input DeleteAutomationRuleInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Deleting automation rule: %s", input.Name)

	// Get rule to confirm it exists and get ID
	rule, err := s.storage.GetAutomationRuleByName(ctx, input.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get automation rule: %w", err)
	}

	// Require confirmation via elicitation. If the elicitation RPC itself
	// fails we abort — proceeding without user confirmation would silently
	// bypass the safety gate for a destructive operation.
	confirmed, err := ElicitDeleteConfirmation(ctx, getSession(request), "automation rule", input.Name)
	if err != nil {
		log.Printf("Elicitation error: %v", err)
		return nil, nil, fmt.Errorf("delete confirmation unavailable: %w", err)
	}
	if !confirmed {
		return nil, map[string]interface{}{
			"cancelled": true,
			"message":   "Deletion cancelled by user",
		}, nil
	}

	// Notify automation engine before deletion
	if s.automationEngine != nil && s.automationEngine.IsRunning() {
		s.automationEngine.NotifyRuleChange(rule.ID, true)
	}

	if err := s.storage.DeleteAutomationRule(ctx, rule.ID); err != nil {
		return nil, nil, fmt.Errorf("failed to delete automation rule: %w", err)
	}

	result := map[string]interface{}{
		"deleted": true,
		"name":    input.Name,
		"message": fmt.Sprintf("Automation rule '%s' deleted successfully", input.Name),
	}

	s.recordAction("delete_automation_rule", "Delete automation rule", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleEnableAutomationRule enables an automation rule.
func (s *Server) handleEnableAutomationRule(ctx context.Context, request *mcpsdk.CallToolRequest, input EnableAutomationRuleInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Enabling automation rule: %s", input.Name)

	rule, err := s.storage.GetAutomationRuleByName(ctx, input.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get automation rule: %w", err)
	}

	if err := s.storage.SetAutomationRuleEnabled(ctx, rule.ID, true); err != nil {
		return nil, nil, fmt.Errorf("failed to enable automation rule: %w", err)
	}

	// Notify automation engine if running
	if s.automationEngine != nil && s.automationEngine.IsRunning() {
		s.automationEngine.NotifyRuleChange(rule.ID, false)
	}

	result := map[string]interface{}{
		"id":      rule.ID,
		"name":    input.Name,
		"enabled": true,
		"message": fmt.Sprintf("Automation rule '%s' enabled", input.Name),
	}

	s.recordAction("enable_automation_rule", "Enable automation rule", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleDisableAutomationRule disables an automation rule.
func (s *Server) handleDisableAutomationRule(ctx context.Context, request *mcpsdk.CallToolRequest, input EnableAutomationRuleInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Disabling automation rule: %s", input.Name)

	rule, err := s.storage.GetAutomationRuleByName(ctx, input.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get automation rule: %w", err)
	}

	if err := s.storage.SetAutomationRuleEnabled(ctx, rule.ID, false); err != nil {
		return nil, nil, fmt.Errorf("failed to disable automation rule: %w", err)
	}

	// Notify automation engine if running
	if s.automationEngine != nil && s.automationEngine.IsRunning() {
		s.automationEngine.NotifyRuleChange(rule.ID, false)
	}

	result := map[string]interface{}{
		"id":      rule.ID,
		"name":    input.Name,
		"enabled": false,
		"message": fmt.Sprintf("Automation rule '%s' disabled", input.Name),
	}

	s.recordAction("disable_automation_rule", "Disable automation rule", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleTriggerAutomationRule manually triggers an automation rule.
func (s *Server) handleTriggerAutomationRule(ctx context.Context, request *mcpsdk.CallToolRequest, input TriggerAutomationRuleInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Triggering automation rule: %s", input.Name)

	if s.automationEngine == nil {
		return nil, nil, fmt.Errorf("automation engine is not available")
	}

	if !s.automationEngine.IsRunning() {
		return nil, nil, fmt.Errorf("automation engine is not running")
	}

	if err := s.automationEngine.TriggerRuleByName(input.Name); err != nil {
		return nil, nil, fmt.Errorf("failed to trigger automation rule: %w", err)
	}

	result := map[string]interface{}{
		"triggered": true,
		"name":      input.Name,
		"message":   fmt.Sprintf("Automation rule '%s' triggered", input.Name),
	}

	s.recordAction("trigger_automation_rule", "Trigger automation rule", input, result, true, time.Since(start))
	return nil, result, nil
}

// handleListRuleExecutions lists automation rule execution history.
func (s *Server) handleListRuleExecutions(ctx context.Context, request *mcpsdk.CallToolRequest, input ListRuleExecutionsInput) (*mcpsdk.CallToolResult, any, error) {
	start := time.Now()
	log.Printf("Listing rule executions (rule=%s, limit=%d)", input.RuleName, input.Limit)

	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var executions []storage.RuleExecution
	var err error

	if input.RuleName != "" {
		// Get rule by name first
		rule, err := s.storage.GetAutomationRuleByName(ctx, input.RuleName)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get automation rule: %w", err)
		}
		executions, err = s.storage.GetRuleExecutions(ctx, rule.ID, limit)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list rule executions: %w", err)
		}
	} else {
		executions, err = s.storage.GetRecentRuleExecutions(ctx, limit)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list rule executions: %w", err)
		}
	}

	// Convert to response format
	execList := make([]map[string]interface{}, len(executions))
	for i, exec := range executions {
		execItem := map[string]interface{}{
			"id":           exec.ID,
			"rule_id":      exec.RuleID,
			"rule_name":    exec.RuleName,
			"trigger_type": exec.TriggerType,
			"status":       exec.Status,
			"started_at":   exec.StartedAt.Format(time.RFC3339),
			"duration_ms":  exec.DurationMs,
		}
		if exec.CompletedAt != nil {
			execItem["completed_at"] = exec.CompletedAt.Format(time.RFC3339)
		}
		if exec.Error != "" {
			execItem["error"] = exec.Error
		}
		if len(exec.ActionResults) > 0 {
			execItem["action_count"] = len(exec.ActionResults)
		}
		execList[i] = execItem
	}

	result := map[string]interface{}{
		"executions": execList,
		"count":      len(executions),
		"message":    fmt.Sprintf("Found %d rule executions", len(executions)),
	}

	s.recordAction("list_rule_executions", "List rule executions", input, result, true, time.Since(start))
	return nil, result, nil
}
