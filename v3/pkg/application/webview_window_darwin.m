//go:build darwin
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import "webview_window_darwin.h"
#import "../events/events_darwin.h"
extern void processMessage(unsigned int, const char*);
extern void processURLRequest(unsigned int, void *);
extern void processWindowKeyDownEvent(unsigned int, const char*);
extern bool hasListeners(unsigned int);
@implementation WebviewWindow
- (WebviewWindow*) initWithContentRect:(NSRect)contentRect styleMask:(NSUInteger)windowStyle backing:(NSBackingStoreType)bufferingType defer:(BOOL)deferCreation;
{
    self = [super initWithContentRect:contentRect styleMask:windowStyle backing:bufferingType defer:deferCreation];
    [self setAlphaValue:1.0];
    [self setBackgroundColor:[NSColor clearColor]];
    [self setOpaque:NO];
    [self setMovableByWindowBackground:YES];
    return self;
}
- (void)keyDown:(NSEvent *)event {
    NSUInteger modifierFlags = event.modifierFlags;
    // Create an array to hold the modifier strings
    NSMutableArray *modifierStrings = [NSMutableArray array];
    // Check for modifier flags and add corresponding strings to the array
    if (modifierFlags & NSEventModifierFlagShift) {
        [modifierStrings addObject:@"shift"];
    }
    if (modifierFlags & NSEventModifierFlagControl) {
        [modifierStrings addObject:@"ctrl"];
    }
    if (modifierFlags & NSEventModifierFlagOption) {
        [modifierStrings addObject:@"option"];
    }
    if (modifierFlags & NSEventModifierFlagCommand) {
        [modifierStrings addObject:@"cmd"];
    }
    NSString *keyString = [self keyStringFromEvent:event];
    if (keyString.length > 0) {
        [modifierStrings addObject:keyString];
    }
    // Combine the modifier strings with the key character
    NSString *keyEventString = [modifierStrings componentsJoinedByString:@"+"];
    const char* utf8String = [keyEventString UTF8String];
    WebviewWindowDelegate *delegate = (WebviewWindowDelegate*)self.delegate;
    processWindowKeyDownEvent(delegate.windowId, utf8String);
}
- (NSString *)keyStringFromEvent:(NSEvent *)event {
    // Get the pressed key
    // Check for special keys like escape and tab
    NSString *characters = [event characters];
    if (characters.length == 0) {
        return @"";
    }
    if ([characters isEqualToString:@"\r"]) {
        return @"enter";
    }
    if ([characters isEqualToString:@"\b"]) {
        return @"backspace";
    }
    if ([characters isEqualToString:@"\e"]) {
        return @"escape";
    }
    // page down
    if ([characters isEqualToString:@"\x0B"]) {
        return @"page down";
    }
    // page up
    if ([characters isEqualToString:@"\x0E"]) {
        return @"page up";
    }
    // home
    if ([characters isEqualToString:@"\x01"]) {
        return @"home";
    }
    // end
    if ([characters isEqualToString:@"\x04"]) {
        return @"end";
    }
    // clear
    if ([characters isEqualToString:@"\x0C"]) {
        return @"clear";
    }
    switch ([event keyCode]) {
        // Function keys
        case 122: return @"f1";
        case 120: return @"f2";
        case 99: return @"f3";
        case 118: return @"f4";
        case 96: return @"f5";
        case 97: return @"f6";
        case 98: return @"f7";
        case 100: return @"f8";
        case 101: return @"f9";
        case 109: return @"f10";
        case 103: return @"f11";
        case 111: return @"f12";
        case 105: return @"f13";
        case 107: return @"f14";
        case 113: return @"f15";
        case 106: return @"f16";
        case 64: return @"f17";
        case 79: return @"f18";
        case 80: return @"f19";
        case 90: return @"f20";
        // Letter keys
        case 0: return @"a";
        case 11: return @"b";
        case 8: return @"c";
        case 2: return @"d";
        case 14: return @"e";
        case 3: return @"f";
        case 5: return @"g";
        case 4: return @"h";
        case 34: return @"i";
        case 38: return @"j";
        case 40: return @"k";
        case 37: return @"l";
        case 46: return @"m";
        case 45: return @"n";
        case 31: return @"o";
        case 35: return @"p";
        case 12: return @"q";
        case 15: return @"r";
        case 1: return @"s";
        case 17: return @"t";
        case 32: return @"u";
        case 9: return @"v";
        case 13: return @"w";
        case 7: return @"x";
        case 16: return @"y";
        case 6: return @"z";
        // Number keys
        case 29: return @"0";
        case 18: return @"1";
        case 19: return @"2";
        case 20: return @"3";
        case 21: return @"4";
        case 23: return @"5";
        case 22: return @"6";
        case 26: return @"7";
        case 28: return @"8";
        case 25: return @"9";
        // Other special keys
        case 51: return @"delete";
        case 117: return @"forward delete";
        case 123: return @"left";
        case 124: return @"right";
        case 126: return @"up";
        case 125: return @"down";
        case 48: return @"tab";
        case 53: return @"escape";
        case 49: return @"space";
        // Punctuation and other keys (for a standard US layout)
        case 33: return @"[";
        case 30: return @"]";
        case 43: return @",";
        case 27: return @"-";
        case 39: return @"'";
        case 44: return @"/";
        case 47: return @".";
        case 41: return @";";
        case 24: return @"=";
        case 50: return @"`";
        case 42: return @"\\";

        default: return @"";
    }
}
- (BOOL)canBecomeKeyWindow {
    return YES;
}
- (BOOL) canBecomeMainWindow {
    return YES;
}
- (BOOL) acceptsFirstResponder {
    return YES;
}
- (BOOL) becomeFirstResponder {
    return YES;
}
- (BOOL) resignFirstResponder {
    return YES;
}
- (void) setDelegate:(id<NSWindowDelegate>) delegate {
    [delegate retain];
    [super setDelegate: delegate];
}
- (void) dealloc {
    // Remove the script handler, otherwise WebviewWindowDelegate won't get deallocated
    // See: https://stackoverflow.com/questions/26383031/wkwebview-causes-my-view-controller-to-leak
    [self.webView.configuration.userContentController removeScriptMessageHandlerForName:@"external"];
    if (self.delegate) {
        [self.delegate release];
    }
    [super dealloc];
}
- (void)windowDidZoom:(NSNotification *)notification {
    NSWindow *window = notification.object;
    WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)[window delegate];
    if ([window isZoomed]) {
        if (hasListeners(EventWindowMaximise)) {
            processWindowEvent(delegate.windowId, EventWindowMaximise);
        }
    } else {
        if (hasListeners(EventWindowUnMaximise)) {
            processWindowEvent(delegate.windowId, EventWindowUnMaximise);
        }
    }
}
- (void)performZoomIn:(id)sender {
    [super zoom:sender];
    if (hasListeners(EventWindowZoomIn)) {
        WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)[sender delegate];
        processWindowEvent(delegate.windowId, EventWindowZoomIn);
    }
}
- (void)performZoomOut:(id)sender {
    [super zoom:sender];
    if (hasListeners(EventWindowZoomOut)) {
        WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)[sender delegate];
        processWindowEvent(delegate.windowId, EventWindowZoomOut);
    }
}
- (void)performZoomReset:(id)sender {
    [self setFrame:[self frameRectForContentRect:[[self screen] visibleFrame]] display:YES];
    if (hasListeners(EventWindowZoomReset)) {
        WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)[sender delegate];
        processWindowEvent(delegate.windowId, EventWindowZoomReset);
    }
}
@end
@implementation WebviewWindowDelegate
- (BOOL)windowShouldClose:(NSWindow *)sender {
    WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)[sender delegate];
    processWindowEvent(delegate.windowId, EventWindowShouldClose);
    return false;
}
- (void) dealloc {
    // Makes sure to remove the retained properties so the reference counter of the retains are decreased
    self.leftMouseEvent = nil;
    [super dealloc];
}
- (void) startDrag:(WebviewWindow*)window {
    [window performWindowDragWithEvent:self.leftMouseEvent];
}
// Handle script messages from the external bridge
- (void)userContentController:(nonnull WKUserContentController *)userContentController didReceiveScriptMessage:(nonnull WKScriptMessage *)message {
    NSString *m = message.body;
    const char *_m = [m UTF8String];
    processMessage(self.windowId, _m);
}
- (void)handleLeftMouseDown:(NSEvent *)event {
    self.leftMouseEvent = event;
    NSWindow *window = [event window];
    WebviewWindowDelegate* delegate = (WebviewWindowDelegate*)[window delegate];
    if( self.invisibleTitleBarHeight > 0 ) {
        NSPoint location = [event locationInWindow];
        NSRect frame = [window frame];
        if( location.y > frame.size.height - self.invisibleTitleBarHeight ) {
            [window performWindowDragWithEvent:event];
            return;
        }
    }
}
- (void)handleLeftMouseUp:(NSWindow *)window {
    self.leftMouseEvent = nil;
}
- (void)webView:(nonnull WKWebView *)webView startURLSchemeTask:(nonnull id<WKURLSchemeTask>)urlSchemeTask {
    processURLRequest(self.windowId, urlSchemeTask);
}
- (void)webView:(nonnull WKWebView *)webView stopURLSchemeTask:(nonnull id<WKURLSchemeTask>)urlSchemeTask {
    NSInputStream *stream = urlSchemeTask.request.HTTPBodyStream;
    if (stream) {
        NSStreamStatus status = stream.streamStatus;
        if (status != NSStreamStatusClosed && status != NSStreamStatusNotOpen) {
            [stream close];
        }
    }
}
- (NSApplicationPresentationOptions)window:(NSWindow *)window willUseFullScreenPresentationOptions:(NSApplicationPresentationOptions)proposedOptions {
    if (self.showToolbarWhenFullscreen) {
        return proposedOptions;
    } else {
        return proposedOptions | NSApplicationPresentationAutoHideToolbar;
    }
}
- (void)windowDidChangeOcclusionState:(NSNotification *)notification {
    NSWindow *window = notification.object;
    BOOL isVisible = ([window occlusionState] & NSWindowOcclusionStateVisible) != 0;
    if (hasListeners(isVisible ? EventWindowShow : EventWindowHide)) {
        processWindowEvent(self.windowId, isVisible ? EventWindowShow : EventWindowHide);
    }
}
- (void)windowDidBecomeKey:(NSNotification *)notification {
    if( hasListeners(EventWindowDidBecomeKey) ) {
        processWindowEvent(self.windowId, EventWindowDidBecomeKey);
    }
}

