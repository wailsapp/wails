//go:build darwin && !ios

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
- (void)startDrag:(WebviewWindow*)window;

@end

void windowSetScreen(void* window, void* screen, int yOffset);

// Liquid Glass support functions
bool isLiquidGlassSupported();
void windowSetLiquidGlass(void* window, int style, int material, double cornerRadius, 
                          int r, int g, int b, int a, 
                          const char* groupID, double groupSpacing);
void windowRemoveVisualEffects(void* window);
void configureWebViewForLiquidGlass(void* window);

// Implemented by the privatemacapis build-tag variants. The default
// implementation uses public WebKit APIs only; the opt-in implementation
// restores WKWebView transparency and undocumented Liquid Glass grouping.
void wailsSetWebViewTransparent(void* window);
void wailsSetWebViewBackgroundColour(void* window, int r, int g, int b, int alpha);
void wailsConfigurePrivateLiquidGlass(void* glassView, int style, const char* groupID, double groupSpacing);

#endif /* WebviewWindowDelegate_h */
