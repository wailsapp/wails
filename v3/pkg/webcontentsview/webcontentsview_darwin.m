#import "webcontentsview_darwin.h"

void* createWebContentsView(int x, int y, int w, int h, WebContentsViewPreferences prefs) {
    NSRect frame = NSMakeRect(x, y, w, h);
    WKWebViewConfiguration* config = [[WKWebViewConfiguration alloc] init];
    
    WKPreferences *preferences = [[WKPreferences alloc] init];
    
    @try {
        if (@available(macOS 10.11, *)) {
            [preferences setValue:@(prefs.devTools) forKey:@"developerExtrasEnabled"];
        }
        
        if (@available(macOS 11.0, *)) {
            WKWebpagePreferences *webpagePreferences = [[WKWebpagePreferences alloc] init];
            webpagePreferences.allowsContentJavaScript = prefs.javascript;
            config.defaultWebpagePreferences = webpagePreferences;
        } else {
            preferences.javaScriptEnabled = prefs.javascript;
        }
        
        if (!prefs.webSecurity) {
            @try {
                [config.preferences setValue:@YES forKey:@"allowFileAccessFromFileURLs"];
                [config.preferences setValue:@YES forKey:@"allowUniversalAccessFromFileURLs"];
            } @catch (NSException *e) {}
        }
        config.preferences = preferences;
    } @catch (NSException *e) {}

    WKWebView* webView = [[WKWebView alloc] initWithFrame:frame configuration:config];
    [webView setValue:@(NO) forKey:@"drawsTransparentBackground"];
    
    return webView;
}

void webContentsViewSetBounds(void* view, int x, int y, int w, int h) {
    WKWebView* webView = (WKWebView*)view;
    NSView* superview = [webView superview];
    
    if (superview != nil) {
        // macOS standard coordinates: 0,0 is bottom-left
        CGFloat superHeight = superview.bounds.size.height;
        CGFloat cocoaY = superHeight - y - h;
        [webView setFrame:NSMakeRect(x, cocoaY, w, h)];
    } else {
        [webView setFrame:NSMakeRect(x, y, w, h)];
    }
}

void webContentsViewSetURL(void* view, const char* url) {
    WKWebView* webView = (WKWebView*)view;
    NSString* nsURL = [NSString stringWithUTF8String:url];
    NSURLRequest* request = [NSURLRequest requestWithURL:[NSURL URLWithString:nsURL]];
    dispatch_async(dispatch_get_main_queue(), ^{
        [webView loadRequest:request];
    });
}

void webContentsViewExecJS(void* view, const char* js) {
    WKWebView* webView = (WKWebView*)view;
    NSString* script = [NSString stringWithUTF8String:js];
    dispatch_async(dispatch_get_main_queue(), ^{
        [webView evaluateJavaScript:script completionHandler:nil];
    });
}

void windowAddWebContentsView(void* nsWindow, void* view) {
    NSWindow* window = (NSWindow*)nsWindow;
    WKWebView* webView = (WKWebView*)view;
    
    // Add directly to the window's root view so we avoid clipping issues
    [window.contentView addSubview:webView];
    
    // Force to front by giving it a high zPosition
    [webView setWantsLayer:YES];
    webView.layer.zPosition = 9999.0;
}

void windowRemoveWebContentsView(void* nsWindow, void* view) {
    WKWebView* webView = (WKWebView*)view;
    [webView removeFromSuperview];
}
