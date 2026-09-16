//go:build darwin && !ios && !server && !wails_native

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>
#if __has_include(<UniformTypeIdentifiers/UniformTypeIdentifiers.h>)
#import <UniformTypeIdentifiers/UniformTypeIdentifiers.h>
#endif
#import "webview_window_drag_out_darwin.h"
#import "webview_window_darwin.h"

// Go callbacks (webview_window_drag_out_darwin.go).
extern void macosDragOutEnded(unsigned int windowId, unsigned int sessionId, unsigned int operation);
extern char* macosDragOutWritePromise(unsigned int windowId, unsigned int sessionId, int index, const char* path);

// WailsDragOutSource is the NSDraggingSource for one StartDrag call and the
// delegate for its file promises. It is retained in dragOutSources while
// the session runs and for a grace period afterwards, because the promise
// provider only holds it weakly and Finder can ask for the file after the
// session has reported its end.
@interface WailsDragOutSource : NSObject <NSDraggingSource, NSFilePromiseProviderDelegate>
@property unsigned int windowId;
@property unsigned int sessionId;
@property unsigned int operations;
@property (retain) NSArray<NSString*>* promiseNames;
@end

static NSMutableDictionary<NSNumber*, WailsDragOutSource*>* dragOutSources = nil;
static NSOperationQueue* dragOutPromiseQueue = nil;
static NSEvent* dragOutLastMouseEvent = nil;
static id dragOutMonitor = nil;

// dragOutInstallMonitor records the latest left mouse down/dragged event so
// StartDrag has a fallback when NSApp.currentEvent has moved on to something
// that is not a mouse event.
static void dragOutInstallMonitor(void) {
    if (dragOutMonitor != nil) {
        return;
    }
    NSEventMask mask = NSEventMaskLeftMouseDown | NSEventMaskLeftMouseDragged | NSEventMaskLeftMouseUp;
    dragOutMonitor = [[NSEvent addLocalMonitorForEventsMatchingMask:mask handler:^NSEvent*(NSEvent* event) {
        [dragOutLastMouseEvent release];
        dragOutLastMouseEvent = (event.type == NSEventTypeLeftMouseUp) ? nil : [event retain];
        return event;
    }] retain];
}

static BOOL dragOutEventIsGesture(NSEvent* event, NSWindow* window) {
    if (event == nil) {
        return NO;
    }
    switch (event.type) {
        case NSEventTypeLeftMouseDown:
        case NSEventTypeLeftMouseDragged:
            break;
        default:
            return NO;
    }
    if (event.window != window) {
        return NO;
    }
    return ([NSEvent pressedMouseButtons] & 1) != 0;
}

// dragOutFrame places an image so that (offsetX, offsetY) from its top-left
// corner sits under the cursor, in the view's coordinate system.
static NSRect dragOutFrame(NSView* view, NSPoint location, NSImage* image, int offsetX, int offsetY) {
    NSSize size = image.size;
    if (size.width <= 0 || size.height <= 0) {
        size = NSMakeSize(32, 32);
    }
    CGFloat x = location.x - offsetX;
    CGFloat y = view.isFlipped ? location.y - offsetY : location.y - (size.height - offsetY);
    return NSMakeRect(x, y, size.width, size.height);
}

static NSImage* dragOutTextImage(NSString* text) {
    NSString* shown = text;
    if (shown.length > 60) {
        shown = [[shown substringToIndex:60] stringByAppendingString:@"..."];
    }
    NSDictionary* attributes = @{
        NSFontAttributeName: [NSFont systemFontOfSize:13],
        NSForegroundColorAttributeName: [NSColor labelColor]
    };
    NSSize textSize = [shown sizeWithAttributes:attributes];
    NSSize size = NSMakeSize(ceil(textSize.width) + 16, ceil(textSize.height) + 8);
    NSImage* image = [[[NSImage alloc] initWithSize:size] autorelease];
    [image lockFocus];
    [[[NSColor windowBackgroundColor] colorWithAlphaComponent:0.92] setFill];
    [[NSBezierPath bezierPathWithRoundedRect:NSMakeRect(0, 0, size.width, size.height) xRadius:6 yRadius:6] fill];
    [shown drawAtPoint:NSMakePoint(8, 4) withAttributes:attributes];
    [image unlockFocus];
    return image;
}

