---
title: "Gin をルーティングに使用する"
description: "Gin Web フレームワークを Wails v3 アプリケーションに統合するための包括的なガイド"
slug: "guides/gin-routing"
sourcePath: "guides/gin-routing.md"
---

このガイドでは、[Gin Web フレームワーク](https://github.com/gin-gonic/gin)を Wails v3 に統合する方法を説明します。Gin は Go で記述された高性能な HTTP Web フレームワークで、Web アプリケーションや API を簡単に構築できます。

## はじめに

Wails v3 は柔軟なアセットシステムを備えており、Gin などの一般的な Web フレームワークを含む、任意の HTTP ハンドラーを使用できます。この統合により、次のことが可能になります。

- Gin のルーティングとミドルウェアを使用して Web コンテンツを配信する
- Wails アプリケーションからアクセスできる RESTful API を作成する
- Wails のデスクトップ統合を維持しながら Gin の機能を使用する

## Wails で Gin をセットアップする

Gin を Wails に統合するには、Gin ルーターを作成し、Wails アプリケーションのアセットハンドラーとして設定する必要があります。以下に手順を示します。

### 1. 依存関係をインストールする

まず、Gin パッケージがインストールされていることを確認します。

```bash
go get -u github.com/gin-gonic/gin
```

### 2. Gin 用のミドルウェアを作成する

Wails と Gin の統合を処理するミドルウェア関数を作成します。

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

このミドルウェアは、すべての HTTP リクエストを Gin ルーターに渡します。

### 3. Gin ルーターを設定する

ルート、ミドルウェア、ハンドラーを指定して Gin ルーターをセットアップします。

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

### 4. Wails アプリケーションに統合する

Gin ルーターをアセットハンドラーとして使用するように Wails アプリケーションを設定します。

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

## 静的コンテンツを配信する

Go の `embed` パッケージを使用して、静的ファイルをバイナリに埋め込みます。

```go
//go:embed static
var staticFiles embed.FS

ginEngine.StaticFS("/static", http.FS(staticFiles))
ginEngine.GET("/", func(c *gin.Context) {
    file, _ := staticFiles.ReadFile("static/index.html")
    c.Data(http.StatusOK, "text/html; charset=utf-8", file)
})
```

開発時には、`ginEngine.Static("/static", "./static")` を使用してディスク上のファイルを直接配信します。

## カスタムミドルウェア

Gin では、さまざまな用途に応じたカスタムミドルウェアを作成できます。以下はロギングミドルウェアの例です。

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

## API リクエストを処理する

Gin を使用すると、RESTful API を簡単に作成できます。API エンドポイントは次のように定義します。

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

## Wails の機能を使用する

Gin で配信する Web コンテンツは、`@wailsio/runtime` パッケージを使用して Wails の機能と連携できます。

### イベントを処理する

Go でイベントハンドラーを登録します。

```go
app.Event.On("my-event", func(event *application.CustomEvent) {
    log.Printf("Received event: %v", event.Data)
})
```

ランタイムを使用して JavaScript から呼び出します。

```html
<script type="module">
    import { Events } from '@wailsio/runtime';

    Events.Emit("my-event", { message: "Hello from frontend" });
    Events.On("response-event", (data) => console.log(data));
</script>
```

## 高度な設定

### Gin のモードをカスタマイズする

本番環境では Gin をリリースモードに設定します。

```go
gin.SetMode(gin.ReleaseMode)  // Use gin.DebugMode for development
ginEngine := gin.New()
```

### WebSocket を処理する

Gorilla WebSocket などのライブラリを使用して、WebSocket を Gin に統合できます。

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

## ベストプラクティス

- <strong>Go の embed パッケージを使用する：</strong>配布しやすくするため、静的ファイルをバイナリに埋め込みます。
- <strong>関心事を分離する：</strong>API ロジックと UI ロジックを分離します。
- <strong>エラー処理：</strong>Gin のルートとフロントエンドコードの両方に、適切なエラー処理を実装します。
- <strong>セキュリティ：</strong>特にユーザー入力を処理する際は、セキュリティ上の考慮事項に注意します。
- <strong>パフォーマンス：</strong>パフォーマンスを向上させるため、本番環境では Gin のリリースモードを使用します。
- <strong>テスト：</strong>Gin のテストユーティリティを使用して、Gin のルートに対するテストを作成します。

## まとめ

Gin と Wails を統合すると、Web 技術を使用したデスクトップアプリケーションを構築するための強力な組み合わせが得られます。Gin のパフォーマンスと豊富な機能は Wails のデスクトップ統合機能を補完するため、双方の長所を活用した高度なアプリケーションを作成できます。

詳細については、[Gin のドキュメント](https://github.com/gin-gonic/gin)および[Wails のドキュメント](https://wails.io)を参照してください。
