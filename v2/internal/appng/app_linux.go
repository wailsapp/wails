//go:build linux && !bindings
// +build linux,!bindings

package appng

import (
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
)

func PreflightChecks(options *options.App, logger *logger.Logger) error {

	_ = options

	return nil
}
