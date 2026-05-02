//go:build ios

package events

/*
#cgo CFLAGS:  -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit

#include <stdlib.h>
#include <stdbool.h>

#define MAX_EVENTS 100

bool hasListener[MAX_EVENTS] = {false};

void registerListener(unsigned int event) {
	hasListener[event] = true;
}

bool hasListeners(unsigned int event) {
	//return hasListener[event];
	return true;
}

*/
import "C"