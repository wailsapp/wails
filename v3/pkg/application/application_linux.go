//go:build linux && cgo && !gtk3 && !android && !server

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
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v3/internal/operatingsystem"
	"github.com/wailsapp/wails/v3/pkg/events"
)

var invalidAppNameChars = regexp.MustCompile(`[^a-zA-Z0-9_-]`)
var leadingDigits = regexp.MustCompile(`^[0-9]+`)

func sanitizeAppName(name string) string {
	name = invalidAppNameChars.ReplaceAllString(name, "_")
	name = leadingDigits.ReplaceAllString(name, "_$0")
	for strings.Contains(name, "__") {
		name = strings.ReplaceAll(name, "__", "_")
	}
	name = strings.Trim(name, "_")
	if name == "" {
		name = "wailsapp"
	}
	return strings.ToLower(name)
}

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

func isNVIDIAGPU() bool {
	if _, err := os.Stat("/sys/module/nvidia"); err == nil {
		return true
	}
	return false
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

func (a *linuxApp) GetFlags(options Options) map[string]any {
	if options.Flags == nil {
		options.Flags = make(map[string]any)
	}
	return options.Flags
}

func (a *linuxApp) name() string {
	return appName()
}

func (a *linuxApp) run() error {
	if len(os.Args) == 2 {
		arg1 := os.Args[1]
		if strings.Contains(arg1, "://") {
			eventContext := newApplicationEventContext()
			eventContext.setURL(arg1)
			applicationEvents <- &ApplicationEvent{
				Id:  uint(events.Common.ApplicationLaunchedWithUrl),
				ctx: eventContext,
			}
		} else if a.parent.options.FileAssociations != nil {
			ext := filepath.Ext(arg1)
			if slices.Contains(a.parent.options.FileAssociations, ext) {
				eventContext := newApplicationEventContext()
				eventContext.setOpenedWithFile(arg1)
				applicationEvents <- &ApplicationEvent{
					Id:  uint(events.Common.ApplicationOpenedWithFile),
					ctx: eventContext,
				}
			}
		}
	}

	a.parent.Event.OnApplicationEvent(events.Linux.ApplicationStartup, func(evt *ApplicationEvent) {
		if err := a.processAndCacheScreens(); err != nil {
			a.parent.handleError(err)
		}
	})
	a.setupCommonEvents()
	// Started here, not from init(): init() is not part of the platformApp
	// interface and nothing calls it, so a monitor started there never runs.
	a.monitorThemeChanges()
	a.monitorPowerEvents()
	return appRun(a.application)
}

func (a *linuxApp) destroy() {
	if !globalApplication.shouldQuit() {
		return
	}
	globalApplication.cleanup()
	appDestroy(a.application)
}

func (a *linuxApp) getApplicationMenu() *Menu {
	return nil
}

func (a *linuxApp) setApplicationMenu(menu *Menu) {}

func (a *linuxApp) hide() {
	a.hideAllWindows()
}

func (a *linuxApp) show() {
	a.showAllWindows()
}

func (a *linuxApp) on(eventID uint) {
}

func (a *linuxApp) isOnMainThread() bool {
	return isOnMainThread()
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

func (a *linuxApp) registerWindow(window pointer, id uint) {
	a.windowMapLock.Lock()
	a.windowMap[windowPointer(window)] = id
	a.windowMapLock.Unlock()
}

func (a *linuxApp) unregisterWindow(window windowPointer) {
	a.windowMapLock.Lock()
	delete(a.windowMap, window)
	remainingWindows := len(a.windowMap)
	a.windowMapLock.Unlock()

	if remainingWindows == 0 && !a.parent.options.Linux.DisableQuitOnLastWindowClosed {
		a.destroy()
	}
}

func newPlatformApp(parent *App) *linuxApp {
	name := sanitizeAppName(parent.options.Name)
	app := &linuxApp{
		parent:      parent,
		application: appNew(name),
		activated:   make(chan struct{}),
		windowMap:   map[windowPointer]uint{},
	}

	if parent.options.Linux.ProgramName != "" {
		setProgramName(parent.options.Linux.ProgramName)
	}

	return app
}

func (a *linuxApp) markActivated() {
	a.activatedOnce.Do(func() {
		close(a.activated)
	})
}

func (a *linuxApp) waitForActivation() {
	<-a.activated
}

func (a *linuxApp) getIconForFile(filename string) ([]byte, error) {
	if filename == "" {
		return nil, nil
	}

	ext := filepath.Ext(filename)
	iconMap := map[string]string{
		".txt":  "text-x-generic",
		".pdf":  "application-pdf",
		".doc":  "x-office-document",
		".docx": "x-office-document",
		".xls":  "x-office-spreadsheet",
		".xlsx": "x-office-spreadsheet",
		".ppt":  "x-office-presentation",
		".pptx": "x-office-presentation",
		".zip":  "package-x-generic",
		".tar":  "package-x-generic",
		".gz":   "package-x-generic",
		".jpg":  "image-x-generic",
		".jpeg": "image-x-generic",
		".png":  "image-x-generic",
		".gif":  "image-x-generic",
		".mp3":  "audio-x-generic",
		".wav":  "audio-x-generic",
		".mp4":  "video-x-generic",
		".avi":  "video-x-generic",
		".html": "text-html",
		".css":  "text-css",
		".js":   "text-javascript",
		".json": "text-json",
		".xml":  "text-xml",
	}

	iconName := "application-x-generic"
	if name, ok := iconMap[ext]; ok {
		iconName = name
	}

	return getIconBytes(iconName)
}

func getIconBytes(iconName string) ([]byte, error) {
	return nil, fmt.Errorf("icon lookup is not currently implemented for the GTK4 build path; build with -tags gtk3 for the legacy implementation")
}

func (a *linuxApp) getAccentColor() string {
	return "rgb(0,122,255)"
}

func (a *linuxApp) isVisible() bool {
	windows := a.getWindows()
	for _, window := range windows {
		if C.gtk_widget_is_visible((*C.GtkWidget)(window)) != 0 {
			return true
		}
	}
	return false
}

func getNativeApplication() *linuxApp {
	return globalApplication.impl.(*linuxApp)
}

// logPlatformInfo logs the platform information to the console
func (a *App) logPlatformInfo() {
	info, err := operatingsystem.Info()
	if err != nil {
		a.error("error getting OS info: %w", err)
		return
	}

	platformInfo := info.AsLogSlice()
	platformInfo = append(platformInfo, "GTK", fmt.Sprintf("%d.%d.%d",
		C.get_compiled_gtk_major_version(),
		C.get_compiled_gtk_minor_version(),
		C.get_compiled_gtk_micro_version()))
	platformInfo = append(platformInfo, "WebKitGTK", fmt.Sprintf("%d.%d.%d",
		C.get_compiled_webkit_major_version(),
		C.get_compiled_webkit_minor_version(),
		C.get_compiled_webkit_micro_version()))

	a.info("Platform Info:", platformInfo...)
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

func fatalHandler(errFunc func(error)) {
	// Stub for windows function
	return
}
