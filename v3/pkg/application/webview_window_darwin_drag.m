//go:build darwin

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>
#import "webview_window_darwin_drag.h"

#import "../events/events_darwin.h"

extern void processDragItems(unsigned int windowId, char** arr, int length, int x, int y);

@implementation WebviewDrag

// initWithFrame:
- (instancetype)initWithFrame:(NSRect)frameRect {
    self = [super initWithFrame:frameRect];
    if (self) {
        NSLog(@"WebviewDrag: initWithFrame - Registering for dragged types. WindowID (at init if available, might be set later): %u", self.windowId); // self.windowId might not be set here yet.
        [self registerForDraggedTypes:@[NSFilenamesPboardType]];
    }
    return self;
}

// draggingEntered:
- (NSDragOperation)draggingEntered:(id<NSDraggingInfo>)sender {
    NSLog(@"WebviewDrag: draggingEntered called. WindowID: %u", self.windowId);
    NSPasteboard *pasteboard = [sender draggingPasteboard];
    if ([[pasteboard types] containsObject:NSFilenamesPboardType]) {
        NSLog(@"WebviewDrag: draggingEntered - Found NSFilenamesPboardType. Firing EventWindowFileDraggingEntered.");
        processWindowEvent(self.windowId, EventWindowFileDraggingEntered);
        return NSDragOperationCopy;
    }
    NSLog(@"WebviewDrag: draggingEntered - NSFilenamesPboardType NOT found.");
    return NSDragOperationNone;
}

// draggingExited:
- (void)draggingExited:(id<NSDraggingInfo>)sender {
    NSLog(@"WebviewDrag: draggingExited called. WindowID: %u", self.windowId); // Added log
    processWindowEvent(self.windowId, EventWindowFileDraggingExited);
}

// prepareForDragOperation:
- (BOOL)prepareForDragOperation:(id<NSDraggingInfo>)sender {
    NSLog(@"WebviewDrag: prepareForDragOperation called. WindowID: %u", self.windowId); // Added log
    return YES;
}

// performDragOperation:
- (BOOL)performDragOperation:(id<NSDraggingInfo>)sender {
    NSLog(@"WebviewDrag: performDragOperation called. WindowID: %u", self.windowId);
    NSPasteboard *pasteboard = [sender draggingPasteboard];
    processWindowEvent(self.windowId, EventWindowFileDraggingPerformed);
    if ([[pasteboard types] containsObject:NSFilenamesPboardType]) {
        NSArray *files = [pasteboard propertyListForType:NSFilenamesPboardType];
        NSUInteger count = [files count];
        NSLog(@"WebviewDrag: performDragOperation - File count: %lu", (unsigned long)count);
        if (count == 0) {
            NSLog(@"WebviewDrag: performDragOperation - No files found in pasteboard, though type was present.");
            return NO;
        }

        char** cArray = (char**)malloc(count * sizeof(char*));
        for (NSUInteger i = 0; i < count; i++) {
            NSString* str = files[i];
            NSLog(@"WebviewDrag: performDragOperation - File %lu: %@", (unsigned long)i, str);
            cArray[i] = (char*)[str UTF8String];
        }
        
        NSPoint dropPointInWindow = [sender draggingLocation];
        NSPoint dropPointInView = [self convertPoint:dropPointInWindow fromView:nil];
        
        // Get the window's content view height
        NSView *contentView = [self.window contentView];
        CGFloat contentHeight = contentView.frame.size.height;
        
        NSLog(@"WebviewDrag: Self height: %.2f, Content view height: %.2f", self.frame.size.height, contentHeight);
        
        int x = (int)dropPointInView.x;
        // Use the content view height for conversion
        int y = (int)(contentHeight - dropPointInView.y);
        
        processDragItems(self.windowId, cArray, (int)count, x, y);
        free(cArray);
        NSLog(@"WebviewDrag: performDragOperation - Returned from processDragItems.");
        return YES;
    }
    NSLog(@"WebviewDrag: performDragOperation - NSFilenamesPboardType NOT found. Returning NO.");
    return NO;
}

@end
