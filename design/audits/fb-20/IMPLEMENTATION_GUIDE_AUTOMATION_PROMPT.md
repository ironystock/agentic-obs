# Implementation Guide: automation-setup Prompt (FB-20 Enhancement)

**Purpose**: Add MCP workflow prompt for the new FB-20 Automation Rules feature
**Status**: Recommended
**Effort**: 1-2 hours
**Risk**: Low
**Priority**: High

---

## Overview

The FB-20 feature added 9 powerful automation tools but no corresponding workflow prompt. This guide provides the specification and implementation steps for adding the `automation-setup` prompt.

---

## Specification

### Prompt Definition

**File**: `internal/mcp/prompts.go`

Add after the existing 13 prompts in `registerPrompts()`:

```go
// Prompt 14: Automation Setup
s.mcpServer.AddPrompt(
    &mcpsdk.Prompt{
        Name:        "automation-setup",
        Description: "Set up and configure automation rules for hands-free OBS control with event triggers and schedules",
        Arguments: []*mcpsdk.PromptArgument{
            {
                Name:        "rule_type",
                Description: "Optional: 'event' for event-triggered rules or 'schedule' for time-based rules",
                Required:    false,
            },
            {
                Name:        "trigger_event",
                Description: "Optional: Specific event to trigger on (e.g., 'stream_start', 'scene_switch', 'recording_stop')",
                Required:    false,
            },
        },
    },
    s.handleAutomationSetup,
)
```

### Handler Implementation

**File**: `internal/mcp/prompts.go`

Add handler function at end of file (after `handleVisualSetup`):

```go
// handleAutomationSetup provides automation rule setup workflow
func (s *Server) handleAutomationSetup(ctx context.Context, req *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	log.Println("Handling automation-setup prompt")

	ruleType := ""
	triggerEvent := ""
	if req != nil && req.Params.Arguments != nil {
		if val, ok := req.Params.Arguments["rule_type"]; ok {
			ruleType = val
		}
		if val, ok := req.Params.Arguments["trigger_event"]; ok {
			triggerEvent = val
		}
	}

	promptText := `Help me set up automation rules for hands-free OBS control:

1. **Understand Automation Concepts**
   - Event-triggered rules: Execute when specific OBS events occur
   - Scheduled rules: Execute at specific times or intervals
   - Examples: Auto-record on stream start, auto-switch scenes, disable sources on stop
   - Use cases: Reduce manual operations, create consistent workflows, automate repetitive tasks

2. **List Existing Automation Rules**
   - Use list_automation_rules to see all configured rules
   - Report rule names, types (event/schedule), current status (enabled/disabled)
   - Show last execution time and success/failure statistics
   - Identify rules that are inactive or need updates`

	// Add rule type specific guidance
	if ruleType == "event" {
		promptText += `

3. **Create Event-Triggered Rule**
   - Event-triggered rules respond to OBS state changes
   - Use create_automation_rule with trigger_type='event'
   - Supported events include:
     * 'stream_start' - When streaming begins
     * 'stream_stop' - When streaming ends
     * 'recording_start' - When recording starts
     * 'recording_stop' - When recording stops
     * 'recording_pause' - When recording is paused
     * 'recording_resume' - When recording resumes
     * 'scene_switch' - When scene changes
     * 'source_visible' - When source becomes visible
     * 'source_hidden' - When source becomes hidden
   - Configure action: Which tool and parameters to execute
   - Set condition: Additional checks before execution (optional)
   - Example: Auto-record when stream starts`
	} else if ruleType == "schedule" {
		promptText += `

3. **Create Scheduled Rule**
   - Scheduled rules execute at specific times
   - Use create_automation_rule with trigger_type='schedule'
   - Specify schedule:
     * 'once': Execute at specific date/time
     * 'daily': Execute every day at specific time
     * 'weekly': Execute on specific days at specific time
     * 'interval': Execute every N minutes/hours/days
   - Example schedules:
     * Daily at 10:00 AM: Start stream preparation
     * Weekdays at 9:00 PM: Stop stream and save configuration
     * Every 30 minutes: Check recording status
   - Configure action: What to execute when triggered
   - Set enable/disable state`
	} else {
		promptText += `

