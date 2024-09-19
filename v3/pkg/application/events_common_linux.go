//go:build linux

package application

import "github.com/wailsapp/wails/v3/pkg/events"

var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
	events.Linux.ApplicationStartup: events.Common.ApplicationStarted,
	events.Linux.SystemThemeChanged: events.Common.ThemeChanged,
}

func (a *linuxApp) setupCommonEvents() {
	for sourceEvent, targetEvent := range commonApplicationEventMap {
		sourceEvent := sourceEvent
		targetEvent := targetEvent
		a.parent.OnApplicationEvent(sourceEvent, func(event *ApplicationEvent) {
			event.Id = uint(targetEvent)
			applicationEvents <- event
		})
	}
}
