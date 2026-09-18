//go:build darwin && !ios && !server

#import <Cocoa/Cocoa.h>
#import <AVFoundation/AVFoundation.h>
#import <CoreLocation/CoreLocation.h>
#import <ApplicationServices/ApplicationServices.h>
#import <IOKit/hidsystem/IOHIDLib.h>
#import <UserNotifications/UserNotifications.h>

#include "permissions_manager_darwin.h"

// Delivered to Go when an asynchronous request completes.
extern void permissionsRequestResult(unsigned int requestID, int status);
// Go side of the WKUIDelegate media-capture hook.
extern int permissionsMediaCaptureDecision(unsigned int windowId, int captureType);

static bool wailsIsBundled(void) {
    return [NSBundle mainBundle].bundleIdentifier != nil;
}

// MARK: - Camera / microphone

static int wailsAVStatus(AVMediaType mediaType) {
    if (@available(macOS 10.14, *)) {
        switch ([AVCaptureDevice authorizationStatusForMediaType:mediaType]) {
            case AVAuthorizationStatusNotDetermined: return WailsPermissionStatusNotDetermined;
            case AVAuthorizationStatusRestricted:    return WailsPermissionStatusRestricted;
            case AVAuthorizationStatusDenied:        return WailsPermissionStatusDenied;
            case AVAuthorizationStatusAuthorized:    return WailsPermissionStatusAuthorized;
        }
        return WailsPermissionStatusUnsupported;
    }
    // Pre-Mojave macOS did not gate capture devices.
    return WailsPermissionStatusAuthorized;
}

static int wailsAVRequest(AVMediaType mediaType, unsigned int requestID) {
    if (@available(macOS 10.14, *)) {
        int current = wailsAVStatus(mediaType);
        if (current != WailsPermissionStatusNotDetermined) {
            return current;
        }
        [AVCaptureDevice requestAccessForMediaType:mediaType completionHandler:^(BOOL granted) {
            permissionsRequestResult(requestID, granted ? WailsPermissionStatusAuthorized : WailsPermissionStatusDenied);
        }];
        return WailsPermissionRequestPending;
    }
    return WailsPermissionStatusAuthorized;
}

// MARK: - Screen recording

static int wailsScreenRecordingStatus(void) {
    if (@available(macOS 10.15, *)) {
        return CGPreflightScreenCaptureAccess() ? WailsPermissionStatusAuthorized : WailsPermissionStatusDenied;
    }
    return WailsPermissionStatusAuthorized;
}

static int wailsScreenRecordingRequest(void) {
    if (@available(macOS 10.15, *)) {
        // Prompts once per app; a granted answer only applies after relaunch,
        // which is why the return value is the preflight result, not "yes".
        return CGRequestScreenCaptureAccess() ? WailsPermissionStatusAuthorized : WailsPermissionStatusDenied;
    }
    return WailsPermissionStatusAuthorized;
}

// MARK: - Accessibility

static int wailsAccessibilityStatus(void) {
    return AXIsProcessTrusted() ? WailsPermissionStatusAuthorized : WailsPermissionStatusDenied;
}

static int wailsAccessibilityRequest(void) {
    NSDictionary *options = @{(__bridge NSString *)kAXTrustedCheckOptionPrompt: @YES};
    return AXIsProcessTrustedWithOptions((__bridge CFDictionaryRef)options) ? WailsPermissionStatusAuthorized : WailsPermissionStatusDenied;
}

// MARK: - Input monitoring

static int wailsInputMonitoringStatus(void) {
    if (@available(macOS 10.15, *)) {
        switch (IOHIDCheckAccess(kIOHIDRequestTypeListenEvent)) {
            case kIOHIDAccessTypeGranted: return WailsPermissionStatusAuthorized;
            case kIOHIDAccessTypeDenied:  return WailsPermissionStatusDenied;
            case kIOHIDAccessTypeUnknown: return WailsPermissionStatusNotDetermined;
        }
        return WailsPermissionStatusUnsupported;
    }
    return WailsPermissionStatusAuthorized;
}

static int wailsInputMonitoringRequest(void) {
    if (@available(macOS 10.15, *)) {
        int current = wailsInputMonitoringStatus();
        if (current != WailsPermissionStatusNotDetermined) {
            return current;
        }
        return IOHIDRequestAccess(kIOHIDRequestTypeListenEvent) ? WailsPermissionStatusAuthorized : WailsPermissionStatusDenied;
    }
    return WailsPermissionStatusAuthorized;
}

