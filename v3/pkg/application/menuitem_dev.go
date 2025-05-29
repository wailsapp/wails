//go:build !production || devtools

package application

func NewOpenDevToolsMenuItem() *MenuItem {
	return NewMenuItem("Open Developer Tools").
		SetAcceleratorItem("Alt+Command+I").
		OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.OpenDevTools()
			}
		})
}
