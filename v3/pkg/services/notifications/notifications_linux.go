//go:build linux

package notifications

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type linuxNotifier struct {
	conn              *dbus.Conn
	categories        map[string]NotificationCategory
	categoriesLock    sync.RWMutex
	notifications     map[uint32]*notificationData
	notificationsLock sync.RWMutex
	scheduledTimers   map[string]*time.Timer
	scheduledLock     sync.Mutex
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
func New() *NotificationService {
	notificationServiceOnce.Do(func() {
		impl := &linuxNotifier{
			categories:      make(map[string]NotificationCategory),
			notifications:   make(map[uint32]*notificationData),
			scheduledTimers: make(map[string]*time.Timer),
		}

		NotificationService_ = &NotificationService{
			impl: impl,
		}
	})

	return NotificationService_
}

// Startup is called when the service is loaded.
func (ln *linuxNotifier) Startup(ctx context.Context, options application.ServiceOptions) error {
	ln.appName = application.Get().Config().Name

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

	if err := ln.setupSignalHandling(signalCtx); err != nil {
		return fmt.Errorf("failed to set up notification signal handling: %w", err)
	}

	return nil
}

// Shutdown will save categories and close the D-Bus connection when the service unloads.
func (ln *linuxNotifier) Shutdown() error {
	// Cancel pending scheduled timers before tearing down the D-Bus connection
	// so their callbacks cannot fire against a closed connection.
	ln.scheduledLock.Lock()
	for id, t := range ln.scheduledTimers {
		t.Stop()
		delete(ln.scheduledTimers, id)
	}
	ln.scheduledLock.Unlock()

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
	if delay, ok := scheduleDelay(options.Schedule); ok && delay > 0 {
		ln.parkSchedule(options.ID, delay, func() error {
			opts := options
			opts.Schedule = nil
			return ln.SendNotification(opts)
		})
		return nil
	}

	// An immediate send must invalidate any previously parked timer for the
	// same ID so it can't fire later and spawn a duplicate notification.
	ln.cancelScheduled(options.ID)

	defaultActionID := "default"
	actions := []string{defaultActionID, "Default"}
	actionMap := map[string]string{
		defaultActionID: DefaultActionIdentifier,
	}

	return ln.notify(options, actions, actionMap)
}

// SendNotificationWithActions sends a notification with additional actions.
func (ln *linuxNotifier) SendNotificationWithActions(options NotificationOptions) error {
	ln.categoriesLock.RLock()
	category, exists := ln.categories[options.CategoryID]
	ln.categoriesLock.RUnlock()

	if options.CategoryID == "" || !exists {
		// Fall back to basic notification
		return ln.SendNotification(options)
	}

	if delay, ok := scheduleDelay(options.Schedule); ok && delay > 0 {
		ln.parkSchedule(options.ID, delay, func() error {
			opts := options
			opts.Schedule = nil
			return ln.SendNotificationWithActions(opts)
		})
		return nil
	}

	// Same rationale as in SendNotification: kill any parked timer so it
	// cannot fire after we deliver immediately.
	ln.cancelScheduled(options.ID)

	var actions []string
	actionMap := make(map[string]string)

	defaultActionID := "default"
	actions = append(actions, defaultActionID, "Default")
	actionMap[defaultActionID] = DefaultActionIdentifier

	for _, action := range category.Actions {
		actions = append(actions, action.ID, action.Title)
		actionMap[action.ID] = action.ID
	}

	return ln.notify(options, actions, actionMap)
}

// notify is the shared D-Bus Notify caller. It wires options.Sound /
// Attachments / ThreadID / InterruptionLevel onto freedesktop hints, looks up
// any existing dbusID for options.ID to enable update-by-id via replaces_id,
// and stores the resulting dbusID for later signal handling and removal.
func (ln *linuxNotifier) notify(options NotificationOptions, actions []string, actionMap map[string]string) error {
	body := options.Body
	if options.Subtitle != "" {
		body = options.Subtitle + "\n" + body
	}

	hints := map[string]dbus.Variant{}
	hints["x-notification-id"] = dbus.MakeVariant(options.ID)

	if options.CategoryID != "" {
		hints["x-category-id"] = dbus.MakeVariant(options.CategoryID)
	}

	if options.Data != nil {
		if userData, err := json.Marshal(options.Data); err == nil {
			hints["x-user-data"] = dbus.MakeVariant(string(userData))
		}
	}

	if options.Sound != nil {
		if options.Sound.Silent {
			hints["suppress-sound"] = dbus.MakeVariant(true)
		} else if options.Sound.Name != "" {
			hints["sound-name"] = dbus.MakeVariant(options.Sound.Name)
		}
	}

	for _, a := range options.Attachments {
		if a.Path == "" {
			continue
		}
		// freedesktop spec only honours one image hint; use the first entry
		// whose extension is a recognised image type. Non-image paths are
		// skipped to prevent arbitrary files from being surfaced to the daemon.
		if isImagePath(a.Path) {
			hints["image-path"] = dbus.MakeVariant(a.Path)
			break
		}
	}

	if options.ThreadID != "" {
		// Use category as the closest stable freedesktop hint for grouping
		// when no explicit category id has been set.
		if _, has := hints["x-category-id"]; !has {
			hints["category"] = dbus.MakeVariant(options.ThreadID)
		}
		hints["x-wails-thread-id"] = dbus.MakeVariant(options.ThreadID)
	}

	switch options.InterruptionLevel {
	case InterruptionLevelPassive:
		hints["urgency"] = dbus.MakeVariant(byte(0))
	case InterruptionLevelTimeSensitive, InterruptionLevelCritical:
		hints["urgency"] = dbus.MakeVariant(byte(2))
	case InterruptionLevelActive, "":
		// default urgency
	}

	// Look up an existing tracked dbusID for this options.ID so D-Bus
	// replaces in place rather than spawning a duplicate. Zero means new.
	var replacesID uint32
	ln.notificationsLock.RLock()
	for id, n := range ln.notifications {
		if n.ID == options.ID {
			replacesID = id
			break
		}
	}
	ln.notificationsLock.RUnlock()

	obj := ln.conn.Object(dbusNotificationInterface, dbusNotificationPath)
	call := obj.Call(
		dbusNotificationInterface+".Notify",
		0,
		ln.appName,
		replacesID,
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
	// Drop the previous tracking entry under the old dbusID if replacing.
	if replacesID != 0 && replacesID != dbusID {
		delete(ln.notifications, replacesID)
	}
	ln.notifications[dbusID] = notification
	ln.notificationsLock.Unlock()

	return nil
}

// parkSchedule defers `fire` by `delay`. Any previously parked timer for the
// same id is cancelled and removed eagerly, and the parked callback only
// removes its map entry if it is still the live owner — that guards against a
// just-fired-but-not-yet-cleaned-up callback racing in to delete the entry
// for a freshly-parked successor.
func (ln *linuxNotifier) parkSchedule(id string, delay time.Duration, fire func() error) {
	ln.scheduledLock.Lock()
	if existing, ok := ln.scheduledTimers[id]; ok {
		existing.Stop()
		delete(ln.scheduledTimers, id)
	}
	var timer *time.Timer
	timer = time.AfterFunc(delay, func() {
		if err := fire(); err != nil {
			fmt.Printf("scheduled notification %s push failed: %v\n", id, err)
		}
		ln.scheduledLock.Lock()
		defer ln.scheduledLock.Unlock()
		if cur, ok := ln.scheduledTimers[id]; ok && cur == timer {
			delete(ln.scheduledTimers, id)
		}
	})
	ln.scheduledTimers[id] = timer
	ln.scheduledLock.Unlock()
}

// cancelScheduled stops and removes any parked timer for id. Safe to call when
// no timer is parked.
func (ln *linuxNotifier) cancelScheduled(id string) {
	ln.scheduledLock.Lock()
	defer ln.scheduledLock.Unlock()
	if t, ok := ln.scheduledTimers[id]; ok {
		t.Stop()
		delete(ln.scheduledTimers, id)
	}
}

// UpdateNotification re-posts a notification by ID. The Send path passes
// replaces_id when an existing notification with the same options.ID is
// tracked, so D-Bus updates in place.
func (ln *linuxNotifier) UpdateNotification(options NotificationOptions) error {
	if options.CategoryID != "" {
		return ln.SendNotificationWithActions(options)
	}
	return ln.SendNotification(options)
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

// RemoveAllPendingNotifications cancels any in-process scheduled timers AND
// closes every currently-tracked delivered notification.
func (ln *linuxNotifier) RemoveAllPendingNotifications() error {
	ln.scheduledLock.Lock()
	for id, t := range ln.scheduledTimers {
		t.Stop()
		delete(ln.scheduledTimers, id)
	}
	ln.scheduledLock.Unlock()

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

// RemovePendingNotification cancels a scheduled timer for the given identifier
// and/or closes the currently-tracked delivered notification.
func (ln *linuxNotifier) RemovePendingNotification(identifier string) error {
	ln.scheduledLock.Lock()
	if t, ok := ln.scheduledTimers[identifier]; ok {
		t.Stop()
		delete(ln.scheduledTimers, identifier)
	}
	ln.scheduledLock.Unlock()

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

// isImagePath reports whether path has an extension recognised as an image by
// freedesktop notification daemons. Only image files are valid for the
// image-path hint; passing arbitrary paths (e.g. text files) is rejected.
func isImagePath(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".svg", ".ico", ".webp", ".tiff", ".tif":
		return true
	}
	return false
}
