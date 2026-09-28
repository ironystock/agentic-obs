# ADR-009: Advanced Scene Switcher Tool Group

**Status:** Accepted
**Date:** 2026-05-28
**Tracking:** FB-51 (Sprint 1.5 candidate)

## Context

Advanced Scene Switcher (ASS) is a rule-based OBS plugin that adds a
macro engine and a variable system on top of core OBS. Streamers who
use it build automation that core `obs-websocket` cannot express — e.g.
"when `current_game` equals 'Valorant', switch to the Gaming
scene and mute the alert source". Natural-language control of OBS via
Claude Code is valuable, but only if Claude can also reach ASS macros
and variables; a missing ASS surface means half the user's automation
is invisible to the LLM.

ASS exposes a plugin-specific surface through `obs-websocket`'s
`CallVendorRequest` opcode (op 8). It registers itself as vendor
`"AdvancedSceneSwitcher"` and accepts three `requestType` values. None
of the existing OBS MCP servers (`royshil/obs-mcp`, `ironystock/agentic-obs`,
`sbroenne/mcp-server-obs`) wrap this surface — they all stop at the
core protocol.

### ASS Vendor API

All requests are fire-and-forget; ASS returns no structured data.
Variable values are strings-only on the wire (ASS rejects non-strings).
Names are case-sensitive.

| requestType | Required params | Optional params | Purpose |
|---|---|---|---|
| `AdvancedSceneSwitcherMessage` | `message: string` | — | Fire a "Websocket message received" event; macros with a matching condition react |
| `AdvancedSceneSwitcherRunMacro` | `name: string` | `variables: [{name, value}]` | Execute a named macro directly, atomically pre-setting variables in the same call |
| `AdvancedSceneSwitcherSetVariables` | `variables: [{name, value}]` | — | Bulk-set plugin variables without firing anything |

**Hard constraint:** the ASS vendor API has no introspect/list requests.
There is no wire-level way to enumerate macro names or variable names.
Callers must know what they want to call.

### Why a separate tool group, not folded into "Automation"

The existing Automation group wraps agentic-obs's own SQLite-backed rule
engine — a different system with different concepts (rules, schedules,
event triggers). ASS is a distinct plugin with its own macro/variable
namespace. Keeping them separate:

- Lets users enable/disable ASS independently (no ASS plugin? disable the group).
- Avoids conceptual confusion between agentic-obs automation rules and ASS macros.
- Follows the existing pattern of one group per conceptual domain.

## Decision

Add a 10th tool group, `AdvancedSceneSwitcher`, exposing four MCP tools:

| Tool | Vendor request | Notes |
|---|---|---|
| `ass_run_macro` | `AdvancedSceneSwitcherRunMacro` | Primary entry point; accepts optional variables for atomic set+run |
| `ass_send_message` | `AdvancedSceneSwitcherMessage` | Trigger macros via websocket-message condition |
| `ass_set_variables` | `AdvancedSceneSwitcherSetVariables` | Bulk variable update without firing any macro |
| `ass_set_variable` | `AdvancedSceneSwitcherSetVariables` (one element) | Ergonomic single-variable shortcut; same vendor call |

Variable values accept `any` at the MCP layer and are coerced to string
before forwarding. This lets the LLM write `{"value": 42}` naturally
without thinking about JSON typing.

The group defaults to **enabled** — ASS tools are low-risk (no
destructive OBS state changes) and the feature is opt-in at the user's
OBS setup level.

A `CallVendorRequest` helper is added to `internal/obs/client.go` as
a general-purpose goobs wrapper, not ASS-specific, so future vendor
plugins can use the same path.

## Consequences

- Users without ASS installed will see the four tools available but
  calls will fail with an OBS error. They can disable the group via
  `set_tool_config`.
- No macro or variable enumeration is possible from the MCP side. The
  tool descriptions state this explicitly.
- The `ass_run_macro` atomic variable-passing approach avoids a
  set-then-run race that would exist if callers used `ass_set_variables`
  followed by a separate macro trigger.
- Integration tests require a live OBS + ASS installation; unit tests
  cover validation paths and the pure payload-conversion helper.

## References

- [ASS WebSocket API](https://github.com/WarmUpTill/SceneSwitcher/wiki/Websockets)
- [goobs CallVendorRequest](https://pkg.go.dev/github.com/andreykaipov/goobs/api/requests/general)
- [ADR-004: Tool Groups](004-tool-groups.md) — the group pattern this extends
