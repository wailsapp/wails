//go:build linux
// +build linux

package linux

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v2/internal/frontend"
)

var (
	conn                       *dbus.Conn
	categories                 map[string]frontend.NotificationCategory = make(map[string]frontend.NotificationCategory)
	categoriesLock             sync.RWMutex
	notifications              map[uint32]*notificationData = make(map[uint32]*notificationData)
	notificationsLock          sync.RWMutex
	notificationResultCallback func(result frontend.NotificationResult)
	callbackLock               sync.RWMutex
	appName                    string
	cancel                     context.CancelFunc
)

type notificationData struct {
	ID         string
	Title      string
	Subtitle   string
	Body       string
	CategoryID string
	Data       map[string]interface{}
	DBusID     uint32
	ActionMap  map[string]string
}

const (
	dbusNotificationInterface = "org.freedesktop.Notifications"
	dbusNotificationPath      = "/org/freedesktop/Notifications"
	DefaultActionIdentifier   = "DEFAULT_ACTION"
)

// Creates a new Notifications Service.
func (f *Frontend) InitializeNotifications() error {
	// Clean up any previous initialization
	f.CleanupNotifications()

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable: %w", err)
	}
	appName = filepath.Base(exe)

	_conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return fmt.Errorf("failed to connect to session bus: %w", err)
	}
	conn = _conn

	if err := f.loadCategories(); err != nil {
		fmt.Printf("Failed to load notification categories: %v\n", err)
	}

	var signalCtx context.Context
	signalCtx, cancel = context.WithCancel(context.Background())

	if err := f.setupSignalHandling(signalCtx); err != nil {
		return fmt.Errorf("failed to set up notification signal handling: %w", err)
	}

	return nil
}

// CleanupNotifications cleans up notification resources
func (f *Frontend) CleanupNotifications() {
	if cancel != nil {
		cancel()
		cancel = nil
	}

	if conn != nil {
		conn.Close()
		conn = nil
	}
}

func (f *Frontend) IsNotificationAvailable() bool {
	return true
}

// RequestNotificationAuthorization is a Linux stub that always returns true, nil.
// (authorization is macOS-specific)
func (f *Frontend) RequestNotificationAuthorization() (bool, error) {
	return true, nil
}

// CheckNotificationAuthorization is a Linux stub that always returns true.
// (authorization is macOS-specific)
func (f *Frontend) CheckNotificationAuthorization() (bool, error) {
	return true, nil
}

// SendNotification sends a basic notification with a unique identifier, title, subtitle, and body.
func (f *Frontend) SendNotification(options frontend.NotificationOptions) error {
	if conn == nil {
		return fmt.Errorf("notifications not initialized")
	}

	hints := map[string]dbus.Variant{}

	body := options.Body
	if options.Subtitle != "" {
		body = options.Subtitle + "\n" + body
	}

	defaultActionID := "default"
	actions := []string{defaultActionID, "Default"}

	actionMap := map[string]string{
		defaultActionID: DefaultActionIdentifier,
	}

	hints["x-notification-id"] = dbus.MakeVariant(options.ID)

	if options.Data != nil {
		userData, err := json.Marshal(options.Data)
		if err == nil {
			hints["x-user-data"] = dbus.MakeVariant(string(userData))
		}
	}

	// Call the Notify method on the D-Bus interface
	obj := conn.Object(dbusNotificationInterface, dbusNotificationPath)
	call := obj.Call(
		dbusNotificationInterface+".Notify",
		0,
		appName,
		uint32(0),
		"", // Icon
		options.Title,
		body,
		actions,
		hints,
		int32(-1),
	)

	if call.Err != nil {
		return fmt.Errorf("failed to send notification: %w", call.Err)
	}

	var dbusID uint32
	if err := call.Store(&dbusID); err != nil {
		return fmt.Errorf("failed to store notification ID: %w", err)
	}

	notification := &notificationData{
		ID:        options.ID,
		Title:     options.Title,
		Subtitle:  options.Subtitle,
		Body:      options.Body,
		Data:      options.Data,
		DBusID:    dbusID,
		ActionMap: actionMap,
	}

	notificationsLock.Lock()
	notifications[dbusID] = notification
	notificationsLock.Unlock()

	return nil
}

