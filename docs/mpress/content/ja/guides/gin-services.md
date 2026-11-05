---
title: "サービスで Gin を使用する"
description: "Gin Web フレームワークを Wails v3 Services と統合するためのガイド"
slug: "guides/gin-services"
sourcePath: "guides/gin-services.md"
---

## サービスで Gin を使用する

Gin Web フレームワークは、Go で HTTP サービスを構築する際によく使われています。Wails v3 では、Gin ベースのサービスをアプリケーションに簡単に統合でき、HTTP リクエストの処理、RESTful API の実装、Web コンテンツの配信を強力に支援します。

このガイドでは、Wails アプリケーション内の特定のルートにマウントできる Gin ベースのサービスを作成します。以下の方法を示す完全なサンプルを構築します。

1. Gin ベースのサービスを作成する
2. Wails Service インターフェースを実装する
3. ルートとミドルウェアを設定する
4. Wails イベントシステムと統合する
5. フロントエンドからサービスを操作する

## 前提条件

始める前に、以下を用意してください。

- Wails v3 がインストールされていること
- Go と Gin フレームワークの基礎知識
- HTTP の概念と RESTful API に関する知識

プロジェクトに Gin フレームワークを追加する必要があります。

```bash
go get github.com/gin-gonic/gin
```

## Gin ベースのサービスを作成する

まず、Wails Service インターフェースを実装する Gin サービスを作成します。このサービスではユーザーのコレクションを管理し、ユーザーレコードを取得および作成するための API エンドポイントを提供します。

### 1. データモデルを定義する

最初に、サービスで使用するデータ構造を定義します。

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

### 2. サービスの構造体を作成する

次に、Gin ルーターと、サービスで保持する必要があるすべての状態を格納するサービス構造体を定義します。

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

### 3. Service インターフェースを実装する

Wails Service インターフェースに必要なメソッドを実装します。

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

### 3. http.Handler インターフェースを実装する

サービスを特定のルートにマウントできるように、`http.Handler` インターフェースを実装します。このインターフェースの唯一のメソッドである `ServeHTTP` は、サービスに対するすべての HTTP リクエストの入口です。このメソッドはリクエスト処理を Gin ルーターに委譲するため、ルーティングやミドルウェアに関する Gin の強力な機能をすべて利用できます。

```go
// ServeHTTP implements the http.Handler interface
func (s *GinService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// All requests go to the Gin router
	s.ginEngine.ServeHTTP(w, r)
}
```

### 4. ルートを設定する

整理しやすいように、API ルートは別のメソッドで定義します。この方法によりコードを簡潔に保ち、API の構造を把握しやすくなります。Gin ルーターは、関連するエンドポイントの整理に役立つルートグループのサポートなど、ルートを定義するための流れるような API を提供します。

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

### 5. カスタムミドルウェアを作成する

カスタム Gin ミドルウェアを作成して、サービスを拡張できます。Gin のミドルウェア関数はルーターに追加された順に実行され、ログ記録、認証、エラー処理などを実行できます。この例では、リクエストの詳細を記録するシンプルなロギングミドルウェアを示します。

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

## サービスを登録する

Gin ベースのサービスを Wails アプリケーションで使用するには、そのサービスをアプリケーションに登録し、マウント先のルートを指定する必要があります。これは Wails アプリケーションのインスタンスを作成するときに行います。指定したルートは、Gin ルーターで定義したすべてのエンドポイントのベースパスになります。

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

この例では、Gin サービスを `/api` ルートにマウントします。つまり、Gin ルーターに `/info` として定義されたエンドポイントがある場合、アプリケーションでは `/api/info` でアクセスできます。この方法により、API エンドポイントを論理的に整理し、アプリケーションのほかの部分との競合を回避できます。

## Wails イベントシステムと統合する

Gin を Wails Services とともに使用する強力な機能の一つは、Wails イベントシステムとシームレスに統合できることです。これにより、バックエンドサービスとフロントエンドの間でリアルタイム通信が可能になります。

サービスの `ServiceStartup` メソッドで、フロントエンドからのイベントを処理するイベントハンドラーを登録できます。

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

Gin のルートやサービスのほかの部分から、フロントエンドにイベントを送出することもできます。

```go
// After creating a new user
s.app.Event.Emit("user-created", newUser)
```

## フロントエンドとの統合

フロントエンドから Gin サービスを操作するには、Wails ランタイムをインポートし、API エンドポイントに HTTP リクエストを送信し、リアルタイム通信に Wails イベントシステムを使用します。

本番環境では、`/wails/runtime.js` を直接インポートする代わりに、`@wailsio/runtime` パッケージを使用することを推奨します。これにより、型安全性、IDE サポートの向上、npm によるバージョン管理、および最新の JavaScript ツールとの互換性が確保されます。

パッケージをインストールします。

```bash
npm install @wailsio/runtime
```

次に、コード内で使用します。

```javascript
import * as wails from '@wailsio/runtime';

// Event emission
wails.Events.Emit('gin-api-event', eventData);
```

フロントエンドとの統合を設定する例を以下に示します。

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

## まとめ

Gin Web フレームワークを Wails v3 Services と統合すると、モジュール化され、保守しやすい Web アプリケーションを構築するための強力で柔軟な手法を利用できます。Gin のルーティング機能およびミドルウェア機能を Wails イベントシステムと組み合わせることで、関心事を明確に分離した、機能豊富でインタラクティブなアプリケーションを作成できます。

このガイドの完全なサンプルコードは、Wails リポジトリ内の `v3/examples/gin-service` にあります。
