//go:build ios

#import "webview_window_ios.h"
#import "application_ios.h"
#import "application_ios_delegate.h"
#import <stdlib.h>

// Buffer console messages until a WKWebView exists
static NSMutableArray<NSString *> *pendingConsoleJS;

// Subclass that optionally hides the input accessory toolbar based on global flag
@interface WailsWebView : WKWebView @end
@implementation WailsWebView
- (UIView *)inputAccessoryView {
    if (ios_is_input_accessory_disabled()) {
        return nil;
    }
    return [super inputAccessoryView];
}
@end

// MARK: - WailsSchemeHandler
@implementation WailsSchemeHandler

- (instancetype)initWithWindowID:(unsigned int)windowID {
    self = [super init];
    if (self) {
        _windowID = windowID;
    }
    return self;
}

- (void)webView:(WKWebView *)webView startURLSchemeTask:(id<WKURLSchemeTask>)urlSchemeTask {
    NSURL *url = urlSchemeTask.request.URL;
    NSLog(@"[WailsSchemeHandler] start task: %@", url.absoluteString);
    ServeAssetRequest(self.windowID, (__bridge void*)urlSchemeTask);
}

- (void)webView:(WKWebView *)webView stopURLSchemeTask:(id<WKURLSchemeTask>)urlSchemeTask {
    NSLog(@"[WailsSchemeHandler] stop task");
}

@end

// MARK: - WailsMessageHandler
@implementation WailsMessageHandler

- (instancetype)initWithWindowID:(unsigned int)windowID {
    self = [super init];
    if (self) {
        _windowID = windowID;
    }
    return self;
}

