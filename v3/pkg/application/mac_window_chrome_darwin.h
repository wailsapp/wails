//go:build darwin && !ios && !server

#ifndef WailsMacWindowChromeDarwin_h
#define WailsMacWindowChromeDarwin_h

#import <Cocoa/Cocoa.h>

// Shared NSWindow chrome operations. These deliberately accept NSWindow,
// rather than WebviewWindow, so toolbars work in both normal and native builds.
void windowSetToolbarStyle(void* nsWindow, int style);
void windowSetHideToolbarSeparator(void* nsWindow, bool hideSeparator);
void windowSetShowToolbarWhenFullscreen(void* nsWindow, bool setting);

#endif
