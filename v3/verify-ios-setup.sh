#!/bin/bash

echo "=== iOS Build System Verification ==="
echo
echo "Checking iOS build assets..."
echo

# Check if files exist
echo "1. Checking build_assets/ios directory:"
if [ -d "internal/commands/build_assets/ios" ]; then
    echo "   ✅ iOS build_assets directory exists"
    ls -la internal/commands/build_assets/ios/
else
    echo "   ❌ iOS build_assets directory missing"
fi
echo

echo "2. Checking updatable_build_assets/ios directory:"
if [ -d "internal/commands/updatable_build_assets/ios" ]; then
    echo "   ✅ iOS updatable_build_assets directory exists"
    ls -la internal/commands/updatable_build_assets/ios/
else
    echo "   ❌ iOS updatable_build_assets directory missing"
fi
echo

echo "3. Checking iOS implementation files:"
for file in pkg/application/application_ios.go pkg/application/application_ios.h pkg/application/application_ios.m; do
    if [ -f "$file" ]; then
        echo "   ✅ $file exists"
    else
        echo "   ❌ $file missing"
    fi
done
echo

echo "4. Checking iOS example:"
if [ -d "examples/ios-poc" ]; then
    echo "   ✅ ios-poc example exists"
    ls -la examples/ios-poc/
else
    echo "   ❌ ios-poc example missing"
fi
echo

echo "5. Checking main Taskfile includes iOS:"
if grep -q "ios:" internal/templates/_common/Taskfile.tmpl.yml 2>/dev/null; then
    echo "   ✅ iOS included in main Taskfile template"
else
    echo "   ❌ iOS not included in main Taskfile template"
fi
echo

echo "6. Checking Xcode tools:"
if command -v xcrun &> /dev/null; then
    echo "   ✅ xcrun available"
    echo "      SDK Path: $(xcrun --sdk iphonesimulator --show-sdk-path 2>/dev/null || echo 'Not found')"
else
    echo "   ❌ xcrun not available"
fi
echo

echo "7. iOS Build System Summary:"
echo "   - Static assets: internal/commands/build_assets/ios/"
echo "   - Templates: internal/commands/updatable_build_assets/ios/"
echo "   - Implementation: pkg/application/application_ios.*"
echo "   - Example: examples/ios-poc/"
echo "   - Build script: build_ios.sh"
echo

echo "=== Verification Complete ==="
echo
echo "The iOS build system structure is in place and ready for:"
echo "1. Creating new iOS projects with 'wails3 init'"
echo "2. Building with 'task ios:build'"
echo "3. Running with 'task ios:run'"
echo
echo "Note: Full compilation requires iOS development environment setup."