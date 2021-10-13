//
//  WailsContext.m
//  test
//
//  Created by Lea Anthony on 10/10/21.
//

#import <Foundation/Foundation.h>
#import "WailsContext.h"
#import "WindowDelegate.h"

@implementation WailsWindow

- (BOOL)canBecomeKeyWindow
{
    return YES;
}

@end

@implementation WailsContext

- (void) SetSize:(int)width :(int)height {

    if (self.shuttingDown) return;

    NSRect frame = [self.mainWindow frame];
    frame.origin.y += frame.size.height - height;
    frame.size.width = width;
    frame.size.height = height;
    [self.mainWindow setFrame:frame display:TRUE animate:FALSE];
}

- (void) SetPosition:(int)x :(int)y {
    
    if (self.shuttingDown) return;
    
    NSScreen* screen = [self getCurrentScreen];
    NSRect windowFrame = [self.mainWindow frame];
    NSRect screenFrame = [screen frame];
    windowFrame.origin.x += screenFrame.origin.x + (float)x;
    windowFrame.origin.y += (screenFrame.origin.y + screenFrame.size.height) - windowFrame.size.height - (float)y;

    [self.mainWindow setFrame:windowFrame display:TRUE animate:FALSE];
}

- (void) SetMinWindowSize:(int)minWidth :(int)minHeight {
    
    if (self.shuttingDown) return;

    NSSize size = { minWidth, minHeight };

    self.minSize = size;
    
    [self.mainWindow setMinSize:size];
    
    [self adjustWindowSize];
}


- (void) SetMaxWindowSize:(int)maxWidth :(int)maxHeight {
    
    if (self.shuttingDown) return;

    NSSize size = { FLT_MAX, FLT_MAX };
    
    size.width = maxWidth > 0 ? maxWidth : FLT_MAX;
    size.height = maxHeight > 0 ? maxHeight : FLT_MAX;

    self.maxSize = size;
    
    [self.mainWindow setMinSize:size];
    
    [self adjustWindowSize];
}


- (void) adjustWindowSize {
    
    if (self.shuttingDown) return;
    
    NSRect currentFrame = [self.mainWindow frame];
    
    if ( currentFrame.size.width > self.maxSize.width ) currentFrame.size.width = self.maxSize.width;
    if ( currentFrame.size.width < self.minSize.width ) currentFrame.size.width = self.minSize.width;
    if ( currentFrame.size.height > self.maxSize.height ) currentFrame.size.height = self.maxSize.height;
    if ( currentFrame.size.height < self.minSize.height ) currentFrame.size.height = self.minSize.height;

    [self.mainWindow setFrame:currentFrame display:TRUE animate:FALSE];
        
}

- (void) dealloc {
    [super dealloc];
    [self.appdelegate release];
    [self.mainWindow release];
}
- (NSScreen*) getCurrentScreen {
    NSScreen* screen = [self.mainWindow screen];
    if( screen == NULL ) {
        screen = [NSScreen mainScreen];
    }
    return screen;
}

- (void) SetTitle:(const char *)title {
    NSString *_title = [NSString stringWithUTF8String:title];
    [self.mainWindow setTitle:_title];
}

- (void) Center {
    [self.mainWindow center];
}

- (void) CreateWindow:(int)width :(int)height :(bool)frameless :(bool)resizable :(bool)fullscreen :(bool)fullSizeContent :(bool)hideTitleBar :(bool)titlebarAppearsTransparent :(bool)hideTitle :(bool)useToolbar :(bool)hideToolbarSeparator :(bool)webviewIsTransparent :(bool)hideWindowOnClose :(const char *)appearance :(bool)windowIsTranslucent {
    
    NSWindowStyleMask styleMask = NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskMiniaturizable;
 
    if (frameless) {
        styleMask = NSWindowStyleMaskBorderless;
    } else {
        if (resizable) {
            styleMask |= NSWindowStyleMaskResizable;
        }
    }
    if (fullscreen) {
        styleMask |= NSWindowStyleMaskFullScreen;
    }

    if( fullSizeContent || frameless || titlebarAppearsTransparent ) {
        styleMask |= NSWindowStyleMaskFullSizeContentView;
    }
    
    self.mainWindow = [[[WailsWindow alloc] initWithContentRect:NSMakeRect(0, 0, width, height)
        styleMask:styleMask backing:NSBackingStoreBuffered defer:NO]
            autorelease];
    
    if (frameless) {
        return;
    }
    
    if (useToolbar) {
        id toolbar = [[NSToolbar alloc] initWithIdentifier:@"wails.toolbar"];
        [toolbar autorelease];
        [toolbar setShowsBaselineSeparator:!hideToolbarSeparator];
        [self.mainWindow setToolbar:toolbar];
    }
    
    [self.mainWindow setTitleVisibility:hideTitle];
    [self.mainWindow setTitlebarAppearsTransparent:titlebarAppearsTransparent];
    [self.mainWindow canBecomeKeyWindow];
    
    if (windowIsTranslucent) {
        id contentView = [self.mainWindow contentView];
        NSVisualEffectView *effectView = [NSVisualEffectView alloc];
        NSRect bounds = [contentView bounds];
        [effectView initWithFrame:bounds];
        [effectView setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
        [effectView setBlendingMode:NSVisualEffectBlendingModeBehindWindow];
        [effectView setState:NSVisualEffectStateActive];
        [contentView addSubview:effectView positioned:NSWindowBelow relativeTo:nil];
    }
    
    if (appearance != nil) {
        NSString *name = [NSString stringWithUTF8String:appearance];
        NSAppearance *nsAppearance = [NSAppearance appearanceNamed:name];
        [self.mainWindow setAppearance:nsAppearance];
    }
    
    // Set up min/max
    NSSize maxSize = { FLT_MAX, FLT_MAX };
    self.maxSize = maxSize;
    NSSize minSize = { 0, 0 };
    self.minSize = minSize;
    [self adjustWindowSize];
    
    WindowDelegate *windowDelegate = [WindowDelegate new];
    windowDelegate.hideOnClose = hideWindowOnClose;
    [self.mainWindow setDelegate:windowDelegate];
    
    // Webview stuff here!
    
}

- (bool) isFullScreen {
    long mask = [self.mainWindow styleMask];
    return (mask & NSWindowStyleMaskFullScreen) == NSWindowStyleMaskFullScreen;
}

// Fullscreen sets the main window to be fullscreen
- (void) Fullscreen {
    if( ! [self isFullScreen] ) {
        [self.mainWindow toggleFullScreen:nil];
    }
}

// UnFullscreen resets the main window after a fullscreen
- (void) UnFullscreen {
    if( [self isFullScreen] ) {
        [self.mainWindow toggleFullScreen:nil];
    }
}

- (void) Minimise {
    [self.mainWindow miniaturize:nil];
}

- (void) UnMinimise {
    [self.mainWindow deminiaturize:nil];
}

@end

