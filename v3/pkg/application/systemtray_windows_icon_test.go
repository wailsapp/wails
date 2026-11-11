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
