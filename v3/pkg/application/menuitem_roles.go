package application

import "runtime"

func NewSpeechMenu() *MenuItem {
	speechMenu := NewMenu()
	speechMenu.AddRole(StartSpeaking)
	speechMenu.AddRole(StopSpeaking)
	subMenu := NewSubMenuItem("Speech")
	subMenu.submenu = speechMenu
	return subMenu
}

func NewHideMenuItem() *MenuItem {
	return NewMenuItem("Hide " + globalApplication.options.Name).
		SetAcceleratorItem("CmdOrCtrl+h").
		SetRoleItem(Hide)
}

func NewHideOthersMenuItem() *MenuItem {
	return NewMenuItem("Hide Others").
		SetAcceleratorItem("CmdOrCtrl+OptionOrAlt+h").
		SetRoleItem(HideOthers)
}

func NewFrontMenuItem() *MenuItem {
	return NewMenuItem("Bring All to Front")
}

func NewUnhideMenuItem() *MenuItem {
	return NewMenuItem("Show All")
}

func NewUndoMenuItem() *MenuItem {
	result := NewMenuItem("Undo").
		SetAcceleratorItem("CmdOrCtrl+z")
	if runtime.GOOS != "darwin" {
		result.OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.undo()
			}
		})
	}
	return result
}

// NewRedoMenuItem creates a new menu item for redoing the last action
func NewRedoMenuItem() *MenuItem {
	result := NewMenuItem("Redo").
		SetAcceleratorItem("CmdOrCtrl+Shift+z")
	if runtime.GOOS != "darwin" {
		result.OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.redo()
			}
		})
	}
	return result
}

func NewCutMenuItem() *MenuItem {
	result := NewMenuItem("Cut").
		SetAcceleratorItem("CmdOrCtrl+x")

	if runtime.GOOS != "darwin" {
		result.OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.cut()
			}
		})
	}
	return result
}

func NewCopyMenuItem() *MenuItem {
	result := NewMenuItem("Copy").
		SetAcceleratorItem("CmdOrCtrl+c")

	if runtime.GOOS != "darwin" {
		result.OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.copy()
			}
		})
	}
	return result
}

func NewPasteMenuItem() *MenuItem {
	result := NewMenuItem("Paste").
		SetAcceleratorItem("CmdOrCtrl+v")

	if runtime.GOOS != "darwin" {
		result.OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.paste()
			}
		})
	}
	return result
}

func NewPasteAndMatchStyleMenuItem() *MenuItem {
	return NewMenuItem("Paste and Match Style").
		SetAcceleratorItem("CmdOrCtrl+OptionOrAlt+Shift+v")
}

func NewDeleteMenuItem() *MenuItem {
	result := NewMenuItem("Delete").
		SetAcceleratorItem("backspace")

	if runtime.GOOS != "darwin" {
		result.OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.delete()
			}
		})
	}
	return result
}

func NewQuitMenuItem() *MenuItem {
	label := "Quit"
	if runtime.GOOS == "darwin" {
		if globalApplication.options.Name != "" {
			label += " " + globalApplication.options.Name
		}
	}
	return NewMenuItem(label).
		SetAcceleratorItem("CmdOrCtrl+q").
		OnClickItem(func(ctx *Context) {
			globalApplication.Quit()
		})
}

func NewSelectAllMenuItem() *MenuItem {
	result := NewMenuItem("Select All").
		SetAcceleratorItem("CmdOrCtrl+a")

	if runtime.GOOS != "darwin" {
		result.OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.selectAll()
			}
		})
	}
	return result
}

func NewAboutMenuItem() *MenuItem {
	label := "About"
	if globalApplication.options.Name != "" {
		label += " " + globalApplication.options.Name
	}
	return NewMenuItem(label).
		OnClickItem(func(ctx *Context) {
			globalApplication.ShowAboutDialog()
		})
}

func NewCloseMenuItem() *MenuItem {
	return NewMenuItem("Close").
		SetAcceleratorItem("CmdOrCtrl+w").
		OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Close()
			}
		})
}

func NewReloadMenuItem() *MenuItem {
	return NewMenuItem("Reload").
		SetAcceleratorItem("CmdOrCtrl+r").
		OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Reload()
			}
		})
}

