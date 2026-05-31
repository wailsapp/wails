// Package notifications provides cross-platform notification capabilities for desktop applications.
// It supports macOS, Windows, and Linux with a consistent API while handling platform-specific
// differences internally. Key features include:
//   - Basic notifications with title, subtitle, and body
//   - Interactive notifications with buttons and actions
//   - Notification categories for reusing configurations
//   - User feedback handling with a unified callback system
//
// Platform-specific notes:
//   - macOS: Requires a properly bundled and signed application
//   - Windows: Uses Windows Toast notifications
//   - Linux: Uses D-Bus and does not support text inputs
package notifications

import (
	"context"
	"fmt"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type platformNotifier interface {
	// Lifecycle methods
	Startup(ctx context.Context, options application.ServiceOptions) error
	Shutdown() error

	// Core notification methods
	RequestNotificationAuthorization() (bool, error)
	CheckNotificationAuthorization() (bool, error)
	SendNotification(options NotificationOptions) error
	SendNotificationWithActions(options NotificationOptions) error

	// Category management
	RegisterNotificationCategory(category NotificationCategory) error
	RemoveNotificationCategory(categoryID string) error

	// Notification management
	RemoveAllPendingNotifications() error
	RemovePendingNotification(identifier string) error
	RemoveAllDeliveredNotifications() error
	RemoveDeliveredNotification(identifier string) error
	RemoveNotification(identifier string) error
}

// Service represents the notifications service
type NotificationService struct {
	impl platformNotifier

	// notificationResponseCallback is called when a notification result is received.
	// Only one callback can be assigned at a time.
	notificationResultCallback func(result NotificationResult)

	callbackLock sync.RWMutex
}

var (
	notificationServiceOnce sync.Once
	NotificationService_    *NotificationService
	notificationServiceLock sync.RWMutex
)

// NotificationAction represents an action button for a notification.
type NotificationAction struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Destructive bool   `json:"destructive,omitempty"` // (macOS-specific)
}

// NotificationCategory groups actions for notifications.
type NotificationCategory struct {
	ID               string               `json:"id,omitempty"`
	Actions          []NotificationAction `json:"actions,omitempty"`
	HasReplyField    bool                 `json:"hasReplyField,omitempty"`
	ReplyPlaceholder string               `json:"replyPlaceholder,omitempty"`
	ReplyButtonTitle string               `json:"replyButtonTitle,omitempty"`
}

// NotificationOptions contains configuration for a notification
type NotificationOptions struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	Subtitle   string                 `json:"subtitle,omitempty"` // (macOS and Linux only)
	Body       string                 `json:"body,omitempty"`
	CategoryID string                 `json:"categoryId,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

const DefaultActionIdentifier = "DEFAULT_ACTION"

// NotificationResponse represents the response sent by interacting with a notification.
type NotificationResponse struct {
	ID               string                 `json:"id,omitempty"`
	ActionIdentifier string                 `json:"actionIdentifier,omitempty"`
	CategoryID       string                 `json:"categoryIdentifier,omitempty"`
	Title            string                 `json:"title,omitempty"`
	Subtitle         string                 `json:"subtitle,omitempty"` // (macOS and Linux only)
	Body             string                 `json:"body,omitempty"`
	UserText         string                 `json:"userText,omitempty"`
	UserInfo         map[string]interface{} `json:"userInfo,omitempty"`
}

// NotificationResult represents the result of a notification response,
// returning the response or any errors that occurred.
type NotificationResult struct {
	Response NotificationResponse
	Error    error
}

// ServiceName returns the name of the service.
func (ns *NotificationService) ServiceName() string {
	return "github.com/wailsapp/wails/v3/services/notifications"
}

// OnNotificationResponse registers a callback function that will be called when
// a notification response is received from the user.
//
//wails:ignore
func (ns *NotificationService) OnNotificationResponse(callback func(result NotificationResult)) {
	ns.callbackLock.Lock()
	defer ns.callbackLock.Unlock()

	ns.notificationResultCallback = callback
}

// handleNotificationResponse is an internal method to handle notification responses
// and invoke the registered callback if one exists.
func (ns *NotificationService) handleNotificationResult(result NotificationResult) {
	ns.callbackLock.RLock()
	callback := ns.notificationResultCallback
	ns.callbackLock.RUnlock()

	if callback != nil {
		callback(result)
	}
}

// ServiceStartup is called when the service is loaded.
func (ns *NotificationService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	return ns.impl.Startup(ctx, options)
}

// ServiceShutdown is called when the service is unloaded.
func (ns *NotificationService) ServiceShutdown() error {
	return ns.impl.Shutdown()
}

// Public methods that delegate to the implementation.
func (ns *NotificationService) RequestNotificationAuthorization() (bool, error) {
	return ns.impl.RequestNotificationAuthorization()
}

func (ns *NotificationService) CheckNotificationAuthorization() (bool, error) {
	return ns.impl.CheckNotificationAuthorization()
}

func (ns *NotificationService) SendNotification(options NotificationOptions) error {
	if err := validateNotificationOptions(options); err != nil {
		return err
	}
	return ns.impl.SendNotification(options)
}

func (ns *NotificationService) SendNotificationWithActions(options NotificationOptions) error {
	if err := validateNotificationOptions(options); err != nil {
		return err
	}
	return ns.impl.SendNotificationWithActions(options)
}

func (ns *NotificationService) RegisterNotificationCategory(category NotificationCategory) error {
	return ns.impl.RegisterNotificationCategory(category)
}

func (ns *NotificationService) RemoveNotificationCategory(categoryID string) error {
	return ns.impl.RemoveNotificationCategory(categoryID)
}

func (ns *NotificationService) RemoveAllPendingNotifications() error {
	return ns.impl.RemoveAllPendingNotifications()
}

func (ns *NotificationService) RemovePendingNotification(identifier string) error {
	return ns.impl.RemovePendingNotification(identifier)
}

func (ns *NotificationService) RemoveAllDeliveredNotifications() error {
	return ns.impl.RemoveAllDeliveredNotifications()
}

func (ns *NotificationService) RemoveDeliveredNotification(identifier string) error {
	return ns.impl.RemoveDeliveredNotification(identifier)
}

func (ns *NotificationService) RemoveNotification(identifier string) error {
	return ns.impl.RemoveNotification(identifier)
}

func getNotificationService() *NotificationService {
	notificationServiceLock.RLock()
	defer notificationServiceLock.RUnlock()
	return NotificationService_
}

// validateNotificationOptions validates an ID and Title are provided for notifications.
func validateNotificationOptions(options NotificationOptions) error {
	if options.ID == "" {
		return fmt.Errorf("notification ID cannot be empty")
	}

	if options.Title == "" {
		return fmt.Errorf("notification title cannot be empty")
	}

	return nil
}
