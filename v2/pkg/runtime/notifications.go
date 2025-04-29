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
	frontend := getFrontend(ctx)
	return frontend.InitializeNotifications()
}

func IsNotificationAvailable(ctx context.Context) bool {
	frontend := getFrontend(ctx)
	return frontend.IsNotificationAvailable()
}

func RequestNotificationAuthorization(ctx context.Context) (bool, error) {
	frontend := getFrontend(ctx)
	return frontend.RequestNotificationAuthorization()
}

func CheckNotificationAuthorization(ctx context.Context) (bool, error) {
	frontend := getFrontend(ctx)
	return frontend.CheckNotificationAuthorization()
}

func SendNotification(ctx context.Context, options frontend.NotificationOptions) error {
	frontend := getFrontend(ctx)
	return frontend.SendNotification(options)
}

func SendNotificationWithActions(ctx context.Context, options frontend.NotificationOptions) error {
	frontend := getFrontend(ctx)
	return frontend.SendNotificationWithActions(options)
}

func RegisterNotificationCategory(ctx context.Context, category frontend.NotificationCategory) error {
	frontend := getFrontend(ctx)
	return frontend.RegisterNotificationCategory(category)
}

func RemoveNotificationCategory(ctx context.Context, categoryId string) error {
	frontend := getFrontend(ctx)
	return frontend.RemoveNotificationCategory(categoryId)
}

func RemoveAllPendingNotifications(ctx context.Context) error {
	frontend := getFrontend(ctx)
	return frontend.RemoveAllPendingNotifications()
}

func RemovePendingNotification(ctx context.Context, identifier string) error {
	frontend := getFrontend(ctx)
	return frontend.RemovePendingNotification(identifier)
}

func RemoveAllDeliveredNotifications(ctx context.Context) error {
	frontend := getFrontend(ctx)
	return frontend.RemoveAllDeliveredNotifications()
}

func RemoveDeliveredNotification(ctx context.Context, identifier string) error {
	frontend := getFrontend(ctx)
	return frontend.RemoveDeliveredNotification(identifier)
}

func RemoveNotification(ctx context.Context, identifier string) error {
	frontend := getFrontend(ctx)
	return frontend.RemoveNotification(identifier)
}

func OnNotificationResponse(ctx context.Context, callback func(result frontend.NotificationResult)) {
	frontend := getFrontend(ctx)
	frontend.OnNotificationResponse(callback)
}
