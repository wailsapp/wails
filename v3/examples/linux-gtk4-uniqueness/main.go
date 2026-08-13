package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Reproduction case: on Linux/GTK4 the second launch of an application never
// opens a window and never exits.
//
// Expected: every launch opens its own window, exactly as the same binary does
// on the GTK3 backend and on macOS and Windows. Nothing here asks for single
// instance behaviour -- Options.SingleInstance is deliberately not set.
//
// Observed: the first launch is fine. Every launch after it comes up as far as
// app.Run(), prints nothing further, draws no window, and sits there until it
// is killed. The heartbeat below keeps printing, so the process is alive and
// the Go side is healthy; it is blocked inside the cgo call to
// g_application_run.
//
// Cause: appNew in pkg/application/linux_cgo.go builds the GtkApplication with
// APPLICATION_DEFAULT_FLAGS, which linux_cgo.h defines as
// G_APPLICATION_DEFAULT_FLAGS, so the application is unique. GTK derives the
// bus name from the application id, which is "org.wails." plus a sanitised
// Options.Name. When a second process finds that name owned, g_application_run
// registers as a remote instance: it forwards "activate" to the primary and
// returns without ever emitting it locally. linuxApp.waitForActivation blocks
// on that signal before any window is built, so the process wedges.
//
// The GTK3 backend passes G_APPLICATION_NON_UNIQUE and has no activation gate,
// which is why this only appears once an app moves to GTK4.
//
// Run it:
//
//	go run .            # leave it running, window appears
//	go run .            # second process: no window, never returns
//
// Note the second process keeps logging "still waiting" while showing nothing.

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	logf := func(format string, args ...any) {
		log.Printf("[pid %d] "+format, append([]any{os.Getpid()}, args...)...)
	}

	logf("start")

	app := application.New(application.Options{
		Name:        "linux-gtk4-uniqueness",
		Description: "Second launch never opens a window on GTK4",
		Assets: application.AssetOptions{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<html><body style="font-family:sans-serif;margin:3rem">
					<h1>Window opened</h1>
					<p>If you can read this, this process got past g_application_run.</p>
				</body></html>`))
			}),
		},
	})

	// Never fires in the wedged process, because the activate signal it is
	// waiting for was delivered to the first instance instead.
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		logf("ApplicationStarted fired")
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "linux-gtk4-uniqueness",
		Width:  700,
		Height: 400,
		URL:    "/",
	})

	// Proof that the process is alive and only the GTK side is stuck.
	go func() {
		for range time.Tick(2 * time.Second) {
			logf("still waiting -- no window, no ApplicationStarted")
		}
	}()

	logf("calling app.Run()")
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
	logf("app.Run() returned")
}
