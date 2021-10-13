//
//  WailsContext.h
//  test
//
//  Created by Lea Anthony on 10/10/21.
//

#ifndef WailsContext_h
#define WailsContext_h

#import <Cocoa/Cocoa.h>

@interface WailsWindow : NSWindow
- (BOOL)canBecomeKeyWindow;
@end

@interface WailsContext : NSObject

@property (retain) WailsWindow* mainWindow;
@property (nonatomic, assign) id appdelegate;

@property bool hideOnClose;
@property bool shuttingDown;

@property NSSize maxSize;
@property NSSize minSize;

@property bool alwaysOnTop;

- (void) CreateWindow:(int)width :(int)height :(bool)frameless :(bool)resizable :(bool)fullscreen :(bool)fullSizeContent :(bool)hideTitleBar :(bool)titlebarAppearsTransparent  :(bool)hideTitle :(bool)useToolbar :(bool)hideToolbarSeparator :(bool)webviewIsTransparent :(bool)hideWindowOnClose :(const char *)appearance :(bool)windowIsTranslucent;
- (void) SetSize:(int)width :(int)height;
- (void) SetPosition:(int)x :(int) y;
- (void) SetMinWindowSize:(int)minWidth :(int)minHeight;
- (void) SetMaxWindowSize:(int)maxWidth :(int)maxHeight;
- (void) SetTitle:(const char*)title;
- (void) Center;
- (void) Fullscreen;
- (void) UnFullscreen;
- (void) Minimise;
- (void) UnMinimise;

@end


#endif /* WailsContext_h */
