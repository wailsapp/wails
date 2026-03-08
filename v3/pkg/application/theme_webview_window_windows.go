//go:build windows

package application

func (w *windowsWebviewWindow) setTheme(theme WinTheme) {
}

func (w *windowsWebviewWindow) getTheme() WinTheme {
	return WinThemeApplication
}
