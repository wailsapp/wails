//go:build ios

#import <UIKit/UIKit.h>
#import <AVFoundation/AVFoundation.h>
#import <LocalAuthentication/LocalAuthentication.h>
#import <UserNotifications/UserNotifications.h>
#import <Security/Security.h>
#import "mobile_features_ios.h"
#import "mobile_features_ios_internal.h"
#import "application_ios_delegate.h"

// appDelegate is defined in application_ios.m.
extern WailsAppDelegate *appDelegate;

// Global appearance state, read back by the WailsViewController overrides.
static int g_mfOrientationMode = 0;   // 0=auto, 1=portrait, 2=landscape
static int g_mfStatusBarStyle = 0;    // 0=default, 1=light content, 2=dark content
static BOOL g_mfStatusBarHidden = NO;

// Run a block on the main thread without deadlocking when already on it.
static void mfRunOnMain(void (^block)(void)) {
    if ([NSThread isMainThread]) {
        block();
    } else {
        dispatch_async(dispatch_get_main_queue(), block);
    }
}

// The view controller that should present modal UI (share sheet, pickers).
static UIViewController* mfTopViewController(void) {
    UIViewController *vc = appDelegate.window.rootViewController;
    while (vc.presentedViewController != nil && !vc.presentedViewController.isBeingDismissed) {
        vc = vc.presentedViewController;
    }
    return vc;
}

