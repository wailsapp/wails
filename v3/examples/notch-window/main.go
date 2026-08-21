//go:build darwin

package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func main() {
	app := application.New(application.Options{
		Name:        "Notch Window Example",
		Description: "Demonstrates the high-level NewNotchWindow API",
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
		Assets: application.AssetOptions{
			Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "text/html; charset=utf-8")
				_, _ = fmt.Fprintf(writer, page, html.EscapeString(request.URL.Query().Get("label")))
			}),
		},
	})

	app.Window.NewNotchWindow(application.NotchWindowOptions{
		Width:    660,
		Height:   92,
		Animated: true,
		WindowOptions: application.WebviewWindowOptions{
			Name:         "primary-notch-window",
			URL:          "/?label=Primary window",
			HideOnEscape: true,
		},
	})

	secondary := app.Window.NewNotchWindow(application.NotchWindowOptions{
		Width:          540,
		Height:         72,
		Animated:       true,
		AnimationSpeed: 600 * time.Millisecond,
		WindowOptions: application.WebviewWindowOptions{
			Name:         "secondary-notch-window",
			URL:          "/?label=Independent overlay",
			Hidden:       true,
			HideOnEscape: true,
		},
	})

	// Show a second independently managed window so the example makes
	// creation-order overlaying visible without any frontend bindings.
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		time.AfterFunc(1500*time.Millisecond, func() {
			log.Print("showing the independent overlay")
			secondary.Show()
		})
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

const page = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
      * { box-sizing: border-box; }
      html, body { width: 100%%; height: 100%%; margin: 0; background: transparent; }
      body {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 20px;
        padding: 14px 18px;
        color: #f5f5f7;
        font-family: -apple-system, BlinkMacSystemFont, "SF Pro Display", sans-serif;
      }
      .identity { min-width: 0; }
      .eyebrow { color: #ff7a1a; font-size: 11px; font-weight: 750; letter-spacing: .12em; text-transform: uppercase; }
      h1 { margin: 6px 0 0; overflow: hidden; font-size: 19px; text-overflow: ellipsis; white-space: nowrap; }
      .hint { flex: 0 0 auto; color: #929298; font-size: 12px; }
      kbd { padding: 3px 7px; border: 1px solid #353539; border-radius: 6px; background: #1d1d20; color: #d8d8dc; }
    </style>
  </head>
  <body>
    <div class="identity">
      <div class="eyebrow">Wails</div>
      <h1>%s</h1>
    </div>
    <div class="hint"><kbd>Esc</kbd> to hide</div>
  </body>
</html>`
