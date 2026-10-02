package serviceprojectionconflict

import (
	"github.com/wailsapp/wails/v3/internal/generator/testdata/service_projection/backend"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type EchoAPI interface {
	Echo(string) string
}

type VersionAPI interface {
	Version() string
}

func RegisterEcho() application.Service {
	return application.NewServiceAs[EchoAPI](&backend.Service{})
}

func RegisterVersion() application.Service {
	return application.NewServiceAs[VersionAPI](&backend.Service{})
}
