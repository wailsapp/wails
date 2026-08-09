//go:build server && !android && !ios

package webcontentsview

import (
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// serverWebContentsView preserves the public API for Wails server builds,
// where no native window or browser engine exists. It intentionally exposes no
// native view.
type serverWebContentsView struct {
	parent *WebContentsView
}

func newWebContentsViewImpl(parent *WebContentsView) webContentsViewImpl {
	return &serverWebContentsView{parent: parent}
}

func (w *serverWebContentsView) setBounds(application.Rect) {}
func (w *serverWebContentsView) setURL(string)              {}
func (w *serverWebContentsView) setHTML(string)             {}
func (w *serverWebContentsView) execJS(string)              {}
func (w *serverWebContentsView) goBack()                    {}
func (w *serverWebContentsView) getURL() string             { return w.parent.optionsSnapshot().URL }
func (w *serverWebContentsView) takeSnapshot() string       { return "" }
func (w *serverWebContentsView) setVisible(bool)            {}
func (w *serverWebContentsView) attach(application.Window)  {}
func (w *serverWebContentsView) detach()                    {}
func (w *serverWebContentsView) destroy()                   {}
func (w *serverWebContentsView) nativeView() unsafe.Pointer { return nil }
