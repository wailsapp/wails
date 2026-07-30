package main

import (
	"context"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

const splashHTML = `<!doctype html>
<html><head><meta charset="utf-8"><style>
* { box-sizing: border-box; }
html, body { width: 100%; height: 100%; margin: 0; }
body { display: grid; place-items: center; color: #f7fafc; background: radial-gradient(circle at top, #384766, #111827 72%); font-family: system-ui, sans-serif; user-select: none; }
main { text-align: center; }
.mark { width: 74px; height: 74px; margin: 0 auto 22px; border: 5px solid rgba(255,255,255,.2); border-top-color: #70d6ff; border-radius: 50%; animation: spin .9s linear infinite; }
h1 { margin: 0 0 8px; font-size: 27px; letter-spacing: .02em; }
p { margin: 0; color: #b8c4da; font-size: 14px; }
@keyframes spin { to { transform: rotate(360deg); } }
</style></head><body><main><div class="mark"></div><h1>Starting Wails</h1><p>Initialising application services…</p></main></body></html>`

const mainHTML = `<!doctype html>
<html><head><meta charset="utf-8"><style>
html, body { height: 100%; margin: 0; }
body { display: grid; place-items: center; background: #f7fafc; color: #172033; font-family: system-ui, sans-serif; }
main { max-width: 540px; padding: 48px; text-align: center; }
h1 { font-size: 34px; margin-bottom: 12px; }
p { color: #526079; line-height: 1.6; }
</style></head><body><main><h1>Application initialised</h1><p>The native event loop remained responsive while the example service spent five seconds starting.</p></main></body></html>`

type slowService struct{}

func (*slowService) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	log.Println("slow service: starting five-second initialisation")
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		log.Println("slow service: initialisation complete")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	app := application.New(application.Options{
		Name:        "Splash Screen Lifecycle Demo",
		Description: "Shows a responsive splash screen while services initialise",
		Services: []application.Service{
			application.NewService(&slowService{}),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	splash := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:            "splash",
		Width:           460,
		Height:          300,
		Hidden:          true,
		Frameless:       true,
		DisableResize:   true,
		AlwaysOnTop:     true,
		InitialPosition: application.WindowCentered,
		HTML:            splashHTML,
	})

	mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:            "main",
		Title:           "Wails Lifecycle Demo",
		Width:           800,
		Height:          520,
		Hidden:          true,
		InitialPosition: application.WindowCentered,
		HTML:            mainHTML,
	})

	app.Event.OnApplicationEvent(events.Common.ApplicationStarting, func(*application.ApplicationEvent) {
		log.Println("ApplicationStarting: showing splash")
		splash.Show()
	})

	app.Event.OnApplicationEvent(events.Common.ApplicationInitialized, func(*application.ApplicationEvent) {
		log.Println("ApplicationInitialized: replacing splash with main window")
		splash.Close()
		mainWindow.Show()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
