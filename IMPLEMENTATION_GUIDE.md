# Implementation Guide - Using Sub-Agents

This guide explains how to efficiently implement the agentic-obs project using Claude Code's sub-agent system for parallel development.

## Overview

The project is divided into independent modules that can be implemented in parallel using specialized sub-agents. This approach maximizes development speed and maintains code quality.

## Implementation Strategy

### Phase 1: Core Infrastructure (Parallel)

These components can be built simultaneously by different sub-agents:

#### Agent 1: Storage Layer (`internal/storage/`)
**Files:** `db.go`, `state.go`, `scenes.go`

**Prompt for Agent:**
```
Implement the SQLite storage layer for agentic-obs:

1. Create db.go with:
   - Database initialization
   - Schema migrations for config, state, and scene_presets tables
   - Connection pooling (modernc.org/sqlite)
   - Error handling with contextual messages

2. Create state.go with:
   - Save/load configuration (OBS host, port, password)
   - Key-value state management
   - First-run detection

3. Create scenes.go with:
   - Scene preset CRUD operations
   - JSON serialization for preset data

Use modernc.org/sqlite (pure Go, no CGO).
Schema is defined in PROJECT_PLAN.md section "Database Schema (Initial)".
Follow Go best practices, add error wrapping, use context.Context.

Return: Complete implementation of all three files with inline comments.
```

---

#### Agent 2: OBS Client Wrapper (`internal/obs/`)
**Files:** `client.go`, `commands.go`, `events.go`

**Prompt for Agent:**
```
Implement the OBS WebSocket client wrapper for agentic-obs using github.com/andreykaipov/goobs v1.5.6:

1. Create client.go with:
   - Client struct with connection state
   - Connect/Disconnect with auto-reconnect logic
   - Health check / connection status
   - Constructor that accepts host, port, password

2. Create commands.go with:
   - GetSceneList() - returns all scenes
   - GetSceneByName(name) - returns scene details with sources
   - SetCurrentScene(name) - switches active scene
   - StartRecording(), StopRecording(), GetRecordingStatus()
   - StartStreaming(), StopStreaming(), GetStreamingStatus()

3. Create events.go with:
   - Event subscription setup
   - Handlers for: SceneCreated, SceneRemoved, CurrentProgramSceneChanged
   - Callback interface for notification dispatch to MCP server

Connection details: localhost:4455 (default)
Event handlers should log events and trigger MCP resource notifications.
Use contextual error messages (e.g., "OBS connection failed. Is OBS Studio running?").

Return: Complete implementation with error handling and inline docs.
```

---

#### Agent 3: MCP Server Foundation (`internal/mcp/`)
**Files:** `server.go`, `tools.go`, `resources.go`

**Prompt for Agent:**
```
Implement the MCP server for agentic-obs using github.com/modelcontextprotocol/go-sdk v1.1.0:

1. Create server.go with:
   - Server initialization with stdio transport
   - Lifecycle management (start, stop, graceful shutdown)
   - Integration points for OBS client and storage
   - Resource notification dispatch (SendResourceListChanged, SendResourceUpdated)

2. Create resources.go with:
   - handleResourcesList - returns all scenes as resources (URI: obs://scene_{name})
   - handleResourceRead - returns scene details as JSON
   - Scene resource struct with uri, name, mimeType, description fields

3. Create tools.go with P0 tools:
   - set_current_scene (input: scene_name)
   - create_scene (input: scene_name)
   - remove_scene (input: scene_name)
   - start_recording, stop_recording, get_recording_status
   - start_streaming, stop_streaming, get_streaming_status
   - get_obs_status

Tool handlers should call OBS client commands and return structured responses.
Resources should use URI pattern: obs://scene_{scene_name}
All tools need JSON schema for input validation.

Return: Complete implementation with MCP protocol compliance.
```

---

### Phase 2: Integration & Wiring

After the three parallel agents complete, use a single integration agent:

#### Agent 4: Main Entry Point & Integration (`main.go`, `config/`)
**Files:** `main.go`, `config/config.go`

**Prompt for Agent:**
```
Integrate the agentic-obs components and create the main entry point:

1. Create config/config.go with:
   - Configuration struct (OBS host, port, password, DB path)
   - Load from storage or environment variables
   - Auto-detection logic (try localhost:4455, prompt on failure)
   - Save successful config to SQLite

2. Create main.go with:
   - Initialize storage layer
   - Load or detect OBS configuration
   - Initialize OBS client with auto-reconnect
   - Initialize MCP server with stdio transport
   - Wire OBS events to MCP resource notifications
   - Setup graceful shutdown (SIGINT, SIGTERM)
   - Error handling with contextual messages

Integration flow:
1. Init storage → 2. Load/detect config → 3. Connect to OBS → 4. Start MCP server → 5. Run event loop

The three internal packages (storage, obs, mcp) are already implemented.
Wire them together so:
- OBS event handlers trigger MCP notifications
- MCP tool handlers call OBS commands
- Config persists in storage

Return: Complete main.go and config.go with startup logging.
```

