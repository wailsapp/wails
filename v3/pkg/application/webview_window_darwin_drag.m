//go:build darwin && !ios && !server && !wails_native

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>
#import "webview_window_darwin_drag.h"

#import "../events/events_darwin.h"

extern void processDragItems(unsigned int windowId, char** arr, int length, int x, int y);
extern void macosOnDragEnter(unsigned int windowId);
extern void macosOnDragExit(unsigned int windowId);
extern void macosOnDragOver(unsigned int windowId, int x, int y);

// Non-file drop support (webview_window_drag_out_darwin.go).
extern int macosDropTypesForWindow(unsigned int windowId);
extern bool macosOnDrop(unsigned int windowId, const char* text, char** urls, int urlCount,
                        char** files, int fileCount, void** images, int* imageLengths, int imageCount,
                        int x, int y);

// DropType bits shared with the Go side.
enum {
    WebviewDropMaskFiles = 1,
    WebviewDropMaskText = 2,
    WebviewDropMaskURLs = 4,
    WebviewDropMaskImages = 8
};

@implementation WebviewDrag

// initWithFrame:
- (instancetype)initWithFrame:(NSRect)frameRect {
    self = [super initWithFrame:frameRect];
    if (self) {
        _dropTypeMask = WebviewDropMaskFiles;
        [self registerForDraggedTypes:@[NSFilenamesPboardType]];
    }
    return self;
}

// setWindowId: resolves the window's DropTypes and registers exactly the
// pasteboard types they need. Only files are registered by default so the
// page's own HTML5 drop handling keeps receiving text, links and images;
// requesting a non-file type routes those drops to Go instead.
- (void)setWindowId:(unsigned int)windowId {
    _windowId = windowId;
    _dropTypeMask = (unsigned int)macosDropTypesForWindow(windowId);
    NSMutableArray* types = [NSMutableArray array];
    if (_dropTypeMask & WebviewDropMaskFiles) {
        [types addObject:NSFilenamesPboardType];
    }
    if (_dropTypeMask & WebviewDropMaskText) {
        [types addObject:NSPasteboardTypeString];
    }
    if (_dropTypeMask & WebviewDropMaskURLs) {
        [types addObject:NSPasteboardTypeURL];
        [types addObject:NSPasteboardTypeFileURL];
    }
    if (_dropTypeMask & WebviewDropMaskImages) {
        [types addObject:NSPasteboardTypePNG];
        [types addObject:NSPasteboardTypeTIFF];
    }
    [self unregisterDraggedTypes];
    if (types.count > 0) {
        [self registerForDraggedTypes:types];
    }
}

// wantsFileDrop: reports whether this drag carries files and the window
// accepts file drops through the overlay.
- (BOOL)wantsFileDrop:(NSPasteboard*)pasteboard {
    return (self.dropTypeMask & WebviewDropMaskFiles) != 0
        && [[pasteboard types] containsObject:NSFilenamesPboardType];
}

// wantsDataDrop: reports whether this drag carries a requested non-file
// type (text, URL or image).
- (BOOL)wantsDataDrop:(NSPasteboard*)pasteboard {
    NSArray* types = [pasteboard types];
    if ((self.dropTypeMask & WebviewDropMaskText) && [types containsObject:NSPasteboardTypeString]) {
        return YES;
    }
    if ((self.dropTypeMask & WebviewDropMaskURLs)
        && ([types containsObject:NSPasteboardTypeURL] || [types containsObject:NSPasteboardTypeFileURL])) {
        return YES;
    }
    if ((self.dropTypeMask & WebviewDropMaskImages)
        && ([types containsObject:NSPasteboardTypePNG] || [types containsObject:NSPasteboardTypeTIFF])) {
        return YES;
    }
    return NO;
}

