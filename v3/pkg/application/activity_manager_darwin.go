//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13

#include <stdlib.h>
#include "activity_manager_darwin.h"
*/
import "C"

import (
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"
)

type macosActivity struct {
	app *App
}

func newActivityImpl(app *App) activityImpl {
	return &macosActivity{app: app}
}

func (m *macosActivity) publish(id uint64, activity UserActivity) error {
	encoded, err := json.Marshal(activity)
	if err != nil {
		return fmt.Errorf("activity: cannot encode activity: %w", err)
	}
	cJSON := C.CString(string(encoded))
	defer C.free(unsafe.Pointer(cJSON))
	InvokeSync(func() {
		C.wailsActivityPublish(C.ulonglong(id), cJSON)
	})
	return nil
}

func (m *macosActivity) update(id uint64, userInfo map[string]any) error {
	if userInfo == nil {
		userInfo = map[string]any{}
	}
	encoded, err := json.Marshal(userInfo)
	if err != nil {
		return fmt.Errorf("activity: cannot encode UserInfo: %w", err)
	}
	cJSON := C.CString(string(encoded))
	defer C.free(unsafe.Pointer(cJSON))
	ok := InvokeSyncWithResult(func() bool {
		return bool(C.wailsActivityUpdate(C.ulonglong(id), cJSON))
	})
	if !ok {
		return ErrActivityInvalidated
	}
	return nil
}

func (m *macosActivity) invalidate(id uint64) {
	InvokeSync(func() {
		C.wailsActivityInvalidate(C.ulonglong(id))
	})
}

// decodeUserActivity parses the JSON produced by wailsUserActivityJSON.
func decodeUserActivity(raw string) (UserActivity, bool) {
	var activity UserActivity
	if raw == "" || raw == "null" {
		return activity, false
	}
	if err := json.Unmarshal([]byte(raw), &activity); err != nil {
		return activity, false
	}
	return activity, true
}

// The four exports below back the NSApplicationDelegate user activity
// methods in application_darwin_delegate.m. They run on the main thread
// while AppKit waits for the answer.

//export activityWillContinue
func activityWillContinue(activityType *C.char) C.bool {
	if globalApplication == nil || globalApplication.Activity == nil {
		return C.bool(false)
	}
	return C.bool(globalApplication.Activity.willContinue(C.GoString(activityType)))
}

//export activityContinue
func activityContinue(activityJSON *C.char) C.bool {
	activity, ok := decodeUserActivity(C.GoString(activityJSON))
	if !ok || globalApplication == nil || globalApplication.Activity == nil {
		return C.bool(false)
	}
	handled := false
	// Universal links share the custom URL scheme path so apps have one
	// place to handle every URL that opens them.
	if activity.Type == UserActivityTypeBrowsingWeb && activity.WebpageURL != "" {
		cURL := C.CString(activity.WebpageURL)
		HandleOpenURL(cURL)
		C.free(unsafe.Pointer(cURL))
		handled = true
	}
	if globalApplication.Activity.continueActivity(activity) {
		handled = true
	}
	return C.bool(handled)
}

//export activityDidFail
func activityDidFail(activityType *C.char, message *C.char) {
	if globalApplication == nil || globalApplication.Activity == nil {
		return
	}
	kind := C.GoString(activityType)
	err := errors.New(C.GoString(message))
	go globalApplication.Activity.failed(kind, err)
}

//export activityDidUpdate
func activityDidUpdate(activityJSON *C.char) {
	activity, ok := decodeUserActivity(C.GoString(activityJSON))
	if !ok || globalApplication == nil || globalApplication.Activity == nil {
		return
	}
	go globalApplication.Activity.updated(activity)
}
