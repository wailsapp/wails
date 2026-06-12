package main

import (
	"embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
//
//go:embed all:frontend/dist
var assets embed.FS

// main is the shared entry point for desktop, iOS and Android. On Android the
// Go code is compiled as a c-shared library, so main is invoked via
// RegisterAndroidMain (see main_android.go); on iOS it is invoked through the
// generated build overlay; on desktop it runs directly.
func main() {
	app := application.New(application.Options{
		Name:        "Wails Mobile Kitchen Sink",
		Description: "Demonstrates the Wails runtime across iOS, Android and desktop",
		Services: []application.Service{
			application.NewService(&SystemService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		// Navigation is handled by the in-page tab bar so the UX is identical
		// on every platform; native iOS tabs are intentionally left off here.
		IOS:     application.IOSOptions{},
		Android: application.AndroidOptions{},
	})

	// Register iOS runtime event handlers (Go path for WKWebView toggles).
	// Compiled to a no-op on non-iOS platforms.
	registerIOSRuntimeEventHandlers(app)

	// JS -> Go event: the frontend emits "ping" and Go replies with "pong"
	// carrying a timestamp, demonstrating bidirectional events on every
	// platform.
	app.Event.On("ping", func(e *application.CustomEvent) {
		app.Event.Emit("pong", time.Now().Format(time.RFC1123))
	})

	// Native system events (battery, network, theme, screen lock, low memory)
	// arrive as application events. The per-platform ios:/android: events are
	// mapped to common: events, so cross-platform code listens on those. These
	// are Go-side only; the payload (where present) is on the event context.
	onSystemEvent := func(label string) func(*application.ApplicationEvent) {
		return func(e *application.ApplicationEvent) {
			app.Logger.Info("system event", "event", label, "data", e.Context().Data())
		}
	}
	app.Event.OnApplicationEvent(events.Common.BatteryChanged, onSystemEvent("battery"))
	app.Event.OnApplicationEvent(events.Common.NetworkChanged, onSystemEvent("network"))
	app.Event.OnApplicationEvent(events.Common.ThemeChanged, onSystemEvent("theme"))
	app.Event.OnApplicationEvent(events.Common.ScreenLocked, onSystemEvent("screen-locked"))
	app.Event.OnApplicationEvent(events.Common.ScreenUnlocked, onSystemEvent("screen-unlocked"))
	app.Event.OnApplicationEvent(events.Common.LowMemory, onSystemEvent("low-memory"))

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Wails Mobile Kitchen Sink",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	// Go -> JS event: emit the current time once a second. The frontend
	// listens for "time" and updates a live clock.
	go func() {
		for {
			app.Event.Emit("time", time.Now().Format(time.RFC1123))
			time.Sleep(time.Second)
		}
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
