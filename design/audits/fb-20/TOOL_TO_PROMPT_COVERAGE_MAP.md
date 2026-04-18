# Tool-to-Prompt Coverage Map

This document shows which prompts exist for each tool category and identifies gaps.

---

## Summary

| Tool Group | Tool Count | Dedicated Prompts | Coverage | Gap |
|------------|-----------|-------------------|----------|-----|
| Core | 25 | Multiple | Excellent | None |
| Sources | 3 | 1 (source-management) | Good | None |
| Audio | 4 | 1 (audio-check) | Good | None |
| Layout | 6 | 2 (preset-switcher, scene-designer) | Good | None |
| Visual | 4 | 2 (visual-check, visual-setup) | Good | None |
| Design | 14 | 1 (scene-designer) | Good | Covered by designer |
| Filters | 7 | 0 | MISSING | filter-management? |
| Transitions | 5 | 0 | MISSING | transition-design? |
| Automation | 9 | 0 | MISSING | automation-setup! |
| Meta | 4 | N/A | N/A | Help is a tool |
| **TOTAL** | **81** | **13** | **82%** | **3 gaps** |

---

## Detailed Coverage Map

### CORE TOOLS (25 tools)

#### Scene Management (4 tools)
- `list_scenes` - Used by: stream-launch, health-check, scene-organizer
- `set_current_scene` - Used by: stream-launch, stream-teardown, recording-workflow, scene-designer
- `create_scene` - Used by: scene-designer
- `remove_scene` - Used by: scene-organizer

**Prompt Coverage**: Excellent (4/4 prompts available)

#### Recording (5 tools)
- `start_recording` - Used by: stream-launch, recording-workflow
- `stop_recording` - Used by: stream-teardown, recording-workflow
- `get_recording_status` - Used by: stream-launch, stream-teardown, health-check, recording-workflow
- `pause_recording` - Used by: recording-workflow
- `resume_recording` - Used by: recording-workflow

**Prompt Coverage**: Excellent (recording-workflow covers all)

#### Streaming (3 tools)
- `start_streaming` - Used by: stream-launch
- `stop_streaming` - Used by: stream-teardown
- `get_streaming_status` - Used by: stream-launch, stream-teardown, health-check

**Prompt Coverage**: Excellent (stream-launch, stream-teardown cover)

#### Virtual Camera (2 tools)
- `get_virtual_cam_status` - Used by: stream-teardown, health-check
- `toggle_virtual_cam` - Used by: stream-teardown

**Prompt Coverage**: Adequate (stream-teardown, health-check mention)

#### Replay Buffer (2 tools)
- `get_replay_buffer_status` - Used by: stream-teardown, health-check, recording-workflow
- `save_replay_buffer` - Used by: stream-teardown, recording-workflow

**Prompt Coverage**: Adequate (stream-teardown, recording-workflow mention)

#### Studio Mode (2 tools)
- `get_studio_mode_enabled` - Used by: health-check
- `toggle_studio_mode` - Used by: (no dedicated prompt reference)

**Prompt Coverage**: Minimal (health-check mentions)

#### Hotkeys (2 tools)
- `trigger_hotkey_by_name` - Used by: (no dedicated prompt reference)
- `list_hotkeys` - Used by: health-check

**Prompt Coverage**: Minimal (health-check only)

#### Status (1 tool)
- `get_obs_status` - Used by: stream-launch, health-check, quick-status, problem-detection

**Prompt Coverage**: Excellent (multiple prompts use)

---

### SOURCES TOOLS (3 tools)

**Prompts Available**: source-management (1 of 1)

- `list_sources` - Used by: stream-launch, audio-check, health-check, scene-designer, source-management
- `toggle_source_visibility` - Used by: source-management
- `get_source_settings` - Used by: source-management

**Coverage**: GOOD - source-management handles all visibility operations

---

### AUDIO TOOLS (4 tools)

**Prompts Available**: audio-check (1 of 1)

