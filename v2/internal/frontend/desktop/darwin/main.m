//go:build ignore
//  main.m
//  test
//
//  Created by Lea Anthony on 10/10/21.
//

// ****** This file is used for testing purposes only ******

#import <Foundation/Foundation.h>
#import "Application.h"

void processMessage(const char*t) {
    NSLog(@"processMessage called");

}

void processMessageDialogResponse(int t) {
    NSLog(@"processMessage called");
}

void processOpenFileDialogResponse(const char *t) {
    NSLog(@"processMessage called %s", t);
}

void processURLRequest(void *ctx, const char* url) {
    NSLog(@"processURLRequest called");
    const char myByteArray[] = { 0x3c,0x68,0x31,0x3e,0x48,0x65,0x6c,0x6c,0x6f,0x20,0x57,0x6f,0x72,0x6c,0x64,0x21,0x3c,0x2f,0x68,0x31,0x3e };
    ProcessURLResponse(ctx, url, "text/html", myByteArray, 21);
}

int main(int argc, const char * argv[]) {
    // insert code here...
    int frameless = 0;
    int resizable = 1;
    int fullscreen = 0;
    int fullSizeContent = 1;
    int hideTitleBar = 0;
    int titlebarAppearsTransparent = 1;
    int hideTitle = 0;
    int useToolbar = 1;
    int hideToolbarSeparator = 1;
    int webviewIsTransparent = 0;
    int alwaysOnTop = 1;
    int hideWindowOnClose = 0;
    const char* appearance = "NSAppearanceNameDarkAqua";
    int windowIsTranslucent = 1;
    int debug = 1;
    WailsContext *result = Create("OI OI!",400,400, frameless,  resizable, fullscreen, fullSizeContent, hideTitleBar, titlebarAppearsTransparent, hideTitle, useToolbar, hideToolbarSeparator, webviewIsTransparent, alwaysOnTop, hideWindowOnClose, appearance, windowIsTranslucent, debug);
    SetRGBA(result, 255, 0, 0, 255);
    
    
    
    Run((void*)CFBridgingRetain(result));
    return 0;
}
