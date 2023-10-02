//go:build linux

package application

import "github.com/wailsapp/wails/v3/pkg/events"

var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
	events.Linux.SystemThemeChanged: events.Common.ThemeChanged,
}

func (m *linuxApp) setupCommonEvents() {
	for sourceEvent, targetEvent := range commonApplicationEventMap {
		sourceEvent := sourceEvent
		targetEvent := targetEvent
		m.parent.On(sourceEvent, func(event *Event) {
			event.Id = uint(targetEvent)
			applicationEvents <- event
		})
	}
}
