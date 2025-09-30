//go:build linux && webkit_6
// +build linux,webkit_6

package linux

/*
#cgo pkg-config: gtk4 webkitgtk-6.0

#include "gtk/gtk.h"
#include <string.h>

static GActionMap *toActionMap(GtkWindow *window) { return (G_ACTION_MAP(window)); }
static GAction *toGAction(GSimpleAction *action) { return (G_ACTION(action)); }
static GMenuModel *toGMenuModel(GMenu *menu) { return (G_MENU_MODEL(menu)); }

extern void handleMenuItemClick(char* aid);
extern void handleMenuCheckItemClick(char* aid, int checked);
extern void handleMenuRadioItemClick(char* radioId, char* prev, char* curr);

static void onAction(GAction *action, GVariant *param) {
	GVariantType *stateType = g_action_get_state_type(action);

	if(stateType != NULL) {
		GVariant *state = g_action_get_state(action);
		gchar *stateStr = g_variant_type_dup_string(stateType);

		if(strcmp(stateStr, "s") == 0) {
			g_simple_action_set_state(G_SIMPLE_ACTION(action), param);

			handleMenuRadioItemClick(
				g_action_get_name(action),
				g_variant_get_string(state, NULL),
				g_variant_get_string(param, NULL));

		} else if(strcmp(stateStr, "b") == 0) {
			gboolean checked = !g_variant_get_boolean(state);
			GVariant *newState = g_variant_new_boolean(checked);

			g_simple_action_set_state(G_SIMPLE_ACTION(action), newState);

			handleMenuCheckItemClick(g_action_get_name(action), checked);
		}

		if(state != NULL) {
			g_variant_unref(state);
		}

		if(stateStr != NULL) {
			g_free(stateStr);
		}
	} else {
		handleMenuItemClick(g_action_get_name(action));
	}
}

gulong connectClick(GSimpleAction *action) {
	return g_signal_connect(action, "activate", G_CALLBACK(onAction), NULL);
}

void setAccels(GtkApplication *app, char *actionName, char *accels) {
	gtk_application_set_accels_for_action(app, actionName, (const char *[]) { accels, NULL });
}
*/
import "C"
import (
	"strings"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

var menuIdCounter int
var menuItemToId map[*menu.MenuItem]int
var menuIdToItem map[int]*menu.MenuItem
var gtkMenuCache map[*menu.MenuItem]*C.GMenu
var gActionIdToMenuItem sync.Map

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

	// Clean up previous menu
	for _, gmenu := range gtkMenuCache {
		if gmenu != nil {
			C.g_object_unref(C.gpointer(gmenu))
		}
	}

	if w.menubar != nil {
		C.gtk_widget_unparent(w.menubar)
		w.menubar = nil
	}

	menuItemToId = make(map[*menu.MenuItem]int)
	menuIdToItem = make(map[int]*menu.MenuItem)
	gtkMenuCache = make(map[*menu.MenuItem]*C.GMenu)

	processMenu(w, inmenu)
}

func processMenu(window *Window, menu *menu.Menu) {
	gmenu := C.g_menu_new()
	defer C.g_object_unref(C.gpointer(gmenu))

	for _, menuItem := range menu.Items {
		itemLabel := C.CString(menuItem.Label)
		defer C.free(unsafe.Pointer(itemLabel))

		submenu := processSubmenu(window, menuItem)
		defer C.g_object_unref(C.gpointer(submenu))

		C.g_menu_append_submenu(gmenu, itemLabel, C.toGMenuModel(submenu))
	}

	window.menubar = C.gtk_popover_menu_bar_new_from_model(C.toGMenuModel(gmenu))
}

func processSubmenu(window *Window, menuItem *menu.MenuItem) *C.GMenu {
	existingMenu := gtkMenuCache[menuItem]

	if existingMenu != nil {
		return existingMenu
	}

	submenu := C.g_menu_new()

	for _, subItem := range menuItem.SubMenu.Items {
		menuID := menuIdCounter
		menuIdToItem[menuID] = subItem
		menuItemToId[subItem] = menuID
		menuIdCounter++

		processMenuItem(window, submenu, subItem)
	}

	gtkMenuCache[menuItem] = submenu

	return submenu
}

var currentRadioActionId string

func processMenuItem(window *Window, parent *C.GMenu, menuItem *menu.MenuItem) {
	if menuItem.Hidden {
		return
	}

	if menuItem.Type != menu.RadioType {
		currentRadioActionId = ""
	}

	var action *C.GSimpleAction

	itemId := strings.ReplaceAll(strings.ToLower(menuItem.Label), " ", "-")
	actionName := itemId

	switch menuItem.Type {
	case menu.SubmenuType:
		itemLabel := C.CString(menuItem.Label)
		defer C.free(unsafe.Pointer(itemLabel))

		submenu := processSubmenu(window, menuItem)
		defer C.g_object_unref(C.gpointer(submenu))

		C.g_menu_append_submenu(parent, itemLabel, C.toGMenuModel(submenu))
		return

	case menu.SeparatorType:
		return

	case menu.CheckboxType:
		name := C.CString(actionName)
		defer C.free(unsafe.Pointer(name))

		action = C.g_simple_action_new_stateful(name, nil, C.g_variant_new_boolean(gtkBool(menuItem.Checked)))

	case menu.RadioType:
		if currentRadioActionId == "" {
			currentRadioActionId = itemId
		}

		if menuItem.Checked {
			name := C.CString(currentRadioActionId)
			defer C.free(unsafe.Pointer(name))

			vType := C.CString("s")
			defer C.free(unsafe.Pointer(vType))

			paramType := C.g_variant_type_new(vType)
			defer C.g_variant_type_free(paramType)

			vContent := C.CString(itemId)
			defer C.free(unsafe.Pointer(vContent))

			action = C.g_simple_action_new_stateful(name, paramType, C.g_variant_new_string(vContent))

			C.g_action_map_add_action(C.toActionMap(window.asGTKWindow()), C.toGAction(action))
		}

		// Use currentRadioActionId as the Action Name and itemId as the Target
		actionName = currentRadioActionId + "::" + itemId

	default:
		name := C.CString(actionName)
		defer C.free(unsafe.Pointer(name))

		action = C.g_simple_action_new(name, nil)
	}

	if currentRadioActionId == "" {
		C.g_action_map_add_action(C.toActionMap(window.asGTKWindow()), C.toGAction(action))
	}

	if action != nil {
		defer C.g_object_unref(C.gpointer(action))

		if menuItem.Disabled {
			C.g_simple_action_set_enabled(action, gtkBool(false))
		}

		if menuItem.Click != nil {
			C.connectClick(action)
		}
	}

	gActionIdToMenuItem.Store(actionName, menuItem)

	detActionName := C.CString("win." + actionName)
	defer C.free(unsafe.Pointer(detActionName))

	itemName := C.CString(menuItem.Label)
	defer C.free(unsafe.Pointer(itemName))

	item := C.g_menu_item_new(itemName, detActionName)
	defer C.g_object_unref(C.gpointer(item))
	C.g_menu_append_item(parent, item)

	if menuItem.Accelerator != nil {
		key, mods := acceleratorToGTK(menuItem.Accelerator)

		accelName := C.gtk_accelerator_name(key, mods)
		defer C.free(unsafe.Pointer(accelName))

		C.setAccels(window.gtkApp, detActionName, accelName)
	}
}
