//go:build !production || devtools

package application

func newShowDevToolsMenuItem() *MenuItem {
	return newMenuItem("Show Developer Tools").
		SetAccelerator("Alt+Command+I").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ToggleDevTools()
			}
		})
}
