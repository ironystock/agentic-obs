# agentic-obs Project Plan

## Executive Summary

This document outlines the architecture, decisions, and implementation roadmap for the **agentic-obs** MCP server - a Go-based bridge between AI assistants and OBS Studio.

**Created:** 2025-12-14
**Updated:** 2025-12-14
**Status:** ✅ Phase 1 Complete - All 19 Tools Implemented

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
│   │   └── resources.go      # Resource handlers (scenes as resources)
│   │
│   ├── obs/
│   │   ├── client.go         # WebSocket client wrapper
│   │   ├── commands.go       # OBS operation implementations
│   │   └── events.go         # Event handling and notification dispatch
│   │
│   └── storage/
│       ├── db.go             # Database setup & migrations
│       ├── scenes.go         # Scene preset persistence
│       └── state.go          # Config & state management
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
- [x] Core tools implemented (P0: 10 tools)
- [x] **P1 tools implemented (9 additional tools)**
- [x] Auto-detection setup flow
- [x] Error handling with context
- [x] **Comprehensive documentation (QUICKSTART, TOOLS.md, examples/)**
- [x] **Natural language prompt examples**

**Success Criteria:** ✅ All Met
- ✅ MCP server responds to tool calls and resource requests via stdio
- ✅ Successfully connects to OBS Studio
- ✅ Scenes exposed as resources with `resources/list` and `resources/read`
- ✅ Notifications sent when scenes change
- ✅ Can switch scenes, start/stop recording via tools
- ✅ Config persists between runs
- ✅ **19 total tools operational (Scene, Recording, Streaming, Source, Audio, Status)**

---

### Phase 2: Enhancement (Partially Complete)

**Status:** P1 Tools Complete ✅, Other Items Future

**Deliverables:**
- [x] **Additional tools (P1: sources, audio) - COMPLETE**
- [ ] Scene preset management
- [ ] Interactive setup (TUI + Web)
- [ ] Comprehensive error handling
- [ ] Unit and integration tests
- [ ] Performance optimization

**Success Criteria:**
- ✅ All P1 tools implemented (19 total)
- ✅ First-run experience is smooth (auto-detection working)
- [ ] Test coverage >70%

---

### Phase 3: Advanced Features (Future)

**Deliverables:**
- [ ] Multi-instance OBS support
- [ ] Additional resource types (sources, filters, transitions, audio)
- [ ] Resource subscriptions (explicit client opt-in)
- [ ] State caching and analytics
- [ ] Advanced preset/macro system
- [ ] Health monitoring
- [ ] Performance metrics

**Success Criteria:**
- Supports multiple OBS instances
- Multiple resource types with notifications
- Production-ready stability

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

### Phase 1
- [x] Project structure created
- [x] Dependencies configured
- [ ] Server responds to MCP tool calls and resource requests
- [ ] Successfully connects to OBS
- [ ] Scenes exposed as resources
- [ ] Notifications working for scene changes
- [ ] Core tools (P0) functional

### Long-term
- Multi-instance support
- Production deployment examples
- Community contributions
- Documentation completeness

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

---

**Document Version:** 1.1
**Last Updated:** 2025-12-14
**Next Review:** After Phase 1 completion

**Changelog:**
- v1.1 (2025-12-14): Added scenes-as-resources architecture, MCP notifications, resource implementation patterns
- v1.0 (2025-12-14): Initial project plan with baseline architecture
