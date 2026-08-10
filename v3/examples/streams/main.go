package main

import (
	"embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Streams Demo",
		Description: "A bidirectional stream between Go and JavaScript",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// One handler, one connection, one goroutine. The handler runs for as long
	// as the connection is open, so returning from it closes the connection.
	app.HandleStream("hello", func(c *application.StreamConn) {
		defer c.Close()
		log.Println("[Go] frontend connected")

		// Go -> JS: send the time once a second until the connection goes away.
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if err := c.Send([]byte("tick " + time.Now().Format("15:04:05"))); err != nil {
						return // frontend is gone
					}
				case <-c.Context().Done():
					return
				}
			}
		}()

		// JS -> Go: read frames until the connection closes, echoing each one.
		for {
			frame, err := c.Receive()
			if err != nil {
				log.Println("[Go] frontend disconnected")
				return
			}
			log.Printf("[Go] received: %s", frame)
			if err := c.Send([]byte("echo: " + string(frame))); err != nil {
				return
			}
		}
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Streams Demo",
		Width:  420,
		Height: 360,
		URL:    "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
