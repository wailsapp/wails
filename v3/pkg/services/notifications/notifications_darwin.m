#import "notifications_darwin.h"
#import <Cocoa/Cocoa.h>
#import <UserNotifications/UserNotifications.h>

extern void didReceiveNotificationResponse(const char *jsonPayload);

@interface NotificationsDelegate : NSObject <UNUserNotificationCenterDelegate>
@end

@implementation NotificationsDelegate

- (void)userNotificationCenter:(UNUserNotificationCenter *)center
       willPresentNotification:(UNNotification *)notification
         withCompletionHandler:(void (^)(UNNotificationPresentationOptions options))completionHandler {
    UNNotificationPresentationOptions options = UNNotificationPresentationOptionList | 
                                                UNNotificationPresentationOptionBanner | 
                                                UNNotificationPresentationOptionSound;
    completionHandler(options);
}

- (void)userNotificationCenter:(UNUserNotificationCenter *)center
didReceiveNotificationResponse:(UNNotificationResponse *)response
         withCompletionHandler:(void (^)(void))completionHandler {

    NSMutableDictionary *payload = [NSMutableDictionary dictionary];
    
    [payload setObject:response.notification.request.identifier forKey:@"identifier"];
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
    if (!error) {
        NSString *jsonString = [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
        didReceiveNotificationResponse([jsonString UTF8String]);
    }
    
    completionHandler();
}

@end

static NotificationsDelegate *delegateInstance = nil;

static void ensureDelegateInitialized(void) {
    if (!delegateInstance) {
        delegateInstance = [[NotificationsDelegate alloc] init];
        UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
        center.delegate = delegateInstance;
    }
}

bool checkBundleIdentifier(void) {
    NSBundle *main = [NSBundle mainBundle];
    if (main.bundleIdentifier == nil) {
        NSLog(@"Error: Cannot use notifications in development mode.\n"
              "  Notifications require the app to be properly bundled with a bundle identifier.\n"
              "  To test notifications:\n"
              "  1. Build and package your app using 'wails3 package'\n"
              "  2. Sign the packaged .app\n"
              "  3. Run the signed .app bundle");
        return false;
    }
    return true;
}

bool requestUserNotificationAuthorization(void *completion) {
    if (!checkBundleIdentifier()) {
        return false;
    }
    
    ensureDelegateInitialized();
    
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    UNAuthorizationOptions options = UNAuthorizationOptionAlert | UNAuthorizationOptionSound | UNAuthorizationOptionBadge;
    
    [center requestAuthorizationWithOptions:options completionHandler:^(BOOL granted, NSError * _Nullable error) {
        if (completion != NULL) {
            void (^callback)(NSError *, BOOL) = completion;
            callback(error, granted);
        }
    }];
    return true;
}

bool checkNotificationAuthorization(void) {
    ensureDelegateInitialized();
    
    __block BOOL isAuthorized = NO;
    dispatch_semaphore_t semaphore = dispatch_semaphore_create(0);
    
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center getNotificationSettingsWithCompletionHandler:^(UNNotificationSettings * _Nonnull settings) {
        isAuthorized = (settings.authorizationStatus == UNAuthorizationStatusAuthorized);
        dispatch_semaphore_signal(semaphore);
    }];
    
    // Wait for response with a timeout
    dispatch_semaphore_wait(semaphore, dispatch_time(DISPATCH_TIME_NOW, 3 * NSEC_PER_SEC));
    return isAuthorized;
}

void sendNotification(const char *identifier, const char *title, const char *subtitle, const char *body, void *completion) {
    ensureDelegateInitialized();
    
    NSString *nsIdentifier = [NSString stringWithUTF8String:identifier];
    NSString *nsTitle = [NSString stringWithUTF8String:title];
    NSString *nsSubtitle = [NSString stringWithUTF8String:subtitle];
    NSString *nsBody = [NSString stringWithUTF8String:body];
    
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    
    UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];
    content.title = nsTitle;
    content.subtitle = nsSubtitle;
    content.body = nsBody;
    content.sound = [UNNotificationSound defaultSound];
    
    UNTimeIntervalNotificationTrigger *trigger = [UNTimeIntervalNotificationTrigger triggerWithTimeInterval:1 repeats:NO];
    
    UNNotificationRequest *request = [UNNotificationRequest requestWithIdentifier:nsIdentifier content:content trigger:trigger];
    
    [center addNotificationRequest:request withCompletionHandler:^(NSError * _Nullable error) {
        if (completion != NULL) {
            void (^callback)(NSError *) = completion;
            callback(error);
        }
    }];
}

