//go:build windows && !server

package application

import (
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/w32"
)

// windowsGlobalShortcuts implements globalShortcutImpl using the Win32
// RegisterHotKey API. Hot keys are registered against the application's hidden
// main-thread window so that WM_HOTKEY messages are delivered to the same
// message loop the rest of the application already pumps (see wndProc's
// WM_HOTKEY case, which calls back into the manager's dispatch).
//
// RegisterHotKey is thread-affine: the WM_HOTKEY message is posted to the
// thread that owns the window passed in. Registration therefore must happen on
// the main UI thread - the GlobalShortcutManager guarantees this by wrapping
// every call in InvokeSync.
type windowsGlobalShortcuts struct {
	manager *GlobalShortcutManager
	ids     map[int]struct{}
}

func newGlobalShortcutImpl(manager *GlobalShortcutManager) globalShortcutImpl {
	return &windowsGlobalShortcuts{
		manager: manager,
		ids:     make(map[int]struct{}),
	}
}

func (g *windowsGlobalShortcuts) hwnd() (w32.HWND, error) {
	app, ok := globalApplication.impl.(*windowsApp)
	if !ok || app == nil || app.mainThreadWindowHWND == 0 {
		return 0, fmt.Errorf("global shortcuts require the application to be running")
	}
	return app.mainThreadWindowHWND, nil
}

func (g *windowsGlobalShortcuts) register(id int, accel *accelerator) error {
	vk, ok := winKeyCodes[accel.Key]
	if !ok {
		return fmt.Errorf("key %q is not supported as a global shortcut on Windows", accel.Key)
	}

	// MOD_NOREPEAT prevents auto-repeat from spamming the callback while the
	// keys are held down.
	mods := uint(w32.MOD_NOREPEAT)
	for _, m := range accel.Modifiers {
		switch m {
		case CmdOrCtrlKey, ControlKey:
			mods |= w32.MOD_CONTROL
		case OptionOrAltKey:
			mods |= w32.MOD_ALT
		case ShiftKey:
			mods |= w32.MOD_SHIFT
		case SuperKey:
			mods |= w32.MOD_WIN
		}
	}

	hwnd, err := g.hwnd()
	if err != nil {
		return err
	}

	if !w32.RegisterHotKey(hwnd, id, mods, vk) {
		// RegisterHotKey returns false (ERROR_HOTKEY_ALREADY_REGISTERED) when
		// the combination is already owned - either by this process or, more
		// commonly, by another application.
		return fmt.Errorf("the shortcut is already registered (possibly by another application)")
	}
	g.ids[id] = struct{}{}
	return nil
}

func (g *windowsGlobalShortcuts) unregister(id int) error {
	if _, ok := g.ids[id]; !ok {
		return nil
	}
	delete(g.ids, id)
	hwnd, err := g.hwnd()
	if err != nil {
		return err
	}
	if !w32.UnregisterHotKey(hwnd, id) {
		return fmt.Errorf("UnregisterHotKey failed for shortcut id %d", id)
	}
	return nil
}

func (g *windowsGlobalShortcuts) unregisterAll() error {
	var firstErr error
	for id := range g.ids {
		if err := g.unregister(id); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// winKeyCodes maps Wails accelerator key names (already lower-cased by
// parseAccelerator) to Windows virtual-key codes. Letters and digits map to
// their ASCII-uppercase value (VK_A == 'A' == 0x41, VK_0 == '0' == 0x30).
var winKeyCodes = map[string]uint{
	// Letters
	"a": 0x41, "b": 0x42, "c": 0x43, "d": 0x44, "e": 0x45, "f": 0x46,
	"g": 0x47, "h": 0x48, "i": 0x49, "j": 0x4A, "k": 0x4B, "l": 0x4C,
	"m": 0x4D, "n": 0x4E, "o": 0x4F, "p": 0x50, "q": 0x51, "r": 0x52,
	"s": 0x53, "t": 0x54, "u": 0x55, "v": 0x56, "w": 0x57, "x": 0x58,
	"y": 0x59, "z": 0x5A,
	// Number row
	"0": 0x30, "1": 0x31, "2": 0x32, "3": 0x33, "4": 0x34,
	"5": 0x35, "6": 0x36, "7": 0x37, "8": 0x38, "9": 0x39,
	// Punctuation (OEM keys, US layout)
	";": 0xBA, "=": 0xBB, ",": 0xBC, "-": 0xBD, ".": 0xBE, "/": 0xBF,
	"`": 0xC0, "[": 0xDB, "\\": 0xDC, "]": 0xDD, "'": 0xDE, "+": 0xBB,
	// Named keys
	"backspace": 0x08,
	"tab":       0x09,
	"return":    0x0D,
	"enter":     0x0D,
	"escape":    0x1B,
	"space":     0x20,
	"page up":   0x21,
	"page down": 0x22,
	"end":       0x23,
	"home":      0x24,
	"left":      0x25,
	"up":        0x26,
	"right":     0x27,
	"down":      0x28,
	"delete":    0x2E,
	"numlock":   0x90,
	// Function keys
	"f1": 0x70, "f2": 0x71, "f3": 0x72, "f4": 0x73, "f5": 0x74, "f6": 0x75,
	"f7": 0x76, "f8": 0x77, "f9": 0x78, "f10": 0x79, "f11": 0x7A, "f12": 0x7B,
	"f13": 0x7C, "f14": 0x7D, "f15": 0x7E, "f16": 0x7F, "f17": 0x80, "f18": 0x81,
	"f19": 0x82, "f20": 0x83, "f21": 0x84, "f22": 0x85, "f23": 0x86, "f24": 0x87,
}
