//go:build bench

package application

import (
	"fmt"
	"sync"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/events"
)

// Note: This file uses internal package access to benchmark window internals
// without requiring GUI initialization.

// BenchmarkWindowEventRegistration measures the cost of registering window event listeners
func BenchmarkWindowEventRegistration(b *testing.B) {
	b.Run("SingleListener", func(b *testing.B) {
		for b.Loop() {
			w := &WebviewWindow{
				eventListeners: make(map[uint][]*WindowEventListener),
			}
			w.OnWindowEvent(events.Common.WindowFocus, func(event *WindowEvent) {})
		}
	})

	b.Run("MultipleListenersSameEvent", func(b *testing.B) {
		for b.Loop() {
			w := &WebviewWindow{
				eventListeners: make(map[uint][]*WindowEventListener),
			}
			for i := 0; i < 10; i++ {
				w.OnWindowEvent(events.Common.WindowFocus, func(event *WindowEvent) {})
			}
		}
	})

	b.Run("MultipleListenersDifferentEvents", func(b *testing.B) {
		eventTypes := []events.WindowEventType{
			events.Common.WindowFocus,
			events.Common.WindowLostFocus,
			events.Common.WindowShow,
			events.Common.WindowHide,
			events.Common.WindowDidMove,
		}
		for b.Loop() {
			w := &WebviewWindow{
				eventListeners: make(map[uint][]*WindowEventListener),
			}
			for _, evt := range eventTypes {
				w.OnWindowEvent(evt, func(event *WindowEvent) {})
			}
		}
	})
}

// BenchmarkWindowHookRegistration measures the cost of registering window event hooks
func BenchmarkWindowHookRegistration(b *testing.B) {
	b.Run("SingleHook", func(b *testing.B) {
		eventID := uint(events.Common.WindowClosing)
		for b.Loop() {
			w := &WebviewWindow{
				eventHooks: make(map[uint][]*WindowEventListener),
			}
			w.eventHooksLock.Lock()
			w.eventHooks[eventID] = append(w.eventHooks[eventID], &WindowEventListener{
				callback: func(event *WindowEvent) {},
			})
			w.eventHooksLock.Unlock()
		}
	})

	b.Run("MultipleHooks", func(b *testing.B) {
		eventID := uint(events.Common.WindowClosing)
		for b.Loop() {
			w := &WebviewWindow{
				eventHooks: make(map[uint][]*WindowEventListener),
			}
			for i := 0; i < 5; i++ {
				w.eventHooksLock.Lock()
				w.eventHooks[eventID] = append(w.eventHooks[eventID], &WindowEventListener{
					callback: func(event *WindowEvent) {},
				})
				w.eventHooksLock.Unlock()
			}
		}
	})
}

// BenchmarkWindowEventDispatch measures the internal event dispatch mechanism
func BenchmarkWindowEventDispatch(b *testing.B) {
	listenerCounts := []int{0, 1, 5, 10, 50}

	for _, count := range listenerCounts {
		b.Run(fmt.Sprintf("Listeners%d", count), func(b *testing.B) {
			w := &WebviewWindow{
				eventListeners: make(map[uint][]*WindowEventListener),
			}

			eventID := uint(events.Common.WindowFocus)

			// Register listeners
			for i := 0; i < count; i++ {
				w.eventListeners[eventID] = append(w.eventListeners[eventID], &WindowEventListener{
					callback: func(event *WindowEvent) {
						_ = event.IsCancelled()
					},
				})
			}

			b.ResetTimer()
			for b.Loop() {
				w.eventListenersLock.RLock()
				listeners := w.eventListeners[eventID]
				w.eventListenersLock.RUnlock()
				_ = listeners
			}
		})
	}
}

