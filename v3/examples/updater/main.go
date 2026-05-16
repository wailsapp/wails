// Updater example: a minimal Wails v3 app that ships in-app updates from a
// GitHub release. Three pieces wire together:
//
//  1. app.Updater.Init configures the source + current version.
//  2. A menu item triggers app.Updater.CheckAndInstall, which opens the
//     default update window and walks the full flow.
//  3. The frontend subscribes to updater:* events to show its own progress
//     UI alongside the framework window.
//
// Out of the box, the example reports as v1.0.0 and points at
// wailsapp/updater-demo — a public repo that publishes a v2.0.0 release
// for each supported platform, so "Check for Updates…" finds, downloads,
// verifies, swaps, and restarts into a real upgraded binary.
//
// Pass APP_VERSION / GH_REPOSITORY env vars to point this at your own repo.
package main

import (
	"context"
	"embed"
	"log"
	"os"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/updater"
	"github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

//go:embed assets
var assets embed.FS

// publicKey is the trust root used to verify release signatures. Set this
// to your project's Ed25519 / Ed25519ph / ECDSA-P256 public key (PEM or
// raw). When the public key is empty and a provider's release ships a
// signature, the install fails closed — by design.
var publicKey []byte

func main() {
	version := envOr("APP_VERSION", "1.0.0")
	repo := envOr("GH_REPOSITORY", "wailsapp/updater-demo")

	app := application.New(application.Options{
		Name:        "Updater Example",
		Description: "Demonstrates app.Updater wired against GitHub Releases.",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Build a GitHub provider. Add a Token if your repo is private or you
	// need to bump rate limits. The demo repo ships a SHA256SUMS sidecar
	// for digest-based verification — the framework refuses to install an
	// artifact whose hash doesn't match.
	gh, err := github.New(github.Config{
		Repository:    repo,
		ChecksumAsset: "SHA256SUMS",
	})
	if err != nil {
		log.Fatalf("github.New: %v", err)
	}

	if err := app.Updater.Init(updater.Config{
		CurrentVersion: version,
		Providers:      []updater.Provider{gh},
		PublicKey:      publicKey,
	}); err != nil {
		log.Fatalf("updater.Init: %v", err)
	}

	// Main window.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Updater Example",
		Width:  640,
		Height: 480,
		URL:    "/",
	})

	// App menu with the canonical "Check for Updates…" item.
	menu := app.Menu.New()
	app.Menu.SetApplicationMenu(menu)
	appMenu := menu.AddSubmenu("App")
	appMenu.Add("Check for Updates…").OnClick(func(*application.Context) {
		go func() {
			if err := app.Updater.CheckAndInstall(context.Background()); err != nil {
				app.Logger.Error("update flow", "error", err)
			}
		}()
	})
	appMenu.AddSeparator()
	appMenu.Add("Quit").OnClick(func(*application.Context) { app.Quit() })

	// Optional: log every updater event server-side too.
	for _, name := range []string{
		updater.EventUpdateAvailable,
		updater.EventDownloadProgress,
		updater.EventUpdateReady,
		updater.EventError,
	} {
		evt := name
		app.Event.On(evt, func(e *application.CustomEvent) {
			app.Logger.Info("updater", "event", evt, "data", e.Data)
		})
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func envOr(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
