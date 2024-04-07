//go:build linux

package application

import "C"
import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func init() {
	// FIXME: This should be handled appropriately in the individual files most likely.
	// Set GDK_BACKEND=x11 if currently unset and XDG_SESSION_TYPE is unset, unspecified or x11 to prevent warnings
	_ = os.Setenv("GDK_BACKEND", "x11")
}

type linuxApp struct {
	application pointer
	parent      *App

	startupActions []func()

	// Native -> uint
	windowMap     map[windowPointer]uint
	windowMapLock sync.Mutex

	theme string
}

func (a *linuxApp) GetFlags(options Options) map[string]any {
	if options.Flags == nil {
		options.Flags = make(map[string]any)
	}
	return options.Flags
}

func getNativeApplication() *linuxApp {
	return globalApplication.impl.(*linuxApp)
}

func (a *linuxApp) hide() {
	a.hideAllWindows()
}

func (a *linuxApp) show() {
	a.showAllWindows()
}

func (a *linuxApp) on(eventID uint) {
	// TODO: Test register/unregister events
	//C.registerApplicationEvent(l.application, C.uint(eventID))
}

func (a *linuxApp) setIcon(icon []byte) {

	log.Println("linuxApp.setIcon", "not implemented")
}

func (a *linuxApp) name() string {
	return appName()
}

type rnr struct {
	f func()
}

func (r rnr) run() {
	r.f()
}

func (a *linuxApp) setApplicationMenu(menu *Menu) {
	// FIXME: How do we avoid putting a menu?
	if menu == nil {
		// Create a default menu
		menu = defaultApplicationMenu()
		globalApplication.ApplicationMenu = menu
	}
}

func (a *linuxApp) run() error {

	a.parent.On(events.Linux.ApplicationStartup, func(evt *Event) {
		fmt.Println("events.Linux.ApplicationStartup received!")
	})
	a.setupCommonEvents()
	a.monitorThemeChanges()
	return appRun(a.application)
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

func (a *linuxApp) destroy() {
	if !globalApplication.shouldQuit() {
		return
	}
	globalApplication.cleanup()
	appDestroy(a.application)
}

func (a *linuxApp) isOnMainThread() bool {
	return isOnMainThread()
}

// register our window to our parent mapping
func (a *linuxApp) registerWindow(window pointer, id uint) {
	a.windowMapLock.Lock()
	a.windowMap[windowPointer(window)] = id
	a.windowMapLock.Unlock()
}

func (a *linuxApp) isDarkMode() bool {
	return strings.Contains(a.theme, "dark")
}

func (a *linuxApp) monitorThemeChanges() {
	go func() {
		conn, err := dbus.ConnectSessionBus()
		if err != nil {
			a.parent.info("[WARNING] Failed to connect to session bus; monitoring for theme changes will not function:", err)
			return
		}
		defer conn.Close()

		if err = conn.AddMatchSignal(
			dbus.WithMatchObjectPath("/org/freedesktop/portal/desktop"),
		); err != nil {
			panic(err)
		}

		c := make(chan *dbus.Signal, 10)
		conn.Signal(c)

		getTheme := func(body []interface{}) (string, bool) {
			if len(body) < 2 {
				return "", false
			}
			if body[0].(string) != "org.gnome.desktop.interface" {
				return "", false
			}
			if body[1].(string) == "color-scheme" {
				return body[2].(dbus.Variant).Value().(string), true
			}
			return "", false
		}

		for v := range c {
			theme, ok := getTheme(v.Body)
			if !ok {
				continue
			}

			if theme != a.theme {
				a.theme = theme
				event := newApplicationEvent(events.Common.ThemeChanged)
				event.Context().setIsDarkMode(a.isDarkMode())
				applicationEvents <- event
			}

		}
	}()
}

func newPlatformApp(parent *App) *linuxApp {
	name := strings.ToLower(strings.Replace(parent.options.Name, " ", "", -1))
	if name == "" {
		name = "undefined"
	}
	app := &linuxApp{
		parent:      parent,
		application: appNew(name),
		windowMap:   map[windowPointer]uint{},
	}
	return app
}

// logPlatformInfo logs the platform information to the console
func (a *App) logPlatformInfo() {
	info, err := operatingsystem.Info()
	if err != nil {
		a.error("Error getting OS info", "error", err.Error())
		return
	}

	wkVersion := operatingsystem.GetWebkitVersion()
	platformInfo := info.AsLogSlice()
	platformInfo = append(platformInfo, "Webkit2Gtk", wkVersion)

	a.info("Platform Info:", platformInfo...)
}

//export processWindowEvent
func processWindowEvent(windowID C.uint, eventID C.uint) {
	windowEvents <- &windowEvent{
		WindowID: uint(windowID),
		EventID:  uint(eventID),
	}
}
