//go:build darwin && !ios

#ifndef WebviewWindowDelegate_h
#define WebviewWindowDelegate_h

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@protocol WailsWebviewWindow <NSObject>
@property (assign) WKWebView* webView;
@property BOOL disableEscapeExitsFullscreen;
@end

@interface WebviewWindow : NSWindow <WailsWebviewWindow>
- (BOOL) canBecomeKeyWindow;
- (BOOL) canBecomeMainWindow;
- (BOOL) acceptsFirstResponder;
- (BOOL) becomeFirstResponder;
- (BOOL) resignFirstResponder;
- (WebviewWindow*) initWithContentRect:(NSRect)contentRect styleMask:(NSUInteger)windowStyle backing:(NSBackingStoreType)bufferingType defer:(BOOL)deferCreation;

@property (assign) WKWebView* webView; // We already retain WKWebView since it's part of the Window.
@property BOOL disableEscapeExitsFullscreen;

@end

@interface WebviewWindowDelegate : NSObject <NSWindowDelegate, WKScriptMessageHandler, WKNavigationDelegate, WKURLSchemeHandler, NSDraggingDestination, WKUIDelegate>

@property unsigned int windowId;
@property (retain) NSEvent* leftMouseEvent;
@property unsigned int invisibleTitleBarHeight;
@property BOOL showToolbarWhenFullscreen;
@property NSWindowStyleMask previousStyleMask; // Used to restore the window style mask when using frameless

- (void)handleLeftMouseUp:(NSWindow *)window;
- (void)handleLeftMouseDown:(NSEvent*)event;
- (void)startDrag:(NSWindow*)window;

@end


NSString* acceleratorStringFromKeyEvent(NSEvent* event);
BOOL dispatchKeyEquivalent(NSEvent* event, NSWindow* window);

void windowSetScreen(void* window, void* screen, int yOffset);

// Liquid Glass support functions
bool isLiquidGlassSupported();
void windowSetLiquidGlass(void* window, int style, int material, double cornerRadius, 
                          int r, int g, int b, int a, 
                          const char* groupID, double groupSpacing);
void windowRemoveVisualEffects(void* window);
void configureWebViewForLiquidGlass(void* window);

// Shared by webview_window_split_darwin.m and
// webview_window_titlebar_accessory_darwin.m; see the .m for details.
WKWebView* buildAuxiliaryWebview(void* nsWindow, double width, double height);

#endif /* WebviewWindowDelegate_h */
