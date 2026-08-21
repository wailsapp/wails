//go:build darwin

package main

import (
	"embed"
	"log"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

const showShortcut = "CmdOrCtrl+Shift+N"

type NotificationController struct {
	mu     sync.RWMutex
	app    *application.App
	window *application.NotchWindow
}

func (controller *NotificationController) configure(app *application.App, window *application.NotchWindow) {
	controller.mu.Lock()
	defer controller.mu.Unlock()
	controller.app = app
	controller.window = window
}

func (controller *NotificationController) currentWindow() *application.NotchWindow {
	controller.mu.RLock()
	defer controller.mu.RUnlock()
	return controller.window
}

// Hide hides the notification while keeping its webview alive for the next
// presentation.
func (controller *NotificationController) Hide() {
	if window := controller.currentWindow(); window != nil {
		window.Hide()
	}
}

// Quit terminates the accessory application instead of merely hiding its
// reusable notification window.
func (controller *NotificationController) Quit() {
	controller.mu.RLock()
	app := controller.app
	controller.mu.RUnlock()
	if app != nil {
		app.Quit()
	}
}

func main() {
	controller := &NotificationController{}
	app := application.New(application.Options{
		Name:        "Notch Notification Example",
		Description: "Demonstrates a rich, stateful NewNotchWindow notification",
		Services: []application.Service{
			application.NewService(controller),
		},
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
	})

	notification := app.Window.NewNotchWindow(application.NotchWindowOptions{
		Width:    424,
		Height:   300,
		Animated: true,
		WindowOptions: application.WebviewWindowOptions{
			Name:         "system-monitor-notification",
			URL:          "/",
			HideOnEscape: true,
		},
	})
	controller.configure(app, notification)

	if err := app.GlobalShortcut.Register(showShortcut, func() {
		notification.Show()
	}); err != nil {
		log.Printf("registering %s: %v", showShortcut, err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
