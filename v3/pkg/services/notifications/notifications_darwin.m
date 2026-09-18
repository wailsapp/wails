//go:build darwin && !ios

#import "notifications_darwin.h"
#include <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>

#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
#import <UserNotifications/UserNotifications.h>
#endif

bool isNotificationAvailable(void) {
    if (@available(macOS 11.0, *)) {
        return YES;
    } else {
        return NO;
    }
}

bool checkBundleIdentifier(void) {
    NSBundle *main = [NSBundle mainBundle];
    if (main.bundleIdentifier == nil) {
        return NO;
    }
    return YES;
}

extern void captureResult(int channelID, bool success, const char* error);
extern void didReceiveNotificationResponse(const char *jsonPayload, const char* error);

@interface NotificationsDelegate : NSObject <UNUserNotificationCenterDelegate>
@end

@implementation NotificationsDelegate

- (void)userNotificationCenter:(UNUserNotificationCenter *)center
       willPresentNotification:(UNNotification *)notification
         withCompletionHandler:(void (^)(UNNotificationPresentationOptions options))completionHandler {
    UNNotificationPresentationOptions options = 0;
    
    if (@available(macOS 11.0, *)) {
        // These options are only available in macOS 11.0+
        options = UNNotificationPresentationOptionList | 
                  UNNotificationPresentationOptionBanner | 
                  UNNotificationPresentationOptionSound;
    }
    
    completionHandler(options);
}

- (void)userNotificationCenter:(UNUserNotificationCenter *)center
didReceiveNotificationResponse:(UNNotificationResponse *)response
         withCompletionHandler:(void (^)(void))completionHandler {

    NSMutableDictionary *payload = [NSMutableDictionary dictionary];
    
    [payload setObject:response.notification.request.identifier forKey:@"id"];
    [payload setObject:response.actionIdentifier forKey:@"actionIdentifier"];
    [payload setObject:response.notification.request.content.title ?: @"" forKey:@"title"];
    [payload setObject:response.notification.request.content.body ?: @"" forKey:@"body"];
    
    if (response.notification.request.content.categoryIdentifier) {
        [payload setObject:response.notification.request.content.categoryIdentifier forKey:@"categoryIdentifier"];
    }

    if (response.notification.request.content.subtitle) {
        [payload setObject:response.notification.request.content.subtitle forKey:@"subtitle"];
    }
    
    if (response.notification.request.content.userInfo) {
        [payload setObject:response.notification.request.content.userInfo forKey:@"userInfo"];
    }
    
    if ([response isKindOfClass:[UNTextInputNotificationResponse class]]) {
        UNTextInputNotificationResponse *textResponse = (UNTextInputNotificationResponse *)response;
        [payload setObject:textResponse.userText forKey:@"userText"];
    }
    
    NSError *error = nil;
    NSData *jsonData = [NSJSONSerialization dataWithJSONObject:payload options:0 error:&error];
    if (error) {
        NSString *errorMsg = [NSString stringWithFormat:@"Error: %@", [error localizedDescription]];
        didReceiveNotificationResponse(NULL, [errorMsg UTF8String]);
    } else {
        NSString *jsonString = [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
        // The Go side copies the C string during the call
        didReceiveNotificationResponse([jsonString UTF8String], NULL);
        [jsonString release];
    }
    
    completionHandler();
}

@end

static NotificationsDelegate *delegateInstance = nil;
static dispatch_once_t onceToken;

bool ensureDelegateInitialized(void) {
    __block BOOL success = YES;

    dispatch_once(&onceToken, ^{
        delegateInstance = [[NotificationsDelegate alloc] init];
        UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
        center.delegate = delegateInstance;
    });

    if (!delegateInstance) {
        success = NO;
    }

    return success;
}

void requestNotificationAuthorization(int channelID) {
    if (!ensureDelegateInitialized()) {
        NSString *errorMsg = @"Notification delegate has been lost. Reinitialize the notification service.";
        captureResult(channelID, false, [errorMsg UTF8String]);
        return;
    }
    
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    UNAuthorizationOptions options = UNAuthorizationOptionAlert | UNAuthorizationOptionSound | UNAuthorizationOptionBadge;
    
    [center requestAuthorizationWithOptions:options completionHandler:^(BOOL granted, NSError * _Nullable error) {
        if (error) {
            NSString *errorMsg = [NSString stringWithFormat:@"Error: %@", [error localizedDescription]];
            captureResult(channelID, false, [errorMsg UTF8String]);
        } else {
            captureResult(channelID, granted, NULL);
        }
    }];
}

void checkNotificationAuthorization(int channelID) {
    if (!ensureDelegateInitialized()) {
        NSString *errorMsg = @"Notification delegate has been lost. Reinitialize the notification service.";
        captureResult(channelID, false, [errorMsg UTF8String]);
        return;
    }
    
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center getNotificationSettingsWithCompletionHandler:^(UNNotificationSettings *settings) {
        BOOL isAuthorized = (settings.authorizationStatus == UNAuthorizationStatusAuthorized);
        captureResult(channelID, isAuthorized, NULL);
    }];
}

