//go:build linux && purego && !gtk3 && !android && !server

package application

// GTK4 purego tests that drive the real native objects: GMenu/GAction for the
// menus (no display needed, but the action activation path goes through the
// same GLib signal machinery a user click uses) and GtkApplicationWindow /
// WebKitWebView for the window setters. The window and screen tests need a
// display and call gtk_init_check; they skip without one and fail when a
// display is advertised but GTK cannot start on it. Run them with:
//
//	xvfb-run -a dbus-run-session -- go test ./pkg/application

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/ebitengine/purego"
)

// Test-only bindings. The production bindings in linux_purego_generated.go
// cover what the backend calls; these read state back so the tests can check
// the outcome of a call rather than trusting that it returned.
var (
	testBindOnce sync.Once

	tGtkInitCheck                     func() int32
	tGtkWidgetActivateActionVariant   func(uintptr, string, uintptr) int32
	tGApplicationRegister             func(uintptr, uintptr, uintptr) int32
	tGMenuModelGetNItems              func(uintptr) int32
	tGMenuModelGetItemAttributeValue  func(uintptr, int32, string, uintptr) uintptr
	tGMenuModelGetItemLink            func(uintptr, int32, string) uintptr
	tGActionGroupGetActionState       func(uintptr, string) uintptr
	tGActionGroupGetActionEnabled     func(uintptr, string) int32
	tGActionGroupActivateAction       func(uintptr, string, uintptr)
	tGtkWindowGetTitle                func(uintptr) string
	tGtkWindowGetDecorated            func(uintptr) int32
	tGtkWindowGetTitlebar             func(uintptr) uintptr
	tGtkWidgetGetOpacity              func(uintptr) float64
	tGtkWidgetGetSizeRequest          func(uintptr, uintptr, uintptr)
	tGtkWidgetGetFirstChild           func(uintptr) uintptr
	tGtkApplicationGetAccelsForAction func(uintptr, string) uintptr
	tWebkitGetBackgroundColor         func(uintptr, uintptr)
	tGtkPopoverMenuBarGetMenuModel    func(uintptr) uintptr
	tGtkWidgetGetNextSibling          func(uintptr) uintptr
	tGtkWidgetObserveControllers      func(uintptr) uintptr
	tGtkDropTargetGetType             func() uint64
)

func bindTestNatives() {
	testBindOnce.Do(func() {
		purego.RegisterLibFunc(&tGtkInitCheck, libGtk, "gtk_init_check")
		purego.RegisterLibFunc(&tGtkWidgetActivateActionVariant, libGtk, "gtk_widget_activate_action_variant")
		purego.RegisterLibFunc(&tGApplicationRegister, libGio, "g_application_register")
		purego.RegisterLibFunc(&tGMenuModelGetNItems, libGio, "g_menu_model_get_n_items")
		purego.RegisterLibFunc(&tGMenuModelGetItemAttributeValue, libGio, "g_menu_model_get_item_attribute_value")
		purego.RegisterLibFunc(&tGMenuModelGetItemLink, libGio, "g_menu_model_get_item_link")
		purego.RegisterLibFunc(&tGActionGroupGetActionState, libGio, "g_action_group_get_action_state")
		purego.RegisterLibFunc(&tGActionGroupGetActionEnabled, libGio, "g_action_group_get_action_enabled")
		purego.RegisterLibFunc(&tGActionGroupActivateAction, libGio, "g_action_group_activate_action")
		purego.RegisterLibFunc(&tGtkWindowGetTitle, libGtk, "gtk_window_get_title")
		purego.RegisterLibFunc(&tGtkWindowGetDecorated, libGtk, "gtk_window_get_decorated")
		purego.RegisterLibFunc(&tGtkWindowGetTitlebar, libGtk, "gtk_window_get_titlebar")
		purego.RegisterLibFunc(&tGtkWidgetGetOpacity, libGtk, "gtk_widget_get_opacity")
		purego.RegisterLibFunc(&tGtkWidgetGetSizeRequest, libGtk, "gtk_widget_get_size_request")
		purego.RegisterLibFunc(&tGtkWidgetGetFirstChild, libGtk, "gtk_widget_get_first_child")
		purego.RegisterLibFunc(&tGtkApplicationGetAccelsForAction, libGtk, "gtk_application_get_accels_for_action")
		purego.RegisterLibFunc(&tWebkitGetBackgroundColor, libWebKit, "webkit_web_view_get_background_color")
		purego.RegisterLibFunc(&tGtkPopoverMenuBarGetMenuModel, libGtk, "gtk_popover_menu_bar_get_menu_model")
		purego.RegisterLibFunc(&tGtkWidgetGetNextSibling, libGtk, "gtk_widget_get_next_sibling")
		purego.RegisterLibFunc(&tGtkWidgetObserveControllers, libGtk, "gtk_widget_observe_controllers")
		purego.RegisterLibFunc(&tGtkDropTargetGetType, libGtk, "gtk_drop_target_get_type")
	})
}

// gtkChildEnv marks the isolated child process started by runGTKTestInChild.
const gtkChildEnv = "WAILS_PUREGO_GTK_CHILD"

// runGTKTestInChild runs the calling GTK test in a fresh process and reports
// its result to t. It returns true only inside that child, where the test body
// should run.
//
// WebKit binds its work to the thread that first used it, and each Go test
// runs on whichever OS thread the scheduler picks. In one process, a second
// webview created on another thread blocks inside WebKit (observed: the
// test hangs in g_object_new_with_properties). A process per test keeps each
// test on a single thread. When WAILS_PUREGO_CHILD_COVER names a directory,
// each child also writes a coverage profile there, to be merged with the
// parent's profile afterwards (go test does not collect them itself).
func runGTKTestInChild(t *testing.T) bool {
	t.Helper()
	if os.Getenv(gtkChildEnv) != "" {
		return true
	}
	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		t.Skip("no display: run under xvfb-run -a")
	}
	args := []string{"-test.run=^" + regexp.QuoteMeta(t.Name()) + "$", "-test.count=1", "-test.v=true"}
	// The child's coverage counts too: it writes into the parent's
	// -test.gocoverdir (raw covdata, merged by go tool covdata), or a profile
	// per test into WAILS_PUREGO_CHILD_COVER.
	if f := flag.Lookup("test.gocoverdir"); f != nil && f.Value.String() != "" {
		args = append(args, "-test.gocoverdir="+f.Value.String())
	} else if dir := os.Getenv("WAILS_PUREGO_CHILD_COVER"); dir != "" {
		args = append(args, "-test.coverprofile="+dir+"/"+t.Name()+".out")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], args...)
	cmd.Env = append(os.Environ(), gtkChildEnv+"=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s failed in its own process: %v\n%s", t.Name(), err, out)
	}
	if bytes.Contains(out, []byte("--- SKIP")) {
		t.Skipf("%s skipped in its own process:\n%s", t.Name(), out)
	}
	return false
}

var (
	gtkInitOnce sync.Once
	gtkInitOK   bool
	testAppSeq  atomic.Uint32
)

// requireGTKDisplay initialises GTK on the calling thread. It skips when no
// display is available and fails when one is but GTK still cannot start.
func requireGTKDisplay(t *testing.T) {
	t.Helper()
	if linuxLibsErr != nil {
		t.Fatal(linuxLibsErr)
	}
	bindTestNatives()
	gtkInitOnce.Do(func() { gtkInitOK = tGtkInitCheck() != 0 })
	if gtkInitOK {
		return
	}
	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		t.Skip("no display: run under xvfb-run -a")
	}
	t.Fatal("display is set but gtk_init_check failed")
}

