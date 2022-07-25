//go:build hybrid
// +build hybrid

package hybrid

import (
	"context"

	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop"
	"github.com/wailsapp/wails/v2/internal/frontend/devserver"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

// New returns a new Hybrid frontend
// A hybrid Frontend implementation contains a devserver.Frontend wrapping a desktop.Frontend
func NewFrontend(ctx context.Context, appoptions *options.App, logger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher) frontend.Frontend {
	appFrontend := desktop.NewFrontend(ctx, appoptions, logger, appBindings, dispatcher)
	return devserver.NewFrontend(ctx, appoptions, logger, appBindings, dispatcher, nil, appFrontend)
}
