# agentic-obs Project Plan

## Executive Summary

This document outlines the architecture, decisions, and implementation roadmap for the **agentic-obs** MCP server - a Go-based bridge between AI assistants and OBS Studio.

**Created:** 2025-12-14
**Updated:** 2025-12-18
**Status:** ✅ Phase 7 Complete - 45 Tools, 4 Resources, 13 Prompts, MCP Completions, Help Tool, Claude Skills

---

## Project Goals

Enable AI assistants (like Claude) to programmatically control OBS Studio through a standardized Model Context Protocol (MCP) interface, allowing natural language interactions for scene management, recording, streaming, and source control.

---

## Technology Stack

| Component | Technology | Version | Rationale |
|-----------|-----------|---------|-----------|
| **Language** | Go | 1.25.5 | Latest stable, excellent concurrency, cross-platform |
| **MCP Framework** | go-sdk | 1.1.0 | Official MCP implementation, stdio transport support |
| **OBS Client** | goobs | 1.5.6 | Active maintenance, protocol 5.5.6 support |
| **Database** | modernc.org/sqlite | latest | Pure Go (no CGO), easy cross-compilation |
| **Transport** | stdio | - | Standard MCP transport for AI assistants |

---

## Architecture

### Communication Flow

```
┌─────────────────────┐
│   MCP Client        │
│  (Claude/AI Agent)  │
└──────────┬──────────┘
           │ stdio (JSON-RPC)
           ↓
┌─────────────────────┐
│   MCP Server        │
│  (agentic-obs)      │
│                     │
│  ┌───────────────┐  │
│  │ MCP Tools     │  │
│  └───────┬───────┘  │
│          │          │
│  ┌───────↓───────┐  │
│  │ OBS Client    │  │
│  └───────┬───────┘  │
│          │          │
│  ┌───────↓───────┐  │
│  │ SQLite DB     │  │
│  └───────────────┘  │
└──────────┬──────────┘
           │ WebSocket (port 4455)
           ↓
┌─────────────────────┐
│   OBS Studio        │
│  (obs-websocket)    │
└─────────────────────┘
```

### Project Structure

```
agentic-obs/
├── main.go                    # Entry point, server initialization
├── go.mod / go.sum           # Dependencies
├── CLAUDE.MD / README.md     # Documentation
│
├── config/
│   └── config.go             # App configuration, env vars
│
├── internal/
│   ├── mcp/
│   │   ├── server.go         # MCP server lifecycle
│   │   ├── tools.go          # Tool definitions & handlers
│   │   ├── resources.go      # Resource handlers (scenes as resources)
│   │   ├── interfaces.go     # OBSClient interface for testing
│   │   └── testutil/         # Mock implementations for testing
│   │
│   ├── obs/
│   │   ├── client.go         # WebSocket client wrapper
│   │   ├── commands.go       # OBS operation implementations
│   │   └── events.go         # Event handling and notification dispatch
│   │
│   ├── storage/
│   │   ├── db.go             # Database setup & migrations
│   │   ├── scenes.go         # Scene preset persistence
│   │   ├── screenshots.go    # Screenshot source & image persistence
│   │   └── state.go          # Config & state management
│   │
│   ├── http/
│   │   └── server.go         # HTTP server for screenshot serving
│   │
│   └── screenshot/
│       └── manager.go        # Background screenshot capture manager
│
└── scripts/
    └── setup.sh              # Dev environment helpers
```

---

## Key Design Decisions

### 1. Storage Strategy: Moderate (Expandable)

**Decision:** SQLite stores connection config and user-defined scene/source presets.

**Rationale:**
- Practical value without over-engineering
- Fast access to frequently-used configurations
- Foundation for future expansion (analytics, state caching)

**Future:** Expand to comprehensive state caching, event history, performance analytics.

---

### 2. SQLite Driver: Pure Go (modernc.org/sqlite)

**Decision:** Use `modernc.org/sqlite` instead of `mattn/go-sqlite3`.

**Rationale:**
- No CGO requirement → easier cross-compilation
- Simpler build process for end users
- Performance difference negligible for config/preset storage
- Single-user, single-tenant use case

**Trade-off:** ~2x slower on writes vs CGO version, but acceptable for this use case.