// MARK: - Location

// WailsLocationDelegate owns the one CLLocationManager the process uses.
// CoreLocation only delivers authorisation changes while the manager that
// asked is alive, so it lives for the whole process.
@interface WailsLocationDelegate : NSObject <CLLocationManagerDelegate>
@property (nonatomic, strong) CLLocationManager *manager;
@property (nonatomic, assign) unsigned int pendingRequestID;
@end

static WailsLocationDelegate *wailsLocationDelegate = nil;

static int wailsLocationStatusFromCL(CLAuthorizationStatus status) {
    switch (status) {
        case kCLAuthorizationStatusNotDetermined: return WailsPermissionStatusNotDetermined;
        case kCLAuthorizationStatusRestricted:    return WailsPermissionStatusRestricted;
        case kCLAuthorizationStatusDenied:        return WailsPermissionStatusDenied;
        default:                                  return WailsPermissionStatusAuthorized;
    }
}

@implementation WailsLocationDelegate

- (instancetype)init {
    self = [super init];
    if (self) {
        _manager = [[CLLocationManager alloc] init];
        _manager.delegate = self;
    }
    return self;
}

- (CLAuthorizationStatus)currentStatus {
    if (@available(macOS 11.0, *)) {
        return self.manager.authorizationStatus;
    }
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
    return [CLLocationManager authorizationStatus];
#pragma clang diagnostic pop
}

- (void)deliver:(CLAuthorizationStatus)status {
    if (self.pendingRequestID == 0 || status == kCLAuthorizationStatusNotDetermined) {
        return;
    }
    unsigned int requestID = self.pendingRequestID;
    self.pendingRequestID = 0;
    permissionsRequestResult(requestID, wailsLocationStatusFromCL(status));
}

- (void)locationManagerDidChangeAuthorization:(CLLocationManager *)manager API_AVAILABLE(macos(11.0)) {
    [self deliver:manager.authorizationStatus];
}

- (void)locationManager:(CLLocationManager *)manager didChangeAuthorizationStatus:(CLAuthorizationStatus)status {
    [self deliver:status];
}

@end

static WailsLocationDelegate *wailsLocation(void) {
    if (wailsLocationDelegate == nil) {
        wailsLocationDelegate = [[WailsLocationDelegate alloc] init];
    }
    return wailsLocationDelegate;
}

static int wailsLocationStatus(void) {
    // Reading the status must not create a manager off the main thread, so
    // fall back to the class method when no manager exists yet.
    if (wailsLocationDelegate != nil) {
        return wailsLocationStatusFromCL([wailsLocationDelegate currentStatus]);
    }
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
    return wailsLocationStatusFromCL([CLLocationManager authorizationStatus]);
#pragma clang diagnostic pop
}

static int wailsLocationRequest(unsigned int requestID) {
    WailsLocationDelegate *location = wailsLocation();
    int current = wailsLocationStatusFromCL([location currentStatus]);
    if (current != WailsPermissionStatusNotDetermined) {
        return current;
    }
    if (@available(macOS 10.15, *)) {
        location.pendingRequestID = requestID;
        [location.manager requestWhenInUseAuthorization];
        return WailsPermissionRequestPending;
    }
    // Before Catalina there is no request API; starting updates prompts.
    location.pendingRequestID = requestID;
    [location.manager startUpdatingLocation];
    return WailsPermissionRequestPending;
}

// MARK: - Notifications

static int wailsNotificationStatusFromUN(NSInteger status) API_AVAILABLE(macos(10.14)) {
    switch ((UNAuthorizationStatus)status) {
        case UNAuthorizationStatusNotDetermined: return WailsPermissionStatusNotDetermined;
        case UNAuthorizationStatusDenied:        return WailsPermissionStatusDenied;
        default:                                 return WailsPermissionStatusAuthorized;
    }
}

static int wailsNotificationsStatus(void) {
    if (!wailsIsBundled()) {
        // UNUserNotificationCenter throws without a bundle identifier.
        return WailsPermissionStatusUnsupported;
    }
    if (@available(macOS 10.14, *)) {
        __block int result = WailsPermissionStatusUnsupported;
        dispatch_semaphore_t done = dispatch_semaphore_create(0);
        [[UNUserNotificationCenter currentNotificationCenter] getNotificationSettingsWithCompletionHandler:^(UNNotificationSettings *settings) {
            result = wailsNotificationStatusFromUN(settings.authorizationStatus);
            dispatch_semaphore_signal(done);
        }];
        // The completion runs on a background queue, so waiting here is safe
        // from any thread including the main thread.
        dispatch_semaphore_wait(done, dispatch_time(DISPATCH_TIME_NOW, 3 * NSEC_PER_SEC));
        dispatch_release(done);
        return result;
    }
    return WailsPermissionStatusUnsupported;
}

