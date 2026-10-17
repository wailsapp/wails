package serviceprojection

import (
	"github.com/wailsapp/wails/v3/internal/generator/testdata/service_projection/backend"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type StatusAPI interface {
	Version() string
}

type FrontendAPI interface {
	StatusAPI
	Echo(value string) string
}

func Register() application.Service {
	return application.NewServiceAsWithOptions[FrontendAPI](&backend.Service{}, application.ServiceOptions{
		Name: "projected service",
	})
}

func RegisterDefault() application.Service {
	return application.NewServiceAs[FrontendAPI](&backend.Service{})
}
