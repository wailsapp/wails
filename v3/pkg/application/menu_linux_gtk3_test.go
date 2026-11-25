//go:build linux && cgo && gtk3 && wailsintegration && !android && !server

package application

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"testing"
	"time"
)

// Each native run owns a GTK thread and treats critical warnings as failures.
func TestGTK3MenuRebuild(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGTK3MenuRebuildProcess$", "-test.v")
	cmd.Env = append(os.Environ(), "WAILS_TEST_GTK3_MENU=1", "G_DEBUG=fatal-criticals")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("native menu regression: %v\n%s", err, output)
	}
}

type gtkMenuTestApp struct{ platformApp }

func (*gtkMenuTestApp) isOnMainThread() bool { return true }

func TestGTK3MenuRebuildProcess(t *testing.T) {
	if os.Getenv("WAILS_TEST_GTK3_MENU") != "1" {
		t.Skip("run by TestGTK3MenuRebuild")
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if !gtkMenuTestInit() {
		t.Fatal("GTK display unavailable; run with xvfb-run -a")
	}
	globalApplication = &App{impl: &gtkMenuTestApp{}, running: true, Logger: slog.Default()}

	menu := NewMenu()
	outer := menu.AddSubmenu("Outer")
	inner := outer.AddSubmenu("Inner")
	item := inner.Add("Action").OnClick(func(*Context) {})
	menu.Update()
	root := menu.impl.(*linuxMenu).native
	gtkMenuTestOwn(root)
	defer func() {
		if root != nil {
			gtkMenuTestRelease(root)
		}
	}()

	assertLabels := func(menu *Menu, want ...string) {
		t.Helper()
		if got := gtkMenuTestLabels(menu); !slices.Equal(got, want) {
			t.Fatalf("native labels = %q, want %q", got, want)
		}
	}
	for i := 0; i < 10; i++ {
		menu.Update()
		assertLabels(menu, "Outer")
		assertLabels(outer, "Inner")
		assertLabels(inner, "Action")
		gtkMenuTestActivate(item)
		select {
		case id := <-menuItemClicked:
			if id != item.id {
				t.Fatalf("activation id = %d, want %d", id, item.id)
			}
		case <-time.After(time.Second):
			t.Fatal("rebuilt menu item did not activate")
		}
	}
	// SetLabel on an attached item must change the native label: each item
	// holds a box (image, label), which gtk_menu_item_set_label ignores.
	menu.FindByLabel("Outer").SetLabel("Renamed")
	assertLabels(menu, "Renamed")
	item.SetLabel("Action renamed")
	assertLabels(inner, "Action renamed")
	item.SetLabel("Action")
	inner.Add("Added")
	menu.Update()
	assertLabels(inner, "Action", "Added")
	inner.RemoveMenuItem(item)
	menu.Update()
	assertLabels(inner, "Added")
	if item.impl != nil {
		t.Fatal("removed item retains a native implementation")
	}
	if _, ok := gtkSignalToMenuItem[item.id]; ok {
		t.Fatal("removed item retains an activation handler")
	}

	// An independently rebuilt submenu must also survive a later parent rebuild.
	inner.Clear()
	inner.Add("After clear")
	inner.Update()
	menu.Update()
	assertLabels(inner, "After clear")

	// Removed branches can remain referenced by application code.
	branch := menu.ItemAt(0)
	menu.RemoveMenuItem(branch)
	menu.Update()
	assertLabels(menu)
	if outer.impl.(*linuxMenu).native != nil || inner.impl.(*linuxMenu).native != nil || branch.impl != nil {
		t.Fatal("removed branch retains destroyed native implementations")
	}
	outer.Update()
	assertLabels(outer, "Inner")
	assertLabels(inner, "After clear")
	menu.items = append(menu.items, branch)
	menu.Update()
	assertLabels(inner, "After clear")

	check := inner.AddCheckbox("Checked", true).OnClick(func(*Context) {})
	first := inner.AddRadio("First", true)
	second := inner.AddRadio("Second", false)
	menu.Update()
	menu.Update()
	if !check.impl.(*linuxMenuItem).isChecked() || !first.impl.(*linuxMenuItem).isChecked() || second.impl.(*linuxMenuItem).isChecked() {
		t.Fatal("rebuild changed check or radio state")
	}
	inner.Add("Callback").OnClick(func(*Context) {})
	menu.Update()
	menu.Destroy()
	menu.Destroy()
	if len(menuDestroyCallbacks) != 1 {
		t.Fatalf("destroy callbacks = %d, want only the root menu", len(menuDestroyCallbacks))
	}
	if len(gtkSignalToMenuItem) != 0 {
		t.Fatal("menu teardown retains activation handlers")
	}
	menu.AddSubmenu("Recreated").Add("Final")
	menu.Update()
	assertLabels(menu, "Recreated")
	// Releasing the root must invalidate every surviving descendant as well.
	gtkMenuTestRelease(root)
	root = nil
	if menu.impl.(*linuxMenu).native != nil || len(menuDestroyCallbacks) != 0 {
		t.Fatal("root destruction retains native implementations or callbacks")
	}
}

// Menu accelerators are GTK accel-group entries on the item's activate signal.
// The keys go through gtk_accel_groups_activate, the same lookup GTK's window
// key handling uses, so a match here is a match for a real key press.
func TestGTK3MenuAccelerator(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGTK3MenuAcceleratorProcess$", "-test.v")
	cmd.Env = append(os.Environ(), "WAILS_TEST_GTK3_ACCEL=1", "G_DEBUG=fatal-criticals")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("native accelerator regression: %v\n%s", err, output)
	}
}

