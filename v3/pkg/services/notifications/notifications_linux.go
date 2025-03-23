//go:build linux

package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type linuxNotifier struct {
	conn             *dbus.Conn
	categories       map[string]NotificationCategory
	categoriesLock   sync.RWMutex
	activeNotifs     map[uint32]string
	notificationMeta map[uint32]map[string]interface{}
	activeNotifsLock sync.RWMutex
	appName          string
	cancel           context.CancelFunc
}

const (
	dbusNotificationInterface = "org.freedesktop.Notifications"
	dbusNotificationPath      = "/org/freedesktop/Notifications"
)

// Creates a new Notifications Service.
func New() *Service {
	notificationServiceOnce.Do(func() {
		impl := &linuxNotifier{
			categories:       make(map[string]NotificationCategory),
			activeNotifs:     make(map[uint32]string),
			notificationMeta: make(map[uint32]map[string]interface{}),
		}

		NotificationService = &Service{
			impl: impl,
		}
	})

	return NotificationService
}

// Startup is called when the service is loaded
func (ln *linuxNotifier) Startup(ctx context.Context, options application.ServiceOptions) error {
	app := application.Get()
	if app != nil && app.Config().Name != "" {
		ln.appName = app.Config().Name
	}

	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return fmt.Errorf("failed to connect to session bus: %w", err)
	}
	ln.conn = conn

	if err := ln.loadCategories(); err != nil {
		fmt.Printf("Failed to load notification categories: %v\n", err)
	}

	var signalCtx context.Context
	signalCtx, ln.cancel = context.WithCancel(context.Background())

	// Set up signal handling for notification actions
	if err := ln.setupSignalHandling(signalCtx); err != nil {
		return fmt.Errorf("failed to set up notification signal handling: %w", err)
	}

	return nil
}

// Shutdown will save categories and close the D-Bus connection when the service unloads
func (ln *linuxNotifier) Shutdown() error {
	ln.cancel()

	if err := ln.saveCategories(); err != nil {
		fmt.Printf("Failed to save notification categories: %v\n", err)
	}

	return ln.conn.Close()
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
	hints := map[string]dbus.Variant{}

	// Use subtitle as part of the body if provided
	body := options.Body
	if options.Subtitle != "" {
		body = options.Subtitle + "\n" + body
	}

	if options.Data != nil {
		jsonData, err := json.Marshal(options.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal notification data: %w", err)
		}
		hints["x-user-info"] = dbus.MakeVariant(string(jsonData))
	}

	metadataJSON, err := json.Marshal(map[string]interface{}{
		"id":       options.ID,
		"title":    options.Title,
		"subtitle": options.Subtitle,
		"body":     options.Body,
		"data":     options.Data,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification metadata: %w", err)
	}
	hints["x-wails-metadata"] = dbus.MakeVariant(string(metadataJSON))

	actions := []string{}

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
		uint32(0),
	)

	if call.Err != nil {
		return fmt.Errorf("failed to send notification: %w", call.Err)
	}

	var notifID uint32
	if err := call.Store(&notifID); err != nil {
		return fmt.Errorf("failed to store notification ID: %w", err)
	}

	ln.activeNotifsLock.Lock()
	ln.activeNotifs[notifID] = options.ID

	metadata := map[string]interface{}{
		"id":       options.ID,
		"title":    options.Title,
		"subtitle": options.Subtitle,
		"body":     options.Body,
		"data":     options.Data,
	}

	ln.notificationMeta[notifID] = metadata
	ln.activeNotifsLock.Unlock()

	return nil
}

