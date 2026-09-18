//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework Quartz -framework QuickLookThumbnailing

#include <stdlib.h>
#include "quicklook_manager_darwin.h"
*/
import "C"

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
	"unsafe"
)

// quickLookThumbnailTimeout bounds how long Thumbnail waits for the system
// to render a file.
const quickLookThumbnailTimeout = 60 * time.Second

type quickLookThumbnailOutcome struct {
	png []byte
	err error
}

// quickLookThumbnailRequests maps in-flight thumbnail requests to the
// channel their completion handler answers on.
var (
	quickLookThumbnailLock     sync.Mutex
	quickLookThumbnailNextID   uint64
	quickLookThumbnailRequests = map[uint64]chan quickLookThumbnailOutcome{}
)

//export quickLookThumbnailResult
func quickLookThumbnailResult(requestID C.ulonglong, png unsafe.Pointer, length C.int, errorMessage *C.char) {
	quickLookThumbnailLock.Lock()
	answer, ok := quickLookThumbnailRequests[uint64(requestID)]
	delete(quickLookThumbnailRequests, uint64(requestID))
	quickLookThumbnailLock.Unlock()
	if !ok {
		return
	}
	outcome := quickLookThumbnailOutcome{}
	if errorMessage != nil {
		outcome.err = errors.New(C.GoString(errorMessage))
	} else {
		outcome.png = C.GoBytes(png, length)
	}
	// The channel is buffered and read exactly once, so this never blocks
	// the Quick Look queue.
	answer <- outcome
}

func quickLookPreview(paths []string) error {
	payload, err := json.Marshal(paths)
	if err != nil {
		return fmt.Errorf("quick look: %w", err)
	}
	cPayload := C.CString(string(payload))
	defer C.free(unsafe.Pointer(cPayload))
	shown := InvokeSyncWithResult(func() bool {
		return bool(C.quickLookPreview(cPayload))
	})
	if !shown {
		return errors.New("quick look: the preview panel could not be shown")
	}
	return nil
}

func quickLookClosePreview() {
	InvokeSync(func() {
		C.quickLookClosePreview()
	})
}

func quickLookIsPreviewOpen() bool {
	return InvokeSyncWithResult(func() bool {
		return bool(C.quickLookIsPreviewOpen())
	})
}

func quickLookThumbnail(path string, opts ThumbnailOptions) ([]byte, error) {
	if !bool(C.quickLookThumbnailAvailable()) {
		return nil, ErrQuickLookNotSupported
	}
	answer := make(chan quickLookThumbnailOutcome, 1)
	quickLookThumbnailLock.Lock()
	quickLookThumbnailNextID++
	requestID := quickLookThumbnailNextID
	quickLookThumbnailRequests[requestID] = answer
	quickLookThumbnailLock.Unlock()

	cPath := C.CString(path)
	InvokeAsync(func() {
		defer C.free(unsafe.Pointer(cPath))
		C.quickLookThumbnail(C.ulonglong(requestID), cPath, C.int(opts.Width), C.int(opts.Height), C.double(opts.Scale), C.bool(opts.IconMode))
	})

	select {
	case outcome := <-answer:
		if outcome.err != nil {
			return nil, fmt.Errorf("quick look: thumbnail for %q: %w", path, outcome.err)
		}
		return outcome.png, nil
	case <-time.After(quickLookThumbnailTimeout):
		quickLookThumbnailLock.Lock()
		delete(quickLookThumbnailRequests, requestID)
		quickLookThumbnailLock.Unlock()
		return nil, fmt.Errorf("quick look: thumbnail for %q timed out", path)
	}
}
