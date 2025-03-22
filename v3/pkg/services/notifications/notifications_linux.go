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
	// Categories
	categories     map[string]NotificationCategory
	categoriesLock sync.RWMutex

	// App info
	appName string

	// Notification system
	sync.Mutex
	method       string
	dbusConn     *dbus.Conn
	sendPath     string
	activeNotifs map[string]uint32               // Maps our notification IDs to system IDs
	contexts     map[string]*notificationContext // Stores notification contexts by our ID

	// Listener management
	listenerCtx     context.Context
	listenerCancel  context.CancelFunc
	listenerRunning bool

	// Initialization
	initOnce    sync.Once
	initialized bool
}

type notificationContext struct {
	ID       string
	SystemID uint32
	Actions  map[string]string      // Maps action keys to display labels
	UserData map[string]interface{} // The original user data
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

	notifyChannelBufferSize = 25
)

type closedReason uint32

// New creates a new Notifications Service
func New() *Service {
	notificationServiceOnce.Do(func() {
		impl := &linuxNotifier{
			categories:   make(map[string]NotificationCategory),
			activeNotifs: make(map[string]uint32),
			contexts:     make(map[string]*notificationContext),
		}

		NotificationService = &Service{
			impl: impl,
		}
	})

	return NotificationService
}

// Startup is called when the service is loaded
func (ln *linuxNotifier) Startup(ctx context.Context, options application.ServiceOptions) error {
	ln.appName = application.Get().Config().Name

	if err := ln.loadCategories(); err != nil {
		fmt.Printf("Failed to load notification categories: %v\n", err)
	}

	var err error
	ln.initOnce.Do(func() {
		err = ln.initNotificationSystem()
		ln.initialized = err == nil
	})

	return err
}

// initNotificationSystem initializes the notification system
func (ln *linuxNotifier) initNotificationSystem() error {
	ln.Lock()
	defer ln.Unlock()

	// Cancel any existing listener
	if ln.listenerCancel != nil {
		ln.listenerCancel()
		ln.listenerCancel = nil
	}

	// Create a new context for the listener
	ln.listenerCtx, ln.listenerCancel = context.WithCancel(context.Background())

	// Reset state
	ln.activeNotifs = make(map[string]uint32)
	ln.contexts = make(map[string]*notificationContext)
	ln.listenerRunning = false

	// Try dbus first
	dbusConn, err := ln.initDBus()
	if err == nil {
		ln.dbusConn = dbusConn
		ln.method = MethodDbus

		// Start the dbus signal listener
		go ln.startDBusListener(ln.listenerCtx)
		ln.listenerRunning = true
		return nil
	}

	// Try notify-send as fallback
	sendPath, err := ln.initNotifySend()
	if err == nil {
		ln.sendPath = sendPath
		ln.method = MethodNotifySend
		return nil
	}

	// No method available
	ln.method = ""
	ln.sendPath = ""
	return errors.New("no notification method is available")
}

// initDBus attempts to initialize D-Bus notifications
func (ln *linuxNotifier) initDBus() (*dbus.Conn, error) {
	conn, err := dbus.SessionBusPrivate()
	if err != nil {
		return nil, err
	}

	if err = conn.Auth(nil); err != nil {
		conn.Close()
		return nil, err
	}

	if err = conn.Hello(); err != nil {
		conn.Close()
		return nil, err
	}

	obj := conn.Object(dbusNotificationsInterface, dbusObjectPath)
	call := obj.Call(callGetCapabilities, 0)
	if call.Err != nil {
		conn.Close()
		return nil, call.Err
	}

	var ret []string
	err = call.Store(&ret)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Add a listener for notification signals
	err = conn.AddMatchSignal(
		dbus.WithMatchObjectPath(dbusObjectPath),
		dbus.WithMatchInterface(dbusNotificationsInterface),
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

// initNotifySend attempts to find notify-send binary
func (ln *linuxNotifier) initNotifySend() (string, error) {
	// Try standard notify-send
	send, err := exec.LookPath("notify-send")
	if err == nil {
		return send, nil
	}

	// Try sw-notify-send (in some distros)
	send, err = exec.LookPath("sw-notify-send")
	if err == nil {
		return send, nil
	}

	return "", errors.New("notify-send not found")
}

// Shutdown is called when the service is unloaded
func (ln *linuxNotifier) Shutdown() error {
	ln.Lock()

	// Cancel the listener context if it's running
	if ln.listenerCancel != nil {
		ln.listenerCancel()
		ln.listenerCancel = nil
	}

	// Close the connection
	if ln.dbusConn != nil {
		ln.dbusConn.Close()
		ln.dbusConn = nil
	}

	// Clear state
	ln.activeNotifs = make(map[string]uint32)
	ln.contexts = make(map[string]*notificationContext)
	ln.method = ""
	ln.sendPath = ""
	ln.initialized = false

	ln.Unlock()

	return ln.saveCategories()
}

// startDBusListener listens for DBus signals for notification actions and closures
func (ln *linuxNotifier) startDBusListener(ctx context.Context) {
	signal := make(chan *dbus.Signal, notifyChannelBufferSize)
	ln.dbusConn.Signal(signal)

	defer func() {
		ln.Lock()
		ln.listenerRunning = false
		ln.Unlock()
		ln.dbusConn.RemoveSignal(signal) // Remove signal handler
		close(signal)                    // Clean up channel
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
				ln.handleNotificationClosed(systemID, reason)
			case signalActionInvoked:
				systemID := s.Body[0].(uint32)
				actionKey := s.Body[1].(string)
				ln.handleActionInvoked(systemID, actionKey)
			}
		}
	}
}