- (void)userContentController:(WKUserContentController *)userContentController didReceiveScriptMessage:(WKScriptMessage *)message {
    // Support both plain string messages and structured objects
    if ([message.body isKindOfClass:[NSString class]]) {
        NSString *msg = (NSString *)message.body;
        HandleJSMessage(self.windowID, (char *)[msg UTF8String]);
        return;
    }
    NSError *error = nil;
    NSData *jsonData = [NSJSONSerialization dataWithJSONObject:message.body options:0 error:&error];
    if (!error && jsonData) {
        NSString *jsonString = [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
        HandleJSMessage(self.windowID, (char *)[jsonString UTF8String]);
    } else {
        // Fallback: attempt to stringify non-serializable payloads
        NSString *desc = [NSString stringWithFormat:@"%@", message.body];
        HandleJSMessage(self.windowID, (char *)[desc UTF8String]);
    }
}

@end

// MARK: - WailsViewController
@implementation WailsViewController

- (instancetype)initWithWindowID:(unsigned int)windowID {
    self = [super init];
    if (self) {
        _windowID = windowID;
    }
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];

    WKWebViewConfiguration *config = [[WKWebViewConfiguration alloc] init];
    config.suppressesIncrementalRendering = YES;
    // Application name for UA (default to "wails.io" if not set)
    const char* appNameForUA = ios_get_app_name_for_user_agent();
    config.applicationNameForUserAgent = appNameForUA ? [NSString stringWithUTF8String:appNameForUA] : @"wails.io";
    // Enable JavaScript using modern API (javaScriptEnabled is deprecated)
    if (@available(iOS 14.0, *)) {
        config.defaultWebpagePreferences.allowsContentJavaScript = YES;
    } else {
        // Fallback for very old iOS versions
        #pragma clang diagnostic push
        #pragma clang diagnostic ignored "-Wdeprecated-declarations"
        config.preferences.javaScriptEnabled = YES;
        #pragma clang diagnostic pop
    }
    // Media playback
    config.allowsInlineMediaPlayback = ios_is_inline_media_playback_enabled();
    if (ios_is_autoplay_without_user_action_enabled()) {
        config.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeNone;
    } else {
        config.mediaTypesRequiringUserActionForPlayback = WKAudiovisualMediaTypeAll;
    }

    // URL scheme handler and script bridge
    self.schemeHandler = [[WailsSchemeHandler alloc] initWithWindowID:self.windowID];
    [config setURLSchemeHandler:self.schemeHandler forURLScheme:@"wails"];
    self.messageHandler = [[WailsMessageHandler alloc] initWithWindowID:self.windowID];
    // Register both handler names used by runtimes: "external" (current runtime) and "wails" (legacy)
    [config.userContentController addScriptMessageHandler:self.messageHandler name:@"external"];
    [config.userContentController addScriptMessageHandler:self.messageHandler name:@"wails"];

    self.webView = [[WailsWebView alloc] initWithFrame:self.view.bounds configuration:config];
    // Custom user agent if provided
    const char* userAgent = ios_get_user_agent();
    if (userAgent) {
        self.webView.customUserAgent = [NSString stringWithUTF8String:userAgent];
    }
    self.webView.autoresizingMask = UIViewAutoresizingFlexibleWidth | UIViewAutoresizingFlexibleHeight;
    self.webView.navigationDelegate = self;
    // Back/forward gestures
    self.webView.allowsBackForwardNavigationGestures = ios_is_back_forward_gestures_enabled();
    // Link preview
    self.webView.allowsLinkPreview = ios_is_link_preview_disabled() ? NO : YES;

    // Configure scrolling & bounce & indicators
    UIScrollView *sv = self.webView.scrollView;
    bool scrollDisabled = ios_is_scroll_disabled();
    bool bounceDisabled = ios_is_bounce_disabled();
    bool indicatorsDisabled = ios_is_scroll_indicators_disabled();
    sv.scrollEnabled = scrollDisabled ? NO : YES;
    sv.bounces = bounceDisabled ? NO : YES;
    sv.alwaysBounceVertical = bounceDisabled ? NO : YES;
    sv.alwaysBounceHorizontal = bounceDisabled ? NO : YES;
    sv.showsVerticalScrollIndicator = indicatorsDisabled ? NO : YES;
    sv.showsHorizontalScrollIndicator = indicatorsDisabled ? NO : YES;
    sv.contentInset = UIEdgeInsetsZero;
    sv.scrollIndicatorInsets = UIEdgeInsetsZero;
    if (@available(iOS 11.0, *)) {
        sv.contentInsetAdjustmentBehavior = UIScrollViewContentInsetAdjustmentNever;
    }

    // Inspector
    BOOL inspectorOn = ios_is_inspectable_disabled() ? NO : YES;
    if (@available(iOS 16.4, *)) {
        self.webView.inspectable = inspectorOn;
    } else {
        @try { [self.webView setValue:@(inspectorOn) forKey:@"inspectable"]; } @catch (__unused NSException *e) {}
    }

    [self.view addSubview:self.webView];

    // Initial load triggers our scheme handler
    NSURLRequest *request = [NSURLRequest requestWithURL:[NSURL URLWithString:@"wails://localhost/"]];
    [self.webView loadRequest:request];

    // Flush any pending console logs now that a webview exists
    dispatch_async(dispatch_get_main_queue(), ^{
        if (pendingConsoleJS.count > 0) {
            for (NSString *js in pendingConsoleJS) {
                [self.webView evaluateJavaScript:js completionHandler:nil];
            }
            [pendingConsoleJS removeAllObjects];
        }
    });

    // Enable native tabs if globally enabled
    BOOL tabsEnabled = ios_native_tabs_is_enabled();
    NSLog(@"[WailsViewController] viewDidLoad: ios_native_tabs_is_enabled=%d", tabsEnabled);
    if (tabsEnabled) {
        [self enableNativeTabs:YES];
    }
}

- (void)viewDidLayoutSubviews {
    [super viewDidLayoutSubviews];
    // Layout webView and optional tabBar respecting safe area
    UIEdgeInsets safe = UIEdgeInsetsZero;
    if (@available(iOS 11.0, *)) {
        safe = self.view.safeAreaInsets;
    }
    CGFloat width = self.view.bounds.size.width;
    CGFloat height = self.view.bounds.size.height;
    CGFloat tabH = 0;
    if (self.tabBar && !self.tabBar.isHidden) {
        CGSize size = [self.tabBar sizeThatFits:CGSizeMake(width, CGFLOAT_MAX)];
        tabH = size.height;
        self.tabBar.frame = CGRectMake(0, height - safe.bottom - tabH, width, tabH);
    }
    CGFloat webTop = safe.top;
    CGFloat webBottom = safe.bottom + tabH;
    self.webView.frame = UIEdgeInsetsInsetRect(self.view.bounds, UIEdgeInsetsMake(webTop, safe.left, webBottom, safe.right));
}

