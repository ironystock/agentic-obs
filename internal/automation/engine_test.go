package automation

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/ironystock/agentic-obs/internal/obs"
	"github.com/ironystock/agentic-obs/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockOBSClient implements OBSClient for testing.
type MockOBSClient struct {
	mu            sync.Mutex
	actions       []string
	currentScene  string
	muted         map[string]bool
	failNextCall  bool
	eventCallback obs.EventCallback
}

func NewMockOBSClient() *MockOBSClient {
	return &MockOBSClient{
		currentScene: "Default",
		muted:        make(map[string]bool),
	}
}

func (m *MockOBSClient) SetCurrentScene(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.failNextCall {
		m.failNextCall = false
		return assert.AnError
	}
	m.actions = append(m.actions, "set_scene:"+name)
	m.currentScene = name
	return nil
}

func (m *MockOBSClient) StartRecording() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "start_recording")
	return nil
}

func (m *MockOBSClient) StopRecording() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "stop_recording")
	return "/output/recording.mp4", nil
}

func (m *MockOBSClient) PauseRecording() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "pause_recording")
	return nil
}

func (m *MockOBSClient) ResumeRecording() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "resume_recording")
	return nil
}

func (m *MockOBSClient) StartStreaming() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "start_streaming")
	return nil
}

func (m *MockOBSClient) StopStreaming() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "stop_streaming")
	return nil
}

func (m *MockOBSClient) GetInputMute(inputName string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.muted[inputName], nil
}

func (m *MockOBSClient) ToggleInputMute(inputName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.muted[inputName] = !m.muted[inputName]
	m.actions = append(m.actions, "toggle_mute:"+inputName)
	return nil
}

func (m *MockOBSClient) SetInputVolume(inputName string, volumeDb *float64, volumeMul *float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "set_volume:"+inputName)
	return nil
}

func (m *MockOBSClient) ToggleSourceVisibility(sceneName string, sourceID int) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "toggle_visibility")
	return true, nil
}

func (m *MockOBSClient) ToggleVirtualCam() (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "toggle_virtual_cam")
	return true, nil
}

func (m *MockOBSClient) StartVirtualCam() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "start_virtual_cam")
	return nil
}

func (m *MockOBSClient) StopVirtualCam() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "stop_virtual_cam")
	return nil
}

func (m *MockOBSClient) ToggleReplayBuffer() (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "toggle_replay_buffer")
	return true, nil
}

func (m *MockOBSClient) SaveReplayBuffer() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "save_replay")
	return nil
}

func (m *MockOBSClient) SetCurrentPreviewScene(sceneName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "set_preview_scene:"+sceneName)
	return nil
}

func (m *MockOBSClient) TriggerStudioModeTransition() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "trigger_transition")
	return nil
}

func (m *MockOBSClient) TriggerHotkeyByName(hotkeyName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, "trigger_hotkey:"+hotkeyName)
	return nil
}

func (m *MockOBSClient) SetEventCallback(callback obs.EventCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.eventCallback = callback
}

func (m *MockOBSClient) GetActions() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]string, len(m.actions))
	copy(result, m.actions)
	return result
}

func (m *MockOBSClient) ClearActions() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = nil
}

// Test helpers

func testAutomationDB(t *testing.T) (*storage.DB, func()) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "automation-test.db")
	db, err := storage.New(context.Background(), storage.Config{Path: dbPath})
	require.NoError(t, err)
	return db, func() { db.Close() }
}

func TestEngineStartStop(t *testing.T) {
	db, cleanup := testAutomationDB(t)
	defer cleanup()

	mock := NewMockOBSClient()
	engine := NewAutomationEngine(db, mock)

	t.Run("starts successfully", func(t *testing.T) {
		err := engine.Start()
		require.NoError(t, err)
		assert.True(t, engine.IsRunning())
	})

	t.Run("handles double start", func(t *testing.T) {
		err := engine.Start()
		require.NoError(t, err)
	})

	t.Run("stops successfully", func(t *testing.T) {
		engine.Stop()
		assert.False(t, engine.IsRunning())
	})

	t.Run("handles double stop", func(t *testing.T) {
		engine.Stop()
		assert.False(t, engine.IsRunning())
	})
}

