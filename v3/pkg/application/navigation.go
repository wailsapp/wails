package application

// NavigationRequest describes a navigation that a window's OnBeforeNavigate
// is being asked to allow or refuse.
type NavigationRequest struct {
	// URL is the absolute URL the frame is about to load.
	URL string

	// IsMainFrame reports whether the navigation replaces the window's
	// top-level document rather than the content of an iframe. A request to
	// open a new window (target="_blank", window.open) is a new top-level
	// document too, and is reported as the main frame.
	IsMainFrame bool
}

// shouldNavigate applies the window's OnBeforeNavigate to request. A window
// without one allows everything, as the webview does by default.
func (w *WebviewWindow) shouldNavigate(request NavigationRequest) bool {
	policy := w.options.OnBeforeNavigate
	if policy == nil {
		return true
	}
	return policy(w, request)
}

// interceptsNavigation reports whether the window identified by windowID has
// an OnBeforeNavigate, so the platform layer can let navigations through
// without a round trip when it does not.
func interceptsNavigation(windowID uint) bool {
	window := navigationWindow(windowID)
	return window != nil && window.options.OnBeforeNavigate != nil
}

// resolveNavigation asks the window identified by windowID whether request
// may go ahead, and passes the answer to decide exactly once.
//
// The webview delivers navigation requests on the main thread and holds the
// navigation until it is answered. The policy runs on a goroutine instead, so
// that it can call window and application methods that themselves wait on
// the main thread. Every path that does not reach a positive answer refuses:
// a window that has gone away, a policy that returns false, and one that
// panics.
func resolveNavigation(windowID uint, request NavigationRequest, decide func(allow bool)) {
	go func() {
		allow := false
		defer func() { decide(allow) }()
		defer handlePanic()

		window := navigationWindow(windowID)
		if window == nil {
			return
		}
		allow = window.shouldNavigate(request)
	}()
}

func navigationWindow(windowID uint) *WebviewWindow {
	if globalApplication == nil {
		return nil
	}
	window, ok := globalApplication.Window.GetByID(windowID)
	if !ok {
		return nil
	}
	webviewWindow, _ := window.(*WebviewWindow)
	return webviewWindow
}
