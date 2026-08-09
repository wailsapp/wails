//go:build windows && !server

package webcontentsview

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/webview2/pkg/edge"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

type windowsWebContentsView struct {
	parent     *WebContentsView
	chromium   *edge.Chromium
	hwnd       w32.HWND
	nativeHWND w32.HWND
	pendingJS  []string
}

func newWebContentsViewImpl(parent *WebContentsView) webContentsViewImpl {
	chromium := edge.NewChromium()
	// The host Wails window and this native panel construct independent
	// WebView2 environments. WebView2 only permits environments sharing a
	// user-data folder when every environment option is identical; the Wails
	// host applies its application-level browser flags while this panel must not
	// inherit private host configuration implicitly. Use a stable, app- and
	// view-scoped folder to avoid an incompatible second environment and to keep
	// panel sessions independent of the host frontend's browser session.
	chromium.DataPath = windowsWebContentsUserDataPath(parent)

	result := &windowsWebContentsView{
		parent:   parent,
		chromium: chromium,
	}

	return result
}

func windowsWebContentsUserDataPath(parent *WebContentsView) string {
	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return ""
	}

	appIdentity := "wails"
	if executable, err := os.Executable(); err == nil {
		appIdentity = executable
	}
	options := parent.optionsSnapshot()
	viewIdentity := options.Name
	if viewIdentity == "" {
		viewIdentity = fmt.Sprintf("unnamed-%d", parent.id)
	}
	appHash := sha256.Sum256([]byte(appIdentity))
	viewHash := sha256.Sum256([]byte(viewIdentity))
	return filepath.Join(
		configDirectory,
		"Wails",
		fmt.Sprintf("app-%x", appHash[:8]),
		"WebContentsView",
		fmt.Sprintf("view-%x", viewHash[:8]),
	)
}

func (w *windowsWebContentsView) configure() {
	options := w.parent.optionsSnapshot()
	settings, err := w.chromium.GetSettings()
	if err == nil {
		if !options.WebPreferences.DevTools.IsSet() || options.WebPreferences.DevTools.Get() {
			settings.PutAreDevToolsEnabled(true)
			settings.PutAreDefaultContextMenusEnabled(true)
		} else {
			settings.PutAreDevToolsEnabled(false)
			settings.PutAreDefaultContextMenusEnabled(false)
		}

		if !options.WebPreferences.Javascript.IsSet() || options.WebPreferences.Javascript.Get() {
			settings.PutIsScriptEnabled(true)
		} else {
			settings.PutIsScriptEnabled(false)
		}

		if options.WebPreferences.ZoomFactor > 0 {
			w.chromium.PutZoomFactor(options.WebPreferences.ZoomFactor)
		}

		if options.WebPreferences.UserAgent != "" {
			settings.PutUserAgent(options.WebPreferences.UserAgent)
		}
	}

}

func (w *windowsWebContentsView) setBounds(bounds application.Rect) {
	if w.chromium != nil && w.chromium.IsReady() {
		edgeBounds := edge.Rect{
			Left:   int32(bounds.X),
			Top:    int32(bounds.Y),
			Right:  int32(bounds.X + bounds.Width),
			Bottom: int32(bounds.Y + bounds.Height),
		}
		application.InvokeSync(func() {
			w.chromium.ResizeWithBounds(&edgeBounds)
			w.bringNativeSurfaceToFront()
		})
	}
}

func (w *windowsWebContentsView) setURL(url string) {
	if w.chromium != nil && w.chromium.IsReady() {
		application.InvokeSync(func() {
			log.Printf("[WebContentsView] navigating native browser to %s", url)
			w.chromium.Navigate(url)
		})
	}
}

func (w *windowsWebContentsView) setHTML(html string) {
	if w.chromium != nil && w.chromium.IsReady() {
		application.InvokeSync(func() {
			w.chromium.NavigateToString(html)
		})
	}
}

func (w *windowsWebContentsView) goBack() {
	if w.chromium != nil && w.chromium.IsReady() {
		application.InvokeSync(func() {
			_ = w.chromium.GoBack()
		})
	}
}

func (w *windowsWebContentsView) getURL() string {
	if w.chromium != nil && w.chromium.IsReady() {
		var url string
		application.InvokeSync(func() {
			url, _ = w.chromium.Source()
		})
		if url != "" {
			return url
		}
	}
	return w.parent.optionsSnapshot().URL
}

func (w *windowsWebContentsView) execJS(js string) {
	if w.chromium == nil || !w.chromium.IsReady() {
		w.pendingJS = append(w.pendingJS, js)
		return
	}
	application.InvokeSync(func() {
		w.chromium.Eval(js)
	})
}

