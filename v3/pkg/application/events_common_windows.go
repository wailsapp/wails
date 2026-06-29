//go:build windows

package application

import "github.com/wailsapp/wails/v3/pkg/events"

var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
	events.Windows.SystemThemeChanged:    events.Common.ThemeChanged,
	events.Windows.ApplicationStarted:    events.Common.ApplicationStarted,
	events.Windows.APMSuspend:            events.Common.SystemWillSleep,
	events.Windows.APMResumeAutomatic:    events.Common.SystemDidWake,
}

func (m *windowsApp) setupCommonEvents() {
	for sourceEvent, targetEvent := range commonApplicationEventMap {
		sourceEvent := sourceEvent
		targetEvent := targetEvent
		m.parent.Event.OnApplicationEvent(sourceEvent, func(event *ApplicationEvent) {
			applicationEvents <- &ApplicationEvent{
				Id:  uint(targetEvent),
				ctx: event.ctx,
			}
		})
	}
}
