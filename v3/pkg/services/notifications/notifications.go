package notifications

import "sync"

// Service represents the notifications service
type Service struct {
	// notificationResponseCallback is called when a notification response is received
	notificationResponseCallback func(response NotificationResponse)

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

// ServiceName returns the name of the service.
func (ns *Service) ServiceName() string {
	return "github.com/wailsapp/wails/v3/services/notifications"
}

// OnNotificationResponse registers a callback function that will be called when
// a notification response is received from the user.
//
//wails:ignore
func (ns *Service) OnNotificationResponse(callback func(response NotificationResponse)) {
	ns.callbackLock.Lock()
	defer ns.callbackLock.Unlock()

	ns.notificationResponseCallback = callback
}

// handleNotificationResponse is an internal method to handle notification responses
// and invoke the registered callback if one exists
func (ns *Service) handleNotificationResponse(response NotificationResponse) {
	ns.callbackLock.RLock()
	callback := ns.notificationResponseCallback
	ns.callbackLock.RUnlock()

	if callback != nil {
		callback(response)
	}
}
