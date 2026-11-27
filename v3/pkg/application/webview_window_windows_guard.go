//go:build windows

package application

import "github.com/wailsapp/wails/v3/internal/webview2/pkg/edge"

// NavigationGuard is consulted before every WebView2 navigation starts.
//
// target is the URI the WebView is about to navigate to. userInitiated reports
// whether the navigation was started by a user action (clicking a link, or
// window.open) rather than by the page's own load or a redirect.
//
// Returning true means the guard has handled the navigation itself, for example
// by opening the URI in the system browser; Wails then cancels the navigation
// and the WebView stays on the current document. Returning false lets the
// navigation proceed normally.
type NavigationGuard func(target string, userInitiated bool) bool

// SetNavigationGuard installs a guard consulted before every navigation.
//
// This runs on the WebView2 event thread, i.e. before any page script gets a
// chance to run, so it keeps working on pages loaded from a different origin,
// where the Wails runtime is not present. Pass nil to remove the guard.
func (w *WebviewWindow) SetNavigationGuard(guard NavigationGuard) {
	if w == nil || w.impl == nil || w.isDestroyed() {
		return
	}
	impl, ok := w.impl.(*windowsWebviewWindow)
	if !ok {
		return
	}
	impl.navigationGuard = guard
	if chromium := impl.chromium; chromium != nil {
		wireNavigationGuard(chromium, &impl.navigationGuard)
	}
}

// RunExternalJS executes js in the window's page, bypassing the runtimeReady
// gate that ExecJS observes.
//
// ExecJS queues scripts until the page posts "wails:runtime:ready", which only
// pages that actually loaded the Wails runtime do. A window that navigates to a
// different origin - a remote web app hosted elsewhere, a third-party login page
// - never posts it, so ExecJS silently never runs anything there. Use this
// method when the target page is not Wails-hosted and the script has to run
// regardless.
//
// The script runs on the UI thread and is fire-and-forget: failures are reported
// through the WebView2 log, not to the caller.
func (w *WebviewWindow) RunExternalJS(js string) {
	if w == nil || w.impl == nil || w.isDestroyed() {
		return
	}
	impl, ok := w.impl.(*windowsWebviewWindow)
	if !ok || impl.chromium == nil {
		return
	}
	chromium := impl.chromium
	globalApplication.dispatchOnMainThread(func() {
		chromium.Eval(js)
	})
}

// wireNavigationGuard connects guard to the WebView2 navigation callback.
//
// guard is passed by pointer so that a guard installed after setupChromium, or
// reinstalled after a WebView2 process recovery rebuilds the controller, is
// picked up without re-registering anything.
func wireNavigationGuard(chromium *edge.Chromium, guard *NavigationGuard) {
	if chromium == nil {
		return
	}
	chromium.NavigationStartingWithArgsCallback = func(_ *edge.ICoreWebView2, args *edge.ICoreWebView2NavigationStartingEventArgs) {
		if args == nil || guard == nil || *guard == nil {
			return
		}
		target, err := args.GetUri()
		if err != nil || target == "" {
			return
		}
		initiated, _ := args.GetIsUserInitiated()
		if (*guard)(target, initiated) {
			_ = args.PutCancel(true)
		}
	}
}
