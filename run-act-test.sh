#!/bin/bash

echo "🧪 TESTING WAILS NIGHTLY RELEASE WITH ACT"
echo "=========================================="
echo ""

# Check if act is installed
if ! command -v act &> /dev/null; then
    echo "❌ act is not installed!"
    echo ""
    echo "Install act first:"
    echo "  brew install act"
    echo "  # or"
    echo "  curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash"
    echo ""
    exit 1
fi

echo "✅ act is installed: $(act --version)"
echo ""

# Show available workflows
echo "📋 Available workflows:"
act -l
echo ""

# Test the nightly release workflow
echo "🚀 Testing nightly release workflow with dry_run=true..."
echo "======================================================="

# Run with dry_run=true (default in our workflow)
act workflow_dispatch \
    -W .github/workflows/nightly-release-v3.yml \
    -j nightly-release \
    --input dry_run=true \
    --input force_release=false \
    -v

echo ""
echo "🎉 Act test completed!"
echo ""
echo "Check the output above for:"
echo "  ✅ 1. CHANGES DETECTED: true/false"
echo "  ✅ 2. CHANGELOG VALIDATION: PASSED"
echo "  ✅ 3. RELEASE NOTES EXTRACTED TO MEMORY"
echo "  ✅ 4. GITHUB PRERELEASE DATA with IS_PRERELEASE=true, IS_LATEST=false"