- (void)windowDidBecomeMain:(NSNotification *)notification {
    if( hasListeners(EventWindowDidBecomeMain) ) {
        processWindowEvent(self.windowId, EventWindowDidBecomeMain);
    }
}

- (void)windowDidBeginSheet:(NSNotification *)notification {
    if( hasListeners(EventWindowDidBeginSheet) ) {
        processWindowEvent(self.windowId, EventWindowDidBeginSheet);
    }
}

- (void)windowDidChangeAlpha:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeAlpha) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeAlpha);
    }
}

- (void)windowDidChangeBackingLocation:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeBackingLocation) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeBackingLocation);
    }
}

- (void)windowDidChangeBackingProperties:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeBackingProperties) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeBackingProperties);
    }
}

- (void)windowDidChangeCollectionBehavior:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeCollectionBehavior) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeCollectionBehavior);
    }
}

- (void)windowDidChangeEffectiveAppearance:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeEffectiveAppearance) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeEffectiveAppearance);
    }
}

- (void)windowDidChangeOrderingMode:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeOrderingMode) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeOrderingMode);
    }
}

- (void)windowDidChangeScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeScreen) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeScreen);
    }
}

- (void)windowDidChangeScreenParameters:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeScreenParameters) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeScreenParameters);
    }
}

