# CLAUDE.MD

This file provides context about the project for AI assistants like Claude Code.

## Project Overview

**Project Name:** agentic-obs

**Description:** A Model Context Protocol (MCP) server that provides AI assistants with programmatic control over OBS Studio via the OBS WebSocket API. This server enables AI agents to manage scenes, sources, recording, streaming, and other OBS functionality through standardized MCP tools.

## Project Structure

```
agentic-obs/
├── main.go                 # Entry point, stdio MCP server or TUI dashboard
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── CLAUDE.MD              # This file - AI assistant context
├── README.md              # User-facing documentation
├── config/
│   └── config.go          # Configuration management (OBS connection, etc.)
├── internal/
│   ├── mcp/
│   │   ├── server.go      # MCP server initialization and lifecycle
│   │   ├── tools.go       # MCP tool registration and handlers (44 tools)
│   │   ├── resources.go   # MCP resource handlers (scenes as resources)
│   │   └── interfaces.go  # OBSClient interface for testing
│   ├── obs/
│   │   ├── client.go      # OBS WebSocket client wrapper
│   │   ├── commands.go    # OBS command implementations (including scene design)
│   │   └── events.go      # OBS event handling and notification dispatch
│   ├── storage/
│   │   ├── db.go          # SQLite database setup and migrations
│   │   ├── scenes.go      # Scene preset persistence
│   │   ├── screenshots.go # Screenshot source and image persistence
│   │   └── state.go       # Configuration and state management
│   ├── http/
│   │   └── server.go      # HTTP server for screenshot image serving
│   ├── screenshot/
│   │   └── manager.go     # Background screenshot capture manager
│   └── tui/
│       └── app.go         # Terminal UI dashboard (bubbletea)
└── scripts/
    └── setup.sh           # Development setup helpers
```

## Key Technologies

- **Go 1.25.5** - Latest stable Go version
- **MCP Go SDK v1.1.0** - `github.com/modelcontextprotocol/go-sdk`
- **goobs v1.5.6** - `github.com/andreykaipov/goobs` (OBS WebSocket 5.5.6 client)
- **modernc.org/sqlite** - Pure Go SQLite driver (CGo-free)
- **bubbletea** - `github.com/charmbracelet/bubbletea` (TUI framework)
- **lipgloss** - `github.com/charmbracelet/lipgloss` (TUI styling)
- **stdio transport** - MCP communication over standard input/output

## Development Setup

### Prerequisites

- Go 1.25.5 or later
- OBS Studio 28+ (with built-in WebSocket server)
- Git

### Installation

**Option 1: Install with Go (Recommended)**
```bash
go install github.com/ironystock/agentic-obs@latest
```

**Option 2: Build from Source**
```bash
# Clone the repository
git clone https://github.com/ironystock/agentic-obs.git
cd agentic-obs

# Install dependencies
go mod download

# Build the server
go build -o agentic-obs main.go

# Run the server
./agentic-obs
```

### Running the Project

```bash
# Run MCP server mode (default)
agentic-obs

# Run TUI dashboard mode
agentic-obs --tui
agentic-obs -t

# Run directly from source
go run main.go
go run main.go --tui

# Build and run
go build && ./agentic-obs

# Install to GOPATH/bin
go install
```

## Architecture & Patterns

### Communication Flow
```
MCP Client (Claude/AI Agent)
    ↕ stdio (JSON-RPC)
MCP Server (this project)
    ├─ Tools (44 total, 6 tool groups)
    ├─ Resources (scenes, screenshots, presets)
    ├─ Prompts (10 workflows)
    ↕ WebSocket (port 4455)
OBS Studio (obs-websocket)
    ↕ SQLite (local)
Configuration & Presets

TUI Dashboard (--tui mode)
    ├─ Status view (OBS connection, server info)
    ├─ Config view (settings display)
    └─ History view (action log)
```

### Design Decisions

1. **Storage Strategy (Moderate)**: SQLite stores connection config and scene presets (source visibility states). Scene presets enable save/restore of source configurations.

2. **SQLite Driver**: Using `modernc.org/sqlite` (pure Go) for simplified cross-compilation and single-user performance requirements.

3. **OBS Connection (Persistent 1:1)**: Maintains a single persistent connection to one OBS instance with auto-reconnect. Future: multi-instance support.

4. **Error Handling (Contextual)**: Errors include diagnostic context and actionable suggestions for users.

5. **Scenes as MCP Resources**: OBS scenes are exposed as MCP resources, enabling server-initiated notifications when scenes are created, deleted, or modified. This allows AI clients to stay synchronized with OBS state changes.

6. **MCP Resource Notifications**: Server monitors OBS events (SceneCreated, SceneRemoved, CurrentProgramSceneChanged) and sends `notifications/resources/updated` and `notifications/resources/list_changed` to clients.

7. **Authentication**: OBS WebSocket password stored in SQLite (local, single-user). Auto-detected on first run with defaults (localhost:4455).

