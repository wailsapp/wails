package application

// focusHandoff tracks whether keyboard focus has been deliberately handed from
// the WebView to the native host window by tabbing out of the page
// (wailsapp/wails#3783). It holds no OS state so the policy can be unit tested
// on any platform; it is only used from the window's UI thread.
type focusHandoff struct {
	// onHost is true while the user has tabbed out of the WebView and focus
	// rests on the native window.
	onHost bool
}

// tabOut records that the WebView asked to move focus away (Tab / Shift+Tab off
// its last / first element). It reports whether the request should be treated
// as handled; programmatic requests are left to WebView2.
func (f *focusHandoff) tabOut(next, previous bool) bool {
	if !next && !previous {
		return false
	}
	f.onHost = true
	return true
}

// shouldRestoreWebview reports whether the WebView should reclaim focus when
// the host window receives it (WM_SETFOCUS). While focus is deliberately on
// the host it must not be stolen back, or Tab would loop inside the WebView.
func (f *focusHandoff) shouldRestoreWebview() bool {
	return !f.onHost
}

// tabIn is called on Tab / Shift+Tab while the host window has focus. When
// focus is on the host it returns the direction to enter the WebView
// (previous=true for Shift+Tab) and clears the hand-off. Ctrl-modified Tab
// remains available to key bindings without changing the hand-off.
func (f *focusHandoff) tabIn(shift, ctrl bool) (previous bool, ok bool) {
	if ctrl || !f.onHost {
		return false, false
	}
	f.onHost = false
	return shift, true
}

// reset clears the hand-off, e.g. when the window is deactivated so that
// re-activating it restores focus to the WebView.
func (f *focusHandoff) reset() {
	f.onHost = false
}
