# agentic-obs

A Model Context Protocol (MCP) server that enables AI assistants to control OBS Studio through the OBS WebSocket API.

## Overview

This MCP server provides AI agents (like Claude) with programmatic control over OBS Studio, enabling automated scene switching, recording control, streaming management, and more through natural language interactions.

## Features

- **Scene Management**: List, switch, create, and remove OBS scenes
- **Scene Presets**: Save and restore source visibility configurations
- **Recording Control**: Start, stop, pause, resume, and monitor recording
- **Streaming Control**: Start, stop, and monitor streaming
- **Source Management**: List and control source visibility
- **Audio Control**: Manage input mute and volume levels
- **Status Monitoring**: Query OBS connection and operational status

## Prerequisites

- **Go 1.25.5+** - [Download](https://go.dev/dl/)
- **OBS Studio 28+** - Includes built-in WebSocket server
- **Git** - For version control

## Installation

### Option 1: Install with Go (Recommended)

```bash
go install github.com/ironystock/agentic-obs@latest
```

This installs the `agentic-obs` binary to your `$GOPATH/bin` directory.

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/ironystock/agentic-obs.git
cd agentic-obs

# Install dependencies
go mod download

# Build the server
go build -o agentic-obs main.go
```

### Configure OBS Studio

1. Open OBS Studio
2. Go to **Tools → WebSocket Server Settings**
3. Enable the WebSocket server (default port: 4455)
4. Set a password (optional but recommended)
5. Note your connection details

## Usage

### Running the Server

```bash
# If installed via go install (ensure $GOPATH/bin is in your PATH)
agentic-obs

# Or run directly from source
go run main.go

# Or use a built binary
./agentic-obs
```

On first run, the server will:
1. Auto-detect OBS on `localhost:4455`
2. Prompt for connection details if auto-detection fails
3. Save successful configuration to SQLite

### Connecting to MCP Clients

This server uses stdio transport. Configure your MCP client to execute the `agentic-obs` command.

Example Claude Desktop configuration (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "obs": {
      "command": "agentic-obs"
    }
  }
}
```

**Note:** If you built from source or the binary isn't in your PATH, use the full path:
```json
{
  "mcpServers": {
    "obs": {
      "command": "/full/path/to/agentic-obs"
    }
  }
}
```

## Available MCP Tools

### Scene Management (4 tools)

| Tool | Description |
|------|-------------|
| `list_scenes` | List all available scenes and identify current scene |
| `set_current_scene` | Switch to a specific scene |
| `create_scene` | Create a new scene |
| `remove_scene` | Remove a scene |

### Recording Control (5 tools)

| Tool | Description |
|------|-------------|
| `start_recording` | Start recording |
| `stop_recording` | Stop recording |
| `pause_recording` | Pause current recording |
| `resume_recording` | Resume paused recording |
| `get_recording_status` | Check recording status and details |

### Streaming Control (3 tools)

| Tool | Description |
|------|-------------|
| `start_streaming` | Start streaming |
| `stop_streaming` | Stop streaming |
| `get_streaming_status` | Check streaming status |

### Source Management (3 tools)

| Tool | Description |
|------|-------------|
| `list_sources` | List all input sources |
| `toggle_source_visibility` | Show/hide a source in a scene |
| `get_source_settings` | Get source configuration |

### Audio Control (4 tools)

| Tool | Description |
|------|-------------|
| `get_input_mute` | Check if audio input is muted |
| `toggle_input_mute` | Toggle audio input mute state |
| `set_input_volume` | Set audio input volume (dB or multiplier) |
| `get_input_volume` | Get current volume level (dB and multiplier) |

### Scene Presets (6 tools)

| Tool | Description |
|------|-------------|
| `save_scene_preset` | Save current scene source states as a named preset |
| `list_scene_presets` | List all saved presets, optionally filter by scene |
| `get_preset_details` | Get detailed information about a specific preset |
| `apply_scene_preset` | Apply a saved preset to restore source visibility |
| `rename_scene_preset` | Rename an existing preset |
| `delete_scene_preset` | Delete a saved preset |

### Screenshot Sources (4 tools)

| Tool | Description |
|------|-------------|
| `create_screenshot_source` | Create a periodic screenshot capture source for AI visual monitoring |
| `remove_screenshot_source` | Stop and remove a screenshot capture source |
| `list_screenshot_sources` | List all configured sources with status and HTTP URLs |
| `configure_screenshot_cadence` | Update the capture interval for a screenshot source |

### Status & Monitoring (1 tool)

| Tool | Description |
|------|-------------|
| `get_obs_status` | Get overall OBS status and connection info |

**Total: 30 tools**

## Development

### Project Structure

```
agentic-obs/
├── main.go                 # Entry point
├── config/                 # Configuration management
├── internal/
│   ├── mcp/               # MCP server implementation
│   ├── obs/               # OBS WebSocket client
│   ├── storage/           # SQLite persistence
│   ├── http/              # HTTP server for screenshots
│   └── screenshot/        # Background capture manager
└── scripts/               # Development helpers
```

### Adding a New Tool

1. Define tool schema in `internal/mcp/tools.go`
2. Implement handler function
3. Register tool in server initialization
4. Add OBS command wrapper in `internal/obs/commands.go` if needed

### Running Tests

```bash
go test ./...
```

## Configuration

Configuration is stored in SQLite (`agentic-obs.db`) and includes:

- OBS WebSocket connection details (host, port, password)
- Scene and source presets
- User preferences

## Troubleshooting

### Connection Issues

**Problem**: "Failed to connect to OBS"

**Solutions**:
- Ensure OBS Studio is running
- Verify WebSocket server is enabled in OBS (Tools → WebSocket Server Settings)
- Check that port 4455 is not blocked by firewall
- Confirm password matches if authentication is enabled

### Permission Issues

**Problem**: Database write errors

**Solutions**:
- Ensure the directory is writable
- Check file permissions on `agentic-obs.db`

## Future Enhancements

- Automation rules and macros
- Multi-instance OBS support
- Interactive setup UI (TUI and web)
- Real-time event notifications

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

[Add your license here]

## Resources

- [Model Context Protocol](https://modelcontextprotocol.io)
- [OBS WebSocket Protocol](https://github.com/obsproject/obs-websocket)
- [goobs Library](https://github.com/andreykaipov/goobs)

---

**Built with:**
- Go 1.25.5
- MCP Go SDK v1.1.0
- goobs v1.5.6
