//go:build darwin

#ifndef NOTIFICATIONS_DARWIN_H
#define NOTIFICATIONS_DARWIN_H

#import <Foundation/Foundation.h>

bool checkBundleIdentifier(void);
bool requestUserNotificationAuthorization(void *completion);
bool checkNotificationAuthorization(void);
void sendNotification(const char *identifier, const char *title, const char *subtitle, const char *body, void *completion);
void sendNotificationWithActions(const char *identifier, const char *title, const char *subtitle, 
                               const char *body, const char *categoryId, const char *actions_json, void *completion);
void registerNotificationCategory(const char *categoryId, const char *actions_json, bool hasReplyField, 
                                const char *replyPlaceholder, const char *replyButtonTitle);
void removeAllPendingNotifications(void);
void removePendingNotification(const char *identifier);
void removeAllDeliveredNotifications(void);
void removeDeliveredNotification(const char *identifier);

#endif /* NOTIFICATIONS_DARWIN_H */