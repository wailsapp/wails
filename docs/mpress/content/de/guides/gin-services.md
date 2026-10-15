---
title: "Gin für Services verwenden"
description: "Eine Anleitung zur Integration des Webframeworks Gin in Wails-v3-Services"
slug: "guides/gin-services"
sourcePath: "guides/gin-services.md"
---

## Gin für Services verwenden

Das Webframework Gin ist eine beliebte Wahl für die Entwicklung von HTTP-Services in Go. Mit Wails v3 können Sie Gin-basierte Services einfach in Ihre Anwendung integrieren. So können Sie HTTP-Anfragen effizient verarbeiten, RESTful APIs implementieren und Webinhalte bereitstellen.

Diese Anleitung führt Sie durch die Erstellung eines Gin-basierten Service, der unter einer bestimmten Route in Ihrer Wails-Anwendung eingebunden werden kann. An einem vollständigen Beispiel wird gezeigt, wie Sie:

1. einen Gin-basierten Service erstellen
2. die Wails-Schnittstelle Service implementieren
3. Routen und Middleware einrichten
4. den Service in das Wails-Ereignissystem integrieren
5. vom Frontend aus mit dem Service interagieren

## Voraussetzungen

Stellen Sie zunächst sicher, dass Sie über Folgendes verfügen:

- eine Installation von Wails v3
- Grundkenntnisse in Go und im Gin-Framework
- Vertrautheit mit HTTP-Konzepten und RESTful APIs

Sie müssen das Gin-Framework zu Ihrem Projekt hinzufügen:

```bash
go get github.com/gin-gonic/gin
```

## Einen Gin-basierten Service erstellen

Erstellen Sie zunächst einen Gin-Service, der die Wails-Schnittstelle Service implementiert. Der Service verwaltet eine Sammlung von Benutzern und stellt API-Endpunkte zum Abrufen und Erstellen von Benutzerdatensätzen bereit.

### 1. Datenmodelle definieren

Definieren Sie zunächst die Datenstrukturen, mit denen der Service arbeitet:

```go
package services

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// EventData represents data sent in events
type EventData struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
```

### 2. Service-Struktur erstellen

Definieren Sie anschließend die Service-Struktur, die den Gin-Router und den vom Service zu verwaltenden Zustand enthält:

```go
// GinService implements a Wails service that uses Gin for HTTP handling
type GinService struct {
	ginEngine *gin.Engine
	users     []User
	nextID    int
	mu        sync.RWMutex
	app       *application.App
}

// NewGinService creates a new GinService instance
func NewGinService() *GinService {
	// Create a new Gin router
	ginEngine := gin.New()

	// Add middlewares
	ginEngine.Use(gin.Recovery())
	ginEngine.Use(LoggingMiddleware())

	service := &GinService{
		ginEngine: ginEngine,
		users: []User{
			{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Now().Add(-72 * time.Hour)},
			{ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: time.Now().Add(-48 * time.Hour)},
			{ID: 3, Name: "Charlie", Email: "charlie@example.com", CreatedAt: time.Now().Add(-24 * time.Hour)},
		},
		nextID: 4,
	}

	// Define routes
	service.setupRoutes()

	return service
}
```

### 3. Service-Schnittstelle implementieren

Implementieren Sie die erforderlichen Methoden der Wails-Schnittstelle Service:

```go
// ServiceName returns the name of the service
func (s *GinService) ServiceName() string {
	return "Gin API Service"
}

// ServiceStartup is called when the service starts
func (s *GinService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// Store the application instance for later use
	s.app = application.Get()

	// Register an event handler that can be triggered from the frontend
	s.app.Event.On("gin-api-event", func(event *application.CustomEvent) {
		// Log the event data
		s.app.Logger.Info("Received event from frontend", "data", event.Data)

		// Emit an event back to the frontend
		s.app.Event.Emit("gin-api-response",
			map[string]interface{}{
                "message": "Response from Gin API Service",
                "time":    time.Now().Format(time.RFC3339),
            },
		)
	})

	return nil
}

// ServiceShutdown is called when the service shuts down.
// IMPORTANT: the interface is `ServiceShutdown() error` (no ctx); a method that
// takes a context.Context does NOT satisfy the interface and will never be called.
func (s *GinService) ServiceShutdown() error {
	// Clean up resources if needed
	return nil
}
```

