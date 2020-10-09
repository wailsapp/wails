package options

import (
	"log"

	"github.com/imdario/mergo"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

// App contains options for creating the App
type App struct {
	Title         string
	Width         int
	Height        int
	DisableResize bool
	Fullscreen    bool
	MinWidth      int
	MinHeight     int
	MaxWidth      int
	MaxHeight     int
	StartHidden   bool
	DevTools      bool
	RGBA          int
	Mac           *mac.Options
	Logger        logger.Logger
	LogLevel      logger.LogLevel
}

// MergeDefaults will set the minimum default values for an application
func (a *App) MergeDefaults() {
	err := mergo.Merge(a, Default)
	if err != nil {
		log.Fatal(err)
	}
}
