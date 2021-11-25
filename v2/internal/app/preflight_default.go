//go:build !windows
// +build !windows

package app

import "github.com/wailsapp/wails/v2/pkg/options"

func (a *App) PreflightChecks(options *options.App) error {
	return nil
}
