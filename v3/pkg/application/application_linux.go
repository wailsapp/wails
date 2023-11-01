//go:build linux

package application

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

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
	//	log.Println("linuxApp.on()", eventID)
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

	// Add a hook to the ApplicationDidFinishLaunching event
	// FIXME: add Wails specific events - i.e. Shouldn't platform specific ones be translated to Wails events?
	m.parent.On(events.Mac.ApplicationDidFinishLaunching, func(evt *Event) {
		// Do we need to do anything now?
		fmt.Println("events.Mac.ApplicationDidFinishLaunching received!")
	})

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
	// FIXME: How do we detect this?
	// Maybe this helps: https://askubuntu.com/questions/1469869/how-does-firefox-detect-light-dark-theme-change-on-kde-systems
	return false
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

/*
//export processApplicationEvent
func processApplicationEvent(eventID C.uint) {
	// TODO: add translation to Wails events
	//       currently reusing Mac specific values
	applicationEvents <- uint(eventID)
}

//export processWindowEvent
func processWindowEvent(windowID C.uint, eventID C.uint) {
	windowEvents <- &windowEvent{
		WindowID: uint(windowID),
		EventID:  uint(eventID),
	}
}

//export processMessage
func processMessage(windowID C.uint, message *C.char) {
	windowMessageBuffer <- &windowMessage{
		windowId: uint(windowID),
		message:  C.GoString(message),
	}
}

//export processDragItems
func processDragItems(windowID C.uint, arr **C.char, length C.int) {
	var filenames []string
	// Convert the C array to a Go slice
	goSlice := (*[1 << 30]*C.char)(unsafe.Pointer(arr))[:length:length]
	for _, str := range goSlice {
		filenames = append(filenames, C.GoString(str))
	}
	windowDragAndDropBuffer <- &dragAndDropMessage{
		windowId:  uint(windowID),
		filenames: filenames,
	}
}

//export processMenuItemClick
func processMenuItemClick(menuID identifier) {
	menuItemClicked <- uint(menuID)
}

func setIcon(icon []byte) {
	if icon == nil {
		return
	}
	//C.setApplicationIcon(unsafe.Pointer(&icon[0]), C.int(len(icon)))
}
*/

// logPlatformInfo logs the platform information to the console
func (a *App) logPlatformInfo() {
	info, err := operatingsystem.Info()
	if err != nil {
		a.error("Error getting OS info", "error", err.Error())
		return
	}

	a.info("Platform Info:", info.AsLogSlice()...)
}
