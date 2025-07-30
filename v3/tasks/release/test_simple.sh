#!/bin/bash

echo "Testing release.go --create-release-notes functionality"
echo "======================================================"

# Save current directory
ORIGINAL_DIR=$(pwd)

# Go to v3 root (where UNRELEASED_CHANGELOG.md should be)
cd ../..

# Create a test UNRELEASED_CHANGELOG.md
cat > UNRELEASED_CHANGELOG.md << 'EOF'
# Unreleased Changes

<!-- 
This file is used to collect changelog entries for the next v3-alpha release.
-->

## Added
<!-- New features, capabilities, or enhancements -->
- Add Windows dark theme support for menus and menubar
- Add `--create-release-notes` flag to release script

## Changed
<!-- Changes in existing functionality -->
- Update Go version to 1.23 in workflow
- Improve error handling in release process

## Fixed
<!-- Bug fixes -->
- Fix nightly release workflow changelog extraction
- Fix Go cache configuration in GitHub Actions

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->

---

### Example Entries:

**Added:**
- Example content
EOF

echo ""
echo "Test 1: Running with valid content"
echo "-----------------------------------"

# Run the release script
cd tasks/release
if go run release.go --create-release-notes; then
    echo "✅ Command succeeded"
    
    # Check if release_notes.md was created
    if [ -f "../../release_notes.md" ]; then
        echo "✅ release_notes.md was created"
        echo ""
        echo "Content:"
        echo "--------"
        cat ../../release_notes.md
        echo ""
        echo "--------"
    else
        echo "❌ release_notes.md was NOT created"
    fi
else
    echo "❌ Command failed"
fi

echo ""
echo "Test 2: Check --check-only flag"
echo "--------------------------------"

# Test the check-only flag
if go run release.go --check-only; then
    echo "✅ --check-only detected content"
else
    echo "❌ --check-only did not detect content"
fi

echo ""
echo "Test 3: Check --extract-changelog flag"
echo "--------------------------------------"

# Test the extract-changelog flag
OUTPUT=$(go run release.go --extract-changelog 2>&1)
if [ $? -eq 0 ]; then
    echo "✅ --extract-changelog succeeded"
    echo "Output:"
    echo "-------"
    echo "$OUTPUT"
    echo "-------"
else
    echo "❌ --extract-changelog failed"
    echo "Error: $OUTPUT"
fi

# Clean up
cd ../..
rm -f release_notes.md

# Restore original UNRELEASED_CHANGELOG.md if it exists
if [ -f "UNRELEASED_CHANGELOG.md.backup" ]; then
    mv UNRELEASED_CHANGELOG.md.backup UNRELEASED_CHANGELOG.md
fi

cd "$ORIGINAL_DIR"

echo ""
echo "======================================================"
echo "Testing complete!"