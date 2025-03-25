//go:build linux

package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type linuxNotifier struct {
	conn              *dbus.Conn
	categories        map[string]NotificationCategory
	categoriesLock    sync.RWMutex
	notifications     map[uint32]*notificationData
	notificationsLock sync.RWMutex
	appName           string
	cancel            context.CancelFunc
}

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
)

// Creates a new Notifications Service.
func New() *Service {
	notificationServiceOnce.Do(func() {
		impl := &linuxNotifier{
			categories:    make(map[string]NotificationCategory),
			notifications: make(map[uint32]*notificationData),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		conn, err := dbus.ConnectSessionBus()
		if err != nil {
			fmt.Printf("Warning: Failed to connect to D-Bus session bus: %v\n", err)
			fmt.Printf("Notifications will be unavailable\n")
		} else {
			impl.conn = conn

			obj := conn.Object(dbusNotificationInterface, dbusNotificationPath)
			call := obj.CallWithContext(ctx, dbusNotificationInterface+".GetCapabilities", 0)

			var capabilities []string
			err := call.Store(&capabilities)

			if err != nil {
				fmt.Printf("Warning: D-Bus notification service not ready: %v\n", err)
			} else {
				fmt.Printf("D-Bus notification service is ready with capabilities: %v\n", capabilities)
			}
		}

		NotificationService = &Service{
			impl: impl,
		}
	})

	return NotificationService
}

// Helper method to check if D-Bus connection is available
func (ln *linuxNotifier) checkConnection() error {
	if ln.conn == nil {
		return fmt.Errorf("D-Bus connection is not initialized, notifications are unavailable")
	}

	return nil
}

// Startup is called when the service is loaded.
func (ln *linuxNotifier) Startup(ctx context.Context, options application.ServiceOptions) error {
	ln.appName = application.Get().Config().Name

	if ln.conn == nil {
		conn, err := dbus.ConnectSessionBus()
		if err != nil {
			fmt.Printf("Warning: Failed to connect to D-Bus session bus: %v\n", err)
			fmt.Printf("Notifications will be unavailable\n")

			return nil
		}
		ln.conn = conn
	}

	if err := ln.loadCategories(); err != nil {
		fmt.Printf("Failed to load notification categories: %v\n", err)
	}

	var signalCtx context.Context
	signalCtx, ln.cancel = context.WithCancel(context.Background())

	if err := ln.setupSignalHandling(signalCtx); err != nil {
		fmt.Printf("Warning: Failed to set up notification signal handling: %v\n", err)
	}

	return nil
}

// Shutdown will save categories and close the D-Bus connection when the service unloads.
func (ln *linuxNotifier) Shutdown() error {
	if ln.cancel != nil {
		ln.cancel()
	}

	if err := ln.saveCategories(); err != nil {
		fmt.Printf("Failed to save notification categories: %v\n", err)
	}

	if ln.conn != nil {
		return ln.conn.Close()
	}
	return nil
}

// RequestNotificationAuthorization is a Linux stub that always returns true, nil.
// (authorization is macOS-specific)
func (ln *linuxNotifier) RequestNotificationAuthorization() (bool, error) {
	return true, nil
}

// CheckNotificationAuthorization is a Linux stub that always returns true.
// (authorization is macOS-specific)
func (ln *linuxNotifier) CheckNotificationAuthorization() (bool, error) {
	return true, nil
}

// SendNotification sends a basic notification with a unique identifier, title, subtitle, and body.
func (ln *linuxNotifier) SendNotification(options NotificationOptions) error {
	if err := ln.checkConnection(); err != nil {
		return err
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
	obj := ln.conn.Object(dbusNotificationInterface, dbusNotificationPath)
	call := obj.Call(
		dbusNotificationInterface+".Notify",
		0,
		ln.appName,
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

	ln.notificationsLock.Lock()
	ln.notifications[dbusID] = notification
	ln.notificationsLock.Unlock()

	return nil
}

// SendNotificationWithActions sends a notification with additional actions.
func (ln *linuxNotifier) SendNotificationWithActions(options NotificationOptions) error {
	if err := ln.checkConnection(); err != nil {
		return err
	}

	ln.categoriesLock.RLock()
	category, exists := ln.categories[options.CategoryID]
	ln.categoriesLock.RUnlock()

	if options.CategoryID == "" || !exists {
		// Fall back to basic notification
		return ln.SendNotification(options)
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

	obj := ln.conn.Object(dbusNotificationInterface, dbusNotificationPath)
	call := obj.Call(
		dbusNotificationInterface+".Notify",
		0,
		ln.appName,
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

	ln.notificationsLock.Lock()
	ln.notifications[dbusID] = notification
	ln.notificationsLock.Unlock()

	return nil
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
func (ln *linuxNotifier) RegisterNotificationCategory(category NotificationCategory) error {
	ln.categoriesLock.Lock()
	ln.categories[category.ID] = category
	ln.categoriesLock.Unlock()

	if err := ln.saveCategories(); err != nil {
		fmt.Printf("Failed to save notification categories: %v\n", err)
	}

	return nil
}

// RemoveNotificationCategory removes a previously registered NotificationCategory.
func (ln *linuxNotifier) RemoveNotificationCategory(categoryId string) error {
	ln.categoriesLock.Lock()
	delete(ln.categories, categoryId)
	ln.categoriesLock.Unlock()

	if err := ln.saveCategories(); err != nil {
		fmt.Printf("Failed to save notification categories: %v\n", err)
	}

	return nil
}

// RemoveAllPendingNotifications attempts to remove all active notifications.
func (ln *linuxNotifier) RemoveAllPendingNotifications() error {
	if err := ln.checkConnection(); err != nil {
		return err
	}

	ln.notificationsLock.Lock()
	dbusIDs := make([]uint32, 0, len(ln.notifications))
	for id := range ln.notifications {
		dbusIDs = append(dbusIDs, id)
	}
	ln.notificationsLock.Unlock()

	for _, id := range dbusIDs {
		ln.closeNotification(id)
	}

	return nil
}

// RemovePendingNotification removes a pending notification.
func (ln *linuxNotifier) RemovePendingNotification(identifier string) error {
	if err := ln.checkConnection(); err != nil {
		return err
	}

	var dbusID uint32
	found := false

	ln.notificationsLock.Lock()
	for id, notif := range ln.notifications {
		if notif.ID == identifier {
			dbusID = id
			found = true
			break
		}
	}
	ln.notificationsLock.Unlock()

	if !found {
		return nil
	}

	return ln.closeNotification(dbusID)
}

// RemoveAllDeliveredNotifications functionally equivalent to RemoveAllPendingNotification on Linux.
func (ln *linuxNotifier) RemoveAllDeliveredNotifications() error {
	return ln.RemoveAllPendingNotifications()
}

// RemoveDeliveredNotification functionally equivalent RemovePendingNotification on Linux.
func (ln *linuxNotifier) RemoveDeliveredNotification(identifier string) error {
	return ln.RemovePendingNotification(identifier)
}

// RemoveNotification removes a notification by identifier.
func (ln *linuxNotifier) RemoveNotification(identifier string) error {
	return ln.RemovePendingNotification(identifier)
}

// Helper method to close a notification.
func (ln *linuxNotifier) closeNotification(id uint32) error {
	if err := ln.checkConnection(); err != nil {
		return err
	}

	obj := ln.conn.Object(dbusNotificationInterface, dbusNotificationPath)
	call := obj.Call(dbusNotificationInterface+".CloseNotification", 0, id)

	if call.Err != nil {
		return fmt.Errorf("failed to close notification: %w", call.Err)
	}

	return nil
}

func (ln *linuxNotifier) getConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}

	appConfigDir := filepath.Join(configDir, ln.appName)
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create app config directory: %w", err)
	}

	return appConfigDir, nil
}

// Save notification categories.
func (ln *linuxNotifier) saveCategories() error {
	configDir, err := ln.getConfigDir()
	if err != nil {
		return err
	}

	categoriesFile := filepath.Join(configDir, "notification-categories.json")

	ln.categoriesLock.RLock()
	categoriesData, err := json.MarshalIndent(ln.categories, "", "  ")
	ln.categoriesLock.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal notification categories: %w", err)
	}

	if err := os.WriteFile(categoriesFile, categoriesData, 0644); err != nil {
		return fmt.Errorf("failed to write notification categories to disk: %w", err)
	}

	return nil
}

// Load notification categories.
func (ln *linuxNotifier) loadCategories() error {
	configDir, err := ln.getConfigDir()
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

	categories := make(map[string]NotificationCategory)
	if err := json.Unmarshal(categoriesData, &categories); err != nil {
		return fmt.Errorf("failed to unmarshal notification categories: %w", err)
	}

	ln.categoriesLock.Lock()
	ln.categories = categories
	ln.categoriesLock.Unlock()

	return nil
}

// Setup signal handling for notification actions.
func (ln *linuxNotifier) setupSignalHandling(ctx context.Context) error {
	if err := ln.checkConnection(); err != nil {
		return err
	}

	if err := ln.conn.AddMatchSignal(
		dbus.WithMatchInterface(dbusNotificationInterface),
		dbus.WithMatchMember("ActionInvoked"),
	); err != nil {
		return err
	}

	if err := ln.conn.AddMatchSignal(
		dbus.WithMatchInterface(dbusNotificationInterface),
		dbus.WithMatchMember("NotificationClosed"),
	); err != nil {
		return err
	}

	c := make(chan *dbus.Signal, 10)
	ln.conn.Signal(c)

	go ln.handleSignals(ctx, c)

	return nil
}

// Handle incoming D-Bus signals.
func (ln *linuxNotifier) handleSignals(ctx context.Context, c chan *dbus.Signal) {
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
				ln.handleActionInvoked(signal)
			case dbusNotificationInterface + ".NotificationClosed":
				ln.handleNotificationClosed(signal)
			}
		}
	}
}

