//go:build ios

#import <UIKit/UIKit.h>
#import <AVFoundation/AVFoundation.h>
#import <LocalAuthentication/LocalAuthentication.h>
#import <UserNotifications/UserNotifications.h>
#import <Security/Security.h>
#import <CoreLocation/CoreLocation.h>
#import <CoreMotion/CoreMotion.h>
#import <SystemConfiguration/SystemConfiguration.h>
#import <netinet/in.h>
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

// Battery: the torch, accelerometer and proximity sensor are paused when the app
// leaves the foreground (the torch in particular is hardware state that would
// otherwise keep draining while the app is suspended). g_mf*Wanted records the
// user's intent so the accelerometer/proximity can be resumed on return.
static BOOL g_mfTorchOn = NO;
static BOOL g_mfMotionWanted = NO;
static BOOL g_mfProximityWanted = NO;
static void mfEnsureLifecycleObservers(void); // defined at end of file

// mfApplyTorch toggles the hardware torch and records the state. Returns NO if no
// usable torch device exists. Emits nothing — the caller reports as needed.
static BOOL mfApplyTorch(BOOL on) {
    AVCaptureDevice *device = [AVCaptureDevice defaultDeviceWithMediaType:AVMediaTypeVideo];
    if (device == nil || !device.hasTorch || !device.isTorchAvailable) return NO;
    NSError *error = nil;
    if (![device lockForConfiguration:&error]) return NO;
    if (on) {
        [device setTorchModeOnWithLevel:AVCaptureMaxAvailableTorchLevel error:nil];
    } else {
        device.torchMode = AVCaptureTorchModeOff;
    }
    [device unlockForConfiguration];
    g_mfTorchOn = on;
    return YES;
}

void ios_set_torch(bool enabled) {
    mfRunOnMain(^{
        mfEnsureLifecycleObservers();
        extern void iosEmitNativeEvent(const char* name, const char* json);
        if (!mfApplyTorch(enabled ? YES : NO)) {
            iosEmitNativeEvent("common:torch", "{\"on\":false,\"available\":false}");
            return;
        }
        iosEmitNativeEvent("common:torch", enabled ? "{\"on\":true,\"available\":true}"
                                                   : "{\"on\":false,\"available\":true}");
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
            mfEmit(@"common:biometric", @{@"ok": @NO,
                @"error": err ? err.localizedDescription : @"biometrics unavailable"});
            return;
        }
    }
    [ctx evaluatePolicy:policy localizedReason:reason reply:^(BOOL success, NSError *error) {
        if (success) {
            mfEmit(@"common:biometric", @{@"ok": @YES});
        } else {
            mfEmit(@"common:biometric", @{@"ok": @NO,
                @"error": error ? error.localizedDescription : @"failed"});
        }
    }];
}

// MARK: - Local notifications

// By default iOS does NOT display a banner for a local notification while the
// app that posted it is in the FOREGROUND — it is delivered silently to
// Notification Center instead. To show it (and to react to taps) the app must
// register a UNUserNotificationCenterDelegate. Apple requires the delegate be
// set before the app finishes launching, so ios_notifications_init() is called
// from the WailsAppDelegate's didFinishLaunchingWithOptions.
@interface MFNotificationDelegate : NSObject <UNUserNotificationCenterDelegate>
@end

@implementation MFNotificationDelegate
// Called when a notification fires while the app is in the foreground. Without
// this returning presentation options, nothing is shown on screen.
- (void)userNotificationCenter:(UNUserNotificationCenter *)center
       willPresentNotification:(UNNotification *)notification
         withCompletionHandler:(void (^)(UNNotificationPresentationOptions))completionHandler {
    // Diagnostic: lets the frontend confirm this delegate is actually running on
    // the device (i.e. the build includes the foreground-presentation fix).
    mfEmit(@"common:notification", @{@"ok": @YES, @"presented": @YES});
    if (@available(iOS 14.0, *)) {
        completionHandler(UNNotificationPresentationOptionBanner
                          | UNNotificationPresentationOptionList
                          | UNNotificationPresentationOptionSound);
    } else {
        completionHandler(UNNotificationPresentationOptionAlert
                          | UNNotificationPresentationOptionSound);
    }
}

