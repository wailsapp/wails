package runtime

import (
	"context"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

// NotificationOptions contains configuration for a notification.
type NotificationOptions = frontend.NotificationOptions

// NotificationAction represents an action button for a notification.
type NotificationAction = frontend.NotificationAction

// NotificationCategory groups actions for notifications.
type NotificationCategory = frontend.NotificationCategory

// NotificationResponse represents the response sent by interacting with a notification.
type NotificationResponse = frontend.NotificationResponse

// NotificationResult represents the result of a notification response,
// returning the response or any errors that occurred.
type NotificationResult = frontend.NotificationResult

func InitializeNotifications(ctx context.Context) error {
	appFrontend := getFrontend(ctx)
	return appFrontend.InitializeNotifications()
}

func CleanupNotifications(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	appFrontend.CleanupNotifications()
}

func IsNotificationAvailable(ctx context.Context) bool {
	appFrontend := getFrontend(ctx)
	return appFrontend.IsNotificationAvailable()
}

func RequestNotificationAuthorization(ctx context.Context) (bool, error) {
	appFrontend := getFrontend(ctx)
	return appFrontend.RequestNotificationAuthorization()
}

func CheckNotificationAuthorization(ctx context.Context) (bool, error) {
	appFrontend := getFrontend(ctx)
	return appFrontend.CheckNotificationAuthorization()
}

func SendNotification(ctx context.Context, options NotificationOptions) error {
	appFrontend := getFrontend(ctx)
	return appFrontend.SendNotification(options)
}

func SendNotificationWithActions(ctx context.Context, options NotificationOptions) error {
	appFrontend := getFrontend(ctx)
	return appFrontend.SendNotificationWithActions(options)
}

func RegisterNotificationCategory(ctx context.Context, category NotificationCategory) error {
	appFrontend := getFrontend(ctx)
	return appFrontend.RegisterNotificationCategory(category)
}

func RemoveNotificationCategory(ctx context.Context, categoryId string) error {
	appFrontend := getFrontend(ctx)
	return appFrontend.RemoveNotificationCategory(categoryId)
}

func RemoveAllPendingNotifications(ctx context.Context) error {
	appFrontend := getFrontend(ctx)
	return appFrontend.RemoveAllPendingNotifications()
}

func RemovePendingNotification(ctx context.Context, identifier string) error {
	appFrontend := getFrontend(ctx)
	return appFrontend.RemovePendingNotification(identifier)
}

func RemoveAllDeliveredNotifications(ctx context.Context) error {
	appFrontend := getFrontend(ctx)
	return appFrontend.RemoveAllDeliveredNotifications()
}

func RemoveDeliveredNotification(ctx context.Context, identifier string) error {
	appFrontend := getFrontend(ctx)
	return appFrontend.RemoveDeliveredNotification(identifier)
}

func RemoveNotification(ctx context.Context, identifier string) error {
	appFrontend := getFrontend(ctx)
	return appFrontend.RemoveNotification(identifier)
}

func OnNotificationResponse(ctx context.Context, callback func(result NotificationResult)) {
	appFrontend := getFrontend(ctx)
	appFrontend.OnNotificationResponse(callback)
}
