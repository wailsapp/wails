//go:build bench

package application

import (
	"testing"
	"time"
)

// Note: SystemTray benchmarks are limited since actual system tray operations
// require platform-specific GUI initialization. These benchmarks focus on
// the Go-side logic that can be tested without a running GUI.

// BenchmarkSystemTrayCreation measures the cost of creating SystemTray instances
func BenchmarkSystemTrayCreation(b *testing.B) {
	for b.Loop() {
		tray := newSystemTray(1)
		_ = tray
	}
}

// BenchmarkSystemTrayConfiguration measures configuration operations
func BenchmarkSystemTrayConfiguration(b *testing.B) {
	b.Run("SetLabel", func(b *testing.B) {
		tray := newSystemTray(1)
		// impl is nil, so this just sets the field
		for b.Loop() {
			tray.SetLabel("Test Label")
		}
	})

	b.Run("SetTooltip", func(b *testing.B) {
		tray := newSystemTray(1)
		for b.Loop() {
			tray.SetTooltip("Test Tooltip")
		}
	})

	b.Run("SetIcon", func(b *testing.B) {
		tray := newSystemTray(1)
		icon := make([]byte, 1024) // 1KB icon data
		for b.Loop() {
			tray.SetIcon(icon)
		}
	})

	b.Run("SetDarkModeIcon", func(b *testing.B) {
		tray := newSystemTray(1)
		icon := make([]byte, 1024)
		for b.Loop() {
			tray.SetDarkModeIcon(icon)
		}
	})

	b.Run("SetTemplateIcon", func(b *testing.B) {
		tray := newSystemTray(1)
		icon := make([]byte, 1024)
		for b.Loop() {
			tray.SetTemplateIcon(icon)
		}
	})

	b.Run("SetIconPosition", func(b *testing.B) {
		tray := newSystemTray(1)
		for b.Loop() {
			tray.SetIconPosition(NSImageLeading)
		}
	})

	b.Run("ChainedConfiguration", func(b *testing.B) {
		icon := make([]byte, 1024)
		for b.Loop() {
			tray := newSystemTray(1)
			tray.SetIcon(icon).
				SetDarkModeIcon(icon).
				SetIconPosition(NSImageLeading)
		}
	})
}

// BenchmarkClickHandlerExecution measures handler registration and invocation
func BenchmarkClickHandlerExecution(b *testing.B) {
	b.Run("RegisterClickHandler", func(b *testing.B) {
		for b.Loop() {
			tray := newSystemTray(1)
			tray.OnClick(func() {})
		}
	})

	b.Run("RegisterAllHandlers", func(b *testing.B) {
		for b.Loop() {
			tray := newSystemTray(1)
			tray.OnClick(func() {})
			tray.OnRightClick(func() {})
			tray.OnDoubleClick(func() {})
			tray.OnRightDoubleClick(func() {})
			tray.OnMouseEnter(func() {})
			tray.OnMouseLeave(func() {})
		}
	})

	b.Run("InvokeClickHandler", func(b *testing.B) {
		tray := newSystemTray(1)
		counter := 0
		tray.OnClick(func() {
			counter++
		})

		b.ResetTimer()
		for b.Loop() {
			if tray.clickHandler != nil {
				tray.clickHandler()
			}
		}
	})

	b.Run("InvokeAllHandlers", func(b *testing.B) {
		tray := newSystemTray(1)
		counter := 0
		handler := func() { counter++ }
		tray.OnClick(handler)
		tray.OnRightClick(handler)
		tray.OnDoubleClick(handler)
		tray.OnRightDoubleClick(handler)
		tray.OnMouseEnter(handler)
		tray.OnMouseLeave(handler)

		b.ResetTimer()
		for b.Loop() {
			if tray.clickHandler != nil {
				tray.clickHandler()
			}
			if tray.rightClickHandler != nil {
				tray.rightClickHandler()
			}
			if tray.doubleClickHandler != nil {
				tray.doubleClickHandler()
			}
			if tray.rightDoubleClickHandler != nil {
				tray.rightDoubleClickHandler()
			}
			if tray.mouseEnterHandler != nil {
				tray.mouseEnterHandler()
			}
			if tray.mouseLeaveHandler != nil {
				tray.mouseLeaveHandler()
			}
		}
	})
}

// BenchmarkWindowAttachment measures window attachment configuration
func BenchmarkWindowAttachment(b *testing.B) {
	b.Run("AttachWindow", func(b *testing.B) {
		// We can't create real windows, but we can test the attachment logic
		for b.Loop() {
			tray := newSystemTray(1)
			// AttachWindow accepts nil gracefully
			tray.AttachWindow(nil)
		}
	})

	b.Run("WindowOffset", func(b *testing.B) {
		tray := newSystemTray(1)
		for b.Loop() {
			tray.WindowOffset(10)
		}
	})

	b.Run("WindowDebounce", func(b *testing.B) {
		tray := newSystemTray(1)
		for b.Loop() {
			tray.WindowDebounce(200 * time.Millisecond)
		}
	})

	b.Run("ChainedAttachment", func(b *testing.B) {
		for b.Loop() {
			tray := newSystemTray(1)
			tray.AttachWindow(nil).
				WindowOffset(10).
				WindowDebounce(200 * time.Millisecond)
		}
	})
}

