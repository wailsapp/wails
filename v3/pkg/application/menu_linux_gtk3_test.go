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
		default:
			t.Fatal("rebuilt menu item did not activate")
		}
	}
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
