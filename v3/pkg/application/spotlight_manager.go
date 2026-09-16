package application

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// ErrSpotlightNotSupported is returned by SpotlightManager methods on
// platforms without a system search index.
var ErrSpotlightNotSupported = errors.New("spotlight indexing is not supported on this platform")

// spotlightDefaultContentType is the uniform type used when a SearchableItem
// does not name one.
const spotlightDefaultContentType = "public.item"

// Activity types delivered through application:continueUserActivity: when
// the user selects one of the application's items in Spotlight
// (CSSearchableItemActionType) or asks to continue a search in the app
// (CSQueryContinuationActionType).
const (
	spotlightItemActionType        = "com.apple.corespotlightitem"
	spotlightQueryContinuationType = "com.apple.corespotlightcontinuation"
	spotlightActivityIdentifierKey = "kCSSearchableItemActivityIdentifier"
	spotlightQueryStringKey        = "kCSSearchQueryString"
)

// SearchableItem describes one entry in the system search index.
type SearchableItem struct {
	// ID identifies the item within the application; indexing an item with
	// an existing ID replaces it. Required.
	ID string
	// Domain groups items so they can be removed together with
	// DeleteDomain, for example "notes" or "notes.archived" (a dot separates
	// nested domains, and deleting a parent removes its children).
	Domain string
	// Title is shown as the result's main line. Required.
	Title string
	// Description is shown under the title.
	Description string
	// Keywords are extra search terms not visible in the result.
	Keywords []string
	// ContentType is a uniform type identifier such as "public.plain-text"
	// or "public.image" that picks the result's icon and category. Empty
	// means "public.item".
	ContentType string
	// ThumbnailPNG is an optional image shown next to the result.
	ThumbnailPNG []byte
	// URL is an optional link stored with the item, typically a deep link
	// the application resolves in OnOpen.
	URL string
	// ExpiresAt removes the item from the index automatically after the
	// given time. Zero means the item does not expire.
	ExpiresAt time.Time
}

func (item SearchableItem) validate() error {
	if strings.TrimSpace(item.ID) == "" {
		return errors.New("spotlight: item ID is empty")
	}
	if strings.TrimSpace(item.Title) == "" {
		return fmt.Errorf("spotlight: item %q has no title", item.ID)
	}
	if len(item.ThumbnailPNG) > 0 && !isPNG(item.ThumbnailPNG) {
		return fmt.Errorf("spotlight: item %q thumbnail is not a PNG", item.ID)
	}
	return nil
}

// isPNG reports whether data starts with the PNG signature.
func isPNG(data []byte) bool {
	const signature = "\x89PNG\r\n\x1a\n"
	return len(data) >= len(signature) && string(data[:len(signature)]) == signature
}

func validateSearchableItems(items []SearchableItem) error {
	if len(items) == 0 {
		return errors.New("spotlight: no items to index")
	}
	seen := map[string]bool{}
	for _, item := range items {
		if err := item.validate(); err != nil {
			return err
		}
		if seen[item.ID] {
			return fmt.Errorf("spotlight: item ID %q appears more than once", item.ID)
		}
		seen[item.ID] = true
	}
	return nil
}

// spotlightWireItem is the JSON form handed to the native indexer.
type spotlightWireItem struct {
	ID           string   `json:"id"`
	Domain       string   `json:"domain,omitempty"`
	Title        string   `json:"title"`
	Description  string   `json:"description,omitempty"`
	Keywords     []string `json:"keywords,omitempty"`
	ContentType  string   `json:"contentType"`
	ThumbnailPNG string   `json:"thumbnail,omitempty"`
	URL          string   `json:"url,omitempty"`
	ExpiresAt    int64    `json:"expires,omitempty"`
}

