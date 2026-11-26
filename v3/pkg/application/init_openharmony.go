//go:build openharmony && cgo && !server

package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// commonApplicationEventMap forwards OHOS platform events to their common
// counterparts (same pattern as events_common_android.go). setupCommonEvents
// copies the source event's context, so any data attached on the OHOS side
// (battery level, theme, …) rides along to the common event.
var commonApplicationEventMap = map[events.ApplicationEventType]events.ApplicationEventType{
	events.Openharmony.AbilityCreated:       events.Common.ApplicationStarted,
	events.Openharmony.ApplicationLowMemory: events.Common.LowMemory,
	events.Openharmony.BatteryChanged:       events.Common.BatteryChanged,
	events.Openharmony.NetworkChanged:       events.Common.NetworkChanged,
	events.Openharmony.ThemeChanged:         events.Common.ThemeChanged,
	events.Openharmony.ScreenLocked:         events.Common.ScreenLocked,
	events.Openharmony.ScreenUnlocked:       events.Common.ScreenUnlocked,
}

// openharmonySystemEventTypes maps the name the ArkTS bridge passes to
// WailsOhEmitSystemEvent onto the typed openharmony: application event (mirrors
// androidSystemEventTypes).
var openharmonySystemEventTypes = map[string]events.ApplicationEventType{
	"batteryChanged": events.Openharmony.BatteryChanged,
	"networkChanged": events.Openharmony.NetworkChanged,
	"themeChanged":   events.Openharmony.ThemeChanged,
	"screenLocked":   events.Openharmony.ScreenLocked,
	"screenUnlocked": events.Openharmony.ScreenUnlocked,
}

// registerOpenharmonyListener registers a Go-side listener for a native OHOS
// event id (mirrors registerAndroidListener). The bridge notifies Go via
// WailsOhEmitSystemEvent; the shell forwards every lifecycle event it cares
// about, so registration itself is a no-op on the native side.
func registerOpenharmonyListener(eventID uint) {
	_ = eventID
}

func openharmonyDebugLogf(format string, args ...any) {
	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		slog.Debug(fmt.Sprintf(format, args...))
	}
}
