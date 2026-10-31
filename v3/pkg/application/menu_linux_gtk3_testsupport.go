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

// gtkMenuTestWindow returns a shown window with the menu bar packed in, as a
// real window has it: GTK only activates accelerators of drawable widgets.
func gtkMenuTestWindow(menuBar pointer) pointer {
	window := pointer(C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL))
	windowAddMenuAccelGroup(window)
	box := C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0)
	C.gtk_box_pack_start((*C.GtkBox)(unsafe.Pointer(box)), (*C.GtkWidget)(menuBar), 0, 0, 0)
	C.gtk_container_add((*C.GtkContainer)(unsafe.Pointer(window)), (*C.GtkWidget)(unsafe.Pointer(box)))
	C.gtk_widget_show_all((*C.GtkWidget)(unsafe.Pointer(window)))
	return window
}

// gtkMenuTestAccelActivate presses a chord in the window: it reports whether
// a GTK accelerator matched, which runs the item's activate signal.
func gtkMenuTestAccelActivate(window pointer, keyval uint32, mods uint32) bool {
	return C.gtk_accel_groups_activate((*C.GObject)(window), C.guint(keyval), C.GdkModifierType(mods)) != 0
}

// gtkMenuTestAccelLabelWidth returns the width of the shortcut text the item's
// accel label draws: zero when the label shows no shortcut. GTK 3 keeps the
// shortcut text private, so the width stands in for it.
func gtkMenuTestAccelLabelWidth(item *MenuItem) int {
	box := C.gtk_bin_get_child((*C.GtkBin)(item.impl.(*linuxMenuItem).native))
	if box == nil {
		return 0
	}
	children := C.gtk_container_get_children((*C.GtkContainer)(unsafe.Pointer(box)))
	defer C.g_list_free(children)
	last := C.g_list_last(children)
	if last == nil {
		return 0
	}
	// The width is computed when the label is sized: ask for its size first.
	var minimum, natural C.gint
	C.gtk_widget_get_preferred_width((*C.GtkWidget)(last.data), &minimum, &natural)
	return int(C.gtk_accel_label_get_accel_width((*C.GtkAccelLabel)(last.data)))
}

// gtkMenuTestAccelCount counts the accelerators bound to the item's widget.
func gtkMenuTestAccelCount(item *MenuItem) int {
	closures := C.gtk_widget_list_accel_closures((*C.GtkWidget)(item.impl.(*linuxMenuItem).native))
	defer C.g_list_free(closures)
	return int(C.g_list_length(closures))
}

// gtkMenuTestClearAccelerator calls the native removal directly, as an
// unset accelerator would.
func gtkMenuTestClearAccelerator(item *MenuItem) {
	item.impl.(*linuxMenuItem).setAccelerator(nil)
}

// gtkMenuTestRealisedWindow builds a shown window laid out as windowNew lays
// it out: a box with a menu bar slot above the webview. A label stands in for
// the webview, since WebKit needs a GtkApplication and the runtime SetMenu
// tests only touch the menu bar and its accelerators.
func gtkMenuTestRealisedWindow() *linuxWebviewWindow {
	window := pointer(C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL))
	windowAddMenuAccelGroup(window)
	vbox := pointer(C.gtk_box_new(C.GTK_ORIENTATION_VERTICAL, 0))
	C.gtk_container_add((*C.GtkContainer)(window), (*C.GtkWidget)(vbox))
	webview := pointer(C.gtk_label_new(nil))
	C.gtk_box_pack_start((*C.GtkBox)(vbox), (*C.GtkWidget)(webview), 1, 1, 0)
	C.gtk_widget_show_all((*C.GtkWidget)(window))
	return &linuxWebviewWindow{window: window, vbox: vbox, webview: webview, parent: &WebviewWindow{}}
}

// gtkMenuTestBoxChildren lists the children of a GtkBox in order.
func gtkMenuTestBoxChildren(box pointer) []pointer {
	children := C.gtk_container_get_children((*C.GtkContainer)(box))
	defer C.g_list_free(children)
	var out []pointer
	for child := children; child != nil; child = child.next {
		out = append(out, pointer(child.data))
	}
	return out
}

// gtkMenuTestNative returns the native menu bar widget of a menu.
func gtkMenuTestNative(menu *Menu) pointer {
	return menu.impl.(*linuxMenu).native
}

func gtkMenuTestItemLabel(item *MenuItem) string {
	box := C.gtk_bin_get_child((*C.GtkBin)(item.impl.(*linuxMenuItem).native))
	if box == nil {
		return ""
	}
	children := C.gtk_container_get_children((*C.GtkContainer)(unsafe.Pointer(box)))
	defer C.g_list_free(children)
	last := C.g_list_last(children)
	if last == nil {
		return ""
	}
	return C.GoString(C.gtk_label_get_text((*C.GtkLabel)(last.data)))
}
