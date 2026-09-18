//go:build darwin && !ios && !server

#ifndef WailsWebviewWindowRestoration_h
#define WailsWebviewWindowRestoration_h

#import <Cocoa/Cocoa.h>

// Window state restoration. Every function expects the main thread (the Go
// side dispatches through InvokeSync) and tolerates a NULL window.

// The NSCoder key under which the JSON string set with
// windowRestorationSetData is encoded in the window's restorable state, and
// decoded again by WailsWindowRestoration at the next launch.
#define WailsRestorationDataKey "wails.restorationData"

// WailsWindowRestoration is the NSWindowRestoration class every restorable
// Wails window names; restoreWindowWithIdentifier:state:completionHandler:
// forwards to Go through processMacWindowRestore(requestID, identifier,
// json) and keeps the completion handler until windowRestorationComplete.
@interface WailsWindowRestoration : NSObject <NSWindowRestoration>
@end

// Makes the window restorable under identifier with WailsWindowRestoration as
// its restoration class. An empty identifier removes it from restoration.
void windowRestorationConfigure(void* nsWindow, const char* identifier);

// Stores the JSON string that window:willEncodeRestorableState: writes for
// the window and asks AppKit to re-encode. NULL or empty clears it.
void windowRestorationSetData(void* nsWindow, const char* json);
// Returns a malloc'd copy of the stored JSON, or NULL; the caller frees it.
char* windowRestorationData(void* nsWindow);

// Finishes a restore request with the recreated window, or with an error
// message when nsWindow is NULL.
void windowRestorationComplete(unsigned long long requestID, void* nsWindow, const char* error);

enum {
    WailsRestorationStateOK = 0,
    WailsRestorationStateNoWebView = -1,
    WailsRestorationStateInvalid = -2,
    WailsRestorationStateUnsupported = -3,
};

// WKWebView.interactionState (macOS 12+). On success *bytes is a malloc'd
// buffer of *length bytes the caller frees (NULL with length 0 when WebKit
// has no state yet).
int windowRestorationInteractionState(void* nsWindow, void** bytes, int* length);
int windowRestorationSetInteractionState(void* nsWindow, const void* bytes, int length);

#endif
