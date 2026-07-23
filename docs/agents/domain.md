# Domain Docs

How the engineering skills should consume this repo's domain documentation when exploring the codebase.

## Before exploring, read these

- **`CONTEXT.md`** at the repo root (created lazily by `/domain-modeling` — may not exist yet).
- **`design/decisions/`** — this repo's ADR home. Read the ADRs that touch the area you're about to work in, and `design/decisions/README.md` for the index. (This is the repo's real ADR location; it does **not** use `docs/adr/`.)

If any of these files don't exist, **proceed silently**. Don't flag their absence; don't suggest creating them upfront. The `/domain-modeling` skill (reached via `/grill-with-docs` and `/improve-codebase-architecture`) creates them lazily when terms or decisions actually get resolved.

## File structure

Single-context repo:

```
/
├── CONTEXT.md                 ← domain glossary (created lazily)
├── design/
│   └── decisions/             ← ADRs (001-…, 002-…) + README index
└── internal/                  ← Go source
```

## Use the glossary's vocabulary

When your output names a domain concept (in an issue title, a refactor proposal, a hypothesis, a test name), use the term as defined in `CONTEXT.md`. Don't drift to synonyms the glossary explicitly avoids.

If the concept you need isn't in the glossary yet, that's a signal — either you're inventing language the project doesn't use (reconsider) or there's a real gap (note it for `/domain-modeling`).

## Flag ADR conflicts

If your output contradicts an existing ADR, surface it explicitly rather than silently overriding:

> _Contradicts ADR-005 (OBS password stored unencrypted) — but worth reopening because…_

New ADRs follow the existing numbering and format in `design/decisions/` (Status / Date / Context / Decision / Consequences / Alternatives), and get added to `design/decisions/README.md`.
