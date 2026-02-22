#import "webcontentsview_darwin.h"

void* createWebContentsView(int x, int y, int w, int h, WebContentsViewPreferences prefs) {
    NSRect frame = NSMakeRect(x, y, w, h);
    WKWebViewConfiguration* config = [[WKWebViewConfiguration alloc] init];
    
    WKPreferences *preferences = [[WKPreferences alloc] init];
    
    // Apply preferences
    if (@available(macOS 10.11, *)) {
        [preferences setValue:@(prefs.devTools) forKey:@"developerExtrasEnabled"];
    }
    
    // JavaScript
    if (@available(macOS 11.0, *)) {
        WKWebpagePreferences *webpagePreferences = [[WKWebpagePreferences alloc] init];
        webpagePreferences.allowsContentJavaScript = prefs.javascript;
        config.defaultWebpagePreferences = webpagePreferences;
    } else {
        preferences.javaScriptEnabled = prefs.javascript;
    }
    
    // WebSecurity / CORS
    if (!prefs.webSecurity) {
        [preferences setValue:@YES forKey:@"allowFileAccessFromFileURLs"];
        [preferences setValue:@YES forKey:@"allowUniversalAccessFromFileURLs"];
    }
    
    // Images
    [preferences setValue:@(prefs.images) forKey:@"loadsImagesAutomatically"];
    
    // Plugins
    [preferences setValue:@(prefs.plugins) forKey:@"plugInsEnabled"];
    
    // Fonts
    if (prefs.defaultFontSize > 0) {
        [preferences setValue:@(prefs.defaultFontSize) forKey:@"defaultFontSize"];
    }
    if (prefs.minimumFontSize > 0) {
        preferences.minimumFontSize = prefs.minimumFontSize;
    }
    
    config.preferences = preferences;

    WKWebView* webView = [[WKWebView alloc] initWithFrame:frame configuration:config];
    
    // Zoom factor
    if (prefs.zoomFactor != 1.0 && prefs.zoomFactor > 0) {
        if (@available(macOS 11.0, *)) {
            [webView setPageZoom:prefs.zoomFactor];
        } else {
            [webView setMagnification:prefs.zoomFactor];
        }
    }
    
    return webView;
}

void webContentsViewSetBounds(void* view, int x, int y, int w, int h) {
    WKWebView* webView = (WKWebView*)view;
    [webView setFrame:NSMakeRect(x, y, w, h)];
}

void webContentsViewSetURL(void* view, const char* url) {
    WKWebView* webView = (WKWebView*)view;
    NSString* nsURL = [NSString stringWithUTF8String:url];
    NSURLRequest* request = [NSURLRequest requestWithURL:[NSURL URLWithString:nsURL]];
    [webView loadRequest:request];
}

void webContentsViewExecJS(void* view, const char* js) {
    WKWebView* webView = (WKWebView*)view;
    NSString* script = [NSString stringWithUTF8String:js];
    [webView evaluateJavaScript:script completionHandler:nil];
}

void windowAddWebContentsView(void* nsWindow, void* view) {
    NSWindow* window = (NSWindow*)nsWindow;
    WKWebView* webView = (WKWebView*)view;
    [window.contentView addSubview:webView positioned:NSWindowAbove relativeTo:nil];
}

void windowRemoveWebContentsView(void* nsWindow, void* view) {
    WKWebView* webView = (WKWebView*)view;
    [webView removeFromSuperview];
}
