package application

import (
	"testing"
	"time"
)

// newTestGlobalShortcutManager returns a manager that is safe to use without a
// running application: because the app is never started, registrations take the
// "pending" path and never invoke the (absent) platform implementation.
func newTestGlobalShortcutManager() *GlobalShortcutManager {
	return newGlobalShortcutManager(&App{})
}

func TestGlobalShortcutRegisterAndQuery(t *testing.T) {
	m := newTestGlobalShortcutManager()

	if err := m.Register("Ctrl+Shift+A", func() {}); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	if !m.IsRegistered("Ctrl+Shift+A") {
		t.Fatal("expected Ctrl+Shift+A to be registered")
	}
	// Modifier order should not matter: this is the same shortcut.
	if !m.IsRegistered("Shift+Ctrl+A") {
		t.Fatal("expected Shift+Ctrl+A to match Ctrl+Shift+A")
	}
	if m.IsRegistered("Ctrl+Shift+B") {
		t.Fatal("did not expect Ctrl+Shift+B to be registered")
	}
}

func TestGlobalShortcutDuplicateIsRejected(t *testing.T) {
	m := newTestGlobalShortcutManager()

	if err := m.Register("Ctrl+Shift+A", func() {}); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	// Registering the same accelerator again (in any modifier order) must be
	// rejected and must not overwrite the original binding.
	if err := m.Register("Shift+Ctrl+A", func() {}); err == nil {
		t.Fatal("expected duplicate registration to be rejected")
	}
	if got := m.GetAll(); len(got) != 1 {
		t.Fatalf("expected exactly one shortcut after duplicate, got %v", got)
	}
}

func TestGlobalShortcutInvalidInputs(t *testing.T) {
	m := newTestGlobalShortcutManager()

	if err := m.Register("Ctrl+A", nil); err == nil {
		t.Fatal("expected error for nil callback")
	}
	if err := m.Register("Ctrl+NotAKey", func() {}); err == nil {
		t.Fatal("expected error for invalid accelerator")
	}
	if m.IsRegistered("Ctrl+NotAKey") {
		t.Fatal("invalid accelerator should never be reported as registered")
	}
}

func TestGlobalShortcutUnregister(t *testing.T) {
	m := newTestGlobalShortcutManager()

	_ = m.Register("Ctrl+Shift+A", func() {})
	_ = m.Register("Ctrl+Shift+B", func() {})

	if err := m.Unregister("Ctrl+Shift+A"); err != nil {
		t.Fatalf("unregister failed: %v", err)
	}
	if m.IsRegistered("Ctrl+Shift+A") {
		t.Fatal("expected Ctrl+Shift+A to be unregistered")
	}
	if !m.IsRegistered("Ctrl+Shift+B") {
		t.Fatal("expected Ctrl+Shift+B to remain registered")
	}
	if err := m.Unregister("Ctrl+Shift+A"); err == nil {
		t.Fatal("expected error unregistering an absent shortcut")
	}

	// The pending queue should no longer reference the removed shortcut.
	if len(m.pending) != 1 {
		t.Fatalf("expected one pending shortcut, got %d", len(m.pending))
	}
}

func TestGlobalShortcutGetAllSorted(t *testing.T) {
	m := newTestGlobalShortcutManager()

	_ = m.Register("Ctrl+Shift+C", func() {})
	_ = m.Register("Ctrl+Shift+A", func() {})
	_ = m.Register("Ctrl+Shift+B", func() {})

	all := m.GetAll()
	if len(all) != 3 {
		t.Fatalf("expected 3 shortcuts, got %d", len(all))
	}
	for i := 1; i < len(all); i++ {
		if all[i-1] > all[i] {
			t.Fatalf("GetAll not sorted: %v", all)
		}
	}
}

func TestGlobalShortcutDispatch(t *testing.T) {
	m := newTestGlobalShortcutManager()

	fired := make(chan struct{}, 1)
	if err := m.Register("Ctrl+Shift+A", func() { fired <- struct{}{} }); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Find the id assigned to the shortcut and dispatch it as the platform
	// layer would.
	m.mu.Lock()
	var id int
	for sid := range m.byID {
		id = sid
	}
	m.mu.Unlock()

	m.dispatch(id)
	select {
	case <-fired:
	case <-time.After(2 * time.Second):
		t.Fatal("callback did not fire after dispatch")
	}

	// Dispatching an unknown id must be a no-op (and must not panic).
	m.dispatch(99999)
}