// newTestLinuxApp returns a GtkApplication with a registered (started) GApplication
// and installs a globalApplication whose main-thread check reports true on this
// goroutine. That makes InvokeSync/InvokeAsync run inline, so the production
// call paths can be exercised without a running GTK main loop.
func newTestLinuxApp(t *testing.T) *linuxApp {
	t.Helper()
	requireGTKDisplay(t)
	runtime.LockOSThread()
	t.Cleanup(runtime.UnlockOSThread)

	prevMain := mainThreadId
	mainThreadId = g_thread_self()
	t.Cleanup(func() { mainThreadId = prevMain })

	appID := fmt.Sprintf("org.wails.test.coverage%d", testAppSeq.Add(1))
	gtkApp := gtk_application_new(appID, gApplicationDefaultFlags)
	if gtkApp == 0 {
		t.Fatal("gtk_application_new returned NULL")
	}
	tGApplicationRegister(gtkApp, 0, 0)
	t.Cleanup(func() { g_object_unref(gtkApp) })

	app := &linuxApp{
		application: pointer(gtkApp),
		activated:   make(chan struct{}),
		windowMap:   map[windowPointer]uint{},
	}
	prevGlobal := globalApplication
	globalApplication = &App{impl: app, running: true, Logger: slog.Default()}
	globalApplication.Screen = newScreenManager(globalApplication)
	app.parent = globalApplication
	t.Cleanup(func() { globalApplication = prevGlobal })
	return app
}

// newTestWindow creates a real GtkApplicationWindow with a WebKitWebView via
// the production windowNew.
func newTestWindow(t *testing.T, app *linuxApp, id uint, menu pointer, style LinuxMenuStyle) *linuxWebviewWindow {
	t.Helper()
	window, webview, vbox, menubar := windowNew(app.application, menu, style, id, WebviewGpuPolicyNever)
	if window == nilPointer || webview == nilPointer || vbox == nilPointer {
		t.Fatalf("windowNew returned window=%#x webview=%#x vbox=%#x", window, webview, vbox)
	}
	app.registerWindow(window, id)
	w := &linuxWebviewWindow{
		id:          id,
		application: app.application,
		window:      window,
		webview:     webview,
		vbox:        vbox,
		menubar:     menubar,
		parent:      &WebviewWindow{id: id},
	}
	t.Cleanup(func() {
		app.unregisterWindowQuiet(window)
		gtk_window_destroy(uintptr(window))
	})
	return w
}

// unregisterWindowQuiet removes a window from the map without the quit-on-last
// -window side effect, for test teardown.
func (a *linuxApp) unregisterWindowQuiet(window pointer) {
	a.windowMapLock.Lock()
	delete(a.windowMap, windowPointer(window))
	a.windowMapLock.Unlock()
}

func testGStringVariant(v uintptr) string {
	if v == 0 {
		return ""
	}
	defer g_variant_unref(v)
	return goString(g_variant_get_string(v, 0))
}

func testMenuItemLabel(model uintptr, index int32) string {
	return testGStringVariant(tGMenuModelGetItemAttributeValue(model, index, "label", 0))
}

// testMenuItem builds a production MenuItem with its native widget.
func testMenuItem(label string) *MenuItem {
	item := NewMenuItem(label)
	item.impl = newMenuItemImpl(item)
	return item
}

func testCheckMenuItem(label string, checked bool) *MenuItem {
	item := NewMenuItemCheckbox(label, checked)
	item.impl = newCheckMenuItemImpl(item)
	return item
}

func testRadioMenuItem(label string, groupId, checkedId uint, id uint) *MenuItem {
	item := &MenuItem{id: id, label: label, itemType: radio}
	item.impl = newRadioMenuItemImpl(item, groupId, checkedId)
	return item
}

func testMenuImpl(native pointer) *Menu {
	return &Menu{impl: &linuxMenu{native: native}}
}

// Window tests -------------------------------------------------------------

func TestGTKWindowTitle(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1001, nilPointer, LinuxMenuStyleMenuBar)

	w.setTitle("Wails coverage")
	if got := tGtkWindowGetTitle(uintptr(w.window)); got != "Wails coverage" {
		t.Fatalf("title = %q", got)
	}

	// Frameless windows keep their title unset: setTitle is a no-op there.
	w.parent.options.Frameless = true
	w.setTitle("hidden")
	if got := tGtkWindowGetTitle(uintptr(w.window)); got != "Wails coverage" {
		t.Fatalf("frameless window title changed to %q", got)
	}
}

func TestGTKWindowSizeReadsBackDefaultSize(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1002, nilPointer, LinuxMenuStyleMenuBar)

	// Unrealised windows report no allocation, so size() falls back to the
	// default size, which is what setSize/setDefaultSize write.
	w.setSize(640, 480)
	if width, height := w.size(); width != 640 || height != 480 {
		t.Fatalf("size after setSize = %dx%d", width, height)
	}
	w.setDefaultSize(300, 200)
	if width, height := w.size(); width != 300 || height != 200 {
		t.Fatalf("size after setDefaultSize = %dx%d", width, height)
	}
}

func TestGTKWindowGeometryHints(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1003, nilPointer, LinuxMenuStyleMenuBar)

	windowSetGeometryHints(w.window, 120, 90, 0, 0)
	var minW, minH int32
	tGtkWidgetGetSizeRequest(uintptr(w.window), uintptr(unsafe.Pointer(&minW)), uintptr(unsafe.Pointer(&minH)))
	if minW != 120 || minH != 90 {
		t.Fatalf("min size request = %dx%d, want 120x90", minW, minH)
	}

	// A zero minimum does not change the existing request.
	windowSetGeometryHints(w.window, 0, 0, 800, 600)
	tGtkWidgetGetSizeRequest(uintptr(w.window), uintptr(unsafe.Pointer(&minW)), uintptr(unsafe.Pointer(&minH)))
	if minW != 120 || minH != 90 {
		t.Fatalf("zero minimum overwrote request: %dx%d", minW, minH)
	}
}

// TestGTKWindowBorderless covers the decoration toggle. setResizable and
// setFrameless also push a JavaScript update into the page; that needs a
// WebKit web process, which the sandboxed test container cannot start, so
// those two setters are not exercised here (see the coverage report).
func TestGTKWindowBorderless(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1004, nilPointer, LinuxMenuStyleMenuBar)

	w.setBorderless(true)
	if tGtkWindowGetDecorated(uintptr(w.window)) != 0 {
		t.Fatal("borderless window still decorated")
	}
	w.setBorderless(false)
	if tGtkWindowGetDecorated(uintptr(w.window)) == 0 {
		t.Fatal("bordered window not decorated")
	}
}

func TestGTKWindowOpacity(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1006, nilPointer, LinuxMenuStyleMenuBar)

	w.setOpacity(0.25)
	// GTK stores opacity as a byte, so 0.25 reads back as 64/255.
	if got := tGtkWidgetGetOpacity(uintptr(w.window)); got < 0.245 || got > 0.255 {
		t.Fatalf("opacity = %v, want 0.25", got)
	}
}

func TestGTKWindowShowHide(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1007, nilPointer, LinuxMenuStyleMenuBar)

	w.windowShow()
	if !w.isVisible() {
		t.Fatal("window not visible after windowShow")
	}
	w.windowHide()
	if w.isVisible() {
		t.Fatal("window still visible after windowHide")
	}
	w.present()
	if !w.isVisible() {
		t.Fatal("window not visible after present")
	}
}

