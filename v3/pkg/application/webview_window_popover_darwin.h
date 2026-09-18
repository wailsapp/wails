//go:build darwin && !ios && !server

#ifndef WailsWebviewWindowPopover_h
#define WailsWebviewWindowPopover_h

#import <Cocoa/Cocoa.h>

// NSPopover hosting a native content strip (a MacAccessory controller built
// with macAccessoryCreate). Every function expects the main thread (the Go
// side dispatches through InvokeSync) and tolerates NULL handles.

enum {
    WailsPopoverShown = 0,
    WailsPopoverNoAnchor = -1,
    WailsPopoverAnchorUnavailable = -2,
};

// Creates a popover with a content view of width x height points. behavior
// is an NSPopoverBehavior value. contentController may be NULL for an empty
// popover; otherwise it is retained as a child of the content controller
// and its view is pinned to the leading and trailing edges and centred
// vertically. Closing is reported to Go through processMacPopoverClosed(id).
// The caller owns the returned handle and releases it with popoverRelease.
void* popoverCreate(unsigned long long id, double width, double height, int behavior,
                    bool animates, void* contentController);

// Shows the popover anchored to a rectangle in the window's content view,
// given in points with the origin at the top left. edge is an NSRectEdge.
int popoverShowInWindow(void* popover, void* nsWindow, double x, double y, double width, double height, int edge);

// Shows the popover anchored to the toolbar item with the identifier in the
// window's toolbar: showRelativeToToolbarItem: on macOS 14 and newer, the
// item's custom view before that.
int popoverShowFromToolbarItem(void* popover, void* nsWindow, const char* identifier, int edge);

// Shows the popover anchored to a status item's button.
int popoverShowFromStatusItem(void* popover, void* nsStatusItem, int edge);

void popoverClose(void* popover);
bool popoverIsShown(void* popover);
void popoverSetContentSize(void* popover, double width, double height);
void popoverSetBehavior(void* popover, int behavior);
// Closes and releases the popover. The content controller passed to
// popoverCreate loses the popover's reference.
void popoverRelease(void* popover);

#endif
