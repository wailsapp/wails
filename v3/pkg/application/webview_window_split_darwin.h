//go:build darwin && !ios && !server

#ifndef WailsSplitWindow_h
#define WailsSplitWindow_h

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

// Creates an NSSplitViewController and installs it as nsWindow's
// contentViewController. Panes must then be added, in order, via
// splitWindowAddPane -- the window's pre-existing WKWebView (window.webView,
// built by windowNew) is untouched until the first splitWindowAddPane call
// that passes reuseExistingWebview=true.
void* splitWindowInstall(void* nsWindow, bool vertical);

void splitWindowSetAutosaveName(void* splitViewController, const char* name);

// Adds one pane, in order. If reuseExistingWebview is true, the pane hosts
// nsWindow's existing WKWebView (window.webView) instead of creating a new
// one -- exactly one pane per split window should pass true (normally pane
// 0), so the window's pre-existing navigation/runtime IPC wiring keeps
// targeting a real on-screen webview instead of one that's been silently
// orphaned. New pane webviews reuse the window's own WebviewWindowDelegate
// (same windowId) for the wails:// scheme handler and the "external"
// script message handler, so asset loading and Call/Events/Clipboard IPC
// from a pane route exactly like they do from the window's main webview --
// there is no separate per-pane identity. One side effect: a pane finishing
// navigation also fires the window's WebViewDidFinishNavigation event, same
// as the main webview finishing.
// Returns the pane's NSSplitViewItem*.
void* splitWindowAddPane(void* splitViewController, void* nsWindow, bool reuseExistingWebview,
    const char* url, int behavior,
    double minThickness, double maxThickness, double preferredThicknessFraction, double holdingPriority,
    bool collapsible, bool startCollapsed);

// Returns the WKWebView* hosted by a pane's NSSplitViewItem.
void* splitPaneWebview(void* splitViewItem);

void splitPaneSetCollapsed(void* splitViewItem, bool collapsed);
bool splitPaneIsCollapsed(void* splitViewItem);

void splitPaneWebviewSetURL(void* webView, const char* url);
void splitPaneWebviewExecJS(void* webView, const char* js);
void splitPaneWebviewReload(void* webView);

#endif /* WailsSplitWindow_h */