// SendNotificationWithActions sends a notification with additional actions.
// (Inputs are only supported on macOS and Windows)
func (ln *linuxNotifier) SendNotificationWithActions(options NotificationOptions) error {
	ln.categoriesLock.RLock()
	category, exists := ln.categories[options.CategoryID]
	ln.categoriesLock.RUnlock()

	if options.CategoryID == "" || !exists {
		// Fall back to basic notification
		return ln.SendNotification(options)
	}

	// Use subtitle as part of the body if provided
	body := options.Body
	if options.Subtitle != "" {
		body = options.Subtitle + "\n" + body
	}

	var actions []string
	for _, action := range category.Actions {
		actions = append(actions, action.ID, action.Title)
	}

	hints := map[string]dbus.Variant{}

	if options.Data != nil {
		jsonData, err := json.Marshal(options.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal notification data: %w", err)
		}
		hints["x-user-info"] = dbus.MakeVariant(string(jsonData))
	}

	hints["category"] = dbus.MakeVariant(options.CategoryID)

	metadataJSON, err := json.Marshal(map[string]interface{}{
		"id":         options.ID,
		"title":      options.Title,
		"subtitle":   options.Subtitle,
		"body":       options.Body,
		"categoryId": options.CategoryID,
		"data":       options.Data,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification metadata: %w", err)
	}
	hints["x-wails-metadata"] = dbus.MakeVariant(string(metadataJSON))

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
		uint32(0),
	)

	if call.Err != nil {
		return fmt.Errorf("failed to send notification: %w", call.Err)
	}

	var notifID uint32
	if err := call.Store(&notifID); err != nil {
		return fmt.Errorf("failed to store notification ID: %w", err)
	}

	ln.activeNotifsLock.Lock()

	ln.activeNotifs[notifID] = options.ID

	metadata := map[string]interface{}{
		"id":         options.ID,
		"title":      options.Title,
		"subtitle":   options.Subtitle,
		"body":       options.Body,
		"categoryId": options.CategoryID,
		"data":       options.Data,
	}

	ln.notificationMeta[notifID] = metadata
	ln.activeNotifsLock.Unlock()

	return nil
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
// Registering a category with the same name as a previously registered NotificationCategory will override it.
func (ln *linuxNotifier) RegisterNotificationCategory(category NotificationCategory) error {
	ln.categoriesLock.Lock()
	defer ln.categoriesLock.Unlock()

	ln.categories[category.ID] = category

	if err := ln.saveCategories(); err != nil {
		fmt.Printf("Failed to save notification categories: %v\n", err)
	}

	return nil
}

// RemoveNotificationCategory removes a previously registered NotificationCategory.
func (ln *linuxNotifier) RemoveNotificationCategory(categoryId string) error {
	ln.categoriesLock.Lock()
	defer ln.categoriesLock.Unlock()

	delete(ln.categories, categoryId)

	if err := ln.saveCategories(); err != nil {
		fmt.Printf("Failed to save notification categories: %v\n", err)
	}

	return nil
}

// RemoveAllPendingNotifications is not directly supported in Linux D-Bus
// but we can try to remove all active notifications.
func (ln *linuxNotifier) RemoveAllPendingNotifications() error {
	ln.activeNotifsLock.RLock()
	notifIDs := make([]uint32, 0, len(ln.activeNotifs))
	for id := range ln.activeNotifs {
		notifIDs = append(notifIDs, id)
	}
	ln.activeNotifsLock.RUnlock()

	for _, id := range notifIDs {
		ln.closeNotification(id)
	}

	return nil
}

// RemovePendingNotification removes a pending notification.
func (ln *linuxNotifier) RemovePendingNotification(identifier string) error {
	var notifID uint32
	found := false

	ln.activeNotifsLock.RLock()
	for id, ident := range ln.activeNotifs {
		if ident == identifier {
			notifID = id
			found = true
			break
		}
	}
	ln.activeNotifsLock.RUnlock()

	if !found {
		return nil
	}

	return ln.closeNotification(notifID)
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
	obj := ln.conn.Object(dbusNotificationInterface, dbusNotificationPath)
	call := obj.Call(dbusNotificationInterface+".CloseNotification", 0, id)

	if call.Err != nil {
		return fmt.Errorf("failed to close notification: %w", call.Err)
	}

	ln.activeNotifsLock.Lock()
	delete(ln.activeNotifs, id)
	delete(ln.notificationMeta, id)
	ln.activeNotifsLock.Unlock()

	return nil
}

// Get the config directory for the app.
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
	if len(signal.Body) < 1 {
		return
	}

	notifID, ok := signal.Body[0].(uint32)
	if !ok {
		return
	}

	actionID, ok := signal.Body[1].(string)
	if !ok {
		return
	}

	ln.activeNotifsLock.RLock()
	identifier, idExists := ln.activeNotifs[notifID]
	metadata, metaExists := ln.notificationMeta[notifID]
	ln.activeNotifsLock.RUnlock()

	if !idExists || !metaExists {
		return
	}

	response := NotificationResponse{
		ID:               identifier,
		ActionIdentifier: actionID,
	}

	if title, ok := metadata["title"].(string); ok {
		response.Title = title
	}

	if subtitle, ok := metadata["subtitle"].(string); ok {
		response.Subtitle = subtitle
	}

	if body, ok := metadata["body"].(string); ok {
		response.Body = body
	}

	if categoryID, ok := metadata["categoryId"].(string); ok {
		response.CategoryID = categoryID
	}

	if userData, ok := metadata["data"].(map[string]interface{}); ok {
		response.UserInfo = userData
	}

	if actionID == DefaultActionIdentifier {
		response.ActionIdentifier = DefaultActionIdentifier
	}

	result := NotificationResult{
		Response: response,
	}

	if ns := getNotificationService(); ns != nil {
		ns.handleNotificationResult(result)
	}

	ln.activeNotifsLock.Lock()
	delete(ln.activeNotifs, notifID)
	delete(ln.notificationMeta, notifID)
	ln.activeNotifsLock.Unlock()
}

// Handle NotificationClosed signal.
// The second parameter contains the close reason:
// 1 - expired timeout
// 2 - dismissed by user (click on X)
// 3 - closed by CloseNotification call
// 4 - undefined/reserved
func (ln *linuxNotifier) handleNotificationClosed(signal *dbus.Signal) {
	if len(signal.Body) < 1 {
		return
	}

	notifID, ok := signal.Body[0].(uint32)
	if !ok {
		return
	}

	var reason uint32 = 0
	if len(signal.Body) > 1 {
		if r, ok := signal.Body[1].(uint32); ok {
			reason = r
		}
	}

	if reason != 1 && reason != 3 {
		ln.activeNotifsLock.RLock()
		identifier, idExists := ln.activeNotifs[notifID]
		metadata, metaExists := ln.notificationMeta[notifID]
		ln.activeNotifsLock.RUnlock()

		if idExists && metaExists {
			response := NotificationResponse{
				ID:               identifier,
				ActionIdentifier: DefaultActionIdentifier,
			}

			if title, ok := metadata["title"].(string); ok {
				response.Title = title
			}

			if subtitle, ok := metadata["subtitle"].(string); ok {
				response.Subtitle = subtitle
			}

			if body, ok := metadata["body"].(string); ok {
				response.Body = body
			}

			if categoryID, ok := metadata["categoryId"].(string); ok {
				response.CategoryID = categoryID
			}

			if userData, ok := metadata["data"].(map[string]interface{}); ok {
				response.UserInfo = userData
			}

			result := NotificationResult{
				Response: response,
			}

			if ns := getNotificationService(); ns != nil {
				ns.handleNotificationResult(result)
			}
		}
	}

	ln.activeNotifsLock.Lock()
	delete(ln.activeNotifs, notifID)
	delete(ln.notificationMeta, notifID)
	ln.activeNotifsLock.Unlock()
}
