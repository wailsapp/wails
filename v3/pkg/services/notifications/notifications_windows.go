//go:build windows

package notifications

import (
	"context"

	"git.sr.ht/~jackmordaunt/go-toast"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var NotificationCategories map[string]NotificationCategory = make(map[string]NotificationCategory)

func New() *Service {
	return &Service{}
}

// ServiceName returns the name of the service
func (ns *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/services/notifications"
}

// ServiceStartup is called when the service is loaded
func (ns *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	toast.SetActivationCallback(func(args string, data []toast.UserData) {
		response := NotificationResponse{
			Name: "notification",
			Data: NotificationResponseData{
				ActionIdentifier: args,
			},
		}

		if userText, found := getUserText(data); found {
			response.Data.UserText = userText
		}

		application.Get().EmitEvent("notificationResponse", response)
	})
	return nil
}

// ServiceShutdown is called when the service is unloaded
func (ns *Service) ServiceShutdown() error {
	return nil
}

// On Windows this does not apply, return true
func CheckBundleIdentifier() bool {
	return true
}

// On Windows this does not apply, return true
func (ns *Service) RequestUserNotificationAuthorization() (bool, error) {
	return true, nil
}

// On Windows this does not apply, return true
func (ns *Service) CheckNotificationAuthorization() bool {
	return true
}

func (ns *Service) SendNotification(identifier, title, _, body string) error {
	n := toast.Notification{
		AppID: identifier,
		Title: title,
		Body:  body,
	}

	err := n.Push()
	if err != nil {
		return err
	}
	return nil
}

func (ns *Service) SendNotificationWithActions(options NotificationOptions) error {
	nCategory := NotificationCategories[options.CategoryID]

	n := toast.Notification{
		AppID: options.ID,
		Title: options.Title,
		Body:  options.Body,
	}

	for _, action := range nCategory.Actions {
		n.Actions = append(n.Actions, toast.Action{
			Content:   action.Title,
			Arguments: action.ID,
		})
	}

	if nCategory.HasReplyField {
		n.Inputs = append(n.Inputs, toast.Input{
			ID:          "userText",
			Title:       nCategory.ReplyButtonTitle,
			Placeholder: nCategory.ReplyPlaceholder,
		})

		n.Actions = append(n.Actions, toast.Action{
			Content: nCategory.ReplyButtonTitle,
		})
	}

	err := n.Push()
	if err != nil {
		return err
	}
	return nil
}

func (ns *Service) RegisterNotificationCategory(category NotificationCategory) error {
	NotificationCategories[category.ID] = NotificationCategory{
		ID:               category.ID,
		Actions:          category.Actions,
		HasReplyField:    bool(category.HasReplyField),
		ReplyPlaceholder: category.ReplyPlaceholder,
		ReplyButtonTitle: category.ReplyButtonTitle,
	}

	return nil
}

// RemoveAllPendingNotifications removes all pending notifications
func (ns *Service) RemoveAllPendingNotifications() error {
	return nil
}

// RemovePendingNotification removes a specific pending notification
func (ns *Service) RemovePendingNotification(_ string) error {
	return nil
}

// RemoveAllDeliveredNotifications removes all delivered notifications
func (ns *Service) RemoveAllDeliveredNotifications() error {
	return nil
}

// RemoveDeliveredNotification removes a specific delivered notification
func (ns *Service) RemoveDeliveredNotification(_ string) error {
	return nil
}

func getUserText(data []toast.UserData) (string, bool) {
	for _, d := range data {
		if d.Key == "userText" {
			return d.Value, true
		}
	}
	return "", false
}
