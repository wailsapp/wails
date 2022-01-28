//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include "gtk/gtk.h"

static GtkMenuItem *toGtkMenuItem(void *pointer) { return (GTK_MENU_ITEM(pointer)); }
static GtkMenuShell *toGtkMenuShell(void *pointer) { return (GTK_MENU_SHELL(pointer)); }
static GtkCheckMenuItem *toGtkCheckMenuItem(void *pointer) { return (GTK_CHECK_MENU_ITEM(pointer)); }
static GtkRadioMenuItem *toGtkRadioMenuItem(void *pointer) { return (GTK_RADIO_MENU_ITEM(pointer)); }

extern void handleMenuItemClick(void*);

void blockClick(GtkWidget* menuItem, gulong handler_id) {
	g_signal_handler_block (menuItem, handler_id);
}

void unblockClick(GtkWidget* menuItem, gulong handler_id) {
	g_signal_handler_unblock (menuItem, handler_id);
}

gulong connectClick(GtkWidget* menuItem) {
	return g_signal_connect(menuItem, "activate", G_CALLBACK(handleMenuItemClick), (void*)menuItem);
}

void addAccelerator(GtkWidget* menuItem, GtkAccelGroup* group, guint key, GdkModifierType mods) {
	gtk_widget_add_accelerator(menuItem, "activate", group, key, mods, GTK_ACCEL_VISIBLE);
}
*/
import "C"
import "github.com/wailsapp/wails/v2/pkg/menu"
import "unsafe"

var menuIdCounter int
var menuItemToId map[*menu.MenuItem]int
var menuIdToItem map[int]*menu.MenuItem
var gtkCheckboxCache map[*menu.MenuItem][]*C.GtkWidget
var gtkMenuCache map[*menu.MenuItem]*C.GtkWidget
var gtkRadioMenuCache map[*menu.MenuItem][]*C.GtkWidget
var gtkSignalHandlers map[*C.GtkWidget]C.gulong
var gtkSignalToMenuItem map[*C.GtkWidget]*menu.MenuItem

func (f *Frontend) MenuSetApplicationMenu(menu *menu.Menu) {
	f.mainWindow.SetApplicationMenu(menu)
}

func (f *Frontend) MenuUpdateApplicationMenu() {
	f.mainWindow.SetApplicationMenu(f.mainWindow.applicationMenu)
}

func (w *Window) SetApplicationMenu(inmenu *menu.Menu) {
	if inmenu == nil {
		return
	}

	// Setup accelerator group
	w.accels = C.gtk_accel_group_new()
	C.gtk_window_add_accel_group(w.asGTKWindow(), w.accels)

	menuItemToId = make(map[*menu.MenuItem]int)
	menuIdToItem = make(map[int]*menu.MenuItem)
	gtkCheckboxCache = make(map[*menu.MenuItem][]*C.GtkWidget)
	gtkMenuCache = make(map[*menu.MenuItem]*C.GtkWidget)
	gtkRadioMenuCache = make(map[*menu.MenuItem][]*C.GtkWidget)
	gtkSignalHandlers = make(map[*C.GtkWidget]C.gulong)
	gtkSignalToMenuItem = make(map[*C.GtkWidget]*menu.MenuItem)

	// Increase ref count?
	w.menubar = C.gtk_menu_bar_new()

	processMenu(w, inmenu)

	C.gtk_widget_show(w.menubar)
}

func processMenu(window *Window, menu *menu.Menu) {
	for _, menuItem := range menu.Items {
		submenu := processSubmenu(menuItem, window.accels)
		C.gtk_menu_shell_append(C.toGtkMenuShell(unsafe.Pointer(window.menubar)), submenu)
	}
}

func processSubmenu(menuItem *menu.MenuItem, group *C.GtkAccelGroup) *C.GtkWidget {
	existingMenu := gtkMenuCache[menuItem]
	if existingMenu != nil {
		return existingMenu
	}
	gtkMenu := C.gtk_menu_new()
	submenu := GtkMenuItemWithLabel(menuItem.Label)
	for _, menuItem := range menuItem.SubMenu.Items {
		menuID := menuIdCounter
		menuIdToItem[menuID] = menuItem
		menuItemToId[menuItem] = menuID
		menuIdCounter++
		processMenuItem(gtkMenu, menuItem, group)
	}
	C.gtk_menu_item_set_submenu(C.toGtkMenuItem(unsafe.Pointer(submenu)), gtkMenu)
	gtkMenuCache[menuItem] = existingMenu
	return submenu
}

var currentRadioGroup *C.GSList

func processMenuItem(parent *C.GtkWidget, menuItem *menu.MenuItem, group *C.GtkAccelGroup) {
	if menuItem.Hidden {
		return
	}

	if menuItem.Type != menu.RadioType {
		currentRadioGroup = nil
	}

	if menuItem.Type == menu.SeparatorType {
		result := C.gtk_separator_menu_item_new()
		C.gtk_menu_shell_append(C.toGtkMenuShell(unsafe.Pointer(parent)), result)
		return
	}

	var result *C.GtkWidget

	switch menuItem.Type {
	case menu.TextType:
		result = GtkMenuItemWithLabel(menuItem.Label)
	case menu.CheckboxType:
		result = GtkCheckMenuItemWithLabel(menuItem.Label)
		if menuItem.Checked {
			C.gtk_check_menu_item_set_active(C.toGtkCheckMenuItem(unsafe.Pointer(result)), 1)
		}
		gtkCheckboxCache[menuItem] = append(gtkCheckboxCache[menuItem], result)

	case menu.RadioType:
		result = GtkRadioMenuItemWithLabel(menuItem.Label, currentRadioGroup)
		currentRadioGroup = C.gtk_radio_menu_item_get_group(C.toGtkRadioMenuItem(unsafe.Pointer(result)))
		if menuItem.Checked {
			C.gtk_check_menu_item_set_active(C.toGtkCheckMenuItem(unsafe.Pointer(result)), 1)
		}
		gtkRadioMenuCache[menuItem] = append(gtkRadioMenuCache[menuItem], result)
	case menu.SubmenuType:
		result = processSubmenu(menuItem, group)
	}
	C.gtk_menu_shell_append(C.toGtkMenuShell(unsafe.Pointer(parent)), result)
	C.gtk_widget_show(result)

	if menuItem.Click != nil {
		handler := C.connectClick(result)
		gtkSignalHandlers[result] = handler
		gtkSignalToMenuItem[result] = menuItem
	}

	if menuItem.Disabled {
		C.gtk_widget_set_sensitive(result, 0)
	}

	if menuItem.Accelerator != nil {
		key, mods := acceleratorToGTK(menuItem.Accelerator)
		C.addAccelerator(result, group, key, mods)
	}
}
