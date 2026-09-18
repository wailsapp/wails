//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.13

#import <Cocoa/Cocoa.h>
#include <stdbool.h>
#include <stdlib.h>

static bool wailsSuddenTerminationDefault(void) {
	@autoreleasepool {
		id value = [[NSBundle mainBundle] objectForInfoDictionaryKey:@"NSSupportsSuddenTermination"];
		return value != nil && [value respondsToSelector:@selector(boolValue)] && [value boolValue];
	}
}

static void wailsDisableTermination(const char *reason) {
	@autoreleasepool {
		NSProcessInfo *info = [NSProcessInfo processInfo];
		[info disableSuddenTermination];
		[info disableAutomaticTermination:[NSString stringWithUTF8String:reason]];
	}
}

static void wailsEnableTermination(const char *reason) {
	@autoreleasepool {
		NSProcessInfo *info = [NSProcessInfo processInfo];
		[info enableAutomaticTermination:[NSString stringWithUTF8String:reason]];
		[info enableSuddenTermination];
	}
}

static void wailsSetSuddenTerminationEnabled(bool enabled) {
	if (enabled) {
		[[NSProcessInfo processInfo] enableSuddenTermination];
	} else {
		[[NSProcessInfo processInfo] disableSuddenTermination];
	}
}
*/
import "C"

import "unsafe"

type macosLifecycle struct {
	app *App
}

func newLifecycleImpl(app *App) lifecycleImpl {
	return &macosLifecycle{app: app}
}

func (l *macosLifecycle) suddenTerminationDefault() bool {
	return bool(C.wailsSuddenTerminationDefault())
}

func (l *macosLifecycle) disableTermination(reason string) {
	cReason := C.CString(reason)
	defer C.free(unsafe.Pointer(cReason))
	C.wailsDisableTermination(cReason)
}

func (l *macosLifecycle) enableTermination(reason string) {
	cReason := C.CString(reason)
	defer C.free(unsafe.Pointer(cReason))
	C.wailsEnableTermination(cReason)
}

func (l *macosLifecycle) setSuddenTerminationEnabled(enabled bool) {
	C.wailsSetSuddenTerminationEnabled(C.bool(enabled))
}
