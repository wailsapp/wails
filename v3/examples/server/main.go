package main

import (
	"embed"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/dist
var assets embed.FS

// GreetService is a simple service that provides greeting functionality.
type GreetService struct{}

// Greet returns a greeting message.
func (g *GreetService) Greet(name string) string {
	if name == "" {
		name = "World"
	}
	return "Hello, " + name + "!"
}

func main() {
	// Create a logger for better visibility
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	app := application.New(application.Options{
		Name:        "Server Mode Example",
		Description: "A Wails application running in server mode",
		Logger:      logger,
		LogLevel:    slog.LevelInfo,

		// Server mode is enabled by building with -tags server
		// Host/port can be overridden via WAILS_SERVER_HOST and WAILS_SERVER_PORT env vars
		Server: application.ServerOptions{
			Host: "localhost",
			Port: 8080,
		},

		// Register services (bindings work the same as desktop mode)
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},

		// Serve frontend assets
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	log.Println("Starting Wails application in server mode...")
	log.Println("Access at: http://localhost:8080")
	log.Println("Health check: http://localhost:8080/health")
	log.Println("Press Ctrl+C to stop")

	// Listen for broadcast events from browsers
	app.Event.On("broadcast", func(event *application.CustomEvent) {
		log.Printf("Received broadcast from %s: %v\n", event.Sender, event.Data)
	})

	// Emit periodic events to test WebSocket broadcasting
	go func() {
		time.Sleep(2 * time.Second) // Wait for server to start
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			app.Event.Emit("server-tick", time.Now().Format(time.RFC3339))
			log.Println("Emitted server-tick event")
		}
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
