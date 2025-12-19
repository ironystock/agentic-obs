# Architecture Decision Records (ADRs)

This directory documents significant architectural decisions for the agentic-obs project.

## ADR Index

| ADR | Title | Status | Date |
|-----|-------|--------|------|
| [001](001-sqlite-pure-go.md) | Pure Go SQLite Driver | Accepted | 2025-12-14 |
| [002](002-obs-connection.md) | Persistent 1:1 OBS Connection | Accepted | 2025-12-14 |
| [003](003-scenes-as-resources.md) | Scenes as MCP Resources | Accepted | 2025-12-14 |
| [004](004-tool-groups.md) | Configurable Tool Groups | Accepted | 2025-12-15 |
| [005](005-auth-storage.md) | SQLite Password Storage | Accepted | 2025-12-14 |
| [006](006-auto-detect-setup.md) | Auto-Detect Setup Flow | Accepted | 2025-12-14 |
| [007](007-web-ui-interfaces.md) | Web UI Interface Pattern | Accepted | 2025-12-18 |

## ADR Template

When adding a new ADR, use this template:

```markdown
# ADR-NNN: Title

**Status:** Proposed | Accepted | Deprecated | Superseded
**Date:** YYYY-MM-DD
**Supersedes:** ADR-XXX (if applicable)

## Context

What is the issue that we're seeing that is motivating this decision?

## Decision

What is the change that we're proposing and/or doing?

## Consequences

### Positive
- Benefit 1
- Benefit 2

### Negative
- Trade-off 1
- Trade-off 2

### Neutral
- Side effect 1
```

## Status Definitions

- **Proposed**: Under discussion, not yet decided
- **Accepted**: Decision has been made and is in effect
- **Deprecated**: Decision is no longer relevant but kept for history
- **Superseded**: Replaced by a newer ADR (link to replacement)
