//go:build darwin && !ios

#ifndef NOTIFICATIONS_DARWIN_H
#define NOTIFICATIONS_DARWIN_H

#import <Foundation/Foundation.h>

bool isNotificationAvailable(void);
bool checkBundleIdentifier(void);
bool ensureDelegateInitialized(void);
void requestNotificationAuthorization(int channelID);
void checkNotificationAuthorization(int channelID);
// sendNotification[WithActions] take a JSON-encoded NotificationOptions blob
// (id, title, subtitle, body, categoryId, data, plus future fields like sound,
// attachments, threadId, interruptionLevel, schedule). Passing one blob keeps
// the C signature stable as new fields are added on the Go side.
void sendNotification(int channelID, const char *options_json);
void sendNotificationWithActions(int channelID, const char *options_json);
void registerNotificationCategory(int channelID, const char *categoryId, const char *actions_json, bool hasReplyField, const char *replyPlaceholder, const char *replyButtonTitle);
void removeNotificationCategory(int channelID, const char *categoryId);
void removeAllPendingNotifications(void);
void removePendingNotification(const char *identifier);
void removeAllDeliveredNotifications(void);
void removeDeliveredNotification(const char *identifier);

#endif /* NOTIFICATIONS_DARWIN_H */