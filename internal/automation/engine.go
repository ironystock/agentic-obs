package automation

import (
	"context"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/ironystock/agentic-obs/internal/storage"
)

// AutomationEngine manages automation rules and their execution.
type AutomationEngine struct {
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	storage   *storage.DB
	executor  *Executor
	scheduler *ScheduleManager

	rules     map[int64]*Rule     // In-memory rule cache
	cooldowns map[int64]time.Time // Last execution time per rule

	eventChan chan EventPayload
	wg        sync.WaitGroup // Tracks in-flight processEvents + executeRule goroutines
	running   bool
}

// NewAutomationEngine creates a new automation engine.
func NewAutomationEngine(db *storage.DB, obsClient OBSClient) *AutomationEngine {
	ctx, cancel := context.WithCancel(context.Background())

	engine := &AutomationEngine{
		ctx:       ctx,
		cancel:    cancel,
		storage:   db,
		executor:  NewExecutor(obsClient),
		rules:     make(map[int64]*Rule),
		cooldowns: make(map[int64]time.Time),
		eventChan: make(chan EventPayload, 100),
	}

	return engine
}

// Start loads rules and begins processing events.
func (e *AutomationEngine) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.running {
		return nil
	}

	// Load enabled rules from storage
	if err := e.loadRulesLocked(); err != nil {
		return err
	}

	// Start scheduler
	e.scheduler = NewScheduleManager(e.executeScheduledRule)
	e.scheduler.Start()

	// Schedule all schedule-type rules
	for _, rule := range e.rules {
		if rule.Enabled && rule.TriggerType == TriggerTypeSchedule {
			if err := e.scheduler.Schedule(rule); err != nil {
				log.Printf("[Automation] Warning: failed to schedule rule '%s': %v", rule.Name, err)
			}
		}
	}

	// Start event processing goroutine
	e.wg.Add(1)
	go e.processEvents()

	e.running = true
	log.Printf("[Automation] Engine started with %d rules", len(e.rules))
	return nil
}

// Stop gracefully shuts down the engine. Waits for all in-flight
// event dispatch and rule execution goroutines before returning so
// execution records never get stranded in the "running" state.
func (e *AutomationEngine) Stop() {
	e.mu.Lock()
	if !e.running {
		e.mu.Unlock()
		return
	}
	e.running = false
	e.mu.Unlock()

	e.cancel()

	if e.scheduler != nil {
		e.scheduler.Stop()
	}

	close(e.eventChan)
	e.wg.Wait()
	log.Println("[Automation] Engine stopped")
}

// IsRunning returns whether the engine is running.
func (e *AutomationEngine) IsRunning() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.running
}

// HandleEvent dispatches an OBS event to matching rules.
func (e *AutomationEngine) HandleEvent(payload EventPayload) {
	if payload.Timestamp.IsZero() {
		payload.Timestamp = time.Now()
	}

	select {
	case e.eventChan <- payload:
	default:
		log.Printf("[Automation] Event buffer full, dropping event: %s", payload.EventType)
	}
}

// TriggerRule manually triggers a rule by ID.
func (e *AutomationEngine) TriggerRule(ruleID int64) error {
	e.mu.RLock()
	rule, exists := e.rules[ruleID]
	e.mu.RUnlock()

	if !exists {
		return &RuleNotFoundError{ID: ruleID}
	}

	e.wg.Add(1)
	go e.executeRule(rule, nil)
	return nil
}

// TriggerRuleByName manually triggers a rule by name.
func (e *AutomationEngine) TriggerRuleByName(name string) error {
	e.mu.RLock()
	var rule *Rule
	for _, r := range e.rules {
		if r.Name == name {
			rule = r
			break
		}
	}
	e.mu.RUnlock()

	if rule == nil {
		return &RuleNotFoundError{Name: name}
	}

	e.wg.Add(1)
	go e.executeRule(rule, nil)
	return nil
}

// ReloadRules reloads rules from storage.
func (e *AutomationEngine) ReloadRules() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.loadRulesLocked()
}

// GetRuleCount returns the number of loaded rules.
func (e *AutomationEngine) GetRuleCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.rules)
}

