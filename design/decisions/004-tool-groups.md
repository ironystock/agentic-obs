# ADR-004: Configurable Tool Groups

**Status:** Accepted
**Date:** 2025-12-15

## Context

As the server grew to 45+ tools, exposing all tools to every AI client became problematic:
- **Context pollution**: AI models have limited context; many tools may be irrelevant
- **Security concerns**: Some users may want to limit destructive capabilities
- **Use case specialization**: A recording assistant doesn't need scene design tools

Users needed a way to customize which tool categories are available.

## Decision

Organize tools into logical groups that can be individually enabled/disabled:

| Group | Tools | Description |
|-------|-------|-------------|
| **Core** | 13 | Scenes, recording, streaming, status |
| **Visual** | 4 | Screenshot sources |
| **Layout** | 6 | Scene presets |
| **Audio** | 4 | Input volume and mute control |
| **Sources** | 3 | Source visibility and settings |
| **Design** | 14 | Source creation and transforms |
| **Help** | 1 | Always enabled |

Configuration is stored in SQLite and prompted during first-run setup.

## Consequences

### Positive
- **Reduced context**: Clients only see relevant tools
- **Customizable**: Users control their AI assistant's capabilities
- **Security**: Can disable design tools to prevent accidental scene changes
- **Persistence**: Preferences survive server restarts

### Negative
- **Complexity**: Conditional tool registration logic
- **Discoverability**: Users may not know disabled tools exist
- **Testing burden**: Must test all combinations

### Neutral
- Default configuration enables all groups
- Help tool is always enabled regardless of settings

## Configuration Structure

```go
type ToolGroupConfig struct {
    Core    bool // Core OBS tools (scenes, recording, streaming, status)
    Visual  bool // Visual monitoring tools (screenshots)
    Layout  bool // Layout management tools (scene presets)
    Audio   bool // Audio control tools
    Sources bool // Source management tools
    Design  bool // Scene design tools (source creation, transforms)
}
```

## First-Run Experience

```
=== Feature Configuration ===
Choose which features to enable (press Enter for defaults):

--- Tool Groups ---
Core OBS control (scenes, recording, streaming) [Y/n]:
Visual monitoring (screenshot capture) [Y/n]:
Layout management (scene presets) [Y/n]:
Audio control (mute, volume) [Y/n]:
Source management (visibility, settings) [Y/n]:
Scene design (create sources, transforms) [Y/n]:
```

## References
- `config/config.go` - ToolGroupConfig struct
- `internal/mcp/tools.go` - Conditional registration
