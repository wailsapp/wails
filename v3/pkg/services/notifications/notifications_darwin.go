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
	"time"
	"unsafe"
)

type notificationChannel struct {
	Success bool
	Error   error
}

var (
	notificationChannels     = make(map[int]chan notificationChannel)
	notificationChannelsLock sync.Mutex
	nextChannelID            int
)

const AppleDefaultActionIdentifier = "com.apple.UNNotificationDefaultActionIdentifier"

// Creates a new Notifications Service.
// Your app must be packaged and signed for this feature to work.
func New() *Service {
	notificationServiceOnce.Do(func() {
		if !CheckBundleIdentifier() {
			panic("\nError: Cannot use the notification API in development mode on macOS.\n" +
				"Notifications require the app to be properly bundled with a bundle identifier and signed.\n" +
				"To use the notification API on macOS:\n" +
				"  1. Build and package your app using 'wails3 package'\n" +
				"  2. Sign the packaged .app\n" +
				"  3. Run the signed .app bundle")
		}

		if NotificationService == nil {
			NotificationService = &Service{}
		}
	})

	return NotificationService
}

func CheckBundleIdentifier() bool {
	return bool(C.checkBundleIdentifier())
}

// RequestNotificationAuthorization requests permission for notifications.
// Default timeout is 15 minutes
func (ns *Service) RequestNotificationAuthorization() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*900)
	defer cancel()

	id, resultCh := registerChannel()

	C.requestNotificationAuthorization(C.int(id))

	select {
	case result := <-resultCh:
		return result.Success, result.Error
	case <-ctx.Done():
		cleanupChannel(id)
		return false, fmt.Errorf("notification authorization timed out after 15 minutes: %w", ctx.Err())
	}
}

// CheckNotificationAuthorization checks current notification permission status.
func (ns *Service) CheckNotificationAuthorization() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	id, resultCh := registerChannel()

	C.checkNotificationAuthorization(C.int(id))

	select {
	case result := <-resultCh:
		return result.Success, result.Error
	case <-ctx.Done():
		cleanupChannel(id)
		return false, fmt.Errorf("notification authorization timed out after 15s: %w", ctx.Err())
	}
}