### 3. Schnittstelle http.Handler implementieren

Damit Ihr Service unter einer bestimmten Route eingebunden werden kann, implementieren Sie die Schnittstelle `http.Handler`. Ihre einzige Methode, `ServeHTTP`, ist der Einstiegspunkt für alle HTTP-Anfragen an den Service. Sie delegiert die Verarbeitung der Anfragen an den Gin-Router, sodass Sie sämtliche leistungsfähigen Routing- und Middleware-Funktionen von Gin verwenden können.

```go
// ServeHTTP implements the http.Handler interface
func (s *GinService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// All requests go to the Gin router
	s.ginEngine.ServeHTTP(w, r)
}
```

### 4. Routen einrichten

Definieren Sie Ihre API-Routen zur besseren Strukturierung in einer separaten Methode. Dadurch bleibt der Code übersichtlich und die Struktur der API ist leichter verständlich. Der Gin-Router bietet eine Fluent API zum Definieren von Routen und unterstützt dabei auch Routengruppen, mit denen sich zusammengehörige Endpunkte übersichtlich anordnen lassen.

```go
// setupRoutes configures the API routes
func (s *GinService) setupRoutes() {
	// Basic info endpoint
	s.ginEngine.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Gin API Service",
			"version": "1.0.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Users group
	users := s.ginEngine.Group("/users")
	{
		// Get all users
		users.GET("", func(c *gin.Context) {
			s.mu.RLock()
			defer s.mu.RUnlock()
			c.JSON(http.StatusOK, s.users)
		})

		// Get user by ID
		users.GET("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}

			s.mu.RLock()
			defer s.mu.RUnlock()

			for _, user := range s.users {
				if user.ID == id {
					c.JSON(http.StatusOK, user)
					return
				}
			}

			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		})

		// Create a new user
		users.POST("", func(c *gin.Context) {
			var newUser User
			if err := c.ShouldBindJSON(&newUser); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			s.mu.Lock()
			defer s.mu.Unlock()

			// Set the ID and creation time
			newUser.ID = s.nextID
			newUser.CreatedAt = time.Now()
			s.nextID++

			// Add to the users slice
			s.users = append(s.users, newUser)

			c.JSON(http.StatusCreated, newUser)

			// Emit an event to notify about the new user
			s.app.Event.Emit("user-created", newUser)
		})

		// Delete a user
		users.DELETE("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}

			s.mu.Lock()
			defer s.mu.Unlock()

			for i, user := range s.users {
				if user.ID == id {
					// Remove the user from the slice
					s.users = append(s.users[:i], s.users[i+1:]...)
					c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
					return
				}
			}

			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		})
	}
}
```

### 5. Benutzerdefinierte Middleware erstellen

Sie können benutzerdefinierte Gin-Middleware erstellen, um Ihren Service zu erweitern. Middleware-Funktionen werden in Gin in der Reihenfolge ausgeführt, in der sie dem Router hinzugefügt wurden, und können Aufgaben wie Protokollierung, Authentifizierung und Fehlerbehandlung übernehmen. Dieses Beispiel zeigt eine einfache Protokollierungs-Middleware, die Anfragedetails aufzeichnet.

```go
// LoggingMiddleware is a Gin middleware that logs request details
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request details
		log.Printf("[GIN] %s %s %d %s", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency)
	}
}
```

## Service registrieren

Um Ihren Gin-basierten Service in einer Wails-Anwendung zu verwenden, müssen Sie ihn bei der Anwendung registrieren und die Route angeben, unter der er eingebunden werden soll. Dies geschieht beim Erstellen der Wails-Anwendungsinstanz. Die angegebene Route wird zum Basispfad für alle in Ihrem Gin-Router definierten Endpunkte.

```go
app := application.New(application.Options{
    Name:        "Gin Service Demo",
    Description: "A demo of using Gin in Wails services",
    Mac: application.MacOptions{
        ApplicationShouldTerminateAfterLastWindowClosed: true,
    },
    LogLevel: slog.LevelDebug,
    Services: []application.Service{
        application.NewServiceWithOptions(services.NewGinService(), application.ServiceOptions{
            Route: "/api",
        }),
    },
    Assets: application.AssetOptions{
        Handler: application.BundledAssetFileServer(assets),
    },
})
```