func NewForceReloadMenuItem() *MenuItem {
	return NewMenuItem("Force Reload").
		SetAcceleratorItem("CmdOrCtrl+Shift+r").
		OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ForceReload()
			}
		})
}

func NewToggleFullscreenMenuItem() *MenuItem {
	result := NewMenuItem("Toggle Full Screen").
		SetAcceleratorItem("Ctrl+Command+F").
		OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ToggleFullscreen()
			}
		})
	if runtime.GOOS != "darwin" {
		result.SetAcceleratorItem("F11")
	}
	return result
}

func NewZoomResetMenuItem() *MenuItem {
	// reset zoom menu item
	return NewMenuItem("Actual Size").
		SetAcceleratorItem("CmdOrCtrl+0").
		OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ZoomReset()
			}
		})
}

func NewZoomInMenuItem() *MenuItem {
	return NewMenuItem("Zoom In").
		SetAcceleratorItem("CmdOrCtrl+plus").
		OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ZoomIn()
			}
		})
}

func NewZoomOutMenuItem() *MenuItem {
	return NewMenuItem("Zoom Out").
		SetAcceleratorItem("CmdOrCtrl+-").
		OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ZoomOut()
			}
		})
}

func NewMinimiseMenuItem() *MenuItem {
	return NewMenuItem("Minimize").
		SetAcceleratorItem("CmdOrCtrl+M").
		OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Minimise()
			}
		})
}

func NewZoomMenuItem() *MenuItem {
	return NewMenuItem("Zoom").
		OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Zoom()
			}
		})
}

func NewFullScreenMenuItem() *MenuItem {
	return NewMenuItem("Fullscreen").
		OnClickItem(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Fullscreen()
			}
		})
}

func NewPrintMenuItem() *MenuItem {
	return NewMenuItem("Print").
		SetAcceleratorItem("CmdOrCtrl+p")
}

func NewPageLayoutMenuItem() *MenuItem {
	return NewMenuItem("Page Setup...").
		SetAcceleratorItem("CmdOrCtrl+Shift+p")
}

func NewShowAllMenuItem() *MenuItem {
	return NewMenuItem("Show All")
}

func NewBringAllToFrontMenuItem() *MenuItem {
	return NewMenuItem("Bring All to Front")
}

func NewNewFileMenuItem() *MenuItem {
	return NewMenuItem("New File").
		SetAcceleratorItem("CmdOrCtrl+n")
}

func NewOpenMenuItem() *MenuItem {
	return NewMenuItem("Open...").
		SetAcceleratorItem("CmdOrCtrl+o").
		SetRoleItem(Open)
}

func NewSaveMenuItem() *MenuItem {
	return NewMenuItem("Save").
		SetAcceleratorItem("CmdOrCtrl+s")
}

func NewSaveAsMenuItem() *MenuItem {
	return NewMenuItem("Save As...").
		SetAcceleratorItem("CmdOrCtrl+Shift+s")
}

func NewStartSpeakingMenuItem() *MenuItem {
	return NewMenuItem("Start Speaking").
		SetAcceleratorItem("CmdOrCtrl+OptionOrAlt+Shift+.")
}

func NewStopSpeakingMenuItem() *MenuItem {
	return NewMenuItem("Stop Speaking").
		SetAcceleratorItem("CmdOrCtrl+OptionOrAlt+Shift+,")
}

func NewRevertMenuItem() *MenuItem {
	return NewMenuItem("Revert").
		SetAcceleratorItem("CmdOrCtrl+r")
}

func NewFindMenuItem() *MenuItem {
	return NewMenuItem("Find...").
		SetAcceleratorItem("CmdOrCtrl+f")
}

func NewFindAndReplaceMenuItem() *MenuItem {
	return NewMenuItem("Find and Replace...").
		SetAcceleratorItem("CmdOrCtrl+Shift+f")
}

func NewFindNextMenuItem() *MenuItem {
	return NewMenuItem("Find Next").
		SetAcceleratorItem("CmdOrCtrl+g")
}

func NewFindPreviousMenuItem() *MenuItem {
	return NewMenuItem("Find Previous").
		SetAcceleratorItem("CmdOrCtrl+Shift+g")
}

func NewHelpMenuItem() *MenuItem {
	return NewMenuItem("Help").
		SetAcceleratorItem("CmdOrCtrl+?")
}