void sendNotificationWithActions(const char *identifier, const char *title, const char *subtitle, 
                             const char *body, const char *categoryId, const char *actions_json, void *completion) {
    ensureDelegateInitialized();
    
    NSString *nsIdentifier = [NSString stringWithUTF8String:identifier];
    NSString *nsTitle = [NSString stringWithUTF8String:title];
    NSString *nsSubtitle = subtitle ? [NSString stringWithUTF8String:subtitle] : @"";
    NSString *nsBody = [NSString stringWithUTF8String:body];
    NSString *nsCategoryId = [NSString stringWithUTF8String:categoryId];
    
    NSMutableDictionary *customData = [NSMutableDictionary dictionary];
    if (actions_json) {
        NSString *actionsJsonStr = [NSString stringWithUTF8String:actions_json];
        NSData *jsonData = [actionsJsonStr dataUsingEncoding:NSUTF8StringEncoding];
        NSError *error = nil;
        NSDictionary *parsedData = [NSJSONSerialization JSONObjectWithData:jsonData options:0 error:&error];
        if (!error && parsedData) {
            [customData addEntriesFromDictionary:parsedData];
        }
    }
    
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    
    UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];
    content.title = nsTitle;
    if (![nsSubtitle isEqualToString:@""]) {
        content.subtitle = nsSubtitle;
    }
    content.body = nsBody;
    content.sound = [UNNotificationSound defaultSound];
    content.categoryIdentifier = nsCategoryId;
    
    // Add custom data if available
    if (customData.count > 0) {
        content.userInfo = customData;
    }
    
    UNTimeIntervalNotificationTrigger *trigger = [UNTimeIntervalNotificationTrigger triggerWithTimeInterval:1 repeats:NO];
    
    UNNotificationRequest *request = [UNNotificationRequest requestWithIdentifier:nsIdentifier content:content trigger:trigger];
    
    [center addNotificationRequest:request withCompletionHandler:^(NSError * _Nullable error) {
        if (completion != NULL) {
            void (^callback)(NSError *) = completion;
            callback(error);
        }
    }];
}

void registerNotificationCategory(const char *categoryId, const char *actions_json, bool hasReplyField, 
                                const char *replyPlaceholder, const char *replyButtonTitle) {
    ensureDelegateInitialized();
    
    NSString *nsCategoryId = [NSString stringWithUTF8String:categoryId];
    NSString *actionsJsonStr = actions_json ? [NSString stringWithUTF8String:actions_json] : @"[]";
    
    NSData *jsonData = [actionsJsonStr dataUsingEncoding:NSUTF8StringEncoding];
    NSError *error = nil;
    NSArray *actionsArray = [NSJSONSerialization JSONObjectWithData:jsonData options:0 error:&error];
    
    if (error) {
        NSLog(@"Error parsing notification actions JSON: %@", error);
        return;
    }
    
    NSMutableArray *actions = [NSMutableArray array];
    
    for (NSDictionary *actionDict in actionsArray) {
        NSString *actionId = actionDict[@"id"];
        NSString *actionTitle = actionDict[@"title"];
        BOOL destructive = [actionDict[@"destructive"] boolValue];
        BOOL authRequired = [actionDict[@"authenticationRequired"] boolValue];
        
        if (actionId && actionTitle) {
            UNNotificationActionOptions options = UNNotificationActionOptionNone;
            if (destructive) options |= UNNotificationActionOptionDestructive;
            if (authRequired) options |= UNNotificationActionOptionAuthenticationRequired;
            
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
    
    UNNotificationCategory *category = [UNNotificationCategory 
                                      categoryWithIdentifier:nsCategoryId
                                      actions:actions
                                      intentIdentifiers:@[]
                                      options:UNNotificationCategoryOptionNone];
    
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center getNotificationCategoriesWithCompletionHandler:^(NSSet<UNNotificationCategory *> *categories) {
        NSMutableSet *updatedCategories = [NSMutableSet setWithSet:categories];
        [updatedCategories addObject:category];
        [center setNotificationCategories:updatedCategories];
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