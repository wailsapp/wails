//go:build linux && purego && !gtk3 && !android && !server

package application

import (
	"runtime"
	"testing"
)

// The DND tests run enableDND and disableDND on a real GtkLabel and count the
// event controllers attached to it. Enabling installs a motion controller and a
// GtkDropTarget; disabling installs nothing, which mirrors the empty C.disableDND
// body in linux_cgo.c. Like the other GTK tests they run in their own process
// and skip without a display.

// dndTestLabel returns a real GtkLabel owned by the caller, who must unref it.
func dndTestLabel(t *testing.T) uintptr {
	t.Helper()
	requireGTKDisplay(t)
	label := g_object_ref_sink(gtk_label_new("dnd"))
	if label == 0 {
		t.Fatal("gtk_label_new returned NULL")
	}
	return label
}

// controllerCount returns the number of event controllers on widget.
func controllerCount(t *testing.T, widget uintptr) uint32 {
	t.Helper()
	controllers := tGtkWidgetObserveControllers(widget)
	if controllers == 0 {
		t.Fatal("gtk_widget_observe_controllers returned NULL")
	}
	defer g_object_unref(controllers)
	return g_list_model_get_n_items(controllers)
}

// hasDropTarget reports whether widget has a GtkDropTarget controller.
func hasDropTarget(t *testing.T, widget uintptr) bool {
	t.Helper()
	controllers := tGtkWidgetObserveControllers(widget)
	if controllers == 0 {
		t.Fatal("gtk_widget_observe_controllers returned NULL")
	}
	defer g_object_unref(controllers)
	dropTarget := tGtkDropTargetGetType()
	for i := range g_list_model_get_n_items(controllers) {
		controller := g_list_model_get_item(controllers, i)
		isTarget := g_type_check_instance_is_a(controller, dropTarget) != 0
		g_object_unref(controller)
		if isTarget {
			return true
		}
	}
	return false
}

func TestDNDEnableInstallsDropTarget(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	runtime.LockOSThread()
	t.Cleanup(runtime.UnlockOSThread)
	label := dndTestLabel(t)
	t.Cleanup(func() { g_object_unref(label) })

	w := &linuxWebviewWindow{webview: pointer(label), parent: &WebviewWindow{id: 42}}
	before := controllerCount(t, label)
	if hasDropTarget(t, label) {
		t.Fatal("a fresh label must not have a GtkDropTarget")
	}
	w.enableDND()
	// enableDNDGo adds a GtkDropControllerMotion and a GtkDropTarget.
	if got := controllerCount(t, label) - before; got != 2 {
		t.Fatalf("enableDND added %d controllers, want 2 (motion and drop target)", got)
	}
	if !hasDropTarget(t, label) {
		t.Fatal("enableDND did not install a GtkDropTarget on the widget")
	}
}

func TestDNDDisableAddsNoControllers(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	runtime.LockOSThread()
	t.Cleanup(runtime.UnlockOSThread)
	label := dndTestLabel(t)
	t.Cleanup(func() { g_object_unref(label) })

	w := &linuxWebviewWindow{webview: pointer(label), parent: &WebviewWindow{id: 42}}
	before := controllerCount(t, label)
	w.disableDND()
	if got := controllerCount(t, label) - before; got != 0 {
		t.Fatalf("disableDND added %d controllers, want 0", got)
	}
	if hasDropTarget(t, label) {
		t.Fatal("disableDND installed a GtkDropTarget")
	}
}