---

### 3. OBS Connection: Persistent 1:1 (Multi-instance Future)

**Decision:** Maintain single persistent connection to one OBS instance with auto-reconnect.

**Rationale:**
- Simplifies Phase 1 implementation
- Covers 95% of use cases
- Reduces connection overhead for frequent operations

**Future Enhancement:** Multi-instance support requires architectural refactor:
- Connection pooling
- Instance ID namespacing
- Per-instance tool routing

**TODO:** Decide on multi-instance architecture before Phase 2.

---

### 4. Error Handling: Contextual & Actionable

**Decision:** Enrich errors with diagnostic context and suggested fixes.

**Examples:**
- ❌ "Connection failed"
- ✅ "OBS connection failed. Is OBS Studio running with WebSocket server enabled on port 4455?"

**Rationale:**
- Better UX for AI agents and end users
- Reduces troubleshooting friction
- Enables self-service problem resolution

---

### 5. Scenes as MCP Resources (With Notifications)

**Decision:** Expose OBS scenes as MCP resources, enabling server-initiated notifications for scene changes.

**Resource Design:**
- **URI Pattern:** `obs://scene_{scene_name}`
- **Resource Type:** `obs-scene`
- **Content:** JSON representation of scene configuration
- **Operations:** `resources/list`, `resources/read`

**Notification Events:**
- `notifications/resources/list_changed` - Scene created/deleted
- `notifications/resources/updated` - Scene modified or activated

**Rationale:**
- MCP resources are designed for dynamic, server-owned data
- Notifications keep AI clients synchronized with OBS state
- Aligns with MCP protocol design patterns
- Enables reactive workflows (e.g., "When I switch to Gaming scene, start recording")

**Implementation:**
- Monitor OBS events: `SceneCreated`, `SceneRemoved`, `CurrentProgramSceneChanged`
- Dispatch MCP notifications to connected clients
- Maintain event subscription in persistent OBS connection

**Future:** Expand to additional resource types (sources, filters, transitions, audio inputs).

---

### 6. MCP Tools vs Resources

**Decision:** Scenes exposed as resources (not tools), with complementary tools for actions.

**Resources (Read-only state):**
- `resources/list` - List all scenes
- `resources/read` - Get scene details

**Tools (Actions/Commands):**
- `set_current_scene` - Switch active scene
- `create_scene` - Create new scene
- `remove_scene` - Delete scene
- Recording, streaming, source, audio tools

**Rationale:**
- Resources represent queryable state with notifications
- Tools represent actions that modify state
- Clear separation of concerns
- Enables MCP clients to "watch" resources and "invoke" tools

---

### 7. Authentication: SQLite Storage (Local Security)

**Decision:** Store OBS WebSocket password in SQLite on first successful connection.

**Security Model:**
- Local, single-user deployment
- Database file permissions protect credentials
- No network exposure of credentials

**Options Considered:**
- Environment variables (good for automation, but not persistent)
- Interactive prompt every run (poor UX)
- Encrypted storage (over-engineering for local use)

**Accepted Risk:** Unencrypted local storage acceptable for single-user, localhost use case.

---

### 8. Setup Experience: Auto-Detect with Fallback

**Decision:** Try default `localhost:4455`, prompt on failure, persist successful config.

**Flow:**
1. Check SQLite for saved config
2. If not found, try `localhost:4455` with no password
3. If fails, interactive prompt for host/port/password
4. On success, save to SQLite for future runs

**Future Enhancement:** Interactive installation UX
- **TUI:** Terminal-based configuration wizard
- **Web:** Browser-based setup with OBS discovery

**TODO:** Design and implement interactive setup (TUI + Web).

---

## Phase 1: Baseline Implementation

### MCP Resources

| Resource Type | URI Pattern | Operations | Priority |
|---------------|-------------|------------|----------|
| **Scenes** | `obs://scene_{name}` | `list`, `read` | P0 |

**Notifications:**
- `notifications/resources/list_changed` - Scene added/removed
- `notifications/resources/updated` - Scene modified or activated

### Core MCP Tools

