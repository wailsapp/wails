package main

import (
	"embed"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed static
var staticFiles embed.FS

// GinMiddleware creates a middleware that passes requests to Gin if they're not handled by Wails
func GinMiddleware(ginEngine *gin.Engine) application.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Let Wails handle the `/wails` route
			if strings.HasPrefix(r.URL.Path, "/wails") {
				next.ServeHTTP(w, r)
				return
			}
			// Let Gin handle everything else
			ginEngine.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware is a Gin middleware that logs request details
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Log request details
		log.Printf("[GIN] %s | %s | %s | %d | %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			latency,
		)
	}
}

func main() {
	// Create a new Gin router
	ginEngine := gin.New() // Using New() instead of Default() to add our own middleware

	// Add middlewares
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(LoggingMiddleware())

	// Serve embedded static files
	ginEngine.StaticFS("/static", http.FS(staticFiles))

	// Define routes
	ginEngine.GET("/", func(c *gin.Context) {
		file, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading index.html")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", file)
	})

	ginEngine.GET("/api/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from Gin API!",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Create a new Wails application
	app := application.New(application.Options{
		Name:        "Gin Example",
		Description: "A demo of using Gin with Wails",
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Assets: application.AssetOptions{
			Handler:    ginEngine,
			Middleware: GinMiddleware(ginEngine),
		},
	})

	// Register event handler and store cleanup function
	removeGinHandler := app.Events.On("gin-button-clicked", func(event *application.CustomEvent) {
		log.Printf("Received event from frontend: %v", event.Data)
	})
	// Note: In production, call removeGinHandler() during cleanup
	_ = removeGinHandler

	// Create window
	app.Windows.NewWithOptions(application.WebviewWindowOptions{
		Title:  "Wails + Gin Example",
		Width:  900,
		Height: 700,
		URL:    "/",
	})

	// Run the app
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
