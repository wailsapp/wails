//go:build darwin && !ios && !server

#import "webview_window_mac_extras_darwin.h"
#import <WebKit/WebKit.h>
#import <objc/runtime.h>
#include <string.h>
#include <stdlib.h>

extern void processMacWindowExtrasExport(unsigned long long requestID, void* bytes, int length, const char* error);

// The WebviewWindow class declares a webView property; the NativeWindow
// class has none. Resolving it dynamically keeps this file free of the
// WebView-only headers so it also compiles in a wails_native build.
static WKWebView* extrasWebView(void* nsWindow) {
    if (nsWindow == NULL) return nil;
    id window = (id)nsWindow;
    if (![window respondsToSelector:@selector(webView)]) return nil;
    id view = [window valueForKey:@"webView"];
    if (![view isKindOfClass:[WKWebView class]]) return nil;
    return (WKWebView*)view;
}

static char* extrasCopyString(NSString* value) {
    if (value == nil) return NULL;
    const char* utf8 = [value UTF8String];
    if (utf8 == NULL) return NULL;
    size_t length = strlen(utf8) + 1;
    char* copy = malloc(length);
    if (copy != NULL) memcpy(copy, utf8, length);
    return copy;
}

static NSString* extrasString(const char* value) {
    if (value == NULL) return @"";
    return [NSString stringWithUTF8String:value];
}

// Represented file and document edited

void windowExtrasSetRepresentedFile(void* nsWindow, const char* path) {
    if (nsWindow == NULL) return;
    NSWindow* window = (NSWindow*)nsWindow;
    NSString* value = extrasString(path);
    if (value.length == 0) {
        window.representedURL = nil;
        return;
    }
    window.representedURL = [NSURL fileURLWithPath:value];
}

char* windowExtrasRepresentedFile(void* nsWindow) {
    if (nsWindow == NULL) return NULL;
    NSWindow* window = (NSWindow*)nsWindow;
    NSURL* url = window.representedURL;
    if (url != nil && url.isFileURL) return extrasCopyString(url.path);
    NSString* filename = window.representedFilename;
    if (filename.length == 0) return NULL;
    return extrasCopyString(filename);
}

void windowExtrasSetDocumentEdited(void* nsWindow, bool edited) {
    if (nsWindow == NULL) return;
    [(NSWindow*)nsWindow setDocumentEdited:edited];
}

bool windowExtrasIsDocumentEdited(void* nsWindow) {
    if (nsWindow == NULL) return false;
    return [(NSWindow*)nsWindow isDocumentEdited];
}

// Subtitle

bool windowExtrasSetSubtitle(void* nsWindow, const char* subtitle) {
    if (nsWindow == NULL) return true;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) {
        ((NSWindow*)nsWindow).subtitle = extrasString(subtitle);
        return true;
    }
#endif
    return false;
}

// Cascading. The cascade point is shared by every window in the process,
// like NSWindowController's document cascading.

static NSPoint extrasCascadePoint;
static BOOL extrasHasCascadePoint = NO;

// cascadeTopLeftFromPoint: places the receiver at the given point and
// returns the point for the next window; with NSZeroPoint it leaves the
// receiver where it is and just returns that next point. Cascading "from"
// another window therefore means placing at the point that window returns.

void windowExtrasCascadeNext(void* nsWindow, void* anchor) {
    if (nsWindow == NULL) return;
    NSWindow* window = (NSWindow*)nsWindow;
    if (!extrasHasCascadePoint) {
        // The first cascade steps away from the anchor (the most recently
        // created window, or the main window) when there is one. With no
        // other window, the window keeps its position and only the next
        // point is recorded.
        NSWindow* other = (NSWindow*)anchor;
        if (other == nil || other == window) {
            other = [NSApp mainWindow] ?: [NSApp keyWindow];
        }
        NSPoint start = NSZeroPoint;
        if (other != nil && other != window) {
            start = [other cascadeTopLeftFromPoint:NSZeroPoint];
        }
        extrasCascadePoint = [window cascadeTopLeftFromPoint:start];
        extrasHasCascadePoint = YES;
        return;
    }
    extrasCascadePoint = [window cascadeTopLeftFromPoint:extrasCascadePoint];
}

void windowExtrasCascadeFrom(void* nsWindow, void* from) {
    if (nsWindow == NULL || from == NULL || nsWindow == from) return;
    NSWindow* window = (NSWindow*)nsWindow;
    NSPoint next = [(NSWindow*)from cascadeTopLeftFromPoint:NSZeroPoint];
    extrasCascadePoint = [window cascadeTopLeftFromPoint:next];
    extrasHasCascadePoint = YES;
}

// Frame autosave

bool windowExtrasSetFrameAutosaveName(void* nsWindow, const char* name) {
    if (nsWindow == NULL) return false;
    NSWindow* window = (NSWindow*)nsWindow;
    NSString* value = extrasString(name);
    if (value.length == 0) {
        [window setFrameAutosaveName:@""];
        return false;
    }
    NSString* key = [@"NSWindow Frame " stringByAppendingString:value];
    bool saved = [[NSUserDefaults standardUserDefaults] objectForKey:key] != nil;
    [window setFrameAutosaveName:value];
    if (saved) {
        [window setFrameUsingName:value];
    }
    return saved;
}

