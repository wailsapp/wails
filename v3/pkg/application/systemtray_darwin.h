//go:build darwin && !ios

#include <Cocoa/Cocoa.h>

@interface StatusItemController : NSObject <NSMenuDelegate>
@property long id;
@property (assign) NSStatusItem *statusItem;
@property (assign) NSMenu *cachedMenu;
@property (strong) id eventMonitor;
// YES once the controller observes NSStatusItem.visible (see
// systemTraySetRemovable); the observer is removed in systemTrayDestroy.
@property BOOL observingVisibility;
- (void)statusItemClicked:(id)sender;
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

// systemTrayCoerceEventType maps the raw NSEvent.type observed inside
// -[StatusItemController statusItemClicked:] to a mouse-button event type
// that Go's processClick can dispatch on. On current macOS betas the
// currentEvent at action time is no longer the originating mouse-down
// (typically a later NSEventTypeMouseMoved), so callers fall back to
// [NSEvent pressedMouseButtons]. Exposed for regression testing (#5752).
int systemTrayCoerceEventType(int rawEventType, unsigned long pressedMouseButtons);

// Tooltip on the status item's button (NSButton.toolTip). NULL clears it.
void systemTraySetTooltip(void* nsStatusItem, const char *tooltip);

// SF Symbol icon (macOS 11+; ignored on older systems). pointSize 0 keeps
// the status bar default; weight 0 keeps the default weight, otherwise it is
// a MacSymbolWeight value (1 = ultraLight ... 9 = black).
void systemTraySetSymbol(void* nsStatusItem, const char *symbolName, double pointSize, int weight, int position);

// User removal (NSStatusItemBehaviorRemovalAllowed) and autosave name.
// Installs a KVO observer on NSStatusItem.visible that reports changes to Go
// through systrayVisibilityCallback.
void systemTraySetRemovable(void* nsStatusItem, bool allowed, const char *autosaveName);

// NSStatusItem.visible
bool systemTrayIsVisible(void* nsStatusItem);