func TestGTK3MenuAcceleratorProcess(t *testing.T) {
	if os.Getenv("WAILS_TEST_GTK3_ACCEL") != "1" {
		t.Skip("run by TestGTK3MenuAccelerator")
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if !gtkMenuTestInit() {
		t.Fatal("GTK display unavailable; run with xvfb-run -a")
	}
	globalApplication = &App{impl: &gtkMenuTestApp{}, running: true, Logger: slog.Default()}

	const (
		ctrl  = uint32(1 << 2)
		shift = uint32(1 << 0)
	)
	menu := NewMenu()
	file := menu.AddSubmenu("File")
	action := file.Add("Action").SetAccelerator("CmdOrCtrl+M").OnClick(func(*Context) {})
	redo := file.Add("Redo").SetAccelerator("CmdOrCtrl+Shift+Z").OnClick(func(*Context) {})
	menu.Update()
	window := gtkMenuTestWindow(pointer(menu.impl.(*linuxMenu).native))

	// press reports whether the chord ran the item, and drains its signal.
	clicked := func(id uint) bool {
		select {
		case got := <-menuItemClicked:
			return got == id
		case <-time.After(200 * time.Millisecond):
			return false
		}
	}
	activates := func(keyval, mods uint32, item *MenuItem) {
		t.Helper()
		if !gtkMenuTestAccelActivate(window, keyval, mods) {
			t.Fatalf("chord %#x/%#x matched no accelerator", keyval, mods)
		}
		if !clicked(item.id) {
			t.Fatalf("chord %#x/%#x did not activate the item", keyval, mods)
		}
	}
	noMatch := func(keyval, mods uint32) {
		t.Helper()
		if gtkMenuTestAccelActivate(window, keyval, mods) {
			t.Fatalf("chord %#x/%#x matched an accelerator it should not", keyval, mods)
		}
	}

	// Window menus: the accelerator fires the item once, and its label shows it.
	activates('m', ctrl, action)
	if w := gtkMenuTestAccelLabelWidth(action); w <= 0 {
		t.Fatalf("accel label shows no shortcut (width %d)", w)
	}
	// GTK reports a shifted letter as its upper-case keysym.
	activates('Z', ctrl|shift, redo)

	// Update rebuilds the items; the new ones must keep the accelerator, and
	// the old ones must not still fire it.
	for i := 0; i < 3; i++ {
		menu.Update()
		activates('m', ctrl, action)
		if n := gtkMenuTestAccelCount(action); n != 1 {
			t.Fatalf("after Update %d the item has %d accelerators, want 1", i, n)
		}
	}

	// SetAccelerator after creation replaces the old chord.
	action.SetAccelerator("CmdOrCtrl+Shift+K")
	noMatch('m', ctrl)
	activates('K', ctrl|shift, action)
	if w := gtkMenuTestAccelLabelWidth(action); w <= 0 {
		t.Fatalf("accel label after SetAccelerator shows no shortcut (width %d)", w)
	}
	if n := gtkMenuTestAccelCount(action); n != 1 {
		t.Fatalf("SetAccelerator left %d accelerators, want 1", n)
	}

	// Clearing the accelerator on the native item removes it.
	gtkMenuTestClearAccelerator(action)
	noMatch('k', ctrl|shift)
	if w := gtkMenuTestAccelLabelWidth(action); w != 0 {
		t.Fatalf("cleared accel label still shows a shortcut (width %d)", w)
	}
	if n := gtkMenuTestAccelCount(action); n != 0 {
		t.Fatalf("cleared item has %d accelerators, want 0", n)
	}

	// A chord that is not bound matches nothing.
	noMatch('q', ctrl)
	noMatch('z', ctrl)

	// RemoveAccelerator is the public way to clear a binding: it must drop the
	// native one, not only the Go field.
	redo.RemoveAccelerator()
	noMatch('Z', ctrl|shift)
	if n := gtkMenuTestAccelCount(redo); n != 0 {
		t.Fatalf("RemoveAccelerator left %d native accelerators, want 0", n)
	}
	if w := gtkMenuTestAccelLabelWidth(redo); w != 0 {
		t.Fatalf("RemoveAccelerator left an accel label (width %d)", w)
	}
}

// TestGTK3SetMenuOnRunningWindow runs the runtime SetMenu checks in their own
// process, as TestGTK3MenuRebuild does.
func TestGTK3SetMenuOnRunningWindow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGTK3SetMenuOnRunningWindowProcess$", "-test.v")
	cmd.Env = append(os.Environ(), "WAILS_TEST_GTK3_SETMENU=1", "G_DEBUG=fatal-criticals")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("runtime SetMenu regression: %v\n%s", err, output)
	}
}

