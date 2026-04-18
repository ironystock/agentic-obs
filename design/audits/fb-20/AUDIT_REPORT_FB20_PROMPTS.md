# Prompt Audit Report: FB-20 Automation Rules Feature

**Date**: 2025-12-23
**Auditor**: Claude Code (MCP Prompts Specialist)
**Branch**: feature/fb-25-26-virtual-cam-studio-mode

---

## Executive Summary

The agentic-obs project currently defines **13 MCP prompts** as documented. The actual prompt inventory in `internal/mcp/prompts.go` matches this count perfectly. However, the FB-20 Automation Rules feature (9 new tools added) is **NOT represented in the current prompt collection**.

**Recommendation**: Add an `automation-setup` prompt to guide users through creating and configuring automation rules and scheduled tasks.

---

## 1. Prompt Count Verification

### Current Status: VERIFIED

| File | Location | Count | Status |
|------|----------|-------|--------|
| `internal/mcp/prompts.go` | Actual prompt definitions | 13 | PASS |
| `internal/mcp/help_content.go` | HelpPromptCount constant | 13 | PASS |
| `CLAUDE.md` | Documentation | 13 | PASS |
| `README.md` | User documentation | 13 | PASS |

**All counts match: 13 prompts as documented.**

### Prompt List (Verified)

1. `stream-launch` - Pre-stream checklist
2. `stream-teardown` - Post-stream cleanup
3. `audio-check` - Audio verification
4. `visual-check` - Screenshot-based visual analysis (requires screenshot_source arg)
5. `health-check` - Comprehensive OBS diagnostic
6. `problem-detection` - Automated issue detection (requires screenshot_source arg)
7. `preset-switcher` - Scene preset management (optional preset_name arg)
8. `recording-workflow` - Recording session management
9. `scene-organizer` - Scene structure analysis
10. `quick-status` - Brief status summary
11. `scene-designer` - Visual layout creation (requires scene_name arg, optional action arg)
12. `source-management` - Source visibility management (requires scene_name arg)
13. `visual-setup` - Screenshot monitoring setup (optional monitor_scene arg)

All 13 prompts are correctly registered in `internal/mcp/prompts.go` lines 14-165.

---

## 2. FB-20 Automation Rules Feature Assessment

### New Tools Added (FB-20)

The FB-20 feature adds **9 new automation tools** to the Core toolset:

| Tool Name | Purpose | Type |
|-----------|---------|------|
| `list_automation_rules` | List all automation rules with status | Query |
| `get_automation_rule` | Get detailed rule configuration | Query |
| `create_automation_rule` | Create event-triggered or scheduled rules | Create |
| `update_automation_rule` | Modify existing rules | Update |
| `delete_automation_rule` | Remove rules (with confirmation) | Delete |
| `enable_automation_rule` | Activate a rule | Control |
| `disable_automation_rule` | Deactivate a rule | Control |
| `trigger_automation_rule` | Manually trigger rule for testing | Control |
| `list_rule_executions` | View execution history | Query |

**Location**: `internal/mcp/tools.go` - Automation tools section
**Tool Count Update**: Confirmed in `help_content.go` line 38: `HelpAutomationToolCount = 9`

### Capabilities

Automation rules support:
- **Event-triggered actions**: Respond to OBS state changes
- **Scheduled tasks**: Run actions at specific times
- **Rule management**: Full CRUD operations
- **Execution history**: Track rule performance
- **Manual testing**: Trigger rules on demand

---

## 3. Prompt Gap Analysis

### Current Prompt Coverage by Domain

| Domain | Prompts | Coverage |
|--------|---------|----------|
| Streaming/Recording | 3 | Good (stream-launch, stream-teardown, recording-workflow) |
| Diagnostics | 3 | Good (health-check, audio-check, visual-check, problem-detection) |
| Scene/Source Management | 5 | Good (scene-designer, source-management, scene-organizer, preset-switcher, visual-setup) |
| Status/Quick Actions | 1 | Good (quick-status) |
| **Automation** | 0 | **MISSING** |

### Gap Identified

**The Automation Rules domain has ZERO dedicated prompts.**

Despite having 9 powerful automation tools (create rules, manage rules, execute rules, view history), there is no workflow prompt to guide users through:
- Setting up automation rules from scratch
- Understanding trigger types and action types
- Creating event-based or scheduled automation
- Testing and debugging rules
- Monitoring rule execution
- Best practices for automation workflows

---

## 4. Recommended New Prompt: `automation-setup`

### Specification

```go
{
    Name: "automation-setup",
    Description: "Set up and configure automation rules for hands-free OBS control with event triggers and schedules",
    Arguments: []mcpsdk.PromptArgument{
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
}
```

### Prompt Responsibilities

The `automation-setup` prompt should guide users through:

1. **Understand Automation Rules**
   - What are event-triggered rules?
   - What are scheduled rules?
   - Real-world use cases (auto-record on stream start, scene transitions, etc.)

