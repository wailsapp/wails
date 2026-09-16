//go:build darwin && !ios && !server

#ifndef WailsWebviewWindowMacExtras_h
#define WailsWebviewWindowMacExtras_h

#import <Cocoa/Cocoa.h>

// Window extras: represented file, document edited, subtitle, cascading,
// frame autosave, window button offsets, attention requests, printing and
// WebView export. Every function expects the main thread (the Go side
// dispatches through InvokeSync or InvokeAsync) and tolerates a NULL window.

// An empty path removes the represented file.
void windowExtrasSetRepresentedFile(void* nsWindow, const char* path);
// Returns a malloc'd copy of the represented file path, or NULL when the
// window represents no file; the caller frees it.
char* windowExtrasRepresentedFile(void* nsWindow);

void windowExtrasSetDocumentEdited(void* nsWindow, bool edited);
bool windowExtrasIsDocumentEdited(void* nsWindow);

// Returns false when the running macOS release has no window subtitles.
bool windowExtrasSetSubtitle(void* nsWindow, const char* subtitle);

// Cascades the window from the application's cascade point and advances
// it. The first time, the point is taken after the anchor window (or the
// main window when anchor is NULL); with neither the window stays put.
void windowExtrasCascadeNext(void* nsWindow, void* anchor);
// Cascades the window from another window's top-left corner and records
// the resulting point as the application's cascade point.
void windowExtrasCascadeFrom(void* nsWindow, void* from);

// Sets the frame autosave name. Returns true when a frame saved under that
// name existed and was restored, in which case option positioning should
// be skipped. An empty name stops autosaving.
bool windowExtrasSetFrameAutosaveName(void* nsWindow, const char* name);

// Moves the three standard window buttons by (x, y) points from their
// default position and keeps them there across titlebar layouts. With
// enabled false the default position is restored.
void windowExtrasSetWindowButtonsOffset(void* nsWindow, int x, int y, bool enabled);
// Current origin of the close button in window coordinates, for probes.
void windowExtrasCloseButtonOrigin(void* nsWindow, double* x, double* y);

// requestUserAttention: returns the request identifier used by
// cancelUserAttentionRequest:.
long windowExtrasRequestAttention(bool critical);
void windowExtrasCancelAttention(long identifier);

enum {
    WailsWindowExtrasPrintStarted = 0,
    WailsWindowExtrasPrintNoWebView = -1,
    WailsWindowExtrasPrintPrinterNotFound = -2,
    WailsWindowExtrasPrintUnsupported = -3,
};

// orientation: 0 automatic, 1 portrait, 2 landscape. hasMargins false
// keeps the print settings' margins; scale 0 keeps the scale; empty
// printer and paper names keep the current ones.
int windowExtrasPrint(void* nsWindow, int orientation, bool hasMargins,
                      double top, double left, double bottom, double right,
                      bool silent, const char* printer, double scale, const char* paper);

// Asynchronous exports. The result is delivered to Go through
// processMacWindowExtrasExport(requestID, bytes, length, error) on the main
// thread; bytes is a malloc'd buffer the Go side copies and frees, or NULL
// with a non-NULL error message. Return values: 0 started, -1 no WebView,
// -3 unsupported on this macOS release.
int windowExtrasExportPDF(void* nsWindow, unsigned long long requestID,
                          bool hasRect, double x, double y, double width, double height);
int windowExtrasSnapshot(void* nsWindow, unsigned long long requestID,
                         bool hasRect, double x, double y, double width, double height,
                         int snapshotWidth);

// Frame accessors for probes and cascade checks.
void windowExtrasFrame(void* nsWindow, double* x, double* y, double* width, double* height);

#endif