- (void)enableNativeTabs:(BOOL)enabled {
    dispatch_async(dispatch_get_main_queue(), ^{
        NSLog(@"[WailsViewController] enableNativeTabs called with enabled=%d, existingTabBar=%@", enabled, self.tabBar ? @"YES" : @"NO");
        if (enabled) {
            if (!self.tabBar) {
                UITabBar *tb = [[UITabBar alloc] init];
                tb.delegate = self;
                if (@available(iOS 13.0, *)) {
                    UITabBarAppearance *appearance = [[UITabBarAppearance alloc] init];
                    [appearance configureWithDefaultBackground];
                    tb.standardAppearance = appearance;
                    if (@available(iOS 15.0, *)) {
                        tb.scrollEdgeAppearance = appearance;
                    }
                }
                // Build items from configured JSON, fallback to defaults
                const char* cjson = ios_native_tabs_get_items_json();
                NSMutableArray<UITabBarItem*> *items = nil;
                if (cjson) {
                    NSString *jsonStr = [NSString stringWithUTF8String:cjson];
                    free((void*)cjson);
                    if (jsonStr.length) {
                        NSData *data = [jsonStr dataUsingEncoding:NSUTF8StringEncoding];
                        NSError *err = nil;
                        id obj = [NSJSONSerialization JSONObjectWithData:data options:0 error:&err];
                        if (!err && [obj isKindOfClass:[NSArray class]]) {
                            NSArray *arr = (NSArray*)obj;
                            NSLog(@"[WailsViewController] Building tab items from JSON, count=%lu", (unsigned long)arr.count);
                            items = [NSMutableArray arrayWithCapacity:arr.count];
                            NSInteger tag = 0;
                            for (id entry in arr) {
                                if (![entry isKindOfClass:[NSDictionary class]]) continue;
                                NSDictionary *d = (NSDictionary*)entry;
                                NSString *title = [d[@"Title"] isKindOfClass:[NSString class]] ? d[@"Title"] : @"";
                                UIImage *img = nil;
                                if (@available(iOS 13.0, *)) {
                                    NSString *symbol = [d[@"SystemImage"] isKindOfClass:[NSString class]] ? d[@"SystemImage"] : nil;
                                    if (symbol.length) {
                                        img = [UIImage systemImageNamed:symbol];
                                    }
                                }
                                UITabBarItem *it = [[UITabBarItem alloc] initWithTitle:(title ?: @"") image:img tag:tag++];
                                [items addObject:it];
                            }
                        }
                        else if (err) {
                            NSLog(@"[WailsViewController] ERROR parsing NativeTabsItems JSON: %@", err);
                        }
                    }
                    else {
                        NSLog(@"[WailsViewController] NativeTabsItems JSON string is empty");
                    }
                }
                if (items != nil && items.count > 0) {
                    tb.items = items;
                    tb.selectedItem = items.firstObject;
                    NSLog(@"[WailsViewController] TabBar created with %lu item(s) from config", (unsigned long)items.count);
                } else {
                    // Default 3 items
                    UITabBarItem *item0 = [[UITabBarItem alloc] initWithTitle:@"Bindings" image:nil tag:0];
                    UITabBarItem *item1 = [[UITabBarItem alloc] initWithTitle:@"Go Runtime" image:nil tag:1];
                    UITabBarItem *item2 = [[UITabBarItem alloc] initWithTitle:@"JS Runtime" image:nil tag:2];
                    tb.items = @[item0, item1, item2];
                    tb.selectedItem = item0;
                    NSLog(@"[WailsViewController] TabBar created with default items (3)" );
                }
                self.tabBar = tb;
                [self.view addSubview:self.tabBar];
                NSLog(@"[WailsViewController] TabBar added as subview");
            }
            self.tabBar.hidden = NO;
            NSLog(@"[WailsViewController] TabBar set hidden=NO");
        } else {
            if (self.tabBar) {
                self.tabBar.hidden = YES;
                NSLog(@"[WailsViewController] TabBar set hidden=YES");
            }
        }
        [self.view setNeedsLayout];
        [self.view layoutIfNeeded];
        NSLog(@"[WailsViewController] Requested layout update (enableNativeTabs)");
    });
}

