//go:build ios

#import <UIKit/UIKit.h>
#import <WebKit/WebKit.h>
#import "application_ios.h"
#import "application_ios_delegate.h"
#import "webview_window_ios.h"
#import <sys/utsname.h>
#import <stdlib.h>
#import <string.h>

// Forward declarations for Go callbacks
void ServeAssetRequest(unsigned int windowID, void* urlSchemeTask);
void HandleJSMessage(unsigned int windowID, char* message);


// Global references - declare after interface
WailsAppDelegate *appDelegate = nil;
unsigned int nextWindowID = 1;
static bool g_disableInputAccessory = false; // default: enabled (shown)
// New global flags with sensible iOS defaults
static bool g_disableScroll = false;              // default: scrolling enabled
static bool g_disableBounce = false;              // default: bounce enabled
static bool g_disableScrollIndicators = false;    // default: indicators shown
static bool g_enableBackForwardGestures = false;  // default: gestures disabled
static bool g_disableLinkPreview = false;         // default: link preview enabled
static bool g_enableInlineMediaPlayback = false;  // default: inline playback disabled
static bool g_enableAutoplayNoUserAction = false; // default: autoplay requires user action
static bool g_disableInspectable = false;         // default: inspector enabled
static NSString* g_userAgent = nil;
static NSString* g_appNameForUA = nil;            // default applied in code when nil
static bool g_enableNativeTabs = false;           // default: off
static NSString* g_nativeTabsItemsJSON = nil;     // JSON array of items

// App-wide background colour storage (RGBA 0-255)
static bool g_appBGSet = false;
static unsigned char g_appBG_R = 255;
static unsigned char g_appBG_G = 255;
static unsigned char g_appBG_B = 255;
static unsigned char g_appBG_A = 255;

// Note: The WailsAppDelegate implementation resides in application_ios_delegate.m

// C interface implementation
void ios_app_init(void) {
    // This will be called from Go's init
    // Explicitly reference WailsAppDelegate to ensure the class is linked and registered.
    (void)[WailsAppDelegate class];
    // The actual UI startup happens via UIApplicationMain in main.m
}

void ios_app_run(void) {
    // This function is no longer used - UIApplicationMain is called from main.m
    // The WailsAppDelegate is automatically instantiated by UIApplicationMain
    extern void LogInfo(const char* source, const char* message);
    LogInfo("ios_app_run", "⚠️ This function should not be called - UIApplicationMain runs from main.m");
}

void ios_app_quit(void) {
    dispatch_async(dispatch_get_main_queue(), ^{
        exit(0);
    });
}

bool ios_is_dark_mode(void) {
    if (@available(iOS 13.0, *)) {
        UIUserInterfaceStyle style = [[UITraitCollection currentTraitCollection] userInterfaceStyle];
        return style == UIUserInterfaceStyleDark;
    }
    return false;
}

void ios_set_disable_input_accessory(bool disabled) {
    g_disableInputAccessory = disabled;
}

bool ios_is_input_accessory_disabled(void) {
    return g_disableInputAccessory;
}

// Scrolling & bounce & indicators
void ios_set_disable_scroll(bool disabled) { g_disableScroll = disabled; }
bool ios_is_scroll_disabled(void) { return g_disableScroll; }
void ios_set_disable_bounce(bool disabled) { g_disableBounce = disabled; }
bool ios_is_bounce_disabled(void) { return g_disableBounce; }
void ios_set_disable_scroll_indicators(bool disabled) { g_disableScrollIndicators = disabled; }
bool ios_is_scroll_indicators_disabled(void) { return g_disableScrollIndicators; }

// Navigation gestures
void ios_set_enable_back_forward_gestures(bool enabled) { g_enableBackForwardGestures = enabled; }
bool ios_is_back_forward_gestures_enabled(void) { return g_enableBackForwardGestures; }

// Link previews
void ios_set_disable_link_preview(bool disabled) { g_disableLinkPreview = disabled; }
bool ios_is_link_preview_disabled(void) { return g_disableLinkPreview; }

// Media playback
void ios_set_enable_inline_media_playback(bool enabled) { g_enableInlineMediaPlayback = enabled; }
bool ios_is_inline_media_playback_enabled(void) { return g_enableInlineMediaPlayback; }
void ios_set_enable_autoplay_without_user_action(bool enabled) { g_enableAutoplayNoUserAction = enabled; }
bool ios_is_autoplay_without_user_action_enabled(void) { return g_enableAutoplayNoUserAction; }

// Inspector
void ios_set_disable_inspectable(bool disabled) { g_disableInspectable = disabled; }
bool ios_is_inspectable_disabled(void) { return g_disableInspectable; }

