//go:build android

package events

// Android events implementation

const MAX_EVENTS = 100

var hasListener [MAX_EVENTS]bool

func registerListener(event uint) {
	if event < MAX_EVENTS {
		hasListener[event] = true
	}
}

func hasListeners(event uint) bool {
	if event < MAX_EVENTS {
		return hasListener[event]
	}
	return false
}
