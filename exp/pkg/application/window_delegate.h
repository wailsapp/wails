//go:build darwin

#ifndef WindowDelegate_h
#define WindowDelegate_h

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface WindowDelegate : NSObject <NSWindowDelegate, WKScriptMessageHandler, WKNavigationDelegate>

@property bool hideOnClose;
@property (retain) WKWebView* webView;
@property unsigned int windowId;

@end


#endif /* WindowDelegate_h */
