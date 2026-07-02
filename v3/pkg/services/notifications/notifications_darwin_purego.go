//go:build darwin && !ios && purego

package notifications

// CGO-free macOS notifications service implementation.
//
// This mirrors the behaviour of the cgo notifications_darwin.go +
// notifications_darwin.m by driving UNUserNotificationCenter directly through
// the Objective-C runtime (via github.com/ebitengine/purego) instead of
// compiling Objective-C. It keeps the notifications service buildable with
// CGO_ENABLED=0.
//
// The Objective-C bits reimplemented here are:
//   - a UNUserNotificationCenterDelegate registered at runtime with IMP
//     callbacks (willPresent + didReceiveResponse) that forward into the same
//     Go handlers the cgo //export callbacks used (captureResult,
//     didReceiveNotificationResponse).
//   - authorization request / status checks, notification send (with and
//     without actions/schedule/attachments/sound/interruption level), category
//     registration/removal, and pending/delivered removal.
//
// Completion-handler blocks handed to us by UNUserNotificationCenter are
// invoked via the Objective-C block ABI (invoke pointer at offset 16).

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type darwinNotifier struct {
	channels      map[int]chan notificationChannel
	channelsLock  sync.Mutex
	nextChannelID int
}

type notificationChannel struct {
	Success bool
	Error   error
}

type ChannelHandler interface {
	GetChannel(id int) (chan notificationChannel, bool)
}

const AppleDefaultActionIdentifier = "com.apple.UNNotificationDefaultActionIdentifier"

// Creates a new Notifications Service.
// Your app must be packaged and signed for this feature to work.
func New() *NotificationService {
	notificationServiceOnce.Do(func() {
		impl := &darwinNotifier{
			channels:      make(map[int]chan notificationChannel),
			nextChannelID: 0,
		}

		NotificationService_ = &NotificationService{
			impl: impl,
		}
	})

	return NotificationService_
}

func (dn *darwinNotifier) Startup(ctx context.Context, options application.ServiceOptions) error {
	if !isNotificationAvailable() {
		return fmt.Errorf("notifications are not available on this system")
	}
	if !checkBundleIdentifier() {
		return fmt.Errorf("notifications require a valid bundle identifier")
	}
	if !ensureDelegateInitialized() {
		return fmt.Errorf("failed to initialize notification center delegate")
	}
	return nil
}

func (dn *darwinNotifier) Shutdown() error {
	return nil
}

// nsOperatingSystemVersion mirrors Foundation's NSOperatingSystemVersion
// (three NSInteger fields), returned by value from
// -[NSProcessInfo operatingSystemVersion].
type nsOperatingSystemVersion struct {
	Major int
	Minor int
	Patch int
}

// macOSAtLeast reports whether the running OS is at least major.minor.
func macOSAtLeast(major, minor int) bool {
	pi := class("NSProcessInfo").send("processInfo")
	v := get[nsOperatingSystemVersion](pi, "operatingSystemVersion")
	if v.Major != major {
		return v.Major > major
	}
	return v.Minor >= minor
}

// isNotificationAvailable checks if notifications are available on the system.
// The cgo implementation requires macOS 11.0+ (@available(macOS 11.0, *)), so
// mirror that here in addition to checking that UNUserNotificationCenter
// actually exists.
func isNotificationAvailable() bool {
	return macOSAtLeast(11, 0) && !class("UNUserNotificationCenter").isNil()
}

func checkBundleIdentifier() bool {
	bundle := class("NSBundle").send("mainBundle")
	return !bundle.send("bundleIdentifier").isNil()
}

// RequestNotificationAuthorization requests permission for notifications.
// Default timeout is 3 minutes
func (dn *darwinNotifier) RequestNotificationAuthorization() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	id, resultCh := dn.registerChannel()

	requestNotificationAuthorization(id)

	select {
	case result := <-resultCh:
		return result.Success, result.Error
	case <-ctx.Done():
		dn.cleanupChannel(id)
		return false, fmt.Errorf("notification authorization timed out after 3 minutes: %w", ctx.Err())
	}
}