static NSDictionary* mfParseJSON(const char* json) {
    if (json == NULL) return @{};
    NSString *str = [NSString stringWithUTF8String:json];
    NSData *data = [str dataUsingEncoding:NSUTF8StringEncoding];
    if (data == nil) return @{};
    id obj = [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
    return [obj isKindOfClass:[NSDictionary class]] ? (NSDictionary *)obj : @{};
}

// MARK: - Share sheet

void ios_share(const char* json) {
    NSDictionary *opts = mfParseJSON(json);
    mfRunOnMain(^{
        NSMutableArray *items = [NSMutableArray array];
        NSString *text = opts[@"text"];
        if ([text isKindOfClass:[NSString class]] && text.length) {
            [items addObject:text];
        }
        NSString *urlStr = opts[@"url"];
        if ([urlStr isKindOfClass:[NSString class]] && urlStr.length) {
            NSURL *url = [NSURL URLWithString:urlStr];
            if (url) [items addObject:url];
        }
        if (items.count == 0) return;
        UIActivityViewController *avc =
            [[UIActivityViewController alloc] initWithActivityItems:items applicationActivities:nil];
        UIViewController *top = mfTopViewController();
        // iPad requires an anchor for the popover; centre it on the view.
        if (avc.popoverPresentationController) {
            avc.popoverPresentationController.sourceView = top.view;
            avc.popoverPresentationController.sourceRect =
                CGRectMake(CGRectGetMidX(top.view.bounds), CGRectGetMidY(top.view.bounds), 0, 0);
            avc.popoverPresentationController.permittedArrowDirections = 0;
        }
        [top presentViewController:avc animated:YES completion:nil];
    });
}

// MARK: - Open URL externally

void ios_open_url(const char* curl) {
    if (curl == NULL) return;
    NSString *str = [NSString stringWithUTF8String:curl];
    NSURL *url = [NSURL URLWithString:str];
    if (url == nil) return;
    mfRunOnMain(^{
        [[UIApplication sharedApplication] openURL:url options:@{} completionHandler:nil];
    });
}

// MARK: - Keep screen awake

void ios_set_keep_awake(bool enabled) {
    mfRunOnMain(^{
        [UIApplication sharedApplication].idleTimerDisabled = enabled ? YES : NO;
    });
}

// MARK: - Torch / flashlight

void ios_set_torch(bool enabled) {
    mfRunOnMain(^{
        AVCaptureDevice *device = [AVCaptureDevice defaultDeviceWithMediaType:AVMediaTypeVideo];
        if (device == nil || !device.hasTorch || !device.isTorchAvailable) {
            extern void iosEmitNativeEvent(const char* name, const char* json);
            iosEmitNativeEvent("native:torch", "{\"on\":false,\"available\":false}");
            return;
        }
        NSError *error = nil;
        if ([device lockForConfiguration:&error]) {
            if (enabled) {
                [device setTorchModeOnWithLevel:AVCaptureMaxAvailableTorchLevel error:nil];
            } else {
                device.torchMode = AVCaptureTorchModeOff;
            }
            [device unlockForConfiguration];
            extern void iosEmitNativeEvent(const char* name, const char* json);
            iosEmitNativeEvent("native:torch", enabled ? "{\"on\":true,\"available\":true}"
                                                       : "{\"on\":false,\"available\":true}");
        }
    });
}

// dupCString clones an NSString into a malloc'd C string (caller frees).
static const char* mfDup(NSString *str) {
    if (str == nil) return NULL;
    const char* utf8 = [str UTF8String];
    if (utf8 == NULL) return NULL;
    size_t len = strlen(utf8) + 1;
    char* out = (char*)malloc(len);
    if (out) memcpy(out, utf8, len);
    return out;
}

static void mfRunOnMainSync(void (^block)(void)) {
    if ([NSThread isMainThread]) {
        block();
    } else {
        dispatch_sync(dispatch_get_main_queue(), block);
    }
}

// MARK: - Safe-area insets

const char* ios_safe_area_json(void) {
    __block NSString *json = nil;
    mfRunOnMainSync(^{
        UIEdgeInsets safe = UIEdgeInsetsZero;
        if (@available(iOS 11.0, *)) {
            UIView *root = appDelegate.window.rootViewController.view;
            safe = root ? root.safeAreaInsets : appDelegate.window.safeAreaInsets;
        }
        json = [NSString stringWithFormat:
            @"{\"top\":%d,\"bottom\":%d,\"left\":%d,\"right\":%d}",
            (int)safe.top, (int)safe.bottom, (int)safe.left, (int)safe.right];
    });
    return mfDup(json);
}

// MARK: - Brightness

void ios_set_brightness(double value) {
    if (value < 0) value = 0;
    if (value > 1) value = 1;
    mfRunOnMain(^{
        [UIScreen mainScreen].brightness = value;
    });
}

double ios_get_brightness(void) {
    __block double v = 0;
    mfRunOnMainSync(^{
        v = (double)[UIScreen mainScreen].brightness;
    });
    return v;
}

// MARK: - App info

const char* ios_app_info_json(void) {
    NSDictionary *info = [[NSBundle mainBundle] infoDictionary];
    NSString *name = info[@"CFBundleDisplayName"] ?: info[@"CFBundleName"] ?: @"";
    NSString *version = info[@"CFBundleShortVersionString"] ?: @"";
    NSString *build = info[@"CFBundleVersion"] ?: @"";
    NSString *bundleId = [[NSBundle mainBundle] bundleIdentifier] ?: @"";
    NSString *json = [NSString stringWithFormat:
        @"{\"name\":\"%@\",\"version\":\"%@\",\"build\":\"%@\",\"bundleId\":\"%@\"}",
        name, version, build, bundleId];
    return mfDup(json);
}

// MARK: - Orientation

UIInterfaceOrientationMask mfSupportedOrientations(void) {
    switch (g_mfOrientationMode) {
        case 1: return UIInterfaceOrientationMaskPortrait;
        case 2: return UIInterfaceOrientationMaskLandscape;
        default: return UIInterfaceOrientationMaskAll;
    }
}

void ios_set_orientation(const char* cmode) {
    NSString *mode = cmode ? [NSString stringWithUTF8String:cmode] : @"auto";
    if ([mode isEqualToString:@"portrait"]) g_mfOrientationMode = 1;
    else if ([mode isEqualToString:@"landscape"]) g_mfOrientationMode = 2;
    else g_mfOrientationMode = 0;
    mfRunOnMain(^{
        UIViewController *vc = appDelegate.window.rootViewController;
        if (@available(iOS 16.0, *)) {
            [vc setNeedsUpdateOfSupportedInterfaceOrientations];
            UIWindowScene *scene = (UIWindowScene *)appDelegate.window.windowScene;
            if (scene) {
                UIWindowSceneGeometryPreferencesIOS *prefs =
                    [[UIWindowSceneGeometryPreferencesIOS alloc]
                        initWithInterfaceOrientations:mfSupportedOrientations()];
                [scene requestGeometryUpdateWithPreferences:prefs errorHandler:^(NSError *e){ (void)e; }];
            }
        } else {
            [UIViewController attemptRotationToDeviceOrientation];
        }
    });
}

const char* ios_get_orientation(void) {
    __block NSString *out = @"unknown";
    mfRunOnMainSync(^{
        UIInterfaceOrientation o = UIInterfaceOrientationUnknown;
        if (@available(iOS 13.0, *)) {
            UIWindowScene *scene = (UIWindowScene *)appDelegate.window.windowScene;
            if (scene) o = scene.interfaceOrientation;
        }
        switch (o) {
            case UIInterfaceOrientationPortrait:
            case UIInterfaceOrientationPortraitUpsideDown:
                out = @"portrait"; break;
            case UIInterfaceOrientationLandscapeLeft:
            case UIInterfaceOrientationLandscapeRight:
                out = @"landscape"; break;
            default: out = @"unknown"; break;
        }
    });
    return mfDup(out);
}

// MARK: - Status bar

UIStatusBarStyle mfStatusBarStyle(void) {
    switch (g_mfStatusBarStyle) {
        case 1: return UIStatusBarStyleLightContent;
        case 2: if (@available(iOS 13.0, *)) return UIStatusBarStyleDarkContent; return UIStatusBarStyleDefault;
        default: return UIStatusBarStyleDefault;
    }
}

BOOL mfStatusBarHidden(void) { return g_mfStatusBarHidden; }

void ios_set_status_bar(const char* json) {
    NSDictionary *opts = mfParseJSON(json);
    NSString *style = opts[@"style"];
    if ([style isEqualToString:@"light"]) g_mfStatusBarStyle = 1;
    else if ([style isEqualToString:@"dark"]) g_mfStatusBarStyle = 2;
    else g_mfStatusBarStyle = 0;
    id hidden = opts[@"hidden"];
    if ([hidden isKindOfClass:[NSNumber class]]) g_mfStatusBarHidden = [hidden boolValue];
    mfRunOnMain(^{
        UIViewController *vc = appDelegate.window.rootViewController;
        [vc setNeedsStatusBarAppearanceUpdate];
    });
}

// Emit a custom event with a dictionary payload (JSON-encoded safely).
static void mfEmit(NSString *name, NSDictionary *dict) {
    extern void iosEmitNativeEvent(const char* name, const char* json);
    NSData *data = [NSJSONSerialization dataWithJSONObject:dict options:0 error:nil];
    NSString *s = data ? [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding] : @"{}";
    iosEmitNativeEvent([name UTF8String], [s UTF8String]);
}

// MARK: - Biometric authentication

void ios_biometric_authenticate(const char* creason) {
    NSString *reason = creason ? [NSString stringWithUTF8String:creason] : @"Authenticate";
    LAContext *ctx = [[LAContext alloc] init];
    NSError *err = nil;
    LAPolicy policy = LAPolicyDeviceOwnerAuthenticationWithBiometrics;
    if (![ctx canEvaluatePolicy:policy error:&err]) {
        // Fall back to device passcode if biometrics are unavailable/unenrolled.
        policy = LAPolicyDeviceOwnerAuthentication;
        if (![ctx canEvaluatePolicy:policy error:&err]) {
            mfEmit(@"native:biometric", @{@"ok": @NO,
                @"error": err ? err.localizedDescription : @"biometrics unavailable"});
            return;
        }
    }
    [ctx evaluatePolicy:policy localizedReason:reason reply:^(BOOL success, NSError *error) {
        if (success) {
            mfEmit(@"native:biometric", @{@"ok": @YES});
        } else {
            mfEmit(@"native:biometric", @{@"ok": @NO,
                @"error": error ? error.localizedDescription : @"failed"});
        }
    }];
}

// MARK: - Local notifications

void ios_post_notification(const char* json) {
    NSDictionary *opts = mfParseJSON(json);
    NSString *title = [opts[@"title"] isKindOfClass:[NSString class]] ? opts[@"title"] : @"Notification";
    NSString *body = [opts[@"body"] isKindOfClass:[NSString class]] ? opts[@"body"] : @"";
    double delay = [opts[@"delay"] isKindOfClass:[NSNumber class]] ? [opts[@"delay"] doubleValue] : 2.0;
    if (delay < 1) delay = 1; // UNTimeIntervalNotificationTrigger requires > 0
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center requestAuthorizationWithOptions:
        (UNAuthorizationOptionAlert | UNAuthorizationOptionSound | UNAuthorizationOptionBadge)
        completionHandler:^(BOOL granted, NSError *error) {
            if (!granted) {
                mfEmit(@"native:notification", @{@"ok": @NO,
                    @"error": error ? error.localizedDescription : @"not authorized"});
                return;
            }
            UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];
            content.title = title;
            content.body = body;
            content.sound = [UNNotificationSound defaultSound];
            UNTimeIntervalNotificationTrigger *trigger =
                [UNTimeIntervalNotificationTrigger triggerWithTimeInterval:delay repeats:NO];
            UNNotificationRequest *req = [UNNotificationRequest
                requestWithIdentifier:[[NSUUID UUID] UUIDString] content:content trigger:trigger];
            [center addNotificationRequest:req withCompletionHandler:^(NSError *e) {
                if (e) {
                    mfEmit(@"native:notification", @{@"ok": @NO, @"error": e.localizedDescription});
                } else {
                    mfEmit(@"native:notification", @{@"ok": @YES, @"scheduled": @(delay)});
                }
            }];
        }];
}

