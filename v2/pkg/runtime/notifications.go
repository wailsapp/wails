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

// InitializeNotifications initializes the notification service for the application.
// This must be called before sending any notifications. On macOS, this also ensures
// the notification delegate is properly initialized.
func InitializeNotifications(ctx context.Context) error {
	fe := getFrontend(ctx)
	return fe.InitializeNotifications()
}

// CleanupNotifications cleans up notification resources and releases any held connections.
// This should be called when shutting down the application to properly release resources
// (primarily needed on Linux to close D-Bus connections).
func CleanupNotifications(ctx context.Context) {
	fe := getFrontend(ctx)
	fe.CleanupNotifications()
}

// IsNotificationAvailable checks if notifications are available on the current platform.
func IsNotificationAvailable(ctx context.Context) bool {
	fe := getFrontend(ctx)
	return fe.IsNotificationAvailable()
}

// RequestNotificationAuthorization requests notification authorization from the user.
// On macOS, this prompts the user to allow notifications. On other platforms, this
// always returns true. Returns true if authorization was granted, false otherwise.
func RequestNotificationAuthorization(ctx context.Context) (bool, error) {
	fe := getFrontend(ctx)
	return fe.RequestNotificationAuthorization()
}

// CheckNotificationAuthorization checks the current notification authorization status.
// On macOS, this checks if the app has notification permissions. On other platforms,
// this always returns true.
func CheckNotificationAuthorization(ctx context.Context) (bool, error) {
	fe := getFrontend(ctx)
	return fe.CheckNotificationAuthorization()
}

// SendNotification sends a basic notification with the given options.
// The notification will display with the provided title, subtitle (if supported),
// and body text.
func SendNotification(ctx context.Context, options NotificationOptions) error {
	fe := getFrontend(ctx)
	return fe.SendNotification(options)
}

// SendNotificationWithActions sends a notification with action buttons.
// A NotificationCategory must be registered first using RegisterNotificationCategory.
// The options.CategoryID must match a previously registered category ID.
// If the category is not found, a basic notification will be sent instead.
func SendNotificationWithActions(ctx context.Context, options NotificationOptions) error {
	fe := getFrontend(ctx)
	return fe.SendNotificationWithActions(options)
}

// RegisterNotificationCategory registers a notification category that can be used
// with SendNotificationWithActions. Categories define the action buttons and optional
// reply fields that will appear on notifications.
func RegisterNotificationCategory(ctx context.Context, category NotificationCategory) error {
	fe := getFrontend(ctx)
	return fe.RegisterNotificationCategory(category)
}

// RemoveNotificationCategory removes a previously registered notification category.
func RemoveNotificationCategory(ctx context.Context, categoryId string) error {
	fe := getFrontend(ctx)
	return fe.RemoveNotificationCategory(categoryId)
}

// RemoveAllPendingNotifications removes all pending notifications from the notification center.
// On Windows, this is a no-op as the platform manages notification lifecycle automatically.
func RemoveAllPendingNotifications(ctx context.Context) error {
	fe := getFrontend(ctx)
	return fe.RemoveAllPendingNotifications()
}

// RemovePendingNotification removes a specific pending notification by its identifier.
// On Windows, this is a no-op as the platform manages notification lifecycle automatically.
func RemovePendingNotification(ctx context.Context, identifier string) error {
	fe := getFrontend(ctx)
	return fe.RemovePendingNotification(identifier)
}

// RemoveAllDeliveredNotifications removes all delivered notifications from the notification center.
// On Windows, this is a no-op as the platform manages notification lifecycle automatically.
func RemoveAllDeliveredNotifications(ctx context.Context) error {
	fe := getFrontend(ctx)
	return fe.RemoveAllDeliveredNotifications()
}

// RemoveDeliveredNotification removes a specific delivered notification by its identifier.
// On Windows, this is a no-op as the platform manages notification lifecycle automatically.
func RemoveDeliveredNotification(ctx context.Context, identifier string) error {
	fe := getFrontend(ctx)
	return fe.RemoveDeliveredNotification(identifier)
}

// RemoveNotification removes a notification by its identifier.
// This is a convenience function that works across platforms. On macOS, use the
// more specific RemovePendingNotification or RemoveDeliveredNotification functions.
func RemoveNotification(ctx context.Context, identifier string) error {
	fe := getFrontend(ctx)
	return fe.RemoveNotification(identifier)
}

// OnNotificationResponse registers a callback function that will be invoked when
// a user interacts with a notification (e.g., clicks an action button or the notification itself).
// The callback receives a NotificationResult containing the response details or any errors.
func OnNotificationResponse(ctx context.Context, callback func(result NotificationResult)) {
	fe := getFrontend(ctx)
	fe.OnNotificationResponse(callback)
}
