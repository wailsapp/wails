//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "presentation_options_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func macPresentationOptionsSupported() bool { return true }

func macSetPresentationOptions(options MacPresentationOptions) error {
	message := InvokeSyncWithResult(func() *C.char {
		return C.presentationOptionsSet(C.ulong(options))
	})
	if message == nil {
		return nil
	}
	defer C.free(unsafe.Pointer(message))
	return fmt.Errorf("%w: %s", ErrMacPresentationOptionsInvalid, C.GoString(message))
}

func macPresentationOptions() MacPresentationOptions {
	return InvokeSyncWithResult(func() MacPresentationOptions {
		return MacPresentationOptions(C.presentationOptionsGet())
	})
}

// macApplyPresentationOptionsAtLaunch applies MacOptions.PresentationOptions
// once the application has finished launching. It runs on the application
// thread from the ApplicationDidFinishLaunching hook in application_darwin.go.
// An invalid combination is reported through the application error handler
// and the default presentation is kept.
func macApplyPresentationOptionsAtLaunch(app *App) {
	if app == nil {
		return
	}
	options := app.options.Mac.PresentationOptions
	if options == MacPresentationDefault {
		return
	}
	if err := options.Validate(); err != nil {
		app.handleError(fmt.Errorf("MacOptions.PresentationOptions: %w", err))
		return
	}
	if message := C.presentationOptionsSet(C.ulong(options)); message != nil {
		app.handleError(fmt.Errorf("MacOptions.PresentationOptions: %w: %s", ErrMacPresentationOptionsInvalid, C.GoString(message)))
		C.free(unsafe.Pointer(message))
	}
}
