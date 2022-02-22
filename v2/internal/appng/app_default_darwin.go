//go:build !dev && !production && !bindings && darwin
// +build !dev,!production,!bindings,darwin

package appng

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/options"
)

// App defines a Wails application structure
type App struct{}

func (a *App) Run() error {
	return nil
}

// CreateApp creates the app!
func CreateApp(_ *options.App) (*App, error) {
	return nil, fmt.Errorf(`Wails applications will not build without the correct build tags.`)
}
