// Package main demonstrates the hide/show window crash issue on macOS
// GitHub Issue: https://github.com/wailsapp/wails/issues/4389
//
// This test reproduces the crash that occurs when:
// 1. A window is shown via system tray click
// 2. The window is hidden via system tray click
// 3. A second show attempt crashes because the window was destroyed during hide
//
// The root cause is that when ApplicationShouldTerminateAfterLastWindowClosed is true
// (the default), hiding the last visible window with orderOut:nil triggers the
// window close event sequence, destroying the window.
package main

import (
	_ "embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

//go:embed icon.png
var iconData []byte

func main() {
	app := application.New(application.Options{
		Name:        "Systray Hide/Show Test",
		Description: "Test for macOS hide/show window crash (Issue #4389)",
		Mac: application.MacOptions{
			// NOTE: Setting this to false is the workaround for the issue
			// ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	// Create the main window
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Hide/Show Test Window",
		Width:  600,
		Height: 400,
		HTML: `<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            padding: 40px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            margin: 0;
            height: 100vh;
            box-sizing: border-box;
        }
        h1 { margin-top: 0; }
        .info {
            background: rgba(255,255,255,0.1);
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
        }
        code {
            background: rgba(0,0,0,0.2);
            padding: 2px 6px;
            border-radius: 4px;
        }
    </style>
</head>
<body>
    <h1>🧪 Hide/Show Test Window</h1>
    <div class="info">
        <h3>Test Instructions:</h3>
        <ol>
            <li>Click the system tray icon to hide this window</li>
            <li>Click the system tray icon again to show this window</li>
            <li>Repeat step 1 - the app should crash on macOS</li>
        </ol>
    </div>
    <div class="info">
        <h3>Expected Behavior:</h3>
        <p>The window should toggle visibility without crashing.</p>
    </div>
    <div class="info">
        <h3>Root Cause:</h3>
        <p>When <code>ApplicationShouldTerminateAfterLastWindowClosed</code> is true (default),
        hiding the last window triggers the close event, destroying the window.</p>
    </div>
</body>
</html>`,
	})

	// Create system tray
	systray := app.SystemTray.New()

	// Set appropriate icon for platform
	if runtime.GOOS == "darwin" {
		systray.SetTemplateIcon(icons.SystrayMacTemplate)
	} else {
		systray.SetIcon(iconData)
	}

	// Toggle window visibility on tray click
	systray.OnClick(func() {
		log.Println("[DEBUG] System tray clicked")
		if window.IsVisible() {
			log.Printf("[DEBUG] Window is visible, hiding it (windowId=%d)", window.ID())
			window.Hide()
			log.Println("[DEBUG] Window hide completed")
		} else {
			log.Printf("[DEBUG] Window is hidden, showing it (windowId=%d)", window.ID())
			window.Show()
			log.Println("[DEBUG] Window show completed")
		}
	})

	// Create the menu
	menu := app.NewMenu()
	menu.Add("Show/Hide Window").OnClick(func(ctx *application.Context) {
		if window.IsVisible() {
			window.Hide()
		} else {
			window.Show()
		}
	})
	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})
	systray.SetMenu(menu)

	log.Println("Starting app - click the system tray icon to toggle window visibility")

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
