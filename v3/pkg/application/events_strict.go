//go:build !production && strictevents

package application

import (
	"sync"
)

var knownUnregisteredEvents sync.Map

func eventRegistered(name string) {
	knownUnregisteredEvents.Delete(name)
}

func warnAboutUnregisteredEvent(name string) {
	// Perform a load first to avoid thrashing the map with unnecessary swaps.
	if _, known := knownUnregisteredEvents.Load(name); known {
		return
	}

	// Perform a swap to synchronize with concurrent emissions.
	if _, known := knownUnregisteredEvents.Swap(name, true); known {
		return
	}

	globalApplication.warning("unregistered event name '%s'", name)
}
