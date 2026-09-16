//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework AVFoundation -framework CoreLocation -framework ApplicationServices -framework IOKit -framework UserNotifications -mmacosx-version-min=10.13

#include "permissions_manager_darwin.h"
*/
import "C"

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// permissionRequestTimeout bounds how long Request waits for the system.
// TCC prompts can sit on screen for a long time, so this is generous; it
// mainly catches requests macOS drops silently because the Info.plist usage
// description is missing.
const permissionRequestTimeout = 5 * time.Minute

type macosPermissions struct {
	app      *App
	lock     sync.Mutex
	nextID   uint
	pending  map[uint]chan PermissionStatus
	fdaProbe string
}

func newPermissionsImpl(app *App) permissionsImpl {
	return &macosPermissions{
		app:     app,
		pending: map[uint]chan PermissionStatus{},
	}
}

// macosPermissionsInstance returns the live manager so the C callback can
// find the pending request.
func macosPermissionsInstance() *macosPermissions {
	if globalApplication == nil || globalApplication.Permissions == nil {
		return nil
	}
	impl, _ := globalApplication.Permissions.impl.(*macosPermissions)
	return impl
}

func permissionKindCode(kind PermissionKind) (C.int, bool) {
	switch kind {
	case PermissionKindCamera:
		return C.WailsPermissionCamera, true
	case PermissionKindMicrophone:
		return C.WailsPermissionMicrophone, true
	case PermissionKindScreenRecording:
		return C.WailsPermissionScreenRecording, true
	case PermissionKindAccessibility:
		return C.WailsPermissionAccessibility, true
	case PermissionKindLocation:
		return C.WailsPermissionLocation, true
	case PermissionKindNotifications:
		return C.WailsPermissionNotifications, true
	case PermissionKindInputMonitoring:
		return C.WailsPermissionInputMonitoring, true
	case PermissionKindFullDiskAccess:
		return C.WailsPermissionFullDiskAccess, true
	}
	return 0, false
}

func permissionStatusFromCode(code C.int) PermissionStatus {
	switch code {
	case C.WailsPermissionStatusNotDetermined:
		return PermissionStatusNotDetermined
	case C.WailsPermissionStatusDenied:
		return PermissionStatusDenied
	case C.WailsPermissionStatusAuthorized:
		return PermissionStatusAuthorized
	case C.WailsPermissionStatusRestricted:
		return PermissionStatusRestricted
	}
	return PermissionStatusUnsupported
}

func (p *macosPermissions) status(kind PermissionKind) PermissionStatus {
	if kind == PermissionKindFullDiskAccess {
		return p.fullDiskAccessStatus()
	}
	code, ok := permissionKindCode(kind)
	if !ok {
		return PermissionStatusUnsupported
	}
	return permissionStatusFromCode(C.wailsPermissionStatus(code))
}

// fullDiskAccessStatus probes a TCC-protected file. Reading it succeeds
// only with Full Disk Access; EPERM means the app does not have it.
func (p *macosPermissions) fullDiskAccessStatus() PermissionStatus {
	home, err := os.UserHomeDir()
	if err != nil {
		return PermissionStatusUnsupported
	}
	probes := []string{
		filepath.Join(home, "Library", "Application Support", "com.apple.TCC", "TCC.db"),
		filepath.Join(home, "Library", "Safari", "Bookmarks.plist"),
	}
	if p.fdaProbe != "" {
		probes = []string{p.fdaProbe}
	}
	for _, probe := range probes {
		file, err := os.Open(probe)
		if err == nil {
			_ = file.Close()
			return PermissionStatusAuthorized
		}
		if errors.Is(err, os.ErrPermission) {
			return PermissionStatusDenied
		}
	}
	return PermissionStatusUnsupported
}

func (p *macosPermissions) request(kind PermissionKind) (PermissionStatus, error) {
	if kind == PermissionKindFullDiskAccess {
		return p.fullDiskAccessStatus(), ErrPermissionNotRequestable
	}
	code, ok := permissionKindCode(kind)
	if !ok {
		return PermissionStatusUnsupported, ErrPermissionUnknownKind
	}
	if p.app != nil && p.app.impl != nil && p.app.impl.isOnMainThread() {
		return p.status(kind), ErrPermissionRequestOnMainThread
	}

	p.lock.Lock()
	p.nextID++
	requestID := p.nextID
	answer := make(chan PermissionStatus, 1)
	p.pending[requestID] = answer
	p.lock.Unlock()

	result := InvokeSyncWithResult(func() C.int {
		return C.wailsPermissionRequest(code, C.uint(requestID))
	})
	if result != C.WailsPermissionRequestPending {
		p.forget(requestID)
		return permissionStatusFromCode(result), nil
	}

	select {
	case status := <-answer:
		return status, nil
	case <-time.After(permissionRequestTimeout):
		p.forget(requestID)
		return p.status(kind), ErrPermissionRequestTimeout
	}
}

func (p *macosPermissions) forget(requestID uint) {
	p.lock.Lock()
	delete(p.pending, requestID)
	p.lock.Unlock()
}

func (p *macosPermissions) deliver(requestID uint, status PermissionStatus) {
	p.lock.Lock()
	answer, ok := p.pending[requestID]
	delete(p.pending, requestID)
	p.lock.Unlock()
	if ok {
		answer <- status
	}
}

func (p *macosPermissions) openSystemSettings(kind PermissionKind) error {
	code, ok := permissionKindCode(kind)
	if !ok {
		return ErrPermissionUnknownKind
	}
	opened := InvokeSyncWithResult(func() bool {
		return bool(C.wailsPermissionOpenSettings(code))
	})
	if !opened {
		return errors.New("permissions: could not open System Settings")
	}
	return nil
}

//export permissionsRequestResult
func permissionsRequestResult(requestID C.uint, status C.int) {
	impl := macosPermissionsInstance()
	if impl == nil {
		return
	}
	impl.deliver(uint(requestID), permissionStatusFromCode(status))
}

//export permissionsMediaCaptureDecision
func permissionsMediaCaptureDecision(windowID C.uint, captureType C.int) C.int {
	// WKMediaCaptureType: 0 camera, 1 microphone, 2 camera and microphone.
	needVideo := captureType == 0 || captureType == 2
	needAudio := captureType == 1 || captureType == 2
	perms := windowPermissions(uint(windowID))
	return C.int(mediaCapturePermissionDecision(perms, needAudio, needVideo))
}