3. **Create Automation Rules - Choose Rule Type**
   - **Event-triggered rules**:
     * Respond to OBS events (stream start, scene switch, recording stop, etc.)
     * Immediate reaction to state changes
     * Examples: Auto-record on stream start, switch scene on event
   - **Scheduled rules**:
     * Execute at specific times (daily, weekly, intervals)
     * Time-based automation
     * Examples: Daily backup, weekly cleanup, hourly monitoring
   - Use create_automation_rule with appropriate trigger_type
   - Or use automation-setup prompt with rule_type='event' or rule_type='schedule' for targeted guidance`
	}

	// Add trigger-event specific guidance if provided
	if triggerEvent != "" {
		promptText += fmt.Sprintf(`

4. **Set Up Rule for '%s' Event**
   - Trigger event: %s
   - Use create_automation_rule with:
     * trigger_type: 'event'
     * trigger_event: '%s'
   - Configure target action:
     * What should happen when %s occurs?
     * Which tool should execute?
     * What parameters are needed?
   - Set enabled: true to activate immediately
   - Example actions based on trigger:
     * On stream_start: start_recording, create_screenshot_source, switch_scene
     * On stream_stop: stop_recording, stop_streaming, switch_offline_scene
     * On scene_switch: disable_certain_sources, apply_preset, notify_chat`, triggerEvent, triggerEvent, triggerEvent, triggerEvent)
	} else {
		promptText += `

4. **Configure Rule Actions**
   - Define what happens when rule triggers
   - Available actions (examples):
     * Recording control: start_recording, stop_recording, pause_recording
     * Streaming control: start_streaming, stop_streaming
     * Scene management: set_current_scene, apply_scene_preset
     * Source control: toggle_source_visibility, enable_source, disable_source
     * Audio control: set_input_volume, toggle_input_mute
     * Screenshot management: create_screenshot_source, remove_screenshot_source
   - Set action parameters (scene names, source names, values)
   - Configure condition checks (optional additional validation)
   - Order multiple actions if needed (sequential execution)`
	}

	promptText += `

5. **Test the Rule**
   - Use trigger_automation_rule to manually test
   - Verify the rule executes correctly
   - Check action was performed as expected
   - Confirm no errors or warnings
   - Review execution details

6. **Monitor Rule Execution**
   - Use list_rule_executions to view execution history
   - Check execution timestamps
   - Review success/failure status
   - Identify any failed executions
   - Debug failed rules with detailed logs

7. **Manage Rules**
   - Use list_automation_rules to see all rules
   - Use get_automation_rule to review specific rule configuration
   - Use update_automation_rule to modify existing rules
   - Use enable_automation_rule to activate disabled rules
   - Use disable_automation_rule to temporarily disable rules
   - Use delete_automation_rule to remove unwanted rules (with confirmation)

8. **Automation Best Practices**
   - Use descriptive rule names (e.g., 'auto-record-on-stream', not 'rule1')
   - Start with simple rules, add complexity gradually
   - Test rules manually before relying on them
   - Check execution history regularly for failures
   - Avoid conflicting rules (multiple rules triggering same action)
   - Document rule purpose and expected behavior
   - Group related rules logically
   - Disable rules during manual testing
   - Review rules when updating OBS configuration
   - Keep rule conditions as specific as possible

9. **Common Automation Workflows**

   **Workflow: Auto-Record on Stream Start**
   - create_automation_rule(
       name: 'auto-record-stream',
       trigger_type: 'event',
       trigger_event: 'stream_start',
       action: 'start_recording'
     )
   - Result: Recording automatically starts when you begin streaming

   **Workflow: Auto-Stop on Stream End**
   - create_automation_rule(
       name: 'auto-stop-stream',
       trigger_type: 'event',
       trigger_event: 'stream_stop',
       actions: ['stop_recording', 'set_current_scene:Offline']
     )
   - Result: Stop recording and switch to offline scene when stream ends

   **Workflow: Daily Stream Check**
   - create_automation_rule(
       name: 'daily-stream-prep',
       trigger_type: 'schedule',
       schedule: 'daily',
       schedule_time: '09:00',
       action: 'health_check'
     )
   - Result: Run health check every day at 9 AM

   **Workflow: Scene-Specific Actions**
   - create_automation_rule(
       name: 'switch-on-gaming',
       trigger_type: 'event',
       trigger_event: 'scene_switch',
       condition: 'scene_name == Gaming',
       actions: ['disable_source:webcam', 'set_input_volume:mic:+3dB']
     )
   - Result: Adjust settings automatically when switching to Gaming scene

