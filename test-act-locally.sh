#!/bin/bash

echo "🧪 LOCAL ACT TESTING SIMULATION"
echo "================================="
echo ""
echo "This script simulates what 'act' would do to test the nightly release workflow"
echo ""

# Simulate the key steps from the GitHub Actions workflow
echo "1. Checking current directory and git status..."
pwd
git status --porcelain

echo ""
echo "2. Testing Go version..."
go version

echo ""
echo "3. Simulating workflow environment setup..."
export GITHUB_WORKSPACE=$(pwd)
export GITHUB_REF="refs/heads/v3-alpha"
export RUNNER_TEMP="/tmp"

echo "GITHUB_WORKSPACE: $GITHUB_WORKSPACE"
echo "GITHUB_REF: $GITHUB_REF"

echo ""
echo "4. Testing the release script..."
cd v3/tasks/release
echo "Current directory: $(pwd)"
echo ""
echo "Running: go run release.go"
echo "========================================"
go run release.go

echo ""
echo "========================================"
echo "5. Checking generated files..."
if [ -f "release-notes.txt" ]; then
    echo "✅ release-notes.txt created"
    echo "Content preview:"
    head -10 release-notes.txt
else
    echo "❌ release-notes.txt not found"
fi

echo ""
echo "6. Test completed!"
echo "========================================"
echo ""
echo "To test with actual act (after installing it):"
echo "  cd /Users/leaanthony/GolandProjects/wails"
echo "  act workflow_dispatch -j nightly-release --input dry_run=true"
echo ""
echo "Or to test specific event:"
echo "  act workflow_dispatch -W .github/workflows/nightly-release-v3.yml"