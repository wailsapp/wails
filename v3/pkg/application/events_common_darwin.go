//go:build darwin

package application

import "github.com/wailsapp/wails/v3/pkg/events"

var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
	events.Mac.ApplicationDidFinishLaunching: events.Common.ApplicationStarted,
}

func (m *macosApp) setupCommonEvents() {
	for sourceEvent, targetEvent := range commonApplicationEventMap {
		m.parent.On(sourceEvent, func() {
			applicationEvents <- uint(targetEvent)
		})
	}
}
