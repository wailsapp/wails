package webcontentsview

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"sync/atomic"
	"unsafe"
)

// WebContentsViewOptions represents the options for creating a WebContentsView.
type WebContentsViewOptions struct {
	Name           string
	URL            string
	HTML           string
	Bounds         application.Rect
	Layout         *Layout
	WebPreferences WebPreferences
}

// WebContentsView represents a native webview that can be embedded into a window.
type WebContentsView struct {
	options                 WebContentsViewOptions
	id                      uint
	impl                    webContentsViewImpl
	automation              *automationState
	window                  application.Window
	clearLayoutResizeHandle func()
}

var webContentsViewID uintptr

// NewWebContentsView creates a new WebContentsView with the given options.
func NewWebContentsView(options WebContentsViewOptions) *WebContentsView {
	result := &WebContentsView{
		id:      uint(atomic.AddUintptr(&webContentsViewID, 1)),
		options: options,
	}
	result.automation = newAutomationState(result)
	registerAutomationView(result)
	result.impl = newWebContentsViewImpl(result)
	return result
}

// SetBounds sets the position and size of the WebContentsView relative to its parent.
func (v *WebContentsView) SetBounds(bounds application.Rect) {
	v.options.Layout = nil
	v.stopLayoutSync()
	v.updateBounds(bounds)
	v.impl.setBounds(bounds)
}

// SetLayout sets the position and size of the WebContentsView relative to its parent
// using absolute and percentage-based values.
func (v *WebContentsView) SetLayout(layout Layout) {
	v.options.Layout = &layout
	v.startLayoutSync()
	v.applyLayout()
}

// SetURL loads the given URL into the WebContentsView.
func (v *WebContentsView) SetURL(url string) {
	v.options.URL = url
	v.updateAutomationURL(url)
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
	v.window = window
	v.impl.attach(window)
	v.startLayoutSync()
	v.applyLayout()
	v.updateAutomationAttached(true)
}

// Detach removes the WebContentsView from the Wails Window.
func (v *WebContentsView) Detach() {
	v.stopLayoutSync()
	v.window = nil
	v.impl.detach()
	v.updateAutomationAttached(false)
}

func (v *WebContentsView) startLayoutSync() {
	v.stopLayoutSync()
	if v.window == nil || v.options.Layout == nil {
		return
	}
	v.clearLayoutResizeHandle = v.window.OnWindowEvent(events.Common.WindowDidResize, func(event *application.WindowEvent) {
		v.applyLayout()
	})
}

func (v *WebContentsView) stopLayoutSync() {
	if v.clearLayoutResizeHandle != nil {
		v.clearLayoutResizeHandle()
		v.clearLayoutResizeHandle = nil
	}
}

func (v *WebContentsView) applyLayout() {
	if v.window == nil || v.options.Layout == nil {
		return
	}

	application.InvokeSync(func() {
		reference := v.impl.layoutReferenceRect()
		if reference.Width == 0 && reference.Height == 0 {
			return
		}
		bounds := v.options.Layout.Resolve(reference)
		v.updateBounds(bounds)
		v.impl.setBounds(bounds)
	})
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
	layoutReferenceRect() application.Rect
	nativeView() unsafe.Pointer
}
