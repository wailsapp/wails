//go:build linux && purego

package application

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/ebitengine/purego"
	"github.com/wailsapp/wails/v2/pkg/assetserver/webview"
)

const (
	gtk3 = "libgtk-3.so"
	gtk4 = "libgtk-4.so"
)

var (
	gtk     uintptr
	version int
	webkit  uintptr
)

func init() {
	// needed for GTK4 to function
	_ = os.Setenv("GDK_BACKEND", "x11")
	var err error
	/*
		gtk, err = purego.Dlopen(gtk4, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err == nil {
			version = 4
			return
		}

		log.Println("Failed to open GTK4: Falling back to GTK3")
	*/
	gtk, err = purego.Dlopen(gtk3, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}
	version = 3

	var webkit4 string = "libwebkit2gtk-4.1.so"
	webkit, err = purego.Dlopen(webkit4, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}
}

type linuxApp struct {
	appName         string
	application     uintptr
	applicationMenu uintptr
	parent          *App

	// Native -> uint
	windows     map[uintptr]uint
	windowsLock sync.Mutex
}

func (m *linuxApp) hide() {
	//	C.hide()
}

func (m *linuxApp) show() {
	//	C.show()
}

func (m *linuxApp) on(eventID uint) {
	log.Println("linuxApp.on()", eventID)

	// TODO: Setup signal handling as appropriate
	// Note: GTK signals seem to be strings!
}

func (m *linuxApp) setIcon(icon []byte) {
	//	C.setApplicationIcon(unsafe.Pointer(&icon[0]), C.int(len(icon)))
}

func (m *linuxApp) name() string {
	return m.appName
}

func (m *linuxApp) getCurrentWindowID() uint {
	var getCurrentWindow func(uintptr) uintptr
	purego.RegisterLibFunc(&getCurrentWindow, gtk, "gtk_application_get_active_window")
	window := getCurrentWindow(m.application)
	if window == 0 {
		return 1
	}
	m.windowsLock.Lock()
	defer m.windowsLock.Unlock()
	if identifier, ok := m.windows[window]; ok {
		return identifier
	}

	return 1
}

func (m *linuxApp) setApplicationMenu(menu *Menu) {
	if menu == nil {
		// Create a default menu
		menu = defaultApplicationMenu()
	}
	globalApplication.dispatchOnMainThread(func() {
		menu.Update()
		m.applicationMenu = (menu.impl).(*linuxMenu).native
	})
}

func (m *linuxApp) activate() {
	fmt.Println("linuxApp.activated!", m.application)
	var hold func(uintptr)
	purego.RegisterLibFunc(&hold, gtk, "g_application_hold")

	hold(m.application)

	//	time.Sleep(50 * time.Millisecond)
	//	m.parent.activate()
}

func (m *linuxApp) run() error {
	// Add a hook to the ApplicationDidFinishLaunching event
	// FIXME: add Wails specific events - i.e. Shouldn't platform specific ones be translated to Wails events?
	/*	m.parent.On(events.Mac.ApplicationDidFinishLaunching, func() {
			// Do we need to do anything now?
			fmt.Println("ApplicationDidFinishLaunching!")
		})
	*/
	m.parent.OnWindowCreation(func(window *WebviewWindow) {
		fmt.Println("OnWindowCreation: ", window)

	})

	var g_signal_connect func(uintptr, string, uintptr, uintptr, bool, int) int
	purego.RegisterLibFunc(&g_signal_connect, gtk, "g_signal_connect_data")
	g_signal_connect(m.application, "activate", purego.NewCallback(m.activate), m.application, false, 0)

	var run func(uintptr, int, []string) int
	purego.RegisterLibFunc(&run, gtk, "g_application_run")

	// FIXME: Convert status to 'error' if needed
	status := run(m.application, 0, []string{})
	fmt.Println("status", status)

	var release func(uintptr)
	purego.RegisterLibFunc(&release, gtk, "g_application_release")
	release(m.application)

	purego.RegisterLibFunc(&release, gtk, "g_object_unref")
	release(m.application)

	return nil
}

func (m *linuxApp) destroy() {
	var quit func(uintptr)
	purego.RegisterLibFunc(&quit, gtk, "g_application_quit")
	quit(m.application)
}

func (m *linuxApp) registerWindow(address uintptr, window uint) {
	m.windowsLock.Lock()
	m.windows[address] = window
	m.windowsLock.Unlock()
}

func newPlatformApp(parent *App) *linuxApp {
	name := strings.ToLower(parent.options.Name)
	if name == "" {
		name = "undefined"
	}
	identifier := fmt.Sprintf("org.wails.%s", strings.Replace(name, " ", "-", -1))

	var gtkNew func(string, uint) uintptr
	purego.RegisterLibFunc(&gtkNew, gtk, "gtk_application_new")
	app := &linuxApp{
		appName:     identifier,
		parent:      parent,
		application: gtkNew(identifier, 0),
		windows:     map[uintptr]uint{},
	}
	return app
}

func processApplicationEvent(eventID uint) {
	// TODO: add translation to Wails events
	//       currently reusing Mac specific values
	applicationEvents <- eventID
}

func processWindowEvent(windowID uint, eventID uint) {
	windowEvents <- &WindowEvent{
		WindowID: windowID,
		EventID:  eventID,
	}
}

func processMessage(windowID uint, message string) {
	windowMessageBuffer <- &windowMessage{
		windowId: windowID,
		message:  message,
	}
}

func processURLRequest(windowID uint, wkUrlSchemeTask uintptr) {
	fmt.Println("processURLRequest", windowID, wkUrlSchemeTask)
	webviewRequests <- &webViewAssetRequest{
		Request:    webview.NewRequest(wkUrlSchemeTask),
		windowId:   windowID,
		windowName: globalApplication.getWindowForID(windowID).Name(),
	}
}

func processDragItems(windowID uint, arr []string, length int) {
	windowDragAndDropBuffer <- &dragAndDropMessage{
		windowId:  windowID,
		filenames: arr,
	}
}

func processMenuItemClick(menuID uint) {
	menuItemClicked <- menuID
}

func setIcon(icon []byte) {
	if icon == nil {
		return
	}
	fmt.Println("setIcon")
	/*
	   GdkPixbufLoader *loader = gdk_pixbuf_loader_new();
	   if (!loader)
	   {
	       return;
	   }
	   if (gdk_pixbuf_loader_write(loader, buf, len, NULL) && gdk_pixbuf_loader_close(loader, NULL))
	   {
	       GdkPixbuf *pixbuf = gdk_pixbuf_loader_get_pixbuf(loader);
	       if (pixbuf)
	       {
	           gtk_window_set_icon(window, pixbuf);
	       }
	   }
	   g_object_unref(loader);*/
}