2. **List Existing Rules**
   - Use `list_automation_rules` to show current automation
   - Report rule status, trigger types, and last execution
   - Identify gaps in automation coverage

3. **Create New Rules**
   - Use `create_automation_rule` to set up automation
   - Support event-based rules (stream start/stop, scene changes, source updates)
   - Support schedule-based rules (time-based triggers)
   - Configure rule actions (start/stop recording, switch scenes, etc.)

4. **Rule Configuration**
   - Set trigger conditions
   - Configure target actions
   - Set enable/disable state
   - Add execution conditions

5. **Test Rules**
   - Use `trigger_automation_rule` to manually test
   - Verify rule executes correctly
   - Check action outcomes

6. **Monitor Rule Execution**
   - Use `list_rule_executions` to view history
   - Track success/failure rate
   - Debug failed executions

7. **Best Practices**
   - Avoid conflicting rules
   - Order rules by priority
   - Use meaningful rule names
   - Document rule purpose
   - Regular testing and cleanup

### Workflow Example

```
User: "Help me set up automation for my stream"
↓
automation-setup prompt (no arguments)
↓
Assistant:
1. Lists existing rules with list_automation_rules
2. Asks what automation user needs (record on stream start? scene transitions?)
3. Creates rule with create_automation_rule
4. Tests with trigger_automation_rule
5. Verifies with list_rule_executions
6. Recommends best practices
```

---

## 5. Tool-to-Prompt Mapping Analysis

### Current State

| Tool Category | Tools | Prompts | Status |
|---------------|-------|---------|--------|
| Core (excluding Automation) | 25 | 7 | Good (stream-launch, stream-teardown, quick-status, health-check, recording-workflow, scene-organizer, etc.) |
| Sources | 3 | 1 | Fair (source-management handles visibility) |
| Audio | 4 | 1 | Fair (audio-check covers verification) |
| Layout | 6 | 2 | Fair (preset-switcher, scene-designer) |
| Visual | 4 | 2 | Good (visual-check, problem-detection, visual-setup) |
| Design | 14 | 1 | Good (scene-designer has detailed guidance) |
| Filters | 7 | 0 | **MISSING** |
| Transitions | 5 | 0 | **MISSING** |
| **Automation** | **9** | **0** | **MISSING** |
| Meta | 4 | 0 | N/A (help is a tool, not prompts) |

### Coverage Summary

- **Well-represented**: Streaming, Recording, Visual Monitoring, Scene Design
- **Partially represented**: Audio, Sources, Presets
- **Not represented**: Filters, Transitions, **Automation**

---

## 6. Documentation Synchronization Check

### Files Verified

| File | Prompt Count | Status | Location |
|------|--------------|--------|----------|
| `internal/mcp/prompts.go` | 13 | Accurate | Lines 14-165 (13 s.mcpServer.AddPrompt calls) |
| `internal/mcp/help_content.go` | 13 | Accurate | Line 26: `HelpPromptCount = 13` |
| `CLAUDE.md` | 13 | Accurate | Line 177: Lists all 13 prompt names |
| `README.md` | 13 | Accurate | Line 24: "13 MCP Prompts" |
| `internal/mcp/help_content.go` | Lists 13 in GetPromptsHelp | Accurate | Lines 309-375 |

**All documentation is synchronized with actual prompt count.**

### Automation Documentation

| Location | Content | Status |
|----------|---------|--------|
| `help_content.go` line 38 | `HelpAutomationToolCount = 9` | Accurate |
| `help_content.go` line 204-214 | Lists automation tools | Accurate |
| `CLAUDE.md` line 165 | Lists automation tools | Accurate |
| `README.md` line 22 | Mentions "Automation Rules" feature | Present |
| Prompts section | No automation prompt mentioned | **MISSING** |

---

## 7. Implementation Recommendations

### Priority 1: Add `automation-setup` Prompt

Add the new prompt to `internal/mcp/prompts.go`:

```go
// Prompt 14: Automation Setup (NEW - FB-20)
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
                Description: "Optional: Specific event to trigger on (e.g., 'stream_start', 'scene_switch')",
                Required:    false,
            },
        },
    },
    s.handleAutomationSetup,
)
```

### Priority 2: Update Help Constants

Update `HelpPromptCount` in `help_content.go`:

```go
HelpPromptCount = 14  // Workflow prompts (was 13, +1 for automation-setup)
```

### Priority 3: Add Handler Implementation

Implement `handleAutomationSetup()` in `prompts.go` to provide comprehensive automation guidance.

### Priority 4: Update Documentation

After implementation:
- [ ] Update `HelpPromptCount` to 14 in `help_content.go`
- [ ] Add automation-setup to prompt list in `GetPromptsHelp()`
- [ ] Update `CLAUDE.md` to list 14 prompts
- [ ] Update `README.md` to show 14 prompts
- [ ] Add automation-setup to prompt examples in workflows