// SendNotification sends a basic notification with a unique identifier, title, subtitle, and body.
func (ns *Service) SendNotification(options NotificationOptions) error {
	if err := validateNotificationOptions(options); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, resultCh := registerChannel()

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
		if err != nil {
			return fmt.Errorf("failed to marshal notification data: %w", err)
		}
		cDataJSON = C.CString(string(jsonData))
		defer C.free(unsafe.Pointer(cDataJSON))
	}

	C.sendNotification(C.int(id), cIdentifier, cTitle, cSubtitle, cBody, cDataJSON)

	select {
	case result := <-resultCh:
		if !result.Success {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("sending notification failed")
		}
		return nil
	case <-ctx.Done():
		cleanupChannel(id)
		return fmt.Errorf("sending notification timed out: %w", ctx.Err())
	}
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
// A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
// If a NotificationCategory is not registered a basic notification will be sent.
func (ns *Service) SendNotificationWithActions(options NotificationOptions) error {
	if err := validateNotificationOptions(options); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, resultCh := registerChannel()

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

	var cDataJSON *C.char
	if options.Data != nil {
		jsonData, err := json.Marshal(options.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal notification data: %w", err)
		}
		cDataJSON = C.CString(string(jsonData))
		defer C.free(unsafe.Pointer(cDataJSON))
	}

	C.sendNotificationWithActions(C.int(id), cIdentifier, cTitle, cSubtitle, cBody, cCategoryID, cDataJSON)
	select {
	case result := <-resultCh:
		if !result.Success {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("sending notification failed")
		}
		return nil
	case <-ctx.Done():
		cleanupChannel(id)
		return fmt.Errorf("sending notification timed out: %w", ctx.Err())
	}
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
// Registering a category with the same name as a previously registered NotificationCategory will override it.
func (ns *Service) RegisterNotificationCategory(category NotificationCategory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, resultCh := registerChannel()

	cCategoryID := C.CString(category.ID)
	defer C.free(unsafe.Pointer(cCategoryID))

	actionsJSON, err := json.Marshal(category.Actions)
	if err != nil {
		return fmt.Errorf("failed to marshal notification category: %w", err)
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

	C.registerNotificationCategory(C.int(id), cCategoryID, cActionsJSON, C.bool(category.HasReplyField),
		cReplyPlaceholder, cReplyButtonTitle)

	select {
	case result := <-resultCh:
		if !result.Success {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("category registration failed")
		}
		return nil
	case <-ctx.Done():
		cleanupChannel(id)
		return fmt.Errorf("category registration timed out: %w", ctx.Err())
	}
}

// RemoveNotificationCategory remove a previously registered NotificationCategory.
func (ns *Service) RemoveNotificationCategory(categoryId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, resultCh := registerChannel()

	cCategoryID := C.CString(categoryId)
	defer C.free(unsafe.Pointer(cCategoryID))

	C.removeNotificationCategory(C.int(id), cCategoryID)

	select {
	case result := <-resultCh:
		if !result.Success {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("category removal failed")
		}
		return nil
	case <-ctx.Done():
		cleanupChannel(id)
		return fmt.Errorf("category removal timed out: %w", ctx.Err())
	}
}

// RemoveAllPendingNotifications removes all pending notifications.
func (ns *Service) RemoveAllPendingNotifications() error {
	C.removeAllPendingNotifications()
	return nil
}

// RemovePendingNotification removes a pending notification matching the unique identifier.
func (ns *Service) RemovePendingNotification(identifier string) error {
	cIdentifier := C.CString(identifier)
	defer C.free(unsafe.Pointer(cIdentifier))
	C.removePendingNotification(cIdentifier)
	return nil
}

// RemoveAllDeliveredNotifications removes all delivered notifications.
func (ns *Service) RemoveAllDeliveredNotifications() error {
	C.removeAllDeliveredNotifications()
	return nil
}

// RemoveDeliveredNotification removes a delivered notification matching the unique identifier.
func (ns *Service) RemoveDeliveredNotification(identifier string) error {
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

//export captureResult
func captureResult(channelID C.int, success C.bool, errorMsg *C.char) {
	resultCh, exists := getChannel(int(channelID))
	if !exists {
		return
	}

	var err error
	if errorMsg != nil {
		err = fmt.Errorf("%s", C.GoString(errorMsg))
	}

	resultCh <- notificationChannel{
		Success: bool(success),
		Error:   err,
	}

	close(resultCh)
}

//export didReceiveNotificationResponse
func didReceiveNotificationResponse(jsonPayload *C.char, err *C.char) {
	result := NotificationResult{}

	if err != nil {
		errMsg := C.GoString(err)
		result.Error = fmt.Errorf("notification response error: %s", errMsg)
		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(result)
		}
		return
	}

	if jsonPayload == nil {
		result.Error = fmt.Errorf("received nil JSON payload in notification response")
		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(result)
		}
		return
	}

	payload := C.GoString(jsonPayload)

	var response NotificationResponse
	if err := json.Unmarshal([]byte(payload), &response); err != nil {
		result.Error = fmt.Errorf("failed to unmarshal notification response: %w", err)
		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(result)
		}
		return
	}

	if response.ActionIdentifier == AppleDefaultActionIdentifier {
		response.ActionIdentifier = DefaultActionIdentifier
	}

	result.Response = response
	if ns := getNotificationService(); ns != nil {
		ns.handleNotificationResult(result)
	}
}

func registerChannel() (int, chan notificationChannel) {
	notificationChannelsLock.Lock()
	defer notificationChannelsLock.Unlock()

	id := nextChannelID
	nextChannelID++

	resultCh := make(chan notificationChannel, 1)

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

func cleanupChannel(id int) {
	notificationChannelsLock.Lock()
	defer notificationChannelsLock.Unlock()

	if ch, exists := notificationChannels[id]; exists {
		delete(notificationChannels, id)
		close(ch)
	}
}
