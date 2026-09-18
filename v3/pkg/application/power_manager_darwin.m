//go:build darwin && !ios && !server

#import <Cocoa/Cocoa.h>
#import <IOKit/ps/IOPowerSources.h>
#import <IOKit/ps/IOPSKeys.h>

#include "power_manager_darwin.h"

// powerNotification is the Go side; kind is 0 for a power-state change and
// 1 for a thermal-state change.
extern void powerNotification(int kind);

void* wailsPowerBeginActivity(const char* reason, bool display) {
    @autoreleasepool {
        NSActivityOptions options = NSActivityIdleSystemSleepDisabled;
        if (display) {
            options |= NSActivityIdleDisplaySleepDisabled;
        }
        NSString *reasonString = [NSString stringWithUTF8String:reason];
        id<NSObject> token = [[NSProcessInfo processInfo] beginActivityWithOptions:options reason:reasonString];
        return (void*)[token retain];
    }
}

void wailsPowerEndActivity(void* token) {
    @autoreleasepool {
        id<NSObject> activity = (id<NSObject>)token;
        if (activity == nil) {
            return;
        }
        [[NSProcessInfo processInfo] endActivity:activity];
        [activity release];
    }
}

static int wailsThermalState(void) {
    if (@available(macOS 10.10.3, *)) {
        switch ([NSProcessInfo processInfo].thermalState) {
            case NSProcessInfoThermalStateNominal:  return 1;
            case NSProcessInfoThermalStateFair:     return 2;
            case NSProcessInfoThermalStateSerious:  return 3;
            case NSProcessInfoThermalStateCritical: return 4;
        }
    }
    return 0;
}

static void wailsFillBattery(WailsPowerState *state) {
    state->onBattery = false;
    state->batteryLevel = -1;
    state->charging = false;

    CFTypeRef info = IOPSCopyPowerSourcesInfo();
    if (info == NULL) {
        return;
    }
    CFArrayRef sources = IOPSCopyPowerSourcesList(info);
    if (sources == NULL) {
        CFRelease(info);
        return;
    }
    CFIndex count = CFArrayGetCount(sources);
    for (CFIndex i = 0; i < count; i++) {
        CFDictionaryRef description = IOPSGetPowerSourceDescription(info, CFArrayGetValueAtIndex(sources, i));
        if (description == NULL) {
            continue;
        }
        CFStringRef type = CFDictionaryGetValue(description, CFSTR(kIOPSTypeKey));
        if (type == NULL || !CFEqual(type, CFSTR(kIOPSInternalBatteryType))) {
            continue;
        }
        CFNumberRef current = CFDictionaryGetValue(description, CFSTR(kIOPSCurrentCapacityKey));
        CFNumberRef max = CFDictionaryGetValue(description, CFSTR(kIOPSMaxCapacityKey));
        double currentValue = 0, maxValue = 0;
        if (current != NULL && max != NULL &&
            CFNumberGetValue(current, kCFNumberDoubleType, &currentValue) &&
            CFNumberGetValue(max, kCFNumberDoubleType, &maxValue) && maxValue > 0) {
            state->batteryLevel = currentValue / maxValue;
        }
        CFBooleanRef charging = CFDictionaryGetValue(description, CFSTR(kIOPSIsChargingKey));
        state->charging = charging != NULL && CFBooleanGetValue(charging);
        CFStringRef powerState = CFDictionaryGetValue(description, CFSTR(kIOPSPowerSourceStateKey));
        state->onBattery = powerState != NULL && CFEqual(powerState, CFSTR(kIOPSBatteryPowerValue));
        break;
    }
    CFRelease(sources);
    CFRelease(info);
}

WailsPowerState wailsPowerState(void) {
    @autoreleasepool {
        WailsPowerState state;
        state.lowPowerMode = false;
        if (@available(macOS 12.0, *)) {
            state.lowPowerMode = [NSProcessInfo processInfo].lowPowerModeEnabled;
        }
        state.thermalState = wailsThermalState();
        wailsFillBattery(&state);
        return state;
    }
}

// WailsPowerObserver forwards NSProcessInfo notifications to Go. It is a
// standalone object rather than an AppDelegate selector so the power code
// stays self-contained.
@interface WailsPowerObserver : NSObject
@end

@implementation WailsPowerObserver
- (void)powerStateDidChange:(NSNotification *)notification {
    powerNotification(0);
}
- (void)thermalStateDidChange:(NSNotification *)notification {
    powerNotification(1);
}
@end

static WailsPowerObserver *wailsPowerObserver = nil;

void wailsPowerObserverStart(void) {
    if (wailsPowerObserver != nil) {
        return;
    }
    wailsPowerObserver = [[WailsPowerObserver alloc] init];
    NSNotificationCenter *center = [NSNotificationCenter defaultCenter];
    if (@available(macOS 12.0, *)) {
        [center addObserver:wailsPowerObserver selector:@selector(powerStateDidChange:) name:NSProcessInfoPowerStateDidChangeNotification object:nil];
    }
    if (@available(macOS 10.10.3, *)) {
        [center addObserver:wailsPowerObserver selector:@selector(thermalStateDidChange:) name:NSProcessInfoThermalStateDidChangeNotification object:nil];
    }
}