// loadRulesLocked loads enabled rules from storage (must hold lock).
func (e *AutomationEngine) loadRulesLocked() error {
	dbRules, err := e.storage.ListAutomationRules(e.ctx, true)
	if err != nil {
		return err
	}

	// Clear and rebuild
	e.rules = make(map[int64]*Rule)

	for _, dbRule := range dbRules {
		rule := convertStorageRule(dbRule)
		e.rules[rule.ID] = rule
	}

	log.Printf("[Automation] Loaded %d enabled rules", len(e.rules))
	return nil
}

// processEvents is the main event processing loop.
func (e *AutomationEngine) processEvents() {
	defer e.wg.Done()
	for payload := range e.eventChan {
		e.dispatchEvent(payload)
	}
}

// dispatchEvent finds and executes matching rules for an event.
//
// Cooldown is recorded at dispatch time (not at execute-end) so that a burst
// of events arriving faster than executeRule can complete cannot re-trigger
// the same rule. Because cooldown check + record must be atomic, the match
// loop runs under a write lock.
func (e *AutomationEngine) dispatchEvent(payload EventPayload) {
	e.mu.Lock()

	// Find matching rules, recording cooldown atomically for each match.
	var matching []*Rule
	now := time.Now()
	for _, rule := range e.rules {
		if !rule.Enabled {
			continue
		}
		if rule.TriggerType != TriggerTypeEvent {
			continue
		}
		if rule.GetEventType() != payload.EventType {
			continue
		}
		if !e.matchesFilter(rule.GetEventFilter(), payload.Data) {
			continue
		}
		if !e.checkCooldownLocked(rule) {
			log.Printf("[Automation] Rule '%s' skipped (cooldown)", rule.Name)
			continue
		}
		if rule.CooldownMs > 0 {
			e.cooldowns[rule.ID] = now
		}
		matching = append(matching, rule)
	}

	e.mu.Unlock()

	// Sort by priority (higher first)
	sort.Slice(matching, func(i, j int) bool {
		return matching[i].Priority > matching[j].Priority
	})

	// Execute matching rules
	for _, rule := range matching {
		e.wg.Add(1)
		go e.executeRule(rule, &payload)
	}
}

// matchesFilter checks if event data matches the rule's event filter.
func (e *AutomationEngine) matchesFilter(filter map[string]interface{}, data map[string]interface{}) bool {
	if filter == nil || len(filter) == 0 {
		return true
	}

	for key, expected := range filter {
		actual, exists := data[key]
		if !exists {
			return false
		}
		// Simple equality check
		if actual != expected {
			return false
		}
	}

	return true
}

// checkCooldownLocked returns true if the rule can be executed (not in cooldown).
// Caller must hold e.mu (read or write).
func (e *AutomationEngine) checkCooldownLocked(rule *Rule) bool {
	if rule.CooldownMs <= 0 {
		return true
	}

	lastRun, exists := e.cooldowns[rule.ID]
	if !exists {
		return true
	}

	cooldown := time.Duration(rule.CooldownMs) * time.Millisecond
	return time.Since(lastRun) >= cooldown
}

// executeScheduledRule is called by the scheduler.
func (e *AutomationEngine) executeScheduledRule(rule *Rule) {
	log.Printf("[Automation] Scheduled trigger for rule '%s'", rule.Name)
	e.wg.Add(1)
	e.executeRule(rule, nil)
}

