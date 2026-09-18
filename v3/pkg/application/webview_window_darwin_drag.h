//go:build darwin && !ios && !wails_native

#import <AppKit/AppKit.h>

@interface WebviewDrag : NSView <NSDraggingDestination>
// windowId is nonatomic because setWindowId: is user defined (it registers
// the window's drop types) while the getter stays synthesized.
@property (nonatomic) unsigned int windowId;
// dropTypeMask is the DropType bitmask (1 files, 2 text, 4 URLs, 8 images)
// the view registered pasteboard types for. It is resolved from the Go
// window options when windowId is set.
@property unsigned int dropTypeMask;
@end
