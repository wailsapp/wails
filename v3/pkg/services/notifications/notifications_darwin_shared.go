//go:build darwin && !ios

package notifications

import (
	"context"
	"fmt"
	"time"
)

func (dn *darwinNotifier) GetChannel(id int) (chan notificationChannel, bool) {
	dn.channelsLock.Lock()
	defer dn.channelsLock.Unlock()

	ch, exists := dn.channels[id]
	if exists {
		delete(dn.channels, id)
	}
	return ch, exists
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

// SendNotification sends a basic notification with a unique identifier, title, subtitle, and body.
func (dn *darwinNotifier) SendNotification(options NotificationOptions) error {
	return dn.dispatchNotification(options, false)
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
// A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
// If a NotificationCategory is not registered a basic notification will be sent.
func (dn *darwinNotifier) SendNotificationWithActions(options NotificationOptions) error {
	return dn.dispatchNotification(options, true)
}

func (dn *darwinNotifier) Shutdown() error {
	return nil
}

// UpdateNotification re-posts a notification with the same identifier.
// UNUserNotificationCenter auto-replaces by identifier, so the existing
// notification's content is updated in place (or freshly delivered if no
// notification with that identifier is currently pending or delivered).
func (dn *darwinNotifier) UpdateNotification(options NotificationOptions) error {
	return dn.dispatchNotification(options, options.CategoryID != "")
}

func (dn *darwinNotifier) cleanupChannel(id int) {
	dn.channelsLock.Lock()
	defer dn.channelsLock.Unlock()

	if ch, exists := dn.channels[id]; exists {
		delete(dn.channels, id)
		close(ch)
	}
}

func (dn *darwinNotifier) registerChannel() (int, chan notificationChannel) {
	dn.channelsLock.Lock()
	defer dn.channelsLock.Unlock()

	id := dn.nextChannelID
	dn.nextChannelID++

	resultCh := make(chan notificationChannel, 1)

	dn.channels[id] = resultCh
	return id, resultCh
}

// Creates a new Notifications Service.
// Your app must be packaged and signed for this feature to work.
func New() *NotificationService {
	notificationServiceOnce.Do(func() {
		impl := &darwinNotifier{
			channels:      make(map[int]chan notificationChannel),
			nextChannelID: 0,
		}

		NotificationService_ = &NotificationService{
			impl: impl,
		}
	})

	return NotificationService_
}

func (dn *darwinNotifier) CheckNotificationAuthorization() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	id, resultCh := dn.registerChannel()

	checkNotificationAuthorization(id)

	select {
	case result := <-resultCh:
		return result.Success, result.Error
	case <-ctx.Done():
		dn.cleanupChannel(id)
		return false, fmt.Errorf("notification authorization timed out after 15s: %w", ctx.Err())
	}
}

func (dn *darwinNotifier) RequestNotificationAuthorization() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	id, resultCh := dn.registerChannel()

	requestNotificationAuthorization(id)

	select {
	case result := <-resultCh:
		return result.Success, result.Error
	case <-ctx.Done():
		dn.cleanupChannel(id)
		return false, fmt.Errorf("notification authorization timed out after 3 minutes: %w", ctx.Err())
	}
}

func (dn *darwinNotifier) RemoveNotificationCategory(categoryId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, resultCh := dn.registerChannel()
	removeNotificationCategory(id, categoryId)

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
