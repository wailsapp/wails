//go:build windows
// +build windows

package windows

import (
	"sync"

	wintoast "git.sr.ht/~jackmordaunt/go-toast/v2/wintoast"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/winc"

	"fmt"
	"os"
	"path/filepath"

	"git.sr.ht/~jackmordaunt/go-toast/v2"
	"github.com/google/uuid"
	"golang.org/x/sys/windows/registry"
)

var (
	categories     map[string]frontend.NotificationCategory
	categoriesLock sync.RWMutex
	appName        string
	appGUID        string
	iconPath       string = ""
	exePath        string
)

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

	// 1. Determine executable path and app name
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable: %w", err)
	}
	exePath = exe
	appName = filepath.Base(exePath)

	// 2. Extract the application's own icon and save it as PNG
	iconPath = filepath.Join(os.TempDir(), appName+".png")

	hIcon, err := winc.ExtractIcon(exePath, 0)
	if err != nil {
		return fmt.Errorf("failed to extract icon: %w", err)
	}
	defer hIcon.Destroy()

	if err := winc.SaveHIconAsPNG(hIcon.Handle(), iconPath); err != nil {
		return fmt.Errorf("failed to save icon as PNG: %w", err)
	}

	// 3. Load or create a GUID for the AppUserModelID
	guid, err := loadOrCreateGUID(appName)
	if err != nil {
		return fmt.Errorf("failed to load/create notification GUID: %w", err)
	}
	appGUID = guid

	// 4. Register toast library
	toast.SetAppData(toast.AppData{
		AppID:         appName,
		GUID:          appGUID,
		IconPath:      iconPath,
		ActivationExe: exePath,
	})
	if err := registerFactoryInternal(wintoast.ClassFactory); err != nil {
		return fmt.Errorf("failed to register toast factory: %w", err)
	}

	return nil
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

// Helper functions

func loadOrCreateGUID(appID string) (string, error) {
	keyPath := ToastRegistryPath + appID
	// Try to open existing GUID
	k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
	if err == nil {
		defer k.Close()
		guid, _, err := k.GetStringValue(ToastRegistryGuidKey)
		if err == nil && guid != "" {
			return guid, nil
		}
	}
	// Create new GUID
	newGUID := uuid.New().String()
	k2, _, err := registry.CreateKey(registry.CURRENT_USER, keyPath, registry.ALL_ACCESS)
	if err != nil {
		return "", err
	}
	defer k2.Close()
	if err := k2.SetStringValue(ToastRegistryGuidKey, newGUID); err != nil {
		return "", err
	}
	return newGUID, nil
}
