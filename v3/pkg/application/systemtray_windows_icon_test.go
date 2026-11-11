//go:build windows && !server

package application

import (
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/icons"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

// A SetIcon made on the main thread is applied on a goroutine, so an older
// call can reach iconLock after a newer one. The older one must not win.
func TestWindowsSystemTrayDropsAnOlderIconUpdate(t *testing.T) {
	s := &windowsSystemTray{parent: &SystemTray{}}

	s.applyIcon(icons.SystrayLight, 2)
	newest := s.lightModeIcon
	if newest == 0 {
		t.Fatal("the newest icon update was not applied")
	}

	s.applyIcon(icons.SystrayDark, 1)
	if s.lightModeIcon != newest {
		t.Fatal("an older icon update replaced a newer one")
	}
	if s.lightApplied != 2 {
		t.Fatalf("lightApplied = %d, want 2", s.lightApplied)
	}
}

// The same rule for the dark mode icon, which has its own sequence.
func TestWindowsSystemTrayDropsAnOlderDarkModeIconUpdate(t *testing.T) {
	s := &windowsSystemTray{parent: &SystemTray{}}

	s.applyDarkModeIcon(icons.SystrayDark, 2)
	newest := s.darkModeIcon
	s.applyDarkModeIcon(icons.SystrayLight, 1)
	if s.darkModeIcon != newest || s.darkApplied != 2 {
		t.Fatalf("an older dark mode icon update replaced a newer one (applied=%d)", s.darkApplied)
	}
}

// recordDestroyedIcons swaps destroyIcon for a recorder for the test.
func recordDestroyedIcons(t *testing.T) *[]w32.HICON {
	var destroyed []w32.HICON
	original := destroyIcon
	destroyIcon = func(h w32.HICON) bool {
		destroyed = append(destroyed, h)
		return true
	}
	t.Cleanup(func() { destroyIcon = original })
	return &destroyed
}

func countIcon(handles []w32.HICON, h w32.HICON) int {
	n := 0
	for _, x := range handles {
		if x == h {
			n++
		}
	}
	return n
}

// The handle the shell is showing must not be destroyed when SetIcon replaces
// the light icon: updateIcon still needs it, either to roll back to when
// NIM_MODIFY fails or to release once when it succeeds. Destroying it in
// applyIcon as well freed it twice.
func TestWindowsSystemTraySetIconKeepsTheShownIcon(t *testing.T) {
	destroyed := recordDestroyedIcons(t)
	shown := w32.HICON(0x1234) // a stand-in handle; destroyIcon is recorded, not called
	s := &windowsSystemTray{
		parent:             &SystemTray{},
		lightModeIcon:      shown,
		lightModeIconOwned: true,
		currentIcon:        shown,
		currentIconOwned:   true,
	}

	s.applyIcon(icons.SystrayDark, 1)
	if n := countIcon(*destroyed, shown); n != 0 {
		t.Fatalf("the icon the shell is showing was destroyed %d time(s) by SetIcon", n)
	}

	s.releaseIcons()
	if n := countIcon(*destroyed, shown); n != 1 {
		t.Fatalf("the shown icon was destroyed %d times, want once", n)
	}
}

// The same rule for SetDarkModeIcon.
func TestWindowsSystemTraySetDarkModeIconKeepsTheShownIcon(t *testing.T) {
	destroyed := recordDestroyedIcons(t)
	shown := w32.HICON(0x1234)
	s := &windowsSystemTray{
		parent:            &SystemTray{},
		darkModeIcon:      shown,
		darkModeIconOwned: true,
		currentIcon:       shown,
		currentIconOwned:  true,
	}

	s.applyDarkModeIcon(icons.SystrayLight, 1)
	if n := countIcon(*destroyed, shown); n != 0 {
		t.Fatalf("the icon the shell is showing was destroyed %d time(s) by SetDarkModeIcon", n)
	}

	s.releaseIcons()
	if n := countIcon(*destroyed, shown); n != 1 {
		t.Fatalf("the shown icon was destroyed %d times, want once", n)
	}
}

// While a worker is stuck applying one main-thread update (Shell_NotifyIcon
// blocked on Explorer), further updates must not start more workers: only
// the newest is kept, from a copy of the caller's bytes.
func TestQueueIconKeepsOnlyTheNewestUpdateWhileBlocked(t *testing.T) {
	var p pendingIcon
	release := make(chan struct{})
	started := make(chan struct{}, 1)
	type applied struct {
		icon []byte
		seq  uint64
	}
	done := make(chan applied, 10)
	apply := func(icon []byte, seq uint64) {
		select {
		case started <- struct{}{}:
		default:
		}
		<-release
		done <- applied{icon, seq}
	}

	buf := []byte{1}
	queueIcon(&p, buf, 1, apply)
	<-started // the first worker is now blocked in apply

	for seq := uint64(2); seq <= 50; seq++ {
		buf[0] = byte(seq)
		queueIcon(&p, buf, seq, apply)
	}
	buf[0] = 0 // the caller reuses its buffer
	close(release)

	first := <-done
	second := <-done
	if first.seq != 1 || second.seq != 50 {
		t.Fatalf("applied seqs %d, %d; want 1 then 50", first.seq, second.seq)
	}
	if second.icon[0] != 50 {
		t.Fatalf("applied icon byte %d; want a copy holding 50", second.icon[0])
	}
	select {
	case extra := <-done:
		t.Fatalf("an extra update was applied: seq %d", extra.seq)
	case <-time.After(100 * time.Millisecond):
	}
	// The worker clears running after apply returns, so wait for it rather
	// than read it once.
	deadline := time.Now().Add(5 * time.Second)
	for {
		p.mu.Lock()
		running := p.running
		p.mu.Unlock()
		if !running {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the worker did not exit after draining")
		}
		time.Sleep(time.Millisecond)
	}
}

// SystemTray.Destroy reaches the platform destroy on every call. When the
// first call left the cleanup to a goroutine (iconLock busy), hwnd keeps the
// destroyed window's handle, so a second call must return before it passes
// that stale handle to Shell_NotifyIcon and DestroyWindow again.
func TestWindowsSystemTrayDestroyTwiceDoesNothingTheSecondTime(t *testing.T) {
	s := &windowsSystemTray{parent: &SystemTray{}, hwnd: w32.HWND(0x1234)}
	s.destroyed.Store(true) // the state the first destroy leaves

	// Any work past the guard reaches the application; without one it panics,
	// which shows the second call did not return early.
	saved := globalApplication
	globalApplication = nil
	t.Cleanup(func() { globalApplication = saved })
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("a second destroy did work again: %v", r)
		}
	}()

	s.destroy()
}

