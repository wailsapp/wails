//go:build bench

package application_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// mockWindowDispatcher implements a no-op dispatcher for benchmarking
type mockWindowDispatcher struct {
	count atomic.Int64
}

func (m *mockWindowDispatcher) dispatchEventToWindows(event *application.CustomEvent) {
	m.count.Add(1)
}

// BenchmarkEventEmit measures event emission with varying listener counts
func BenchmarkEventEmit(b *testing.B) {
	listenerCounts := []int{0, 1, 10, 100}

	for _, count := range listenerCounts {
		b.Run(fmt.Sprintf("Listeners%d", count), func(b *testing.B) {
			dispatcher := &mockWindowDispatcher{}
			processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)

			// Register listeners
			for i := 0; i < count; i++ {
				processor.On("benchmark-event", func(event *application.CustomEvent) {
					// Minimal work
					_ = event.Data
				})
			}

			event := &application.CustomEvent{
				Name: "benchmark-event",
				Data: "test payload",
			}

			b.ResetTimer()
			for b.Loop() {
				_ = processor.Emit(event)
			}
		})
	}
}

// BenchmarkEventRegistration measures the cost of registering event listeners
func BenchmarkEventRegistration(b *testing.B) {
	dispatcher := &mockWindowDispatcher{}

	b.Run("SingleRegistration", func(b *testing.B) {
		for b.Loop() {
			processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)
			processor.On("test-event", func(event *application.CustomEvent) {})
		}
	})

	b.Run("MultipleRegistrations", func(b *testing.B) {
		for b.Loop() {
			processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)
			for i := 0; i < 10; i++ {
				processor.On(fmt.Sprintf("test-event-%d", i), func(event *application.CustomEvent) {})
			}
		}
	})

	b.Run("SameEventMultipleListeners", func(b *testing.B) {
		for b.Loop() {
			processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)
			for i := 0; i < 10; i++ {
				processor.On("test-event", func(event *application.CustomEvent) {})
			}
		}
	})
}

// BenchmarkEventUnregistration measures the cost of unregistering event listeners
func BenchmarkEventUnregistration(b *testing.B) {
	dispatcher := &mockWindowDispatcher{}

	b.Run("SingleUnregister", func(b *testing.B) {
		for b.Loop() {
			processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)
			cancel := processor.On("test-event", func(event *application.CustomEvent) {})
			cancel()
		}
	})

	b.Run("UnregisterFromMany", func(b *testing.B) {
		processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)
		// Pre-register many listeners
		cancels := make([]func(), 100)
		for i := 0; i < 100; i++ {
			cancels[i] = processor.On("test-event", func(event *application.CustomEvent) {})
		}

		b.ResetTimer()
		for i := 0; b.Loop(); i++ {
			// Re-register to have something to unregister
			if i%100 == 0 {
				for j := 0; j < 100; j++ {
					cancels[j] = processor.On("test-event", func(event *application.CustomEvent) {})
				}
			}
			cancels[i%100]()
		}
	})

	b.Run("OffAllListeners", func(b *testing.B) {
		for b.Loop() {
			processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)
			for i := 0; i < 10; i++ {
				processor.On("test-event", func(event *application.CustomEvent) {})
			}
			processor.Off("test-event")
		}
	})
}

// BenchmarkHookExecution measures the cost of hook execution during emit
func BenchmarkHookExecution(b *testing.B) {
	hookCounts := []int{0, 1, 5, 10}

	for _, count := range hookCounts {
		b.Run(fmt.Sprintf("Hooks%d", count), func(b *testing.B) {
			dispatcher := &mockWindowDispatcher{}
			processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)

			// Register hooks
			for i := 0; i < count; i++ {
				processor.RegisterHook("benchmark-event", func(event *application.CustomEvent) {
					// Minimal work - don't cancel
					_ = event.Data
				})
			}

			event := &application.CustomEvent{
				Name: "benchmark-event",
				Data: "test payload",
			}

			b.ResetTimer()
			for b.Loop() {
				_ = processor.Emit(event)
			}
		})
	}
}

// BenchmarkConcurrentEmit measures event emission under concurrent load
func BenchmarkConcurrentEmit(b *testing.B) {
	concurrencyLevels := []int{1, 4, 16}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Goroutines%d", concurrency), func(b *testing.B) {
			dispatcher := &mockWindowDispatcher{}
			processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)

			// Register a few listeners
			for i := 0; i < 5; i++ {
				processor.On("benchmark-event", func(event *application.CustomEvent) {
					_ = event.Data
				})
			}

			b.ResetTimer()
			b.SetParallelism(concurrency)
			b.RunParallel(func(pb *testing.PB) {
				event := &application.CustomEvent{
					Name: "benchmark-event",
					Data: "test payload",
				}
				for pb.Next() {
					_ = processor.Emit(event)
				}
			})
		})
	}
}

