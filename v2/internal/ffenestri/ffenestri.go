package ffenestri

import (
	"runtime"
	"strings"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/features"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher"
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

// DEBUG is the global Ffenestri debug flag.
// TODO: move to compile time.
var DEBUG bool = true

// Config defines how our application should be configured
type Config struct {
	Title       string
	Width       int
	Height      int
	MinWidth    int
	MinHeight   int
	MaxWidth    int
	MaxHeight   int
	DevTools    bool
	Resizable   bool
	Fullscreen  bool
	Frameless   bool
	StartHidden bool
}

var defaultConfig = &Config{
	Title:       "My Wails App",
	Width:       800,
	Height:      600,
	DevTools:    true,
	Resizable:   true,
	Fullscreen:  false,
	Frameless:   false,
	StartHidden: false,
}

// Application is our main application object
type Application struct {
	config *Config
	memory []unsafe.Pointer

	// This is the main app pointer
	app unsafe.Pointer

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
func NewApplicationWithConfig(config *Config, logger *logger.Logger) *Application {
	return &Application{
		config: config,
		logger: logger.CustomLogger("Ffenestri"),
	}
}

// NewApplication creates a new Application with the default config
func NewApplication(logger *logger.Logger) *Application {
	return &Application{
		config: defaultConfig,
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

// Run the application
func (a *Application) Run(incomingDispatcher Dispatcher, bindings string, features *features.Features) error {
	title := a.string2CString(a.config.Title)
	width := C.int(a.config.Width)
	height := C.int(a.config.Height)
	resizable := a.bool2Cint(a.config.Resizable)
	devtools := a.bool2Cint(a.config.DevTools)
	fullscreen := a.bool2Cint(a.config.Fullscreen)
	startHidden := a.bool2Cint(a.config.StartHidden)
	app := C.NewApplication(title, width, height, resizable, devtools, fullscreen, startHidden)

	// Save app reference
	a.app = unsafe.Pointer(app)

	// Set Min Window Size
	minWidth := C.int(a.config.MinWidth)
	minHeight := C.int(a.config.MinHeight)
	C.SetMinWindowSize(a.app, minWidth, minHeight)

	// Set Max Window Size
	maxWidth := C.int(a.config.MaxWidth)
	maxHeight := C.int(a.config.MaxHeight)
	C.SetMaxWindowSize(a.app, maxWidth, maxHeight)

	// Set debug if needed
	C.SetDebug(app, a.bool2Cint(DEBUG))

	// Set Frameless
	if a.config.Frameless {
		C.DisableFrame(a.app)
	}

	// Escape bindings so C doesn't freak out
	bindings = strings.ReplaceAll(bindings, `"`, `\"`)

	// Set bindings
	C.SetBindings(app, a.string2CString(bindings))

	// Process feature flags
	a.processFeatureFlags(features)

	// save the dispatcher in a package variable so that the C callbacks
	// can access it
	dispatcher = incomingDispatcher.RegisterClient(newClient(a))

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

func (a *Application) processFeatureFlags(features *features.Features) {

	// Process generic features

	// Process OS Specific flags
	a.processOSFeatureFlags(features)
}
