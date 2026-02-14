//go:build linux
// +build linux

package linux

/*
#cgo linux pkg-config: gtk+-3.0
#cgo !webkit2_41 pkg-config: webkit2gtk-4.0
#cgo webkit2_41 pkg-config: webkit2gtk-4.1

#include "gtk/gtk.h"

static GtkCheckMenuItem *toGtkCheckMenuItem(void *pointer) { return (GTK_CHECK_MENU_ITEM(pointer)); }

extern void blockClick(GtkWidget* menuItem, gulong handler_id);
extern void unblockClick(GtkWidget* menuItem, gulong handler_id);
*/
import "C"
import (
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/menu"
)

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

func GtkRadioMenuItemWithLabel(label string, group *C.GSList) *C.GtkWidget {
	cLabel := C.CString(label)
	result := C.gtk_radio_menu_item_new_with_label(group, cLabel)
	C.free(unsafe.Pointer(cLabel))
	return result
}

//export handleMenuItemClick
func handleMenuItemClick(gtkWidget unsafe.Pointer) {
	// Make sure to execute the final callback on a new goroutine otherwise if the callback e.g. tries to open a dialog, the
	// main thread will get blocked and so the message loop blocks. As a result the app will block and shows a
	// "not responding" dialog.

	widget := (*C.GtkWidget)(gtkWidget)
	item := appMenuCache.gtkSignalToMenuItem[widget]
	if item == nil {
		item = trayMenuCache.gtkSignalToMenuItem[widget]
	}
	if item == nil {
		return
	}

	switch item.Type {
	case menu.CheckboxType:
		item.Checked = !item.Checked
		checked := C.int(0)
		if item.Checked {
			checked = C.int(1)
		}

		updateCheckbox := func(cache *menuCache) {
			for _, gtkCheckbox := range cache.gtkCheckboxCache[item] {
				handler := cache.gtkSignalHandlers[gtkCheckbox]
				C.blockClick(gtkCheckbox, handler)
				C.gtk_check_menu_item_set_active(C.toGtkCheckMenuItem(unsafe.Pointer(gtkCheckbox)), checked)
				C.unblockClick(gtkCheckbox, handler)
			}
		}
		updateCheckbox(appMenuCache)
		updateCheckbox(trayMenuCache)

		go item.Click(&menu.CallbackData{MenuItem: item})
	case menu.RadioType:
		active := C.gtk_check_menu_item_get_active(C.toGtkCheckMenuItem(gtkWidget))
		if int(active) == 1 {
			updateRadio := func(cache *menuCache) {
				for _, gtkRadioItem := range cache.gtkRadioMenuCache[item] {
					handler := cache.gtkSignalHandlers[gtkRadioItem]
					C.blockClick(gtkRadioItem, handler)
					C.gtk_check_menu_item_set_active(C.toGtkCheckMenuItem(unsafe.Pointer(gtkRadioItem)), 1)
					C.unblockClick(gtkRadioItem, handler)
				}
			}
			updateRadio(appMenuCache)
			updateRadio(trayMenuCache)
			item.Checked = true
			go item.Click(&menu.CallbackData{MenuItem: item})
		} else {
			item.Checked = false
		}
	default:
		go item.Click(&menu.CallbackData{MenuItem: item})
	}
}
