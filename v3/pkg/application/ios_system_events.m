//go:build ios
#import <Foundation/Foundation.h>
#import <UIKit/UIKit.h>
#import <Network/Network.h>
#import "application_ios.h"
#import "../events/events_ios.h"

// processApplicationEvent is a Go //export (application_ios.go). It takes an
// event ID and an optional data pointer; we pass a JSON object string, which Go
// decodes into the ApplicationEvent context. Declared locally (matching the
// other .m files) so we don't pull in cgo's _cgo_export.h.
extern void processApplicationEvent(unsigned int eventID, void* data);

// Emit an application event with an optional JSON payload string.
static void emitEvent(unsigned int eventID, NSString *json) {
    processApplicationEvent(eventID, json ? (void *)[json UTF8String] : NULL);
}

// ---- Battery -------------------------------------------------------------

static NSString *batteryStateString(UIDeviceBatteryState s) {
    switch (s) {
        case UIDeviceBatteryStateCharging:  return @"charging";
        case UIDeviceBatteryStateFull:      return @"full";
        case UIDeviceBatteryStateUnplugged: return @"unplugged";
        default:                            return @"unknown";
    }
}

static void emitBattery(void) {
    UIDevice *dev = [UIDevice currentDevice];
    float level = dev.batteryLevel; // 0..1, or -1 when unknown (e.g. Simulator)
    NSString *state = batteryStateString(dev.batteryState);
    BOOL lowPower = [NSProcessInfo processInfo].lowPowerModeEnabled;
    emitEvent(EventBatteryChanged, [NSString stringWithFormat:
        @"{\"level\":%.2f,\"state\":\"%@\",\"lowPowerMode\":%@}",
        level, state, lowPower ? @"true" : @"false"]);
}

// ---- Theme ---------------------------------------------------------------
// "isDarkMode" matches the key the desktop platforms set (ApplicationEventContext.IsDarkMode()).

static void emitTheme(void) {
    emitEvent(EventThemeChanged, ios_is_dark_mode() ? @"{\"isDarkMode\":true}" : @"{\"isDarkMode\":false}");
}

// ---- Network (NWPathMonitor) --------------------------------------------

static nw_path_monitor_t g_pathMonitor = nil;

static void emitNetwork(nw_path_t path) {
    BOOL connected = (nw_path_get_status(path) == nw_path_status_satisfied);
    NSString *type = @"none";
    if (connected) {
        if (nw_path_uses_interface_type(path, nw_interface_type_wifi)) {
            type = @"wifi";
        } else if (nw_path_uses_interface_type(path, nw_interface_type_cellular)) {
            type = @"cellular";
        } else if (nw_path_uses_interface_type(path, nw_interface_type_wired)) {
            type = @"wired";
        } else {
            type = @"other";
        }
    }
    BOOL expensive = nw_path_is_expensive(path);   // cellular / hotspot
    BOOL constrained = NO;                          // Low Data Mode
    if (@available(iOS 13.0, *)) {
        constrained = nw_path_is_constrained(path);
    }
    emitEvent(EventNetworkChanged, [NSString stringWithFormat:
        @"{\"connected\":%@,\"type\":\"%@\",\"expensive\":%@,\"constrained\":%@}",
        connected ? @"true" : @"false", type,
        expensive ? @"true" : @"false", constrained ? @"true" : @"false"]);
}

// ---- Setup ---------------------------------------------------------------

void ios_start_system_event_monitors(void) {
    static dispatch_once_t once;
    dispatch_once(&once, ^{
        dispatch_async(dispatch_get_main_queue(), ^{
            NSNotificationCenter *nc = [NSNotificationCenter defaultCenter];

            // Battery: monitoring must be enabled before level/state are valid.
            [UIDevice currentDevice].batteryMonitoringEnabled = YES;
            [nc addObserverForName:UIDeviceBatteryLevelDidChangeNotification
                            object:nil queue:nil
                        usingBlock:^(NSNotification *n){ emitBattery(); }];
            [nc addObserverForName:UIDeviceBatteryStateDidChangeNotification
                            object:nil queue:nil
                        usingBlock:^(NSNotification *n){ emitBattery(); }];
            [nc addObserverForName:NSProcessInfoPowerStateDidChangeNotification
                            object:nil queue:nil
                        usingBlock:^(NSNotification *n){ emitBattery(); }];

            // Refresh battery/theme snapshots when the app comes to the front.
            // (Lifecycle itself is delivered by the generated UIApplication
            // delegate events, so we don't re-emit those here.)
            [nc addObserverForName:UIApplicationDidBecomeActiveNotification
                            object:nil queue:nil
                        usingBlock:^(NSNotification *n){ emitBattery(); emitTheme(); }];

            // Lock / unlock: approximated via data-protection availability.
            // Only fires when the device has a passcode set.
            [nc addObserverForName:UIApplicationProtectedDataWillBecomeUnavailable
                            object:nil queue:nil
                        usingBlock:^(NSNotification *n){ emitEvent(EventScreenLocked, nil); }];
            [nc addObserverForName:UIApplicationProtectedDataDidBecomeAvailable
                            object:nil queue:nil
                        usingBlock:^(NSNotification *n){ emitEvent(EventScreenUnlocked, nil); }];

            // Network reachability / interface type. The update handler fires
            // immediately with the current path, giving listeners an initial
            // value as well as live changes.
            g_pathMonitor = nw_path_monitor_create();
            nw_path_monitor_set_queue(g_pathMonitor, dispatch_get_main_queue());
            nw_path_monitor_set_update_handler(g_pathMonitor, ^(nw_path_t path){
                emitNetwork(path);
            });
            nw_path_monitor_start(g_pathMonitor);

            // Initial snapshot so listeners mounted before any change still get
            // current battery/theme.
            emitBattery();
            emitTheme();
        });
    });
}
