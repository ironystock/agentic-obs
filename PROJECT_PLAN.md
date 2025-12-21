# agentic-obs Project Status

**Status:** Active Development
**Version:** 0.11.0
**Updated:** 2025-12-21

---

## Quick Links

| Document | Description |
|----------|-------------|
| [README.md](README.md) | User documentation, installation, usage |
| [CLAUDE.md](CLAUDE.md) | AI assistant context, development guidelines |
| [CHANGELOG.md](CHANGELOG.md) | Version history, phase details |
| [design/ARCHITECTURE.md](design/ARCHITECTURE.md) | System diagrams, component responsibilities |
| [design/ROADMAP.md](design/ROADMAP.md) | Future enhancements, research topics |
| [design/decisions/](design/decisions/) | Architecture Decision Records (ADRs) |

---

## Current Metrics

| Metric | Count |
|--------|-------|
| **MCP Tools** | 69 |
| **MCP Resources** | 4 |
| **MCP Prompts** | 13 |
| **Claude Skills** | 4 |

---

## Technology Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| **Language** | Go | 1.25.5 |
| **MCP SDK** | go-sdk | 1.1.0 |
| **OBS Client** | goobs | 1.5.6 |
| **Database** | modernc.org/sqlite | latest |
| **TUI** | bubbletea | 1.3.3 |

---

## Project Structure

```
agentic-obs/
├── main.go                 # Entry point
├── config/                 # Configuration management
├── internal/
│   ├── mcp/               # MCP server (tools, resources, prompts)
│   ├── obs/               # OBS WebSocket client
│   ├── storage/           # SQLite persistence
│   ├── http/              # Web dashboard & API
│   ├── screenshot/        # Background capture manager
│   └── tui/               # Terminal dashboard
├── skills/                 # Claude Skills packages
├── examples/               # Usage examples
├── docs/                   # Additional documentation
└── design/                 # Architecture & decisions
```

---

## Development

### Build & Run

```bash
# Build
go build -o agentic-obs

# Run MCP server (default)
./agentic-obs

# Run TUI dashboard
./agentic-obs --tui
```

### Testing

```bash
go test ./...
```

### Adding Features

See [CLAUDE.md](CLAUDE.md) for:
- Adding new MCP tools
- Adding new MCP resources
- Adding new MCP prompts

---

## References

- [MCP Specification](https://modelcontextprotocol.io)
- [OBS WebSocket Protocol](https://github.com/obsproject/obs-websocket/blob/master/docs/generated/protocol.md)
- [goobs Documentation](https://pkg.go.dev/github.com/andreykaipov/goobs)
