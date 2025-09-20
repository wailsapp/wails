//go:build ios

package application

/*
#include <stdlib.h>

// Forward declarations of C functions for window management
void* ios_create_webview_with_id(unsigned int wailsID);
void ios_window_exec_js(void* viewController, const char* js);
unsigned int ios_window_get_id(void* viewController);
void ios_window_release_handle(void* viewController);
void ios_window_load_url(void* viewController, const char* url);
void ios_window_set_html(void* viewController, const char* html);
void ios_window_set_background_color(void* viewController, unsigned char r, unsigned char g, unsigned char b, unsigned char a);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// iosWebviewWindow implements the webviewWindowImpl interface for iOS
type iosWebviewWindow struct {
	windowID     uint32         // Wails window ID for tracking
	nativeHandle unsafe.Pointer // Native WailsViewController* pointer
	parent       *WebviewWindow
}

func (w *iosWebviewWindow) center() {}

func (w *iosWebviewWindow) close() {}

func (w *iosWebviewWindow) destroy() {
	// Release the native handle
	if w.nativeHandle != nil {
		C.ios_window_release_handle(w.nativeHandle)
		w.nativeHandle = nil
	}
	w.parent.markAsDestroyed()
}

func (w *iosWebviewWindow) execJS(js string) {
	// Direct call to the native window's JavaScript execution
	if w.nativeHandle != nil {
		// Call the Objective-C method directly on this window's view controller
		C.ios_window_exec_js(w.nativeHandle, C.CString(js))
	}
}

func (w *iosWebviewWindow) flash(_ bool) {}

func (w *iosWebviewWindow) focus() {}

func (w *iosWebviewWindow) forceReload() {}

func (w *iosWebviewWindow) fullscreen() {}

func (w *iosWebviewWindow) getScreen() (*Screen, error) {
	screens, err := getScreens()
	if err != nil || len(screens) == 0 {
		return nil, err
	}
	return screens[0], nil
}

func (w *iosWebviewWindow) getZoom() float64 {
	return 1.0
}

func (w *iosWebviewWindow) handleDragAndDropMessage(_ string) {}

func (w *iosWebviewWindow) hasParent() bool {
	return false
}

func (w *iosWebviewWindow) height() int {
	return 2532 // Default iPhone height
}

func (w *iosWebviewWindow) hide() {}

func (w *iosWebviewWindow) isAlwaysOnTop() bool {
	return false
}

func (w *iosWebviewWindow) isCloseRequested() bool {
	return false
}

func (w *iosWebviewWindow) setCloseRequested(_ bool) {}

func (w *iosWebviewWindow) isFocused() bool {
	return true
}

func (w *iosWebviewWindow) isFullscreen() bool {
	return true
}

func (w *iosWebviewWindow) isMaximised() bool {
	return true
}

func (w *iosWebviewWindow) isMinimised() bool {
	return false
}

func (w *iosWebviewWindow) isNormal() bool {
	return false
}

func (w *iosWebviewWindow) isVisible() bool {
	return true
}

func (w *iosWebviewWindow) maximise() {}

func (w *iosWebviewWindow) minimise() {}

func (w *iosWebviewWindow) openContextMenu(_ *Menu, _ *ContextMenuData) {}

func (w *iosWebviewWindow) openDevTools() {}

func (w *iosWebviewWindow) print() error {
	return nil
}

func (w *iosWebviewWindow) reload() {}

func (w *iosWebviewWindow) relativePosition() (int, int) {
	return 0, 0
}

func (w *iosWebviewWindow) resizable() bool {
	return false
}

func (w *iosWebviewWindow) restore() {}

func (w *iosWebviewWindow) setAbsolutePosition(_ int, _ int) {}

func (w *iosWebviewWindow) setAlwaysOnTop(_ bool) {}

func (w *iosWebviewWindow) setBackgroundColour(col RGBA) {
    if w.nativeHandle == nil {
        return
    }
    C.ios_window_set_background_color(
        w.nativeHandle,
        C.uchar(col.Red), C.uchar(col.Green), C.uchar(col.Blue), C.uchar(col.Alpha),
    )
}

func (w *iosWebviewWindow) setEnabled(_ bool) {}

func (w *iosWebviewWindow) setFrameless(_ bool) {}

func (w *iosWebviewWindow) setFullscreenButtonEnabled(_ bool) {}

func (w *iosWebviewWindow) setMaxSize(_ int, _ int) {}

func (w *iosWebviewWindow) setMinSize(_ int, _ int) {}

func (w *iosWebviewWindow) setRelativePosition(_ int, _ int) {}

func (w *iosWebviewWindow) setResizable(_ bool) {}

func (w *iosWebviewWindow) setSize(_ int, _ int) {}

func (w *iosWebviewWindow) setTitle(_ string) {}

func (w *iosWebviewWindow) setZoom(_ float64) {}

func (w *iosWebviewWindow) show() {}

func (w *iosWebviewWindow) size() (int, int) {
	return 1170, 2532 // Default iPhone size
}

func (w *iosWebviewWindow) toggleDevTools() {}

func (w *iosWebviewWindow) unfullscreen() {}

func (w *iosWebviewWindow) unmaximise() {}

func (w *iosWebviewWindow) unminimise() {}

func (w *iosWebviewWindow) width() int {
	return 1170 // Default iPhone width
}

func (w *iosWebviewWindow) zoom() {}

func (w *iosWebviewWindow) zoomIn() {}

func (w *iosWebviewWindow) zoomOut() {}

func (w *iosWebviewWindow) zoomReset() {}

func (w *iosWebviewWindow) setParent(_ *WebviewWindow) error {
	return nil
}

func (w *iosWebviewWindow) run() {
	fmt.Printf("🔥 iosWebviewWindow.run() called! nativeHandle: %v\n", w.nativeHandle)
	// Create the native WebView when the window runs
	if w.nativeHandle == nil {
		// Get the Wails window ID from the parent
		wailsID := w.parent.ID()
		fmt.Printf("🔥 Creating native WebView with Wails ID: %d\n", wailsID)
		// Create the native WebView with the Wails window ID
		w.nativeHandle = C.ios_create_webview_with_id(C.uint(wailsID))
		if w.nativeHandle != nil {
			// Store the window ID (should match what we passed in)
			w.windowID = uint32(wailsID)
			fmt.Printf("🔥 Native WebView created successfully! Handle: %v\n", w.nativeHandle)
			// Apply initial background colour if set (default white otherwise)
			rgba := w.parent.options.BackgroundColour
			C.ios_window_set_background_color(
				w.nativeHandle,
				C.uchar(rgba.Red), C.uchar(rgba.Green), C.uchar(rgba.Blue), C.uchar(rgba.Alpha),
			)
		} else {
			fmt.Printf("🔴 FAILED to create native WebView!\n")
		}
	} else {
		fmt.Printf("🔥 Native WebView already exists!\n")
	}
}

func (w *iosWebviewWindow) setIgnoreMouseEvents(_ bool) {}

func (w *iosWebviewWindow) setOpacity(_ float32) {}

func (w *iosWebviewWindow) setTheme(_ Theme) {}

func (w *iosWebviewWindow) setPinned(_ bool) {}

func (w *iosWebviewWindow) startResize(_ string) error {
	return nil
}

func (w *iosWebviewWindow) startDrag() error {
	return nil
}

func (w *iosWebviewWindow) enableDevTools() {}

func (w *iosWebviewWindow) disableContextMenu() {}

func (w *iosWebviewWindow) disableDefaultContextMenu() {}

func (w *iosWebviewWindow) setShouldClose(_ func() bool) {}

func (w *iosWebviewWindow) absolutePosition() (int, int) {
	return 0, 0
}

func (w *iosWebviewWindow) startMove() {}

func (w *iosWebviewWindow) windowMenu() *Menu {
	return nil
}

func (w *iosWebviewWindow) setWindowMenu(_ *Menu) {}

func (w *iosWebviewWindow) isIgnoreMouseEvents() bool {
	return false
}

func (w *iosWebviewWindow) flashCancel() {}

func (w *iosWebviewWindow) setFocusable(_ bool) {}

func (w *iosWebviewWindow) bounds() Rect {
	return Rect{
		X:      0,
		Y:      0,
		Width:  1170,
		Height: 2532,
	}
}

func (w *iosWebviewWindow) copy() {
	// iOS copy implementation
}

func (w *iosWebviewWindow) cut() {
	// iOS cut implementation
}

func (w *iosWebviewWindow) paste() {
	// iOS paste implementation
}

func (w *iosWebviewWindow) selectAll() {
	// iOS select all implementation
}

func (w *iosWebviewWindow) undo() {
	// iOS undo implementation
}

func (w *iosWebviewWindow) redo() {
	// iOS redo implementation
}

func (w *iosWebviewWindow) delete() {
	// iOS delete implementation
}

func (w *iosWebviewWindow) getBorderSizes() *LRTB {
	return &LRTB{}
}

func (w *iosWebviewWindow) handleKeyEvent(acceleratorString string) {
	// iOS handle key event
}

func (w *iosWebviewWindow) hideMenuBar() {
	// iOS doesn't have menu bar
}

func (w *iosWebviewWindow) unhideMenuBar() {
	// iOS doesn't have menu bar
}

func (w *iosWebviewWindow) toggleMenuBar() {
	// iOS doesn't have menu bar
}

func (w *iosWebviewWindow) isMenuBarHidden() bool {
	return true // iOS doesn't have menu bar
}

func (w *iosWebviewWindow) nativeWindow() unsafe.Pointer {
	return w.nativeHandle
}

func (w *iosWebviewWindow) on(eventID uint) {
	// iOS event handling
}

func (w *iosWebviewWindow) position() (int, int) {
	return 0, 0
}

func (w *iosWebviewWindow) physicalBounds() Rect {
	return Rect{
		X:      0,
		Y:      0,
		Width:  1170,
		Height: 2532,
	}
}

func (w *iosWebviewWindow) setBounds(bounds Rect) {
	// iOS set bounds - not applicable on mobile
}

func (w *iosWebviewWindow) setMinimiseButtonState(_ ButtonState) {
	// iOS doesn't have minimize buttons like desktop platforms
}

func (w *iosWebviewWindow) setMaximiseButtonState(_ ButtonState) {
	// iOS doesn't have maximize buttons like desktop platforms
}

func (w *iosWebviewWindow) setCloseButtonState(_ ButtonState) {
	// iOS doesn't have close buttons like desktop platforms
}

func (w *iosWebviewWindow) setContentProtection(_ bool) {
	// iOS content protection - could be implemented with UIScreen captured notifications
}

func (w *iosWebviewWindow) setHTML(html string) {
	if w.nativeHandle == nil || html == "" {
		return
	}
	cstr := C.CString(html)
	C.ios_window_set_html(w.nativeHandle, cstr)
	C.free(unsafe.Pointer(cstr))
}

func (w *iosWebviewWindow) setMenu(_ *Menu) {
	// iOS doesn't support window menus like desktop platforms
}

func (w *iosWebviewWindow) setPhysicalBounds(_ Rect) {
	// iOS doesn't support arbitrary window bounds - apps are fullscreen
}

func (w *iosWebviewWindow) setPosition(_ int, _ int) {
	// iOS doesn't support window positioning - apps are fullscreen
}

func (w *iosWebviewWindow) setURL(url string) {
	if w.nativeHandle == nil || url == "" {
		return
	}
	cstr := C.CString(url)
	C.ios_window_load_url(w.nativeHandle, cstr)
	C.free(unsafe.Pointer(cstr))
}

func (w *iosWebviewWindow) showMenuBar() {
	// iOS doesn't have menu bars like desktop platforms
}

func (w *iosWebviewWindow) snapAssist() {
	// iOS doesn't support window snap assist like Windows
}

func newWindowImpl(parent *WebviewWindow) *iosWebviewWindow {
	// Create iOS WebView implementation but don't create native view yet
	// It will be created when run() is called
	return &iosWebviewWindow{
		parent: parent,
	}
}

func newWebviewWindow(options WebviewWindowOptions) *WebviewWindow {
	result := &WebviewWindow{
		options:        options,
		eventListeners: make(map[uint][]*WindowEventListener),
		eventHooks:     make(map[uint][]*WindowEventListener),
		keyBindings:    make(map[string]func(Window)),
		menuBindings:   make(map[string]*MenuItem),
	}
	result.id = result.ID()
	result.impl = newWindowImpl(result)
	return result
}
