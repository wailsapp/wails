//go:build linux && !android && !server

package application

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// resolveSystrayMenuItem returns the *systrayMenuItem impl bound to item by
// linuxSystemTray.setMenu, failing the test (with name in the message) if the
// type assertion fails.
func resolveSystrayMenuItem(t *testing.T, item *MenuItem, name string) *systrayMenuItem {
	t.Helper()
	sm, ok := item.impl.(*systrayMenuItem)
	if !ok {
		t.Fatalf("%s has no systrayMenuItem impl", name)
	}
	return sm
}

// Drives setMenu in a tight loop against the dbusmenu callbacks that the
// godbus worker goroutine would dispatch in production. Without a lock on
// itemMap the runtime aborts with "concurrent map read and map write".
func TestLinuxSystemTrayConcurrentSetMenu(t *testing.T) {
	tray := &linuxSystemTray{parent: &SystemTray{}}
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

	// SecondaryActivate reads s.menu — race that read against setMenu's
	// write. A single goroutine is enough; we only need the s.menu access
	// pattern to be observable to -race, not to stress concurrent
	// SecondaryActivate calls (which would surface a separate, pre-existing
	// race on lastClickX/Y that is out of scope for this PR).
	wg.Add(1)
	go func() {
		defer wg.Done()
		for time.Now().Before(deadline) {
			_ = tray.SecondaryActivate(0, 0)
		}
	}()

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

	sm1 := resolveSystrayMenuItem(t, item1, "item1")
	sm2 := resolveSystrayMenuItem(t, item2, "item2")

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
			// Mimic godbus's post-return reflection-based serialisation:
			// after GetLayout releases its read lock, the worker goroutine
			// iterates V1 and the V2 children. If GetLayout aliases the
			// live map, this iteration races setLabel / setChecked.
			_, layout, _ := tray.GetLayout(0, -1, nil)
			for k := range layout.V1 {
				_ = k
			}
			for _, child := range layout.V2 {
				if cm, ok := child.Value().(*dbusMenu); ok {
					for k := range cm.V1 {
						_ = k
					}
				}
			}
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

// Verifies that GetLayout never returns a (layout, revision) pair where the
// revision lags behind the layout that produced it. If refreshLocked were
// moved out of the writer's locked region (the "menuVersion bumped after
// Unlock" regression), a reader could RLock between the V1 mutation and the
// version bump, observing a label that the writer just wrote together with
// the previous revision number — breaking the monotonic-revision contract
// that dbusmenu clients use to invalidate their layout cache.
//
// The check works because we run a single writer that calls setLabel(i) in
// strict sequence; refreshLocked then bumps menuVersion by exactly one per
// call, so the revision after the i-th call is `2 + i` (initial Store(1) +
// the bump from the priming setMenu = 2, then +1 per setLabel). A reader
// observing label "label-i" must therefore see revision >= 2 + i.
func TestLinuxSystemTrayLayoutRevisionMonotonic(t *testing.T) {
	tray := &linuxSystemTray{}
	tray.menuVersion.Store(1)

	m := NewMenu()
	item := m.Add("label-0")
	tray.setMenu(m)

	sm := resolveSystrayMenuItem(t, item, "item")

	deadline := time.Now().Add(300 * time.Millisecond)
	var wg sync.WaitGroup
	var failures atomic.Uint64
	var firstFailure atomic.Value // stores string

	writerDone := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(writerDone)
		for i := 1; time.Now().Before(deadline); i++ {
			sm.setLabel(fmt.Sprintf("label-%d", i))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-writerDone:
				return
			default:
			}
			rev, layout, _ := tray.GetLayout(0, -1, nil)
			for _, child := range layout.V2 {
				cm, ok := child.Value().(*dbusMenu)
				if !ok || cm.V0 != int32(item.id) {
					continue
				}
				lblVar, ok := cm.V1["label"]
				if !ok {
					continue
				}
				lbl, ok := lblVar.Value().(string)
				if !ok {
					continue
				}
				var n int
				if _, err := fmt.Sscanf(lbl, "label-%d", &n); err != nil || n == 0 {
					continue
				}
				expected := uint32(2 + n)
				if rev < expected {
					failures.Add(1)
					firstFailure.CompareAndSwap(nil,
						fmt.Sprintf("label=%q rev=%d expected>=%d", lbl, rev, expected))
				}
			}
		}
	}()
	wg.Wait()

	if n := failures.Load(); n > 0 {
		t.Fatalf("non-atomic (layout, revision) snapshot detected %d times; first: %v",
			n, firstFailure.Load())
	}
}
