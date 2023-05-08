//go:build linux && !purego

package application

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include <gtk/gtk.h>
#include <gdk/gdk.h>
#include <webkit2/webkit2.h>
#include <stdio.h>
#include <limits.h>
#include <stdint.h>

typedef struct App {
  void *app;
} App;

extern void processApplicationEvent(uint);

extern void activateLinux(gpointer data);

static void activate (GtkApplication* app, gpointer data) {
   // FIXME: should likely emit a WAILS specific code
   // events.Mac.EventApplicationDidFinishLaunching == 1032
   //processApplicationEvent(1032);

   activateLinux(data);
}

static GtkApplication* init(char* name) {
   return gtk_application_new(name, G_APPLICATION_DEFAULT_FLAGS);
}

static int run(void *app, void *data) {
  g_signal_connect (app, "activate", G_CALLBACK (activate), data);
  g_application_hold(app);  // allows it to run without a window
  int status = g_application_run (G_APPLICATION (app), 0, NULL);
  g_application_release(app);
  g_object_unref (app);
  return status;
}

*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/events"
)

func init() {
	// Set GDK_BACKEND=x11 if currently unset and XDG_SESSION_TYPE is unset, unspecified or x11 to prevent warnings
	_ = os.Setenv("GDK_BACKEND", "x11")
}

type linuxApp struct {
	application     unsafe.Pointer
	applicationMenu unsafe.Pointer
	parent          *App

	startupActions []func()

	// Native -> uint
	windows     map[*C.GtkWindow]uint
	windowsLock sync.Mutex
}

func getNativeApplication() *linuxApp {
	return globalApplication.impl.(*linuxApp)
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
	/* // FIXME: WIP
	   loader := C.gdk_pixbuf_loader_new()

	   	if loader == nil {
	   		return
	   	}

	   loaded := C.gdk_pixbuf_loader_write(loader, (*C.guchar)(&icon[0]), (C.gsize)(len(icon)), 0)

	   	if loaded == C.bool(1) && C.gdk_pixbuf_loader_close(loader, 0) {
	   		pixbuf := C.gdk_pixbuf_loader_get_pixbuf(loader)
	   		if pixbuf != nil {
	   			ww := m.parent.CurrentWindow()
	   			window := ww.impl.window
	   			C.gtk_window_set_icon(window, pixbuf)
	   		}
	   	}

	   C.g_object_unref(loader)
	*/
}

func (m *linuxApp) name() string {
	// appName := C.getAppName()
	// defer C.free(unsafe.Pointer(appName))
	// return C.GoString(appName)
	return ""
}

func (m *linuxApp) getCurrentWindowID() uint {
	// TODO: Add extra metadata to window
	window := (*C.GtkWindow)(C.gtk_application_get_active_window((*C.GtkApplication)(m.application)))
	if window == nil {
		return uint(1)
	}
	m.windowsLock.Lock()
	defer m.windowsLock.Unlock()
	identifier, ok := m.windows[window]
	if ok {
		return identifier
	}
	return uint(1)
}

func (m *linuxApp) setApplicationMenu(menu *Menu) {
	if menu == nil {
		// Create a default menu
		menu = defaultApplicationMenu()
	}
	globalApplication.dispatchOnMainThread(func() {
		fmt.Println("setApplicationMenu")

		menu.Update()
		m.applicationMenu = (menu.impl).(*linuxMenu).native
	})
}

func (m *linuxApp) run() error {

	// Add a hook to the ApplicationDidFinishLaunching event
	// FIXME: add Wails specific events - i.e. Shouldn't platform specific ones be translated to Wails events?
	m.parent.On(events.Mac.ApplicationDidFinishLaunching, func() {
		// Do we need to do anything now?
		fmt.Println("events.Mac.ApplicationDidFinishLaunching received!")
	})

	var app C.App
	app.app = unsafe.Pointer(m)
	C.run(m.application, m.application)
	return nil
}

func (m *linuxApp) destroy() {
	C.g_application_quit((*C.GApplication)(m.application))
}

// register our window to our parent mapping
func (m *linuxApp) registerWindow(window *C.GtkWindow, id uint) {
	m.windowsLock.Lock()
	m.windows[window] = id
	m.windowsLock.Unlock()
}

func newPlatformApp(parent *App) *linuxApp {
	name := strings.ToLower(strings.Replace(parent.options.Name, " ", "", -1))
	if name == "" {
		name = "undefined"
	}
	nameC := C.CString(fmt.Sprintf("org.wails.%s", name))
	app := &linuxApp{
		parent:      parent,
		application: unsafe.Pointer(C.init(nameC)),
		//		name:        fmt.Sprintf("org.wails.%s", name),
		windows: map[*C.GtkWindow]uint{},
	}
	C.free(unsafe.Pointer(nameC))
	return app
}

// executeStartupActions is called by `activateLinux` below to execute
// code which needs to be run after the 'activate' signal is received
func (m *linuxApp) executeStartupActions() {
	for _, fn := range m.startupActions {
		fn()
	}
}

//export activateLinux
func activateLinux(data unsafe.Pointer) {
	//	globalApplication.activate()
	app := (globalApplication.impl).(*linuxApp)
	app.executeStartupActions()
}

//export processApplicationEvent
func processApplicationEvent(eventID C.uint) {
	// TODO: add translation to Wails events
	//       currently reusing Mac specific values
	applicationEvents <- uint(eventID)
}

//export processWindowEvent
func processWindowEvent(windowID C.uint, eventID C.uint) {
	windowEvents <- &WindowEvent{
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
func processMenuItemClick(menuID C.uint) {
	menuItemClicked <- uint(menuID)
}

func setIcon(icon []byte) {
	if icon == nil {
		return
	}
	//C.setApplicationIcon(unsafe.Pointer(&icon[0]), C.int(len(icon)))
}
