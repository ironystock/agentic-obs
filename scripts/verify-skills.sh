#!/bin/bash
# Claude Skills Validation Script
# Validates the structure and content of skills in the skills/ directory
#
# Usage: ./scripts/verify-skills.sh
# Exit codes: 0 = all checks pass, 1 = issues found

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Expected skills
EXPECTED_SKILLS=(
    "streaming-assistant"
    "scene-designer"
    "audio-engineer"
    "preset-manager"
)

# Required frontmatter fields
REQUIRED_FRONTMATTER=(
    "name"
    "description"
)

# Required sections (at least one of these patterns should match)
REQUIRED_SECTIONS=(
    "## When to Use This Skill"
    "## Workflow Steps"
    "## Core Responsibilities"
)

echo "=========================================="
echo "Claude Skills Validation"
echo "=========================================="
echo ""
echo "Expected skills: ${#EXPECTED_SKILLS[@]}"
echo "  - streaming-assistant"
echo "  - scene-designer"
echo "  - audio-engineer"
echo "  - preset-manager"
echo ""

ISSUES_FOUND=0

# Check that skills directory exists
echo "--- Directory Structure ---"
echo -n "Checking: skills/ directory exists... "
if [ -d "skills" ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}MISSING${NC}"
    echo "The skills/ directory does not exist!"
    exit 1
fi

# Check that skills/README.md exists
echo -n "Checking: skills/README.md exists... "
if [ -f "skills/README.md" ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}MISSING${NC}"
    ISSUES_FOUND=1
fi

echo ""
echo "--- Skill Directories ---"

# Check each expected skill directory exists
for skill in "${EXPECTED_SKILLS[@]}"; do
    echo -n "Checking: skills/$skill/ exists... "
    if [ -d "skills/$skill" ]; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}MISSING${NC}"
        ISSUES_FOUND=1
    fi
done

echo ""
echo "--- Skill Files ---"

# Check each skill has a SKILL.md file
for skill in "${EXPECTED_SKILLS[@]}"; do
    echo -n "Checking: skills/$skill/SKILL.md exists... "
    if [ -f "skills/$skill/SKILL.md" ]; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}MISSING${NC}"
        ISSUES_FOUND=1
    fi
done

echo ""
echo "--- Frontmatter Validation ---"

# Check frontmatter in each SKILL.md
for skill in "${EXPECTED_SKILLS[@]}"; do
    SKILL_FILE="skills/$skill/SKILL.md"

    if [ ! -f "$SKILL_FILE" ]; then
        continue
    fi

    echo "Validating: $skill"

    # Check for frontmatter delimiters
    echo -n "  Frontmatter delimiters (---)... "
    FRONTMATTER_COUNT=$(grep -c "^---$" "$SKILL_FILE" 2>/dev/null || echo "0")
    if [ "$FRONTMATTER_COUNT" -ge 2 ]; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}MISSING${NC}"
        echo -e "    ${YELLOW}→${NC} SKILL.md must have YAML frontmatter with --- delimiters"
        ISSUES_FOUND=1
    fi

    # Check for required frontmatter fields
    for field in "${REQUIRED_FRONTMATTER[@]}"; do
        echo -n "  Frontmatter field: $field... "
        if grep -q "^$field:" "$SKILL_FILE" 2>/dev/null; then
            echo -e "${GREEN}OK${NC}"
        else
            echo -e "${RED}MISSING${NC}"
            ISSUES_FOUND=1
        fi
    done

    # Verify name matches directory
    echo -n "  Frontmatter name matches directory... "
    FRONTMATTER_NAME=$(grep "^name:" "$SKILL_FILE" 2>/dev/null | sed 's/name: *//' | tr -d '\r' || echo "")
    if [ "$FRONTMATTER_NAME" = "$skill" ]; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}MISMATCH${NC}"
        echo -e "    ${YELLOW}→${NC} Expected: $skill, Found: $FRONTMATTER_NAME"
        ISSUES_FOUND=1
    fi
done

echo ""
echo "--- Content Structure ---"

