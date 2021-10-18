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

WailsContext* Create(const char* title, int width, int height, int frameless, int resizable, int fullscreen, int fullSizeContent, int hideTitleBar, int titlebarAppearsTransparent, int hideTitle, int useToolbar, int hideToolbarSeparator, int webviewIsTransparent, int alwaysOnTop, int hideWindowOnClose, const char *appearance, int windowIsTranslucent, int debug);
void Run(void*);

void SetTitle(void* ctx, const char *title);
void Center(void* ctx);
void SetSize(void* ctx, int width, int height);
void SetMinSize(void* ctx, int width, int height);
void SetMaxSize(void* ctx, int width, int height);
void SetPosition(void* ctx, int x, int y);
void Fullscreen(void* ctx);
void UnFullscreen(void* ctx);
void Minimise(void* ctx);
void UnMinimise(void* ctx);
void Maximise(void* ctx);
void UnMaximise(void* ctx);
void Hide(void* ctx);
void Show(void* ctx);
void SetRGBA(void* ctx, int r, int g, int b, int a);
void ExecJS(void* ctx, const char*);
void Quit(void*);

const char* GetSize(void *ctx);
const char* GetPos(void *ctx);

void ProcessURLResponse(void *inctx, const char *url, const char *contentType, const char *data, int datalength);

#endif /* Application_h */
