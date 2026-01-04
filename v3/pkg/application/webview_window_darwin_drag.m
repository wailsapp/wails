//go:build darwin && !ios

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>
#import "webview_window_darwin_drag.h"

#import "../events/events_darwin.h"

extern void processDragItems(unsigned int windowId, char** arr, int length, int x, int y);
extern void macosOnDragEnter(unsigned int windowId);
extern void macosOnDragExit(unsigned int windowId);
extern void macosOnDragOver(unsigned int windowId, int x, int y);

@implementation WebviewDrag

// initWithFrame:
- (instancetype)initWithFrame:(NSRect)frameRect {
    self = [super initWithFrame:frameRect];
    if (self) {
        [self registerForDraggedTypes:@[NSFilenamesPboardType]];
    }
    return self;
}

// draggingEntered:
- (NSDragOperation)draggingEntered:(id<NSDraggingInfo>)sender {
    NSPasteboard *pasteboard = [sender draggingPasteboard];
    if ([[pasteboard types] containsObject:NSFilenamesPboardType]) {
        processWindowEvent(self.windowId, EventWindowFileDraggingEntered);
        // Notify JS for hover effects
        macosOnDragEnter(self.windowId);
        return NSDragOperationCopy;
    }
    return NSDragOperationNone;
}

// draggingUpdated:
- (NSDragOperation)draggingUpdated:(id<NSDraggingInfo>)sender {
    NSPasteboard *pasteboard = [sender draggingPasteboard];
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
