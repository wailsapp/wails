package options

import (
	"log"
	"runtime"

	wailsruntime "github.com/wailsapp/wails/v2/internal/runtime"
	"github.com/wailsapp/wails/v2/pkg/menu"

	"github.com/imdario/mergo"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

// App contains options for creating the App
type App struct {
	Title             string
	Width             int
	Height            int
	DisableResize     bool
	Fullscreen        bool
	MinWidth          int
	MinHeight         int
	MaxWidth          int
	MaxHeight         int
	StartHidden       bool
	HideWindowOnClose bool
	DevTools          bool
	RGBA              int
	ContextMenus      []*menu.ContextMenu
	TrayMenus         []*menu.TrayMenu
	Menu              *menu.Menu
	Mac               *mac.Options
	Logger            logger.Logger `json:"-"`
	LogLevel          logger.LogLevel
	Startup           func(*wailsruntime.Runtime) `json:"-"`
	Shutdown          func()                      `json:"-"`
	Bind              []interface{}
}

// MergeDefaults will set the minimum default values for an application
func MergeDefaults(appoptions *App) {
	err := mergo.Merge(appoptions, Default)
	if err != nil {
		log.Fatal(err)
	}

}

func GetTrayMenus(appoptions *App) []*menu.TrayMenu {
	var result []*menu.TrayMenu
	switch runtime.GOOS {
	case "darwin":
		if appoptions.Mac != nil {
			result = appoptions.Mac.TrayMenus
		}
		//case "linux":
		//	if appoptions.Linux != nil {
		//		result = appoptions.Linux.TrayMenu
		//	}
		//case "windows":
		//	if appoptions.Windows != nil {
		//		result = appoptions.Windows.TrayMenu
		//	}
	}

	if result == nil {
		result = appoptions.TrayMenus
	}

	return result
}

func GetApplicationMenu(appoptions *App) *menu.Menu {
	var result *menu.Menu
	switch runtime.GOOS {
	case "darwin":
		if appoptions.Mac != nil {
			result = appoptions.Mac.Menu
		}
		//case "linux":
		//	if appoptions.Linux != nil {
		//		result = appoptions.Linux.TrayMenu
		//	}
		//case "windows":
		//	if appoptions.Windows != nil {
		//		result = appoptions.Windows.TrayMenu
		//	}
	}

	if result == nil {
		result = appoptions.Menu
	}

	return result
}

func GetContextMenus(appoptions *App) []*menu.ContextMenu {
	var result []*menu.ContextMenu

	switch runtime.GOOS {
	case "darwin":
		if appoptions.Mac != nil {
			result = appoptions.Mac.ContextMenus
		}
		//case "linux":
		//	if appoptions.Linux != nil {
		//		result = appoptions.Linux.TrayMenu
		//	}
		//case "windows":
		//	if appoptions.Windows != nil {
		//		result = appoptions.Windows.TrayMenu
		//	}
	}

	if result == nil {
		result = appoptions.ContextMenus
	}

	return result
}
