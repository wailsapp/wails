package main

import (
	_ "embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

// Manual reproduction harness for #4424.
//
// The original bug surfaces as menu/systray clicks freezing while the
// backend is busy emitting events. The freeze happens because the event
// dispatch path holds windowsLock during a per-window ExecJS call, which
// blocks any goroutine that needs windowsLock for write — including the
// main thread when it processes the menu activation.
//
// To reproduce, EventIPCTransport.DispatchWailsEvent must actually reach
// the per-window ExecJS path; that requires (a) at least one webview
// window registered on the App and (b) the event-emit goroutine running
// frequently enough that a menu click is likely to land inside one of
// the per-window dispatches.
//
// Manual procedure:
//  1. Run the binary.
//  2. A visible window opens showing the incrementing counter (proves
//     events round-trip through the per-window dispatch path).
//  3. Right-click the systray icon and click each menu item repeatedly.
//  4. With the fix in place, menu items respond instantly even while
//     the counter is ticking. Without the fix, clicks queue up or are
//     ignored until the event burst stops.
//
//go:embed assets/index.html
var indexHTML []byte

func main() {
	app := application.New(application.Options{
		Name:        "Event Block Test",
		Description: "Manual reproduction for #4424 — event emitter blocks menu clicks on Linux",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	// A visible window is required: EventIPCTransport.DispatchWailsEvent
	// iterates over app.windows and calls ExecJS on each. With no window
	// registered, the dispatch loop is a no-op and the bug cannot reproduce.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Event Block Test (#4424)",
		Width:  480,
		Height: 240,
		HTML:   string(indexHTML),
	})

	systemTray := app.SystemTray.New()
	systemTray.SetIcon(icons.WailsLogoBlack)

	menu := app.NewMenu()
	menu.Add("About").OnClick(func(ctx *application.Context) {
		log.Println("About clicked!")
		app.Dialog.Info().
			SetTitle("About").
			SetMessage("Manual repro harness for issue #4424").
			Show()
	})
	menu.Add("Show Time").OnClick(func(ctx *application.Context) {
		log.Printf("Time clicked! Current time: %s\n", time.Now().Format(time.RFC3339))
	})
	menu.AddSeparator()
	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		log.Println("Quit clicked!")
		app.Quit()
	})
	systemTray.SetMenu(menu)

	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		go func() {
			log.Println("Starting event emitter at 100ms intervals.")
			log.Println("TEST: right-click the systray icon and click menu items repeatedly.")
			log.Println("PASS: menu items respond immediately while the counter ticks.")
			log.Println("FAIL: menu clicks stall, queue up, or are dropped until the burst stops.")

			counter := 0
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()
			for range ticker.C {
				counter++
				app.Event.Emit("backendmessage", map[string]interface{}{
					"counter": counter,
					"time":    time.Now().Format(time.RFC3339Nano),
				})
			}
		}()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
