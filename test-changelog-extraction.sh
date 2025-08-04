#!/bin/bash

# Test script for changelog extraction logic
set -e

echo "🧪 Testing Changelog Extraction Logic"
echo "======================================"

# Test v2 changelog extraction
echo "📋 Testing v2 changelog extraction..."
CHANGELOG_FILE="website/src/pages/changelog.mdx"

if [ ! -f "$CHANGELOG_FILE" ]; then
    echo "❌ v2 changelog file not found: $CHANGELOG_FILE"
    exit 1
fi

echo "✅ v2 changelog file found"

# Extract unreleased section
awk '
/^## \[Unreleased\]/ { found=1; next }
found && /^## / { exit }
found && !/^$/ { print }
' $CHANGELOG_FILE > test_v2_notes.md

echo "📝 v2 extracted content:"
echo "------------------------"
if [ -s test_v2_notes.md ]; then
    head -10 test_v2_notes.md
    echo "..."
    echo "✅ v2 changelog extraction successful ($(wc -l < test_v2_notes.md) lines)"
else
    echo "⚠️  v2 unreleased section is empty"
fi

echo ""

# Test v3 changelog extraction (when on v3-alpha branch)
echo "📋 Testing v3 changelog extraction..."

# Check if we can access v3 changelog
if git show v3-alpha:docs/src/content/docs/changelog.mdx > /dev/null 2>&1; then
    echo "✅ v3 changelog accessible from v3-alpha branch"
    
    # Extract from v3-alpha branch
    git show v3-alpha:docs/src/content/docs/changelog.mdx | awk '
    /^## \[Unreleased\]/ { found=1; next }
    found && /^## / { exit }
    found && !/^$/ { print }
    ' > test_v3_notes.md
    
    echo "📝 v3 extracted content:"
    echo "------------------------"
    if [ -s test_v3_notes.md ]; then
        head -10 test_v3_notes.md
        echo "..."
        echo "✅ v3 changelog extraction successful ($(wc -l < test_v3_notes.md) lines)"
    else
        echo "⚠️  v3 unreleased section is empty"
    fi
else
    echo "⚠️  v3 changelog not accessible (expected if not on v3-alpha branch)"
fi

echo ""
echo "🧹 Cleaning up test files..."
rm -f test_v2_notes.md test_v3_notes.md

echo "✅ Changelog extraction test completed!"