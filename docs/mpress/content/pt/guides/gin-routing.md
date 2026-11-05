---
title: "Uso do Gin para roteamento"
description: "Um guia abrangente para integrar o framework web Gin a aplicações Wails v3"
slug: "guides/gin-routing"
sourcePath: "guides/gin-routing.md"
---

Este guia demonstra como integrar o [framework web Gin](https://github.com/gin-gonic/gin) ao Wails v3. O Gin é um framework web HTTP de alto desempenho escrito em Go que facilita a criação de aplicações web e APIs.

## Introdução

O Wails v3 oferece um sistema flexível de recursos que permite usar qualquer handler HTTP, incluindo frameworks web populares como o Gin. Essa integração permite:

- Servir conteúdo web usando o roteamento e os middlewares do Gin
- Criar APIs RESTful acessíveis pela sua aplicação Wails
- Usar os recursos do Gin mantendo a integração do Wails com o ambiente de desktop

## Configuração do Gin com o Wails

Para integrar o Gin ao Wails, crie um roteador Gin e configure-o como o handler de recursos da sua aplicação Wails. Veja o passo a passo:

### 1. Instale as dependências

Primeiro, verifique se o pacote Gin está instalado:

```bash
go get -u github.com/gin-gonic/gin
```

### 2. Crie um middleware para o Gin

Crie uma função de middleware para gerenciar a integração entre o Wails e o Gin:

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

Esse middleware encaminha todas as requisições HTTP ao roteador Gin.

### 3. Configure seu roteador Gin

Configure seu roteador Gin com rotas, middlewares e handlers:

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

### 4. Integre-o à aplicação Wails

Configure sua aplicação Wails para usar o roteador Gin como handler de recursos:

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

## Como servir conteúdo estático

Use o pacote `embed` do Go para incorporar arquivos estáticos ao seu binário:

```go
//go:embed static
var staticFiles embed.FS

ginEngine.StaticFS("/static", http.FS(staticFiles))
ginEngine.GET("/", func(c *gin.Context) {
    file, _ := staticFiles.ReadFile("static/index.html")
    c.Data(http.StatusOK, "text/html; charset=utf-8", file)
})
```

Durante o desenvolvimento, sirva os arquivos diretamente do disco usando `ginEngine.Static("/static", "./static")`.

## Middleware personalizado

O Gin permite criar middlewares personalizados para diversas finalidades. Veja um exemplo de middleware de registro de logs:

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

## Tratamento de requisições de API

O Gin facilita a criação de APIs RESTful. Veja como definir endpoints de API:

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

## Uso dos recursos do Wails

Seu conteúdo web servido pelo Gin pode interagir com os recursos do Wails usando o pacote `@wailsio/runtime`.

### Tratamento de eventos

Registre handlers de eventos em Go:

```go
app.Event.On("my-event", func(event *application.CustomEvent) {
    log.Printf("Received event: %v", event.Data)
})
```

Faça a chamada a partir do JavaScript usando o runtime:

```html
<script type="module">
    import { Events } from '@wailsio/runtime';

    Events.Emit("my-event", { message: "Hello from frontend" });
    Events.On("response-event", (data) => console.log(data));
</script>
```

## Configuração avançada

### Personalização do modo do Gin

Defina o Gin no modo de release para produção:

```go
gin.SetMode(gin.ReleaseMode)  // Use gin.DebugMode for development
ginEngine := gin.New()
```

### Tratamento de WebSockets

Você pode integrar WebSockets ao Gin usando bibliotecas como a Gorilla WebSocket:

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

## Práticas recomendadas

- **Use o pacote embed do Go:** Incorpore arquivos estáticos ao seu binário para facilitar a distribuição.
- **Separe as responsabilidades:** Mantenha a lógica da API separada da lógica da interface do usuário.
- **Tratamento de erros:** Implemente o tratamento adequado de erros tanto nas rotas do Gin quanto no código do frontend.
- **Segurança:** Esteja atento às questões de segurança, especialmente ao processar entradas do usuário.
- **Desempenho:** Use o modo de release do Gin em produção para obter melhor desempenho.
- **Testes:** Escreva testes para suas rotas do Gin usando os utilitários de teste do Gin.

## Conclusão

A integração do Gin ao Wails oferece uma combinação poderosa para criar aplicações de desktop com tecnologias web. O desempenho e o conjunto de recursos do Gin complementam os recursos de integração do Wails com o ambiente de desktop, permitindo criar aplicações sofisticadas que aproveitam o melhor das duas tecnologias.

Para obter mais informações, consulte a [documentação do Gin](https://github.com/gin-gonic/gin) e a [documentação do Wails](https://wails.io).
