//go:build linux

package notifications

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type linuxNotifier struct {
	categories          map[string]NotificationCategory
	categoriesLock      sync.RWMutex
	appName             string
	internal            *internalNotifier
	notificationInitErr error
	initOnce            sync.Once
}

const (
	dbusObjectPath             = "/org/freedesktop/Notifications"
	dbusNotificationsInterface = "org.freedesktop.Notifications"
	signalNotificationClosed   = "org.freedesktop.Notifications.NotificationClosed"
	signalActionInvoked        = "org.freedesktop.Notifications.ActionInvoked"
	callGetCapabilities        = "org.freedesktop.Notifications.GetCapabilities"
	callCloseNotification      = "org.freedesktop.Notifications.CloseNotification"

	MethodNotifySend = "notify-send"
	MethodDbus       = "dbus"
	MethodKdialog    = "kdialog"

	notifyChannelBufferSize = 25
)

type closedReason uint32

// internalNotifier handles the actual notification sending via dbus or command line
type notificationContext struct {
	ID       string
	SystemID uint32
	Actions  map[string]string      // Maps action keys to display labels
	UserData map[string]interface{} // The original user data
}

type internalNotifier struct {
	sync.Mutex
	method          string
	dbusConn        *dbus.Conn
	sendPath        string
	activeNotifs    map[string]uint32               // Maps our notification IDs to system IDs
	contexts        map[string]*notificationContext // Stores notification contexts by our ID
	listenerCtx     context.Context
	listenerCancel  context.CancelFunc
	listenerRunning bool
}

// New creates a new Notifications Service
func New() *Service {
	notificationServiceOnce.Do(func() {
		impl := &linuxNotifier{
			categories: make(map[string]NotificationCategory),
		}

		NotificationService = &Service{
			impl: impl,
		}
	})
	return NotificationService
}

// Startup is called when the service is loaded
func (ls *linuxNotifier) Startup(ctx context.Context) error {
	ls.appName = application.Get().Config().Name

	if err := ls.loadCategories(); err != nil {
		fmt.Printf("Failed to load notification categories: %v\n", err)
	}

	// Initialize the internal notifier
	ls.internal = &internalNotifier{
		activeNotifs: make(map[string]uint32),
		contexts:     make(map[string]*notificationContext),
	}

	var err error
	ls.initOnce.Do(func() {
		// Initialize notification system
		err = ls.initNotificationSystem()
	})

	return err
}

// Shutdown is called when the service is unloaded
func (ls *linuxNotifier) Shutdown() error {
	if ls.internal != nil {
		ls.internal.Lock()
		defer ls.internal.Unlock()

		// Cancel the listener context if it's running
		if ls.internal.listenerCancel != nil {
			ls.internal.listenerCancel()
			ls.internal.listenerCancel = nil
		}

		// Close the connection
		if ls.internal.dbusConn != nil {
			ls.internal.dbusConn.Close()
			ls.internal.dbusConn = nil
		}

		// Clear state
		ls.internal.activeNotifs = make(map[string]uint32)
		ls.internal.contexts = make(map[string]*notificationContext)
		ls.internal.method = "none"
		ls.internal.sendPath = ""
	}

	return ls.saveCategories()
}

