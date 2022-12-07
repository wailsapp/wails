//go:build darwin


#ifndef WindowDelegate_h
#define WindowDelegate_h

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface WindowDelegate : NSObject <NSWindowDelegate>

@property bool hideOnClose;
@property (retain) WKWebView* webView;

@end


#endif /* WindowDelegate_h */
