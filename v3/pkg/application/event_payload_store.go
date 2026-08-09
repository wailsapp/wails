package application

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// Oversized Go→JS events are not spliced into the JavaScript source passed to
// the webview's eval. Above a platform-specific size WebKit switches from inline
// IPC message data to an out-of-line shared-memory transfer, and one of the two
// processes then retains those regions for as long as the app keeps emitting.
// Measured with v3/tests/event-performance:
//
//   - macOS 26.4.1 / WebKit-Cocoa: the retention lands in the HOST process,
//     visible in vmmap as "owned unmapped" — one region per oversized eval.
//     At a constant 4 MB/s, 8 KB payloads held the host flat at ~36 MB while
//     16 KB payloads climbed to 363 MB. 100 ev/s of 1 MB reached 11.5 GB.
//   - Ubuntu 26.04 / WebKitGTK 2.52.3: the mirror image — the host stays flat
//     and the WEB process grows instead, reaching 6.2 GB on the same scenario.
//     Its switchover sits higher, between 64 KB and 128 KB.
//
// Instead, the payload is parked here under an unguessable id and the webview
// is told to fetch it from the asset server, which delivers via a URL scheme
// task and never touches the eval IPC path.
const (
	// maxInlineEventPayload is the largest marshalled event JSON that may be
	// spliced directly into an eval. 8192 is the largest size measured at 0%
	// retention on macOS, whose knee (8-16 KB) is the lower of the two measured
	// platforms; WebKitGTK's is 64-128 KB, so this is correct there too, just
	// conservative. WebView2 is unmeasured — run the iso-* scenarios in
	// v3/tests/event-performance before assuming this value transfers.
	maxInlineEventPayload = 8192

	// eventPayloadTTL bounds how long an unfetched payload is held. A page
	// reload or a window close between dispatch and fetch would otherwise
	// strand it forever — which would just move the leak into Go.
	eventPayloadTTL = 30 * time.Second

	// eventPayloadStoreMaxBytes caps total parked bytes. On overflow the caller
	// falls back to inline delivery: that event then pays the out-of-line
	// retention cost, which is strictly better than dropping it or letting the
	// store grow without bound.
	eventPayloadStoreMaxBytes = 64 * 1024 * 1024

	eventPayloadPath = "/wails/eventpayload/"
)

type parkedEventPayload struct {
	data     []byte
	windowID uint
	created  time.Time
}

// eventPayloadStore holds oversized event payloads awaiting a one-shot fetch
// from the webview.
type eventPayloadStore struct {
	mu      sync.Mutex
	items   map[string]parkedEventPayload
	bytes   int
	janitor sync.Once
	stop    chan struct{}
}

func newEventPayloadStore() *eventPayloadStore {
	return &eventPayloadStore{
		items: make(map[string]parkedEventPayload),
		stop:  make(chan struct{}),
	}
}

// put parks a payload and returns its id. ok is false when the store is full,
// in which case the caller must fall back to inline delivery.
func (s *eventPayloadStore) put(windowID uint, data []byte) (id string, ok bool) {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", false
	}
	id = hex.EncodeToString(buf[:])

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.bytes+len(data) > eventPayloadStoreMaxBytes {
		return "", false
	}
	s.items[id] = parkedEventPayload{data: data, windowID: windowID, created: time.Now()}
	s.bytes += len(data)

	s.janitor.Do(func() { go s.reap() })
	return id, true
}

// take returns the payload for id exactly once. windowID must match the window
// the payload was dispatched to; a zero windowID skips the check, for platforms
// that do not tag asset requests with a window id.
func (s *eventPayloadStore) take(id string, windowID uint) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, found := s.items[id]
	if !found {
		return nil, false
	}
	if windowID != 0 && item.windowID != windowID {
		return nil, false
	}
	delete(s.items, id)
	s.bytes -= len(item.data)
	return item.data, true
}

// dropWindow discards anything parked for a window that is going away.
func (s *eventPayloadStore) dropWindow(windowID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, item := range s.items {
		if item.windowID == windowID {
			delete(s.items, id)
			s.bytes -= len(item.data)
		}
	}
}

func (s *eventPayloadStore) reap() {
	ticker := time.NewTicker(eventPayloadTTL)
	defer ticker.Stop()
	for {
		select {
		case <-s.stop:
			return
		case now := <-ticker.C:
			s.mu.Lock()
			for id, item := range s.items {
				if now.Sub(item.created) > eventPayloadTTL {
					delete(s.items, id)
					s.bytes -= len(item.data)
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *eventPayloadStore) close() {
	s.janitor.Do(func() {}) // ensure reap is never started after close
	select {
	case <-s.stop:
	default:
		close(s.stop)
	}
}