# Check for required sections in each SKILL.md
for skill in "${EXPECTED_SKILLS[@]}"; do
    SKILL_FILE="skills/$skill/SKILL.md"

    if [ ! -f "$SKILL_FILE" ]; then
        continue
    fi

    echo "Validating: $skill"

    # Check for at least one required section pattern
    echo -n "  Required sections present... "
    SECTION_FOUND=0
    for section in "${REQUIRED_SECTIONS[@]}"; do
        if grep -q "^$section" "$SKILL_FILE" 2>/dev/null; then
            SECTION_FOUND=1
            break
        fi
    done

    if [ $SECTION_FOUND -eq 1 ]; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}MISSING${NC}"
        echo -e "    ${YELLOW}→${NC} Must have one of: ${REQUIRED_SECTIONS[*]}"
        ISSUES_FOUND=1
    fi

    # Check for at least one heading
    echo -n "  Markdown headings present... "
    HEADING_COUNT=$(grep -c "^#" "$SKILL_FILE" 2>/dev/null || echo "0")
    if [ "$HEADING_COUNT" -gt 0 ]; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${RED}MISSING${NC}"
        echo -e "    ${YELLOW}→${NC} SKILL.md should have markdown headings"
        ISSUES_FOUND=1
    fi

    # Check for non-empty content (more than just frontmatter)
    echo -n "  Substantial content present... "
    CONTENT_LINES=$(wc -l < "$SKILL_FILE" 2>/dev/null || echo "0")
    if [ "$CONTENT_LINES" -gt 50 ]; then
        echo -e "${GREEN}OK${NC}"
    else
        echo -e "${YELLOW}WARNING${NC}"
        echo -e "    ${YELLOW}→${NC} SKILL.md seems short ($CONTENT_LINES lines), consider adding more detail"
    fi
done

echo ""
echo "--- skills/README.md References ---"

# Check that skills/README.md references all skills
if [ -f "skills/README.md" ]; then
    for skill in "${EXPECTED_SKILLS[@]}"; do
        echo -n "Checking: README references $skill... "
        if grep -q "$skill" "skills/README.md" 2>/dev/null; then
            echo -e "${GREEN}OK${NC}"
        else
            echo -e "${RED}MISSING${NC}"
            echo -e "  ${YELLOW}→${NC} skills/README.md should reference $skill"
            ISSUES_FOUND=1
        fi
    done
fi

echo ""
echo "--- Additional Checks ---"

# Check for unexpected skill directories
echo -n "Checking: No unexpected skill directories... "
EXTRA_SKILLS=$(find skills -mindepth 1 -maxdepth 1 -type d ! -name ".*" | while read dir; do
    skill_name=$(basename "$dir")
    is_expected=0
    for expected in "${EXPECTED_SKILLS[@]}"; do
        if [ "$skill_name" = "$expected" ]; then
            is_expected=1
            break
        fi
    done
    if [ $is_expected -eq 0 ]; then
        echo "$skill_name"
    fi
done)

if [ -z "$EXTRA_SKILLS" ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${YELLOW}WARNING${NC}"
    echo "$EXTRA_SKILLS" | while read -r skill; do
        echo -e "  ${YELLOW}→${NC} Unexpected skill directory: $skill"
    done
fi

# Check for TODO or FIXME in skill files
echo -n "Checking: No TODO/FIXME markers in skills... "
TODO_RESULTS=$(grep -rn "TODO\|FIXME" skills/ --include="*.md" 2>/dev/null || true)
if [ -z "$TODO_RESULTS" ]; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${YELLOW}FOUND (review if still valid)${NC}"
    echo "$TODO_RESULTS" | while read -r line; do
        echo -e "  ${YELLOW}→${NC} $line"
    done
fi

# Summary
echo ""
echo "=========================================="
if [ $ISSUES_FOUND -eq 0 ]; then
    echo -e "${GREEN}All skill validation checks passed!${NC}"
    echo ""
    echo "Summary:"
    echo "  - All ${#EXPECTED_SKILLS[@]} expected skill directories exist"
    echo "  - All SKILL.md files have valid frontmatter"
    echo "  - All required sections are present"
    echo "  - skills/README.md references all skills"
    exit 0
else
    echo -e "${RED}Skill validation issues found - please review above${NC}"
    echo ""
    echo "Common fixes:"
    echo "  - Ensure all skill directories have SKILL.md files"
    echo "  - Add YAML frontmatter with --- delimiters"
    echo "  - Include 'name' and 'description' in frontmatter"
    echo "  - Add required sections (## When to Use This Skill, etc.)"
    echo "  - Reference all skills in skills/README.md"
    exit 1
fi
