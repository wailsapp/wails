//go:build !windows

package tui

import (
	"fmt"
	"sync"
	"time"

	"github.com/kungfusheep/glyph"
)

// runSpinner animates a glyph inline spinner while fn runs.
func runSpinner(message string, fn func() error) error {
	var (
		mu    sync.Mutex
		frame int
	)

	displayMsg := message
	app := glyph.NewInlineApp()
	app.Height(1)
	// glyph reads the *int frame pointer itself without taking our mutex.
	app.SetView(
		glyph.HBox(
			glyph.Spinner(&frame).Frames(glyph.SpinnerBraille).FG(glyph.Cyan),
			glyph.SpaceW(1),
			glyph.Text(&displayMsg),
		),
	)

	stopTicker := make(chan struct{})
	tickerDone := make(chan struct{})
	go func() {
		defer close(tickerDone)
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stopTicker:
				return
			case <-ticker.C:
				mu.Lock()
				frame++
				mu.Unlock()
				app.RequestRender()
			}
		}
	}()

	// workerDone receives fn's result; it is buffered so the worker never blocks.
	type result struct {
		err      error
		panicked bool
		panicVal any
	}
	workerDone := make(chan result, 1)
	go func() {
		var res result
		defer func() {
			if r := recover(); r != nil {
				res.panicked = true
				res.panicVal = r
			}
			// Always stop the app, even if fn panicked.
			app.Stop()
			workerDone <- res
		}()
		res.err = fn()
	}()

	_ = app.RunNonInteractive()

	// RunNonInteractive may return before fn has finished (e.g. on its own
	// error), so always wait for the worker before reading its result.
	res := <-workerDone
	close(stopTicker)
	<-tickerDone

	if res.panicked {
		// Terminal state is restored; re-panic on the calling goroutine.
		panic(res.panicVal)
	}

	if res.err != nil {
		fmt.Println(ansi("91", "  ✗") + " " + message)
	} else {
		fmt.Println(ansi("92", "  ✓") + " " + message)
	}
	return res.err
}
