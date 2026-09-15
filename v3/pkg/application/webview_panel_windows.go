//go:build windows && !server

package application

import (
	"fmt"
	"math"
	"net/url"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/internal/assetserver/webview"
	"github.com/wailsapp/wails/v3/internal/webview2/pkg/edge"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

type windowsPanelImpl struct {
	panel          *WebviewPanel
	parent         *windowsWebviewWindow
	chromium       *edge.Chromium
	hwnd           w32.HWND
	initialRequest bool
	initialURL     string
}

func newPanelImpl(panel *WebviewPanel) webviewPanelImpl {
	parent, ok := panel.parent.impl.(*windowsWebviewWindow)
	if !ok || parent.hwnd == 0 {
		return nil
	}
	return &windowsPanelImpl{panel: panel, parent: parent}
}

func (p *windowsPanelImpl) scale() float64 {
	if w32.HasGetDpiForWindowFunc() {
		if dpi := w32.GetDpiForWindow(p.parent.hwnd); dpi != 0 {
			return float64(dpi) / 96
		}
	}
	return 1
}

func (p *windowsPanelImpl) create() {
	options := p.panel.snapshotOptions()
	style := uint(w32.WS_CHILD | w32.WS_CLIPCHILDREN | w32.WS_CLIPSIBLINGS)
	if *options.Visible {
		style |= w32.WS_VISIBLE
	}
	p.hwnd = w32.CreateWindowEx(0, w32.MustStringToUTF16Ptr("STATIC"), nil, style,
		0, 0, 1, 1, p.parent.hwnd, 0, w32.GetModuleHandle(""), nil)
	if p.hwnd == 0 {
		globalApplication.error("create panel child window: %v", w32.GetLastError())
		return
	}
	p.chromium = edge.NewChromium()
	if globalApplication.options.ErrorHandler != nil {
		p.chromium.SetErrorCallback(globalApplication.options.ErrorHandler)
	}
	// Scope persistent browser data to the application, rather than sharing a
	// global panel-1 profile with unrelated Wails applications.
	base := globalApplication.options.Windows.WebviewUserDataPath
	if base == "" {
		exe, err := os.Executable()
		if err != nil {
			p.destroy()
			globalApplication.error("panel executable path: %v", err)
			return
		}
		base = filepath.Join(os.Getenv("AppData"), filepath.Base(exe))
	}
	p.chromium.DataPath = filepath.Join(base, "panels", fmt.Sprint(p.panel.id))
	p.chromium.BrowserPath = globalApplication.options.Windows.WebviewBrowserPath
	p.chromium.WebResourceRequestedCallback = p.processRequest
	if !p.chromium.Embed(p.hwnd) {
		p.destroy()
		return
	}
	// Embed pumps messages: destruction can be requested during creation.
	if p.panel.isDestroyed() {
		return
	}
	p.chromium.AddWebResourceRequestedFilter("*", edge.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)
	settings, err := p.chromium.GetSettings()
	if err != nil {
		globalApplication.error("panel settings: %v", err)
		p.destroy()
		return
	}
	enabled := globalApplication.isDebugMode
	if options.DevToolsEnabled != nil {
		enabled = *options.DevToolsEnabled
	}
	if err := settings.PutAreDevToolsEnabled(enabled); err != nil {
		globalApplication.error("panel devtools: %v", err)
	}
	if err := settings.PutAreDefaultContextMenusEnabled(enabled); err != nil {
		globalApplication.error("panel context menu: %v", err)
	}
	if options.UserAgent != "" {
		if err := settings.PutUserAgent(options.UserAgent); err != nil {
			globalApplication.error("panel user agent: %v", err)
		}
	}
	p.setZoom(options.Zoom)
	colour := options.BackgroundColour
	if options.Transparent {
		colour = RGBA{}
	}
	p.chromium.SetBackgroundColour(colour.Red, colour.Green, colour.Blue, colour.Alpha)
	p.setBounds(Rect{X: options.X, Y: options.Y, Width: options.Width, Height: options.Height})
	p.initialURL, p.initialRequest = panelRequestURL(options.URL), len(options.Headers) > 0
	if options.URL != "" {
		p.chromium.Navigate(options.URL)
	}
	if *options.Visible {
		p.show()
	} else {
		p.hide()
	}
	if enabled && options.OpenInspectorOnStartup {
		p.openDevTools()
	}
}

// processRequest supplies initial navigation headers and local assets without
// installing the parent window's message bridge in external panel content.
func (p *windowsPanelImpl) processRequest(req *edge.ICoreWebView2WebResourceRequest, args *edge.ICoreWebView2WebResourceRequestedEventArgs) {
	uri, err := req.GetUri()
	if err != nil {
		return
	}
	if p.initialRequest && panelRequestURL(uri) == p.initialURL {
		p.initialRequest = false
		if headers, err := req.GetHeaders(); err == nil {
			for key, value := range p.panel.snapshotOptions().Headers {
				if err := headers.SetHeader(key, value); err != nil {
					globalApplication.error("panel request header: %v", err)
				}
			}
			headers.Release()
		}
	}
	parsed, err := url.Parse(uri)
	if err != nil || parsed.Scheme != "http" || parsed.Hostname() != "wails.localhost" || globalApplication.assets == nil {
		return
	}
	request, err := webview.NewRequest(p.chromium.Environment(), args, InvokeSync)
	if err != nil {
		globalApplication.error("panel asset request: %v", err)
		return
	}
	webviewRequests <- &webViewAssetRequest{Request: request, windowId: p.parent.parent.id, windowName: p.parent.parent.Name()}
}

// A browser adds '/' to an empty HTTP path and strips fragments from requests.
func panelRequestURL(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	parsed.Fragment, parsed.RawFragment = "", ""
	if parsed.Path == "" {
		parsed.Path = "/"
	}
	return parsed.String()
}

func (p *windowsPanelImpl) destroy() {
	if p.chromium != nil {
		p.chromium.ShuttingDown()
		if controller := p.chromium.GetController(); controller != nil {
			controller.Close()
		}
		p.chromium = nil
	}
	if p.hwnd != 0 {
		w32.DestroyWindow(p.hwnd)
		p.hwnd = 0
	}
}

func (p *windowsPanelImpl) setBounds(bounds Rect) {
	if p.hwnd == 0 {
		return
	}
	scale := p.scale()
	physical := func(v int) int { return int(math.Round(float64(v) * scale)) }
	w32.SetWindowPos(p.hwnd, 0, physical(bounds.X), physical(bounds.Y), physical(bounds.Width), physical(bounds.Height), w32.SWP_NOZORDER|w32.SWP_NOACTIVATE)
	if p.chromium != nil {
		p.chromium.Resize()
	}
}

func (p *windowsPanelImpl) bounds() Rect {
	if p.hwnd == 0 {
		return Rect{}
	}
	rect := w32.GetWindowRect(p.hwnd)
	if rect == nil {
		return Rect{}
	}
	x, y := w32.ClientToScreen(p.parent.hwnd, 0, 0)
	scale := p.scale()
	dip := func(v int) int { return int(math.Round(float64(v) / scale)) }
	return Rect{X: dip(int(rect.Left) - x), Y: dip(int(rect.Top) - y), Width: dip(int(rect.Right - rect.Left)), Height: dip(int(rect.Bottom - rect.Top))}
}

func (p *windowsPanelImpl) setZIndex(_ int) {
	panels := p.panel.sortedSiblings()
	// Reorder every sibling so values, rather than the last setter, determine order.
	for _, panel := range panels {
		panel.destroyedLock.RLock()
		native, ok := panel.impl.(*windowsPanelImpl)
		panel.destroyedLock.RUnlock()
		if ok && native.hwnd != 0 {
			w32.SetWindowPos(native.hwnd, w32.HWND_TOP, 0, 0, 0, 0, w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOACTIVATE)
		}
	}
}

func (p *windowsPanelImpl) setURL(url string) {
	if p.chromium != nil {
		p.initialRequest = false
		p.chromium.Navigate(url)
	}
}
func (p *windowsPanelImpl) reload() {
	if p.chromium != nil {
		p.chromium.Eval("window.location.reload();")
	}
}
func (p *windowsPanelImpl) forceReload() { p.reload() }
func (p *windowsPanelImpl) show() {
	if p.hwnd != 0 {
		w32.ShowWindow(p.hwnd, w32.SW_SHOWNA)
	}
}
func (p *windowsPanelImpl) hide() {
	if p.hwnd != 0 {
		w32.ShowWindow(p.hwnd, w32.SW_HIDE)
	}
}
func (p *windowsPanelImpl) isVisible() bool {
	return p.hwnd != 0 && uint32(w32.GetWindowLong(p.hwnd, w32.GWL_STYLE))&w32.WS_VISIBLE != 0
}
func (p *windowsPanelImpl) setZoom(zoom float64) {
	if p.chromium != nil {
		p.chromium.PutZoomFactor(zoom)
	}
}
func (p *windowsPanelImpl) getZoom() float64 {
	if p.chromium != nil {
		if controller := p.chromium.GetController(); controller != nil {
			if zoom, err := controller.GetZoomFactor(); err == nil {
				return zoom
			}
		}
	}
	return p.panel.snapshotOptions().Zoom
}
func (p *windowsPanelImpl) openDevTools() {
	if p.chromium != nil {
		p.chromium.OpenDevToolsWindow()
	}
}
func (p *windowsPanelImpl) focus() {
	if p.chromium != nil {
		p.chromium.Focus()
	}
}
func (p *windowsPanelImpl) isFocused() bool {
	if p.hwnd == 0 {
		return false
	}
	for hwnd := w32.GetFocus(); hwnd != 0; hwnd = w32.GetParent(hwnd) {
		if hwnd == p.hwnd {
			return true
		}
	}
	return false
}
