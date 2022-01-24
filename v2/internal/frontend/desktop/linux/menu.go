//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include "gtk/gtk.h"

static GtkMenuItem *toGtkMenuItem(void *pointer) { return (GTK_MENU_ITEM(pointer)); }
static GtkMenuShell *toGtkMenuShell(void *pointer) { return (GTK_MENU_SHELL(pointer)); }
static GtkCheckMenuItem *toGtkCheckMenuItem(void *pointer) { return (GTK_CHECK_MENU_ITEM(pointer)); }

extern void handleMenuItemClick(int);

void clickCallback(void *dummy, gpointer data) {
	handleMenuItemClick(GPOINTER_TO_INT(data));
}

void connectClick(GtkWidget* menuItem, int data) {
	g_signal_connect(menuItem, "activate", G_CALLBACK(clickCallback), GINT_TO_POINTER(data));
}
*/
import "C"
import "github.com/wailsapp/wails/v2/pkg/menu"
import "unsafe"

var menuIdCounter int
var menuItemToId map[*menu.MenuItem]int
var menuIdToItem map[int]*menu.MenuItem
var gtkCheckboxCache map[*menu.MenuItem][]*C.GtkWidget

func (f *Frontend) MenuSetApplicationMenu(menu *menu.Menu) {
	f.mainWindow.SetApplicationMenu(menu)
}

func (f *Frontend) MenuUpdateApplicationMenu() {
	//processMenu(f.mainWindow, f.mainWindow.applicationMenu)
}

func (w *Window) SetApplicationMenu(inmenu *menu.Menu) {
	if inmenu == nil {
		return
	}

	menuItemToId = make(map[*menu.MenuItem]int)
	menuIdToItem = make(map[int]*menu.MenuItem)
	gtkCheckboxCache = make(map[*menu.MenuItem][]*C.GtkWidget)

	// Increase ref count?
	w.menubar = C.gtk_menu_bar_new()

	processMenu(w, inmenu)

	C.gtk_widget_show(w.menubar)
}

func processMenu(window *Window, menu *menu.Menu) {
	for _, menuItem := range menu.Items {
		gtkMenu := C.gtk_menu_new()
		submenu := GtkMenuItemWithLabel(menuItem.Label)
		for _, menuItem := range menuItem.SubMenu.Items {
			menuID := menuIdCounter
			menuIdToItem[menuID] = menuItem
			menuItemToId[menuItem] = menuID
			menuIdCounter++
			processMenuItem(gtkMenu, menuItem, menuID)
		}
		C.gtk_menu_item_set_submenu(C.toGtkMenuItem(unsafe.Pointer(submenu)), gtkMenu)
		C.gtk_menu_shell_append(C.toGtkMenuShell(unsafe.Pointer(window.menubar)), submenu)
	}
}

func processMenuItem(parent *C.GtkWidget, menuItem *menu.MenuItem, menuID int) {
	if menuItem.Hidden {
		return
	}
	switch menuItem.Type {
	case menu.SeparatorType:
		separator := C.gtk_separator_menu_item_new()
		C.gtk_menu_shell_append(C.toGtkMenuShell(unsafe.Pointer(parent)), separator)
	case menu.TextType:
		gtkMenuItem := GtkMenuItemWithLabel(menuItem.Label)
		//if menuItem.Accelerator != nil {
		//	setAccelerator(gtkMenuItem, menuItem.Accelerator)
		//}

		C.gtk_menu_shell_append(C.toGtkMenuShell(unsafe.Pointer(parent)), gtkMenuItem)
		C.gtk_widget_show(gtkMenuItem)

		if menuItem.Click != nil {
			C.connectClick(gtkMenuItem, C.int(menuID))
		}

		if menuItem.Disabled {
			C.gtk_widget_set_sensitive(gtkMenuItem, 0)
		}

	case menu.CheckboxType:

		gtkMenuItem := GtkCheckMenuItemWithLabel(menuItem.Label)
		if menuItem.Checked {
			C.gtk_check_menu_item_set_active(C.toGtkCheckMenuItem(unsafe.Pointer(gtkMenuItem)), 1)
		}
		C.gtk_menu_shell_append(C.toGtkMenuShell(unsafe.Pointer(parent)), gtkMenuItem)
		C.gtk_widget_show(gtkMenuItem)

		if menuItem.Click != nil {
			C.connectClick(gtkMenuItem, C.int(menuID))
		}

		if menuItem.Disabled {
			C.gtk_widget_set_sensitive(gtkMenuItem, 0)
		}

		gtkCheckboxCache[menuItem] = append(gtkCheckboxCache[menuItem], gtkMenuItem)

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
