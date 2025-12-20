# agentic-obs

A Model Context Protocol (MCP) server that enables AI assistants to control OBS Studio through the OBS WebSocket API.

## Overview

This MCP server provides AI agents (like Claude) with programmatic control over OBS Studio, enabling automated scene switching, recording control, streaming management, and more through natural language interactions.

## Features

- **57 MCP Tools**: Comprehensive control over OBS Studio operations in 8 tool groups
- **Scene Management**: List, switch, create, and remove OBS scenes
- **Scene Presets**: Save and restore source visibility configurations
- **Recording Control**: Start, stop, pause, resume, and monitor recording
- **Streaming Control**: Start, stop, and monitor streaming
- **Source Management**: List and control source visibility
- **Audio Control**: Manage input mute and volume levels
- **Screenshot Sources**: AI visual monitoring with periodic image capture
- **Agentic Scene Design**: Create and manipulate sources (text, image, color, browser, media)
- **Help & Discovery**: Built-in help tool with topic-based guidance
- **Status Monitoring**: Query OBS connection and operational status
- **4 MCP Resources**: Scenes, screenshots, screenshot URLs, and presets exposed as resources
- **13 MCP Prompts**: Pre-built workflows for common tasks and diagnostics
- **MCP Completions**: Autocomplete for prompt arguments and resource URIs
- **Claude Skills**: Shareable skill packages for advanced AI orchestration
- **TUI Dashboard**: Terminal interface for status, config, and history (`--tui` flag)
- **Web Dashboard**: Browser-based dashboard at `http://localhost:8765/`

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

# Build with Make (recommended - includes version info)
make build

# Or build directly with Go
go build -o agentic-obs .

# Verify the build
./agentic-obs --version
```

For cross-platform builds, release automation, and advanced build options, see the BUILD documentation.

### Configure OBS Studio

1. Open OBS Studio
2. Go to **Tools → WebSocket Server Settings**
3. Enable the WebSocket server (default port: 4455)
4. Set a password (optional but recommended)
5. Note your connection details

## Usage

### Running the Server

```bash
# MCP server mode (default) - if installed via go install
agentic-obs

# TUI dashboard mode - terminal interface for monitoring
agentic-obs --tui
agentic-obs -t

# Or run directly from source
go run main.go
go run main.go --tui

# Or use a built binary
./agentic-obs
```

### TUI Dashboard

The TUI dashboard provides a terminal-based interface with four views:
- **Status**: OBS connection status, server info, statistics
- **Config**: Current configuration settings
- **History**: Action history log with scrolling
- **Docs**: Embedded documentation with terminal rendering

Navigate with `1/2/3/4` keys or Tab, press `q` to quit.

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

### Help & Discovery (1 tool)

| Tool | Description |
|------|-------------|
| `help` | Get detailed help on tools, resources, prompts, workflows, or troubleshooting |

### Scene Design (14 tools)

Enable AI to create and manipulate OBS sources programmatically.

#### Source Creation

| Tool | Description |
|------|-------------|
| `create_text_source` | Create a text/label source with customizable font and color |
| `create_image_source` | Create an image source from a file path |
| `create_color_source` | Create a solid color source |
| `create_browser_source` | Create a browser source for web content |
| `create_media_source` | Create a media/video source from a file |

#### Layout Control

| Tool | Description |
|------|-------------|
| `set_source_transform` | Set position, scale, and rotation of a source |
| `get_source_transform` | Get current transform properties |
| `set_source_crop` | Set crop values for a source |
| `set_source_bounds` | Set bounds type and size for a source |
| `set_source_order` | Set z-order index (front-to-back ordering) |

#### Advanced

| Tool | Description |
|------|-------------|
| `set_source_locked` | Lock/unlock a source to prevent changes |
| `duplicate_source` | Duplicate a source within or across scenes |
| `remove_source` | Remove a source from a scene |
| `list_input_kinds` | List all available input source types |

### Filters (7 tools)

| Tool | Description |
|------|-------------|
| `list_source_filters` | List all filters applied to a source |
| `get_source_filter` | Get filter details and settings |
| `create_source_filter` | Add a new filter to a source (color correction, noise suppression, etc.) |
| `remove_source_filter` | Remove a filter from a source |
| `toggle_source_filter` | Enable/disable a filter |
| `set_source_filter_settings` | Modify filter configuration |
| `list_filter_kinds` | List all available filter types |

### Transitions (5 tools)

| Tool | Description |
|------|-------------|
| `list_transitions` | List all available transitions and current one |
| `get_current_transition` | Get current transition details |
| `set_current_transition` | Change the active transition (Cut, Fade, Swipe, etc.) |
| `set_transition_duration` | Set transition duration in milliseconds |
| `trigger_transition` | Trigger studio mode transition (preview to program) |

**Total: 57 tools in 8 groups** (Core, Sources, Audio, Layout, Visual, Design, Filters, Transitions) + Help

## MCP Resources

The server exposes OBS data as MCP resources for efficient access and monitoring:

| Resource Type | URI Pattern | Content Type | Description |
|---------------|-------------|--------------|-------------|
| Scenes | `obs://scene/{name}` | `obs-scene` (JSON) | Scene configuration with sources and settings |
| Screenshots | `obs://screenshot/{name}` | `image/png` or `image/jpeg` | Binary screenshot images from capture sources |
| Screenshot URLs | `obs://screenshot-url/{name}` | `text/plain` | HTTP URL for screenshot image access |
| Presets | `obs://preset/{name}` | `obs-preset` (JSON) | Scene preset configurations with source visibility |