// CheckNotificationAuthorization checks current notification permission status.
func (dn *darwinNotifier) CheckNotificationAuthorization() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	id, resultCh := dn.registerChannel()

	checkNotificationAuthorization(id)

	select {
	case result := <-resultCh:
		return result.Success, result.Error
	case <-ctx.Done():
		dn.cleanupChannel(id)
		return false, fmt.Errorf("notification authorization timed out after 15s: %w", ctx.Err())
	}
}

// SendNotification sends a basic notification with a unique identifier, title, subtitle, and body.
func (dn *darwinNotifier) SendNotification(options NotificationOptions) error {
	return dn.dispatchNotification(options, false)
}

// SendNotificationWithActions sends a notification with additional actions and inputs.
// A NotificationCategory must be registered with RegisterNotificationCategory first. The `CategoryID` must match the registered category.
// If a NotificationCategory is not registered a basic notification will be sent.
func (dn *darwinNotifier) SendNotificationWithActions(options NotificationOptions) error {
	return dn.dispatchNotification(options, true)
}

// dispatchNotification is the shared body of SendNotification and
// SendNotificationWithActions.
func (dn *darwinNotifier) dispatchNotification(options NotificationOptions, withActions bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, resultCh := dn.registerChannel()
	sendNotification(id, options, withActions)

	select {
	case result := <-resultCh:
		if !result.Success {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("sending notification failed")
		}
		return nil
	case <-ctx.Done():
		dn.cleanupChannel(id)
		return fmt.Errorf("sending notification timed out: %w", ctx.Err())
	}
}

// UpdateNotification re-posts a notification with the same identifier.
// UNUserNotificationCenter auto-replaces by identifier, so the existing
// notification's content is updated in place (or freshly delivered if no
// notification with that identifier is currently pending or delivered).
func (dn *darwinNotifier) UpdateNotification(options NotificationOptions) error {
	return dn.dispatchNotification(options, options.CategoryID != "")
}

// RegisterNotificationCategory registers a new NotificationCategory to be used with SendNotificationWithActions.
// Registering a category with the same name as a previously registered NotificationCategory will override it.
func (dn *darwinNotifier) RegisterNotificationCategory(category NotificationCategory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, resultCh := dn.registerChannel()
	registerNotificationCategory(id, category)

	select {
	case result := <-resultCh:
		if !result.Success {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("category registration failed")
		}
		return nil
	case <-ctx.Done():
		dn.cleanupChannel(id)
		return fmt.Errorf("category registration timed out: %w", ctx.Err())
	}
}

// RemoveNotificationCategory remove a previously registered NotificationCategory.
func (dn *darwinNotifier) RemoveNotificationCategory(categoryId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, resultCh := dn.registerChannel()
	removeNotificationCategory(id, categoryId)

	select {
	case result := <-resultCh:
		if !result.Success {
			if result.Error != nil {
				return result.Error
			}
			return fmt.Errorf("category removal failed")
		}
		return nil
	case <-ctx.Done():
		dn.cleanupChannel(id)
		return fmt.Errorf("category removal timed out: %w", ctx.Err())
	}
}

// RemoveAllPendingNotifications removes all pending notifications.
func (dn *darwinNotifier) RemoveAllPendingNotifications() error {
	notificationCenter().send("removeAllPendingNotificationRequests")
	return nil
}

// RemovePendingNotification removes a pending notification matching the unique identifier.
func (dn *darwinNotifier) RemovePendingNotification(identifier string) error {
	notificationCenter().send("removePendingNotificationRequestsWithIdentifiers:", nsArrayOfStrings(identifier))
	return nil
}

// RemoveAllDeliveredNotifications removes all delivered notifications.
func (dn *darwinNotifier) RemoveAllDeliveredNotifications() error {
	notificationCenter().send("removeAllDeliveredNotifications")
	return nil
}

