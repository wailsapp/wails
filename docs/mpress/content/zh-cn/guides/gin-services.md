---
title: "在服务中使用 Gin"
description: "将 Gin Web 框架与 Wails v3 服务集成的指南"
slug: "guides/gin-services"
sourcePath: "guides/gin-services.md"
---

## 在服务中使用 Gin

Gin Web 框架是使用 Go 构建 HTTP 服务时的热门选择。借助 Wails v3，你可以轻松地将基于 Gin 的服务集成到应用程序中，从而以功能强大的方式处理 HTTP 请求、实现 RESTful API 并提供 Web 内容。

本指南将逐步介绍如何创建基于 Gin 的服务，并将其挂载到 Wails 应用程序中的特定路由。我们将构建一个完整示例，演示如何：

1. 创建基于 Gin 的服务
2. 实现 Wails Service 接口
3. 设置路由和中间件
4. 与 Wails 事件系统集成
5. 从前端与服务交互

## 前提条件

开始之前，请确保已具备以下条件：

- 已安装 Wails v3
- 具备 Go 和 Gin 框架的基础知识
- 熟悉 HTTP 概念和 RESTful API

你需要将 Gin 框架添加到项目中：

```bash
go get github.com/gin-gonic/gin
```

## 创建基于 Gin 的服务

首先，创建一个实现 Wails Service 接口的 Gin 服务。该服务将管理用户集合，并提供用于检索和创建用户记录的 API 端点。

### 1. 定义数据模型

首先，定义服务要使用的数据结构：

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

### 2. 创建服务结构体

接下来，定义服务结构体，用于保存 Gin 路由器以及服务需要维护的任何状态：

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

### 3. 实现 Service 接口

实现 Wails Service 接口所需的方法：

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

### 3. 实现 http.Handler 接口

要使服务可以挂载到特定路由，请实现`http.Handler`接口。该接口唯一的方法`ServeHTTP`是所有发往服务的 HTTP 请求的入口。它将请求处理委托给 Gin 路由器，使你能够使用 Gin 在路由和中间件方面的全部强大功能。

```go
// ServeHTTP implements the http.Handler interface
func (s *GinService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// All requests go to the Gin router
	s.ginEngine.ServeHTTP(w, r)
}
```

### 4. 设置路由

为改善代码组织，请在单独的方法中定义 API 路由。这种方式可以保持代码整洁，并使 API 的结构更易理解。Gin 路由器提供流式 API 来定义路由，还支持路由组，可用于组织相关端点。

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

### 5. 创建自定义中间件

你可以创建自定义 Gin 中间件来增强服务。Gin 中的中间件函数会按其添加到路由器的顺序执行，可完成日志记录、身份验证和错误处理等任务。此示例展示了一个记录请求详细信息的简单日志中间件。

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

## 注册服务

要在 Wails 应用程序中使用基于 Gin 的服务，需要向应用程序注册该服务，并指定其挂载路由。这项设置在创建 Wails 应用程序实例时完成。指定的路由将成为 Gin 路由器中定义的所有端点的基础路径。

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

在此示例中，Gin 服务挂载在`/api`路由上。这意味着，如果 Gin 路由器中定义了`/info`端点，那么可在应用程序中通过`/api/info`访问它。这种方式让你能够合理组织 API 端点，并避免与应用程序的其他部分发生冲突。

## 与 Wails 事件系统集成

将 Gin 与 Wails 服务结合使用的一项强大功能，是能够与 Wails 事件系统无缝集成。这样便可在后端服务与前端之间进行实时通信。

可以在服务的`ServiceStartup`方法中注册事件处理程序，以处理来自前端的事件：

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

还可以从 Gin 路由或服务的其他部分向前端发出事件：

```go
// After creating a new user
s.app.Event.Emit("user-created", newUser)
```

## 前端集成

要从前端与 Gin 服务交互，需要导入 Wails 运行时、向 API 端点发送 HTTP 请求，并使用 Wails 事件系统进行实时通信。

对于生产环境，建议使用`@wailsio/runtime`包，而不是直接导入`/wails/runtime.js`。这样可以确保类型安全、提供更好的 IDE 支持、通过 npm 进行版本管理，并与现代 JavaScript 工具兼容。

安装该包：

```bash
npm install @wailsio/runtime
```

然后在代码中使用它：

```javascript
import * as wails from '@wailsio/runtime';

// Event emission
wails.Events.Emit('gin-api-event', eventData);
```

以下示例演示了如何设置前端集成：

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

## 结语

将 Gin Web 框架与 Wails v3 服务集成，为构建模块化、易维护的 Web 应用程序提供了一种强大而灵活的方法。结合使用 Gin 的路由和中间件功能与 Wails 事件系统，你可以创建关注点清晰分离、功能丰富且交互性强的应用程序。

本指南的完整示例代码可在 Wails 仓库的`v3/examples/gin-service`下找到。