// parseOptions decodes the JSON-encoded NotificationOptions blob handed in
// from Go. Returns nil and writes into *parseError on failure.
static NSDictionary* parseOptions(const char *options_json, NSError **parseError) {
    if (!options_json) {
        if (parseError) {
            *parseError = [NSError errorWithDomain:@"WailsNotifications"
                                              code:1
                                          userInfo:@{NSLocalizedDescriptionKey: @"options_json was NULL"}];
        }
        return nil;
    }
    NSString *str = [NSString stringWithUTF8String:options_json];
    NSData *data = [str dataUsingEncoding:NSUTF8StringEncoding];
    NSDictionary *dict = [NSJSONSerialization JSONObjectWithData:data options:0 error:parseError];
    if (![dict isKindOfClass:[NSDictionary class]]) {
        if (parseError && !*parseError) {
            *parseError = [NSError errorWithDomain:@"WailsNotifications"
                                              code:2
                                          userInfo:@{NSLocalizedDescriptionKey: @"options_json was not a JSON object"}];
        }
        return nil;
    }
    return dict;
}

// stringOrEmpty returns the value for `key` if it is a non-null NSString, else @"".
static NSString* stringOrEmpty(NSDictionary *dict, NSString *key) {
    id val = dict[key];
    if ([val isKindOfClass:[NSString class]]) {
        return (NSString *)val;
    }
    return @"";
}

// applySoundToContent honours the optional Sound field on NotificationOptions.
//   nil  -> default sound (kept on content)
//   {silent: true} -> no sound
//   {name: "Ping"} -> [UNNotificationSound soundNamed:@"Ping"]
static void applySoundToContent(UNMutableNotificationContent *content, NSDictionary *options) {
    id raw = options[@"sound"];
    if (![raw isKindOfClass:[NSDictionary class]]) {
        return;
    }
    NSDictionary *sound = (NSDictionary *)raw;
    id silent = sound[@"silent"];
    if ([silent respondsToSelector:@selector(boolValue)] && [silent boolValue]) {
        content.sound = nil;
        return;
    }
    NSString *name = stringOrEmpty(sound, @"name");
    if (name.length > 0) {
        content.sound = [UNNotificationSound soundNamed:name];
    }
}

// applyAttachmentsToContent translates each NotificationAttachment entry into
// a UNNotificationAttachment. Failed attachments are skipped silently to avoid
// breaking delivery for one bad entry; the failure is logged.
static void applyAttachmentsToContent(UNMutableNotificationContent *content, NSDictionary *options) {
    id raw = options[@"attachments"];
    if (![raw isKindOfClass:[NSArray class]]) {
        return;
    }
    NSMutableArray<UNNotificationAttachment *> *attachments = [NSMutableArray array];
    for (id entry in (NSArray *)raw) {
        if (![entry isKindOfClass:[NSDictionary class]]) {
            continue;
        }
        NSDictionary *att = (NSDictionary *)entry;
        NSString *path = stringOrEmpty(att, @"path");
        if (path.length == 0) {
            continue;
        }
        NSString *attID = stringOrEmpty(att, @"id");
        if (attID.length == 0) {
            attID = [[NSUUID UUID] UUIDString];
        }
        NSURL *url = [path hasPrefix:@"file://"]
            ? [NSURL URLWithString:path]
            : [NSURL fileURLWithPath:path];
        if (!url) continue;

        // The Type field is overloaded: on Windows it carries placement
        // hints like "hero" / "appLogoOverride" / "inline"; on macOS it is
        // an optional UTI hint such as "public.png" / "public.audio".
        // Only forward the value as a UTI hint when it actually looks like
        // one (UTI strings always contain a "."). Otherwise let the
        // notification center auto-infer the type from the file extension.
        NSDictionary *attOptions = nil;
        NSString *uti = stringOrEmpty(att, @"type");
        if (uti.length > 0 && [uti containsString:@"."]) {
            attOptions = @{UNNotificationAttachmentOptionsTypeHintKey: uti};
        }

        NSError *err = nil;
        UNNotificationAttachment *a = [UNNotificationAttachment
            attachmentWithIdentifier:attID URL:url options:attOptions error:&err];
        if (err || !a) {
            NSLog(@"wails/notifications: failed to attach %@: %@", path, err);
            continue;
        }
        [attachments addObject:a];
    }
    if (attachments.count > 0) {
        content.attachments = attachments;
    }
}