// Window button offsets. AppKit repositions the standard buttons whenever it
// lays the titlebar out (resize, toolbar changes, fullscreen transitions),
// so an observer reapplies the offset after each layout. A reset is detected
// by the button no longer sitting where the observer last put it; that new
// position is AppKit's default, which the offset is added to.

static const void* WailsWindowButtonsObserverKey = &WailsWindowButtonsObserverKey;

@interface WailsWindowButtonsObserver : NSObject
@property (assign) NSWindow* window;
@property NSPoint offset;
@property BOOL enabled;
@property BOOL applying;
@end

@implementation WailsWindowButtonsObserver {
    NSPoint lastApplied[3];
    NSPoint defaults[3];
    BOOL hasLastApplied[3];
}

- (instancetype)initWithWindow:(NSWindow*)window {
    self = [super init];
    if (self) {
        _window = window;
        NSNotificationCenter* center = [NSNotificationCenter defaultCenter];
        for (NSNotificationName name in @[NSWindowDidResizeNotification,
                                          NSWindowDidUpdateNotification,
                                          NSWindowDidExitFullScreenNotification,
                                          NSWindowDidBecomeKeyNotification]) {
            [center addObserver:self selector:@selector(windowDidLayout:) name:name object:window];
        }
    }
    return self;
}

- (void)dealloc {
    [[NSNotificationCenter defaultCenter] removeObserver:self];
    [super dealloc];
}

- (void)windowDidLayout:(NSNotification*)notification {
    [self apply];
}

- (NSButton*)buttonAt:(int)index {
    static const NSWindowButton kinds[3] = {NSWindowCloseButton, NSWindowMiniaturizeButton, NSWindowZoomButton};
    return [self.window standardWindowButton:kinds[index]];
}

- (void)apply {
    if (self.applying || self.window == nil) return;
    if (self.window.styleMask & NSWindowStyleMaskFullScreen) return;
    self.applying = YES;
    for (int i = 0; i < 3; i++) {
        NSButton* button = [self buttonAt:i];
        if (button == nil) continue;
        NSPoint origin = button.frame.origin;
        if (hasLastApplied[i] && NSEqualPoints(origin, lastApplied[i])) continue;
        defaults[i] = origin;
        if (!self.enabled) {
            hasLastApplied[i] = NO;
            continue;
        }
        NSPoint target = NSMakePoint(origin.x + self.offset.x, origin.y - self.offset.y);
        [button setFrameOrigin:target];
        lastApplied[i] = target;
        hasLastApplied[i] = YES;
    }
    self.applying = NO;
}

- (void)restoreDefaults {
    self.applying = YES;
    for (int i = 0; i < 3; i++) {
        NSButton* button = [self buttonAt:i];
        if (button == nil || !hasLastApplied[i]) continue;
        if (NSEqualPoints(button.frame.origin, lastApplied[i])) {
            [button setFrameOrigin:defaults[i]];
        }
        hasLastApplied[i] = NO;
    }
    self.applying = NO;
}
@end