| Category | Tools | Priority |
|----------|-------|----------|
| **Scene Actions** | `set_current_scene`, `create_scene`, `remove_scene` | P0 |
| **Recording** | `start_recording`, `stop_recording`, `get_recording_status` | P0 |
| **Streaming** | `start_streaming`, `stop_streaming`, `get_streaming_status` | P0 |
| **Sources** | `list_sources`, `toggle_source_visibility`, `get_source_settings` | P1 |
| **Audio** | `get_input_mute`, `toggle_input_mute`, `set_input_volume` | P1 |
| **Status** | `get_obs_status` (connection & operational state) | P0 |

### Database Schema (Initial)

```sql
-- Connection configuration
CREATE TABLE config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Scene presets (user-defined)
CREATE TABLE scene_presets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    scene_name TEXT NOT NULL,
    sources TEXT, -- JSON array of source visibility states
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Application state
CREATE TABLE state (
    key TEXT PRIMARY KEY,
    value TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Known Unknowns & Future Research

### 1. Multi-Instance Architecture

**Question:** How to support multiple OBS instances through a single MCP server?

**Options:**
- A. Separate MCP servers per OBS instance (simple, but harder to coordinate)
- B. Single server with instance ID namespacing (e.g., `set_scene_instance1`)
- C. Connection pool with tool parameter for target instance

**Research Needed:**
- MCP tool naming conventions for multi-instance scenarios
- Connection pool management patterns in Go
- User experience for instance selection in AI conversations

**Decision Point:** Before Phase 2 implementation.

---

### 2. Additional Resource Types

**Question:** What other OBS entities should be exposed as MCP resources?

**Candidates:**
- **Sources:** Individual scene sources (cameras, displays, media)
- **Filters:** Audio/video filters applied to sources
- **Transitions:** Scene transition effects
- **Audio Inputs:** Microphones, desktop audio, etc.

**Resource URI Patterns:**
- `obs://source_{scene_name}_{source_name}`
- `obs://filter_{source_name}_{filter_name}`
- `obs://transition_{transition_name}`
- `obs://audio_{input_name}`

**Research Needed:**
- Which resources provide most value for AI interactions
- Notification frequency and performance impact
- Resource hierarchy and relationships
- Client subscription preferences

**Decision Point:** After Phase 1 scene resources are validated.

---

### 3. Interactive Setup UX

**Question:** What's the best first-run experience for non-technical users?

**Options:**
- **TUI (Terminal):** Go libraries like `bubbletea`, `tview`
- **Web:** Local HTTP server with browser UI
- **Hybrid:** TUI by default, web option with `--setup-web` flag

**User Personas:**
- Technical users: Prefer CLI/TUI
- Content creators: Prefer GUI/web
- Automated deployments: Prefer env vars/config files

**Research Needed:**
- Go TUI framework comparison
- Auto-browser-launch patterns
- OBS instance auto-discovery protocols

**Decision Point:** Phase 2 planning.

---

## Implementation Roadmap

### Phase 1: Foundation ✅ COMPLETE

**Status:** ✅ Completed 2025-12-14

**Deliverables:**
- [x] Go 1.25.5 installed
- [x] Project structure created
- [x] Dependencies installed
- [x] Documentation written (CLAUDE.MD, README.md, PROJECT_PLAN.md)
- [x] Basic MCP server implementation
- [x] Scene resources with list/read operations
- [x] OBS event monitoring (SceneCreated, SceneRemoved, CurrentProgramSceneChanged)
- [x] Resource notification dispatch
- [x] SQLite storage layer
- [x] OBS client wrapper with event handling
- [x] Core tools implemented (19 tools)
- [x] Auto-detection setup flow
- [x] Error handling with context

**Tools Implemented (19):**
- Scene Management: `list_scenes`, `set_current_scene`, `create_scene`, `remove_scene`
- Recording: `start_recording`, `stop_recording`, `get_recording_status`, `pause_recording`, `resume_recording`
- Streaming: `start_streaming`, `stop_streaming`, `get_streaming_status`
- Sources: `list_sources`, `toggle_source_visibility`, `get_source_settings`
- Audio: `get_input_mute`, `toggle_input_mute`, `set_input_volume`
- Status: `get_obs_status`

---

### Phase 2: Scene Presets & Testing ✅ COMPLETE

