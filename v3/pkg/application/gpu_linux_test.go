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
	hybrid := map[string]fakeCard{"card0": {"nvidia", "disconnected", false}, "card1": {"amdgpu", "connected", false}}
	offload := map[string]string{"__NV_PRIME_RENDER_OFFLOAD": "1", "__GLX_VENDOR_LIBRARY_NAME": "nvidia"}
	cases := []struct {
		name   string
		loaded bool
		cards  map[string]fakeCard
		env    map[string]string
		want   bool
	}{
		{"no nvidia module", false, map[string]fakeCard{"card1": {"amdgpu", "connected", false}}, nil, false},
		{"nvidia drives the screen", true, map[string]fakeCard{"card0": {"nvidia", "connected", false}}, nil, true},
		{"hybrid laptop on the AMD iGPU", true, hybrid, nil, false},
		{"hybrid laptop launched on the dGPU", true, hybrid, offload, true},
		{"external screen on the dGPU, desktop on the boot iGPU", true, map[string]fakeCard{"card0": {"nvidia", "connected", false}, "card1": {"amdgpu", "connected", true}}, nil, false},
		{"boot VGA is NVIDIA", true, map[string]fakeCard{"card0": {"nvidia", "connected", true}, "card1": {"amdgpu", "connected", false}}, nil, true},
		{"nothing known", true, nil, nil, true},
	}
	for _, c := range cases {
		getenv := func(k string) string { return c.env[k] }
		if got := nvidiaRenders(fakeSysfs(t, c.loaded, c.cards), getenv); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}
