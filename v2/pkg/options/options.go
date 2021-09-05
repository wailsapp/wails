package options

import (
	"context"
	"embed"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"log"

	"github.com/wailsapp/wails/v2/pkg/menu"

	"github.com/imdario/mergo"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

// App contains options for creating the App
type App struct {
	Title             string
	Width             int
	Height            int
	DisableResize     bool
	Fullscreen        bool
	Frameless         bool
	MinWidth          int
	MinHeight         int
	MaxWidth          int
	MaxHeight         int
	StartHidden       bool
	HideWindowOnClose bool
	RGBA              int
	Assets            *embed.FS
	Menu              *menu.Menu
	Logger            logger.Logger `json:"-"`
	LogLevel          logger.LogLevel
	OnStartup         func(ctx context.Context) `json:"-"`
	OnDomReady        func(ctx context.Context) `json:"-"`
	OnShutdown        func(ctx context.Context) `json:"-"`
	Bind              []interface{}

	//ContextMenus []*menu.ContextMenu
	//TrayMenus    []*menu.TrayMenu
	Windows *windows.Options
	//Mac     *mac.Options
}

// MergeDefaults will set the minimum default values for an application
func MergeDefaults(appoptions *App) {
	err := mergo.Merge(appoptions, Default)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure max and min are valid
	if appoptions.MinWidth > 0 && appoptions.MaxWidth > 0 {
		if appoptions.MinWidth > appoptions.MaxWidth {
			appoptions.MinWidth = appoptions.MaxWidth
		}
	}
	if appoptions.MinHeight > 0 && appoptions.MaxHeight > 0 {
		if appoptions.MinHeight > appoptions.MaxHeight {
			appoptions.MinHeight = appoptions.MaxHeight
		}
	}
	// Ensure width and height are limited if max/min is set
	if appoptions.Width < appoptions.MinWidth {
		appoptions.Width = appoptions.MinWidth
	}
	if appoptions.MaxWidth > 0 && appoptions.Width > appoptions.MaxWidth {
		appoptions.Width = appoptions.MaxWidth
	}
	if appoptions.Height < appoptions.MinHeight {
		appoptions.Height = appoptions.MinHeight
	}
	if appoptions.MaxHeight > 0 && appoptions.Height > appoptions.MaxHeight {
		appoptions.Height = appoptions.MaxHeight
	}

}