// runIconUpdate sets up a tray whose shell shows an owned handle in one mode
// (the other mode unset), then replaces that mode's icon with NIM_MODIFY
// answering ok. After the update both modes hold the new handle, so the
// result does not depend on the machine's theme.
func runIconUpdate(t *testing.T, dark bool, ok bool) (s *windowsSystemTray, shown w32.HICON, destroyed *[]w32.HICON, sent []w32.HICON) {
	t.Helper()
	destroyed = recordDestroyedIcons(t)

	original := shellNotifyIcon
	shellNotifyIcon = func(msg uintptr, nid *w32.NOTIFYICONDATA) bool {
		if msg == w32.NIM_MODIFY {
			sent = append(sent, nid.HIcon)
		}
		return ok
	}
	t.Cleanup(func() { shellNotifyIcon = original })

	savedApp := globalApplication
	globalApplication = &App{options: Options{WarningHandler: func(string) {}}}
	t.Cleanup(func() { globalApplication = savedApp })

	shown = w32.HICON(0x1234)
	s = &windowsSystemTray{
		parent:           &SystemTray{},
		hwnd:             w32.HWND(0x1),
		currentIcon:      shown,
		currentIconOwned: true,
	}
	if dark {
		s.darkModeIcon, s.darkModeIconOwned = shown, true
		s.applyDarkModeIcon(icons.SystrayLight, 1)
	} else {
		s.lightModeIcon, s.lightModeIconOwned = shown, true
		s.applyIcon(icons.SystrayDark, 1)
	}
	return s, shown, destroyed, sent
}

