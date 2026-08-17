//go:build darwin

package main

import (
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
	app := application.New(application.Options{
		NativeOnly:  true,
		Name:        "Native Notes",
		Description: "A WebView-free Wails AppKit editor",
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	benchmark, err := loadBenchmarkConfig()
	if err != nil {
		log.Fatal(err)
	}

	editorApp, err := newNativeEditorApp(app, benchmark)
	if err != nil {
		log.Fatal(err)
	}

	tray := app.SystemTray.New().
		SetTemplateIcon(icons.SystrayMacTemplate)
	tray.SetTooltip("Native Notes")
	tray.OnClick(editorApp.show)

	trayMenu := app.NewMenu()
	trayMenu.Add("Open Native Notes").OnClick(func(*application.Context) {
		editorApp.show()
	})
	trayMenu.Add("Save").OnClick(func(*application.Context) {
		if err := editorApp.save(); err != nil {
			log.Printf("save: %v", err)
		}
	})
	trayMenu.AddSeparator()
	trayMenu.Add("Quit").OnClick(func(*application.Context) { app.Quit() })
	tray.SetMenu(trayMenu)

	installApplicationMenu(app, editorApp)
	app.Event.OnApplicationEvent(events.Mac.ApplicationDidFinishLaunching, func(*application.ApplicationEvent) {
		if benchmark.forceGC {
			// Benchmark-only diagnostic: distinguish live editor memory from
			// reclaimable Go heap pages left by loading and bridging a document.
			debug.FreeOSMemory()
		}
		if benchmark.readyFile != "" {
			if err := os.WriteFile(benchmark.readyFile, []byte("ready\n"), 0o644); err != nil {
				log.Printf("benchmark ready file: %v", err)
			}
		}
		if benchmark.autoQuit > 0 {
			time.AfterFunc(benchmark.autoQuit, app.Quit)
		}
	})
	log.Printf("editing native text files in %s", editorApp.directory)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func installApplicationMenu(app *application.App, editor *nativeEditorApp) {
	menu := app.NewMenu()
	menu.AddRole(application.AppMenu)
	file := menu.AddSubmenu("File")
	file.Add("Save").SetAccelerator("CmdOrCtrl+s").OnClick(func(*application.Context) {
		if err := editor.save(); err != nil {
			log.Printf("save: %v", err)
		}
	})
	file.AddSeparator()
	file.AddRole(application.CloseWindow)
	menu.AddRole(application.EditMenu)
	view := menu.AddSubmenu("View")
	view.Add("Toggle Sidebar").SetAccelerator("Ctrl+Cmd+s").OnClick(func(*application.Context) {
		editor.sidebarPane.Toggle()
	})
	menu.AddRole(application.WindowMenu)
	app.Menu.Set(menu)
}
