package main

import (
	"embed"
	"log"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:assets
var assets embed.FS

// This example demonstrates global (system-wide) shortcuts. Unlike menu
// accelerators or key bindings, a global shortcut fires even when the
// application does not have focus.
//
// Manual test:
//  1. Run the example. The window lists the registered shortcuts.
//  2. Press a shortcut and watch the log update.
//  3. Switch focus to another application (or hide the window with its
//     shortcut) and press the shortcuts again. They still fire.
func main() {
	app := application.New(application.Options{
		Name:        "Global Shortcuts Demo",
		Description: "A demo of the global shortcuts API",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		// Note: we deliberately do NOT set
		// ApplicationShouldTerminateAfterLastWindowClosed here. On macOS, Hide()
		// uses orderOut: which makes the window non-visible, and AppKit treats
		// the last non-visible window as "closed" and would terminate the app.
		// That defeats the hide-then-resummon flow this example demonstrates.
	})

	menu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	}

	window := app.Window.New()

	// emit forwards a fired shortcut to the frontend log.
	emit := func(name, accelerator string) {
		log.Printf("global shortcut fired: %s (%s)", name, accelerator)
		app.Event.Emit("shortcut", map[string]string{"name": name, "accelerator": accelerator})
	}

	const (
		showShortcut = "CmdOrCtrl+Shift+G"
		hideShortcut = "CmdOrCtrl+Shift+H"
		pingShortcut = "CmdOrCtrl+Shift+K"
	)

	register := func(accelerator string, fn func()) {
		if err := app.GlobalShortcut.Register(accelerator, fn); err != nil {
			log.Printf("failed to register %q: %v", accelerator, err)
		}
	}

	// Bring the window to the front from anywhere.
	register(showShortcut, func() {
		emit("show", showShortcut)
		window.Show()
		window.Focus()
	})

	// Hide the window from anywhere.
	register(hideShortcut, func() {
		emit("hide", hideShortcut)
		window.Hide()
	})

	// Fire without changing the window, to show the shortcut works regardless
	// of focus.
	register(pingShortcut, func() {
		emit("ping", pingShortcut)
	})

	// Registering the same shortcut twice within the same application returns an
	// error and leaves the original binding untouched.
	if err := app.GlobalShortcut.Register(showShortcut, func() {}); err != nil {
		log.Printf("re-registering %q correctly rejected: %v", showShortcut, err)
	}

	// Tell the frontend which accelerators were registered (they differ per
	// platform because of CmdOrCtrl). Wait for the window to be ready first.
	go func() {
		time.Sleep(500 * time.Millisecond)
		app.Event.Emit("shortcuts:registered", map[string]string{
			"show": showShortcut,
			"hide": hideShortcut,
			"ping": pingShortcut,
		})
	}()

	log.Printf("registered global shortcuts: %v", app.GlobalShortcut.GetAll())

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