**Status:** ✅ Completed 2025-12-15

**Deliverables:**
- [x] Scene preset management (save/restore source visibility states)
- [x] 7 new tools: scene presets (6) + `get_input_volume`
- [x] OBSClient interface for dependency injection
- [x] Mock OBS client for testing
- [x] Comprehensive test coverage for storage layer
- [x] MCP tool handler tests

**Tools Added (7):**
- Scene Presets: `save_scene_preset`, `list_scene_presets`, `get_preset_details`, `apply_scene_preset`, `rename_scene_preset`, `delete_scene_preset`
- Audio: `get_input_volume`

**Total Tools:** 26

---

### Phase 3: Agentic Screenshot Sources ✅ COMPLETE

**Status:** ✅ Completed 2025-12-15

**Deliverables:**
- [x] Screenshot source management (periodic capture from OBS sources)
- [x] 4 new tools for screenshot control
- [x] HTTP server for serving screenshots at `http://localhost:8765/screenshot/{name}`
- [x] Background capture manager with configurable cadence
- [x] SQLite storage for screenshot sources and images
- [x] Automatic cleanup (keeps latest 10 screenshots per source)
- [x] Security hardening (path traversal prevention, generic error messages)

**Tools Added (4):**
- `create_screenshot_source` - Create periodic screenshot capture
- `remove_screenshot_source` - Stop and remove a screenshot source
- `list_screenshot_sources` - List all sources with status and URLs
- `configure_screenshot_cadence` - Update capture interval

**Total Tools:** 30

---

### Phase 4: MCP Resources & Prompts ✅ COMPLETE

**Deliverables:**
- [x] Screenshot resource (`obs://screenshot/{name}`) - Binary image blob
- [x] Preset resource (`obs://preset/{name}`) - JSON configuration
- [x] 10 MCP prompts for workflow automation
- [x] Prompt argument handling (required/optional)
- [x] Comprehensive test coverage (57 new tests)
- [x] Documentation updates (CLAUDE.md, README.md, docs/TOOLS.md)

**MCP Prompts Implemented:**
| Prompt | Arguments | Purpose |
|--------|-----------|---------|
| `stream-launch` | - | Pre-stream checklist |
| `stream-teardown` | - | End-stream cleanup |
| `audio-check` | - | Audio verification |
| `visual-check` | screenshot_source | Visual layout analysis |
| `health-check` | - | OBS diagnostic |
| `problem-detection` | screenshot_source | Issue detection |
| `preset-switcher` | preset_name (opt) | Preset management |
| `recording-workflow` | - | Recording session |
| `scene-organizer` | - | Scene organization |
| `quick-status` | - | Brief status |

**Success Criteria:** ✅ All met
- 3 MCP resource types (scenes, screenshots, presets)
- 10 MCP prompts with argument handling
- All tests passing

**Total Resources:** 4 (scenes, screenshots, screenshot-url, presets)
**Total Prompts:** 10

---

### Phase 5A: Setup & Configuration Enhancements ✅ COMPLETE

**Deliverables:**
- [x] Tool group preferences - Feature-based groupings (Core, Visual, Layout, Audio, Sources)
- [x] First-run setup prompts for tool groups and webserver
- [x] Optional webserver configuration (HTTP server can be disabled)
- [x] Screenshot-URL resource (`obs://screenshot-url/{name}`) - lightweight JSON alternative
- [x] Conditional tool registration based on user preferences
- [x] Persistent tool group and webserver configuration in SQLite

**New Resource:** `obs://screenshot-url/{sourceName}` - Returns JSON with HTTP URL instead of binary data

**Success Criteria:**
- [x] Users can select tool groups during first-run setup
- [x] Tool registration respects group preferences
- [x] HTTP server can be disabled via config
- [x] Screenshot-URL resource returns JSON with URL
- [x] All preferences persist across restarts

---

### Phase 6.1: Web Frontend ✅ COMPLETE

**Deliverables:**
- [x] Web dashboard at `http://localhost:8765/`
- [x] API endpoints: `/api/status`, `/api/history`, `/api/history/stats`, `/api/screenshots`, `/api/config`
- [x] Real-time status display with auto-refresh
- [x] Screenshot sources gallery with live preview
- [x] Action history viewer with filtering
- [x] Configuration management via web interface
- [x] Dark-themed responsive UI (vanilla JS, no build step)
- [x] Action history database table with statistics

