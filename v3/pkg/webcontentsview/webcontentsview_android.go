//go:build android

package webcontentsview

import (
	"unsafe"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type androidWebContentsView struct {
	parent *WebContentsView
}

func newWebContentsViewImpl(parent *WebContentsView) webContentsViewImpl {
	return &androidWebContentsView{parent: parent}
}

func (w *androidWebContentsView) setBounds(bounds application.Rect) {}
func (w *androidWebContentsView) setURL(url string)     {}
func (w *androidWebContentsView) execJS(js string)      {}
func (w *androidWebContentsView) attach(window application.Window) {}
func (w *androidWebContentsView) detach() {}
func (w *androidWebContentsView) nativeView() unsafe.Pointer { return nil }
