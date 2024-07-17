//go:build darwin

#import <AppKit/AppKit.h>
#import <WebKit/WebKit.h>

#define DragAndDropTypeNone 0
#define DragAndDropTypeWebview 1
#define DragAndDropTypeWindow 2

@interface WailsWebView : WKWebView
@property unsigned int windowId;
@property unsigned int dragAndDropType;
@end