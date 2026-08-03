//go:build !darwin || ios || server

package application

import "unsafe"

// AddTitlebarAccessory is macOS only; on every other platform the returned
// accessory's native handles stay nil forever, so SetURL/ExecJS calls
// warn-and-no-op rather than doing anything.
func addTitlebarAccessoryNow(w *WebviewWindow, options MacTitlebarAccessoryOptions, accessory *MacTitlebarAccessory) error {
	return nil
}

func removeTitlebarAccessoryNow(w *WebviewWindow, controller unsafe.Pointer) {
}

func macTitlebarAccessoryWebviewSetURL(webview unsafe.Pointer, url string) {
}

func macTitlebarAccessoryWebviewExecJS(webview unsafe.Pointer, js string) {
}
