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

## Feature Backlog (FB Items)

Tracked features with unique identifiers for reference.

### Completed

| ID | Name | Description | Completed |
|----|------|-------------|-----------|
| FB-2 | Help Tool | MCP help tool with 6 topics + per-tool help | Phase 7 |
| FB-9 | Claude Skills | 4 skill packages (streaming-assistant, scene-designer, audio-engineer, preset-manager) | Phase 7 |
| FB-10 | Additional Prompts | 3 new prompts (scene-designer, source-management, visual-setup) | Phase 7 |
| FB-11 | MCP Completions | Autocomplete for prompts/resources with 5s TTL caching | Phase 7 |
| FB-12 | MCP-UI Go SDK | Protocol implementation (SEP-1865) | Phase 8 |
| FB-13 | MCP-UI Integration | Rich UI resources for agentic-obs (4 phases) | Phase 8 |
| FB-15 | SDK Extraction | mcpui-go as standalone package v0.1.0 | Post-Phase 8 |

### Active Backlog

| ID | Name | Priority | Complexity | Dependencies | Description |
|----|------|----------|------------|--------------|-------------|
| FB-17 | Config Sync | High | Medium | - | Env vars, version management, validation |
| FB-16 | Skills Completion | Medium | Low | - | Missing SKILL.md files for audio-engineer, preset-manager |
| FB-18 | Build System | Medium | Medium | FB-17 | Makefile, goreleaser, version injection |
| FB-14 | Brand & Design | Medium | Medium | UX-SPEC | Visual identity implementation (blocked by UX-SPEC.md) |
| FB-1 | Embedded Docs | Medium | Low-Med | - | Docs in HTTP/TUI via go:embed |
| FB-3 | Elucidation | High | High | FB-2 ✅ | Intent disambiguation framework |
| FB-4 | SDK Migration | Medium | Medium | - | Process for tracking go-sdk updates |
| FB-5 | Static Website | Low | Medium | FB-14 | Project documentation site |
| FB-6 | Network API | Low | High | - | Non-localhost HTTP exposure |
| FB-7 | Multi-Instance | Medium | Very High | - | Multiple OBS support |
| FB-8 | Remote Hosted | Low | Very High | FB-6, FB-7 | Cloud-hosted server |

### Other Backlog Items

**Documentation & Developer Experience:**
- [ ] API documentation generation (godoc)
- [ ] Integration test suite with mock OBS
- [ ] Example Claude Code skills showcase

**Performance & Reliability:**
- [ ] Connection health metrics
- [ ] Screenshot capture performance profiling
- [ ] Memory usage optimization for long-running servers

**Security Hardening:**
- [ ] Optional password encryption (SQLCipher)
- [ ] API authentication for web UI
- [ ] Rate limiting on HTTP endpoints

**Platform Support:**
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
