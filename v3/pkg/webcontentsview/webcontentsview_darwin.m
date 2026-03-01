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
    
    if (prefs.userAgent != NULL) {
        NSString* customUA = [NSString stringWithUTF8String:prefs.userAgent];
        [webView setCustomUserAgent:customUA];
    }
    
    return webView;
}

void webContentsViewSetBounds(void* view, int x, int y, int w, int h) {
    WKWebView* webView = (WKWebView*)view;
    NSView* superview = [webView superview];
    
    if (superview != nil) {
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
    [window.contentView addSubview:webView];
    [webView setWantsLayer:YES];
    webView.layer.zPosition = 9999.0;
}

void windowRemoveWebContentsView(void* nsWindow, void* view) {
    WKWebView* webView = (WKWebView*)view;
    [webView removeFromSuperview];
}

void webContentsViewGoBack(void* view) {
    WKWebView* webView = (WKWebView*)view;
    dispatch_async(dispatch_get_main_queue(), ^{
        if ([webView canGoBack]) {
            [webView goBack];
        }
    });
}

const char* webContentsViewGetURL(void* view) {
    __block const char* result = NULL;
    WKWebView* webView = (WKWebView*)view;
    
    if ([NSThread isMainThread]) {
        if (webView.URL != nil) {
            result = strdup(webView.URL.absoluteString.UTF8String);
        }
    } else {
        dispatch_sync(dispatch_get_main_queue(), ^{
            if (webView.URL != nil) {
                result = strdup(webView.URL.absoluteString.UTF8String);
            }
        });
    }
    return result;
}

extern void browserViewSnapshotCallback(uintptr_t callbackID, const char* base64);

void webContentsViewTakeSnapshot(void* view, uintptr_t callbackID) {
    WKWebView* webView = (WKWebView*)view;
    dispatch_async(dispatch_get_main_queue(), ^{
        @try {
            if (@available(macOS 10.13, *)) {
                WKSnapshotConfiguration *config = [[WKSnapshotConfiguration alloc] init];
                [webView takeSnapshotWithConfiguration:config completionHandler:^(NSImage *image, NSError *error) {
                    if (error != nil || image == nil) {
                        browserViewSnapshotCallback(callbackID, NULL);
                        return;
                    }
                    
                    @try {
                        CGImageRef cgRef = [image CGImageForProposedRect:NULL context:nil hints:nil];
                        if (cgRef == NULL) {
                            browserViewSnapshotCallback(callbackID, NULL);
                            return;
                        }
                        
                        NSBitmapImageRep *newRep = [[NSBitmapImageRep alloc] initWithCGImage:cgRef];
                        if (newRep == nil) {
                            browserViewSnapshotCallback(callbackID, NULL);
                            return;
                        }
                        
                        [newRep setSize:[image size]];
                        NSData *pngData = [newRep representationUsingType:NSBitmapImageFileTypePNG properties:@{}];
                        
                        if (pngData == nil) {
                            browserViewSnapshotCallback(callbackID, NULL);
                            return;
                        }
                        
                        NSString *base64String = [pngData base64EncodedStringWithOptions:0];
                        if (base64String == nil) {
                            browserViewSnapshotCallback(callbackID, NULL);
                            return;
                        }
                        
                        NSString *fullBase64 = [NSString stringWithFormat:@"data:image/png;base64,%@", base64String];
                        browserViewSnapshotCallback(callbackID, [fullBase64 UTF8String]);
                    } @catch (NSException *innerException) {
                        NSLog(@"Error processing snapshot image: %@", innerException.reason);
                        browserViewSnapshotCallback(callbackID, NULL);
                    }
                }];
            } else {
                browserViewSnapshotCallback(callbackID, NULL);
            }
        } @catch (NSException *e) {
            NSLog(@"Exception in snapshot configuration: %@", e.reason);
            browserViewSnapshotCallback(callbackID, NULL);
        }
    });
}
