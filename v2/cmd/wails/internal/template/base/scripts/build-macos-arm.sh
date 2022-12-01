#! /bin/bash

echo -e "Start running the script..."
cd ../

echo -e "Start building the app for macos platform..."
wails build --clean --platform darwin/arm64

echo -e "End running the script!"
