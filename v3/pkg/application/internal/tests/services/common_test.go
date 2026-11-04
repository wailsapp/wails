package services

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func TestMatchingOptionsWithoutReporter(t *testing.T) {
	options := application.ServiceOptions{Name: "service", Route: "/service"}
	// Some lifecycle fixtures omit Config.T when no assertion failure is expected.
	for _, service := range []interface {
		Configure(Config)
		ServiceStartup(context.Context, application.ServiceOptions) error
	}{&Startupper{}, &StartupShutdowner{}} {
		service.Configure(Config{Seq: &atomic.Int64{}, Options: options})
		if err := service.ServiceStartup(context.Background(), options); err != nil {
			t.Fatal(err)
		}
	}
}
