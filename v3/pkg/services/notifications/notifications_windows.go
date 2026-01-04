//go:build windows

package notifications

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	_ "unsafe"

	"encoding/json"

	"git.sr.ht/~jackmordaunt/go-toast/v2"
	wintoast "git.sr.ht/~jackmordaunt/go-toast/v2/wintoast"
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
	exePath        string
}

const (
	ToastRegistryPath                  = `Software\Classes\AppUserModelId\`
	ToastRegistryGuidKey               = "CustomActivator"
	NotificationCategoriesRegistryPath = `SOFTWARE\%s\NotificationCategories`
	NotificationCategoriesRegistryKey  = "Categories"
)

// NotificationPayload combines the action ID and user data into a single structure
type NotificationPayload struct {
	Action  string              `json:"action"`
	Options NotificationOptions `json:"payload,omitempty"`
}

// Creates a new Notifications Service.
func New() *NotificationService {
	notificationServiceOnce.Do(func() {
		impl := &windowsNotifier{
			categories: make(map[string]NotificationCategory),
		}

		NotificationService_ = &NotificationService{
			impl: impl,
		}
	})

	return NotificationService_
}

//go:linkname registerFactoryInternal git.sr.ht/~jackmordaunt/go-toast/v2/wintoast.registerClassFactory
func registerFactoryInternal(factory *wintoast.IClassFactory) error

// Startup is called when the service is loaded
// Sets an activation callback to emit an event when notifications are interacted with.
func (wn *windowsNotifier) Startup(ctx context.Context, options application.ServiceOptions) error {
	wn.categoriesLock.Lock()
	defer wn.categoriesLock.Unlock()

	app := application.Get()
	cfg := app.Config()

	wn.appName = cfg.Name

	guid, err := wn.getGUID()
	if err != nil {
		return err
	}
	wn.appGUID = guid

	wn.iconPath = filepath.Join(os.TempDir(), wn.appName+wn.appGUID+".png")

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	wn.exePath = exe

	// Create the registry key for the toast activator
	key, _, err := registry.CreateKey(registry.CURRENT_USER,
		`Software\Classes\CLSID\`+wn.appGUID+`\LocalServer32`, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to create CLSID key: %w", err)
	}

	if err := key.SetStringValue("", fmt.Sprintf("\"%s\" %%1", wn.exePath)); err != nil {
		return fmt.Errorf("failed to set CLSID server path: %w", err)
	}
	key.Close()

	toast.SetAppData(toast.AppData{
		AppID:         wn.appName,
		GUID:          guid,
		IconPath:      wn.iconPath,
		ActivationExe: wn.exePath,
	})

	toast.SetActivationCallback(func(args string, data []toast.UserData) {
		result := NotificationResult{}

		actionIdentifier, options, err := parseNotificationResponse(args)

		if err != nil {
			result.Error = err
		} else {
			// Subtitle is retained but was not shown with the notification
			response := NotificationResponse{
				ID:               options.ID,
				ActionIdentifier: actionIdentifier,
				Title:            options.Title,
				Subtitle:         options.Subtitle,
				Body:             options.Body,
				CategoryID:       options.CategoryID,
				UserInfo:         options.Data,
			}

			if userText, found := wn.getUserText(data); found {
				response.UserText = userText
			}

			result.Response = response
		}

		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(result)
		}
	})

	// Register the class factory for the toast activator
	if err := registerFactoryInternal(wintoast.ClassFactory); err != nil {
		return fmt.Errorf("CoRegisterClassObject failed: %w", err)
	}

	return wn.loadCategoriesFromRegistry()
}

// Shutdown will attempt to save the categories to the registry when the service unloads
func (wn *windowsNotifier) Shutdown() error {
	wn.categoriesLock.Lock()
	defer wn.categoriesLock.Unlock()

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
// (subtitle is only available on macOS and Linux)
func (wn *windowsNotifier) SendNotification(options NotificationOptions) error {
	if err := wn.saveIconToDir(); err != nil {
		fmt.Printf("Error saving icon: %v\n", err)
	}

	n := toast.Notification{
		Title:               options.Title,
		Body:                options.Body,
		ActivationType:      toast.Foreground,
		ActivationArguments: DefaultActionIdentifier,
	}

	encodedPayload, err := wn.encodePayload(DefaultActionIdentifier, options)
	if err != nil {
		return fmt.Errorf("failed to encode notification payload: %w", err)
	}
	n.ActivationArguments = encodedPayload

	return n.Push()
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
// A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
// If a NotificationCategory is not registered a basic notification will be sent.
// (subtitle is only available on macOS and Linux)
func (wn *windowsNotifier) SendNotificationWithActions(options NotificationOptions) error {
	if err := wn.saveIconToDir(); err != nil {
		fmt.Printf("Error saving icon: %v\n", err)
	}

	wn.categoriesLock.RLock()
	nCategory, categoryExists := wn.categories[options.CategoryID]
	wn.categoriesLock.RUnlock()

	if options.CategoryID == "" || !categoryExists {
		fmt.Printf("Category '%s' not found, sending basic notification without actions\n", options.CategoryID)
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

	encodedPayload, err := wn.encodePayload(n.ActivationArguments, options)
	if err != nil {
		return fmt.Errorf("failed to encode notification payload: %w", err)
	}
	n.ActivationArguments = encodedPayload

	for index := range n.Actions {
		encodedPayload, err := wn.encodePayload(n.Actions[index].Arguments, options)
		if err != nil {
			return fmt.Errorf("failed to encode notification payload: %w", err)
		}
		n.Actions[index].Arguments = encodedPayload
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
// (macOS and Linux only)
func (wn *windowsNotifier) RemoveAllPendingNotifications() error {
	return nil
}

// RemovePendingNotification is a Windows stub that always returns nil.
// (macOS and Linux only)
func (wn *windowsNotifier) RemovePendingNotification(_ string) error {
	return nil
}

// RemoveAllDeliveredNotifications is a Windows stub that always returns nil.
// (macOS and Linux only)
func (wn *windowsNotifier) RemoveAllDeliveredNotifications() error {
	return nil
}

// RemoveDeliveredNotification is a Windows stub that always returns nil.
// (macOS and Linux only)
func (wn *windowsNotifier) RemoveDeliveredNotification(_ string) error {
	return nil
}

// RemoveNotification is a Windows stub that always returns nil.
// (Linux-specific)
func (wn *windowsNotifier) RemoveNotification(identifier string) error {
	return nil
}

// encodePayload combines an action ID and user data into a single encoded string
func (wn *windowsNotifier) encodePayload(actionID string, options NotificationOptions) (string, error) {
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
func decodePayload(encodedString string) (string, NotificationOptions, error) {
	jsonData, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return encodedString, NotificationOptions{}, fmt.Errorf("failed to decode base64 payload: %w", err)
	}

	var payload NotificationPayload
	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return encodedString, NotificationOptions{}, fmt.Errorf("failed to unmarshal notification payload: %w", err)
	}

	return payload.Action, payload.Options, nil
}

// parseNotificationResponse updated to use structured payload decoding
func parseNotificationResponse(response string) (action string, options NotificationOptions, err error) {
	actionID, options, err := decodePayload(response)

	if err != nil {
		fmt.Printf("Warning: Failed to decode notification response: %v\n", err)
		return response, NotificationOptions{}, err
	}

	return actionID, options, nil
}

func (wn *windowsNotifier) saveIconToDir() error {
	icon, err := application.NewIconFromResource(w32.GetModuleHandle(""), uint16(3))
	if err != nil {
		return fmt.Errorf("failed to retrieve application icon: %w", err)
	}

	return w32.SaveHIconAsPNG(icon, wn.iconPath)
}

func (wn *windowsNotifier) saveCategoriesToRegistry() error {
	// We assume lock is held by caller

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
	// We assume lock is held by caller

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
