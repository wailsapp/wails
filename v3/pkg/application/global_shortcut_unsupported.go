//go:build ios || android || server || (linux && !cgo)

package application

import "errors"

// errGlobalShortcutsUnsupported is returned on platforms where system-wide
// global shortcuts are not available (mobile, headless/server builds).
var errGlobalShortcutsUnsupported = errors.New("global shortcuts are not supported on this platform")

type unsupportedGlobalShortcuts struct{}

func newGlobalShortcutImpl(_ *GlobalShortcutManager) globalShortcutImpl {
	return &unsupportedGlobalShortcuts{}
}

func (unsupportedGlobalShortcuts) register(_ int, _ *accelerator) error {
	return errGlobalShortcutsUnsupported
}

func (unsupportedGlobalShortcuts) unregister(_ int) error { return nil }

func (unsupportedGlobalShortcuts) unregisterAll() error { return nil }
