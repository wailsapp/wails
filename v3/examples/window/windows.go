//go:build windows

package main

import "github.com/wailsapp/wails/v3/pkg/w32"

func init() {
	getExStyle = func() int {
		return w32.WS_EX_TOOLWINDOW | w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOPMOST
	}
}
