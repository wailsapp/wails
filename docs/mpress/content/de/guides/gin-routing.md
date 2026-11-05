---
title: "Gin für das Routing verwenden"
description: "Umfassende Anleitung zur Integration des Gin-Webframeworks in Wails-v3-Anwendungen"
slug: "guides/gin-routing"
sourcePath: "guides/gin-routing.md"
---

Diese Anleitung zeigt, wie Sie das [Gin-Webframework](https://github.com/gin-gonic/gin) in Wails v3 integrieren. Gin ist ein in Go geschriebenes, leistungsstarkes HTTP-Webframework, mit dem sich Webanwendungen und APIs einfach erstellen lassen.

## Einführung

Wails v3 bietet ein flexibles Asset-System, mit dem Sie beliebige HTTP-Handler verwenden können, darunter verbreitete Webframeworks wie Gin. Diese Integration ermöglicht Ihnen Folgendes:

- Webinhalte mit dem Routing und der Middleware von Gin bereitstellen
- RESTful APIs erstellen, auf die Ihre Wails-Anwendung zugreifen kann
- Gin-Funktionen nutzen und zugleich die Desktop-Integration von Wails beibehalten

## Gin mit Wails einrichten

Um Gin in Wails zu integrieren, müssen Sie einen Gin-Router erstellen und ihn in Ihrer Wails-Anwendung als Asset-Handler konfigurieren. Gehen Sie dazu wie folgt vor:

### 1. Abhängigkeiten installieren

Stellen Sie zunächst sicher, dass das Gin-Paket installiert ist:

```bash
go get -u github.com/gin-gonic/gin
```

### 2. Middleware für Gin erstellen

Erstellen Sie eine Middleware-Funktion, welche die Integration zwischen Wails und Gin übernimmt:

```go
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
```

Diese Middleware leitet alle HTTP-Anfragen an den Gin-Router weiter.

### 3. Gin-Router konfigurieren

Richten Sie Ihren Gin-Router mit Routen, Middleware und Handlern ein:

```go
// Create a new Gin router
ginEngine := gin.New() // Using New() instead of Default() to add custom middleware

// Add middlewares
ginEngine.Use(gin.Recovery())
ginEngine.Use(LoggingMiddleware()) // Your custom middleware

// Define routes
ginEngine.GET("/", func(c *gin.Context) {
    // Serve your main page
})

ginEngine.GET("/api/hello", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "message": "Hello from Gin API!",
        "time":    time.Now().Format(time.RFC3339),
    })
})
```

### 4. In die Wails-Anwendung integrieren

Konfigurieren Sie Ihre Wails-Anwendung so, dass sie den Gin-Router als Asset-Handler verwendet:

```go
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

// Create window
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:    "Wails + Gin Example",
    Width:    900,
    Height:   700,
    URL:      "/", // This will load the route handled by Gin
})
```

## Statische Inhalte bereitstellen

Verwenden Sie das Go-Paket `embed`, um statische Dateien in Ihre Binärdatei einzubetten:

```go
//go:embed static
var staticFiles embed.FS

ginEngine.StaticFS("/static", http.FS(staticFiles))
ginEngine.GET("/", func(c *gin.Context) {
    file, _ := staticFiles.ReadFile("static/index.html")
    c.Data(http.StatusOK, "text/html; charset=utf-8", file)
})
```

Stellen Sie Dateien während der Entwicklung mit `ginEngine.Static("/static", "./static")` direkt vom Datenträger bereit.

## Benutzerdefinierte Middleware

Mit Gin können Sie benutzerdefinierte Middleware für verschiedene Zwecke erstellen. Das folgende Beispiel zeigt eine Protokollierungs-Middleware:

```go
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
```

## API-Anfragen verarbeiten

Mit Gin lassen sich RESTful APIs einfach erstellen. So definieren Sie API-Endpunkte:

```go
// GET endpoint
ginEngine.GET("/api/users", func(c *gin.Context) {
    c.JSON(http.StatusOK, users)
})

// POST endpoint with JSON binding
ginEngine.POST("/api/users", func(c *gin.Context) {
    var newUser User
    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Process the new user...
    c.JSON(http.StatusCreated, newUser)
})

// Path parameters
ginEngine.GET("/api/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    // Find user by ID...
    c.JSON(http.StatusOK, user)
})

// Query parameters
ginEngine.GET("/api/search", func(c *gin.Context) {
    query := c.DefaultQuery("q", "")
    limit := c.DefaultQuery("limit", "10")
    // Perform search...
    c.JSON(http.StatusOK, results)
})
```

## Wails-Funktionen verwenden

Ihre über Gin bereitgestellten Webinhalte können über das Paket `@wailsio/runtime` mit Wails-Funktionen interagieren.

### Ereignisbehandlung

Registrieren Sie Ereignis-Handler in Go:

```go
app.Event.On("my-event", func(event *application.CustomEvent) {
    log.Printf("Received event: %v", event.Data)
})
```

Rufen Sie sie über die Runtime aus JavaScript auf:

```html
<script type="module">
    import { Events } from '@wailsio/runtime';

    Events.Emit("my-event", { message: "Hello from frontend" });
    Events.On("response-event", (data) => console.log(data));
</script>
```

## Erweiterte Konfiguration

### Gin-Modus anpassen

Aktivieren Sie für den Produktivbetrieb den Release-Modus von Gin:

```go
gin.SetMode(gin.ReleaseMode)  // Use gin.DebugMode for development
ginEngine := gin.New()
```

### WebSockets verarbeiten

Sie können WebSockets mithilfe von Bibliotheken wie Gorilla WebSocket in Gin integrieren:

```go
import "github.com/gorilla/websocket"

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all connections
    },
}

// In your route handler:
ginEngine.GET("/ws", func(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()
    
    // Handle WebSocket connection...
})
```

## Bewährte Verfahren

- **Go-Paket embed verwenden:** Betten Sie statische Dateien in Ihre Binärdatei ein, um die Verteilung zu vereinfachen.
- **Zuständigkeiten trennen:** Halten Sie Ihre API-Logik von Ihrer UI-Logik getrennt.
- **Fehlerbehandlung:** Implementieren Sie sowohl in den Gin-Routen als auch im Frontend-Code eine angemessene Fehlerbehandlung.
- **Sicherheit:** Berücksichtigen Sie Sicherheitsaspekte, insbesondere bei der Verarbeitung von Benutzereingaben.
- **Leistung:** Verwenden Sie für eine höhere Leistung im Produktivbetrieb den Release-Modus von Gin.
- **Tests:** Schreiben Sie mithilfe der Testwerkzeuge von Gin Tests für Ihre Gin-Routen.

## Fazit

Die Integration von Gin in Wails bietet eine leistungsstarke Kombination zum Erstellen von Desktop-Anwendungen mit Webtechnologien. Die Leistung und der Funktionsumfang von Gin ergänzen die Desktop-Integrationsfunktionen von Wails. So können Sie anspruchsvolle Anwendungen erstellen, die die Vorteile beider Technologien vereinen.

Weitere Informationen finden Sie in der [Gin-Dokumentation](https://github.com/gin-gonic/gin) und der [Wails-Dokumentation](https://wails.io).
