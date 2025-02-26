//go:build linux

package notifications

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"git.sr.ht/~whereswaldon/shout"
	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var NotificationLock sync.RWMutex
var NotificationCategories = make(map[string]NotificationCategory)
var Notifier shout.Notifier
var appName = application.Get().Config().Name

// Creates a new Notifications Service.
func New() *Service {
	if NotificationService == nil {
		NotificationService = &Service{}
	}
	return NotificationService
}

// ServiceStartup is called when the service is loaded
func (ns *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	if err := loadCategories(); err != nil {
		fmt.Printf("Failed to load notification categories: %v\n", err)
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		return fmt.Errorf("failed to connect to D-Bus session bus: %v", err)
	}

	var iconPath string

	Notifier, err = shout.NewNotifier(conn, appName, iconPath, func(notificationID, action string, platformData map[string]dbus.Variant, target, notifierResponse dbus.Variant, err error) {
		if err != nil {
			return
		}

		response := NotificationResponse{
			ID:               notificationID,
			ActionIdentifier: action,
		}

		if target.Signature().String() == "s" {
			var targetStr string
			if err := target.Store(&targetStr); err == nil {
				var userInfo map[string]interface{}
				userInfoStr, err := base64.StdEncoding.DecodeString(targetStr)
				if err != nil {
					if err := json.Unmarshal([]byte(targetStr), &userInfo); err == nil {
						response.UserInfo = userInfo
					}
				}
				if err := json.Unmarshal(userInfoStr, &userInfo); err == nil {
					response.UserInfo = userInfo
				}
			}
		}

		// if notifierResponse.Signature().String() == "s" {
		// 	var userText string
		// 	if err := notifierResponse.Store(&userText); err == nil {
		// 		response.UserText = userText
		// 	}
		// }

		if NotificationService != nil {
			NotificationService.handleNotificationResponse(response)
		}
	})

	if err != nil {
		return fmt.Errorf("failed to create notifier: %v", err)
	}

	return nil
}

// ServiceShutdown is called when the service is unloaded
func (ns *Service) ServiceShutdown() error {
	return saveCategories()
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
func (ns *Service) CheckNotificationAuthorization() (bool, error) {
	return true, nil
}

// SendNotification sends a basic notification with a unique identifier, title, subtitle, and body.
func (ns *Service) SendNotification(options NotificationOptions) error {
	notification := shout.Notification{
		Title:         options.Title,
		Body:          options.Body,
		Priority:      shout.Normal,
		DefaultAction: DefaultActionIdentifier,
	}

	if options.Data != nil {
		jsonData, err := json.Marshal(options.Data)
		if err == nil {
			notification.DefaultActionTarget = dbus.MakeVariant(base64.StdEncoding.EncodeToString(jsonData))
		}
	}

	return Notifier.Send(options.ID, notification)
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
func (ns *Service) SendNotificationWithActions(options NotificationOptions) error {
	NotificationLock.RLock()
	category, exists := NotificationCategories[options.CategoryID]
	NotificationLock.RUnlock()

	if !exists {
		return ns.SendNotification(options)
	}

	notification := shout.Notification{
		Title:         options.Title,
		Body:          options.Body,
		Priority:      shout.Normal,
		DefaultAction: DefaultActionIdentifier,
	}

	if options.Data != nil {
		jsonData, err := json.Marshal(options.Data)
		if err == nil {
			notification.DefaultActionTarget = dbus.MakeVariant(base64.StdEncoding.EncodeToString(jsonData))
		}
	}

	for _, action := range category.Actions {
		notification.Buttons = append(notification.Buttons, shout.Button{
			Label:  action.Title,
			Action: action.ID,
			Target: "", // Will be set below if we have user data
		})
	}

	if options.Data != nil {
		jsonData, err := json.Marshal(options.Data)
		if err == nil {
			for index := range notification.Buttons {
				notification.Buttons[index].Target = string(jsonData)
			}
		}
	}

	return Notifier.Send(options.ID, notification)
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
func (ns *Service) RegisterNotificationCategory(category NotificationCategory) error {
	NotificationLock.Lock()
	NotificationCategories[category.ID] = category
	NotificationLock.Unlock()
	return saveCategories()
}

// RemoveNotificationCategory removes a previously registered NotificationCategory.
func (ns *Service) RemoveNotificationCategory(categoryId string) error {
	NotificationLock.Lock()
	delete(NotificationCategories, categoryId)
	NotificationLock.Unlock()
	return saveCategories()
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

// RemoveNotification removes a notification by ID (Linux-specific)
func (ns *Service) RemoveNotification(identifier string) error {
	return Notifier.Remove(identifier)
}

// getConfigFilePath returns the path to the configuration file for storing notification categories
func getConfigFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %v", err)
	}

	appConfigDir := filepath.Join(configDir, appName)
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %v", err)
	}

	return filepath.Join(appConfigDir, "notification-categories.json"), nil
}

// saveCategories saves the notification categories to a file.
func saveCategories() error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	NotificationLock.RLock()
	data, err := json.Marshal(NotificationCategories)
	NotificationLock.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal notification categories: %v", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write notification categories to file: %v", err)
	}

	return nil
}

// loadCategories loads notification categories from a file.
func loadCategories() error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read notification categories file: %v", err)
	}

	if len(data) == 0 {
		return nil
	}

	categories := make(map[string]NotificationCategory)
	if err := json.Unmarshal(data, &categories); err != nil {
		return fmt.Errorf("failed to unmarshal notification categories: %v", err)
	}

	NotificationLock.Lock()
	NotificationCategories = categories
	NotificationLock.Unlock()

	return nil
}
