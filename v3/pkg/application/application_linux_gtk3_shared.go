//go:build linux && gtk3 && !android && !server

package application

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// logPlatformInfo logs the platform information to the console
func (a *App) logPlatformInfo() {
	info, err := operatingsystem.Info()
	if err != nil {
		a.error("error getting OS info: %w", err)
		return
	}

	wkVersion := operatingsystem.GetWebkitVersion()
	platformInfo := info.AsLogSlice()
	platformInfo = append(platformInfo, "Webkit2Gtk", wkVersion)

	a.info("Platform Info:", platformInfo...)
}

func (a *linuxApp) run() error {

	if len(os.Args) == 2 { // Case: program + 1 argument
		arg1 := os.Args[1]
		// Check if the argument is likely a URL from a custom protocol invocation
		if strings.Contains(arg1, "://") {
			a.parent.debug("Application launched with argument, potentially a URL from custom protocol", "url", arg1)
			eventContext := newApplicationEventContext()
			eventContext.setURL(arg1)
			applicationEvents <- &ApplicationEvent{
				Id:  uint(events.Common.ApplicationLaunchedWithUrl),
				ctx: eventContext,
			}
		} else {
			// Check if the argument matches any file associations
			matched := false
			if a.parent.options.FileAssociations != nil {
				ext := filepath.Ext(arg1)
				if slices.Contains(a.parent.options.FileAssociations, ext) {
					a.parent.debug("File opened via file association", "file", arg1, "extension", ext)
					eventContext := newApplicationEventContext()
					eventContext.setOpenedWithFile(arg1)
					applicationEvents <- &ApplicationEvent{
						Id:  uint(events.Common.ApplicationOpenedWithFile),
						ctx: eventContext,
					}
					matched = true
				}
			}
			if !matched {
				a.parent.debug("Application launched with single argument (not a URL), potential file open?", "arg", arg1)
			}
		}
	} else if len(os.Args) > 2 {
		// Log if multiple arguments are passed
		a.parent.debug("Application launched with multiple arguments", "args", os.Args[1:])
	}

	a.parent.Event.OnApplicationEvent(events.Linux.ApplicationStartup, func(evt *ApplicationEvent) {
		if err := a.processAndCacheScreens(); err != nil {
			a.parent.handleError(err)
		}
	})
	a.setupCommonEvents()
	a.monitorThemeChanges()
	a.monitorPowerEvents()
	return appRun(a.application)
}

func (a *linuxApp) setApplicationMenu(menu *Menu) {
	// FIXME: How do we avoid putting a menu?
	if menu == nil {
		// Create a default menu
		menu = DefaultApplicationMenu()
		globalApplication.applicationMenu = menu
	}
}

func (a *linuxApp) unregisterWindow(w windowPointer) {
	a.windowMapLock.Lock()
	delete(a.windowMap, w)
	a.windowMapLock.Unlock()

	// If this was the last window...
	if len(a.windowMap) == 0 && !a.parent.options.Linux.DisableQuitOnLastWindowClosed {
		a.destroy()
	}
}

func newPlatformApp(parent *App) *linuxApp {
	appID, err := applicationID(parent.options)
	if err != nil {
		parent.error("invalid Linux.ApplicationID: %w; falling back to %q", err, appID)
	}

	app := &linuxApp{
		parent:      parent,
		application: appNew(appID),
		windowMap:   map[windowPointer]uint{},
	}

	if name := programName(parent.options, appID); name != "" {
		setProgramName(name)
	}

	return app
}

func (r rnr) run() {
	r.f()
}
