//go:build ios

#import <UIKit/UIKit.h>
#import <WebKit/WebKit.h>
#import <UniformTypeIdentifiers/UniformTypeIdentifiers.h>
#import "application_ios.h"
#import "application_ios_delegate.h"
#import "webview_window_ios.h"
#import <sys/utsname.h>
#import <stdlib.h>
#import <string.h>
#import <os/log.h>

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

// Framework diagnostic logging (off unless enabled from Go in debug builds)
static bool g_verboseLogging = false;
void ios_set_verbose_logging(bool enabled) { g_verboseLogging = enabled; }
bool ios_is_verbose_logging(void) { return g_verboseLogging; }

// Note: The WailsAppDelegate implementation resides in application_ios_delegate.m

// C interface implementation
void ios_app_init(void) {
    // This will be called from Go's init
    // Explicitly reference WailsAppDelegate to ensure the class is linked and registered.
    (void)[WailsAppDelegate class];
    // The actual UI startup happens via UIApplicationMain in main.m
}

void ios_app_run(void) {
    // No-op: UIApplicationMain is invoked from main.m (the C entry point) on the
    // real OS main thread. The Go runtime is started by the app delegate's
    // didFinishLaunchingWithOptions, i.e. only after UIKit has launched.
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
    forEachViewController(^(WailsViewController *vc){
        [vc enableNativeTabs:(enabled ? YES : NO)];
    });
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
        return;
    }
    g_nativeTabsItemsJSON = [NSString stringWithUTF8String:json];
    // Apply to existing controllers if visible
    forEachViewController(^(WailsViewController *vc){
        if (vc.tabBar && !vc.tabBar.isHidden) {
            // Re-enable to rebuild items from JSON
            [vc enableNativeTabs:YES];
        }
    });
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
        WailsVLog(@"[ios_haptics_impact] requested style=%@", style);
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
                WailsVLog(@"[ios_haptics_impact] Simulator detected: no physical haptic feedback will be felt.");
            #endif
        } else {
            WailsVLog(@"[ios_haptics_impact] iOS version < 13.0: no haptic API available.");
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

// Run a block on the main thread, synchronously, without deadlocking when
// already on the main thread.
static void runOnMainSync(void (^block)(void)) {
    if ([NSThread isMainThread]) {
        block();
    } else {
        dispatch_sync(dispatch_get_main_queue(), block);
    }
}

// Returns the view controller that should present modal UI (alerts, pickers).
static UIViewController* wailsTopViewController(void) {
    UIViewController *vc = appDelegate.window.rootViewController;
    while (vc.presentedViewController != nil && !vc.presentedViewController.isBeingDismissed) {
        vc = vc.presentedViewController;
    }
    return vc;
}

// MARK: - Screen metrics

const char* ios_screen_info_json(void) {
    __block NSString *json = nil;
    runOnMainSync(^{
        UIScreen *screen = [UIScreen mainScreen];
        CGRect bounds = screen.bounds;             // points, current orientation
        CGRect native = screen.nativeBounds;       // pixels, portrait-up
        CGFloat scale = screen.scale;
        UIEdgeInsets safe = UIEdgeInsetsZero;
        if (appDelegate && appDelegate.window) {
            safe = appDelegate.window.safeAreaInsets;
        }
        json = [NSString stringWithFormat:
            @"{\"pointWidth\":%d,\"pointHeight\":%d,\"pixelWidth\":%d,\"pixelHeight\":%d,\"scale\":%.2f,"
             "\"safeTop\":%d,\"safeBottom\":%d,\"safeLeft\":%d,\"safeRight\":%d}",
            (int)bounds.size.width, (int)bounds.size.height,
            (int)native.size.width, (int)native.size.height,
            (double)scale,
            (int)safe.top, (int)safe.bottom, (int)safe.left, (int)safe.right];
    });
    return dupCString(json);
}

// MARK: - Clipboard

void ios_clipboard_set_text(const char* text) {
    NSString *str = text ? [NSString stringWithUTF8String:text] : @"";
    runOnMainSync(^{
        [UIPasteboard generalPasteboard].string = str;
    });
}

const char* ios_clipboard_get_text(void) {
    __block NSString *str = nil;
    runOnMainSync(^{
        str = [UIPasteboard generalPasteboard].string;
    });
    return dupCString(str);
}

// MARK: - Message dialogs

void ios_show_message_dialog(const char* ctitle, const char* cmessage, const char* cbuttonsJSON, unsigned int callbackID) {
    NSString *title = ctitle ? [NSString stringWithUTF8String:ctitle] : nil;
    NSString *message = cmessage ? [NSString stringWithUTF8String:cmessage] : nil;
    NSString *buttonsJSON = cbuttonsJSON ? [NSString stringWithUTF8String:cbuttonsJSON] : nil;
    dispatch_async(dispatch_get_main_queue(), ^{
        UIAlertController *alert = [UIAlertController alertControllerWithTitle:title
                                                                       message:message
                                                                preferredStyle:UIAlertControllerStyleAlert];
        NSArray *buttons = nil;
        if (buttonsJSON.length) {
            NSData *data = [buttonsJSON dataUsingEncoding:NSUTF8StringEncoding];
            id obj = [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
            if ([obj isKindOfClass:[NSArray class]]) {
                buttons = (NSArray *)obj;
            }
        }
        if (buttons.count == 0) {
            // No buttons configured: show a plain OK
            [alert addAction:[UIAlertAction actionWithTitle:@"OK" style:UIAlertActionStyleDefault
                                                    handler:^(UIAlertAction *action) {
                iosDialogCallback(callbackID, -1);
            }]];
        } else {
            NSInteger idx = 0;
            for (id entry in buttons) {
                if (![entry isKindOfClass:[NSDictionary class]]) { idx++; continue; }
                NSDictionary *d = (NSDictionary *)entry;
                NSString *label = [d[@"label"] isKindOfClass:[NSString class]] ? d[@"label"] : @"";
                BOOL isCancel = [d[@"isCancel"] boolValue];
                BOOL isDefault = [d[@"isDefault"] boolValue];
                NSInteger buttonIndex = idx;
                UIAlertActionStyle style = isCancel ? UIAlertActionStyleCancel : UIAlertActionStyleDefault;
                UIAlertAction *action = [UIAlertAction actionWithTitle:label style:style
                                                               handler:^(UIAlertAction *a) {
                    iosDialogCallback(callbackID, (int)buttonIndex);
                }];
                [alert addAction:action];
                if (isDefault) {
                    alert.preferredAction = action;
                }
                idx++;
            }
        }
        [wailsTopViewController() presentViewController:alert animated:YES completion:nil];
    });
}

// MARK: - Document picker

@interface WailsDocumentPickerDelegate : NSObject <UIDocumentPickerDelegate>
@property (nonatomic, assign) unsigned int callbackID;
@property (nonatomic, assign) BOOL opensInPlace;
@end

// Keep delegates alive until the picker completes
static NSMutableDictionary<NSNumber*, WailsDocumentPickerDelegate*> *g_pickerDelegates;

@implementation WailsDocumentPickerDelegate
- (void)documentPicker:(UIDocumentPickerViewController *)controller didPickDocumentsAtURLs:(NSArray<NSURL *> *)urls {
    for (NSURL *url in urls) {
        // Files are imported with asCopy:YES, so they are sandbox copies that
        // need no security-scoped access. Only directories are opened in place
        // and require it; that access is held for the app session (bookmark
        // persistence is not implemented).
        if (self.opensInPlace) {
            [url startAccessingSecurityScopedResource];
        }
        iosOpenFileCallback(self.callbackID, (char *)[url.path UTF8String]);
    }
    iosOpenFileCallbackEnd(self.callbackID);
    [g_pickerDelegates removeObjectForKey:@(self.callbackID)];
}
- (void)documentPickerWasCancelled:(UIDocumentPickerViewController *)controller {
    iosOpenFileCallbackEnd(self.callbackID);
    [g_pickerDelegates removeObjectForKey:@(self.callbackID)];
}
@end

void ios_show_document_picker(unsigned int callbackID, bool directories, bool multiple) {
    dispatch_async(dispatch_get_main_queue(), ^{
        UIDocumentPickerViewController *picker;
        if (directories) {
            // Folders open in place (asCopy is not supported for folders)
            picker = [[UIDocumentPickerViewController alloc] initForOpeningContentTypes:@[UTTypeFolder]];
        } else {
            // Import file copies into the app sandbox so the app can read
            // them without managing security-scoped access.
            picker = [[UIDocumentPickerViewController alloc] initForOpeningContentTypes:@[UTTypeItem] asCopy:YES];
        }
        picker.allowsMultipleSelection = multiple ? YES : NO;
        WailsDocumentPickerDelegate *delegate = [[WailsDocumentPickerDelegate alloc] init];
        delegate.callbackID = callbackID;
        delegate.opensInPlace = directories ? YES : NO;
        if (!g_pickerDelegates) {
            g_pickerDelegates = [NSMutableDictionary dictionary];
        }
        g_pickerDelegates[@(callbackID)] = delegate;
        picker.delegate = delegate;
        [wailsTopViewController() presentViewController:picker animated:YES completion:nil];
    });
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
