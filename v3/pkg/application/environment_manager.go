package application

import (
	"runtime"

	"github.com/wailsapp/wails/v3/internal/fileexplorer"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
)

// EnvironmentManager manages environment-related operations
type EnvironmentManager struct {
	app *App
}

// newEnvironmentManager creates a new EnvironmentManager instance
func newEnvironmentManager(app *App) *EnvironmentManager {
	return &EnvironmentManager{
		app: app,
	}
}

// Info returns environment information
func (em *EnvironmentManager) Info() EnvironmentInfo {
	info, _ := operatingsystem.Info()
	result := EnvironmentInfo{
		OS:     runtime.GOOS,
		Arch:   runtime.GOARCH,
		Debug:  em.app.isDebugMode,
		OSInfo: info,
	}
	result.PlatformInfo = em.app.platformEnvironment()
	return result
}

// IsDarkMode returns true if the system is in dark mode
func (em *EnvironmentManager) IsDarkMode() bool {
	if em.app.impl == nil {
		return false
	}
	return em.app.impl.isDarkMode()
}

// GetAccentColor returns the system accent color
func (em *EnvironmentManager) GetAccentColor() string {
	if em.app.impl == nil {
		return "rgb(0,122,255)"
	}
	return em.app.impl.getAccentColor()
}

// OpenFileManager opens the file manager at the specified path, optionally selecting the file
func (em *EnvironmentManager) OpenFileManager(path string, selectFile bool) error {
	return InvokeSyncWithError(func() error {
		return fileexplorer.OpenFileManager(path, selectFile)
	})
}
