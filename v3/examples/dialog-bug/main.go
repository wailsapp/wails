package main

import (
	"fmt"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
	app := application.New(application.Options{
		Name:        "Dialog Bug",
		Description: "Minimal repro for dialog behaviour",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Dialog Bug",
		Width:  900,
		Height: 600,
		HTML: `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Dialog Bug</title>
  <style>
    body { font-family: system-ui, -apple-system, Segoe UI, Roboto, sans-serif; margin: 24px; }
    code { background: #f4f4f5; padding: 2px 6px; border-radius: 6px; }
  </style>
</head>
<body>
  <h1>Dialog Bug</h1>
  <p>This window exists only to attach native dialogs to it.</p>
  <p>Check your terminal logs after clicking a button.</p>
  <p><code>MessageDialog</code> should not quit the app unless your code calls <code>app.Quit()</code>.</p>
</body>
</html>`,
	})

	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		dialog := app.Dialog.Question().
			SetTitle("Migration Warning").
			SetMessage("2 migrations are available. Do you want to run them now?").
			AttachToWindow(window).
			AddDefaultButton("Run").
			AddCancelButton("Cancel")

		result, err := dialog.Result()
		if err != nil {
			log.Printf("dialog error: %v\n", err)
			return
		}

		log.Printf("dialog result: %s\n", result)
		_ = app.Dialog.Info().
			AttachToWindow(window).
			SetTitle("Result").
			SetMessage(fmt.Sprintf("You clicked: %s", result)).
			Show()
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
