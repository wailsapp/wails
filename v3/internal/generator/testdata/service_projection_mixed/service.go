package serviceprojectionmixed

import (
	"github.com/wailsapp/wails/v3/internal/generator/testdata/service_projection/backend"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type FrontendAPI interface {
	Echo(string) string
}

func RegisterProjected() application.Service {
	return application.NewServiceAs[FrontendAPI](&backend.Service{})
}

func RegisterUnprojected() application.Service {
	return application.NewService(&backend.Service{})
}