// BenchmarkEventToJSON measures CustomEvent JSON serialization
func BenchmarkEventToJSON(b *testing.B) {
	b.Run("SimpleData", func(b *testing.B) {
		event := &application.CustomEvent{
			Name: "test-event",
			Data: "simple string payload",
		}
		for b.Loop() {
			_ = event.ToJSON()
		}
	})

	b.Run("ComplexData", func(b *testing.B) {
		event := &application.CustomEvent{
			Name: "test-event",
			Data: map[string]interface{}{
				"id":      12345,
				"name":    "Test Event",
				"tags":    []string{"tag1", "tag2", "tag3"},
				"enabled": true,
				"nested": map[string]interface{}{
					"value": 3.14159,
				},
			},
		}
		for b.Loop() {
			_ = event.ToJSON()
		}
	})

	b.Run("WithSender", func(b *testing.B) {
		event := &application.CustomEvent{
			Name:   "test-event",
			Data:   "payload",
			Sender: "main-window",
		}
		for b.Loop() {
			_ = event.ToJSON()
		}
	})
}

// BenchmarkAtomicCancel measures the atomic cancel/check operations
func BenchmarkAtomicCancel(b *testing.B) {
	b.Run("Cancel", func(b *testing.B) {
		for b.Loop() {
			event := &application.CustomEvent{
				Name: "test",
				Data: nil,
			}
			event.Cancel()
		}
	})

	b.Run("IsCancelled", func(b *testing.B) {
		event := &application.CustomEvent{
			Name: "test",
			Data: nil,
		}
		for b.Loop() {
			_ = event.IsCancelled()
		}
	})

	b.Run("CancelAndCheck", func(b *testing.B) {
		for b.Loop() {
			event := &application.CustomEvent{
				Name: "test",
				Data: nil,
			}
			event.Cancel()
			_ = event.IsCancelled()
		}
	})
}

// BenchmarkEventProcessorCreation measures processor instantiation
func BenchmarkEventProcessorCreation(b *testing.B) {
	dispatcher := &mockWindowDispatcher{}
	for b.Loop() {
		_ = application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)
	}
}

// BenchmarkOnceEvent measures the Once registration and auto-unregistration
func BenchmarkOnceEvent(b *testing.B) {
	dispatcher := &mockWindowDispatcher{}

	b.Run("RegisterAndTrigger", func(b *testing.B) {
		for b.Loop() {
			processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)
			var wg sync.WaitGroup
			wg.Add(1)
			processor.Once("once-event", func(event *application.CustomEvent) {
				wg.Done()
			})
			_ = processor.Emit(&application.CustomEvent{Name: "once-event", Data: nil})
			wg.Wait()
		}
	})
}

// BenchmarkOnMultipleEvent measures the OnMultiple registration
func BenchmarkOnMultipleEvent(b *testing.B) {
	dispatcher := &mockWindowDispatcher{}

	b.Run("ThreeEvents", func(b *testing.B) {
		for b.Loop() {
			processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)
			var wg sync.WaitGroup
			wg.Add(3)
			processor.OnMultiple("multi-event", func(event *application.CustomEvent) {
				wg.Done()
			}, 3)
			for i := 0; i < 3; i++ {
				_ = processor.Emit(&application.CustomEvent{Name: "multi-event", Data: nil})
			}
			wg.Wait()
		}
	})
}

// BenchmarkMixedEventOperations simulates realistic event usage patterns
func BenchmarkMixedEventOperations(b *testing.B) {
	dispatcher := &mockWindowDispatcher{}

	b.Run("RegisterEmitUnregister", func(b *testing.B) {
		processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)

		for b.Loop() {
			cancel := processor.On("mixed-event", func(event *application.CustomEvent) {
				_ = event.Data
			})
			_ = processor.Emit(&application.CustomEvent{Name: "mixed-event", Data: "test"})
			cancel()
		}
	})

	b.Run("HookAndEmit", func(b *testing.B) {
		processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)
		processor.RegisterHook("hooked-event", func(event *application.CustomEvent) {
			// Validation hook
			if event.Data == nil {
				event.Cancel()
			}
		})
		processor.On("hooked-event", func(event *application.CustomEvent) {
			_ = event.Data
		})

		event := &application.CustomEvent{Name: "hooked-event", Data: "valid"}

		b.ResetTimer()
		for b.Loop() {
			_ = processor.Emit(event)
		}
	})
}

// BenchmarkEventNameLookup measures the map lookup performance for event names
func BenchmarkEventNameLookup(b *testing.B) {
	dispatcher := &mockWindowDispatcher{}
	processor := application.NewWailsEventProcessor(dispatcher.dispatchEventToWindows)

	// Register events with different name lengths
	shortName := "evt"
	mediumName := "application:user:action"
	longName := "com.mycompany.myapp.module.submodule.event.type.action"

	processor.On(shortName, func(event *application.CustomEvent) {})
	processor.On(mediumName, func(event *application.CustomEvent) {})
	processor.On(longName, func(event *application.CustomEvent) {})

	b.Run("ShortName", func(b *testing.B) {
		event := &application.CustomEvent{Name: shortName, Data: nil}
		for b.Loop() {
			_ = processor.Emit(event)
		}
	})

	b.Run("MediumName", func(b *testing.B) {
		event := &application.CustomEvent{Name: mediumName, Data: nil}
		for b.Loop() {
			_ = processor.Emit(event)
		}
	})

	b.Run("LongName", func(b *testing.B) {
		event := &application.CustomEvent{Name: longName, Data: nil}
		for b.Loop() {
			_ = processor.Emit(event)
		}
	})
}
