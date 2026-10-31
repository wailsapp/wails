//go:build darwin && purego && !ios && !server

package application

// CGO-free re-implementation of autostart_darwin_smappservice.go.
//
// The launchd-plist path (autostart_darwin.go) is shared and already pure Go;
// only the SMAppService (macOS 13+) login-item registration touches
// Objective-C. Here we drive [SMAppService mainAppService] and its
// register/unregister/status API through the purego runtime helpers instead of
// cgo.

import (
	"errors"
	"sync"

	"github.com/ebitengine/purego"
)

var (
	errSMAppServiceUnavailable      = errors.New("SMAppService unavailable on this macOS")
	errSMAppServiceNotRegistered    = errors.New("SMAppService not registered")
	errSMAppServiceRequiresApproval = errors.New("SMAppService requires user approval in System Settings")
)

// SMAppServiceStatus enum values (ServiceManagement/SMAppService.h).
const (
	smAppServiceStatusNotRegistered    = 0
	smAppServiceStatusEnabled          = 1
	smAppServiceStatusRequiresApproval = 2
	smAppServiceStatusNotFound         = 3
)

const frameworkServiceManagement = "/System/Library/Frameworks/ServiceManagement.framework/ServiceManagement"

var serviceManagementOnce sync.Once

// mainAppService returns [SMAppService mainAppService], or a nil id if the
// SMAppService class is unavailable (pre-macOS 13). The ServiceManagement
// framework is dlopen'd on first use so objc_getClass can resolve the class.
func mainAppService() id {
	serviceManagementOnce.Do(func() {
		_, _ = purego.Dlopen(frameworkServiceManagement, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	})
	cls := class("SMAppService")
	if cls.isNil() {
		return id(0)
	}
	return cls.send("mainAppService")
}

func smAppServiceRegister() error {
	var out error
	withAutoreleasePool(func() {
		svc := mainAppService()
		if svc.isNil() {
			out = errSMAppServiceUnavailable
			return
		}
		var errObj id
		if get[bool](svc, "registerAndReturnError:", &errObj) {
			return
		}
		if !errObj.isNil() {
			out = errors.New(errObj.send("localizedDescription").string())
			return
		}
		out = errors.New("SMAppService register failed")
	})
	return out
}

func smAppServiceUnregister() error {
	var out error
	withAutoreleasePool(func() {
		svc := mainAppService()
		if svc.isNil() {
			out = errSMAppServiceUnavailable
			return
		}
		status := get[int](svc, "status")
		if status == smAppServiceStatusNotRegistered || status == smAppServiceStatusNotFound {
			out = errSMAppServiceNotRegistered
			return
		}
		var errObj id
		if get[bool](svc, "unregisterAndReturnError:", &errObj) {
			return
		}
		if !errObj.isNil() {
			out = errors.New(errObj.send("localizedDescription").string())
			return
		}
		out = errors.New("SMAppService unregister failed")
	})
	return out
}

func smAppServiceIsEnabled() (bool, error) {
	var (
		enabled bool
		out     error
	)
	withAutoreleasePool(func() {
		svc := mainAppService()
		if svc.isNil() {
			out = errSMAppServiceUnavailable
			return
		}
		switch get[int](svc, "status") {
		case smAppServiceStatusEnabled:
			enabled = true
		case smAppServiceStatusRequiresApproval:
			out = errSMAppServiceRequiresApproval
		default:
			// NotRegistered / NotFound -> not enabled, no error.
		}
	})
	return enabled, out
}
