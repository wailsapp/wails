//go:build linux && cgo && gtk3 && !android && !server && !purego

package application

/*
	#include "gtk/gtk.h"
	#include "webkit2/webkit2.h"
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
)

func init() {
	// FIXME: This should be handled appropriately in the individual files most likely.
	// Set GDK_BACKEND=x11 if currently unset and XDG_SESSION_TYPE is unset, unspecified or x11 to prevent warnings
	if os.Getenv("GDK_BACKEND") == "" &&
		(os.Getenv("XDG_SESSION_TYPE") == "" || os.Getenv("XDG_SESSION_TYPE") == "unspecified" || os.Getenv("XDG_SESSION_TYPE") == "x11") {
		_ = os.Setenv("GDK_BACKEND", "x11")
	}

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

	startupActions []func()

	// Native -> uint
	windowMap     map[windowPointer]uint
	windowMapLock sync.Mutex

	icon pointer
}

type rnr struct {
	f func()
}

//export processWindowEvent
func processWindowEvent(windowID C.uint, eventID C.uint) {
	windowEvents <- &windowEvent{
		WindowID: uint(windowID),
		EventID:  uint(eventID),
	}
}

func buildVersionString(major, minor, micro C.uint) string {
	return fmt.Sprintf("%d.%d.%d", uint(major), uint(minor), uint(micro))
}

func (a *App) platformEnvironment() map[string]any {
	result := map[string]any{}
	result["gtk3-compiled"] = buildVersionString(
		C.get_compiled_gtk_major_version(),
		C.get_compiled_gtk_minor_version(),
		C.get_compiled_gtk_micro_version(),
	)
	result["gtk3-runtime"] = buildVersionString(
		C.gtk_get_major_version(),
		C.gtk_get_minor_version(),
		C.gtk_get_micro_version(),
	)

	result["webkit2gtk-compiled"] = buildVersionString(
		C.get_compiled_webkit_major_version(),
		C.get_compiled_webkit_minor_version(),
		C.get_compiled_webkit_micro_version(),
	)
	result["webkit2gtk-runtime"] = buildVersionString(
		C.webkit_get_major_version(),
		C.webkit_get_minor_version(),
		C.webkit_get_micro_version(),
	)

	result["compositor"] = detectCompositor()
	result["wayland"] = isWayland()
	result["focusFollowsMouse"] = detectFocusFollowsMouse()

	return result
}
