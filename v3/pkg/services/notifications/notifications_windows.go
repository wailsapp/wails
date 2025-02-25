//go:build windows

package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"git.sr.ht/~jackmordaunt/go-toast/v2"
	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/sys/windows/registry"
)

var NotificationService *Service
var NotificationLock sync.RWMutex
var NotificationCategories = make(map[string]NotificationCategory)

const (
	dataSeparator = ":::"
)

func New() *Service {
	if NotificationService == nil {
		NotificationService = &Service{}
	}
	return NotificationService
}

// ServiceName returns the name of the service
func (ns *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/services/notifications"
}

// ServiceStartup is called when the service is loaded
// Sets an activation callback to emit an event when notifications are interacted with.
func (ns *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	// Need to get App Name and generate a UUID on first launch
	//
	// toast.SetAppData(toast.AppData{
	// 	AppID:         "Notifications",
	// 	GUID:          "{8F2E1A3D-C497-42B6-9E5D-72F8A169B051}",
	// 	IconPath:      "C:\\Users\\Zach\\Development\\notifications_demo\\build\\appicon.ico",
	// 	ActivationExe: "C:\\Users\\Zach\\Development\\notifications_demo\\bin\\Notifications.exe",
	// })
	//

	toast.SetActivationCallback(func(args string, data []toast.UserData) {
		actionIdentifier, userInfo := parseNotificationResponse(args)
		response := NotificationResponse{
			ActionIdentifier: actionIdentifier,
		}

		if userInfo != "" {
			jsonStr := strings.ReplaceAll(userInfo, "'", "\"")
			var userInfoMap map[string]interface{}
			if err := json.Unmarshal([]byte(jsonStr), &userInfoMap); err == nil {
				response.UserInfo = userInfoMap
			}
		}

		if userText, found := getUserText(data); found {
			response.UserText = userText
		}

		if NotificationService != nil {
			NotificationService.handleNotificationResponse(response)
		}
	})

	return loadCategoriesFromRegistry()
}

// ServiceShutdown is called when the service is unloaded
func (ns *Service) ServiceShutdown() error {
	return saveCategoriesToRegistry()
}

// CheckBundleIdentifier is a Windows stub that always returns true.
// (bundle identifiers are macOS-specific)
func CheckBundleIdentifier() bool {
	return true
}

// RequestUserNotificationAuthorization is a Windows stub that always returns true, nil.
// (user authorization is macOS-specific)
func (ns *Service) RequestUserNotificationAuthorization() (bool, error) {
	return true, nil
}

// CheckNotificationAuthorization is a Windows stub that always returns true.
// (user authorization is macOS-specific)
func (ns *Service) CheckNotificationAuthorization() bool {
	return true
}

