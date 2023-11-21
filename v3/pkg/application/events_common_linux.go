//go:build linux

package application

import "github.com/wailsapp/wails/v3/pkg/events"

var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
	events.Linux.ApplicationStarted: events.Common.ApplicationStarted,
	events.Linux.SystemThemeChanged: events.Common.ThemeChanged,
}
