//go:build windows

package application

import "github.com/wailsapp/wails/v3/pkg/events"

var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
	events.Windows.SystemThemeChanged: events.Common.ThemeChanged,
}

func (m *windowsApp) setupCommonEvents() {
	for sourceEvent, targetEvent := range commonApplicationEventMap {
		m.parent.On(sourceEvent, func(event *Event) {
			event.Id = uint(targetEvent)
			applicationEvents <- event
		})
	}
}
