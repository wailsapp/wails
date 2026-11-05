---
title: "서비스에 Gin 사용하기"
description: "Gin 웹 프레임워크를 Wails v3 서비스와 통합하는 방법을 설명하는 가이드"
slug: "guides/gin-services"
sourcePath: "guides/gin-services.md"
---

## 서비스에 Gin 사용하기

Gin 웹 프레임워크는 Go로 HTTP 서비스를 구축할 때 널리 사용됩니다. Wails v3에서는 Gin 기반 서비스를 애플리케이션에 손쉽게 통합하여 HTTP 요청을 처리하고, RESTful API를 구현하며, 웹 콘텐츠를 제공할 수 있습니다.

이 가이드에서는 Wails 애플리케이션의 특정 라우트에 마운트할 수 있는 Gin 기반 서비스를 만드는 방법을 단계별로 설명합니다. 다음 작업을 수행하는 완전한 예제를 만들어 보겠습니다.

1. Gin 기반 서비스 만들기
2. Wails Service 인터페이스 구현하기
3. 라우트와 미들웨어 설정하기
4. Wails 이벤트 시스템과 통합하기
5. 프런트엔드에서 서비스와 상호 작용하기

## 사전 요구 사항

시작하기 전에 다음 사항을 준비했는지 확인하세요.

- Wails v3 설치
- Go 및 Gin 프레임워크에 대한 기본 지식
- HTTP 개념 및 RESTful API에 대한 이해

프로젝트에 Gin 프레임워크를 추가해야 합니다.

```bash
go get github.com/gin-gonic/gin
```

## Gin 기반 서비스 만들기

먼저 Wails Service 인터페이스를 구현하는 Gin 서비스를 만들어 보겠습니다. 이 서비스는 사용자 컬렉션을 관리하고 사용자 레코드를 조회하고 생성하는 API 엔드포인트를 제공합니다.

### 1. 데이터 모델 정의하기

먼저 서비스에서 사용할 데이터 구조를 정의하세요.

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

### 2. 서비스 구조체 만들기

다음으로 Gin 라우터와 서비스에서 유지해야 하는 모든 상태를 담을 서비스 구조체를 정의하세요.

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

### 3. Service 인터페이스 구현하기

Wails Service 인터페이스에 필요한 메서드를 구현하세요.

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

### 3. http.Handler 인터페이스 구현하기

서비스를 특정 라우트에 마운트할 수 있도록 `http.Handler` 인터페이스를 구현하세요. 단일 메서드인 `ServeHTTP`은 서비스로 들어오는 모든 HTTP 요청의 진입점입니다. 이 메서드는 요청 처리를 Gin 라우터에 위임하므로 라우팅과 미들웨어를 위한 Gin의 모든 강력한 기능을 사용할 수 있습니다.

```go
// ServeHTTP implements the http.Handler interface
func (s *GinService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// All requests go to the Gin router
	s.ginEngine.ServeHTTP(w, r)
}
```

### 4. 라우트 설정하기

API 라우트를 별도의 메서드에 정의하여 체계적으로 구성하세요. 이 방식은 코드를 깔끔하게 유지하고 API 구조를 더 쉽게 이해할 수 있게 해 줍니다. Gin 라우터는 라우트 정의를 위한 플루언트 API를 제공하며, 관련 엔드포인트를 체계적으로 구성하는 데 유용한 라우트 그룹도 지원합니다.

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

### 5. 사용자 지정 미들웨어 만들기

사용자 지정 Gin 미들웨어를 만들어 서비스를 개선할 수 있습니다. Gin의 미들웨어 함수는 라우터에 추가된 순서대로 실행되며 로깅, 인증, 오류 처리 등의 작업을 수행할 수 있습니다. 이 예제에서는 요청 세부 정보를 기록하는 간단한 로깅 미들웨어를 보여 줍니다.

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

## 서비스 등록하기

Wails 애플리케이션에서 Gin 기반 서비스를 사용하려면 애플리케이션에 서비스를 등록하고 마운트할 라우트를 지정해야 합니다. 이 작업은 Wails 애플리케이션 인스턴스를 만들 때 수행합니다. 지정한 라우트는 Gin 라우터에 정의된 모든 엔드포인트의 기본 경로가 됩니다.

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

이 예제에서는 Gin 서비스를 `/api` 라우트에 마운트합니다. 즉, Gin 라우터에 `/info` 엔드포인트가 정의되어 있다면 애플리케이션에서 `/api/info` 경로로 접근할 수 있습니다. 이 방식을 사용하면 API 엔드포인트를 논리적으로 구성하고 애플리케이션의 다른 부분과 충돌하지 않도록 할 수 있습니다.

## Wails 이벤트 시스템과 통합하기

Wails 서비스와 함께 Gin을 사용할 때의 강력한 기능 중 하나는 Wails 이벤트 시스템과 원활하게 통합할 수 있다는 점입니다. 이를 통해 백엔드 서비스와 프런트엔드가 실시간으로 통신할 수 있습니다.

서비스의 `ServiceStartup` 메서드에서 프런트엔드의 이벤트를 처리할 이벤트 핸들러를 등록할 수 있습니다.

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

Gin 라우트나 서비스의 다른 부분에서 프런트엔드로 이벤트를 내보낼 수도 있습니다.

```go
// After creating a new user
s.app.Event.Emit("user-created", newUser)
```

## 프런트엔드 통합

프런트엔드에서 Gin 서비스와 상호 작용하려면 Wails 런타임을 가져오고, API 엔드포인트에 HTTP 요청을 보내며, 실시간 통신에 Wails 이벤트 시스템을 사용해야 합니다.

프로덕션 환경에서는 `/wails/runtime.js`을 직접 가져오는 대신 `@wailsio/runtime` 패키지를 사용하는 것이 좋습니다. 이렇게 하면 타입 안전성, 향상된 IDE 지원, npm을 통한 버전 관리, 최신 JavaScript 도구와의 호환성을 확보할 수 있습니다.

패키지를 설치하세요.

```bash
npm install @wailsio/runtime
```

그런 다음 코드에서 사용하세요.

```javascript
import * as wails from '@wailsio/runtime';

// Event emission
wails.Events.Emit('gin-api-event', eventData);
```

다음은 프런트엔드 통합을 설정하는 방법의 예입니다.

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

## 마무리

Gin 웹 프레임워크를 Wails v3 서비스와 통합하면 모듈화되고 유지 관리하기 쉬운 웹 애플리케이션을 강력하고 유연한 방식으로 구축할 수 있습니다. Gin의 라우팅 및 미들웨어 기능을 Wails 이벤트 시스템과 함께 활용하면 관심사를 명확하게 분리한 풍부한 대화형 애플리케이션을 만들 수 있습니다.

이 가이드의 전체 예제 코드는 Wails 저장소의 `v3/examples/gin-service`에서 확인할 수 있습니다.