func TestGTKWindowZoomIsValidatedAndScaled(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1008, nilPointer, LinuxMenuStyleMenuBar)

	if got := w.getZoom(); got != 1 {
		t.Fatalf("initial zoom = %v, want 1", got)
	}
	w.setZoom(0.5) // framed windows may zoom below 1 (zoomFactor)
	if got := w.getZoom(); got != 0.5 {
		t.Fatalf("zoom after setZoom(0.5) = %v, want 0.5", got)
	}
	w.setZoom(0) // non-positive factors are ignored
	if got := w.getZoom(); got != 0.5 {
		t.Fatalf("zoom after setZoom(0) = %v, want 0.5", got)
	}
	w.setZoom(2)
	if got := w.getZoom(); got != 2 {
		t.Fatalf("zoom after setZoom(2) = %v", got)
	}
	w.zoomIn()
	if got := w.getZoom(); got < 2.19 || got > 2.21 {
		t.Fatalf("zoomIn from 2 = %v, want 2.2", got)
	}
	w.zoomOut()
	if got := w.getZoom(); got < 1.99 || got > 2.01 {
		t.Fatalf("zoomOut back = %v, want 2", got)
	}
	w.zoomReset()
	if got := w.getZoom(); got != 1 {
		t.Fatalf("zoomReset = %v", got)
	}
}

func TestGTKWindowBackgroundColour(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1009, nilPointer, LinuxMenuStyleMenuBar)

	w.setBackgroundColour(NewRGBA(255, 0, 0, 255))
	var rgba gdkRGBA
	tWebkitGetBackgroundColor(w.webKitWebView(), uintptr(unsafe.Pointer(&rgba)))
	if rgba.red != 1 || rgba.green != 0 || rgba.blue != 0 || rgba.alpha != 1 {
		t.Fatalf("background = %+v, want opaque red", rgba)
	}
}

func TestGTKWindowMenuStyles(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	menu := menuNew()
	t.Cleanup(func() { g_object_unref(uintptr(menu)) })

	// Primary menu: the menu lives in a header bar set as the window titlebar.
	primary := newTestWindow(t, app, 1010, pointer(menu), LinuxMenuStylePrimaryMenu)
	if tGtkWindowGetTitlebar(uintptr(primary.window)) == 0 {
		t.Fatal("primary menu style did not install a titlebar")
	}

	// Menu bar: no titlebar, and the menu bar is the first child of the box.
	bar := newTestWindow(t, app, 1011, pointer(menu), LinuxMenuStyleMenuBar)
	if tGtkWindowGetTitlebar(uintptr(bar.window)) != 0 {
		t.Fatal("menu bar style installed a titlebar")
	}
	if first := tGtkWidgetGetFirstChild(uintptr(bar.vbox)); first == 0 {
		t.Fatal("menu bar not prepended to the window box")
	}

	// No menu: the webview is the only child of the box.
	plain := newTestWindow(t, app, 1012, nilPointer, LinuxMenuStyleMenuBar)
	if first := tGtkWidgetGetFirstChild(uintptr(plain.vbox)); first != uintptr(plain.webview) {
		t.Fatalf("first child of a menu-less window = %#x, want the webview %#x", first, plain.webview)
	}
}

func TestGTKWindowStateDefaults(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1013, nilPointer, LinuxMenuStyleMenuBar)

	if w.isFullscreen() || w.isMaximised() || w.isMinimised() || w.isFocused() {
		t.Fatal("fresh window reports a non-default state")
	}
	if w.isVisible() {
		t.Fatal("window that was never shown reports visible")
	}
	// Without a surface there is no monitor to report.
	if _, err := w.getScreen(); err == nil {
		t.Fatal("getScreen on an unmapped window returned no error")
	}
	x, y, width, height, scale := w.getCurrentMonitorGeometry()
	if x != -1 || y != -1 || width != -1 || height != -1 || scale != 1 {
		t.Fatalf("unmapped monitor geometry = %d,%d %dx%d @%v", x, y, width, height, scale)
	}
	if w.webKitWebView() != uintptr(w.webview) {
		t.Fatal("webKitWebView is not the window's webview")
	}
}

func TestGTKRegisteredWindowsAndVisibility(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	a := newTestWindow(t, app, 1014, nilPointer, LinuxMenuStyleMenuBar)
	b := newTestWindow(t, app, 1015, nilPointer, LinuxMenuStyleMenuBar)

	windows := app.getWindows()
	if !slices.Contains(windows, a.window) || !slices.Contains(windows, b.window) {
		t.Fatalf("getWindows = %#v, missing created windows", windows)
	}

	app.hideAllWindows()
	if app.isVisible() {
		t.Fatal("app visible after hideAllWindows")
	}
	app.showAllWindows()
	if !app.isVisible() {
		t.Fatal("app not visible after showAllWindows")
	}

	// The active window is one of ours once shown, or the 1 fallback.
	if got := app.getCurrentWindowID(); got != 1 && got != 1014 && got != 1015 {
		t.Fatalf("getCurrentWindowID = %d", got)
	}
}

func TestGTKUnregisterKeepsAppWhileWindowsRemain(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	a := newTestWindow(t, app, 1016, nilPointer, LinuxMenuStyleMenuBar)
	newTestWindow(t, app, 1017, nilPointer, LinuxMenuStyleMenuBar)

	// With another window still registered the application must not run
	// destroy(); the map drops only the window that went away.
	app.unregisterWindow(windowPointer(a.window))
	app.windowMapLock.Lock()
	remaining := len(app.windowMap)
	app.windowMapLock.Unlock()
	if remaining != 1 {
		t.Fatalf("window map holds %d windows, want 1", remaining)
	}
}

func TestGTKAppVersionStrings(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	info := map[string]string{}
	app.appendGTKVersion(info)
	if !strings.HasPrefix(info["GTK"], "4.") || !strings.HasPrefix(info["WebKit"], "2.") {
		t.Fatalf("version info = %v", info)
	}
	if env := (&App{}).platformEnvironment(); env["gtk4-compiled"] == nil {
		t.Fatalf("platformEnvironment has no gtk4-compiled key: %v", env)
	}
}

// Screen tests ------------------------------------------------------------