// BenchmarkKeyBindingLookup measures key binding lookup performance
func BenchmarkKeyBindingLookup(b *testing.B) {
	bindingCounts := []int{1, 10, 50, 100}

	for _, count := range bindingCounts {
		b.Run(fmt.Sprintf("Bindings%d", count), func(b *testing.B) {
			w := &WebviewWindow{
				keyBindings: make(map[string]func(Window)),
			}

			// Register bindings
			for i := 0; i < count; i++ {
				key := fmt.Sprintf("ctrl+shift+%c", 'a'+i%26)
				w.keyBindings[key] = func(Window) {}
			}

			// Lookup key that exists
			lookupKey := "ctrl+shift+m"
			w.keyBindings[lookupKey] = func(Window) {}

			b.ResetTimer()
			for b.Loop() {
				w.keyBindingsLock.RLock()
				_ = w.keyBindings[lookupKey]
				w.keyBindingsLock.RUnlock()
			}
		})
	}

	b.Run("MissLookup", func(b *testing.B) {
		w := &WebviewWindow{
			keyBindings: make(map[string]func(Window)),
		}

		// Register some bindings
		for i := 0; i < 50; i++ {
			key := fmt.Sprintf("ctrl+shift+%c", 'a'+i%26)
			w.keyBindings[key] = func(Window) {}
		}

		lookupKey := "ctrl+alt+nonexistent"

		b.ResetTimer()
		for b.Loop() {
			w.keyBindingsLock.RLock()
			_ = w.keyBindings[lookupKey]
			w.keyBindingsLock.RUnlock()
		}
	})
}

// BenchmarkConcurrentWindowOps measures concurrent access patterns
func BenchmarkConcurrentWindowOps(b *testing.B) {
	b.Run("ConcurrentEventLookup", func(b *testing.B) {
		w := &WebviewWindow{
			eventListeners: make(map[uint][]*WindowEventListener),
		}

		eventID := uint(events.Common.WindowFocus)
		for i := 0; i < 10; i++ {
			w.eventListeners[eventID] = append(w.eventListeners[eventID], &WindowEventListener{
				callback: func(event *WindowEvent) {},
			})
		}

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				w.eventListenersLock.RLock()
				_ = w.eventListeners[eventID]
				w.eventListenersLock.RUnlock()
			}
		})
	})

	b.Run("ConcurrentKeyBindingLookup", func(b *testing.B) {
		w := &WebviewWindow{
			keyBindings: make(map[string]func(Window)),
		}

		for i := 0; i < 50; i++ {
			key := fmt.Sprintf("ctrl+shift+%c", 'a'+i%26)
			w.keyBindings[key] = func(Window) {}
		}

		keys := []string{"ctrl+shift+a", "ctrl+shift+m", "ctrl+shift+z"}

		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				w.keyBindingsLock.RLock()
				_ = w.keyBindings[keys[i%len(keys)]]
				w.keyBindingsLock.RUnlock()
				i++
			}
		})
	})

	b.Run("MixedReadWrite", func(b *testing.B) {
		w := &WebviewWindow{
			eventListeners: make(map[uint][]*WindowEventListener),
		}

		eventID := uint(events.Common.WindowFocus)

		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				if i%10 == 0 {
					// Write operation (10% of ops)
					w.eventListenersLock.Lock()
					w.eventListeners[eventID] = append(w.eventListeners[eventID], &WindowEventListener{
						callback: func(event *WindowEvent) {},
					})
					w.eventListenersLock.Unlock()
				} else {
					// Read operation (90% of ops)
					w.eventListenersLock.RLock()
					_ = w.eventListeners[eventID]
					w.eventListenersLock.RUnlock()
				}
				i++
			}
		})
	})
}

// BenchmarkWindowEventCreation measures WindowEvent allocation
func BenchmarkWindowEventCreation(b *testing.B) {
	for b.Loop() {
		event := NewWindowEvent()
		_ = event
	}
}

// BenchmarkWindowEventCancellation measures cancel/check operations
func BenchmarkWindowEventCancellation(b *testing.B) {
	b.Run("Cancel", func(b *testing.B) {
		for b.Loop() {
			event := NewWindowEvent()
			event.Cancel()
		}
	})

	b.Run("IsCancelled", func(b *testing.B) {
		event := NewWindowEvent()
		for b.Loop() {
			_ = event.IsCancelled()
		}
	})

	b.Run("CancelledCheck", func(b *testing.B) {
		event := NewWindowEvent()
		event.Cancel()
		for b.Loop() {
			_ = event.IsCancelled()
		}
	})
}

