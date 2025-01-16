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
		SetAccelerator("CmdOrCtrl+h").
		SetRole(Hide)
}

func NewHideOthersMenuItem() *MenuItem {
	return NewMenuItem("Hide Others").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+h").
		SetRole(HideOthers)
}

func NewFrontMenuItem() *MenuItem {
	return NewMenuItem("Bring All to Front")
}

func NewUnhideMenuItem() *MenuItem {
	return NewMenuItem("Show All")
}

func NewUndoMenuItem() *MenuItem {
	result := NewMenuItem("Undo").
		SetAccelerator("CmdOrCtrl+z")
	if runtime.GOOS != "darwin" {
		result.OnClick(func(ctx *Context) {
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
		SetAccelerator("CmdOrCtrl+Shift+z")
	if runtime.GOOS != "darwin" {
		result.OnClick(func(ctx *Context) {
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
		SetAccelerator("CmdOrCtrl+x")

	if runtime.GOOS != "darwin" {
		result.OnClick(func(ctx *Context) {
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
		SetAccelerator("CmdOrCtrl+c")

	if runtime.GOOS != "darwin" {
		result.OnClick(func(ctx *Context) {
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
		SetAccelerator("CmdOrCtrl+v")

	if runtime.GOOS != "darwin" {
		result.OnClick(func(ctx *Context) {
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
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+v")
}

func NewDeleteMenuItem() *MenuItem {
	result := NewMenuItem("Delete").
		SetAccelerator("backspace")

	if runtime.GOOS != "darwin" {
		result.OnClick(func(ctx *Context) {
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
		SetAccelerator("CmdOrCtrl+q").
		OnClick(func(ctx *Context) {
			globalApplication.Quit()
		})
}

func NewSelectAllMenuItem() *MenuItem {
	result := NewMenuItem("Select All").
		SetAccelerator("CmdOrCtrl+a")

	if runtime.GOOS != "darwin" {
		result.OnClick(func(ctx *Context) {
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
		OnClick(func(ctx *Context) {
			globalApplication.ShowAboutDialog()
		})
}

func NewCloseMenuItem() *MenuItem {
	return NewMenuItem("Close").
		SetAccelerator("CmdOrCtrl+w").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Close()
			}
		})
}

func NewReloadMenuItem() *MenuItem {
	return NewMenuItem("Reload").
		SetAccelerator("CmdOrCtrl+r").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Reload()
			}
		})
}

func NewForceReloadMenuItem() *MenuItem {
	return NewMenuItem("Force Reload").
		SetAccelerator("CmdOrCtrl+Shift+r").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ForceReload()
			}
		})
}

func NewToggleFullscreenMenuItem() *MenuItem {
	result := NewMenuItem("Toggle Full Screen").
		SetAccelerator("Ctrl+Command+F").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ToggleFullscreen()
			}
		})
	if runtime.GOOS != "darwin" {
		result.SetAccelerator("F11")
	}
	return result
}

func NewZoomResetMenuItem() *MenuItem {
	// reset zoom menu item
	return NewMenuItem("Actual Size").
		SetAccelerator("CmdOrCtrl+0").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ZoomReset()
			}
		})
}

func NewZoomInMenuItem() *MenuItem {
	return NewMenuItem("Zoom In").
		SetAccelerator("CmdOrCtrl+plus").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ZoomIn()
			}
		})
}

func NewZoomOutMenuItem() *MenuItem {
	return NewMenuItem("Zoom Out").
		SetAccelerator("CmdOrCtrl+-").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.ZoomOut()
			}
		})
}

func NewMinimiseMenuItem() *MenuItem {
	return NewMenuItem("Minimize").
		SetAccelerator("CmdOrCtrl+M").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Minimise()
			}
		})
}

func NewZoomMenuItem() *MenuItem {
	return NewMenuItem("Zoom").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Zoom()
			}
		})
}

func NewFullScreenMenuItem() *MenuItem {
	return NewMenuItem("Fullscreen").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.CurrentWindow()
			if currentWindow != nil {
				currentWindow.Fullscreen()
			}
		})
}

func NewPrintMenuItem() *MenuItem {
	return NewMenuItem("Print").
		SetAccelerator("CmdOrCtrl+p")
}

func NewPageLayoutMenuItem() *MenuItem {
	return NewMenuItem("Page Setup...").
		SetAccelerator("CmdOrCtrl+Shift+p")
}

func NewShowAllMenuItem() *MenuItem {
	return NewMenuItem("Show All")
}

func NewBringAllToFrontMenuItem() *MenuItem {
	return NewMenuItem("Bring All to Front")
}

func NewNewFileMenuItem() *MenuItem {
	return NewMenuItem("New File").
		SetAccelerator("CmdOrCtrl+n")
}

func NewOpenMenuItem() *MenuItem {
	return NewMenuItem("Open...").
		SetAccelerator("CmdOrCtrl+o").
		SetRole(Open)
}

func NewSaveMenuItem() *MenuItem {
	return NewMenuItem("Save").
		SetAccelerator("CmdOrCtrl+s")
}

func NewSaveAsMenuItem() *MenuItem {
	return NewMenuItem("Save As...").
		SetAccelerator("CmdOrCtrl+Shift+s")
}

func NewStartSpeakingMenuItem() *MenuItem {
	return NewMenuItem("Start Speaking").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+.")
}

func NewStopSpeakingMenuItem() *MenuItem {
	return NewMenuItem("Stop Speaking").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+,")
}

func NewRevertMenuItem() *MenuItem {
	return NewMenuItem("Revert").
		SetAccelerator("CmdOrCtrl+r")
}

func NewFindMenuItem() *MenuItem {
	return NewMenuItem("Find...").
		SetAccelerator("CmdOrCtrl+f")
}

func NewFindAndReplaceMenuItem() *MenuItem {
	return NewMenuItem("Find and Replace...").
		SetAccelerator("CmdOrCtrl+Shift+f")
}

func NewFindNextMenuItem() *MenuItem {
	return NewMenuItem("Find Next").
		SetAccelerator("CmdOrCtrl+g")
}

func NewFindPreviousMenuItem() *MenuItem {
	return NewMenuItem("Find Previous").
		SetAccelerator("CmdOrCtrl+Shift+g")
}

func NewHelpMenuItem() *MenuItem {
	return NewMenuItem("Help").
		SetAccelerator("CmdOrCtrl+?")
}
