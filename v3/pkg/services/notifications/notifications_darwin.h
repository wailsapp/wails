//go:build darwin

#ifndef NOTIFICATIONS_DARWIN_H
#define NOTIFICATIONS_DARWIN_H

#import <Foundation/Foundation.h>

bool isNotificationAvailable(void);
bool checkBundleIdentifier(void);
bool ensureDelegateInitialized(void);
void requestNotificationAuthorization(int channelID);
void checkNotificationAuthorization(int channelID);
void sendNotification(int channelID, const char *identifier, const char *title, const char *subtitle, const char *body, const char *data_json);
void sendNotificationWithActions(int channelID, const char *identifier, const char *title, const char *subtitle, const char *body, const char *categoryId, const char *actions_json);
void registerNotificationCategory(int channelID, const char *categoryId, const char *actions_json, bool hasReplyField, const char *replyPlaceholder, const char *replyButtonTitle);
void removeNotificationCategory(int channelID, const char *categoryId);
void removeAllPendingNotifications(void);
void removePendingNotification(const char *identifier);
void removeAllDeliveredNotifications(void);
void removeDeliveredNotification(const char *identifier);

#endif /* NOTIFICATIONS_DARWIN_H */