func TestEngineEventTrigger(t *testing.T) {
	db, cleanup := testAutomationDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create a rule
	rule := storage.AutomationRule{
		Name:        "scene-change-mute",
		Enabled:     true,
		TriggerType: TriggerTypeEvent,
		TriggerConfig: map[string]interface{}{
			"event_type": EventSceneChanged,
			"event_filter": map[string]interface{}{
				"scene_name": "BRB",
			},
		},
		Actions: []storage.RuleAction{
			{Type: ActionTypeSetMute, Parameters: map[string]interface{}{
				"input_name": "Microphone",
				"muted":      true,
			}},
		},
	}

	_, err := db.CreateAutomationRule(ctx, rule)
	require.NoError(t, err)

	mock := NewMockOBSClient()
	engine := NewAutomationEngine(db, mock)

	err = engine.Start()
	require.NoError(t, err)
	defer engine.Stop()

	t.Run("triggers on matching event", func(t *testing.T) {
		mock.ClearActions()

		engine.HandleEvent(EventPayload{
			EventType: EventSceneChanged,
			Data: map[string]interface{}{
				"scene_name": "BRB",
			},
		})

		// Wait for async processing
		time.Sleep(100 * time.Millisecond)

		actions := mock.GetActions()
		assert.Contains(t, actions, "toggle_mute:Microphone")
	})

	t.Run("does not trigger on non-matching event", func(t *testing.T) {
		mock.ClearActions()

		engine.HandleEvent(EventPayload{
			EventType: EventSceneChanged,
			Data: map[string]interface{}{
				"scene_name": "Gaming", // Different scene
			},
		})

		time.Sleep(100 * time.Millisecond)

		actions := mock.GetActions()
		assert.Empty(t, actions)
	})

	t.Run("does not trigger on different event type", func(t *testing.T) {
		mock.ClearActions()

		engine.HandleEvent(EventPayload{
			EventType: EventRecordingStarted,
			Data:      map[string]interface{}{},
		})

		time.Sleep(100 * time.Millisecond)

		actions := mock.GetActions()
		assert.Empty(t, actions)
	})
}

func TestEngineManualTrigger(t *testing.T) {
	db, cleanup := testAutomationDB(t)
	defer cleanup()

	ctx := context.Background()

	rule := storage.AutomationRule{
		Name:          "manual-trigger-test",
		Enabled:       true,
		TriggerType:   TriggerTypeManual,
		TriggerConfig: map[string]interface{}{},
		Actions: []storage.RuleAction{
			{Type: ActionTypeStartRecording},
		},
	}

	ruleID, err := db.CreateAutomationRule(ctx, rule)
	require.NoError(t, err)

	mock := NewMockOBSClient()
	engine := NewAutomationEngine(db, mock)

	err = engine.Start()
	require.NoError(t, err)
	defer engine.Stop()

	t.Run("triggers by ID", func(t *testing.T) {
		mock.ClearActions()

		err := engine.TriggerRule(ruleID)
		require.NoError(t, err)

		time.Sleep(100 * time.Millisecond)

		actions := mock.GetActions()
		assert.Contains(t, actions, "start_recording")
	})

	t.Run("triggers by name", func(t *testing.T) {
		mock.ClearActions()

		err := engine.TriggerRuleByName("manual-trigger-test")
		require.NoError(t, err)

		time.Sleep(100 * time.Millisecond)

		actions := mock.GetActions()
		assert.Contains(t, actions, "start_recording")
	})

	t.Run("returns error for non-existent rule", func(t *testing.T) {
		err := engine.TriggerRule(99999)
		assert.Error(t, err)
		assert.IsType(t, &RuleNotFoundError{}, err)
	})
}

