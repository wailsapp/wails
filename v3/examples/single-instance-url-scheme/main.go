package main

import (
	"embed"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets/index.html
var assets embed.FS

// Reproduction case for https://github.com/wailsapp/wails/issues/5089.
//
// Expected: when a second instance is launched with a custom URL scheme on
// macOS (while a first instance is already running), the URL should reach
// the first instance — either via SecondInstanceData (e.g. in Args) or via
// the ApplicationLaunchedWithUrl event.
//
// Observed (bug): on macOS, the URL is delivered to the second process by
// LaunchServices through an Apple Event (kAEGetURL) AFTER application.New
// has already detected the single-instance lock and called os.Exit. The
// URL never appears in os.Args (macOS does not pass URL-scheme launches via
// argv), and ApplicationLaunchedWithUrl does not fire because the second
// process exits before NSApplication.run installs the Apple Event handler.
// Result: URL is lost.

const scheme = "wails-single-url"

func main() {
	logFile, err := os.OpenFile("/tmp/wails-single-instance-url.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		log.SetOutput(logFile)
		defer logFile.Close()
	}
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Printf("==== process start pid=%d argv=%v ====", os.Getpid(), os.Args)

	var window *application.WebviewWindow

	app := application.New(application.Options{
		Name:        "Single Instance URL Scheme Repro",
		LogLevel:    slog.LevelDebug,
		Description: "Reproduces wails#5089",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "com.wails.example.single-instance-url-scheme",
			OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
				url, found := findSchemeURL(data.Args)
				log.Printf("[first] OnSecondInstanceLaunch fired")
				log.Printf("[first]   Args           = %v", data.Args)
				log.Printf("[first]   WorkingDir     = %s", data.WorkingDir)
				log.Printf("[first]   AdditionalData = %v", data.AdditionalData)
				log.Printf("[first]   url-in-args?   = %v  (url=%q)", found, url)
				if window != nil {
					window.EmitEvent("secondInstance", map[string]any{
						"args":   data.Args,
						"url":    url,
						"found":  found,
						"source": "OnSecondInstanceLaunch",
					})
					window.Restore()
					window.Focus()
				}
			},
			AdditionalData: map[string]string{
				"launchtime": time.Now().Format(time.RFC3339Nano),
			},
		},
	})

	app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl,
		func(e *application.ApplicationEvent) {
			url := e.Context().URL()
			log.Printf("[first] ApplicationLaunchedWithUrl fired url=%q", url)
			if window != nil {
				window.EmitEvent("launchedWithUrl", url)
			}
		})

	window = app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Single Instance URL Scheme Repro",
		Width:  900,
		Height: 600,
		URL:    "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func findSchemeURL(args []string) (string, bool) {
	for _, a := range args {
		if strings.HasPrefix(a, scheme+"://") {
			return a, true
		}
	}
	return "", false
}

