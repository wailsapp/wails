//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework CoreSpotlight -framework UniformTypeIdentifiers

#include <stdlib.h>
#include "spotlight_manager_darwin.h"
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

// spotlightOperationTimeout bounds how long an index or delete call waits
// for CoreSpotlight to confirm.
const spotlightOperationTimeout = 60 * time.Second

// spotlightRequests maps in-flight index operations to the channel their
// completion handler answers on.
var (
	spotlightRequestsLock sync.Mutex
	spotlightNextRequest  uint64
	spotlightRequests     = map[uint64]chan error{}
)

//export spotlightOperationResult
func spotlightOperationResult(requestID C.ulonglong, errorMessage *C.char) {
	spotlightRequestsLock.Lock()
	answer, ok := spotlightRequests[uint64(requestID)]
	delete(spotlightRequests, uint64(requestID))
	spotlightRequestsLock.Unlock()
	if !ok {
		return
	}
	var err error
	if errorMessage != nil {
		err = errors.New(C.GoString(errorMessage))
	}
	// Buffered and read exactly once, so the CoreSpotlight queue never blocks.
	answer <- err
}

//export spotlightHandleContinueActivityC
func spotlightHandleContinueActivityC(activityType *C.char, userInfoJSON *C.char) C.bool {
	return C.bool(spotlightHandleContinueActivity(C.GoString(activityType), C.GoString(userInfoJSON)))
}

// spotlightRun starts a native operation on the main thread and waits for
// its completion callback on the calling goroutine.
func spotlightRun(what string, start func(requestID uint64)) error {
	answer := make(chan error, 1)
	spotlightRequestsLock.Lock()
	spotlightNextRequest++
	requestID := spotlightNextRequest
	spotlightRequests[requestID] = answer
	spotlightRequestsLock.Unlock()

	InvokeAsync(func() {
		start(requestID)
	})

	select {
	case err := <-answer:
		if err != nil {
			return fmt.Errorf("spotlight: %s: %w", what, err)
		}
		return nil
	case <-time.After(spotlightOperationTimeout):
		spotlightRequestsLock.Lock()
		delete(spotlightRequests, requestID)
		spotlightRequestsLock.Unlock()
		return fmt.Errorf("spotlight: %s timed out", what)
	}
}

func spotlightIsAvailable() bool {
	return bool(C.spotlightIsAvailable())
}

func spotlightIndex(itemsJSON string) error {
	cItems := C.CString(itemsJSON)
	return spotlightRun("index", func(requestID uint64) {
		defer C.free(unsafe.Pointer(cItems))
		C.spotlightIndexItems(C.ulonglong(requestID), cItems)
	})
}

func spotlightDelete(ids []string) error {
	payload, err := json.Marshal(ids)
	if err != nil {
		return fmt.Errorf("spotlight: %w", err)
	}
	cIDs := C.CString(string(payload))
	return spotlightRun("delete", func(requestID uint64) {
		defer C.free(unsafe.Pointer(cIDs))
		C.spotlightDeleteItems(C.ulonglong(requestID), cIDs)
	})
}

func spotlightDeleteDomain(domain string) error {
	payload, err := json.Marshal([]string{domain})
	if err != nil {
		return fmt.Errorf("spotlight: %w", err)
	}
	cDomains := C.CString(string(payload))
	return spotlightRun("delete domain", func(requestID uint64) {
		defer C.free(unsafe.Pointer(cDomains))
		C.spotlightDeleteDomains(C.ulonglong(requestID), cDomains)
	})
}

func spotlightDeleteAll() error {
	return spotlightRun("delete all", func(requestID uint64) {
		C.spotlightDeleteAllItems(C.ulonglong(requestID))
	})
}