- (void)windowDidChangeScreenProfile:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeScreenProfile) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeScreenProfile);
    }
}

- (void)windowDidChangeScreenSpace:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeScreenSpace) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeScreenSpace);
    }
}

- (void)windowDidChangeScreenSpaceProperties:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeScreenSpaceProperties) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeScreenSpaceProperties);
    }
}

- (void)windowDidChangeSharingType:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeSharingType) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeSharingType);
    }
}

- (void)windowDidChangeSpace:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeSpace) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeSpace);
    }
}

- (void)windowDidChangeSpaceOrderingMode:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeSpaceOrderingMode) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeSpaceOrderingMode);
    }
}

- (void)windowDidChangeTitle:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeTitle) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeTitle);
    }
}

- (void)windowDidChangeToolbar:(NSNotification *)notification {
    if( hasListeners(EventWindowDidChangeToolbar) ) {
        processWindowEvent(self.windowId, EventWindowDidChangeToolbar);
    }
}

- (void)windowDidDeminiaturize:(NSNotification *)notification {
    if( hasListeners(EventWindowDidDeminiaturize) ) {
        processWindowEvent(self.windowId, EventWindowDidDeminiaturize);
    }
}

- (void)windowDidEndSheet:(NSNotification *)notification {
    if( hasListeners(EventWindowDidEndSheet) ) {
        processWindowEvent(self.windowId, EventWindowDidEndSheet);
    }
}

