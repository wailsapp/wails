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
#define unicode(input) [NSString stringWithFormat:@"%C", input]
#define STREQ(a,b) strcmp(a, b) == 0


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

@property bool debug;

@property (retain) WKUserContentController* userContentController;
@property (retain) NSMutableDictionary *urlRequests;

@property (retain) NSMenu* applicationMenu;

@property (retain) NSImage* aboutImage;
@property const char* aboutTitle;
@property const char* aboutDescription;
@property (retain) NSString* appName;
@property (retain) NSString* title;

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
- (void) Quit;

- (void) MessageDialog :(const char*)dialogType :(const char*)title :(const char*)message :(const char*)button1 :(const char*)button2 :(const char*)button3 :(const char*)button4 :(const char*)defaultButton :(const char*)cancelButton;
- (void) OpenFileDialog :(const char*)title :(const char*)defaultFilename :(const char*)defaultDirectory :(bool)allowDirectories :(bool)allowFiles :(bool)canCreateDirectories :(bool)treatPackagesAsDirectories :(bool)resolveAliases :(bool)showHiddenFiles :(bool)allowMultipleSelection :(const char*)filters;
- (void) SaveFileDialog :(const char*)title :(const char*)defaultFilename :(const char*)defaultDirectory :(bool)canCreateDirectories :(bool)treatPackagesAsDirectories :(bool)showHiddenFiles :(const char*)filters;

- (void) loadRequest:(NSString*)url;
- (void) processURLResponse:(NSString *)url :(NSString *)contentType :(NSData*)data;
- (void) ExecJS:(NSString*)script;
- (NSScreen*) getCurrentScreen;

- (void) SetAbout :(const char*)title :(const char*)description :(void*)imagedata :(int)datalen;

@end


#endif /* WailsContext_h */
