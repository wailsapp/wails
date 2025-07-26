//go:build windows
// +build windows

package windows

import (
	"encoding/base64"
	"encoding/json"
	"sync"

	wintoast "git.sr.ht/~jackmordaunt/go-toast/v2/wintoast"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc/w32"

	"fmt"
	"os"
	"path/filepath"
	_ "unsafe" // for go:linkname

	"git.sr.ht/~jackmordaunt/go-toast/v2"
	"golang.org/x/sys/windows/registry"
)

var (
	categories     map[string]frontend.NotificationCategory
	categoriesLock sync.RWMutex
	appName        string
	appGUID        string
	iconPath       string = ""
	exePath        string

	notificationResultCallback func(result frontend.NotificationResult)
	callbackLock               sync.RWMutex
)

const DefaultActionIdentifier = "DEFAULT_ACTION"

const (
	ToastRegistryPath                  = `Software\Classes\AppUserModelId\`
	ToastRegistryGuidKey               = "CustomActivator"
	NotificationCategoriesRegistryPath = `SOFTWARE\%s\NotificationCategories`
	NotificationCategoriesRegistryKey  = "Categories"
)

// NotificationPayload combines the action ID and user data into a single structure
type NotificationPayload struct {
	Action  string                       `json:"action"`
	Options frontend.NotificationOptions `json:"payload,omitempty"`
}

func (f *Frontend) InitializeNotifications() error {
	categories = make(map[string]frontend.NotificationCategory)
	categoriesLock.Lock()
	defer categoriesLock.Unlock()

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable: %w", err)
	}
	exePath = exe
	appName = filepath.Base(exePath)

	appGUID, err = getGUID()
	if err != nil {
		return err
	}

	iconPath = filepath.Join(os.TempDir(), appName+appGUID+".png")

	// Create the registry key for the toast activator
	key, _, err := registry.CreateKey(registry.CURRENT_USER,
		`Software\Classes\CLSID\`+appGUID+`\LocalServer32`, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to create CLSID key: %w", err)
	}

	if err := key.SetStringValue("", fmt.Sprintf("\"%s\" %%1", exePath)); err != nil {
		return fmt.Errorf("failed to set CLSID server path: %w", err)
	}
	key.Close()

	toast.SetAppData(toast.AppData{
		AppID:         appName,
		GUID:          appGUID,
		IconPath:      iconPath,
		ActivationExe: exePath,
	})

	toast.SetActivationCallback(func(args string, data []toast.UserData) {
		result := frontend.NotificationResult{}

		actionIdentifier, options, err := parseNotificationResponse(args)

		if err != nil {
			result.Error = err
		} else {
			// Subtitle is retained but was not shown with the notification
			response := frontend.NotificationResponse{
				ID:               options.ID,
				ActionIdentifier: actionIdentifier,
				Title:            options.Title,
				Subtitle:         options.Subtitle,
				Body:             options.Body,
				CategoryID:       options.CategoryID,
				UserInfo:         options.Data,
			}

			if userText, found := getUserText(data); found {
				response.UserText = userText
			}

			result.Response = response
		}

		handleNotificationResult(result)
	})

	// Register the class factory for the toast activator
	if err := registerFactoryInternal(wintoast.ClassFactory); err != nil {
		return fmt.Errorf("CoRegisterClassObject failed: %w", err)
	}

	return loadCategoriesFromRegistry()
}

//go:linkname registerFactoryInternal git.sr.ht/~jackmordaunt/go-toast/v2/wintoast.registerClassFactory
func registerFactoryInternal(factory *wintoast.IClassFactory) error

func (f *Frontend) IsNotificationAvailable() bool {
	return true
}

func (f *Frontend) RequestNotificationAuthorization() (bool, error) {
	return true, nil
}

func (f *Frontend) CheckNotificationAuthorization() (bool, error) {
	return true, nil
}

// SendNotification sends a basic notification with a name, title, and body. All other options are ignored on Windows.
// (subtitle is only available on macOS and Linux)
func (f *Frontend) SendNotification(options frontend.NotificationOptions) error {
	if err := f.saveIconToDir(); err != nil {
		fmt.Printf("Error saving icon: %v\n", err)
	}

	n := toast.Notification{
		Title:               options.Title,
		Body:                options.Body,
		ActivationType:      toast.Foreground,
		ActivationArguments: DefaultActionIdentifier,
	}

	if options.Data != nil {
		encodedPayload, err := encodePayload(DefaultActionIdentifier, options)
		if err != nil {
			return fmt.Errorf("failed to encode notification payload: %w", err)
		}
		n.ActivationArguments = encodedPayload
	}

	return n.Push()
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
// A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
// If a NotificationCategory is not registered a basic notification will be sent.
// (subtitle is only available on macOS and Linux)
func (f *Frontend) SendNotificationWithActions(options frontend.NotificationOptions) error {
	if err := f.saveIconToDir(); err != nil {
		fmt.Printf("Error saving icon: %v\n", err)
	}

	categoriesLock.RLock()
	nCategory, categoryExists := categories[options.CategoryID]
	categoriesLock.RUnlock()

	if options.CategoryID == "" || !categoryExists {
		fmt.Printf("Category '%s' not found, sending basic notification without actions\n", options.CategoryID)
		return f.SendNotification(options)
	}

	n := toast.Notification{
		Title:               options.Title,
		Body:                options.Body,
		ActivationType:      toast.Foreground,
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
		encodedPayload, err := encodePayload(n.ActivationArguments, options)
		if err != nil {
			return fmt.Errorf("failed to encode notification payload: %w", err)
		}
		n.ActivationArguments = encodedPayload

		for index := range n.Actions {
			encodedPayload, err := encodePayload(n.Actions[index].Arguments, options)
			if err != nil {
				return fmt.Errorf("failed to encode notification payload: %w", err)
			}
			n.Actions[index].Arguments = encodedPayload
		}
	}

	return n.Push()
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
// Registering a category with the same name as a previously registered NotificationCategory will override it.
func (f *Frontend) RegisterNotificationCategory(category frontend.NotificationCategory) error {
	categoriesLock.Lock()
	defer categoriesLock.Unlock()

	categories[category.ID] = frontend.NotificationCategory{
		ID:               category.ID,
		Actions:          category.Actions,
		HasReplyField:    bool(category.HasReplyField),
		ReplyPlaceholder: category.ReplyPlaceholder,
		ReplyButtonTitle: category.ReplyButtonTitle,
	}

	return saveCategoriesToRegistry()
}

// RemoveNotificationCategory removes a previously registered NotificationCategory.
func (f *Frontend) RemoveNotificationCategory(categoryId string) error {
	categoriesLock.Lock()
	defer categoriesLock.Unlock()

	delete(categories, categoryId)

	return saveCategoriesToRegistry()
}

// RemoveAllPendingNotifications is a Windows stub that always returns nil.
// (macOS and Linux only)
func (f *Frontend) RemoveAllPendingNotifications() error {
	return nil
}

// RemovePendingNotification is a Windows stub that always returns nil.
// (macOS and Linux only)
func (f *Frontend) RemovePendingNotification(_ string) error {
	return nil
}

// RemoveAllDeliveredNotifications is a Windows stub that always returns nil.
// (macOS and Linux only)
func (f *Frontend) RemoveAllDeliveredNotifications() error {
	return nil
}

// RemoveDeliveredNotification is a Windows stub that always returns nil.
// (macOS and Linux only)
func (f *Frontend) RemoveDeliveredNotification(_ string) error {
	return nil
}

// RemoveNotification is a Windows stub that always returns nil.
// (Linux-specific)
func (f *Frontend) RemoveNotification(identifier string) error {
	return nil
}

func (f *Frontend) OnNotificationResponse(callback func(result frontend.NotificationResult)) {
	callbackLock.Lock()
	defer callbackLock.Unlock()

	notificationResultCallback = callback
}

func (f *Frontend) saveIconToDir() error {
	hIcon := w32.ExtractIcon(exePath, 0)
	if hIcon == 0 {
		return fmt.Errorf("ExtractIcon failed for %s", exePath)
	}

	if err := winc.SaveHIconAsPNG(hIcon, iconPath); err != nil {
		return fmt.Errorf("SaveHIconAsPNG failed: %w", err)
	}

	return nil
}

func saveCategoriesToRegistry() error {
	// We assume lock is held by caller

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

	data, err := json.Marshal(categories)
	if err != nil {
		return err
	}

	return key.SetStringValue(NotificationCategoriesRegistryKey, string(data))
}

func loadCategoriesFromRegistry() error {
	// We assume lock is held by caller

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

	_categories := make(map[string]frontend.NotificationCategory)
	if err := json.Unmarshal([]byte(data), &_categories); err != nil {
		return fmt.Errorf("failed to parse notification categories from registry: %w", err)
	}

	categories = _categories

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

// encodePayload combines an action ID and user data into a single encoded string
func encodePayload(actionID string, options frontend.NotificationOptions) (string, error) {
	payload := NotificationPayload{
		Action:  actionID,
		Options: options,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return actionID, err
	}

	encodedPayload := base64.StdEncoding.EncodeToString(jsonData)
	return encodedPayload, nil
}

// decodePayload extracts the action ID and user data from an encoded payload
func decodePayload(encodedString string) (string, frontend.NotificationOptions, error) {
	jsonData, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return encodedString, frontend.NotificationOptions{}, fmt.Errorf("failed to decode base64 payload: %w", err)
	}

	var payload NotificationPayload
	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return encodedString, frontend.NotificationOptions{}, fmt.Errorf("failed to unmarshal notification payload: %w", err)
	}

	return payload.Action, payload.Options, nil
}

// parseNotificationResponse updated to use structured payload decoding
func parseNotificationResponse(response string) (action string, options frontend.NotificationOptions, err error) {
	actionID, options, err := decodePayload(response)

	if err != nil {
		fmt.Printf("Warning: Failed to decode notification response: %v\n", err)
		return response, frontend.NotificationOptions{}, err
	}

	return actionID, options, nil
}

func handleNotificationResult(result frontend.NotificationResult) {
	callbackLock.Lock()
	callback := notificationResultCallback
	callbackLock.Unlock()

	if callback != nil {
		go callback(result)
	}
}

// Helper functions

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
