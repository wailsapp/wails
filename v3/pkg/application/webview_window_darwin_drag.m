//go:build darwin

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>
#import "webview_window_darwin_drag.h"

#import "../events/events_darwin.h"

extern void processDragItems(unsigned int windowId, char** arr, int length, int x, int y);

@implementation WailsWebView

- (instancetype)initWithFrame:(NSRect)frameRect {
    self = [super initWithFrame:frameRect];
    if (self) {
        [self registerForDraggedTypes:@[NSFilenamesPboardType]];
    }

    return self;
}

- (void)draggingExited:(id<NSDraggingInfo>)sender {
    processWindowEvent(self.windowId, EventWindowFileDraggingExited);
}

- (BOOL)prepareForDragOperation:(id<NSDraggingInfo>)sender {

   if ( self.dragAndDropType == DragAndDropTypeNone ) {
    return [super prepareForDragOperation: sender];
   }

  if ( self.dragAndDropType == DragAndDropTypeWindow ) {
    return YES;
  }

  return [super prepareForDragOperation: sender];
}

- (BOOL)performDragOperation:(id<NSDraggingInfo>)sender {
    if ( self.dragAndDropType == DragAndDropTypeNone ) {
        return [super performDragOperation: sender];
    }

    // Get the mouse x and y
    NSPoint mouseLocation = [sender draggingLocation];
    // Translate mouse x and y to be relative to the window
    mouseLocation = [self convertPoint:mouseLocation fromView:nil];
    int x = (int)mouseLocation.x;
    int y = (int)mouseLocation.y;
    NSPasteboard *pasteboard = [sender draggingPasteboard];
    NSArray<NSPasteboardType> * types = [pasteboard types];
    if( !types ) {
        return [super draggingUpdated: sender];
    }
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
    }

    if ( self.dragAndDropType == DragAndDropTypeWindow ) {
        return YES;
    }
    return [super performDragOperation: sender];;
}


- (NSDragOperation)draggingUpdated:(id <NSDraggingInfo>)sender {
    if ( self.dragAndDropType == DragAndDropTypeNone ) {
        return [super draggingUpdated: sender];
    }

    NSPasteboard *pboard = [sender draggingPasteboard];

    // if no types, then we'll just let the WKWebView handle the drag-n-drop as normal
    NSArray<NSPasteboardType> * types = [pboard types];
    if( !types ) {
        return [super draggingUpdated: sender];
    }

    if ( self.dragAndDropType == DragAndDropTypeWindow ) {
        // we should call super as otherwise events will not pass
        [super draggingUpdated: sender];

        // pass NSDragOperationGeneric = 4 to show regular hover for drag and drop. As we want to ignore webkit behaviours that depends on webpage
        return 4;
    }

    return [super draggingUpdated: sender];
}

- (NSDragOperation)draggingEntered:(id<NSDraggingInfo>)sender {

    if ( self.dragAndDropType == DragAndDropTypeNone ) {
        return [super draggingEntered: sender];
    }
    processWindowEvent(self.windowId, EventWindowFileDraggingEntered);
    NSPasteboard *pboard = [sender draggingPasteboard];

    // if no types, then we'll just let the WKWebView handle the drag-n-drop as normal
    NSArray<NSPasteboardType> * types = [pboard types];
    if( !types ) {
        return [super draggingEntered: sender];
    }
    if ( self.dragAndDropType == DragAndDropTypeWindow ) {
        // we should call supper as otherwise events will not pass
        [super draggingEntered: sender];

        // pass NSDragOperationGeneric = 4 to show regular hover for drag and drop. As we want to ignore webkit behaviours that depends on webpage
        return 4;
    }
    return [super draggingEntered: sender];
}
@end

