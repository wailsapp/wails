package options

import (
	"context"
	"io/fs"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"github.com/wailsapp/wails/v2/pkg/menu"

	"github.com/imdario/mergo"
	"github.com/wailsapp/wails/v2/pkg/logger"
)

type WindowStartState int

const (
	Normal     WindowStartState = 0
	Maximised  WindowStartState = 1
	Minimised  WindowStartState = 2
	Fullscreen WindowStartState = 3
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
	Assets            fs.FS
	Menu              *menu.Menu
	Logger            logger.Logger `json:"-"`
	LogLevel          logger.LogLevel
	OnStartup         func(ctx context.Context)                `json:"-"`
	OnDomReady        func(ctx context.Context)                `json:"-"`
	OnShutdown        func(ctx context.Context)                `json:"-"`
	OnBeforeClose     func(ctx context.Context) (prevent bool) `json:"-"`
	Bind              []interface{}
	WindowStartState  WindowStartState

	//ContextMenus []*menu.ContextMenu
	//TrayMenus    []*menu.TrayMenu
	Windows *windows.Options
	Mac     *mac.Options
	Linux   *linux.Options
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

	switch runtime.GOOS {
	case "darwin":
		if appoptions.Menu == nil {
			appoptions.Menu = defaultMacMenu
		}
	}

}