func (w *windowsWebContentsView) takeSnapshot() string {
	if w.chromium == nil || !w.chromium.IsReady() {
		return ""
	}
	type screenshotResult struct {
		data []byte
		err  error
	}
	resultCh := make(chan screenshotResult, 1)
	var requestErr error
	application.InvokeSync(func() {
		requestErr = w.chromium.CapturePreview(func(data []byte, err error) {
			resultCh <- screenshotResult{data: data, err: err}
		})
	})
	if requestErr != nil {
		return ""
	}
	select {
	case result := <-resultCh:
		if result.err != nil || len(result.data) == 0 {
			return ""
		}
		return "data:image/png;base64," + base64.StdEncoding.EncodeToString(result.data)
	case <-time.After(30 * time.Second):
		return ""
	}
}

func (w *windowsWebContentsView) setVisible(visible bool) {
	if w.chromium == nil || !w.chromium.IsReady() {
		return
	}
	application.InvokeSync(func() {
		if visible {
			_ = w.chromium.Show()
			w.bringNativeSurfaceToFront()
		} else {
			_ = w.chromium.Hide()
		}
	})
}

func (w *windowsWebContentsView) attach(window application.Window) {
	if window.NativeWindow() == nil || w.chromium == nil {
		return
	}
	if w.chromium.IsReady() {
		options := w.parent.optionsSnapshot()
		visible := w.parent.visibleSnapshot()
		w.setBounds(options.Bounds)
		application.InvokeSync(func() {
			if visible {
				_ = w.chromium.Show()
			} else {
				_ = w.chromium.Hide()
			}
		})
		return
	}
	application.InvokeSync(func() {
		w.hwnd = w32.HWND(uintptr(window.NativeWindow()))
		knownChromeChildren := chromeWidgetChildren(w.hwnd)
		if !w.chromium.Embed(uintptr(w.hwnd)) || !w.chromium.IsReady() {
			log.Printf("[WebContentsView] native WebView2 did not become ready")
			return
		}
		w.nativeHWND = newChromeWidgetChild(w.hwnd, knownChromeChildren)
		if w.nativeHWND != 0 {
			log.Printf("[WebContentsView] native browser surface ready: hwnd=0x%x", uintptr(w.nativeHWND))
			w.bringNativeSurfaceToFront()
		} else {
			log.Printf("[WebContentsView] native browser surface ready; no Chrome child HWND was discovered")
		}
		w.configure()
		options := w.parent.optionsSnapshot()
		visible := w.parent.visibleSnapshot()
		w.setBounds(options.Bounds)
		if visible {
			_ = w.chromium.Show()
		} else {
			_ = w.chromium.Hide()
		}

		if options.URL != "" {
			log.Printf("[WebContentsView] navigating native browser to initial URL %s", options.URL)
			w.chromium.Navigate(options.URL)
		} else if options.HTML != "" {
			w.chromium.NavigateToString(options.HTML)
		}
		for _, js := range w.pendingJS {
			w.chromium.Eval(js)
		}
		w.pendingJS = nil
	})
}

// chromeWidgetChildren returns WebView2's HWND-hosted renderer windows below
// a Wails window. A normal Wails window owns one such child already; a native
// WebContentsView creates another. The latter must be promoted after creation
// or the host webview can cover it completely despite correct controller bounds.
func chromeWidgetChildren(parent w32.HWND) map[w32.HWND]struct{} {
	children := make(map[w32.HWND]struct{})
	w32.EnumChildWindows(parent, func(hwnd w32.HWND, _ w32.LPARAM) w32.LRESULT {
		if strings.HasPrefix(w32.GetClassName(hwnd), "Chrome_WidgetWin_") {
			children[hwnd] = struct{}{}
		}
		return 1
	})
	return children
}

func newChromeWidgetChild(parent w32.HWND, existing map[w32.HWND]struct{}) w32.HWND {
	for hwnd := range chromeWidgetChildren(parent) {
		if _, found := existing[hwnd]; !found {
			return hwnd
		}
	}
	return 0
}

func (w *windowsWebContentsView) bringNativeSurfaceToFront() {
	if w.nativeHWND == 0 {
		return
	}
	if !w32.SetWindowPos(
		w.nativeHWND,
		w32.HWND_TOP,
		0,
		0,
		0,
		0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOACTIVATE,
	) {
		log.Printf("[WebContentsView] could not promote native browser surface: hwnd=0x%x", uintptr(w.nativeHWND))
	}
}

func (w *windowsWebContentsView) detach() {
	if w.chromium != nil && w.chromium.IsReady() {
		application.InvokeSync(func() {
			_ = w.chromium.Hide()
		})
	}
}

func (w *windowsWebContentsView) destroy() {
	if w.chromium == nil {
		return
	}
	application.InvokeSync(w.chromium.Close)
	w.hwnd = 0
	w.nativeHWND = 0
}

func (w *windowsWebContentsView) nativeView() unsafe.Pointer {
	return unsafe.Pointer(w.chromium)
}
