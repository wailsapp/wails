//go:build darwin && !ios && !server

#import "webview_window_titlebar_accessory_darwin.h"
#import "webview_window_darwin.h"

void* titlebarAccessoryAdd(void* nsWindow, int position, double width, double height, const char* url) {
    WKWebView* webView = buildAuxiliaryWebview(nsWindow, width, height); // +1
    // Unlike a split pane, a titlebar accessory always sits directly on the
    // titlebar's own background -- an opaque WKWebView here isn't a style
    // option, it renders as a visibly wrong white/black rectangle floating
    // in the titlebar. Always clear it, not just for material-backed panes.
    if (@available(macOS 12.0, *)) {
        webView.underPageBackgroundColor = NSColor.clearColor;
    }

    NSTitlebarAccessoryViewController* controller = [[NSTitlebarAccessoryViewController alloc] init];
    [controller autorelease];
    // position values match MacTitlebarAccessoryPosition's iota order:
    // 0=Leading, 1=Trailing.
    controller.layoutAttribute = (position == 1) ? NSLayoutAttributeTrailing : NSLayoutAttributeLeading;
    controller.view = webView;

    [(WebviewWindow*)nsWindow addTitlebarAccessoryViewController:controller];

    if (url != NULL && strlen(url) > 0) {
        // Loading only after the accessory is actually installed, same
        // reasoning as splitWindowAddPane: a webview with no window/view
        // hierarchy yet silently suspends WKWebView's loading.
        NSURL* nsURL = [NSURL URLWithString:[NSString stringWithUTF8String:url]];
        [webView loadRequest:[NSURLRequest requestWithURL:nsURL]];
    }

    [webView release]; // release our +1; controller.view now owns the only reference that matters.

    return (void*)controller; // retained by nsWindow.titlebarAccessoryViewControllers
}

void* titlebarAccessoryWebview(void* controller) {
    return (void*)((NSTitlebarAccessoryViewController*)controller).view;
}

void titlebarAccessoryRemove(void* nsWindow, void* controller) {
    NSWindow* window = (NSWindow*)nsWindow;
    NSInteger index = [window.titlebarAccessoryViewControllers indexOfObject:(NSTitlebarAccessoryViewController*)controller];
    if (index != NSNotFound) {
        [window removeTitlebarAccessoryViewControllerAtIndex:index];
    }
}

void titlebarAccessoryWebviewSetURL(void* webView, const char* url) {
    NSURL* nsURL = [NSURL URLWithString:[NSString stringWithUTF8String:url]];
    [(WKWebView*)webView loadRequest:[NSURLRequest requestWithURL:nsURL]];
}

void titlebarAccessoryWebviewExecJS(void* webView, const char* js) {
    [(WKWebView*)webView evaluateJavaScript:[NSString stringWithUTF8String:js] completionHandler:nil];
}
