package debounce

import (
	"sync"
	"time"
)

// New returns a debounced function that calls f after it stops being invoked
// for the given duration. The last invocation wins if called with different functions.
func New(after time.Duration) func(f func()) {
	d := &debouncer{after: after}
	return func(f func()) {
		d.add(f)
	}
}

type debouncer struct {
	mu         sync.Mutex
	after      time.Duration
	timer      *time.Timer
	generation uint64
}

func (d *debouncer) add(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
	}
	d.generation++
	gen := d.generation
	d.timer = time.AfterFunc(d.after, func() {
		d.mu.Lock()
		if d.generation != gen {
			d.mu.Unlock()
			return
		}
		d.mu.Unlock()
		f()
	})
}
