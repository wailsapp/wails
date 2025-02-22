//go:build darwin

package notifications

/*
#cgo CFLAGS: -mmacosx-version-min=10.14 -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.14 -framework UserNotifications
#import "./notifications_darwin.h"
*/
import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Service struct {
}

// NotificationAction represents an action button for a notification
type NotificationAction = struct {
	ID                     string `json:"id"`
	Title                  string `json:"title"`
	Destructive            bool   `json:"destructive,omitempty"`
	AuthenticationRequired bool   `json:"authenticationRequired,omitempty"`
}

// NotificationCategory groups actions for notifications
type NotificationCategory = struct {
	ID               string               `json:"id"`
	Actions          []NotificationAction `json:"actions"`
	HasReplyField    bool                 `json:"hasReplyField,omitempty"`
	ReplyPlaceholder string               `json:"replyPlaceholder,omitempty"`
	ReplyButtonTitle string               `json:"replyButtonTitle,omitempty"`
}

// NotificationOptions contains configuration for a notification
type NotificationOptions = struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	Subtitle   string                 `json:"subtitle,omitempty"`
	Body       string                 `json:"body"`
	CategoryID string                 `json:"categoryId,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

func New() *Service {
	return &Service{}
}

func CheckBundleIdentifier() bool {
	return bool(C.checkBundleIdentifier())
}

// ServiceName returns the name of the service
func (ns *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/services/notifications"
}

// ServiceStartup is called when the service is loaded
func (ns *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	return nil
}

// ServiceShutdown is called when the service is unloaded
func (ns *Service) ServiceShutdown() error {
	return nil
}

// RequestUserNotificationAuthorization requests permission for notifications.
func (ns *Service) RequestUserNotificationAuthorization() (bool, error) {
	if !CheckBundleIdentifier() {
		return false, fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	result := C.requestUserNotificationAuthorization(nil)
	return result == true, nil
}

// CheckNotificationAuthorization checks current permission status
func (ns *Service) CheckNotificationAuthorization() (bool, error) {
	if !CheckBundleIdentifier() {
		return false, fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	return bool(C.checkNotificationAuthorization()), nil
}

// SendNotification sends a notification with the given identifier, title, subtitle, and body.
func (ns *Service) SendNotification(identifier, title, subtitle, body string) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	cIdentifier := C.CString(identifier)
	cTitle := C.CString(title)
	cSubtitle := C.CString(subtitle)
	cBody := C.CString(body)
	defer C.free(unsafe.Pointer(cIdentifier))
	defer C.free(unsafe.Pointer(cTitle))
	defer C.free(unsafe.Pointer(cSubtitle))
	defer C.free(unsafe.Pointer(cBody))

	C.sendNotification(cIdentifier, cTitle, cSubtitle, cBody, nil)
	return nil
}

// SendNotificationWithActions sends a notification with the specified actions
func (ns *Service) SendNotificationWithActions(options NotificationOptions) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	cIdentifier := C.CString(options.ID)
	cTitle := C.CString(options.Title)
	cSubtitle := C.CString(options.Subtitle)
	cBody := C.CString(options.Body)
	cCategoryID := C.CString(options.CategoryID)
	defer C.free(unsafe.Pointer(cIdentifier))
	defer C.free(unsafe.Pointer(cTitle))
	defer C.free(unsafe.Pointer(cSubtitle))
	defer C.free(unsafe.Pointer(cBody))
	defer C.free(unsafe.Pointer(cCategoryID))

	var cActionsJSON *C.char
	if options.Data != nil {
		jsonData, err := json.Marshal(options.Data)
		if err == nil {
			cActionsJSON = C.CString(string(jsonData))
			defer C.free(unsafe.Pointer(cActionsJSON))
		}
	}

	C.sendNotificationWithActions(cIdentifier, cTitle, cSubtitle, cBody, cCategoryID, cActionsJSON, nil)
	return nil
}

// RegisterNotificationCategory registers a category with actions and optional reply field
func (ns *Service) RegisterNotificationCategory(category NotificationCategory) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	cCategoryID := C.CString(category.ID)
	defer C.free(unsafe.Pointer(cCategoryID))

	actionsJSON, err := json.Marshal(category.Actions)
	if err != nil {
		return err
	}
	cActionsJSON := C.CString(string(actionsJSON))
	defer C.free(unsafe.Pointer(cActionsJSON))

	var cReplyPlaceholder, cReplyButtonTitle *C.char
	if category.HasReplyField {
		cReplyPlaceholder = C.CString(category.ReplyPlaceholder)
		cReplyButtonTitle = C.CString(category.ReplyButtonTitle)
		defer C.free(unsafe.Pointer(cReplyPlaceholder))
		defer C.free(unsafe.Pointer(cReplyButtonTitle))
	}

	C.registerNotificationCategory(cCategoryID, cActionsJSON, C.bool(category.HasReplyField),
		cReplyPlaceholder, cReplyButtonTitle)
	return nil
}

// RemoveAllPendingNotifications removes all pending notifications
func (ns *Service) RemoveAllPendingNotifications() error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	C.removeAllPendingNotifications()
	return nil
}

// RemovePendingNotification removes a specific pending notification
func (ns *Service) RemovePendingNotification(identifier string) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	cIdentifier := C.CString(identifier)
	defer C.free(unsafe.Pointer(cIdentifier))
	C.removePendingNotification(cIdentifier)
	return nil
}

func (ns *Service) RemoveAllDeliveredNotifications() error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	C.removeAllDeliveredNotifications()
	return nil
}

func (ns *Service) RemoveDeliveredNotification(identifier string) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	cIdentifier := C.CString(identifier)
	defer C.free(unsafe.Pointer(cIdentifier))
	C.removeDeliveredNotification(cIdentifier)
	return nil
}