// handleNotificationClosed processes notification closed signals
func (ln *linuxNotifier) handleNotificationClosed(systemID uint32, reason string) {
	// Find our notification ID for this system ID
	var notifID string
	var userData map[string]interface{}

	ln.Lock()
	for id, sysID := range ln.activeNotifs {
		if sysID == systemID {
			notifID = id
			// Get the user data from context if available
			if ctx, exists := ln.contexts[id]; exists {
				userData = ctx.UserData
			}
			break
		}
	}
	ln.Unlock()

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
		ln.Lock()
		delete(ln.contexts, notifID)
		delete(ln.activeNotifs, notifID)
		ln.Unlock()
	}
}

// handleActionInvoked processes action invoked signals
func (ln *linuxNotifier) handleActionInvoked(systemID uint32, actionKey string) {
	// Find our notification ID and context for this system ID
	var notifID string
	var ctx *notificationContext

	ln.Lock()
	for id, sysID := range ln.activeNotifs {
		if sysID == systemID {
			notifID = id
			ctx = ln.contexts[id]
			break
		}
	}
	ln.Unlock()

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
		ln.Lock()
		delete(ln.contexts, notifID)
		delete(ln.activeNotifs, notifID)
		ln.Unlock()
	}
}

// CheckBundleIdentifier is a Linux stub that always returns true.
func (ln *linuxNotifier) CheckBundleIdentifier() bool {
	return true
}

// RequestNotificationAuthorization is a Linux stub that always returns true.
func (ln *linuxNotifier) RequestNotificationAuthorization() (bool, error) {
	return true, nil
}

// CheckNotificationAuthorization is a Linux stub that always returns true.
func (ln *linuxNotifier) CheckNotificationAuthorization() (bool, error) {
	return true, nil
}

// SendNotification sends a basic notification with a unique identifier, title, subtitle, and body.
func (ln *linuxNotifier) SendNotification(options NotificationOptions) error {
	if !ln.initialized {
		return errors.New("notification service not initialized")
	}

	if err := validateNotificationOptions(options); err != nil {
		return err
	}

	ln.Lock()
	defer ln.Unlock()

	var (
		systemID uint32
		err      error
	)

	switch ln.method {
	case MethodDbus:
		systemID, err = ln.sendViaDbus(options, nil)
	case MethodNotifySend:
		systemID, err = ln.sendViaNotifySend(options)
	default:
		err = errors.New("no notification method is available")
	}

	if err == nil && systemID > 0 {
		// Store the system ID mapping
		ln.activeNotifs[options.ID] = systemID

		// Create and store the notification context
		ctx := &notificationContext{
			ID:       options.ID,
			SystemID: systemID,
			UserData: options.Data,
		}
		ln.contexts[options.ID] = ctx
	}

	return err
}

