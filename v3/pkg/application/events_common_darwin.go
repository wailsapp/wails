//go:build darwin && !ios

package application

import "github.com/wailsapp/wails/v3/pkg/events"

var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
	events.Mac.ApplicationDidFinishLaunching: events.Common.ApplicationStarted,
	events.Mac.ApplicationDidChangeTheme:     events.Common.ThemeChanged,
}

func (m *macosApp) setupCommonEvents() {
	for sourceEvent, targetEvent := range commonApplicationEventMap {
		sourceEvent := sourceEvent
		targetEvent := targetEvent
		m.parent.Event.OnApplicationEvent(sourceEvent, func(event *ApplicationEvent) {
			event.Id = uint(targetEvent)
			applicationEvents <- event
		})
	}

	// Handle dock icon click (applicationShouldHandleReopen) to show windows
	// when there are no visible windows. This provides the expected macOS UX
	// where clicking the dock icon shows a hidden app's window.
	// Issue #4583: Apps with StartHidden: true should show when dock icon is clicked.
	m.parent.Event.OnApplicationEvent(events.Mac.ApplicationShouldHandleReopen, func(event *ApplicationEvent) {
		if !event.Context().HasVisibleWindows() {
			// Show all windows that are not visible
			for _, window := range m.parent.Window.GetAll() {
				if !window.IsVisible() {
					window.Show()
				}
			}
		}
	})
}