In diesem Beispiel wird der Gin-Service unter der Route `/api` eingebunden. Wenn in Ihrem Gin-Router also ein Endpunkt als `/info` definiert ist, ist er in Ihrer Anwendung unter `/api/info` erreichbar. Mit diesem Ansatz können Sie Ihre API-Endpunkte logisch strukturieren und Konflikte mit anderen Teilen Ihrer Anwendung vermeiden.

## Integration in das Wails-Ereignissystem

Ein großer Vorteil der Verwendung von Gin mit Wails Services ist die nahtlose Integration in das Wails-Ereignissystem. Dies ermöglicht die Echtzeitkommunikation zwischen Ihrem Backend-Service und dem Frontend.

In der Methode `ServiceStartup` Ihres Service können Sie Ereignishandler registrieren, um Ereignisse aus dem Frontend zu verarbeiten:

```go
s.app.Event.On("gin-api-event", func(event *application.CustomEvent) {
	// Log the event data
	s.app.Logger.Info("Received event from frontend", "data", event.Data)

	// Emit an event back to the frontend
	s.app.Event.Emit("gin-api-response",
		map[string]interface{}{
			"message": "Response from Gin API Service",
			"time":    time.Now().Format(time.RFC3339),
		},
	)
})
```

Sie können außerdem aus Ihren Gin-Routen oder anderen Teilen Ihres Service Ereignisse an das Frontend senden:

```go
// After creating a new user
s.app.Event.Emit("user-created", newUser)
```

## Frontend-Integration

Um vom Frontend aus mit Ihrem Gin-Service zu interagieren, müssen Sie die Wails-Runtime importieren, HTTP-Anfragen an Ihre API-Endpunkte senden und das Wails-Ereignissystem für die Echtzeitkommunikation verwenden.

Für den Produktiveinsatz wird empfohlen, das Paket `@wailsio/runtime` zu verwenden, anstatt `/wails/runtime.js` direkt zu importieren. Dies gewährleistet Typsicherheit, bessere IDE-Unterstützung, Versionsverwaltung über npm und Kompatibilität mit modernen JavaScript-Werkzeugen.

Installieren Sie das Paket:

```bash
npm install @wailsio/runtime
```

Verwenden Sie es anschließend in Ihrem Code:

```javascript
import * as wails from '@wailsio/runtime';

// Event emission
wails.Events.Emit('gin-api-event', eventData);
```

Das folgende Beispiel zeigt, wie Sie die Frontend-Integration einrichten:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Gin Service Example</title>
    <!-- Styles omitted for brevity -->