// BenchmarkWindowOptionsInit measures window options initialization patterns
func BenchmarkWindowOptionsInit(b *testing.B) {
	b.Run("DefaultOptions", func(b *testing.B) {
		for b.Loop() {
			opts := WebviewWindowOptions{}
			_ = opts
		}
	})

	b.Run("CommonOptions", func(b *testing.B) {
		for b.Loop() {
			opts := WebviewWindowOptions{
				Title:     "Test Window",
				Width:     800,
				Height:    600,
				MinWidth:  400,
				MinHeight: 300,
			}
			_ = opts
		}
	})

	b.Run("FullOptions", func(b *testing.B) {
		for b.Loop() {
			opts := WebviewWindowOptions{
				Title:             "Full Test Window",
				Width:             1024,
				Height:            768,
				MinWidth:          400,
				MinHeight:         300,
				MaxWidth:          1920,
				MaxHeight:         1080,
				URL:               "http://localhost:8080",
				Frameless:         false,
				DisableResize:     false,
				AlwaysOnTop:       false,
				Hidden:            false,
				EnableDragAndDrop: true,
				BackgroundColour:  RGBA{Red: 255, Green: 255, Blue: 255, Alpha: 255},
			}
			_ = opts
		}
	})
}

// BenchmarkMenuBindingLookup measures menu binding lookups
func BenchmarkMenuBindingLookup(b *testing.B) {
	w := &WebviewWindow{
		menuBindings: make(map[string]*MenuItem),
	}

	// Register some menu bindings
	for i := 0; i < 50; i++ {
		id := fmt.Sprintf("menu-item-%d", i)
		w.menuBindings[id] = &MenuItem{id: uint(i)}
	}

	lookupID := "menu-item-25"

	b.Run("Hit", func(b *testing.B) {
		for b.Loop() {
			w.menuBindingsLock.RLock()
			_ = w.menuBindings[lookupID]
			w.menuBindingsLock.RUnlock()
		}
	})

	b.Run("Miss", func(b *testing.B) {
		missID := "nonexistent-menu-item"
		for b.Loop() {
			w.menuBindingsLock.RLock()
			_ = w.menuBindings[missID]
			w.menuBindingsLock.RUnlock()
		}
	})
}

// BenchmarkWindowDestroyedCheck measures the destroyed flag check pattern
func BenchmarkWindowDestroyedCheck(b *testing.B) {
	w := &WebviewWindow{}

	b.Run("NotDestroyed", func(b *testing.B) {
		for b.Loop() {
			w.destroyedLock.RLock()
			_ = w.destroyed
			w.destroyedLock.RUnlock()
		}
	})

	b.Run("Destroyed", func(b *testing.B) {
		w.destroyed = true
		for b.Loop() {
			w.destroyedLock.RLock()
			_ = w.destroyed
			w.destroyedLock.RUnlock()
		}
	})
}

// BenchmarkCancellerManagement measures canceller function management
func BenchmarkCancellerManagement(b *testing.B) {
	b.Run("AddCanceller", func(b *testing.B) {
		for b.Loop() {
			w := &WebviewWindow{
				cancellers: make([]func(), 0),
			}
			for i := 0; i < 10; i++ {
				w.cancellersLock.Lock()
				w.cancellers = append(w.cancellers, func() {})
				w.cancellersLock.Unlock()
			}
		}
	})

	b.Run("ExecuteCancellers", func(b *testing.B) {
		w := &WebviewWindow{
			cancellers: make([]func(), 100),
		}
		for i := 0; i < 100; i++ {
			w.cancellers[i] = func() {}
		}

		b.ResetTimer()
		for b.Loop() {
			w.cancellersLock.RLock()
			cancellers := w.cancellers
			w.cancellersLock.RUnlock()
			for _, cancel := range cancellers {
				cancel()
			}
		}
	})
}

// BenchmarkRWMutexPatterns compares different locking patterns
func BenchmarkRWMutexPatterns(b *testing.B) {
	b.Run("RLockUnlock", func(b *testing.B) {
		var mu sync.RWMutex
		for b.Loop() {
			mu.RLock()
			mu.RUnlock()
		}
	})

	b.Run("LockUnlock", func(b *testing.B) {
		var mu sync.RWMutex
		for b.Loop() {
			mu.Lock()
			mu.Unlock()
		}
	})

	b.Run("DeferredRLock", func(b *testing.B) {
		var mu sync.RWMutex
		for b.Loop() {
			func() {
				mu.RLock()
				defer mu.RUnlock()
			}()
		}
	})
}
