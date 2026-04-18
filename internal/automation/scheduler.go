package automation

import (
	"fmt"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

// ScheduleManager handles cron-like scheduling for automation rules.
type ScheduleManager struct {
	mu       sync.RWMutex
	cron     *cron.Cron
	ruleJobs map[int64]cron.EntryID // rule ID → cron entry
	executor func(rule *Rule)
	running  bool
}

// NewScheduleManager creates a new schedule manager.
func NewScheduleManager(executor func(rule *Rule)) *ScheduleManager {
	return &ScheduleManager{
		cron:     cron.New(cron.WithSeconds()),
		ruleJobs: make(map[int64]cron.EntryID),
		executor: executor,
	}
}

// Start begins the scheduler.
func (sm *ScheduleManager) Start() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.running {
		return
	}

	sm.cron.Start()
	sm.running = true
	log.Println("[Scheduler] Started")
}

// Stop halts the scheduler.
func (sm *ScheduleManager) Stop() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if !sm.running {
		return
	}

	ctx := sm.cron.Stop()
	<-ctx.Done()
	sm.running = false
	log.Println("[Scheduler] Stopped")
}

// Schedule adds a rule to the scheduler.
func (sm *ScheduleManager) Schedule(rule *Rule) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	schedule := rule.GetSchedule()
	if schedule == "" {
		return fmt.Errorf("rule '%s' has no schedule", rule.Name)
	}

	// Convert 5-field cron to 6-field (add seconds)
	// The robfig/cron/v3 library with WithSeconds() expects 6 fields
	cronExpr := "0 " + schedule

	// Remove any existing schedule for this rule
	if entryID, exists := sm.ruleJobs[rule.ID]; exists {
		sm.cron.Remove(entryID)
	}

	// Create a closure that captures the rule
	ruleCopy := *rule // Copy to avoid issues with pointer reuse
	job := func() {
		sm.executor(&ruleCopy)
	}

	entryID, err := sm.cron.AddFunc(cronExpr, job)
	if err != nil {
		return fmt.Errorf("invalid cron expression '%s': %w", schedule, err)
	}

	sm.ruleJobs[rule.ID] = entryID
	log.Printf("[Scheduler] Scheduled rule '%s' with cron '%s'", rule.Name, schedule)
	return nil
}

// Unschedule removes a rule from the scheduler.
func (sm *ScheduleManager) Unschedule(ruleID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if entryID, exists := sm.ruleJobs[ruleID]; exists {
		sm.cron.Remove(entryID)
		delete(sm.ruleJobs, ruleID)
		log.Printf("[Scheduler] Unscheduled rule ID %d", ruleID)
	}
}

// GetScheduledCount returns the number of scheduled rules.
func (sm *ScheduleManager) GetScheduledCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.ruleJobs)
}

// ValidateCronExpression checks if a cron expression is valid.
func ValidateCronExpression(expr string) error {
	// Test parsing with seconds prefix
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse("0 " + expr)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}
	return nil
}
