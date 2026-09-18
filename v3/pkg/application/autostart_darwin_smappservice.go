//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.15 -x objective-c -Wno-unguarded-availability-new
#cgo LDFLAGS: -framework Foundation -framework ServiceManagement

#include <stdlib.h> // free
#include <string.h> // strdup
#import <Foundation/Foundation.h>
#import <ServiceManagement/ServiceManagement.h>

// Return codes shared with the Go side.
enum {
	SMAS_OK              = 0,
	SMAS_UNAVAILABLE     = 1, // SMAppService class not present (pre macOS 13)
	SMAS_NOT_REGISTERED  = 2, // unregister called when nothing was registered
	SMAS_REQUIRES_APPROVAL = 3, // user disabled it in System Settings
	SMAS_ERROR           = 4, // generic failure; *outMsg populated
};

static int smAppServiceRegister(char** outMsg) {
	if (@available(macOS 13.0, *)) {
		@autoreleasepool {
			SMAppService *svc = [SMAppService mainAppService];
			NSError *err = nil;
			if ([svc registerAndReturnError:&err]) {
				return SMAS_OK;
			}
			if (err != nil) {
				*outMsg = strdup([[err localizedDescription] UTF8String]);
			}
			return SMAS_ERROR;
		}
	}
	return SMAS_UNAVAILABLE;
}

static int smAppServiceUnregister(char** outMsg) {
	if (@available(macOS 13.0, *)) {
		@autoreleasepool {
			SMAppService *svc = [SMAppService mainAppService];
			if (svc.status == SMAppServiceStatusNotRegistered ||
			    svc.status == SMAppServiceStatusNotFound) {
				return SMAS_NOT_REGISTERED;
			}
			NSError *err = nil;
			if ([svc unregisterAndReturnError:&err]) {
				return SMAS_OK;
			}
			if (err != nil) {
				*outMsg = strdup([[err localizedDescription] UTF8String]);
			}
			return SMAS_ERROR;
		}
	}
	return SMAS_UNAVAILABLE;
}

// smAppServiceStatus: 0 = unavailable, 1 = not registered / not found,
//                     2 = enabled, 3 = requires approval.
static int smAppServiceStatus(void) {
	if (@available(macOS 13.0, *)) {
		@autoreleasepool {
			SMAppService *svc = [SMAppService mainAppService];
			switch (svc.status) {
				case SMAppServiceStatusEnabled:           return 2;
				case SMAppServiceStatusRequiresApproval:  return 3;
				default:                                  return 1;
			}
		}
	}
	return 0;
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

var (
	errSMAppServiceUnavailable    = errors.New("SMAppService unavailable on this macOS")
	errSMAppServiceNotRegistered  = errors.New("SMAppService not registered")
	errSMAppServiceRequiresApproval = errors.New("SMAppService requires user approval in System Settings")
)

func smAppServiceRegister() error {
	var cMsg *C.char
	rc := C.smAppServiceRegister(&cMsg)
	if cMsg != nil {
		defer C.free(unsafe.Pointer(cMsg))
	}
	switch rc {
	case 0:
		return nil
	case 1:
		return errSMAppServiceUnavailable
	default:
		if cMsg != nil {
			return errors.New(C.GoString(cMsg))
		}
		return errors.New("SMAppService register failed")
	}
}

func smAppServiceUnregister() error {
	var cMsg *C.char
	rc := C.smAppServiceUnregister(&cMsg)
	if cMsg != nil {
		defer C.free(unsafe.Pointer(cMsg))
	}
	switch rc {
	case 0:
		return nil
	case 1:
		return errSMAppServiceUnavailable
	case 2:
		return errSMAppServiceNotRegistered
	default:
		if cMsg != nil {
			return errors.New(C.GoString(cMsg))
		}
		return errors.New("SMAppService unregister failed")
	}
}

func smAppServiceIsEnabled() (bool, error) {
	switch C.smAppServiceStatus() {
	case 0:
		return false, errSMAppServiceUnavailable
	case 2:
		return true, nil
	case 3:
		return false, errSMAppServiceRequiresApproval
	default:
		return false, nil
	}
}
