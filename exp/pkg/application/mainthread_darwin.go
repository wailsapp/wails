//go:build darwin

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa

#include "Cocoa/Cocoa.h"

extern void dispatchCallback(unsigned int);

void dispatch(unsigned int id) {
	dispatch_async(dispatch_get_main_queue(), ^{
		dispatchCallback(id);
	});
}

*/
import "C"