**New Files:**
- `internal/http/handlers.go` - API handlers
- `internal/http/static/index.html` - Dashboard UI
- `internal/storage/history.go` - Action history storage

**API Endpoints:**
| Route | Method | Description |
|-------|--------|-------------|
| `/` | GET | Web dashboard |
| `/api/status` | GET | Server status JSON |
| `/api/history` | GET | Action history (supports `?limit=N`, `?tool=name`) |
| `/api/history/stats` | GET | Action statistics |
| `/api/screenshots` | GET | Screenshot sources with URLs |
| `/api/config` | GET/POST | Read/update configuration |

---

### Phase 6.2-6.4: Future Enhancements (Planned)

**Phase 6.2 - TUI App:**
- [ ] Terminal interface using Bubbletea
- [ ] Mirror web frontend functionality

**Phase 6.3 - Scene Design:**
- [ ] Source creation tools (text, image, color, browser, media)
- [ ] Layout control tools (transform, bounds, crop, order)
- [ ] Advanced tools (duplicate, lock, list input kinds)

**Phase 6.4 - Advanced:**
- [ ] Automation rules and macros (event-triggered actions)
- [ ] Multi-instance OBS support
- [ ] Additional resource types (sources, filters, transitions, audio inputs)
- [ ] Resource subscriptions (explicit client opt-in)
- [ ] Action logging hooks in MCP tool handlers

**Success Criteria:**
- TUI offers command-line management
- AI can programmatically design OBS scenes

---

## Development Workflow

### Adding a New MCP Tool

1. **Define Schema** (`internal/mcp/tools.go`)
   ```go
   &mcp.Tool{
       Name: "my_tool",
       Description: "Does something useful",
       InputSchema: /* JSON schema */
   }
   ```

2. **Implement Handler** (`internal/mcp/tools.go`)
   ```go
   func handleMyTool(ctx context.Context, req *mcp.Request) (*mcp.Response, error) {
       // Implementation
   }
   ```

3. **Add OBS Command** (`internal/obs/commands.go`)
   ```go
   func (c *Client) MyCommand(params) (result, error) {
       // goobs call
   }
   ```

4. **Register** (`internal/mcp/server.go`)
   ```go
   mcp.AddTool(server, toolDef, handleMyTool)
   ```

### Adding a New MCP Resource

1. **Define Resource Schema** (`internal/mcp/resources.go`)
   ```go
   // Resource URI: obs://scene_{scene_name}
   type SceneResource struct {
       URI         string `json:"uri"`
       Name        string `json:"name"`
       MimeType    string `json:"mimeType"`
       Description string `json:"description"`
   }
   ```

2. **Implement Handlers** (`internal/mcp/resources.go`)
   ```go
   func handleResourcesList(ctx context.Context) ([]*Resource, error) {
       // Return list of all scene resources
   }

   func handleResourceRead(ctx context.Context, uri string) (*ResourceContent, error) {
       // Return specific scene configuration
   }
   ```

3. **Monitor OBS Events** (`internal/obs/events.go`)
   ```go
   func (c *Client) handleSceneCreated(event *events.SceneCreated) {
       c.notifyResourceListChanged()
   }

   func (c *Client) handleCurrentProgramSceneChanged(event *events.CurrentProgramSceneChanged) {
       c.notifyResourceUpdated(fmt.Sprintf("obs://scene_%s", event.SceneName))
   }
   ```

4. **Send Notifications** (`internal/mcp/server.go`)
   ```go
   server.SendResourceListChanged()
   server.SendResourceUpdated(uri)
   ```

### Testing Strategy

- **Unit Tests:** Storage layer, utility functions
- **Integration Tests:** OBS client (requires running OBS)
- **E2E Tests:** Full MCP server with mock client
- **Resource Tests:** Verify notifications sent on state changes

---

## Success Metrics

