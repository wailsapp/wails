//go:build darwin

package application

/*

#cgo CFLAGS:  -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13
#include "application.h"
*/
import "C"
import "github.com/wailsapp/wails/exp/pkg/options"

func New(options *options.Application) *App {
	C.Init()
	if options.Mac != nil {
		C.SetActivationPolicy(C.int(options.Mac.ActivationPolicy))
	}
	return &App{
		options:              options,
		systemEventListeners: make(map[string][]func()),
	}
}

func (a *App) Run() error {

	go func() {
		for {
			event := <-systemEvents
			a.handleSystemEvent(event)
		}
	}()

	// run windows
	for _, window := range a.windows {
		err := window.Run()
		if err != nil {
			return err
		}
	}

	C.Run()
	return nil
}

func (a *App) handleSystemEvent(event string) {
	listeners, ok := a.systemEventListeners[event]
	if !ok {
		return
	}
	for _, listener := range listeners {
		go listener()
	}
}

//export systemEventHandler
func systemEventHandler(name *C.char) {
	goString := C.GoString(name)
	systemEvents <- goString
}
