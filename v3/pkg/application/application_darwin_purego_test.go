//go:build darwin && purego

package application

import (
	"testing"

	"github.com/wailsapp/wails/v3/pkg/events"
)

func TestMacosAppOnRecordsListener(t *testing.T) {
	eventID := uint(events.Mac.ApplicationDidFinishLaunching)
	other := uint(events.Mac.ApplicationDidHide)
	if other == eventID {
		t.Skip("test event ids must differ")
	}

	appEventListeners.Delete(eventID)
	appEventListeners.Delete(other)
	t.Cleanup(func() {
		appEventListeners.Delete(eventID)
		appEventListeners.Delete(other)
	})

	app := &macosApp{}
	app.on(eventID)

	if _, ok := appEventListeners.Load(eventID); !ok {
		t.Fatalf("macosApp.on(%d) did not record the listener", eventID)
	}
	if _, ok := appEventListeners.Load(other); ok {
		t.Fatalf("macosApp.on(%d) recorded unrelated event %d", eventID, other)
	}

	// hasAppListener mirrors the cgo hasListeners, which returns true today.
	if !hasAppListener(eventID) {
		t.Fatalf("hasAppListener(%d) = false, want true to mirror the cgo hasListeners", eventID)
	}
}
