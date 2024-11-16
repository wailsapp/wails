#!/usr/bin/env bash
# Copyright (c) 2018-Present Lea Anthony
# SPDX-License-Identifier: MIT

# Fail script on any error
set -euxo pipefail

# Define variables
APP_DIR="${APP_NAME}.AppDir"

# Create AppDir structure
mkdir -p "${APP_DIR}/usr/bin"
cp -r "${APP_BINARY}" "${APP_DIR}/usr/bin/"
cp "${ICON_PATH}" "${APP_DIR}/"
cp "${DESKTOP_FILE}" "${APP_DIR}/"

# Download linuxdeploy and make it executable
wget -q -4 -N https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-x86_64.AppImage
chmod +x linuxdeploy-x86_64.AppImage

# Run linuxdeploy to bundle the application
./linuxdeploy-x86_64.AppImage --appdir "${APP_DIR}" --output appimage

# Rename the generated AppImage
mv "${APP_NAME}*.AppImage" "${APP_NAME}.AppImage"

