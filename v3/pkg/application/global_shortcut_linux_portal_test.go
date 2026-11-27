//go:build linux && !android && !server

package application

import (
	"bufio"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
)

// The portal backend is what Wayland sessions use for global shortcuts. These
// tests run it against a fake org.freedesktop.portal.GlobalShortcuts on a
// private bus: the fake answers each request the way xdg-desktop-portal does
// (a Request object path back, then a Response signal on it) and presses
// shortcuts by emitting Activated.

// privateSessionBus starts a dbus-daemon for the test and points the session
// bus address at it, so nothing reaches the desktop's own bus.
func privateSessionBus(t *testing.T) {
	t.Helper()
	daemon, err := exec.LookPath("dbus-daemon")
	if err != nil {
		t.Skip("dbus-daemon is not installed")
	}
	cmd := exec.Command(daemon, "--session", "--nofork", "--print-address=1")
	out, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Skipf("dbus-daemon: %v", err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})
	address, err := bufio.NewReader(out).ReadString('\n')
	if err != nil {
		t.Fatalf("dbus-daemon printed no address: %v", err)
	}
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", strings.TrimSpace(address))
}

type portalBinding struct {
	ID    string
	Props map[string]dbus.Variant
}

// portalBehaviour is how the fake portal answers.
type portalBehaviour struct {
	createCode uint32 // CreateSession's response code: 0 granted
	bindCode   uint32 // BindShortcuts's response code
	moveHandle bool   // return a request path other than the predicted one
}

// fakePortal is the portal side of the GlobalShortcuts interface.
type fakePortal struct {
	portalBehaviour
	conn  *dbus.Conn
	binds chan []portalBinding

	mu      sync.Mutex
	session dbus.ObjectPath
}

func newFakePortal(t *testing.T, behaviour portalBehaviour) *fakePortal {
	t.Helper()
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	f := &fakePortal{portalBehaviour: behaviour, conn: conn, binds: make(chan []portalBinding, 8)}
	if err := conn.Export(f, portalPath, portalShortcutIf); err != nil {
		t.Fatal(err)
	}
	if reply, err := conn.RequestName(portalService, dbus.NameFlagDoNotQueue); err != nil || reply != dbus.RequestNameReplyPrimaryOwner {
		t.Fatalf("claiming %s: %v, %v", portalService, reply, err)
	}
	return f
}

// request returns the Request object path for a call and, shortly after,
// emits its Response, as the portal does.
func (f *fakePortal) request(sender dbus.Sender, options map[string]dbus.Variant, code uint32, results map[string]dbus.Variant) dbus.ObjectPath {
	token, _ := options["handle_token"].Value().(string)
	path := "/org/freedesktop/portal/desktop/request/" +
		strings.ReplaceAll(strings.TrimPrefix(string(sender), ":"), ".", "_") + "/" + token
	if f.moveHandle {
		path += "_moved"
	}
	go func() {
		// The caller registers a moved path only once the call returns.
		time.Sleep(50 * time.Millisecond)
		_ = f.conn.Emit(dbus.ObjectPath(path), portalRequestIf+".Response", code, results)
	}()
	return dbus.ObjectPath(path)
}

func (f *fakePortal) CreateSession(sender dbus.Sender, options map[string]dbus.Variant) (dbus.ObjectPath, *dbus.Error) {
	token, _ := options["session_handle_token"].Value().(string)
	session := dbus.ObjectPath("/org/freedesktop/portal/desktop/session/e2e/" + token)
	f.mu.Lock()
	f.session = session
	f.mu.Unlock()
	return f.request(sender, options, f.createCode, map[string]dbus.Variant{
		"session_handle": dbus.MakeVariant(string(session)),
	}), nil
}

func (f *fakePortal) BindShortcuts(sender dbus.Sender, session dbus.ObjectPath, shortcuts []portalBinding, parent string, options map[string]dbus.Variant) (dbus.ObjectPath, *dbus.Error) {
	f.binds <- shortcuts
	return f.request(sender, options, f.bindCode, map[string]dbus.Variant{}), nil
}

// press emits Activated for the shortcut the portal knows as id.
func (f *fakePortal) press(id string) error {
	f.mu.Lock()
	session := f.session
	f.mu.Unlock()
	return f.conn.Emit(portalPath, portalShortcutIf+".Activated", session, id, uint64(0), map[string]dbus.Variant{})
}

func (f *fakePortal) awaitBind(t *testing.T) []portalBinding {
	t.Helper()
	select {
	case b := <-f.binds:
		return b
	case <-time.After(5 * time.Second):
		t.Fatal("the backend did not call BindShortcuts")
		return nil
	}
}