func spotlightItemsJSON(items []SearchableItem) (string, error) {
	wire := make([]spotlightWireItem, 0, len(items))
	for _, item := range items {
		contentType := strings.TrimSpace(item.ContentType)
		if contentType == "" {
			contentType = spotlightDefaultContentType
		}
		entry := spotlightWireItem{
			ID:          item.ID,
			Domain:      item.Domain,
			Title:       item.Title,
			Description: item.Description,
			Keywords:    item.Keywords,
			ContentType: contentType,
			URL:         item.URL,
		}
		if len(item.ThumbnailPNG) > 0 {
			entry.ThumbnailPNG = base64.StdEncoding.EncodeToString(item.ThumbnailPNG)
		}
		if !item.ExpiresAt.IsZero() {
			entry.ExpiresAt = item.ExpiresAt.Unix()
		}
		wire = append(wire, entry)
	}
	data, err := json.Marshal(wire)
	if err != nil {
		return "", fmt.Errorf("spotlight: %w", err)
	}
	return string(data), nil
}

// SpotlightOpenHandler receives the item ID the user selected in Spotlight
// (query is then empty) or, for "Search in App" continuations, the query
// text (id is then empty).
type SpotlightOpenHandler func(ctx *Context, id string, query string)

// SpotlightManager indexes application content for system search.
//
// Platform support is documented on each method. Methods without support on
// the current platform are safe no-ops or return a documented error.
type SpotlightManager struct {
	app *App

	lock         sync.Mutex
	openHandlers map[uint]SpotlightOpenHandler
	nextHandler  uint
	// bridged records that the ActivityManager continuation bridge has been
	// registered; it is installed on the first OnOpen call so applications
	// that never use Spotlight keep their user activity behaviour unchanged.
	bridged bool
}

func newSpotlightManager(app *App) *SpotlightManager {
	return &SpotlightManager{app: app, openHandlers: map[uint]SpotlightOpenHandler{}}
}

// IsAvailable reports whether the system search index accepts items.
// Indexing needs a bundled application: results from an unbundled binary
// never show up in Spotlight even when the index calls succeed.
//
// macOS: CSSearchableIndex isIndexingAvailable. Other platforms: false.
func (s *SpotlightManager) IsAvailable() bool {
	return spotlightIsAvailable()
}

// Index adds items to the search index, replacing any with the same ID, and
// waits for the index to confirm.
//
// macOS 10.13+: CSSearchableIndex. Other platforms return ErrSpotlightNotSupported.
func (s *SpotlightManager) Index(items []SearchableItem) error {
	if err := validateSearchableItems(items); err != nil {
		return err
	}
	payload, err := spotlightItemsJSON(items)
	if err != nil {
		return err
	}
	return spotlightIndex(payload)
}

// Delete removes the items with the given IDs.
//
// macOS 10.13+: CSSearchableIndex. Other platforms return ErrSpotlightNotSupported.
func (s *SpotlightManager) Delete(ids []string) error {
	if len(ids) == 0 {
		return errors.New("spotlight: no item IDs to delete")
	}
	for _, id := range ids {
		if strings.TrimSpace(id) == "" {
			return errors.New("spotlight: item ID is empty")
		}
	}
	return spotlightDelete(ids)
}

// DeleteDomain removes every item in the domain and its nested domains.
//
// macOS 10.13+: CSSearchableIndex. Other platforms return ErrSpotlightNotSupported.
func (s *SpotlightManager) DeleteDomain(domain string) error {
	if strings.TrimSpace(domain) == "" {
		return errors.New("spotlight: domain is empty")
	}
	return spotlightDeleteDomain(domain)
}

// DeleteAll removes every item the application indexed.
//
// macOS 10.13+: CSSearchableIndex. Other platforms return ErrSpotlightNotSupported.
func (s *SpotlightManager) DeleteAll() error {
	return spotlightDeleteAll()
}

