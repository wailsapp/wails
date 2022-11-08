package runtime

import (
	"context"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

// LinuxNotificationSound represents a sound to be played when a notification shows up
type LinuxNotificationSound = frontend.LinuxNotificationSound

// LinuxNotificationAction represents a clickable action on a notification
type LinuxNotificationAction = frontend.LinuxNotificationAction

// LinuxNotificationOptions contains the options for the linux specific notification options
type LinuxNotificationOptions = frontend.LinuxNotificationOptions

// WindowsNotificationAction represents a Notification action for a notification
type WindowsNotificationAction = frontend.WindowsNotificationAction

// WindowsNotificationOptions contains the options for the Windows specific notification options
type WindowsNotificationOptions = frontend.WindowsNotificationOptions

// MacNotificationAction represents a Notification action for a notification
type MacNotificationAction = frontend.MacNotificationAction

// MacNotificationOptions contains the options for the MacOs specific notification options
type MacNotificationOptions = frontend.MacNotificationOptions

// NotificationOptions contains the options for desktop notification options
type NotificationOptions = frontend.NotificationOptions

// SendNotification creates and sends a desktop notification
func SendNotification(ctx context.Context, notificationOptions NotificationOptions) error {
	appFrontend := getFrontend(ctx)

	if notificationOptions.AppID == "" {
		notificationOptions.AppID = appFrontend.AppID()
	}

	if notificationOptions.LinuxOptions == nil {
		notificationOptions.LinuxOptions = &LinuxNotificationOptions{
			Urgency: -1,
		}
	}

	if notificationOptions.WindowsOptions == nil {
		notificationOptions.WindowsOptions = &WindowsNotificationOptions{}
	}

	if notificationOptions.MacOptions == nil {
		notificationOptions.MacOptions = &MacNotificationOptions{}
	}

	return appFrontend.SendNotification(notificationOptions)
}
