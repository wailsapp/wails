//go:build android

package events

// Android events implementation

var Android = newAndroidEvents()

type androidEvents struct {
	ActivityCreated ApplicationEventType
}

func newAndroidEvents() androidEvents {
	return androidEvents{
		ActivityCreated: 1259,
	}
}

var hasListener = make(map[uint]bool)

func registerListener(event uint) {
	hasListener[event] = true
}

func hasListeners(event uint) bool {
	return hasListener[event]
}