// RemoveDeliveredNotification removes a delivered notification matching the unique identifier.
func (dn *darwinNotifier) RemoveDeliveredNotification(identifier string) error {
	notificationCenter().send("removeDeliveredNotificationsWithIdentifiers:", nsArrayOfStrings(identifier))
	return nil
}

// RemoveNotification is a macOS stub that always returns nil.
// Use one of the following instead:
// RemoveAllPendingNotifications
// RemovePendingNotification
// RemoveAllDeliveredNotifications
// RemoveDeliveredNotification
// (Linux-specific)
func (dn *darwinNotifier) RemoveNotification(identifier string) error {
	return nil
}

// captureResult forwards a result for the channel registered under channelID.
// This is the Go handler that the notification-center completion-handler blocks
// call back into (the cgo //export captureResult equivalent). An empty errMsg
// means no error.
func captureResult(channelID int, success bool, errMsg string) {
	ns := getNotificationService()
	if ns == nil {
		return
	}

	handler, ok := ns.impl.(ChannelHandler)
	if !ok {
		return
	}

	resultCh, exists := handler.GetChannel(channelID)
	if !exists {
		return
	}

	var err error
	if errMsg != "" {
		err = fmt.Errorf("%s", errMsg)
	}

	resultCh <- notificationChannel{
		Success: success,
		Error:   err,
	}

	close(resultCh)
}

// didReceiveNotificationResponse is the Go handler that the delegate's
// didReceiveNotificationResponse IMP calls back into (the cgo //export
// equivalent). An empty errMsg means no error; an empty jsonPayload means a nil
// payload was received.
func didReceiveNotificationResponse(jsonPayload string, errMsg string) {
	result := NotificationResult{}

	if errMsg != "" {
		result.Error = fmt.Errorf("notification response error: %s", errMsg)
		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(result)
		}
		return
	}

	if jsonPayload == "" {
		result.Error = fmt.Errorf("received nil JSON payload in notification response")
		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(result)
		}
		return
	}

	var response NotificationResponse
	if err := json.Unmarshal([]byte(jsonPayload), &response); err != nil {
		result.Error = fmt.Errorf("failed to unmarshal notification response: %w", err)
		if ns := getNotificationService(); ns != nil {
			ns.handleNotificationResult(result)
		}
		return
	}

	if response.ActionIdentifier == AppleDefaultActionIdentifier {
		response.ActionIdentifier = DefaultActionIdentifier
	}

	result.Response = response
	if ns := getNotificationService(); ns != nil {
		ns.handleNotificationResult(result)
	}
}

// Helper methods

func (dn *darwinNotifier) registerChannel() (int, chan notificationChannel) {
	dn.channelsLock.Lock()
	defer dn.channelsLock.Unlock()

	id := dn.nextChannelID
	dn.nextChannelID++

	resultCh := make(chan notificationChannel, 1)

	dn.channels[id] = resultCh
	return id, resultCh
}

func (dn *darwinNotifier) GetChannel(id int) (chan notificationChannel, bool) {
	dn.channelsLock.Lock()
	defer dn.channelsLock.Unlock()

	ch, exists := dn.channels[id]
	if exists {
		delete(dn.channels, id)
	}
	return ch, exists
}

func (dn *darwinNotifier) cleanupChannel(id int) {
	dn.channelsLock.Lock()
	defer dn.channelsLock.Unlock()

	if ch, exists := dn.channels[id]; exists {
		delete(dn.channels, id)
		close(ch)
	}
}

// ===========================================================================
// Objective-C runtime helpers (package-local; not shared across packages).
// ===========================================================================

const (
	frameworkFoundation        = "/System/Library/Frameworks/Foundation.framework/Foundation"
	frameworkAppKit            = "/System/Library/Frameworks/AppKit.framework/AppKit"
	frameworkUserNotifications = "/System/Library/Frameworks/UserNotifications.framework/UserNotifications"
)

var frameworksOnce sync.Once

