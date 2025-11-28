//go:build android

package application

import "github.com/wailsapp/wails/v3/pkg/events"

// Map platform events → common events (same pattern as macOS & others)
var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
	events.Android.ActivityCreated: events.Common.ApplicationStarted,
}

// setupCommonEvents forwards Android platform events to their common counterparts
func (a *androidApp) setupCommonEvents() {
	for sourceEvent, targetEvent := range commonApplicationEventMap {
		sourceEvent := sourceEvent
		targetEvent := targetEvent
		a.parent.Event.OnApplicationEvent(sourceEvent, func(event *ApplicationEvent) {
			event.Id = uint(targetEvent)
			androidLogf("info", "[events_common_android.go] Forwarding Android event %d → common %d", sourceEvent, targetEvent)
			applicationEvents <- event
		})
	}
}