**Usage:**
- `resources/list` - List all available resources
- `resources/read` - Get detailed resource content
- Resources support notifications for real-time updates
- **Completions**: Autocomplete available for resource URI names

## MCP Prompts

Pre-built workflow prompts guide AI assistants through common OBS tasks:

| Prompt | Arguments | Purpose |
|--------|-----------|---------|
| `stream-launch` | none | Pre-stream checklist and setup |
| `stream-teardown` | none | End-stream cleanup workflow |
| `audio-check` | none | Audio verification and diagnostics |
| `visual-check` | screenshot_source | Visual layout analysis |
| `health-check` | none | Comprehensive OBS diagnostics |
| `problem-detection` | screenshot_source | Automated issue detection |
| `preset-switcher` | preset_name (optional) | Scene preset management |
| `recording-workflow` | none | Recording session guidance |
| `scene-organizer` | none | Scene organization and cleanup |
| `quick-status` | none | Brief status summary |
| `scene-designer` | scene_name, action (optional) | Visual layout creation with Design tools |
| `source-management` | scene_name | Manage source visibility and properties |
| `visual-setup` | monitor_scene (optional) | Configure screenshot monitoring |

Prompts combine multiple tools into cohesive workflows with best practices built-in.
**Completions**: Autocomplete available for prompt arguments (preset names, screenshot sources, scene names).

## Development

### Project Structure

```
agentic-obs/
├── main.go                 # Entry point (MCP server or TUI)
├── config/                 # Configuration management
├── internal/
│   ├── mcp/               # MCP server implementation (57 tools)
│   ├── obs/               # OBS WebSocket client
│   ├── storage/           # SQLite persistence
│   ├── http/              # HTTP server for screenshots and dashboard
│   ├── screenshot/        # Background capture manager
│   └── tui/               # Terminal UI dashboard
├── skills/                 # Claude Skills packages
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

## Claude Skills

The `skills/` directory contains shareable Claude Skills packages that teach AI assistants how to orchestrate agentic-obs tools effectively:

| Skill | Purpose |
|-------|---------|
| `streaming-assistant` | Complete streaming workflows (pre-stream, live, teardown) |
| `scene-designer` | Visual layout creation with 14 Design tools |
| `audio-engineer` | Audio optimization and troubleshooting |
| `preset-manager` | Preset lifecycle and organization |

Skills use progressive disclosure for token-efficient guidance. See `skills/README.md` for installation.

## Future Enhancements

- Automation rules and macros
- Multi-instance OBS support
- Real-time event notifications
- Additional resource types (filters, transitions)

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
- bubbletea/lipgloss (TUI)
