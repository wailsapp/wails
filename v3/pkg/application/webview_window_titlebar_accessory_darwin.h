//go:build darwin && !ios && !server

#ifndef WailsTitlebarAccessory_h
#define WailsTitlebarAccessory_h

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

// Builds an NSTitlebarAccessoryViewController hosting its own WKWebView
// (via buildAuxiliaryWebview, so it shares the window's WebviewWindowDelegate
// and windowId-keyed asset/IPC routing) and adds it to nsWindow. position
// values must match MacTitlebarAccessoryPosition's iota order in Go:
// 0=Leading, 1=Trailing. Returns the NSTitlebarAccessoryViewController*.
void* titlebarAccessoryAdd(void* nsWindow, int position, double width, double height, const char* url);

// Returns the WKWebView* hosted by a titlebar accessory controller.
void* titlebarAccessoryWebview(void* controller);

// Removes a titlebar accessory previously added via titlebarAccessoryAdd.
void titlebarAccessoryRemove(void* nsWindow, void* controller);

void titlebarAccessoryWebviewSetURL(void* webView, const char* url);
void titlebarAccessoryWebviewExecJS(void* webView, const char* js);

#endif /* WailsTitlebarAccessory_h */
