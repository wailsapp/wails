//go:build darwin && !ios && !server

#ifndef WailsWebviewWindowSheet_h
#define WailsWebviewWindowSheet_h

#import <Cocoa/Cocoa.h>

// Window sheets. Every function expects the main thread (the Go side
// dispatches through InvokeSync) and tolerates NULL windows.

enum {
    WailsWindowSheetBegan = 0,
    WailsWindowSheetNoWindow = -1,
    WailsWindowSheetSelf = -2,
};

// Attaches sheet to parent with beginSheet:completionHandler: (or
// beginCriticalSheet: when critical). The completion handler reports the
// response code to Go through processMacSheetEnded(sheet, code).
int windowSheetBegin(void* parent, void* sheet, bool critical);

// Ends the sheet through its sheetParent. Returns false when the window is
// not currently a sheet.
bool windowSheetEnd(void* sheet, long code);

// NSWindow.sheetParent and NSWindow.attachedSheet, or NULL.
void* windowSheetParent(void* sheet);
void* windowSheetAttached(void* parent);

#endif
