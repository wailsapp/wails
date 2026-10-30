//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit -framework AppKit
#import <Foundation/Foundation.h>
#import <CoreGraphics/CoreGraphics.h>
#import <Cocoa/Cocoa.h>
#import <AppKit/AppKit.h>
#include <stdlib.h>
#include <string.h>

typedef struct Screen {
	const char* id;
	const char* name;
	int p_width;
	int p_height;
	int width;
	int height;
	int x;
	int y;
	int w_width;
	int w_height;
	int w_x;
	int w_y;
	float scaleFactor;
	double rotation;
	bool isPrimary;
} Screen;


int GetNumScreens(){
	return [[NSScreen screens] count];
}

// primaryScreenHeight returns the height (in points) of the primary screen,
// used to flip NSScreen's Y-up coordinate space to the Y-down convention.
// Callers resolve it once and pass it into processScreen so the screen list is
// not re-enumerated ([NSScreen screens]) once per screen.
CGFloat primaryScreenHeight(){
	NSScreen* primaryScreen = [[NSScreen screens] firstObject];
	if (primaryScreen == NULL) {
		primaryScreen = [NSScreen mainScreen];
	}
	if (primaryScreen == NULL) {
		return 0;
	}
	return [primaryScreen frame].size.height;
}

// strdupOrNull copies a C string with strdup, tolerating NULL. Used to copy
// the autoreleased buffers returned by -[NSString UTF8String] into malloc'd
// memory the caller owns: the autoreleased buffer only lives until the
// enclosing autorelease pool drains, which can happen before Go reads the
// string in cScreenToScreen (use-after-free, see #5556). The Go side frees
// the copies after conversion.
static const char* strdupOrNull(const char* s) {
	return s != NULL ? strdup(s) : NULL;
}

Screen processScreen(NSScreen* screen, CGFloat primaryHeight){
	// Zero-initialise so every field has a defined value: not all callers set
	// isPrimary (only getAllScreens does), and the availability-gated name
	// assignment may be skipped, leaving Go to read uninitialised pointers.
	Screen returnScreen = {0};

	// A window can have no associated screen (e.g. minimised or on a
	// disconnected display), so ((NSWindow*)window).screen is NULL. Fall back
	// to the main screen rather than messaging a nil NSScreen and returning a
	// garbage Screen, matching the guard used elsewhere (systemtray_darwin.m).
	if (screen == NULL) {
		screen = [NSScreen mainScreen];
	}
	if (screen == NULL) {
		// No screens attached at all; returnScreen stays zero-initialised.
		return returnScreen;
	}

	returnScreen.scaleFactor = screen.backingScaleFactor;

	// NSScreen's native coordinate space is Y-up with (0,0) at the bottom-left
	// of the primary screen. We normalise to Y-down with (0,0) at the top-left
	// of the primary screen so that Bounds matches windowGetPosition /
	// windowSetPosition and the public conventions used by Windows, GTK,
	// Electron and the web. Screens above the primary therefore have negative
	// Y after the flip; Bounds.Y is the screen's top edge. primaryHeight is
	// resolved once by the caller (see primaryScreenHeight).

	// screen bounds
	returnScreen.height = screen.frame.size.height;
	returnScreen.width = screen.frame.size.width;
	returnScreen.x = screen.frame.origin.x;
	returnScreen.y = primaryHeight - screen.frame.origin.y - screen.frame.size.height;

	// work area
	NSRect workArea = [screen visibleFrame];
	returnScreen.w_height = workArea.size.height;
	returnScreen.w_width = workArea.size.width;
	returnScreen.w_x = workArea.origin.x;
	returnScreen.w_y = primaryHeight - workArea.origin.y - workArea.size.height;


	// adapted from https://stackoverflow.com/a/1237490/4188138
	NSDictionary* screenDictionary = [screen deviceDescription];
	NSNumber* screenID = [screenDictionary objectForKey:@"NSScreenNumber"];
	CGDirectDisplayID displayID = [screenID unsignedIntValue];
	returnScreen.id = strdupOrNull([[NSString stringWithFormat:@"%u", displayID] UTF8String]);

	// Get physical monitor size
	NSValue *sizeValue = [screenDictionary objectForKey:@"NSDeviceSize"];
	NSSize physicalSize = sizeValue.sizeValue;
	returnScreen.p_height = physicalSize.height;
	returnScreen.p_width = physicalSize.width;

	// Get the rotation
	double rotation = CGDisplayRotation(displayID);
	returnScreen.rotation = rotation;

	returnScreen.name = NULL;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101500
	if( @available(macOS 10.15, *) ){
		returnScreen.name = strdupOrNull([screen.localizedName UTF8String]);
	}
#endif
	return returnScreen;
}

