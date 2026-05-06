//go:build linux && !android && !server

package application

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Drives setMenu in a tight loop against the dbusmenu callbacks that the
// godbus worker goroutine would dispatch in production. Without a lock on
// itemMap the runtime aborts with "concurrent map read and map write".
func TestLinuxSystemTrayConcurrentSetMenu(t *testing.T) {
	tray := &linuxSystemTray{}
	tray.menuVersion.Store(1)
	tray.setMenu(buildSystrayRaceMenu(0))

	const readers = 4
	deadline := time.Now().Add(500 * time.Millisecond)

	var iters atomic.Uint64
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; time.Now().Before(deadline); i++ {
			tray.setMenu(buildSystrayRaceMenu(i))
			iters.Add(1)
		}
	}()

	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ids := []int32{0, 1, 2, 3, 4, 5, 6, 7}
			for time.Now().Before(deadline) {
				_, _, _ = tray.GetLayout(0, -1, nil)
				_, _ = tray.GetGroupProperties(ids, nil)
				_, _ = tray.GetProperty(0, "label")
			}
		}()
	}

	wg.Wait()

	if iters.Load() == 0 {
		t.Fatalf("writer made no progress")
	}

	rev, layout, _ := tray.GetLayout(0, -1, nil)
	if rev == 0 {
		t.Fatalf("expected non-zero menu revision after %d setMenu calls", iters.Load())
	}
	if layout.V0 != 0 {
		t.Fatalf("expected root layout id 0, got %d", layout.V0)
	}
}

// Same race against a stable menu: the systrayMenuItem setters write
// dbusItem.V1 directly, so reads on the panel side need the lock too.
func TestLinuxSystemTrayConcurrentItemSetters(t *testing.T) {
	tray := &linuxSystemTray{}
	tray.menuVersion.Store(1)

	m := NewMenu()
	item1 := m.Add("label")
	item2 := m.AddCheckbox("check", false)
	tray.setMenu(m)

	sm1, ok := item1.impl.(*systrayMenuItem)
	if !ok {
		t.Fatalf("item1 has no systrayMenuItem impl")
	}
	sm2, ok := item2.impl.(*systrayMenuItem)
	if !ok {
		t.Fatalf("item2 has no systrayMenuItem impl")
	}

	deadline := time.Now().Add(300 * time.Millisecond)
	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		i := 0
		for time.Now().Before(deadline) {
			sm1.setLabel(fmt.Sprintf("label-%d", i))
			i++
		}
	}()
	go func() {
		defer wg.Done()
		i := 0
		for time.Now().Before(deadline) {
			sm2.setChecked(i%2 == 0)
			i++
		}
	}()
	go func() {
		defer wg.Done()
		i := 0
		for time.Now().Before(deadline) {
			sm1.setDisabled(i%2 == 0)
			i++
		}
	}()
	ids := []int32{int32(item1.id), int32(item2.id)}
	go func() {
		defer wg.Done()
		for time.Now().Before(deadline) {
			_, _, _ = tray.GetLayout(0, -1, nil)
		}
	}()
	go func() {
		defer wg.Done()
		for time.Now().Before(deadline) {
			_, _ = tray.GetGroupProperties(ids, nil)
		}
	}()

	wg.Wait()
}

func buildSystrayRaceMenu(seed int) *Menu {
	m := NewMenu()
	m.Add(fmt.Sprintf("item-%d", seed))
	m.AddCheckbox("checked", seed%2 == 0)
	m.AddSeparator()
	sub := m.AddSubmenu("more")
	sub.Add(fmt.Sprintf("nested-%d", seed))
	return m
}
