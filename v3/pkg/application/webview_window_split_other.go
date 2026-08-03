//go:build !darwin || ios || server

package application

import (
	"fmt"
	"unsafe"
)

func newSplitWindow(options MacSplitWindowOptions) (*MacSplitWindow, error) {
	return nil, fmt.Errorf("MacSplitWindow is only supported on macOS")
}

func macSplitPaneIsCollapsed(item unsafe.Pointer) bool {
	return false
}

func macSplitPaneSetCollapsed(item unsafe.Pointer, collapsed bool) {
}

func macSplitPaneWebviewSetURL(webview unsafe.Pointer, url string) {
}

func macSplitPaneWebviewExecJS(webview unsafe.Pointer, js string) {
}

func macSplitPaneWebviewReload(webview unsafe.Pointer) {
}