8. **Setup Experience**: Auto-detect defaults, interactive prompt on failure, persist successful config. Future: TUI and web-based setup.

### MCP Resources

The server exposes three types of MCP resources for client access:

**1. OBS Scenes** - Scene configurations and source layouts
- **Resource URI Pattern**: `obs://scene/{sceneName}`
- **Resource Type**: `obs-scene`
- **Content**: JSON representation of scene configuration (sources, settings)
- **Notifications**:
  - `notifications/resources/list_changed` - When scenes are created/deleted
  - `notifications/resources/updated` - When a specific scene is modified or becomes active

**2. Screenshot Images** - Binary screenshot data
- **Resource URI Pattern**: `obs://screenshot/{sourceName}`
- **Resource Type**: `image/png` or `image/jpeg`
- **Content**: Binary image data (PNG or JPEG format)
- **Description**: Direct access to screenshot images as MCP resources

**3. Scene Presets** - Saved source visibility configurations
- **Resource URI Pattern**: `obs://preset/{presetName}`
- **Resource Type**: `obs-preset`
- **Content**: JSON representation of preset configuration (scene, sources, visibility states)
- **Description**: Access to saved scene preset configurations

**Resource Operations:**
- `resources/list` - List all available resources of all types
- `resources/read` - Get detailed resource content (JSON or binary)
- `resources/subscribe` - Subscribe to resource change notifications (future)

### MCP Tools Implemented (44 tools in 6 groups)

**Core Tools (13)** - Scene management, recording, streaming, status:
- **Scene Management**: `list_scenes`, `set_current_scene`, `create_scene`, `remove_scene`
- **Recording**: `start_recording`, `stop_recording`, `get_recording_status`, `pause_recording`, `resume_recording`
- **Streaming**: `start_streaming`, `stop_streaming`, `get_streaming_status`
- **Status**: `get_obs_status`

**Sources Tools (3)** - Source information and visibility:
- `list_sources`, `toggle_source_visibility`, `get_source_settings`

**Audio Tools (4)** - Audio input control:
- `get_input_mute`, `toggle_input_mute`, `set_input_volume`, `get_input_volume`

**Layout Tools (6)** - Scene presets:
- `save_scene_preset`, `list_scene_presets`, `get_preset_details`, `apply_scene_preset`, `rename_scene_preset`, `delete_scene_preset`

**Visual Tools (4)** - Screenshot sources:
- `create_screenshot_source`, `remove_screenshot_source`, `list_screenshot_sources`, `configure_screenshot_cadence`

**Design Tools (14)** - Agentic scene design and source manipulation:
- **Source Creation**: `create_text_source`, `create_image_source`, `create_color_source`, `create_browser_source`, `create_media_source`
- **Layout Control**: `set_source_transform`, `get_source_transform`, `set_source_crop`, `set_source_bounds`, `set_source_order`
- **Advanced**: `set_source_locked`, `duplicate_source`, `remove_source`, `list_input_kinds`

**Note:** Scenes are also exposed as MCP resources (`resources/list`), enabling notification support. Tool groups can be enabled/disabled in configuration.

### MCP Prompts Implemented (10 prompts)

Pre-built prompts for common OBS workflows and tasks:

| Prompt | Arguments | Description |
|--------|-----------|-------------|
| `stream-launch` | none | Pre-stream checklist and setup guidance |
| `stream-teardown` | none | End-stream cleanup and shutdown workflow |
| `audio-check` | none | Audio verification and diagnostics |
| `visual-check` | screenshot_source (required) | Visual layout analysis using screenshot sources |
| `health-check` | none | Comprehensive OBS diagnostic and status check |
| `problem-detection` | screenshot_source (required) | Automated issue detection from visual monitoring |
| `preset-switcher` | preset_name (optional) | Scene preset management and switching |
| `recording-workflow` | none | Complete recording session management |
| `scene-organizer` | none | Scene organization and cleanup guidance |
| `quick-status` | none | Brief status summary for rapid checks |

**Prompts enhance the AI assistant's capability to:**
- Guide users through complex multi-step workflows
- Provide context-aware recommendations
- Automate common tasks with best practices built-in
- Offer structured checklists and verification steps
- Combine multiple tools into cohesive operations

## Important Context for AI Assistants

### Code Style

- Follow standard Go conventions (gofmt, go vet)
- Use meaningful variable names
- Keep functions focused and concise
- Error handling: always check errors, provide context
- Use internal/ for private packages not meant for external import

### Project Phases

**Phase 1 (Complete):**
- Basic MCP server with stdio transport
- Scenes exposed as MCP resources with notifications
- OBS event monitoring and notification dispatch
- Core OBS tools (recording, streaming, sources, audio)
- SQLite config storage
- Persistent OBS connection with auto-reconnect
- Auto-detection of OBS connection

**Phase 2 (Complete):**
- Scene preset management (save/restore source visibility states)
- 7 new tools: scene presets (6) + get_input_volume
- Comprehensive test coverage with mock OBS client

