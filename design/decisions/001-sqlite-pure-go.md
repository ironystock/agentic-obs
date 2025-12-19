# ADR-001: Pure Go SQLite Driver

**Status:** Accepted
**Date:** 2025-12-14

## Context

The agentic-obs server requires persistent storage for:
- OBS connection configuration (host, port, password)
- Scene presets (source visibility states)
- Screenshot sources and captured images
- Application state and action history

Two main SQLite driver options exist for Go:
1. `mattn/go-sqlite3` - CGO-based, wraps the official SQLite C library
2. `modernc.org/sqlite` - Pure Go implementation, no CGO required

## Decision

Use `modernc.org/sqlite` as the SQLite driver.

## Consequences

### Positive
- **No CGO requirement**: Simplified cross-compilation to Windows, macOS, and Linux
- **Simpler build process**: End users don't need C compiler toolchain installed
- **Single binary distribution**: No shared library dependencies
- **Consistent behavior**: Same code path on all platforms

### Negative
- **~2x slower writes**: Pure Go implementation is slower than C library
- **Larger binary size**: Pure Go SQLite adds ~8MB to the binary

### Neutral
- Performance difference negligible for this use case (config/preset storage, not high-throughput)
- Single-user, single-tenant deployment means no concurrent write pressure

## Alternatives Considered

### mattn/go-sqlite3
- Faster performance (native C)
- Requires CGO and C compiler
- Cross-compilation complexity
- **Rejected**: Build complexity outweighs performance benefits for this use case

## References
- [Go SQLite Benchmarks](https://github.com/cvilsmeier/go-sqlite-bench)
- [modernc.org/sqlite documentation](https://pkg.go.dev/modernc.org/sqlite)
