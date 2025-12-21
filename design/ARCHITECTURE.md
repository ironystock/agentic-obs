# Architecture

This document describes the system architecture of agentic-obs.

## System Overview

agentic-obs is an MCP (Model Context Protocol) server that bridges AI assistants with OBS Studio. It provides 72 tools, 4 resource types, and 13 prompts for programmatic OBS control.

```
┌─────────────────────────────────────────────────────────────────┐
│                        AI Assistant                              │
│                    (Claude, GPT, etc.)                          │
└────────────────────────────┬────────────────────────────────────┘
                             │ MCP Protocol (stdio/JSON-RPC)
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                      agentic-obs                                 │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    MCP Layer                             │    │
│  │  ┌──────────┐  ┌───────────┐  ┌─────────┐  ┌──────────┐ │    │
│  │  │  Tools   │  │ Resources │  │ Prompts │  │Completions│ │    │
│  │  │  (72)    │  │   (4)     │  │  (13)   │  │          │ │    │
│  │  └────┬─────┘  └─────┬─────┘  └────┬────┘  └────┬─────┘ │    │
│  └───────┼──────────────┼─────────────┼───────────┼────────┘    │
│          │              │             │           │              │
│  ┌───────↓──────────────↓─────────────↓───────────↓────────┐    │
│  │                   OBS Client                             │    │
│  │            (WebSocket connection manager)                │    │
│  └────────────────────────┬────────────────────────────────┘    │
│                           │                                      │
│  ┌────────────────────────↓────────────────────────────────┐    │
│  │                    Storage Layer                         │    │
│  │  ┌──────────┐  ┌────────────┐  ┌──────────────────────┐ │    │
│  │  │  Config  │  │   Presets  │  │   Screenshots/History│ │    │
│  │  └──────────┘  └────────────┘  └──────────────────────┘ │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    HTTP Server                           │    │
│  │  ┌──────────┐  ┌────────────┐  ┌────────────────────┐   │    │
│  │  │ Web UI   │  │  REST API  │  │ Screenshot Serving │   │    │
│  │  └──────────┘  └────────────┘  └────────────────────┘   │    │
│  └─────────────────────────────────────────────────────────┘    │
└────────────────────────────┬────────────────────────────────────┘
                             │ WebSocket (port 4455)
                             ↓
┌─────────────────────────────────────────────────────────────────┐
│                        OBS Studio                                │
│                    (obs-websocket 5.x)                          │
└─────────────────────────────────────────────────────────────────┘
```

## Component Responsibilities

| Component | Package | Responsibility |
|-----------|---------|----------------|
| **MCP Server** | `internal/mcp/server.go` | Lifecycle, stdio transport, notification dispatch |
| **Tools** | `internal/mcp/tools.go` | 72 tool handlers organized in 8 groups + meta |
| **Resources** | `internal/mcp/resources.go` | Scene, screenshot, preset resource handlers |
| **Prompts** | `internal/mcp/prompts.go` | 13 workflow prompt definitions |
| **Completions** | `internal/mcp/completions.go` | Autocomplete for arguments and URIs |
| **OBS Client** | `internal/obs/client.go` | WebSocket connection, reconnection, events |
| **OBS Commands** | `internal/obs/commands.go` | OBS API method implementations |
| **Storage** | `internal/storage/db.go` | SQLite database, migrations |
| **HTTP Server** | `internal/http/server.go` | REST API, static files, screenshot serving |
| **Screenshot Mgr** | `internal/screenshot/manager.go` | Background capture, cadence management |

## Data Flow

### Tool Invocation

```
AI Client                MCP Server              OBS Client           OBS Studio
    │                        │                       │                     │
    │─── tools/call ────────>│                       │                     │
    │                        │── parse input ───────>│                     │
    │                        │                       │── WebSocket req ───>│
    │                        │                       │<── WebSocket res ───│
    │                        │<── result ────────────│                     │
    │<── tool result ────────│                       │                     │
    │                        │                       │                     │
```

### Resource Notification

```
OBS Studio              OBS Client              MCP Server           AI Client
    │                        │                       │                     │
    │── event (scene ───────>│                       │                     │
    │   changed)             │── handleEvent ───────>│                     │
    │                        │                       │── notify/updated ──>│
    │                        │                       │   (or list_changed) │
    │                        │                       │                     │
```

### Screenshot Capture

```
Screenshot Manager           OBS Client           Storage              HTTP Server
    │                            │                   │                      │
    │── ticker fires ───────────>│                   │                      │
    │                            │── screenshot ─────│                      │
    │<── image data ─────────────│                   │                      │
    │                            │                   │                      │
    │── store image ─────────────────────────────────>│                     │
    │                            │                   │                      │
    │                            │                   │<── GET /screenshot ──│
    │                            │                   │── image data ───────>│
```

## Package Dependencies

