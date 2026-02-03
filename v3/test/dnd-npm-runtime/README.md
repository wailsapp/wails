# DND NPM Runtime Test

This test verifies that drag-and-drop functionality works correctly when using the `@wailsio/runtime` npm module instead of the bundled `/wails/runtime.js`.

## Background

There was a bug where the Go backend called `window.wails.Window.HandlePlatformFileDrop()` for native file drops (macOS/Linux), but the npm module only registers the handler at `window._wails.handlePlatformFileDrop`.

The bundled runtime sets `window.wails = Runtime`, so the call worked. But with the npm module, `window.wails` is an empty object.

## The Fix

Changed `v3/pkg/application/webview_window.go` to call the internal path that both runtimes set up:

```go
// Before (only worked with bundled runtime):
"window.wails.Window.HandlePlatformFileDrop(%s, %d, %d);"

// After (works with both):
"window._wails.handlePlatformFileDrop(%s, %d, %d);"
```

## Running the Test

```bash
cd frontend
npm install
npm run build
cd ..
go run .
```

Then drag files from Finder/Explorer onto the drop zone. Files should be categorized and displayed.

## What This Tests

1. `@wailsio/runtime` npm module initialization
2. Event system (`Events.On('files-dropped', ...)`)
3. Native file drop handling on macOS/Linux via `window._wails.handlePlatformFileDrop`
4. Drop target detection with `data-file-drop-target` attribute
5. Visual feedback with `.file-drop-target-active` class