// dragOutIconForType looks the UTType class up at runtime so the binary
// does not link UniformTypeIdentifiers (macOS 11+) while the deployment
// target is older.
static NSImage* dragOutIconForType(NSString* uti) {
#if __has_include(<UniformTypeIdentifiers/UniformTypeIdentifiers.h>)
    if (@available(macOS 11.0, *)) {
        Class utTypeClass = NSClassFromString(@"UTType");
        if (utTypeClass != nil && [utTypeClass respondsToSelector:@selector(typeWithIdentifier:)]) {
            UTType* type = [utTypeClass performSelector:@selector(typeWithIdentifier:) withObject:uti];
            if (type != nil) {
                return [[NSWorkspace sharedWorkspace] iconForContentType:type];
            }
        }
    }
#endif
    return [NSImage imageNamed:NSImageNameMultipleDocuments];
}

@implementation WailsDragOutSource

- (void)dealloc {
    [_promiseNames release];
    [super dealloc];
}

- (NSDragOperation)draggingSession:(NSDraggingSession*)session sourceOperationMaskForDraggingContext:(NSDraggingContext)context {
    return (NSDragOperation)self.operations;
}

- (BOOL)ignoreModifierKeysForDraggingSession:(NSDraggingSession*)session {
    return NO;
}

- (void)draggingSession:(NSDraggingSession*)session endedAtPoint:(NSPoint)screenPoint operation:(NSDragOperation)operation {
    macosDragOutEnded(self.windowId, self.sessionId, (unsigned int)operation);
    NSNumber* key = @(self.sessionId);
    // Keep the source (and its promise delegate role) alive for a while:
    // a destination may still be writing promised files.
    dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(120 * NSEC_PER_SEC)), dispatch_get_main_queue(), ^{
        [dragOutSources removeObjectForKey:key];
    });
}

- (NSString*)filePromiseProvider:(NSFilePromiseProvider*)filePromiseProvider fileNameForType:(NSString*)fileType {
    NSUInteger index = [(NSNumber*)filePromiseProvider.userInfo unsignedIntegerValue];
    if (index < self.promiseNames.count) {
        return self.promiseNames[index];
    }
    return @"file";
}

- (void)filePromiseProvider:(NSFilePromiseProvider*)filePromiseProvider writePromiseToURL:(NSURL*)url completionHandler:(void (^)(NSError* _Nullable))completionHandler {
    int index = [(NSNumber*)filePromiseProvider.userInfo intValue];
    char* message = macosDragOutWritePromise(self.windowId, self.sessionId, index, [url.path UTF8String]);
    if (message == NULL) {
        completionHandler(nil);
        return;
    }
    NSString* description = [NSString stringWithUTF8String:message];
    free(message);
    NSError* error = [NSError errorWithDomain:@"io.wails.dragout" code:1 userInfo:@{NSLocalizedDescriptionKey: description ?: @"promise failed"}];
    completionHandler(error);
}

- (NSOperationQueue*)operationQueueForFilePromiseProvider:(NSFilePromiseProvider*)filePromiseProvider {
    if (dragOutPromiseQueue == nil) {
        dragOutPromiseQueue = [[NSOperationQueue alloc] init];
        dragOutPromiseQueue.name = @"io.wails.dragout.promises";
        dragOutPromiseQueue.qualityOfService = NSQualityOfServiceUserInitiated;
    }
    return dragOutPromiseQueue;
}

@end