func TestEngineCooldown(t *testing.T) {
	db, cleanup := testAutomationDB(t)
	defer cleanup()

	ctx := context.Background()

	rule := storage.AutomationRule{
		Name:        "cooldown-test",
		Enabled:     true,
		TriggerType: TriggerTypeEvent,
		TriggerConfig: map[string]interface{}{
			"event_type": EventSceneChanged,
		},
		Actions: []storage.RuleAction{
			{Type: ActionTypeStartRecording},
		},
		CooldownMs: 500, // 500ms cooldown
	}

	_, err := db.CreateAutomationRule(ctx, rule)
	require.NoError(t, err)

	mock := NewMockOBSClient()
	engine := NewAutomationEngine(db, mock)

	err = engine.Start()
	require.NoError(t, err)
	defer engine.Stop()

	event := EventPayload{
		EventType: EventSceneChanged,
		Data:      map[string]interface{}{},
	}

	// First trigger should execute
	engine.HandleEvent(event)
	time.Sleep(50 * time.Millisecond)
	assert.Len(t, mock.GetActions(), 1)

	// Second trigger within cooldown should be skipped
	engine.HandleEvent(event)
	time.Sleep(50 * time.Millisecond)
	assert.Len(t, mock.GetActions(), 1) // Still 1

	// Wait for cooldown to expire
	time.Sleep(500 * time.Millisecond)

	// Third trigger should execute
	engine.HandleEvent(event)
	time.Sleep(50 * time.Millisecond)
	assert.Len(t, mock.GetActions(), 2)
}

// TestEngineOnErrorStop verifies that an action with OnError="stop" halts
// the rule's action chain, while the default OnError="continue" proceeds.
func TestEngineOnErrorStop(t *testing.T) {
	db, cleanup := testAutomationDB(t)
	defer cleanup()
	ctx := context.Background()

	// Rule with three actions. Action #1 (set_scene) will fail because we
	// prime mock.failNextCall. With OnError="stop" on the failing action,
	// action #2 must NOT run.
	rule := storage.AutomationRule{
		Name:          "onerror-stop",
		Enabled:       true,
		TriggerType:   TriggerTypeManual,
		TriggerConfig: map[string]interface{}{},
		Actions: []storage.RuleAction{
			{Type: ActionTypeStartRecording},
			{
				Type:       ActionTypeSetScene,
				Parameters: map[string]interface{}{"scene_name": "Fails"},
				OnError:    ActionErrorStop,
			},
			{Type: ActionTypeStopRecording},
		},
	}
	id, err := db.CreateAutomationRule(ctx, rule)
	require.NoError(t, err)

	mock := NewMockOBSClient()
	mock.failNextCall = true // causes the set_scene action to fail

	engine := NewAutomationEngine(db, mock)
	require.NoError(t, engine.Start())
	defer engine.Stop()

	require.NoError(t, engine.TriggerRule(id))

	// Wait for execution record to be finalized so the test isn't timing-flaky.
	assert.Eventually(t, func() bool {
		execs, err := db.GetRuleExecutions(ctx, id, 10)
		return err == nil && len(execs) == 1 && execs[0].Status == storage.ExecutionStatusFailed
	}, 2*time.Second, 20*time.Millisecond)

	actions := mock.GetActions()
	assert.Contains(t, actions, "start_recording", "action #0 should have run")
	assert.NotContains(t, actions, "stop_recording", "action #2 must NOT run after OnError=stop halts chain")

	// Sanity: there is exactly one failed execution with three action results
	// (action #0 success, action #1 failure, action #2 absent).
	execs, err := db.GetRuleExecutions(ctx, id, 10)
	require.NoError(t, err)
	require.Len(t, execs, 1)
	assert.Equal(t, storage.ExecutionStatusFailed, execs[0].Status)
	assert.Len(t, execs[0].ActionResults, 2, "only the first two actions should have results")
	assert.True(t, execs[0].ActionResults[0].Success, "action #0 succeeded")
	assert.False(t, execs[0].ActionResults[1].Success, "action #1 failed")
}

func TestEngineDroppedEventsCounter(t *testing.T) {
	db, cleanup := testAutomationDB(t)
	defer cleanup()

	mock := NewMockOBSClient()
	engine := NewAutomationEngine(db, mock)

	// Don't start the engine — processEvents must NOT drain the channel
	// so we can force a buffer overflow. eventChan is buffered to 100.

	assert.Equal(t, uint64(0), engine.DroppedEventsTotal(),
		"counter should start at zero")

	// Fill the buffer exactly to capacity; none should drop yet.
	for i := 0; i < 100; i++ {
		engine.HandleEvent(EventPayload{EventType: EventSceneChanged})
	}
	assert.Equal(t, uint64(0), engine.DroppedEventsTotal(),
		"no drops expected while buffer has room")

	// Overflow by 5.
	for i := 0; i < 5; i++ {
		engine.HandleEvent(EventPayload{EventType: EventSceneChanged})
	}
	assert.Equal(t, uint64(5), engine.DroppedEventsTotal(),
		"overflow events should increment counter")
}

