package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sync"
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
	// arrive in Go as common: application events (mapped from the per-platform
	// ios:/android: events), with their payload on the event context. They are
	// Go-only, so the app forwards them to the frontend as custom events here.
	// forward re-emits a Go application event to the frontend as a "sys:*"
	// custom event, trailing-debounced: the OS fires theme/network callbacks in
	// rapid bursts that often interleave transient values (e.g. one dark-mode
	// toggle yields several ThemeChanged events; disabling wifi flaps
	// wifi→none→cellular). We forward only the value that the burst settles on,
	// once it has been quiet for debounce, and skip it if it matches the last
	// value we forwarded.
	const debounce = 300 * time.Millisecond
	forward := func(jsName string) func(*application.ApplicationEvent) {
		var (
			mu     sync.Mutex
			timer  *time.Timer
			latest map[string]any
			last   string
		)
		return func(e *application.ApplicationEvent) {
			mu.Lock()
			latest = e.Context().Data()
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(debounce, func() {
				mu.Lock()
				data := latest
				key, _ := json.Marshal(data)
				dup := string(key) == last
				last = string(key)
				mu.Unlock()
				if dup {
					return
				}
				app.Logger.Info("system event", "event", jsName, "data", data)
				app.Event.Emit(jsName, data)
			})
			mu.Unlock()
		}
	}
	app.Event.OnApplicationEvent(events.Common.NetworkChanged, forward("sys:network"))
	app.Event.OnApplicationEvent(events.Common.ThemeChanged, forward("sys:theme"))
	// Low memory is a pulse, not a state, so it is not de-duplicated.
	app.Event.OnApplicationEvent(events.Common.LowMemory, func(e *application.ApplicationEvent) {
		app.Logger.Info("system event", "event", "sys:memory")
		app.Event.Emit("sys:memory", map[string]any{})
	})
	app.Event.OnApplicationEvent(events.Common.ScreenLocked, func(e *application.ApplicationEvent) {
		app.Event.Emit("sys:lock", map[string]any{"locked": true})
	})
	app.Event.OnApplicationEvent(events.Common.ScreenUnlocked, func(e *application.ApplicationEvent) {
		app.Event.Emit("sys:lock", map[string]any{"locked": false})
	})

	// Battery is reported by the OS far more often than is useful (on Android,
	// ACTION_BATTERY_CHANGED also fires on temperature/voltage changes). Throttle
	// what we forward: every 10% while above 10%, then every 1% from 10% down,
	// plus immediately on any charge-state / low-power change.
	var (
		batteryMu       sync.Mutex
		lastBatteryPct  = -1
		lastBatteryMeta = ""
	)
	app.Event.OnApplicationEvent(events.Common.BatteryChanged, func(e *application.ApplicationEvent) {
		data := e.Context().Data()
		pct := -1
		if lv, ok := data["level"].(float64); ok && lv >= 0 {
			pct = int(math.Round(lv * 100))
		}
		state, _ := data["state"].(string)
		low, _ := data["lowPowerMode"].(bool)
		meta := fmt.Sprintf("%s|%t", state, low)

		batteryMu.Lock()
		report := lastBatteryPct < 0 || meta != lastBatteryMeta
		if !report && pct >= 0 {
			if pct <= 10 || lastBatteryPct <= 10 {
				report = pct != lastBatteryPct // 1% steps once at/below 10%
			} else {
				delta := pct - lastBatteryPct
				if delta < 0 {
					delta = -delta
				}
				report = delta >= 10 // 10% steps above 10%
			}
		}
		if report {
			lastBatteryPct = pct
			lastBatteryMeta = meta
		}
		batteryMu.Unlock()

		if report {
			app.Logger.Info("system event", "event", "sys:battery", "pct", pct, "data", data)
			app.Event.Emit("sys:battery", data)
		}
	})

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
