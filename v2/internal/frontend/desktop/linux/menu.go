//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include "gtk/gtk.h"

static GtkMenuItem *toGtkMenuItem(void *pointer) { return (GTK_MENU_ITEM(pointer)); }
static GtkMenuShell *toGtkMenuShell(void *pointer) { return (GTK_MENU_SHELL(pointer)); }

*/
import "C"
import "github.com/wailsapp/wails/v2/pkg/menu"
import "unsafe"

func (f *Frontend) MenuSetApplicationMenu(menu *menu.Menu) {
	f.mainWindow.SetApplicationMenu(menu)
}

func (f *Frontend) MenuUpdateApplicationMenu() {
	//processMenu(f.mainWindow, f.mainWindow.applicationMenu)
}

func (w *Window) SetApplicationMenu(menu *menu.Menu) {
	if menu == nil {
		return
	}

	// TODO: Remove existing menu if exists

	// Increase ref count?
	w.menubar = C.gtk_menu_bar_new()

	processMenu(w, menu)

	C.gtk_widget_show(w.menubar)
}

func processMenu(window *Window, menu *menu.Menu) {
	for _, menuItem := range menu.Items {
		gtkMenu := C.gtk_menu_new()
		cLabel := C.CString(menuItem.Label)
		submenu := C.gtk_menu_item_new_with_label(cLabel)
		for _, menuItem := range menuItem.SubMenu.Items {
			processMenuItem(gtkMenu, menuItem)
		}
		C.gtk_menu_item_set_submenu(C.toGtkMenuItem(unsafe.Pointer(submenu)), gtkMenu)
		C.gtk_menu_shell_append(C.toGtkMenuShell(unsafe.Pointer(window.menubar)), submenu)
	}
}

func processMenuItem(parent *C.GtkWidget, menuItem *menu.MenuItem) {
	if menuItem.Hidden {
		return
	}
	switch menuItem.Type {
	//case menu.SeparatorType:
	//	parent.AddSeparator()
	case menu.TextType:
		//shortcut := acceleratorToWincShortcut(menuItem.Accelerator)
		cLabel := C.CString(menuItem.Label)
		textMenuItem := C.gtk_menu_item_new_with_label(cLabel)
		C.gtk_menu_shell_append(C.toGtkMenuShell(unsafe.Pointer(parent)), textMenuItem)
		C.gtk_widget_show(textMenuItem)
		//newItem := parent.AddItem(menuItem.Label, shortcut)

		//if menuItem.Tooltip != "" {
		//	newItem.SetToolTip(menuItem.Tooltip)
		//}

		//if menuItem.Click != nil {
		//	newItem.OnClick().Bind(func(e *winc.Event) {
		//		menuItem.Click(&menu.CallbackData{
		//			MenuItem: menuItem,
		//		})
		//	})
		//}

		//newItem.SetEnabled(!menuItem.Disabled)

		//case menu.CheckboxType:
		//	shortcut := acceleratorToWincShortcut(menuItem.Accelerator)
		//	newItem := parent.AddItem(menuItem.Label, shortcut)
		//	newItem.SetCheckable(true)
		//	newItem.SetChecked(menuItem.Checked)
		//	//if menuItem.Tooltip != "" {
		//	//	newItem.SetToolTip(menuItem.Tooltip)
		//	//}
		//	if menuItem.Click != nil {
		//		newItem.OnClick().Bind(func(e *winc.Event) {
		//			toggleCheckBox(menuItem)
		//			menuItem.Click(&menu.CallbackData{
		//				MenuItem: menuItem,
		//			})
		//		})
		//	}
		//	newItem.SetEnabled(!menuItem.Disabled)
		//	addCheckBoxToMap(menuItem, newItem)
		//case menu.RadioType:
		//	shortcut := acceleratorToWincShortcut(menuItem.Accelerator)
		//	newItem := parent.AddItemRadio(menuItem.Label, shortcut)
		//	newItem.SetCheckable(true)
		//	newItem.SetChecked(menuItem.Checked)
		//	//if menuItem.Tooltip != "" {
		//	//	newItem.SetToolTip(menuItem.Tooltip)
		//	//}
		//	if menuItem.Click != nil {
		//		newItem.OnClick().Bind(func(e *winc.Event) {
		//			toggleRadioItem(menuItem)
		//			menuItem.Click(&menu.CallbackData{
		//				MenuItem: menuItem,
		//			})
		//		})
		//	}
		//	newItem.SetEnabled(!menuItem.Disabled)
		//	addRadioItemToMap(menuItem, newItem)
		//case menu.SubmenuType:
		//	submenu := parent.AddSubMenu(menuItem.Label)
		//	for _, menuItem := range menuItem.SubMenu.Items {
		//		processMenuItem(submenu, menuItem)
		//	}
	}
}
