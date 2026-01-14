//go:build windows

package application

import (
	"fmt"
	"strings"

	"github.com/wailsapp/go-webview2/pkg/edge"
	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

type windowsPanelImpl struct {
	panel    *WebviewPanel
	parent   *windowsWebviewWindow
	chromium *edge.Chromium
	hwnd     w32.HWND // Child window handle to host the WebView2

	// Track navigation state
	navigationCompleted bool
}

func newPanelImpl(panel *WebviewPanel) webviewPanelImpl {
	parentWindow := panel.parent
	if parentWindow == nil || parentWindow.impl == nil {
		return nil
	}

	windowsParent, ok := parentWindow.impl.(*windowsWebviewWindow)
	if !ok {
		return nil
	}

	return &windowsPanelImpl{
		panel:  panel,
		parent: windowsParent,
	}
}

func (p *windowsPanelImpl) create() {
	options := p.panel.options

	// Create a child window to host the WebView2
	// We use WS_CHILD style to make it a child of the parent window
	style := uint(w32.WS_CHILD | w32.WS_CLIPSIBLINGS)
	if options.Visible == nil || *options.Visible {
		style |= w32.WS_VISIBLE
	}

	// Convert DIP coordinates to physical pixels
	bounds := DipToPhysicalRect(Rect{
		X:      options.X,
		Y:      options.Y,
		Width:  options.Width,
		Height: options.Height,
	})

	// Create the child window
	p.hwnd = w32.CreateWindowEx(
		0,
		w32.MustStringToUTF16Ptr("STATIC"), // Using STATIC class for the container
		nil,
		style,
		bounds.X,
		bounds.Y,
		bounds.Width,
		bounds.Height,
		p.parent.hwnd,
		0,
		w32.GetModuleHandle(""),
		nil,
	)

	if p.hwnd == 0 {
		globalApplication.error("failed to create panel child window")
		return
	}

	// Setup WebView2 (Chromium)
	p.setupChromium()
}

func (p *windowsPanelImpl) setupChromium() {
	p.chromium = edge.NewChromium()

	if globalApplication.options.ErrorHandler != nil {
		p.chromium.SetErrorCallback(globalApplication.options.ErrorHandler)
	}

	// Configure chromium
	p.chromium.DataPath = globalApplication.options.Windows.WebviewUserDataPath
	p.chromium.BrowserPath = globalApplication.options.Windows.WebviewBrowserPath

	// Set up callbacks
	p.chromium.MessageCallback = p.processMessage
	p.chromium.NavigationCompletedCallback = p.navigationCompletedCallback

	// Embed the WebView2 into our child window
	p.chromium.Embed(p.hwnd)
	p.chromium.Resize()

	// Configure settings
	settings, err := p.chromium.GetSettings()
	if err != nil {
		globalApplication.error("failed to get chromium settings: %v", err)
		return
	}

	debugMode := globalApplication.isDebugMode

	// Disable context menus unless in debug mode or explicitly enabled
	devToolsEnabled := debugMode
	if p.panel.options.DevToolsEnabled != nil {
		devToolsEnabled = *p.panel.options.DevToolsEnabled
	}
	err = settings.PutAreDefaultContextMenusEnabled(devToolsEnabled)
	if err != nil {
		globalApplication.error("failed to configure context menus: %v", err)
	}

	err = settings.PutAreDevToolsEnabled(devToolsEnabled)
	if err != nil {
		globalApplication.error("failed to configure devtools: %v", err)
	}

	// Set zoom if specified
	if p.panel.options.Zoom > 0 && p.panel.options.Zoom != 1.0 {
		p.chromium.PutZoomFactor(p.panel.options.Zoom)
	}

	// Set background colour
	if p.panel.options.Transparent {
		p.chromium.SetBackgroundColour(0, 0, 0, 0)
	} else {
		p.chromium.SetBackgroundColour(
			p.panel.options.BackgroundColour.Red,
			p.panel.options.BackgroundColour.Green,
			p.panel.options.BackgroundColour.Blue,
			p.panel.options.BackgroundColour.Alpha,
		)
	}

	// Navigate to initial content
	if p.panel.options.HTML != "" {
		p.loadHTMLWithScripts()
	} else if p.panel.options.URL != "" {
		startURL, err := assetserver.GetStartURL(p.panel.options.URL)
		if err != nil {
			globalApplication.error("failed to get start URL: %v", err)
			return
		}
		p.chromium.Navigate(startURL)
	}

	// Open inspector if requested
	if debugMode && p.panel.options.OpenInspectorOnStartup {
		p.chromium.OpenDevToolsWindow()
	}
}

func (p *windowsPanelImpl) loadHTMLWithScripts() {
	var script string
	if p.panel.options.JS != "" {
		script = p.panel.options.JS
	}
	if p.panel.options.CSS != "" {
		// Escape CSS for safe injection into JavaScript string
		escapedCSS := strings.ReplaceAll(p.panel.options.CSS, `\`, `\\`)
		escapedCSS = strings.ReplaceAll(escapedCSS, `"`, `\"`)
		escapedCSS = strings.ReplaceAll(escapedCSS, "\n", `\n`)
		escapedCSS = strings.ReplaceAll(escapedCSS, "\r", `\r`)
		script += fmt.Sprintf(
			"; addEventListener(\"DOMContentLoaded\", (event) => { document.head.appendChild(document.createElement('style')).innerHTML=\"%s\"; });",
			escapedCSS,
		)
	}
	if script != "" {
		p.chromium.Init(script)
	}
	p.chromium.NavigateToString(p.panel.options.HTML)
}

func (p *windowsPanelImpl) processMessage(message string, _ *edge.ICoreWebView2, _ *edge.ICoreWebView2WebMessageReceivedEventArgs) {
	// For now, just log panel messages
	// In future, we could route these to the parent window or handle panel-specific messages
	globalApplication.debug("Panel message received", "panel", p.panel.name, "message", message)
}

func (p *windowsPanelImpl) navigationCompletedCallback(_ *edge.ICoreWebView2, _ *edge.ICoreWebView2NavigationCompletedEventArgs) {
	p.navigationCompleted = true

	// Execute any pending JS
	if p.panel.options.JS != "" && p.panel.options.HTML == "" {
		p.execJS(p.panel.options.JS)
	}
	if p.panel.options.CSS != "" && p.panel.options.HTML == "" {
		// Escape CSS for safe injection into JavaScript string
		escapedCSS := strings.ReplaceAll(p.panel.options.CSS, `\`, `\\`)
		escapedCSS = strings.ReplaceAll(escapedCSS, `'`, `\'`)
		escapedCSS = strings.ReplaceAll(escapedCSS, "\n", `\n`)
		escapedCSS = strings.ReplaceAll(escapedCSS, "\r", `\r`)
		js := fmt.Sprintf(
			"(function() { var style = document.createElement('style'); style.appendChild(document.createTextNode('%s')); document.head.appendChild(style); })();",
			escapedCSS,
		)
		p.execJS(js)
	}

	// Mark runtime as loaded
	p.panel.markRuntimeLoaded()
}

func (p *windowsPanelImpl) destroy() {
	if p.chromium != nil {
		p.chromium.ShuttingDown()
	}
	if p.hwnd != 0 {
		w32.DestroyWindow(p.hwnd)
		p.hwnd = 0
	}
	p.chromium = nil
}

func (p *windowsPanelImpl) setBounds(bounds Rect) {
	if p.hwnd == 0 {
		return
	}

	// Convert DIP to physical pixels
	physicalBounds := DipToPhysicalRect(bounds)

	// Move and resize the child window
	w32.SetWindowPos(
		p.hwnd,
		0,
		physicalBounds.X,
		physicalBounds.Y,
		physicalBounds.Width,
		physicalBounds.Height,
		w32.SWP_NOZORDER|w32.SWP_NOACTIVATE,
	)

	// Resize the WebView2 to fill the child window
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

	// Get parent window position to calculate relative position
	parentRect := w32.GetWindowRect(p.parent.hwnd)
	if parentRect == nil {
		return Rect{}
	}

	// Calculate position relative to parent's client area
	parentClientX, parentClientY := w32.ClientToScreen(p.parent.hwnd, 0, 0)

	physicalBounds := Rect{
		X:      int(rect.Left) - parentClientX,
		Y:      int(rect.Top) - parentClientY,
		Width:  int(rect.Right - rect.Left),
		Height: int(rect.Bottom - rect.Top),
	}

	return PhysicalToDipRect(physicalBounds)
}

func (p *windowsPanelImpl) setZIndex(zIndex int) {
	if p.hwnd == 0 {
		return
	}

	// Use SetWindowPos to change z-order.
	// Note: This is a binary implementation - panels are either on top (zIndex > 0)
	// or at the bottom (zIndex <= 0). Granular z-index ordering is not supported
	// on Windows because child windows share a z-order space and precise positioning
	// would require tracking all panels and re-ordering them relative to each other.
	var insertAfter uintptr
	if zIndex > 0 {
		insertAfter = w32.HWND_TOP
	} else {
		insertAfter = w32.HWND_BOTTOM
	}

	w32.SetWindowPos(
		p.hwnd,
		insertAfter,
		0, 0, 0, 0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOACTIVATE,
	)
}

func (p *windowsPanelImpl) setURL(url string) {
	if p.chromium == nil {
		return
	}
	startURL, err := assetserver.GetStartURL(url)
	if err != nil {
		globalApplication.error("failed to get start URL: %v", err)
		return
	}
	p.navigationCompleted = false
	p.chromium.Navigate(startURL)
}

func (p *windowsPanelImpl) setHTML(html string) {
	if p.chromium == nil {
		return
	}
	p.chromium.NavigateToString(html)
}

func (p *windowsPanelImpl) execJS(js string) {
	if p.chromium == nil {
		return
	}
	globalApplication.dispatchOnMainThread(func() {
		p.chromium.Eval(js)
	})
}

func (p *windowsPanelImpl) reload() {
	p.execJS("window.location.reload();")
}

func (p *windowsPanelImpl) forceReload() {
	// WebView2 doesn't have a cache-bypass reload, so just reload normally
	p.reload()
}

func (p *windowsPanelImpl) show() {
	if p.hwnd == 0 {
		return
	}
	w32.ShowWindow(p.hwnd, w32.SW_SHOW)
}

func (p *windowsPanelImpl) hide() {
	if p.hwnd == 0 {
		return
	}
	w32.ShowWindow(p.hwnd, w32.SW_HIDE)
}

func (p *windowsPanelImpl) isVisible() bool {
	if p.hwnd == 0 {
		return false
	}
	style := uint32(w32.GetWindowLong(p.hwnd, w32.GWL_STYLE))
	return style&w32.WS_VISIBLE != 0
}

func (p *windowsPanelImpl) setZoom(zoom float64) {
	if p.chromium == nil {
		return
	}
	p.chromium.PutZoomFactor(zoom)
}

func (p *windowsPanelImpl) getZoom() float64 {
	if p.chromium == nil {
		return 1.0
	}
	controller := p.chromium.GetController()
	if controller == nil {
		return 1.0
	}
	factor, err := controller.GetZoomFactor()
	if err != nil {
		return 1.0
	}
	return factor
}

func (p *windowsPanelImpl) openDevTools() {
	if p.chromium == nil {
		return
	}
	p.chromium.OpenDevToolsWindow()
}

func (p *windowsPanelImpl) focus() {
	if p.hwnd == 0 {
		return
	}
	w32.SetFocus(p.hwnd)
	if p.chromium != nil {
		p.chromium.Focus()
	}
}

func (p *windowsPanelImpl) isFocused() bool {
	if p.hwnd == 0 {
		return false
	}
	return w32.GetFocus() == p.hwnd
}
