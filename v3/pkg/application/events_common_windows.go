//go:build windows

package application

import "github.com/wailsapp/wails/v3/pkg/events"

var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
	events.Windows.SystemThemeChanged: events.Common.ThemeChanged,
	events.Windows.ApplicationStarted: events.Common.ApplicationStarted,
}
