//go:build windows

package webcontentsview

import (
	"unsafe"

	"github.com/wailsapp/go-webview2/pkg/edge"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

type windowsWebContentsView struct {
	parent   *WebContentsView
	chromium *edge.Chromium
	hwnd     w32.HWND
}

func newWebContentsViewImpl(parent *WebContentsView) webContentsViewImpl {
	chromium := edge.NewChromium()

	result := &windowsWebContentsView{
		parent:   parent,
		chromium: chromium,
	}
	
	settings, err := chromium.GetSettings()
	if err == nil {
		if parent.options.WebPreferences.DevTools != application.Disabled {
			settings.PutAreDevToolsEnabled(true)
			settings.PutAreDefaultContextMenusEnabled(true)
		} else {
			settings.PutAreDevToolsEnabled(false)
			settings.PutAreDefaultContextMenusEnabled(false)
		}

		if parent.options.WebPreferences.Javascript != application.Disabled {
			settings.PutIsScriptEnabled(true)
		} else {
			settings.PutIsScriptEnabled(false)
		}
		
		if parent.options.WebPreferences.ZoomFactor > 0 {
			chromium.PutZoomFactor(parent.options.WebPreferences.ZoomFactor)
		}

		if parent.options.WebPreferences.UserAgent != "" {
			settings.PutUserAgent(parent.options.WebPreferences.UserAgent)
		}
	}

	return result
}

func (w *windowsWebContentsView) setBounds(bounds application.Rect) {
	if w.chromium != nil {
		edgeBounds := edge.Rect{
			Left:   int32(bounds.X),
			Top:    int32(bounds.Y),
			Right:  int32(bounds.X + bounds.Width),
			Bottom: int32(bounds.Y + bounds.Height),
		}
		w.chromium.ResizeWithBounds(edgeBounds)
	}
}

func (w *windowsWebContentsView) setURL(url string) {
	if w.chromium != nil {
		w.chromium.Navigate(url)
	}
}

func (w *windowsWebContentsView) execJS(js string) {
	if w.chromium != nil {
		w.chromium.Eval(js)
	}
}

func (w *windowsWebContentsView) attach(window application.Window) {
	if window.NativeWindow() != nil {
		w.hwnd = w32.HWND(window.NativeWindow())
		w.chromium.Embed(w.hwnd)
		
		w.chromium.Resize()
		w.chromium.Show()

		if w.parent.options.URL != "" {
			w.chromium.Navigate(w.parent.options.URL)
		}
	}
}

func (w *windowsWebContentsView) detach() {
	if w.chromium != nil {
		w.chromium.Hide()
	}
}

func (w *windowsWebContentsView) nativeView() unsafe.Pointer {
	return unsafe.Pointer(w.chromium)
}
