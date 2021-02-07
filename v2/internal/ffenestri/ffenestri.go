package ffenestri

import (
	"runtime"
	"strings"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/menumanager"

	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher"
	"github.com/wailsapp/wails/v2/pkg/options"
)

/*

#cgo linux CFLAGS: -DFFENESTRI_LINUX=1
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#cgo darwin CFLAGS: -DFFENESTRI_DARWIN=1
#cgo darwin LDFLAGS: -framework WebKit -lobjc

#include <stdlib.h>
#include "ffenestri.h"


*/
import "C"

// Application is our main application object
type Application struct {
	config *options.App
	memory []unsafe.Pointer

	// This is the main app pointer
	app *C.struct_Application

	// Manages menus
	menuManager *menumanager.Manager

	// Logger
	logger logger.CustomLogger
}

func (a *Application) saveMemoryReference(mem unsafe.Pointer) {
	a.memory = append(a.memory, mem)
}

func (a *Application) string2CString(str string) *C.char {
	result := C.CString(str)
	a.saveMemoryReference(unsafe.Pointer(result))
	return result
}

func init() {
	runtime.LockOSThread()
}

// NewApplicationWithConfig creates a new application based on the given config
func NewApplicationWithConfig(config *options.App, logger *logger.Logger, menuManager *menumanager.Manager) *Application {
	return &Application{
		config:      config,
		logger:      logger.CustomLogger("Ffenestri"),
		menuManager: menuManager,
	}
}

// NewApplication creates a new Application with the default config
func NewApplication(logger *logger.Logger) *Application {
	return &Application{
		config: options.Default,
		logger: logger.CustomLogger("Ffenestri"),
	}
}

func (a *Application) freeMemory() {
	for _, mem := range a.memory {
		// fmt.Printf("Freeing memory: %+v\n", mem)
		C.free(mem)
	}
}

// bool2Cint converts a Go boolean to a C integer
func (a *Application) bool2Cint(value bool) C.int {
	if value {
		return C.int(1)
	}
	return C.int(0)
}

// dispatcher is the interface to send messages to
var dispatcher *messagedispatcher.DispatchClient

// Dispatcher is what we register out client with
type Dispatcher interface {
	RegisterClient(client messagedispatcher.Client) *messagedispatcher.DispatchClient
}

// DispatchClient is the means for passing messages to the backend
type DispatchClient interface {
	SendMessage(string)
}

func intToColour(colour int) (C.int, C.int, C.int, C.int) {
	var alpha = C.int(colour & 0xFF)
	var blue = C.int((colour >> 8) & 0xFF)
	var green = C.int((colour >> 16) & 0xFF)
	var red = C.int((colour >> 24) & 0xFF)
	return red, green, blue, alpha
}

// Run the application
func (a *Application) Run(incomingDispatcher Dispatcher, bindings string, debug bool) error {
	title := a.string2CString(a.config.Title)
	width := C.int(a.config.Width)
	height := C.int(a.config.Height)
	resizable := a.bool2Cint(!a.config.DisableResize)
	devtools := a.bool2Cint(a.config.DevTools)
	fullscreen := a.bool2Cint(a.config.Fullscreen)
	startHidden := a.bool2Cint(a.config.StartHidden)
	logLevel := C.int(a.config.LogLevel)
	hideWindowOnClose := a.bool2Cint(a.config.HideWindowOnClose)
	app := C.NewApplication(title, width, height, resizable, devtools, fullscreen, startHidden, logLevel, hideWindowOnClose)

	// Save app reference
	a.app = (*C.struct_Application)(app)

	// Set Min Window Size
	minWidth := C.int(a.config.MinWidth)
	minHeight := C.int(a.config.MinHeight)
	C.SetMinWindowSize(a.app, minWidth, minHeight)

	// Set Max Window Size
	maxWidth := C.int(a.config.MaxWidth)
	maxHeight := C.int(a.config.MaxHeight)
	C.SetMaxWindowSize(a.app, maxWidth, maxHeight)

	// Set debug if needed
	C.SetDebug(app, a.bool2Cint(debug))

	// TODO: Move frameless to Linux options
	// if a.config.Frameless {
	// 	C.DisableFrame(a.app)
	// }

	if a.config.RGBA != 0 {
		r, g, b, alpha := intToColour(a.config.RGBA)
		C.SetColour(a.app, r, g, b, alpha)
	}

	// Escape bindings so C doesn't freak out
	bindings = strings.ReplaceAll(bindings, `"`, `\"`)

	// Set bindings
	C.SetBindings(app, a.string2CString(bindings))

	// save the dispatcher in a package variable so that the C callbacks
	// can access it
	dispatcher = incomingDispatcher.RegisterClient(newClient(a))

	// Process platform settings
	err := a.processPlatformSettings()
	if err != nil {
		return err
	}

	// Check we could initialise the application
	if app != nil {
		// Yes - Save memory reference and run app, cleaning up afterwards
		a.saveMemoryReference(unsafe.Pointer(app))
		C.Run(app, 0, nil)
	} else {
		// Oh no! We couldn't initialise the application
		a.logger.Fatal("Cannot initialise Application.")
	}

	a.freeMemory()
	return nil
}

// messageFromWindowCallback is called by any messages sent in
// webkit to window.external.invoke. It relays the message on to
// the dispatcher.
//export messageFromWindowCallback
func messageFromWindowCallback(data *C.char) {
	dispatcher.DispatchMessage(C.GoString(data))
}