- (void)windowDidEnterFullScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowDidEnterFullScreen) ) {
        processWindowEvent(self.windowId, EventWindowDidEnterFullScreen);
    }
}

- (void)windowMaximise:(NSNotification *)notification {
    if( hasListeners(EventWindowMaximise) ) {
        processWindowEvent(self.windowId, EventWindowMaximise);
    }
}

- (void)windowUnMaximise:(NSNotification *)notification {
    if( hasListeners(EventWindowUnMaximise) ) {
        processWindowEvent(self.windowId, EventWindowUnMaximise);
    }
}

- (void)windowDidEnterVersionBrowser:(NSNotification *)notification {
    if( hasListeners(EventWindowDidEnterVersionBrowser) ) {
        processWindowEvent(self.windowId, EventWindowDidEnterVersionBrowser);
    }
}

- (void)windowDidExitFullScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowDidExitFullScreen) ) {
        processWindowEvent(self.windowId, EventWindowDidExitFullScreen);
    }
}

- (void)windowDidExitVersionBrowser:(NSNotification *)notification {
    if( hasListeners(EventWindowDidExitVersionBrowser) ) {
        processWindowEvent(self.windowId, EventWindowDidExitVersionBrowser);
    }
}

- (void)windowDidExpose:(NSNotification *)notification {
    if( hasListeners(EventWindowDidExpose) ) {
        processWindowEvent(self.windowId, EventWindowDidExpose);
    }
}

- (void)windowDidFocus:(NSNotification *)notification {
    if( hasListeners(EventWindowDidFocus) ) {
        processWindowEvent(self.windowId, EventWindowDidFocus);
    }
}

- (void)windowDidMiniaturize:(NSNotification *)notification {
    if( hasListeners(EventWindowDidMiniaturize) ) {
        processWindowEvent(self.windowId, EventWindowDidMiniaturize);
    }
}

- (void)windowDidMove:(NSNotification *)notification {
    if( hasListeners(EventWindowDidMove) ) {
        processWindowEvent(self.windowId, EventWindowDidMove);
    }
}

- (void)windowDidOrderOffScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowDidOrderOffScreen) ) {
        processWindowEvent(self.windowId, EventWindowDidOrderOffScreen);
    }
}

- (void)windowDidOrderOnScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowDidOrderOnScreen) ) {
        processWindowEvent(self.windowId, EventWindowDidOrderOnScreen);
    }
}

- (void)windowDidResignKey:(NSNotification *)notification {
    if( hasListeners(EventWindowDidResignKey) ) {
        processWindowEvent(self.windowId, EventWindowDidResignKey);
    }
}

- (void)windowDidResignMain:(NSNotification *)notification {
    if( hasListeners(EventWindowDidResignMain) ) {
        processWindowEvent(self.windowId, EventWindowDidResignMain);
    }
}

- (void)windowDidResize:(NSNotification *)notification {
    if( hasListeners(EventWindowDidResize) ) {
        processWindowEvent(self.windowId, EventWindowDidResize);
    }
}

- (void)windowDidUpdate:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdate) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdate);
    }
}

- (void)windowDidUpdateAlpha:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdateAlpha) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdateAlpha);
    }
}

- (void)windowDidUpdateCollectionBehavior:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdateCollectionBehavior) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdateCollectionBehavior);
    }
}

- (void)windowDidUpdateCollectionProperties:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdateCollectionProperties) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdateCollectionProperties);
    }
}

- (void)windowDidUpdateShadow:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdateShadow) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdateShadow);
    }
}

- (void)windowDidUpdateTitle:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdateTitle) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdateTitle);
    }
}