10. **Troubleshooting Automation**

    **Issue: Rule not triggering**
    - Verify rule is enabled with list_automation_rules
    - Check trigger event name is correct
    - Review rule conditions are met
    - Use trigger_automation_rule to manually test
    - Check execution history for error messages
    - Verify OBS actions are configured correctly

    **Issue: Rule triggers at wrong time**
    - Review schedule configuration for scheduled rules
    - Check trigger event conditions
    - Verify system time is correct for time-based rules
    - Review execution history for timing details

    **Issue: Rule action fails silently**
    - Use list_rule_executions to check status
    - Review error details in execution history
    - Manually execute the action to verify it works
    - Check required parameters are specified
    - Verify target scenes/sources exist

11. **Advanced Automation Patterns**

    **Multiple Actions on Single Trigger**
    - Chain actions together in single rule
    - Actions execute in order
    - If one fails, subsequent actions may not execute
    - Use separate rules for independent actions

    **Conditional Automation**
    - Add conditions to rules to check additional state
    - Example: Record only if scene is 'Gaming'
    - Conditions must be met for action to execute

    **Scheduled Batches**
    - Create rules that execute multiple actions
    - Example: Daily maintenance (save presets, check health, cleanup)
    - Interval-based rules for regular monitoring

Provide step-by-step guidance for setting up automation, with examples for both event-triggered and scheduled rules.`

	if triggerEvent != "" {
		promptText += fmt.Sprintf("\n\nTarget trigger event: %s", triggerEvent)
	}

	return &mcpsdk.GetPromptResult{
		Description: "Set up and configure automation rules for hands-free OBS control",
		Messages: []*mcpsdk.PromptMessage{{
			Role: "user",
			Content: &mcpsdk.TextContent{
				Text: promptText,
			},
		}},
	}, nil
}
```

---

## Documentation Updates

### 1. Update HelpPromptCount

**File**: `internal/mcp/help_content.go` (Line 26)

Change:
```go
HelpPromptCount   = 13 // Workflow prompts
```

To:
```go
HelpPromptCount   = 14 // Workflow prompts
```

### 2. Update GetPromptsHelp()

**File**: `internal/mcp/help_content.go` (Lines 309-375)

In the `GetPromptsHelp()` function, add after the **visual-setup** prompt section:

```go
## Automation

**automation-setup** - Create and manage automation rules
- Set up event-triggered rules (respond to OBS events)
- Configure scheduled rules (time-based automation)
- Test and monitor rule execution
- Hands-free OBS control via rules and macros
```

Also update the initial help text (line 309):

Change:
```go
help := fmt.Sprintf(`# MCP Prompts (%d workflows)
```

The count will automatically update from 13 to 14 when HelpPromptCount is changed.

### 3. Update CLAUDE.md

**File**: `CLAUDE.md` (Line 9 and 177)

Change line 9 from:
```
**Current Status:** 81 Tools | 4 Resources | 13 Prompts | 4 Skills
```

To:
```
**Current Status:** 81 Tools | 4 Resources | 14 Prompts | 4 Skills
```

Update line 177 from:
```
`stream-launch`, `stream-teardown`, `audio-check`, `visual-check`, `health-check`, `problem-detection`, `preset-switcher`, `recording-workflow`, `scene-organizer`, `quick-status`, `scene-designer`, `source-management`, `visual-setup`
```

To:
```
`stream-launch`, `stream-teardown`, `audio-check`, `visual-check`, `health-check`, `problem-detection`, `preset-switcher`, `recording-workflow`, `scene-organizer`, `quick-status`, `scene-designer`, `source-management`, `visual-setup`, `automation-setup`
```

### 4. Update README.md

**File**: `README.md` (Line 24)

Change:
```
- **13 MCP Prompts**: Pre-built workflows for common tasks and diagnostics
```

To:
```
- **14 MCP Prompts**: Pre-built workflows for common tasks and diagnostics
```

---

## Testing Checklist

### Unit Tests

```go
// Test automation-setup with no arguments
func TestHandleAutomationSetup_NoArgs(t *testing.T) {
    // Should explain automation concepts
    // Should mention list_automation_rules
    // Should offer to create new rule
}

