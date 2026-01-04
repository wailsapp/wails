//go:build linux && !android && !gtk3

package application

import (
	"fmt"
	"runtime"
)

type linuxMenuItem struct {
	menuItem   *MenuItem
	native     pointer
	handlerId  uint
	parentMenu pointer
	menuIndex  int
	isHidden   bool
}

func (l linuxMenuItem) setTooltip(tooltip string) {
	InvokeSync(func() {
		menuItemSetToolTip(l.native, tooltip)
	})
}

func (l linuxMenuItem) destroy() {
	InvokeSync(func() {
		menuItemDestroy(l.native)
	})
}

func (l linuxMenuItem) blockSignal() {
}

func (l linuxMenuItem) setBitmap(data []byte) {
	InvokeSync(func() {
		menuItemSetBitmap(l.native, data)
	})
}

func (l linuxMenuItem) unBlockSignal() {
}

func (l linuxMenuItem) setLabel(s string) {
	InvokeSync(func() {
		menuItemSetLabel(l.native, s)
	})
}

func (l linuxMenuItem) isChecked() bool {
	return menuItemChecked(l.native)
}

func (l linuxMenuItem) setDisabled(disabled bool) {
	InvokeSync(func() {
		menuItemSetDisabled(l.native, disabled)
	})
}

func (l linuxMenuItem) setChecked(checked bool) {
	InvokeSync(func() {
		menuItemSetChecked(l.native, checked)
	})
}

func (l *linuxMenuItem) setHidden(hidden bool) {
	if l.isHidden == hidden {
		return
	}
	InvokeSync(func() {
		menuItemSetHidden(l, hidden)
	})
	l.isHidden = hidden
}

func (l linuxMenuItem) setAccelerator(accelerator *accelerator) {
	if accelerator == nil || l.menuItem == nil {
		return
	}
	InvokeSync(func() {
		setMenuItemAccelerator(l.menuItem.id, accelerator)
	})
}

func newMenuItemImpl(item *MenuItem) *linuxMenuItem {
	result := &linuxMenuItem{
		menuItem: item,
	}
	switch item.itemType {
	case text:
		result.native = menuItemNewWithId(item.label, item.bitmap, item.id)
	case submenu:
		result.native = menuItemNewWithId(item.label, item.bitmap, item.id)
	default:
		panic(fmt.Sprintf("Unknown menu type for newMenuItemImpl: %v", item.itemType))
	}
	if item.accelerator != nil {
		result.setAccelerator(item.accelerator)
	}
	result.setDisabled(result.menuItem.disabled)
	return result
}

func newCheckMenuItemImpl(item *MenuItem) *linuxMenuItem {
	result := &linuxMenuItem{
		menuItem: item,
		native:   menuCheckItemNewWithId(item.label, item.bitmap, item.id, item.checked),
	}
	if item.accelerator != nil {
		result.setAccelerator(item.accelerator)
	}
	result.setDisabled(result.menuItem.disabled)
	return result
}

func newRadioMenuItemImpl(item *MenuItem, groupId uint, checkedId uint) *linuxMenuItem {
	result := &linuxMenuItem{
		menuItem: item,
		native:   menuRadioItemNewWithGroup(item.label, item.id, groupId, checkedId),
	}
	if item.accelerator != nil {
		result.setAccelerator(item.accelerator)
	}
	result.setDisabled(result.menuItem.disabled)
	return result
}

func newSpeechMenu() *MenuItem {
	speechMenu := NewMenu()
	speechMenu.Add("Start Speaking").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+.").
		OnClick(func(ctx *Context) {})
	speechMenu.Add("Stop Speaking").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+,").
		OnClick(func(ctx *Context) {})
	subMenu := NewSubMenuItem("Speech")
	subMenu.submenu = speechMenu
	return subMenu
}

func newFrontMenuItem() *MenuItem {
	panic("implement me")
}

func newHideMenuItem() *MenuItem {
	return NewMenuItem("Hide " + globalApplication.options.Name).
		SetAccelerator("CmdOrCtrl+h").
		OnClick(func(ctx *Context) {})
}

func newHideOthersMenuItem() *MenuItem {
	return NewMenuItem("Hide Others").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+h").
		OnClick(func(ctx *Context) {})
}

func newUnhideMenuItem() *MenuItem {
	return NewMenuItem("Show All").
		OnClick(func(ctx *Context) {})
}

func newUndoMenuItem() *MenuItem {
	return NewMenuItem("Undo").
		SetAccelerator("CmdOrCtrl+z").
		OnClick(func(ctx *Context) {})
}

func newRedoMenuItem() *MenuItem {
	return NewMenuItem("Redo").
		SetAccelerator("CmdOrCtrl+Shift+z").
		OnClick(func(ctx *Context) {})
}

func newCutMenuItem() *MenuItem {
	return NewMenuItem("Cut").
		SetAccelerator("CmdOrCtrl+x").
		OnClick(func(ctx *Context) {})
}

func newCopyMenuItem() *MenuItem {
	return NewMenuItem("Copy").
		SetAccelerator("CmdOrCtrl+c").
		OnClick(func(ctx *Context) {})
}

