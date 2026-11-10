//go:build linux && !android

package application

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

// compositorEnv clears every variable the compositor detection reads, then
// sets the given ones.
func compositorEnv(t *testing.T, set map[string]string) {
	t.Helper()
	for _, name := range []string{"HYPRLAND_INSTANCE_SIGNATURE", "SWAYSOCK", "I3SOCK", "XDG_CURRENT_DESKTOP", "XDG_SESSION_TYPE", "WAYLAND_DISPLAY"} {
		t.Setenv(name, "")
		os.Unsetenv(name)
	}
	for name, value := range set {
		t.Setenv(name, value)
	}
}

func TestDetectCompositor(t *testing.T) {
	for _, c := range []struct {
		env    map[string]string
		want   string
		tiling bool
	}{
		{map[string]string{"HYPRLAND_INSTANCE_SIGNATURE": "abc", "SWAYSOCK": "/s"}, "hyprland", true},
		{map[string]string{"SWAYSOCK": "/s", "I3SOCK": "/i"}, "sway", true},
		{map[string]string{"I3SOCK": "/i", "XDG_CURRENT_DESKTOP": "GNOME"}, "i3", true},
		{map[string]string{"XDG_CURRENT_DESKTOP": "KDE"}, "kde", false},
		{nil, "unknown", false},
	} {
		compositorEnv(t, c.env)
		if got := detectCompositor(); got != c.want {
			t.Errorf("%v: detectCompositor() = %q, want %q", c.env, got, c.want)
		}
		if got := isTilingWM(); got != c.tiling {
			t.Errorf("%v: isTilingWM() = %v, want %v", c.env, got, c.tiling)
		}
		if got := detectFocusFollowsMouse(); got != c.tiling {
			t.Errorf("%v: detectFocusFollowsMouse() = %v, want %v", c.env, got, c.tiling)
		}
	}
}

func TestIsWayland(t *testing.T) {
	for _, c := range []struct {
		env  map[string]string
		want bool
	}{
		{map[string]string{"XDG_SESSION_TYPE": "wayland"}, true},
		{map[string]string{"WAYLAND_DISPLAY": "wayland-0"}, true},
		{map[string]string{"XDG_SESSION_TYPE": "x11"}, false},
		{nil, false},
	} {
		compositorEnv(t, c.env)
		if got := isWayland(); got != c.want {
			t.Errorf("%v: isWayland() = %v, want %v", c.env, got, c.want)
		}
	}
}

func TestParseCursorPos(t *testing.T) {
	for _, c := range []struct {
		in   string
		x, y int
		ok   bool
	}{
		{"120, 340", 120, 340, true},
		{"-5, 0", -5, 0, true},
		{"120,340", 0, 0, false},
		{"a, 1", 0, 0, false},
		{"1, b", 0, 0, false},
		{"", 0, 0, false},
	} {
		x, y, ok := parseCursorPos(c.in)
		if x != c.x || y != c.y || ok != c.ok {
			t.Errorf("parseCursorPos(%q) = %d, %d, %v; want %d, %d, %v", c.in, x, y, ok, c.x, c.y, c.ok)
		}
	}
}

// listenUnix serves each connection on a Unix socket at path with handle.
func listenUnix(t *testing.T, path string, handle func(net.Conn)) {
	t.Helper()
	l, err := net.Listen("unix", path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			handle(conn)
			conn.Close()
		}
	}()
}

func TestHyprlandCursorPosition(t *testing.T) {
	runtimeDir := t.TempDir()
	compositorEnv(t, map[string]string{"HYPRLAND_INSTANCE_SIGNATURE": "sig", "XDG_RUNTIME_DIR": runtimeDir})
	if _, _, ok := getCursorPositionFromCompositor(); ok {
		t.Fatal("a position without Hyprland's socket")
	}
	dir := filepath.Join(runtimeDir, "hypr", "sig")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	commands := make(chan string, 1)
	listenUnix(t, filepath.Join(dir, ".socket.sock"), func(conn net.Conn) {
		buf := make([]byte, 64)
		n, _ := conn.Read(buf)
		commands <- string(buf[:n])
		conn.Write([]byte("640, 480\n"))
	})
	x, y, ok := getCursorPositionFromCompositor()
	if !ok || x != 640 || y != 480 {
		t.Errorf("cursor %d, %d, %v; want 640, 480, true", x, y, ok)
	}
	if got := <-commands; got != "cursorpos" {
		t.Errorf("Hyprland was sent %q, want cursorpos", got)
	}

	os.Unsetenv("HYPRLAND_INSTANCE_SIGNATURE")
	if _, err := hyprlandIPC("cursorpos"); err == nil {
		t.Error("hyprlandIPC without HYPRLAND_INSTANCE_SIGNATURE: no error")
	}
}

func TestSwayCursorPosition(t *testing.T) {
	// Sway's IPC has no pointer position: the caller falls back without
	// asking it.
	compositorEnv(t, map[string]string{"SWAYSOCK": filepath.Join(t.TempDir(), "sway.sock")})
	if x, y, ok := getCursorPositionFromCompositor(); ok {
		t.Errorf("sway: cursor %d, %d reported", x, y)
	}
}
