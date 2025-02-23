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

// Creates a new Notifications Service.
// Your app must be packaged and signed for this feature to work.
func New() *Service {
	return &Service{}
}

// ServiceName returns the name of the service.
func (ns *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/services/notifications"
}

// ServiceStartup is called when the service is loaded.
func (ns *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("notifications require a bundled application with a unique bundle identifier")
	}
	return nil
}

// ServiceShutdown is called when the service is unloaded.
func (ns *Service) ServiceShutdown() error {
	return nil
}

func CheckBundleIdentifier() bool {
	return bool(C.checkBundleIdentifier())
}

// RequestUserNotificationAuthorization requests permission for notifications.
func (ns *Service) RequestUserNotificationAuthorization() (bool, error) {
	if !CheckBundleIdentifier() {
		return false, fmt.Errorf("notifications require a bundled application with a unique bundle identifier")
	}
	result := C.requestUserNotificationAuthorization(nil)
	return result == true, nil
}

// CheckNotificationAuthorization checks current notification permission status.
func (ns *Service) CheckNotificationAuthorization() (bool, error) {
	if !CheckBundleIdentifier() {
		return false, fmt.Errorf("notifications require a bundled application with a unique bundle identifier")
	}
	return bool(C.checkNotificationAuthorization()), nil
}

// SendNotification sends a basic notification with a unique identifier, title, subtitle, and body.
func (ns *Service) SendNotification(options NotificationOptions) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("notifications require a bundled application with a unique bundle identifier")
	}
	cIdentifier := C.CString(options.ID)
	cTitle := C.CString(options.Title)
	cSubtitle := C.CString(options.Subtitle)
	cBody := C.CString(options.Body)
	defer C.free(unsafe.Pointer(cIdentifier))
	defer C.free(unsafe.Pointer(cTitle))
	defer C.free(unsafe.Pointer(cSubtitle))
	defer C.free(unsafe.Pointer(cBody))

	var cDataJSON *C.char
	if options.Data != nil {
		jsonData, err := json.Marshal(options.Data)
		if err == nil {
			cDataJSON = C.CString(string(jsonData))
			defer C.free(unsafe.Pointer(cDataJSON))
		}
	}

	C.sendNotification(cIdentifier, cTitle, cSubtitle, cBody, cDataJSON, nil)
	return nil
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
// A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
// If a NotificationCategory is not registered a basic notification will be sent.
func (ns *Service) SendNotificationWithActions(options NotificationOptions) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("notifications require a bundled application with a unique bundle identifier")
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

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
// Registering a category with the same name as a previously registered NotificationCategory will override it.
func (ns *Service) RegisterNotificationCategory(category NotificationCategory) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("notifications require a bundled application with a unique bundle identifier")
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

// RemoveNotificationCategory remove a previously registered NotificationCategory.
func (ns *Service) RemoveNotificationCategory(categoryId string) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("notifications require a bundled application with a unique bundle identifier")
	}
	cCategoryID := C.CString(categoryId)
	defer C.free(unsafe.Pointer(cCategoryID))

	C.removeNotificationCategory(cCategoryID)
	return nil
}

// RemoveAllPendingNotifications removes all pending notifications.
func (ns *Service) RemoveAllPendingNotifications() error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("notifications require a bundled application with a unique bundle identifier")
	}
	C.removeAllPendingNotifications()
	return nil
}

// RemovePendingNotification removes a pending notification matching the unique identifier.
func (ns *Service) RemovePendingNotification(identifier string) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("notifications require a bundled application with a unique bundle identifier")
	}
	cIdentifier := C.CString(identifier)
	defer C.free(unsafe.Pointer(cIdentifier))
	C.removePendingNotification(cIdentifier)
	return nil
}

// RemoveAllDeliveredNotifications removes all delivered notifications.
func (ns *Service) RemoveAllDeliveredNotifications() error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("notifications require a bundled application with a unique bundle identifier")
	}
	C.removeAllDeliveredNotifications()
	return nil
}

// RemoveDeliveredNotification removes a delivered notification matching the unique identifier.
func (ns *Service) RemoveDeliveredNotification(identifier string) error {
	if !CheckBundleIdentifier() {
		return fmt.Errorf("notifications require a bundled application with a unique bundle identifier")
	}
	cIdentifier := C.CString(identifier)
	defer C.free(unsafe.Pointer(cIdentifier))
	C.removeDeliveredNotification(cIdentifier)
	return nil
}

//export didReceiveNotificationResponse
func didReceiveNotificationResponse(jsonPayload *C.char) {
	payload := C.GoString(jsonPayload)

	var response NotificationResponseData
	if err := json.Unmarshal([]byte(payload), &response); err != nil {
		return
	}

	application.Get().EmitEvent("notificationResponse", NotificationResponse{
		Name: "notification",
		Data: response,
	})
}
