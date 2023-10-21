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

#if __has_include(<UniformTypeIdentifiers/UTType.h>)
#define USE_NEW_FILTERS
#import <UniformTypeIdentifiers/UTType.h>
#endif

#define ON_MAIN_THREAD(str) dispatch_async(dispatch_get_main_queue(), ^{ str; });
#define unicode(input) [NSString stringWithFormat:@"%C", input]

@interface WailsWindow : NSWindow

@property NSSize userMinSize;
@property NSSize userMaxSize;

- (BOOL) canBecomeKeyWindow;
- (void) applyWindowConstraints;
- (void) disableWindowConstraints;
@end

@interface WailsContext : NSObject <WKURLSchemeHandler,WKScriptMessageHandler,WKNavigationDelegate,WKUIDelegate>

@property (retain) WailsWindow* mainWindow;
@property (retain) WKWebView* webview;
@property (nonatomic, assign) id appdelegate;

@property bool hideOnClose;
@property bool shuttingDown;
@property bool startHidden;
@property bool startFullscreen;

@property bool singleInstanceLockEnabled;
@property (retain) NSString* singleInstanceUniqueId;

@property (retain) NSEvent* mouseEvent;

@property bool alwaysOnTop;

@property bool devtoolsEnabled;
@property bool defaultContextMenuEnabled;

@property (retain) WKUserContentController* userContentController;

@property (retain) NSMenu* applicationMenu;

@property (retain) NSImage* aboutImage;
@property (retain) NSString* aboutTitle;
@property (retain) NSString* aboutDescription;

struct Preferences {
  bool *tabFocusesLinks;
  bool *textInteractionEnabled;
  bool *fullscreenEnabled;
};

- (void) CreateWindow:(int)width :(int)height :(bool)frameless :(bool)resizable :(bool)fullscreen :(bool)fullSizeContent :(bool)hideTitleBar :(bool)titlebarAppearsTransparent  :(bool)hideTitle :(bool)useToolbar :(bool)hideToolbarSeparator :(bool)webviewIsTransparent :(bool)hideWindowOnClose :(NSString *)appearance :(bool)windowIsTranslucent :(int)minWidth :(int)minHeight :(int)maxWidth :(int)maxHeight :(bool)fraudulentWebsiteWarningEnabled :(struct Preferences)preferences;
- (void) SetSize:(int)width :(int)height;
- (void) SetPosition:(int)x :(int) y;
- (void) SetMinSize:(int)minWidth :(int)minHeight;
- (void) SetMaxSize:(int)maxWidth :(int)maxHeight;
- (void) SetTitle:(NSString*)title;
- (void) SetAlwaysOnTop:(int)onTop;
- (void) Center;
- (void) Fullscreen;
- (void) UnFullscreen;
- (bool) IsFullScreen;
- (void) Minimise;
- (void) UnMinimise;
- (bool) IsMinimised;
- (void) Maximise;
- (void) ToggleMaximise;
- (void) UnMaximise;
- (bool) IsMaximised;
- (void) SetBackgroundColour:(int)r :(int)g :(int)b :(int)a;
- (void) HideMouse;
- (void) ShowMouse;
- (void) Hide;
- (void) Show;
- (void) HideApplication;
- (void) ShowApplication;
- (void) Quit;

-(void) MessageDialog :(NSString*)dialogType :(NSString*)title :(NSString*)message :(NSString*)button1 :(NSString*)button2 :(NSString*)button3 :(NSString*)button4 :(NSString*)defaultButton :(NSString*)cancelButton :(void*)iconData :(int)iconDataLength;
- (void) OpenFileDialog :(NSString*)title :(NSString*)defaultFilename :(NSString*)defaultDirectory :(bool)allowDirectories :(bool)allowFiles :(bool)canCreateDirectories :(bool)treatPackagesAsDirectories :(bool)resolveAliases :(bool)showHiddenFiles :(bool)allowMultipleSelection :(NSString*)filters;
- (void) SaveFileDialog :(NSString*)title :(NSString*)defaultFilename :(NSString*)defaultDirectory :(bool)canCreateDirectories :(bool)treatPackagesAsDirectories :(bool)showHiddenFiles :(NSString*)filters;

- (void) loadRequest:(NSString*)url;
- (void) ExecJS:(NSString*)script;
- (NSScreen*) getCurrentScreen;

- (void) SetAbout :(NSString*)title :(NSString*)description :(void*)imagedata :(int)datalen;
- (void) dealloc;

@end


#endif /* WailsContext_h */
