---
title: "Utiliser Gin pour le routage"
description: "Guide complet sur l’intégration du framework web Gin aux applications Wails v3"
slug: "guides/gin-routing"
sourcePath: "guides/gin-routing.md"
---

Ce guide explique comment intégrer le [framework web Gin](https://github.com/gin-gonic/gin) à Wails v3. Gin est un framework web HTTP hautes performances écrit en Go qui facilite la création d’applications web et d’API.

## Introduction

Wails v3 fournit un système de ressources flexible qui permet d’utiliser n’importe quel gestionnaire HTTP, y compris des frameworks web populaires comme Gin. Cette intégration permet de :

- Servir du contenu web à l’aide du routage et des middlewares de Gin
- Créer des API REST accessibles depuis votre application Wails
- Utiliser les fonctionnalités de Gin tout en conservant l’intégration de Wails à l’environnement de bureau

## Configurer Gin avec Wails

Pour intégrer Gin à Wails, créez un routeur Gin et configurez-le comme gestionnaire de ressources dans votre application Wails. Procédez comme suit :

### 1. Installer les dépendances

Tout d’abord, vérifiez que le paquet Gin est installé :

```bash
go get -u github.com/gin-gonic/gin
```

### 2. Créer un middleware pour Gin

Créez une fonction middleware qui assurera l’intégration entre Wails et Gin :

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

Ce middleware transmet toutes les requêtes HTTP au routeur Gin.

### 3. Configurer votre routeur Gin

Configurez votre routeur Gin avec des routes, des middlewares et des gestionnaires :

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

### 4. Intégrer le routeur à l’application Wails

Configurez votre application Wails afin qu’elle utilise le routeur Gin comme gestionnaire de ressources :

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

## Servir du contenu statique

Utilisez le paquet `embed` de Go pour incorporer les fichiers statiques dans votre binaire :

```go
//go:embed static
var staticFiles embed.FS

ginEngine.StaticFS("/static", http.FS(staticFiles))
ginEngine.GET("/", func(c *gin.Context) {
    file, _ := staticFiles.ReadFile("static/index.html")
    c.Data(http.StatusOK, "text/html; charset=utf-8", file)
})
```

Pendant le développement, servez les fichiers directement depuis le disque à l’aide de `ginEngine.Static("/static", "./static")`.

## Middleware personnalisé

Gin permet de créer des middlewares personnalisés à différentes fins. Voici un exemple de middleware de journalisation :

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

## Traiter les requêtes d’API

Gin facilite la création d’API REST. Voici comment définir des points de terminaison d’API :

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

## Utiliser les fonctionnalités de Wails

Votre contenu web servi par Gin peut interagir avec les fonctionnalités de Wails à l’aide du paquet `@wailsio/runtime`.

### Gérer les événements

Enregistrez les gestionnaires d’événements en Go :

```go
app.Event.On("my-event", func(event *application.CustomEvent) {
    log.Printf("Received event: %v", event.Data)
})
```

Effectuez l’appel depuis JavaScript à l’aide du runtime :

```html
<script type="module">
    import { Events } from '@wailsio/runtime';

    Events.Emit("my-event", { message: "Hello from frontend" });
    Events.On("response-event", (data) => console.log(data));
</script>
```

## Configuration avancée

### Personnaliser le mode de Gin

En production, réglez Gin sur le mode release :

```go
gin.SetMode(gin.ReleaseMode)  // Use gin.DebugMode for development
ginEngine := gin.New()
```

### Gérer les WebSockets

Vous pouvez intégrer les WebSockets à Gin à l’aide de bibliothèques telles que Gorilla WebSocket :

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

## Bonnes pratiques

- **Utilisez le paquet embed de Go :** incorporez les fichiers statiques dans votre binaire afin de faciliter sa distribution.
- **Séparez les responsabilités :** séparez la logique de votre API de celle de votre interface utilisateur.
- **Gérez les erreurs :** mettez en œuvre une gestion appropriée des erreurs dans les routes Gin comme dans le code frontend.
- **Sécurité :** tenez compte des impératifs de sécurité, en particulier lors du traitement des saisies utilisateur.
- **Performances :** utilisez le mode release de Gin en production pour améliorer les performances.
- **Tests :** écrivez des tests pour vos routes Gin à l’aide des utilitaires de test de Gin.

## Conclusion

L’intégration de Gin à Wails offre une combinaison puissante pour créer des applications de bureau avec des technologies web. Les performances et les fonctionnalités de Gin complètent les capacités d’intégration de Wails à l’environnement de bureau, ce qui permet de créer des applications sophistiquées tirant le meilleur parti de ces deux technologies.

Pour plus d’informations, consultez la [documentation de Gin](https://github.com/gin-gonic/gin) et la [documentation de Wails](https://wails.io).
