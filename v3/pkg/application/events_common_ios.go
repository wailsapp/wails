//go:build ios

package application

import "github.com/wailsapp/wails/v3/pkg/events"

// Map platform events → common events (same pattern as macOS & others).
// setupCommonEvents copies the source event's context, so any data attached on
// the iOS side (battery level, theme, …) rides along to the common event.
var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
    events.IOS.ApplicationDidFinishLaunching:      events.Common.ApplicationStarted,
    events.IOS.ApplicationDidReceiveMemoryWarning: events.Common.LowMemory,
    events.IOS.BatteryChanged:                     events.Common.BatteryChanged,
    events.IOS.NetworkChanged:                     events.Common.NetworkChanged,
    events.IOS.ThemeChanged:                       events.Common.ThemeChanged,
    events.IOS.ScreenLocked:                       events.Common.ScreenLocked,
    events.IOS.ScreenUnlocked:                     events.Common.ScreenUnlocked,
}

// setupCommonEvents forwards iOS platform events to their common counterparts
func (i *iosApp) setupCommonEvents() {
    for sourceEvent, targetEvent := range commonApplicationEventMap {
        sourceEvent := sourceEvent
        targetEvent := targetEvent
        i.parent.Event.OnApplicationEvent(sourceEvent, func(event *ApplicationEvent) {
            // Log the forwarding so we can see every emitted event in iOS NSLog
            iosConsoleLogf("info", " [events_common_ios.go] Forwarding iOS event %d → common %d", sourceEvent, targetEvent)
            applicationEvents <- &ApplicationEvent{
                Id:  uint(targetEvent),
                ctx: event.ctx,
            }
        })
    }
}