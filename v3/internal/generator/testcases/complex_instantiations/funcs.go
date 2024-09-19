package main

import "github.com/wailsapp/wails/v3/pkg/application"

func ServiceInitialiser[T any]() func(*T, ...application.ServiceOptions) application.Service {
	return application.NewService[T]
}

func CustomNewServices[T any, U any]() []application.Service {
	return []application.Service{
		application.NewService(new(T)),
		application.NewService(new(U)),
	}
}