// BenchmarkMenuConfiguration measures menu setup operations
func BenchmarkMenuConfiguration(b *testing.B) {
	b.Run("SetNilMenu", func(b *testing.B) {
		tray := newSystemTray(1)
		for b.Loop() {
			tray.SetMenu(nil)
		}
	})

	b.Run("SetSimpleMenu", func(b *testing.B) {
		menu := NewMenu()
		menu.Add("Item 1")
		menu.Add("Item 2")
		menu.Add("Item 3")

		tray := newSystemTray(1)
		b.ResetTimer()
		for b.Loop() {
			tray.SetMenu(menu)
		}
	})

	b.Run("SetComplexMenu", func(b *testing.B) {
		menu := NewMenu()
		for i := 0; i < 20; i++ {
			menu.Add("Item")
		}
		submenu := NewMenu()
		for i := 0; i < 10; i++ {
			submenu.Add("Subitem")
		}

		tray := newSystemTray(1)
		b.ResetTimer()
		for b.Loop() {
			tray.SetMenu(menu)
		}
	})
}

// BenchmarkIconSizes measures icon handling with different sizes
func BenchmarkIconSizes(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"16x16", 16 * 16 * 4},     // 1KB - small icon
		{"32x32", 32 * 32 * 4},     // 4KB - medium icon
		{"64x64", 64 * 64 * 4},     // 16KB - large icon
		{"128x128", 128 * 128 * 4}, // 64KB - retina icon
		{"256x256", 256 * 256 * 4}, // 256KB - high-res icon
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			icon := make([]byte, size.size)
			tray := newSystemTray(1)

			b.ResetTimer()
			for b.Loop() {
				tray.SetIcon(icon)
			}
		})
	}
}

// BenchmarkWindowAttachConfigInit measures WindowAttachConfig initialization
func BenchmarkWindowAttachConfigInit(b *testing.B) {
	b.Run("DefaultConfig", func(b *testing.B) {
		for b.Loop() {
			config := WindowAttachConfig{
				Window:   nil,
				Offset:   0,
				Debounce: 200 * time.Millisecond,
			}
			_ = config
		}
	})

	b.Run("FullConfig", func(b *testing.B) {
		for b.Loop() {
			config := WindowAttachConfig{
				Window:       nil,
				Offset:       10,
				Debounce:     300 * time.Millisecond,
				justClosed:   false,
				hasBeenShown: true,
			}
			_ = config
		}
	})
}

// BenchmarkSystemTrayShowHide measures show/hide state changes
// Note: These operations are no-ops when impl is nil, but we measure the check overhead
func BenchmarkSystemTrayShowHide(b *testing.B) {
	b.Run("Show", func(b *testing.B) {
		tray := newSystemTray(1)
		for b.Loop() {
			tray.Show()
		}
	})

	b.Run("Hide", func(b *testing.B) {
		tray := newSystemTray(1)
		for b.Loop() {
			tray.Hide()
		}
	})

	b.Run("ToggleShowHide", func(b *testing.B) {
		tray := newSystemTray(1)
		for b.Loop() {
			tray.Show()
			tray.Hide()
		}
	})
}

// BenchmarkIconPositionConstants measures icon position constant access
func BenchmarkIconPositionConstants(b *testing.B) {
	positions := []IconPosition{
		NSImageNone,
		NSImageOnly,
		NSImageLeft,
		NSImageRight,
		NSImageBelow,
		NSImageAbove,
		NSImageOverlaps,
		NSImageLeading,
		NSImageTrailing,
	}

	b.Run("SetAllPositions", func(b *testing.B) {
		tray := newSystemTray(1)
		for b.Loop() {
			for _, pos := range positions {
				tray.SetIconPosition(pos)
			}
		}
	})
}

// BenchmarkLabelOperations measures label getter/setter performance
func BenchmarkLabelOperations(b *testing.B) {
	b.Run("SetLabel", func(b *testing.B) {
		tray := newSystemTray(1)
		for b.Loop() {
			tray.SetLabel("System Tray Label")
		}
	})

	b.Run("GetLabel", func(b *testing.B) {
		tray := newSystemTray(1)
		tray.SetLabel("System Tray Label")
		b.ResetTimer()
		for b.Loop() {
			_ = tray.Label()
		}
	})

	b.Run("SetGetLabel", func(b *testing.B) {
		tray := newSystemTray(1)
		for b.Loop() {
			tray.SetLabel("Label")
			_ = tray.Label()
		}
	})
}

// BenchmarkDefaultClickHandler measures the default click handler logic
func BenchmarkDefaultClickHandler(b *testing.B) {
	b.Run("NoAttachedWindow", func(b *testing.B) {
		tray := newSystemTray(1)
		// With no menu and no attached window, defaultClickHandler returns early
		for b.Loop() {
			tray.defaultClickHandler()
		}
	})

	b.Run("WithNilWindow", func(b *testing.B) {
		tray := newSystemTray(1)
		tray.attachedWindow.Window = nil
		for b.Loop() {
			tray.defaultClickHandler()
		}
	})
}
