//go:build windows

package notifications

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"git.sr.ht/~jackmordaunt/go-toast/v2"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/sys/windows/registry"
)

var NotificationLock sync.RWMutex
var NotificationCategories = make(map[string]NotificationCategory)

// NotificationPayload combines the action ID and user data into a single structure
type NotificationPayload struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

// Creates a new Notifications Service.
func New() *Service {
	if NotificationService == nil {
		NotificationService = &Service{}
	}
	return NotificationService
}

// ServiceStartup is called when the service is loaded
// Sets an activation callback to emit an event when notifications are interacted with.
func (ns *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	appName := application.Get().Config().Name

	guid, err := getGUID(appName)
	if err != nil {
		return err
	}

	toast.SetAppData(toast.AppData{
		AppID:    appName,
		GUID:     guid,
		IconPath: filepath.Join(os.TempDir(), appName+guid+".png"),
	})

	toast.SetActivationCallback(func(args string, data []toast.UserData) {
		actionIdentifier, userInfo := parseNotificationResponse(args)
		response := NotificationResponse{
			ActionIdentifier: actionIdentifier,
		}

		if userInfo != "" {
			var userInfoMap map[string]interface{}
			if err := json.Unmarshal([]byte(userInfo), &userInfoMap); err == nil {
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
	if err := saveIconToDir(); err != nil {
		fmt.Printf("Error saving icon: %v\n", err)
	}

	n := toast.Notification{
		Title:               options.Title,
		Body:                options.Body,
		ActivationArguments: DefaultActionIdentifier,
	}

	if options.Data != nil {
		encodedPayload, err := encodePayload(DefaultActionIdentifier, options.Data)
		if err == nil {
			n.ActivationArguments = encodedPayload
		}
	}

	return n.Push()
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
// A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
// If a NotificationCategory is not registered a basic notification will be sent.
// (subtitle and category id are only available on macOS)
func (ns *Service) SendNotificationWithActions(options NotificationOptions) error {
	if err := saveIconToDir(); err != nil {
		fmt.Printf("Error saving icon: %v\n", err)
	}

	NotificationLock.RLock()
	nCategory := NotificationCategories[options.CategoryID]
	NotificationLock.RUnlock()

	n := toast.Notification{
		Title:               options.Title,
		Body:                options.Body,
		ActivationArguments: DefaultActionIdentifier,
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
			Placeholder: nCategory.ReplyPlaceholder,
		})

		n.Actions = append(n.Actions, toast.Action{
			Content:   nCategory.ReplyButtonTitle,
			Arguments: "TEXT_REPLY",
			InputID:   "userText",
		})
	}

	if options.Data != nil {
		n.ActivationArguments, _ = encodePayload(n.ActivationArguments, options.Data)

		for index := range n.Actions {
			n.Actions[index].Arguments, _ = encodePayload(n.Actions[index].Arguments, options.Data)
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

// RemoveNotification is a Windows stub that always returns nil.
// (Linux-specific)
func (ns *Service) RemoveNotification(identifier string) error {
	return nil
}

// encodePayload combines an action ID and user data into a single encoded string
func encodePayload(actionID string, data map[string]interface{}) (string, error) {
	payload := NotificationPayload{
		Action: actionID,
		Data:   data,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return actionID, err
	}

	encodedPayload := base64.StdEncoding.EncodeToString(jsonData)
	return encodedPayload, nil
}

// decodePayload extracts the action ID and user data from an encoded payload
func decodePayload(encodedString string) (string, map[string]interface{}, error) {
	jsonData, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return encodedString, nil, nil
	}

	var payload NotificationPayload
	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return encodedString, nil, nil
	}

	return payload.Action, payload.Data, nil
}

// parseNotificationResponse updated to use structured payload decoding
func parseNotificationResponse(response string) (action string, data string) {
	actionID, userData, _ := decodePayload(response)

	if userData != nil {
		userDataJSON, err := json.Marshal(userData)
		if err == nil {
			return actionID, string(userDataJSON)
		}
	}

	return actionID, ""
}

func saveIconToDir() error {
	options := application.Get().Config()
	appName := options.Name
	icon := options.Icon

	if len(icon) == 0 {
		return fmt.Errorf("failed to retrieve icon from application")
	}

	guid, err := getGUID(appName)
	if err != nil {
		return fmt.Errorf("failed to retrieve application guid from registry")
	}

	iconPath := filepath.Join(os.TempDir(), appName+guid+".png")

	return os.WriteFile(iconPath, icon, 0644)
}

func saveCategoriesToRegistry() error {
	appName := application.Get().Config().Name
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
	appName := application.Get().Config().Name
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

func getGUID(name string) (string, error) {
	keyPath := `Software\Classes\AppUserModelId\` + name

	k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
	if err == nil {
		guid, _, err := k.GetStringValue("CustomActivator")
		k.Close()
		if err == nil && guid != "" {
			return guid, nil
		}
	}

	guid := generateGUID()

	k, _, err = registry.CreateKey(registry.CURRENT_USER, keyPath, registry.WRITE)
	if err != nil {
		return "", fmt.Errorf("failed to create registry key: %w", err)
	}
	defer k.Close()

	if err := k.SetStringValue("CustomActivator", guid); err != nil {
		return "", fmt.Errorf("failed to write GUID to registry: %w", err)
	}

	return guid, nil
}

func generateGUID() string {
	guid := uuid.New()
	return fmt.Sprintf("{%s}", guid.String())
}
