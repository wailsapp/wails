package webcontentsview

import (
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
	options WebContentsViewOptions
	id      uint
	impl    webContentsViewImpl
}

var webContentsViewID uintptr

// NewWebContentsView creates a new WebContentsView with the given options.
func NewWebContentsView(options WebContentsViewOptions) *WebContentsView {
	result := &WebContentsView{
		id:      uint(atomic.AddUintptr(&webContentsViewID, 1)),
		options: options,
	}
	result.impl = newWebContentsViewImpl(result)
	return result
}

// SetBounds sets the position and size of the WebContentsView relative to its parent.
func (v *WebContentsView) SetBounds(bounds application.Rect) {
	v.impl.setBounds(bounds)
}

// SetURL loads the given URL into the WebContentsView.
func (v *WebContentsView) SetURL(url string) {
	v.impl.setURL(url)
}

// ExecJS executes the given javascript in the WebContentsView.
func (v *WebContentsView) ExecJS(js string) {
	v.impl.execJS(js)
}

// GoBack navigates to the previous page in history.
func (v *WebContentsView) GoBack() {
	v.impl.goBack()
}

// GetURL returns the current URL of the view.
func (v *WebContentsView) GetURL() string {
	return v.impl.getURL()
}

// TakeSnapshot returns a base64 encoded PNG of the current view.
func (v *WebContentsView) TakeSnapshot() string {
	return v.impl.takeSnapshot()
}

// Attach binds the WebContentsView to a Wails Window.
func (v *WebContentsView) Attach(window application.Window) {
	v.impl.attach(window)
}

// Detach removes the WebContentsView from the Wails Window.
func (v *WebContentsView) Detach() {
	v.impl.detach()
}

// webContentsViewImpl is the interface that platform-specific implementations must satisfy.
type webContentsViewImpl interface {
	setBounds(bounds application.Rect)
	setURL(url string)
	execJS(js string)
	goBack()
	getURL() string
	takeSnapshot() string
	attach(window application.Window)
	detach()
	nativeView() unsafe.Pointer
}