// Called when the user taps a delivered notification.
- (void)userNotificationCenter:(UNUserNotificationCenter *)center
didReceiveNotificationResponse:(UNNotificationResponse *)response
         withCompletionHandler:(void (^)(void))completionHandler {
    mfEmit(@"common:notification", @{@"ok": @YES, @"tapped": @YES});
    completionHandler();
}
@end

static MFNotificationDelegate *g_mfNotificationDelegate = nil;

void ios_notifications_init(void) {
    mfRunOnMain(^{
        if (g_mfNotificationDelegate == nil) {
            g_mfNotificationDelegate = [[MFNotificationDelegate alloc] init];
        }
        [UNUserNotificationCenter currentNotificationCenter].delegate = g_mfNotificationDelegate;
    });
}

void ios_post_notification(const char* json) {
    // Defensive: make sure the foreground-presentation delegate is registered
    // before we schedule, even if ios_notifications_init() wasn't called at
    // launch (idempotent — just (re)assigns the singleton delegate).
    ios_notifications_init();
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
                mfEmit(@"common:notification", @{@"ok": @NO,
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
                    mfEmit(@"common:notification", @{@"ok": @NO, @"error": e.localizedDescription});
                } else {
                    mfEmit(@"common:notification", @{@"ok": @YES, @"scheduled": @(delay)});
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

// MARK: - Haptics

void ios_haptic(const char* ctype) {
    NSString *type = ctype ? [NSString stringWithUTF8String:ctype] : @"impact-medium";
    mfRunOnMain(^{
        if ([type isEqualToString:@"selection"]) {
            UISelectionFeedbackGenerator *g = [[UISelectionFeedbackGenerator alloc] init];
            [g prepare];
            [g selectionChanged];
        } else if ([type isEqualToString:@"success"] || [type isEqualToString:@"warning"]
                   || [type isEqualToString:@"error"]) {
            UINotificationFeedbackType t = UINotificationFeedbackTypeSuccess;
            if ([type isEqualToString:@"warning"]) t = UINotificationFeedbackTypeWarning;
            else if ([type isEqualToString:@"error"]) t = UINotificationFeedbackTypeError;
            UINotificationFeedbackGenerator *g = [[UINotificationFeedbackGenerator alloc] init];
            [g prepare];
            [g notificationOccurred:t];
        } else {
            UIImpactFeedbackStyle s = UIImpactFeedbackStyleMedium;
            if ([type isEqualToString:@"impact-light"]) s = UIImpactFeedbackStyleLight;
            else if ([type isEqualToString:@"impact-heavy"]) s = UIImpactFeedbackStyleHeavy;
            UIImpactFeedbackGenerator *g = [[UIImpactFeedbackGenerator alloc] initWithStyle:s];
            [g prepare];
            [g impactOccurred];
        }
    });
}

// MARK: - Geolocation

@interface MFLocationDelegate : NSObject <CLLocationManagerDelegate>
@end
@implementation MFLocationDelegate
- (void)locationManager:(CLLocationManager *)manager didUpdateLocations:(NSArray<CLLocation *> *)locations {
    CLLocation *loc = locations.lastObject;
    if (!loc) return;
    mfEmit(@"common:location", @{
        @"lat": @(loc.coordinate.latitude),
        @"lng": @(loc.coordinate.longitude),
        @"accuracy": @(loc.horizontalAccuracy),
    });
}
- (void)locationManager:(CLLocationManager *)manager didFailWithError:(NSError *)error {
    mfEmit(@"common:location", @{@"error": error ? error.localizedDescription : @"location failed"});
}
- (void)locationManagerDidChangeAuthorization:(CLLocationManager *)manager {
    if (@available(iOS 14.0, *)) {
        CLAuthorizationStatus s = manager.authorizationStatus;
        if (s == kCLAuthorizationStatusAuthorizedWhenInUse || s == kCLAuthorizationStatusAuthorizedAlways) {
            [manager requestLocation];
        } else if (s == kCLAuthorizationStatusDenied || s == kCLAuthorizationStatusRestricted) {
            mfEmit(@"common:location", @{@"error": @"location permission denied"});
        }
    }
}
@end

static CLLocationManager *g_mfLocationManager = nil;
static MFLocationDelegate *g_mfLocationDelegate = nil;

void ios_get_location(void) {
    mfRunOnMain(^{
        if (g_mfLocationManager == nil) {
            g_mfLocationDelegate = [[MFLocationDelegate alloc] init];
            g_mfLocationManager = [[CLLocationManager alloc] init];
            g_mfLocationManager.delegate = g_mfLocationDelegate;
            g_mfLocationManager.desiredAccuracy = kCLLocationAccuracyHundredMeters;
        }
        CLAuthorizationStatus status;
        if (@available(iOS 14.0, *)) {
            status = g_mfLocationManager.authorizationStatus;
        } else {
            status = [CLLocationManager authorizationStatus];
        }
        if (status == kCLAuthorizationStatusNotDetermined) {
            // The authorization callback issues requestLocation once granted.
            [g_mfLocationManager requestWhenInUseAuthorization];
        } else if (status == kCLAuthorizationStatusAuthorizedWhenInUse
                   || status == kCLAuthorizationStatusAuthorizedAlways) {
            [g_mfLocationManager requestLocation];
        } else {
            mfEmit(@"common:location", @{@"error": @"location permission denied"});
        }
    });
}

// MARK: - Accelerometer / device motion

static CMMotionManager *g_mfMotionManager = nil;

void ios_set_motion(bool enabled) {
    mfRunOnMain(^{
        mfEnsureLifecycleObservers();
        g_mfMotionWanted = enabled ? YES : NO;
        if (g_mfMotionManager == nil) {
            g_mfMotionManager = [[CMMotionManager alloc] init];
            g_mfMotionManager.accelerometerUpdateInterval = 0.1; // ~10 Hz
        }
        if (enabled) {
            if (!g_mfMotionManager.isAccelerometerAvailable) {
                mfEmit(@"common:motion", @{@"available": @NO});
                return;
            }
            [g_mfMotionManager startAccelerometerUpdatesToQueue:[NSOperationQueue mainQueue]
                withHandler:^(CMAccelerometerData *data, NSError *error) {
                    if (!data) return;
                    mfEmit(@"common:motion", @{
                        @"x": @(data.acceleration.x),
                        @"y": @(data.acceleration.y),
                        @"z": @(data.acceleration.z),
                    });
                }];
        } else {
            [g_mfMotionManager stopAccelerometerUpdates];
        }
    });
}

// MARK: - Proximity sensor

static id g_mfProximityObserver = nil;

void ios_set_proximity(bool enabled) {
    mfRunOnMain(^{
        mfEnsureLifecycleObservers();
        g_mfProximityWanted = enabled ? YES : NO;
        UIDevice *device = [UIDevice currentDevice];
        if (enabled) {
            device.proximityMonitoringEnabled = YES;
            if (!device.proximityMonitoringEnabled) {
                mfEmit(@"common:proximity", @{@"available": @NO});
                return;
            }
            if (g_mfProximityObserver == nil) {
                g_mfProximityObserver = [[NSNotificationCenter defaultCenter]
                    addObserverForName:UIDeviceProximityStateDidChangeNotification
                    object:nil queue:[NSOperationQueue mainQueue]
                    usingBlock:^(NSNotification *note) {
                        mfEmit(@"common:proximity", @{@"near": @([UIDevice currentDevice].proximityState)});
                    }];
            }
        } else {
            if (g_mfProximityObserver) {
                [[NSNotificationCenter defaultCenter] removeObserver:g_mfProximityObserver];
                g_mfProximityObserver = nil;
            }
            device.proximityMonitoringEnabled = NO;
        }
    });
}

// MARK: - Text-to-speech

static AVSpeechSynthesizer *g_mfSpeech = nil;

void ios_speak(const char* ctext) {
    NSString *text = ctext ? [NSString stringWithUTF8String:ctext] : @"";
    if (text.length == 0) return;
    mfRunOnMain(^{
        if (g_mfSpeech == nil) g_mfSpeech = [[AVSpeechSynthesizer alloc] init];
        AVSpeechUtterance *u = [AVSpeechUtterance speechUtteranceWithString:text];
        u.voice = [AVSpeechSynthesisVoice voiceWithLanguage:@"en-US"];
        [g_mfSpeech speakUtterance:u];
    });
}

void ios_stop_speak(void) {
    mfRunOnMain(^{
        [g_mfSpeech stopSpeakingAtBoundary:AVSpeechBoundaryImmediate];
    });
}

// MARK: - Storage info

const char* ios_storage_json(void) {
    NSError *err = nil;
    NSDictionary *attrs = [[NSFileManager defaultManager]
        attributesOfFileSystemForPath:NSHomeDirectory() error:&err];
    long long freeBytes = [attrs[NSFileSystemFreeSize] longLongValue];
    long long totalBytes = [attrs[NSFileSystemSize] longLongValue];
    NSString *json = [NSString stringWithFormat:@"{\"free\":%lld,\"total\":%lld}", freeBytes, totalBytes];
    return mfDup(json);
}

const char* ios_storage_path(void) {
    NSArray<NSString *> *paths = NSSearchPathForDirectoriesInDomains(
        NSApplicationSupportDirectory, NSUserDomainMask, YES);
    NSString *dir = paths.firstObject;
    if (dir == nil) {
        return mfDup(@"");
    }
    // Unlike Android's getFilesDir(), the Application Support directory is not
    // created automatically. Ensure it exists so callers can open databases
    // there immediately. If creation fails, return "" rather than a path that
    // isn't there (createDirectoryAtPath succeeds when the directory already
    // exists, so this is not an error on subsequent calls).
    NSError *err = nil;
    BOOL ok = [[NSFileManager defaultManager] createDirectoryAtPath:dir
                                       withIntermediateDirectories:YES
                                                        attributes:nil
                                                             error:&err];
    if (!ok) {
        return mfDup(@"");
    }
    return mfDup(dir);
}

// MARK: - Power / battery state

const char* ios_power_json(void) {
    __block NSString *json = nil;
    mfRunOnMainSync(^{
        UIDevice *device = [UIDevice currentDevice];
        device.batteryMonitoringEnabled = YES;
        float level = device.batteryLevel; // -1 when unknown (e.g. simulator)
        UIDeviceBatteryState state = device.batteryState;
        BOOL charging = (state == UIDeviceBatteryStateCharging || state == UIDeviceBatteryStateFull);
        BOOL lowPower = [NSProcessInfo processInfo].isLowPowerModeEnabled;
        json = [NSString stringWithFormat:
            @"{\"level\":%.2f,\"charging\":%@,\"lowPower\":%@}",
            level, charging ? @"true" : @"false", lowPower ? @"true" : @"false"];
    });
    return mfDup(json);
}

// MARK: - Network reachability

const char* ios_network_json(void) {
    struct sockaddr_in zeroAddress;
    bzero(&zeroAddress, sizeof(zeroAddress));
    zeroAddress.sin_len = sizeof(zeroAddress);
    zeroAddress.sin_family = AF_INET;
    SCNetworkReachabilityRef ref =
        SCNetworkReachabilityCreateWithAddress(NULL, (const struct sockaddr *)&zeroAddress);
    SCNetworkReachabilityFlags flags = 0;
    BOOL ok = ref && SCNetworkReachabilityGetFlags(ref, &flags);
    if (ref) CFRelease(ref);
    BOOL reachable = ok && (flags & kSCNetworkReachabilityFlagsReachable)
        && !(flags & kSCNetworkReachabilityFlagsConnectionRequired);
    NSString *type = @"none";
    if (reachable) {
        type = (flags & kSCNetworkReachabilityFlagsIsWWAN) ? @"cellular" : @"wifi";
    }
    NSString *json = [NSString stringWithFormat:
        @"{\"connected\":%@,\"type\":\"%@\"}", reachable ? @"true" : @"false", type];
    return mfDup(json);
}

// MARK: - Keyboard insets

static id g_mfKeyboardShowObs = nil;
static id g_mfKeyboardHideObs = nil;

void ios_set_keyboard_watch(bool enabled) {
    mfRunOnMain(^{
        NSNotificationCenter *nc = [NSNotificationCenter defaultCenter];
        if (enabled) {
            if (g_mfKeyboardShowObs == nil) {
                g_mfKeyboardShowObs = [nc addObserverForName:UIKeyboardWillChangeFrameNotification
                    object:nil queue:[NSOperationQueue mainQueue]
                    usingBlock:^(NSNotification *note) {
                        CGRect frame = [note.userInfo[UIKeyboardFrameEndUserInfoKey] CGRectValue];
                        CGFloat screenH = [UIScreen mainScreen].bounds.size.height;
                        CGFloat height = MAX(0, screenH - frame.origin.y);
                        mfEmit(@"common:keyboard", @{@"visible": @(height > 0), @"height": @((int)height)});
                    }];
            }
            if (g_mfKeyboardHideObs == nil) {
                g_mfKeyboardHideObs = [nc addObserverForName:UIKeyboardWillHideNotification
                    object:nil queue:[NSOperationQueue mainQueue]
                    usingBlock:^(NSNotification *note) {
                        mfEmit(@"common:keyboard", @{@"visible": @NO, @"height": @0});
                    }];
            }
        } else {
            if (g_mfKeyboardShowObs) { [nc removeObserver:g_mfKeyboardShowObs]; g_mfKeyboardShowObs = nil; }
            if (g_mfKeyboardHideObs) { [nc removeObserver:g_mfKeyboardHideObs]; g_mfKeyboardHideObs = nil; }
        }
    });
}

// MARK: - Screen-capture detection
//
// iOS provides no public API to BLOCK screenshots; instead it notifies after a
// screenshot is taken and exposes UIScreen.isCaptured for mirroring/recording.
// We report both so the frontend can react (e.g. blur sensitive content).

static id g_mfScreenshotObs = nil;
static id g_mfCapturedObs = nil;

static void mfEmitCaptured(void) {
    BOOL captured = NO;
    if (@available(iOS 11.0, *)) captured = [UIScreen mainScreen].isCaptured;
    mfEmit(@"common:screenCapture", @{@"recording": @(captured), @"protected": @NO});
}

void ios_set_screen_protect(bool enabled) {
    mfRunOnMain(^{
        NSNotificationCenter *nc = [NSNotificationCenter defaultCenter];
        if (enabled) {
            if (g_mfScreenshotObs == nil) {
                g_mfScreenshotObs = [nc addObserverForName:UIApplicationUserDidTakeScreenshotNotification
                    object:nil queue:[NSOperationQueue mainQueue]
                    usingBlock:^(NSNotification *note) {
                        mfEmit(@"common:screenCapture", @{@"screenshot": @YES, @"protected": @NO});
                    }];
            }
            if (@available(iOS 11.0, *)) {
                if (g_mfCapturedObs == nil) {
                    g_mfCapturedObs = [nc addObserverForName:UIScreenCapturedDidChangeNotification
                        object:nil queue:[NSOperationQueue mainQueue]
                        usingBlock:^(NSNotification *note) { mfEmitCaptured(); }];
                }
                mfEmitCaptured(); // report the current state immediately
            }
        } else {
            if (g_mfScreenshotObs) { [nc removeObserver:g_mfScreenshotObs]; g_mfScreenshotObs = nil; }
            if (g_mfCapturedObs) { [nc removeObserver:g_mfCapturedObs]; g_mfCapturedObs = nil; }
        }
    });
}

// MARK: - Foreground/background lifecycle (battery)
//
// When the app leaves the foreground we stop the accelerometer, disable the
// proximity sensor and switch the torch off so none of them keep draining the
// battery while the app is backgrounded/suspended. On return we restart only the
// sensors the user had switched on (the torch is left off — re-enabling hardware
// light without a user action would be surprising).

static void mfHandleEnterBackground(void) {
    if (g_mfMotionManager && g_mfMotionManager.isAccelerometerActive) {
        [g_mfMotionManager stopAccelerometerUpdates];
    }
    if ([UIDevice currentDevice].proximityMonitoringEnabled) {
        [UIDevice currentDevice].proximityMonitoringEnabled = NO;
    }
    if (g_mfTorchOn && mfApplyTorch(NO)) {
        extern void iosEmitNativeEvent(const char* name, const char* json);
        iosEmitNativeEvent("common:torch", "{\"on\":false,\"available\":true}");
    }
}

static void mfHandleEnterForeground(void) {
    if (g_mfMotionWanted && g_mfMotionManager
        && g_mfMotionManager.isAccelerometerAvailable
        && !g_mfMotionManager.isAccelerometerActive) {
        [g_mfMotionManager startAccelerometerUpdatesToQueue:[NSOperationQueue mainQueue]
            withHandler:^(CMAccelerometerData *data, NSError *error) {
                if (!data) return;
                mfEmit(@"common:motion", @{
                    @"x": @(data.acceleration.x),
                    @"y": @(data.acceleration.y),
                    @"z": @(data.acceleration.z),
                });
            }];
    }
    if (g_mfProximityWanted) {
        [UIDevice currentDevice].proximityMonitoringEnabled = YES;
    }
}

static void mfEnsureLifecycleObservers(void) {
    static BOOL registered = NO;
    if (registered) return; // only ever called on the main thread
    registered = YES;
    NSNotificationCenter *nc = [NSNotificationCenter defaultCenter];
    [nc addObserverForName:UIApplicationDidEnterBackgroundNotification object:nil
        queue:[NSOperationQueue mainQueue] usingBlock:^(NSNotification *n){ mfHandleEnterBackground(); }];
    [nc addObserverForName:UIApplicationWillEnterForegroundNotification object:nil
        queue:[NSOperationQueue mainQueue] usingBlock:^(NSNotification *n){ mfHandleEnterForeground(); }];
}

// MARK: - Camera capture (Phase E)

// Downscale a UIImage into a base64 JPEG data URL for display in the webview.
static NSString* mfImageThumbnailDataURL(UIImage *image) {
    if (!image) return nil;
    CGFloat maxDim = 640;
    CGFloat scale = MIN((CGFloat)1.0, maxDim / MAX(image.size.width, image.size.height));
    CGSize sz = CGSizeMake(image.size.width * scale, image.size.height * scale);
    UIGraphicsBeginImageContextWithOptions(sz, YES, 1.0);
    [image drawInRect:CGRectMake(0, 0, sz.width, sz.height)];
    UIImage *small = UIGraphicsGetImageFromCurrentImageContext();
    UIGraphicsEndImageContext();
    NSData *jpeg = UIImageJPEGRepresentation(small, 0.7);
    if (!jpeg) return nil;
    return [@"data:image/jpeg;base64," stringByAppendingString:[jpeg base64EncodedStringWithOptions:0]];
}

@interface MFCameraDelegate : NSObject <UIImagePickerControllerDelegate, UINavigationControllerDelegate>
@end

static MFCameraDelegate *g_mfCameraDelegate = nil; // retained while the picker is presented

@implementation MFCameraDelegate
- (void)imagePickerController:(UIImagePickerController *)picker
        didFinishPickingMediaWithInfo:(NSDictionary<UIImagePickerControllerInfoKey, id> *)info {
    [picker dismissViewControllerAnimated:YES completion:nil];
    NSString *mediaType = info[UIImagePickerControllerMediaType];
    NSFileManager *fm = [NSFileManager defaultManager];
    if ([mediaType isEqualToString:@"public.movie"]) {
        NSURL *src = info[UIImagePickerControllerMediaURL];
        NSString *dst = [NSTemporaryDirectory() stringByAppendingPathComponent:
            [NSString stringWithFormat:@"capture_%@.mov", [[NSUUID UUID] UUIDString]]];
        [fm removeItemAtPath:dst error:nil];
        [fm copyItemAtURL:src toURL:[NSURL fileURLWithPath:dst] error:nil];
        unsigned long long size = [[fm attributesOfItemAtPath:dst error:nil] fileSize];
        // Stream the clip from the temp dir via the wails:// scheme handler (Range
        // support, any length) rather than inlining it as a data URL.
        NSString *streamUrl = [@"/__capture__/" stringByAppendingString:[dst lastPathComponent]];
        mfEmit(@"common:capture", @{@"type": @"video", @"path": dst,
                                    @"size": @(size), @"streamUrl": streamUrl});
    } else {
        UIImage *image = info[UIImagePickerControllerOriginalImage];
        NSData *jpeg = UIImageJPEGRepresentation(image, 0.9);
        NSString *dst = [NSTemporaryDirectory() stringByAppendingPathComponent:
            [NSString stringWithFormat:@"capture_%@.jpg", [[NSUUID UUID] UUIDString]]];
        [jpeg writeToFile:dst atomically:YES];
        NSMutableDictionary *payload = [@{@"type": @"photo", @"path": dst,
                                          @"size": @(jpeg.length)} mutableCopy];
        NSString *thumb = mfImageThumbnailDataURL(image);
        if (thumb) payload[@"thumb"] = thumb;
        mfEmit(@"common:capture", payload);
    }
    g_mfCameraDelegate = nil;
}
- (void)imagePickerControllerDidCancel:(UIImagePickerController *)picker {
    [picker dismissViewControllerAnimated:YES completion:nil];
    mfEmit(@"common:capture", @{@"cancelled": @YES});
    g_mfCameraDelegate = nil;
}
@end

static void mfPresentCamera(BOOL video) {
    mfRunOnMain(^{
        if (![UIImagePickerController isSourceTypeAvailable:UIImagePickerControllerSourceTypeCamera]) {
            mfEmit(@"common:capture", @{@"error": @"camera not available (e.g. simulator)"});
            return;
        }
        UIImagePickerController *picker = [[UIImagePickerController alloc] init];
        picker.sourceType = UIImagePickerControllerSourceTypeCamera;

        // Photos use the picker's default media type (public.image), which never
        // throws. Video requires explicitly opting into public.movie — but
        // -setMediaTypes: throws an uncaught NSException ("No available types for
        // source 1") when the camera doesn't currently offer movie capture. That
        // happens on the iOS Simulator, whose camera reports an empty media-type
        // list regardless of authorization. Guard it so a missing video type
        // degrades to a graceful error instead of a SIGABRT. Real devices list
        // public.movie and record normally.
        if (video) {
            @try {
                picker.mediaTypes = @[@"public.movie"];
                picker.cameraCaptureMode = UIImagePickerControllerCameraCaptureModeVideo;
            } @catch (NSException *ex) {
                mfEmit(@"common:capture", @{@"error":
                    @"video recording isn't available from this camera (the iOS Simulator camera can't record video — use a physical device)"});
                return;
            }
        }

        g_mfCameraDelegate = [[MFCameraDelegate alloc] init];
        picker.delegate = g_mfCameraDelegate;
        [mfTopViewController() presentViewController:picker animated:YES completion:nil];
    });
}

void ios_capture_photo(void) { mfPresentCamera(NO); }
void ios_capture_video(void) { mfPresentCamera(YES); }

// MARK: - Background-task window (Phase E)

// UIBackgroundTaskInvalid is an extern const (not a compile-time constant), so a
// separate BOOL tracks whether a task is active rather than initializing to it.
static UIBackgroundTaskIdentifier g_mfBgTask;
static BOOL g_mfBgTaskActive = NO;

static void mfEndBgTask(void) {
    if (g_mfBgTaskActive) {
        [[UIApplication sharedApplication] endBackgroundTask:g_mfBgTask];
        g_mfBgTaskActive = NO;
    }
}

void ios_begin_background_task(int seconds) {
    mfRunOnMain(^{
        mfEndBgTask();
        g_mfBgTask = [[UIApplication sharedApplication] beginBackgroundTaskWithName:@"wails.bgtask"
            expirationHandler:^{
                mfEmit(@"ios:backgroundTask", @{@"message": @"expired"});
                mfEndBgTask();
            }];
        if (g_mfBgTask == UIBackgroundTaskInvalid) {
            mfEmit(@"ios:backgroundTask", @{@"message": @"unavailable"});
            return;
        }
        g_mfBgTaskActive = YES;
        NSTimeInterval remaining = [UIApplication sharedApplication].backgroundTimeRemaining;
        NSString *msg = remaining > 1e8
            ? @"started (full time while in foreground)"
            : [NSString stringWithFormat:@"started (~%.0fs remaining)", remaining];
        mfEmit(@"ios:backgroundTask", @{@"message": msg, @"granted": @(seconds)});
        int secs = seconds > 0 ? seconds : 20;
        dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)secs * NSEC_PER_SEC),
            dispatch_get_main_queue(), ^{
                if (g_mfBgTaskActive) {
                    mfEmit(@"ios:backgroundTask", @{@"message": @"finished"});
                    mfEndBgTask();
                }
            });
    });
}

void ios_end_background_task(void) {
    mfRunOnMain(^{ mfEndBgTask(); });
}
