#!/bin/bash
set -e

# Build configuration
APP_NAME="mac-window-tabs"
BUNDLE_ID="com.example.macwindowtabs"
VERSION="0.1.0"
BUILD_NUMBER="0.1.0"
BUILD_DIR="build/ios"
TARGET="simulator"

echo "Building iOS app: $APP_NAME"
echo "Bundle ID: $BUNDLE_ID"
echo "Version: $VERSION ($BUILD_NUMBER)"
echo "Target: $TARGET"

# Ensure build directory exists
mkdir -p "$BUILD_DIR"

# Determine SDK and target architecture
if [ "$TARGET" = "simulator" ]; then
    SDK="iphonesimulator"
    ARCH="arm64-apple-ios15.0-simulator"
elif [ "$TARGET" = "device" ]; then
    SDK="iphoneos"
    ARCH="arm64-apple-ios15.0"
else
    echo "Unknown target: $TARGET"
    exit 1
fi

# Get SDK path
SDK_PATH=$(xcrun --sdk $SDK --show-sdk-path)

# Compile the application
echo "Compiling with SDK: $SDK"
xcrun -sdk $SDK clang \
    -target $ARCH \
    -isysroot "$SDK_PATH" \
    -framework Foundation \
    -framework UIKit \
    -framework WebKit \
    -framework CoreGraphics \
    -o "$BUILD_DIR/$APP_NAME" \
    "$BUILD_DIR/main.m"

# Create app bundle
echo "Creating app bundle..."
APP_BUNDLE="$BUILD_DIR/$APP_NAME.app"
rm -rf "$APP_BUNDLE"
mkdir -p "$APP_BUNDLE"

# Move executable
mv "$BUILD_DIR/$APP_NAME" "$APP_BUNDLE/"

# Copy Info.plist
cp "$BUILD_DIR/Info.plist" "$APP_BUNDLE/"

# Sign the app
echo "Signing app..."
codesign --force --sign - "$APP_BUNDLE"

echo "Build complete: $APP_BUNDLE"

# Deploy to simulator if requested
if [ "$TARGET" = "simulator" ]; then
    echo "Deploying to simulator..."
    xcrun simctl terminate booted "$BUNDLE_ID" 2>/dev/null || true
    xcrun simctl install booted "$APP_BUNDLE"
    xcrun simctl launch booted "$BUNDLE_ID"
    echo "App launched on simulator"
fi