//go:build android

package webcontentsview

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"unsafe"
)

type androidWebContentsView struct {
	parent *WebContentsView
}

func newWebContentsViewImpl(parent *WebContentsView) webContentsViewImpl {
	return &androidWebContentsView{parent: parent}
}

func (w *androidWebContentsView) setBounds(bounds application.Rect) {}
func (w *androidWebContentsView) setURL(url string)                 {}
func (w *androidWebContentsView) setHTML(html string)               {}
func (w *androidWebContentsView) execJS(js string)                  {}
func (w *androidWebContentsView) goBack()                           {}
func (w *androidWebContentsView) getURL() string                    { return "" }
func (w *androidWebContentsView) takeSnapshot() string              { return "" }
func (w *androidWebContentsView) setVisible(bool)                   {}

func (w *androidWebContentsView) attach(window application.Window) {}
func (w *androidWebContentsView) detach()                          {}
func (w *androidWebContentsView) destroy()                         {}
func (w *androidWebContentsView) nativeView() unsafe.Pointer       { return nil }
