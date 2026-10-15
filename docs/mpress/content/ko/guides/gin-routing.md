---
title: "라우팅에 Gin 사용하기"
description: "Gin 웹 프레임워크를 Wails v3 애플리케이션과 통합하는 종합 가이드"
slug: "guides/gin-routing"
sourcePath: "guides/gin-routing.md"
---

이 가이드에서는 [Gin 웹 프레임워크](https://github.com/gin-gonic/gin)를 Wails v3와 통합하는 방법을 설명합니다. Gin은 Go로 작성된 고성능 HTTP 웹 프레임워크로, 웹 애플리케이션과 API를 쉽게 구축할 수 있게 해 줍니다.

## 소개

Wails v3는 Gin과 같은 널리 사용되는 웹 프레임워크를 비롯해 모든 HTTP 핸들러를 사용할 수 있는 유연한 애셋 시스템을 제공합니다. 이를 통합하면 다음 작업을 수행할 수 있습니다.

- Gin의 라우팅과 미들웨어를 사용하여 웹 콘텐츠 제공
- Wails 애플리케이션에서 접근할 수 있는 RESTful API 생성
- Wails 데스크톱 통합을 유지하면서 Gin의 기능 사용

## Wails에서 Gin 설정하기

Gin을 Wails와 통합하려면 Gin 라우터를 생성하고 Wails 애플리케이션의 애셋 핸들러로 구성해야 합니다. 다음 단계에 따라 진행하세요.

### 1. 종속성 설치

먼저 Gin 패키지가 설치되어 있는지 확인하세요.

```bash
go get -u github.com/gin-gonic/gin
```

### 2. Gin용 미들웨어 생성

Wails와 Gin 간의 통합을 처리할 미들웨어 함수를 생성하세요.

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

이 미들웨어는 모든 HTTP 요청을 Gin 라우터로 전달합니다.

### 3. Gin 라우터 구성

라우트, 미들웨어, 핸들러를 사용하여 Gin 라우터를 설정하세요.

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

### 4. Wails 애플리케이션과 통합

Gin 라우터를 애셋 핸들러로 사용하도록 Wails 애플리케이션을 구성하세요.

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

## 정적 콘텐츠 제공

Go의 `embed` 패키지를 사용하여 정적 파일을 바이너리에 임베드하세요.

```go
//go:embed static
var staticFiles embed.FS

ginEngine.StaticFS("/static", http.FS(staticFiles))
ginEngine.GET("/", func(c *gin.Context) {
    file, _ := staticFiles.ReadFile("static/index.html")
    c.Data(http.StatusOK, "text/html; charset=utf-8", file)
})
```

개발 중에는 `ginEngine.Static("/static", "./static")`을 사용하여 디스크에서 직접 파일을 제공하세요.

## 사용자 정의 미들웨어

Gin에서는 다양한 용도의 사용자 정의 미들웨어를 생성할 수 있습니다. 다음은 로깅 미들웨어의 예입니다.

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

## API 요청 처리

Gin을 사용하면 RESTful API를 쉽게 생성할 수 있습니다. API 엔드포인트를 정의하는 방법은 다음과 같습니다.

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

## Wails 기능 사용

Gin이 제공하는 웹 콘텐츠에서는 `@wailsio/runtime` 패키지를 사용하여 Wails 기능과 상호 작용할 수 있습니다.

### 이벤트 처리

Go에서 이벤트 핸들러를 등록하세요.

```go
app.Event.On("my-event", func(event *application.CustomEvent) {
    log.Printf("Received event: %v", event.Data)
})
```

런타임을 사용하여 JavaScript에서 호출하세요.

```html
<script type="module">
    import { Events } from '@wailsio/runtime';

    Events.Emit("my-event", { message: "Hello from frontend" });
    Events.On("response-event", (data) => console.log(data));
</script>
```

## 고급 구성

### Gin 모드 사용자 정의

프로덕션 환경에서는 Gin을 릴리스 모드로 설정하세요.

```go
gin.SetMode(gin.ReleaseMode)  // Use gin.DebugMode for development
ginEngine := gin.New()
```

### WebSocket 처리

Gorilla WebSocket과 같은 라이브러리를 사용하여 WebSocket을 Gin과 통합할 수 있습니다.

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

## 모범 사례

- **Go의 embed 패키지 사용:** 더 쉽게 배포할 수 있도록 정적 파일을 바이너리에 임베드하세요.
- **관심사 분리:** API 로직과 UI 로직을 분리하세요.
- **오류 처리:** Gin 라우트와 프런트엔드 코드 모두에서 적절한 오류 처리를 구현하세요.
- **보안:** 특히 사용자 입력을 처리할 때 보안 고려 사항에 유의하세요.
- **성능:** 프로덕션 환경에서는 성능 향상을 위해 Gin의 릴리스 모드를 사용하세요.
- **테스트:** Gin의 테스트 유틸리티를 사용하여 Gin 라우트에 대한 테스트를 작성하세요.

## 결론

Gin과 Wails를 통합하면 웹 기술로 데스크톱 애플리케이션을 구축할 때 강력한 조합을 활용할 수 있습니다. Gin의 성능과 기능은 Wails의 데스크톱 통합 기능을 보완하므로, 두 기술의 장점을 모두 활용하는 정교한 애플리케이션을 만들 수 있습니다.

자세한 내용은 [Gin 문서](https://github.com/gin-gonic/gin)와 [Wails 문서](https://wails.io)를 참조하세요.
