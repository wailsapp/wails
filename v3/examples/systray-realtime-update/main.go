package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	app := application.New(application.Options{
		Name:        "Systray Realtime Update Demo",
		Description: "Demonstrates real-time systray menu updates while menu is open",
		Assets:      application.AlphaAssets,
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	// Counter for dynamic updates
	var clickCount atomic.Int32

	// Create the systray menu
	menu := app.NewMenu()

	// Add a dynamic counter item
	counterItem := menu.Add("Clicks: 0")
	counterItem.OnClick(func(ctx *application.Context) {
		count := clickCount.Add(1)
		counterItem.SetLabel(fmt.Sprintf("Clicks: %d", count))
		menu.Update()
		log.Printf("Counter clicked! New count: %d", count)
	})

	menu.AddSeparator()

	// Add a timestamp item that updates every second
	timestampItem := menu.Add("Time: --:--:--")
	timestampItem.OnClick(func(ctx *application.Context) {
		log.Println("Timestamp item clicked!")
	})

	menu.AddSeparator()

	// Add checkbox items to test state updates
	checkbox1 := menu.AddCheckbox("Option A", false)
	checkbox1.OnClick(func(ctx *application.Context) {
		log.Printf("Option A toggled: %v", ctx.IsChecked())
		menu.Update()
	})

	checkbox2 := menu.AddCheckbox("Option B", true)
	checkbox2.OnClick(func(ctx *application.Context) {
		log.Printf("Option B toggled: %v", ctx.IsChecked())
		menu.Update()
	})

	menu.AddSeparator()

	// Status item that changes based on checkbox states
	statusItem := menu.Add("Status: Ready")
	statusItem.OnClick(func(ctx *application.Context) {
		log.Println("Status clicked!")
	})

	menu.AddSeparator()

	// Quit item
	quitItem := menu.Add("Quit")
	quitItem.OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	// Create the system tray
	systray := app.NewSystemTray()
	systray.SetMenu(menu)
	systray.SetLabel("RT")

	// Callbacks for menu open/close
	systray.OnMenuOpen(func() {
		log.Println(">>> Menu opened")
	})

	systray.OnMenuClose(func() {
		log.Println("<<< Menu closed")
	})

	// Background goroutine that updates the menu every second
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		updateCount := 0
		for range ticker.C {
			updateCount++

			// Update timestamp
			now := time.Now().Format("15:04:05")
			timestampItem.SetLabel(fmt.Sprintf("Time: %s", now))

			// Update status based on checkbox states
			var status string
			if checkbox1.Checked() && checkbox2.Checked() {
				status = "Status: All enabled"
			} else if checkbox1.Checked() {
				status = "Status: Only A"
			} else if checkbox2.Checked() {
				status = "Status: Only B"
			} else {
				status = "Status: None enabled"
			}
			statusItem.SetLabel(status)

			// Trigger menu update - this should work even while menu is open!
			menu.Update()

			if updateCount%5 == 0 {
				log.Printf("Background update #%d (menu should update in real-time)", updateCount)
			}
		}
	}()

	// Run the app
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
