# agentic-obs Documentation

Complete documentation for controlling OBS Studio with AI through the Model Context Protocol.

## Getting Started

- **[Quick Start Guide](QUICKSTART.md)** - Get up and running in 10 minutes
- **[Tool Reference](TOOLS.md)** - Comprehensive reference for all 30 MCP tools
- **[Screenshot Guide](SCREENSHOTS.md)** - AI visual monitoring of your stream

## Examples

- **[Natural Language Prompts](../examples/prompts/)** - Conversational examples for AI assistants
- **[API Examples](../examples/api/)** - Technical JSON-RPC request/response examples

## Documentation

| Document | Description |
|----------|-------------|
| [QUICKSTART.md](QUICKSTART.md) | Step-by-step installation and setup guide |
| [TOOLS.md](TOOLS.md) | Complete reference for all 30 tools with examples |
| [SCREENSHOTS.md](SCREENSHOTS.md) | Detailed guide to AI visual monitoring |
| [../README.md](../README.md) | Project overview and features |
| [../CLAUDE.md](../CLAUDE.md) | AI assistant context and architecture |
| [../PROJECT_PLAN.md](../PROJECT_PLAN.md) | Detailed project roadmap and design decisions |

## Quick Links

### For End Users
1. [Installation](QUICKSTART.md#installation)
2. [OBS Setup](QUICKSTART.md#configure-obs-websocket)
3. [Claude Desktop Config](QUICKSTART.md#connecting-to-claude-desktop)
4. [First Commands](QUICKSTART.md#your-first-commands)

### For Developers
1. [Architecture](../CLAUDE.md#architecture--patterns)
2. [Adding Tools](../CLAUDE.md#adding-a-new-mcp-tool)
3. [Code Structure](../CLAUDE.md#project-structure)
4. [Design Decisions](../PROJECT_PLAN.md)

### Troubleshooting
- [Common Issues](QUICKSTART.md#common-first-time-issues)
- [Connection Problems](QUICKSTART.md#failed-to-connect-to-obs)
- [GitHub Issues](https://github.com/ironystock/agentic-obs/issues)

## Tool Categories

### Scene Management (4 tools)
Control OBS scenes - list, switch, create, and remove scenes

### Scene Presets (6 tools)
Save and restore source visibility configurations for quick layout switching

### Recording Control (5 tools)
Manage recordings - start, stop, pause, resume, and check status

### Streaming Control (3 tools)
Control streaming - start, stop, and monitor streams

### Source Management (3 tools)
Manage input sources - list, toggle visibility, configure settings

### Audio Control (4 tools)
Control audio inputs - mute state and volume levels

### Screenshot Sources (4 tools)
Enable AI visual monitoring of your stream output - see [SCREENSHOTS.md](SCREENSHOTS.md)

### Status & Monitoring (1 tool)
Get comprehensive OBS status and connection info

---

**Total: 30 MCP tools available**

For detailed information on each tool, see [TOOLS.md](TOOLS.md).

---

## MCP Resources

Resources provide structured access to OBS data via URI patterns:

| Resource | URI Pattern | Content Type | Description |
|----------|-------------|--------------|-------------|
| **Scenes** | `obs://scene/{name}` | JSON | Scene configuration and sources |
| **Screenshots** | `obs://screenshot/{name}` | Binary (PNG/JPEG) | Latest captured screenshot |
| **Presets** | `obs://preset/{name}` | JSON | Saved source visibility states |

**Total: 3 MCP resources available**

---

## MCP Prompts

Prompts provide pre-built workflow templates for common OBS operations:

| Prompt | Arguments | Description |
|--------|-----------|-------------|
| `stream-launch` | - | Pre-stream checklist and setup verification |
| `stream-teardown` | - | End-stream cleanup sequence |
| `audio-check` | - | Comprehensive audio source verification |
| `visual-check` | screenshot_source | AI visual analysis of stream layout |
| `health-check` | - | Overall OBS status diagnostic |
| `problem-detection` | screenshot_source | Identify visual issues in stream |
| `preset-switcher` | preset_name (optional) | Scene preset management |
| `recording-workflow` | - | Recording session management |
| `scene-organizer` | - | Scene inventory and organization |
| `quick-status` | - | Brief OBS status summary |

**Total: 10 MCP prompts available**
