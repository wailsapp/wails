//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.13 -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework WebKit
#include "webview_window_split_darwin.h"
#include <stdlib.h>
*/
import "C"

import (
	"sync/atomic"
	"unsafe"

	"github.com/wailsapp/wails/v3/internal/assetserver"
)

var splitPaneIDCounter uint64

func nextSplitPaneID() uint {
	return uint(atomic.AddUint64(&splitPaneIDCounter, 1))
}

func newSplitWindow(options MacSplitWindowOptions) (*MacSplitWindow, error) {
	panes := make([]*MacSplitPane, len(options.Panes))
	for i, po := range options.Panes {
		panes[i] = &MacSplitPane{
			id:       nextSplitPaneID(),
			name:     po.Name,
			behavior: po.Behavior,
			content:  po.Content,
		}
	}

	firstContent := options.Panes[0].Content.(MacWebviewContent)

	windowOptions := WebviewWindowOptions{
		Title:  options.Title,
		Width:  options.Width,
		Height: options.Height,
		// MacSplitWindow.Show() controls visibility explicitly; showing
		// before the split view is installed would flash a plain window.
		Hidden: true,
		URL:    firstContent.URL,
		Mac: MacWindow{
			TitleBar:           options.TitleBar,
			Appearance:         options.Appearance,
			Backdrop:           options.Backdrop,
			WindowLevel:        options.WindowLevel,
			CollectionBehavior: options.CollectionBehavior,
			WebviewPreferences: firstContent.WebviewPreferences,
		},
	}

	outer := NewWindow(windowOptions)
	outer.macSplitPending = &macSplitPendingConfig{
		vertical:     options.Orientation == SplitHorizontal,
		autosaveName: options.AutosaveName,
		paneOptions:  options.Panes,
		panes:        panes,
	}

	// Mirrors WindowManager.NewWithOptions: register the window and either
	// run it now or defer to app.Run(), same as every other window. Done by
	// hand (instead of calling NewWithOptions) so macSplitPending is set
	// before the window is reachable via a.windows/pendingRun -- otherwise
	// Run() could race installSplitPanes reading it.
	globalApplication.windowsLock.Lock()
	globalApplication.windows[outer.id] = outer
	globalApplication.windowsLock.Unlock()
	for _, hook := range globalApplication.windowCreatedCallbacks {
		hook(outer)
	}
	globalApplication.runOrDeferToAppRun(outer)

	return &MacSplitWindow{window: outer, panes: panes}, nil
}

// installSplitPanes runs on the main thread as part of run(), right after
// windowNew() has built the window's own NSWindow+WKWebView. It converts
// the window's content to an NSSplitViewController and adds each pane in
// order; pane 0 reuses the window's existing webview (see the .h comment on
// splitWindowAddPane for why), panes 1..N get freshly built ones.
func (w *macosWebviewWindow) installSplitPanes(pending *macSplitPendingConfig) {
	splitVC := C.splitWindowInstall(w.nsWindow, C.bool(pending.vertical))

	if pending.autosaveName != "" {
		nameC := C.CString(pending.autosaveName)
		C.splitWindowSetAutosaveName(splitVC, nameC)
		C.free(unsafe.Pointer(nameC))
	}

	for i, po := range pending.paneOptions {
		pane := pending.panes[i]
		content := po.Content.(MacWebviewContent) // validated in NewSplitWindow

		// Pane 0 reuses the window's own webview, whose content is already
		// loaded by the normal run() flow via a correctly-resolved
		// assetserver.GetStartURL call; skip re-resolving/reloading it here.
		// Panes 1..N need the same wails://-scheme resolution applied by
		// hand, or a bare path like "/content.html" has no scheme for
		// WKWebView to route to the custom scheme handler and silently
		// fails to load (confirmed empirically -- no error, no scheme
		// handler callback, nothing).
		resolvedURL := content.URL
		if i != 0 {
			if resolved, err := assetserver.GetStartURL(content.URL); err == nil {
				resolvedURL = resolved
			} else {
				globalApplication.error("MacSplitWindow pane %q: resolving URL %q: %s", pane.name, content.URL, err)
			}
		}

		urlC := C.CString(resolvedURL)
		item := C.splitWindowAddPane(splitVC, w.nsWindow, C.bool(i == 0),
			urlC, C.int(po.Behavior),
			C.double(po.MinThickness), C.double(po.MaxThickness),
			C.double(po.PreferredThicknessFraction), C.double(po.HoldingPriority),
			C.bool(po.Collapsible), C.bool(po.StartCollapsed))
		C.free(unsafe.Pointer(urlC))

		pane.handleLock.Lock()
		pane.nativeItem = item
		pane.nativeWebview = C.splitPaneWebview(item)
		pane.handleLock.Unlock()
	}
}

func macSplitPaneIsCollapsed(item unsafe.Pointer) bool {
	return bool(C.splitPaneIsCollapsed(item))
}

func macSplitPaneSetCollapsed(item unsafe.Pointer, collapsed bool) {
	C.splitPaneSetCollapsed(item, C.bool(collapsed))
}

func macSplitPaneWebviewSetURL(webview unsafe.Pointer, url string) {
	resolvedURL := url
	if resolved, err := assetserver.GetStartURL(url); err == nil {
		resolvedURL = resolved
	} else {
		globalApplication.error("MacWebviewPane.SetURL: resolving URL %q: %s", url, err)
	}
	urlC := C.CString(resolvedURL)
	defer C.free(unsafe.Pointer(urlC))
	C.splitPaneWebviewSetURL(webview, urlC)
}

func macSplitPaneWebviewExecJS(webview unsafe.Pointer, js string) {
	jsC := C.CString(js)
	defer C.free(unsafe.Pointer(jsC))
	C.splitPaneWebviewExecJS(webview, jsC)
}

func macSplitPaneWebviewReload(webview unsafe.Pointer) {
	C.splitPaneWebviewReload(webview)
}