// applyInterruptionLevelToContent maps the InterruptionLevel string onto
// UNNotificationInterruptionLevel (macOS 12+). On older macOS it is a no-op.
static void applyInterruptionLevelToContent(UNMutableNotificationContent *content, NSDictionary *options) {
    NSString *level = stringOrEmpty(options, @"interruptionLevel");
    if (level.length == 0) {
        return;
    }
    if (@available(macOS 12.0, *)) {
        if ([level isEqualToString:@"passive"]) {
            content.interruptionLevel = UNNotificationInterruptionLevelPassive;
        } else if ([level isEqualToString:@"active"]) {
            content.interruptionLevel = UNNotificationInterruptionLevelActive;
        } else if ([level isEqualToString:@"timeSensitive"]) {
            content.interruptionLevel = UNNotificationInterruptionLevelTimeSensitive;
        } else if ([level isEqualToString:@"critical"]) {
            // Requires the Critical Alert entitlement; without it macOS
            // silently downgrades to active.
            content.interruptionLevel = UNNotificationInterruptionLevelCritical;
        }
    }
}

// buildTriggerFromSchedule returns a UNCalendar/UNTimeInterval trigger built
// from the optional Schedule field, or nil for immediate delivery.
//   {delaySeconds: 30} -> UNTimeIntervalNotificationTrigger 30s
//   {at: 1717181600}    -> UNCalendarNotificationTrigger at the corresponding date
static UNNotificationTrigger* buildTriggerFromSchedule(NSDictionary *options) {
    id raw = options[@"schedule"];
    if (![raw isKindOfClass:[NSDictionary class]]) {
        return nil;
    }
    NSDictionary *schedule = (NSDictionary *)raw;

    id delayObj = schedule[@"delaySeconds"];
    if ([delayObj respondsToSelector:@selector(doubleValue)]) {
        double delay = [delayObj doubleValue];
        if (delay > 0) {
            return [UNTimeIntervalNotificationTrigger triggerWithTimeInterval:delay repeats:NO];
        }
    }

    id atObj = schedule[@"at"];
    if ([atObj respondsToSelector:@selector(doubleValue)]) {
        double at = [atObj doubleValue];
        if (at > 0) {
            NSDate *date = [NSDate dateWithTimeIntervalSince1970:at];
            NSCalendar *cal = [NSCalendar currentCalendar];
            NSDateComponents *components = [cal components:
                NSCalendarUnitYear | NSCalendarUnitMonth | NSCalendarUnitDay |
                NSCalendarUnitHour | NSCalendarUnitMinute | NSCalendarUnitSecond
                fromDate:date];
            return [UNCalendarNotificationTrigger triggerWithDateMatchingComponents:components repeats:NO];
        }
    }

    return nil;
}

