//go:build linux && cgo && !gtk3 && !android && !server && !purego

package application

/*
#include <gtk/gtk.h>
#include <webkit/webkit.h>
static guint get_compiled_gtk_major_version() { return GTK_MAJOR_VERSION; }
static guint get_compiled_gtk_minor_version() { return GTK_MINOR_VERSION; }
static guint get_compiled_gtk_micro_version() { return GTK_MICRO_VERSION; }
static guint get_compiled_webkit_major_version() { return WEBKIT_MAJOR_VERSION; }
static guint get_compiled_webkit_minor_version() { return WEBKIT_MINOR_VERSION; }
static guint get_compiled_webkit_micro_version() { return WEBKIT_MICRO_VERSION; }
*/
import "C"
import (
	"fmt"
	"os"
	"sync"

	"github.com/wailsapp/wails/v3/internal/operatingsystem"
)

func init() {
	// Disable the DMA-BUF renderer when WebKit renders on NVIDIA (isNVIDIAGPU), on any session type, to prevent blank windows and
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
		C.get_compiled_gtk_major_version(),
		C.get_compiled_gtk_minor_version(),
		C.get_compiled_gtk_micro_version())
	result["WebKit"] = fmt.Sprintf("%d.%d.%d",
		C.get_compiled_webkit_major_version(),
		C.get_compiled_webkit_minor_version(),
		C.get_compiled_webkit_micro_version())
}

func (a *linuxApp) init(_ *App, options Options) {
	osInfo, _ := operatingsystem.Info()
	a.parent.info("Compiled with GTK %d.%d.%d",
		C.get_compiled_gtk_major_version(),
		C.get_compiled_gtk_minor_version(),
		C.get_compiled_gtk_micro_version())
	a.parent.info("Compiled with WebKitGTK %d.%d.%d",
		C.get_compiled_webkit_major_version(),
		C.get_compiled_webkit_minor_version(),
		C.get_compiled_webkit_micro_version())
	a.parent.info("Using %s", osInfo.Name)

	if options.Icon != nil {
		a.setIcon(options.Icon)
	}
}

func buildVersionString(major, minor, micro C.guint) string {
	return fmt.Sprintf("%d.%d.%d", uint(major), uint(minor), uint(micro))
}

func (a *App) platformEnvironment() map[string]any {
	result := map[string]any{}
	result["gtk4-compiled"] = buildVersionString(
		C.get_compiled_gtk_major_version(),
		C.get_compiled_gtk_minor_version(),
		C.get_compiled_gtk_micro_version(),
	)
	result["gtk4-runtime"] = buildVersionString(
		C.gtk_get_major_version(),
		C.gtk_get_minor_version(),
		C.gtk_get_micro_version(),
	)

	result["webkitgtk6-compiled"] = buildVersionString(
		C.get_compiled_webkit_major_version(),
		C.get_compiled_webkit_minor_version(),
		C.get_compiled_webkit_micro_version(),
	)
	result["webkitgtk6-runtime"] = buildVersionString(
		C.webkit_get_major_version(),
		C.webkit_get_minor_version(),
		C.webkit_get_micro_version(),
	)

	result["compositor"] = detectCompositor()
	result["wayland"] = isWayland()
	result["focusFollowsMouse"] = detectFocusFollowsMouse()

	return result
}

// gtkVersion and webkitVersion report the versions this binary was compiled
// against. The purego backend cannot read the compiled headers, so it reports
// the versions loaded at runtime instead; that is the one deliberate difference.
func gtkVersion() (major, minor, micro uint) {
	return uint(C.get_compiled_gtk_major_version()), uint(C.get_compiled_gtk_minor_version()), uint(C.get_compiled_gtk_micro_version())
}

func webkitVersion() (major, minor, micro uint) {
	return uint(C.get_compiled_webkit_major_version()), uint(C.get_compiled_webkit_minor_version()), uint(C.get_compiled_webkit_micro_version())
}