### Phase 1 ✅
- [x] Project structure created
- [x] Dependencies configured
- [x] Server responds to MCP tool calls and resource requests
- [x] Successfully connects to OBS
- [x] Scenes exposed as resources
- [x] Notifications working for scene changes
- [x] Core tools (19) functional

### Phase 2 ✅
- [x] Scene preset management (6 tools)
- [x] Testing infrastructure (interfaces, mocks)
- [x] Storage layer test coverage

### Phase 3 ✅
- [x] Screenshot source management (4 tools)
- [x] HTTP server for image serving
- [x] Background capture manager
- [x] Security review passed

### Phase 4 ✅
- [x] Screenshot resource (binary blob via MCP)
- [x] Preset resource (JSON via MCP)
- [x] 10 MCP prompts for workflows
- [x] Prompt argument validation
- [x] 57 new tests (resources + prompts)

### Long-term
- [ ] Multi-instance support
- [ ] Production deployment examples
- [ ] Community contributions
- [ ] Documentation completeness

---

## Resources & References

**MCP Protocol:**
- [MCP Specification](https://modelcontextprotocol.io)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [MCP Core Architecture](https://modelcontextprotocol.io/docs/concepts/architecture)

**OBS Integration:**
- [goobs Documentation](https://pkg.go.dev/github.com/andreykaipov/goobs)
- [OBS WebSocket Protocol](https://github.com/obsproject/obs-websocket/blob/master/docs/generated/protocol.md)
- [OBS Studio](https://obsproject.com)

**Go & SQLite:**
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)
- [Go SQLite Benchmarks](https://github.com/cvilsmeier/go-sqlite-bench)

**Discussions:**
- [MCP Notifications](https://github.com/modelcontextprotocol/modelcontextprotocol/discussions/1192)
- [WebSocket Transport Proposal](https://github.com/modelcontextprotocol/modelcontextprotocol/issues/1288)

---

## Appendix: Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2025-12-14 | Use Go 1.25.5 | Latest stable, best performance |
| 2025-12-14 | modernc.org/sqlite over mattn | Pure Go, easier cross-compilation |
| 2025-12-14 | Persistent 1:1 OBS connection | Simplicity for Phase 1, defer multi-instance |
| 2025-12-14 | Contextual error messages | Better UX for AI agents and users |
| 2025-12-14 | **Scenes as MCP Resources** | Enable notifications for scene changes, aligns with MCP design |
| 2025-12-14 | Tools vs Resources separation | Resources for state, tools for actions |
| 2025-12-14 | Resource notifications in Phase 1 | Critical for keeping clients synchronized with OBS |
| 2025-12-14 | SQLite password storage | Acceptable for local, single-user deployment |
| 2025-12-14 | Auto-detect setup with fallback | Best balance of convenience and flexibility |
| 2025-12-15 | OBSClient interface for testing | Enables mock injection without running OBS |
| 2025-12-15 | Scene presets store visibility states | Minimal but useful; covers common "show/hide sources" workflow |
| 2025-12-15 | HTTP server for screenshots | Browser sources can fetch images; simpler than MCP binary resources |
| 2025-12-15 | Background capture workers | Per-source goroutines with configurable cadence |
| 2025-12-15 | Keep 10 screenshots per source | Balance between history and storage; ~13-26MB per source |
| 2025-12-15 | Path traversal validation | Security hardening for HTTP endpoint source names |

---

**Document Version:** 1.6
**Last Updated:** 2025-12-18
**Next Review:** Before Phase 6.2 planning

**Changelog:**
- v1.6 (2025-12-18): Phase 6.1 complete - Web dashboard, API endpoints, action history, config management
- v1.5 (2025-12-15): Phase 5A complete - Tool group preferences, optional webserver, screenshot-url resource, first-run setup
- v1.4 (2025-12-15): Phase 4 complete - MCP resources (screenshots, presets), 10 MCP prompts, comprehensive tests
- v1.3 (2025-12-15): Phase 3 complete - agentic screenshot sources, HTTP server, background capture manager
- v1.2 (2025-12-15): Phase 2 complete - scene presets, testing infrastructure, 26 tools
- v1.1 (2025-12-14): Added scenes-as-resources architecture, MCP notifications, resource implementation patterns
- v1.0 (2025-12-14): Initial project plan with baseline architecture
