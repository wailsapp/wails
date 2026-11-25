//go:build linux && !gtk3 && !android && !server

package application

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func (a *linuxApp) getApplicationMenu() *Menu {
	return nil
}

func (a *linuxApp) getIconForFile(filename string) ([]byte, error) {
	if filename == "" {
		return nil, nil
	}

	ext := filepath.Ext(filename)
	iconMap := map[string]string{
		".txt":  "text-x-generic",
		".pdf":  "application-pdf",
		".doc":  "x-office-document",
		".docx": "x-office-document",
		".xls":  "x-office-spreadsheet",
		".xlsx": "x-office-spreadsheet",
		".ppt":  "x-office-presentation",
		".pptx": "x-office-presentation",
		".zip":  "package-x-generic",
		".tar":  "package-x-generic",
		".gz":   "package-x-generic",
		".jpg":  "image-x-generic",
		".jpeg": "image-x-generic",
		".png":  "image-x-generic",
		".gif":  "image-x-generic",
		".mp3":  "audio-x-generic",
		".wav":  "audio-x-generic",
		".mp4":  "video-x-generic",
		".avi":  "video-x-generic",
		".html": "text-html",
		".css":  "text-css",
		".js":   "text-javascript",
		".json": "text-json",
		".xml":  "text-xml",
	}

	iconName := "application-x-generic"
	if name, ok := iconMap[ext]; ok {
		iconName = name
	}

	return getIconBytes(iconName)
}

func (a *linuxApp) markActivated() {
	a.activatedOnce.Do(func() {
		close(a.activated)
	})
}

func (a *linuxApp) waitForActivation() {
	<-a.activated
}

func (a *linuxApp) setApplicationMenu(menu *Menu) {}

func (a *linuxApp) run() error {
	if len(os.Args) == 2 {
		arg1 := os.Args[1]
		if strings.Contains(arg1, "://") {
			eventContext := newApplicationEventContext()
			eventContext.setURL(arg1)
			applicationEvents <- &ApplicationEvent{
				Id:  uint(events.Common.ApplicationLaunchedWithUrl),
				ctx: eventContext,
			}
		} else if a.parent.options.FileAssociations != nil {
			ext := filepath.Ext(arg1)
			if slices.Contains(a.parent.options.FileAssociations, ext) {
				eventContext := newApplicationEventContext()
				eventContext.setOpenedWithFile(arg1)
				applicationEvents <- &ApplicationEvent{
					Id:  uint(events.Common.ApplicationOpenedWithFile),
					ctx: eventContext,
				}
			}
		}
	}

	a.parent.Event.OnApplicationEvent(events.Linux.ApplicationStartup, func(evt *ApplicationEvent) {
		if err := a.processAndCacheScreens(); err != nil {
			a.parent.handleError(err)
		}
	})
	a.setupCommonEvents()
	// Started here, not from init(): init() is not part of the platformApp
	// interface and nothing calls it, so a monitor started there never runs.
	a.monitorThemeChanges()
	a.monitorPowerEvents()
	return appRun(a.application)
}

func (a *linuxApp) unregisterWindow(window windowPointer) {
	a.windowMapLock.Lock()
	delete(a.windowMap, window)
	remainingWindows := len(a.windowMap)
	a.windowMapLock.Unlock()

	if remainingWindows == 0 && !a.parent.options.Linux.DisableQuitOnLastWindowClosed {
		a.destroy()
	}
}

func getIconBytes(iconName string) ([]byte, error) {
	return nil, fmt.Errorf("icon lookup is not currently implemented for the GTK4 build path; build with -tags gtk3 for the legacy implementation")
}

func newPlatformApp(parent *App) *linuxApp {
	appID, err := applicationID(parent.options)
	if err != nil {
		parent.error("invalid Linux.ApplicationID: %w; falling back to %q", err, appID)
	}

	app := &linuxApp{
		parent:      parent,
		application: appNew(appID),
		activated:   make(chan struct{}),
		windowMap:   map[windowPointer]uint{},
	}

	if name := programName(parent.options, appID); name != "" {
		setProgramName(name)
	}

	return app
}

func (a *linuxApp) isVisible() bool {
	windows := a.getWindows()
	for _, window := range windows {
		if widgetIsVisible(window) {
			return true
		}
	}
	return false
}

// logPlatformInfo logs the platform information to the console
func (a *App) logPlatformInfo() {
	info, err := operatingsystem.Info()
	if err != nil {
		a.error("error getting OS info: %w", err)
		return
	}

	platformInfo := info.AsLogSlice()
	gtkMajor, gtkMinor, gtkMicro := gtkVersion()
	platformInfo = append(platformInfo, "GTK", fmt.Sprintf("%d.%d.%d", gtkMajor, gtkMinor, gtkMicro))
	webkitMajor, webkitMinor, webkitMicro := webkitVersion()
	platformInfo = append(platformInfo, "WebKitGTK", fmt.Sprintf("%d.%d.%d", webkitMajor, webkitMinor, webkitMicro))

	a.info("Platform Info:", platformInfo...)
}