// createNotificationContentFromOptions builds the UNMutableNotificationContent
// from a parsed NotificationOptions dict. Reads title/subtitle/body/userInfo
// plus the optional sound, attachments, threadId, and interruptionLevel
// fields. Schedule is handled separately by buildTriggerFromSchedule.
static UNMutableNotificationContent* createNotificationContentFromOptions(NSDictionary *options) {
    UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];
    content.title = stringOrEmpty(options, @"title");
    NSString *subtitle = stringOrEmpty(options, @"subtitle");
    if (subtitle.length > 0) {
        content.subtitle = subtitle;
    }
    content.body = stringOrEmpty(options, @"body");
    content.sound = [UNNotificationSound defaultSound];

    id userInfo = options[@"data"];
    if ([userInfo isKindOfClass:[NSDictionary class]]) {
        content.userInfo = (NSDictionary *)userInfo;
    }

    NSString *threadId = stringOrEmpty(options, @"threadId");
    if (threadId.length > 0) {
        content.threadIdentifier = threadId;
    }

    applySoundToContent(content, options);
    applyAttachmentsToContent(content, options);
    applyInterruptionLevelToContent(content, options);

    return content;
}

void sendNotification(int channelID, const char *options_json) {
    if (!ensureDelegateInitialized()) {
        NSString *errorMsg = @"Notification delegate has been lost. Reinitialize the notification service.";
        captureResult(channelID, false, [errorMsg UTF8String]);
        return;
    }

    NSError *parseError = nil;
    NSDictionary *options = parseOptions(options_json, &parseError);
    if (parseError || !options) {
        NSString *errorMsg = [NSString stringWithFormat:@"Error: %@", [parseError localizedDescription]];
        captureResult(channelID, false, [errorMsg UTF8String]);
        return;
    }

    NSString *identifier = stringOrEmpty(options, @"id");
    UNMutableNotificationContent *content = createNotificationContentFromOptions(options);

    UNNotificationTrigger *trigger = buildTriggerFromSchedule(options);

    UNNotificationRequest *request = [UNNotificationRequest requestWithIdentifier:identifier content:content trigger:trigger];
    // The request keeps its own copy of the content
    [content release];

    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center addNotificationRequest:request withCompletionHandler:^(NSError * _Nullable error) {
        if (error) {
            NSString *errorMsg = [NSString stringWithFormat:@"Error: %@", [error localizedDescription]];
            captureResult(channelID, false, [errorMsg UTF8String]);
        } else {
            captureResult(channelID, true, NULL);
        }
    }];
}

void sendNotificationWithActions(int channelID, const char *options_json) {
    if (!ensureDelegateInitialized()) {
        NSString *errorMsg = @"Notification delegate has been lost. Reinitialize the notification service.";
        captureResult(channelID, false, [errorMsg UTF8String]);
        return;
    }

    NSError *parseError = nil;
    NSDictionary *options = parseOptions(options_json, &parseError);
    if (parseError || !options) {
        NSString *errorMsg = [NSString stringWithFormat:@"Error: %@", [parseError localizedDescription]];
        captureResult(channelID, false, [errorMsg UTF8String]);
        return;
    }

    NSString *identifier = stringOrEmpty(options, @"id");
    NSString *categoryId = stringOrEmpty(options, @"categoryId");

    UNMutableNotificationContent *content = createNotificationContentFromOptions(options);
    if (categoryId.length > 0) {
        content.categoryIdentifier = categoryId;
    }

    UNNotificationTrigger *trigger = buildTriggerFromSchedule(options);

    UNNotificationRequest *request = [UNNotificationRequest requestWithIdentifier:identifier content:content trigger:trigger];
    // The request keeps its own copy of the content
    [content release];

    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center addNotificationRequest:request withCompletionHandler:^(NSError * _Nullable error) {
        if (error) {
            NSString *errorMsg = [NSString stringWithFormat:@"Error: %@", [error localizedDescription]];
            captureResult(channelID, false, [errorMsg UTF8String]);
        } else {
            captureResult(channelID, true, NULL);
        }
    }];
}

