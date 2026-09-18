//go:build windows && !server

package application

import (
	"testing"
	"unicode/utf16"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/w32"
)

type nativeMenuItemSnapshot struct {
	label   string
	id      uint32
	state   uint32
	submenu w32.HMENU
}

func snapshotNativeMenuItem(t *testing.T, menu w32.HMENU, position uint32) nativeMenuItemSnapshot {
	t.Helper()
	buffer := make([]uint16, 128)
	info := w32.MENUITEMINFO{
		CbSize:     uint32(unsafe.Sizeof(w32.MENUITEMINFO{})),
		FMask:      w32.MIIM_STRING | w32.MIIM_ID | w32.MIIM_STATE | w32.MIIM_SUBMENU,
		DwTypeData: &buffer[0],
		Cch:        uint32(len(buffer)),
	}
	if !w32.GetMenuItemInfo(menu, position, true, &info) {
		t.Fatalf("GetMenuItemInfo(position=%d) failed", position)
	}
	return nativeMenuItemSnapshot{
		label:   string(utf16.Decode(buffer[:info.Cch])),
		id:      info.WID,
		state:   info.FState,
		submenu: info.HSubMenu,
	}
}

func TestAssignCommandIDRejectsMenuWithoutItems(t *testing.T) {
	menu := w32.NewPopupMenu()
	if menu == 0 {
		t.Fatal("NewPopupMenu failed")
	}
	defer w32.DestroyMenu(menu)

	if err := assignCommandID(menu, 1); err == nil {
		t.Fatal("assignCommandID succeeded for an empty menu")
	}
}

func TestWindowsSubmenuParentRuntimeMutations(t *testing.T) {
	builders := []struct {
		name  string
		build func(*Menu) (w32.HMENU, func())
	}{
		{
			name: "windowsMenu",
			build: func(menu *Menu) (w32.HMENU, func()) {
				builder := newMenuImpl(menu)
				builder.update()
				return builder.hMenu, func() {
					builder.detachedSubmenus.destroyAll()
					w32.DestroyMenu(builder.hMenu)
				}
			},
		},
		{
			name: "Win32Menu",
			build: func(menu *Menu) (w32.HMENU, func()) {
				builder := NewPopupMenu(0, menu)
				return builder.menu, builder.Destroy
			},
		},
	}

	for _, builder := range builders {
		t.Run(builder.name, func(t *testing.T) {
			menu := NewMenu()
			submenu := menu.AddSubmenu("Parent")
			submenu.Add("Child")
			parentItem := menu.ItemAt(0)

			nativeMenu, cleanup := builder.build(menu)
			t.Cleanup(cleanup)

			before := snapshotNativeMenuItem(t, nativeMenu, 0)
			expectedID := uint32(parentItem.impl.(*windowsMenuItem).id)
			if before.id != expectedID {
				t.Fatalf("native submenu parent ID = %d, want %d", before.id, expectedID)
			}

			parentItem.SetLabel("Updated")
			parentItem.SetEnabled(false)
			parentItem.SetChecked(true)
			afterMutation := snapshotNativeMenuItem(t, nativeMenu, 0)
			if afterMutation.label != "Updated" {
				t.Errorf("label = %q, want Updated", afterMutation.label)
			}
			if afterMutation.state&w32.MFS_DISABLED == 0 {
				t.Error("submenu parent is not disabled")
			}
			if afterMutation.state&w32.MFS_CHECKED == 0 {
				t.Error("submenu parent is not checked")
			}

			parentItem.SetHidden(true)
			if count := w32.GetMenuItemCount(nativeMenu); count != 0 {
				t.Fatalf("menu item count after hide = %d, want 0", count)
			}
			parentItem.SetHidden(false)
			if count := w32.GetMenuItemCount(nativeMenu); count != 1 {
				t.Fatalf("menu item count after show = %d, want 1", count)
			}
			restored := snapshotNativeMenuItem(t, nativeMenu, 0)
			if restored.submenu != before.submenu {
				t.Errorf("restored submenu handle = %#x, want %#x", restored.submenu, before.submenu)
			}
		})
	}
}

func TestWindowsMenuBuildersReleaseDetachedSubmenus(t *testing.T) {
	builders := []struct {
		name  string
		build func(*Menu) (rebuild func(), destroy func())
	}{
		{
			name: "windowsMenu",
			build: func(menu *Menu) (func(), func()) {
				builder := newMenuImpl(menu)
				builder.update()
				return builder.update, func() {
					builder.detachedSubmenus.destroyAll()
					w32.DestroyMenu(builder.hMenu)
				}
			},
		},
		{
			name: "Win32Menu",
			build: func(menu *Menu) (func(), func()) {
				builder := NewPopupMenu(0, menu)
				return builder.Update, builder.Destroy
			},
		},
	}

	for _, builder := range builders {
		t.Run(builder.name+"/initially hidden", func(t *testing.T) {
			menu := NewMenu()
			submenu := menu.AddSubmenu("Parent")
			submenu.Add("Child")
			parentItem := menu.ItemAt(0).SetHidden(true)

			rebuild, destroy := builder.build(menu)
			t.Cleanup(destroy)
			detached := parentItem.impl.(*windowsMenuItem).submenu
			if count := w32.GetMenuItemCount(detached); count != 1 {
				t.Fatalf("detached submenu item count = %d, want 1", count)
			}

			rebuild()
			if count := w32.GetMenuItemCount(detached); count != -1 {
				t.Errorf("old detached submenu remains valid after rebuild: count=%d", count)
			}
		})

		t.Run(builder.name+"/hidden at runtime", func(t *testing.T) {
			menu := NewMenu()
			submenu := menu.AddSubmenu("Parent")
			submenu.Add("Child")
			parentItem := menu.ItemAt(0)

			rebuild, destroy := builder.build(menu)
			t.Cleanup(destroy)
			detached := parentItem.impl.(*windowsMenuItem).submenu
			parentItem.SetHidden(true)

			rebuild()
			if count := w32.GetMenuItemCount(detached); count != -1 {
				t.Errorf("runtime-detached submenu remains valid after rebuild: count=%d", count)
			}
		})
	}
}

func TestWin32MenuDestroyReleasesDetachedSubmenu(t *testing.T) {
	menu := NewMenu()
	submenu := menu.AddSubmenu("Parent")
	submenu.Add("Child")
	parentItem := menu.ItemAt(0).SetHidden(true)

	builder := NewPopupMenu(0, menu)
	detached := parentItem.impl.(*windowsMenuItem).submenu
	if count := w32.GetMenuItemCount(detached); count != 1 {
		t.Fatalf("detached submenu item count = %d, want 1", count)
	}

	builder.Destroy()
	if count := w32.GetMenuItemCount(detached); count != -1 {
		t.Errorf("detached submenu remains valid after destroy: count=%d", count)
	}
}
