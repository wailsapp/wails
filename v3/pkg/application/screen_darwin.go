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

// copyCString duplicates an NSString's UTF-8 representation onto the C heap.
// UTF8String returns a pointer owned by the (autoreleased) NSString; storing it
// in a struct that Go reads later is a dangling pointer once the string is
// deallocated (observed as SIGSEGV in C.GoString during display-reconfiguration
// event storms, e.g. on sleep/wake). The caller (Go side) frees the copy.
static const char* copyCString(NSString* s) {
	const char* utf8 = [s UTF8String];
	return strdup(utf8 ? utf8 : "");
}

Screen processScreen(NSScreen* screen){
	Screen returnScreen = {0};
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
	returnScreen.id = copyCString([NSString stringWithFormat:@"%d", displayID]);

	// Get physical monitor size
	NSValue *sizeValue = [screenDictionary objectForKey:@"NSDeviceSize"];
	NSSize physicalSize = sizeValue.sizeValue;
	returnScreen.p_height = physicalSize.height;
	returnScreen.p_width = physicalSize.width;

	// Get the rotation
	double rotation = CGDisplayRotation(displayID);
	returnScreen.rotation = rotation;

	NSString* localizedName = nil;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 101500
	if( @available(macOS 10.15, *) ){
		localizedName = screen.localizedName;
	}
#endif
	// Always an owned, non-NULL pointer (empty string when unavailable),
	// so the Go side can unconditionally GoString + free.
	returnScreen.name = copyCString(localizedName);
	return returnScreen;
}

// Get primary screen
Screen GetPrimaryScreen(){
	// Get primary screen
	NSScreen *mainScreen = [NSScreen mainScreen];
	return processScreen(mainScreen);
}

// getAllScreens returns the screen array and reports its length via outCount.
// The count must be captured in the same snapshot as the allocation: fetching
// it separately (the previous GetNumScreens round-trip) races with display
// reconfiguration and can over-read the malloc'd array.
Screen* getAllScreens(int* outCount) {
	NSArray<NSScreen *> *screens = [NSScreen screens];
	int count = (int)screens.count;
	*outCount = count;
	Screen* returnScreens = malloc(sizeof(Screen) * count);
	for (int i = 0; i < count; i++) {
		NSScreen* screen = [screens objectAtIndex:i];
		returnScreens[i] = processScreen(screen);
		returnScreens[i].isPrimary = (i == 0);
	}
	return returnScreens;
}

Screen getScreenForWindow(void* window){
	NSScreen* screen = ((NSWindow*)window).screen;
	return processScreen(screen);
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
	// NSScreen.frame and visibleFrame return values in points (already DIPs).
	// applyDPIScaling in screenmanager.go expects Physical* fields to be in
	// device pixels and produces Bounds/WorkArea in DIPs by dividing by
	// ScaleFactor. Pre-multiply the point values by backingScaleFactor so the
	// division lands back on the original point values. Without this, bounds
	// on Retina displays are halved (e.g. 1496×967 becomes 748×484).
	sf := float64(screen.scaleFactor)
	toPhysical := func(points C.int) int { return int(float64(points) * sf) }

	// id/name are strdup'd on the C heap (see copyCString); free them once
	// converted so every processScreen caller stays leak-free.
	defer C.free(unsafe.Pointer(screen.id))
	defer C.free(unsafe.Pointer(screen.name))

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
		ID:          C.GoString(screen.id),
		Name:        C.GoString(screen.name),
		IsPrimary:   bool(screen.isPrimary),
		Rotation:    float32(screen.rotation),
	}
}

func (m *macosApp) processAndCacheScreens() error {
	var cCount C.int
	cScreens := C.getAllScreens(&cCount)
	defer C.free(unsafe.Pointer(cScreens))
	numScreens := int(cCount)
	screens := make([]*Screen, numScreens)
	cScreenHeaders := (*[1 << 30]C.Screen)(unsafe.Pointer(cScreens))[:numScreens:numScreens]
	for i := 0; i < numScreens; i++ {
		screens[i] = cScreenToScreen(cScreenHeaders[i])
	}
	return m.parent.Screen.LayoutScreens(screens)
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
