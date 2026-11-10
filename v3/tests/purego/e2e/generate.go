package main

// On Windows the probe must carry the manifest that apps built with the wails3
// build tasks carry: without Common Controls 6 the message dialogs fall back
// to MessageBox, whose Yes/No dialogs ignore Escape, and the probe would test
// a configuration apps do not ship. Run go generate before building for
// Windows; the .syso files are not committed.
//
//go:generate go run ../../../cmd/wails3 generate syso -manifest e2e.exe.manifest -icon ../../../internal/commands/build_assets/windows/icon.ico -arch amd64 -out e2e_windows_amd64.syso
//go:generate go run ../../../cmd/wails3 generate syso -manifest e2e.exe.manifest -icon ../../../internal/commands/build_assets/windows/icon.ico -arch arm64 -out e2e_windows_arm64.syso
