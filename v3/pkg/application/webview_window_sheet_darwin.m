//go:build darwin && !ios && !server

#import "webview_window_sheet_darwin.h"

extern void processMacSheetEnded(void* sheet, long code);

int windowSheetBegin(void* parentPtr, void* sheetPtr, bool critical) {
    if (parentPtr == NULL || sheetPtr == NULL) return WailsWindowSheetNoWindow;
    if (parentPtr == sheetPtr) return WailsWindowSheetSelf;
    NSWindow* parent = (NSWindow*)parentPtr;
    NSWindow* sheet = (NSWindow*)sheetPtr;
    void (^completion)(NSModalResponse) = ^(NSModalResponse returnCode) {
        processMacSheetEnded(sheetPtr, (long)returnCode);
    };
    if (critical) {
        [parent beginCriticalSheet:sheet completionHandler:completion];
    } else {
        [parent beginSheet:sheet completionHandler:completion];
    }
    return WailsWindowSheetBegan;
}

bool windowSheetEnd(void* sheetPtr, long code) {
    if (sheetPtr == NULL) return false;
    NSWindow* sheet = (NSWindow*)sheetPtr;
    NSWindow* parent = sheet.sheetParent;
    if (parent == nil) return false;
    [parent endSheet:sheet returnCode:(NSModalResponse)code];
    return true;
}

void* windowSheetParent(void* sheetPtr) {
    if (sheetPtr == NULL) return NULL;
    return (void*)((NSWindow*)sheetPtr).sheetParent;
}

void* windowSheetAttached(void* parentPtr) {
    if (parentPtr == NULL) return NULL;
    return (void*)((NSWindow*)parentPtr).attachedSheet;
}
