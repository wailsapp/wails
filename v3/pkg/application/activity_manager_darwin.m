//go:build darwin && !ios && !server

#import <Cocoa/Cocoa.h>

#include "activity_manager_darwin.h"

// Live activities keyed by the Go-side id. Only touched on the main thread.
static NSMutableDictionary<NSNumber *, NSUserActivity *> *wailsActivities(void) {
    static NSMutableDictionary *activities = nil;
    static dispatch_once_t once;
    dispatch_once(&once, ^{
        activities = [[NSMutableDictionary alloc] init];
    });
    return activities;
}

static NSUserActivity *wailsCurrentActivity = nil;

// wailsJSONSafe converts values NSJSONSerialization cannot encode (URLs,
// dates, data, arbitrary objects) into strings so a userInfo dictionary
// received from another device always reaches Go intact.
static id wailsJSONSafe(id value) {
    if (value == nil || value == [NSNull null]) {
        return [NSNull null];
    }
    if ([value isKindOfClass:[NSString class]] || [value isKindOfClass:[NSNumber class]]) {
        return value;
    }
    if ([value isKindOfClass:[NSArray class]]) {
        NSMutableArray *result = [NSMutableArray array];
        for (id item in (NSArray *)value) {
            [result addObject:wailsJSONSafe(item)];
        }
        return result;
    }
    if ([value isKindOfClass:[NSDictionary class]]) {
        NSMutableDictionary *result = [NSMutableDictionary dictionary];
        [(NSDictionary *)value enumerateKeysAndObjectsUsingBlock:^(id key, id object, BOOL *stop) {
            result[[key description]] = wailsJSONSafe(object);
        }];
        return result;
    }
    if ([value isKindOfClass:[NSURL class]]) {
        return [(NSURL *)value absoluteString] ?: @"";
    }
    if ([value isKindOfClass:[NSDate class]]) {
        NSISO8601DateFormatter *formatter = [[[NSISO8601DateFormatter alloc] init] autorelease];
        return [formatter stringFromDate:(NSDate *)value];
    }
    if ([value isKindOfClass:[NSData class]]) {
        return [(NSData *)value base64EncodedStringWithOptions:0];
    }
    return [value description] ?: @"";
}

static char *wailsJSONCStringOrNull(id object) {
    if (object == nil) {
        return strdup("null");
    }
    NSData *data = [NSJSONSerialization dataWithJSONObject:object options:0 error:nil];
    if (data == nil) {
        return strdup("null");
    }
    NSString *json = [[[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding] autorelease];
    return strdup([json UTF8String]);
}

static NSDictionary *wailsParseJSONDictionary(const char *json) {
    if (json == NULL) {
        return nil;
    }
    NSData *data = [NSData dataWithBytes:json length:strlen(json)];
    id parsed = [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
    return [parsed isKindOfClass:[NSDictionary class]] ? parsed : nil;
}

// wailsUserActivityJSON serialises an NSUserActivity in the UserActivity
// shape. It is shared with the application delegate. Caller frees.
char *wailsUserActivityJSON(NSUserActivity *activity) {
    @autoreleasepool {
        if (activity == nil) {
            return strdup("null");
        }
        NSMutableDictionary *result = [NSMutableDictionary dictionary];
        result[@"type"] = activity.activityType ?: @"";
        result[@"title"] = activity.title ?: @"";
        result[@"userInfo"] = wailsJSONSafe(activity.userInfo ?: @{});
        result[@"webpageURL"] = activity.webpageURL.absoluteString ?: @"";
        result[@"eligibleForHandoff"] = @(activity.eligibleForHandoff);
        result[@"eligibleForSearch"] = @(activity.eligibleForSearch);
        // eligibleForPrediction is iOS and watchOS only; macOS has no Siri
        // suggestions for activities.
        result[@"eligibleForPrediction"] = @NO;
        NSMutableArray *keywords = [NSMutableArray array];
        for (NSString *keyword in activity.keywords) {
            [keywords addObject:keyword];
        }
        result[@"keywords"] = keywords;
        return wailsJSONCStringOrNull(result);
    }
}

void wailsActivityPublish(unsigned long long id, const char *json) {
    @autoreleasepool {
        NSDictionary *spec = wailsParseJSONDictionary(json);
        if (spec == nil) {
            return;
        }
        NSString *type = spec[@"type"];
        NSUserActivity *activity = [[NSUserActivity alloc] initWithActivityType:type];
        NSString *title = spec[@"title"];
        if ([title isKindOfClass:[NSString class]] && title.length > 0) {
            activity.title = title;
        }
        NSDictionary *userInfo = spec[@"userInfo"];
        if ([userInfo isKindOfClass:[NSDictionary class]]) {
            activity.userInfo = userInfo;
        }
        NSString *webpage = spec[@"webpageURL"];
        if ([webpage isKindOfClass:[NSString class]] && webpage.length > 0) {
            activity.webpageURL = [NSURL URLWithString:webpage];
        }
        activity.eligibleForHandoff = [spec[@"eligibleForHandoff"] boolValue];
        activity.eligibleForSearch = [spec[@"eligibleForSearch"] boolValue];
        NSArray *keywords = spec[@"keywords"];
        if ([keywords isKindOfClass:[NSArray class]] && keywords.count > 0) {
            activity.keywords = [NSSet setWithArray:keywords];
        }
        wailsActivities()[@(id)] = activity;
        [activity becomeCurrent];
        wailsCurrentActivity = activity;
        [activity release];
    }
}

bool wailsActivityUpdate(unsigned long long id, const char *userInfoJSON) {
    @autoreleasepool {
        NSUserActivity *activity = wailsActivities()[@(id)];
        if (activity == nil) {
            return false;
        }
        NSDictionary *userInfo = wailsParseJSONDictionary(userInfoJSON) ?: @{};
        activity.userInfo = userInfo;
        activity.needsSave = YES;
        return true;
    }
}

void wailsActivityInvalidate(unsigned long long id) {
    @autoreleasepool {
        NSUserActivity *activity = wailsActivities()[@(id)];
        if (activity == nil) {
            return;
        }
        if (wailsCurrentActivity == activity) {
            wailsCurrentActivity = nil;
        }
        [activity invalidate];
        [wailsActivities() removeObjectForKey:@(id)];
    }
}

char *wailsActivityCurrentJSON(void) {
    return wailsUserActivityJSON(wailsCurrentActivity);
}
