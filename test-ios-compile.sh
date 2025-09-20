#!/bin/bash
set -e

echo "=== Testing iOS Build System ==="
echo

# Step 1: Install wails3
echo "Step 1: Installing wails3..."
cd /Users/leaanthony/test/wails/v3
go install ./cmd/wails3
echo "✓ wails3 installed"
echo

# Step 2: Create test project
echo "Step 2: Creating test project..."
TEST_DIR="./tmp-ios-test"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Initialize project
wails3 init -n TestIOSApp -t vanilla
cd TestIOSApp
echo "✓ Project created"
echo

# Step 3: Check iOS files
echo "Step 3: Checking iOS build files..."
if [ -d "build/ios" ]; then
    echo "iOS build directory exists:"
    ls -la build/ios/
else
    echo "ERROR: iOS build directory not found!"
    exit 1
fi
echo

# Step 4: Build for iOS
echo "Step 4: Building iOS app..."
wails3 task ios:build
echo "✓ iOS build completed"
echo

# Step 5: Check build output
echo "Step 5: Checking build output..."
if [ -f "bin/TestIOSApp" ]; then
    echo "✓ Binary created: bin/TestIOSApp"
    ls -la bin/
else
    echo "ERROR: Binary not created!"
    exit 1
fi
echo

# Step 6: Attempt to run on simulator
echo "Step 6: Attempting to run on simulator..."
wails3 task ios:run
echo

echo "=== iOS Build Test Complete ==="