// When NIM_MODIFY succeeds the shell shows the new icon and the old one is
// released exactly once.
func TestWindowsSystemTrayIconUpdateReleasesTheOldIconOnceOnSuccess(t *testing.T) {
	for _, dark := range []bool{false, true} {
		s, shown, destroyed, sent := runIconUpdate(t, dark, true)
		if len(sent) != 1 || sent[0] != s.currentIcon || s.currentIcon == shown {
			t.Fatalf("dark=%v: NIM_MODIFY sent %v, currentIcon %v; want the new icon", dark, sent, s.currentIcon)
		}
		if !s.currentIconOwned {
			t.Fatalf("dark=%v: the new current icon is not marked owned", dark)
		}
		if n := countIcon(*destroyed, shown); n != 1 {
			t.Fatalf("dark=%v: the old icon was destroyed %d times, want once", dark, n)
		}
		if n := countIcon(*destroyed, s.currentIcon); n != 0 {
			t.Fatalf("dark=%v: the icon now shown was destroyed", dark)
		}
	}
}

// When NIM_MODIFY fails (the icon is not registered, e.g. while Explorer
// restarts) the tray rolls back to the old icon, which must still be valid,
// and releases each handle exactly once when the tray is torn down.
func TestWindowsSystemTrayIconUpdateKeepsTheOldIconOnFailure(t *testing.T) {
	for _, dark := range []bool{false, true} {
		s, shown, destroyed, sent := runIconUpdate(t, dark, false)
		if len(sent) != 1 || sent[0] == shown {
			t.Fatalf("dark=%v: NIM_MODIFY sent %v; want one attempt with the new icon", dark, sent)
		}
		newIcon := sent[0]
		if s.currentIcon != shown || !s.currentIconOwned {
			t.Fatalf("dark=%v: currentIcon %v (owned %v); want the rollback to the old icon", dark, s.currentIcon, s.currentIconOwned)
		}
		if n := countIcon(*destroyed, shown); n != 0 {
			t.Fatalf("dark=%v: the icon still shown was destroyed %d time(s)", dark, n)
		}

		s.releaseIcons()
		if n := countIcon(*destroyed, shown); n != 1 {
			t.Fatalf("dark=%v: the old icon was destroyed %d times at teardown, want once", dark, n)
		}
		if n := countIcon(*destroyed, newIcon); n != 1 {
			t.Fatalf("dark=%v: the new icon was destroyed %d times at teardown, want once", dark, n)
		}
	}
}

// SystemTray.Run sets the implementation and then queues run on the main
// thread. A Destroy that lands in between must stop run from creating the
// window and adding an icon nobody can remove.
func TestWindowsSystemTrayRunAfterDestroyDoesNothing(t *testing.T) {
	s := &windowsSystemTray{parent: &SystemTray{}}
	s.destroyed.Store(true)

	// Creating the window reads the application options; without a guard this
	// panics, which shows run did not return early.
	saved := globalApplication
	globalApplication = nil
	t.Cleanup(func() { globalApplication = saved })
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("run went ahead on a destroyed tray: %v", r)
		}
	}()

	s.run()
	if s.hwnd != 0 {
		t.Fatal("run created a window for a destroyed tray")
	}
	if !s.iconLock.TryLock() {
		t.Fatal("run left iconLock held")
	}
	s.iconLock.Unlock()
}
