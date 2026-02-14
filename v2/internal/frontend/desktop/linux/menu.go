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
	"runtime"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

type menuCache struct {
	menuIdCounter       int
	menuItemToId        map[*menu.MenuItem]int
	menuIdToItem        map[int]*menu.MenuItem
	gtkCheckboxCache    map[*menu.MenuItem][]*C.GtkWidget
	gtkMenuCache        map[*menu.MenuItem]*C.GtkWidget
	gtkRadioMenuCache   map[*menu.MenuItem][]*C.GtkWidget
	gtkSignalHandlers   map[*C.GtkWidget]C.gulong
	gtkSignalToMenuItem map[*C.GtkWidget]*menu.MenuItem
}

func newMenuCache() *menuCache {
	return &menuCache{
		menuItemToId:        make(map[*menu.MenuItem]int),
		menuIdToItem:        make(map[int]*menu.MenuItem),
		gtkCheckboxCache:    make(map[*menu.MenuItem][]*C.GtkWidget),
		gtkMenuCache:        make(map[*menu.MenuItem]*C.GtkWidget),
		gtkRadioMenuCache:   make(map[*menu.MenuItem][]*C.GtkWidget),
		gtkSignalHandlers:   make(map[*C.GtkWidget]C.gulong),
		gtkSignalToMenuItem: make(map[*C.GtkWidget]*menu.MenuItem),
	}
}

var appMenuCache = newMenuCache()
var trayMenuCache = newMenuCache()

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

	invokeOnMainThread(func() {
		trayMenuCache = newMenuCache()
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
		var imageBytes []byte
		if trayMenu.Image != "" {
			// Try file
			if _, err := os.Stat(trayMenu.Image); err == nil {
				data, err := os.ReadFile(trayMenu.Image)
				if err == nil && len(data) > 0 {
					imageBytes = data
				}
			} else {
				// Try base64
				data, err := base64.StdEncoding.DecodeString(trayMenu.Image)
				if err == nil && len(data) > 0 {
					imageBytes = data
				}
			}
			if len(imageBytes) > 0 {
				imageData = (*C.guchar)(unsafe.Pointer(&imageBytes[0]))
				imageLen = C.gsize(len(imageBytes))
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
				processMenuItem(gtkMenu, item, f.mainWindow.trayAccelGroup, trayMenuCache)
			}
		}

		C.TraySetSystemTray(C.toGtkWindow(f.mainWindow.gtkWindow), label, imageData, imageLen, tooltip, gtkMenu)
		runtime.KeepAlive(imageBytes)
	})
}

func (w *Window) SetApplicationMenu(inmenu *menu.Menu) {
	if inmenu == nil {
		return
	}

	// Setup accelerator group
	w.accels = C.gtk_accel_group_new()
	C.gtk_window_add_accel_group(w.asGTKWindow(), w.accels)

	appMenuCache = newMenuCache()

	// Increase ref count?
	w.menubar = C.gtk_menu_bar_new()

	processMenu(w, inmenu)

	C.gtk_widget_show(w.menubar)
}

func processMenu(window *Window, menu *menu.Menu) {
	for _, menuItem := range menu.Items {
		if menuItem.SubMenu != nil {
			submenu := processSubmenu(menuItem, window.accels, appMenuCache)
			C.gtk_menu_shell_append(C.toGtkMenuShell(unsafe.Pointer(window.menubar)), submenu)
		}
	}
}

func processSubmenu(menuItem *menu.MenuItem, group *C.GtkAccelGroup, cache *menuCache) *C.GtkWidget {
	existingMenu := cache.gtkMenuCache[menuItem]
	if existingMenu != nil {
		return existingMenu
	}
	gtkMenu := C.gtk_menu_new()
	submenu := GtkMenuItemWithLabel(menuItem.Label)
	for _, menuItem := range menuItem.SubMenu.Items {
		menuID := cache.menuIdCounter
		cache.menuIdToItem[menuID] = menuItem
		cache.menuItemToId[menuItem] = menuID
		cache.menuIdCounter++
		processMenuItem(gtkMenu, menuItem, group, cache)
	}
	C.gtk_menu_item_set_submenu(C.toGtkMenuItem(unsafe.Pointer(submenu)), gtkMenu)
	cache.gtkMenuCache[menuItem] = gtkMenu
	return submenu
}

var currentRadioGroup *C.GSList

func processMenuItem(parent *C.GtkWidget, menuItem *menu.MenuItem, group *C.GtkAccelGroup, cache *menuCache) {
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
		cache.gtkCheckboxCache[menuItem] = append(cache.gtkCheckboxCache[menuItem], result)

	case menu.RadioType:
		result = GtkRadioMenuItemWithLabel(menuItem.Label, currentRadioGroup)
		currentRadioGroup = C.gtk_radio_menu_item_get_group(C.toGtkRadioMenuItem(unsafe.Pointer(result)))
		if menuItem.Checked {
			C.gtk_check_menu_item_set_active(C.toGtkCheckMenuItem(unsafe.Pointer(result)), 1)
		}
		cache.gtkRadioMenuCache[menuItem] = append(cache.gtkRadioMenuCache[menuItem], result)
	case menu.SubmenuType:
		result = processSubmenu(menuItem, group, cache)
	}
	C.gtk_menu_shell_append(C.toGtkMenuShell(unsafe.Pointer(parent)), result)
	C.gtk_widget_show(result)

	if menuItem.Click != nil {
		handler := C.connectClick(result)
		cache.gtkSignalHandlers[result] = handler
		cache.gtkSignalToMenuItem[result] = menuItem
	}

	if menuItem.Disabled {
		C.gtk_widget_set_sensitive(result, 0)
	}

	if menuItem.Accelerator != nil && group != nil {
		key, mods := acceleratorToGTK(menuItem.Accelerator)
		C.addAccelerator(result, group, key, mods)
	}
}
