//go:build linux && !android && (cgo || purego)

package webview

import "context"

func (r *request) Close() error {
	var err error
	if r.body != nil {
		err = r.body.Close()
	}
	r.Response().Finish()
	r.lifetime.close()
	return err
}

func (r *request) Context() context.Context { return r.lifetime.Context }

func (r *request) Response() ResponseWriter {
	if r.rw != nil {
		return r.rw
	}

	r.rw = &responseWriter{req: r.req}
	return r.rw
}

func (l *requestLifetime) abandon() {
	l.state.Lock()
	if l.state.nativeReleased {
		l.state.Unlock()
		return
	}
	l.state.nativeReleased = true
	l.state.Unlock()
	l.cancel()
	// Keep the cgo handle alive because the native toggle reference remains
	// attached until WebKit tears down the object with the application. There is
	// no live GTK loop left on which it could be removed safely.
	l.removeActive()
}

func (l *requestLifetime) close() {
	l.once.Do(func() {
		if !invokeOnMainSyncIfEnabled(l.releaseNative) {
			// The GTK loop is gone, so no native callback can arrive after this
			// point. Do not release the GObject from the worker; process shutdown
			// will reclaim it with the rest of the WebKit instance.
			l.abandon()
		}
	})
}

func (l *requestLifetime) removeActive() {
	activeRequestLifetimes.Lock()
	delete(activeRequestLifetimes.values, l)
	activeRequestLifetimes.Unlock()
}

// CloseActiveRequests releases all native request references while the GTK
// thread is still active. The application calls this immediately before it
// quits the main loop, covering handlers that outlive window destruction.
func CloseActiveRequests() {
	activeRequestLifetimes.Lock()
	lifetimes := make([]*requestLifetime, 0, len(activeRequestLifetimes.values))
	for lifetime := range activeRequestLifetimes.values {
		lifetimes = append(lifetimes, lifetime)
	}
	activeRequestLifetimes.Unlock()

	for _, lifetime := range lifetimes {
		lifetime.releaseNative()
	}
}

func (l *requestLifetime) releaseNative() {
	l.state.Lock()
	if l.state.nativeReleased {
		l.state.Unlock()
		return
	}
	l.state.nativeReleased = true
	l.state.Unlock()

	// Remove the native callback before deleting its Go handle. Doing both on
	// the GTK thread serializes this with WebKit's ref changes.
	releaseRequestRef(l.object, uintptr(l.handle))
	l.cancel()
	l.deleteHandle()
	l.removeActive()
}
