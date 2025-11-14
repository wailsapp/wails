package main

import (
	"embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

//go:embed assets
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Issue #4650 Reproduction",
		Description: "Demonstrates macOS window behavior bugs",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
			// Note: ActivationPolicyAccessory is suspected to be related to Issue #1
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	// Create system tray
	// BUG #1: On macOS, when window.Hide() is called, this tray icon disappears too
	systemTray := app.SystemTray.New()

	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
	}
	systemTray.SetLabel("Issue 4650")

	// Create window with frameless configuration
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Issue #4650 - macOS Window Bugs",
		Width:  800,
		Height: 600,
		// BUG #2: Frameless=true causes white flicker during maximize on macOS
		Frameless: true,
		// Dark background to make the white flicker very visible
		BackgroundColour: application.NewRGB(30, 30, 30),
		Mac: application.MacWindow{
			// Translucent backdrop makes the flicker more noticeable
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
			InvisibleTitleBarHeight: 50,
		},
		Windows: application.WindowsWindow{
			// Hide from taskbar to test tray-only mode
			HiddenOnTaskbar: true,
		},
	})

	// Create an App service to expose methods to frontend
	appService := &AppService{
		window: window,
	}

	// Register the service
	app.RegisterService(appService)

	// Prevent window close from quitting the app
	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		window.Hide()
		e.Cancel()
	})

	// System tray menu
	menu := app.NewMenu()
	menu.Add("Show Window").OnClick(func(ctx *application.Context) {
		window.Show()
	})

	menu.Add("Hide Window (Bug: Tray disappears on macOS)").OnClick(func(ctx *application.Context) {
		// BUG #1: On macOS, this causes the system tray icon to disappear
		// Expected: Only window hides, tray should remain visible
		window.Hide()
	})

	menu.AddSeparator()

	menu.Add("Toggle Maximize (Bug: White flicker on macOS)").OnClick(func(ctx *application.Context) {
		// BUG #2: On macOS with Frameless=true, this shows a white flash
		// Expected: Smooth transition with dark background maintained
		window.ToggleMaximise()
	})

	menu.AddSeparator()

	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	systemTray.SetMenu(menu)

	// Also show window when clicking tray icon
	systemTray.OnClick(func() {
		window.Show()
	})

	// Run the application
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// AppService provides methods callable from the frontend
type AppService struct {
	window application.Window
}

// HideWindow demonstrates Issue #1
// On macOS, this causes both the window AND tray icon to disappear
func (a *AppService) HideWindow() {
	a.window.Hide()
}

// ShowWindow shows the window
func (a *AppService) ShowWindow() {
	a.window.Show()
}

// ToggleMaximise demonstrates Issue #2
// On macOS with Frameless=true, this causes a white flash
func (a *AppService) ToggleMaximise() {
	a.window.ToggleMaximise()
}

// MinimizeWindow - for comparison, minimize works correctly
func (a *AppService) MinimizeWindow() {
	a.window.Minimise()
}

// IsMaximised returns whether the window is maximised
func (a *AppService) IsMaximised() bool {
	return a.window.IsMaximised()
}
