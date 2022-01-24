package linux

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include "gtk/gtk.h"

static GtkCheckMenuItem *toGtkCheckMenuItem(void *pointer) { return (GTK_CHECK_MENU_ITEM(pointer)); }

*/
import "C"
import "unsafe"
import "github.com/wailsapp/wails/v2/pkg/menu"

func GtkMenuItemWithLabel(label string) *C.GtkWidget {
	cLabel := C.CString(label)
	result := C.gtk_menu_item_new_with_label(cLabel)
	C.free(unsafe.Pointer(cLabel))
	return result
}

func GtkCheckMenuItemWithLabel(label string) *C.GtkWidget {
	cLabel := C.CString(label)
	result := C.gtk_check_menu_item_new_with_label(cLabel)
	C.free(unsafe.Pointer(cLabel))
	return result
}

var menuItemBeingProcessed *menu.MenuItem

//export handleMenuItemClick
func handleMenuItemClick(menuID int) {
	item := menuIdToItem[menuID]
	// This is here because setting multiple checkboxes to active triggers the gtk action for each one,
	// meaning multiple calls to this method. As it's all in the same thread, we can short circuit the other calls
	// by keeping a "lock". When we deal with multiple windows, we will need to take the window ID into consideration too
	if item == menuItemBeingProcessed {
		return
	}
	menuItemBeingProcessed = item
	if item.Type == menu.CheckboxType {
		item.Checked = !item.Checked
		checked := C.int(0)
		if item.Checked {
			checked = C.int(1)
		}
		for _, gtkCheckbox := range gtkCheckboxCache[item] {
			C.gtk_check_menu_item_set_active(C.toGtkCheckMenuItem(unsafe.Pointer(gtkCheckbox)), checked)
		}
	}
	item.Click(&menu.CallbackData{MenuItem: item})
	menuItemBeingProcessed = nil
}
