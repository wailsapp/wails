package application

import (
	"context"
	"reflect"
)

type ServiceName interface {
	ServiceName() string
}

type ServiceStartup interface {
	ServiceStartup(ctx context.Context, options ServiceOptions) error
}

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