func TestGTK3SetMenuOnRunningWindowProcess(t *testing.T) {
	if os.Getenv("WAILS_TEST_GTK3_SETMENU") != "1" {
		t.Skip("run by TestGTK3SetMenuOnRunningWindow")
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if !gtkMenuTestInit() {
		t.Fatal("GTK display unavailable; run with xvfb-run -a")
	}
	globalApplication = &App{impl: &gtkMenuTestApp{}, running: true, Logger: slog.Default()}

	w := gtkMenuTestRealisedWindow()
	box := w.vbox
	webview := w.webview
	const keyOne, keyTwo, ctrl = 0x31, 0x32, 1 << 2 // GDK_KEY_1, GDK_KEY_2, GDK_CONTROL_MASK

	// expectBar: the box is [menu bar, webview] and the bar shows want.
	expectBar := func(menu *Menu, want ...string) {
		t.Helper()
		kids := gtkMenuTestBoxChildren(box)
		if len(kids) != 2 || kids[1] != webview || kids[0] != w.menubar {
			t.Fatalf("box children = %#x, want [menu bar %#x, webview %#x]", kids, w.menubar, webview)
		}
		var view *Menu
		for candidate := range gtk3MenuViews[menu] {
			if gtkMenuTestNative(candidate) == w.menubar {
				view = candidate
				break
			}
		}
		if view == nil {
			t.Fatal("menu bar has no matching window view")
		}
		if got := gtkMenuTestLabels(view); !slices.Equal(got, want) {
			t.Fatalf("menu bar labels = %q, want %q", got, want)
		}
	}
	expectNoBar := func() {
		t.Helper()
		if kids := gtkMenuTestBoxChildren(box); len(kids) != 1 || kids[0] != webview {
			t.Fatalf("box children = %#x, want only the webview %#x", kids, webview)
		}
	}
	// chord reports whether the window's accelerator for key fires.
	chord := func(key uint32) bool {
		return gtkMenuTestAccelActivate(w.window, key, ctrl)
	}

	alpha := NewMenu()
	alpha.AddSubmenu("Alpha").Add("One").SetAccelerator("ctrl+1")
	alpha.Update()
	beta := NewMenu()
	beta.AddSubmenu("Beta").Add("Two").SetAccelerator("ctrl+2")
	beta.Update()

	expectNoBar()
	w.setMenu(alpha)
	expectBar(alpha, "Alpha")
	if !chord(keyOne) {
		t.Fatal("alpha accelerator is not registered on the running window")
	}

	// Replacing the bar: the new labels show, the old bar is gone, and the
	// replaced menu's accelerator no longer fires.
	w.setMenu(beta)
	expectBar(beta, "Beta")
	if chord(keyOne) {
		t.Fatal("replaced menu's accelerator still fires")
	}
	if !chord(keyTwo) {
		t.Fatal("new menu's accelerator is not registered")
	}

	// Repeated replacement neither adds widgets nor leaves a stale bar.
	for i := 0; i < 10; i++ {
		w.setMenu(alpha)
		expectBar(alpha, "Alpha")
		w.setMenu(beta)
		expectBar(beta, "Beta")
	}

	// nil removes the bar and its accelerators.
	w.setMenu(nil)
	expectNoBar()
	if chord(keyTwo) {
		t.Fatal("removed menu's accelerator still fires")
	}

	// The window takes a menu again after its bar was removed.
	w.setMenu(alpha)
	expectBar(alpha, "Alpha")
	if !chord(keyOne) {
		t.Fatal("alpha accelerator not registered after re-adding the menu")
	}
}

// The application menu is inherited by both windows; removing it from one
// must neither keep its shortcut there nor disable it in the other window.
func TestGTK3InheritedMenuAccelerators(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGTK3InheritedMenuAcceleratorsProcess$", "-test.v")
	cmd.Env = append(os.Environ(), "WAILS_TEST_GTK3_INHERITED=1", "G_DEBUG=fatal-criticals")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("inherited menu: %v\n%s", err, output)
	}
}