// SendNotificationWithActions sends a notification with additional actions.
func (ln *linuxNotifier) SendNotificationWithActions(options NotificationOptions) error {
	if !ln.initialized {
		return errors.New("notification service not initialized")
	}

	if err := validateNotificationOptions(options); err != nil {
		return err
	}

	ln.categoriesLock.RLock()
	category, exists := ln.categories[options.CategoryID]
	ln.categoriesLock.RUnlock()

	if !exists {
		return ln.SendNotification(options)
	}

	ln.Lock()
	defer ln.Unlock()

	var (
		systemID uint32
		err      error
	)

	switch ln.method {
	case MethodDbus:
		systemID, err = ln.sendViaDbus(options, &category)
	case MethodNotifySend:
		// notify-send doesn't support actions, fall back to basic notification
		systemID, err = ln.sendViaNotifySend(options)
	default:
		err = errors.New("no notification method is available")
	}

	if err == nil && systemID > 0 {
		// Store the system ID mapping
		ln.activeNotifs[options.ID] = systemID

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

		ln.contexts[options.ID] = ctx
	}

	return err
}

// sendViaDbus sends a notification via dbus
func (ln *linuxNotifier) sendViaDbus(options NotificationOptions, category *NotificationCategory) (result uint32, err error) {
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
	obj := ln.dbusConn.Object(dbusNotificationsInterface, dbusObjectPath)
	dbusArgs := []interface{}{
		ln.appName,    // App name
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
func (ln *linuxNotifier) sendViaNotifySend(options NotificationOptions) (uint32, error) {
	args := []string{
		options.Title,
		options.Body,
		"--urgency=normal",
	}

	// Execute the command
	cmd := exec.Command(ln.sendPath, args...)
	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("notify-send error: %v", err)
	}

	// notify-send doesn't return IDs, so we use 0
	return 0, nil
}

// RegisterNotificationCategory registers a new NotificationCategory
func (ln *linuxNotifier) RegisterNotificationCategory(category NotificationCategory) error {
	ln.categoriesLock.Lock()
	ln.categories[category.ID] = category
	ln.categoriesLock.Unlock()

	return ln.saveCategories()
}

// RemoveNotificationCategory removes a previously registered NotificationCategory
func (ln *linuxNotifier) RemoveNotificationCategory(categoryId string) error {
	ln.categoriesLock.Lock()
	delete(ln.categories, categoryId)
	ln.categoriesLock.Unlock()

	return ln.saveCategories()
}

// RemoveAllPendingNotifications is a Linux stub that always returns nil
func (ln *linuxNotifier) RemoveAllPendingNotifications() error {
	return nil
}

// RemovePendingNotification is a Linux stub that always returns nil
func (ln *linuxNotifier) RemovePendingNotification(_ string) error {
	return nil
}

// RemoveAllDeliveredNotifications is a Linux stub that always returns nil
func (ln *linuxNotifier) RemoveAllDeliveredNotifications() error {
	return nil
}

// RemoveDeliveredNotification is a Linux stub that always returns nil
func (ln *linuxNotifier) RemoveDeliveredNotification(_ string) error {
	return nil
}

// RemoveNotification removes a notification by ID (Linux-specific)
func (ln *linuxNotifier) RemoveNotification(identifier string) error {
	if !ln.initialized || ln.method != MethodDbus || ln.dbusConn == nil {
		return errors.New("dbus not available for closing notifications")
	}

	// Get the system ID for this notification
	ln.Lock()
	systemID, exists := ln.activeNotifs[identifier]
	ln.Unlock()

	if !exists {
		return nil // Already closed or unknown
	}

	// Call CloseNotification on dbus
	obj := ln.dbusConn.Object(dbusNotificationsInterface, dbusObjectPath)
	call := obj.Call(callCloseNotification, 0, systemID)

	return call.Err
}

// getConfigFilePath returns the path to the configuration file
func (ln *linuxNotifier) getConfigFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %v", err)
	}

	appConfigDir := filepath.Join(configDir, ln.appName)
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %v", err)
	}

	return filepath.Join(appConfigDir, "notification-categories.json"), nil
}

// saveCategories saves the notification categories to a file
func (ln *linuxNotifier) saveCategories() error {
	filePath, err := ln.getConfigFilePath()
	if err != nil {
		return err
	}

	ln.categoriesLock.RLock()
	data, err := json.Marshal(ln.categories)
	ln.categoriesLock.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal notification categories: %v", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write notification categories to file: %v", err)
	}

	return nil
}

// loadCategories loads notification categories from a file
func (ln *linuxNotifier) loadCategories() error {
	filePath, err := ln.getConfigFilePath()
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

	ln.categoriesLock.Lock()
	ln.categories = categories
	ln.categoriesLock.Unlock()

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
