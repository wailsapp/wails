//go:build darwin && !ios

#include <Cocoa/Cocoa.h>

@interface StatusItemController : NSObject <NSMenuDelegate>
@property long id;
@property (assign) NSStatusItem *statusItem;
@property (assign) NSMenu *cachedMenu;
@property (strong) NSGestureRecognizer *gestureObserver;
- (void)statusItemClicked:(id)sender;
@end

// Observes mouse-down on the status-item button through the gesture-recognizer
// system so the pre-click callback can attach a menu before the button starts
// native tracking. Replaces an NSEvent local monitor, an input path that
// misses Sidecar/touch-synthesised events on macOS 27 (TN3212). Never claims
// the gesture, so the button's own click handling is unaffected.
@interface WailsStatusItemGestureObserver : NSGestureRecognizer
@property (assign) StatusItemController *controller;
@end

void* systemTrayNew(long id);
void systemTraySetLabel(void* nsStatusItem, char *label);
void systemTraySetANSILabel(void* nsStatusItem, void* attributedString);
void systemTraySetLabelColor(void* nsStatusItem, char *fg, char *bg);
void* createAttributedString(char *title, char *FG, char *BG);
void* appendAttributedString(void* original, char* label, char* fg, char* bg);
NSImage* imageFromBytes(const unsigned char *bytes, int length);
void systemTraySetIcon(void* nsStatusItem, void* nsImage, int position, bool isTemplate);
void systemTrayDestroy(void* nsStatusItem);
void showMenu(void* nsStatusItem, void *nsMenu);
void systemTraySetCachedMenu(void* nsStatusItem, void *nsMenu);
void systemTrayGetBounds(void* nsStatusItem, NSRect *rect, void **screen);
NSRect NSScreen_frame(void* screen);
void windowSetScreen(void* window, void* screen, int yOffset);
int statusBarHeight();
void systemTrayPositionWindow(void* nsStatusItem, void* nsWindow, int offset);