// executeRule runs a single automation rule. Every call path into this
// method must precede it with e.wg.Add(1); executeRule matches with Done.
func (e *AutomationEngine) executeRule(rule *Rule, payload *EventPayload) {
	defer e.wg.Done()
	startTime := time.Now()
	log.Printf("[Automation] Executing rule '%s' (ID: %d)", rule.Name, rule.ID)

	// Create execution record
	exec := storage.RuleExecution{
		RuleID:      rule.ID,
		RuleName:    rule.Name,
		TriggerType: rule.TriggerType,
		StartedAt:   startTime,
		Status:      storage.ExecutionStatusRunning,
	}
	if payload != nil {
		exec.TriggerData = payload.Data
	}

	execID, err := e.storage.CreateRuleExecution(e.ctx, exec)
	if err != nil {
		log.Printf("[Automation] Warning: failed to create execution record: %v", err)
	}
	exec.ID = execID

	// Execute actions sequentially
	var actionResults []storage.ActionResult
	var execError error

	for i, action := range rule.Actions {
		result := e.executor.ExecuteAction(action, i)
		actionResults = append(actionResults, storage.ActionResult{
			ActionType: result.ActionType,
			Index:      result.Index,
			Success:    result.Success,
			Error:      result.Error,
			DurationMs: result.DurationMs,
		})

		if !result.Success && action.GetOnError() == ActionErrorStop {
			execError = &ActionError{
				ActionType: action.Type,
				Index:      i,
				Message:    result.Error,
			}
			break
		}
	}

	// Update execution record
	completedAt := time.Now()
	exec.CompletedAt = &completedAt
	exec.DurationMs = time.Since(startTime).Milliseconds()
	exec.ActionResults = actionResults

	if execError != nil {
		exec.Status = storage.ExecutionStatusFailed
		exec.Error = execError.Error()
		log.Printf("[Automation] Rule '%s' failed: %v", rule.Name, execError)
	} else {
		exec.Status = storage.ExecutionStatusCompleted
		log.Printf("[Automation] Rule '%s' completed in %dms", rule.Name, exec.DurationMs)
	}

	if execID > 0 {
		if err := e.storage.UpdateRuleExecution(e.ctx, exec); err != nil {
			log.Printf("[Automation] Warning: failed to update execution record: %v", err)
		}
	}

	// Update rule run stats
	if err := e.storage.UpdateRuleRunStats(e.ctx, rule.ID, startTime); err != nil {
		log.Printf("[Automation] Warning: failed to update rule stats: %v", err)
	}
}

// NotifyRuleChange should be called when rules are modified via MCP.
func (e *AutomationEngine) NotifyRuleChange(ruleID int64, deleted bool) {
	if deleted {
		e.mu.Lock()
		if rule, exists := e.rules[ruleID]; exists {
			if rule.TriggerType == TriggerTypeSchedule && e.scheduler != nil {
				e.scheduler.Unschedule(ruleID)
			}
			delete(e.rules, ruleID)
		}
		delete(e.cooldowns, ruleID)
		e.mu.Unlock()
		return
	}

	// Reload the specific rule
	dbRule, err := e.storage.GetAutomationRule(e.ctx, ruleID)
	if err != nil {
		log.Printf("[Automation] Warning: failed to reload rule %d: %v", ruleID, err)
		return
	}

	rule := convertStorageRule(dbRule)

	e.mu.Lock()
	defer e.mu.Unlock()

	// Update scheduler if needed
	if e.scheduler != nil {
		// Remove old schedule
		e.scheduler.Unschedule(ruleID)

		// Add new schedule if applicable
		if rule.Enabled && rule.TriggerType == TriggerTypeSchedule {
			if err := e.scheduler.Schedule(rule); err != nil {
				log.Printf("[Automation] Warning: failed to schedule rule '%s': %v", rule.Name, err)
			}
		}
	}

	if rule.Enabled {
		e.rules[ruleID] = rule
	} else {
		delete(e.rules, ruleID)
	}
}

// convertStorageRule converts a storage rule to an automation rule.
func convertStorageRule(dbRule *storage.AutomationRule) *Rule {
	actions := make([]Action, len(dbRule.Actions))
	for i, a := range dbRule.Actions {
		actions[i] = Action{
			Type:       a.Type,
			Parameters: a.Parameters,
			OnError:    a.OnError,
		}
	}

	return &Rule{
		ID:            dbRule.ID,
		Name:          dbRule.Name,
		Description:   dbRule.Description,
		Enabled:       dbRule.Enabled,
		TriggerType:   dbRule.TriggerType,
		TriggerConfig: dbRule.TriggerConfig,
		Actions:       actions,
		CooldownMs:    dbRule.CooldownMs,
		Priority:      dbRule.Priority,
		CreatedAt:     dbRule.CreatedAt,
		UpdatedAt:     dbRule.UpdatedAt,
		LastRun:       dbRule.LastRun,
		RunCount:      dbRule.RunCount,
	}
}

// Error types

// RuleNotFoundError is returned when a rule is not found.
type RuleNotFoundError struct {
	ID   int64
	Name string
}

func (e *RuleNotFoundError) Error() string {
	if e.Name != "" {
		return "automation rule '" + e.Name + "' not found"
	}
	return "automation rule not found"
}

// ActionError is returned when an action fails and stops execution.
type ActionError struct {
	ActionType string
	Index      int
	Message    string
}

func (e *ActionError) Error() string {
	return e.Message
}
