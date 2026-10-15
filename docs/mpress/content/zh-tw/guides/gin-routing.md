---
title: "使用 Gin 進行路由處理"
description: "將 Gin Web 框架與 Wails v3 應用程式整合的完整指南"
slug: "guides/gin-routing"
sourcePath: "guides/gin-routing.md"
---

本指南說明如何將[Gin Web 框架](https://github.com/gin-gonic/gin)與 Wails v3 整合。Gin 是以 Go 編寫的高效能 HTTP Web 框架，可讓您輕鬆建置 Web 應用程式和 API。

## 簡介

Wails v3 提供彈性的資產系統，可讓您使用任何 HTTP 處理常式，包括 Gin 等熱門 Web 框架。透過這項整合，您可以：

- 使用 Gin 的路由和中介軟體提供 Web 內容
- 建立可從 Wails 應用程式存取的 RESTful API
- 在維持 Wails 桌面整合的同時使用 Gin 的功能

## 設定 Gin 與 Wails

若要整合 Gin 與 Wails，您需要建立 Gin 路由器，並將其設定為 Wails 應用程式的資產處理常式。以下是逐步指南：

### 1. 安裝相依套件

首先，請確認已安裝 Gin 套件：

```bash
go get -u github.com/gin-gonic/gin
```

### 2. 建立 Gin 中介軟體

建立一個中介軟體函式，負責處理 Wails 與 Gin 之間的整合：

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

此中介軟體會將所有 HTTP 請求傳遞至 Gin 路由器。

### 3. 設定 Gin 路由器

使用路由、中介軟體和處理常式設定 Gin 路由器：

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

### 4. 與 Wails 應用程式整合

設定 Wails 應用程式，使其將 Gin 路由器用作資產處理常式：

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

## 提供靜態內容

使用 Go 的`embed`套件將靜態檔案嵌入二進位檔：

```go
//go:embed static
var staticFiles embed.FS

ginEngine.StaticFS("/static", http.FS(staticFiles))
ginEngine.GET("/", func(c *gin.Context) {
    file, _ := staticFiles.ReadFile("static/index.html")
    c.Data(http.StatusOK, "text/html; charset=utf-8", file)
})
```

開發期間，請使用`ginEngine.Static("/static", "./static")`直接從磁碟提供檔案。

## 自訂中介軟體

Gin 可讓您針對各種用途建立自訂中介軟體。以下是記錄中介軟體的範例：

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

## 處理 API 請求

Gin 可讓您輕鬆建立 RESTful API。以下是定義 API 端點的方法：

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

## 使用 Wails 功能

由 Gin 提供的 Web 內容可使用`@wailsio/runtime`套件與 Wails 功能互動。

### 事件處理

在 Go 中註冊事件處理常式：

```go
app.Event.On("my-event", func(event *application.CustomEvent) {
    log.Printf("Received event: %v", event.Data)
})
```

使用 runtime 從 JavaScript 呼叫：

```html
<script type="module">
    import { Events } from '@wailsio/runtime';

    Events.Emit("my-event", { message: "Hello from frontend" });
    Events.On("response-event", (data) => console.log(data));
</script>
```

## 進階設定

### 自訂 Gin 模式

在正式環境中將 Gin 設為 release 模式：

```go
gin.SetMode(gin.ReleaseMode)  // Use gin.DebugMode for development
ginEngine := gin.New()
```

### 處理 WebSocket

您可以使用 Gorilla WebSocket 等程式庫，將 WebSocket 與 Gin 整合：

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

## 最佳實務

- <strong>使用 Go 的 embed 套件：</strong>將靜態檔案嵌入二進位檔，以便於散布。
- <strong>分離關注點：</strong>將 API 邏輯與 UI 邏輯分開。
- <strong>錯誤處理：</strong>在 Gin 路由和前端程式碼中實作適當的錯誤處理。
- <strong>安全性：</strong>請留意安全性考量，尤其是在處理使用者輸入時。
- <strong>效能：</strong>在正式環境中使用 Gin 的 release 模式，以獲得更佳效能。
- <strong>測試：</strong>使用 Gin 的測試工具為 Gin 路由撰寫測試。

## 結論

整合 Gin 與 Wails，能為使用 Web 技術建置桌面應用程式提供強大的組合。Gin 的效能和功能集可與 Wails 的桌面整合能力相輔相成，讓您能結合兩者的優勢，建立功能完善的應用程式。

如需更多資訊，請參閱[Gin 文件](https://github.com/gin-gonic/gin)和[Wails 文件](https://wails.io)。
