//go:build darwin

package application

/*
#cgo CFLAGS: -mmacosx-version-min=10.14 -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=10.14 -framework UserNotifications
#import "notifications_darwin.h"
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"unsafe"
)

// NotificationAction represents a button in a notification
type NotificationAction struct {
	ID                     string `json:"id"`
	Title                  string `json:"title"`
	Destructive            bool   `json:"destructive,omitempty"`
	AuthenticationRequired bool   `json:"authenticationRequired,omitempty"`
}

// NotificationCategory groups actions for notifications
type NotificationCategory struct {
	ID               string               `json:"id"`
	Actions          []NotificationAction `json:"actions"`
	HasReplyField    bool                 `json:"hasReplyField,omitempty"`
	ReplyPlaceholder string               `json:"replyPlaceholder,omitempty"`
	ReplyButtonTitle string               `json:"replyButtonTitle,omitempty"`
}

// NotificationOptions contains configuration for a notification
type NotificationOptions struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	Subtitle   string                 `json:"subtitle,omitempty"`
	Body       string                 `json:"body"`
	CategoryID string                 `json:"categoryId,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// Check if the app has a valid bundle identifier
func CheckBundleIdentifier() bool {
	return bool(C.checkBundleIdentifier())
}

// RequestUserNotificationAuthorization requests permission for notifications.
func RequestUserNotificationAuthorization() (bool, error) {
	if !CheckBundleIdentifier() {
		return false, fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	result := C.requestUserNotificationAuthorization(nil)
	return result == true, nil
}

// CheckNotificationAuthorization checks current permission status
func CheckNotificationAuthorization() (bool, error) {
	if !CheckBundleIdentifier() {
		return false, fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	return bool(C.checkNotificationAuthorization()), nil
}

// SendNotification sends a notification with the given identifier, title, subtitle, and body.
func SendNotification(identifier, title, subtitle, body string) error {
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
func SendNotificationWithActions(options NotificationOptions) error {
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
func RegisterNotificationCategory(category NotificationCategory) error {
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
func RemoveAllPendingNotifications() error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	C.removeAllPendingNotifications()
	return nil
}

// RemovePendingNotification removes a specific pending notification
func RemovePendingNotification(identifier string) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	cIdentifier := C.CString(identifier)
	defer C.free(unsafe.Pointer(cIdentifier))
	C.removePendingNotification(cIdentifier)
	return nil
}

func RemoveAllDeliveredNotifications() error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	C.removeAllDeliveredNotifications()
	return nil
}

func RemoveDeliveredNotification(identifier string) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("Notifications require a bundled application with a unique bundle identifier")
	}
	cIdentifier := C.CString(identifier)
	defer C.free(unsafe.Pointer(cIdentifier))
	C.removeDeliveredNotification(cIdentifier)
	return nil
}