// User agent customization
void ios_set_user_agent(const char* ua) {
    if (ua == NULL) { g_userAgent = nil; return; }
    g_userAgent = [NSString stringWithUTF8String:ua];
}
const char* ios_get_user_agent(void) {
    if (g_userAgent == nil) return NULL;
    return [g_userAgent UTF8String];
}
void ios_set_app_name_for_user_agent(const char* name) {
    if (name == NULL) { g_appNameForUA = nil; return; }
    g_appNameForUA = [NSString stringWithUTF8String:name];
}
const char* ios_get_app_name_for_user_agent(void) {
    if (g_appNameForUA == nil) return NULL;
    return [g_appNameForUA UTF8String];
}

// Live runtime mutations (apply to existing WKWebView instances)
static void forEachViewController(void (^block)(WailsViewController *vc)) {
    if (!appDelegate || !appDelegate.viewControllers) return;
    void (^applyBlock)(void) = ^{
        for (WailsViewController *vc in appDelegate.viewControllers) {
            if (!vc || !vc.webView) continue;
            block(vc);
        }
    };
    if ([NSThread isMainThread]) {
        applyBlock();
    } else {
        dispatch_async(dispatch_get_main_queue(), applyBlock);
    }
}

void ios_runtime_set_scroll_enabled(bool enabled) {
    g_disableScroll = !enabled;
    forEachViewController(^(WailsViewController *vc){
        vc.webView.scrollView.scrollEnabled = enabled ? YES : NO;
    });
}

void ios_runtime_set_bounce_enabled(bool enabled) {
    g_disableBounce = !enabled;
    forEachViewController(^(WailsViewController *vc){
        UIScrollView *sv = vc.webView.scrollView;
        sv.bounces = enabled ? YES : NO;
        sv.alwaysBounceVertical = enabled ? YES : NO;
        sv.alwaysBounceHorizontal = enabled ? YES : NO;
    });
}

void ios_runtime_set_scroll_indicators_enabled(bool enabled) {
    g_disableScrollIndicators = !enabled;
    forEachViewController(^(WailsViewController *vc){
        UIScrollView *sv = vc.webView.scrollView;
        sv.showsVerticalScrollIndicator = enabled ? YES : NO;
        sv.showsHorizontalScrollIndicator = enabled ? YES : NO;
    });
}

void ios_runtime_set_back_forward_gestures_enabled(bool enabled) {
    g_enableBackForwardGestures = enabled;
    forEachViewController(^(WailsViewController *vc){
        vc.webView.allowsBackForwardNavigationGestures = enabled ? YES : NO;
    });
}

void ios_runtime_set_link_preview_enabled(bool enabled) {
    g_disableLinkPreview = !enabled;
    forEachViewController(^(WailsViewController *vc){
        vc.webView.allowsLinkPreview = enabled ? YES : NO;
    });
}

void ios_runtime_set_inspectable_enabled(bool enabled) {
    g_disableInspectable = !enabled;
    forEachViewController(^(WailsViewController *vc){
        BOOL inspectorOn = enabled ? YES : NO;
        if (@available(iOS 16.4, *)) {
            vc.webView.inspectable = inspectorOn;
        } else {
            @try { [vc.webView setValue:@(inspectorOn) forKey:@"inspectable"]; } @catch (__unused NSException *e) {}
        }
    });
}

void ios_runtime_set_custom_user_agent(const char* ua) {
    ios_set_user_agent(ua);
    NSString *uaStr = (ua ? [NSString stringWithUTF8String:ua] : nil);
    forEachViewController(^(WailsViewController *vc){
        vc.webView.customUserAgent = uaStr;
    });
}

// Forward declaration used by getters
static const char* dupCString(NSString *str);

// Native bottom tab bar controls
void ios_native_tabs_set_enabled(bool enabled) {
    g_enableNativeTabs = enabled;
    NSLog(@"[ios_native_tabs_set_enabled] enabled=%d", enabled);
    __block NSUInteger count = 0;
    forEachViewController(^(WailsViewController *vc){
        count++;
        [vc enableNativeTabs:(enabled ? YES : NO)];
    });
    NSLog(@"[ios_native_tabs_set_enabled] applied to %lu view controller(s)", (unsigned long)count);
}

bool ios_native_tabs_is_enabled(void) {
    return g_enableNativeTabs;
}

void ios_native_tabs_select_index(int index) {
    forEachViewController(^(WailsViewController *vc){
        [vc selectNativeTabIndex:(NSInteger)index];
    });
}

