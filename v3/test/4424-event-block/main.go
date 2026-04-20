package main

import (
	_ "embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/icons"
)

func main() {
	app := application.New(application.Options{
		Name:        "Event Block Test",
		Description: "Test for issue #4424 - Event emitter blocks menu button press",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	// Create a system tray with a menu
	systemTray := app.SystemTray.New()
	systemTray.SetIcon(icons.WailsLogoBlack)

	menu := app.NewMenu()
	menu.Add("About").OnClick(func(ctx *application.Context) {
		log.Println("About clicked!")
		application.InfoDialog().
			SetTitle("About").
			SetMessage("This is a test for issue #4424").
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

	// Start emitting events after application starts
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		go func() {
			log.Println("Starting event emitter...")
			log.Println("TEST: Right-click the systray icon and click menu items.")
			log.Println("They should continue to work even with events being emitted.")
			log.Println("BUG: On Linux, menu clicks stop working after events start.")

			counter := 0
			for {
				counter++
				log.Printf("Emitting event #%d...\n", counter)
				app.Event.Emit("backendmessage", map[string]interface{}{
					"counter": counter,
					"time":    time.Now().Format(time.RFC3339),
				})
				time.Sleep(3 * time.Second)
			}
		}()
	})

	// Run the application
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
