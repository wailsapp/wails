//go:build linux && cgo && !android && !server

package application

import (
	"os"
	"path/filepath"
	"strings"
)

// isNVIDIAGPU reports whether WebKit will render on NVIDIA's proprietary driver, whose DMA-BUF import fails in
// WebKitGTK. A loaded nvidia module alone is not enough: a laptop with an AMD or Intel iGPU and an NVIDIA dGPU (PRIME)
// renders on the iGPU unless the app was launched on the dGPU, and turning the DMA-BUF renderer off there drops
// rendering to a slow copy path.
func isNVIDIAGPU() bool {
	return nvidiaRenders("/sys", os.Getenv)
}

// nvidiaRenders reads the sysfs tree at root and the process environment. NVIDIA renders when the nvidia module is
// loaded and the app was offloaded to it (switcheroo's "Launch using Discrete Graphics Card" and prime-run set
// __NV_PRIME_RENDER_OFFLOAD, which is all EGL needs), glvnd offers EGL no vendor but NVIDIA's, or the session renders
// on it. A Wayland compositor renders on its primary GPU, which mutter and KWin take to be the boot VGA device, and
// WebKit renders on the compositor's device; the boot VGA mark is trusted only when that card drives a connected
// screen, and otherwise the cards driving the connected screens decide. Xorg renders where its configuration says
// (reverse PRIME puts NVIDIA behind the iGPU's screen), which cannot be read here, so an X11 session keeps the
// workaround whenever the module is loaded.
func nvidiaRenders(root string, getenv func(string) string) bool {
	if _, err := os.Stat(filepath.Join(root, "module", "nvidia")); err != nil {
		return false
	}
	// glvnd loads only the EGL vendors the list names, so a list settles it before offload can: no NVIDIA vendor
	// means NVIDIA cannot render EGL, only NVIDIA's means it must.
	switch eglVendors(getenv("__EGL_VENDOR_LIBRARY_FILENAMES")) {
	case vendorsNone:
		return false
	case vendorsOnlyNVIDIA:
		return true
	}
	if getenv("__NV_PRIME_RENDER_OFFLOAD") == "1" {
		return true
	}
	if getenv("XDG_SESSION_TYPE") != "wayland" && getenv("WAYLAND_DISPLAY") == "" {
		return true
	}
	drm := filepath.Join(root, "class", "drm")
	driverOf := func(card string) string {
		driver, err := filepath.EvalSymlinks(filepath.Join(drm, card, "device", "driver"))
		if err != nil {
			return ""
		}
		return filepath.Base(driver)
	}
	screens := map[string]bool{}
	statuses, _ := filepath.Glob(filepath.Join(drm, "card*-*", "status"))
	for _, status := range statuses {
		if b, err := os.ReadFile(status); err == nil && strings.TrimSpace(string(b)) == "connected" {
			card, _, _ := strings.Cut(filepath.Base(filepath.Dir(status)), "-")
			screens[card] = true
		}
	}
	for card := range screens {
		if b, err := os.ReadFile(filepath.Join(drm, card, "device", "boot_vga")); err == nil && strings.TrimSpace(string(b)) == "1" {
			return driverOf(card) == "nvidia"
		}
	}
	for card := range screens {
		if driver := driverOf(card); driver != "" && driver != "nvidia" {
			return false
		}
	}
	return true
}

type vendorList int

const (
	vendorsUnset vendorList = iota
	vendorsNone
	vendorsMixed
	vendorsOnlyNVIDIA
)

// eglVendors reads a glvnd EGL vendor list: unset, naming no NVIDIA vendor file, mixed, or naming only NVIDIA's.
func eglVendors(files string) vendorList {
	if files == "" {
		return vendorsUnset
	}
	nvidia, other := false, false
	for _, f := range strings.Split(files, ":") {
		if strings.Contains(filepath.Base(f), "nvidia") {
			nvidia = true
		} else {
			other = true
		}
	}
	switch {
	case !nvidia:
		return vendorsNone
	case other:
		return vendorsMixed
	}
	return vendorsOnlyNVIDIA
}
