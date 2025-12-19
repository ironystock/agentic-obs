# ADR-006: Auto-Detect Setup Flow

**Status:** Accepted
**Date:** 2025-12-14

## Context

First-run experience significantly impacts user adoption. Users should be able to start using the server quickly without complex configuration.

OBS Studio's WebSocket server has predictable defaults:
- Host: `localhost`
- Port: `4455`
- Password: Often empty for local use

## Decision

Implement a three-stage setup flow:

1. **Check saved config**: Load from SQLite if previously configured
2. **Try defaults**: Attempt `localhost:4455` with no password
3. **Interactive prompt**: On failure, prompt for host/port/password
4. **Persist success**: Save working configuration to SQLite

```
┌─────────────────────────────────────┐
│         Start Server                │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│   Check SQLite for saved config     │
└──────────────┬──────────────────────┘
               ↓
        ┌──────┴──────┐
        │ Config found?│
        └──────┬──────┘
          Yes  │  No
               ↓
┌─────────────────────────────────────┐
│   Try localhost:4455 (no password)  │
└──────────────┬──────────────────────┘
               ↓
        ┌──────┴──────┐
        │  Connected? │
        └──────┬──────┘
          Yes  │  No
               ↓
┌─────────────────────────────────────┐
│   Interactive prompt for config     │
└──────────────┬──────────────────────┘
               ↓
┌─────────────────────────────────────┐
│   Save successful config to SQLite  │
└─────────────────────────────────────┘
```

## Consequences

### Positive
- **Zero config for most users**: Default OBS setup "just works"
- **Graceful fallback**: Clear prompts when defaults fail
- **One-time setup**: Configuration persists for future runs
- **Self-documenting**: Prompts guide users through setup

### Negative
- **Console dependency**: Interactive prompts don't work in all environments
- **Blocking start**: Server waits for user input on first run
- **No headless mode**: Automation requires environment variables

### Neutral
- Environment variables can override for automated deployments
- Setup only occurs once per installation

## Interactive Prompt Format

```
OBS configuration needed. Starting interactive setup...
Note: Make sure OBS Studio is running with WebSocket server enabled.
      (Tools > WebSocket Server Settings in OBS Studio)

OBS WebSocket Host [localhost]:
OBS WebSocket Port [4455]:
OBS WebSocket Password (leave empty if none):
```

## Future Enhancements

- **TUI Setup**: Rich terminal interface with validation
- **Web Setup**: Browser-based configuration wizard
- **OBS Discovery**: Auto-detect running OBS instances on network

## References
- `config/config.go` - DetectOrPrompt function
- `main.go` - Startup flow