**Phase 3 (Complete):**
- Agentic screenshot sources for AI visual monitoring
- 4 new tools: `create_screenshot_source`, `remove_screenshot_source`, `list_screenshot_sources`, `configure_screenshot_cadence`
- HTTP server for serving screenshots at `http://localhost:8765/screenshot/{name}`
- Background capture manager with configurable cadence
- Automatic cleanup to keep storage bounded

**Phase 4 (Complete):**
- MCP Resources expansion: 3 resource types (scenes, screenshots, presets)
- MCP Prompts: 10 workflow prompts for common tasks
- Screenshots and presets exposed as MCP resources
- Resource URIs: `obs://scene/{name}`, `obs://screenshot/{name}`, `obs://preset/{name}`
- Workflow prompts for streaming, recording, diagnostics, and management

**Phase 5A (Complete):**
- Tool groups with enable/disable configuration (6 groups: Core, Visual, Layout, Audio, Sources, Design)
- Optional HTTP server (can be disabled in config)
- Screenshot URL resource type (`obs://screenshot-url/{name}`)

**Phase 6.1 (Complete):**
- Web dashboard at `http://localhost:8765/` with status, config, and history views
- REST API endpoints for status, actions, and configuration
- Action history tracking in database

**Phase 6.2 (Complete):**
- TUI dashboard mode (`--tui` flag)
- Status, Config, and History views using bubbletea/lipgloss
- Tab navigation with keyboard shortcuts
- Auto-refresh and scrollable history

**Phase 6.3 (Complete):**
- Agentic scene design with 14 new Design tools
- Source creation: text, image, color, browser, media sources
- Layout control: transform, crop, bounds, z-order
- Advanced: lock, duplicate, remove sources, list input kinds

**Future Enhancements:**
- Multi-instance OBS support (requires architecture refactor)
- Additional resource types (filters, transitions)
- Automation rules and macros
- Resource subscriptions (explicit client opt-in to notifications)

### Testing

- Unit tests for storage layer
- Integration tests for OBS client (requires running OBS)
- Mock MCP client for tool handler testing

### Common Tasks

**Adding a new MCP tool:**
1. Define tool schema in `internal/mcp/tools.go`
2. Implement handler function
3. Register tool in server initialization
4. Add OBS command in `internal/obs/commands.go` if needed

**Adding a new MCP resource:**
1. Define resource schema in `internal/mcp/resources.go`
2. Implement read/list handlers
3. Add OBS event handlers in `internal/obs/events.go`
4. Send notifications when resource state changes

**Updating dependencies:**
```bash
go get -u github.com/modelcontextprotocol/go-sdk
go get -u github.com/andreykaipov/goobs
go mod tidy
```

### Known Issues or Constraints

- Single OBS instance support only (multi-instance planned for future)
- SQLite password storage is local but unencrypted (acceptable for single-user, local deployment)
- Resource notifications require persistent connection to OBS (auto-reconnect on failure)

### Areas Needing Attention

**TODO - Future Enhancements:**
- [x] **Agentic Screenshot Sources** (Phase 3): Enable AI to "see" stream via periodic screenshot capture with HTTP serving ✓
- [x] **MCP Resources & Prompts** (Phase 4): Expand resources and add workflow prompts ✓
- [x] **Tool Groups & Optional Web Server** (Phase 5A): Configurable tool groups, optional HTTP ✓
- [x] **Web Dashboard** (Phase 6.1): Web-based status, config, and history views ✓
- [x] **TUI Dashboard** (Phase 6.2): Terminal-based dashboard with `--tui` flag ✓
- [x] **Agentic Scene Design** (Phase 6.3): AI can create and manipulate sources ✓
- [ ] Automation rules and macros: Event-triggered actions and multi-step sequences
- [ ] Multi-target OBS support (architecture decision needed)
- [ ] Additional resource types (filters, transitions, audio inputs)
- [ ] Resource subscriptions (explicit client opt-in to notifications)

### Future Research Topics

1. **Additional Resource Types**: Expand beyond scenes to expose sources, filters, transitions, and audio inputs as MCP resources. Each resource type would support notifications for state changes.

2. **Multi-Instance Architecture**: Research connection pooling patterns and how to expose multiple OBS instances through a single MCP server (namespace by instance ID?). Each instance would need isolated resource URIs (e.g., `obs://instance1/scene/Gaming`).

## Additional Resources

- [MCP Specification](https://modelcontextprotocol.io)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [goobs Documentation](https://pkg.go.dev/github.com/andreykaipov/goobs)
- [OBS WebSocket Protocol](https://github.com/obsproject/obs-websocket/blob/master/docs/generated/protocol.md)
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)

---

**Last Updated:** 2025-12-18
**Go Version:** 1.25.5
**MCP SDK Version:** 1.1.0
**goobs Version:** 1.5.6
**bubbletea Version:** 1.3.3
