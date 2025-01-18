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
}

// DefaultServiceOptions specifies the default values of service options,
// used when no [ServiceOptions] instance is provided to [NewService].
var DefaultServiceOptions = ServiceOptions{
	Route: "",
}

// NewService returns a Service value wrapping the given pointer.
// If T is not a concrete named type, the returned value is invalid.
//
// It accepts optionally at most one instance of [ServiceOptions].
func NewService[T any](instance *T, options ...ServiceOptions) Service {
	switch len(options) {
	case 0:
		return Service{instance, DefaultServiceOptions}
	case 1:
		return Service{instance, options[0]}
	default:
		panic("NewService accepts at most one instance of ServiceOptions")
	}
}

// Instance returns the service instance provided to [NewService].
func (s Service) Instance() any {
	return s.instance
}

// Service instances that implement ServiceName
// may specify a custom name
// for logging and debugging purposes.
//
// If a non-empty name is provided with [ServiceOptions],
// it takes precedence over the one returned by the ServiceName method.
type ServiceName interface {
	ServiceName() string
}

// Service instances that implement ServiceStartup
// will be notified at application startup
// and receive a copy of the options specified at creation time.
//
// The context will be valid as long as the application is running,
// and will be canceled right before shutdown.
//
// If the return value is non-nil, it is logged along with the service name,
// the startup process aborts and the application quits.
// When that happens, service instances that have been already initialised
// receive a shutdown notification.
type ServiceStartup interface {
	ServiceStartup(ctx context.Context, options ServiceOptions) error
}

// Service instances that implement ServiceShutdown
// will be notified at application shutdown.
// If the return value is non-nil, it is logged along with the service name.
type ServiceShutdown interface {
	ServiceShutdown() error
}

func getServiceName(service any) string {
	// First check it conforms to ServiceName interface
	if serviceName, ok := service.(ServiceName); ok {
		return serviceName.ServiceName()
	}
	// Next, get the name from the type
	return reflect.TypeOf(service).String()
}
