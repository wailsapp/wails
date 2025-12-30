//go:build ios

package notifications

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type iosNotifier struct{}

// New creates a new Notifications Service.
// On iOS, this returns a stub implementation.
// iOS notification functionality will be implemented via native bridges.
func New() *NotificationService {
	notificationServiceOnce.Do(func() {
		impl := &iosNotifier{}

		NotificationService_ = &NotificationService{
			impl: impl,
		}
	})

	return NotificationService_
}

func (n *iosNotifier) Startup(ctx context.Context, options application.ServiceOptions) error {
	// iOS notification startup - implementation pending native bridge
	return nil
}

func (n *iosNotifier) Shutdown() error {
	// iOS notification shutdown - implementation pending native bridge
	return nil
}

func (n *iosNotifier) RequestNotificationAuthorization(callback func(bool, error)) {
	// iOS notification authorization would go here via native bridge
	if callback != nil {
		callback(true, nil)
	}
}

func (n *iosNotifier) CheckNotificationAuthorization(callback func(bool, error)) {
	// iOS notification authorization check would go here via native bridge
	if callback != nil {
		callback(true, nil)
	}
}

func (n *iosNotifier) SendNotification(options NotificationOptions) error {
	// iOS notification would go here via native bridge
	return nil
}

func (n *iosNotifier) SendNotificationWithActions(options NotificationOptions) error {
	// iOS notification with actions would go here via native bridge
	return nil
}

func (n *iosNotifier) RegisterNotificationCategory(category NotificationCategory) error {
	// iOS notification category registration would go here via native bridge
	return nil
}

func (n *iosNotifier) RemoveNotificationCategory(categoryID string) error {
	// iOS notification category removal would go here via native bridge
	return nil
}

func (n *iosNotifier) RemoveAllPendingNotifications() error {
	// iOS pending notifications removal would go here via native bridge
	return nil
}

func (n *iosNotifier) RemovePendingNotification(identifier string) error {
	// iOS pending notification removal would go here via native bridge
	return nil
}

func (n *iosNotifier) RemoveAllDeliveredNotifications() error {
	// iOS delivered notifications removal would go here via native bridge
	return nil
}

func (n *iosNotifier) RemoveDeliveredNotification(identifier string) error {
	// iOS delivered notification removal would go here via native bridge
	return nil
}

func (n *iosNotifier) RemoveNotification(identifier string) error {
	// iOS notification removal would go here via native bridge
	return nil
}