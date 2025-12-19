# Roadmap

Future enhancements and research topics for agentic-obs.

> **Note:** For completed phases, see [CHANGELOG.md](../CHANGELOG.md).

## Planned Enhancements

### Automation Rules & Macros

**Priority:** High
**Complexity:** Medium

Event-triggered actions and multi-step sequences:

- **Trigger types:**
  - OBS events (scene change, recording start/stop)
  - Time-based (scheduled actions)
  - External webhooks

- **Action sequences:**
  - Multi-tool workflows (e.g., "switch to BRB scene, mute mic, start timer")
  - Conditional logic
  - Variable substitution

- **Use cases:**
  - Auto-mute on scene switch
  - Scheduled scene rotations
  - Stream start/stop routines

### Multi-Instance OBS Support

**Priority:** Medium
**Complexity:** High
**See:** [ADR-002: OBS Connection](decisions/002-obs-connection.md)

Support controlling multiple OBS instances from a single MCP server:

- **Architecture Options:**
  - Separate MCP servers per OBS instance (simple, harder to coordinate)
  - Single server with instance ID namespacing (e.g., `obs://instance1/scene/Gaming`)
  - Connection pool with tool parameter for target instance

- **Research Needed:**
  - MCP tool naming conventions for multi-instance scenarios
  - Connection pool management patterns in Go
  - User experience for instance selection in AI conversations

- **Resource URI Changes:**
  - `obs://scene/{name}` → `obs://{instance}/scene/{name}`

### Additional Resource Types

**Priority:** Medium
**Complexity:** Low-Medium

Expand beyond scenes to expose more OBS entities:

| Resource Type | URI Pattern | Notifications |
|---------------|-------------|---------------|
| **Sources** | `obs://source/{scene}/{name}` | Visibility, settings changes |
| **Filters** | `obs://filter/{source}/{name}` | Filter enable/disable, settings |
| **Transitions** | `obs://transition/{name}` | Active transition changes |
| **Audio Inputs** | `obs://audio/{name}` | Volume, mute changes |

**Benefits:**
- Richer AI context about OBS state
- More granular notifications
- Better support for complex workflows

### Resource Subscriptions

**Priority:** Low
**Complexity:** Medium

Explicit client opt-in to resource notifications:

- Currently: Server broadcasts all notifications
- Proposed: Clients subscribe to specific resource patterns

```
// Client requests
resources/subscribe { uri: "obs://scene/*" }
resources/subscribe { uri: "obs://audio/Microphone" }
```

**Benefits:**
- Reduced notification noise
- Lower bandwidth usage
- Client-controlled update frequency

---

## Research Topics

### Interactive Setup UX

**Question:** Best first-run experience for diverse users?

**Options:**
| Approach | Pros | Cons |
|----------|------|------|
| **TUI (Terminal)** | Technical users love it; Go has good libs (bubbletea) | Not intuitive for content creators |
| **Web UI** | Visual, familiar to all users | Requires browser launch |
| **Hybrid** | Best of both; TUI default, `--setup-web` flag | More code to maintain |

**Research Needed:**
- Go TUI framework comparison
- Auto-browser-launch patterns
- OBS instance auto-discovery protocols

### Streaming Platform Integration

**Question:** Should agentic-obs integrate directly with streaming platforms?

**Potential Features:**
- Chat overlay management
- Stream alerts integration
- Viewer count monitoring
- Go live/offline detection

**Considerations:**
- Scope creep vs. focused OBS control
- OAuth complexity for each platform
- Alternative: Recommend companion tools

### AI Vision Integration

**Question:** Can we leverage multimodal AI for visual monitoring?

**Ideas:**
- Screenshot → AI analysis pipeline
- "Is my webcam visible?" detection
- Layout quality scoring
- Brand guideline compliance checking

**Dependencies:**
- Multimodal AI API availability
- Latency requirements
- Cost per analysis

---

## Backlog Items

### Documentation & Developer Experience

- [ ] **FB-15:** Extract mcpui-go as standalone library
- [ ] API documentation generation (godoc)
- [ ] Integration test suite with mock OBS
- [ ] Example Claude Code skills showcase

### Performance & Reliability

- [ ] Connection health metrics
- [ ] Screenshot capture performance profiling
- [ ] Memory usage optimization for long-running servers

### Security Hardening

- [ ] Optional password encryption (SQLCipher)
- [ ] API authentication for web UI
- [ ] Rate limiting on HTTP endpoints

### Platform Support

- [ ] Docker containerization
- [ ] Windows service installation guide
- [ ] macOS launchd integration

---

## How to Contribute

Ideas for new features? Open a GitHub issue with:

1. **Use case:** What problem does this solve?
2. **Proposed solution:** How should it work?
3. **Alternatives considered:** What else could solve this?

For significant changes, we'll create an ADR in [decisions/](decisions/).