</head>
<body>
    <h1>Gin Service Example</h1>
    
    <div class="card">
        <h2>API Endpoints</h2>
        <p>Try the Gin API endpoints mounted at /api:</p>
        
        <button id="getInfo">Get Service Info</button>
        <button id="getUsers">Get All Users</button>
        <button id="getUser">Get User by ID</button>
        <button id="createUser">Create User</button>
        <button id="deleteUser">Delete User</button>

        <div id="apiResult">
            <pre id="apiResponse">Results will appear here...</pre>
        </div>
    </div>
    
    <div class="card">
        <h2>Event Communication</h2>
        <p>Trigger an event to communicate with the Gin service:</p>
        
        <button id="triggerEvent">Trigger Event</button>
        
        <div id="eventResult">
            <pre id="eventResponse">Event responses will appear here...</pre>
        </div>
    </div>
    
    <div class="card" id="createUserForm" style="display: none; border: 2px solid #0078d7;">
        <h2>Create New User</h2>

        <div>
            <label for="userName">Name:</label>
            <input type="text" id="userName" placeholder="Enter name">
        </div>

        <div>
            <label for="userEmail">Email:</label>
            <input type="email" id="userEmail" placeholder="Enter email">
        </div>

        <button id="submitUser">Submit</button>
        <button id="cancelCreate">Cancel</button>
    </div>

    <script type="module">
        // Import the Wails runtime
        // Note: In production, use '@wailsio/runtime' instead
        import * as wails from "/wails/runtime.js";
        
        // Helper function to fetch API endpoints
        async function fetchAPI(endpoint, options = {}) {
            try {
                const response = await fetch(`/api${endpoint}`, options);
                const data = await response.json();
                
                document.getElementById('apiResponse').textContent = JSON.stringify(data, null, 2);
                return data;
            } catch (error) {
                document.getElementById('apiResponse').textContent = `Error: ${error.message}`;
                console.error('API Error:', error);
            }
        }
        
        // Event listeners for API buttons
        document.getElementById('getInfo').addEventListener('click', () => {
            fetchAPI('/info');
        });
        
        document.getElementById('getUsers').addEventListener('click', () => {
            fetchAPI('/users');
        });

        document.getElementById('getUser').addEventListener('click', async () => {
            const userId = prompt('Enter user ID:');
            if (userId) {
                await fetchAPI(`/users/${userId}`);
            }
        });

        document.getElementById('createUser').addEventListener('click', () => {
            const form = document.getElementById('createUserForm');
            form.style.display = 'block';
            form.scrollIntoView({ behavior: 'smooth' });
        });

        document.getElementById('cancelCreate').addEventListener('click', () => {
            document.getElementById('createUserForm').style.display = 'none';
        });

        document.getElementById('submitUser').addEventListener('click', async () => {
            const name = document.getElementById('userName').value;
            const email = document.getElementById('userEmail').value;

            if (!name || !email) {
                alert('Please enter both name and email');
                return;
            }

            try {
                await fetchAPI('/users', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ name, email })
                });

                document.getElementById('createUserForm').style.display = 'none';
                document.getElementById('userName').value = '';
                document.getElementById('userEmail').value = '';

                // Automatically fetch the updated user list
                await fetchAPI('/users');

                // Show a success message
                const apiResponse = document.getElementById('apiResponse');
                const currentData = JSON.parse(apiResponse.textContent);
                apiResponse.textContent = JSON.stringify({
                    message: "User created successfully!",
                    users: currentData
                }, null, 2);
            } catch (error) {
                console.error('Error creating user:', error);
            }
        });

        document.getElementById('deleteUser').addEventListener('click', async () => {
            const userId = prompt('Enter user ID to delete:');
            if (userId) {
                try {
                    await fetchAPI(`/users/${userId}`, {
                        method: 'DELETE'
                    });

                    // Show success message
                    document.getElementById('apiResponse').textContent = JSON.stringify({
                        message: `User with ID ${userId} deleted successfully`
                    }, null, 2);

                    // Refresh the user list
                    setTimeout(() => fetchAPI('/users'), 1000);
                } catch (error) {
                    console.error('Error deleting user:', error);
                }
            }
        });

        // Using Wails Events API for event communication
        document.getElementById('triggerEvent').addEventListener('click', async () => {
            // Display the event being sent
            document.getElementById('eventResponse').textContent = JSON.stringify({
                status: "Sending event to backend...",
                data: { timestamp: new Date().toISOString() }
            }, null, 2);
            
            // Use the Wails runtime to emit an event
            const eventData = {
                message: "Hello from the frontend!",
                timestamp: new Date().toISOString()
            };
            wails.Events.Emit('gin-api-event', eventData);
        });
        
        // Set up event listener for responses from the backend
        window.addEventListener('DOMContentLoaded', () => {
            // Register event listener using Wails runtime
            wails.Events.On("gin-api-response", (data) => {
                document.getElementById('eventResponse').textContent = JSON.stringify(data, null, 2);
            });
            
            // Also listen for user-created events
            wails.Events.On("user-created", (data) => {
                document.getElementById('eventResponse').textContent = JSON.stringify({
                    event: "user-created",
                    user: data
                }, null, 2);
            });

            // Initial API call to get service info
            fetchAPI('/info');
        });
    </script>
</body>
</html>
```

## Fazit

Die Integration des Webframeworks Gin in Wails-v3-Services bietet einen leistungsfähigen und flexiblen Ansatz für die Entwicklung modularer, wartbarer Webanwendungen. Indem Sie die Routing- und Middleware-Funktionen von Gin zusammen mit dem Wails-Ereignissystem nutzen, können Sie umfangreiche, interaktive Anwendungen mit einer klaren Trennung der Zuständigkeiten erstellen.

Den vollständigen Beispielcode zu dieser Anleitung finden Sie im Wails-Repository unter `v3/examples/gin-service`.
