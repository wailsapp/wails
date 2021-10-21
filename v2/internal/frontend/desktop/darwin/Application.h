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

void MessageDialog(void *inctx, const char* dialogType, const char* title, const char* message, const char* button1, const char* button2, const char* button3, const char* button4, const char* defaultButton, const char* cancelButton);
void OpenFileDialog(void *inctx, const char* title, const char* defaultFilename, const char* defaultDirectory, int allowDirectories, int allowFiles, int canCreateDirectories, int treatPackagesAsDirectories, int resolveAliases, int showHiddenFiles, int allowMultipleSelection, const char* filters);
void SaveFileDialog(void *inctx, const char* title, const char* defaultFilename, const char* defaultDirectory, int canCreateDirectories, int treatPackagesAsDirectories, int showHiddenFiles, const char* filters);
#endif /* Application_h */