static int wailsNotificationsRequest(unsigned int requestID) {
    if (!wailsIsBundled()) {
        return WailsPermissionStatusUnsupported;
    }
    if (@available(macOS 10.14, *)) {
        UNAuthorizationOptions options = UNAuthorizationOptionAlert | UNAuthorizationOptionSound | UNAuthorizationOptionBadge;
        [[UNUserNotificationCenter currentNotificationCenter] requestAuthorizationWithOptions:options completionHandler:^(BOOL granted, NSError *error) {
            permissionsRequestResult(requestID, granted ? WailsPermissionStatusAuthorized : WailsPermissionStatusDenied);
        }];
        return WailsPermissionRequestPending;
    }
    return WailsPermissionStatusUnsupported;
}

// MARK: - Entry points

int wailsPermissionStatus(int kind) {
    @autoreleasepool {
        switch (kind) {
            case WailsPermissionCamera:          return wailsAVStatus(AVMediaTypeVideo);
            case WailsPermissionMicrophone:      return wailsAVStatus(AVMediaTypeAudio);
            case WailsPermissionScreenRecording: return wailsScreenRecordingStatus();
            case WailsPermissionAccessibility:   return wailsAccessibilityStatus();
            case WailsPermissionLocation:        return wailsLocationStatus();
            case WailsPermissionNotifications:   return wailsNotificationsStatus();
            case WailsPermissionInputMonitoring: return wailsInputMonitoringStatus();
            default:                             return WailsPermissionStatusUnsupported;
        }
    }
}

int wailsPermissionRequest(int kind, unsigned int requestID) {
    @autoreleasepool {
        switch (kind) {
            case WailsPermissionCamera:          return wailsAVRequest(AVMediaTypeVideo, requestID);
            case WailsPermissionMicrophone:      return wailsAVRequest(AVMediaTypeAudio, requestID);
            case WailsPermissionScreenRecording: return wailsScreenRecordingRequest();
            case WailsPermissionAccessibility:   return wailsAccessibilityRequest();
            case WailsPermissionLocation:        return wailsLocationRequest(requestID);
            case WailsPermissionNotifications:   return wailsNotificationsRequest(requestID);
            case WailsPermissionInputMonitoring: return wailsInputMonitoringRequest();
            default:                             return WailsPermissionStatusUnsupported;
        }
    }
}

static NSString *wailsSettingsURL(int kind) {
    NSString *privacy = @"x-apple.systempreferences:com.apple.preference.security?";
    switch (kind) {
        case WailsPermissionCamera:          return [privacy stringByAppendingString:@"Privacy_Camera"];
        case WailsPermissionMicrophone:      return [privacy stringByAppendingString:@"Privacy_Microphone"];
        case WailsPermissionScreenRecording: return [privacy stringByAppendingString:@"Privacy_ScreenCapture"];
        case WailsPermissionAccessibility:   return [privacy stringByAppendingString:@"Privacy_Accessibility"];
        case WailsPermissionLocation:        return [privacy stringByAppendingString:@"Privacy_LocationServices"];
        case WailsPermissionInputMonitoring: return [privacy stringByAppendingString:@"Privacy_ListenEvent"];
        case WailsPermissionFullDiskAccess:  return [privacy stringByAppendingString:@"Privacy_AllFiles"];
        case WailsPermissionNotifications:
            if (@available(macOS 13.0, *)) {
                return @"x-apple.systempreferences:com.apple.Notifications-Settings.extension";
            }
            return @"x-apple.systempreferences:com.apple.preference.notifications";
        default:
            return nil;
    }
}

bool wailsPermissionOpenSettings(int kind) {
    @autoreleasepool {
        NSString *url = wailsSettingsURL(kind);
        if (url == nil) {
            return false;
        }
        return [[NSWorkspace sharedWorkspace] openURL:[NSURL URLWithString:url]];
    }
}

int wailsMediaCapturePermissionDecision(unsigned int windowId, int captureType) {
    return permissionsMediaCaptureDecision(windowId, captureType);
}
