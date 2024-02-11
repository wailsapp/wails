#ifndef WailsWebView_h
#define WailsWebView_h

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

// We will override WKWebView, so we can detect file drop in obj-c
//  and grab their file path, to then inject into JS
@interface WailsWebView : WKWebView
@property bool disableWebViewDragAndDrop;
@property bool enableDragAndDrop;
@end

#endif /* WailsWebView_h */