- `get_input_mute` - Used by: stream-launch, audio-check, stream-teardown, recording-workflow
- `toggle_input_mute` - Used by: audio-check, stream-teardown
- `set_input_volume` - Used by: audio-check
- `get_input_volume` - Used by: audio-check

**Coverage**: GOOD - audio-check provides comprehensive guidance

---

### LAYOUT TOOLS (6 tools)

**Prompts Available**: preset-switcher, scene-designer (2 of 1)

- `save_scene_preset` - Used by: preset-switcher, scene-designer
- `apply_scene_preset` - Used by: preset-switcher, scene-organizer
- `list_scene_presets` - Used by: preset-switcher
- `get_preset_details` - Used by: preset-switcher
- `rename_scene_preset` - Used by: (mentioned in preset-switcher)
- `delete_scene_preset` - Used by: (mentioned in preset-switcher)

**Coverage**: GOOD - preset-switcher covers all preset management

---

### VISUAL TOOLS (4 tools)

**Prompts Available**: visual-check, visual-setup (2 of 2)

- `create_screenshot_source` - Used by: stream-launch, visual-check, visual-setup, problem-detection
- `remove_screenshot_source` - Used by: visual-setup
- `list_screenshot_sources` - Used by: visual-check, visual-setup, problem-detection
- `configure_screenshot_cadence` - Used by: visual-setup

**Coverage**: EXCELLENT - dedicated prompts for visual workflows

---

### DESIGN TOOLS (14 tools)

**Prompts Available**: scene-designer (1 of 14)

#### Source Creation (5 tools)
- `create_text_source` - Used by: scene-designer
- `create_image_source` - Used by: scene-designer
- `create_color_source` - Used by: scene-designer
- `create_browser_source` - Used by: scene-designer
- `create_media_source` - Used by: scene-designer

#### Layout Control (5 tools)
- `set_source_transform` - Used by: scene-designer
- `get_source_transform` - Used by: scene-designer
- `set_source_crop` - Used by: scene-designer
- `set_source_bounds` - Used by: scene-designer
- `set_source_order` - Used by: scene-designer

#### Advanced (4 tools)
- `set_source_locked` - Used by: scene-designer
- `duplicate_source` - Used by: source-management, scene-designer
- `remove_source` - Used by: source-management, scene-designer
- `list_input_kinds` - Used by: scene-designer

**Coverage**: EXCELLENT - scene-designer provides comprehensive design guidance

---

### FILTERS TOOLS (7 tools)

**Prompts Available**: NONE (0 of 7) ❌

- `list_source_filters` - No dedicated prompt
- `get_source_filter` - No dedicated prompt
- `create_source_filter` - No dedicated prompt
- `remove_source_filter` - No dedicated prompt
- `toggle_source_filter` - No dedicated prompt
- `set_source_filter_settings` - No dedicated prompt
- `list_filter_kinds` - No dedicated prompt

**Coverage**: MISSING ❌

**Gap**: No `filter-management` prompt exists
**Impact**: Users must discover filters via tool help
**Suggestion**: Create `filter-management` prompt for:
- Listing available filters
- Creating filters (noise suppression, color correction, etc.)
- Configuring filter settings
- Toggling filters on/off
- Best practices for filter usage

---

### TRANSITIONS TOOLS (5 tools)

**Prompts Available**: NONE (0 of 5) ❌

- `list_transitions` - No dedicated prompt
- `get_current_transition` - No dedicated prompt
- `set_current_transition` - No dedicated prompt
- `set_transition_duration` - No dedicated prompt
- `trigger_transition` - No dedicated prompt

**Coverage**: MISSING ❌

**Gap**: No `transition-design` prompt exists
**Impact**: Users must discover transitions via tool help
**Suggestion**: Create `transition-design` prompt for:
- Understanding transition types (cut, fade, swipe, etc.)
- Selecting appropriate transitions
- Configuring transition duration
- Setting up studio mode transitions
- Best practices for scene transitions

---