// initNotificationSystem initializes the notification system, choosing the best available method
func (ls *linuxNotifier) initNotificationSystem() error {
	var err error

	// Cancel any existing listener before starting a new one
	if ls.internal.listenerCancel != nil {
		ls.internal.listenerCancel()
	}

	// Create a new context for the listener
	ls.internal.listenerCtx, ls.internal.listenerCancel = context.WithCancel(context.Background())

	// Reset state
	ls.internal.activeNotifs = make(map[string]uint32)
	ls.internal.contexts = make(map[string]*notificationContext)
	ls.internal.listenerRunning = false

	checkDbus := func() (*dbus.Conn, error) {
		conn, err := dbus.SessionBusPrivate()
		if err != nil {
			return conn, err
		}

		if err = conn.Auth(nil); err != nil {
			return conn, err
		}

		if err = conn.Hello(); err != nil {
			return conn, err
		}

		obj := conn.Object(dbusNotificationsInterface, dbusObjectPath)
		call := obj.Call(callGetCapabilities, 0)
		if call.Err != nil {
			return conn, call.Err
		}

		var ret []string
		err = call.Store(&ret)
		if err != nil {
			return conn, err
		}

		// Add a listener for notification signals
		err = conn.AddMatchSignal(
			dbus.WithMatchObjectPath(dbusObjectPath),
			dbus.WithMatchInterface(dbusNotificationsInterface),
		)
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	// Try dbus first
	ls.internal.dbusConn, err = checkDbus()
	if err == nil {
		ls.internal.method = MethodDbus
		// Start the dbus signal listener with context
		go ls.startDBusListener(ls.internal.listenerCtx)
		ls.internal.listenerRunning = true
		return nil
	}
	if ls.internal.dbusConn != nil {
		ls.internal.dbusConn.Close()
		ls.internal.dbusConn = nil
	}

	// Try notify-send
	send, err := exec.LookPath("notify-send")
	if err == nil {
		ls.internal.sendPath = send
		ls.internal.method = MethodNotifySend
		return nil
	}

	// Try sw-notify-send
	send, err = exec.LookPath("sw-notify-send")
	if err == nil {
		ls.internal.sendPath = send
		ls.internal.method = MethodNotifySend
		return nil
	}

	// No method available
	ls.internal.method = "none"
	ls.internal.sendPath = ""

	return errors.New("no notification method is available")
}

// startDBusListener listens for DBus signals for notification actions and closures
func (ls *linuxNotifier) startDBusListener(ctx context.Context) {
	signal := make(chan *dbus.Signal, notifyChannelBufferSize)
	ls.internal.dbusConn.Signal(signal)

	defer func() {
		ls.internal.Lock()
		ls.internal.listenerRunning = false
		ls.internal.Unlock()
		ls.internal.dbusConn.RemoveSignal(signal) // Remove signal handler
		close(signal)                             // Clean up channel
	}()

	for {
		select {
		case <-ctx.Done():
			// Context was cancelled, exit gracefully
			return

		case s := <-signal:
			if s == nil {
				// Channel closed or nil signal
				continue
			}

			if len(s.Body) < 2 {
				continue
			}

			switch s.Name {
			case signalNotificationClosed:
				systemID := s.Body[0].(uint32)
				reason := closedReason(s.Body[1].(uint32)).string()
				ls.handleNotificationClosed(systemID, reason)
			case signalActionInvoked:
				systemID := s.Body[0].(uint32)
				actionKey := s.Body[1].(string)
				ls.handleActionInvoked(systemID, actionKey)
			}
		}
	}
}

// handleNotificationClosed processes notification closed signals
func (ls *linuxNotifier) handleNotificationClosed(systemID uint32, reason string) {
	// Find our notification ID for this system ID
	var notifID string
	var userData map[string]interface{}

	ls.internal.Lock()
	for id, sysID := range ls.internal.activeNotifs {
		if sysID == systemID {
			notifID = id
			// Get the user data from context if available
			if ctx, exists := ls.internal.contexts[id]; exists {
				userData = ctx.UserData
			}
			break
		}
	}
	ls.internal.Unlock()

	if notifID != "" {
		response := NotificationResponse{
			ID:               notifID,
			ActionIdentifier: DefaultActionIdentifier,
			UserInfo:         userData,
		}

		// Add reason to UserInfo or create it if none exists
		if response.UserInfo == nil {
			response.UserInfo = map[string]interface{}{
				"reason": reason,
			}
		} else {
			response.UserInfo["reason"] = reason
		}

		result := NotificationResult{}
		result.Response = response
		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(result)
		}

		// Clean up the context
		ls.internal.Lock()
		delete(ls.internal.contexts, notifID)
		delete(ls.internal.activeNotifs, notifID)
		ls.internal.Unlock()
	}
}

// handleActionInvoked processes action invoked signals
func (ls *linuxNotifier) handleActionInvoked(systemID uint32, actionKey string) {
	// Find our notification ID and context for this system ID
	var notifID string
	var ctx *notificationContext

	ls.internal.Lock()
	for id, sysID := range ls.internal.activeNotifs {
		if sysID == systemID {
			notifID = id
			ctx = ls.internal.contexts[id]
			break
		}
	}
	ls.internal.Unlock()

	if notifID != "" {
		if actionKey == "default" {
			actionKey = DefaultActionIdentifier
		}

		// First, send the action response with the user data
		response := NotificationResponse{
			ID:               notifID,
			ActionIdentifier: actionKey,
		}

		// Include the user data if we have it
		if ctx != nil {
			response.UserInfo = ctx.UserData
		}

		result := NotificationResult{}
		result.Response = response
		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(result)
		}

		// Then, trigger a closed event with "activated-by-user" reason
		closeResponse := NotificationResponse{
			ID:               notifID,
			ActionIdentifier: DefaultActionIdentifier,
		}

		// Include the same user data in the close response
		if ctx != nil {
			closeResponse.UserInfo = ctx.UserData
		} else {
			closeResponse.UserInfo = map[string]interface{}{}
		}

		// Add the reason to the user info
		closeResponse.UserInfo["reason"] = closedReason(5).string() // "activated-by-user"

		closeResult := NotificationResult{}
		closeResult.Response = closeResponse
		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(closeResult)
		}

		// Clean up the context
		ls.internal.Lock()
		delete(ls.internal.contexts, notifID)
		delete(ls.internal.activeNotifs, notifID)
		ls.internal.Unlock()
	}
}