// SendNotificationWithActions sends a notification with additional actions.
func (f *Frontend) SendNotificationWithActions(options frontend.NotificationOptions) error {
	if conn == nil {
		return fmt.Errorf("notifications not initialized")
	}

	categoriesLock.RLock()
	category, exists := categories[options.CategoryID]
	categoriesLock.RUnlock()

	if options.CategoryID == "" || !exists {
		// Fall back to basic notification
		return f.SendNotification(options)
	}

	body := options.Body
	if options.Subtitle != "" {
		body = options.Subtitle + "\n" + body
	}

	var actions []string
	actionMap := make(map[string]string)

	defaultActionID := "default"
	actions = append(actions, defaultActionID, "Default")
	actionMap[defaultActionID] = DefaultActionIdentifier

	for _, action := range category.Actions {
		actions = append(actions, action.ID, action.Title)
		actionMap[action.ID] = action.ID
	}

	hints := map[string]dbus.Variant{}

	hints["x-notification-id"] = dbus.MakeVariant(options.ID)

	hints["x-category-id"] = dbus.MakeVariant(options.CategoryID)

	if options.Data != nil {
		userData, err := json.Marshal(options.Data)
		if err == nil {
			hints["x-user-data"] = dbus.MakeVariant(string(userData))
		}
	}

	obj := conn.Object(dbusNotificationInterface, dbusNotificationPath)
	call := obj.Call(
		dbusNotificationInterface+".Notify",
		0,
		appName,
		uint32(0),
		"", // Icon
		options.Title,
		body,
		actions,
		hints,
		int32(-1),
	)

	if call.Err != nil {
		return fmt.Errorf("failed to send notification: %w", call.Err)
	}

	var dbusID uint32
	if err := call.Store(&dbusID); err != nil {
		return fmt.Errorf("failed to store notification ID: %w", err)
	}

	notification := &notificationData{
		ID:         options.ID,
		Title:      options.Title,
		Subtitle:   options.Subtitle,
		Body:       options.Body,
		CategoryID: options.CategoryID,
		Data:       options.Data,
		DBusID:     dbusID,
		ActionMap:  actionMap,
	}

	notificationsLock.Lock()
	notifications[dbusID] = notification
	notificationsLock.Unlock()

	return nil
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
func (f *Frontend) RegisterNotificationCategory(category frontend.NotificationCategory) error {
	categoriesLock.Lock()
	categories[category.ID] = category
	categoriesLock.Unlock()

	if err := f.saveCategories(); err != nil {
		fmt.Printf("Failed to save notification categories: %v\n", err)
	}

	return nil
}

// RemoveNotificationCategory removes a previously registered NotificationCategory.
func (f *Frontend) RemoveNotificationCategory(categoryId string) error {
	categoriesLock.Lock()
	delete(categories, categoryId)
	categoriesLock.Unlock()

	if err := f.saveCategories(); err != nil {
		fmt.Printf("Failed to save notification categories: %v\n", err)
	}

	return nil
}

// RemoveAllPendingNotifications attempts to remove all active notifications.
func (f *Frontend) RemoveAllPendingNotifications() error {
	notificationsLock.Lock()
	dbusIDs := make([]uint32, 0, len(notifications))
	for id := range notifications {
		dbusIDs = append(dbusIDs, id)
	}
	notificationsLock.Unlock()

	for _, id := range dbusIDs {
		f.closeNotification(id)
	}

	return nil
}

// RemovePendingNotification removes a pending notification.
func (f *Frontend) RemovePendingNotification(identifier string) error {
	var dbusID uint32
	found := false

	notificationsLock.Lock()
	for id, notif := range notifications {
		if notif.ID == identifier {
			dbusID = id
			found = true
			break
		}
	}
	notificationsLock.Unlock()

	if !found {
		return nil
	}

	return f.closeNotification(dbusID)
}

// RemoveAllDeliveredNotifications functionally equivalent to RemoveAllPendingNotification on Linux.
func (f *Frontend) RemoveAllDeliveredNotifications() error {
	return f.RemoveAllPendingNotifications()
}

// RemoveDeliveredNotification functionally equivalent RemovePendingNotification on Linux.
func (f *Frontend) RemoveDeliveredNotification(identifier string) error {
	return f.RemovePendingNotification(identifier)
}

// RemoveNotification removes a notification by identifier.
func (f *Frontend) RemoveNotification(identifier string) error {
	return f.RemovePendingNotification(identifier)
}

func (f *Frontend) OnNotificationResponse(callback func(result frontend.NotificationResult)) {
	callbackLock.Lock()
	defer callbackLock.Unlock()

	notificationResultCallback = callback
}

// Helper method to close a notification.
func (f *Frontend) closeNotification(id uint32) error {
	if conn == nil {
		return fmt.Errorf("notifications not initialized")
	}

	obj := conn.Object(dbusNotificationInterface, dbusNotificationPath)
	call := obj.Call(dbusNotificationInterface+".CloseNotification", 0, id)

	if call.Err != nil {
		return fmt.Errorf("failed to close notification: %w", call.Err)
	}

	return nil
}

func (f *Frontend) getConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}

	appConfigDir := filepath.Join(configDir, appName)
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create app config directory: %w", err)
	}

	return appConfigDir, nil
}