func loadFrameworks() {
	frameworksOnce.Do(func() {
		for _, fw := range []string{frameworkFoundation, frameworkAppKit, frameworkUserNotifications} {
			_, _ = purego.Dlopen(fw, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		}
	})
}

// id is a thin wrapper around objc.ID for fluent message sends.
type id objc.ID

func (o id) raw() objc.ID { return objc.ID(o) }
func (o id) isNil() bool  { return uintptr(o) == 0 }

func (o id) send(sel string, args ...any) id {
	return id(objc.ID(o).Send(sel_(sel), args...))
}

func get[T any](o id, sel string, args ...any) T {
	return objc.Send[T](objc.ID(o), sel_(sel), args...)
}

func class(name string) id {
	loadFrameworks()
	return id(objc.ID(objc.GetClass(name)))
}

var (
	selMu    sync.RWMutex
	selCache = map[string]objc.SEL{}
)

func sel_(name string) objc.SEL {
	selMu.RLock()
	s, ok := selCache[name]
	selMu.RUnlock()
	if ok {
		return s
	}
	selMu.Lock()
	defer selMu.Unlock()
	if s, ok = selCache[name]; ok {
		return s
	}
	s = objc.RegisterName(name)
	selCache[name] = s
	return s
}

func nsString(s string) id {
	return class("NSString").send("stringWithUTF8String:", s)
}

// goString converts an NSString id back to a Go string.
func (o id) string() string {
	if o.isNil() {
		return ""
	}
	cstr := get[uintptr](o, "UTF8String")
	if cstr == 0 {
		return ""
	}
	return goStringFromC(cstr)
}

// goStringFromC copies a NUL-terminated C string at the given address into a Go
// string without cgo.
func goStringFromC(p uintptr) string {
	if p == 0 {
		return ""
	}
	var n int
	for {
		b := *(*byte)(unsafe.Pointer(p + uintptr(n)))
		if b == 0 {
			break
		}
		n++
	}
	if n == 0 {
		return ""
	}
	buf := make([]byte, n)
	copy(buf, unsafe.Slice((*byte)(unsafe.Pointer(p)), n))
	return string(buf)
}

func nsData(b []byte) id {
	if len(b) == 0 {
		return class("NSData").send("data")
	}
	return class("NSData").send("dataWithBytes:length:", unsafe.Pointer(&b[0]), uint(len(b)))
}

func nsArrayOfStrings(items ...string) id {
	arr := class("NSMutableArray").send("array")
	for _, s := range items {
		arr.send("addObject:", nsString(s))
	}
	return arr
}

// nsJSONObject converts a Go value into the equivalent Foundation object graph
// (NSDictionary / NSArray / ...) by round-tripping through JSON. Returns a nil
// id on failure.
func nsJSONObject(v any) id {
	b, err := json.Marshal(v)
	if err != nil {
		return id(0)
	}
	data := nsData(b)
	return class("NSJSONSerialization").send("JSONObjectWithData:options:error:", data, uint(0), unsafe.Pointer(nil))
}

// notificationCenter returns [UNUserNotificationCenter currentNotificationCenter].
func notificationCenter() id {
	return class("UNUserNotificationCenter").send("currentNotificationCenter")
}

// errorDescription returns the localizedDescription of an NSError, prefixed to
// match the "Error: %@" format the cgo implementation used.
func errorDescription(err id) string {
	if err.isNil() {
		return ""
	}
	return "Error: " + err.send("localizedDescription").string()
}

// invokeBlockVoid invokes an Objective-C block that takes no extra arguments,
// via the block ABI (invoke pointer lives at offset 16 in the block literal).
func invokeBlockVoid(block objc.ID) {
	if block == 0 {
		return
	}
	invoke := *(*uintptr)(unsafe.Pointer(uintptr(block) + 16))
	purego.SyscallN(invoke, uintptr(block))
}

// invokeBlockUint invokes an Objective-C block taking a single NSUInteger arg.
func invokeBlockUint(block objc.ID, arg uintptr) {
	if block == 0 {
		return
	}
	invoke := *(*uintptr)(unsafe.Pointer(uintptr(block) + 16))
	purego.SyscallN(invoke, uintptr(block), arg)
}

// ===========================================================================
// UserNotifications constants
// ===========================================================================

const (
	unAuthorizationOptionBadge = 1 << 0
	unAuthorizationOptionSound = 1 << 1
	unAuthorizationOptionAlert = 1 << 2

	unAuthorizationStatusAuthorized = 2

	unNotificationActionOptionNone        = 0
	unNotificationActionOptionDestructive = 1 << 1

	unNotificationCategoryOptionNone = 0

	// Presentation options (macOS 11+): List | Banner | Sound.
	unNotificationPresentationOptionSound  = 1 << 1
	unNotificationPresentationOptionList   = 1 << 3
	unNotificationPresentationOptionBanner = 1 << 4

	// UNNotificationInterruptionLevel (macOS 12+).
	unInterruptionLevelPassive       = 0
	unInterruptionLevelActive        = 1
	unInterruptionLevelTimeSensitive = 2
	unInterruptionLevelCritical      = 3

	// NSCalendarUnit: Year | Month | Day | Hour | Minute | Second.
	nsCalendarUnitDateComponents = (1 << 2) | (1 << 3) | (1 << 4) | (1 << 5) | (1 << 6) | (1 << 7)

	nsUTF8StringEncoding = 4
)

// ===========================================================================
// Delegate registration
// ===========================================================================

var (
	delegateOnce     sync.Once
	delegateInstance id
)

func ensureDelegateInitialized() bool {
	delegateOnce.Do(func() {
		loadFrameworks()
		cls, err := objc.RegisterClass(
			"WailsNotificationDelegate",
			objc.GetClass("NSObject"),
			nil,
			nil,
			[]objc.MethodDef{
				{
					Cmd: sel_("userNotificationCenter:willPresentNotification:withCompletionHandler:"),
					Fn:  willPresentNotificationIMP,
				},
				{
					Cmd: sel_("userNotificationCenter:didReceiveNotificationResponse:withCompletionHandler:"),
					Fn:  didReceiveNotificationResponseIMP,
				},
			},
		)
		if err != nil {
			return
		}
		delegateInstance = id(objc.ID(cls)).send("alloc").send("init")
		notificationCenter().send("setDelegate:", delegateInstance)
	})

	return !delegateInstance.isNil()
}

// willPresentNotificationIMP mirrors -userNotificationCenter:willPresentNotification:withCompletionHandler:
// and asks the system to present the notification as a banner, in the list, and
// with sound while the app is in the foreground.
func willPresentNotificationIMP(self objc.ID, cmd objc.SEL, center objc.ID, notification objc.ID, completionHandler objc.ID) {
	options := uintptr(unNotificationPresentationOptionList |
		unNotificationPresentationOptionBanner |
		unNotificationPresentationOptionSound)
	invokeBlockUint(completionHandler, options)
}

// didReceiveNotificationResponseIMP mirrors -userNotificationCenter:didReceiveNotificationResponse:withCompletionHandler:.
// It builds the same JSON payload the cgo .m produced and forwards it to the Go
// handler.
func didReceiveNotificationResponseIMP(self objc.ID, cmd objc.SEL, center objc.ID, response objc.ID, completionHandler objc.ID) {
	resp := id(response)
	request := resp.send("notification").send("request")
	content := request.send("content")

	payload := class("NSMutableDictionary").send("dictionary")

	setObj := func(value id, key string) {
		payload.send("setObject:forKey:", value, nsString(key))
	}

	setObj(request.send("identifier"), "id")
	setObj(resp.send("actionIdentifier"), "actionIdentifier")

	if title := content.send("title"); !title.isNil() {
		setObj(title, "title")
	} else {
		setObj(nsString(""), "title")
	}
	if body := content.send("body"); !body.isNil() {
		setObj(body, "body")
	} else {
		setObj(nsString(""), "body")
	}
	if cat := content.send("categoryIdentifier"); !cat.isNil() {
		setObj(cat, "categoryIdentifier")
	}
	if sub := content.send("subtitle"); !sub.isNil() {
		setObj(sub, "subtitle")
	}
	if userInfo := content.send("userInfo"); !userInfo.isNil() {
		setObj(userInfo, "userInfo")
	}

	if textCls := class("UNTextInputNotificationResponse"); !textCls.isNil() {
		if get[bool](resp, "isKindOfClass:", textCls.raw()) {
			setObj(resp.send("userText"), "userText")
		}
	}

	jsonData := class("NSJSONSerialization").send("dataWithJSONObject:options:error:", payload, uint(0), unsafe.Pointer(nil))
	if jsonData.isNil() {
		didReceiveNotificationResponse("", "Error: failed to serialize notification response")
	} else {
		jsonString := class("NSString").send("alloc").send("initWithData:encoding:", jsonData, uint(nsUTF8StringEncoding))
		didReceiveNotificationResponse(jsonString.string(), "")
	}

	invokeBlockVoid(completionHandler)
}

// ===========================================================================
// Authorization
// ===========================================================================

func requestNotificationAuthorization(channelID int) {
	if !ensureDelegateInitialized() {
		captureResult(channelID, false, "Notification delegate has been lost. Reinitialize the notification service.")
		return
	}

	center := notificationCenter()
	options := uintptr(unAuthorizationOptionAlert | unAuthorizationOptionSound | unAuthorizationOptionBadge)

	handler := objc.NewBlock(func(_ objc.Block, granted bool, err objc.ID) {
		if e := id(err); !e.isNil() {
			captureResult(channelID, false, errorDescription(e))
		} else {
			captureResult(channelID, granted, "")
		}
	})
	center.send("requestAuthorizationWithOptions:completionHandler:", options, handler)
	// The framework copies the block during the call above; drop our +1.
	handler.Release()
}

func checkNotificationAuthorization(channelID int) {
	if !ensureDelegateInitialized() {
		captureResult(channelID, false, "Notification delegate has been lost. Reinitialize the notification service.")
		return
	}

	center := notificationCenter()
	handler := objc.NewBlock(func(_ objc.Block, settings objc.ID) {
		status := get[int](id(settings), "authorizationStatus")
		captureResult(channelID, status == unAuthorizationStatusAuthorized, "")
	})
	center.send("getNotificationSettingsWithCompletionHandler:", handler)
	// The framework copies the block during the call above; drop our +1.
	handler.Release()
}

// ===========================================================================
// Sending notifications
// ===========================================================================

// applySound honours the optional Sound field on NotificationOptions.
func applySound(content id, sound *NotificationSound) {
	if sound == nil {
		return
	}
	if sound.Silent {
		content.send("setSound:", objc.ID(0))
		return
	}
	if sound.Name != "" {
		content.send("setSound:", class("UNNotificationSound").send("soundNamed:", nsString(sound.Name)))
	}
}

// applyAttachments translates each NotificationAttachment into a
// UNNotificationAttachment. Failed attachments are skipped silently.
func applyAttachments(content id, attachments []NotificationAttachment) {
	if len(attachments) == 0 {
		return
	}
	arr := class("NSMutableArray").send("array")
	added := false
	for _, att := range attachments {
		if att.Path == "" {
			continue
		}
		attID := att.ID
		if attID == "" {
			attID = class("NSUUID").send("UUID").send("UUIDString").string()
		}

		var url id
		if len(att.Path) >= 7 && att.Path[:7] == "file://" {
			url = class("NSURL").send("URLWithString:", nsString(att.Path))
		} else {
			url = class("NSURL").send("fileURLWithPath:", nsString(att.Path))
		}
		if url.isNil() {
			continue
		}

		// The Type field is an optional UTI hint on macOS; only forward it as a
		// type hint when it actually looks like a UTI (contains a ".").
		var attOptions id
		if containsDot(att.Type) {
			attOptions = class("NSMutableDictionary").send("dictionary")
			// "typeHint" is the runtime value of the framework constant
			// UNNotificationAttachmentOptionsTypeHintKey.
			attOptions.send("setObject:forKey:", nsString(att.Type), nsString("typeHint"))
		}

		a := class("UNNotificationAttachment").send("attachmentWithIdentifier:URL:options:error:",
			nsString(attID), url, attOptions, unsafe.Pointer(nil))
		if a.isNil() {
			continue
		}
		arr.send("addObject:", a)
		added = true
	}
	if added {
		content.send("setAttachments:", arr)
	}
}

// applyInterruptionLevel maps the InterruptionLevel string onto
// UNNotificationInterruptionLevel (macOS 12+). On older macOS it is a no-op.
func applyInterruptionLevel(content id, level string) {
	if level == "" {
		return
	}
	if !get[bool](content, "respondsToSelector:", sel_("setInterruptionLevel:")) {
		return
	}
	var value int
	switch level {
	case "passive":
		value = unInterruptionLevelPassive
	case "active":
		value = unInterruptionLevelActive
	case "timeSensitive":
		value = unInterruptionLevelTimeSensitive
	case "critical":
		value = unInterruptionLevelCritical
	default:
		return
	}
	content.send("setInterruptionLevel:", value)
}

// buildTrigger returns a UNCalendar/UNTimeInterval trigger built from the
// optional Schedule field, or a nil id for immediate delivery.
func buildTrigger(schedule *NotificationSchedule) id {
	if schedule == nil {
		return id(0)
	}
	if schedule.DelaySeconds > 0 {
		return class("UNTimeIntervalNotificationTrigger").send("triggerWithTimeInterval:repeats:",
			float64(schedule.DelaySeconds), false)
	}
	if schedule.At > 0 {
		date := class("NSDate").send("dateWithTimeIntervalSince1970:", float64(schedule.At))
		cal := class("NSCalendar").send("currentCalendar")
		components := cal.send("components:fromDate:", uint(nsCalendarUnitDateComponents), date)
		return class("UNCalendarNotificationTrigger").send("triggerWithDateMatchingComponents:repeats:",
			components, false)
	}
	return id(0)
}

// createNotificationContent builds the UNMutableNotificationContent from a
// NotificationOptions value.
func createNotificationContent(options NotificationOptions) id {
	content := class("UNMutableNotificationContent").send("alloc").send("init")
	content.send("setTitle:", nsString(options.Title))
	if options.Subtitle != "" {
		content.send("setSubtitle:", nsString(options.Subtitle))
	}
	content.send("setBody:", nsString(options.Body))
	content.send("setSound:", class("UNNotificationSound").send("defaultSound"))

	if len(options.Data) > 0 {
		if userInfo := nsJSONObject(options.Data); !userInfo.isNil() {
			content.send("setUserInfo:", userInfo)
		}
	}

	if options.ThreadID != "" {
		content.send("setThreadIdentifier:", nsString(options.ThreadID))
	}

	applySound(content, options.Sound)
	applyAttachments(content, options.Attachments)
	applyInterruptionLevel(content, options.InterruptionLevel)

	return content
}

func sendNotification(channelID int, options NotificationOptions, withActions bool) {
	if !ensureDelegateInitialized() {
		captureResult(channelID, false, "Notification delegate has been lost. Reinitialize the notification service.")
		return
	}

	content := createNotificationContent(options)
	if withActions && options.CategoryID != "" {
		content.send("setCategoryIdentifier:", nsString(options.CategoryID))
	}

	trigger := buildTrigger(options.Schedule)

	request := class("UNNotificationRequest").send("requestWithIdentifier:content:trigger:",
		nsString(options.ID), content, trigger)

	center := notificationCenter()
	handler := objc.NewBlock(func(_ objc.Block, err objc.ID) {
		if e := id(err); !e.isNil() {
			captureResult(channelID, false, errorDescription(e))
		} else {
			captureResult(channelID, true, "")
		}
	})
	center.send("addNotificationRequest:withCompletionHandler:", request, handler)
	// The framework copies the block during the call above; drop our +1.
	handler.Release()
}

// ===========================================================================
// Categories
// ===========================================================================

func registerNotificationCategory(channelID int, category NotificationCategory) {
	if !ensureDelegateInitialized() {
		captureResult(channelID, false, "Notification delegate has been lost. Reinitialize the notification service.")
		return
	}

	nsCategoryID := nsString(category.ID)
	actions := class("NSMutableArray").send("array")

	for _, action := range category.Actions {
		if action.ID == "" || action.Title == "" {
			continue
		}
		opts := unNotificationActionOptionNone
		if action.Destructive {
			opts |= unNotificationActionOptionDestructive
		}
		a := class("UNNotificationAction").send("actionWithIdentifier:title:options:",
			nsString(action.ID), nsString(action.Title), opts)
		actions.send("addObject:", a)
	}

	if category.HasReplyField {
		textAction := class("UNTextInputNotificationAction").send(
			"actionWithIdentifier:title:options:textInputButtonTitle:textInputPlaceholder:",
			nsString("TEXT_REPLY"),
			nsString(category.ReplyButtonTitle),
			unNotificationActionOptionNone,
			nsString(category.ReplyButtonTitle),
			nsString(category.ReplyPlaceholder),
		)
		actions.send("addObject:", textAction)
	}

	newCategory := class("UNNotificationCategory").send(
		"categoryWithIdentifier:actions:intentIdentifiers:options:",
		nsCategoryID,
		actions,
		class("NSArray").send("array"),
		unNotificationCategoryOptionNone,
	)

	center := notificationCenter()

	// The completion handler runs asynchronously, after the caller's
	// autorelease pool may have drained: retain the captured (autoreleased)
	// objects now and release them inside the block once it is done with them.
	nsCategoryID.send("retain")
	newCategory.send("retain")

	handler := objc.NewBlock(func(_ objc.Block, categories objc.ID) {
		updated := class("NSMutableSet").send("setWithSet:", id(categories))

		// Remove any existing category with the same identifier.
		var existing id
		enumerator := updated.send("objectEnumerator")
		for {
			obj := enumerator.send("nextObject")
			if obj.isNil() {
				break
			}
			if get[bool](obj.send("identifier"), "isEqualToString:", nsCategoryID) {
				existing = obj
				break
			}
		}
		if !existing.isNil() {
			updated.send("removeObject:", existing)
		}

		updated.send("addObject:", newCategory)
		center.send("setNotificationCategories:", updated)

		nsCategoryID.send("release")
		newCategory.send("release")

		captureResult(channelID, true, "")
	})
	center.send("getNotificationCategoriesWithCompletionHandler:", handler)
	// The framework copies the block during the call above; drop our +1.
	handler.Release()
}

func removeNotificationCategory(channelID int, categoryID string) {
	nsCategoryID := nsString(categoryID)
	center := notificationCenter()

	// The completion handler runs asynchronously, after the caller's
	// autorelease pool may have drained: retain the captured (autoreleased)
	// NSString now and release it inside the block once it is done with it.
	nsCategoryID.send("retain")

	handler := objc.NewBlock(func(_ objc.Block, categories objc.ID) {
		updated := class("NSMutableSet").send("setWithSet:", id(categories))

		var toRemove id
		enumerator := updated.send("objectEnumerator")
		for {
			obj := enumerator.send("nextObject")
			if obj.isNil() {
				break
			}
			if get[bool](obj.send("identifier"), "isEqualToString:", nsCategoryID) {
				toRemove = obj
				break
			}
		}

		nsCategoryID.send("release")

		if !toRemove.isNil() {
			updated.send("removeObject:", toRemove)
			center.send("setNotificationCategories:", updated)
			captureResult(channelID, true, "")
		} else {
			captureResult(channelID, false, fmt.Sprintf("Category '%s' not found", categoryID))
		}
	})
	center.send("getNotificationCategoriesWithCompletionHandler:", handler)
	// The framework copies the block during the call above; drop our +1.
	handler.Release()
}

// containsDot reports whether s contains a ".".
func containsDot(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			return true
		}
	}
	return false
}