// dropPointForSender: converts the drag location to content coordinates
// with a top-left origin, matching the file drop pipeline.
- (NSPoint)dropPointForSender:(id<NSDraggingInfo>)sender {
    NSPoint dropPointInWindow = [sender draggingLocation];
    NSPoint dropPointInView = [self convertPoint:dropPointInWindow fromView:nil];
    NSView *contentView = [self.window contentView];
    CGFloat contentHeight = contentView.frame.size.height;
    return NSMakePoint(dropPointInView.x, contentHeight - dropPointInView.y);
}

static char** webviewDropCopyStrings(NSArray<NSString*>* strings, int* outCount) {
    *outCount = 0;
    if (strings.count == 0) {
        return NULL;
    }
    char** result = (char**)calloc(strings.count, sizeof(char*));
    if (result == NULL) {
        return NULL;
    }
    int count = 0;
    for (NSString* string in strings) {
        const char* utf8 = [string UTF8String];
        if (utf8 != NULL) {
            result[count++] = strdup(utf8);
        }
    }
    *outCount = count;
    return result;
}

static void webviewDropFreeStrings(char** strings, int count) {
    if (strings == NULL) {
        return;
    }
    for (int i = 0; i < count; i++) {
        free(strings[i]);
    }
    free(strings);
}

// deliverDataDrop: collects the requested representations and hands them
// to Go. Returns whether a Go listener took the drop.
- (BOOL)deliverDataDrop:(id<NSDraggingInfo>)sender {
    NSPasteboard *pasteboard = [sender draggingPasteboard];
    NSArray* types = [pasteboard types];

    NSString* text = nil;
    if (self.dropTypeMask & WebviewDropMaskText) {
        text = [pasteboard stringForType:NSPasteboardTypeString];
    }

    NSMutableArray<NSString*>* urls = [NSMutableArray array];
    if (self.dropTypeMask & WebviewDropMaskURLs) {
        NSArray* objects = [pasteboard readObjectsForClasses:@[[NSURL class]] options:@{}];
        for (NSURL* url in objects) {
            if (url.absoluteString != nil) {
                [urls addObject:url.absoluteString];
            }
        }
    }

    NSMutableArray<NSString*>* files = [NSMutableArray array];
    if ([types containsObject:NSFilenamesPboardType]) {
        NSArray* names = [pasteboard propertyListForType:NSFilenamesPboardType];
        for (id name in names) {
            if ([name isKindOfClass:[NSString class]]) {
                [files addObject:name];
            }
        }
    }

    NSData* image = nil;
    if (self.dropTypeMask & WebviewDropMaskImages) {
        image = [pasteboard dataForType:NSPasteboardTypePNG];
        if (image == nil) {
            NSData* tiff = [pasteboard dataForType:NSPasteboardTypeTIFF];
            if (tiff != nil) {
                NSBitmapImageRep* rep = [NSBitmapImageRep imageRepWithData:tiff];
                image = [rep representationUsingType:NSBitmapImageFileTypePNG properties:@{}];
            }
        }
    }

    int urlCount = 0;
    char** cURLs = webviewDropCopyStrings(urls, &urlCount);
    int fileCount = 0;
    char** cFiles = webviewDropCopyStrings(files, &fileCount);
    void* imagePointers[1] = { NULL };
    int imageLengths[1] = { 0 };
    int imageCount = 0;
    if (image != nil && image.length > 0) {
        imagePointers[0] = (void*)image.bytes;
        imageLengths[0] = (int)image.length;
        imageCount = 1;
    }

    NSPoint point = [self dropPointForSender:sender];
    bool handled = macosOnDrop(self.windowId, text != nil ? [text UTF8String] : NULL,
                               cURLs, urlCount, cFiles, fileCount,
                               imagePointers, imageLengths, imageCount,
                               (int)point.x, (int)point.y);
    webviewDropFreeStrings(cURLs, urlCount);
    webviewDropFreeStrings(cFiles, fileCount);
    return handled ? YES : NO;
}