func TestGTKScreens(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	requireGTKDisplay(t)
	screens, err := getScreens(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(screens) == 0 {
		t.Fatal("no screens reported")
	}
	s := screens[0]
	if !s.IsPrimary || s.ID != "0" || s.Size.Width <= 0 || s.Size.Height <= 0 {
		t.Fatalf("primary screen = %+v", s)
	}
	if s.PhysicalBounds.Width != int(float64(s.Size.Width)*float64(s.ScaleFactor)) {
		t.Fatalf("physical width %d inconsistent with %d at scale %v",
			s.PhysicalBounds.Width, s.Size.Width, s.ScaleFactor)
	}
	if s.Bounds != s.WorkArea {
		t.Logf("work area differs from bounds: %+v vs %+v", s.Bounds, s.WorkArea)
	}

	primary, err := getPrimaryScreen()
	if err != nil || primary.Size != s.Size || !primary.IsPrimary {
		t.Fatalf("getPrimaryScreen = %+v, %v", primary, err)
	}

	if byIndex := getScreenByIndex(gdk_display_get_default(), 0); byIndex == nil || byIndex.Size != s.Size {
		t.Fatalf("getScreenByIndex(0) = %+v", byIndex)
	}
	if getScreenByIndex(gdk_display_get_default(), 999) != nil {
		t.Fatal("getScreenByIndex past the end returned a screen")
	}

	cx, cy, centreScreen := getMousePosition()
	if centreScreen == nil || cx != s.X+s.Size.Width/2 || cy != s.Y+s.Size.Height/2 {
		t.Fatalf("getMousePosition = %d,%d (%+v), want the centre of %+v", cx, cy, centreScreen, s)
	}
}

func TestGTKMonitorScaleAndGeometry(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	requireGTKDisplay(t)
	display := gdk_display_get_default()
	monitor := g_list_model_get_item(gdk_display_get_monitors(display), 0)
	if monitor == 0 {
		t.Fatal("no monitor")
	}
	defer g_object_unref(monitor)

	if scale := monitorScale(monitor); scale < 1 {
		t.Fatalf("monitorScale = %v", scale)
	}
	geometry := monitorGeometry(monitor)
	if geometry.width <= 0 || geometry.height <= 0 {
		t.Fatalf("monitorGeometry = %+v", geometry)
	}
}

// Menu tests --------------------------------------------------------------

func TestMenuAppendHonoursHiddenItems(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	parentNative := pointer(menuNew())
	t.Cleanup(func() { g_object_unref(uintptr(parentNative)) })
	parent := testMenuImpl(parentNative)

	a := testMenuItem("alpha")
	b := testMenuItem("beta")
	b.hidden = true
	menuAppend(parent, a, false)
	menuAppend(parent, b, true)

	if n := tGMenuModelGetNItems(uintptr(parentNative)); n != 1 {
		t.Fatalf("menu holds %d items, want 1 (beta is hidden)", n)
	}
	if got := testMenuItemLabel(uintptr(parentNative), 0); got != "alpha" {
		t.Fatalf("item 0 = %q", got)
	}

	// Showing the hidden item inserts it at its registered position.
	bImpl := b.impl.(*linuxMenuItem)
	menuItemSetHidden(bImpl, false)
	if n := tGMenuModelGetNItems(uintptr(parentNative)); n != 2 {
		t.Fatalf("after show, menu holds %d items, want 2", n)
	}
	if got := testMenuItemLabel(uintptr(parentNative), 1); got != "beta" {
		t.Fatalf("item 1 = %q, want beta", got)
	}

	// Hiding removes it again, and the other item keeps its place.
	menuItemSetHidden(bImpl, true)
	if n := tGMenuModelGetNItems(uintptr(parentNative)); n != 1 {
		t.Fatalf("after hide, menu holds %d items, want 1", n)
	}
	if got := testMenuItemLabel(uintptr(parentNative), 0); got != "alpha" {
		t.Fatalf("remaining item = %q", got)
	}
}

func TestMenuClearResetsPositions(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	parentNative := pointer(menuNew())
	t.Cleanup(func() { g_object_unref(uintptr(parentNative)) })
	parent := testMenuImpl(parentNative)

	first := testMenuItem("one")
	menuAppend(parent, first, false)
	menuAppend(parent, testMenuItem("two"), false)
	if n := tGMenuModelGetNItems(uintptr(parentNative)); n != 2 {
		t.Fatalf("before clear: %d items", n)
	}

	menuClear(parent)
	if n := tGMenuModelGetNItems(uintptr(parentNative)); n != 0 {
		t.Fatalf("after clear: %d items", n)
	}

	// Positions restart at zero, so a rebuilt item is shown at index 0.
	rebuilt := testMenuItem("again")
	menuAppend(parent, rebuilt, false)
	if idx := rebuilt.impl.(*linuxMenuItem).menuIndex; idx != 0 {
		t.Fatalf("rebuilt item index = %d, want 0", idx)
	}
	if got := testMenuItemLabel(uintptr(parentNative), 0); got != "again" {
		t.Fatalf("item 0 after rebuild = %q", got)
	}
}

func TestMenuNilAndUnownedNoOps(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	// Items without a native widget and menus without a native model are
	// skipped rather than dereferenced.
	menuAppend(&Menu{impl: &linuxMenu{}}, testMenuItem("x"), false)
	menuClear(&Menu{})
	menuAppendSection(&Menu{impl: &linuxMenu{}}, menuNewSection())
	menuItemSetLabel(0, "ignored")
	menuItemDestroy(0)
	menuItemSetHidden(&linuxMenuItem{}, true) // no parent menu
	menuSetSubmenu(&MenuItem{impl: &linuxMenuItem{}}, &Menu{impl: &linuxMenu{}})
	menuItemSetChecked(0, true)
	menuItemSetDisabled(0, true)
	if menuItemChecked(0) {
		t.Fatal("NULL widget reported checked")
	}
	if menuItemChecked(pointer(menuNew())) {
		t.Fatal("unregistered widget reported checked")
	}
	if got := menuGetRadioGroup(nil); got != nil {
		t.Fatal("radio group is not nil")
	}
}

func TestMenuSectionsAndSubmenus(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	parentNative := pointer(menuNew())
	t.Cleanup(func() { g_object_unref(uintptr(parentNative)) })
	parent := testMenuImpl(parentNative)

	section := menuNewSection()
	menuAppendSection(parent, section)
	if n := tGMenuModelGetNItems(uintptr(parentNative)); n != 1 {
		t.Fatalf("section not appended: %d items", n)
	}
	if link := tGMenuModelGetItemLink(uintptr(parentNative), 0, "section"); link == 0 {
		t.Fatal("appended item is not a section")
	}

	hiddenItem := testMenuItem("hidden in section")
	hiddenItem.hidden = true
	menuAppendItemToSection(section, hiddenItem)
	shown := testMenuItem("in section")
	menuAppendItemToSection(section, shown)
	sectionModel := tGMenuModelGetItemLink(uintptr(parentNative), 0, "section")
	if n := tGMenuModelGetNItems(sectionModel); n != 1 {
		t.Fatalf("section holds %d visible items, want 1", n)
	}

	sub := testMenuImpl(pointer(menuNew()))
	t.Cleanup(func() { g_object_unref(uintptr(sub.impl.(*linuxMenu).native)) })
	// The submenu is attached before the item is appended, matching the
	// production order in linuxMenu.processMenu: GMenu copies an item when it
	// is appended, so later changes to it are not seen by the menu.
	submenuItem := testMenuItem("More")
	menuSetSubmenu(submenuItem, sub)
	menuAppend(parent, submenuItem, false)
	if link := tGMenuModelGetItemLink(uintptr(parentNative), 1, "submenu"); link == 0 {
		t.Fatal("submenu not attached")
	}
}

func TestMenuCheckboxStateAndEnabled(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	initAppActionGroup()
	item := testCheckMenuItem("Option", false)
	t.Cleanup(func() { removeTestAction(item.id) })
	widget := item.impl.(*linuxMenuItem).native

	if menuItemChecked(widget) {
		t.Fatal("new checkbox starts checked")
	}
	menuItemSetChecked(widget, true)
	if !menuItemChecked(widget) {
		t.Fatal("checkbox not checked after menuItemSetChecked(true)")
	}
	name, _ := lookupActionName(item.id)
	if !getActionState(name) {
		t.Fatal("GAction state not updated")
	}

	menuItemSetDisabled(widget, true)
	if tGActionGroupGetActionEnabled(appActionGroup, name) != 0 {
		t.Fatal("disabled checkbox action still enabled")
	}
	menuItemSetDisabled(widget, false)
	if tGActionGroupGetActionEnabled(appActionGroup, name) == 0 {
		t.Fatal("re-enabled checkbox action still disabled")
	}
}

func TestMenuTextItemActivationEmitsClick(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	initAppActionGroup()
	item := testMenuItem("Click me")
	t.Cleanup(func() { removeTestAction(item.id) })
	attachMenuHandler(item)
	t.Cleanup(func() {
		gtkSignalToMenuItemLock.Lock()
		delete(gtkSignalToMenuItem, item.id)
		gtkSignalToMenuItemLock.Unlock()
	})
	name, _ := lookupActionName(item.id)

	tGActionGroupActivateAction(appActionGroup, name, 0)
	select {
	case id := <-menuItemClicked:
		if id != item.id {
			t.Fatalf("clicked item %d, want %d", id, item.id)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("activating the item emitted no click")
	}
}

func TestMenuCheckboxActivationTogglesState(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	initAppActionGroup()
	item := testCheckMenuItem("Toggle", false)
	t.Cleanup(func() { removeTestAction(item.id) })
	attachMenuHandler(item)
	t.Cleanup(func() {
		gtkSignalToMenuItemLock.Lock()
		delete(gtkSignalToMenuItem, item.id)
		gtkSignalToMenuItemLock.Unlock()
	})
	name, _ := lookupActionName(item.id)
	widget := item.impl.(*linuxMenuItem).native

	// A user click activates the GAction; the handler flips the checked state
	// once and reports the click once.
	tGActionGroupActivateAction(appActionGroup, name, 0)
	drainClicks(t, item.id)
	if !menuItemChecked(widget) {
		t.Fatal("one activation did not check the item")
	}
	tGActionGroupActivateAction(appActionGroup, name, 0)
	drainClicks(t, item.id)
	if menuItemChecked(widget) {
		t.Fatal("second activation did not uncheck the item")
	}
}

func TestMenuRadioGroupActivation(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	initAppActionGroup()
	const groupID, first, second = uint(930001), uint(930011), uint(930012)
	a := testRadioMenuItem("First", groupID, first, first)
	b := testRadioMenuItem("Second", groupID, first, second)
	t.Cleanup(func() {
		removeTestAction(a.id)
		removeTestAction(b.id)
	})
	for _, it := range []*MenuItem{a, b} {
		attachMenuHandler(it)
	}
	t.Cleanup(func() {
		gtkSignalToMenuItemLock.Lock()
		delete(gtkSignalToMenuItem, a.id)
		delete(gtkSignalToMenuItem, b.id)
		gtkSignalToMenuItemLock.Unlock()
	})

	// Selecting the second radio item sets the shared action's state to its
	// target and reports the click for that item.
	action := fmt.Sprintf("radio_group_%d", groupID)
	group := appActionGroup
	target := g_variant_new_string(fmt.Sprintf("%d", second))
	tGActionGroupActivateAction(group, action, target)
	drainClicks(t, second)
	if got := testGStringVariant(tGActionGroupGetActionState(group, action)); got != fmt.Sprintf("%d", second) {
		t.Fatalf("radio group state = %q, want %d", got, second)
	}
}

// TestMenuRadioItemCheckedReadsBack checks that a radio item reads back as
// checked exactly when the group's shared action holds its target, and that
// selecting one item unchecks the others.
func TestMenuRadioItemCheckedReadsBack(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	initAppActionGroup()
	const groupID, first, second = uint(930101), uint(930111), uint(930112)
	a := testRadioMenuItem("A", groupID, first, first)
	b := testRadioMenuItem("B", groupID, first, second)
	t.Cleanup(func() {
		removeTestAction(a.id)
		removeTestAction(b.id)
	})
	widgetA := a.impl.(*linuxMenuItem).native
	widgetB := b.impl.(*linuxMenuItem).native

	// The group was created with first selected, so only A reads back checked.
	if !menuItemChecked(widgetA) || menuItemChecked(widgetB) {
		t.Fatal("initial selection: want A checked and B unchecked")
	}

	// Selecting B moves the group state to B's target.
	action := fmt.Sprintf("radio_group_%d", groupID)
	tGActionGroupActivateAction(appActionGroup, action, g_variant_new_string(fmt.Sprintf("%d", second)))
	if !menuItemChecked(widgetB) {
		t.Fatal("radio item not reported checked after it was selected")
	}
	if menuItemChecked(widgetA) {
		t.Fatal("previously selected radio item still reads checked")
	}

	// Setting A checked selects it; clearing the selected member is a no-op.
	a.impl.setChecked(true)
	if !menuItemChecked(widgetA) || menuItemChecked(widgetB) {
		t.Fatal("setChecked(true) on A did not select it")
	}
	// B is not selected, so clearing it changes nothing.
	b.impl.setChecked(false)
	if menuItemChecked(widgetB) || !menuItemChecked(widgetA) {
		t.Fatal("setChecked(false) on an unselected radio item changed the selection")
	}
	if got := testGStringVariant(tGActionGroupGetActionState(appActionGroup, action)); got != fmt.Sprintf("%d", first) {
		t.Fatalf("group state = %q, want %d", got, first)
	}
}

// TestMenuRadioReactivateSelectedEmitsClick checks that choosing a radio item
// that is already selected still reports a click, as the macOS and Windows
// backends do, and that the selection does not change.
func TestMenuRadioReactivateSelectedEmitsClick(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	initAppActionGroup()
	const groupID, first, second = uint(930201), uint(930211), uint(930212)
	a := testRadioMenuItem("First", groupID, first, first)
	b := testRadioMenuItem("Second", groupID, first, second)
	for _, it := range []*MenuItem{a, b} {
		attachMenuHandler(it)
	}
	t.Cleanup(func() {
		gtkSignalToMenuItemLock.Lock()
		delete(gtkSignalToMenuItem, a.id)
		delete(gtkSignalToMenuItem, b.id)
		gtkSignalToMenuItemLock.Unlock()
	})
	action := fmt.Sprintf("radio_group_%d", groupID)
	activate := func(id uint) {
		tGActionGroupActivateAction(appActionGroup, action, g_variant_new_string(fmt.Sprintf("%d", id)))
	}

	// The already-selected item reports a click on every activation.
	activate(first)
	drainClicks(t, first)
	activate(first)
	drainClicks(t, first)
	if got := testGStringVariant(tGActionGroupGetActionState(appActionGroup, action)); got != fmt.Sprintf("%d", first) {
		t.Fatalf("state after re-selecting = %q, want %d", got, first)
	}

	// Choosing the other item moves the selection and reports it.
	activate(second)
	drainClicks(t, second)
	activate(second)
	drainClicks(t, second)
	if got := testGStringVariant(tGActionGroupGetActionState(appActionGroup, action)); got != fmt.Sprintf("%d", second) {
		t.Fatalf("state after selecting second = %q, want %d", got, second)
	}
}

func TestSetActionEnabledAndStateOnUnknownNames(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	initAppActionGroup()
	setActionEnabled("no-such-action", false)
	setActionState("no-such-action", true)
	if getActionState("no-such-action") {
		t.Fatal("unknown action reported a true state")
	}
}

func TestMenuItemAcceleratorIsRegistered(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	item := testMenuItem("Save")
	t.Cleanup(func() { removeTestAction(item.id) })
	name, _ := lookupActionName(item.id)

	setMenuItemAccelerator(item.id, &accelerator{Key: "s", Modifiers: []modifier{CmdOrCtrlKey, ShiftKey}})
	if got := acceleratorsFor(uintptr(app.application), "app."+name); len(got) != 1 || got[0] != "<Shift><Control>s" {
		t.Fatalf("accelerators = %q, want [<Shift><Control>s]", got)
	}

	// An unknown key produces no accelerator, and nothing is registered.
	setMenuItemAccelerator(item.id, &accelerator{Key: "nope"})
	if got := acceleratorsFor(uintptr(app.application), "app."+name); len(got) != 1 {
		t.Fatalf("unknown key changed the accelerator list: %q", got)
	}
	// nil and unknown ids are ignored.
	setMenuItemAccelerator(item.id, nil)
	setMenuItemAccelerator(0xdead, &accelerator{Key: "s"})
}

// RemoveAccelerator must take the accelerator out of GTK, not only the Go
// field: the shortcut keeps working until GTK drops it.
func TestMenuItemRemoveAcceleratorClearsNativeBinding(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	item := testMenuItem("Save")
	t.Cleanup(func() { removeTestAction(item.id) })
	name, _ := lookupActionName(item.id)
	action := "app." + name

	item.SetAccelerator("CmdOrCtrl+s")
	if got := acceleratorsFor(uintptr(app.application), action); len(got) != 1 {
		t.Fatalf("before RemoveAccelerator: accelerators = %q, want one", got)
	}

	item.RemoveAccelerator()
	if got := acceleratorsFor(uintptr(app.application), action); len(got) != 0 {
		t.Fatalf("after RemoveAccelerator: GTK still holds %q", got)
	}
	if got := item.GetAccelerator(); got != "" {
		t.Fatalf("GetAccelerator after RemoveAccelerator = %q, want empty", got)
	}

	// Setting it again after removal binds the new chord, and only that one.
	item.SetAccelerator("CmdOrCtrl+Shift+k")
	if got := acceleratorsFor(uintptr(app.application), action); len(got) != 1 || got[0] != "<Shift><Control>k" {
		t.Fatalf("after re-set: accelerators = %q, want [<Shift><Control>k]", got)
	}
}

// helpers ----------------------------------------------------------------

// acceleratorsFor reads the accelerator list GTK holds for an action.
func acceleratorsFor(app uintptr, action string) []string {
	arr := tGtkApplicationGetAccelsForAction(app, action)
	if arr == 0 {
		return nil
	}
	var out []string
	for i := uintptr(0); ; i++ {
		p := *(*uintptr)(unsafe.Pointer(arr + i*unsafe.Sizeof(uintptr(0))))
		if p == 0 {
			break
		}
		out = append(out, goString(p))
	}
	return out
}

// removeTestAction drops the GAction and the id mapping a test created.
func removeTestAction(id uint) {
	menuItemsLock.Lock()
	delete(menuItemActions, id)
	menuItemsLock.Unlock()
}

// drainClicks waits for one click for id and fails on anything else.
func drainClicks(t *testing.T, id uint) {
	t.Helper()
	select {
	case got := <-menuItemClicked:
		if got != id {
			t.Fatalf("click for %d, want %d", got, id)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("no click reported for %d", id)
	}
}

// Menu construction and item setters -----------------------------------------

func TestGTKMenuBuildsNativeModel(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	initAppActionGroup()

	menu := NewMenu()
	menu.Add("Open").OnClick(func(*Context) {})
	menu.AddCheckbox("Check", true).OnClick(func(*Context) {})
	menu.AddRadio("R1", false).OnClick(func(*Context) {})
	menu.AddRadio("R2", true).OnClick(func(*Context) {})
	sub := menu.AddSubmenu("Sub")
	sub.Add("Inside").OnClick(func(*Context) {})

	impl := newMenuImpl(menu)
	menu.impl = impl
	impl.processMenu(menu)
	native := uintptr(impl.native)
	if n := tGMenuModelGetNItems(native); n != 5 {
		t.Fatalf("top-level model holds %d items, want 5", n)
	}
	want := []string{"Open", "Check", "R1", "R2", "Sub"}
	for i, label := range want {
		if got := testMenuItemLabel(native, int32(i)); got != label {
			t.Errorf("item %d = %q, want %q", i, got, label)
		}
	}
	if link := tGMenuModelGetItemLink(native, 4, "submenu"); link == 0 {
		t.Fatal("Sub has no submenu link")
	} else if n := tGMenuModelGetNItems(link); n != 1 || testMenuItemLabel(link, 0) != "Inside" {
		t.Fatalf("submenu holds %d items", n)
	}

	// Rebuilding replaces the items instead of appending duplicates (#5464).
	impl.processMenu(menu)
	if n := tGMenuModelGetNItems(native); n != 5 {
		t.Fatalf("after reprocess the model holds %d items, want 5", n)
	}
}

func TestGTKMenuWithSeparatorsUsesSections(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	initAppActionGroup()

	menu := NewMenu()
	menu.Add("First").OnClick(func(*Context) {})
	menu.AddSeparator()
	menu.Add("Second").OnClick(func(*Context) {})

	impl := newMenuImpl(menu)
	menu.impl = impl
	impl.processMenu(menu)
	native := uintptr(impl.native)
	if n := tGMenuModelGetNItems(native); n != 2 {
		t.Fatalf("separator should split the menu into 2 sections, got %d", n)
	}
	for i, label := range []string{"First", "Second"} {
		section := tGMenuModelGetItemLink(native, int32(i), "section")
		if section == 0 {
			t.Fatalf("section %d missing", i)
		}
		if got := testMenuItemLabel(section, 0); got != label {
			t.Errorf("section %d item = %q, want %q", i, got, label)
		}
	}
}

func TestGTKMenuItemSettersReachNativeWidget(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	initAppActionGroup()

	parentNative := pointer(menuNew())
	t.Cleanup(func() { g_object_unref(uintptr(parentNative)) })
	parent := testMenuImpl(parentNative)
	item := testMenuItem("Start")
	menuAppend(parent, item, false)
	impl := item.impl.(*linuxMenuItem)
	name, _ := lookupActionName(item.id)

	impl.setDisabled(true) // InvokeSync runs inline on the test thread
	if tGActionGroupGetActionEnabled(appActionGroup, name) != 0 {
		t.Fatal("setDisabled(true) left the action enabled")
	}
	impl.setDisabled(false)

	impl.setHidden(true)
	if n := tGMenuModelGetNItems(uintptr(parentNative)); n != 0 {
		t.Fatalf("hidden item still in the menu (%d items)", n)
	}
	impl.setHidden(false)
	if n := tGMenuModelGetNItems(uintptr(parentNative)); n != 1 {
		t.Fatalf("shown item missing (%d items)", n)
	}
	// Tooltips and bitmaps are not supported by GMenu; the calls are no-ops.
	impl.setTooltip("tip")
	impl.setBitmap([]byte{1, 2, 3})
}

// X11, clipboard, file chooser, screens ---------------------------------------

func TestGTKX11MoveAndPosition(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1020, nilPointer, LinuxMenuStyleMenuBar)
	w.windowShow()
	// Map the window: X11 needs its surface before it can be moved or queried.
	for i := 0; i < 200 && gtk_native_get_surface(uintptr(w.window)) == 0; i++ {
		g_main_context_iteration(0, 0)
	}
	if _, _, ok := x11WindowForGtkWindow(uintptr(w.window)); !ok {
		t.Skip("no X11 surface for the window (not an X11 display)")
	}
	for i := 0; i < 20; i++ {
		g_main_context_iteration(0, 0)
	}

	w.move(120, 90)
	for i := 0; i < 20; i++ {
		g_main_context_iteration(0, 0)
	}
	if x, y := w.position(); x != 120 || y != 90 {
		t.Fatalf("position after move(120,90) = %d,%d", x, y)
	}
}

func TestGTKAlwaysOnTopIsRecorded(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1021, nilPointer, LinuxMenuStyleMenuBar)

	// The requested state is stored on the window so windowShow can re-apply
	// it: 1 for on, 2 for off, 0 for never set.
	w.setAlwaysOnTop(true)
	if got := g_object_get_data(uintptr(w.window), "wails-always-on-top"); got != 1 {
		t.Fatalf("always-on-top sentinel = %d after enabling, want 1", got)
	}
	w.setAlwaysOnTop(false)
	if got := g_object_get_data(uintptr(w.window), "wails-always-on-top"); got != 2 {
		t.Fatalf("always-on-top sentinel = %d after disabling, want 2", got)
	}
}

func TestGTKMappedWindowReportsMonitor(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1022, nilPointer, LinuxMenuStyleMenuBar)
	w.windowShow()
	for i := 0; i < 200 && gtk_native_get_surface(uintptr(w.window)) == 0; i++ {
		g_main_context_iteration(0, 0)
	}
	screen, err := w.getScreen()
	if err != nil {
		t.Fatalf("getScreen on a mapped window: %v", err)
	}
	if screen.Size.Width <= 0 || screen.ID != "1022" {
		t.Fatalf("screen = %+v", screen)
	}
	x, y, width, height, scale := w.getCurrentMonitorGeometry()
	if width != screen.Size.Width || height != screen.Size.Height || scale < 1 {
		t.Fatalf("monitor geometry %d,%d %dx%d @%v does not match screen %+v", x, y, width, height, scale, screen)
	}
}

