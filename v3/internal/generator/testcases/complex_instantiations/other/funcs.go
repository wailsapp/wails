package other

import "github.com/wailsapp/wails/v3/pkg/application"

func CustomNewService[T any](srv T) application.Service {
	return application.NewService(&srv)
}

func ServiceInitialiser[T any]() func(*T) application.Service {
	return application.NewService[T]
}

func CustomNewServices[T any, U any]() []application.Service {
	return []application.Service{
		application.NewService(new(T)),
		application.NewService(new(U)),
	}
}