// MARK: - Secure storage (Keychain)

static NSMutableDictionary* mfKeychainQuery(NSString *key) {
    NSMutableDictionary *q = [NSMutableDictionary dictionary];
    q[(__bridge id)kSecClass] = (__bridge id)kSecClassGenericPassword;
    q[(__bridge id)kSecAttrService] = @"wails.mobile.securestore";
    q[(__bridge id)kSecAttrAccount] = key ?: @"";
    return q;
}

void ios_secure_set(const char* ckey, const char* cvalue) {
    NSString *key = ckey ? [NSString stringWithUTF8String:ckey] : @"";
    NSString *value = cvalue ? [NSString stringWithUTF8String:cvalue] : @"";
    if (key.length == 0) return;
    NSMutableDictionary *q = mfKeychainQuery(key);
    SecItemDelete((__bridge CFDictionaryRef)q);
    q[(__bridge id)kSecValueData] = [value dataUsingEncoding:NSUTF8StringEncoding];
    q[(__bridge id)kSecAttrAccessible] = (__bridge id)kSecAttrAccessibleWhenUnlockedThisDeviceOnly;
    SecItemAdd((__bridge CFDictionaryRef)q, NULL);
}

const char* ios_secure_get(const char* ckey) {
    NSString *key = ckey ? [NSString stringWithUTF8String:ckey] : @"";
    NSMutableDictionary *q = mfKeychainQuery(key);
    q[(__bridge id)kSecReturnData] = @YES;
    q[(__bridge id)kSecMatchLimit] = (__bridge id)kSecMatchLimitOne;
    CFTypeRef result = NULL;
    OSStatus status = SecItemCopyMatching((__bridge CFDictionaryRef)q, &result);
    if (status == errSecSuccess && result) {
        NSData *data = (__bridge_transfer NSData *)result;
        NSString *s = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
        return mfDup(s ?: @"");
    }
    return mfDup(@"");
}

void ios_secure_delete(const char* ckey) {
    NSString *key = ckey ? [NSString stringWithUTF8String:ckey] : @"";
    SecItemDelete((__bridge CFDictionaryRef)mfKeychainQuery(key));
}
