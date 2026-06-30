//go:build linux && !android && !server

package application

import "github.com/wailsapp/wails/v3/pkg/events"

var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
	events.Linux.ApplicationStartup: events.Common.ApplicationStarted,
	events.Linux.SystemThemeChanged: events.Common.ThemeChanged,
	events.Linux.SystemWillSleep:    events.Common.SystemWillSleep,
	events.Linux.SystemDidWake:      events.Common.SystemDidWake,
}

func (a *linuxApp) setupCommonEvents() {
	for sourceEvent, targetEvent := range commonApplicationEventMap {
		sourceEvent := sourceEvent
		targetEvent := targetEvent
		a.parent.Event.OnApplicationEvent(sourceEvent, func(event *ApplicationEvent) {
			applicationEvents <- &ApplicationEvent{
				Id:  uint(targetEvent),
				ctx: event.ctx,
			}
		})
	}
}
