package options

import (
	"github.com/wailsapp/wails/v2/pkg/menu"
	"log"
	"runtime"

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
	ContextMenus  *menu.ContextMenus
	Tray          *menu.TrayOptions
	Menu          *menu.Menu
	Mac           *mac.Options
	Logger        logger.Logger `json:"-"`
	LogLevel      logger.LogLevel
}

// MergeDefaults will set the minimum default values for an application
func MergeDefaults(appoptions *App) {
	err := mergo.Merge(appoptions, Default)
	if err != nil {
		log.Fatal(err)
	}

}

func GetTray(appoptions *App) *menu.TrayOptions {
	var result *menu.TrayOptions
	switch runtime.GOOS {
	case "darwin":
		if appoptions.Mac != nil {
			result = appoptions.Mac.Tray
		}
		//case "linux":
		//	if appoptions.Linux != nil {
		//		result = appoptions.Linux.Tray
		//	}
		//case "windows":
		//	if appoptions.Windows != nil {
		//		result = appoptions.Windows.Tray
		//	}
	}

	if result == nil {
		result = appoptions.Tray
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
		//		result = appoptions.Linux.Tray
		//	}
		//case "windows":
		//	if appoptions.Windows != nil {
		//		result = appoptions.Windows.Tray
		//	}
	}

	if result == nil {
		result = appoptions.Menu
	}

	return result
}

func GetContextMenus(appoptions *App) *menu.ContextMenus {
	var result *menu.ContextMenus

	result = appoptions.ContextMenus
	var contextMenuOverrides *menu.ContextMenus
	switch runtime.GOOS {
	case "darwin":
		if appoptions.Mac != nil {
			contextMenuOverrides = appoptions.Mac.ContextMenus
		}
		//case "linux":
		//	if appoptions.Linux != nil {
		//		result = appoptions.Linux.Tray
		//	}
		//case "windows":
		//	if appoptions.Windows != nil {
		//		result = appoptions.Windows.Tray
		//	}
	}

	// Overwrite defaults with OS Specific context menus
	if contextMenuOverrides != nil {
		for id, contextMenu := range contextMenuOverrides.Items {
			result.AddMenu(id, contextMenu)
		}
	}

	return result
}