```
main.go
    ↓
┌───┴───┐
│config │ ← Configuration structs, loading
└───┬───┘
    ↓
┌─────────────────────────────────────────────────────┐
│                  internal/mcp                        │
│  server.go → tools.go → resources.go → prompts.go   │
└───────────────────────┬─────────────────────────────┘
                        ↓
          ┌─────────────┼─────────────┐
          ↓             ↓             ↓
    internal/obs   internal/storage   internal/http
          ↓             ↓
      goobs lib    modernc.org/sqlite
```

## MCP Protocol Integration

### Tools (72 total)

| Group | Tools | Description |
|-------|-------|-------------|
| **Core** | 25 | Scene management, recording, streaming, virtual cam, replay buffer, studio mode, hotkeys |
| **Sources** | 3 | Source visibility and settings |
| **Audio** | 4 | Volume and mute control |
| **Layout** | 6 | Scene preset management |
| **Visual** | 4 | Screenshot source control |
| **Design** | 14 | Source creation and transforms |
| **Filters** | 7 | Filter creation and management |
| **Transitions** | 5 | Scene transition control |
| **Meta** | 4 | Help, tool config (always enabled) |

### Resources (4 types)

| Type | URI Pattern | Content |
|------|-------------|---------|
| **Scenes** | `obs://scene/{name}` | Scene configuration JSON |
| **Screenshots** | `obs://screenshot/{name}` | Binary image data |
| **Screenshot URLs** | `obs://screenshot-url/{name}` | HTTP URL for image |
| **Presets** | `obs://preset/{name}` | Preset configuration JSON |

### Prompts (13 workflows)

Pre-built conversation starters for common tasks:
- Stream management: `stream-launch`, `stream-teardown`
- Diagnostics: `health-check`, `audio-check`, `visual-check`, `problem-detection`
- Recording: `recording-workflow`
- Scene design: `scene-designer`, `source-management`, `scene-organizer`
- Quick access: `quick-status`, `preset-switcher`, `visual-setup`

## Web UI Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    HTTP Server (port 8765)                       │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    Static Files                          │    │
│  │  /              → index.html (dashboard)                 │    │
│  │  /ui/scenes     → scene_preview.html                     │    │
│  │  /ui/audio      → audio_mixer.html                       │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    REST API                              │    │
│  │  GET  /api/status          → OBS connection status       │    │
│  │  GET  /api/history         → Action history              │    │
│  │  GET  /api/history/stats   → Statistics                  │    │
│  │  GET  /api/screenshots     → Screenshot sources          │    │
│  │  GET  /api/config          → Configuration               │    │
│  │  POST /api/config          → Update configuration        │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    UI API                                │    │
│  │  GET  /ui/api/scenes              → Scene list + thumbs  │    │
│  │  POST /ui/api/scenes/{name}       → Switch scene         │    │
│  │  GET  /ui/api/audio               → Audio inputs         │    │
│  │  POST /ui/api/audio/{name}/mute   → Toggle mute          │    │
│  │  POST /ui/api/audio/{name}/volume → Set volume           │    │
│  │  GET  /ui/scene-thumbnail/{name}  → Scene thumbnail      │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │               Interface Pattern                          │    │
│  │  StatusProvider ──→ GetStatus(), GetScenes(), etc.       │    │
│  │  ActionExecutor ──→ SetCurrentScene(), ToggleMute(), etc.│    │
│  │                                                          │    │
│  │  MCP Server implements both interfaces                   │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

## Storage Schema

```sql
-- Configuration key-value store
CREATE TABLE config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Scene presets (source visibility snapshots)
CREATE TABLE scene_presets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    scene_name TEXT NOT NULL,
    sources TEXT,  -- JSON array
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- Screenshot sources (periodic capture config)
CREATE TABLE screenshot_sources (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    source_name TEXT NOT NULL,
    cadence_ms INTEGER NOT NULL DEFAULT 5000,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Screenshot images (binary blobs)
CREATE TABLE screenshots (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_id INTEGER NOT NULL,
    format TEXT NOT NULL,
    data BLOB NOT NULL,
    captured_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (source_id) REFERENCES screenshot_sources(id)
);

-- Action history (for auditing)
CREATE TABLE action_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    action TEXT NOT NULL,
    tool_name TEXT,
    input TEXT,
    output TEXT,
    success BOOLEAN,
    duration_ms INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Threading Model

- **Main goroutine**: MCP server stdio transport
- **OBS event goroutine**: Receives WebSocket events, dispatches notifications
- **Screenshot workers**: Per-source goroutines for periodic capture
- **Thumbnail cache cleanup**: Background goroutine for expired cache entries
- **HTTP server**: Standard Go HTTP server (connection-per-goroutine)

## Security Considerations

See [decisions/005-auth-storage.md](decisions/005-auth-storage.md) for authentication storage rationale.

**Key Points:**
- OBS password stored unencrypted in SQLite (local deployment only)
- HTTP server binds to localhost by default
- Path traversal prevention on screenshot endpoints
- No external network access required
