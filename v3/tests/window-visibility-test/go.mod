module window-visibility-test

go 1.25.0

replace github.com/wailsapp/wails/v3 => ../../

require github.com/wailsapp/wails/v3 v3.0.0-alpha.62

require (
	github.com/adrg/xdg v0.5.3 // indirect
	github.com/coder/websocket v1.8.14 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/godbus/dbus/v5 v5.2.2 // indirect
	github.com/jchv/go-winloader v0.0.0-20250406163304-c1995be93bd1 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/wailsapp/wails/webview2 v1.0.27 // indirect
	golang.org/x/sys v0.43.0 // indirect
)

// Add any other dependencies that might be needed
// These will be resolved when the user runs go mod tidy
