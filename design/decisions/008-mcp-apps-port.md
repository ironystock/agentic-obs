# ADR-008: Porting UI Resources to the MCP Apps Spec

**Status:** Proposed
**Date:** 2026-04-18
**Tracking:** FB-32 (Sprint 0.5 — Upstream Alignment)

## Context

`internal/mcp/ui_resources.go` currently depends on `mcpui-go` (FB-15),
an extracted standalone package implementing SEP-1865. That SEP became
the **MCP Apps Extension**, stabilized on 2026-01-26 at
`apps.extensions.modelcontextprotocol.io` in the
[`ext-apps`](https://github.com/modelcontextprotocol/ext-apps)
repository. The stable spec is a refinement (not a replacement) of
SEP-1865, but it pins several wire-level details that mcpui-go and our
consumers did not previously commit to.

### What the 2026-01-26 spec pins

| Concern | Spec (2026-01-26) | mcpui-go v0.1.0 today |
|---|---|---|
| URI scheme | `ui://` | `ui://` ✅ matches |
| Resource MIME type | `text/html;profile=mcp-app` | `mcpui.MIMETypeHTML` = `text/html` ❌ missing profile suffix |
| Client capability | `capabilities.extensions["io.modelcontextprotocol/ui"] = { mimeTypes: [...] }` at initialize | Not advertised ❌ |
| Tool→UI binding | `tool._meta.ui.resourceUri` + `visibility` | No tool-side declaration; UI resources are free-standing ❌ |
| Host→app push | `ui/notifications/tool-input`, `ui/notifications/tool-result`, `ui/notifications/tool-cancelled` | Not implemented ❌ |
| App→host calls | `postMessage` transport carrying standard `tools/call` JSON-RPC | Not implemented ❌ |
| Resource sandbox | `_meta.ui.csp` (connect/resource/frame/baseUri domains), `_meta.ui.permissions` (camera, microphone, geolocation, clipboardWrite), `_meta.ui.domain`, `_meta.ui.prefersBorder` | Not implemented ❌ |
| Display modes | `availableDisplayModes: ["inline"\|"fullscreen"\|"pip"]` | Not implemented ❌ |

### Constraints

- **SDK posture.** `github.com/modelcontextprotocol/go-sdk@v1.5.0`
  (FB-4 ✅) does not yet ship MCP-Apps-specific types. It does ship the
  primitives the spec layers on: `Annotations`, `Icons`,
  `FormElicitationCapabilities`, `URLElicitationCapabilities`. The spec
  lives in the `ext-apps` repo with JS/TS reference servers, not Go —
  so we are implementing an extension, not consuming a new SDK feature.
- **Breaking change.** Any host that renders our current `ui://`
  resources and does NOT negotiate `io.modelcontextprotocol/ui` will
  still see `text/html` (no profile suffix). Flipping the MIME type
  unilaterally risks breaking SEP-1865-only hosts that lookahead on
  exact MIME.
- **mcpui-go ownership.** We extracted mcpui-go to its own repo
  (FB-15). Anything we change upstream in mcpui-go benefits other
  consumers; anything we bolt onto agentic-obs as an ad-hoc shim does
  not. Prefer upstream changes.

## Decision

Adopt the MCP Apps 2026-01-26 spec in **three staged passes**, each
shippable on its own. Do NOT do a big-bang rewrite — the port has
real backwards-compat surface.

### Stage 1 — wire the extension capability (additive, low risk)

1. In `mcpui-go` upstream, add `MIMETypeMCPApp = "text/html;profile=mcp-app"`
   alongside the existing `MIMETypeHTML`. Do NOT change the default.
2. In agentic-obs, advertise the client capability on server init:
   ```go
   // internal/mcp/server.go
   capabilities.Extensions["io.modelcontextprotocol/ui"] = map[string]any{
       "mimeTypes": []string{"text/html;profile=mcp-app", "text/html"},
   }
   ```
3. For each of our 4 UI resources (`status/dashboard`, `scene/preview`,
   `audio/mixer`, `screenshot/gallery`), serve the profile-suffixed
   MIME when the client advertises support, fall back to plain
   `text/html` otherwise.

**Exit criteria:** MCP-Apps-aware hosts render our resources with
sandbox defaults; SEP-1865-only hosts render them unchanged.

### Stage 2 — tool→UI declarative binding (additive, medium risk)

1. Add `_meta.ui.resourceUri` to tools that have a canonical UI
   surface. Candidates (small list — most of our 81 tools are
   data-only):
   - `get_obs_status` → `ui://status/dashboard`
   - `list_scenes` → `ui://scene/preview`
   - `list_sources` with audio filter → `ui://audio/mixer`
   - `list_screenshot_sources` → `ui://screenshot/gallery`
2. Set `visibility: ["model", "app"]` so both Claude-the-reasoner and
   the embedded UI iframe receive the tool output.
3. Populate `_meta.ui.csp` with our known-safe domains (connect:
   localhost:8765 for the HTTP dashboard; resource: self; frame: none).
   Populate `_meta.ui.permissions` with an empty object (no camera /
   mic / geolocation for these views).

**Exit criteria:** A conformant host like ChatGPT Apps or the VS Code
MCP Apps client can auto-render the dashboard without user action
after calling `get_obs_status`.

### Stage 3 — bidirectional notifications (new surface, higher risk)

1. Implement `ui/notifications/tool-input` in the server so the UI
   iframe receives the tool arguments after `ui/initialize`.
2. Implement `ui/notifications/tool-result` so tool output is pushed
   to the iframe on completion.
3. Handle `ui/notifications/tool-cancelled` for user-initiated
   cancellation.
4. Define the `postMessage` contract for app→host `tools/call`. This
   reuses our existing MCP tool dispatcher — the iframe acts as
   another MCP client.

**Exit criteria:** Interactive resource UIs (e.g. the audio mixer
sliders) can call `set_input_volume` directly without going through
Claude-the-reasoner on every knob turn.

### Explicit non-goals for FB-32

- **Stinger / video-embed UI surfaces.** Out of scope; revisit once
  the streaming UX warrants it.
- **`availableDisplayModes` for pip/fullscreen.** Default to `inline`
  until we have a use case.
- **OpenAI / Claude / VSCode-specific extensions.** Stick to the core
  2026-01-26 spec.
- **Sampling integration.** Separately tracked in the FB-4 tangent
  log as "sampling" / "tool-call loop" candidates.

## Consequences

### Positive

- **Hosts get richer UX for free.** Claude Apps, ChatGPT Apps, VS Code,
  and Goose all ship conformant clients by mid-2026; we ride on that.
- **Sandbox posture improves.** Our current `ui://` resources run
  with whatever default the host provides; pinning CSP + permissions
  gives us a principled minimum.
- **Tool→UI binding formalizes** what's currently tribal knowledge
  ("call this tool then somehow render this resource").
- **Upstream mcpui-go benefits.** Other SEP-1865 consumers get the
  MIME-profile support and the notification methods without
  re-implementing them.

### Negative

- **3-stage rollout is real work.** Each stage is its own PR with its
  own test matrix. Total effort probably 2-3 working days.
- **Wire-format breaking risk.** If we ever decide to drop plain
  `text/html` support, we break SEP-1865-only hosts. Mitigate by
  keeping both MIME types advertised indefinitely — the profile
  suffix is the new truth, bare `text/html` is a permanent fallback.
- **Spec drift risk.** The 2026-01-26 spec has a "draft" successor at
  `ext-apps/blob/main/specification/draft/apps.mdx`. Track by pinning
  the stable profile string in code comments with the spec date.
- **Cross-repo coordination.** Stage 1 requires a mcpui-go release
  before agentic-obs can depend on it. Sequence: mcpui-go v0.2.0 →
  agentic-obs bumps mcpui-go → Stage 1 lands.

## Alternatives Considered

### A. Rewrite mcpui-go as an MCP Apps client from scratch

Rejected. The SEP-1865 surface mcpui-go already implements is a strict
subset of MCP Apps 2026-01-26. Incremental extension is cheaper than
replacement, and it preserves the FB-15 extraction investment.

### B. Wait for go-sdk to ship native MCP Apps types

Rejected. The go-sdk maintainers are unlikely to absorb an
extensions-track spec into the core SDK (c.f. the JS reference lives
in a separate `ext-apps` repo, not the main `@modelcontextprotocol`
core). Waiting means waiting indefinitely.

### C. Ship only Stage 1 and defer 2+3

Tempting — Stage 1 is pure additive wire-format work. But without
Stage 2's `_meta.ui.resourceUri`, hosts still don't know which tool
output to bind to which UI resource, and the UX doesn't improve
materially. Stage 1 is necessary but not sufficient; commit to
Stages 1+2 together.

### D. Implement on agentic-obs's side only, bypassing mcpui-go

Rejected. We would accumulate MCP-Apps-specific code in
`internal/mcp/` that other mcpui-go consumers couldn't reuse. Given
we own mcpui-go, the right home for Stage-1/Stage-3 primitives is
upstream.

## Implementation Checklist

This ADR intentionally stops at the decision. Implementation is
gated on:

- [ ] Stage 1: mcpui-go v0.2.0 with `MIMETypeMCPApp` + capability
      negotiation helpers
- [ ] Stage 1: agentic-obs bumps mcpui-go, advertises capability,
      dual-MIME serves resources
- [ ] Stage 2: `_meta.ui.resourceUri` on 4 tools; CSP + permissions
      defaults codified in `ui_resources.go`
- [ ] Stage 3: `ui/notifications/tool-*` methods + postMessage bridge
      (likely needs a `ui_apps_bridge.go` in agentic-obs and
      an equivalent in mcpui-go)
- [ ] Conformance test against `ext-apps/examples/basic-host`

## References

- Spec (stable): https://github.com/modelcontextprotocol/ext-apps/blob/main/specification/2026-01-26/apps.mdx
- Spec (draft): https://github.com/modelcontextprotocol/ext-apps/blob/main/specification/draft/apps.mdx
- Reference examples: https://github.com/modelcontextprotocol/ext-apps/tree/main/examples
- SEP-1865 discussion: https://github.com/modelcontextprotocol/modelcontextprotocol/pull/1865
- mcp-ui JS client SDK: https://github.com/idosal/mcp-ui
- agentic-obs mcpui-go: https://github.com/ironystock/mcpui-go
- FB-4 tangent log (SDK-side primitives this port leans on):
  `design/sprints/0.5-tangent-log.md`
