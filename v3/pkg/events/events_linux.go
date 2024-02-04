//go:build linux

package events

/*
#include "events_linux.h"
#include <stdlib.h>
#include <stdbool.h>

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
