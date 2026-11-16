module window-visibility-test

go 1.26.0

replace github.com/wailsapp/wails/v3 => ../../

require github.com/wailsapp/wails/v3 v3.0.0-alpha.62

require (
	github.com/coder/websocket v1.8.14 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/godbus/dbus/v5 v5.2.2 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/term v0.46.0 // indirect
)

// Add any other dependencies that might be needed
// These will be resolved when the user runs go mod tidy
