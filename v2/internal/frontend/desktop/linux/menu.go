//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0
#cgo !webkit2_41 pkg-config: webkit2gtk-4.0
#cgo webkit2_41 pkg-config: webkit2gtk-4.1

#include "gtk/gtk.h"
#include "window.h"

static GtkMenuItem *toGtkMenuItem(void *pointer) { return (GTK_MENU_ITEM(pointer)); }
static GtkMenu *toGtkMenu(void *pointer) { return (GTK_MENU(pointer)); }
static GtkWindow *toGtkWindow(void *pointer) { return (GTK_WINDOW(pointer)); }
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
import (
	"encoding/base64"
	"os"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

var menuIdCounter int
var menuItemToId map[*menu.MenuItem]int
var menuIdToItem map[int]*menu.MenuItem
var gtkCheckboxCache map[*menu.MenuItem][]*C.GtkWidget
var gtkMenuCache map[*menu.MenuItem]*C.GtkWidget
var gtkRadioMenuCache map[*menu.MenuItem][]*C.GtkWidget
var gtkSignalHandlers map[*C.GtkWidget]C.gulong
var gtkSignalToMenuItem map[*C.GtkWidget]*menu.MenuItem

func initMaps() {
	// Ensure maps are initialized (they may not be if no application menu was set)
	if gtkSignalHandlers == nil {
		gtkSignalHandlers = make(map[*C.GtkWidget]C.gulong)
	}
	if gtkSignalToMenuItem == nil {
		gtkSignalToMenuItem = make(map[*C.GtkWidget]*menu.MenuItem)
	}
	if gtkCheckboxCache == nil {
		gtkCheckboxCache = make(map[*menu.MenuItem][]*C.GtkWidget)
	}
	if gtkRadioMenuCache == nil {
		gtkRadioMenuCache = make(map[*menu.MenuItem][]*C.GtkWidget)
	}
	if gtkMenuCache == nil {
		gtkMenuCache = make(map[*menu.MenuItem]*C.GtkWidget)
	}
	if menuItemToId == nil {
		menuItemToId = make(map[*menu.MenuItem]int)
	}
	if menuIdToItem == nil {
		menuIdToItem = make(map[int]*menu.MenuItem)
	}
}

func (f *Frontend) MenuSetApplicationMenu(menu *menu.Menu) {
	f.mainWindow.SetApplicationMenu(menu)
}

func (f *Frontend) MenuUpdateApplicationMenu() {
	f.mainWindow.SetApplicationMenu(f.mainWindow.applicationMenu)
}

func (f *Frontend) TraySetSystemTray(trayMenu *menu.TrayMenu) {
	if trayMenu == nil {
		return
	}

	initMaps()

	invokeOnMainThread(func() {
		var label *C.char
		if trayMenu.Label != "" {
			label = C.CString(trayMenu.Label)
			defer C.free(unsafe.Pointer(label))
		}

		var tooltip *C.char
		if trayMenu.Tooltip != "" {
			tooltip = C.CString(trayMenu.Tooltip)
			defer C.free(unsafe.Pointer(tooltip))
		}

		var imageData *C.guchar
		var imageLen C.gsize
		if trayMenu.Image != "" {
			// Try file
			if _, err := os.Stat(trayMenu.Image); err == nil {
				data, err := os.ReadFile(trayMenu.Image)
				if err == nil && len(data) > 0 {
					imageData = (*C.guchar)(unsafe.Pointer(&data[0]))
					imageLen = C.gsize(len(data))
				}
			} else {
				// Try base64
				data, err := base64.StdEncoding.DecodeString(trayMenu.Image)
				if err == nil && len(data) > 0 {
					imageData = (*C.guchar)(unsafe.Pointer(&data[0]))
					imageLen = C.gsize(len(data))
				}
			}
		}

		if f.mainWindow.trayAccelGroup != nil {
			C.gtk_window_remove_accel_group(C.toGtkWindow(f.mainWindow.gtkWindow), f.mainWindow.trayAccelGroup)
			C.g_object_unref(C.gpointer(f.mainWindow.trayAccelGroup))
			f.mainWindow.trayAccelGroup = nil
		}

		var gtkMenu *C.GtkWidget
		if trayMenu.Menu != nil {
			gtkMenu = C.gtk_menu_new()
			f.mainWindow.trayAccelGroup = C.gtk_accel_group_new()
			C.gtk_window_add_accel_group(C.toGtkWindow(f.mainWindow.gtkWindow), f.mainWindow.trayAccelGroup)
			C.gtk_menu_set_accel_group(C.toGtkMenu(unsafe.Pointer(gtkMenu)), f.mainWindow.trayAccelGroup)
			for _, item := range trayMenu.Menu.Items {
				processMenuItem(gtkMenu, item, f.mainWindow.trayAccelGroup)
			}
		}

		C.TraySetSystemTray(C.toGtkWindow(f.mainWindow.gtkWindow), label, imageData, imageLen, tooltip, gtkMenu)
	})
}

func (w *Window) SetApplicationMenu(inmenu *menu.Menu) {
	if inmenu == nil {
		return
	}

	// Setup accelerator group
	w.accels = C.gtk_accel_group_new()
	C.gtk_window_add_accel_group(w.asGTKWindow(), w.accels)

	initMaps()

	// Increase ref count?
	w.menubar = C.gtk_menu_bar_new()

	processMenu(w, inmenu)

	C.gtk_widget_show(w.menubar)
}

func processMenu(window *Window, menu *menu.Menu) {
	for _, menuItem := range menu.Items {
		if menuItem.SubMenu != nil {
			submenu := processSubmenu(menuItem, window.accels)
			C.gtk_menu_shell_append(C.toGtkMenuShell(unsafe.Pointer(window.menubar)), submenu)
		}
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
	gtkMenuCache[menuItem] = gtkMenu
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

	if menuItem.Accelerator != nil && group != nil {
		key, mods := acceleratorToGTK(menuItem.Accelerator)
		C.addAccelerator(result, group, key, mods)
	}
}
