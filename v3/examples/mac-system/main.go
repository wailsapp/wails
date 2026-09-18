package main

import (
	"embed"
	"fmt"
	"log"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets
var assets embed.FS

// SystemService is bound to the frontend and wraps the Permissions, Power,
// Lifecycle and Environment managers.
type SystemService struct {
	app *application.App

	lock            sync.Mutex
	releaseSleep    func()
	releaseTermHold func()
}

// PermissionStatuses returns the status of every permission kind.
func (s *SystemService) PermissionStatuses() map[string]string {
	result := map[string]string{}
	for _, kind := range application.AllPermissionKinds() {
		result[kind.String()] = s.app.Permissions.Status(kind).String()
	}
	return result
}

// RequestPermission prompts for a permission and returns the new status.
func (s *SystemService) RequestPermission(kind string) (string, error) {
	status, err := s.app.Permissions.Request(application.PermissionKind(kind))
	if err != nil {
		return status.String(), err
	}
	return status.String(), nil
}

// OpenPermissionSettings opens the matching System Settings pane.
func (s *SystemService) OpenPermissionSettings(kind string) error {
	return s.app.Permissions.OpenSystemSettings(application.PermissionKind(kind))
}

// SetKeepAwake toggles a PreventSleep hold.
func (s *SystemService) SetKeepAwake(on bool, display bool) (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if on && s.releaseSleep == nil {
		release, err := s.app.Power.PreventSleep("mac-system example keep awake", application.PreventSleepOptions{Display: display})
		if err != nil {
			return s.app.Power.SleepHoldCount(), err
		}
		s.releaseSleep = release
	}
	if !on && s.releaseSleep != nil {
		s.releaseSleep()
		s.releaseSleep = nil
	}
	return s.app.Power.SleepHoldCount(), nil
}

// SetTerminationHold toggles a HoldTermination hold.
func (s *SystemService) SetTerminationHold(on bool) int {
	s.lock.Lock()
	defer s.lock.Unlock()
	if on && s.releaseTermHold == nil {
		s.releaseTermHold = s.app.Lifecycle.HoldTermination("mac-system example hold")
	}
	if !on && s.releaseTermHold != nil {
		s.releaseTermHold()
		s.releaseTermHold = nil
	}
	return s.app.Lifecycle.TerminationHoldCount()
}

// SetSuddenTermination toggles sudden termination for the process.
func (s *SystemService) SetSuddenTermination(enabled bool) bool {
	s.app.Lifecycle.SetSuddenTerminationEnabled(enabled)
	return s.app.Lifecycle.SuddenTerminationEnabled()
}

// PowerState returns the current power and thermal state.
func (s *SystemService) PowerState() application.PowerState {
	return s.app.Power.State()
}

// Accessibility returns the current accessibility settings.
func (s *SystemService) Accessibility() application.AccessibilitySettings {
	return s.app.Env.Accessibility()
}

// KeyboardLayout returns the active keyboard input source.
func (s *SystemService) KeyboardLayout() application.KeyboardLayout {
	return s.app.Env.KeyboardLayout()
}

// Locale returns the application's locale.
func (s *SystemService) Locale() application.LocaleInfo {
	return s.app.Env.Locale()
}

func main() {
	service := &SystemService{}

	app := application.New(application.Options{
		Name:        "System Integration",
		Description: "Permissions, power, termination, accessibility, keyboard layout and locale",
		Services: []application.Service{
			application.NewService(service),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	service.app = app

	// Forward the system events to the frontend as custom events carrying
	// the fresh state, so the page updates live.
	forward := func(event events.ApplicationEventType, name string, payload func() any) {
		app.Event.OnApplicationEvent(event, func(*application.ApplicationEvent) {
			app.Event.Emit(name, payload())
		})
	}
	forward(events.Mac.ApplicationDidChangePowerState, "system:power", func() any { return app.Power.State() })
	forward(events.Mac.ApplicationDidChangeThermalState, "system:power", func() any { return app.Power.State() })
	forward(events.Common.AccessibilitySettingsChanged, "system:accessibility", func() any { return app.Env.Accessibility() })
	forward(events.Mac.ApplicationDidChangeKeyboardLayout, "system:keyboard", func() any { return app.Env.KeyboardLayout() })
	forward(events.Mac.ApplicationDidChangeLocale, "system:locale", func() any { return app.Env.Locale() })

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "System Integration",
		Width:  900,
		Height: 760,
		URL:    "/",
		// Web content in this window may use the microphone without the
		// webview's own prompt; the camera still asks. See issue #6067.
		Permissions: map[application.PermissionType]application.Permission{
			application.PermissionMicrophone: application.PermissionAllow,
		},
	})

	if err := app.Run(); err != nil {
		log.Fatal(fmt.Errorf("run: %w", err))
	}
}
