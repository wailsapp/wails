package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
)

// UserActivityTypeBrowsingWeb is the activity type macOS uses for Handoff
// from Safari and for universal links (NSUserActivityTypeBrowsingWeb). Such
// activities carry the page in WebpageURL.
const UserActivityTypeBrowsingWeb = "NSUserActivityTypeBrowsingWeb"

// UserActivity describes what the user is doing right now so it can be
// continued on another device (Handoff), found in Spotlight or offered as a
// Siri suggestion. It mirrors NSUserActivity.
type UserActivity struct {
	// Type is a reverse-DNS identifier such as "com.example.app.editing".
	// Every type the app publishes or continues must be listed under the
	// NSUserActivityTypes key in Info.plist.
	Type string `json:"type"`
	// Title is shown to the user, for example in the Dock's Handoff icon.
	Title string `json:"title"`
	// UserInfo carries the state needed to continue the activity. It must
	// be JSON serialisable and is delivered as-is to the continuing device.
	UserInfo map[string]any `json:"userInfo"`
	// WebpageURL is a fallback for devices without the app: they open this
	// page instead. It is also where universal links arrive.
	WebpageURL string `json:"webpageURL"`
	// EligibleForHandoff advertises the activity to nearby devices.
	EligibleForHandoff bool `json:"eligibleForHandoff"`
	// EligibleForSearch indexes the activity in Spotlight.
	EligibleForSearch bool `json:"eligibleForSearch"`
	// EligibleForPrediction lets Siri suggest the activity as a shortcut.
	// The property exists on iOS and watchOS only; macOS accepts and
	// ignores it.
	EligibleForPrediction bool `json:"eligibleForPrediction"`
	// Keywords help Spotlight match the activity.
	Keywords []string `json:"keywords"`
}

func (a UserActivity) validate() error {
	if strings.TrimSpace(a.Type) == "" {
		return errors.New("activity: Type is required")
	}
	if strings.ContainsAny(a.Type, " \t\r\n") {
		return errors.New("activity: Type must not contain whitespace")
	}
	if a.WebpageURL != "" {
		parsed, err := url.Parse(a.WebpageURL)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			return fmt.Errorf("activity: WebpageURL %q must be an absolute http or https URL", a.WebpageURL)
		}
	}
	if a.Type == UserActivityTypeBrowsingWeb && a.WebpageURL == "" {
		return errors.New("activity: a browsing activity needs a WebpageURL")
	}
	if a.EligibleForHandoff && a.Title == "" && a.WebpageURL == "" {
		return errors.New("activity: a Handoff activity needs a Title or WebpageURL")
	}
	if _, err := json.Marshal(a.UserInfo); err != nil {
		return fmt.Errorf("activity: UserInfo is not JSON serialisable: %w", err)
	}
	return nil
}

// ErrActivityUnsupported is returned by Publish on platforms without
// NSUserActivity.
var ErrActivityUnsupported = errors.New("activity: user activities are not supported on this platform")

// ErrActivityInvalidated is returned by PublishedActivity.Update after
// Invalidate.
var ErrActivityInvalidated = errors.New("activity: the activity has been invalidated")

// activityImpl is the platform side of ActivityManager.
type activityImpl interface {
	publish(id uint64, activity UserActivity) error
	update(id uint64, userInfo map[string]any) error
	invalidate(id uint64)
}

// PublishedActivity is a live NSUserActivity created by Publish. It stays
// current until another activity is published or it is invalidated.
type PublishedActivity struct {
	manager  *ActivityManager
	id       uint64
	lock     sync.Mutex
	activity UserActivity
	invalid  bool
}

// Activity returns a copy of the activity as last published or updated.
func (p *PublishedActivity) Activity() UserActivity {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.activity
}

// Update replaces UserInfo on the live activity and marks it as needing to
// be saved, so Handoff peers see the new state.
func (p *PublishedActivity) Update(userInfo map[string]any) error {
	if _, err := json.Marshal(userInfo); err != nil {
		return fmt.Errorf("activity: UserInfo is not JSON serialisable: %w", err)
	}
	p.lock.Lock()
	if p.invalid {
		p.lock.Unlock()
		return ErrActivityInvalidated
	}
	p.activity.UserInfo = userInfo
	p.lock.Unlock()
	if p.manager.impl == nil {
		return nil
	}
	return p.manager.impl.update(p.id, userInfo)
}

// Invalidate ends the activity. It stops being advertised for Handoff and
// is removed from Spotlight. Calling Invalidate more than once is harmless.
func (p *PublishedActivity) Invalidate() {
	p.lock.Lock()
	if p.invalid {
		p.lock.Unlock()
		return
	}
	p.invalid = true
	p.lock.Unlock()
	p.manager.forget(p.id)
	if p.manager.impl != nil {
		p.manager.impl.invalidate(p.id)
	}
}

// ActivityManager publishes and continues user activities for Handoff,
// Spotlight and universal links.
//
// Publishing needs the NSUserActivityTypes Info.plist key listing every
// activity Type; universal links additionally need the
// com.apple.developer.associated-domains entitlement with an
// "applinks:example.com" entry and the matching apple-app-site-association
// file on that domain.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type ActivityManager struct {
	app  *App
	impl activityImpl

	nextID    atomic.Uint64
	lock      sync.RWMutex
	published map[uint64]*PublishedActivity

	handlerLock     sync.RWMutex
	continueFns     []func(*Context, UserActivity) bool
	willContinueFns []func(string) bool
	failedFns       []func(string, error)
	updatedFns      []func(UserActivity)
}

