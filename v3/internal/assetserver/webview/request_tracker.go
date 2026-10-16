package webview

import (
	"context"
	"crypto/rand"
	"sync"
)

// RequestIDHeader is a private bridge marker, removed before calling handlers.
const RequestIDHeader = "X-Wails-Request-Id"

// Keepalive registration precedes native interception. Bound these hints so
// scripts cannot accumulate unlimited registrations without issuing requests.
const maxKeepaliveRequests = 1024

// RequestTracker connects browser request identities to Go contexts. A tracker
// belongs to one WebView; neither URLs nor native pointers identify requests.
type RequestTracker struct {
	mu         sync.Mutex
	ctx        context.Context
	cancel     context.CancelFunc
	byID       map[string]*trackedRequest
	byToken    map[string]*trackedRequest
	keepalives int
}

type trackedRequest struct {
	id, token  string
	ctx        context.Context
	cancel     context.CancelFunc
	claimed    bool
	persistent bool
	active     int
	failed     bool
}

func NewRequestTracker(parent context.Context) *RequestTracker {
	ctx, cancel := context.WithCancel(parent)
	return &RequestTracker{ctx: ctx, cancel: cancel, byID: make(map[string]*trackedRequest), byToken: make(map[string]*trackedRequest)}
}

// Start runs before the browser resumes an intercepted request. Redirects may
// reuse a network ID, so each interception receives a fresh, unguessable token.
func (t *RequestTracker) Start(id string) string { return t.start(id, false) }

// StartKeepalive spans redirects and outlives the JavaScript response consumer.
func (t *RequestTracker) StartKeepalive(id string) string { return t.start(id, true) }

func (t *RequestTracker) start(id string, persistent bool) string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if id == "" || t.ctx.Err() != nil {
		return ""
	}
	t.remove(t.byID[id])
	if persistent && t.keepalives >= maxKeepaliveRequests {
		return ""
	}
	if persistent {
		t.keepalives++
	}
	ctx, cancel := context.WithCancel(t.ctx)
	r := &trackedRequest{id: id, token: rand.Text(), ctx: ctx, cancel: cancel, persistent: persistent}
	t.byID[id], t.byToken[r.token] = r, r
	return r.token
}

// Cancel also removes unclaimed requests: an abort can arrive before the
// WebResourceRequested callback, or prevent that callback entirely.
func (t *RequestTracker) Cancel(id string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.remove(t.byID[id])
}

// FailKeepalive releases failed delivery once all dispatched handlers finish.
func (t *RequestTracker) FailKeepalive(id string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if entry := t.byID[id]; entry != nil {
		entry.failed = true
		if entry.active == 0 {
			t.remove(entry)
		}
	}
}

// Token returns the marker for a request registered before interception.
func (t *RequestTracker) Token(id string) string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if entry := t.byID[id]; entry != nil {
		return entry.token
	}
	return ""
}

// Close cancels all requests, including requests lacking a browser identity.
func (t *RequestTracker) Close() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.cancel()
	clear(t.byID)
	clear(t.byToken)
	t.keepalives = 0
}

func (t *RequestTracker) remove(r *trackedRequest) {
	if r == nil {
		return
	}
	r.cancel()
	delete(t.byToken, r.token)
	if t.byID[r.id] == r {
		delete(t.byID, r.id)
		if r.persistent {
			t.keepalives--
		}
	}
}

// Wrap consumes the marker and transfers context ownership to the request.
// A missing nonempty token means an earlier abort (or an invalid marker), and
// must not resurrect the request. Closing an old redirect cannot remove its successor.
func (t *RequestTracker) Wrap(r Request) Request {
	header, _ := r.Header()
	token := header.Get(RequestIDHeader)
	header.Del(RequestIDHeader)
	t.mu.Lock()
	entry := t.byToken[token]
	ctx, release := context.WithCancel(t.ctx)
	if entry != nil && entry.persistent {
		release()
		entry.active++
		ctx, release = context.WithCancel(entry.ctx)
		cancel := release
		release = func() {
			cancel()
			t.mu.Lock()
			defer t.mu.Unlock()
			entry.active--
			if entry.failed && entry.active == 0 {
				t.remove(entry)
			}
		}
	} else if entry != nil && !entry.claimed {
		entry.claimed = true
		release()
		ctx = entry.ctx
		release = func() { t.mu.Lock(); defer t.mu.Unlock(); t.remove(entry) }
	} else if token != "" {
		release()
	}
	t.mu.Unlock()
	return &contextRequest{Request: r, ctx: ctx, release: release}
}

type contextRequest struct {
	Request
	ctx     context.Context
	release context.CancelFunc
	once    sync.Once
	err     error
}

func (r *contextRequest) Context() context.Context { return r.ctx }

func (r *contextRequest) Close() error {
	r.once.Do(func() { r.release(); r.err = r.Request.Close() })
	return r.err
}
