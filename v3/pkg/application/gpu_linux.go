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
// loaded and either the app was offloaded to it (switcheroo's "Launch using Discrete Graphics Card", prime-run), or
// the compositor's GPU is NVIDIA's: the boot VGA device when the kernel marks one, else the cards driving the
// connected screens. Without any of that to read, a loaded module counts.
func nvidiaRenders(root string, getenv func(string) string) bool {
	if _, err := os.Stat(filepath.Join(root, "module", "nvidia")); err != nil {
		return false
	}
	if getenv("__NV_PRIME_RENDER_OFFLOAD") == "1" || getenv("__GLX_VENDOR_LIBRARY_NAME") == "nvidia" ||
		strings.Contains(getenv("__EGL_VENDOR_LIBRARY_FILENAMES"), "nvidia") {
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
	cards, _ := filepath.Glob(filepath.Join(drm, "card[0-9]*"))
	for _, card := range cards {
		if strings.Contains(filepath.Base(card), "-") {
			continue
		}
		if b, err := os.ReadFile(filepath.Join(card, "device", "boot_vga")); err == nil && strings.TrimSpace(string(b)) == "1" {
			return driverOf(filepath.Base(card)) == "nvidia"
		}
	}
	statuses, _ := filepath.Glob(filepath.Join(drm, "card*-*", "status"))
	for _, status := range statuses {
		b, err := os.ReadFile(status)
		if err != nil || strings.TrimSpace(string(b)) != "connected" {
			continue
		}
		card, _, _ := strings.Cut(filepath.Base(filepath.Dir(status)), "-")
		if driver := driverOf(card); driver != "" && driver != "nvidia" {
			return false
		}
	}
	return true
}
