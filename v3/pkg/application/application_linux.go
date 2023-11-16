//go:build linux

package application

import (
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
	application     pointer
	applicationMenu pointer
	parent          *App

	startupActions []func()

	// Native -> uint
	windows     map[windowPointer]uint
	windowsLock sync.Mutex

	theme string
}

func (m *linuxApp) GetFlags(options Options) map[string]any {
	if options.Flags == nil {
		options.Flags = make(map[string]any)
	}
	return options.Flags
}

func getNativeApplication() *linuxApp {
	return globalApplication.impl.(*linuxApp)
}

func (m *linuxApp) hide() {
	hideAllWindows(m.application)
}

func (m *linuxApp) show() {
	showAllWindows(m.application)
}

func (m *linuxApp) on(eventID uint) {
	// TODO: What do we need to do here?
	log.Println("linuxApp.on()", eventID)
}

func (m *linuxApp) setIcon(icon []byte) {

	log.Println("linuxApp.setIcon", "not implemented")
}

func (m *linuxApp) name() string {
	return appName()
}

func (m *linuxApp) getCurrentWindowID() uint {
	return getCurrentWindowID(m.application, m.windows)
}

type rnr struct {
	f func()
}

func (r rnr) run() {
	r.f()
}

func (m *linuxApp) getApplicationMenu() pointer {
	if m.applicationMenu != nilPointer {
		return m.applicationMenu
	}

	menu := globalApplication.ApplicationMenu
	if menu != nil {
		InvokeSync(func() {
			menu.Update()
		})
		m.applicationMenu = (menu.impl).(*linuxMenu).native
	}
	return m.applicationMenu
}

func (m *linuxApp) setApplicationMenu(menu *Menu) {
	// FIXME: How do we avoid putting a menu?
	if menu == nil {
		// Create a default menu
		menu = defaultApplicationMenu()
		globalApplication.ApplicationMenu = menu
	}
}

func (m *linuxApp) run() error {
	m.monitorThemeChanges()
	return appRun(m.application)
}

func (m *linuxApp) destroy() {
	appDestroy(m.application)
}

func (m *linuxApp) isOnMainThread() bool {
	return isOnMainThread()
}

// register our window to our parent mapping
func (m *linuxApp) registerWindow(window pointer, id uint) {
	m.windowsLock.Lock()
	m.windows[windowPointer(window)] = id
	m.windowsLock.Unlock()
}

func (m *linuxApp) isDarkMode() bool {
	return strings.Contains(m.theme, "dark")
}

func (m *linuxApp) monitorThemeChanges() {
	go func() {
		conn, err := dbus.ConnectSessionBus()
		if err != nil {
			m.parent.info("[WARNING] Failed to connect to session bus; monitoring for theme changes will not function:", err)
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

			if theme != m.theme {
				m.theme = theme
				event := newApplicationEvent(events.Common.ThemeChanged)
				event.Context().setIsDarkMode(m.isDarkMode())
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
		windows:     map[windowPointer]uint{},
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
