//go:build darwin

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>
#import "webview_window_darwin_drag.h"

#import "../events/events_darwin.h"

extern void processDragItems(unsigned int windowId, char** arr, int length, int x, int y);

@implementation WebviewDrag

- (instancetype)initWithFrame:(NSRect)frameRect {
    self = [super initWithFrame:frameRect];
    if (self) {
        [self registerForDraggedTypes:@[NSFilenamesPboardType]];
    }

    return self;
}

- (NSDragOperation)draggingEntered:(id<NSDraggingInfo>)sender {
    NSPasteboard *pasteboard = [sender draggingPasteboard];
    if ([[pasteboard types] containsObject:NSFilenamesPboardType]) {
        processWindowEvent(self.windowId, EventWindowFileDraggingEntered);
        return NSDragOperationCopy;
    }
    return NSDragOperationNone;
}


- (void)draggingExited:(id<NSDraggingInfo>)sender {
    processWindowEvent(self.windowId, EventWindowFileDraggingExited);
}

- (BOOL)prepareForDragOperation:(id<NSDraggingInfo>)sender {
    return YES;
}

- (BOOL)performDragOperation:(id<NSDraggingInfo>)sender {
    // Get the mouse x and y
    NSPoint mouseLocation = [sender draggingLocation];
    // Translate mouse x and y to be relative to the window
    mouseLocation = [self convertPoint:mouseLocation fromView:nil];
    int x = (int)mouseLocation.x;
    int y = (int)mouseLocation.y;
    // Translate Y coordinate to be relative to the top left corner
    y = self.frame.size.height - y;
    NSPasteboard *pasteboard = [sender draggingPasteboard];
    processWindowEvent(self.windowId, EventWindowFileDraggingPerformed);
    if ([[pasteboard types] containsObject:NSFilenamesPboardType]) {
        NSArray *files = [pasteboard propertyListForType:NSFilenamesPboardType];
		NSUInteger count = [files count];
		char** cArray = (char**)malloc(count * sizeof(char*));
		for (NSUInteger i = 0; i < count; i++) {
			NSString* str = files[i];
			cArray[i] = (char*)[str UTF8String];
		}
		processDragItems(self.windowId, cArray, (int)count, x, y);
		free(cArray);
        return YES;
    }
    return NO;
}


@end

