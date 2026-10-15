//go:build linux && cgo && gtk3 && wailsintegration && !android && !server

package application

/*
#include <gtk/gtk.h>
*/
import "C"
import "unsafe"

func gtkMenuTestInit() bool {
	return C.gtk_init_check(nil, nil) != 0
}

func gtkMenuTestLabels(menu *Menu) []string {
	children := C.gtk_container_get_children((*C.GtkContainer)(menu.impl.(*linuxMenu).native))
	defer C.g_list_free(children)
	var labels []string
	for child := children; child != nil; child = child.next {
		text := ""
		if box := C.gtk_bin_get_child((*C.GtkBin)(child.data)); box != nil {
			if contents := C.gtk_container_get_children((*C.GtkContainer)(unsafe.Pointer(box))); contents != nil {
				if last := C.g_list_last(contents); last != nil {
					text = C.GoString(C.gtk_label_get_text((*C.GtkLabel)(last.data)))
				}
				C.g_list_free(contents)
			}
		}
		labels = append(labels, text)
	}
	return labels
}

func gtkMenuTestActivate(item *MenuItem) {
	C.gtk_menu_item_activate((*C.GtkMenuItem)(item.impl.(*linuxMenuItem).native))
}

func gtkMenuTestOwn(widget pointer) {
	C.g_object_ref_sink(C.gpointer(widget))
}

func gtkMenuTestRelease(widget pointer) {
	C.gtk_widget_destroy((*C.GtkWidget)(widget))
	C.g_object_unref(C.gpointer(unsafe.Pointer(widget)))
}