func TestGTKClipboardRoundTrip(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	const text = "wails clipboard check"
	clipboardSet(text)
	if got := clipboardGet(); got != text {
		t.Fatalf("clipboardGet = %q, want %q", got, text)
	}
}

func TestGTKFileChooserHelpers(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)

	chooser := gtk_file_chooser_widget_new(gtkFileChooserActionOpen)
	if chooser == 0 {
		t.Fatal("chooser widget not created")
	}
	if !setFileChooserFolder(chooser, "/") {
		t.Fatal("setFileChooserFolder(/) failed")
	}
	addFileChooserFilter(chooser, "Text", "*.txt;*.md")
	if got := acceptedFileChooserPaths(chooser, true); len(got) != 0 {
		t.Fatalf("nothing selected but got paths %q", got)
	}
	if got := acceptedFileChooserPaths(0, false); got != nil {
		t.Fatalf("nil chooser returned %q", got)
	}
}

func TestGTKFinishCallbacksDeliverSelection(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)

	register := func(id uint) chan string {
		ch := make(chan string, 100)
		dialogRequestMutex.Lock()
		fileDialogCallbacks[id] = ch
		dialogRequestMutex.Unlock()
		return ch
	}
	drain := func(ch chan string) []string {
		var got []string
		for p := range ch {
			got = append(got, p)
		}
		return got
	}

	// Single result: the finish function returns one GFile.
	id := nextDialogRequestID()
	ch := register(id)
	finishSingleFile(func(_, _, _ uintptr) uintptr {
		return g_file_new_for_path("/single/file.txt")
	}, 0, 0, id)
	if got := drain(ch); !slices.Equal(got, []string{"/single/file.txt"}) {
		t.Fatalf("single result = %q", got)
	}

	// Error result: a GError is reported, so the request is cancelled.
	id = nextDialogRequestID()
	ch = register(id)
	finishSingleFile(func(_, _, gerr uintptr) uintptr {
		*(*uintptr)(unsafe.Pointer(gerr)) = g_error_new_literal(1, 0, "cancelled")
		return 0
	}, 0, 0, id)
	if got := drain(ch); len(got) != 0 {
		t.Fatalf("error result delivered %q", got)
	}

	// No file and no error: also a cancellation.
	id = nextDialogRequestID()
	ch = register(id)
	finishSingleFile(func(_, _, _ uintptr) uintptr { return 0 }, 0, 0, id)
	if got := drain(ch); len(got) != 0 {
		t.Fatalf("null result delivered %q", got)
	}

	// Multiple results: a list model of GFile objects.
	id = nextDialogRequestID()
	ch = register(id)
	store := g_list_store_new(gtype_Object())
	g_list_store_append(store, g_file_new_for_path("/m/one"))
	g_list_store_append(store, g_file_new_for_path("/m/two"))
	finishMultipleFiles(func(_, _, _ uintptr) uintptr { return store }, 0, 0, id)
	if got := drain(ch); !slices.Equal(got, []string{"/m/one", "/m/two"}) {
		t.Fatalf("multiple result = %q", got)
	}
}

