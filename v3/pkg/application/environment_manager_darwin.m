//go:build darwin && !ios && !server

#import <Cocoa/Cocoa.h>
#import <Carbon/Carbon.h>

#include "environment_manager_darwin.h"

// environmentNotification is the Go side; kind is 0 for accessibility
// settings, 1 for keyboard layout and 2 for locale.
extern void environmentNotification(int kind);

WailsAccessibilitySettings wailsAccessibilitySettings(void) {
    @autoreleasepool {
        NSWorkspace *workspace = [NSWorkspace sharedWorkspace];
        WailsAccessibilitySettings settings;
        settings.reduceMotion = workspace.accessibilityDisplayShouldReduceMotion;
        settings.reduceTransparency = workspace.accessibilityDisplayShouldReduceTransparency;
        settings.increaseContrast = workspace.accessibilityDisplayShouldIncreaseContrast;
        settings.differentiateWithoutColor = workspace.accessibilityDisplayShouldDifferentiateWithoutColor;
        settings.invertColors = workspace.accessibilityDisplayShouldInvertColors;
        settings.voiceOverEnabled = false;
        settings.switchControlEnabled = false;
        if (@available(macOS 10.13, *)) {
            settings.voiceOverEnabled = workspace.voiceOverEnabled;
            settings.switchControlEnabled = workspace.switchControlEnabled;
        }
        return settings;
    }
}

static char* wailsJSONCString(NSDictionary *dictionary) {
    NSError *error = nil;
    NSData *data = [NSJSONSerialization dataWithJSONObject:dictionary options:0 error:&error];
    if (data == nil) {
        return strdup("{}");
    }
    NSString *json = [[[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding] autorelease];
    return strdup([json UTF8String]);
}

char* wailsKeyboardLayoutJSON(void) {
    @autoreleasepool {
        TISInputSourceRef source = TISCopyCurrentKeyboardInputSource();
        if (source == NULL) {
            return strdup("{}");
        }
        NSString *identifier = (NSString *)TISGetInputSourceProperty(source, kTISPropertyInputSourceID);
        NSString *name = (NSString *)TISGetInputSourceProperty(source, kTISPropertyLocalizedName);
        NSArray *languages = (NSArray *)TISGetInputSourceProperty(source, kTISPropertyInputSourceLanguages);
        NSDictionary *layout = @{
            @"id": identifier ?: @"",
            @"name": name ?: @"",
            @"languages": languages ?: @[],
        };
        CFRelease(source);
        return wailsJSONCString(layout);
    }
}

char* wailsLocaleJSON(void) {
    @autoreleasepool {
        NSLocale *locale = [NSLocale currentLocale];
        NSString *language = [locale objectForKey:NSLocaleLanguageCode];
        NSString *region = [locale objectForKey:NSLocaleCountryCode];
        NSDictionary *info = @{
            @"identifier": locale.localeIdentifier ?: @"",
            @"language": language ?: @"",
            @"region": region ?: @"",
            @"preferred": [NSLocale preferredLanguages] ?: @[],
        };
        return wailsJSONCString(info);
    }
}

// WailsEnvironmentObserver forwards NSWorkspace, distributed and locale
// notifications to Go as a standalone object so nothing has to be added to
// the AppDelegate.
@interface WailsEnvironmentObserver : NSObject
@end

@implementation WailsEnvironmentObserver
- (void)accessibilityDidChange:(NSNotification *)notification {
    environmentNotification(0);
}
- (void)keyboardLayoutDidChange:(NSNotification *)notification {
    environmentNotification(1);
}
- (void)localeDidChange:(NSNotification *)notification {
    environmentNotification(2);
}
- (void)observeValueForKeyPath:(NSString *)keyPath ofObject:(id)object change:(NSDictionary *)change context:(void *)context {
    // VoiceOver and Switch Control have no notification, only KVO.
    environmentNotification(0);
}
@end

static WailsEnvironmentObserver *wailsEnvironmentObserver = nil;

void wailsEnvironmentObserverStart(void) {
    if (wailsEnvironmentObserver != nil) {
        return;
    }
    wailsEnvironmentObserver = [[WailsEnvironmentObserver alloc] init];

    NSWorkspace *workspace = [NSWorkspace sharedWorkspace];
    [[workspace notificationCenter] addObserver:wailsEnvironmentObserver
                                       selector:@selector(accessibilityDidChange:)
                                           name:NSWorkspaceAccessibilityDisplayOptionsDidChangeNotification
                                         object:nil];
    if (@available(macOS 10.13, *)) {
        [workspace addObserver:wailsEnvironmentObserver forKeyPath:@"voiceOverEnabled" options:0 context:NULL];
        [workspace addObserver:wailsEnvironmentObserver forKeyPath:@"switchControlEnabled" options:0 context:NULL];
    }

    // Input source changes are posted on the distributed centre by the
    // Text Input Services daemon.
    [[NSDistributedNotificationCenter defaultCenter] addObserver:wailsEnvironmentObserver
                                                        selector:@selector(keyboardLayoutDidChange:)
                                                            name:(NSString *)kTISNotifySelectedKeyboardInputSourceChanged
                                                          object:nil];

    [[NSNotificationCenter defaultCenter] addObserver:wailsEnvironmentObserver
                                             selector:@selector(localeDidChange:)
                                                 name:NSCurrentLocaleDidChangeNotification
                                               object:nil];
}
