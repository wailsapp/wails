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

func New(options *options.Application) *App {
	C.Init()
	if options.Mac != nil {
		C.SetActivationPolicy(C.int(options.Mac.ActivationPolicy))
	}
	return &App{
		options:                   options,
		applicationEventListeners: make(map[uint][]func()),
	}
}

func (a *App) run() error {
	C.Run()
	return nil
}

func (a *App) handleApplicationEvent(event uint) {
	listeners, ok := a.applicationEventListeners[event]
	if !ok {
		return
	}
	for _, listener := range listeners {
		go listener()
	}
}

//export applicationEventHandler
func applicationEventHandler(eventID C.uint) {
	applicationEvents <- uint(eventID)
}

//export windowEventHandler
func windowEventHandler(windowID C.uint, eventID C.uint) {
	windowEvents <- &WindowEvent{
		WindowID: uint(windowID),
		EventID:  uint(eventID),
	}
}

//export processMessage
func processMessage(windowID C.uint, message *C.char) {
	messageBuffer <- &windowMessage{
		windowId: uint(windowID),
		message:  C.GoString(message),
	}
}
