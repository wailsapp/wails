// Package notifications provides cross-platform notification capabilities for desktop applications.
// It supports macOS, Windows, and Linux with a consistent API while handling platform-specific
// differences internally. Key features include:
//   - Basic notifications with title, subtitle, and body
//   - Interactive notifications with buttons and actions
//   - Notification categories for reusing configurations
//   - Custom sounds (default / silent / named)
//   - Attached media (images on every platform; audio/video on macOS)
//   - Threading / grouping by ThreadID
//   - Priority via InterruptionLevel (passive/active/timeSensitive/critical)
//   - Scheduled delivery (native on macOS; in-process timer on Windows + Linux)
//   - Updating an in-flight notification by ID
//   - User feedback handling with a unified callback system
//
// Platform-specific notes:
//   - macOS: Requires a properly bundled and signed application. Critical
//     interruption level requires the Critical Alert entitlement; without it
//     the level silently degrades.
//   - Windows: Uses Windows Toast notifications via the wintoast subpackage.
//     Reply fields are supported. Scheduled notifications use an in-process
//     timer and are lost if the app exits before delivery. Update-by-ID
//     redelivers as a new notification (true replace requires upstream
//     wintoast support for tag/group).
//   - Linux: Uses D-Bus org.freedesktop.Notifications. Reply fields are NOT
//     supported (not part of the spec). Subtitle is concatenated into the
//     body. Scheduled notifications use an in-process timer.
//
// See the NotificationOptions godoc and the notifications example app for
// the full per-feature support matrix.
package notifications

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// scheduleDelay returns the time.Duration until delivery for a Schedule.
// The bool is false if the schedule resolves to immediate delivery (nil
// schedule, zero values, or an At in the past). Used by Windows and Linux
// backends, which fall back to in-process time.AfterFunc timers because
// neither has a native deferred-delivery primitive exposed by the libraries
// we currently depend on (wintoast on Windows, godbus on Linux).
func scheduleDelay(s *NotificationSchedule) (time.Duration, bool) {
	if s == nil {
		return 0, false
	}
	if s.DelaySeconds > 0 {
		return time.Duration(s.DelaySeconds) * time.Second, true
	}
	if s.At > 0 {
		until := time.Until(time.Unix(s.At, 0))
		if until > 0 {
			return until, true
		}
	}
	return 0, false
}

