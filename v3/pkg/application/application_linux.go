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

func (l *linuxApp) GetFlags(options Options) map[string]any {
	if options.Flags == nil {
		options.Flags = make(map[string]any)
	}
	return options.Flags
}

func getNativeApplication() *linuxApp {
	return globalApplication.impl.(*linuxApp)
}

func (l *linuxApp) hide() {
	hideAllWindows(l.application)
}

func (l *linuxApp) show() {
	showAllWindows(l.application)
}

func (l *linuxApp) on(eventID uint) {
	// TODO: Test register/unregister events
	//C.registerApplicationEvent(l.application, C.uint(eventID))
}

func (l *linuxApp) setIcon(icon []byte) {

	log.Println("linuxApp.setIcon", "not implemented")
}

func (l *linuxApp) name() string {
	return appName()
}

func (l *linuxApp) getCurrentWindowID() uint {
	return getCurrentWindowID(l.application, l.windowMap)
}

type rnr struct {
	f func()
}

func (r rnr) run() {
	r.f()
}

func (l *linuxApp) setApplicationMenu(menu *Menu) {
	// FIXME: How do we avoid putting a menu?
	if menu == nil {
		// Create a default menu
		menu = defaultApplicationMenu()
		globalApplication.ApplicationMenu = menu
	}
}

func (l *linuxApp) run() error {

	l.parent.On(events.Linux.ApplicationStartup, func(evt *Event) {
		fmt.Println("events.Linux.ApplicationStartup received!")
	})
	l.setupCommonEvents()
	l.monitorThemeChanges()
	return appRun(l.application)
}

func (l *linuxApp) unregisterWindow(w windowPointer) {
	l.windowMapLock.Lock()
	delete(l.windowMap, w)
	l.windowMapLock.Unlock()

	// If this was the last window...
	if len(l.windowMap) == 0 && !l.parent.options.Linux.DisableQuitOnLastWindowClosed {
		l.destroy()
	}
}

func (l *linuxApp) destroy() {
	if !globalApplication.shouldQuit() {
		return
	}
	globalApplication.cleanup()
	appDestroy(l.application)
}

func (l *linuxApp) isOnMainThread() bool {
	return isOnMainThread()
}

// register our window to our parent mapping
func (l *linuxApp) registerWindow(window pointer, id uint) {
	l.windowMapLock.Lock()
	l.windowMap[windowPointer(window)] = id
	l.windowMapLock.Unlock()
}

func (l *linuxApp) isDarkMode() bool {
	return strings.Contains(l.theme, "dark")
}

func (l *linuxApp) monitorThemeChanges() {
	go func() {
		conn, err := dbus.ConnectSessionBus()
		if err != nil {
			l.parent.info("[WARNING] Failed to connect to session bus; monitoring for theme changes will not function:", err)
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

			if theme != l.theme {
				l.theme = theme
				event := newApplicationEvent(events.Common.ThemeChanged)
				event.Context().setIsDarkMode(l.isDarkMode())
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