func newActivityManager(app *App) *ActivityManager {
	return &ActivityManager{
		app:       app,
		impl:      newActivityImpl(app),
		published: map[uint64]*PublishedActivity{},
	}
}

// Publish creates an NSUserActivity from activity and makes it current. The
// returned PublishedActivity can be updated or invalidated later. Publishing
// a new activity supersedes the previous current one.
//
// macOS: NSUserActivity becomeCurrent.
//
// Other platforms: ErrActivityUnsupported.
func (m *ActivityManager) Publish(activity UserActivity) (*PublishedActivity, error) {
	if err := activity.validate(); err != nil {
		return nil, err
	}
	if m.impl == nil {
		return nil, ErrActivityUnsupported
	}
	if activity.UserInfo == nil {
		activity.UserInfo = map[string]any{}
	}
	if activity.Keywords == nil {
		activity.Keywords = []string{}
	}
	published := &PublishedActivity{manager: m, id: m.nextID.Add(1), activity: activity}
	m.lock.Lock()
	m.published[published.id] = published
	m.lock.Unlock()
	if err := m.impl.publish(published.id, activity); err != nil {
		m.forget(published.id)
		return nil, err
	}
	return published, nil
}

func (m *ActivityManager) forget(id uint64) {
	m.lock.Lock()
	delete(m.published, id)
	m.lock.Unlock()
}

// OnContinue registers a handler for activities handed to this app: Handoff
// from another device, Spotlight results and universal links. The handler
// returns true when it took over the activity. Universal links
// (Type == UserActivityTypeBrowsingWeb) are additionally delivered through
// events.Common.ApplicationLaunchedWithUrl, the same event custom URL
// schemes use, so apps can keep a single URL code path.
//
// The handler runs on the main thread while AppKit waits for the answer.
func (m *ActivityManager) OnContinue(handler func(*Context, UserActivity) bool) {
	if handler == nil {
		return
	}
	m.handlerLock.Lock()
	m.continueFns = append(m.continueFns, handler)
	m.handlerLock.Unlock()
}

// OnWillContinue registers a handler that is told an activity of the given
// type is about to arrive (the data may still be in transit). Return true
// to signal the app will handle it. When no handler is registered, the app
// accepts the activity if an OnContinue handler exists.
func (m *ActivityManager) OnWillContinue(handler func(activityType string) bool) {
	if handler == nil {
		return
	}
	m.handlerLock.Lock()
	m.willContinueFns = append(m.willContinueFns, handler)
	m.handlerLock.Unlock()
}

// OnFailed registers a handler for activities that could not be continued,
// for example because the transfer from the other device failed.
func (m *ActivityManager) OnFailed(handler func(activityType string, err error)) {
	if handler == nil {
		return
	}
	m.handlerLock.Lock()
	m.failedFns = append(m.failedFns, handler)
	m.handlerLock.Unlock()
}

// OnUpdated registers a handler that observes the state of the current
// activity whenever the system saves it (application:didUpdateUserActivity:).
func (m *ActivityManager) OnUpdated(handler func(UserActivity)) {
	if handler == nil {
		return
	}
	m.handlerLock.Lock()
	m.updatedFns = append(m.updatedFns, handler)
	m.handlerLock.Unlock()
}

// willContinue is called by the platform layer. It reports whether the app
// wants the activity.
func (m *ActivityManager) willContinue(activityType string) bool {
	m.handlerLock.RLock()
	handlers := append([]func(string) bool{}, m.willContinueFns...)
	hasContinue := len(m.continueFns) > 0
	m.handlerLock.RUnlock()
	if len(handlers) == 0 {
		return hasContinue || activityType == UserActivityTypeBrowsingWeb
	}
	accepted := false
	for _, handler := range handlers {
		if handler(activityType) {
			accepted = true
		}
	}
	return accepted
}

// continueActivity is called by the platform layer with the incoming
// activity and reports whether any handler took it.
func (m *ActivityManager) continueActivity(activity UserActivity) bool {
	if activity.UserInfo == nil {
		activity.UserInfo = map[string]any{}
	}
	if activity.Keywords == nil {
		activity.Keywords = []string{}
	}
	m.handlerLock.RLock()
	handlers := append([]func(*Context, UserActivity) bool{}, m.continueFns...)
	m.handlerLock.RUnlock()
	handled := false
	for _, handler := range handlers {
		if handler(newContext(), activity) {
			handled = true
		}
	}
	return handled
}

func (m *ActivityManager) failed(activityType string, err error) {
	m.handlerLock.RLock()
	handlers := append([]func(string, error){}, m.failedFns...)
	m.handlerLock.RUnlock()
	for _, handler := range handlers {
		handler(activityType, err)
	}
}

func (m *ActivityManager) updated(activity UserActivity) {
	m.handlerLock.RLock()
	handlers := append([]func(UserActivity){}, m.updatedFns...)
	m.handlerLock.RUnlock()
	for _, handler := range handlers {
		handler(activity)
	}
}
