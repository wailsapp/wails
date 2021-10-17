//
//  Application.m
//
//  Created by Lea Anthony on 10/10/21.
//

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import "WailsContext.h"
#import "Application.h"
#import "AppDelegate.h"

WailsContext* Create(const char* title, int width, int height, int frameless, int resizable, int fullscreen, int fullSizeContent, int hideTitleBar, int titlebarAppearsTransparent, int hideTitle, int useToolbar, int hideToolbarSeparator, int webviewIsTransparent, int alwaysOnTop, int hideWindowOnClose, const char *appearance, int windowIsTranslucent, int debug) {
    
    WailsContext *result = [WailsContext new];

    result.debug = debug;
    
    [result CreateWindow:width :height :frameless :resizable :fullscreen :fullSizeContent :hideTitleBar :titlebarAppearsTransparent :hideTitle :useToolbar :hideToolbarSeparator :webviewIsTransparent :hideWindowOnClose :appearance :windowIsTranslucent];
    [result SetTitle:title];
    [result Center];
    
    result.alwaysOnTop = alwaysOnTop;
    result.hideOnClose = hideWindowOnClose;
    
    return result;
}

void ProcessURLResponse(void *inctx, const char *url, const char *contentType, Byte data[], int datalength) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    NSString *nsurl = [[NSString alloc] initWithUTF8String:url];
    NSString *nsContentType = [[NSString alloc] initWithUTF8String:contentType];
    NSData *nsdata = [NSData dataWithBytes:data length:datalength];
    
    [ctx processURLResponse:nsurl :nsContentType :nsdata];
}

void SetTitle(WailsContext *ctx, const char *title) {
    ON_MAIN_THREAD(
       [ctx SetTitle:title];
    );
}


void SetRGBA(void *inctx, int r, int g, int b, int a) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    ON_MAIN_THREAD(
       [ctx SetRGBA:r :g :b :a];
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

void Quit(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    [NSApp stop:ctx];
}


void Run(void *inctx) {
    WailsContext *ctx = (__bridge WailsContext*) inctx;
    [NSApplication sharedApplication];
    AppDelegate* delegate = [AppDelegate new];
    [NSApp setDelegate:(id)delegate];
    ctx.appdelegate = delegate;
    delegate.mainWindow = ctx.mainWindow;
    delegate.alwaysOnTop = ctx.alwaysOnTop;

    [ctx loadRequest:@"wails://wails/ipc.js"];
    
    [NSApp run];
    [ctx release];
    NSLog(@"Here");
}
