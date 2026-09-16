//go:build darwin && !ios && !server && !wails_native

#ifndef WebviewWindowDragOut_h
#define WebviewWindowDragOut_h

#import <Cocoa/Cocoa.h>

// Result codes for dragOutBegin.
enum {
    WailsDragOutStarted = 0,
    WailsDragOutNoGesture = 1,
    WailsDragOutNoWebView = 2,
    WailsDragOutNoItems = 3
};

// dragOutBegin starts an NSDraggingSession from the window's WKWebView using
// the mouse event currently being processed. files, promiseNames and
// promiseTypes are arrays of UTF-8 strings; text and png may be NULL.
// offsetX/offsetY position the drag image relative to the cursor (top-left
// origin). operations is an NSDragOperation mask. It must run on the main
// thread.
int dragOutBegin(void* nsWindow, unsigned int windowId, unsigned int sessionId,
                 const char** files, int fileCount,
                 const char** promiseNames, const char** promiseTypes, int promiseCount,
                 const char* text,
                 const void* png, int pngLength,
                 int offsetX, int offsetY,
                 unsigned int operations);

#endif