void registerNotificationCategory(int channelID, const char *categoryId, const char *actions_json, bool hasReplyField, 
                                const char *replyPlaceholder, const char *replyButtonTitle) {
    if (!ensureDelegateInitialized()) {
        NSString *errorMsg = @"Notification delegate has been lost. Reinitialize the notification service.";
        captureResult(channelID, false, [errorMsg UTF8String]);
        return;
    }
    
    NSString *nsCategoryId = [NSString stringWithUTF8String:categoryId];
    NSString *actionsJsonStr = actions_json ? [NSString stringWithUTF8String:actions_json] : @"[]";
    
    NSData *jsonData = [actionsJsonStr dataUsingEncoding:NSUTF8StringEncoding];
    NSError *error = nil;
    NSArray *actionsArray = [NSJSONSerialization JSONObjectWithData:jsonData options:0 error:&error];
    
    if (error) {
        NSString *errorMsg = [NSString stringWithFormat:@"Error: %@", [error localizedDescription]];
        captureResult(channelID, false, [errorMsg UTF8String]);
        return;
    }
    
    NSMutableArray *actions = [NSMutableArray array];
    
    for (NSDictionary *actionDict in actionsArray) {
        NSString *actionId = actionDict[@"id"];
        NSString *actionTitle = actionDict[@"title"];
        BOOL destructive = [actionDict[@"destructive"] boolValue];
        
        if (actionId && actionTitle) {
            UNNotificationActionOptions options = UNNotificationActionOptionNone;
            if (destructive) options |= UNNotificationActionOptionDestructive;
            
            UNNotificationAction *action = [UNNotificationAction 
                                          actionWithIdentifier:actionId
                                          title:actionTitle
                                          options:options];
            [actions addObject:action];
        }
    }
    
    if (hasReplyField && replyPlaceholder && replyButtonTitle) {
        NSString *placeholder = [NSString stringWithUTF8String:replyPlaceholder];
        NSString *buttonTitle = [NSString stringWithUTF8String:replyButtonTitle];
        
        UNTextInputNotificationAction *textAction = 
            [UNTextInputNotificationAction actionWithIdentifier:@"TEXT_REPLY"
                                                         title:buttonTitle
                                                       options:UNNotificationActionOptionNone
                                          textInputButtonTitle:buttonTitle
                                          textInputPlaceholder:placeholder];
        [actions addObject:textAction];
    }
    
    UNNotificationCategory *newCategory = [UNNotificationCategory 
                                      categoryWithIdentifier:nsCategoryId
                                      actions:actions
                                      intentIdentifiers:@[]
                                      options:UNNotificationCategoryOptionNone];
    
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center getNotificationCategoriesWithCompletionHandler:^(NSSet<UNNotificationCategory *> *categories) {
        NSMutableSet *updatedCategories = [NSMutableSet setWithSet:categories];
        
        // Remove existing category with same ID if it exists
        UNNotificationCategory *existingCategory = nil;
        for (UNNotificationCategory *category in updatedCategories) {
            if ([category.identifier isEqualToString:nsCategoryId]) {
                existingCategory = category;
                break;
            }
        }
        if (existingCategory) {
            [updatedCategories removeObject:existingCategory];
        }
        
        // Add the new category
        [updatedCategories addObject:newCategory];
        [center setNotificationCategories:updatedCategories];

        captureResult(channelID, true, NULL);
    }];
}

void removeNotificationCategory(int channelID, const char *categoryId) {
    NSString *nsCategoryId = [NSString stringWithUTF8String:categoryId];
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    
    [center getNotificationCategoriesWithCompletionHandler:^(NSSet<UNNotificationCategory *> *categories) {
        NSMutableSet *updatedCategories = [NSMutableSet setWithSet:categories];
        
        // Find and remove the category with matching identifier
        UNNotificationCategory *categoryToRemove = nil;
        for (UNNotificationCategory *category in updatedCategories) {
            if ([category.identifier isEqualToString:nsCategoryId]) {
                categoryToRemove = category;
                break;
            }
        }
        
        if (categoryToRemove) {
            [updatedCategories removeObject:categoryToRemove];
            [center setNotificationCategories:updatedCategories];
            captureResult(channelID, true, NULL);
        } else {
            NSString *errorMsg = [NSString stringWithFormat:@"Category '%@' not found", nsCategoryId];
            captureResult(channelID, false, [errorMsg UTF8String]);
        }
    }];
}

void removeAllPendingNotifications(void) {
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center removeAllPendingNotificationRequests];
}

void removePendingNotification(const char *identifier) {
    NSString *nsIdentifier = [NSString stringWithUTF8String:identifier];
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center removePendingNotificationRequestsWithIdentifiers:@[nsIdentifier]];
}

void removeAllDeliveredNotifications(void) {
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center removeAllDeliveredNotifications];
}

void removeDeliveredNotification(const char *identifier) {
    NSString *nsIdentifier = [NSString stringWithUTF8String:identifier];
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center removeDeliveredNotificationsWithIdentifiers:@[nsIdentifier]];
}