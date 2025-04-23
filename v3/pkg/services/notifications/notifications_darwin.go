//go:build darwin

package notifications

/*
#cgo CFLAGS:-x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa

#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
#cgo LDFLAGS: -framework UserNotifications
#endif

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

	"github.com/wailsapp/wails/v3/pkg/application"
)

type darwinNotifier struct {
	channels      map[int]chan notificationChannel
	channelsLock  sync.Mutex
	nextChannelID int
}

type notificationChannel struct {
	Success bool
	Error   error
}

type ChannelHandler interface {
	GetChannel(id int) (chan notificationChannel, bool)
}

const AppleDefaultActionIdentifier = "com.apple.UNNotificationDefaultActionIdentifier"

// Creates a new Notifications Service.
// Your app must be packaged and signed for this feature to work.
func New() *Service {
	notificationServiceOnce.Do(func() {
		impl := &darwinNotifier{
			channels:      make(map[int]chan notificationChannel),
			nextChannelID: 0,
		}

		NotificationService = &Service{
			impl: impl,
		}
	})

	return NotificationService
}

func (dn *darwinNotifier) Startup(ctx context.Context, options application.ServiceOptions) error {
	if !isNotificationAvailable() {
		return fmt.Errorf("notifications are not available on this system")
	}
	if !checkBundleIdentifier() {
		return fmt.Errorf("notifications require a valid bundle identifier")
	}
	if !bool(C.ensureDelegateInitialized()) {
		return fmt.Errorf("failed to initialize notification center delegate")
	}
	return nil
}

func (dn *darwinNotifier) Shutdown() error {
	return nil
}

// isNotificationAvailable checks if notifications are available on the system.
func isNotificationAvailable() bool {
	return bool(C.isNotificationAvailable())
}

func checkBundleIdentifier() bool {
	return bool(C.checkBundleIdentifier())
}

// RequestNotificationAuthorization requests permission for notifications.
// Default timeout is 3 minutes
func (dn *darwinNotifier) RequestNotificationAuthorization() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	id, resultCh := dn.registerChannel()

	C.requestNotificationAuthorization(C.int(id))

	select {
	case result := <-resultCh:
		return result.Success, result.Error
	case <-ctx.Done():
		dn.cleanupChannel(id)
		return false, fmt.Errorf("notification authorization timed out after 3 minutes: %w", ctx.Err())
	}
}

// CheckNotificationAuthorization checks current notification permission status.
func (dn *darwinNotifier) CheckNotificationAuthorization() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	id, resultCh := dn.registerChannel()

	C.checkNotificationAuthorization(C.int(id))

	select {
	case result := <-resultCh:
		return result.Success, result.Error
	case <-ctx.Done():
		dn.cleanupChannel(id)
		return false, fmt.Errorf("notification authorization timed out after 15s: %w", ctx.Err())
	}
}

// SendNotification sends a basic notification with a unique identifier, title, subtitle, and body.
func (dn *darwinNotifier) SendNotification(options NotificationOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	id, resultCh := dn.registerChannel()
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
		dn.cleanupChannel(id)
		return fmt.Errorf("sending notification timed out: %w", ctx.Err())
	}
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
// A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
// If a NotificationCategory is not registered a basic notification will be sent.
func (dn *darwinNotifier) SendNotificationWithActions(options NotificationOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	id, resultCh := dn.registerChannel()
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
		dn.cleanupChannel(id)
		return fmt.Errorf("sending notification timed out: %w", ctx.Err())
	}
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
// Registering a category with the same name as a previously registered NotificationCategory will override it.
func (dn *darwinNotifier) RegisterNotificationCategory(category NotificationCategory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	id, resultCh := dn.registerChannel()
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
		dn.cleanupChannel(id)
		return fmt.Errorf("category registration timed out: %w", ctx.Err())
	}
}

// RemoveNotificationCategory remove a previously registered NotificationCategory.
func (dn *darwinNotifier) RemoveNotificationCategory(categoryId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cCategoryID := C.CString(categoryId)
	defer C.free(unsafe.Pointer(cCategoryID))

	id, resultCh := dn.registerChannel()
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
		dn.cleanupChannel(id)
		return fmt.Errorf("category removal timed out: %w", ctx.Err())
	}
}

// RemoveAllPendingNotifications removes all pending notifications.
func (dn *darwinNotifier) RemoveAllPendingNotifications() error {
	C.removeAllPendingNotifications()
	return nil
}

// RemovePendingNotification removes a pending notification matching the unique identifier.
func (dn *darwinNotifier) RemovePendingNotification(identifier string) error {
	cIdentifier := C.CString(identifier)
	defer C.free(unsafe.Pointer(cIdentifier))
	C.removePendingNotification(cIdentifier)
	return nil
}

// RemoveAllDeliveredNotifications removes all delivered notifications.
func (dn *darwinNotifier) RemoveAllDeliveredNotifications() error {
	C.removeAllDeliveredNotifications()
	return nil
}

// RemoveDeliveredNotification removes a delivered notification matching the unique identifier.
func (dn *darwinNotifier) RemoveDeliveredNotification(identifier string) error {
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
func (dn *darwinNotifier) RemoveNotification(identifier string) error {
	return nil
}

//export captureResult
func captureResult(channelID C.int, success C.bool, errorMsg *C.char) {
	ns := getNotificationService()
	if ns == nil {
		return
	}

	handler, ok := ns.impl.(ChannelHandler)
	if !ok {
		return
	}

	resultCh, exists := handler.GetChannel(int(channelID))
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

// Helper methods

func (dn *darwinNotifier) registerChannel() (int, chan notificationChannel) {
	dn.channelsLock.Lock()
	defer dn.channelsLock.Unlock()

	id := dn.nextChannelID
	dn.nextChannelID++

	resultCh := make(chan notificationChannel, 1)

	dn.channels[id] = resultCh
	return id, resultCh
}

func (dn *darwinNotifier) GetChannel(id int) (chan notificationChannel, bool) {
	dn.channelsLock.Lock()
	defer dn.channelsLock.Unlock()

	ch, exists := dn.channels[id]
	if exists {
		delete(dn.channels, id)
	}
	return ch, exists
}

func (dn *darwinNotifier) cleanupChannel(id int) {
	dn.channelsLock.Lock()
	defer dn.channelsLock.Unlock()

	if ch, exists := dn.channels[id]; exists {
		delete(dn.channels, id)
		close(ch)
	}
}
