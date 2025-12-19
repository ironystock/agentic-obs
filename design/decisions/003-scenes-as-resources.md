# ADR-003: Scenes as MCP Resources

**Status:** Accepted
**Date:** 2025-12-14

## Context

MCP (Model Context Protocol) provides two primary abstractions:
- **Tools**: Actions that modify state or perform operations
- **Resources**: Read-only data that can be listed, read, and monitored for changes

OBS scenes represent dynamic, server-owned state that changes based on user actions in OBS Studio. The question is whether to expose scenes through tools only, or also as resources.

## Decision

Expose OBS scenes as MCP resources with the following characteristics:

- **URI Pattern**: `obs://scene/{sceneName}`
- **Resource Type**: `obs-scene`
- **Content**: JSON representation of scene configuration (sources, settings)
- **Operations**: `resources/list`, `resources/read`

Additionally, send MCP notifications when scenes change:
- `notifications/resources/list_changed` - Scene created/deleted
- `notifications/resources/updated` - Scene modified or activated

## Consequences

### Positive
- **Client synchronization**: AI clients stay updated on scene changes without polling
- **Reactive workflows**: Enables patterns like "When I switch to Gaming scene, start recording"
- **Protocol alignment**: Resources are designed for exactly this use case in MCP
- **Rich metadata**: Resource representation includes source details, not just scene names
- **Clear separation**: Resources for state, tools for actions

### Negative
- **Event monitoring complexity**: Must subscribe to and process OBS events
- **Notification overhead**: Clients receive updates they may not need
- **Stale cache risk**: Client-side caching could show outdated state

### Neutral
- Scene tools (`list_scenes`, `set_current_scene`, etc.) remain for actions
- Resource read provides richer data than `list_scenes` tool

## Implementation Details

**OBS Events Monitored:**
- `SceneCreated` → `list_changed` notification
- `SceneRemoved` → `list_changed` notification
- `CurrentProgramSceneChanged` → `updated` notification for new active scene

**Resource Content Structure:**
```json
{
  "name": "Gaming",
  "index": 0,
  "isCurrent": true,
  "sources": [
    {"name": "Game Capture", "visible": true, "locked": false},
    {"name": "Webcam", "visible": true, "locked": false}
  ]
}
```

## References
- [MCP Resources Specification](https://modelcontextprotocol.io/docs/concepts/resources)
- [MCP Notifications](https://modelcontextprotocol.io/docs/concepts/notifications)
