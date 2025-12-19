#!/bin/bash
# Documentation Consistency Verification Script
# Run this after each phase/release to verify documentation consistency
#
# Usage: ./scripts/verify-docs.sh
# Exit codes: 0 = all checks pass, 1 = issues found

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Current expected values - UPDATE THESE AFTER EACH PHASE
EXPECTED_TOOLS=45
EXPECTED_RESOURCES=4
EXPECTED_PROMPTS=13
EXPECTED_API_ENDPOINTS=8
CURRENT_PHASE=7

echo "=========================================="
echo "Documentation Consistency Verification"
echo "=========================================="
echo ""
echo "Expected values:"
echo "  Tools: $EXPECTED_TOOLS"
echo "  Resources: $EXPECTED_RESOURCES"
echo "  Prompts: $EXPECTED_PROMPTS"
echo "  API Endpoints: $EXPECTED_API_ENDPOINTS"
echo "  Phase: $CURRENT_PHASE"
echo ""

ISSUES_FOUND=0

# Function to check and report
check_pattern() {
    local description="$1"
    local pattern="$2"
    local exclude="$3"

    echo -n "Checking: $description... "

    # Run grep and capture results, excluding the expected value
    if [ -n "$exclude" ]; then
        RESULTS=$(grep -rE "$pattern" . --include="*.md" 2>/dev/null | grep -v "$exclude" | grep -v "verify-docs.sh" | grep -v "DOC_UPDATE_CHECKLIST.md" || true)
    else
        RESULTS=$(grep -rE "$pattern" . --include="*.md" 2>/dev/null | grep -v "verify-docs.sh" | grep -v "DOC_UPDATE_CHECKLIST.md" || true)
    fi

    if [ -z "$RESULTS" ]; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}ISSUES FOUND${NC}"
        echo "$RESULTS" | while read -r line; do
            echo -e "  ${YELLOW}→${NC} $line"
        done
        ISSUES_FOUND=1
    fi
}

# Check for stale phase references
echo "--- Phase Consistency ---"
check_pattern "Stale phase references" "Phase [0-9] Complete" "Phase $CURRENT_PHASE"

# Check for incorrect TOTAL counts (not category counts like "4 tools")
echo ""
echo "--- Metric Consistency ---"
# Only check for "X MCP tools" or "Total: X tools" patterns, not category headers
check_pattern "Incorrect total tool counts" "(Total|total|MCP).*[0-9]+ (tools|Tools)" "$EXPECTED_TOOLS"
check_pattern "Incorrect resource counts" "[0-9]+ (resources|Resources)" "$EXPECTED_RESOURCES"
check_pattern "Incorrect prompt counts" "[0-9]+ (prompts|Prompts)" "$EXPECTED_PROMPTS"
check_pattern "Incorrect API endpoint counts" "[0-9]+ (HTTP )?API endpoints" "$EXPECTED_API_ENDPOINTS"

# Check for unresolved items
echo ""
echo "--- Unresolved Items ---"
echo -n "Checking: TODO items... "
TODO_RESULTS=$(grep -r "TODO" . --include="*.md" 2>/dev/null | grep -v "verify-docs.sh" | grep -v "DOC_UPDATE_CHECKLIST.md" || true)
if [ -z "$TODO_RESULTS" ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${YELLOW}FOUND (review if still valid)${NC}"
    echo "$TODO_RESULTS" | while read -r line; do
        echo -e "  ${YELLOW}→${NC} $line"
    done
fi

echo -n "Checking: 'Coming Soon' references... "
COMING_SOON=$(grep -ri "coming soon" . --include="*.md" 2>/dev/null | grep -v "verify-docs.sh" | grep -v "DOC_UPDATE_CHECKLIST.md" || true)
if [ -z "$COMING_SOON" ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}ISSUES FOUND${NC}"
    echo "$COMING_SOON" | while read -r line; do
        echo -e "  ${YELLOW}→${NC} $line"
    done
    ISSUES_FOUND=1
fi

# Check required files exist
echo ""
echo "--- Required Files ---"
REQUIRED_FILES=(
    "README.md"
    "CLAUDE.md"
    "PROJECT_PLAN.md"
    "docs/README.md"
    "docs/TOOLS.md"
    "docs/SCREENSHOTS.md"
    "docs/QUICKSTART.md"
    "docs/API.md"
    "skills/README.md"
)

for file in "${REQUIRED_FILES[@]}"; do
    echo -n "Checking: $file exists... "
    if [ -f "$file" ]; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}MISSING${NC}"
        ISSUES_FOUND=1
    fi
done

# Check API endpoint documentation matches code
echo ""
echo "--- API Documentation Consistency ---"
API_ENDPOINTS=(
    "/api/status"
    "/api/history"
    "/api/history/stats"
    "/api/screenshots"
    "/api/config"
    "/screenshot/"
)

