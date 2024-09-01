# Services

Services in Wails v3 provide a powerful way to extend the functionality of your application. They allow you to create 
modular, reusable components that can be easily integrated into your Wails application.

## Overview

Services are designed to encapsulate specific functionality and can be registered with the application at startup. 
They can handle various tasks such as file serving, database operations, logging, and more. 
Services can also interact with the application lifecycle and respond to HTTP requests.

## Creating a Service

To create a service, you simply define a struct. Here's a basic structure of a service:

```go
type MyService struct {
    // Your service fields
}

func NewMyService() *MyService {
    // Initialize and return your service
}

func (s *MyService) Greet(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}
```

This service has a single method, `Greet`, which accepts a name and returns a greeting.

## Registering a Service

To register a service with the application, you need to provide an instance of the service to the `Services` field of
the `application.Options` struct (All services need to be wrapped by an `application.NewService` call. Here's an example:

```go
app := application.New(application.Options{
    Services: []application.Service{
        application.NewService(NewMyService()),
    },
})

```

## Optional Methods

Services can implement optional methods to hook into the application lifecycle:

### Name

```go
func (s *Service) Name() string
```

This method returns the name of the service. It is used for logging purposes only.

### OnStartup

```go
func (s *Service) OnStartup(ctx context.Context, options application.ServiceOptions) error
```

This method is called when the application is starting up. You can use it to initialize resources, set up connections, 
or perform any necessary setup tasks. The context is the application context, and the `options` parameter provides
additional information about the service.

### OnShutdown

```go
func (s *Service) OnShutdown() error
```

This method is called when the application is shutting down. Use it to clean up resources, close connections, or 
perform any necessary cleanup tasks.

### ServeHTTP

```go
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

If your service needs to handle HTTP requests, implement this method. It allows your service to act as an HTTP handler.
The route of the handler is defined in the service options:

```go
    application.NewService(fileserver.New(&fileserver.Config{
        RootPath: rootPath,
    }), application.ServiceOptions{
        Route: "/files",
    }),
```

## Example: File Server Service

Let's look at a simplified version of the `fileserver` service as an example:

```go
type Service struct {
    config *Config
    fs     http.Handler
}

func New(config *Config) *Service {
    return &Service{
        config: config,
        fs:     http.FileServer(http.Dir(config.RootPath)),
    }
}

func (s *Service) Name() string {
    return "github.com/wailsapp/wails/v3/services/fileserver"
}

func (s *Service) OnStartup(ctx context.Context, options application.ServiceOptions) error {
    // Any initialization code here
    return nil
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.fs.ServeHTTP(w, r)
}
```

We can now use this service in our application:

```go
app := application.New(application.Options{
    Services: []application.Service{
      application.NewService(fileserver.New(&fileserver.Config{
      RootPath: rootPath,
    }), application.ServiceOptions{
      Route: "/files",
    }),
```
All requests to `/files` will be handled by the `fileserver` service.

## Application Lifecycle and Services

1. During application initialization, services are registered with the application.
2. When the application starts (`app.Run()`), the `OnStartup` method of each service is called with the application
   context and service options.
3. Throughout the application's lifetime, services can perform their specific tasks.
4. If a service implements `ServeHTTP`, it can handle HTTP requests at the specified path.
5. When the application is shutting down, the `OnShutdown` method of each service is called as well as the context being cancelled.
