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

## Sprint Planning

Work is organized into sprints using a semver-inspired cadence:

- **N.0 — Feature sprint** — introduces new capability (tools, resources, prompts, skills).
- **N.5 — Chore / alignment sprint** — dep bumps, deprecations, debt, follow-ups from the preceding feature sprint.

Sprint names follow the format **"Sprint N.N — <Focus>"**, where the focus describes what a user or reviewer would call the work (e.g. "Canvas & Multi-Output"), not the internal category.

FB item numbers and sprint numbers are orthogonal: an **FB item** is a single tracked piece of work; a **sprint** is a timebox grouping of FB items. An FB belongs to exactly one sprint once scheduled — the `Sprint` column in the Active Backlog captures that mapping.

### Active Sprints

| Sprint | Focus | Status | Contains |
|--------|-------|--------|----------|
| **0.5** | Upstream Alignment & FB-20 Follow-ups | In progress | FB-4, FB-32..FB-41, FB-43, FB-44 (FB-30, FB-31 folded into FB-33) |
| **1.0** | Canvas & Multi-Output | Candidate (blocked on 0.5) | FB-42 |

### Sprint 0.5 — Upstream Alignment & FB-20 Follow-ups

**Theme:** realign with what upstream has shipped while we were away (MCP Go SDK, obs-websocket 5.5 → 5.7, MCP Apps spec, Skills vs. current tool inventory), and close the safety/test items deferred from FB-20 (PR #26).

**Order of operations (matters — each unblocks discoveries for the next):**
1. **FB-4** — MCP Go SDK bump. Surfaces what MCP Apps / sampling / annotations look like in the current SDK.
2. **FB-34** — goobs + obs-websocket protocol bump. Surfaces additional 5.6/5.7 requests and events beyond the ones we already know about.
3. **FB-32** — MCP Apps port, informed by FB-4's discoveries.
4. **FB-33** — Skills modernization sweep, informed by FB-20 ✅ + FB-27 ✅ tool inventory.
5. **FB-35, FB-36** — Deprecated-field audit and `RecordFileChanged` wiring (mechanical, follow FB-34).
6. **FB-37..FB-41** — FB-20 safety / test / retention follow-ups (can parallelize).
7. **FB-43, FB-44** — Charm v1 majors (TUI/docs) and `modernc/sqlite` catch-up (storage). Independent of MCP/OBS layers; can parallelize with FB-37..FB-41.

**Fold-in dep bumps** (ride along with FB-4 / FB-34 `go mod tidy`, no dedicated FB): `yuin/goldmark` v1.7→v1.8, `google/jsonschema-go` v0.3→v0.4, `golang.org/x/{net,crypto,oauth2}` security catch-up, `golang-jwt/jwt/v5` v5.2→v5.3. Explicitly `go get` to target after tidy so we control the floor rather than accepting `tidy`'s minimum-required.

**Tangent discipline:** during each dep bump, feature opportunities that surface go into [`sprints/0.5-tangent-log.md`](sprints/0.5-tangent-log.md) rather than expanding the sprint. A **30-minute gate** applies: if the enhancement is literally 30 minutes AND entirely additive, fold it into the dep bump's commit; anything bigger goes to the log. Candidates are promoted to FB numbers at sprint close.

### Sprint 1.0 — Canvas & Multi-Output *(candidate)*

OBS 30's multi-canvas feature (vertical streaming, multi-output from a single OBS instance) shipped in obs-websocket protocol 5.7.0. Tracked as FB-42. Hard-depends on Sprint 0.5's goobs bump (FB-34) — the generated Go types for `GetCanvasList`, `Canvas*` events, and the optional `canvasUuid` parameter don't exist in goobs v1.5.6.

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
| FB-16 | Skills Completion | SKILL.md files for all 4 skill packages | Post-Phase 8 |
| FB-17 | Config Sync | Env vars, version injection, ApplyEnvOverrides() | Post-Phase 8 |
| FB-18 | Build System | Makefile, goreleaser, version.go, BUILD.md | Post-Phase 8 |
| FB-14 | Brand & Design | Logo, favicon, ASCII banners, design system colors | Phase 9 |
| FB-1 | Embedded Docs | Docs in HTTP/TUI via go:embed, goldmark/glamour rendering | Phase 9 |
| FB-3 | MCP Elicitation | User confirmation for high-risk tools (streaming, delete) | Phase 9 |
| FB-23 | Filters Tool Group | 7 tools for source filter management | Phase 10 |
| FB-24 | Transitions | 5 tools for scene transition control | Phase 10 |
| FB-25 | Virtual Cam & Replay | 6 tools for virtual camera and replay buffer | Phase 11 |
| FB-26 | Studio Mode & Hotkeys | 6 tools for studio mode and hotkey control | Phase 11 |
| FB-27 | Dynamic Tool Config | 3 meta-tools for runtime tool group enable/disable | Phase 12 |
| FB-28 | Skills Update | streaming-assistant with virtual cam, replay, studio mode, hotkeys | Phase 12 |
| FB-20 | Automation Rules | 9 automation tools for event-triggered actions and scheduled tasks | Phase 13 |

### Active Backlog

Ordered by sprint, then priority. Items marked `—` in the Sprint column are unscheduled.

| ID | Name | Priority | Complexity | Sprint | Dependencies | Description |
|----|------|----------|------------|--------|--------------|-------------|
| FB-4 | SDK Migration (MCP Go SDK bump) | High | Medium | **0.5** | - | Bump `github.com/modelcontextprotocol/go-sdk` past v1.1.0; absorb API changes; surface MCP Apps / sampling / annotations opportunities |
| FB-34 | goobs + obs-websocket 5.5 → 5.7 bump | High | Medium | **0.5** | - | Bump goobs v1.5.6 → v1.8.3 (protocol 5.5 → 5.7); absorb API changes; unblocks FB-42 |
| FB-32 | MCP Apps Port | High | Medium-High | **0.5** | FB-4 | Port the mcpui-go work (FB-15 ✅) to the new MCP Apps spec at `apps.extensions.modelcontextprotocol.io`; shape depends on FB-4 discoveries |
| FB-33 | Skills Modernization Sweep | Medium | Medium | **0.5** | FB-20 ✅, FB-27 ✅ | Audit all 4 skills against 81 tools + 14 prompts; add automation coverage to `streaming-assistant`; folds FB-30 and FB-31 |
| FB-35 | Deprecated field audit (`currentProgram*` / `currentPreview*`) | Medium | Low | **0.5** | FB-34 | Migrate off scene-response fields flagged for removal in obs-websocket protocol |
| FB-36 | `RecordFileChanged` event wiring | Low | Low | **0.5** | FB-34 | Bridge OBS 30+ file-split event into the automation engine's event bridge |
| FB-37 | Automation cooldown race fix | Medium | Low-Med | **0.5** | FB-20 ✅ | Move `recordCooldown` from execute-end to dispatch-time; fixes intermittent `TestEngineCooldown` at `-count>=3` |
| FB-38 | Automation queue-overflow metric | Low | Low | **0.5** | FB-20 ✅ | Expose `dropped_events_total` counter for the engine's 100-deep `eventChan` |
| FB-39 | Automation concurrency tests | Medium | Medium | **0.5** | FB-20 ✅ | Stress tests for cooldown map + rule cache; may need a CGO-enabled test lane for `-race` |
| FB-40 | Automation `OnError="stop"` test | Low | Low | **0.5** | FB-20 ✅ | Unit test covering action-chain halt when `OnError=stop` |
| FB-41 | Automation execution retention | Medium | Low-Med | **0.5** | FB-20 ✅ | Scheduled sweep via `ClearOldRuleExecutions` to prevent DB bloat over time |
| FB-43 | Charm v1 majors (bubbles + glamour) | Low-Med | Low-Med | **0.5** | - | Bump `charmbracelet/bubbles` v0.21→v1.0 and `charmbracelet/glamour` v0.10→v1.0 together; both hit TUI + docs rendering layer; absorb v1 API breaks |
| FB-44 | modernc/sqlite catch-up | Medium | Low-Med | **0.5** | - | Bump `modernc.org/sqlite` v1.40→v1.49 (9 minor versions); verify storage layer + driver pool behavior unchanged via existing `internal/storage` tests |
| FB-42 | Canvas Support (OBS 30+) | High | Medium-High | **1.0** | FB-34 | New tool group for OBS 30 multi-canvas: `GetCanvasList`, `obs://canvas/*` resource, `Canvas*` event bridge |
| FB-30 | Scene Designer Filters | Medium | Low | 0.5 (via FB-33) | FB-23 ✅ | Add filter section to `scene-designer` skill |
| FB-31 | Studio Mode Skill | Medium | Medium | 0.5 (via FB-33) | FB-26 ✅ | New `studio-mode-operator` skill for preview/program workflow |
| FB-29 | New Prompts (virtual-cam-control, replay-management) | Medium | Low | — | FB-25 ✅, FB-26 ✅ | Add virtual-cam-control, replay-management prompts |
| FB-19 | Release Automation | Medium | Low | — | FB-18 ✅ | GitHub Actions workflow for automated releases |
| FB-21 | Additional Resources | Medium | Low-Med | — | - | Sources, filters, audio as MCP resources |
| FB-22 | Docker Container | Medium | Low | — | - | Containerized deployment option |
| FB-5 | Static Website | Low | Medium | — | FB-14 ✅ | Project documentation site |
| FB-6 | Network API | Low | High | — | - | Non-localhost HTTP exposure |
| FB-7 | Multi-Instance | Medium | Very High | — | - | Multiple OBS support |
| FB-8 | Remote Hosted | Low | Very High | — | FB-6, FB-7 | Cloud-hosted server |

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
