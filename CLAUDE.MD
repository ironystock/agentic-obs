# CLAUDE.MD

This file provides context about the project for AI assistants like Claude Code.

## Project Overview

**Project Name:** agentic-obs

**Description:** A Model Context Protocol (MCP) server that provides AI assistants with programmatic control over OBS Studio via the OBS WebSocket API. This server enables AI agents to manage scenes, sources, recording, streaming, and other OBS functionality through standardized MCP tools.

## Project Structure

```
agentic-obs/
├── main.go                 # Entry point, stdio MCP server initialization
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
├── CLAUDE.MD              # This file - AI assistant context
├── README.md              # User-facing documentation
├── config/
│   └── config.go          # Configuration management (OBS connection, etc.)
├── internal/
│   ├── mcp/
│   │   ├── server.go      # MCP server initialization and lifecycle
│   │   ├── tools.go       # MCP tool registration and handlers
│   │   └── resources.go   # MCP resource handlers (scenes as resources)
│   ├── obs/
│   │   ├── client.go      # OBS WebSocket client wrapper
│   │   ├── commands.go    # OBS command implementations
│   │   └── events.go      # OBS event handling and notification dispatch
│   └── storage/
│       ├── db.go          # SQLite database setup and migrations
│       ├── scenes.go      # Scene preset persistence
│       └── state.go       # Configuration and state management
└── scripts/
    └── setup.sh           # Development setup helpers
```

## Key Technologies

- **Go 1.25.5** - Latest stable Go version
- **MCP Go SDK v1.1.0** - `github.com/modelcontextprotocol/go-sdk`
- **goobs v1.5.6** - `github.com/andreykaipov/goobs` (OBS WebSocket 5.5.6 client)
- **modernc.org/sqlite** - Pure Go SQLite driver (CGo-free)
- **stdio transport** - MCP communication over standard input/output

## Development Setup

### Prerequisites

- Go 1.25.5 or later
- OBS Studio 28+ (with built-in WebSocket server)
- Git

### Installation

```bash
# Clone the repository
git clone <repository-url>
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
# Run directly
go run main.go

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
    ↕ WebSocket (port 4455)
OBS Studio (obs-websocket)
    ↕ SQLite (local)
Configuration & Presets
```

### Design Decisions

1. **Storage Strategy (Moderate)**: SQLite stores connection config and user-defined scene/source presets. Future expansion will include comprehensive state caching and analytics.

2. **SQLite Driver**: Using `modernc.org/sqlite` (pure Go) for simplified cross-compilation and single-user performance requirements.

3. **OBS Connection (Persistent 1:1)**: Maintains a single persistent connection to one OBS instance with auto-reconnect. Future: multi-instance support.

4. **Error Handling (Contextual)**: Errors include diagnostic context and actionable suggestions for users.

5. **Scenes as MCP Resources**: OBS scenes are exposed as MCP resources, enabling server-initiated notifications when scenes are created, deleted, or modified. This allows AI clients to stay synchronized with OBS state changes.

6. **MCP Resource Notifications**: Server monitors OBS events (SceneCreated, SceneRemoved, CurrentProgramSceneChanged) and sends `notifications/resources/updated` and `notifications/resources/list_changed` to clients.

7. **Authentication**: OBS WebSocket password stored in SQLite (local, single-user). Auto-detected on first run with defaults (localhost:4455).

8. **Setup Experience**: Auto-detect defaults, interactive prompt on failure, persist successful config. Future: TUI and web-based setup.

### MCP Resources

**OBS Scenes** are exposed as MCP resources with the following structure:

- **Resource URI Pattern**: `obs://scene/{scene_name}`
- **Resource Type**: `obs-scene`
- **Content**: JSON representation of scene configuration (sources, settings)
- **Notifications**:
  - `notifications/resources/list_changed` - When scenes are created/deleted
  - `notifications/resources/updated` - When a specific scene is modified or becomes active

**Resource Operations:**
- `resources/list` - List all available scenes
- `resources/read` - Get detailed scene configuration
- `resources/subscribe` - Subscribe to scene change notifications (future)

### MCP Tools Implemented

Core tools for OBS control:
- **Scene Management**: `set_current_scene`, `create_scene`, `remove_scene`
- **Recording**: `start_recording`, `stop_recording`, `get_recording_status`
- **Streaming**: `start_streaming`, `stop_streaming`, `get_streaming_status`
- **Sources**: `list_sources`, `toggle_source_visibility`, `get_source_settings`
- **Audio**: `get_input_mute`, `toggle_input_mute`, `set_input_volume`
- **Status**: `get_obs_status` (overall OBS connection and state)

**Note:** Scene listing moved from tools to resources (`resources/list`), enabling notification support.

## Important Context for AI Assistants

### Code Style

- Follow standard Go conventions (gofmt, go vet)
- Use meaningful variable names
- Keep functions focused and concise
- Error handling: always check errors, provide context
- Use internal/ for private packages not meant for external import

### Project Phases

**Phase 1 (Current Baseline):**
- Basic MCP server with stdio transport
- Scenes exposed as MCP resources with notifications
- OBS event monitoring and notification dispatch
- Core OBS tools (recording, streaming, sources, audio)
- SQLite config storage
- Persistent OBS connection with auto-reconnect
- Auto-detection of OBS connection

**Future Enhancements (Post-Baseline):**
- Multi-instance OBS support (requires architecture refactor)
- Comprehensive state caching and analytics
- Additional resource types (sources, filters, transitions)
- Interactive setup (TUI and web interface)
- Advanced preset management
- Performance monitoring

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
- [ ] Multi-target OBS support (architecture decision needed)
- [ ] Interactive installation UX (TUI + Web)
- [ ] Comprehensive state caching (expand from current moderate approach)
- [ ] Advanced preset/macro system
- [ ] Health check and monitoring endpoints
- [ ] Additional resource types (sources, filters, transitions, audio inputs)
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

**Last Updated:** 2025-12-14
**Go Version:** 1.25.5
**MCP SDK Version:** 1.1.0
**goobs Version:** 1.5.6
