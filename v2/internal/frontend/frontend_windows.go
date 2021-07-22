package frontend

import (
	"github.com/wailsapp/wails/v2/internal/frontend/windows"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func NewFrontend(appoptions *options.App, myLogger *logger.Logger) *windows.Frontend {
	return windows.NewFrontend(appoptions, myLogger)
}
