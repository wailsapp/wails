package serviceprojectionunsupported

import (
	"github.com/wailsapp/wails/v3/internal/generator/testdata/service_projection/backend"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type FrontendAPI interface {
	Echo(string) string
}

func RegisterInterfaceTyped() application.Service {
	var instance FrontendAPI = &backend.Service{}
	return application.NewServiceAs[FrontendAPI](instance)
}

func RegisterGeneric[T FrontendAPI](instance T) application.Service {
	return application.NewServiceAs[FrontendAPI](instance)
}

func Examples() (application.Service, application.Service) {
	instance := &backend.Service{}
	return RegisterInterfaceTyped(), RegisterGeneric(instance)
}
