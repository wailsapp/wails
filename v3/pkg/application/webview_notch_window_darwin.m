//go:build darwin && !ios && !server

#import <Cocoa/Cocoa.h>
#import <QuartzCore/QuartzCore.h>
#include <math.h>
#import "webview_notch_window_darwin.h"

// Keep these webview insets synchronized with notchWindowHorizontalInset and
// notchWindowBottomInset in notch_window.go.
static const CGFloat NotchHorizontalInset = 30.0;
static const CGFloat NotchBottomInset = 12.0;
static const CGFloat NotchWingInset = 20.0;
static const CGFloat NotchBottomCornerRadius = 28.0;

@interface WailsNotchBackgroundView : NSView
@end

@implementation WailsNotchBackgroundView

- (BOOL)isOpaque {
    return NO;
}

- (void)drawRect:(NSRect)dirtyRect {
    [super drawRect:dirtyRect];

    CGFloat width = NSWidth(self.bounds);
    CGFloat height = NSHeight(self.bounds);
    CGFloat bodyBottom = NotchBottomCornerRadius;
    NSBezierPath* path = [NSBezierPath bezierPath];

    // The top spans the complete transparent canvas. Each side then curves
    // inward to the narrower body before the tighter lower corners.
    [path moveToPoint:NSMakePoint(0.0, height)];
    [path lineToPoint:NSMakePoint(width, height)];
    [path curveToPoint:NSMakePoint(width - NotchWingInset, height - NotchWingInset)
          controlPoint1:NSMakePoint(width - NotchWingInset * 0.58, height)
          controlPoint2:NSMakePoint(width - NotchWingInset, height - NotchWingInset * 0.42)];
    [path lineToPoint:NSMakePoint(width - NotchWingInset, bodyBottom)];
    [path curveToPoint:NSMakePoint(width - NotchWingInset - NotchBottomCornerRadius, 0.0)
          controlPoint1:NSMakePoint(width - NotchWingInset, NotchBottomCornerRadius * 0.45)
          controlPoint2:NSMakePoint(width - NotchWingInset - NotchBottomCornerRadius * 0.45, 0.0)];
    [path lineToPoint:NSMakePoint(NotchWingInset + NotchBottomCornerRadius, 0.0)];
    [path curveToPoint:NSMakePoint(NotchWingInset, bodyBottom)
          controlPoint1:NSMakePoint(NotchWingInset + NotchBottomCornerRadius * 0.45, 0.0)
          controlPoint2:NSMakePoint(NotchWingInset, NotchBottomCornerRadius * 0.45)];
    [path lineToPoint:NSMakePoint(NotchWingInset, height - NotchWingInset)];
    [path curveToPoint:NSMakePoint(0.0, height)
          controlPoint1:NSMakePoint(NotchWingInset, height - NotchWingInset * 0.42)
          controlPoint2:NSMakePoint(NotchWingInset * 0.58, height)];
    [path closePath];

    [[NSColor blackColor] setFill];
    [path fill];
}

@end

@interface WebviewNotchWindow ()
@property NSUInteger notchAnimationGeneration;
@property BOOL notchRequestedVisible;
@property (retain) NSTrackingArea* notchTrackingArea;
@property (copy) NSString* notchTargetScreenID;
- (NSScreen*)notchTargetScreen;
- (CGFloat)cameraHousingCentreForScreen:(NSScreen*)screen found:(BOOL*)found;
- (NSRect)notchTargetFrame;
- (NSRect)notchHiddenFrameForTarget:(NSRect)target;
@end

@implementation WebviewNotchWindow

- (NSScreen*)notchTargetScreen {
    NSArray<NSScreen*>* screens = [NSScreen screens];
    if (self.notchTargetScreenID.length > 0) {
        for (NSScreen* screen in screens) {
            NSNumber* screenNumber = [screen.deviceDescription objectForKey:@"NSScreenNumber"];
            NSString* screenID = [NSString stringWithFormat:@"%u", [screenNumber unsignedIntValue]];
            if ([screenID isEqualToString:self.notchTargetScreenID]) {
                return screen;
            }
        }
    }
    NSScreen* primaryScreen = [screens firstObject];
    return primaryScreen != nil ? primaryScreen : [NSScreen mainScreen];
}

- (CGFloat)cameraHousingCentreForScreen:(NSScreen*)screen found:(BOOL*)found {
    *found = NO;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 120000
    if (@available(macOS 12.0, *)) {
        NSRect leftArea = screen.auxiliaryTopLeftArea;
        NSRect rightArea = screen.auxiliaryTopRightArea;
        CGFloat left = NSMaxX(leftArea);
        CGFloat right = NSMinX(rightArea);
        if (screen.safeAreaInsets.top > 0.0 &&
            !NSIsEmptyRect(leftArea) && !NSIsEmptyRect(rightArea) && right > left) {
            *found = YES;
            return left + (right - left) / 2.0;
        }
    }
#endif
    return NSMidX(screen.frame);
}

- (NSRect)notchTargetFrame {
    NSScreen* screen = [self notchTargetScreen];
    if (screen == nil) {
        return self.frame;
    }

    BOOL hasCameraHousing = NO;
    CGFloat centre = [self cameraHousingCentreForScreen:screen found:&hasCameraHousing];
    NSRect referenceFrame = hasCameraHousing ? screen.frame : screen.visibleFrame;
    NSRect frame = self.frame;
    frame.origin.x = round(centre - NSWidth(frame) / 2.0);
    frame.origin.y = NSMaxY(referenceFrame) - NSHeight(frame);
    return frame;
}