// Get primary screen
Screen GetPrimaryScreen(){
	// Get primary screen
	NSScreen *mainScreen = [NSScreen mainScreen];
	return processScreen(mainScreen, primaryScreenHeight());
}

// getAllScreens returns a malloc'd array of Screen and, via outCount, the
// number of entries it contains. The count comes from the same
// [NSScreen screens] snapshot used to size the allocation; callers must not
// query the screen count separately — a display change between the two calls
// would over-read the buffer (freeing garbage id/name pointers) or leak the
// strdup'd strings in unvisited tail entries.
Screen* getAllScreens(int* outCount) {
	// The explicit pool releases the autoreleased objects created during
	// enumeration as soon as it ends: without it they leak when this is
	// called from a Go goroutine thread that has no ambient pool. Only the
	// strdup'd strings in the returned structs survive the pool.
	@autoreleasepool {
		NSArray<NSScreen *> *screens = [NSScreen screens];
		NSUInteger count = screens.count;
		if (outCount != NULL) {
			*outCount = (int)count;
		}
		// Reuse the snapshot above instead of re-enumerating [NSScreen screens];
		// screens[0] is the primary screen (matches isPrimary = (i == 0) below).
		CGFloat primaryHeight = count > 0 ? [[screens objectAtIndex:0] frame].size.height : 0;
		Screen* returnScreens = malloc(sizeof(Screen) * count);
		for (NSUInteger i = 0; i < count; i++) {
			returnScreens[i] = processScreen([screens objectAtIndex:i], primaryHeight);
			returnScreens[i].isPrimary = (i == 0);
		}
		return returnScreens;
	}
}

Screen getScreenForWindow(void* window){
	@autoreleasepool {
		NSScreen* screen = ((NSWindow*)window).screen;
		return processScreen(screen, primaryScreenHeight());
	}
}

// Get the screen for the system tray
Screen getScreenForSystemTray(void* nsStatusItem) {
	NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
	NSRect frame = statusItem.button.frame;
	NSArray<NSScreen *> *screens = NSScreen.screens;
	NSScreen *associatedScreen = nil;

	for (NSScreen *screen in screens) {
		if (NSPointInRect(frame.origin, screen.frame)) {
			associatedScreen = screen;
			break;
		}
	}
	return processScreen(associatedScreen, primaryScreenHeight());
}

void* getWindowForSystray(void* nsStatusItem) {
	NSStatusItem *statusItem = (NSStatusItem *)nsStatusItem;
	return statusItem.button.window;
}


*/
import "C"
import "unsafe"

