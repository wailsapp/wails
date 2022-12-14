//go:build darwin

package application

/*

#cgo CFLAGS:  -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13
#include "application.h"
#include <stdlib.h>
*/
import "C"
import (
	"github.com/wailsapp/wails/exp/pkg/options"
)

func New() *App {
	C.Init()
	return newApp()
}

func newApp() *App {
	return &App{
		applicationEventListeners: make(map[uint][]func()),
		systemTrays:               make(map[uint]*SystemTray),
	}
}

func NewWithOptions(options *options.Application) *App {
	C.Init()
	if options.Mac != nil {
		C.SetActivationPolicy(C.int(options.Mac.ActivationPolicy))
	}
	return &App{
		options:                   options,
		applicationEventListeners: make(map[uint][]func()),
		systemTrays:               make(map[uint]*SystemTray),
	}
}

func (a *App) run() error {
	C.Run()
	return nil
}

//export processApplicationEvent
func processApplicationEvent(eventID C.uint) {
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

//export processMenuItemClick
func processMenuItemClick(menuID C.uint) {
	menuItemClicked <- uint(menuID)
}