// CheckBundleIdentifier is a Linux stub that always returns true.
// (bundle identifiers are macOS-specific)
func (ls *linuxNotifier) CheckBundleIdentifier() bool {
	return true
}

// RequestNotificationAuthorization is a Linux stub that always returns true.
// (user authorization is macOS-specific)
func (ls *linuxNotifier) RequestNotificationAuthorization() (bool, error) {
	return true, nil
}

// CheckNotificationAuthorization is a Linux stub that always returns true.
// (user authorization is macOS-specific)
func (ls *linuxNotifier) CheckNotificationAuthorization() (bool, error) {
	return true, nil
}

// SendNotification sends a basic notification with a unique identifier, title, subtitle, and body.
func (ls *linuxNotifier) SendNotification(options NotificationOptions) error {
	if ls.internal == nil {
		return errors.New("notification service not initialized")
	}

	if ls.internal.method == "" || (ls.internal.method == MethodDbus && ls.internal.dbusConn == nil) ||
		(ls.internal.method == MethodNotifySend && ls.internal.sendPath == "") {
		return errors.New("notification system not properly initialized")
	}

	ls.internal.Lock()
	defer ls.internal.Unlock()

	var (
		systemID uint32
		err      error
	)

	switch ls.internal.method {
	case MethodDbus:
		systemID, err = ls.sendViaDbus(options, nil)
	case MethodNotifySend:
		systemID, err = ls.sendViaNotifySend(options)
	default:
		err = errors.New("no notification method is available")
	}

	if err == nil && systemID > 0 {
		// Store the system ID mapping
		ls.internal.activeNotifs[options.ID] = systemID

		// Create and store the notification context
		ctx := &notificationContext{
			ID:       options.ID,
			SystemID: systemID,
			UserData: options.Data,
		}
		ls.internal.contexts[options.ID] = ctx
	}

	return err
}

// SendNotificationWithActions sends a notification with additional actions.
func (ls *linuxNotifier) SendNotificationWithActions(options NotificationOptions) error {
	if ls.internal == nil {
		return errors.New("notification service not initialized")
	}

	if ls.internal.method == "" || (ls.internal.method == MethodDbus && ls.internal.dbusConn == nil) ||
		(ls.internal.method == MethodNotifySend && ls.internal.sendPath == "") {
		return errors.New("notification system not properly initialized")
	}

	ls.categoriesLock.RLock()
	category, exists := ls.categories[options.CategoryID]
	ls.categoriesLock.RUnlock()

	if !exists {
		return ls.SendNotification(options)
	}

	ls.internal.Lock()
	defer ls.internal.Unlock()

	var (
		systemID uint32
		err      error
	)

	switch ls.internal.method {
	case MethodDbus:
		systemID, err = ls.sendViaDbus(options, &category)
	case MethodNotifySend:
		// notify-send doesn't support actions, fall back to basic notification
		systemID, err = ls.sendViaNotifySend(options)
	default:
		err = errors.New("no notification method is available")
	}

	if err == nil && systemID > 0 {
		// Store the system ID mapping
		ls.internal.activeNotifs[options.ID] = systemID

		// Create and store the notification context with actions
		ctx := &notificationContext{
			ID:       options.ID,
			SystemID: systemID,
			UserData: options.Data,
			Actions:  make(map[string]string),
		}

		// Store action mappings
		if exists {
			for _, action := range category.Actions {
				ctx.Actions[action.ID] = action.Title
			}
		}

		ls.internal.contexts[options.ID] = ctx
	}

	return err
}

