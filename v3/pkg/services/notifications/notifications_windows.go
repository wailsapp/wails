//go:build windows

package notifications

import (
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"git.sr.ht/~jackmordaunt/go-toast/v2"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"golang.org/x/sys/windows/registry"
)

var (
	NotificationCategories     = make(map[string]NotificationCategory)
	notificationCategoriesLock sync.RWMutex
	appName                    string
	appGUID                    string
	iconPath                   string
)

const (
	ToastRegistryPath                  = `Software\Classes\AppUserModelId\`
	ToastRegistryGuidKey               = "CustomActivator"
	NotificationCategoriesRegistryPath = `SOFTWARE\%s\NotificationCategories`
	NotificationCategoriesRegistryKey  = "Categories"
)

// NotificationPayload combines the action ID and user data into a single structure
type NotificationPayload struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

// Creates a new Notifications Service.
func New() *Service {
	notificationServiceOnce.Do(func() {
		if NotificationService == nil {
			NotificationService = &Service{}
		}
	})

	return NotificationService
}

// ServiceStartup is called when the service is loaded
// Sets an activation callback to emit an event when notifications are interacted with.
func (ns *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	appName = application.Get().Config().Name

	guid, err := getGUID()
	if err != nil {
		return err
	}
	appGUID = guid

	iconPath = filepath.Join(os.TempDir(), appName+appGUID+".png")

	toast.SetAppData(toast.AppData{
		AppID:    appName,
		GUID:     guid,
		IconPath: iconPath,
	})

	toast.SetActivationCallback(func(args string, data []toast.UserData) {
		result := NotificationResult{}
		actionIdentifier, userInfo := parseNotificationResponse(args)
		response := NotificationResponse{
			ActionIdentifier: actionIdentifier,
		}

		if userInfo != "" {
			var userInfoMap map[string]interface{}
			if err := json.Unmarshal([]byte(userInfo), &userInfoMap); err != nil {
				result.Error = fmt.Errorf("failed to unmarshal notification response: %w", err)

				if ns := getNotificationService(); ns != nil {
					ns.handleNotificationResult(result)
				}
			}
			response.UserInfo = userInfoMap
		}

		if userText, found := getUserText(data); found {
			response.UserText = userText
		}

		result.Response = response
		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(result)
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

// RequestNotificationAuthorization is a Windows stub that always returns true, nil.
// (user authorization is macOS-specific)
func (ns *Service) RequestNotificationAuthorization() (bool, error) {
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
	if err := validateNotificationOptions(options); err != nil {
		return err
	}

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
		if err != nil {
			return fmt.Errorf("failed to encode notification data: %w", err)
		}
		n.ActivationArguments = encodedPayload
	}

	return n.Push()
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
// A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
// If a NotificationCategory is not registered a basic notification will be sent.
// (subtitle and category id are only available on macOS)
func (ns *Service) SendNotificationWithActions(options NotificationOptions) error {
	if err := validateNotificationOptions(options); err != nil {
		return err
	}

	if err := saveIconToDir(); err != nil {
		fmt.Printf("Error saving icon: %v\n", err)
	}

	notificationCategoriesLock.RLock()
	nCategory := NotificationCategories[options.CategoryID]
	notificationCategoriesLock.RUnlock()

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
			encodedPayload, err := encodePayload(n.Actions[index].Arguments, options.Data)
			if err != nil {
				return fmt.Errorf("failed to encode notification data: %w", err)
			}
			n.Actions[index].Arguments = encodedPayload
		}
	}

	return n.Push()
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
// Registering a category with the same name as a previously registered NotificationCategory will override it.
func (ns *Service) RegisterNotificationCategory(category NotificationCategory) error {
	notificationCategoriesLock.Lock()
	NotificationCategories[category.ID] = NotificationCategory{
		ID:               category.ID,
		Actions:          category.Actions,
		HasReplyField:    bool(category.HasReplyField),
		ReplyPlaceholder: category.ReplyPlaceholder,
		ReplyButtonTitle: category.ReplyButtonTitle,
	}
	notificationCategoriesLock.Unlock()

	return saveCategoriesToRegistry()
}

// RemoveNotificationCategory removes a previously registered NotificationCategory.
func (ns *Service) RemoveNotificationCategory(categoryId string) error {
	notificationCategoriesLock.Lock()
	delete(NotificationCategories, categoryId)
	notificationCategoriesLock.Unlock()

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
	icon, err := application.NewIconFromResource(w32.GetModuleHandle(""), uint16(3))
	if err != nil {
		return fmt.Errorf("failed to retrieve application icon: %w", err)
	}

	return saveHIconAsPNG(icon, iconPath)
}

func saveCategoriesToRegistry() error {
	registryPath := fmt.Sprintf(NotificationCategoriesRegistryPath, appName)

	key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		registryPath,
		registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	notificationCategoriesLock.RLock()
	data, err := json.Marshal(NotificationCategories)
	notificationCategoriesLock.RUnlock()
	if err != nil {
		return err
	}

	return key.SetStringValue(NotificationCategoriesRegistryKey, string(data))
}

func loadCategoriesFromRegistry() error {
	registryPath := fmt.Sprintf(NotificationCategoriesRegistryPath, appName)

	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		registryPath,
		registry.QUERY_VALUE,
	)
	if err != nil {
		if err == registry.ErrNotExist {
			// Not an error, no saved categories
			return nil
		}
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	data, _, err := key.GetStringValue(NotificationCategoriesRegistryKey)
	if err != nil {
		if err == registry.ErrNotExist {
			// No value yet, but key exists
			return nil
		}
		return fmt.Errorf("failed to read categories from registry: %w", err)
	}

	categories := make(map[string]NotificationCategory)
	if err := json.Unmarshal([]byte(data), &categories); err != nil {
		return fmt.Errorf("failed to parse notification categories from registry: %w", err)
	}

	notificationCategoriesLock.Lock()
	NotificationCategories = categories
	notificationCategoriesLock.Unlock()

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

func getGUID() (string, error) {
	keyPath := ToastRegistryPath + appName

	k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
	if err == nil {
		guid, _, err := k.GetStringValue(ToastRegistryGuidKey)
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

	if err := k.SetStringValue(ToastRegistryGuidKey, guid); err != nil {
		return "", fmt.Errorf("failed to write GUID to registry: %w", err)
	}

	return guid, nil
}

func generateGUID() string {
	guid := uuid.New()
	return fmt.Sprintf("{%s}", guid.String())
}
