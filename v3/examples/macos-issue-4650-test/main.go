package main

import (
	_ "embed"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

//go:embed index.html
var indexHTML string

func main() {
	app := application.New(application.Options{
		Name:        "macOS Issue #4650 Test",
		Description: "Test application for macOS window behavior fixes",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			// Using ActivationPolicyAccessory to trigger issue #1
			// (tray icon disappearing when window is hidden)
			ActivationPolicy: application.ActivationPolicyAccessory,
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	// Create system tray
	systemTray := app.SystemTray.New()
	systemTray.SetLabel("Test")

	if runtime.GOOS == "darwin" {
		systemTray.SetTemplateIcon(icons.SystrayMacTemplate)
	}

	// Create window with dark background to test issue #2
	// (white flicker when maximizing)
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "macOS Issue #4650 Test",
		Width:  800,
		Height: 600,
		HTML:   indexHTML,
		// Dark background to make white flicker visible
		BackgroundColour: application.NewRGBA(30, 30, 30, 255),
		Mac: application.MacWindowOptions{
			Backdrop: application.MacBackdropNormal,
		},
	})

	// Prevent window close, just hide instead
	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		window.Hide()
		e.Cancel()
	})

	// Bind methods for UI controls
	app.NewService(&TestService{
		window: window,
		app:    app,
	})

	// System tray menu
	menu := app.NewMenu()
	menu.Add("Show Window").OnClick(func(ctx *application.Context) {
		window.Show()
		window.Focus()
	})
	menu.Add("Hide Window").OnClick(func(ctx *application.Context) {
		window.Hide()
	})
	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})
	systemTray.SetMenu(menu)

	// Click tray icon to show/hide window
	systemTray.OnClick(func() {
		if window.IsVisible() {
			window.Hide()
		} else {
			window.Show()
			window.Focus()
		}
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

type TestService struct {
	window application.Window
	app    *application.App
}

func (s *TestService) HideWindow() {
	s.window.Hide()
}

func (s *TestService) ShowWindow() {
	s.window.Show()
	s.window.Focus()
}

func (s *TestService) MaximizeWindow() {
	s.window.Maximize()
}

func (s *TestService) RestoreWindow() {
	s.window.Restore()
}

func (s *TestService) IsMaximized() bool {
	return s.window.IsMaximised()
}

func (s *TestService) IsVisible() bool {
	return s.window.IsVisible()
}
