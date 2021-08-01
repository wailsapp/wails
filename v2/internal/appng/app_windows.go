package appng

import (
	"github.com/wailsapp/wails/v2/internal/binding"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/windows"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func NewFrontend(appoptions *options.App, myLogger *logger.Logger, bindings *binding.Bindings, dispatcher frontend.Dispatcher) frontend.Frontend {
	return windows.NewFrontend(appoptions, myLogger, bindings, dispatcher)
}
