# Security Review — PR #58 (FB-42 Canvas Support)

**Branch:** `feat/fb-42-canvas-support`
**Reviewed:** 2026-04-18
**Reviewer:** `/security-review` skill (Claude Opus 4.7)
**Outcome:** ✅ No vulnerabilities at confidence ≥ 8 — safe to merge

---

## Scope

Reviewed all 27 modified files in PR #58 with focus on the user-flagged areas:

1. URI parsing in `extractCanvasNameFromURI` — path traversal / injection
2. JSON round-trip decoder in `GetCanvasList` — malformed-input handling
3. New HTTP body field `canvas` in `handleUpdateConfig` — input validation
4. New SQLite write path for `tools_enabled_canvas`
5. Canvas event bridge (`subscriptions.Canvases` mask) — trust assumptions about OBS-side payloads
6. Incidental Automation zero-value fix in `main.go` and `handlers.go` — change in attack surface

Repo context applied:
- Localhost-bound HTTP server with no auth (per [ADR-007](../decisions/007-web-ui-interfaces.md))
- OBS WebSocket password stored unencrypted in SQLite (documented constraint in CLAUDE.md)
- Heavy use of Go `internal/` packages (limits external API surface)
- MCP server uses stdio transport for Claude Code

---

## Summary

**No high-confidence security vulnerabilities were newly introduced by this PR.**

All flagged areas use established secure patterns already present in the codebase (parameterized SQL, JSON-encoded outputs, byte-bounded URI parsing, localhost-only HTTP per ADR-007). The PR is consistent with prior tool-group additions and does not introduce new trust boundaries or sinks.

---

## Findings Evaluated (all dismissed, confidence < 8)

### 1. `extractCanvasNameFromURI` — path traversal / injection

**File:** `internal/mcp/resources.go:560-568`
**Verdict:** Not a vulnerability.

Extracted name is used only for in-memory string equality against `canvases[i].Name` from OBS, then echoed via `json.MarshalIndent` (safely escaped) and `fmt.Sprintf("%q", ...)` (escapes). Never reaches a filesystem, SQL, shell, or template sink. A URI like `obs://canvas/../../etc/passwd` yields the literal string, fails the equality check, and returns `not found`. Pattern is identical to the existing accepted `extractScreenshotURLNameFromURI`.

### 2. JSON round-trip in `GetCanvasList`

**File:** `internal/obs/commands.go:1517-1542`
**Verdict:** Not a vulnerability.

`encoding/json` is safe against malformed input — returns errors, doesn't panic, doesn't execute code. No deserialization gadget exists in Go's stdlib `encoding/json`.

### 3. New `canvas` and `automation` fields in `handleUpdateConfig`

**File:** `internal/http/handlers.go:218-219`
**Verdict:** Not a vulnerability.

Boolean inputs read via `getBool(...)` and persisted via parameterized SQL. No injection surface. The HTTP handler is documented as localhost-bound + unauthenticated per ADR-007 (pre-existing posture). The new fields don't widen that posture — anyone reachable on the localhost socket could already toggle any of the other 9 tool groups.

### 4. New SQLite write path `tools_enabled_canvas`

**File:** `internal/storage/state.go:381-383, 422-424`
**Verdict:** Not a vulnerability.

Same `SetState`/`GetState` parameterized pattern as the other 9 tool groups. Constant key, bool value via `boolToStr`. No injection vector.

### 5. Canvas event bridge — trust assumptions about OBS payloads

**Files:** `internal/obs/client.go:135`, `internal/obs/events.go:243-271`, `internal/mcp/server.go:444-462`
**Verdict:** Not a vulnerability.

Canvas event names from OBS flow into `GetResourceURIForCanvas(name)` → `obs://canvas/<name>` → MCP client. OBS is already a fully trusted local component (the server holds its WebSocket password). A malicious OBS would already have arbitrary control over scenes/sources/recordings — no new trust boundary is crossed. All string flows are JSON-encoded throughout.

### 6. Incidental "Automation enabled by default on existing installs"

**Files:** `main.go:101`, `internal/http/handlers.go:218`
**Verdict:** Not a vulnerability.

The fix restores the *intended* default (`Automation: true`, set in `DefaultToolGroupConfig`). Existing installs that had explicitly disabled Automation via `set_tool_config` still load that preference from SQLite via `LoadToolGroupConfig` — the fix only affects the in-memory struct constructed at server launch / on HTTP config updates that omit the field. The Automation tool group itself was already shipped (FB-20). High-risk tools like `delete_automation_rule` already have elicitation-confirmation guards. No new tool exposed; no new attack surface.

---

## Conclusion

**Zero confirmed vulnerabilities at ≥8 confidence. PR #58 is safe to merge from a security standpoint.**