// Save notification categories.
func (f *Frontend) saveCategories() error {
	configDir, err := f.getConfigDir()
	if err != nil {
		return err
	}

	categoriesFile := filepath.Join(configDir, "notification-categories.json")

	categoriesLock.RLock()
	categoriesData, err := json.MarshalIndent(categories, "", "  ")
	categoriesLock.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal notification categories: %w", err)
	}

	if err := os.WriteFile(categoriesFile, categoriesData, 0644); err != nil {
		return fmt.Errorf("failed to write notification categories to disk: %w", err)
	}

	return nil
}

// Load notification categories.
func (f *Frontend) loadCategories() error {
	configDir, err := f.getConfigDir()
	if err != nil {
		return err
	}

	categoriesFile := filepath.Join(configDir, "notification-categories.json")

	if _, err := os.Stat(categoriesFile); os.IsNotExist(err) {
		return nil
	}

	categoriesData, err := os.ReadFile(categoriesFile)
	if err != nil {
		return fmt.Errorf("failed to read notification categories from disk: %w", err)
	}

	_categories := make(map[string]frontend.NotificationCategory)
	if err := json.Unmarshal(categoriesData, &_categories); err != nil {
		return fmt.Errorf("failed to unmarshal notification categories: %w", err)
	}

	categoriesLock.Lock()
	categories = _categories
	categoriesLock.Unlock()

	return nil
}

// Setup signal handling for notification actions.
func (f *Frontend) setupSignalHandling(ctx context.Context) error {
	if err := conn.AddMatchSignal(
		dbus.WithMatchInterface(dbusNotificationInterface),
		dbus.WithMatchMember("ActionInvoked"),
	); err != nil {
		return err
	}

	if err := conn.AddMatchSignal(
		dbus.WithMatchInterface(dbusNotificationInterface),
		dbus.WithMatchMember("NotificationClosed"),
	); err != nil {
		return err
	}

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)

	go f.handleSignals(ctx, c)

	return nil
}

// Handle incoming D-Bus signals.
func (f *Frontend) handleSignals(ctx context.Context, c chan *dbus.Signal) {
	for {
		select {
		case <-ctx.Done():
			return
		case signal, ok := <-c:
			if !ok {
				return
			}

			switch signal.Name {
			case dbusNotificationInterface + ".ActionInvoked":
				f.handleActionInvoked(signal)
			case dbusNotificationInterface + ".NotificationClosed":
				f.handleNotificationClosed(signal)
			}
		}
	}
}

// Handle ActionInvoked signal.
func (f *Frontend) handleActionInvoked(signal *dbus.Signal) {
	if len(signal.Body) < 2 {
		return
	}

	dbusID, ok := signal.Body[0].(uint32)
	if !ok {
		return
	}

	actionID, ok := signal.Body[1].(string)
	if !ok {
		return
	}

	notificationsLock.Lock()
	notification, exists := notifications[dbusID]
	if exists {
		delete(notifications, dbusID)
	}
	notificationsLock.Unlock()

	if !exists {
		return
	}

	appActionID, ok := notification.ActionMap[actionID]
	if !ok {
		appActionID = actionID
	}

	response := frontend.NotificationResponse{
		ID:               notification.ID,
		ActionIdentifier: appActionID,
		Title:            notification.Title,
		Subtitle:         notification.Subtitle,
		Body:             notification.Body,
		CategoryID:       notification.CategoryID,
		UserInfo:         notification.Data,
	}

	result := frontend.NotificationResult{
		Response: response,
	}

	handleNotificationResult(result)
}

func handleNotificationResult(result frontend.NotificationResult) {
	callbackLock.Lock()
	callback := notificationResultCallback
	callbackLock.Unlock()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "panic in notification callback: %v\n", r)
			}
		}()
		callback(result)
	}()
}

// Handle NotificationClosed signal.
// Reason codes:
// 1 - expired timeout
// 2 - dismissed by user (click on X)
// 3 - closed by CloseNotification call
// 4 - undefined/reserved
func (f *Frontend) handleNotificationClosed(signal *dbus.Signal) {
	if len(signal.Body) < 2 {
		return
	}

	dbusID, ok := signal.Body[0].(uint32)
	if !ok {
		return
	}

	reason, ok := signal.Body[1].(uint32)
	if !ok {
		reason = 0 // Unknown reason
	}

	notificationsLock.Lock()
	notification, exists := notifications[dbusID]
	if exists {
		delete(notifications, dbusID)
	}
	notificationsLock.Unlock()

	if !exists {
		return
	}

	if reason == 2 {
		response := frontend.NotificationResponse{
			ID:               notification.ID,
			ActionIdentifier: DefaultActionIdentifier,
			Title:            notification.Title,
			Subtitle:         notification.Subtitle,
			Body:             notification.Body,
			CategoryID:       notification.CategoryID,
			UserInfo:         notification.Data,
		}

		result := frontend.NotificationResult{
			Response: response,
		}

		handleNotificationResult(result)
	}
}