// newPortalManager returns a started manager on the portal backend, as on a
// Wayland session, and the channel its errors arrive on.
func newPortalManager(t *testing.T, register func(m *GlobalShortcutManager)) (*GlobalShortcutManager, chan error) {
	t.Helper()
	t.Setenv("XDG_SESSION_TYPE", "wayland")
	errs := make(chan error, 8)
	m := newGlobalShortcutManager(&App{options: Options{ErrorHandler: func(err error) { errs <- err }}})
	register(m)
	m.flushPending()
	if _, ok := m.impl.(*portalGlobalShortcuts); !ok {
		t.Fatalf("a Wayland session got the %T backend, want the portal", m.impl)
	}
	return m, errs
}

func TestPortalGlobalShortcutsBindAndActivate(t *testing.T) {
	privateSessionBus(t)
	portal := newFakePortal(t, portalBehaviour{})
	pressed := make(chan string, 8)
	m, errs := newPortalManager(t, func(m *GlobalShortcutManager) {
		for _, accel := range []string{"CmdOrCtrl+Shift+A", "Alt+Super+F5"} {
			accel := accel
			if err := m.Register(accel, func() { pressed <- accel }); err != nil {
				t.Fatal(err)
			}
		}
	})

	// The startup registrations arrive as one BindShortcuts.
	bound := portal.awaitBind(t)
	triggers := map[string]string{}
	for _, b := range bound {
		trigger, _ := b.Props["preferred_trigger"].Value().(string)
		triggers[b.ID] = trigger
		if _, ok := b.Props["description"]; !ok {
			t.Errorf("shortcut %s has no description", b.ID)
		}
	}
	want := map[string]string{"0": "CTRL+SHIFT+a", "1": "ALT+LOGO+F5"}
	for id, trigger := range want {
		if triggers[id] != trigger {
			t.Errorf("shortcut %s: preferred_trigger %q, want %q (all: %v)", id, triggers[id], trigger, triggers)
		}
	}

	// The compositor reports a press: the callback runs.
	if err := portal.press("1"); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-pressed:
		if got != "Alt+Super+F5" {
			t.Errorf("pressing shortcut 1 ran %q's callback", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Activated did not run the callback")
	}
	// An id the app never bound does nothing.
	if err := portal.press("99"); err != nil {
		t.Fatal(err)
	}

	// Unregistering binds the set that is left.
	if err := m.impl.unregister(0); err != nil {
		t.Fatal(err)
	}
	if left := portal.awaitBind(t); len(left) != 1 || left[0].ID != "1" {
		t.Errorf("after unregistering 0, BindShortcuts got %v", left)
	}
	if err := m.impl.unregisterAll(); err != nil {
		t.Fatal(err)
	}
	if left := portal.awaitBind(t); len(left) != 0 {
		t.Errorf("after unregisterAll, BindShortcuts got %v", left)
	}
	select {
	case got := <-pressed:
		t.Errorf("unexpected callback %q", got)
	case err := <-errs:
		t.Errorf("unexpected error: %v", err)
	default:
	}
}

func TestPortalGlobalShortcutsMovedRequestHandle(t *testing.T) {
	// A portal may return a Request path other than the predicted one; the
	// Response arrives on the path it returned.
	privateSessionBus(t)
	portal := newFakePortal(t, portalBehaviour{moveHandle: true})
	_, errs := newPortalManager(t, func(m *GlobalShortcutManager) {
		if err := m.Register("Ctrl+B", func() {}); err != nil {
			t.Fatal(err)
		}
	})
	portal.awaitBind(t)
	select {
	case err := <-errs:
		t.Fatalf("unexpected error: %v", err)
	case <-time.After(300 * time.Millisecond):
	}
}

func TestPortalGlobalShortcutsRefused(t *testing.T) {
	for _, c := range []struct {
		name       string
		createCode uint32
		bindCode   uint32
		want       string
	}{
		{"the user cancels the bind dialog", 0, 1, "did not grant shortcuts (response 1)"},
		{"the portal refuses a session", 2, 0, "CreateSession refused (response 2)"},
	} {
		t.Run(c.name, func(t *testing.T) {
			privateSessionBus(t)
			newFakePortal(t, portalBehaviour{createCode: c.createCode, bindCode: c.bindCode})
			_, errs := newPortalManager(t, func(m *GlobalShortcutManager) {
				if err := m.Register("Ctrl+B", func() {}); err != nil {
					t.Fatal(err)
				}
			})
			select {
			case err := <-errs:
				if !strings.Contains(err.Error(), c.want) {
					t.Errorf("error %q, want it to say %q", err, c.want)
				}
			case <-time.After(5 * time.Second):
				t.Fatalf("no error reported; want %q", c.want)
			}
		})
	}
}

func TestPortalGlobalShortcutsNoPortal(t *testing.T) {
	// No portal on the bus: the failure is reported, not hung on.
	privateSessionBus(t)
	_, errs := newPortalManager(t, func(m *GlobalShortcutManager) {
		if err := m.Register("Ctrl+B", func() {}); err != nil {
			t.Fatal(err)
		}
	})
	select {
	case err := <-errs:
		if !strings.Contains(err.Error(), "portal session failed") {
			t.Errorf("error %q, want a portal session failure", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("no error reported without a portal")
	}
}
