//go:build darwin


#ifndef WindowDelegate_h
#define WindowDelegate_h

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

extern void processMessage(unsigned int, const char*);

@interface WindowDelegate : NSObject <NSWindowDelegate, WKScriptMessageHandler>

@property bool hideOnClose;
@property (retain) WKWebView* webView;
@property unsigned int windowId;

@end


#endif /* WindowDelegate_h */