// sendViaDbus sends a notification via dbus
func (ls *linuxNotifier) sendViaDbus(options NotificationOptions, category *NotificationCategory) (result uint32, err error) {
	// Prepare actions
	var actions []string
	if category != nil {
		for _, action := range category.Actions {
			actions = append(actions, action.ID, action.Title)
		}
	}

	// Default timeout (-1 means use system default)
	timeout := int32(-1)

	// Prepare hints
	hints := map[string]dbus.Variant{
		// Normal urgency by default
		"urgency": dbus.MakeVariant(byte(1)),
	}

	// Add user data to hints if available
	if options.Data != nil {
		if userData, err := json.Marshal(options.Data); err == nil {
			hints["x-wails-user-data"] = dbus.MakeVariant(string(userData))
		}
	}

	// Send the notification
	obj := ls.internal.dbusConn.Object(dbusNotificationsInterface, dbusObjectPath)
	dbusArgs := []interface{}{
		ls.appName,    // App name
		uint32(0),     // Replaces ID (0 means new notification)
		"",            // App icon (empty for now)
		options.Title, // Title
		options.Body,  // Body
		actions,       // Actions
		hints,         // Hints
		timeout,       // Timeout
	}

	call := obj.Call("org.freedesktop.Notifications.Notify", 0, dbusArgs...)
	if call.Err != nil {
		return 0, fmt.Errorf("dbus notification error: %v", call.Err)
	}

	err = call.Store(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// sendViaNotifySend sends a notification via notify-send command
func (ls *linuxNotifier) sendViaNotifySend(options NotificationOptions) (uint32, error) {
	args := []string{
		options.Title,
		options.Body,
	}

	// Add icon if eventually supported
	// if options.Icon != "" { ... }

	// Add urgency (normal by default)
	args = append(args, "--urgency=normal")

	// Execute the command
	cmd := exec.Command(ls.internal.sendPath, args...)
	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("notify-send error: %v", err)
	}

	// notify-send doesn't return IDs, so we use 0
	return 0, nil
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
func (ls *linuxNotifier) RegisterNotificationCategory(category NotificationCategory) error {
	ls.categoriesLock.Lock()
	ls.categories[category.ID] = category
	ls.categoriesLock.Unlock()

	return ls.saveCategories()
}

// RemoveNotificationCategory removes a previously registered NotificationCategory.
func (ls *linuxNotifier) RemoveNotificationCategory(categoryId string) error {
	ls.categoriesLock.Lock()
	delete(ls.categories, categoryId)
	ls.categoriesLock.Unlock()

	return ls.saveCategories()
}

// RemoveAllPendingNotifications is a Linux stub that always returns nil.
// (macOS-specific)
func (ls *linuxNotifier) RemoveAllPendingNotifications() error {
	return nil
}

// RemovePendingNotification is a Linux stub that always returns nil.
// (macOS-specific)
func (ls *linuxNotifier) RemovePendingNotification(_ string) error {
	return nil
}

// RemoveAllDeliveredNotifications is a Linux stub that always returns nil.
// (macOS-specific)
func (ls *linuxNotifier) RemoveAllDeliveredNotifications() error {
	return nil
}

// RemoveDeliveredNotification is a Linux stub that always returns nil.
// (macOS-specific)
func (ls *linuxNotifier) RemoveDeliveredNotification(_ string) error {
	return nil
}

// RemoveNotification removes a notification by ID (Linux-specific)
func (ls *linuxNotifier) RemoveNotification(identifier string) error {
	if ls.internal == nil || ls.internal.method != MethodDbus || ls.internal.dbusConn == nil {
		return errors.New("dbus not available for closing notifications")
	}

	// Get the system ID for this notification
	ls.internal.Lock()
	systemID, exists := ls.internal.activeNotifs[identifier]
	ls.internal.Unlock()

	if !exists {
		return nil // Already closed or unknown
	}

	// Call CloseNotification on dbus
	obj := ls.internal.dbusConn.Object(dbusNotificationsInterface, dbusObjectPath)
	call := obj.Call(callCloseNotification, 0, systemID)

	return call.Err
}

// getConfigFilePath returns the path to the configuration file for storing notification categories
func (ls *linuxNotifier) getConfigFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %v", err)
	}

	appConfigDir := filepath.Join(configDir, ls.appName)
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %v", err)
	}

	return filepath.Join(appConfigDir, "notification-categories.json"), nil
}

// saveCategories saves the notification categories to a file.
func (ls *linuxNotifier) saveCategories() error {
	filePath, err := ls.getConfigFilePath()
	if err != nil {
		return err
	}

	ls.categoriesLock.RLock()
	data, err := json.Marshal(ls.categories)
	ls.categoriesLock.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal notification categories: %v", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write notification categories to file: %v", err)
	}

	return nil
}

// loadCategories loads notification categories from a file.
func (ls *linuxNotifier) loadCategories() error {
	filePath, err := ls.getConfigFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read notification categories file: %v", err)
	}

	if len(data) == 0 {
		return nil
	}

	categories := make(map[string]NotificationCategory)
	if err := json.Unmarshal(data, &categories); err != nil {
		return fmt.Errorf("failed to unmarshal notification categories: %v", err)
	}

	ls.categoriesLock.Lock()
	ls.categories = categories
	ls.categoriesLock.Unlock()

	return nil
}

func (r closedReason) string() string {
	switch r {
	case 1:
		return "expired"
	case 2:
		return "dismissed-by-user"
	case 3:
		return "closed-by-call"
	case 4:
		return "unknown"
	case 5:
		return "activated-by-user"
	default:
		return "other"
	}
}
