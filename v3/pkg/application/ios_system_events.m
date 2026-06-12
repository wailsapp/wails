//go:build ios
#import <Foundation/Foundation.h>
#import <UIKit/UIKit.h>
#import <Network/Network.h>
#import "application_ios.h"

// emitSystemEvent is a Go //export (application_ios.go). Declared locally so we
// don't pull in cgo's _cgo_export.h (which would clash with the header). The
// linker matches on symbol name; const-ness is not part of the C symbol.
extern void emitSystemEvent(const char* name, const char* json);

// Emit a "system:*" event with a JSON payload string.
static void emitSys(NSString *name, NSString *json) {
    emitSystemEvent([name UTF8String], json ? [json UTF8String] : "");
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
    NSString *json = [NSString stringWithFormat:
        @"{\"level\":%.2f,\"state\":\"%@\",\"lowPowerMode\":%@}",
        level, state, lowPower ? @"true" : @"false"];
    emitSys(@"system:battery", json);
}

// ---- Theme ---------------------------------------------------------------

static void emitTheme(void) {
    emitSys(@"system:theme", ios_is_dark_mode() ? @"{\"dark\":true}" : @"{\"dark\":false}");
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
    NSString *json = [NSString stringWithFormat:
        @"{\"connected\":%@,\"type\":\"%@\",\"expensive\":%@,\"constrained\":%@}",
        connected ? @"true" : @"false", type,
        expensive ? @"true" : @"false", constrained ? @"true" : @"false"];
    emitSys(@"system:network", json);
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

            // App lifecycle -> system:appstate.
            [nc addObserverForName:UIApplicationDidBecomeActiveNotification
                            object:nil queue:nil
                        usingBlock:^(NSNotification *n){
                emitSys(@"system:appstate", @"{\"state\":\"active\"}");
                // Reopening the app is a good moment to refresh snapshots.
                emitBattery();
                emitTheme();
            }];
            [nc addObserverForName:UIApplicationWillResignActiveNotification
                            object:nil queue:nil
                        usingBlock:^(NSNotification *n){
                emitSys(@"system:appstate", @"{\"state\":\"inactive\"}");
            }];
            [nc addObserverForName:UIApplicationDidEnterBackgroundNotification
                            object:nil queue:nil
                        usingBlock:^(NSNotification *n){
                emitSys(@"system:appstate", @"{\"state\":\"background\"}");
            }];
            [nc addObserverForName:UIApplicationWillEnterForegroundNotification
                            object:nil queue:nil
                        usingBlock:^(NSNotification *n){
                emitSys(@"system:appstate", @"{\"state\":\"foreground\"}");
            }];

            // Memory warning -> system:memory.
            [nc addObserverForName:UIApplicationDidReceiveMemoryWarningNotification
                            object:nil queue:nil
                        usingBlock:^(NSNotification *n){
                emitSys(@"system:memory", @"{}");
            }];

            // Lock / unlock: approximated via data-protection availability.
            // Only fires when the device has a passcode set.
            [nc addObserverForName:UIApplicationProtectedDataWillBecomeUnavailable
                            object:nil queue:nil
                        usingBlock:^(NSNotification *n){
                emitSys(@"system:lock", @"{\"locked\":true}");
            }];
            [nc addObserverForName:UIApplicationProtectedDataDidBecomeAvailable
                            object:nil queue:nil
                        usingBlock:^(NSNotification *n){
                emitSys(@"system:lock", @"{\"locked\":false}");
            }];

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
