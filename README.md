# agentic-obs

A Model Context Protocol (MCP) server that enables AI assistants to control OBS Studio through the OBS WebSocket API.

## Overview

This MCP server provides AI agents (like Claude) with programmatic control over OBS Studio, enabling automated scene switching, recording control, streaming management, and more through natural language interactions.

## Features

- **Scene Management**: List, switch, and manage OBS scenes
- **Recording Control**: Start, stop, and check recording status
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

| Tool | Description |
|------|-------------|
| `list_scenes` | List all available scenes in OBS |
| `get_current_scene` | Get the currently active scene |
| `set_current_scene` | Switch to a specific scene |
| `start_recording` | Start recording |
| `stop_recording` | Stop recording |
| `get_recording_status` | Check if recording is active |
| `start_streaming` | Start streaming |
| `stop_streaming` | Stop streaming |
| `get_streaming_status` | Check if streaming is active |
| `list_sources` | List all sources |
| `toggle_source_visibility` | Show/hide a source |
| `get_source_settings` | Get source configuration |
| `get_input_mute` | Check if input is muted |
| `toggle_input_mute` | Mute/unmute an input |
| `set_input_volume` | Set input volume level |
| `get_obs_status` | Get overall OBS connection status |

## Development

### Project Structure

```
agentic-obs/
├── main.go                 # Entry point
├── config/                 # Configuration management
├── internal/
│   ├── mcp/               # MCP server implementation
│   ├── obs/               # OBS WebSocket client
│   └── storage/           # SQLite persistence
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

- Multi-instance OBS support
- Interactive setup UI (TUI and web)
- Real-time event notifications
- Advanced preset and macro system
- Performance monitoring and analytics

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
