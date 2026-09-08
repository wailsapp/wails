//go:build windows

package commands

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/rjeczalik/notify"
)

// notify's Windows backend casts short event buffers to MAX_LONG_PATH arrays,
// which fails Go's pointer checks. Use fsnotify for native Windows events while
// retaining the dev session's existing event contract.
type manifestWindowsWatch struct {
	watcher *fsnotify.Watcher
	stop    chan struct{}
	done    chan struct{}
}

var manifestWindowsWatches = struct {
	sync.Mutex
	channels map[chan<- notify.EventInfo]*manifestWindowsWatch
}{channels: make(map[chan<- notify.EventInfo]*manifestWindowsWatch)}

type manifestWindowsEvent struct {
	path  string
	event notify.Event
}

func (e manifestWindowsEvent) Path() string        { return e.path }
func (e manifestWindowsEvent) Event() notify.Event { return e.event }
func (e manifestWindowsEvent) Sys() any            { return nil }

func watchManifestDirectory(path string, events chan<- notify.EventInfo) error {
	manifestWindowsWatches.Lock()
	defer manifestWindowsWatches.Unlock()
	watch := manifestWindowsWatches.channels[events]
	if watch == nil {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		watch = &manifestWindowsWatch{watcher: watcher, stop: make(chan struct{}), done: make(chan struct{})}
		manifestWindowsWatches.channels[events] = watch
		go watch.forward(events)
	}
	return watch.watcher.Add(path)
}

func (w *manifestWindowsWatch) forward(events chan<- notify.EventInfo) {
	defer close(w.done)
	for {
		select {
		case <-w.stop:
			return
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			var kind notify.Event
			if event.Op&fsnotify.Create != 0 {
				kind |= notify.Create
			}
			if event.Op&(fsnotify.Write|fsnotify.Chmod) != 0 {
				kind |= notify.Write
			}
			if event.Op&fsnotify.Remove != 0 {
				kind |= notify.Remove
			}
			if event.Op&fsnotify.Rename != 0 {
				kind |= notify.Rename
			}
			if kind == 0 {
				continue
			}
			select {
			case events <- manifestWindowsEvent{event.Name, kind}:
			case <-w.stop:
				return
			}
		case _, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			// An overflow means the precise change is unknown. Request a reload through
			// the manifest path so the session rebuilds watches and reevaluates inputs.
			for _, path := range w.watcher.WatchList() {
				select {
				case events <- manifestWindowsEvent{path + "/wails.hcl", notify.Write}:
				case <-w.stop:
					return
				}
			}
		}
	}
}
func stopManifestWatchChannel(events chan<- notify.EventInfo) {
	manifestWindowsWatches.Lock()
	watch := manifestWindowsWatches.channels[events]
	delete(manifestWindowsWatches.channels, events)
	manifestWindowsWatches.Unlock()
	if watch != nil {
		close(watch.stop)
		_ = watch.watcher.Close()
		<-watch.done
	}
}
