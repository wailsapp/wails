//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#include <stdlib.h>
#include "sound_manager_darwin.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

const macSystemSoundsDir = "/System/Library/Sounds"

func systemSoundsDir() string {
	return macSystemSoundsDir
}

func soundBeep() {
	InvokeAsync(func() {
		C.soundBeep()
	})
}

func soundResultError(code C.int, what string) error {
	switch code {
	case C.WailsSoundOK:
		return nil
	case C.WailsSoundNotFound:
		return fmt.Errorf("sound %s: not found or not decodable", what)
	default:
		return fmt.Errorf("sound %s: playback failed", what)
	}
}

func soundPlayNamed(name string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	code := InvokeSyncWithResult(func() C.int {
		return C.soundPlayNamed(cName)
	})
	return soundResultError(code, fmt.Sprintf("%q", name))
}

func soundPlayFile(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	code := InvokeSyncWithResult(func() C.int {
		return C.soundPlayFile(cPath)
	})
	return soundResultError(code, fmt.Sprintf("%q", path))
}

func soundPlayData(data []byte) error {
	if len(data) == 0 {
		return errors.New("sound data is empty")
	}
	code := InvokeSyncWithResult(func() C.int {
		// NSSound copies the bytes into an NSData before this returns.
		return C.soundPlayData(unsafe.Pointer(&data[0]), C.int(len(data)))
	})
	return soundResultError(code, "data")
}