// TestEngineRetentionSweep verifies the background retention sweeper
// deletes execution history older than the configured retention window.
func TestEngineRetentionSweep(t *testing.T) {
	db, cleanup := testAutomationDB(t)
	defer cleanup()
	ctx := context.Background()

	// Seed: one 48h-old execution (will be pruned) and one current
	// execution (will be kept). Retention = 24h.
	rule := storage.AutomationRule{
		Name:          "retention-test",
		Enabled:       true,
		TriggerType:   TriggerTypeManual,
		TriggerConfig: map[string]interface{}{},
		Actions:       []storage.RuleAction{{Type: ActionTypeStartRecording}},
	}
	ruleID, err := db.CreateAutomationRule(ctx, rule)
	require.NoError(t, err)

	old := storage.RuleExecution{
		RuleID:      ruleID,
		RuleName:    rule.Name,
		TriggerType: TriggerTypeManual,
		StartedAt:   time.Now().Add(-48 * time.Hour),
		Status:      storage.ExecutionStatusCompleted,
	}
	_, err = db.CreateRuleExecution(ctx, old)
	require.NoError(t, err)

	recent := storage.RuleExecution{
		RuleID:      ruleID,
		RuleName:    rule.Name,
		TriggerType: TriggerTypeManual,
		StartedAt:   time.Now(),
		Status:      storage.ExecutionStatusCompleted,
	}
	_, err = db.CreateRuleExecution(ctx, recent)
	require.NoError(t, err)

	mock := NewMockOBSClient()
	engine := NewAutomationEngine(db, mock)
	engine.SetExecutionRetention(24 * time.Hour)

	// Run a one-shot sweep without starting the engine loops.
	deleted, err := engine.RunRetentionSweep()
	require.NoError(t, err)
	assert.Equal(t, int64(1), deleted, "expected the 48h-old execution to be removed")

	remaining, err := db.GetRuleExecutions(ctx, ruleID, 10)
	require.NoError(t, err)
	assert.Len(t, remaining, 1, "recent execution should be retained")
}

// TestEngineRetentionSweeperStartsOnStart verifies the sweeper goroutine
// is launched by Start and shut down by Stop without deadlocking.
func TestEngineRetentionSweeperStartsOnStart(t *testing.T) {
	db, cleanup := testAutomationDB(t)
	defer cleanup()

	mock := NewMockOBSClient()
	engine := NewAutomationEngine(db, mock)
	// Aggressive sweep tempo so the goroutine runs at least once.
	engine.SetRetentionSweepInterval(20 * time.Millisecond)
	engine.SetExecutionRetention(1 * time.Millisecond)

	require.NoError(t, engine.Start())
	time.Sleep(60 * time.Millisecond) // let sweeper tick
	engine.Stop()                     // must not hang
	assert.False(t, engine.IsRunning())
}

func TestEngineMultipleActions(t *testing.T) {
	db, cleanup := testAutomationDB(t)
	defer cleanup()

	ctx := context.Background()

	rule := storage.AutomationRule{
		Name:          "multi-action",
		Enabled:       true,
		TriggerType:   TriggerTypeManual,
		TriggerConfig: map[string]interface{}{},
		Actions: []storage.RuleAction{
			{Type: ActionTypeSetScene, Parameters: map[string]interface{}{"scene_name": "Gaming"}},
			{Type: ActionTypeStartRecording},
			{Type: ActionTypeStartStreaming},
		},
	}

	ruleID, err := db.CreateAutomationRule(ctx, rule)
	require.NoError(t, err)

	mock := NewMockOBSClient()
	engine := NewAutomationEngine(db, mock)

	err = engine.Start()
	require.NoError(t, err)
	defer engine.Stop()

	err = engine.TriggerRule(ruleID)
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	actions := mock.GetActions()
	assert.Len(t, actions, 3)
	assert.Equal(t, "set_scene:Gaming", actions[0])
	assert.Equal(t, "start_recording", actions[1])
	assert.Equal(t, "start_streaming", actions[2])
}

