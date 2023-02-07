//go:build linux
// +build linux

package desktop

import (
	"context"
	"github.com/ciderapp/wails/v2/internal/binding"
	"github.com/ciderapp/wails/v2/internal/frontend"
	"github.com/ciderapp/wails/v2/internal/frontend/desktop/linux"
	"github.com/ciderapp/wails/v2/internal/logger"
	"github.com/ciderapp/wails/v2/pkg/options"
)

func NewFrontend(ctx context.Context, appoptions *options.App, logger *logger.Logger, appBindings *binding.Bindings, dispatcher frontend.Dispatcher) frontend.Frontend {
	return linux.NewFrontend(ctx, appoptions, logger, appBindings, dispatcher)
}
