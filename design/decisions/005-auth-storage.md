# ADR-005: SQLite Password Storage

**Status:** Accepted
**Date:** 2025-12-14

## Context

OBS WebSocket server can require password authentication. The MCP server needs to store this password for automatic reconnection across restarts.

Options considered:
1. **Environment variables**: Standard for secrets, but not persistent
2. **Interactive prompt**: Secure, but poor UX for every run
3. **Encrypted storage**: Most secure, but complex key management
4. **Plain SQLite storage**: Simple, relies on filesystem permissions

## Decision

Store the OBS WebSocket password in SQLite database without encryption.

**Security Model:**
- Local, single-user deployment only
- Database file permissions protect credentials (user-readable only)
- No network exposure of credentials
- Password only transmitted over localhost WebSocket

## Consequences

### Positive
- **Simplicity**: No key management or encryption complexity
- **Persistence**: Password survives restarts without re-prompting
- **Automatic reconnection**: Server can reconnect after OBS restart

### Negative
- **Unencrypted storage**: Password visible in database file
- **Shared machine risk**: Other users with file access can read password
- **Not enterprise-ready**: Doesn't meet compliance requirements

### Neutral
- Acceptable trade-off for local, single-user use case
- OBS WebSocket password is low-sensitivity (controls local OBS only)

## Mitigations

1. Database file created with restrictive permissions (0600)
2. Documentation warns against shared machine deployment
3. Users can use environment variables for sensitive deployments

## Alternatives Rejected

### Environment Variables Only
- No persistence across terminal sessions
- Poor UX for interactive use
- **Rejected**: Convenience outweighs for local deployment

### Encrypted Storage (SQLCipher/Age)
- Complex key management (where to store encryption key?)
- Over-engineering for localhost use case
- **Rejected**: Complexity not justified for threat model

### OS Keychain Integration
- Platform-specific implementations
- Complex dependency management
- **Rejected**: Cross-platform simplicity prioritized

## References
- `internal/storage/state.go` - OBS config storage
- `config/config.go` - Configuration management
