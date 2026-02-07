//go:build darwin
// +build darwin

package darwin

/*
#cgo CFLAGS:-x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa

#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
#cgo LDFLAGS: -framework UserNotifications
#endif

#import "Application.h"
#import "WailsContext.h"
*/
import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

// Package-scoped variable only accessible within this file
var (
	currentFrontend *Frontend
	frontendMutex   sync.RWMutex
	// Notification channels
	channels      map[int]chan notificationChannel
	channelsLock  sync.Mutex
	nextChannelID int

	notificationResultCallback func(result frontend.NotificationResult)
	callbackLock               sync.RWMutex
)

const DefaultActionIdentifier = "DEFAULT_ACTION"
const AppleDefaultActionIdentifier = "com.apple.UNNotificationDefaultActionIdentifier"

// setCurrentFrontend sets the current frontend instance
// This is called when RequestNotificationAuthorization or CheckNotificationAuthorization is called
func setCurrentFrontend(f *Frontend) {
	frontendMutex.Lock()
	defer frontendMutex.Unlock()
	currentFrontend = f
}

// getCurrentFrontend gets the current frontend instance
func getCurrentFrontend() *Frontend {
	frontendMutex.RLock()
	defer frontendMutex.RUnlock()
	return currentFrontend
}

type notificationChannel struct {
	Success bool
	Error   error
}

func (f *Frontend) InitializeNotifications() error {
	if !f.IsNotificationAvailable() {
		return fmt.Errorf("notifications are not available on this system")
	}
	if !f.checkBundleIdentifier() {
		return fmt.Errorf("notifications require a valid bundle identifier")
	}
	if !bool(C.EnsureDelegateInitialized(f.mainWindow.context)) {
		return fmt.Errorf("failed to initialize notification center delegate")
	}

	channels = make(map[int]chan notificationChannel)
	nextChannelID = 0

	setCurrentFrontend(f)

	return nil
}

// CleanupNotifications is a macOS stub that does nothing.
// (Linux-specific cleanup)
func (f *Frontend) CleanupNotifications() {
	// No cleanup needed on macOS
}

func (f *Frontend) IsNotificationAvailable() bool {
	return bool(C.IsNotificationAvailable(f.mainWindow.context))
}

func (f *Frontend) checkBundleIdentifier() bool {
	return bool(C.CheckBundleIdentifier(f.mainWindow.context))
}

func (f *Frontend) RequestNotificationAuthorization() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	id, resultCh := f.registerChannel()

	C.RequestNotificationAuthorization(f.mainWindow.context, C.int(id))

	select {
	case result := <-resultCh:
		close(resultCh)
		return result.Success, result.Error
	case <-ctx.Done():
		f.cleanupChannel(id)
		return false, fmt.Errorf("notification authorization timed out after 3 minutes: %w", ctx.Err())
	}
}

func (f *Frontend) CheckNotificationAuthorization() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	id, resultCh := f.registerChannel()

	C.CheckNotificationAuthorization(f.mainWindow.context, C.int(id))

	select {
	case result := <-resultCh:
		close(resultCh)
		return result.Success, result.Error
	case <-ctx.Done():
		f.cleanupChannel(id)
		return false, fmt.Errorf("notification authorization timed out after 15s: %w", ctx.Err())
	}
}

