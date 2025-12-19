# ADR-007: Web UI Interface Pattern (StatusProvider/ActionExecutor)

**Status:** Accepted
**Date:** 2025-12-18

## Context

Phase 6.1 introduced a web dashboard that needed to:
1. Display OBS status (connection, recording, streaming state)
2. Show scene list with thumbnails
3. Display audio mixer with volume controls
4. Execute actions (switch scenes, toggle mute, adjust volume)

The HTTP server in `internal/http/` needed access to MCP server state and OBS commands, but direct coupling would create circular dependencies and complicate testing.

## Decision

Define two interfaces that the MCP server implements:

```go
// StatusProvider reads OBS state for display
type StatusProvider interface {
    GetStatus() (any, error)
    GetScenes() ([]SceneInfo, error)
    GetAudioInputs() ([]AudioInputInfo, error)
    GetScreenshotSources() ([]ScreenshotSourceInfo, error)
}

// ActionExecutor performs OBS actions
type ActionExecutor interface {
    SetCurrentScene(sceneName string) error
    ToggleInputMute(inputName string) error
    SetInputVolume(inputName string, volumeDb float64) error
    TakeSceneThumbnail(sceneName string) ([]byte, string, error)
}
```

The HTTP server receives these interfaces via `SetStatusProvider()`, enabling loose coupling.

## Consequences

### Positive
- **Decoupling**: HTTP server doesn't import MCP package
- **Testability**: Can mock interfaces for HTTP handler tests
- **Clear contracts**: Interfaces document exactly what UI needs
- **No circular imports**: Clean dependency graph

### Negative
- **Interface maintenance**: Must update interfaces when adding UI features
- **Indirection**: Extra layer between HTTP handlers and OBS client
- **Type assertions**: Status provider returns `any` for flexibility

### Neutral
- Pattern is common in Go for breaking package cycles
- Interface compliance verified at compile time

## Implementation Details

**HTTP Server Setup:**
```go
// In internal/http/server.go
func (s *Server) SetStatusProvider(sp StatusProvider) error {
    if sp == nil {
        return fmt.Errorf("status provider cannot be nil")
    }
    s.statusProvider = sp
    // Check if also implements ActionExecutor
    if ae, ok := sp.(ActionExecutor); ok {
        s.actionExecutor = ae
    }
    return nil
}
```

**MCP Server implements both:**
```go
// In internal/mcp/status_provider.go
var _ agenthttp.StatusProvider = (*Server)(nil)
var _ agenthttp.ActionExecutor = (*Server)(nil)
```

**Compile-time verification** ensures MCP Server satisfies both interfaces.

## Thumbnail Caching

The `TakeSceneThumbnail` method includes intelligent caching:
- 5-second TTL per scene
- Cache invalidation on scene changes
- Background cleanup of expired entries

This prevents excessive OBS API calls when multiple clients request thumbnails.

## References
- `internal/http/server.go` - Interface definitions
- `internal/mcp/status_provider.go` - Interface implementation
- `internal/mcp/server.go` - ThumbnailCache
