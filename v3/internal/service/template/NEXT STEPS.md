# Next Steps

Congratulations on generating your new service. This guide will help you author your service
and provide some tips on how to get started.

## Service Structure

The service is a standard Go module that provides functionality to your application.
It can provide Go methods, that may be called from the frontend. It can also provide
assets to the frontend.

## Service Directory Structure

The service directory structure is as follows:

```
service-name
├── service.go
├── README.md
├── go.mod
├── go.sum
└── service.yml
```

### `service.go`

This file contains the service code. It should contain a struct that implements the `Service` interface
and a `NewService()` method that returns a pointer to the struct. Methods are exported by capitalising
the first letter of the method name. These methods may be called from the frontend. If methods
accept or return structs, these structs must be exported. 

### `service.yml`

This file contains the service metadata. It is important to fill this out correctly
as it will be used by the Wails CLI.

### `README.md`

This file should contain a description of the service and how to use it. It should
also contain a link to the service repository and how to report bugs.

### `go.mod` and `go.sum`

These are standard Go module files. The package name in `go.mod` should match the
name of the service, e.g. `github.com/myuser/wails-service-example`.

## Service Lifecycle Methods

During the service lifecycle, there are two methods that are called, *if* they are implemented:

### `OnStartup() error`

The `OnStartup() error` method is called when the service is loaded. This method should return an error if it fails to initialise.
This method is called synchronously so the application will not start until it returns.

### `OnShutdown()`

The `OnShutdown()` method is called when the application is shutting down. This is a good place to
perform any cleanup. This method is called synchronously so the application will not exit completely until
it returns.

## Service Assets

The service can provide assets to the frontend. These assets are available to the frontend via a configurable path.

To provide assets, simply provide an http handler:

```go
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Serve your assets here
    // e.g. http.FileServer(http.Dir("assets"))
}
```

The path this is mounted on is configurable when the service is registered:

```go
Services: []application.Service{
    application.NewService(&{{.Name}}{}, "/services/myservice"),
},
```

## Promoting your Service

Once you have created your service, you should promote it on the Wails Discord server
in the `#services` channel. 
You should also open a PR to promote your service on the Wails website. 