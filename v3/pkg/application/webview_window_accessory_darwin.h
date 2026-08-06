//go:build darwin && !ios && !server

#ifndef WailsWebviewWindowAccessory_h
#define WailsWebviewWindowAccessory_h

#import <Cocoa/Cocoa.h>

enum {
    WailsMacAccessoryStyleApplied = 0,
    WailsMacAccessoryStyleUnavailable = -1,
    WailsMacAccessoryStyleInvalidController = -2,
};

// Returns 1 for NSTitlebarAccessoryViewController, 2 for
// NSSplitViewItemAccessoryViewController, and 0 for any other object.
int macAccessoryViewControllerKind(void* controller);

// Returns whether preferredScrollEdgeEffectStyle can be used on this
// controller on the running macOS version.
bool macAccessoryViewControllerSupportsScrollEdgeEffectStyle(void* controller);

// Applies a MacScrollEdgeEffectStyle value. Automatic is a successful no-op
// on systems older than macOS 26.1; explicit styles report unavailable.
int macAccessoryViewControllerSetScrollEdgeEffectStyle(void* controller, int style);

// Returns a MacScrollEdgeEffectStyle value, or one of the negative result
// constants above.
int macAccessoryViewControllerScrollEdgeEffectStyle(void* controller);

#endif
