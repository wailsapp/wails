//go:build !(linux && gtk3)

package main

import (
	"fmt"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// GTK4 shows context menus as popovers that do not take the X keyboard focus
// without a window manager, so its context menu is driven with the mouse.
// macOS does not navigate the menu bar and its context menus take the
// keyboard.
const contextMenuByMouse = runtime.GOOS == "linux"

// menuNavigateKeys are the keys that, once F10 has taken the menu bar, choose
// the first item of the first menu. GTK 4.24 and Windows select the first
// menu without opening it, and Down opens it on its first item. GTK 4.22 and
// older open the menu with that item already highlighted.
func menuNavigateKeys(app *application.App) []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"down", "return"}
	case "linux":
		var major, minor int
		version, _ := app.Env.Info().PlatformInfo["gtk4-runtime"].(string)
		if _, err := fmt.Sscanf(version, "%d.%d", &major, &minor); err == nil && (major > 4 || minor >= 24) {
			return []string{"down", "return"}
		}
		return []string{"return"}
	}
	return nil
}
