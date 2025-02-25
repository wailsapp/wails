//go:build linux

package notifications

import (
	"context"

	"git.sr.ht/~whereswaldon/shout"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var NotificationService *Service
var Notifier shout.Notifier

// Creates a new Notifications Service.
func New() *Service {
	if NotificationService == nil {
		NotificationService = &Service{}
	}
	return NotificationService
}

// ServiceStartup is called when the service is loaded
func (ns *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	return nil
}

// ServiceShutdown is called when the service is unloaded
func (ns *Service) ServiceShutdown() error {
	return nil
}

// CheckBundleIdentifier is a Linux stub that always returns true.
// (bundle identifiers are macOS-specific)
func CheckBundleIdentifier() bool {
	return true
}

// RequestUserNotificationAuthorization is a Linux stub that always returns true, nil.
// (user authorization is macOS-specific)
func (ns *Service) RequestUserNotificationAuthorization() (bool, error) {
	return true, nil
}

// CheckNotificationAuthorization is a Linux stub that always returns true.
// (user authorization is macOS-specific)
func (ns *Service) CheckNotificationAuthorization() bool {
	return true
}

// SendNotification sends a basic notification with a name, title, and body. All other options are ignored on Linux.
// (subtitle and category id are only available on macOS)
func (ns *Service) SendNotification(options NotificationOptions) error {
	return nil
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
// A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
// If a NotificationCategory is not registered a basic notification will be sent.
// (subtitle and category id are only available on macOS)
func (ns *Service) SendNotificationWithActions(options NotificationOptions) error {
	return nil
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
func (ns *Service) RegisterNotificationCategory(category NotificationCategory) error {
	return nil
}

// RemoveNotificationCategory removes a previously registered NotificationCategory.
func (ns *Service) RemoveNotificationCategory(categoryId string) error {
	return nil
}

// RemoveAllPendingNotifications is a Linux stub that always returns nil.
// (macOS-specific)
func (ns *Service) RemoveAllPendingNotifications() error {
	return nil
}

// RemovePendingNotification is a Linux stub that always returns nil.
// (macOS-specific)
func (ns *Service) RemovePendingNotification(_ string) error {
	return nil
}

// RemoveAllDeliveredNotifications is a Linux stub that always returns nil.
// (macOS-specific)
func (ns *Service) RemoveAllDeliveredNotifications() error {
	return nil
}

// RemoveDeliveredNotification is a Linux stub that always returns nil.
// (macOS-specific)
func (ns *Service) RemoveDeliveredNotification(_ string) error {
	return nil
}

// RemoveNotification removes a pending notification matching the unique identifier.
// (Linux-specific)
func (ns *Service) RemoveNotification(identifier string) error {
	return nil
}
