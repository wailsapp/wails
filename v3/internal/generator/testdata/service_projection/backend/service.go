package backend

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct{}

func (*Service) Echo(value string) string {
	return value
}

func (*Service) Version() string {
	return "v1"
}

func (*Service) BackendOnly() string {
	return "secret"
}

func (*Service) ServiceStartup(context.Context, application.ServiceOptions) error {
	return nil
}
