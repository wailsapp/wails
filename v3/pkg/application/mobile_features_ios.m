//go:build ios

#import <UIKit/UIKit.h>
#import <AVFoundation/AVFoundation.h>
#import "mobile_features_ios.h"
#import "application_ios_delegate.h"

// appDelegate is defined in application_ios.m.
extern WailsAppDelegate *appDelegate;

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
