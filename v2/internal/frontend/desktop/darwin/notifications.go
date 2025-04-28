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
	"fmt"
	"sync"
	"time"
)

// Package-scoped variable only accessible within this file
var (
	currentFrontend *Frontend
	frontendMutex   sync.RWMutex
	// Notification channels
	channels      map[int]chan notificationChannel
	channelsLock  sync.Mutex
	nextChannelID int
)

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

type ChannelHandler interface {
	GetChannel(id int) (chan notificationChannel, bool)
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

	return nil
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

	setCurrentFrontend(f)
	defer setCurrentFrontend(nil)

	id, resultCh := f.registerChannel()

	C.RequestNotificationAuthorization(f.mainWindow.context, C.int(id))

	select {
	case result := <-resultCh:
		return result.Success, result.Error
	case <-ctx.Done():
		f.cleanupChannel(id)
		return false, fmt.Errorf("notification authorization timed out after 3 minutes: %w", ctx.Err())
	}
}

func (f *Frontend) CheckNotificationAuthorization() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	setCurrentFrontend(f)
	defer setCurrentFrontend(nil)

	id, resultCh := f.registerChannel()

	C.CheckNotificationAuthorization(f.mainWindow.context, C.int(id))

	select {
	case result := <-resultCh:
		return result.Success, result.Error
	case <-ctx.Done():
		f.cleanupChannel(id)
		return false, fmt.Errorf("notification authorization timed out after 15s: %w", ctx.Err())
	}
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
}

// Helper methods

func (f *Frontend) registerChannel() (int, chan notificationChannel) {
	channelsLock.Lock()
	defer channelsLock.Unlock()

	id := nextChannelID
	nextChannelID++

	resultCh := make(chan notificationChannel, 1)

	channels[id] = resultCh
	return id, resultCh
}

func (f *Frontend) GetChannel(id int) (chan notificationChannel, bool) {
	channelsLock.Lock()
	defer channelsLock.Unlock()

	ch, exists := channels[id]
	if exists {
		delete(channels, id)
	}
	return ch, exists
}

func (f *Frontend) cleanupChannel(id int) {
	channelsLock.Lock()
	defer channelsLock.Unlock()

	if ch, exists := channels[id]; exists {
		delete(channels, id)
		close(ch)
	}
}