void windowExtrasSetWindowButtonsOffset(void* nsWindow, int x, int y, bool enabled) {
    if (nsWindow == NULL) return;
    NSWindow* window = (NSWindow*)nsWindow;
    WailsWindowButtonsObserver* observer = objc_getAssociatedObject(window, WailsWindowButtonsObserverKey);
    if (observer == nil) {
        if (!enabled) return;
        observer = [[WailsWindowButtonsObserver alloc] initWithWindow:window];
        objc_setAssociatedObject(window, WailsWindowButtonsObserverKey, observer, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
        [observer release];
    }
    [observer restoreDefaults];
    observer.offset = NSMakePoint(x, y);
    observer.enabled = enabled;
    [observer apply];
}

void windowExtrasCloseButtonOrigin(void* nsWindow, double* x, double* y) {
    if (x != NULL) *x = 0;
    if (y != NULL) *y = 0;
    if (nsWindow == NULL) return;
    NSButton* button = [(NSWindow*)nsWindow standardWindowButton:NSWindowCloseButton];
    if (button == nil) return;
    if (x != NULL) *x = button.frame.origin.x;
    if (y != NULL) *y = button.frame.origin.y;
}

// Attention

long windowExtrasRequestAttention(bool critical) {
    return (long)[NSApp requestUserAttention:critical ? NSCriticalRequest : NSInformationalRequest];
}

void windowExtrasCancelAttention(long identifier) {
    [NSApp cancelUserAttentionRequest:(NSInteger)identifier];
}

// Printing. Credit for the WKWebView print operation set-up:
// https://stackoverflow.com/q/33319295

int windowExtrasPrint(void* nsWindow, int orientation, bool hasMargins,
                      double top, double left, double bottom, double right,
                      bool silent, const char* printer, double scale, const char* paper) {
    if (nsWindow == NULL) return WailsWindowExtrasPrintNoWebView;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) {
        NSWindow* window = (NSWindow*)nsWindow;
        WKWebView* webView = extrasWebView(nsWindow);
        if (webView == nil) return WailsWindowExtrasPrintNoWebView;

        NSPrintInfo* info = [[[NSPrintInfo sharedPrintInfo] copy] autorelease];
        info.horizontalPagination = NSPrintingPaginationModeAutomatic;
        info.verticalPagination = NSPrintingPaginationModeAutomatic;
        info.verticallyCentered = YES;
        info.horizontallyCentered = YES;
        if (orientation == 1) {
            info.orientation = NSPaperOrientationPortrait;
        } else if (orientation == 2) {
            info.orientation = NSPaperOrientationLandscape;
        }
        if (hasMargins) {
            info.topMargin = top;
            info.leftMargin = left;
            info.bottomMargin = bottom;
            info.rightMargin = right;
        }
        NSString* printerName = extrasString(printer);
        if (printerName.length > 0) {
            NSPrinter* selected = [NSPrinter printerWithName:printerName];
            if (selected == nil) return WailsWindowExtrasPrintPrinterNotFound;
            info.printer = selected;
        }
        if (scale > 0) {
            info.scalingFactor = scale;
        }
        NSString* paperName = extrasString(paper);
        if (paperName.length > 0) {
            info.paperName = paperName;
        }

        NSPrintOperation* operation = [webView printOperationWithPrintInfo:info];
        operation.showsPrintPanel = !silent;
        operation.showsProgressPanel = !silent;
        // Without this the operation raises; the rect's values are ignored.
        operation.view.frame = webView.bounds;
        // runOperation does not work with WKWebView; the modal form does.
        [operation runOperationModalForWindow:window delegate:nil didRunSelector:nil contextInfo:NULL];
        return WailsWindowExtrasPrintStarted;
    }
#endif
    return WailsWindowExtrasPrintUnsupported;
}

// Export

static void extrasDeliverData(unsigned long long requestID, NSData* data, NSError* error) {
    if (error != nil || data == nil) {
        NSString* message = error != nil ? error.localizedDescription : @"WebKit returned no data";
        processMacWindowExtrasExport(requestID, NULL, 0, [message UTF8String]);
        return;
    }
    void* bytes = NULL;
    if (data.length > 0) {
        bytes = malloc(data.length);
        if (bytes == NULL) {
            processMacWindowExtrasExport(requestID, NULL, 0, "out of memory copying the export");
            return;
        }
        memcpy(bytes, data.bytes, data.length);
    }
    processMacWindowExtrasExport(requestID, bytes, (int)data.length, NULL);
}

int windowExtrasExportPDF(void* nsWindow, unsigned long long requestID,
                          bool hasRect, double x, double y, double width, double height) {
    WKWebView* webView = extrasWebView(nsWindow);
    if (webView == nil) return -1;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) {
        WKPDFConfiguration* configuration = [[[WKPDFConfiguration alloc] init] autorelease];
        if (hasRect) {
            configuration.rect = NSMakeRect(x, y, width, height);
        }
        [webView createPDFWithConfiguration:configuration completionHandler:^(NSData* data, NSError* error) {
            extrasDeliverData(requestID, data, error);
        }];
        return 0;
    }
#endif
    return -3;
}

int windowExtrasSnapshot(void* nsWindow, unsigned long long requestID,
                         bool hasRect, double x, double y, double width, double height,
                         int snapshotWidth) {
    WKWebView* webView = extrasWebView(nsWindow);
    if (webView == nil) return -1;
    WKSnapshotConfiguration* configuration = [[[WKSnapshotConfiguration alloc] init] autorelease];
    if (hasRect) {
        configuration.rect = NSMakeRect(x, y, width, height);
    }
    if (snapshotWidth > 0) {
        configuration.snapshotWidth = @(snapshotWidth);
    }
    [webView takeSnapshotWithConfiguration:configuration completionHandler:^(NSImage* image, NSError* error) {
        if (error != nil || image == nil) {
            extrasDeliverData(requestID, nil, error);
            return;
        }
        CGImageRef cgImage = [image CGImageForProposedRect:NULL context:nil hints:nil];
        if (cgImage == NULL) {
            processMacWindowExtrasExport(requestID, NULL, 0, "the snapshot image has no bitmap");
            return;
        }
        NSBitmapImageRep* rep = [[[NSBitmapImageRep alloc] initWithCGImage:cgImage] autorelease];
        rep.size = image.size;
        NSData* png = [rep representationUsingType:NSBitmapImageFileTypePNG properties:@{}];
        extrasDeliverData(requestID, png, nil);
    }];
    return 0;
}

void windowExtrasFrame(void* nsWindow, double* x, double* y, double* width, double* height) {
    NSRect frame = NSZeroRect;
    if (nsWindow != NULL) frame = ((NSWindow*)nsWindow).frame;
    if (x != NULL) *x = frame.origin.x;
    if (y != NULL) *y = frame.origin.y;
    if (width != NULL) *width = frame.size.width;
    if (height != NULL) *height = frame.size.height;
}
