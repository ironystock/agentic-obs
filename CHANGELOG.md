# Changelog

All notable changes to agentic-obs are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added
- Documentation restructuring with `design/` directory
- Architecture Decision Records (ADRs)
- docs-maintainer agent for documentation consistency

### Changed
- **FB-15: mcpui-go extraction** - Extracted `pkg/mcpui/` to standalone module
  - New repository: [github.com/ironystock/mcpui-go](https://github.com/ironystock/mcpui-go)
  - Go SDK for MCP-UI protocol with 77.7% test coverage
  - 6 documentation files, 5 runnable examples
  - agentic-obs now depends on external module

---

## [0.12.0] - 2025-12-21

### Phase 12: Dynamic Tool Configuration

**Summary:** FB-27 (Dynamic Tool Config) and FB-28 (Skills Update).

### Added
- **Dynamic Tool Configuration** (FB-27, 3 tools):
  - `get_tool_config` - Query tool group configuration (enabled/disabled state, tool counts)
  - `set_tool_config` - Enable/disable tool groups at runtime (session-only or persistent)
  - `list_tool_groups` - List all tool groups with descriptions and status
- **Meta-tools category** - 4 always-enabled tools that cannot be disabled:
  - `help`, `get_tool_config`, `set_tool_config`, `list_tool_groups`
- Thread-safe tool configuration with `sync.RWMutex`
- Tool group metadata with tool counts and tool name lists
- Persistence option for tool configuration via SQLite

### Changed
- "Help" category renamed to "Meta" tools for clarity
- Help content updated to describe meta-tools functionality
- **streaming-assistant skill** (FB-28): Added FB-25/26 tools and workflows:
  - Virtual camera management (start/stop for video calls)
  - Replay buffer highlight capture workflows
  - Studio mode preview/program transitions
  - Hotkey automation guidance
  - Updated cleanup recommendations

### Tests
- 32 test cases across 9 test functions for tool config handlers
- Tests for getGroupEnabled, setGroupEnabled, convertToStorageConfig helpers
- Tool group metadata validation tests

### Metrics
- **Tools:** 72 (+3)
- **Resources:** 4 (unchanged)
- **Prompts:** 13 (unchanged)
- **Skills:** 4 (unchanged)

---

## [0.11.0] - 2025-12-21

### Phase 11: Virtual Camera, Replay Buffer, Studio Mode & Hotkeys

**Summary:** FB-25 and FB-26 - Advanced OBS output control and preview features.

### Added
- **Virtual Camera & Replay Buffer** (FB-25, 6 tools):
  - `get_virtual_cam_status`, `toggle_virtual_cam` - Virtual camera control
  - `get_replay_buffer_status`, `toggle_replay_buffer` - Replay buffer state
  - `save_replay_buffer`, `get_last_replay` - Capture highlights
- **Studio Mode & Hotkeys** (FB-26, 6 tools):
  - `get_studio_mode_enabled`, `toggle_studio_mode` - Studio mode control
  - `get_preview_scene`, `set_preview_scene` - Preview/program workflow
  - `list_hotkeys`, `trigger_hotkey_by_name` - Hotkey automation

### Changed
- `stream-teardown` prompt: Scene switch now happens BEFORE stopping stream (proper ordering)
- `health-check` prompt: Added virtual cam, replay buffer, studio mode, and hotkey checks
- `recording-workflow` prompt: Integrated replay buffer for highlight capture

### Tests
- 15 new test functions with 36 test cases for FB-25/FB-26 handlers
- 3 integration workflow tests (virtual cam, studio mode, hotkeys)

### Metrics
- **Tools:** 69 (+12)
- **Resources:** 4 (unchanged)
- **Prompts:** 13 (unchanged)

---

## [0.10.0] - 2025-12-20

### Phase 10: Filters & Transitions

**Summary:** FB-23 and FB-24 - Source filter management and scene transition control.

### Added
- **Filters Tool Group** (FB-23, 7 tools):
  - `list_source_filters`, `get_source_filter` - Query filters
  - `create_source_filter`, `remove_source_filter` - Manage filters
  - `toggle_source_filter`, `set_source_filter_settings` - Configure filters
  - `list_filter_kinds` - Discover available filter types
- **Transitions Tool Group** (FB-24, 5 tools):
  - `list_transitions`, `get_current_transition` - Query transitions
  - `set_current_transition`, `set_transition_duration` - Configure transitions
  - `trigger_transition` - Trigger studio mode transition

### Metrics
- **Tools:** 57 (+12)
- **Resources:** 4 (unchanged)
- **Prompts:** 13 (unchanged)

---

## [0.7.0] - 2025-12-18

### Phase 7: MCP Completions, Help Tool & Claude Skills

**Summary:** Autocomplete support, comprehensive help system, and shareable skill packages.

### Added
- **MCP Completions**: Autocomplete for prompt arguments and resource URIs
- **Help Tool**: Topic-based guidance for tools, resources, prompts, workflows, troubleshooting
- **Claude Skills**: 4 shareable skill packages
  - `streaming-assistant` - Stream management workflows
  - `scene-designer` - Visual layout creation
  - `audio-engineer` - Audio control and monitoring
  - `preset-manager` - Scene preset management
- **New Prompts**: `scene-designer`, `source-management`, `visual-setup` (total: 13)

### Metrics
- **Tools:** 45 (unchanged)
- **Resources:** 4 (unchanged)
- **Prompts:** 13 (+3)

---

## [0.6.3] - 2025-12-17

### Phase 6.3: Agentic Scene Design

**Summary:** AI can programmatically create and manipulate OBS sources.

### Added
- **Design Tool Group**: 14 new tools for scene design
  - Source Creation: `create_text_source`, `create_image_source`, `create_color_source`, `create_browser_source`, `create_media_source`
  - Layout Control: `set_source_transform`, `get_source_transform`, `set_source_crop`, `set_source_bounds`, `set_source_order`
  - Advanced: `set_source_locked`, `duplicate_source`, `remove_source`, `list_input_kinds`

### Metrics
- **Tools:** 45 (+14)
- **Resources:** 4 (unchanged)
- **Prompts:** 10 (unchanged)

---

## [0.6.2] - 2025-12-17

### Phase 6.2: TUI Dashboard

**Summary:** Terminal-based dashboard using bubbletea/lipgloss.

### Added
- TUI dashboard mode (`--tui` or `-t` flag)
- Status view with OBS connection info
- Config view with settings display
- History view with scrollable action log
- Tab navigation with keyboard shortcuts
- Auto-refresh capability

### Dependencies
- `github.com/charmbracelet/bubbletea` v1.3.3
- `github.com/charmbracelet/lipgloss` v1.1.0

---

## [0.6.1] - 2025-12-16

### Phase 6.1: Web Dashboard

**Summary:** Web-based dashboard with REST API.

### Added
- Web dashboard at `http://localhost:8765/`
- REST API endpoints:
  - `GET /api/status` - Server status
  - `GET /api/history` - Action history (supports `?limit=N`, `?tool=name`)
  - `GET /api/history/stats` - Action statistics
  - `GET /api/screenshots` - Screenshot sources
  - `GET/POST /api/config` - Configuration management
- Real-time status with auto-refresh
- Screenshot gallery with live preview
- Action history viewer with filtering
- Dark-themed responsive UI
- Action history database table

---

## [0.5.0] - 2025-12-15

### Phase 5A: Setup & Configuration

**Summary:** Tool groups, optional HTTP server, enhanced setup experience.

### Added
- **Tool Groups**: Configurable categories (Core, Visual, Layout, Audio, Sources, Design)
- First-run setup prompts for tool groups and webserver
- Optional HTTP server (can be disabled)
- Screenshot-URL resource (`obs://screenshot-url/{name}`)
- Conditional tool registration based on preferences
- Persistent configuration in SQLite

### Changed
- Tool registration now respects group preferences
- Screenshot access available via URL resource (lightweight alternative)

---

## [0.4.0] - 2025-12-15

### Phase 4: MCP Resources & Prompts

**Summary:** Expanded resources and workflow prompts.

### Added
- **Screenshot Resource**: `obs://screenshot/{name}` - Binary image blob
- **Preset Resource**: `obs://preset/{name}` - JSON configuration
- **10 MCP Prompts**:
  - `stream-launch` - Pre-stream checklist
  - `stream-teardown` - End-stream cleanup
  - `audio-check` - Audio verification
  - `visual-check` - Visual layout analysis
  - `health-check` - OBS diagnostic
  - `problem-detection` - Issue detection
  - `preset-switcher` - Preset management
  - `recording-workflow` - Recording session
  - `scene-organizer` - Scene organization
  - `quick-status` - Brief status
- Prompt argument handling (required/optional)
- 57 new tests

### Metrics
- **Resources:** 4 (+2: screenshots, presets)
- **Prompts:** 10 (new)

---

## [0.3.0] - 2025-12-15

### Phase 3: Agentic Screenshot Sources

**Summary:** Enable AI visual monitoring through periodic screenshot capture.

### Added
- **Screenshot Tools** (4):
  - `create_screenshot_source` - Create periodic capture
  - `remove_screenshot_source` - Stop and remove source
  - `list_screenshot_sources` - List sources with status
  - `configure_screenshot_cadence` - Update capture interval
- HTTP server at `http://localhost:8765/screenshot/{name}`
- Background capture manager with configurable cadence
- SQLite storage for sources and images
- Automatic cleanup (keeps 10 latest per source)
- Security hardening (path traversal prevention)

### Metrics
- **Tools:** 30 (+4)

---

## [0.2.0] - 2025-12-15

### Phase 2: Scene Presets & Testing

**Summary:** Preset management and testing infrastructure.

### Added
- **Scene Preset Tools** (6):
  - `save_scene_preset` - Save source visibility states
  - `list_scene_presets` - List saved presets
  - `get_preset_details` - Get preset configuration
  - `apply_scene_preset` - Restore preset
  - `rename_scene_preset` - Rename preset
  - `delete_scene_preset` - Remove preset
- **Audio Tool**: `get_input_volume`
- OBSClient interface for dependency injection
- Mock OBS client for testing
- Comprehensive storage layer tests
- MCP tool handler tests

### Metrics
- **Tools:** 26 (+7)

---

## [0.1.0] - 2025-12-14

### Phase 1: Foundation

**Summary:** Initial MCP server with core OBS control.

### Added
- Go 1.25.5 project structure
- MCP server with stdio transport
- **Scene Resources**: `obs://scene/{name}` with notifications
- OBS event monitoring (SceneCreated, SceneRemoved, CurrentProgramSceneChanged)
- SQLite storage layer (modernc.org/sqlite, pure Go)
- Auto-detection setup flow
- Contextual error handling

### Core Tools (19)
- **Scene Management**: `list_scenes`, `set_current_scene`, `create_scene`, `remove_scene`
- **Recording**: `start_recording`, `stop_recording`, `get_recording_status`, `pause_recording`, `resume_recording`
- **Streaming**: `start_streaming`, `stop_streaming`, `get_streaming_status`
- **Sources**: `list_sources`, `toggle_source_visibility`, `get_source_settings`
- **Audio**: `get_input_mute`, `toggle_input_mute`, `set_input_volume`
- **Status**: `get_obs_status`

### Metrics
- **Tools:** 19
- **Resources:** 1 (scenes)

---

## Version Summary

| Version | Phase | Tools | Resources | Prompts | Date |
|---------|-------|-------|-----------|---------|------|
| 0.12.0 | 12 | 72 | 4 | 13 | 2025-12-21 |
| 0.11.0 | 11 | 69 | 4 | 13 | 2025-12-21 |
| 0.10.0 | 10 | 57 | 4 | 13 | 2025-12-20 |
| 0.7.0 | 7 | 45 | 4 | 13 | 2025-12-18 |
| 0.6.3 | 6.3 | 45 | 4 | 10 | 2025-12-17 |
| 0.6.2 | 6.2 | 31 | 4 | 10 | 2025-12-17 |
| 0.6.1 | 6.1 | 31 | 4 | 10 | 2025-12-16 |
| 0.5.0 | 5A | 30 | 4 | 10 | 2025-12-15 |
| 0.4.0 | 4 | 30 | 4 | 10 | 2025-12-15 |
| 0.3.0 | 3 | 30 | 2 | - | 2025-12-15 |
| 0.2.0 | 2 | 26 | 1 | - | 2025-12-15 |
| 0.1.0 | 1 | 19 | 1 | - | 2025-12-14 |
