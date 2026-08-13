package main

import (
	"log"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Fix verification for the crash reproduced by the example of the same name in
// the unpatched checkout. Same manifest, same app id, same launch path; the
// only difference is Options.Linux.ApplicationID below.
//
// Without it, appNew builds the GtkApplication id as "org.wails." plus a
// sanitised Options.Name. WebKit derives the accessibility bus name it asks the
// portal to own from that id, a flatpak may only own names prefixed with its
// own app id, so the request is refused and the web process aborts -- which
// reaches Go as a SIGTRAP raised during the cgo call to g_application_run.
//
// With it, the two ids agree and the window opens.
//
// Run it:
//
//	wails3 task bundle
//	flatpak install --user ./flatpak/com.example.WailsFlatpakAppId.flatpak
//
// then launch "Wails Flatpak App ID" from the desktop. Launching from the menu
// matters: the web process only makes the portal call when the accessibility
// bus is reachable inside the sandbox.

// Must match app-id in flatpak/com.example.WailsFlatpakAppId.yml.
const appID = "com.example.WailsFlatpakAppId"

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("[pid %d] start, FLATPAK_ID=%q", os.Getpid(), os.Getenv("FLATPAK_ID"))

	app := application.New(application.Options{
		Name:        "linux-flatpak-appid",
		Description: "Flatpak app id matches the GtkApplication id",
		Linux: application.LinuxOptions{
			// The whole fix. Everything else here is identical to the repro.
			ApplicationID: appID,
		},
		Assets: application.AssetOptions{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`<html><body style="font-family:sans-serif;margin:3rem">
					<h1>Window opened</h1>
					<p>Reaching this means the web process was allowed to start.</p>
				</body></html>`))
			}),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "linux-flatpak-appid",
		Width:  700,
		Height: 400,
		URL:    "/",
	})

	log.Printf("[pid %d] calling app.Run()", os.Getpid())
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
