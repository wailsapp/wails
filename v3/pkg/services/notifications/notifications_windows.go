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

type windowsNotifier struct {
	categories     map[string]NotificationCategory
	categoriesLock sync.RWMutex
	appName        string
	appGUID        string
	iconPath       string
}

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
		impl := &windowsNotifier{
			categories: make(map[string]NotificationCategory),
		}

		NotificationService = &Service{
			impl: impl,
		}
	})

	return NotificationService
}

// Startup is called when the service is loaded
// Sets an activation callback to emit an event when notifications are interacted with.
func (wn *windowsNotifier) Startup(ctx context.Context) error {
	wn.appName = application.Get().Config().Name

	guid, err := wn.getGUID()
	if err != nil {
		return err
	}
	wn.appGUID = guid

	wn.iconPath = filepath.Join(os.TempDir(), wn.appName+wn.appGUID+".png")

	toast.SetAppData(toast.AppData{
		AppID:    wn.appName,
		GUID:     guid,
		IconPath: wn.iconPath,
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

		if userText, found := wn.getUserText(data); found {
			response.UserText = userText
		}

		result.Response = response
		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(result)
		}
	})

	return wn.loadCategoriesFromRegistry()
}

// Shutdown will attempt to save the categories to the registry when the service unloads
func (wn *windowsNotifier) Shutdown() error {
	return wn.saveCategoriesToRegistry()
}

// RequestNotificationAuthorization is a Windows stub that always returns true, nil.
// (user authorization is macOS-specific)
func (wn *windowsNotifier) RequestNotificationAuthorization() (bool, error) {
	return true, nil
}

// CheckNotificationAuthorization is a Windows stub that always returns true.
// (user authorization is macOS-specific)
func (wn *windowsNotifier) CheckNotificationAuthorization() (bool, error) {
	return true, nil
}

// SendNotification sends a basic notification with a name, title, and body. All other options are ignored on Windows.
// (subtitle is only available on macOS)
func (wn *windowsNotifier) SendNotification(options NotificationOptions) error {
	if err := validateNotificationOptions(options); err != nil {
		return err
	}

	if err := wn.saveIconToDir(); err != nil {
		fmt.Printf("Error saving icon: %v\n", err)
	}

	n := toast.Notification{
		Title:               options.Title,
		Body:                options.Body,
		ActivationArguments: DefaultActionIdentifier,
	}

	if options.Data != nil {
		encodedPayload, err := wn.encodePayload(DefaultActionIdentifier, options.Data)
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
// (subtitle is only available on macOS)
func (wn *windowsNotifier) SendNotificationWithActions(options NotificationOptions) error {
	if err := validateNotificationOptions(options); err != nil {
		return err
	}

	if err := wn.saveIconToDir(); err != nil {
		fmt.Printf("Error saving icon: %v\n", err)
	}

	wn.categoriesLock.RLock()
	nCategory := wn.categories[options.CategoryID]
	wn.categoriesLock.RUnlock()

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
		encodedPayload, err := wn.encodePayload(n.ActivationArguments, options.Data)
		if err != nil {
			return fmt.Errorf("failed to encode notification data: %w", err)
		}
		n.ActivationArguments = encodedPayload

		for index := range n.Actions {
			encodedPayload, err := wn.encodePayload(n.Actions[index].Arguments, options.Data)
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
func (wn *windowsNotifier) RegisterNotificationCategory(category NotificationCategory) error {
	wn.categoriesLock.Lock()
	defer wn.categoriesLock.Unlock()

	wn.categories[category.ID] = NotificationCategory{
		ID:               category.ID,
		Actions:          category.Actions,
		HasReplyField:    bool(category.HasReplyField),
		ReplyPlaceholder: category.ReplyPlaceholder,
		ReplyButtonTitle: category.ReplyButtonTitle,
	}

	return wn.saveCategoriesToRegistry()
}

// RemoveNotificationCategory removes a previously registered NotificationCategory.
func (wn *windowsNotifier) RemoveNotificationCategory(categoryId string) error {
	wn.categoriesLock.Lock()
	defer wn.categoriesLock.Unlock()

	delete(wn.categories, categoryId)

	return wn.saveCategoriesToRegistry()
}

// RemoveAllPendingNotifications is a Windows stub that always returns nil.
// (macOS-specific)
func (wn *windowsNotifier) RemoveAllPendingNotifications() error {
	return nil
}

// RemovePendingNotification is a Windows stub that always returns nil.
// (macOS-specific)
func (wn *windowsNotifier) RemovePendingNotification(_ string) error {
	return nil
}

// RemoveAllDeliveredNotifications is a Windows stub that always returns nil.
// (macOS-specific)
func (wn *windowsNotifier) RemoveAllDeliveredNotifications() error {
	return nil
}

// RemoveDeliveredNotification is a Windows stub that always returns nil.
// (macOS-specific)
func (wn *windowsNotifier) RemoveDeliveredNotification(_ string) error {
	return nil
}

// RemoveNotification is a Windows stub that always returns nil.
// (Linux-specific)
func (wn *windowsNotifier) RemoveNotification(identifier string) error {
	return nil
}

// encodePayload combines an action ID and user data into a single encoded string
func (wn *windowsNotifier) encodePayload(actionID string, data map[string]interface{}) (string, error) {
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

func (wn *windowsNotifier) saveIconToDir() error {
	icon, err := application.NewIconFromResource(w32.GetModuleHandle(""), uint16(3))
	if err != nil {
		return fmt.Errorf("failed to retrieve application icon: %w", err)
	}

	return saveHIconAsPNG(icon, wn.iconPath)
}

func (wn *windowsNotifier) saveCategoriesToRegistry() error {
	wn.categoriesLock.Lock()
	defer wn.categoriesLock.Unlock()

	registryPath := fmt.Sprintf(NotificationCategoriesRegistryPath, wn.appName)

	key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		registryPath,
		registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	data, err := json.Marshal(wn.categories)
	if err != nil {
		return err
	}

	return key.SetStringValue(NotificationCategoriesRegistryKey, string(data))
}

func (wn *windowsNotifier) loadCategoriesFromRegistry() error {
	wn.categoriesLock.Lock()
	defer wn.categoriesLock.Unlock()

	registryPath := fmt.Sprintf(NotificationCategoriesRegistryPath, wn.appName)

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

	wn.categories = categories

	return nil
}

func (wn *windowsNotifier) getUserText(data []toast.UserData) (string, bool) {
	for _, d := range data {
		if d.Key == "userText" {
			return d.Value, true
		}
	}
	return "", false
}

func (wn *windowsNotifier) getGUID() (string, error) {
	keyPath := ToastRegistryPath + wn.appName

	k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
	if err == nil {
		guid, _, err := k.GetStringValue(ToastRegistryGuidKey)
		k.Close()
		if err == nil && guid != "" {
			return guid, nil
		}
	}

	guid := wn.generateGUID()

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

func (wn *windowsNotifier) generateGUID() string {
	guid := uuid.New()
	return fmt.Sprintf("{%s}", guid.String())
}
