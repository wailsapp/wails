//go:build windows

package application

import (
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/w32"
)

// showSnapAssist sends Win+Z key combination to trigger SnapAssist for the specified window
func showSnapAssist(window *WebviewWindow) {
	// First, ensure the window is visible and focused to target SnapAssist correctly
	window.Show()    // Ensure window is visible
	window.Restore() // Restore if minimized
	
	// Get the native window handle and set it as foreground
	if hwnd := w32.HWND(window.impl.nativeWindowHandle()); hwnd != 0 {
		w32.SetForegroundWindow(hwnd)
		// Small delay to ensure window focus is established
		time.Sleep(50 * time.Millisecond)
	}
	
	// Send Win+Z key combination using a robotgo-style approach
	sendKeyCombo(w32.VK_Z, w32.VK_LWIN)
}

// sendKeyCombo sends a key combination (key + modifiers)
func sendKeyCombo(key uint16, modifiers ...uint16) {
	// Calculate total inputs needed (press and release for each key)
	numKeys := len(modifiers) + 1
	inputs := make([]w32.INPUT, numKeys*2)
	
	inputIndex := 0
	
	// Press all modifier keys first
	for _, modifier := range modifiers {
		inputs[inputIndex] = w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki: w32.KEYBDINPUT{
				WVk:         modifier,
				WScan:       0,
				DwFlags:     0,
				Time:        0,
				DwExtraInfo: 0,
			},
		}
		inputIndex++
	}
	
	// Press the main key
	inputs[inputIndex] = w32.INPUT{
		Type: w32.INPUT_KEYBOARD,
		Ki: w32.KEYBDINPUT{
			WVk:         key,
			WScan:       0,
			DwFlags:     0,
			Time:        0,
			DwExtraInfo: 0,
		},
	}
	inputIndex++
	
	// Release the main key
	inputs[inputIndex] = w32.INPUT{
		Type: w32.INPUT_KEYBOARD,
		Ki: w32.KEYBDINPUT{
			WVk:         key,
			WScan:       0,
			DwFlags:     w32.KEYEVENTF_KEYUP,
			Time:        0,
			DwExtraInfo: 0,
		},
	}
	inputIndex++
	
	// Release all modifier keys in reverse order
	for i := len(modifiers) - 1; i >= 0; i-- {
		inputs[inputIndex] = w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki: w32.KEYBDINPUT{
				WVk:         modifiers[i],
				WScan:       0,
				DwFlags:     w32.KEYEVENTF_KEYUP,
				Time:        0,
				DwExtraInfo: 0,
			},
		}
		inputIndex++
	}
	
	// Send all input events
	w32.SendInput(len(inputs), unsafe.Pointer(&inputs[0]), int(unsafe.Sizeof(w32.INPUT{})))
}