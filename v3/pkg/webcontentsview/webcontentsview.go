package webcontentsview

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// WebContentsViewOptions represents the options for creating a WebContentsView.
type WebContentsViewOptions struct {
	Name           string
	URL            string
	HTML           string
	Bounds         application.Rect
	WebPreferences WebPreferences
}

// WebContentsView represents a native webview that can be embedded into a window.
type WebContentsView struct {
	mu          sync.RWMutex
	operationMu sync.Mutex
	options     WebContentsViewOptions
	id          uint
	impl        webContentsViewImpl
	visible     bool
	destroyed   bool
}

var webContentsViewID uintptr

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

// NewWebContentsView creates a new WebContentsView with the given options.
func NewWebContentsView(options WebContentsViewOptions) *WebContentsView {
	result := &WebContentsView{
		id:      uint(atomic.AddUintptr(&webContentsViewID, 1)),
		options: options,
		visible: true,
	}
	result.impl = newWebContentsViewImpl(result)
	return result
}

// SetBounds sets the position and size of the WebContentsView relative to its parent.
func (v *WebContentsView) SetBounds(bounds application.Rect) {
	v.mu.Lock()
	if v.destroyed {
		v.mu.Unlock()
		return
	}
	v.options.Bounds = bounds
	v.mu.Unlock()
	v.withLiveImpl(func(impl webContentsViewImpl) { impl.setBounds(bounds) })
}

// SetURL loads the given URL into the WebContentsView.
func (v *WebContentsView) SetURL(url string) {
	v.mu.Lock()
	if v.destroyed {
		v.mu.Unlock()
		return
	}
	v.options.URL = url
	v.options.HTML = ""
	v.mu.Unlock()
	v.withLiveImpl(func(impl webContentsViewImpl) { impl.setURL(url) })
}

// SetHTML loads inline HTML into the WebContentsView. It replaces any pending
// URL configured through WebContentsViewOptions or SetURL.
func (v *WebContentsView) SetHTML(html string) {
	v.mu.Lock()
	if v.destroyed {
		v.mu.Unlock()
		return
	}
	v.options.HTML = html
	v.options.URL = ""
	v.mu.Unlock()
	v.withLiveImpl(func(impl webContentsViewImpl) { impl.setHTML(html) })
}

// ExecJS executes the given javascript in the WebContentsView.
func (v *WebContentsView) ExecJS(js string) {
	v.withLiveImpl(func(impl webContentsViewImpl) { impl.execJS(js) })
}

// GoBack navigates to the previous page in history.
func (v *WebContentsView) GoBack() {
	v.withLiveImpl(func(impl webContentsViewImpl) { impl.goBack() })
}

// GetURL returns the current URL of the view.
func (v *WebContentsView) GetURL() string {
	var url string
	v.withLiveImpl(func(impl webContentsViewImpl) { url = impl.getURL() })
	return url
}

// TakeSnapshot returns a base64 encoded PNG of the current view.
func (v *WebContentsView) TakeSnapshot() string {
	var snapshot string
	v.withLiveImpl(func(impl webContentsViewImpl) { snapshot = impl.takeSnapshot() })
	return snapshot
}

// Show makes an attached WebContentsView visible. Native child views are
// rendered above the host webview, so use this rather than relying on a
// zero-sized frontend placeholder to hide one.
func (v *WebContentsView) Show() {
	v.mu.Lock()
	if v.destroyed {
		v.mu.Unlock()
		return
	}
	v.visible = true
	v.mu.Unlock()
	v.withLiveImpl(func(impl webContentsViewImpl) { impl.setVisible(true) })
}

// Hide hides an attached WebContentsView without removing it from its host
// window. A later Show restores the same native view and browser session.
func (v *WebContentsView) Hide() {
	v.mu.Lock()
	if v.destroyed {
		v.mu.Unlock()
		return
	}
	v.visible = false
	v.mu.Unlock()
	v.withLiveImpl(func(impl webContentsViewImpl) { impl.setVisible(false) })
}

// Attach binds the WebContentsView to a Wails Window.
func (v *WebContentsView) Attach(window application.Window) {
	v.withLiveImpl(func(impl webContentsViewImpl) { impl.attach(window) })
}

// Detach removes the WebContentsView from the Wails Window.
func (v *WebContentsView) Detach() {
	v.withLiveImpl(func(impl webContentsViewImpl) { impl.detach() })
}

// Destroy permanently releases the native browser surface. Remove the view
// from its owner first with WebviewWindow.RemoveChildView so the window no
// longer retains this ChildView; a destroyed view cannot be attached again.
func (v *WebContentsView) Destroy() {
	v.operationMu.Lock()
	defer v.operationMu.Unlock()
	v.mu.Lock()
	if v.destroyed {
		v.mu.Unlock()
		return
	}
	v.visible = false
	v.destroyed = true
	impl := v.impl
	v.mu.Unlock()
	impl.detach()
	impl.destroy()
}

// withLiveImpl serializes access to the platform-native view. State is not
// locked during the callback because platform attach paths obtain snapshots of
// that state while running on the UI thread.
func (v *WebContentsView) withLiveImpl(fn func(webContentsViewImpl)) bool {
	v.operationMu.Lock()
	defer v.operationMu.Unlock()
	v.mu.RLock()
	if v.destroyed {
		v.mu.RUnlock()
		return false
	}
	impl := v.impl
	v.mu.RUnlock()
	fn(impl)
	return true
}

func (v *WebContentsView) optionsSnapshot() WebContentsViewOptions {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.options
}

func (v *WebContentsView) visibleSnapshot() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.visible
}

// webContentsViewImpl is the interface that platform-specific implementations must satisfy.
type webContentsViewImpl interface {
	setBounds(bounds application.Rect)
	setURL(url string)
	setHTML(html string)
	execJS(js string)
	goBack()
	getURL() string
	takeSnapshot() string
	setVisible(visible bool)
	attach(window application.Window)
	detach()
	destroy()
	nativeView() unsafe.Pointer
}
