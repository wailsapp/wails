//go:build !production || devtools

package application

func NewOpenDevToolsMenuItem() *MenuItem {
	return NewMenuItem("Open Developer Tools").
		SetAccelerator("Alt+Command+I").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.Window.Current()
			if currentWindow != nil {
				currentWindow.OpenDevTools()
			}
		})
}