// Test automation-setup with rule_type="event"
func TestHandleAutomationSetup_EventType(t *testing.T) {
    // Should focus on event-triggered automation
    // Should list supported OBS events
    // Should explain event conditions
}

// Test automation-setup with rule_type="schedule"
func TestHandleAutomationSetup_ScheduleType(t *testing.T) {
    // Should focus on scheduled automation
    // Should explain schedule options
    // Should provide schedule examples
}

// Test automation-setup with trigger_event
func TestHandleAutomationSetup_WithTriggerEvent(t *testing.T) {
    // Should provide targeted guidance for that event
    // Should suggest relevant actions
}
```

### Integration Tests

1. Verify automation-setup appears in `prompts/list`
2. Verify automation-setup can be called with no arguments
3. Verify automation-setup can be called with rule_type argument
4. Verify automation-setup can be called with trigger_event argument
5. Verify prompt text references valid automation tools
6. Verify prompt text is helpful and actionable

### Manual Testing

1. Run agentic-obs MCP server
2. Use MCP client to call `prompts/get` with `automation-setup`
3. Verify response is helpful and complete
4. Verify arguments are optional
5. Verify workflow guidance is clear
6. Verify tool references are accurate

---

## Implementation Checklist

- [ ] Add handler function `handleAutomationSetup()` to `prompts.go`
- [ ] Register prompt in `registerPrompts()`
- [ ] Update `HelpPromptCount` to 14 in `help_content.go`
- [ ] Update `GetPromptsHelp()` in `help_content.go`
- [ ] Update prompt count in CLAUDE.md (line 9)
- [ ] Update prompt list in CLAUDE.md (line 177)
- [ ] Update prompt count in README.md (line 24)
- [ ] Write unit tests for handler
- [ ] Run integration tests
- [ ] Manual testing with MCP client
- [ ] Verify all documentation is synchronized
- [ ] Commit changes with clear message

---

## Git Commit Message

```
Add automation-setup prompt for FB-20 Automation Rules feature

- Add handleAutomationSetup() workflow prompt handler
- Guide users through event-triggered and scheduled rule creation
- Support optional rule_type and trigger_event arguments
- Document automation best practices and workflows
- Update HelpPromptCount from 13 to 14
- Update documentation (CLAUDE.md, README.md, help_content.go)

Implements requested enhancement for FB-20 feature gap.
```

---

## Optional Enhancements (Post-Implementation)

1. **Create automation examples** in `examples/prompts/automation-setup/`
   - Example: Event-triggered auto-record
   - Example: Daily stream preparation
   - Example: Scene-specific automation

2. **Add automation documentation**
   - Best practices guide
   - Common patterns and recipes
   - Troubleshooting guide

3. **Consider additional prompts** for related gaps:
   - `filter-management` for 7 filter tools
   - `transition-design` for 5 transition tools

---

## Questions to Consider

1. Should the prompt provide template rules that users can customize?
2. Should there be separate prompts for event vs. schedule rules?
3. Should the prompt integrate with existing pre-configured rules?
4. Should there be an automation monitoring prompt (separate from setup)?

---

## Success Criteria

- automation-setup prompt is registered and callable
- Prompt appears in prompts/list
- Documentation is synchronized across all files
- Prompt count is 14 everywhere (help_content.go, CLAUDE.md, README.md)
- Automation-setup guides users through rule creation workflow
- All automation tools are properly referenced
- Tests pass (unit and integration)
- Manual MCP client testing works correctly

---

**This implementation guide is ready for development.**
