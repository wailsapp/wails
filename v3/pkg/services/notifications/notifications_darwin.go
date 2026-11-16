//go:build darwin && !ios && !purego

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
	"fmt"
	"sync"
	"time"
	"unsafe"

	"encoding/json"

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

// isNotificationAvailable checks if notifications are available on the system.
func isNotificationAvailable() bool {
	return bool(C.isNotificationAvailable())
}

func checkBundleIdentifier() bool {
	return bool(C.checkBundleIdentifier())
}

// dispatchNotification is the shared body of SendNotification and
// SendNotificationWithActions. The whole NotificationOptions struct is
// JSON-encoded and passed to Objective-C as a single string so future fields
// (sound, attachments, threadId, interruptionLevel, schedule) can be added on
// the Go side without changing C signatures.
func (dn *darwinNotifier) dispatchNotification(options NotificationOptions, withActions bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("failed to marshal notification options: %w", err)
	}

	cOptionsJSON := C.CString(string(jsonBytes))
	defer C.free(unsafe.Pointer(cOptionsJSON))

	id, resultCh := dn.registerChannel()
	if withActions {
		C.sendNotificationWithActions(C.int(id), cOptionsJSON)
	} else {
		C.sendNotification(C.int(id), cOptionsJSON)
	}

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

// checkNotificationAuthorization asks the system for the current notification
// authorization; the result is delivered on the channel registered as channelID.
func checkNotificationAuthorization(channelID int) {
	C.checkNotificationAuthorization(C.int(channelID))
}

// requestNotificationAuthorization prompts for notification authorization; the
// result is delivered on the channel registered as channelID.
func requestNotificationAuthorization(channelID int) {
	C.requestNotificationAuthorization(C.int(channelID))
}

// removeNotificationCategory removes categoryID; the result is delivered on the
// channel registered as channelID.
func removeNotificationCategory(channelID int, categoryID string) {
	cCategoryID := C.CString(categoryID)
	defer C.free(unsafe.Pointer(cCategoryID))
	C.removeNotificationCategory(C.int(channelID), cCategoryID)
}
