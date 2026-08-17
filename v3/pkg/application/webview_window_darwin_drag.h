//go:build darwin && !ios && !wails_native

#import <AppKit/AppKit.h>

@interface WebviewDrag : NSView <NSDraggingDestination>
@property unsigned int windowId;
@end
