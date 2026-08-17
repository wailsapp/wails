//go:build darwin && !ios && !server

#ifndef WailsNativeWindowDarwin_h
#define WailsNativeWindowDarwin_h

#import <Cocoa/Cocoa.h>

extern void processNativeWindowClosed(unsigned int windowID);

void* nativeWindowCreate(unsigned int windowID, int width, int height, bool hideOnClose);
void nativeWindowDestroy(void* nsWindow);
void nativeWindowSetTitle(void* nsWindow, const char* title);
void nativeWindowSetResizable(void* nsWindow, bool resizable);
void nativeWindowSetMinSize(void* nsWindow, int width, int height);
void nativeWindowSetMaxSize(void* nsWindow, int width, int height);
void nativeWindowSetAlwaysOnTop(void* nsWindow, bool alwaysOnTop);
void nativeWindowConfigureTitlebar(void* nsWindow, bool appearsTransparent,
    bool fullSizeContent, bool hideTitle, bool hideToolbarSeparator, int toolbarStyle);
void nativeWindowShow(void* nsWindow);
void nativeWindowHide(void* nsWindow);
void nativeWindowFocus(void* nsWindow);
bool nativeWindowIsVisible(void* nsWindow);
void nativeWindowCenter(void* nsWindow);
void nativeWindowSetPosition(void* nsWindow, int x, int y);

#endif
