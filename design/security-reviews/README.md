# Security Reviews

Per-PR security review artifacts produced by Claude Code's `/security-review` skill.

## Why this exists

agentic-obs has followers and forks. Anyone reading the project should be able to see, for any merged PR, what was specifically checked from a security standpoint and what the conclusion was. These files are point-in-time snapshots — they are **not** updated after merge, even if a finding turns out to be wrong or a regression is later discovered. Treat them as audit records.

## When to run a review

`/security-review` is **required pre-merge** for any PR that touches:

- `internal/http/` — HTTP server, routes, request handlers
- `internal/storage/` — SQLite schema or query construction
- `internal/mcp/` — MCP tool input schemas, resource URI parsers, elicitation guards
- `internal/obs/` — anything reading from OBS event payloads (the OBS connection is trusted, but new sinks should be re-evaluated)
- `config/`, `main.go` — startup config or env-var handling
- Any file that introduces a new external input parser (URI, JSON, query string, etc.)

For pure refactors, doc-only changes, or test-only changes, a review is **not** required.

## How to run one

In Claude Code on the PR branch:

```
/security-review
```

The skill reads the current branch's diff and produces a markdown report. Save it here as:

```
design/security-reviews/<branch-slug>_security-review.md
```

Where `<branch-slug>` is the branch name with the `feat/`, `fix/`, or `chore/` prefix stripped. Examples:

- `feat/fb-42-canvas-support` → `fb-42-canvas-support_security-review.md`
- `fix/sqlite-busy-storage` → `sqlite-busy-storage_security-review.md`

Commit the file to the PR branch so it ships with the change being reviewed.

## What a review should contain

At minimum:

- **Header**: branch name, review date, reviewer (skill version / model), outcome (`safe to merge` / `findings — fix required` / `findings — accepted as known limitation`)
- **Scope**: which files / focus areas were specifically examined, plus relevant repo context (security model, threat model)
- **Findings**: every plausible vulnerability evaluated, with verdict and reasoning. Include dismissed ones — the dismissal reasoning is itself valuable for future reviewers.
- **Conclusion**: explicit `safe to merge` / `not safe` statement

Confidence threshold for "must fix": **≥ 8** (per the skill's scoring). Findings 5–7 should be noted but don't block merge unless reviewer judgment escalates them.

## Index of reviews

| Branch | PR | Date | Outcome |
|---|---|---|---|
| `feat/fb-42-canvas-support` | [#58](https://github.com/ironystock/agentic-obs/pull/58) | 2026-04-18 | ✅ Safe to merge — [report](fb-42-canvas-support_security-review.md) |
