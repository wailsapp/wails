# Wails v3 Service Template

This README provides an overview of the Wails v3 service template and explains how to adapt it to create your own custom service.

## Overview

The service template provides a basic structure for creating a Wails v3 service. A service in Wails v3 is a Go package that can be integrated into your Wails application to provide specific functionality, handle HTTP requests, and interact with the frontend.

## Template Structure

The template defines a `MyService` struct and several methods:

### MyService Struct

```go
type MyService struct {
    ctx context.Context
    options application.ServiceOptions
}
```

This is the main service struct. You can rename it to better reflect your service's purpose. The struct holds a context and service options, which are set during startup.

### Name Method

```go
func (p *MyService) Name() string
```

This method returns the name of the service. It's used to identify the service within the Wails application.

### OnStartup Method

```go
func (p *MyService) OnStartup(ctx context.Context, options application.ServiceOptions) error
```

This method is called when the app is starting up. Use it to initialize resources, set up connections, or perform any necessary setup tasks. 
It receives a context and service options, which are stored in the service struct.

### OnShutdown Method

```go
func (p *MyService) OnShutdown() error
```

This method is called when the app is shutting down. Use it to clean up resources, close connections, or perform any necessary cleanup tasks.

### ServeHTTP Method

```go
func (p *MyService) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

This method handles HTTP requests to the service. It's called when the frontend makes an HTTP request to the backend 
at the path specified in the `Route` field of the service options.

### Service Methods

```go
func (p *MyService) Greet(name string) string
```

This is an example of a service method. You can add as many methods as you need. These methods can be called from the frontend.

## Adapting the Template

To create your own service:

1. Rename the `MyService` struct to reflect your service's purpose (e.g., `DatabaseService`, `AuthService`).
2. Update the `Name` method to return your service's unique identifier.
3. Implement the `OnStartup` method to initialize your service. This might include setting up database connections, loading configuration, etc.
4. If needed, implement the `OnShutdown` method to properly clean up resources when the application closes.
5. If your service needs to handle HTTP requests, implement the `ServeHTTP` method. Use this to create API endpoints, serve files, or handle any HTTP interactions.
6. Add your own methods to the service. These can include database operations, business logic, or any functionality your service needs to provide.
7. If your service requires configuration, consider adding a `Config` struct and a `New` function to create and configure your service.

## Example: Database Service

Here's how you might adapt the template for a database service:

```go
type DatabaseService struct {
    ctx context.Context
    options application.ServiceOptions
    db *sql.DB
}

func (s *DatabaseService) Name() string {
    return "github.com/myname/DatabaseService"
}

func (s *DatabaseService) OnStartup(ctx context.Context, options application.ServiceOptions) error {
    s.ctx = ctx
    s.options = options
    // Initialize database connection
    var err error
    s.db, err = sql.Open("mysql", "user:password@/dbname")
    return err
}

func (s *DatabaseService) OnShutdown() error {
    return s.db.Close()
}

func (s *DatabaseService) GetUser(id int) (User, error) {
    // Implement database query
}

// Add more methods as needed
```

## Long-running tasks

If your service needs to perform long-running tasks, consider using goroutines and channels to manage these tasks.
You can use the `context.Context` to listen for when the application shuts down:

```go
func (s *DatabaseService) longRunningTask() {
    for {
        select {
        case <-s.ctx.Done():
            // Cleanup and exit
            return
        // Perform long-running task
        }   
    }
}
```
