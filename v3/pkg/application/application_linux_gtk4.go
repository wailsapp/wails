//go:build linux && cgo && gtk4 && !android

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
	"strings"
	"sync"

	"github.com/godbus/dbus/v5"
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
	if os.Getenv("WEBKIT_DISABLE_DMABUF_RENDERER") == "" &&
		os.Getenv("XDG_SESSION_TYPE") == "wayland" &&
		isNVIDIAGPU() {
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

	theme string

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

	go listenForSystemThemeChanges(a)
}

func listenForSystemThemeChanges(a *linuxApp) {
	conn, err := dbus.SessionBus()
	if err != nil {
		a.parent.error("failed to connect to session bus: %v", err)
		return
	}

	if err = conn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.portal.Settings"),
		dbus.WithMatchMember("SettingChanged"),
	); err != nil {
		return
	}

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)

	for s := range c {
		if len(s.Body) < 3 {
			continue
		}
		namespace, ok := s.Body[0].(string)
		if !ok || namespace != "org.freedesktop.appearance" {
			continue
		}
		key, ok := s.Body[1].(string)
		if !ok || key != "color-scheme" {
			continue
		}
		processApplicationEvent(C.uint(events.Linux.SystemThemeChanged), nil)
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
	return nil, fmt.Errorf("icon lookup not implemented for GTK4")
}

func (a *linuxApp) isDarkMode() bool {
	conn, err := dbus.SessionBus()
	if err != nil {
		return false
	}

	obj := conn.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop")
	call := obj.Call("org.freedesktop.portal.Settings.Read", 0, "org.freedesktop.appearance", "color-scheme")
	if call.Err != nil {
		return false
	}

	var result dbus.Variant
	if err := call.Store(&result); err != nil {
		return false
	}

	innerVariant := result.Value().(dbus.Variant)
	colorScheme, ok := innerVariant.Value().(uint32)
	if !ok {
		return false
	}

	return colorScheme == 1
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

var _ = dbus.SessionBus
var _ = filepath.Ext
var _ = operatingsystem.Info

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
