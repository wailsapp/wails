//go:build linux && !android && !gtk3 && !server

package application

import "testing"

type systrayUpdateTestApp struct{ platformApp }

func (*systrayUpdateTestApp) isOnMainThread() bool { return true }

// A tray's dbusmenu layout is built from the Menu by setMenu, so Menu.Update
// must rebuild it as well or the tray keeps showing the first menu forever.
func TestLinuxSystemTrayMenuUpdate(t *testing.T) {
	prev := globalApplication
	globalApplication = &App{
		impl:        &systrayUpdateTestApp{},
		running:     true,
		systemTrays: map[uint]*SystemTray{},
	}
	t.Cleanup(func() { globalApplication = prev })

	menu := NewMenu()
	menu.Add("First")
	tray := &SystemTray{id: 1, menu: menu}
	impl := &linuxSystemTray{parent: tray}
	impl.menuVersion.Store(1)
	tray.impl = impl
	globalApplication.systemTrays[tray.id] = tray
	impl.setMenu(menu)
	revision := impl.menuVersion.Load()

	menu.Clear()
	added := menu.Add("Second")
	menu.Update()

	if _, ok := impl.itemMap[int32(added.id)]; !ok {
		t.Fatalf("tray itemMap %v lacks item %d added before Menu.Update", trayItemIDs(impl), added.id)
	}
	if got := impl.menuVersion.Load(); got <= revision {
		t.Fatalf("menu revision stayed at %d after Menu.Update", got)
	}
	// Item setters must reach the tray, not the native menu Update also built.
	if _, ok := added.impl.(*systrayMenuItem); !ok {
		t.Fatalf("added item impl is %T, want *systrayMenuItem", added.impl)
	}
}

func trayItemIDs(s *linuxSystemTray) []int32 {
	ids := make([]int32, 0, len(s.itemMap))
	for id := range s.itemMap {
		ids = append(ids, id)
	}
	return ids
}