// SendNotification sends a basic notification with a name, title, and body. All other options are ignored on Windows.
// (subtitle and category id are only available on macOS)
func (ns *Service) SendNotification(options NotificationOptions) error {
	n := toast.Notification{
		AppID:               options.ID,
		Title:               options.Title,
		Body:                options.Body,
		ActivationArguments: DefaultActionIdentifier,
		Audio:               toast.IM,
	}

	if options.Data != nil {
		jsonData, err := json.Marshal(options.Data)
		if err == nil {
			n.ActivationArguments = DefaultActionIdentifier + dataSeparator + strings.ReplaceAll(string(jsonData), "\"", "'")
		}
	}

	err := n.Push()
	if err != nil {
		return err
	}
	return nil
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
// A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
// If a NotificationCategory is not registered a basic notification will be sent.
// (subtitle and category id are only available on macOS)
func (ns *Service) SendNotificationWithActions(options NotificationOptions) error {
	NotificationLock.RLock()
	nCategory := NotificationCategories[options.CategoryID]
	NotificationLock.RUnlock()

	n := toast.Notification{
		AppID:               options.ID,
		Title:               options.Title,
		Body:                options.Body,
		ActivationArguments: DefaultActionIdentifier,
		Audio:               toast.IM,
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
			Content:   nCategory.ReplyButtonTitle,
			Arguments: "TEXT_REPLY",
			InputID:   "userText",
		})
	}

	if options.Data != nil {
		jsonData, err := json.Marshal(options.Data)
		if err == nil {
			n.ActivationArguments = DefaultActionIdentifier + dataSeparator + strings.ReplaceAll(string(jsonData), "\"", "'")
			for index := range n.Actions {
				n.Actions[index].Arguments = n.Actions[index].Arguments + dataSeparator + strings.ReplaceAll(string(jsonData), "\"", "'")
			}
		}
	}

	err := n.Push()
	if err != nil {
		return err
	}
	return nil
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
// Registering a category with the same name as a previously registered NotificationCategory will override it.
func (ns *Service) RegisterNotificationCategory(category NotificationCategory) error {
	NotificationLock.Lock()
	NotificationCategories[category.ID] = NotificationCategory{
		ID:               category.ID,
		Actions:          category.Actions,
		HasReplyField:    bool(category.HasReplyField),
		ReplyPlaceholder: category.ReplyPlaceholder,
		ReplyButtonTitle: category.ReplyButtonTitle,
	}
	NotificationLock.Unlock()

	return saveCategoriesToRegistry()
}

// RemoveNotificationCategory removes a previously registered NotificationCategory.
func (ns *Service) RemoveNotificationCategory(categoryId string) error {
	NotificationLock.Lock()
	delete(NotificationCategories, categoryId)
	NotificationLock.Unlock()

	return saveCategoriesToRegistry()
}

// RemoveAllPendingNotifications is a Windows stub that always returns nil.
// (macOS-specific)
func (ns *Service) RemoveAllPendingNotifications() error {
	return nil
}

// RemovePendingNotification is a Windows stub that always returns nil.
// (macOS-specific)
func (ns *Service) RemovePendingNotification(_ string) error {
	return nil
}

// RemoveAllDeliveredNotifications is a Windows stub that always returns nil.
// (macOS-specific)
func (ns *Service) RemoveAllDeliveredNotifications() error {
	return nil
}

// RemoveDeliveredNotification is a Windows stub that always returns nil.
// (macOS-specific)
func (ns *Service) RemoveDeliveredNotification(_ string) error {
	return nil
}

func parseNotificationResponse(response string) (action string, data string) {
	parts := strings.Split(response, dataSeparator)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

// Is there a better way for me to grab this from the Wails config?
func getExeName() string {
	executable, err := os.Executable()
	if err != nil {
		return ""
	}

	return strings.TrimSuffix(filepath.Base(executable), filepath.Ext(executable))
}

func saveCategoriesToRegistry() error {
	appName := getExeName()
	if appName == "" {
		return fmt.Errorf("failed to save categories to registry: empty executable name")
	}
	registryPath := fmt.Sprintf(`SOFTWARE\%s\NotificationCategories`, appName)

	key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		registryPath,
		registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	NotificationLock.RLock()
	data, err := json.Marshal(NotificationCategories)
	NotificationLock.RUnlock()
	if err != nil {
		return err
	}

	return key.SetStringValue("Categories", string(data))
}

func loadCategoriesFromRegistry() error {
	appName := getExeName()
	if appName == "" {
		return fmt.Errorf("failed to save categories to registry: empty executable name")
	}
	registryPath := fmt.Sprintf(`SOFTWARE\%s\NotificationCategories`, appName)

	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		registryPath,
		registry.QUERY_VALUE,
	)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil
		}
		return err
	}
	defer key.Close()

	data, _, err := key.GetStringValue("Categories")
	if err != nil {
		return err
	}

	categories := make(map[string]NotificationCategory)
	if err := json.Unmarshal([]byte(data), &categories); err != nil {
		return err
	}

	NotificationLock.Lock()
	NotificationCategories = categories
	NotificationLock.Unlock()

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
