package main

import (
	"fmt"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
	app := application.New(application.Options{
		Name:        "Dialog Sync",
		Description: "Synchronous (blocking) message dialog example using Result()",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Dialog Sync",
		Width:  900,
		Height: 600,
		HTML: `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
	<title>Dialog Sync</title>
  <style>
    body { font-family: system-ui, -apple-system, Segoe UI, Roboto, sans-serif; margin: 24px; }
    code { background: #f4f4f5; padding: 2px 6px; border-radius: 6px; }
  </style>
</head>
<body>
	<h1>Dialog Sync</h1>
	<p>This example shows how to display a blocking message dialog using <code>Result()</code>.</p>
	<p>Check your terminal logs after clicking a button.</p>
</body>
</html>`,
	})

	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		result, err := app.Dialog.Question().
			AttachToWindow(window).
			SetTitle("Migration Warning").
			SetMessage("2 migrations are available. Do you want to run them now?").
			AddDefaultButton("Run").
			AddCancelButton("Cancel").
			Result()
		if err != nil {
			log.Printf("dialog error: %v\n", err)
			return
		}

		log.Printf("dialog result: %s\n", result)
		if result == "Run" {
			_ = app.Dialog.Info().
				AttachToWindow(window).
				SetTitle("Migrations").
				SetMessage("Running migrations...").
				Show()
			return
		}
		_ = app.Dialog.Info().
			AttachToWindow(window).
			SetTitle("Migrations").
			SetMessage(fmt.Sprintf("Cancelled (clicked: %s)", result)).
			Show()
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
