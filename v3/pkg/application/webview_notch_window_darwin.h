//go:build darwin && !ios && !server

#ifndef WebviewNotchWindow_h
#define WebviewNotchWindow_h

#import "webview_panel_darwin.h"

// WebviewNotchWindow is the native shaped panel used by NewNotchWindow. The
// WKWebView occupies only the caller-requested inner area; AppKit draws the
// black attachment wings around it.
@interface WebviewNotchWindow : WebviewPanel

- (void)configureWebView:(WKWebView*)webView
            contentWidth:(CGFloat)contentWidth
           contentHeight:(CGFloat)contentHeight
          targetScreenID:(NSString*)targetScreenID;
- (void)showAnimated:(BOOL)animated duration:(NSTimeInterval)duration;
- (void)hideAnimated:(BOOL)animated duration:(NSTimeInterval)duration;

@end

#endif /* WebviewNotchWindow_h */
