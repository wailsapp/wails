---
title: "Utiliser Gin pour les services"
description: "Guide d’intégration du framework web Gin aux services Wails v3"
slug: "guides/gin-services"
sourcePath: "guides/gin-services.md"
---

## Utiliser Gin pour les services

Le framework web Gin est un choix courant pour créer des services HTTP en Go. Avec Wails v3, vous pouvez facilement intégrer des services basés sur Gin à votre application afin de disposer d’une solution puissante pour traiter les requêtes HTTP, implémenter des API RESTful et servir du contenu web.

Ce guide vous accompagne dans la création d’un service basé sur Gin pouvant être monté sur une route précise de votre application Wails. Nous construirons un exemple complet qui montre comment :

1. Créer un service basé sur Gin
2. Implémenter l’interface Service de Wails
3. Configurer les routes et les middlewares
4. Intégrer le système d’événements de Wails
5. Interagir avec le service depuis le frontend

## Prérequis

Avant de commencer, vérifiez que vous disposez des éléments suivants :

- Wails v3 installé
- Connaissances de base de Go et du framework Gin
- Connaissance des concepts HTTP et des API RESTful

Vous devez ajouter le framework Gin à votre projet :

```bash
go get github.com/gin-gonic/gin
```

## Créer un service basé sur Gin

Commençons par créer un service Gin qui implémente l’interface Service de Wails. Notre service gérera une collection d’utilisateurs et fournira des points de terminaison d’API permettant de récupérer et de créer des enregistrements utilisateur.

### 1. Définir vos modèles de données

Commencez par définir les structures de données que votre service utilisera :

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

### 2. Créer la structure de votre service

Définissez ensuite la structure du service qui contiendra votre routeur Gin et tout état que votre service doit conserver :

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

### 3. Implémenter l’interface Service

Implémentez les méthodes requises par l’interface Service de Wails :

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

### 3. Implémenter l’interface http.Handler

Pour permettre le montage de votre service sur une route précise, implémentez l’interface `http.Handler`. Sa méthode unique, `ServeHTTP`, constitue le point d’entrée de toutes les requêtes HTTP adressées à votre service. Elle délègue leur traitement au routeur Gin, ce qui vous permet d’utiliser toutes les puissantes fonctionnalités de routage et de middleware de Gin.

```go
// ServeHTTP implements the http.Handler interface
func (s *GinService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// All requests go to the Gin router
	s.ginEngine.ServeHTTP(w, r)
}
```

### 4. Configurer vos routes

Pour une meilleure organisation, définissez vos routes d’API dans une méthode distincte. Cette approche garde votre code clair et facilite la compréhension de la structure de votre API. Le routeur Gin fournit une API fluide pour définir les routes et prend notamment en charge les groupes de routes, qui permettent d’organiser les points de terminaison associés.

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

### 5. Créer un middleware personnalisé

Vous pouvez créer un middleware Gin personnalisé pour enrichir votre service. Dans Gin, les fonctions de middleware s’exécutent dans l’ordre où elles sont ajoutées au routeur et peuvent effectuer des tâches telles que la journalisation, l’authentification et la gestion des erreurs. Cet exemple présente un middleware de journalisation simple qui consigne les détails des requêtes.

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

## Enregistrer votre service

Pour utiliser votre service basé sur Gin dans une application Wails, vous devez l’enregistrer auprès de l’application et préciser la route sur laquelle il doit être monté. Cette opération s’effectue lors de la création de l’instance de l’application Wails. La route indiquée devient le chemin de base de tous les points de terminaison définis dans votre routeur Gin.

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

Dans cet exemple, le service Gin est monté sur la route `/api`. Ainsi, si un point de terminaison `/info` est défini dans votre routeur Gin, il sera accessible dans votre application à l’adresse `/api/info`. Cette approche vous permet d’organiser logiquement les points de terminaison de votre API et d’éviter les conflits avec les autres parties de votre application.

## Intégrer le système d’événements de Wails

L’un des principaux atouts de l’utilisation de Gin avec les services Wails réside dans la possibilité d’intégrer facilement le système d’événements de Wails. Vous pouvez ainsi établir une communication en temps réel entre votre service backend et le frontend.

Dans la méthode `ServiceStartup` de votre service, vous pouvez enregistrer des gestionnaires d’événements afin de traiter ceux provenant du frontend :

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

Vous pouvez également émettre des événements vers le frontend depuis vos routes Gin ou d’autres parties de votre service :

```go
// After creating a new user
s.app.Event.Emit("user-created", newUser)
```

## Intégration au frontend

Pour interagir avec votre service Gin depuis le frontend, vous devez importer le runtime Wails, envoyer des requêtes HTTP aux points de terminaison de votre API et utiliser le système d’événements de Wails pour la communication en temps réel.

Pour une utilisation en production, il est recommandé d’utiliser le paquet `@wailsio/runtime` plutôt que d’importer directement `/wails/runtime.js`. Vous bénéficiez ainsi de la sûreté des types, d’une meilleure prise en charge par l’IDE, de la gestion des versions via npm et de la compatibilité avec les outils JavaScript modernes.

Installez le paquet :

```bash
npm install @wailsio/runtime
```

Utilisez-le ensuite dans votre code :

```javascript
import * as wails from '@wailsio/runtime';

// Event emission
wails.Events.Emit('gin-api-event', eventData);
```

Voici un exemple de configuration de l’intégration au frontend :

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

## Conclusion

L’intégration du framework web Gin aux services Wails v3 offre une approche puissante et flexible pour créer des applications web modulaires et faciles à maintenir. En associant les fonctionnalités de routage et de middleware de Gin au système d’événements de Wails, vous pouvez créer des applications riches et interactives tout en assurant une séparation claire des responsabilités.

Le code complet de l’exemple présenté dans ce guide se trouve dans le dépôt Wails, sous `v3/examples/gin-service`.