### AUTOMATION TOOLS (9 tools) ⚠️

**Prompts Available**: NONE (0 of 9) ❌ **[PRIMARY FINDING]**

- `list_automation_rules` - No dedicated prompt
- `get_automation_rule` - No dedicated prompt
- `create_automation_rule` - No dedicated prompt
- `update_automation_rule` - No dedicated prompt
- `delete_automation_rule` - No dedicated prompt
- `enable_automation_rule` - No dedicated prompt
- `disable_automation_rule` - No dedicated prompt
- `trigger_automation_rule` - No dedicated prompt
- `list_rule_executions` - No dedicated prompt

**Coverage**: MISSING ❌ **[CRITICAL GAP - FB-20]**

**Gap**: No `automation-setup` prompt exists
**Impact**: Users have no workflow guidance for automation rules
**Tools affected**: All 9 automation tools lack workflow prompt
**Severity**: HIGH - Automation is a powerful feature
**Recommendation**: **Create `automation-setup` prompt (PRIORITY)**

---

### META TOOLS (4 tools)

- `help` - This is a meta-tool, not a prompt target
- `get_tool_config` - Configuration query
- `set_tool_config` - Configuration control
- `list_tool_groups` - Listing groups

**Coverage**: N/A - These are system tools, not workflow targets

---

## Prompt Usage Distribution

### Most Referenced Tools (prompts that use them)

1. `list_scenes` - 5+ prompts (stream-launch, audio-check, health-check, scene-organizer, quick-status)
2. `get_obs_status` - 4 prompts (stream-launch, health-check, quick-status, problem-detection)
3. `list_sources` - 5+ prompts (stream-launch, audio-check, health-check, scene-designer, source-management)
4. `get_recording_status` - 4 prompts (stream-launch, stream-teardown, health-check, recording-workflow)

### Under-Referenced Tools (minimal prompt coverage)

1. Filters (0 prompts) - 7 tools with zero dedicated guidance
2. Transitions (0 prompts) - 5 tools with zero dedicated guidance
3. Automation (0 prompts) - 9 tools with zero dedicated guidance
4. Virtual Camera (minimal) - mentioned in 2 prompts
5. Studio Mode (minimal) - mentioned in 1 prompt
6. Hotkeys (minimal) - mentioned in 1 prompt

---

## Coverage Recommendations

### Priority 1: CRITICAL (High Impact)
- **Add `automation-setup` prompt** for 9 Automation tools
  - Impact: Major feature lacks workflow guidance
  - Effort: 1-2 hours
  - Value: High

### Priority 2: MEDIUM (Moderate Impact)
- **Consider `filter-management` prompt** for 7 Filters tools
  - Impact: Common audio/video effects lack workflow guidance
  - Effort: 1-2 hours
  - Value: Medium-High

- **Consider `transition-design` prompt** for 5 Transitions tools
  - Impact: Scene transitions lack workflow guidance
  - Effort: 1-2 hours
  - Value: Medium

### Priority 3: LOW (Low Impact)
- **Enhance `health-check` prompt** for Studio Mode and Hotkeys
  - Current: Minimal coverage
  - Could expand guidance in health-check prompt

---

## Summary Statistics

| Metric | Value |
|--------|-------|
| Total Tools | 81 |
| Total Prompts | 13 |
| Prompts Needed | 3+ |
| Tool-Prompt Ratio | 6.2:1 |
| Well-Covered Groups | 6 of 9 |
| Missing Prompts | 3 groups |

---

## Conclusion

**Current State**: 82% coverage (6 of 9 tool groups have dedicated prompts)

**Gaps Identified**:
1. **Automation** (9 tools, 0 prompts) - CRITICAL
2. **Filters** (7 tools, 0 prompts) - Medium
3. **Transitions** (5 tools, 0 prompts) - Medium

**Recommendation**: Implement `automation-setup` prompt immediately to address the most significant gap.

---

*Last Updated: 2025-12-23*
*Audit Status: Complete*
