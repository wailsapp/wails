package notifications

type Service struct {
}

// NotificationAction represents an action button for a notification
type NotificationAction = struct {
	ID                     string `json:"id,omitempty"`
	Title                  string `json:"title,omitempty"`
	Destructive            bool   `json:"destructive,omitempty"`
	AuthenticationRequired bool   `json:"authenticationRequired,omitempty"`
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
	Subtitle   string                 `json:"subtitle,omitempty"`
	Body       string                 `json:"body,omitempty"`
	CategoryID string                 `json:"categoryId,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

// NotificationResponseData
type NotificationResponseData = struct {
	ID               string                 `json:"id,omitempty"`
	ActionIdentifier string                 `json:"actionIdentifier,omitempty"`
	CategoryID       string                 `json:"categoryIdentifier,omitempty"`
	Title            string                 `json:"title,omitempty"`
	Subtitle         string                 `json:"subtitle,omitempty"`
	Body             string                 `json:"body,omitempty"`
	UserText         string                 `json:"userText,omitempty"`
	UserInfo         map[string]interface{} `json:"userInfo,omitempty"`
}

// NotificationResponse
type NotificationResponse = struct {
	Name string                   `json:"name"`
	Data NotificationResponseData `json:"data"`
}
