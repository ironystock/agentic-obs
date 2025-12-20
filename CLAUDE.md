# CLAUDE.MD

Context for AI assistants working on the agentic-obs project.

## Project Overview

**agentic-obs** is an MCP (Model Context Protocol) server providing AI assistants with programmatic control over OBS Studio via the WebSocket API.

**Current Status:** 45 Tools | 4 Resources | 13 Prompts | 4 Skills

## Project Structure

```
agentic-obs/
├── main.go                 # Entry point (MCP server or TUI dashboard)
├── config/config.go        # Configuration management
├── internal/
│   ├── mcp/
│   │   ├── server.go       # MCP server lifecycle
│   │   ├── tools.go        # Tool registration (45 tools)
│   │   ├── resources.go    # Resource handlers (4 types)
│   │   ├── prompts.go      # Prompt definitions (13 prompts)
│   │   ├── completions.go  # Autocomplete handler
│   │   ├── help.go         # Help tool
│   │   └── interfaces.go   # OBSClient interface
│   ├── obs/
│   │   ├── client.go       # WebSocket client
│   │   ├── commands.go     # OBS command implementations
│   │   └── events.go       # Event handling
│   ├── storage/            # SQLite persistence
│   ├── http/               # Web dashboard & API
│   ├── screenshot/         # Background capture
│   └── tui/                # Terminal dashboard
├── skills/                 # Claude Skills (4 packages)
├── examples/               # Usage examples
├── docs/                   # Additional documentation
└── design/                 # Architecture & decisions
```

## Key Technologies

| Component | Package | Version |
|-----------|---------|---------|
| MCP SDK | `github.com/modelcontextprotocol/go-sdk` | 1.1.0 |
| OBS Client | `github.com/andreykaipov/goobs` | 1.5.6 |
| SQLite | `modernc.org/sqlite` | latest |
| TUI | `github.com/charmbracelet/bubbletea` | 1.3.3 |

## Communication Flow

```
AI Assistant (Claude)
    ↕ stdio (JSON-RPC)
MCP Server (agentic-obs)
    ├─ Tools (45) ─────────┐
    ├─ Resources (4) ──────┼─→ OBS Client ─→ OBS Studio
    ├─ Prompts (13) ───────┘      ↕ WebSocket (4455)
    └─ Storage ─────────────→ SQLite

Web Dashboard (localhost:8765)
    ↕ HTTP/REST
HTTP Server
    └─ StatusProvider/ActionExecutor interfaces
```

## Development

### Running

```bash
# MCP server (default)
go run main.go

# TUI dashboard
go run main.go --tui

# Build
go build -o agentic-obs
```

### Testing

```bash
go test ./...
```

## Code Style

- Follow standard Go conventions (gofmt, go vet)
- Use meaningful variable names
- Keep functions focused and concise
- Always check errors, provide context
- Use `internal/` for private packages

## Common Tasks

### Adding a New Tool

1. Define schema in `internal/mcp/tools.go`:
```go
{
    Name: "my_tool",
    Description: "Does something useful",
    InputSchema: myToolSchema,
}
```

2. Implement handler:
```go
func (s *Server) handleMyTool(ctx context.Context, params json.RawMessage) (*mcpsdk.CallToolResult, error) {
    // Implementation
}
```

3. Register in `registerToolHandlers()`

4. Add OBS command in `internal/obs/commands.go` if needed

### Adding a New Resource

1. Define in `internal/mcp/resources.go`
2. Implement `handleResourceRead()` case
3. Add to `handleResourcesList()`
4. Add OBS event handlers for notifications

### Adding a New Prompt

1. Define in `internal/mcp/prompts.go`:
```go
{
    Name: "my-prompt",
    Description: "Description",
    Arguments: []mcpsdk.PromptArgument{...},
}
```

2. Add message generation in `handleGetPrompt()`

### Updating Dependencies

```bash
go get -u github.com/modelcontextprotocol/go-sdk
go get -u github.com/andreykaipov/goobs
go mod tidy
```

## MCP Capabilities Summary

### Tools (45 in 6 groups)

| Group | Count | Examples |
|-------|-------|----------|
| Core | 13 | `list_scenes`, `start_recording`, `get_obs_status` |
| Sources | 3 | `list_sources`, `toggle_source_visibility` |
| Audio | 4 | `toggle_input_mute`, `set_input_volume` |
| Layout | 6 | `save_scene_preset`, `apply_scene_preset` |
| Visual | 4 | `create_screenshot_source`, `list_screenshot_sources` |
| Design | 14 | `create_text_source`, `set_source_transform` |
| Help | 1 | `help` (always enabled) |

### Resources (4 types)

| Type | URI Pattern | Content |
|------|-------------|---------|
| Scenes | `obs://scene/{name}` | Scene configuration JSON |
| Screenshots | `obs://screenshot/{name}` | Binary image data |
| Screenshot URLs | `obs://screenshot-url/{name}` | HTTP URL string |
| Presets | `obs://preset/{name}` | Preset configuration JSON |

### Prompts (13 workflows)

`stream-launch`, `stream-teardown`, `audio-check`, `visual-check`, `health-check`, `problem-detection`, `preset-switcher`, `recording-workflow`, `scene-organizer`, `quick-status`, `scene-designer`, `source-management`, `visual-setup`

## Key Design Decisions

For detailed rationale, see [design/decisions/](design/decisions/).

1. **Pure Go SQLite** (`modernc.org/sqlite`) - No CGO for easy cross-compilation
2. **Persistent OBS Connection** - Single 1:1 connection with auto-reconnect
3. **Scenes as Resources** - Enable MCP notifications for state synchronization
4. **Tool Groups** - Configurable categories (Core, Visual, Layout, Audio, Sources, Design)
5. **Interface Pattern** - StatusProvider/ActionExecutor for web UI decoupling

## Documentation Links

| Document | Purpose |
|----------|---------|
| [README.md](README.md) | User documentation |
| [CHANGELOG.md](CHANGELOG.md) | Version history |
| [docs/BUILD.md](docs/BUILD.md) | Build, test, release |
| [design/ARCHITECTURE.md](design/ARCHITECTURE.md) | System diagrams |
| [design/ROADMAP.md](design/ROADMAP.md) | Future plans |
| [design/decisions/](design/decisions/) | ADRs |
| [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) | Common issues |
| [docs/TOOLS.md](docs/TOOLS.md) | Detailed tool reference |

## Testing Approach

- **Storage tests**: `internal/storage/*_test.go`
- **MCP tests**: `internal/mcp/*_test.go` (uses mock OBS client)
- **Mock client**: `internal/mcp/testutil/mock_obs.go`

## Known Constraints

- Single OBS instance only (multi-instance planned)
- OBS password stored unencrypted in SQLite (local deployment)
- Resource notifications require persistent OBS connection

## External References

- [MCP Specification](https://modelcontextprotocol.io)
- [OBS WebSocket Protocol](https://github.com/obsproject/obs-websocket/blob/master/docs/generated/protocol.md)
- [goobs Documentation](https://pkg.go.dev/github.com/andreykaipov/goobs)

---

**Last Updated:** 2025-12-19 | **Go:** 1.25.5 | **MCP SDK:** 1.1.0 | **goobs:** 1.5.6