for endpoint in "${API_ENDPOINTS[@]}"; do
    echo -n "Checking: $endpoint documented... "
    if grep -q "$endpoint" docs/API.md 2>/dev/null; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}MISSING${NC}"
        ISSUES_FOUND=1
    fi
done

# Check help_content.go constants match expected values
echo ""
echo "--- Help Tool Content Consistency ---"

# Check HelpToolCount constant (use word boundary to avoid matching HelpHelpToolCount)
echo -n "Checking: HelpToolCount constant ($EXPECTED_TOOLS)... "
HELP_TOOL_CONST=$(grep -E "^\s*HelpToolCount\s*=" internal/mcp/help_content.go 2>/dev/null | head -1 | awk -F'=' '{print $2}' | awk '{print $1}' || echo "0")
if [ "$HELP_TOOL_CONST" = "$EXPECTED_TOOLS" ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}ISSUES FOUND${NC}"
    echo -e "  ${YELLOW}→${NC} HelpToolCount is $HELP_TOOL_CONST (expected $EXPECTED_TOOLS)"
    ISSUES_FOUND=1
fi

# Check HelpPromptCount constant
echo -n "Checking: HelpPromptCount constant ($EXPECTED_PROMPTS)... "
HELP_PROMPT_CONST=$(grep -E "^\s*HelpPromptCount\s*=" internal/mcp/help_content.go 2>/dev/null | head -1 | awk -F'=' '{print $2}' | awk '{print $1}' || echo "0")
if [ "$HELP_PROMPT_CONST" = "$EXPECTED_PROMPTS" ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}ISSUES FOUND${NC}"
    echo -e "  ${YELLOW}→${NC} HelpPromptCount is $HELP_PROMPT_CONST (expected $EXPECTED_PROMPTS)"
    ISSUES_FOUND=1
fi

# Check HelpResourceCount constant
echo -n "Checking: HelpResourceCount constant ($EXPECTED_RESOURCES)... "
HELP_RESOURCE_CONST=$(grep -E "^\s*HelpResourceCount\s*=" internal/mcp/help_content.go 2>/dev/null | head -1 | awk -F'=' '{print $2}' | awk '{print $1}' || echo "0")
if [ "$HELP_RESOURCE_CONST" = "$EXPECTED_RESOURCES" ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}ISSUES FOUND${NC}"
    echo -e "  ${YELLOW}→${NC} HelpResourceCount is $HELP_RESOURCE_CONST (expected $EXPECTED_RESOURCES)"
    ISSUES_FOUND=1
fi

# Check tool count in help_content.go text
echo -n "Checking: help_content.go tool count text ($EXPECTED_TOOLS)... "
HELP_TOOL_ISSUES=$(grep -E "All Available Tools \([0-9]+ total\)" internal/mcp/help_content.go 2>/dev/null | grep -v "$EXPECTED_TOOLS total" || true)
if [ -z "$HELP_TOOL_ISSUES" ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}ISSUES FOUND${NC}"
    echo -e "  ${YELLOW}→${NC} help_content.go has incorrect tool count text (expected $EXPECTED_TOOLS)"
    ISSUES_FOUND=1
fi

# Check prompt count in help_content.go text
echo -n "Checking: help_content.go prompt count text ($EXPECTED_PROMPTS)... "
HELP_PROMPT_ISSUES=$(grep -E "MCP Prompts \([0-9]+ workflows\)" internal/mcp/help_content.go 2>/dev/null | grep -v "$EXPECTED_PROMPTS workflows" || true)
if [ -z "$HELP_PROMPT_ISSUES" ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}ISSUES FOUND${NC}"
    echo -e "  ${YELLOW}→${NC} help_content.go has incorrect prompt count text (expected $EXPECTED_PROMPTS)"
    ISSUES_FOUND=1
fi

# Check resource count in help_content.go text
echo -n "Checking: help_content.go resource count text ($EXPECTED_RESOURCES)... "
HELP_RESOURCE_ISSUES=$(grep -E "MCP Resources \([0-9]+ types\)" internal/mcp/help_content.go 2>/dev/null | grep -v "$EXPECTED_RESOURCES types" || true)
if [ -z "$HELP_RESOURCE_ISSUES" ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}ISSUES FOUND${NC}"
    echo -e "  ${YELLOW}→${NC} help_content.go has incorrect resource count text (expected $EXPECTED_RESOURCES)"
    ISSUES_FOUND=1
fi

# Summary
echo ""
echo "=========================================="
if [ $ISSUES_FOUND -eq 0 ]; then
    echo -e "${GREEN}All documentation checks passed!${NC}"
    exit 0
else
    echo -e "${RED}Documentation issues found - please review above${NC}"
    exit 1
fi
