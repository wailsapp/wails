//go:build darwin && !ios && !server

#import "webview_window_split_darwin.h"
#import "webview_window_darwin.h"

void* splitWindowInstall(void* nsWindow, bool vertical) {
    NSSplitViewController* splitVC = [[[NSSplitViewController alloc] init] autorelease];
    splitVC.splitView.vertical = vertical;
    // window.contentViewController is a strong property: it retains splitVC
    // for as long as the window keeps it installed, so the autorelease
    // above is safe -- Go only ever holds a non-owning pointer to it.
    ((WebviewWindow*)nsWindow).contentViewController = splitVC;
    return (void*)splitVC;
}

void splitWindowSetAutosaveName(void* splitViewController, const char* name) {
    NSSplitViewController* splitVC = (NSSplitViewController*)splitViewController;
    splitVC.splitView.autosaveName = [NSString stringWithUTF8String:name];
}

// Behavior values must match MacSplitPaneBehavior's iota order in Go:
// 0=Default, 1=Sidebar, 2=ContentList, 3=Inspector.
static NSSplitViewItem* buildSplitViewItem(NSViewController* paneVC, int behavior) {
    switch (behavior) {
        case 1:
            return [NSSplitViewItem sidebarWithViewController:paneVC];
        case 2:
            return [NSSplitViewItem contentListWithViewController:paneVC];
        case 3:
            if (@available(macOS 11.0, *)) {
                return [NSSplitViewItem inspectorWithViewController:paneVC];
            }
            return [NSSplitViewItem splitViewItemWithViewController:paneVC];
        default:
            return [NSSplitViewItem splitViewItemWithViewController:paneVC];
    }
}

void* splitWindowAddPane(void* splitViewController, void* nsWindow, bool reuseExistingWebview,
    const char* url, int behavior,
    double minThickness, double maxThickness, double preferredThicknessFraction, double holdingPriority,
    bool collapsible, bool startCollapsed) {

    NSSplitViewController* splitVC = (NSSplitViewController*)splitViewController;

    WKWebView* webView;
    if (reuseExistingWebview) {
        webView = ((WebviewWindow*)nsWindow).webView;
        // window.webView is `assign` (non-owning): the view hierarchy was
        // its only strong reference. Bump it before detaching so it
        // survives the gap until the new pane view controller retains it.
        [webView retain];
        [webView removeFromSuperview];
    } else {
        webView = buildAuxiliaryWebview(nsWindow, 200, 200); // +1
    }

    // Sidebar/ContentList/Inspector panes get a translucent NSVisualEffectView
    // material behind them automatically from the factory methods below, but
    // WKWebView paints its own opaque white backing by default -- setting the
    // page's CSS background to transparent alone does nothing at the native
    // layer, the material never shows through. This is the other half of
    // "automatic glass with no extra code": it needs to be told to stop
    // drawing its own background, not just asked nicely from CSS.
    if (behavior == 1 || behavior == 2 || behavior == 3) {
        if (@available(macOS 12.0, *)) {
            webView.underPageBackgroundColor = NSColor.clearColor;
        }
    }

    NSViewController* paneVC = [[[NSViewController alloc] init] autorelease];
    paneVC.view = webView;

    NSSplitViewItem* item = buildSplitViewItem(paneVC, behavior);
    item.canCollapse = collapsible;
    if (minThickness > 0) item.minimumThickness = minThickness;
    if (maxThickness > 0) item.maximumThickness = maxThickness;
    if (preferredThicknessFraction > 0) item.preferredThicknessFraction = preferredThicknessFraction;
    if (holdingPriority > 0) item.holdingPriority = holdingPriority;
    if (startCollapsed) item.collapsed = YES;

    [splitVC addSplitViewItem:item]; // splitVC.splitViewItems now owns item

    if (!reuseExistingWebview && url != NULL && strlen(url) > 0) {
        // Loading only after the webview is actually in the split view's
        // hierarchy -- see the comment on buildAuxiliaryWebview. url must
        // already be a fully-qualified wails://... URL (resolved Go-side
        // via assetserver.GetStartURL); a bare path has no scheme for
        // WKWebView to route to the custom scheme handler.
        NSURL* nsURL = [NSURL URLWithString:[NSString stringWithUTF8String:url]];
        [webView loadRequest:[NSURLRequest requestWithURL:nsURL]];
    }

    [webView release]; // release our +1 (new pane) or the temporary retain (reused pane);
                        // paneVC.view now owns the only reference that matters.

    return (void*)item;
}

void* splitPaneWebview(void* splitViewItem) {
    NSSplitViewItem* item = (NSSplitViewItem*)splitViewItem;
    return (void*)item.viewController.view;
}

void splitPaneSetCollapsed(void* splitViewItem, bool collapsed) {
    ((NSSplitViewItem*)splitViewItem).collapsed = collapsed;
}

bool splitPaneIsCollapsed(void* splitViewItem) {
    return ((NSSplitViewItem*)splitViewItem).isCollapsed;
}

void splitPaneWebviewSetURL(void* webView, const char* url) {
    NSURL* nsURL = [NSURL URLWithString:[NSString stringWithUTF8String:url]];
    [(WKWebView*)webView loadRequest:[NSURLRequest requestWithURL:nsURL]];
}

void splitPaneWebviewExecJS(void* webView, const char* js) {
    [(WKWebView*)webView evaluateJavaScript:[NSString stringWithUTF8String:js] completionHandler:nil];
}

void splitPaneWebviewReload(void* webView) {
    [(WKWebView*)webView reload];
}
