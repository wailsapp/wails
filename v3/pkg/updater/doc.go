// Package updater provides an in-app self-update facility for Wails v3
// applications.
//
// The package exposes a single Updater that is reachable as `app.Updater`.
// Configure it once via Init, then call Check / DownloadAndInstall /
// CheckAndInstall to drive the update flow.
//
// Update sources are pluggable through the Provider interface. The Updater
// owns verification, atomic writes, the binary swap and the default window;
// providers only describe how to look up and stream a release.
//
// Subscribe to lifecycle events through the standard Wails event system —
// both Go and JavaScript subscribe the same way:
//
//	app.Event.On(updater.EventDownloadProgress, func(e *application.CustomEvent) {
//	    var p updater.Progress
//	    _ = json.Unmarshal(e.JSON(), &p)
//	})
//
//	wails.Events.On("wails:updater:download-progress", (e) => { /* ... */ })
package updater
