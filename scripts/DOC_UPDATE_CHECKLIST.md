# Documentation Update Checklist

Run this checklist after each phase/release to ensure documentation consistency.

## Current Metrics

Update these after each phase:

| Metric | Current Value |
|--------|---------------|
| Tool Count | 45 |
| Resource Count | 4 |
| Prompt Count | 13 |
| API Endpoints | 8 |
| Current Phase | Phase 7 Complete |

---

## Files to Check

### 1. README.md (Project Root)

- [ ] Features section has correct tool/resource/prompt counts
- [ ] Tool categories list is complete and accurate
- [ ] MCP Resources section matches implementation
- [ ] MCP Prompts section matches implementation
- [ ] Installation instructions are current

### 2. CLAUDE.md (AI Context)

- [ ] Architecture diagram reflects current state
- [ ] MCP Resources section lists all 4 resource types
- [ ] MCP Tools section lists all 45 tools by category
- [ ] MCP Prompts section lists all 13 prompts with arguments
- [ ] Phase status in "Project Phases" section is current
- [ ] "Last Updated" date is correct

### 3. PROJECT_PLAN.md (Roadmap)

- [ ] Header status line reflects current phase
- [ ] Current phase marked with ✅ COMPLETE
- [ ] Success metrics section updated for current phase
- [ ] Changelog has latest entry
- [ ] Document version bumped
- [ ] "Next Review" updated to next phase

### 4. docs/README.md (Documentation Index)

- [ ] Tool count in description is correct
- [ ] Links to all documentation files work
- [ ] Tool Categories table is complete (8 categories)
- [ ] MCP Resources section present with all 3 types
- [ ] MCP Prompts section present with all 10 prompts

### 5. docs/TOOLS.md (Tool Reference)

- [ ] Tool count in header is correct (45 tools)
- [ ] All tools documented with examples
- [ ] New tools have complete documentation (including help tool)
- [ ] MCP Resources section is complete
- [ ] MCP Prompts section is complete

### 6. docs/SCREENSHOTS.md (Screenshot Guide)

- [ ] Screenshot source tools documented
- [ ] HTTP endpoint documented
- [ ] Visual monitoring workflows current

### 7. docs/API.md (HTTP API Reference)

- [ ] All endpoints documented
- [ ] Request/response examples present
- [ ] Validation rules documented
- [ ] Security considerations listed
- [ ] API endpoint count correct

### 8. examples/ Directory

- [ ] New features have example prompts
- [ ] examples/prompts/README.md lists all prompt files
- [ ] Workflow examples updated for new features

### 9. internal/mcp/help.go (Embedded Help Content)

- [ ] Tool count matches expected (45 tools)
- [ ] Resource count matches expected (4 resources)
- [ ] Prompt count matches expected (13 prompts)
- [ ] New tools have help entries in toolHelp map
- [ ] New prompts listed in prompts section
- [ ] Prompt arguments section lists all prompt arguments

---

## Automated Verification

### Quick Check (Recommended)

Run the automated verification script:

```bash
./scripts/verify-docs.sh
```

This script checks:
- Stale phase references
- Incorrect tool/resource/prompt/API endpoint counts
- Unresolved TODO items
- "Coming Soon" for completed features
- Required files exist
- API endpoints are documented

### Pre-commit Hook

Install the pre-commit hook to automatically check documentation on commit:

```bash
cp scripts/pre-commit-docs .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

### CI Integration

Documentation checks run automatically on:
- Pull requests that modify `.md` files
- Pushes to `main` that modify `.md` files

See `.github/workflows/docs-check.yml` for the workflow.

---

## Manual Consistency Checks

If you need to run checks manually:

```bash
# Check for stale phase references (update "Phase 7" to current)
grep -r "Phase [0-9] Complete" . --include="*.md" | grep -v "Phase 7"

# Check for incorrect tool counts (update "45" to current)
grep -rE "[0-9]+ (tools|Tools)" . --include="*.md" | grep -v "45"

# Check for incorrect resource counts
grep -rE "[0-9]+ (resources|Resources)" . --include="*.md" | grep -v "4"

# Check for incorrect prompt counts
grep -rE "[0-9]+ (prompts|Prompts)" . --include="*.md" | grep -v "13"

# Check for incorrect API endpoint counts
grep -rE "[0-9]+ (HTTP )?API endpoints" . --include="*.md" | grep -v "8"

# Find TODO items that should be resolved
grep -r "TODO" . --include="*.md"

