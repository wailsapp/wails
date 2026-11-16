//go:build linux && purego && gtk3 && !android && !server

package application

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
	"unsafe"

	"github.com/ebitengine/purego"
)

type gtk3ViewsTestApp struct{ platformApp }

func (*gtk3ViewsTestApp) isOnMainThread() bool { return true }

func TestPureGoGTK3InheritedMenuViews(t *testing.T) {
	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		t.Skip("GTK display unavailable")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestPureGoGTK3InheritedMenuViewsProcess$", "-test.v")
	cmd.Env = append(os.Environ(), "WAILS_TEST_PUREGO_GTK3_VIEWS=1", "G_DEBUG=fatal-criticals")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("GTK3 menu views: %v\n%s", err, output)
	}
}

func TestPureGoGTK3InheritedMenuViewsProcess(t *testing.T) {
	if os.Getenv("WAILS_TEST_PUREGO_GTK3_VIEWS") != "1" {
		t.Skip("run by TestPureGoGTK3InheritedMenuViews")
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	ensureGTK3Libraries()
	if gtk3LibsErr != nil {
		t.Fatal(gtk3LibsErr)
	}
	lib, err := purego.Dlopen("libgtk-3.so.0", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		t.Fatal(err)
	}
	var initCheck func(uintptr, uintptr) int32
	var windowNew func(int32) uintptr
	var labelNew func(uintptr) uintptr
	var accelActivate func(uintptr, uint32, uint32) int32
	purego.RegisterLibFunc(&initCheck, lib, "gtk_init_check")
	purego.RegisterLibFunc(&windowNew, lib, "gtk_window_new")
	purego.RegisterLibFunc(&labelNew, lib, "gtk_label_new")
	purego.RegisterLibFunc(&accelActivate, lib, "gtk_accel_groups_activate")
	if initCheck(0, 0) == 0 {
		t.Fatal("GTK display unavailable")
	}
	globalApplication = &App{impl: &gtk3ViewsTestApp{}, running: true, Logger: slog.Default()}
	menu := NewMenu()
	file := menu.AddSubmenu("File")
	action := file.Add("Action").SetAccelerator("ctrl+1").OnClick(func(*Context) {})
	check := file.AddCheckbox("Check", false)
	file.AddRadio("First", true)
	radio := file.AddRadio("Second", false).SetAccelerator("ctrl+3").OnClick(func(*Context) {})
	menu.Update()
	globalApplication.applicationMenu = menu
	newWindow := func() *linuxWebviewWindow {
		window := pointer(unsafe.Pointer(windowNew(0)))
		windowAddMenuAccelGroup(window)
		box := pointer(_Cfunc_gtk_box_new(_Ciconst_GTK_ORIENTATION_VERTICAL, 0))
		_Cfunc_gtk_container_add((*_Ctype_GtkContainer)(window), (*_Ctype_GtkWidget)(box))
		label := pointer(unsafe.Pointer(labelNew(0)))
		_Cfunc_gtk_box_pack_start((*_Ctype_GtkBox)(box), (*_Ctype_GtkWidget)(label), 1, 1, 0)
		_Cfunc_gtk_widget_show_all((*_Ctype_GtkWidget)(window))
		w := &linuxWebviewWindow{window: window, vbox: box, webview: label, parent: &WebviewWindow{}}
		w.menubar = windowAttachMenuBar(window, box, windowMenuNative(menu), LinuxMenuStyleMenuBar)
		return w
	}
	first, second := newWindow(), newWindow()
	activate := func(w *linuxWebviewWindow, key uint32) bool {
		matched := accelActivate(uintptr(w.window), key, 1<<2) != 0
		if matched {
			select {
			case <-menuItemClicked:
			default:
			}
		}
		return matched
	}
	if !activate(second, '1') {
		t.Fatal("second window cannot activate inherited shortcut")
	}
	check.SetChecked(true)
	if accelActivate(uintptr(second.window), '3', 1<<2) == 0 {
		t.Fatal("radio shortcut did not activate")
	}
	select {
	case id := <-menuItemClicked:
		if id != radio.id {
			t.Fatalf("radio activation id = %d", id)
		}
	default:
		t.Fatal("radio activation not dispatched")
	}
	radio.handleClick()
	for view := range gtk3MenuViews[menu] {
		submenu := view.items[0].submenu
		if !submenu.FindByLabel("Check").impl.(*linuxMenuItem).isChecked() || !submenu.FindByLabel("Second").impl.(*linuxMenuItem).isChecked() || submenu.FindByLabel("First").impl.(*linuxMenuItem).isChecked() {
			t.Fatal("check or radio state did not propagate")
		}
	}
	first.setMenu(nil)
	if activate(first, '1') || !activate(second, '1') {
		t.Fatal("removal did not isolate window shortcuts")
	}
	first.setMenu(menu)
	action.SetAccelerator("ctrl+2")
	if activate(first, '1') || !activate(first, '2') || !activate(second, '2') {
		t.Fatal("accelerator update did not propagate")
	}
	file.Add("Added")
	file.Update()
	for view := range gtk3MenuViews[menu] {
		if view.items[0].submenu.FindByLabel("Added") == nil {
			t.Fatal("submenu update did not propagate")
		}
	}
	first.setMenu(nil)
	second.setMenu(nil)
	if len(gtk3MenuViews[menu]) != 0 {
		t.Fatal("window menu views leaked")
	}
}