- (void)windowDidUpdateToolbar:(NSNotification *)notification {
    if( hasListeners(EventWindowDidUpdateToolbar) ) {
        processWindowEvent(self.windowId, EventWindowDidUpdateToolbar);
    }
}

- (void)windowWillBecomeKey:(NSNotification *)notification {
    if( hasListeners(EventWindowWillBecomeKey) ) {
        processWindowEvent(self.windowId, EventWindowWillBecomeKey);
    }
}

- (void)windowWillBecomeMain:(NSNotification *)notification {
    if( hasListeners(EventWindowWillBecomeMain) ) {
        processWindowEvent(self.windowId, EventWindowWillBecomeMain);
    }
}

- (void)windowWillBeginSheet:(NSNotification *)notification {
    if( hasListeners(EventWindowWillBeginSheet) ) {
        processWindowEvent(self.windowId, EventWindowWillBeginSheet);
    }
}

- (void)windowWillChangeOrderingMode:(NSNotification *)notification {
    if( hasListeners(EventWindowWillChangeOrderingMode) ) {
        processWindowEvent(self.windowId, EventWindowWillChangeOrderingMode);
    }
}

- (void)windowWillClose:(NSNotification *)notification {
    if( hasListeners(EventWindowWillClose) ) {
        processWindowEvent(self.windowId, EventWindowWillClose);
    }
}

- (void)windowWillDeminiaturize:(NSNotification *)notification {
    if( hasListeners(EventWindowWillDeminiaturize) ) {
        processWindowEvent(self.windowId, EventWindowWillDeminiaturize);
    }
}

- (void)windowWillEnterFullScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowWillEnterFullScreen) ) {
        processWindowEvent(self.windowId, EventWindowWillEnterFullScreen);
    }
}

- (void)windowWillEnterVersionBrowser:(NSNotification *)notification {
    if( hasListeners(EventWindowWillEnterVersionBrowser) ) {
        processWindowEvent(self.windowId, EventWindowWillEnterVersionBrowser);
    }
}

- (void)windowWillExitFullScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowWillExitFullScreen) ) {
        processWindowEvent(self.windowId, EventWindowWillExitFullScreen);
    }
}

- (void)windowWillExitVersionBrowser:(NSNotification *)notification {
    if( hasListeners(EventWindowWillExitVersionBrowser) ) {
        processWindowEvent(self.windowId, EventWindowWillExitVersionBrowser);
    }
}

- (void)windowWillFocus:(NSNotification *)notification {
    if( hasListeners(EventWindowWillFocus) ) {
        processWindowEvent(self.windowId, EventWindowWillFocus);
    }
}

- (void)windowWillMiniaturize:(NSNotification *)notification {
    if( hasListeners(EventWindowWillMiniaturize) ) {
        processWindowEvent(self.windowId, EventWindowWillMiniaturize);
    }
}

- (void)windowWillMove:(NSNotification *)notification {
    if( hasListeners(EventWindowWillMove) ) {
        processWindowEvent(self.windowId, EventWindowWillMove);
    }
}

- (void)windowWillOrderOffScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowWillOrderOffScreen) ) {
        processWindowEvent(self.windowId, EventWindowWillOrderOffScreen);
    }
}

- (void)windowWillOrderOnScreen:(NSNotification *)notification {
    if( hasListeners(EventWindowWillOrderOnScreen) ) {
        processWindowEvent(self.windowId, EventWindowWillOrderOnScreen);
    }
}

- (void)windowWillResignMain:(NSNotification *)notification {
    if( hasListeners(EventWindowWillResignMain) ) {
        processWindowEvent(self.windowId, EventWindowWillResignMain);
    }
}

- (void)windowWillResize:(NSNotification *)notification {
    if( hasListeners(EventWindowWillResize) ) {
        processWindowEvent(self.windowId, EventWindowWillResize);
    }
}

- (void)windowWillUnfocus:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUnfocus) ) {
        processWindowEvent(self.windowId, EventWindowWillUnfocus);
    }
}

- (void)windowWillUpdate:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdate) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdate);
    }
}

- (void)windowWillUpdateAlpha:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateAlpha) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateAlpha);
    }
}

- (void)windowWillUpdateCollectionBehavior:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateCollectionBehavior) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateCollectionBehavior);
    }
}

