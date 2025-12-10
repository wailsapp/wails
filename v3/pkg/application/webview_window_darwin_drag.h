//go:build darwin && !ios

#import <AppKit/AppKit.h>

@interface WebviewDrag : NSView <NSDraggingDestination>
@property unsigned int windowId;
@end