void ios_native_tabs_set_items_json(const char* json) {
    if (json == NULL) {
        g_nativeTabsItemsJSON = nil;
        NSLog(@"[ios_native_tabs_set_items_json] JSON cleared (nil)");
        return;
    }
    g_nativeTabsItemsJSON = [NSString stringWithUTF8String:json];
    NSLog(@"[ios_native_tabs_set_items_json] JSON length=%lu", (unsigned long)g_nativeTabsItemsJSON.length);
    // Apply to existing controllers if visible
    __block NSUInteger refreshed = 0;
    forEachViewController(^(WailsViewController *vc){
        if (vc.tabBar && !vc.tabBar.isHidden) {
            // Re-enable to rebuild items from JSON
            [vc enableNativeTabs:YES];
            refreshed++;
        }
    });
    NSLog(@"[ios_native_tabs_set_items_json] refreshed %lu visible tab bar(s)", (unsigned long)refreshed);
}

const char* ios_native_tabs_get_items_json(void) {
    return dupCString(g_nativeTabsItemsJSON);
}

// App-wide background colour control
void ios_set_app_background_color(unsigned char r, unsigned char g, unsigned char b, unsigned char a, bool isSet) {
    g_appBGSet = isSet;
    if (isSet) {
        g_appBG_R = r; g_appBG_G = g; g_appBG_B = b; g_appBG_A = a;
    }
}

bool ios_get_app_background_color(unsigned char* r, unsigned char* g, unsigned char* b, unsigned char* a) {
    if (!g_appBGSet) return false;
    if (r) *r = g_appBG_R;
    if (g) *g = g_appBG_G;
    if (b) *b = g_appBG_B;
    if (a) *a = g_appBG_A;
    return true;
}


// iOS Runtime bridges
void ios_haptics_impact(const char* cstyle) {
    if (cstyle == NULL) return;
    NSString *style = [NSString stringWithUTF8String:cstyle];
    dispatch_async(dispatch_get_main_queue(), ^{
        NSLog(@"[ios_haptics_impact] requested style=%@", style);
        if (@available(iOS 13.0, *)) {
            UIImpactFeedbackStyle feedbackStyle = UIImpactFeedbackStyleMedium;
            if ([style isEqualToString:@"light"]) feedbackStyle = UIImpactFeedbackStyleLight;
            else if ([style isEqualToString:@"medium"]) feedbackStyle = UIImpactFeedbackStyleMedium;
            else if ([style isEqualToString:@"heavy"]) feedbackStyle = UIImpactFeedbackStyleHeavy;
#if __IPHONE_OS_VERSION_MAX_ALLOWED >= 130000
            else if ([style isEqualToString:@"soft"]) feedbackStyle = UIImpactFeedbackStyleSoft;
            else if ([style isEqualToString:@"rigid"]) feedbackStyle = UIImpactFeedbackStyleRigid;
#endif
            UIImpactFeedbackGenerator *generator = [[UIImpactFeedbackGenerator alloc] initWithStyle:feedbackStyle];
            [generator prepare];
            [generator impactOccurred];
            #if TARGET_OS_SIMULATOR
                NSLog(@"[ios_haptics_impact] Simulator detected: no physical haptic feedback will be felt.");
            #else
                NSLog(@"[ios_haptics_impact] Haptic impact triggered.");
            #endif
        } else {
            NSLog(@"[ios_haptics_impact] iOS version < 13.0: no haptic API available.");
        }
    });
}

static const char* dupCString(NSString *str) {
    if (str == nil) return NULL;
    const char* utf8 = [str UTF8String];
    if (utf8 == NULL) return NULL;
    size_t len = strlen(utf8) + 1;
    char* out = (char*)malloc(len);
    if (out) memcpy(out, utf8, len);
    return out;
}

const char* ios_device_info_json(void) {
    UIDevice *device = [UIDevice currentDevice];
    struct utsname systemInfo;
    uname(&systemInfo);
    NSString *model = [NSString stringWithUTF8String:systemInfo.machine];
    NSString *systemName = device.systemName ?: @"iOS";
    NSString *systemVersion = device.systemVersion ?: @"";
#if TARGET_OS_SIMULATOR
    BOOL isSimulator = YES;
#else
    BOOL isSimulator = NO;
#endif
    NSString *json = [NSString stringWithFormat:
        @"{\"model\":\"%@\",\"systemName\":\"%@\",\"systemVersion\":\"%@\",\"isSimulator\":%@}",
        model, systemName, systemVersion, isSimulator ? @"true" : @"false"
    ];
    return dupCString(json);
}
