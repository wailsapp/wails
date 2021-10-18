//
//  WailsContext.h
//  test
//
//  Created by Lea Anthony on 10/10/21.
//

#ifndef WailsContext_h
#define WailsContext_h

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

#define ON_MAIN_THREAD(str) dispatch_async(dispatch_get_main_queue(), ^{ str; });

@interface WailsWindow : NSWindow
- (BOOL)canBecomeKeyWindow;
@end

@interface WailsContext : NSObject <WKURLSchemeHandler,WKScriptMessageHandler>

@property (retain) WailsWindow* mainWindow;
@property (retain) WKWebView* webview;
@property (nonatomic, assign) id appdelegate;

@property bool hideOnClose;
@property bool shuttingDown;

@property NSSize maxSize;
@property NSSize minSize;

@property (retain) NSEvent* mouseEvent;

@property bool alwaysOnTop;
@property bool maximised;

@property bool debug;

@property (retain) WKUserContentController* userContentController;

@property (retain) NSMutableDictionary *urlRequests;

- (void) CreateWindow:(int)width :(int)height :(bool)frameless :(bool)resizable :(bool)fullscreen :(bool)fullSizeContent :(bool)hideTitleBar :(bool)titlebarAppearsTransparent  :(bool)hideTitle :(bool)useToolbar :(bool)hideToolbarSeparator :(bool)webviewIsTransparent :(bool)hideWindowOnClose :(const char *)appearance :(bool)windowIsTranslucent;
- (void) SetSize:(int)width :(int)height;
- (void) SetPosition:(int)x :(int) y;
- (void) SetMinSize:(int)minWidth :(int)minHeight;
- (void) SetMaxSize:(int)maxWidth :(int)maxHeight;
- (void) SetTitle:(const char*)title;
- (void) Center;
- (void) Fullscreen;
- (void) UnFullscreen;
- (void) Minimise;
- (void) UnMinimise;
- (void) Maximise;
- (void) UnMaximise;
- (void) SetRGBA:(int)r :(int)g :(int)b :(int)a;
- (void) HideMouse;
- (void) ShowMouse;
- (void) Hide;
- (void) Show;


- (void) loadRequest:(NSString*)url;
- (void) processURLResponse:(NSString *)url :(NSString *)contentType :(NSData*)data;
- (void) ExecJS:(const char*)script;
- (NSScreen*) getCurrentScreen;

@end


#endif /* WailsContext_h */