- (void)windowWillUpdateCollectionProperties:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateCollectionProperties) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateCollectionProperties);
    }
}

- (void)windowWillUpdateShadow:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateShadow) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateShadow);
    }
}

- (void)windowWillUpdateTitle:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateTitle) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateTitle);
    }
}

- (void)windowWillUpdateToolbar:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateToolbar) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateToolbar);
    }
}

- (void)windowWillUpdateVisibility:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUpdateVisibility) ) {
        processWindowEvent(self.windowId, EventWindowWillUpdateVisibility);
    }
}

- (void)windowWillUseStandardFrame:(NSNotification *)notification {
    if( hasListeners(EventWindowWillUseStandardFrame) ) {
        processWindowEvent(self.windowId, EventWindowWillUseStandardFrame);
    }
}

- (void)windowFileDraggingEntered:(NSNotification *)notification {
    if( hasListeners(EventWindowFileDraggingEntered) ) {
        processWindowEvent(self.windowId, EventWindowFileDraggingEntered);
    }
}

- (void)windowFileDraggingPerformed:(NSNotification *)notification {
    if( hasListeners(EventWindowFileDraggingPerformed) ) {
        processWindowEvent(self.windowId, EventWindowFileDraggingPerformed);
    }
}

- (void)windowFileDraggingExited:(NSNotification *)notification {
    if( hasListeners(EventWindowFileDraggingExited) ) {
        processWindowEvent(self.windowId, EventWindowFileDraggingExited);
    }
}

- (void)windowShow:(NSNotification *)notification {
    if( hasListeners(EventWindowShow) ) {
        processWindowEvent(self.windowId, EventWindowShow);
    }
}

- (void)windowHide:(NSNotification *)notification {
    if( hasListeners(EventWindowHide) ) {
        processWindowEvent(self.windowId, EventWindowHide);
    }
}

- (void)webView:(nonnull WKWebView *)webview didStartProvisionalNavigation:(WKNavigation *)navigation {
    if( hasListeners(EventWebViewDidStartProvisionalNavigation) ) {
        processWindowEvent(self.windowId, EventWebViewDidStartProvisionalNavigation);
    }
}

- (void)webView:(nonnull WKWebView *)webview didReceiveServerRedirectForProvisionalNavigation:(WKNavigation *)navigation {
    if( hasListeners(EventWebViewDidReceiveServerRedirectForProvisionalNavigation) ) {
        processWindowEvent(self.windowId, EventWebViewDidReceiveServerRedirectForProvisionalNavigation);
    }
}

- (void)webView:(nonnull WKWebView *)webview didFinishNavigation:(WKNavigation *)navigation {
    if( hasListeners(EventWebViewDidFinishNavigation) ) {
        processWindowEvent(self.windowId, EventWebViewDidFinishNavigation);
    }
}

- (void)webView:(nonnull WKWebView *)webview didCommitNavigation:(WKNavigation *)navigation {
    if( hasListeners(EventWebViewDidCommitNavigation) ) {
        processWindowEvent(self.windowId, EventWebViewDidCommitNavigation);
    }
}

// GENERATED EVENTS END
@end
void windowSetScreen(void* window, void* screen, int yOffset) {
    WebviewWindow* nsWindow = (WebviewWindow*)window;
    NSScreen* nsScreen = (NSScreen*)screen;
    
    // Get current frame
    NSRect frame = [nsWindow frame];
    
    // Convert frame to screen coordinates
    NSRect screenFrame = [nsScreen frame];
    NSRect currentScreenFrame = [[nsWindow screen] frame];
    
    // Calculate the menubar height for the target screen
    NSRect visibleFrame = [nsScreen visibleFrame];
    CGFloat menubarHeight = screenFrame.size.height - visibleFrame.size.height;
    
    // Calculate the distance from the top of the current screen
    CGFloat topOffset = currentScreenFrame.origin.y + currentScreenFrame.size.height - frame.origin.y;
    
    // Position relative to new screen's top, accounting for menubar
    frame.origin.x = screenFrame.origin.x + (frame.origin.x - currentScreenFrame.origin.x);
    frame.origin.y = screenFrame.origin.y + screenFrame.size.height - topOffset - menubarHeight - yOffset;
    
    // Set the frame which moves the window to the new screen
    [nsWindow setFrame:frame display:YES];
}
