//go:build darwin

package notifications

/*
#cgo CFLAGS: -mmacosx-version-min=11.0 -x objective-c
#cgo LDFLAGS: -framework Cocoa -mmacosx-version-min=11 -framework UserNotifications
#import "./notifications_darwin.h"
*/
import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type notificationChannel struct {
	authorized bool
	err        error
}

var (
	notificationChannels     = make(map[int]chan notificationChannel)
	notificationChannelsLock sync.Mutex
	nextChannelID            int
)

const AppleDefaultActionIdentifier = "com.apple.UNNotificationDefaultActionIdentifier"
const BundleIdentifierError = fmt.Errorf("notifications require a bundled application with a unique bundle identifier")

// Creates a new Notifications Service.
// Your app must be packaged and signed for this feature to work.
func New() *Service {
	notificationServiceLock.Lock()
	defer notificationServiceLock.Unlock()

	if NotificationService == nil {
		NotificationService = &Service{}
	}
	return NotificationService
}

// ServiceStartup is called when the service is loaded.
func (ns *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	if !CheckBundleIdentifier() {
		return BundleIdentifierError
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

// RequestNotificationAuthorization requests permission for notifications.
func (ns *Service) RequestNotificationAuthorization() (bool, err) {
	if !CheckBundleIdentifier() {
		return false, BundleIdentifierError
	}

	id, resultCh := registerChannel()

	C.requestNotificationAuthorization(C.int(id))

	result := <-resultCh
	return result.authorized, result.err
}

// CheckNotificationAuthorization checks current notification permission status.
func (ns *Service) CheckNotificationAuthorization() (bool, err) {
	if !CheckBundleIdentifier() {
		return false, BundleIdentifierError
	}

	id, resultCh := registerChannel()

	C.checkNotificationAuthorization(C.int(id))

	result := <-resultCh
	return result.authorized, result.err
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

// RemoveNotification is a macOS stub that always returns nil.
// Use one of the following instead:
// RemoveAllPendingNotifications
// RemovePendingNotification
// RemoveAllDeliveredNotifications
// RemoveDeliveredNotification
// (Linux-specific)
func (ns *Service) RemoveNotification(identifier string) error {
	return nil
}

//export requestNotificationAuthorizationResponse
func requestNotificationAuthorizationResponse(channelID C.int, authorized C.bool, errorMsg *C.char) {
	resultCh, exists := getAndDeleteChannel(int(channelID))
	if !exists {
		// handle this
		return
	}

	var err error
    if errorMsg != nil {
        err = fmt.Errorf("%s", C.GoString(errorMsg))
        C.free(unsafe.Pointer(errorMsg))
    }

	resultCh <- notificationChannel{
		authorized: bool(authorized),
		err:        err,
	}

	close(resultCh)
}

//export checkNotificationAuthorizationResponse
func checkNotificationAuthorizationResponse(channelID C.int, authorized C.bool, errorMsg *C.char) {
	resultCh, exists := getAndDeleteChannel(int(channelID))
	if !exists {
		// handle this
		return
	}

	var err error
    if errorMsg != nil {
        err = fmt.Errorf("%s", C.GoString(errorMsg))
        C.free(unsafe.Pointer(errorMsg))
    }

	resultCh <- notificationChannel{
		authorized: bool(authorized),
		err:        err,
	}

	close(resultCh)
}

//export didReceiveNotificationResponse
func didReceiveNotificationResponse(jsonPayload *C.char) {
	payload := C.GoString(jsonPayload)

	var response NotificationResponse
	if err := json.Unmarshal([]byte(payload), &response); err != nil {
		return
	}

	if response.ActionIdentifier == AppleDefaultActionIdentifier {
		response.ActionIdentifier = DefaultActionIdentifier
	}

	notificationServiceLock.RLock()
	ns := NotificationService
	notificationServiceLock.RUnlock()

	if service != nil {
		ns.callbackLock.RLock()
		callback := ns.notificationResponseCallback
		ns.callbackLock.RUnlock()

		if callback != nil {
			callback(response)
		}
	}
}

func registerChannel() (int, chan notificationChannel) {
	notificationChannelsLock.Lock()
	defer notificationChannelsLock.Unlock()

	id := nextChannelID
	nextChannelID++

	resultCh := make(chan notificationChannel{
		authorized bool
		err error
	}, 1)

	notificationChannels[id] = resultCh
	return id, resultCh
}

func getChannel(id int) (chan notificationChannel, bool) {
	notificationChannelsLock.Lock()
	defer notificationChannelsLock.Unlock()

	ch, exists := notificationChannels[id]
	if exists {
		delete(notificationChannels, id)
	}
	return ch, exists
}