// Pass mouse events through to the WKWebView underneath so CSS :hover,
// mousemove, and cursor updates work. Drag-and-drop delivery uses
// registerForDraggedTypes + geometry, not hitTest, so file drop is unaffected.
- (NSView *)hitTest:(NSPoint)point {
    return nil;
}

// draggingEntered:
- (NSDragOperation)draggingEntered:(id<NSDraggingInfo>)sender {
    NSPasteboard *pasteboard = [sender draggingPasteboard];
    if ([self wantsFileDrop:pasteboard]) {
        processWindowEvent(self.windowId, EventWindowFileDraggingEntered);
        // Notify JS for hover effects
        macosOnDragEnter(self.windowId);
        return NSDragOperationCopy;
    }
    if ([self wantsDataDrop:pasteboard]) {
        return NSDragOperationCopy;
    }
    return NSDragOperationNone;
}

// draggingUpdated:
- (NSDragOperation)draggingUpdated:(id<NSDraggingInfo>)sender {
    NSPasteboard *pasteboard = [sender draggingPasteboard];
    if (![self wantsFileDrop:pasteboard]) {
        return [self wantsDataDrop:pasteboard] ? NSDragOperationCopy : NSDragOperationNone;
    }
    if ([[pasteboard types] containsObject:NSFilenamesPboardType]) {
        // Get the current mouse position
        NSPoint dropPointInWindow = [sender draggingLocation];
        NSPoint dropPointInView = [self convertPoint:dropPointInWindow fromView:nil];
        
        // Get the window's content view height for coordinate conversion
        NSView *contentView = [self.window contentView];
        CGFloat contentHeight = contentView.frame.size.height;
        
        int x = (int)dropPointInView.x;
        int y = (int)(contentHeight - dropPointInView.y);
        
        // Notify JS for hover effects
        macosOnDragOver(self.windowId, x, y);
        
        return NSDragOperationCopy;
    }
    return NSDragOperationNone;
}

// draggingExited:
- (void)draggingExited:(id<NSDraggingInfo>)sender {
    processWindowEvent(self.windowId, EventWindowFileDraggingExited);
    // Notify JS to clean up hover effects
    macosOnDragExit(self.windowId);
}

// prepareForDragOperation:
- (BOOL)prepareForDragOperation:(id<NSDraggingInfo>)sender {
    return YES;
}

// performDragOperation:
- (BOOL)performDragOperation:(id<NSDraggingInfo>)sender {
    NSPasteboard *pasteboard = [sender draggingPasteboard];
    if (![self wantsFileDrop:pasteboard]) {
        // Text, URL and image drops go to Go through OnDrop; the file
        // pipeline below is untouched.
        return [self wantsDataDrop:pasteboard] ? [self deliverDataDrop:sender] : NO;
    }
    processWindowEvent(self.windowId, EventWindowFileDraggingPerformed);
    if ([[pasteboard types] containsObject:NSFilenamesPboardType]) {
        NSArray *files = [pasteboard propertyListForType:NSFilenamesPboardType];
        NSUInteger count = [files count];
        if (count == 0) {
            return NO;
        }

        char** cArray = (char**)malloc(count * sizeof(char*));
        for (NSUInteger i = 0; i < count; i++) {
            NSString* str = files[i];
            cArray[i] = (char*)[str UTF8String];
        }
        
        NSPoint dropPointInWindow = [sender draggingLocation];
        NSPoint dropPointInView = [self convertPoint:dropPointInWindow fromView:nil];
        
        // Get the window's content view height
        NSView *contentView = [self.window contentView];
        CGFloat contentHeight = contentView.frame.size.height;
        
        int x = (int)dropPointInView.x;
        // Use the content view height for conversion
        int y = (int)(contentHeight - dropPointInView.y);
        
        processDragItems(self.windowId, cArray, (int)count, x, y);
        free(cArray);
        return YES;
    }
    return NO;
}

@end
