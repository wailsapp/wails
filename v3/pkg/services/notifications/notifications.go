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
//   - Linux: Falls back between D-Bus, notify-send, or other methods and does not support text inputs
package notifications

import (
	"fmt"
	"sync"
)

// Service represents the notifications service
type Service struct {
	// notificationResponseCallback is called when a notification result is received.
	// Only one callback can be assigned at a time.
	notificationResultCallback func(result NotificationResult)

	callbackLock sync.RWMutex
}

var NotificationService *Service
var notificationServiceLock sync.RWMutex

// NotificationAction represents an action button for a notification
type NotificationAction = struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Destructive bool   `json:"destructive,omitempty"` // (macOS-specific)
}

// NotificationCategory groups actions for notifications
type NotificationCategory = struct {
	ID               string               `json:"id,omitempty"`
	Actions          []NotificationAction `json:"actions,omitempty"`
	HasReplyField    bool                 `json:"hasReplyField,omitempty"`
	ReplyPlaceholder string               `json:"replyPlaceholder,omitempty"`
	ReplyButtonTitle string               `json:"replyButtonTitle,omitempty"`
}

// NotificationOptions contains configuration for a notification
type NotificationOptions = struct {
	ID         string                 `json:"id,omitempty"`
	Title      string                 `json:"title,omitempty"`
	Subtitle   string                 `json:"subtitle,omitempty"` // (macOS-specific)
	Body       string                 `json:"body,omitempty"`
	CategoryID string                 `json:"categoryId,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

var DefaultActionIdentifier = "DEFAULT_ACTION"

// NotificationResponse represents a user's response to a notification
type NotificationResponse = struct {
	ID               string                 `json:"id,omitempty"`
	ActionIdentifier string                 `json:"actionIdentifier,omitempty"`
	CategoryID       string                 `json:"categoryIdentifier,omitempty"`
	Title            string                 `json:"title,omitempty"`
	Subtitle         string                 `json:"subtitle,omitempty"` // (macOS-specific)
	Body             string                 `json:"body,omitempty"`
	UserText         string                 `json:"userText,omitempty"`
	UserInfo         map[string]interface{} `json:"userInfo,omitempty"`
}

// NotificationResult
type NotificationResult = struct {
	Response NotificationResponse
	Error    error
}

// ServiceName returns the name of the service.
func (ns *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/services/notifications"
}

func getNotificationService() *Service {
	notificationServiceLock.RLock()
	defer notificationServiceLock.RUnlock()
	return NotificationService
}

// OnNotificationResponse registers a callback function that will be called when
// a notification response is received from the user.
//
//wails:ignore
func (ns *Service) OnNotificationResponse(callback func(result NotificationResult)) {
	ns.callbackLock.Lock()
	defer ns.callbackLock.Unlock()

	ns.notificationResultCallback = callback
}

// handleNotificationResponse is an internal method to handle notification responses
// and invoke the registered callback if one exists
func (ns *Service) handleNotificationResult(result NotificationResult) {
	ns.callbackLock.RLock()
	callback := ns.notificationResultCallback
	ns.callbackLock.RUnlock()

	if callback != nil {
		callback(result)
	}
}

func validateNotificationOptions(options NotificationOptions) error {
	if options.ID == "" {
		return fmt.Errorf("notification ID cannot be empty")
	}

	if options.Title == "" {
		return fmt.Errorf("notification title cannot be empty")
	}

	return nil
}