- (void)selectNativeTabIndex:(NSInteger)index {
    dispatch_async(dispatch_get_main_queue(), ^{
        if (!self.tabBar || self.tabBar.isHidden) return;
        if (index < 0 || index >= (NSInteger)self.tabBar.items.count) return;
        UITabBarItem *item = self.tabBar.items[index];
        self.tabBar.selectedItem = item;
        [self tabBar:self.tabBar didSelectItem:item];
    });
}

#pragma mark - UITabBarDelegate

- (void)tabBar:(UITabBar *)tabBar didSelectItem:(UITabBarItem *)item {
    NSInteger idx = [tabBar.items indexOfObject:item];
    if (idx == NSNotFound) return;
    // Dispatch a CustomEvent to the frontend
    NSString *js = [NSString stringWithFormat:@"window.dispatchEvent(new CustomEvent('nativeTabSelected',{detail:{index:%ld}}));", (long)idx];
    [self executeJavaScript:js];
}

- (void)executeJavaScript:(NSString *)js {
    [self.webView evaluateJavaScript:js completionHandler:^(id result, NSError *error) {
        if (error) {
            NSLog(@"[WailsViewController] JS error: %@", error);
        }
    }];
}

@end

// MARK: - C bridges used by Go

unsigned int ios_create_webview(void) {
    __block unsigned int windowID = nextWindowID++;
    if (!appDelegate || !appDelegate.window) {
        return windowID;
    }
    dispatch_async(dispatch_get_main_queue(), ^{
        WailsViewController *vc = [[WailsViewController alloc] initWithWindowID:windowID];
        if (!appDelegate.viewControllers) appDelegate.viewControllers = [NSMutableArray array];
        [appDelegate.viewControllers addObject:vc];
        appDelegate.window.rootViewController = vc;
        [vc loadView];
        [vc viewDidLoad];
    });
    return windowID;
}

void* ios_create_webview_with_id(unsigned int wailsID) {
    __block WailsViewController *viewController = nil;
    if (!appDelegate || !appDelegate.window) {
        return NULL;
    }
    void (^createBlock)(void) = ^{
        viewController = [[WailsViewController alloc] initWithWindowID:wailsID];
        if (!appDelegate.viewControllers) appDelegate.viewControllers = [NSMutableArray array];
        [appDelegate.viewControllers addObject:viewController];
        if (appDelegate.viewControllers.count == 1) {
            appDelegate.window.rootViewController = viewController;
        }
        [viewController loadView];
        [viewController viewDidLoad];
    };
    if ([NSThread isMainThread]) {
        createBlock();
    } else {
        dispatch_sync(dispatch_get_main_queue(), createBlock);
    }
    return (__bridge_retained void*)viewController;
}

void ios_execute_javascript(unsigned int windowID, const char* js) {
    if (!js) return;
    NSString *jsString = [NSString stringWithUTF8String:js];
    dispatch_async(dispatch_get_main_queue(), ^{
        for (WailsViewController *vc in appDelegate.viewControllers) {
            if (vc.windowID == windowID) { [vc executeJavaScript:jsString]; break; }
        }
    });
}

void ios_window_exec_js(void* viewController, const char* js) {
    if (!viewController || !js) return;
    WailsViewController *vc = (__bridge WailsViewController *)viewController;
    NSString *jsString = [NSString stringWithUTF8String:js];
    dispatch_async(dispatch_get_main_queue(), ^{ [vc executeJavaScript:jsString]; });
}