func TestGTKProcessAndCacheScreens(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	if err := app.processAndCacheScreens(); err != nil {
		t.Fatalf("processAndCacheScreens: %v", err)
	}
	primary := globalApplication.Screen.GetPrimary()
	if primary == nil || !primary.IsPrimary || primary.Size.Width <= 0 {
		t.Fatalf("primary screen after caching = %+v", primary)
	}
	if got, err := app.getPrimaryScreen(); err != nil || got == nil {
		t.Fatalf("linuxApp.getPrimaryScreen = %+v, %v", got, err)
	}
}

// TestGTKMenuItemSetLabelAfterAppend checks that setLabel shows on an item that
// is already in a menu. GMenu copies an item when it is appended, so setLabel
// must replace the model's copy. The copy is found by its id attribute, not by
// position, so hidden siblings and separator sections are covered too.
func TestGTKMenuItemSetLabelAfterAppend(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	newTestLinuxApp(t)
	initAppActionGroup()

	parentNative := pointer(menuNew())
	t.Cleanup(func() { g_object_unref(uintptr(parentNative)) })
	item := testMenuItem("Start")
	menuAppend(testMenuImpl(parentNative), item, false)
	item.impl.(*linuxMenuItem).setLabel("Renamed")
	if got := testMenuItemLabel(uintptr(parentNative), 0); got != "Renamed" {
		t.Fatalf("label after setLabel = %q", got)
	}

	// A radio item is copied on append too.
	radio := testRadioMenuItem("Option", 930301, 930311, 930311)
	menuAppend(testMenuImpl(parentNative), radio, false)
	radio.impl.(*linuxMenuItem).setLabel("Option renamed")
	if got := testMenuItemLabel(uintptr(parentNative), 1); got != "Option renamed" {
		t.Fatalf("radio label after setLabel = %q", got)
	}

	// A hidden item is not in the model, so the visible item after it sits at
	// position 0. Finding the model copy by position would replace the wrong
	// item here.
	hiddenParent := pointer(menuNew())
	t.Cleanup(func() { g_object_unref(uintptr(hiddenParent)) })
	hidden := testMenuItem("Hidden")
	menuAppend(testMenuImpl(hiddenParent), hidden, true)
	visible := testMenuItem("Visible")
	menuAppend(testMenuImpl(hiddenParent), visible, false)
	visible.impl.(*linuxMenuItem).setLabel("Shown")
	if n := tGMenuModelGetNItems(uintptr(hiddenParent)); n != 1 {
		t.Fatalf("model has %d items, want 1 (hidden item must stay out)", n)
	}
	if got := testMenuItemLabel(uintptr(hiddenParent), 0); got != "Shown" {
		t.Fatalf("visible label after setLabel = %q", got)
	}
	// Changing a hidden item's label takes effect when it is shown again.
	hidden.impl.(*linuxMenuItem).setLabel("Later")
	hidden.impl.(*linuxMenuItem).setHidden(false)
	if got := testMenuItemLabel(uintptr(hiddenParent), 0); got != "Later" {
		t.Fatalf("shown item label = %q, want Later", got)
	}

	// An item inside a separator section is replaced inside the section.
	sectionParent := pointer(menuNew())
	t.Cleanup(func() { g_object_unref(uintptr(sectionParent)) })
	section := menuNewSection()
	inSection := testMenuItem("Inside")
	menuAppendItemToSection(section, inSection)
	menuAppendSection(testMenuImpl(sectionParent), section)
	inSection.impl.(*linuxMenuItem).setLabel("Inside renamed")
	if got := testMenuItemLabel(uintptr(section), 0); got != "Inside renamed" {
		// The section is the item's parent model, so its copy is replaced there.
		t.Fatalf("sectioned label after setLabel = %q", got)
	}
}

