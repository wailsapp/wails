---
title: "在服務中使用 Gin"
description: "將 Gin Web 框架與 Wails v3 服務整合的指南"
slug: "guides/gin-services"
sourcePath: "guides/gin-services.md"
---

## 在服務中使用 Gin

Gin Web 框架是使用 Go 建置 HTTP 服務的熱門選擇。透過 Wails v3，您可以輕鬆將以 Gin 為基礎的服務整合至應用程式，並以強大且有效的方式處理 HTTP 要求、實作 RESTful API，以及提供 Web 內容。

本指南將逐步說明如何建立以 Gin 為基礎的服務，並將其掛載至 Wails 應用程式中的特定路由。我們將建置一個完整範例，示範如何：

1. 建立以 Gin 為基礎的服務
2. 實作 Wails Service 介面
3. 設定路由與中介軟體
4. 與 Wails 事件系統整合
5. 從前端與服務互動

## 先決條件

開始前，請確認您已具備：

- 已安裝 Wails v3
- 具備 Go 與 Gin 框架的基本知識
- 熟悉 HTTP 概念與 RESTful API

您需要將 Gin 框架新增至專案：

```bash
go get github.com/gin-gonic/gin
```

## 建立以 Gin 為基礎的服務

首先，建立一個實作 Wails Service 介面的 Gin 服務。此服務將管理使用者集合，並提供用於擷取及建立使用者記錄的 API 端點。

### 1. 定義資料模型

首先，定義服務將使用的資料結構：

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

### 2. 建立服務結構

接著，定義用來保存 Gin 路由器及服務所需維護之任何狀態的服務結構：

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

### 3. 實作 Service 介面

實作 Wails Service 介面所需的方法：

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

### 3. 實作 http.Handler 介面

若要讓服務可掛載至特定路由，請實作`http.Handler`介面。它唯一的方法`ServeHTTP`是所有傳送至服務之 HTTP 要求的進入點。此方法會將要求處理委派給 Gin 路由器，讓您能使用 Gin 所有強大的路由與中介軟體功能。

```go
// ServeHTTP implements the http.Handler interface
func (s *GinService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// All requests go to the Gin router
	s.ginEngine.ServeHTTP(w, r)
}
```

### 4. 設定路由

為了提升組織性，請在個別方法中定義 API 路由。這種做法能保持程式碼整潔，也更容易理解 API 的結構。Gin 路由器提供流暢式 API 來定義路由，並支援路由群組，以便組織相關端點。

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

### 5. 建立自訂中介軟體

您可以建立自訂 Gin 中介軟體來強化服務。Gin 中的中介軟體函式會依新增至路由器的順序執行，可執行記錄、驗證及錯誤處理等工作。此範例示範一個記錄要求詳細資料的簡易日誌中介軟體。

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

## 註冊服務

若要在 Wails 應用程式中使用以 Gin 為基礎的服務，您需要向應用程式註冊該服務，並指定其掛載路由。這項設定會在建立 Wails 應用程式執行個體時完成。您指定的路由將成為 Gin 路由器中所有已定義端點的基礎路徑。

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

在此範例中，Gin 服務會掛載至`/api`路由。這表示如果 Gin 路由器定義了`/info`端點，便可在應用程式中透過`/api/info`存取該端點。這種做法可讓您有條理地組織 API 端點，並避免與應用程式的其他部分發生衝突。

## 與 Wails 事件系統整合

搭配 Wails Services 使用 Gin 的強大功能之一，是能與 Wails 事件系統無縫整合。這可讓後端服務與前端即時通訊。

您可以在服務的`ServiceStartup`方法中註冊事件處理常式，以處理來自前端的事件：

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

您也可以從 Gin 路由或服務的其他部分向前端發出事件：

```go
// After creating a new user
s.app.Event.Emit("user-created", newUser)
```

## 前端整合

若要從前端與 Gin 服務互動，您需要匯入 Wails 執行階段、向 API 端點發出 HTTP 要求，並使用 Wails 事件系統進行即時通訊。

在正式環境中，建議使用`@wailsio/runtime`套件，而不要直接匯入`/wails/runtime.js`。這可確保型別安全、提供更完善的 IDE 支援、透過 npm 管理版本，並與現代 JavaScript 工具相容。

安裝套件：

```bash
npm install @wailsio/runtime
```

接著，在程式碼中使用它：

```javascript
import * as wails from '@wailsio/runtime';

// Event emission
wails.Events.Emit('gin-api-event', eventData);
```

以下範例示範如何設定前端整合：

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

## 結語

將 Gin Web 框架與 Wails v3 Services 整合，能以強大且靈活的方式建置模組化、易於維護的 Web 應用程式。結合 Gin 的路由及中介軟體功能與 Wails 事件系統，您可以建立功能豐富、互動性高且關注點明確分離的應用程式。

本指南的完整範例程式碼位於 Wails 儲存庫的`v3/examples/gin-service`下。
