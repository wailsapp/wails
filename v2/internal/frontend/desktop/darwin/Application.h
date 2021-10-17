//
//  Application.h
//  test
//
//  Created by Lea Anthony on 10/10/21.
//

#ifndef Application_h
#define Application_h

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import "WailsContext.h"

#define ON_MAIN_THREAD(str) dispatch_async(dispatch_get_main_queue(), ^{ str; });

WailsContext* Create(const char* title, int width, int height, int frameless, int resizable, int fullscreen, int fullSizeContent, int hideTitleBar, int titlebarAppearsTransparent, int hideTitle, int useToolbar, int hideToolbarSeparator, int webviewIsTransparent, int alwaysOnTop, int hideWindowOnClose, const char *appearance, int windowIsTranslucent, int debug);
void Run(void*);

void SetTitle(WailsContext *ctx, const char *title);
void Center(WailsContext *ctx);
void SetSize(WailsContext *ctx, int width, int height);
void SetMinWindowSize(WailsContext *ctx, int width, int height);
void SetMaxWindowSize(WailsContext *ctx, int width, int height);
void SetPosition(WailsContext *ctx, int x, int y);
void Fullscreen(WailsContext *ctx);
void UnFullscreen(WailsContext *ctx);
void Minimise(WailsContext *ctx);
void UnMinimise(WailsContext *ctx);

void SetRGBA(void *ctx, int r, int g, int b, int a);

void Quit(void*);

void ProcessURLResponse(void *inctx, const char *url, const char *contentType, Byte data[], int datalength);

#endif /* Application_h */
