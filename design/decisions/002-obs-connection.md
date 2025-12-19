# ADR-002: Persistent 1:1 OBS Connection

**Status:** Accepted
**Date:** 2025-12-14

## Context

The MCP server needs to communicate with OBS Studio via the WebSocket API. Several connection strategies are possible:

1. **Connect per request**: Open connection for each tool call, close after
2. **Persistent single connection**: Maintain one connection throughout server lifetime
3. **Connection pool**: Multiple connections to one or more OBS instances

The server also needs to receive OBS events for MCP resource notifications (scene changes, etc.).

## Decision

Maintain a single persistent WebSocket connection to one OBS instance with automatic reconnection on failure.

## Consequences

### Positive
- **Reduced latency**: No connection overhead per tool call
- **Event support**: Persistent connection enables OBS event monitoring
- **Resource notifications**: Can dispatch MCP notifications when scenes change
- **Simpler implementation**: Single connection state to manage
- **Covers 95% of use cases**: Most users control one local OBS instance

### Negative
- **Single instance only**: Cannot control multiple OBS instances simultaneously
- **Connection dependency**: Server functionality depends on OBS availability
- **Reconnection complexity**: Must handle disconnections gracefully

### Neutral
- Auto-reconnect with exponential backoff handles temporary disconnections
- Connection state is clearly visible in status endpoints

## Future Considerations

Multi-instance support would require architectural changes:
- Connection pooling with instance ID routing
- Namespaced resource URIs (e.g., `obs://instance1/scene/Gaming`)
- Per-instance tool parameters or separate MCP servers

**Decision Point:** Evaluate multi-instance architecture if user demand materializes.

## References
- [goobs WebSocket client](https://pkg.go.dev/github.com/andreykaipov/goobs)
- [OBS WebSocket Protocol](https://github.com/obsproject/obs-websocket/blob/master/docs/generated/protocol.md)