func TestGTK3InheritedMenuAcceleratorsProcess(t *testing.T) {
	if os.Getenv("WAILS_TEST_GTK3_INHERITED") != "1" {
		t.Skip("run by TestGTK3InheritedMenuAccelerators")
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if !gtkMenuTestInit() {
		t.Fatal("GTK display unavailable")
	}
	globalApplication = &App{impl: &gtkMenuTestApp{}, running: true, Logger: slog.Default()}
	menu := NewMenu()
	file := menu.AddSubmenu("File")
	file.Add("Action").SetAccelerator("ctrl+1").OnClick(func(*Context) {})
	check := file.AddCheckbox("Check", false)
	file.AddRadio("First", true)
	radio := file.AddRadio("Second", false).SetAccelerator("ctrl+3").OnClick(func(*Context) {})
	menu.Update()
	globalApplication.applicationMenu = menu
	first := gtkMenuTestRealisedWindow()
	second := gtkMenuTestRealisedWindow()
	first.menubar = windowAttachMenuBar(first.window, first.vbox, windowMenuNative(menu), LinuxMenuStyleMenuBar)
	second.menubar = windowAttachMenuBar(second.window, second.vbox, windowMenuNative(menu), LinuxMenuStyleMenuBar)
	if first.menubar == second.menubar {
		t.Fatal("windows share one native menu bar")
	}
	const key, ctrl = 0x31, 1 << 2
	activate := func(window pointer, key uint32) bool {
		matched := gtkMenuTestAccelActivate(window, key, ctrl)
		if matched {
			select {
			case <-menuItemClicked:
			default:
			}
		}
		return matched
	}
	if !activate(second.window, key) {
		t.Fatal("inherited accelerator did not fire in second window")
	}
	menu.FindByLabel("Action").SetLabel("Renamed")
	check.SetChecked(true)
	for view := range gtk3MenuViews[menu] {
		if !view.items[0].submenu.FindByLabel("Check").impl.(*linuxMenuItem).isChecked() {
			t.Fatal("checked state did not propagate")
		}
		if got := gtkMenuTestItemLabel(view.items[0].submenu.items[0]); got != "Renamed" {
			t.Fatalf("native label = %q", got)
		}
	}
	if !gtkMenuTestAccelActivate(second.window, '3', ctrl) {
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
		if !submenu.FindByLabel("Second").impl.(*linuxMenuItem).isChecked() || submenu.FindByLabel("First").impl.(*linuxMenuItem).isChecked() {
			t.Fatal("radio state did not propagate")
		}
	}
	first.setMenu(nil)
	if activate(first.window, key) {
		t.Error("removed window still has the accelerator")
	}
	if !activate(second.window, key) {
		t.Error("other window lost its inherited accelerator")
	}
	first.setMenu(menu)
	if !activate(first.window, key) {
		t.Error("restored menu lost its accelerator")
	}
	menu.FindByLabel("Renamed").SetAccelerator("ctrl+2")
	if activate(first.window, key) || !activate(first.window, '2') || !activate(second.window, '2') {
		t.Error("accelerator update did not propagate")
	}
	menu.items[0].submenu.Add("Added")
	menu.items[0].submenu.Update()
	for view := range gtk3MenuViews[menu] {
		submenu := view.items[0].submenu
		if submenu.FindByLabel("Added") == nil || len(gtkMenuTestBoxChildren(gtkMenuTestNative(submenu))) != 5 {
			t.Fatal("submenu rebuild did not propagate")
		}
	}
	first.setMenu(nil)
	second.setMenu(nil)
	if len(gtk3MenuViews[menu]) != 0 {
		t.Fatal("destroyed windows left menu views behind")
	}
	action := menu.FindByLabel("Renamed")
	if got := len(gtk3ItemViews[action]); got != 1 {
		t.Fatalf("remaining native item views = %d, want the original only", got)
	}
	menu.Destroy()
	if len(gtk3ItemViews) != 0 {
		t.Fatalf("destroyed menu retained %d item view entries", len(gtk3ItemViews))
	}
}
