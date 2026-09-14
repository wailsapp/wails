// User Panic Handling Example
//
// Wails automatically recovers panics that occur inside bound Service methods
// and internal runtime callbacks — it captures a PanicDetails and routes it to
// the PanicHandler you register in application.Options.
//
// Wails does NOT automatically recover panics in goroutines you spawn
// yourself. If one of those panics, the default Go behaviour applies: the
// whole process crashes.
//
// This example shows how to funnel BOTH panic paths through a single handler:
//   1. Bound-method panic  — wails catches it, calls reportPanic directly.
//   2. User goroutine panic — user defers recoverAndReport(), which builds a
//      PanicDetails manually and calls the same reportPanic function.

package main

import (
	"embed"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed assets/*
var assets embed.FS

var app *application.App

// reportPanic is the single place every panic in this app ends up, regardless
// of whether wails caught it or user code did. Replace the body with real
// reporting (log file, Sentry, crash dialog, etc.) in a production app.
func reportPanic(pd *application.PanicDetails) {
	fmt.Printf("\n*** PANIC ***\n")
	fmt.Printf("Source: %s\n", panicSource(pd))
	fmt.Printf("Time:   %s\n", pd.Time.Format(time.RFC3339))
	fmt.Printf("Error:  %s\n", pd.Error)
	fmt.Printf("Stack:\n%s\n", pd.StackTrace)
}

// panicSource labels whether this PanicDetails came from wails' own recovery
// (trimmed stack trace) or from user-side recoverAndReport (full stack).
// Purely for demo output — wails does not mark the source itself.
func panicSource(pd *application.PanicDetails) string {
	if pd.StackTrace == pd.FullStackTrace {
		return "user goroutine (recovered via recoverAndReport)"
	}
	return "wails runtime (bound method or internal callback)"
}

// recoverAndReport is what user code defers at the top of any goroutine it
// spawns, so that panic recovery flows through the same handler as wails.
//
// Usage:
//   go func() {
//       defer recoverAndReport()
//       // ...work...
//   }()
func recoverAndReport() {
	r := recover()
	if r == nil {
		return
	}
	err, ok := r.(error)
	if !ok {
		err = fmt.Errorf("%v", r)
	}
	stack := string(debug.Stack())
	reportPanic(&application.PanicDetails{
		Error:          err,
		Time:           time.Now(),
		StackTrace:     stack,
		FullStackTrace: stack,
	})
}

// ---- Demo services --------------------------------------------------------

// WindowService is a normal wails-bound service. Panics in GeneratePanic are
// caught by wails automatically because the method is invoked from the
// frontend through wails' binding machinery.
type WindowService struct{}

func (s *WindowService) GeneratePanic() {
	panic(fmt.Errorf("panic from bound service method — wails catches this automatically"))
}

// BackgroundWorker is a plain Go type, NOT a wails service. It owns its own
// goroutine that wails has no hooks into. Without the deferred
// recoverAndReport(), a panic inside tickLoop would take the process down.
type BackgroundWorker struct{}

func (w *BackgroundWorker) Start() {
	go w.tickLoop()
}

func (w *BackgroundWorker) tickLoop() {
	defer recoverAndReport() // MUST be first deferred call in every user goroutine
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	count := 0
	for range ticker.C {
		count++
		if count == 3 {
			panic(fmt.Errorf("ticker exploded after %d ticks", count))
		}
	}
}

// ---- Wire everything up ---------------------------------------------------

func main() {
	app = application.New(application.Options{
		Name:        "User Panic Handler Demo",
		Description: "Handle panics in both wails-bound code and user-spawned goroutines",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Services: []application.Service{
			application.NewService(&WindowService{}),
		},
		// Register the handler wails uses for bound-method + internal panics.
		PanicHandler: reportPanic,
	})

	app.Window.New().SetTitle("User Panic Handler Demo").Show()

	// Start a goroutine wails knows nothing about. It will panic after 3s.
	(&BackgroundWorker{}).Start()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
