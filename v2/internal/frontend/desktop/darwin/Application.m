////
////  Window.m
////  test
////
////  Created by Lea Anthony on 10/10/21.
////
//
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import "WailsContext.h"
#import "Application.h"
#import "AppDelegate.h"
//
//
WailsContext* Create(const char* title, int width, int height, int frameless, int resizable, int fullscreen, int fullSizeContent, int hideTitleBar, int titlebarAppearsTransparent, int hideTitle, int useToolbar, int hideToolbarSeparator, int webviewIsTransparent, int alwaysOnTop, int hideWindowOnClose, const char *appearance, int windowIsTranslucent) {
    
    WailsContext *result = [WailsContext new];
    
    [result CreateWindow:width :height :frameless :resizable :fullscreen :fullSizeContent :hideTitleBar :titlebarAppearsTransparent :hideTitle :useToolbar :hideToolbarSeparator :webviewIsTransparent :hideWindowOnClose :appearance :windowIsTranslucent];
    [result SetTitle:title];
    [result Center];
    
    result.alwaysOnTop = alwaysOnTop;
    result.hideOnClose = hideWindowOnClose;
    
    return result;
}

void SetTitle(WailsContext *ctx, const char *title) {
    ON_MAIN_THREAD(
       [ctx SetTitle:title];
    );
}

void SetSize(WailsContext *ctx, int width, int height) {
    ON_MAIN_THREAD(
       [ctx SetSize:width :height];
    );
}

void SetMinWindowSize(WailsContext *ctx, int width, int height) {
    ON_MAIN_THREAD(
       [ctx SetMinWindowSize:width :height];
    );
}

void SetMaxWindowSize(WailsContext *ctx, int width, int height) {
    ON_MAIN_THREAD(
       [ctx SetMaxWindowSize:width :height];
    );
}

void SetPosition(WailsContext *ctx, int x, int y) {
    ON_MAIN_THREAD(
       [ctx SetSize:x :y];
    );
}

void Center(WailsContext *ctx) {
    ON_MAIN_THREAD(
       [ctx Center];
    );
}

void Fullscreen(WailsContext *ctx) {
    ON_MAIN_THREAD(
       [ctx Fullscreen];
    );
}

void UnFullscreen(WailsContext *ctx) {
    ON_MAIN_THREAD(
       [ctx UnFullscreen];
    );
}

void Minimise(WailsContext *ctx) {
    ON_MAIN_THREAD(
       [ctx Minimise];
    );
}

void UnMinimise(WailsContext *ctx) {
    ON_MAIN_THREAD(
       [ctx UnMinimise];
    );
}

void Run(void *inctx) {
    WailsContext *ctx = (WailsContext*) inctx;
    [NSApplication sharedApplication];
    AppDelegate* delegate = [AppDelegate new];
    [NSApp setDelegate:(id)delegate];
    ctx.appdelegate = delegate;
    delegate.mainWindow = ctx.mainWindow;
    delegate.alwaysOnTop = ctx.alwaysOnTop;
    [NSApp run];
    [ctx release];
    NSLog(@"Here");
}
//
