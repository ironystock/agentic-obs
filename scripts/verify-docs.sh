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
EXPECTED_TOOLS=30
EXPECTED_RESOURCES=3
EXPECTED_PROMPTS=10
CURRENT_PHASE=4

echo "=========================================="
echo "Documentation Consistency Verification"
echo "=========================================="
echo ""
echo "Expected values:"
echo "  Tools: $EXPECTED_TOOLS"
echo "  Resources: $EXPECTED_RESOURCES"
echo "  Prompts: $EXPECTED_PROMPTS"
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