- (NSRect)notchHiddenFrameForTarget:(NSRect)target {
    NSRect hidden = target;
    hidden.origin.y = NSMaxY(target);
    return hidden;
}

- (void)configureWebView:(WKWebView*)webView
            contentWidth:(CGFloat)contentWidth
           contentHeight:(CGFloat)contentHeight
          targetScreenID:(NSString*)targetScreenID {
    // The high-level API owns camera-housing placement. WebviewPanel is
    // movable by its background by default, so disable both native drag paths
    // before exposing the notch window.
    self.movable = NO;
    self.movableByWindowBackground = NO;
    self.notchTargetScreenID = targetScreenID;

    NSView* root = self.contentView;
    WailsNotchBackgroundView* background = [[WailsNotchBackgroundView alloc] initWithFrame:root.bounds];
    background.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
    [root addSubview:background positioned:NSWindowBelow relativeTo:webView];
    [background release];

    // Let web content reach the display's top edge. The camera housing masks
    // the centre naturally, while controls can occupy the usable areas on
    // either side instead of sitting below a native-only top inset.
    webView.frame = NSMakeRect(NotchHorizontalInset, NotchBottomInset, contentWidth, contentHeight);
    webView.autoresizingMask = NSViewNotSizable;
    // WKWebView still has no public background toggle. Wails already uses this
    // WebKit key for transparent and translucent windows.
    [webView setValue:@NO forKey:@"drawsBackground"];

    NSTrackingAreaOptions trackingOptions = NSTrackingMouseEnteredAndExited |
        NSTrackingActiveAlways | NSTrackingInVisibleRect;
    NSTrackingArea* trackingArea = [[NSTrackingArea alloc]
        initWithRect:NSZeroRect
        options:trackingOptions
        owner:self
        userInfo:nil];
    [root addTrackingArea:trackingArea];
    self.notchTrackingArea = trackingArea;
    [trackingArea release];

    [self setFrame:[self notchTargetFrame] display:NO];
}

- (void)performWindowDragWithEvent:(NSEvent*)event {
    // CSS drag regions and Wails' explicit drag bridge both end up here. Keep
    // notch windows anchored even if web content requests a native drag.
}

- (void)mouseEntered:(NSEvent*)event {
    if (!self.visible) {
        return;
    }
    // NSWindowStyleMaskNonactivatingPanel keeps the application inactive while
    // the panel and its WKWebView become key for immediate keyboard input.
    // Ordering first also makes hover choose the front window when independent
    // notch windows overlap.
    [self orderFrontRegardless];
    [self makeKeyWindow];
    if (self.webView != nil) {
        [self makeFirstResponder:self.webView];
    }
}

- (void)mouseExited:(NSEvent*)event {
    // Keep the panel key until the user focuses another window. Hover only
    // acquires focus; it does not imply dismissal or a focus-loss policy.
}

- (void)showAnimated:(BOOL)animated duration:(NSTimeInterval)duration {
    if (self.notchRequestedVisible) {
        [self orderFrontRegardless];
        return;
    }

    self.notchRequestedVisible = YES;
    self.notchAnimationGeneration += 1;
    NSUInteger generation = self.notchAnimationGeneration;
    NSRect target = [self notchTargetFrame];
    if (!animated || duration <= 0.0) {
        self.alphaValue = 1.0;
        [self setFrame:target display:NO];
        [self orderFrontRegardless];
        return;
    }

    [self setFrame:[self notchHiddenFrameForTarget:target] display:NO];
    self.alphaValue = 0.96;
    [self orderFrontRegardless];
    [NSAnimationContext runAnimationGroup:^(NSAnimationContext* context) {
        context.duration = duration;
        context.timingFunction = [CAMediaTimingFunction functionWithControlPoints:0.22 :1.0 :0.36 :1.0];
        [[self animator] setFrame:target display:YES];
        [[self animator] setAlphaValue:1.0];
    } completionHandler:^{
        if (self.notchAnimationGeneration != generation || !self.notchRequestedVisible) {
            return;
        }
        [self setFrame:target display:NO];
        self.alphaValue = 1.0;
    }];
}

- (void)hideAnimated:(BOOL)animated duration:(NSTimeInterval)duration {
    if (!self.notchRequestedVisible && !self.visible) {
        return;
    }

    self.notchRequestedVisible = NO;
    self.notchAnimationGeneration += 1;
    NSUInteger generation = self.notchAnimationGeneration;
    NSRect target = [self notchTargetFrame];
    NSRect hidden = [self notchHiddenFrameForTarget:target];
    if (!animated || duration <= 0.0) {
        [self orderOut:nil];
        [self setFrame:target display:NO];
        self.alphaValue = 1.0;
        return;
    }

    [NSAnimationContext runAnimationGroup:^(NSAnimationContext* context) {
        context.duration = duration * (2.0 / 3.0);
        context.timingFunction = [CAMediaTimingFunction functionWithControlPoints:0.4 :0.0 :1.0 :1.0];
        [[self animator] setFrame:hidden display:YES];
        [[self animator] setAlphaValue:0.9];
    } completionHandler:^{
        if (self.notchAnimationGeneration != generation || self.notchRequestedVisible) {
            return;
        }
        [self orderOut:nil];
        [self setFrame:target display:NO];
        self.alphaValue = 1.0;
    }];
}

- (void)close {
    self.notchRequestedVisible = NO;
    self.notchAnimationGeneration += 1;
    [super close];
}

- (void)dealloc {
    if (self.notchTrackingArea != nil) {
        [self.contentView removeTrackingArea:self.notchTrackingArea];
        self.notchTrackingArea = nil;
    }
    self.notchTargetScreenID = nil;
    [super dealloc];
}

@end
