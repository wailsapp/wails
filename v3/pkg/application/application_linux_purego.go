//go:build linux && purego && !gtk3 && !android && !server

package application

import (
	"fmt"
	"os"
	"sync"

	"github.com/wailsapp/wails/v3/internal/operatingsystem"
)

func init() {
	// Disable DMA-BUF renderer on any session type with NVIDIA to prevent blank windows and
	// "Error 71 (Protocol error)" crashes. NVIDIA proprietary drivers fail gbm_bo_map() when
	// importing DMA-BUF, causing blank/white screens on both X11 and Wayland.
	// See: https://bugs.webkit.org/show_bug.cgi?id=262607
	// See: https://github.com/wailsapp/wails/issues/4985
	if os.Getenv("WEBKIT_DISABLE_DMABUF_RENDERER") == "" && isNVIDIAGPU() {
		_ = os.Setenv("WEBKIT_DISABLE_DMABUF_RENDERER", "1")
	}
}

type linuxApp struct {
	application pointer
	parent      *App

	activated     chan struct{}
	activatedOnce sync.Once

	windowMap     map[windowPointer]uint
	windowMapLock sync.Mutex


	icon pointer
}

func (a *linuxApp) appendGTKVersion(result map[string]string) {
	result["GTK"] = fmt.Sprintf("%d.%d.%d",
		gtk_get_major_version(),
		gtk_get_minor_version(),
		gtk_get_micro_version())
	result["WebKit"] = fmt.Sprintf("%d.%d.%d",
		webkit_get_major_version(),
		webkit_get_minor_version(),
		webkit_get_micro_version())
}

func (a *linuxApp) init(_ *App, options Options) {
	osInfo, _ := operatingsystem.Info()
	a.parent.info("Running GTK %d.%d.%d",
		gtk_get_major_version(),
		gtk_get_minor_version(),
		gtk_get_micro_version())
	a.parent.info("Running WebKitGTK %d.%d.%d",
		webkit_get_major_version(),
		webkit_get_minor_version(),
		webkit_get_micro_version())
	a.parent.info("Using %s", osInfo.Name)

	if options.Icon != nil {
		a.setIcon(options.Icon)
	}
}

func buildVersionString(major, minor, micro uint32) string {
	return fmt.Sprintf("%d.%d.%d", uint(major), uint(minor), uint(micro))
}

func (a *App) platformEnvironment() map[string]any {
	result := map[string]any{}
	result["gtk4-compiled"] = buildVersionString(
		gtk_get_major_version(),
		gtk_get_minor_version(),
		gtk_get_micro_version(),
	)
	result["gtk4-runtime"] = buildVersionString(
		gtk_get_major_version(),
		gtk_get_minor_version(),
		gtk_get_micro_version(),
	)

	result["webkitgtk6-compiled"] = buildVersionString(
		webkit_get_major_version(),
		webkit_get_minor_version(),
		webkit_get_micro_version(),
	)
	result["webkitgtk6-runtime"] = buildVersionString(
		webkit_get_major_version(),
		webkit_get_minor_version(),
		webkit_get_micro_version(),
	)

	result["compositor"] = detectCompositor()
	result["wayland"] = isWayland()
	result["focusFollowsMouse"] = detectFocusFollowsMouse()

	return result
}

// gtkVersion and webkitVersion report the versions loaded at runtime. The cgo
// backend reports the versions it was compiled against instead.
func gtkVersion() (major, minor, micro uint) {
	return uint(gtk_get_major_version()), uint(gtk_get_minor_version()), uint(gtk_get_micro_version())
}

func webkitVersion() (major, minor, micro uint) {
	return uint(webkit_get_major_version()), uint(webkit_get_minor_version()), uint(webkit_get_micro_version())
}