func TestExecutorActions(t *testing.T) {
	mock := NewMockOBSClient()
	executor := NewExecutor(mock)

	tests := []struct {
		name     string
		action   Action
		expected string
	}{
		{
			name:     "set_scene",
			action:   Action{Type: ActionTypeSetScene, Parameters: map[string]interface{}{"scene_name": "Test"}},
			expected: "set_scene:Test",
		},
		{
			name:     "toggle_mute",
			action:   Action{Type: ActionTypeToggleMute, Parameters: map[string]interface{}{"input_name": "Mic"}},
			expected: "toggle_mute:Mic",
		},
		{
			name:     "start_recording",
			action:   Action{Type: ActionTypeStartRecording},
			expected: "start_recording",
		},
		{
			name:     "stop_recording",
			action:   Action{Type: ActionTypeStopRecording},
			expected: "stop_recording",
		},
		{
			name:     "start_streaming",
			action:   Action{Type: ActionTypeStartStreaming},
			expected: "start_streaming",
		},
		{
			name:     "stop_streaming",
			action:   Action{Type: ActionTypeStopStreaming},
			expected: "stop_streaming",
		},
		{
			name:     "toggle_virtual_cam",
			action:   Action{Type: ActionTypeToggleVirtualCam},
			expected: "toggle_virtual_cam",
		},
		{
			name:     "save_replay",
			action:   Action{Type: ActionTypeSaveReplay},
			expected: "save_replay",
		},
		{
			name:     "trigger_hotkey",
			action:   Action{Type: ActionTypeTriggerHotkey, Parameters: map[string]interface{}{"hotkey_name": "OBS_KEY_F1"}},
			expected: "trigger_hotkey:OBS_KEY_F1",
		},
		{
			name:     "trigger_transition",
			action:   Action{Type: ActionTypeTriggerTransition},
			expected: "trigger_transition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ClearActions()

			result := executor.ExecuteAction(tt.action, 0)

			assert.True(t, result.Success)
			assert.Empty(t, result.Error)

			actions := mock.GetActions()
			require.Len(t, actions, 1)
			assert.Equal(t, tt.expected, actions[0])
		})
	}
}

func TestExecutorDelay(t *testing.T) {
	mock := NewMockOBSClient()
	executor := NewExecutor(mock)

	action := Action{
		Type: ActionTypeDelay,
		Parameters: map[string]interface{}{
			"delay_ms": float64(100),
		},
	}

	start := time.Now()
	result := executor.ExecuteAction(action, 0)
	elapsed := time.Since(start)

	assert.True(t, result.Success)
	assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(100))
}

func TestScheduleManager(t *testing.T) {
	executed := make(chan string, 10)
	executor := func(rule *Rule) {
		executed <- rule.Name
	}

	sm := NewScheduleManager(executor)
	sm.Start()
	defer sm.Stop()

	t.Run("schedules and executes rule", func(t *testing.T) {
		rule := &Rule{
			ID:          1,
			Name:        "every-second",
			TriggerType: TriggerTypeSchedule,
			TriggerConfig: map[string]interface{}{
				"schedule": "* * * * *", // Every minute (we'll just test scheduling works)
			},
		}

		err := sm.Schedule(rule)
		require.NoError(t, err)
		assert.Equal(t, 1, sm.GetScheduledCount())
	})

	t.Run("unschedules rule", func(t *testing.T) {
		sm.Unschedule(1)
		assert.Equal(t, 0, sm.GetScheduledCount())
	})

	t.Run("validates cron expression", func(t *testing.T) {
		err := ValidateCronExpression("* * * * *")
		assert.NoError(t, err)

		err = ValidateCronExpression("invalid")
		assert.Error(t, err)
	})
}

func TestRuleHelpers(t *testing.T) {
	rule := Rule{
		TriggerType: TriggerTypeEvent,
		TriggerConfig: map[string]interface{}{
			"event_type": "scene_changed",
			"event_filter": map[string]interface{}{
				"scene_name": "BRB",
			},
		},
	}

	t.Run("GetEventType", func(t *testing.T) {
		assert.Equal(t, "scene_changed", rule.GetEventType())
	})

	t.Run("GetEventFilter", func(t *testing.T) {
		filter := rule.GetEventFilter()
		assert.Equal(t, "BRB", filter["scene_name"])
	})

	scheduleRule := Rule{
		TriggerType: TriggerTypeSchedule,
		TriggerConfig: map[string]interface{}{
			"schedule": "0 * * * *",
		},
	}

	t.Run("GetSchedule", func(t *testing.T) {
		assert.Equal(t, "0 * * * *", scheduleRule.GetSchedule())
	})
}