// Runtime SetMenu on a realised window ------------------------------------------

// testBoxChildren lists the children of a GTK4 box in order.
func testBoxChildren(box uintptr) []uintptr {
	var out []uintptr
	for c := tGtkWidgetGetFirstChild(box); c != 0; c = tGtkWidgetGetNextSibling(c) {
		out = append(out, c)
	}
	return out
}

// testMenuBarLabels reads the top-level labels of the model a menu bar shows.
func testMenuBarLabels(bar uintptr) []string {
	model := tGtkPopoverMenuBarGetMenuModel(bar)
	n := tGMenuModelGetNItems(model)
	var out []string
	for i := int32(0); i < n; i++ {
		out = append(out, testMenuItemLabel(model, i))
	}
	return out
}

// newRuntimeMenu builds a menu with one top-level submenu named label, whose
// only item carries accelerator key (with CmdOrCtrl) and is built by Update.
func newRuntimeMenu(label string, key string) *Menu {
	menu := NewMenu()
	menu.AddSubmenu(label).Add(label + " item").SetAccelerator("ctrl+" + key)
	menu.Update()
	return menu
}

// runtimeMenuItem returns the one item of the top-level submenu of a menu
// built by newRuntimeMenu.
func runtimeMenuItem(menu *Menu) *MenuItem {
	return menu.items[0].submenu.items[0]
}

