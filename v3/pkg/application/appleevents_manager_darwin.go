//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework CoreServices

#include <stdlib.h>
#include "appleevents_manager_darwin.h"
*/
import "C"

import (
	"errors"
	"sync"
	"unsafe"
)

// appleEventDelivery is a suspended Apple Event waiting for its handler.
type appleEventDelivery struct {
	eventJSON    string
	suspensionID unsafe.Pointer
}

// appleEventDeliveries carries suspended events from the main thread to the
// drain loop registered in init, which runs each handler on its own
// goroutine and resumes the event on the main thread when it returns.
var appleEventDeliveries = make(chan appleEventDelivery, 32)

//export appleEventsDeliver
func appleEventsDeliver(eventJSON *C.char, suspensionID unsafe.Pointer) *C.char {
	payload := C.GoString(eventJSON)
	if suspensionID != nil {
		appleEventDeliveries <- appleEventDelivery{eventJSON: payload, suspensionID: suspensionID}
		return nil
	}
	// The event could not be suspended: run the handler inline on the main
	// thread so the sender still gets a reply. Freed by the caller.
	if globalApplication == nil || globalApplication.AppleEvents == nil {
		return nil
	}
	return C.CString(globalApplication.AppleEvents.dispatch(payload))
}

func handleAppleEventDelivery(app *App, delivery appleEventDelivery) {
	reply := app.AppleEvents.dispatch(delivery.eventJSON)
	InvokeAsync(func() {
		cReply := C.CString(reply)
		defer C.free(unsafe.Pointer(cReply))
		C.appleEventsResume(delivery.suspensionID, cReply)
	})
}

func init() {
	registerChromeEventLoop(func(app *App) {
		for {
			delivery := <-appleEventDeliveries
			go handleAppleEventDelivery(app, delivery)
		}
	})
}

// appleEventsNativeRegistrations tracks which class/ID pairs already have the
// shared handler object installed, so re-registering a Go handler does not
// touch NSAppleEventManager again.
var (
	appleEventsNativeLock          sync.Mutex
	appleEventsNativeRegistrations = map[appleEventKey]bool{}
)

// appleEventsRegistration installs one native registration on the main
// thread. It implements runnable so registrations made before App.Run are
// deferred until the application starts.
type appleEventsRegistration struct {
	class string
	id    string
}

func (r appleEventsRegistration) Run() {
	cClass := C.CString(r.class)
	defer C.free(unsafe.Pointer(cClass))
	cID := C.CString(r.id)
	defer C.free(unsafe.Pointer(cID))
	InvokeSync(func() {
		C.appleEventsRegister(cClass, cID)
	})
}

func appleEventsRegisterNative(m *AppleEventsManager, eventClass, eventID string) error {
	if m.app == nil {
		return errors.New("apple events: manager is not attached to an application")
	}
	key := appleEventKey{class: eventClass, id: eventID}
	appleEventsNativeLock.Lock()
	already := appleEventsNativeRegistrations[key]
	appleEventsNativeRegistrations[key] = true
	appleEventsNativeLock.Unlock()
	if already {
		return nil
	}
	m.app.runOrDeferToAppRun(appleEventsRegistration{class: eventClass, id: eventID})
	return nil
}

// appleEventsSendNative runs on the calling goroutine: AESendMessage waits
// on its own reply port, so the main thread stays free to service other
// events (including replies that require it).
func appleEventsSendNative(target, eventClass, eventID, directJSON string) (string, error) {
	cTarget := C.CString(target)
	defer C.free(unsafe.Pointer(cTarget))
	cClass := C.CString(eventClass)
	defer C.free(unsafe.Pointer(cClass))
	cID := C.CString(eventID)
	defer C.free(unsafe.Pointer(cID))
	cDirect := C.CString(directJSON)
	defer C.free(unsafe.Pointer(cDirect))

	var cError *C.char
	cReply := C.appleEventsSend(cTarget, cClass, cID, cDirect, &cError)
	if cReply == nil {
		message := "apple events: send failed"
		if cError != nil {
			message = C.GoString(cError)
			C.free(unsafe.Pointer(cError))
		}
		return "", errors.New(message)
	}
	defer C.free(unsafe.Pointer(cReply))
	return C.GoString(cReply), nil
}
