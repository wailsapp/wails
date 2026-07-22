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

// Observes primary-button mouse input through the gesture-recognizer system
// and forwards it to the window's WebviewWindowDelegate. macOS 27 deprecates
// NSEvent monitors and NSResponder mouse overrides for this (TN3212), and
// gesture recognizers are the only input path that receives Sidecar/touch
// synthesised events. The observer never claims the gesture, so event
// delivery to the webview is unaffected.
@interface WailsWindowMouseGestureObserver : NSGestureRecognizer
@end

// Content view for frameless windows. On macOS 27 the system resolves corner
// radii through the view's cornerConfiguration; this view returns a
// container-concentric configuration so its corners track the window's system
// radius, and applies the resolved radii to its backing layer. On macOS <= 26
// none of the corner-configuration machinery is invoked and windowNew's
// hardcoded radius stands.
@interface WailsFramelessContentView : NSView
@end

void windowSetScreen(void* window, void* screen, int yOffset);

// Liquid Glass support functions
bool isLiquidGlassSupported();
void windowSetLiquidGlass(void* window, int style, int material, double cornerRadius, 
                          int r, int g, int b, int a, 
                          const char* groupID, double groupSpacing);
void windowRemoveVisualEffects(void* window);
void configureWebViewForLiquidGlass(void* window);

#endif /* WebviewWindowDelegate_h */