### Optional: Consider Additional Prompts

While not immediately urgent, also consider:

**`filter-management`** prompt for Filters tools (7 tools, currently 0 prompts)
- Guide through filter creation and configuration
- Example: "add noise suppression to microphone"

**`transition-design`** prompt for Transitions tools (5 tools, currently 0 prompts)
- Guide through setting up scene transitions
- Example: "configure fade transitions between scenes"

---

## 8. Test Coverage Verification

### Current Prompt Tests

Based on `internal/mcp/*_test.go` patterns:
- Mock OBS client in `internal/mcp/testutil/mock_obs.go`
- Prompt handlers tested via mock

### Recommended Tests for New Prompt

When implementing `automation-setup`:

```go
// Test 1: automation-setup with no arguments
// - Should list existing rules
// - Should explain automation concepts
// - Should offer to create new rule

// Test 2: automation-setup with rule_type="event"
// - Should focus on event-triggered automation
// - Should mention relevant OBS events

// Test 3: automation-setup with trigger_event="stream_start"
// - Should provide targeted guidance for that specific event
// - Should suggest relevant action types

// Test 4: automation-setup integration
// - Should reference valid automation tools
// - Should handle missing rules gracefully
```

---

## 9. Findings Summary

### What's Working Well

1. **Prompt count is accurate and synchronized** across all files
2. **All 13 existing prompts are well-implemented** with clear workflows
3. **Help documentation is comprehensive** and up-to-date for core features
4. **Tool-to-prompt mapping is logical** for streaming, recording, and design workflows
5. **Automation tools are fully implemented** with 9 powerful tools added in FB-20

### Issues Found

1. **Zero prompts for Automation Rules** - No workflow guidance despite 9 new tools
2. **Zero prompts for Filters** - 7 filter tools lack workflow prompt
3. **Zero prompts for Transitions** - 5 transition tools lack workflow prompt
4. **Help documentation** mentions automation features but has no prompt workflow

### Risk Assessment

**LOW RISK** - Adding `automation-setup` prompt:
- Follows existing prompt patterns
- Uses already-implemented tools
- No breaking changes
- Enhancement to existing functionality

**IMPACT** - Users attempting to use automation rules:
- Must discover tools via tool help (manual process)
- No guided workflow like other domains
- Reduced discoverability of automation capabilities
- Missed opportunity for AI to guide complex rule creation

---

## 10. Detailed Recommendations

### Immediate Actions

1. **Create `automation-setup` prompt** (Priority: HIGH)
   - Add handler in `prompts.go`
   - Implement workflow covering all 9 automation tools
   - Add to prompt registration

2. **Update `HelpPromptCount` to 14** (Priority: HIGH)
   - Update `internal/mcp/help_content.go` line 26
   - Verify sync with actual prompt count

3. **Update documentation** (Priority: MEDIUM)
   - `CLAUDE.md` - Update prompt count and list
   - `README.md` - Update prompt count
   - `help_content.go` - Add automation-setup to prompt help text

### Future Enhancements

1. **Consider `filter-management` prompt** for comprehensive filter workflows
2. **Consider `transition-design` prompt** for transition setup guidance
3. **Create automation examples** in `examples/prompts/` directory
4. **Add automation best practices** documentation

---

## 11. File References

### Source Files Examined

| File | Purpose | Lines |
|------|---------|-------|
| `E:\code\agentic-obs\internal\mcp\prompts.go` | Prompt definitions | 1-1120 |
| `E:\code\agentic-obs\internal\mcp\help_content.go` | Help constants and content | 1-670 |
| `E:\code\agentic-obs\CLAUDE.md` | AI context documentation | 1-225 |
| `E:\code\agentic-obs\README.md` | User documentation | 1-100+ |
| `E:\code\agentic-obs\internal\mcp\tools.go` | Tool definitions | (automation section verified) |

### Updated By This Audit

Files requiring updates if implementing recommendation:
1. `internal/mcp/prompts.go` - Add handler and registration
2. `internal/mcp/help_content.go` - Update HelpPromptCount, add to help text
3. `CLAUDE.md` - Update prompt count
4. `README.md` - Update prompt count

---

## Conclusion

The current prompt inventory is **accurate and well-synchronized** across all documentation. However, the FB-20 Automation Rules feature represents a significant capability gap - 9 powerful automation tools exist with zero dedicated workflow prompts.

**Strong recommendation to add `automation-setup` prompt** to:
- Guide users through automation rule creation
- Document best practices
- Reduce friction for discovering automation capabilities
- Maintain consistency with other feature domains

This would bring the prompt count from 13 to 14 and ensure all major tool categories (except Meta) have workflow guidance.

---

**Report Status**: Complete
**Recommendations**: Ready for Implementation
**Impact**: Low Risk, High Value
