package options

import (
	"context"
	"embed"
	"log"

	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

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
	AlwaysOnTop       bool
	RGBA              *RGBA
	Assets            embed.FS
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
	Mac     *mac.Options
}

type RGBA struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
	A uint8 `json:"a"`
}

// MergeDefaults will set the minimum default values for an application
func MergeDefaults(appoptions *App) {
	err := mergo.Merge(appoptions, Default)
	if err != nil {
		log.Fatal(err)
	}

	// DEfault colour. Doesn't work well with mergo
	if appoptions.RGBA == nil {
		appoptions.RGBA = &RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		}
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