func cScreenToScreen(screen C.Screen) *Screen {
	// id and name are malloc'd copies made by processScreen (strdupOrNull);
	// this function owns them and must free them exactly once.
	id := C.GoString(screen.id)
	name := C.GoString(screen.name)
	C.free(unsafe.Pointer(screen.id))
	C.free(unsafe.Pointer(screen.name))

	// NSScreen.frame and visibleFrame return points, which processScreen has
	// flipped into the canonical DIP space — exactly what Bounds and WorkArea
	// hold, and the space windowGetPosition reports a window's frame in. They are
	// therefore carried across as they are, and the Physical* rectangles are the
	// same rectangles in device pixels. Scaling the origins here instead, to
	// synthesise the device-pixel input LayoutScreens expects, is what moved a
	// non-primary display out of the space its own windows are measured in; the
	// screens go to layoutDIPScreens (see processAndCacheScreens).
	sf := float64(screen.scaleFactor)
	toPixels := func(points C.int) int { return int(float64(points) * sf) }

	bounds := Rect{
		X:      int(screen.x),
		Y:      int(screen.y),
		Width:  int(screen.width),
		Height: int(screen.height),
	}
	workArea := Rect{
		X:      int(screen.w_x),
		Y:      int(screen.w_y),
		Width:  int(screen.w_width),
		Height: int(screen.w_height),
	}

	return &Screen{
		// Screen.X/Y must mirror Bounds.X/Y: shared code in screenmanager.go
		// (screenNearestPoint, intersects, right, bottom) reads the top-level
		// fields alongside Bounds and assumes they agree.
		X:      bounds.X,
		Y:      bounds.Y,
		Size:   bounds.Size(),
		Bounds: bounds,
		PhysicalBounds: Rect{
			X:      toPixels(screen.x),
			Y:      toPixels(screen.y),
			Width:  toPixels(screen.width),
			Height: toPixels(screen.height),
		},
		WorkArea: workArea,
		PhysicalWorkArea: Rect{
			X:      toPixels(screen.w_x),
			Y:      toPixels(screen.w_y),
			Width:  toPixels(screen.w_width),
			Height: toPixels(screen.w_height),
		},
		ScaleFactor: float32(screen.scaleFactor),
		ID:          id,
		Name:        name,
		IsPrimary:   bool(screen.isPrimary),
		Rotation:    float32(screen.rotation),
	}
}

// allScreens enumerates the attached screens and converts them to Go values.
// It is a free function (rather than inlined in processAndCacheScreens) so
// tests can exercise the C string ownership handover without cgo, which is
// unavailable in test files.
func allScreens() []*Screen {
	var count C.int
	cScreens := C.getAllScreens(&count)
	defer C.free(unsafe.Pointer(cScreens))
	numScreens := int(count)
	screens := make([]*Screen, numScreens)
	cScreenHeaders := (*[1 << 30]C.Screen)(unsafe.Pointer(cScreens))[:numScreens:numScreens]
	for i := 0; i < numScreens; i++ {
		screens[i] = cScreenToScreen(cScreenHeaders[i])
	}
	return screens
}

func (m *macosApp) processAndCacheScreens() error {
	// NSScreen and other AppKit APIs are not thread-safe and must be accessed on
	// the main thread. Application events (including ApplicationDidChangeScreenParameters)
	// are dispatched on background goroutines and can fire several times in quick
	// succession during a display reconfiguration, so without marshalling this
	// enumerates [NSScreen screens] concurrently off the main thread and crashes
	// (SIGSEGV). Running on the main run loop also serialises the burst of events.
	// Guard against InvokeSync deadlocking when we are already on the main thread.
	var screens []*Screen
	if m.isOnMainThread() {
		screens = allScreens()
	} else {
		InvokeSync(func() { screens = allScreens() })
	}
	return m.parent.Screen.layoutDIPScreens(screens)
}

func (m *macosApp) getPrimaryScreen() (*Screen, error) {
	if m.parent.Screen.GetPrimary() == nil {
		if err := m.processAndCacheScreens(); err != nil {
			return nil, err
		}
	}
	return m.parent.Screen.GetPrimary(), nil
}

func (m *macosApp) getScreens() ([]*Screen, error) {
	if len(m.parent.Screen.GetAll()) == 0 {
		if err := m.processAndCacheScreens(); err != nil {
			return nil, err
		}
	}
	return m.parent.Screen.GetAll(), nil
}

func getScreenForWindow(window *macosWebviewWindow) (*Screen, error) {
	cScreen := C.getScreenForWindow(window.nsWindow)
	return cScreenToScreen(cScreen), nil
}

func getScreenForSystray(systray *macosSystemTray) (*Screen, error) {
	// Get the Window for the status item
	// https://stackoverflow.com/a/5875019/4188138
	window := C.getWindowForSystray(systray.nsStatusItem)
	cScreen := C.getScreenForWindow(window)
	return cScreenToScreen(cScreen), nil
}
