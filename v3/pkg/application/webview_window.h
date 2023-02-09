//go:build darwin

#ifndef WebviewWindowDelegate_h
#define WebviewWindowDelegate_h

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface WebviewWindow : NSWindow
- (BOOL) canBecomeKeyWindow;
- (BOOL) canBecomeMainWindow;
- (BOOL) acceptsFirstResponder;
- (BOOL) becomeFirstResponder;
- (BOOL) resignFirstResponder;
- (WebviewWindow*) initWithContentRect:(NSRect)contentRect styleMask:(NSUInteger)windowStyle backing:(NSBackingStoreType)bufferingType defer:(BOOL)deferCreation;

@end

@interface WebviewWindowDelegate : NSObject <NSWindowDelegate, WKScriptMessageHandler, WKNavigationDelegate, WKURLSchemeHandler, NSDraggingDestination>

@property bool hideOnClose;
@property (retain) WKWebView* webView;
@property unsigned int windowId;
@property (retain) NSWindow* window;
@property (retain) NSEvent* leftMouseEvent;
@property unsigned int invisibleTitleBarHeight;
@property NSWindowStyleMask previousStyleMask; // Used to restore the window style mask when using frameless

- (void)handleLeftMouseUp:(NSWindow *)window;
- (void)handleLeftMouseDown:(NSEvent*)event;

@end


#endif /* WebviewWindowDelegate_h */
