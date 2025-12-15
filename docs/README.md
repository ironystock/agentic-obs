# agentic-obs Documentation

Complete documentation for controlling OBS Studio with AI through the Model Context Protocol.

## Getting Started

- **[Quick Start Guide](QUICKSTART.md)** - Get up and running in 10 minutes
- **[Tool Reference](TOOLS.md)** - Comprehensive reference for all 19 MCP tools

## Examples

- **[Natural Language Prompts](../examples/prompts/)** - Conversational examples for AI assistants
- **[API Examples](../examples/api/)** - Technical JSON-RPC request/response examples

## Documentation

| Document | Description |
|----------|-------------|
| [QUICKSTART.md](QUICKSTART.md) | Step-by-step installation and setup guide |
| [TOOLS.md](TOOLS.md) | Complete reference for all 19 tools with examples |
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

### Recording Control (5 tools)
Manage recordings - start, stop, pause, resume, and check status

### Streaming Control (3 tools)
Control streaming - start, stop, and monitor streams

### Source Management (3 tools)
Manage input sources - list, toggle visibility, configure settings

### Audio Control (3 tools)
Control audio inputs - mute state and volume levels

### Status & Monitoring (1 tool)
Get comprehensive OBS status and connection info

---

**Total: 19 MCP tools available**

For detailed information on each tool, see [TOOLS.md](TOOLS.md).
