//go:build darwin

#ifndef NOTIFICATIONS_DARWIN_H
#define NOTIFICATIONS_DARWIN_H

#import <Foundation/Foundation.h>

bool checkBundleIdentifier(void);
bool requestNotificationAuthorization(void *channelPtr);
bool checkNotificationAuthorization(void *channelPtr);
void sendNotification(const char *identifier, const char *title, const char *subtitle, const char *body, const char *data_json, void *completion);
void sendNotificationWithActions(const char *identifier, const char *title, const char *subtitle, const char *body, const char *categoryId, const char *actions_json, void *completion);
void registerNotificationCategory(const char *categoryId, const char *actions_json, bool hasReplyField, const char *replyPlaceholder, const char *replyButtonTitle);
void removeNotificationCategory(const char *categoryId);
void removeAllPendingNotifications(void);
void removePendingNotification(const char *identifier);
void removeAllDeliveredNotifications(void);
void removeDeliveredNotification(const char *identifier);

extern void requestNotificationAuthorizationResponse(int channelID, bool authorized, const char* error);
extern void checkNotificationAuthorizationResponse(int channelID, bool authorized, const char* error);
extern void didReceiveNotificationResponse(const char *jsonPayload);

#endif /* NOTIFICATIONS_DARWIN_H */