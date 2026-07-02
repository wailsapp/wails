//go:build darwin && !ios && !server && !purego

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

// strdupOrNull copies a C string with strdup, tolerating NULL. Used to copy
// the autoreleased buffers returned by -[NSString UTF8String] into malloc'd
// memory the caller owns: the autoreleased buffer only lives until the
// enclosing autorelease pool drains, which can happen before Go reads the
// string in cScreenToScreen (use-after-free, see #5556). The Go side frees
// the copies after conversion.
static const char* strdupOrNull(const char* s) {
	return s != NULL ? strdup(s) : NULL;
}

Screen processScreen(NSScreen* screen){
	Screen returnScreen;
	returnScreen.scaleFactor = screen.backingScaleFactor;

	// NSScreen's native coordinate space is Y-up with (0,0) at the bottom-left
	// of the primary screen. We normalise to Y-down with (0,0) at the top-left
	// of the primary screen so that Bounds matches windowGetPosition /
	// windowSetPosition and the public conventions used by Windows, GTK,
	// Electron and the web. Screens above the primary therefore have negative
	// Y after the flip; Bounds.Y is the screen's top edge.
	NSScreen* primaryScreen = [[NSScreen screens] firstObject];
	if (primaryScreen == NULL) {
		primaryScreen = [NSScreen mainScreen];
	}
	CGFloat primaryHeight = [primaryScreen frame].size.height;

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
	returnScreen.id = strdupOrNull([[NSString stringWithFormat:@"%d", displayID] UTF8String]);

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
	return processScreen(mainScreen);
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
		Screen* returnScreens = malloc(sizeof(Screen) * count);
		for (NSUInteger i = 0; i < count; i++) {
			NSScreen* screen = [screens objectAtIndex:i];
			returnScreens[i] = processScreen(screen);
			returnScreens[i].isPrimary = (i == 0);
		}
		return returnScreens;
	}
}

Screen getScreenForWindow(void* window){
	@autoreleasepool {
		NSScreen* screen = ((NSWindow*)window).screen;
		return processScreen(screen);
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
	return processScreen(associatedScreen);
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

	// NSScreen.frame and visibleFrame return values in points (already DIPs).
	// applyDPIScaling in screenmanager.go expects Physical* fields to be in
	// device pixels and produces Bounds/WorkArea in DIPs by dividing by
	// ScaleFactor. Pre-multiply the point values by backingScaleFactor so the
	// division lands back on the original point values. Without this, bounds
	// on Retina displays are halved (e.g. 1496×967 becomes 748×484).
	sf := float64(screen.scaleFactor)
	toPhysical := func(points C.int) int { return int(float64(points) * sf) }

	return &Screen{
		// Screen.X/Y must mirror Bounds.X/Y: shared code in screenmanager.go
		// (areScreensTouching, calculateScreenPlacement, move) reads the
		// top-level fields alongside Bounds and assumes they agree.
		X: toPhysical(screen.x),
		Y: toPhysical(screen.y),
		Size: Size{
			Width:  int(screen.p_width),
			Height: int(screen.p_height),
		},
		Bounds: Rect{
			X:      toPhysical(screen.x),
			Y:      toPhysical(screen.y),
			Height: toPhysical(screen.height),
			Width:  toPhysical(screen.width),
		},
		PhysicalBounds: Rect{
			X:      toPhysical(screen.x),
			Y:      toPhysical(screen.y),
			Height: toPhysical(screen.height),
			Width:  toPhysical(screen.width),
		},
		WorkArea: Rect{
			X:      toPhysical(screen.w_x),
			Y:      toPhysical(screen.w_y),
			Height: toPhysical(screen.w_height),
			Width:  toPhysical(screen.w_width),
		},
		PhysicalWorkArea: Rect{
			X:      toPhysical(screen.w_x),
			Y:      toPhysical(screen.w_y),
			Height: toPhysical(screen.w_height),
			Width:  toPhysical(screen.w_width),
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
	return m.parent.Screen.LayoutScreens(allScreens())
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
