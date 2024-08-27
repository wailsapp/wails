package application

import "reflect"

type ServiceName interface {
	Name() string
}

type ServiceStartup interface {
	OnStartup() error
}

type ServiceShutdown interface {
	OnShutdown() error
}

func getServiceName(service any) string {
	// First check it conforms to ServiceName interface
	if serviceName, ok := service.(ServiceName); ok {
		return serviceName.Name()
	}
	// Next, get the name from the type
	return reflect.TypeOf(service).String()
}
