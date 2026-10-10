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
		box := C.gtk_bin_get_child((*C.GtkBin)(child.data))
		contents := C.gtk_container_get_children((*C.GtkContainer)(unsafe.Pointer(box)))
		label := (*C.GtkLabel)(C.g_list_last(contents).data)
		labels = append(labels, C.GoString(C.gtk_label_get_text(label)))
		C.g_list_free(contents)
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
