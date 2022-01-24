package linux

/*
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include "gtk/gtk.h"

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

//export handleMenuItemClick
func handleMenuItemClick(menuID int) {
	item := menuIdToItem[menuID]
	if item.Type == menu.CheckboxType {
		item.Checked = !item.Checked
	}
	item.Click(&menu.CallbackData{MenuItem: item})
}