int dragOutBegin(void* nsWindow, unsigned int windowId, unsigned int sessionId,
                 const char** files, int fileCount,
                 const char** promiseNames, const char** promiseTypes, int promiseCount,
                 const char* text,
                 const void* png, int pngLength,
                 int offsetX, int offsetY,
                 unsigned int operations) {
    dragOutInstallMonitor();
    if (dragOutSources == nil) {
        dragOutSources = [[NSMutableDictionary alloc] init];
    }

    NSWindow<WailsWebviewWindow>* window = (NSWindow<WailsWebviewWindow>*)nsWindow;
    WKWebView* webView = window.webView;
    if (webView == nil) {
        return WailsDragOutNoWebView;
    }

    NSEvent* event = [NSApp currentEvent];
    if (!dragOutEventIsGesture(event, window)) {
        event = dragOutLastMouseEvent;
        if (!dragOutEventIsGesture(event, window)) {
            return WailsDragOutNoGesture;
        }
    }

    NSImage* image = nil;
    if (png != NULL && pngLength > 0) {
        image = [[[NSImage alloc] initWithData:[NSData dataWithBytes:png length:pngLength]] autorelease];
    }
    NSPoint location = [webView convertPoint:event.locationInWindow fromView:nil];

    WailsDragOutSource* source = [[WailsDragOutSource alloc] init];
    source.windowId = windowId;
    source.sessionId = sessionId;
    source.operations = operations;

    NSMutableArray<NSDraggingItem*>* items = [NSMutableArray array];

    for (int i = 0; i < fileCount; i++) {
        NSString* path = [NSString stringWithUTF8String:files[i]];
        if (path == nil) {
            continue;
        }
        NSURL* url = [NSURL fileURLWithPath:path];
        NSDraggingItem* item = [[[NSDraggingItem alloc] initWithPasteboardWriter:url] autorelease];
        NSImage* itemImage = image ?: [[NSWorkspace sharedWorkspace] iconForFile:path];
        [item setDraggingFrame:dragOutFrame(webView, location, itemImage, offsetX, offsetY) contents:itemImage];
        [items addObject:item];
    }

    NSMutableArray<NSString*>* names = [NSMutableArray arrayWithCapacity:promiseCount];
    for (int i = 0; i < promiseCount; i++) {
        NSString* name = [NSString stringWithUTF8String:promiseNames[i]] ?: @"file";
        NSString* uti = [NSString stringWithUTF8String:promiseTypes[i]] ?: @"public.data";
        [names addObject:name];
        NSFilePromiseProvider* provider = [[[NSFilePromiseProvider alloc] initWithFileType:uti delegate:source] autorelease];
        provider.userInfo = @(i);
        NSDraggingItem* item = [[[NSDraggingItem alloc] initWithPasteboardWriter:provider] autorelease];
        NSImage* itemImage = image ?: dragOutIconForType(uti);
        [item setDraggingFrame:dragOutFrame(webView, location, itemImage, offsetX, offsetY) contents:itemImage];
        [items addObject:item];
    }
    source.promiseNames = names;

    if (text != NULL && text[0] != '\0') {
        NSString* string = [NSString stringWithUTF8String:text];
        if (string != nil) {
            NSDraggingItem* item = [[[NSDraggingItem alloc] initWithPasteboardWriter:string] autorelease];
            NSImage* itemImage = image ?: dragOutTextImage(string);
            [item setDraggingFrame:dragOutFrame(webView, location, itemImage, offsetX, offsetY) contents:itemImage];
            [items addObject:item];
        }
    }

    if (items.count == 0) {
        [source release];
        return WailsDragOutNoItems;
    }

    NSDraggingSession* session = [webView beginDraggingSessionWithItems:items event:event source:source];
    if (session == nil) {
        [source release];
        return WailsDragOutNoGesture;
    }
    session.animatesToStartingPositionsOnCancelOrFail = YES;
    dragOutSources[@(sessionId)] = source;
    [source release];
    return WailsDragOutStarted;
}
