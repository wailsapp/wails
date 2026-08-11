//go:build android

package application

import "github.com/wailsapp/wails/v3/pkg/events"

// Map platform events → common events (same pattern as macOS & others).
// setupCommonEvents copies the source event's context, so any data attached on
// the Android side (battery level, theme, …) rides along to the common event.
var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
	events.Android.ActivityCreated:      events.Common.ApplicationStarted,
	events.Android.ApplicationLowMemory: events.Common.LowMemory,
	events.Android.BatteryChanged:       events.Common.BatteryChanged,
	events.Android.NetworkChanged:       events.Common.NetworkChanged,
	events.Android.ThemeChanged:         events.Common.ThemeChanged,
	events.Android.ScreenLocked:         events.Common.ScreenLocked,
	events.Android.ScreenUnlocked:       events.Common.ScreenUnlocked,
}

// setupCommonEvents forwards Android platform events to their common counterparts
func (a *androidApp) setupCommonEvents() {
	for sourceEvent, targetEvent := range commonApplicationEventMap {
		sourceEvent := sourceEvent
		targetEvent := targetEvent
		a.parent.Event.OnApplicationEvent(sourceEvent, func(event *ApplicationEvent) {
			androidDebugLogf("[events_common_android.go] forwarding android event %d → common %d", sourceEvent, targetEvent)
			applicationEvents <- &ApplicationEvent{
				Id:  uint(targetEvent),
				ctx: event.ctx,
			}
		})
	}
}
