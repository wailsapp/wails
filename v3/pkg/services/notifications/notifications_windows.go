//go:build windows

package notifications

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	_ "unsafe"

	wintoast "git.sr.ht/~jackmordaunt/go-toast/v2/wintoast"
	"github.com/wailsapp/wails/v3/internal/uuid"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"golang.org/x/sys/windows/registry"
)

type windowsNotifier struct {
	categories      map[string]NotificationCategory
	categoriesLock  sync.RWMutex
	appName         string
	appGUID         string
	iconPath        string
	exePath         string
	scheduledLock   sync.Mutex
	scheduledTimers map[string]*time.Timer
}

const (
	ToastRegistryPath                  = `Software\Classes\AppUserModelId\`
	ToastRegistryGuidKey               = "CustomActivator"
	NotificationCategoriesRegistryPath = `SOFTWARE\%s\NotificationCategories`
	NotificationCategoriesRegistryKey  = "Categories"
)

// NotificationPayload combines the action ID and user data into a single structure
type NotificationPayload struct {
	Action  string              `json:"action"`
	Options NotificationOptions `json:"payload,omitempty"`
}

// Creates a new Notifications Service.
func New() *NotificationService {
	notificationServiceOnce.Do(func() {
		impl := &windowsNotifier{
			categories:      make(map[string]NotificationCategory),
			scheduledTimers: make(map[string]*time.Timer),
		}

		NotificationService_ = &NotificationService{
			impl: impl,
		}
	})

	return NotificationService_
}

//go:linkname registerFactoryInternal git.sr.ht/~jackmordaunt/go-toast/v2/wintoast.registerClassFactory
func registerFactoryInternal(factory *wintoast.IClassFactory) error

// Startup is called when the service is loaded
// Sets an activation callback to emit an event when notifications are interacted with.
func (wn *windowsNotifier) Startup(ctx context.Context, options application.ServiceOptions) error {
	wn.categoriesLock.Lock()
	defer wn.categoriesLock.Unlock()

	app := application.Get()
	cfg := app.Config()

	wn.appName = cfg.Name

	guid, err := wn.getGUID()
	if err != nil {
		return err
	}
	wn.appGUID = guid

	wn.iconPath = filepath.Join(os.TempDir(), wn.appName+wn.appGUID+".png")

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	wn.exePath = exe

	// Create the registry key for the toast activator
	key, _, err := registry.CreateKey(registry.CURRENT_USER,
		`Software\Classes\CLSID\`+wn.appGUID+`\LocalServer32`, registry.ALL_ACCESS)
	if err != nil {
		return fmt.Errorf("failed to create CLSID key: %w", err)
	}

	if err := key.SetStringValue("", fmt.Sprintf("\"%s\" %%1", wn.exePath)); err != nil {
		return fmt.Errorf("failed to set CLSID server path: %w", err)
	}
	key.Close()

	if err := wintoast.SetAppData(wintoast.AppData{
		AppID:         wn.appName,
		GUID:          guid,
		IconPath:      wn.iconPath,
		ActivationExe: wn.exePath,
	}); err != nil {
		return fmt.Errorf("failed to set toast app data: %w", err)
	}

	wintoast.SetActivationCallback(func(_ string, args string, data []wintoast.UserData) {
		result := NotificationResult{}

		actionIdentifier, options, err := parseNotificationResponse(args)

		if err != nil {
			result.Error = err
		} else {
			// Subtitle is retained but was not shown with the notification
			response := NotificationResponse{
				ID:               options.ID,
				ActionIdentifier: actionIdentifier,
				Title:            options.Title,
				Subtitle:         options.Subtitle,
				Body:             options.Body,
				CategoryID:       options.CategoryID,
				UserInfo:         options.Data,
			}

			if userText, found := wn.getUserText(data); found {
				response.UserText = userText
			}

			result.Response = response
		}

		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(result)
		}
	})

	// Register the class factory eagerly so activation callbacks fire even if
	// the user clicks a toast that pre-dates the first call to wintoast.Push
	// (e.g. a notification still pinned in Action Center from a previous run).
	if err := registerFactoryInternal(wintoast.ClassFactory); err != nil {
		return fmt.Errorf("CoRegisterClassObject failed: %w", err)
	}

	return wn.loadCategoriesFromRegistry()
}

// Shutdown will attempt to save the categories to the registry when the service unloads
func (wn *windowsNotifier) Shutdown() error {
	// Cancel pending scheduled timers before tearing down COM/wintoast so their
	// callbacks cannot fire against a deregistered toast activator.
	wn.scheduledLock.Lock()
	for id, t := range wn.scheduledTimers {
		t.Stop()
		delete(wn.scheduledTimers, id)
	}
	wn.scheduledLock.Unlock()

	wn.categoriesLock.Lock()
	defer wn.categoriesLock.Unlock()

	return wn.saveCategoriesToRegistry()
}

// RequestNotificationAuthorization is a Windows stub that always returns true, nil.
// (user authorization is macOS-specific)
func (wn *windowsNotifier) RequestNotificationAuthorization() (bool, error) {
	return true, nil
}

// CheckNotificationAuthorization is a Windows stub that always returns true.
// (user authorization is macOS-specific)
func (wn *windowsNotifier) CheckNotificationAuthorization() (bool, error) {
	return true, nil
}

// SendNotification sends a basic notification with a name, title, and body. All other options are ignored on Windows.
// (subtitle is only available on macOS and Linux)
func (wn *windowsNotifier) SendNotification(options NotificationOptions) error {
	return wn.sendOrSchedule(options, nil)
}

// sendOrSchedule pushes the toast immediately, or parks it on a time.AfterFunc
// timer if options.Schedule is set. The scheduled timer is tracked by ID so
// RemovePendingNotification can cancel it. An immediate send invalidates any
// previously parked timer for the same ID so it cannot fire afterwards and
// spawn a duplicate notification.
func (wn *windowsNotifier) sendOrSchedule(options NotificationOptions, category *NotificationCategory) error {
	if delay, ok := scheduleDelay(options.Schedule); ok && delay > 0 {
		wn.parkSchedule(options.ID, delay, options, category)
		return nil
	}
	wn.cancelScheduled(options.ID)
	return wn.pushNow(options, category)
}

// parkSchedule defers a delivery by `delay`. Any previously parked timer for
// the same id is cancelled and removed eagerly; the parked callback only
// removes its map entry if it is still the live owner so a just-fired-but-
// not-yet-cleaned-up callback cannot delete the entry for a freshly-parked
// successor.
func (wn *windowsNotifier) parkSchedule(id string, delay time.Duration, options NotificationOptions, category *NotificationCategory) {
	wn.scheduledLock.Lock()
	if existing, ok := wn.scheduledTimers[id]; ok {
		existing.Stop()
		delete(wn.scheduledTimers, id)
	}
	opts := options
	opts.Schedule = nil
	cat := category
	var timer *time.Timer
	timer = time.AfterFunc(delay, func() {
		if err := wn.pushNow(opts, cat); err != nil {
			fmt.Printf("scheduled notification %s push failed: %v\n", opts.ID, err)
		}
		wn.scheduledLock.Lock()
		defer wn.scheduledLock.Unlock()
		if cur, ok := wn.scheduledTimers[id]; ok && cur == timer {
			delete(wn.scheduledTimers, id)
		}
	})
	wn.scheduledTimers[id] = timer
	wn.scheduledLock.Unlock()
}

// cancelScheduled stops and removes any parked timer for id. Safe to call when
// no timer is parked.
func (wn *windowsNotifier) cancelScheduled(id string) {
	wn.scheduledLock.Lock()
	defer wn.scheduledLock.Unlock()
	if t, ok := wn.scheduledTimers[id]; ok {
		t.Stop()
		delete(wn.scheduledTimers, id)
	}
}

func (wn *windowsNotifier) pushNow(options NotificationOptions, category *NotificationCategory) error {
	if err := wn.saveIconToDir(); err != nil {
		fmt.Printf("Error saving icon: %v\n", err)
	}

	defaultArgs, err := wn.encodePayload(DefaultActionIdentifier, options)
	if err != nil {
		return fmt.Errorf("failed to encode notification payload: %w", err)
	}

	doc, err := wn.buildToastXML(options, category, defaultArgs)
	if err != nil {
		return err
	}

	return wintoast.Push(wn.appName, doc, wintoast.PowershellFallback)
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
// A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
// If a NotificationCategory is not registered a basic notification will be sent.
// (subtitle is only available on macOS and Linux)
func (wn *windowsNotifier) SendNotificationWithActions(options NotificationOptions) error {
	wn.categoriesLock.RLock()
	nCategory, categoryExists := wn.categories[options.CategoryID]
	wn.categoriesLock.RUnlock()

	if options.CategoryID == "" || !categoryExists {
		fmt.Printf("Category '%s' not found, sending basic notification without actions\n", options.CategoryID)
		return wn.SendNotification(options)
	}

	return wn.sendOrSchedule(options, &nCategory)
}

// UpdateNotification re-posts a toast with the same options. On Windows
// (with the current wintoast surface) this delivers a fresh notification;
// true replace-by-tag/group requires upstream support that isn't exposed yet.
// Documented limitation in the package godoc.
func (wn *windowsNotifier) UpdateNotification(options NotificationOptions) error {
	if options.CategoryID != "" {
		return wn.SendNotificationWithActions(options)
	}
	return wn.SendNotification(options)
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
// Registering a category with the same name as a previously registered NotificationCategory will override it.
func (wn *windowsNotifier) RegisterNotificationCategory(category NotificationCategory) error {
	wn.categoriesLock.Lock()
	defer wn.categoriesLock.Unlock()

	wn.categories[category.ID] = NotificationCategory{
		ID:               category.ID,
		Actions:          category.Actions,
		HasReplyField:    bool(category.HasReplyField),
		ReplyPlaceholder: category.ReplyPlaceholder,
		ReplyButtonTitle: category.ReplyButtonTitle,
	}

	return wn.saveCategoriesToRegistry()
}

// RemoveNotificationCategory removes a previously registered NotificationCategory.
func (wn *windowsNotifier) RemoveNotificationCategory(categoryId string) error {
	wn.categoriesLock.Lock()
	defer wn.categoriesLock.Unlock()

	delete(wn.categories, categoryId)

	return wn.saveCategoriesToRegistry()
}

// RemoveAllPendingNotifications cancels every notification scheduled via
// time.AfterFunc that has not yet fired. Already-delivered toasts are not
// affected on Windows (no API for retracting from Action Center via wintoast).
func (wn *windowsNotifier) RemoveAllPendingNotifications() error {
	wn.scheduledLock.Lock()
	defer wn.scheduledLock.Unlock()
	for id, t := range wn.scheduledTimers {
		t.Stop()
		delete(wn.scheduledTimers, id)
	}
	return nil
}

// RemovePendingNotification cancels a scheduled (not-yet-fired) notification
// matching the given identifier. No-op if no such timer is parked.
func (wn *windowsNotifier) RemovePendingNotification(identifier string) error {
	wn.scheduledLock.Lock()
	defer wn.scheduledLock.Unlock()
	if t, ok := wn.scheduledTimers[identifier]; ok {
		t.Stop()
		delete(wn.scheduledTimers, identifier)
	}
	return nil
}

// RemoveAllDeliveredNotifications is a Windows stub that always returns nil.
// (macOS and Linux only)
func (wn *windowsNotifier) RemoveAllDeliveredNotifications() error {
	return nil
}

// RemoveDeliveredNotification is a Windows stub that always returns nil.
// (macOS and Linux only)
func (wn *windowsNotifier) RemoveDeliveredNotification(_ string) error {
	return nil
}

// RemoveNotification is a Windows stub that always returns nil.
// (Linux-specific)
func (wn *windowsNotifier) RemoveNotification(identifier string) error {
	return nil
}

// toastDoc is the root element of a Windows toast notification XML document.
// Marshalling these structs with encoding/xml lets the Go stdlib handle attribute
// escaping, which the previous text/template approach did not.
// See https://learn.microsoft.com/en-us/uwp/schemas/tiles/toastschema/element-toast
type toastDoc struct {
	XMLName        xml.Name      `xml:"toast"`
	ActivationType string        `xml:"activationType,attr,omitempty"`
	Launch         string        `xml:"launch,attr,omitempty"`
	Duration       string        `xml:"duration,attr,omitempty"`
	Scenario       string        `xml:"scenario,attr,omitempty"`
	Visual         toastVisual   `xml:"visual"`
	Audio          *toastAudio   `xml:"audio,omitempty"`
	Actions        *toastActions `xml:"actions,omitempty"`
	Header         *toastHeader  `xml:"header,omitempty"`
}

type toastVisual struct {
	Binding toastBinding `xml:"binding"`
}

type toastBinding struct {
	Template string       `xml:"template,attr"`
	Texts    []toastText  `xml:"text"`
	Images   []toastImage `xml:"image"`
}

type toastText struct {
	HintMaxLines string `xml:"hint-maxLines,attr,omitempty"`
	Value        string `xml:",chardata"`
}

type toastImage struct {
	Src       string `xml:"src,attr"`
	Placement string `xml:"placement,attr,omitempty"`
	HintCrop  string `xml:"hint-crop,attr,omitempty"`
}

type toastAudio struct {
	Src    string `xml:"src,attr,omitempty"`
	Loop   string `xml:"loop,attr,omitempty"`
	Silent string `xml:"silent,attr,omitempty"`
}

type toastActions struct {
	Inputs  []toastInput  `xml:"input"`
	Actions []toastAction `xml:"action"`
}

type toastInput struct {
	ID                 string `xml:"id,attr"`
	Type               string `xml:"type,attr"`
	Title              string `xml:"title,attr,omitempty"`
	PlaceHolderContent string `xml:"placeHolderContent,attr,omitempty"`
}

type toastAction struct {
	Content        string `xml:"content,attr"`
	Arguments      string `xml:"arguments,attr"`
	ActivationType string `xml:"activationType,attr,omitempty"`
	HintInputID    string `xml:"hint-inputId,attr,omitempty"`
}

type toastHeader struct {
	ID        string `xml:"id,attr"`
	Title     string `xml:"title,attr"`
	Arguments string `xml:"arguments,attr,omitempty"`
}

// resolveSoundSrc maps a Sound.Name onto a toast XML src URI.
// Names that already use ms-winsoundevent:/ms-appx: are passed through.
// Bare names are wrapped in the ms-winsoundevent: scheme for built-in events.
func resolveSoundSrc(name string) string {
	if strings.HasPrefix(name, "ms-winsoundevent:") || strings.HasPrefix(name, "ms-appx:") {
		return name
	}
	return "ms-winsoundevent:Notification." + name
}

// scenarioForLevel maps NotificationOptions.InterruptionLevel onto the
// <toast scenario="..."> attribute. Returns "" for default delivery.
func scenarioForLevel(level string) string {
	switch level {
	case InterruptionLevelTimeSensitive:
		return "reminder"
	case InterruptionLevelCritical:
		return "urgent"
	default:
		return ""
	}
}

// buildToastXML produces a Windows toast XML document for the given options
// and optional category. The defaultArgs string is the activation argument
// emitted when the user clicks the toast body itself (not a specific action).
func (wn *windowsNotifier) buildToastXML(options NotificationOptions, category *NotificationCategory, defaultArgs string) (string, error) {
	t := toastDoc{
		ActivationType: "foreground",
		Launch:         defaultArgs,
		Duration:       "short",
		Scenario:       scenarioForLevel(options.InterruptionLevel),
		Visual: toastVisual{
			Binding: toastBinding{Template: "ToastGeneric"},
		},
		Audio: &toastAudio{
			Src:  "ms-winsoundevent:Notification.Default",
			Loop: "false",
		},
	}

	if options.Sound != nil {
		if options.Sound.Silent {
			t.Audio = &toastAudio{Silent: "true"}
		} else if options.Sound.Name != "" {
			t.Audio = &toastAudio{
				Src:  resolveSoundSrc(options.Sound.Name),
				Loop: "false",
			}
		}
	}

	if options.Title != "" {
		t.Visual.Binding.Texts = append(t.Visual.Binding.Texts, toastText{
			HintMaxLines: "1",
			Value:        options.Title,
		})
	}
	if options.Body != "" {
		t.Visual.Binding.Texts = append(t.Visual.Binding.Texts, toastText{
			Value: options.Body,
		})
	}

	for _, a := range options.Attachments {
		if a.Path == "" {
			continue
		}
		img := toastImage{Src: a.Path}
		switch a.Type {
		case "hero":
			img.Placement = "hero"
		case "appLogoOverride":
			img.Placement = "appLogoOverride"
			img.HintCrop = "circle"
		case "inline", "":
			// no placement attribute = inline body image
		default:
			img.Placement = a.Type
		}
		t.Visual.Binding.Images = append(t.Visual.Binding.Images, img)
	}

	if options.ThreadID != "" {
		t.Header = &toastHeader{
			ID:    options.ThreadID,
			Title: options.ThreadID,
		}
	}

	if category != nil {
		actions := &toastActions{}

		if category.HasReplyField {
			actions.Inputs = append(actions.Inputs, toastInput{
				ID:                 "userText",
				Type:               "text",
				PlaceHolderContent: category.ReplyPlaceholder,
			})
		}

		for _, a := range category.Actions {
			args, err := wn.encodePayload(a.ID, options)
			if err != nil {
				return "", fmt.Errorf("failed to encode action payload: %w", err)
			}
			actions.Actions = append(actions.Actions, toastAction{
				Content:        a.Title,
				Arguments:      args,
				ActivationType: "foreground",
			})
		}

		if category.HasReplyField {
			args, err := wn.encodePayload("TEXT_REPLY", options)
			if err != nil {
				return "", fmt.Errorf("failed to encode reply payload: %w", err)
			}
			actions.Actions = append(actions.Actions, toastAction{
				Content:        category.ReplyButtonTitle,
				Arguments:      args,
				ActivationType: "foreground",
				HintInputID:    "userText",
			})
		}

		if len(actions.Inputs) > 0 || len(actions.Actions) > 0 {
			t.Actions = actions
		}
	}

	body, err := xml.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal toast xml: %w", err)
	}
	return xml.Header + string(body), nil
}

// encodePayload combines an action ID and user data into a single encoded string
func (wn *windowsNotifier) encodePayload(actionID string, options NotificationOptions) (string, error) {
	payload := NotificationPayload{
		Action:  actionID,
		Options: options,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return actionID, err
	}

	encodedPayload := base64.StdEncoding.EncodeToString(jsonData)
	return encodedPayload, nil
}

// decodePayload extracts the action ID and user data from an encoded payload
func decodePayload(encodedString string) (string, NotificationOptions, error) {
	jsonData, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return encodedString, NotificationOptions{}, fmt.Errorf("failed to decode base64 payload: %w", err)
	}

	var payload NotificationPayload
	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return encodedString, NotificationOptions{}, fmt.Errorf("failed to unmarshal notification payload: %w", err)
	}

	return payload.Action, payload.Options, nil
}

// parseNotificationResponse updated to use structured payload decoding
func parseNotificationResponse(response string) (action string, options NotificationOptions, err error) {
	actionID, options, err := decodePayload(response)

	if err != nil {
		fmt.Printf("Warning: Failed to decode notification response: %v\n", err)
		return response, NotificationOptions{}, err
	}

	return actionID, options, nil
}

func (wn *windowsNotifier) saveIconToDir() error {
	icon, err := application.NewIconFromResource(w32.GetModuleHandle(""), uint16(3))
	if err != nil {
		return fmt.Errorf("failed to retrieve application icon: %w", err)
	}

	return w32.SaveHIconAsPNG(icon, wn.iconPath)
}

func (wn *windowsNotifier) saveCategoriesToRegistry() error {
	// We assume lock is held by caller

	registryPath := fmt.Sprintf(NotificationCategoriesRegistryPath, wn.appName)

	key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		registryPath,
		registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	data, err := json.Marshal(wn.categories)
	if err != nil {
		return err
	}

	return key.SetStringValue(NotificationCategoriesRegistryKey, string(data))
}

func (wn *windowsNotifier) loadCategoriesFromRegistry() error {
	// We assume lock is held by caller

	registryPath := fmt.Sprintf(NotificationCategoriesRegistryPath, wn.appName)

	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		registryPath,
		registry.QUERY_VALUE,
	)
	if err != nil {
		if err == registry.ErrNotExist {
			// Not an error, no saved categories
			return nil
		}
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	data, _, err := key.GetStringValue(NotificationCategoriesRegistryKey)
	if err != nil {
		if err == registry.ErrNotExist {
			// No value yet, but key exists
			return nil
		}
		return fmt.Errorf("failed to read categories from registry: %w", err)
	}

	categories := make(map[string]NotificationCategory)
	if err := json.Unmarshal([]byte(data), &categories); err != nil {
		return fmt.Errorf("failed to parse notification categories from registry: %w", err)
	}

	wn.categories = categories

	return nil
}

func (wn *windowsNotifier) getUserText(data []wintoast.UserData) (string, bool) {
	for _, d := range data {
		if d.Key == "userText" {
			return d.Value, true
		}
	}
	return "", false
}

func (wn *windowsNotifier) getGUID() (string, error) {
	keyPath := ToastRegistryPath + wn.appName

	k, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
	if err == nil {
		guid, _, err := k.GetStringValue(ToastRegistryGuidKey)
		k.Close()
		if err == nil && guid != "" {
			return guid, nil
		}
	}

	guid := wn.generateGUID()

	k, _, err = registry.CreateKey(registry.CURRENT_USER, keyPath, registry.WRITE)
	if err != nil {
		return "", fmt.Errorf("failed to create registry key: %w", err)
	}
	defer k.Close()

	if err := k.SetStringValue(ToastRegistryGuidKey, guid); err != nil {
		return "", fmt.Errorf("failed to write GUID to registry: %w", err)
	}

	return guid, nil
}

func (wn *windowsNotifier) generateGUID() string {
	guid := uuid.New()
	return fmt.Sprintf("{%s}", guid.String())
}