---

## Parallel Execution Commands

### Option 1: Sequential (One at a time)

```bash
# Start with storage layer
Task: "Implement storage layer as described in IMPLEMENTATION_GUIDE.md Agent 1"

# Then OBS client
Task: "Implement OBS client wrapper as described in IMPLEMENTATION_GUIDE.md Agent 2"

# Then MCP server
Task: "Implement MCP server as described in IMPLEMENTATION_GUIDE.md Agent 3"

# Finally integration
Task: "Integrate components as described in IMPLEMENTATION_GUIDE.md Agent 4"
```

### Option 2: Parallel (Faster - use 3 agents simultaneously)

Tell Claude Code:
```
Launch 3 agents in parallel to implement:
1. Storage layer (internal/storage/)
2. OBS client wrapper (internal/obs/)
3. MCP server foundation (internal/mcp/)

Use the prompts from IMPLEMENTATION_GUIDE.md for each agent.
```

Claude Code will spawn multiple Task tools in a single message for parallel execution.

---

## Verification After Each Phase

### After Storage Agent:
```bash
go build ./internal/storage/
# Should compile without errors
```

### After OBS Agent:
```bash
go build ./internal/obs/
# Should compile without errors
```

### After MCP Agent:
```bash
go build ./internal/mcp/
# Should compile without errors
```

### After Integration:
```bash
go build -o agentic-obs.exe main.go
# Should produce executable
```

### Final Test:
```bash
# Start OBS Studio with WebSocket enabled
./agentic-obs.exe
# Should connect to OBS and wait for MCP client connection via stdio
```

---

## File Checklist

After all agents complete, you should have:

```
agentic-obs/
├── main.go                      ✓ (Agent 4)
├── config/
│   └── config.go               ✓ (Agent 4)
├── internal/
│   ├── storage/
│   │   ├── db.go               ✓ (Agent 1)
│   │   ├── state.go            ✓ (Agent 1)
│   │   └── scenes.go           ✓ (Agent 1)
│   ├── obs/
│   │   ├── client.go           ✓ (Agent 2)
│   │   ├── commands.go         ✓ (Agent 2)
│   │   └── events.go           ✓ (Agent 2)
│   └── mcp/
│       ├── server.go           ✓ (Agent 3)
│       ├── resources.go        ✓ (Agent 3)
│       └── tools.go            ✓ (Agent 3)
└── ... (docs, go.mod, etc. already exist)
```

---

## Debugging & Iteration

If any agent's implementation has issues:

1. **Review Agent Output:** Check for compilation errors
2. **Incremental Fix:** Ask Claude Code to fix specific issues
3. **Integration Issues:** Use Agent 4 to wire components correctly
4. **Test Individually:** Build each package separately to isolate problems

---

## Next Steps After Phase 1

Once Phase 1 is complete and tested:

### Phase 1.5: P1 Tools (Optional)
Add lower-priority tools:
- Source management tools (list_sources, toggle_source_visibility)
- Audio tools (get_input_mute, toggle_input_mute, set_input_volume)

### Phase 2: Enhancement
- Scene preset management
- Interactive setup (TUI)
- Comprehensive testing

### Phase 3: Advanced Features
- Additional resource types (sources, audio, filters)
- Multi-instance OBS support
- Agentic screenshot sources

---

## Tips for Working with Sub-Agents

1. **Be Specific:** Provide complete context in each agent prompt
2. **Reference Docs:** Point agents to PROJECT_PLAN.md, RESOURCES.md, CLAUDE.MD
3. **Independent Work:** Ensure agents don't need to coordinate (clean interfaces)
4. **Review Output:** Validate each agent's code before integrating
5. **Iterate:** If output isn't perfect, refine with follow-up prompts

---

## Example: Launching Parallel Agents

**Your Prompt to Claude Code:**
```
I want to implement agentic-obs in parallel. Launch 3 agents simultaneously using the prompts in IMPLEMENTATION_GUIDE.md:

Agent 1: Storage layer (internal/storage/)
Agent 2: OBS client wrapper (internal/obs/)
Agent 3: MCP server foundation (internal/mcp/)

Run these in parallel.
```

Claude Code will create 3 Task tool calls in a single message, executing them concurrently.

---

**Document Version:** 1.0
**Last Updated:** 2025-12-14
**Related:** PROJECT_PLAN.md, CLAUDE.MD, RESOURCES.md