// Handle ActionInvoked signal.
func (ln *linuxNotifier) handleActionInvoked(signal *dbus.Signal) {
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

	ln.notificationsLock.Lock()
	notification, exists := ln.notifications[dbusID]
	if exists {
		delete(ln.notifications, dbusID)
	}
	ln.notificationsLock.Unlock()

	if !exists {
		return
	}

	appActionID, ok := notification.ActionMap[actionID]
	if !ok {
		appActionID = actionID
	}

	response := NotificationResponse{
		ID:               notification.ID,
		ActionIdentifier: appActionID,
		Title:            notification.Title,
		Subtitle:         notification.Subtitle,
		Body:             notification.Body,
		CategoryID:       notification.CategoryID,
		UserInfo:         notification.Data,
	}

	result := NotificationResult{
		Response: response,
	}

	if ns := getNotificationService(); ns != nil {
		ns.handleNotificationResult(result)
	}
}

// Handle NotificationClosed signal.
// Reason codes:
// 1 - expired timeout
// 2 - dismissed by user (click on X)
// 3 - closed by CloseNotification call
// 4 - undefined/reserved
func (ln *linuxNotifier) handleNotificationClosed(signal *dbus.Signal) {
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

	ln.notificationsLock.Lock()
	notification, exists := ln.notifications[dbusID]
	if exists {
		delete(ln.notifications, dbusID)
	}
	ln.notificationsLock.Unlock()

	if !exists {
		return
	}

	if reason == 2 {
		response := NotificationResponse{
			ID:               notification.ID,
			ActionIdentifier: DefaultActionIdentifier,
			Title:            notification.Title,
			Subtitle:         notification.Subtitle,
			Body:             notification.Body,
			CategoryID:       notification.CategoryID,
			UserInfo:         notification.Data,
		}

		result := NotificationResult{
			Response: response,
		}

		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(result)
		}
	}
}
