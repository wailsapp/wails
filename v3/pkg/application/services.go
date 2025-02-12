package application

import (
	"context"
	"reflect"
)

// Service wraps a bound type instance.
// The zero value of Service is invalid.
// Valid values may only be obtained by calling [NewService].
type Service struct {
	instance any
	options  ServiceOptions
}

// ServiceOptions provides optional parameters for calls to [NewService].
type ServiceOptions struct {
	// Name can be set to override the name of the service
	// for logging and debugging purposes.
	//
	// If empty, it will default
	// either to the value obtained through the [ServiceName] interface,
	// or to the type name.
	Name string

	// If the service instance implements [http.Handler],
	// it will be mounted on the internal asset server
	// at the prefix specified by Route.
	Route string

	// MarshalError will be called if non-nil
	// to marshal to JSON the error values returned by this service's methods.
	//
	// MarshalError is not allowed to fail,
	// but it may return a nil slice to fall back
	// to the globally configured error handler.
	//
	// If the returned slice is not nil, it must contain valid JSON.
	MarshalError func(error) []byte
}

// DefaultServiceOptions specifies the default values of service options,
// used when no [ServiceOptions] instance is provided to [NewService].
var DefaultServiceOptions = ServiceOptions{}

// NewService returns a Service value wrapping the given pointer.
// If T is not a concrete named type, the returned value is invalid.
func NewService[T any](instance *T) Service {
	return Service{instance, DefaultServiceOptions}
}

// NewServiceWithOptions returns a Service value wrapping the given pointer
// and specifying the given service options.
// If T is not a concrete named type, the returned value is invalid.
func NewServiceWithOptions[T any](instance *T, options ServiceOptions) Service {
	service := NewService(instance) // Delegate to NewService so that the static analyser may detect T. Do not remove this call.
	service.options = options
	return service
}

// Instance returns the service instance provided to [NewService].
func (s Service) Instance() any {
	return s.instance
}

// ServiceName returns the name of the service
//
// This is an *optional* method that may be implemented by service instances.
// It is used for logging and debugging purposes.
//
// If a non-empty name is provided with [ServiceOptions],
// it takes precedence over the one returned by the ServiceName method.
type ServiceName interface {
	ServiceName() string
}

// ServiceStartup is an *optional* method that may be implemented by service instances.
//
// This method will be called during application startup and will receive a copy of the options
// specified at creation time. It can be used for initialising resources.
//
// The context will be valid as long as the application is running,
// and will be canceled right before shutdown.
//
// Services are guaranteed to receive the startup notification
// in the exact order in which they were either
// listed in [Options.Services] or registered with [App.RegisterService],
// with those from [Options.Services] coming first.
//
// If the return value is non-nil, the startup process aborts
// and [App.Run] returns the error wrapped with [fmt.Errorf]
// in a user-friendly message comprising the service name.
// The original error can be retrieved either by calling the Unwrap method
// or through the [errors.As] API.
//
// When that happens, service instances that have been already initialised
// receive a shutdown notification.
type ServiceStartup interface {
	ServiceStartup(ctx context.Context, options ServiceOptions) error
}

// ServiceShutdown is an *optional* method that may be implemented by service instances.
//
// This method will be called during application shutdown. It can be used for cleaning up resources.
// If a service has received a startup notification,
// then it is guaranteed to receive a shutdown notification too,
// except in case of unhandled panics during shutdown.
//
// Services receive shutdown notifications in reverse registration order,
// after all user-provided shutdown hooks have run (see [App.OnShutdown]).
//
// If the return value is non-nil, it is passed to the application's
// configured error handler at [Options.ErrorHandler],
// wrapped with [fmt.Errorf] in a user-friendly message comprising the service name.
// The default behaviour is to log the error along with the service name.
// The original error can be retrieved either by calling the Unwrap method
// or through the [errors.As] API.
type ServiceShutdown interface {
	ServiceShutdown() error
}

func getServiceName(service Service) string {
	if service.options.Name != "" {
		return service.options.Name
	}

	// Check if the service implements the ServiceName interface
	if s, ok := service.Instance().(ServiceName); ok {
		return s.ServiceName()
	}

	// Finally, get the name from the type.
	return reflect.TypeOf(service.Instance()).Elem().String()
}
