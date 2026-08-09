//go:build ios

package webcontentsview

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"unsafe"
)

type iosWebContentsView struct {
	parent *WebContentsView
}

func newWebContentsViewImpl(parent *WebContentsView) webContentsViewImpl {
	return &iosWebContentsView{parent: parent}
}

func (w *iosWebContentsView) setBounds(bounds application.Rect) {}
func (w *iosWebContentsView) setURL(url string)                 {}
func (w *iosWebContentsView) setHTML(html string)               {}
func (w *iosWebContentsView) execJS(js string)                  {}
func (w *iosWebContentsView) goBack()                           {}
func (w *iosWebContentsView) getURL() string                    { return "" }
func (w *iosWebContentsView) takeSnapshot() string              { return "" }
func (w *iosWebContentsView) setVisible(bool)                   {}

func (w *iosWebContentsView) attach(window application.Window) {}
func (w *iosWebContentsView) detach()                          {}
func (w *iosWebContentsView) destroy()                         {}
func (w *iosWebContentsView) nativeView() unsafe.Pointer       { return nil }