// SendNotification sends a basic notification with a unique identifier, title, subtitle, and body.
func (f *Frontend) SendNotification(options frontend.NotificationOptions) error {
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

	id, resultCh := f.registerChannel()
	C.SendNotification(f.mainWindow.context, C.int(id), cIdentifier, cTitle, cSubtitle, cBody, cDataJSON)

	select {
	case result := <-resultCh:
		close(resultCh)
		if !result.Success {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("sending notification failed")
		}
		return nil
	case <-ctx.Done():
		f.cleanupChannel(id)
		return fmt.Errorf("sending notification timed out: %w", ctx.Err())
	}
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
// A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
// If a NotificationCategory is not registered a basic notification will be sent.
func (f *Frontend) SendNotificationWithActions(options frontend.NotificationOptions) error {
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

	id, resultCh := f.registerChannel()
	C.SendNotificationWithActions(f.mainWindow.context, C.int(id), cIdentifier, cTitle, cSubtitle, cBody, cCategoryID, cDataJSON)

	select {
	case result := <-resultCh:
		close(resultCh)
		if !result.Success {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("sending notification failed")
		}
		return nil
	case <-ctx.Done():
		f.cleanupChannel(id)
		return fmt.Errorf("sending notification timed out: %w", ctx.Err())
	}
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
// Registering a category with the same name as a previously registered NotificationCategory will override it.
func (f *Frontend) RegisterNotificationCategory(category frontend.NotificationCategory) error {
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

	id, resultCh := f.registerChannel()
	C.RegisterNotificationCategory(f.mainWindow.context, C.int(id), cCategoryID, cActionsJSON, C.bool(category.HasReplyField),
		cReplyPlaceholder, cReplyButtonTitle)

	select {
	case result := <-resultCh:
		close(resultCh)
		if !result.Success {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("category registration failed")
		}
		return nil
	case <-ctx.Done():
		f.cleanupChannel(id)
		return fmt.Errorf("category registration timed out: %w", ctx.Err())
	}
}

// RemoveNotificationCategory remove a previously registered NotificationCategory.
func (f *Frontend) RemoveNotificationCategory(categoryId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cCategoryID := C.CString(categoryId)
	defer C.free(unsafe.Pointer(cCategoryID))

	id, resultCh := f.registerChannel()
	C.RemoveNotificationCategory(f.mainWindow.context, C.int(id), cCategoryID)

	select {
	case result := <-resultCh:
		close(resultCh)
		if !result.Success {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("category removal failed")
		}
		return nil
	case <-ctx.Done():
		f.cleanupChannel(id)
		return fmt.Errorf("category removal timed out: %w", ctx.Err())
	}
}

// RemoveAllPendingNotifications removes all pending notifications.
func (f *Frontend) RemoveAllPendingNotifications() error {
	C.RemoveAllPendingNotifications(f.mainWindow.context)
	return nil
}

// RemovePendingNotification removes a pending notification matching the unique identifier.
func (f *Frontend) RemovePendingNotification(identifier string) error {
	cIdentifier := C.CString(identifier)
	defer C.free(unsafe.Pointer(cIdentifier))
	C.RemovePendingNotification(f.mainWindow.context, cIdentifier)
	return nil
}

// RemoveAllDeliveredNotifications removes all delivered notifications.
func (f *Frontend) RemoveAllDeliveredNotifications() error {
	C.RemoveAllDeliveredNotifications(f.mainWindow.context)
	return nil
}

// RemoveDeliveredNotification removes a delivered notification matching the unique identifier.
func (f *Frontend) RemoveDeliveredNotification(identifier string) error {
	cIdentifier := C.CString(identifier)
	defer C.free(unsafe.Pointer(cIdentifier))
	C.RemoveDeliveredNotification(f.mainWindow.context, cIdentifier)
	return nil
}

// RemoveNotification is a macOS stub that always returns nil.
// Use one of the following instead:
// RemoveAllPendingNotifications
// RemovePendingNotification
// RemoveAllDeliveredNotifications
// RemoveDeliveredNotification
// (Linux-specific)
func (f *Frontend) RemoveNotification(identifier string) error {
	return nil
}

func (f *Frontend) OnNotificationResponse(callback func(result frontend.NotificationResult)) {
	callbackLock.Lock()
	notificationResultCallback = callback
	callbackLock.Unlock()
}

//export captureResult
func captureResult(channelID C.int, success C.bool, errorMsg *C.char) {
	f := getCurrentFrontend()
	if f == nil {
		return
	}

	resultCh, exists := f.GetChannel(int(channelID))
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
}

//export didReceiveNotificationResponse
func didReceiveNotificationResponse(jsonPayload *C.char, err *C.char) {
	result := frontend.NotificationResult{}

	if err != nil {
		errMsg := C.GoString(err)
		result.Error = fmt.Errorf("notification response error: %s", errMsg)
		handleNotificationResult(result)

		return
	}

	if jsonPayload == nil {
		result.Error = fmt.Errorf("received nil JSON payload in notification response")
		handleNotificationResult(result)
		return
	}

	payload := C.GoString(jsonPayload)

	var response frontend.NotificationResponse
	if err := json.Unmarshal([]byte(payload), &response); err != nil {
		result.Error = fmt.Errorf("failed to unmarshal notification response: %w", err)
		handleNotificationResult(result)
		return
	}

	if response.ActionIdentifier == AppleDefaultActionIdentifier {
		response.ActionIdentifier = DefaultActionIdentifier
	}

	result.Response = response
	handleNotificationResult(result)
}

func handleNotificationResult(result frontend.NotificationResult) {
	callbackLock.Lock()
	callback := notificationResultCallback
	callbackLock.Unlock()

	if callback != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					// Log panic but don't crash the app
					fmt.Fprintf(os.Stderr, "panic in notification callback: %v\n", r)
				}
			}()
			callback(result)
		}()
	}
}

// Helper methods

func (f *Frontend) registerChannel() (int, chan notificationChannel) {
	channelsLock.Lock()
	defer channelsLock.Unlock()

	// Initialize channels map if it's nil
	if channels == nil {
		channels = make(map[int]chan notificationChannel)
		nextChannelID = 0
	}

	id := nextChannelID
	nextChannelID++

	resultCh := make(chan notificationChannel, 1)

	channels[id] = resultCh
	return id, resultCh
}

func (f *Frontend) GetChannel(id int) (chan notificationChannel, bool) {
	channelsLock.Lock()
	defer channelsLock.Unlock()

	if channels == nil {
		return nil, false
	}

	ch, exists := channels[id]
	if exists {
		delete(channels, id)
	}
	return ch, exists
}

func (f *Frontend) cleanupChannel(id int) {
	channelsLock.Lock()
	defer channelsLock.Unlock()

	if channels == nil {
		return
	}

	if ch, exists := channels[id]; exists {
		delete(channels, id)
		close(ch)
	}
}