func TestGTKSetMenuOnRunningWindow(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1030, nilPointer, LinuxMenuStyleMenuBar)
	box, webview := uintptr(w.vbox), uintptr(w.webview)

	// expectBar: the box is [menu bar, webview] and the bar shows want.
	expectBar := func(t *testing.T, want ...string) {
		t.Helper()
		kids := testBoxChildren(box)
		if len(kids) != 2 || kids[1] != webview || kids[0] == 0 || kids[0] != uintptr(w.menubar) {
			t.Fatalf("box children = %#x, want [menu bar %#x, webview %#x]", kids, w.menubar, webview)
		}
		if got := testMenuBarLabels(kids[0]); !slices.Equal(got, want) {
			t.Fatalf("menu bar labels = %q, want %q", got, want)
		}
	}
	// expectNoBar: the webview is the only child of the box.
	expectNoBar := func(t *testing.T) {
		t.Helper()
		if kids := testBoxChildren(box); len(kids) != 1 || kids[0] != webview || w.menubar != nilPointer {
			t.Fatalf("box children = %#x (menubar %#x), want only the webview %#x", kids, w.menubar, webview)
		}
	}
	// accels: the accelerators GTK holds for an item's current action.
	accels := func(item *MenuItem) []string {
		name, ok := lookupActionName(item.id)
		if !ok {
			t.Fatalf("item %q has no action", item.label)
		}
		return acceleratorsFor(uintptr(app.application), "app."+name)
	}

	alpha := newRuntimeMenu("Alpha", "1")
	beta := newRuntimeMenu("Beta", "2")
	alphaItem, betaItem := runtimeMenuItem(alpha), runtimeMenuItem(beta)

	expectNoBar(t)
	w.setMenu(alpha)
	expectBar(t, "Alpha")
	if got := accels(alphaItem); !slices.Equal(got, []string{"<Control>1"}) {
		t.Fatalf("alpha accelerators = %q, want [<Control>1]", got)
	}

	// Replacing the bar: the new labels show, the old bar is gone, and the
	// replaced menu's accelerator is unregistered.
	w.setMenu(beta)
	expectBar(t, "Beta")
	if got := accels(betaItem); !slices.Equal(got, []string{"<Control>2"}) {
		t.Fatalf("beta accelerators = %q, want [<Control>2]", got)
	}
	if got := accels(alphaItem); len(got) != 0 {
		t.Fatalf("replaced menu kept accelerators %q", got)
	}

	// Repeated replacement neither adds widgets nor leaves a stale bar.
	for i := 0; i < 10; i++ {
		w.setMenu(alpha)
		expectBar(t, "Alpha")
		w.setMenu(beta)
		expectBar(t, "Beta")
	}
	// Setting the menu that is already shown rebuilds the same bar.
	w.setMenu(beta)
	expectBar(t, "Beta")
	if got := accels(betaItem); !slices.Equal(got, []string{"<Control>2"}) {
		t.Fatalf("beta accelerators after re-set = %q, want [<Control>2]", got)
	}

	// nil removes the bar and its accelerators.
	w.setMenu(nil)
	expectNoBar(t)
	if got := accels(betaItem); len(got) != 0 {
		t.Fatalf("removed menu kept accelerators %q", got)
	}

	// The window takes a menu again after its bar was removed.
	w.setMenu(alpha)
	expectBar(t, "Alpha")
	if got := accels(alphaItem); !slices.Equal(got, []string{"<Control>1"}) {
		t.Fatalf("alpha accelerators after re-add = %q, want [<Control>1]", got)
	}
}

func TestGTKSetMenuNilRemovesStartupBar(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	alpha := newRuntimeMenu("Startup", "3")
	item := runtimeMenuItem(alpha)
	name, _ := lookupActionName(item.id)

	// A window started with a menu bar, as run() builds it.
	native := alpha.impl.(*linuxMenu).native
	w := newTestWindow(t, app, 1031, pointer(native), LinuxMenuStyleMenuBar)
	w.parent.options.Linux.Menu = alpha
	w.menuInstalled = alpha
	if kids := testBoxChildren(uintptr(w.vbox)); len(kids) != 2 {
		t.Fatalf("startup window box has %d children, want menu bar and webview", len(kids))
	}

	w.setMenu(nil)
	if kids := testBoxChildren(uintptr(w.vbox)); len(kids) != 1 || kids[0] != uintptr(w.webview) {
		t.Fatalf("box children after SetMenu(nil) = %#x, want only the webview", kids)
	}
	if got := acceleratorsFor(uintptr(app.application), "app."+name); len(got) != 0 {
		t.Fatalf("startup menu kept accelerators %q after SetMenu(nil)", got)
	}
}

func TestGTKSetMenuPrimaryStyleHeaderBar(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	w := newTestWindow(t, app, 1032, nilPointer, LinuxMenuStylePrimaryMenu)
	// run() takes the style from the options; the runtime path reads it there.
	w.parent.options.Linux.MenuStyle = LinuxMenuStylePrimaryMenu
	alpha := newRuntimeMenu("Primary", "4")

	// The header bar, not the box, carries the menu in primary style.
	w.setMenu(alpha)
	if tGtkWindowGetTitlebar(uintptr(w.window)) == 0 {
		t.Fatal("primary style runtime SetMenu did not install a header bar")
	}
	if kids := testBoxChildren(uintptr(w.vbox)); len(kids) != 1 {
		t.Fatalf("primary style box has %d children, want only the webview", len(kids))
	}
	w.setMenu(nil)
	if tGtkWindowGetTitlebar(uintptr(w.window)) != 0 {
		t.Fatal("SetMenu(nil) left the header bar in place")
	}
	w.setMenu(alpha)
	w.setMenu(alpha)
	if tGtkWindowGetTitlebar(uintptr(w.window)) == 0 {
		t.Fatal("primary style header bar lost after re-set")
	}
}

func TestGTKMenuActionsAreScopedToWindow(t *testing.T) {
	if !runGTKTestInChild(t) {
		return
	}
	app := newTestLinuxApp(t)
	menu := newRuntimeMenu("Shared", "5")
	item := runtimeMenuItem(menu)
	name, ok := lookupActionName(item.id)
	if !ok {
		t.Fatal("missing menu action")
	}
	native := menu.impl.(*linuxMenu).native
	first := newTestWindow(t, app, 1041, pointer(native), LinuxMenuStyleMenuBar)
	second := newTestWindow(t, app, 1042, pointer(native), LinuxMenuStyleMenuBar)
	first.menuInstalled = menu
	second.menuInstalled = menu
	activate := func(w *linuxWebviewWindow) bool {
		return tGtkWidgetActivateActionVariant(uintptr(w.webview), "app."+name, 0) != 0
	}
	if !activate(first) || !activate(second) {
		t.Fatal("shared menu action not available in both windows")
	}
	first.setMenu(nil)
	if activate(first) {
		t.Fatal("removed menu action still available")
	}
	if !activate(second) {
		t.Fatal("other window lost its menu action")
	}
	first.setMenu(menu)
	if !activate(first) {
		t.Fatal("menu action did not return after restoring menu")
	}
}
