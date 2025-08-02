#!/bin/bash

# Test script for version increment logic
set -e

echo "🧪 Testing Version Increment Logic"
echo "=================================="

# Test v2 version increment
echo "📈 Testing v2 version increment..."

# Get current v2 version
CURRENT_V2=$(cat v2/cmd/wails/internal/version.txt | sed 's/^v//')
echo "Current v2 version: v$CURRENT_V2"

# Parse version parts
IFS='.' read -ra VERSION_PARTS <<< "$CURRENT_V2"
MAJOR=${VERSION_PARTS[0]}
MINOR=${VERSION_PARTS[1]}
PATCH=${VERSION_PARTS[2]}

echo "Parsed: MAJOR=$MAJOR, MINOR=$MINOR, PATCH=$PATCH"

# Test patch increment
PATCH_VERSION="$MAJOR.$MINOR.$((PATCH + 1))"
echo "✅ Patch increment: v$CURRENT_V2 → v$PATCH_VERSION"

# Test minor increment
MINOR_VERSION="$MAJOR.$((MINOR + 1)).0"
echo "✅ Minor increment: v$CURRENT_V2 → v$MINOR_VERSION"

# Test major increment
MAJOR_VERSION="$((MAJOR + 1)).0.0"
echo "✅ Major increment: v$CURRENT_V2 → v$MAJOR_VERSION"

echo ""

# Test v3 version increment (simulate)
echo "📈 Testing v3 version increment..."

# Simulate current v3 version
CURRENT_V3="v3.0.0-alpha.9"
echo "Simulated current v3 version: $CURRENT_V3"

if [[ $CURRENT_V3 =~ v3\.0\.0-alpha\.([0-9]+) ]]; then
    ALPHA_NUM=${BASH_REMATCH[1]}
    NEW_ALPHA_NUM=$((ALPHA_NUM + 1))
    NEW_V3_VERSION="v3.0.0-alpha.$NEW_ALPHA_NUM"
    echo "✅ Alpha increment: $CURRENT_V3 → $NEW_V3_VERSION"
else
    echo "❌ Failed to parse v3 version format"
    exit 1
fi

echo ""

# Test conventional commit detection
echo "🔍 Testing Conventional Commit Detection..."

# Simulate commit messages
COMMITS="
feat: add new dialog API
fix: resolve memory leak
chore: update dependencies
feat!: remove deprecated API
docs: update README
BREAKING CHANGE: remove v1 compatibility
"

echo "Test commits:"
echo "$COMMITS"

# Test release type detection
if echo "$COMMITS" | grep -q "feat!\|fix!\|BREAKING CHANGE:"; then
    RELEASE_TYPE="major"
elif echo "$COMMITS" | grep -q "feat\|BREAKING CHANGE"; then
    RELEASE_TYPE="minor"
else
    RELEASE_TYPE="patch"
fi

echo "✅ Detected release type: $RELEASE_TYPE"

echo ""
echo "✅ Version logic test completed!"