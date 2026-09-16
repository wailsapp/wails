//go:build darwin && !ios && !server

#import "webview_window_popover_darwin.h"
#import <objc/runtime.h>
#include <string.h>

extern void processMacPopoverClosed(unsigned long long popoverID);

@interface WailsPopoverDelegate : NSObject <NSPopoverDelegate>
@property unsigned long long popoverID;
@end

@implementation WailsPopoverDelegate
- (void)popoverDidClose:(NSNotification*)notification {
    processMacPopoverClosed(self.popoverID);
}
@end

// WailsPopoverContentController owns a plain container view showing the
// content strip: the NSStackView of controls that macAccessoryCreate placed
// in the accessory controller's view. Only the stack moves into the
// popover; the controller and its own container view stay out of any
// window, so MacAccessory.Remove (which looks for a titlebar host through
// the controller's view) stays a safe no-op for popover content. The
// controller is retained here as well, so its control targets outlive an
// early Remove until the popover is released.
@interface WailsPopoverContentController : NSViewController
@property (retain) NSViewController* strip;
@property (retain) NSView* stripView;
@property NSSize contentSize;
@end

@implementation WailsPopoverContentController
- (void)loadView {
    NSView* container = [[NSView alloc] initWithFrame:NSMakeRect(0, 0, self.contentSize.width, self.contentSize.height)];
    self.view = container;
    [container release];
    NSView* stripView = self.stripView;
    if (stripView != nil) {
        stripView.translatesAutoresizingMaskIntoConstraints = NO;
        [container addSubview:stripView];
        [NSLayoutConstraint activateConstraints:@[
            [stripView.leadingAnchor constraintEqualToAnchor:container.leadingAnchor],
            [stripView.trailingAnchor constraintEqualToAnchor:container.trailingAnchor],
            [stripView.centerYAnchor constraintEqualToAnchor:container.centerYAnchor],
        ]];
    }
}
- (void)dealloc {
    [_stripView removeFromSuperview];
    [_stripView release];
    [_strip release];
    [super dealloc];
}
@end

void* popoverCreate(unsigned long long id, double width, double height, int behavior,
                    bool animates, void* contentController) {
    NSPopover* popover = [[NSPopover alloc] init];
    WailsPopoverDelegate* delegate = [[WailsPopoverDelegate alloc] init];
    delegate.popoverID = id;
    popover.delegate = delegate;
    // The popover keeps the delegate alive through an associated object so
    // it is released together with the popover.
    objc_setAssociatedObject(popover, "wailsPopoverDelegate", delegate, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    [delegate release];

    WailsPopoverContentController* content = [[WailsPopoverContentController alloc] init];
    content.contentSize = NSMakeSize(width, height);
    if (contentController != NULL) {
        NSViewController* strip = (NSViewController*)contentController;
        content.strip = strip;
        NSView* stack = strip.view.subviews.firstObject;
        if (stack != nil) {
            [stack removeFromSuperview];
            stack.hidden = strip.view.hidden;
            content.stripView = stack;
        }
    }
    popover.contentViewController = content;
    [content release];
    popover.contentSize = NSMakeSize(width, height);
    popover.behavior = (NSPopoverBehavior)behavior;
    popover.animates = animates ? YES : NO;
    return popover;
}

static int popoverShowRelativeToView(NSPopover* popover, NSView* view, NSRect rect, int edge) {
    if (popover == nil || view == nil) return WailsPopoverNoAnchor;
    [popover showRelativeToRect:rect ofView:view preferredEdge:(NSRectEdge)edge];
    return WailsPopoverShown;
}

int popoverShowInWindow(void* popoverPtr, void* nsWindow, double x, double y, double width, double height, int edge) {
    if (popoverPtr == NULL || nsWindow == NULL) return WailsPopoverNoAnchor;
    NSView* content = ((NSWindow*)nsWindow).contentView;
    if (content == nil) return WailsPopoverNoAnchor;
    // The rectangle arrives with a top-left origin; AppKit views have a
    // bottom-left origin unless flipped.
    NSRect rect = NSMakeRect(x, y, width, height);
    if (!content.isFlipped) {
        rect.origin.y = content.bounds.size.height - y - height;
    }
    return popoverShowRelativeToView((NSPopover*)popoverPtr, content, rect, edge);
}

int popoverShowFromToolbarItem(void* popoverPtr, void* nsWindow, const char* identifier, int edge) {
    if (popoverPtr == NULL || nsWindow == NULL || identifier == NULL) return WailsPopoverNoAnchor;
    NSToolbar* toolbar = ((NSWindow*)nsWindow).toolbar;
    if (toolbar == nil) return WailsPopoverNoAnchor;
    NSString* wanted = [NSString stringWithUTF8String:identifier];
    NSToolbarItem* item = nil;
    for (NSToolbarItem* candidate in toolbar.items) {
        if ([candidate.itemIdentifier isEqualToString:wanted]) {
            item = candidate;
            break;
        }
    }
    if (item == nil) return WailsPopoverNoAnchor;
    NSPopover* popover = (NSPopover*)popoverPtr;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 140000
    if (@available(macOS 14.0, *)) {
        [popover showRelativeToToolbarItem:item];
        return WailsPopoverShown;
    }
#endif
    NSView* view = item.view;
    if (view == nil) return WailsPopoverAnchorUnavailable;
    return popoverShowRelativeToView(popover, view, view.bounds, edge);
}

int popoverShowFromStatusItem(void* popoverPtr, void* nsStatusItem, int edge) {
    if (popoverPtr == NULL || nsStatusItem == NULL) return WailsPopoverNoAnchor;
    NSStatusBarButton* button = ((NSStatusItem*)nsStatusItem).button;
    if (button == nil) return WailsPopoverNoAnchor;
    return popoverShowRelativeToView((NSPopover*)popoverPtr, button, button.bounds, edge);
}

void popoverClose(void* popoverPtr) {
    if (popoverPtr == NULL) return;
    NSPopover* popover = (NSPopover*)popoverPtr;
    if (popover.isShown) [popover performClose:nil];
}

bool popoverIsShown(void* popoverPtr) {
    if (popoverPtr == NULL) return false;
    return ((NSPopover*)popoverPtr).isShown ? true : false;
}

void popoverSetContentSize(void* popoverPtr, double width, double height) {
    if (popoverPtr == NULL) return;
    NSPopover* popover = (NSPopover*)popoverPtr;
    popover.contentSize = NSMakeSize(width, height);
    WailsPopoverContentController* content = (WailsPopoverContentController*)popover.contentViewController;
    if ([content isKindOfClass:[WailsPopoverContentController class]]) {
        content.contentSize = popover.contentSize;
    }
}

void popoverSetBehavior(void* popoverPtr, int behavior) {
    if (popoverPtr == NULL) return;
    ((NSPopover*)popoverPtr).behavior = (NSPopoverBehavior)behavior;
}

void popoverRelease(void* popoverPtr) {
    if (popoverPtr == NULL) return;
    NSPopover* popover = (NSPopover*)popoverPtr;
    if (popover.isShown) [popover close];
    popover.delegate = nil;
    popover.contentViewController = nil;
    [popover release];
}
