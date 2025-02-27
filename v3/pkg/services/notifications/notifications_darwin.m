#import "notifications_darwin.h"
#import <Cocoa/Cocoa.h>
#import <UserNotifications/UserNotifications.h>

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
        NSLog(@"Error: Cannot use the notification API in development mode.\n"
              "  Notifications require the app to be properly bundled with a bundle identifier and signed.\n"
              "  To test notifications:\n"
              "  1. Build and package your app using 'wails3 package'\n"
              "  2. Sign the packaged .app\n"
              "  3. Run the signed .app bundle");
        return false;
    }
    return true;
}

bool requestNotificationAuthorization(int channelID) {
    ensureDelegateInitialized();
    
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    UNAuthorizationOptions options = UNAuthorizationOptionAlert | UNAuthorizationOptionSound | UNAuthorizationOptionBadge;
    
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    UNAuthorizationOptions options = UNAuthorizationOptionAlert | UNAuthorizationOptionSound | UNAuthorizationOptionBadge;
    
    [center requestAuthorizationWithOptions:options completionHandler:^(BOOL granted, NSError * _Nullable error) {
        if (error) {
            requestNotificationAuthorizationResponse(channelID, false, [[error localizedDescription] UTF8String]);
        } else {
            requestNotificationAuthorizationResponse(channelID, granted, NULL);
        }
    }];
}

bool checkNotificationAuthorization(int channelID) {
    ensureDelegateInitialized();
    
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    [center getNotificationSettingsWithCompletionHandler:^(UNNotificationSettings *settings) {
        BOOL isAuthorized = (settings.authorizationStatus == UNAuthorizationStatusAuthorized);
        checkNotificationAuthorizationResponse(channelID, isAuthorized, NULL);
    }];
}

void sendNotification(const char *identifier, const char *title, const char *subtitle, const char *body, const char *data_json, void *completion) {
    ensureDelegateInitialized();
    
    NSString *nsIdentifier = [NSString stringWithUTF8String:identifier];
    NSString *nsTitle = [NSString stringWithUTF8String:title];
    NSString *nsSubtitle = [NSString stringWithUTF8String:subtitle];
    NSString *nsBody = [NSString stringWithUTF8String:body];
    
    NSMutableDictionary *customData = [NSMutableDictionary dictionary];
    if (data_json) {
        NSString *dataJsonStr = [NSString stringWithUTF8String:data_json];
        NSData *jsonData = [dataJsonStr dataUsingEncoding:NSUTF8StringEncoding];
        NSError *error = nil;
        NSDictionary *parsedData = [NSJSONSerialization JSONObjectWithData:jsonData options:0 error:&error];
        if (!error && parsedData) {
            [customData addEntriesFromDictionary:parsedData];
        }
    }
    
    UNUserNotificationCenter *center = [UNUserNotificationCenter currentNotificationCenter];
    
    UNMutableNotificationContent *content = [[UNMutableNotificationContent alloc] init];
    content.title = nsTitle;
    content.subtitle = nsSubtitle;
    content.body = nsBody;
    content.sound = [UNNotificationSound defaultSound];
    
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
    }];
}

void removeNotificationCategory(const char *categoryId) {
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