type platformNotifier interface {
	// Lifecycle methods
	Startup(ctx context.Context, options application.ServiceOptions) error
	Shutdown() error

	// Core notification methods
	RequestNotificationAuthorization() (bool, error)
	CheckNotificationAuthorization() (bool, error)
	SendNotification(options NotificationOptions) error
	SendNotificationWithActions(options NotificationOptions) error
	UpdateNotification(options NotificationOptions) error

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

// NotificationOptions contains configuration for a notification.
//
// New optional fields (Sound, Attachments, ThreadID, InterruptionLevel,
// Schedule) gracefully degrade when a platform cannot honour them; see the
// package-level godoc for the per-platform support matrix.
type NotificationOptions struct {
	ID         string                 `json:"id"`
	Title      string                 `json:"title"`
	Subtitle   string                 `json:"subtitle,omitempty"` // (macOS and Linux only)
	Body       string                 `json:"body,omitempty"`
	CategoryID string                 `json:"categoryId,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`

	// Sound controls the sound played on delivery.
	//   nil                                -> platform default sound
	//   &NotificationSound{Silent: true}   -> no sound
	//   &NotificationSound{Name: "Ping"}   -> named/bundled sound
	// On macOS, Name is resolved by [UNNotificationSound soundNamed:] and
	// requires the audio file to live under the bundle's Library/Sounds.
	// On Windows, Name is used as-is if it begins with "ms-winsoundevent:" or
	// "ms-appx:"; otherwise it is wrapped in "ms-winsoundevent:" for built-in
	// event names. On Linux it is forwarded as the freedesktop "sound-name"
	// hint (theme-dependent).
	Sound *NotificationSound `json:"sound,omitempty"`

	// Attachments are media files shown alongside the notification. macOS
	// supports multiple attachments of any media type; Windows and Linux
	// honour the first image-typed attachment (Linux limits to one per spec).
	Attachments []NotificationAttachment `json:"attachments,omitempty"`

	// ThreadID groups related notifications together in Notification Center
	// (macOS) / Action Center (Windows) / the notification daemon (Linux).
	ThreadID string `json:"threadId,omitempty"`

	// InterruptionLevel controls notification priority. One of "passive",
	// "active" (default), "timeSensitive", "critical". Critical requires
	// macOS 12+ and the Critical Alert entitlement. Linux maps to the
	// freedesktop urgency hint; Windows maps to <toast scenario="...">.
	InterruptionLevel string `json:"interruptionLevel,omitempty"`

	// Schedule defers delivery. macOS uses a native trigger and persists
	// across app restarts. Windows and Linux fall back to an in-process
	// time.AfterFunc timer that does NOT survive an app exit.
	Schedule *NotificationSchedule `json:"schedule,omitempty"`
}

// NotificationSound configures audio playback for a notification.
type NotificationSound struct {
	Silent bool   `json:"silent,omitempty"`
	Name   string `json:"name,omitempty"`
}

// NotificationAttachment is a media file shown with the notification.
// Path is an absolute filesystem path (or "file://" URL on macOS).
type NotificationAttachment struct {
	ID   string `json:"id,omitempty"`
	Path string `json:"path"`
	// Type is an optional placement/UTI hint.
	//   On macOS: a UTI like "public.png" / "public.audio" (often inferred).
	//   On Windows: "hero" | "appLogoOverride" | "inline" (default "inline").
	//   On Linux: ignored (always image-path hint).
	Type string `json:"type,omitempty"`
}

// NotificationSchedule defers delivery. Exactly one of DelaySeconds or At
// must be set. At is interpreted as Unix seconds (UTC).
type NotificationSchedule struct {
	DelaySeconds int   `json:"delaySeconds,omitempty"`
	At           int64 `json:"at,omitempty"`
}

// Allowed values for NotificationOptions.InterruptionLevel.
const (
	InterruptionLevelPassive       = "passive"
	InterruptionLevelActive        = "active"
	InterruptionLevelTimeSensitive = "timeSensitive"
	InterruptionLevelCritical      = "critical"
)

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

// UpdateNotification updates an in-flight notification by ID. On macOS this
// is auto-deduplicated by UNUserNotificationCenter; on Linux it uses the
// D-Bus replaces_id parameter. On Windows it currently redelivers as a new
// notification (true replace requires upstream wintoast support for tag/group).
func (ns *NotificationService) UpdateNotification(options NotificationOptions) error {
	if err := validateNotificationOptions(options); err != nil {
		return err
	}
	return ns.impl.UpdateNotification(options)
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

// validateNotificationOptions validates required fields and the shape of the
// new optional fields on NotificationOptions.
func validateNotificationOptions(options NotificationOptions) error {
	if options.ID == "" {
		return fmt.Errorf("notification ID cannot be empty")
	}

	if options.Title == "" {
		return fmt.Errorf("notification title cannot be empty")
	}

	if options.InterruptionLevel != "" {
		switch options.InterruptionLevel {
		case InterruptionLevelPassive, InterruptionLevelActive,
			InterruptionLevelTimeSensitive, InterruptionLevelCritical:
		default:
			return fmt.Errorf("invalid interruption level %q (want passive|active|timeSensitive|critical)", options.InterruptionLevel)
		}
	}

	if options.Schedule != nil {
		if options.Schedule.DelaySeconds < 0 {
			return fmt.Errorf("schedule.delaySeconds cannot be negative")
		}
		if options.Schedule.At < 0 {
			return fmt.Errorf("schedule.at cannot be negative")
		}
		if options.Schedule.DelaySeconds > 0 && options.Schedule.At > 0 {
			return fmt.Errorf("schedule.delaySeconds and schedule.at are mutually exclusive")
		}
		if options.Schedule.DelaySeconds == 0 && options.Schedule.At == 0 {
			return fmt.Errorf("schedule must set either delaySeconds or at")
		}
	}

	for i, a := range options.Attachments {
		if a.Path == "" {
			return fmt.Errorf("attachments[%d].path cannot be empty", i)
		}
		statPath := a.Path
		if strings.HasPrefix(statPath, "file://") {
			// macOS UNNotificationAttachment accepts file:// URLs as well as
			// plain paths. Map back to a local filesystem path so os.Stat
			// doesn't reject the URL form as a literal path.
			//
			// url.Parse rejects Windows backslash paths ("file://C:\path")
			// because it treats the drive letter as a hostname and the rest
			// as an invalid port. Fall back to prefix-stripping in that case.
			u, err := url.Parse(statPath)
			if err != nil {
				// Windows backslash fallback: strip prefix and let filepath
				// normalise the remaining native path.
				p := filepath.FromSlash(strings.TrimPrefix(statPath, "file://"))
				if !filepath.IsAbs(p) {
					return fmt.Errorf("attachments[%d].path: file:// URL must be absolute", i)
				}
				statPath = p
			} else {
				// Reject non-empty hosts other than localhost. A URL like
				// "file://relative/image.png" parses with Host="relative" and
				// Path="/image.png" — it looks absolute but is not a valid local
				// file reference. Only "file:///..." (empty host) and
				// "file://localhost/..." are valid local file:// URLs.
				if u.Host != "" && u.Host != "localhost" {
					return fmt.Errorf("attachments[%d].path: file:// URL must be absolute (use file:///path)", i)
				}
				// Standard three-slash form (file:///C:/path or file:///tmp/path).
				// u.Path is already the URL-decoded path component. On Windows,
				// file:///C:/path parses to path "/C:/path" — drop the leading
				// separator before the drive letter.
				p := filepath.FromSlash(u.Path)
				if len(p) > 2 && p[0] == filepath.Separator && p[2] == ':' {
					p = p[1:]
				}
				if !filepath.IsAbs(p) {
					return fmt.Errorf("attachments[%d].path: file:// URL must be absolute", i)
				}
				statPath = p
			}
		} else if !filepath.IsAbs(statPath) {
			return fmt.Errorf("attachments[%d].path must be an absolute path", i)
		}
		if _, err := os.Stat(statPath); err != nil {
			return fmt.Errorf("attachments[%d].path is not accessible", i)
		}
	}

	return nil
}
