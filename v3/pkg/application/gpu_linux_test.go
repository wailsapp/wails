//go:build linux && cgo && !android && !server

package application

import (
	"os"
	"path/filepath"
	"testing"
)

type fakeCard struct {
	driver, status string
	bootVGA        bool
}

// fakeSysfs builds a sysfs tree with the nvidia module loaded or not, and each card with its driver, one connector in
// the given state and, when set, the boot VGA mark.
func fakeSysfs(t *testing.T, nvidiaLoaded bool, cards map[string]fakeCard) string {
	t.Helper()
	root := t.TempDir()
	write := func(p, body string) {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	link := func(target, p string) {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, p); err != nil {
			t.Fatal(err)
		}
	}
	if nvidiaLoaded {
		write(filepath.Join(root, "module", "nvidia", "refcnt"), "1\n")
	}
	for card, c := range cards {
		drv := filepath.Join(root, "bus", "pci", "drivers", c.driver)
		if err := os.MkdirAll(drv, 0o755); err != nil {
			t.Fatal(err)
		}
		dev := filepath.Join(root, "devices", card)
		link(drv, filepath.Join(dev, "driver"))
		if c.bootVGA {
			write(filepath.Join(dev, "boot_vga"), "1\n")
		}
		link(dev, filepath.Join(root, "class", "drm", card, "device"))
		write(filepath.Join(root, "class", "drm", card+"-eDP-1", "status"), c.status+"\n")
	}
	return root
}

func TestNVIDIARenders(t *testing.T) {
	hybrid := map[string]fakeCard{"card0": {"nvidia", "disconnected", false}, "card1": {"amdgpu", "connected", true}}
	wayland := map[string]string{"XDG_SESSION_TYPE": "wayland"}
	with := func(kv ...string) map[string]string {
		env := map[string]string{"XDG_SESSION_TYPE": "wayland"}
		for i := 0; i < len(kv); i += 2 {
			env[kv[i]] = kv[i+1]
		}
		return env
	}
	mesaAndNVIDIA := "/usr/share/glvnd/egl_vendor.d/50_mesa.json:/usr/share/glvnd/egl_vendor.d/10_nvidia.json"
	cases := []struct {
		name   string
		loaded bool
		cards  map[string]fakeCard
		env    map[string]string
		want   bool
	}{
		{"no nvidia module", false, map[string]fakeCard{"card1": {"amdgpu", "connected", true}}, wayland, false},
		{"nvidia drives the screen", true, map[string]fakeCard{"card0": {"nvidia", "connected", true}}, wayland, true},
		{"hybrid laptop on the AMD iGPU", true, hybrid, wayland, false},
		{"hybrid laptop offloaded to the dGPU", true, hybrid, with("__NV_PRIME_RENDER_OFFLOAD", "1"), true},
		{"GLX vendor alone does not move EGL", true, hybrid, with("__GLX_VENDOR_LIBRARY_NAME", "nvidia"), false},
		{"EGL limited to NVIDIA", true, hybrid, with("__EGL_VENDOR_LIBRARY_FILENAMES", "/usr/share/glvnd/egl_vendor.d/10_nvidia.json"), true},
		{"EGL offered Mesa and NVIDIA", true, hybrid, with("__EGL_VENDOR_LIBRARY_FILENAMES", mesaAndNVIDIA), false},
		{"EGL limited to Mesa outranks offload", true, hybrid, with("__NV_PRIME_RENDER_OFFLOAD", "1", "__EGL_VENDOR_LIBRARY_FILENAMES", "/usr/share/glvnd/egl_vendor.d/50_mesa.json"), false},
		{"mixed EGL list leaves offload to decide", true, hybrid, with("__NV_PRIME_RENDER_OFFLOAD", "1", "__EGL_VENDOR_LIBRARY_FILENAMES", mesaAndNVIDIA), true},
		{"external screen on the dGPU, desktop on the boot iGPU", true, map[string]fakeCard{"card0": {"nvidia", "connected", false}, "card1": {"amdgpu", "connected", true}}, wayland, false},
		{"boot VGA is NVIDIA with an AMD screen too", true, map[string]fakeCard{"card0": {"nvidia", "connected", true}, "card1": {"amdgpu", "connected", false}}, wayland, true},
		{"screenless AMD boot VGA, screens on NVIDIA", true, map[string]fakeCard{"card0": {"nvidia", "connected", false}, "card1": {"amdgpu", "disconnected", true}}, wayland, true},
		{"X11 keeps the workaround", true, hybrid, map[string]string{"XDG_SESSION_TYPE": "x11"}, true},
		{"nothing known", true, nil, wayland, true},
	}
	for _, c := range cases {
		getenv := func(k string) string { return c.env[k] }
		if got := nvidiaRenders(fakeSysfs(t, c.loaded, c.cards), getenv); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}
