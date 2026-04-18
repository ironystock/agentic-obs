# Prompt Audit - Quick Reference Card

## Current State

| Metric | Value | Status |
|--------|-------|--------|
| Total Prompts | 13 | VERIFIED |
| Prompts in Code | 13 | MATCH |
| Prompts in Docs | 13 | SYNC |
| Documentation | Synchronized | PASS |

## Prompt Inventory

1. ✓ stream-launch
2. ✓ stream-teardown
3. ✓ audio-check
4. ✓ visual-check
5. ✓ health-check
6. ✓ problem-detection
7. ✓ preset-switcher
8. ✓ recording-workflow
9. ✓ scene-organizer
10. ✓ quick-status
11. ✓ scene-designer
12. ✓ source-management
13. ✓ visual-setup

## FB-20 Automation Tools

9 automation tools added but **NO prompt created**:
- list_automation_rules
- get_automation_rule
- create_automation_rule
- update_automation_rule
- delete_automation_rule
- enable_automation_rule
- disable_automation_rule
- trigger_automation_rule
- list_rule_executions

## Gap Analysis

| Category | Tools | Prompts | Gap |
|----------|-------|---------|-----|
| Core | 25 | 7 | COVERED |
| Sources | 3 | 1 | COVERED |
| Audio | 4 | 1 | COVERED |
| Layout | 6 | 2 | COVERED |
| Visual | 4 | 2 | COVERED |
| Design | 14 | 1 | COVERED |
| Filters | 7 | 0 | MISSING |
| Transitions | 5 | 0 | MISSING |
| **Automation** | **9** | **0** | **MISSING** |

## Recommendation

### Priority: HIGH
Add `automation-setup` prompt to:
1. Guide automation rule creation
2. Document best practices
3. Support discovery of automation capabilities

### Effort: 1-2 hours
- Implement handler function
- Update documentation
- Run tests

### Files to Change
1. `internal/mcp/prompts.go` - Add handler
2. `internal/mcp/help_content.go` - Update count to 14
3. `CLAUDE.md` - Update count and list
4. `README.md` - Update count

### Result: 13 → 14 prompts

## Detailed Reports

- **Full Report**: `AUDIT_REPORT_FB20_PROMPTS.md`
- **Summary**: `AUDIT_FINDINGS_SUMMARY.txt`
- **Implementation**: `IMPLEMENTATION_GUIDE_AUTOMATION_PROMPT.md`

## Key Findings

✓ All 13 existing prompts verified
✓ Documentation synchronized
✓ No broken references
✗ Automation feature lacks workflow guidance
✗ Filters and Transitions also lack prompts

## Risk Assessment

| Aspect | Level | Notes |
|--------|-------|-------|
| Adding automation-setup | LOW | Follows existing patterns |
| Breaking Changes | NONE | Enhancement only |
| Testing Complexity | LOW | Standard prompt handler |
| Documentation Sync | LOW | Simple count updates |

## Implementation Timeline

- [ ] Create handler: 30 min
- [ ] Test: 20 min
- [ ] Update docs: 10 min
- [ ] Review & commit: 10 min
- **Total: ~70 minutes**

---

**All information current as of 2025-12-23**