void ios_window_load_url(void* viewController, const char* url) {
    if (!viewController || !url) return;
    WailsViewController *vc = (__bridge WailsViewController *)viewController;
    NSString *urlString = [NSString stringWithUTF8String:url];
    dispatch_async(dispatch_get_main_queue(), ^{
        NSURL *nsurl = [NSURL URLWithString:urlString];
        if (!nsurl) return;
        [vc.webView loadRequest:[NSURLRequest requestWithURL:nsurl]];
    });
}

void ios_window_set_html(void* viewController, const char* html) {
    if (!viewController || !html) return;
    WailsViewController *vc = (__bridge WailsViewController *)viewController;
    NSString *htmlString = [NSString stringWithUTF8String:html];
    dispatch_async(dispatch_get_main_queue(), ^{
        [vc.webView loadHTMLString:htmlString baseURL:[NSURL URLWithString:@"wails://localhost/"]];
    });
}

unsigned int ios_window_get_id(void* viewController) {
    if (!viewController) return 0;
    WailsViewController *vc = (__bridge WailsViewController *)viewController;
    return vc.windowID;
}

void ios_window_release_handle(void* viewController) {
    if (!viewController) return;
    CFRelease(viewController);
}

// Broadcast a console message to all active WKWebViews
void ios_console_log(const char* level, const char* message) {
    if (!message) return;
    NSString *lvl = level ? [NSString stringWithUTF8String:level] : @"log";
    NSString *msg = [NSString stringWithUTF8String:message];
    // Mirror to system log for simctl visibility
    NSLog(@"[ios_console_log][%@] %@", lvl, msg);
    // Robustly encode message to avoid JS string escaping issues
    NSData *data = [msg dataUsingEncoding:NSUTF8StringEncoding];
    NSString *b64 = [data base64EncodedStringWithOptions:0];
    NSString *levelJS = ([lvl length] ? [NSString stringWithFormat:@"'%@'", lvl] : @"'log'");
    NSString *js = [NSString stringWithFormat:
                    @"(function(){try{var b=atob('%@');var bytes=new Uint8Array(b.length);for(var i=0;i<b.length;i++){bytes[i]=b.charCodeAt(i);}var msg=new TextDecoder('utf-8').decode(bytes);console[%@](msg);}catch(e){console.log('wails log bridge error:'+e)}})();",
                    b64, levelJS];
    dispatch_async(dispatch_get_main_queue(), ^{
        // Ensure buffer is initialised
        if (pendingConsoleJS == nil) {
            pendingConsoleJS = [NSMutableArray array];
        }
        NSUInteger count = appDelegate.viewControllers.count;
        if (count == 0) {
            // No webviews yet: buffer
            [pendingConsoleJS addObject:js];
            return;
        }
        // Broadcast to all existing webviews
        for (WailsViewController *vc in appDelegate.viewControllers) {
            [vc.webView evaluateJavaScript:js completionHandler:nil];
        }
        // Helpful native debug: number of inspectors/webviews
        NSLog(@"[ios_console_log] Delivered log to %lu webview(s)", (unsigned long)count);
    });
}

// Set background color (applies to VC view, WKWebView, and app window)
void ios_window_set_background_color(void* viewController, unsigned char r, unsigned char g, unsigned char b, unsigned char a) {
    if (!viewController) return;
    WailsViewController *vc = (__bridge WailsViewController *)viewController;
    CGFloat fr = ((CGFloat)r) / 255.0;
    CGFloat fg = ((CGFloat)g) / 255.0;
    CGFloat fb = ((CGFloat)b) / 255.0;
    CGFloat fa = ((CGFloat)a) / 255.0;
    UIColor *color = [UIColor colorWithRed:fr green:fg blue:fb alpha:fa];
    dispatch_async(dispatch_get_main_queue(), ^{
        vc.view.backgroundColor = color;
        if (vc.webView) {
            vc.webView.opaque = (a == 255);
            vc.webView.backgroundColor = color;
            vc.webView.scrollView.backgroundColor = color;
        }
        if (appDelegate && appDelegate.window) {
            appDelegate.window.backgroundColor = color;
        }
    });
}