func newPasteMenuItem() *MenuItem {
	return NewMenuItem("Paste").
		SetAccelerator("CmdOrCtrl+v").
		OnClick(func(ctx *Context) {})
}

func newPasteAndMatchStyleMenuItem() *MenuItem {
	return NewMenuItem("Paste and Match Style").
		SetAccelerator("CmdOrCtrl+OptionOrAlt+Shift+v").
		OnClick(func(ctx *Context) {})
}

func newDeleteMenuItem() *MenuItem {
	return NewMenuItem("Delete").
		SetAccelerator("backspace").
		OnClick(func(ctx *Context) {})
}

func newQuitMenuItem() *MenuItem {
	return NewMenuItem("Quit " + globalApplication.options.Name).
		SetAccelerator("CmdOrCtrl+q").
		OnClick(func(ctx *Context) {
			globalApplication.Quit()
		})
}

func newSelectAllMenuItem() *MenuItem {
	return NewMenuItem("Select All").
		SetAccelerator("CmdOrCtrl+a").
		OnClick(func(ctx *Context) {})
}

func newAboutMenuItem() *MenuItem {
	return NewMenuItem("About " + globalApplication.options.Name).
		OnClick(func(ctx *Context) {
			globalApplication.Menu.ShowAbout()
		})
}

func newCloseMenuItem() *MenuItem {
	return NewMenuItem("Close").
		SetAccelerator("CmdOrCtrl+w").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.Window.Current()
			if currentWindow != nil {
				currentWindow.Close()
			}
		})
}

func newReloadMenuItem() *MenuItem {
	return NewMenuItem("Reload").
		SetAccelerator("CmdOrCtrl+r").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.Window.Current()
			if currentWindow != nil {
				currentWindow.Reload()
			}
		})
}

func newForceReloadMenuItem() *MenuItem {
	return NewMenuItem("Force Reload").
		SetAccelerator("CmdOrCtrl+Shift+r").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.Window.Current()
			if currentWindow != nil {
				currentWindow.ForceReload()
			}
		})
}

func newToggleFullscreenMenuItem() *MenuItem {
	result := NewMenuItem("Toggle Full Screen").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.Window.Current()
			if currentWindow != nil {
				currentWindow.ToggleFullscreen()
			}
		})
	if runtime.GOOS == "darwin" {
		result.SetAccelerator("Ctrl+Command+F")
	} else {
		result.SetAccelerator("F11")
	}
	return result
}

func newZoomResetMenuItem() *MenuItem {
	return NewMenuItem("Actual Size").
		SetAccelerator("CmdOrCtrl+0").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.Window.Current()
			if currentWindow != nil {
				currentWindow.ZoomReset()
			}
		})
}

func newZoomInMenuItem() *MenuItem {
	return NewMenuItem("Zoom In").
		SetAccelerator("CmdOrCtrl+plus").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.Window.Current()
			if currentWindow != nil {
				currentWindow.ZoomIn()
			}
		})
}

func newZoomOutMenuItem() *MenuItem {
	return NewMenuItem("Zoom Out").
		SetAccelerator("CmdOrCtrl+-").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.Window.Current()
			if currentWindow != nil {
				currentWindow.ZoomOut()
			}
		})
}

func newMinimizeMenuItem() *MenuItem {
	return NewMenuItem("Minimize").
		SetAccelerator("CmdOrCtrl+M").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.Window.Current()
			if currentWindow != nil {
				currentWindow.Minimise()
			}
		})
}

func newZoomMenuItem() *MenuItem {
	return NewMenuItem("Zoom").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.Window.Current()
			if currentWindow != nil {
				currentWindow.Zoom()
			}
		})
}

func newFullScreenMenuItem() *MenuItem {
	return NewMenuItem("Fullscreen").
		OnClick(func(ctx *Context) {
			currentWindow := globalApplication.Window.Current()
			if currentWindow != nil {
				currentWindow.Fullscreen()
			}
		})
}

func newPrintMenuItem() *MenuItem {
	panic("Implement me")
}

func newPageLayoutMenuItem() *MenuItem {
	panic("Implement me")
}

func newShowAllMenuItem() *MenuItem {
	panic("Implement me")
}

func newBringAllToFrontMenuItem() *MenuItem {
	panic("Implement me")
}

func newNewFileMenuItem() *MenuItem {
	panic("Implement me")
}

func newOpenMenuItem() *MenuItem {
	panic("Implement me")
}

func newSaveMenuItem() *MenuItem {
	panic("Implement me")
}

func newSaveAsMenuItem() *MenuItem {
	panic("Implement me")
}

func newStartSpeakingMenuItem() *MenuItem {
	panic("Implement me")
}

func newStopSpeakingMenuItem() *MenuItem {
	panic("Implement me")
}

func newRevertMenuItem() *MenuItem {
	panic("Implement me")
}

func newFindMenuItem() *MenuItem {
	panic("Implement me")
}

func newFindAndReplaceMenuItem() *MenuItem {
	panic("Implement me")
}

func newFindNextMenuItem() *MenuItem {
	panic("Implement me")
}

func newFindPreviousMenuItem() *MenuItem {
	panic("Implement me")
}

func newHelpMenuItem() *MenuItem {
	panic("Implement me")
}
