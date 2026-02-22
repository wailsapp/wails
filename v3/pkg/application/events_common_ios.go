//go:build ios

package application

import "github.com/wailsapp/wails/v3/pkg/events"

// Map platform events → common events (same pattern as macOS & others)
var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
    events.IOS.ApplicationDidFinishLaunching: events.Common.ApplicationStarted,
}

// setupCommonEvents forwards iOS platform events to their common counterparts
func (i *iosApp) setupCommonEvents() {
    for sourceEvent, targetEvent := range commonApplicationEventMap {
        sourceEvent := sourceEvent
        targetEvent := targetEvent
        i.parent.Event.OnApplicationEvent(sourceEvent, func(event *ApplicationEvent) {
            event.Id = uint(targetEvent)
            // Log the forwarding so we can see every emitted event in iOS NSLog
            iosConsoleLogf("info", " [events_common_ios.go] Forwarding iOS event %d → common %d", sourceEvent, targetEvent)
            applicationEvents <- event
        })
    }
}