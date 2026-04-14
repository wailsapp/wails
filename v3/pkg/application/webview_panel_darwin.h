//go:build darwin && !ios

#ifndef WebviewPanel_h
#define WebviewPanel_h

#import "webview_window_darwin.h"

@interface WebviewPanel : NSPanel
- (BOOL) canBecomeKeyWindow;
- (BOOL) canBecomeMainWindow;
- (BOOL) acceptsFirstResponder;
- (BOOL) becomeFirstResponder;
- (BOOL) resignFirstResponder;
- (WebviewPanel*) initWithContentRect:(NSRect)contentRect styleMask:(NSUInteger)windowStyle backing:(NSBackingStoreType)bufferingType defer:(BOOL)deferCreation;

@property (assign) WKWebView* webView;

@end

#endif /* WebviewPanel_h */
