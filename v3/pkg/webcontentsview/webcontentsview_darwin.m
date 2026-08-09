#import "webcontentsview_darwin.h"

extern void browserViewSnapshotCallback(uintptr_t callbackID, const char* base64);

void* createWebContentsView(int x, int y, int w, int h, WebContentsViewPreferences prefs) {
    WKWebViewConfiguration* config = [[WKWebViewConfiguration alloc] init];
    WKPreferences* preferences = [[WKPreferences alloc] init];
    @try {
        if (@available(macOS 10.11, *)) {
            [preferences setValue:@(prefs.devTools) forKey:@"developerExtrasEnabled"];
        }
        if (@available(macOS 11.0, *)) {
            WKWebpagePreferences* webpagePreferences = [[WKWebpagePreferences alloc] init];
            webpagePreferences.allowsContentJavaScript = prefs.javascript;
            config.defaultWebpagePreferences = webpagePreferences;
        } else {
            preferences.javaScriptEnabled = prefs.javascript;
        }
        config.preferences = preferences;
    } @catch (NSException* ignored) {
    }
    WKWebView* webView = [[WKWebView alloc] initWithFrame:NSMakeRect(x, y, w, h) configuration:config];
    if (prefs.userAgent != NULL) {
        webView.customUserAgent = [NSString stringWithUTF8String:prefs.userAgent];
    }
    return webView;
}

void webContentsViewSetBounds(void* view, int x, int y, int w, int h) {
    WKWebView* webView = (WKWebView*)view;
    NSView* superview = webView.superview;
    CGFloat cocoaY = superview == nil ? y : superview.bounds.size.height - y - h;
    webView.frame = NSMakeRect(x, cocoaY, w, h);
}

void webContentsViewSetVisible(void* view, bool visible) { [(WKWebView*)view setHidden:!visible]; }

void webContentsViewDestroy(void* view) {
    WKWebView* webView = (WKWebView*)view;
    if (webView == nil) return;
    [webView stopLoading];
    [webView removeFromSuperview];
    [webView release];
}

void webContentsViewSetURL(void* view, const char* url) {
    WKWebView* webView = (WKWebView*)view;
    NSString* value = [NSString stringWithUTF8String:url];
    dispatch_async(dispatch_get_main_queue(), ^{ [webView loadRequest:[NSURLRequest requestWithURL:[NSURL URLWithString:value]]]; });
}

void webContentsViewSetHTML(void* view, const char* html) {
    WKWebView* webView = (WKWebView*)view;
    NSString* value = html == NULL ? @"" : [NSString stringWithUTF8String:html];
    dispatch_async(dispatch_get_main_queue(), ^{ [webView loadHTMLString:value baseURL:nil]; });
}

void webContentsViewExecJS(void* view, const char* js) {
    WKWebView* webView = (WKWebView*)view;
    NSString* value = [NSString stringWithUTF8String:js];
    dispatch_async(dispatch_get_main_queue(), ^{ [webView evaluateJavaScript:value completionHandler:nil]; });
}

void windowAddWebContentsView(void* nsWindow, void* view) {
    WKWebView* webView = (WKWebView*)view;
    [((NSWindow*)nsWindow).contentView addSubview:webView];
    webView.wantsLayer = YES;
    webView.layer.zPosition = 9999.0;
}

void windowRemoveWebContentsView(void* nsWindow, void* view) { [(WKWebView*)view removeFromSuperview]; }

void webContentsViewGoBack(void* view) {
    WKWebView* webView = (WKWebView*)view;
    dispatch_async(dispatch_get_main_queue(), ^{ if (webView.canGoBack) [webView goBack]; });
}

const char* webContentsViewGetURL(void* view) {
    __block const char* result = NULL;
    void (^readURL)(void) = ^{ WKWebView* webView = (WKWebView*)view; if (webView.URL != nil) result = strdup(webView.URL.absoluteString.UTF8String); };
    if ([NSThread isMainThread]) readURL(); else dispatch_sync(dispatch_get_main_queue(), readURL);
    return result;
}

void webContentsViewTakeSnapshot(void* view, uintptr_t callbackID) {
    WKWebView* webView = (WKWebView*)view;
    dispatch_async(dispatch_get_main_queue(), ^{
        if (@available(macOS 10.13, *)) {
            [webView takeSnapshotWithConfiguration:[[WKSnapshotConfiguration alloc] init] completionHandler:^(NSImage* image, NSError* error) {
                if (error != nil || image == nil) { browserViewSnapshotCallback(callbackID, NULL); return; }
                NSBitmapImageRep* bitmap = [[NSBitmapImageRep alloc] initWithCGImage:[image CGImageForProposedRect:NULL context:nil hints:nil]];
                NSData* png = [bitmap representationUsingType:NSBitmapImageFileTypePNG properties:@{}];
                NSString* dataURL = png == nil ? nil : [NSString stringWithFormat:@"data:image/png;base64,%@", [png base64EncodedStringWithOptions:0]];
                browserViewSnapshotCallback(callbackID, dataURL == nil ? NULL : dataURL.UTF8String);
            }];
        } else {
            browserViewSnapshotCallback(callbackID, NULL);
        }
    });
}