# Find "Coming Soon" for completed features
grep -ri "coming soon" . --include="*.md"
```

---

## Verification Checklist

After updates, verify:

- [ ] All files agree on tool count
- [ ] All files agree on resource count
- [ ] All files agree on prompt count
- [ ] All files agree on API endpoint count
- [ ] All files agree on current phase status
- [ ] No "TODO" or "Coming Soon" for completed features
- [ ] Build passes: `go build`
- [ ] Tests pass: `go test ./...`
- [ ] No vet warnings: `go vet ./...`
- [ ] Docs check passes: `./scripts/verify-docs.sh`

---

## When to Run This Checklist

1. **After each phase completion** - Full review
2. **After PR merge** - Quick consistency check
3. **Before releases** - Full review with verification commands
4. **When adding new tools/resources/prompts** - Update metrics and relevant sections

---

## Quick Reference: File Purposes

| File | Purpose | Key Sections |
|------|---------|--------------|
| README.md | User-facing overview | Features, Installation, Usage |
| CLAUDE.md | AI assistant context | Architecture, Tools, Resources, Prompts |
| PROJECT_PLAN.md | Project status & links | Metrics, Quick links |
| CHANGELOG.md | Version history | Phases, releases, metrics |
| design/ARCHITECTURE.md | System diagrams | Communication flow, components |
| design/ROADMAP.md | Future enhancements | Planned features, research |
| design/decisions/*.md | ADRs | Technical decisions |
| docs/README.md | Documentation index | Links, Categories, Resources, Prompts, API |
| docs/TOOLS.md | Tool reference | All 45 tools with examples |
| docs/TROUBLESHOOTING.md | Common issues | Connection, Web UI, audio |
| docs/API.md | HTTP API reference | Endpoints, Validation, Security |
| internal/mcp/help_content.go | Embedded help content | Tool counts, prompts, help text |

---

## Checklists for Adding Features

### Adding a New Tool

1. **Code changes:**
   - [ ] Add tool definition in `internal/mcp/tools.go`
   - [ ] Implement handler function
   - [ ] Add to tool group registration
   - [ ] Add OBS command if needed in `internal/obs/commands.go`

2. **Documentation updates:**
   - [ ] Update tool count in `internal/mcp/help_content.go` (HelpToolCount)
   - [ ] Add help entry in `internal/mcp/help_tools.go`
   - [ ] Add tool to `docs/TOOLS.md`
   - [ ] Update tool count in `CLAUDE.md`
   - [ ] Update tool count in `scripts/verify-docs.sh`
   - [ ] Add to relevant example file in `examples/prompts/`

3. **Verify:**
   - [ ] Run `./scripts/verify-docs.sh`
   - [ ] Run `go test ./...`

### Adding a New Prompt

1. **Code changes:**
   - [ ] Add prompt definition in `internal/mcp/prompts.go`
   - [ ] Add message generation in `handleGetPrompt()`
   - [ ] Add completion support if arguments needed

2. **Documentation updates:**
   - [ ] Update prompt count in `internal/mcp/help_content.go` (HelpPromptCount)
   - [ ] Add prompt help in `internal/mcp/help_content.go`
   - [ ] Update prompt count in `CLAUDE.md`
   - [ ] Update prompt count in `scripts/verify-docs.sh`
   - [ ] Add to `examples/prompts/README.md` mapping table

3. **Verify:**
   - [ ] Run `./scripts/verify-docs.sh`
   - [ ] Run `go test ./...`

### Adding a New Resource

1. **Code changes:**
   - [ ] Add to `handleResourcesList()` in `internal/mcp/resources.go`
   - [ ] Add case in `handleResourceRead()`
   - [ ] Add OBS event handlers for notifications

2. **Documentation updates:**
   - [ ] Update resource count in `internal/mcp/help_content.go` (HelpResourceCount)
   - [ ] Add resource documentation in help content
   - [ ] Update resource count in `CLAUDE.md`
   - [ ] Update resource count in `scripts/verify-docs.sh`
   - [ ] Add to `design/ARCHITECTURE.md` resources table

3. **Verify:**
   - [ ] Run `./scripts/verify-docs.sh`
   - [ ] Run `go test ./...`

### Making an Architectural Decision

1. **Create ADR:**
   - [ ] Create `design/decisions/NNN-title.md` using template
   - [ ] Add to index in `design/decisions/README.md`

2. **Documentation updates:**
   - [ ] Reference ADR in relevant documentation
   - [ ] Update `design/ARCHITECTURE.md` if affects architecture

---

**Last Updated:** 2025-12-19
**Created:** Phase 4 completion
**Updated:** Phase 7+ - Documentation restructuring, design/ directory, checklists