// OnOpen registers a handler called when the user selects one of the
// application's items in Spotlight or chooses "Search in App". It returns a
// function that removes the handler. Handlers run on their own goroutine.
//
// macOS: delivered through application:continueUserActivity: for the
// CSSearchableItemActionType and CSQueryContinuationActionType activities,
// which the ActivityManager routes here; app.Activity.OnContinue handlers
// still see those activities too. Other platforms never call the handler.
func (s *SpotlightManager) OnOpen(handler SpotlightOpenHandler) func() {
	if handler == nil {
		return func() {}
	}
	s.lock.Lock()
	id := s.nextHandler
	s.nextHandler++
	s.openHandlers[id] = handler
	bridge := !s.bridged && s.app != nil && s.app.Activity != nil
	if bridge {
		s.bridged = true
	}
	s.lock.Unlock()
	if bridge {
		s.app.Activity.OnContinue(func(_ *Context, activity UserActivity) bool {
			return s.continueUserActivity(activity)
		})
	}
	return func() {
		s.lock.Lock()
		delete(s.openHandlers, id)
		s.lock.Unlock()
	}
}

func (s *SpotlightManager) hasOpenHandlers() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return len(s.openHandlers) > 0
}

// dispatchOpen runs every OnOpen handler on its own goroutine and reports
// whether any handler was registered.
func (s *SpotlightManager) dispatchOpen(id, query string) bool {
	s.lock.Lock()
	handlers := make([]SpotlightOpenHandler, 0, len(s.openHandlers))
	for _, handler := range s.openHandlers {
		handlers = append(handlers, handler)
	}
	s.lock.Unlock()
	for _, handler := range handlers {
		go func(handler SpotlightOpenHandler) {
			defer handlePanic()
			handler(newContext(), id, query)
		}(handler)
	}
	return len(handlers) > 0
}

// continueUserActivity resolves a user activity to an OnOpen call. It
// returns false when the activity is not a Spotlight continuation or no
// handler is registered, so other continuation handlers keep their turn.
func (s *SpotlightManager) continueUserActivity(activity UserActivity) bool {
	if s == nil {
		return false
	}
	switch activity.Type {
	case spotlightItemActionType:
		id, _ := activity.UserInfo[spotlightActivityIdentifierKey].(string)
		if id == "" {
			return false
		}
		return s.dispatchOpen(id, "")
	case spotlightQueryContinuationType:
		query, _ := activity.UserInfo[spotlightQueryStringKey].(string)
		return s.dispatchOpen("", query)
	default:
		return false
	}
}

// spotlightContinueActivity is continueUserActivity for callers that hold
// the activity type and its userInfo as a JSON object.
func (s *SpotlightManager) spotlightContinueActivity(activityType string, userInfoJSON string) bool {
	if s == nil {
		return false
	}
	var userInfo map[string]any
	if strings.TrimSpace(userInfoJSON) != "" {
		if err := json.Unmarshal([]byte(userInfoJSON), &userInfo); err != nil {
			if s.app != nil {
				s.app.handleError(fmt.Errorf("spotlight: bad user activity payload: %w", err))
			}
			return false
		}
	}
	return s.continueUserActivity(UserActivity{Type: activityType, UserInfo: userInfo})
}

// spotlightHandleContinueActivity is an alternative entry point for
// application:continueUserActivity:restorationHandler: that takes the
// activity's activityType and its userInfo dictionary serialised as a JSON
// object, returning true when a Spotlight OnOpen handler consumed it. The
// normal path needs no call at all: OnOpen registers itself with
// app.Activity.OnContinue, which the delegate already drives. On macOS the C
// helper wailsSpotlightHandleUserActivity(NSUserActivity *) declared in
// spotlight_manager_darwin.h wraps this function for native callers.
func spotlightHandleContinueActivity(activityType string, userInfoJSON string) bool {
	if globalApplication == nil || globalApplication.Spotlight == nil {
		return false
	}
	return globalApplication.Spotlight.spotlightContinueActivity(activityType, userInfoJSON)
}
