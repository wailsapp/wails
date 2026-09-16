package application

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

var testPNG = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")

func TestSearchableItemValidate(t *testing.T) {
	tests := []struct {
		name    string
		item    SearchableItem
		wantErr string
	}{
		{name: "ok", item: SearchableItem{ID: "1", Title: "One"}},
		{name: "with png", item: SearchableItem{ID: "1", Title: "One", ThumbnailPNG: testPNG}},
		{name: "missing id", item: SearchableItem{Title: "One"}, wantErr: "ID is empty"},
		{name: "blank id", item: SearchableItem{ID: "  ", Title: "One"}, wantErr: "ID is empty"},
		{name: "missing title", item: SearchableItem{ID: "1"}, wantErr: "no title"},
		{name: "bad thumbnail", item: SearchableItem{ID: "1", Title: "One", ThumbnailPNG: []byte("JFIF")}, wantErr: "not a PNG"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.item.validate()
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("error = %v, want it to contain %q", err, test.wantErr)
			}
		})
	}

	if err := validateSearchableItems(nil); err == nil {
		t.Fatal("expected an error for no items")
	}
	dup := []SearchableItem{{ID: "1", Title: "A"}, {ID: "1", Title: "B"}}
	if err := validateSearchableItems(dup); err == nil || !strings.Contains(err.Error(), "more than once") {
		t.Fatalf("expected a duplicate ID error, got %v", err)
	}
	if err := validateSearchableItems([]SearchableItem{{ID: "1", Title: "A"}, {ID: "2", Title: "B"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSpotlightItemsJSON(t *testing.T) {
	expires := time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC)
	payload, err := spotlightItemsJSON([]SearchableItem{
		{ID: "a", Domain: "notes", Title: "Alpha", Description: "First", Keywords: []string{"x", "y"}, URL: "app://notes/a", ThumbnailPNG: testPNG, ExpiresAt: expires},
		{ID: "b", Title: "Beta", ContentType: "public.plain-text"},
	})
	if err != nil {
		t.Fatal(err)
	}
	var items []map[string]any
	if err := json.Unmarshal([]byte(payload), &items); err != nil {
		t.Fatalf("payload is not JSON: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	a, b := items[0], items[1]
	if a["id"] != "a" || a["domain"] != "notes" || a["title"] != "Alpha" || a["description"] != "First" || a["url"] != "app://notes/a" {
		t.Fatalf("unexpected first item: %v", a)
	}
	if a["contentType"] != spotlightDefaultContentType {
		t.Fatalf("expected default content type, got %v", a["contentType"])
	}
	if a["expires"] != float64(expires.Unix()) {
		t.Fatalf("unexpected expiry: %v", a["expires"])
	}
	if a["thumbnail"] == nil || a["thumbnail"] == "" {
		t.Fatal("thumbnail not encoded")
	}
	if b["contentType"] != "public.plain-text" {
		t.Fatalf("content type not kept: %v", b["contentType"])
	}
	for _, key := range []string{"domain", "thumbnail", "url", "expires", "keywords"} {
		if _, present := b[key]; present {
			t.Fatalf("empty field %q should be omitted", key)
		}
	}
}

func TestIsPNG(t *testing.T) {
	if !isPNG(testPNG) {
		t.Fatal("PNG signature not recognised")
	}
	if isPNG([]byte("\x89PNG")) || isPNG(nil) || isPNG([]byte("GIF89a..")) {
		t.Fatal("non-PNG data accepted")
	}
}

func TestSpotlightContinueActivity(t *testing.T) {
	s := newSpotlightManager(nil)

	if s.spotlightContinueActivity(spotlightItemActionType, `{"kCSSearchableItemActivityIdentifier":"note-1"}`) {
		t.Fatal("expected false with no handlers")
	}

	type call struct{ id, query string }
	calls := make(chan call, 4)
	unsubscribe := s.OnOpen(func(ctx *Context, id string, query string) {
		if ctx == nil {
			t.Error("context is nil")
		}
		calls <- call{id: id, query: query}
	})
	wait := func() call {
		select {
		case c := <-calls:
			return c
		case <-time.After(2 * time.Second):
			t.Fatal("handler was not called")
			return call{}
		}
	}

	if !s.spotlightContinueActivity(spotlightItemActionType, `{"kCSSearchableItemActivityIdentifier":"note-1"}`) {
		t.Fatal("item activity not handled")
	}
	if got := wait(); got != (call{id: "note-1"}) {
		t.Fatalf("unexpected call %+v", got)
	}

	if !s.spotlightContinueActivity(spotlightQueryContinuationType, `{"kCSSearchQueryString":"alpha"}`) {
		t.Fatal("query activity not handled")
	}
	if got := wait(); got != (call{query: "alpha"}) {
		t.Fatalf("unexpected call %+v", got)
	}

	if s.spotlightContinueActivity(spotlightItemActionType, `{}`) {
		t.Fatal("item activity without an identifier should not be handled")
	}
	if s.spotlightContinueActivity("NSUserActivityTypeBrowsingWeb", `{}`) {
		t.Fatal("unrelated activity should not be handled")
	}
	if s.spotlightContinueActivity(spotlightItemActionType, `not json`) {
		t.Fatal("malformed payload should not be handled")
	}

	unsubscribe()
	if s.spotlightContinueActivity(spotlightItemActionType, `{"kCSSearchableItemActivityIdentifier":"note-1"}`) {
		t.Fatal("expected false after unsubscribe")
	}
	if s.OnOpen(nil) == nil {
		t.Fatal("OnOpen(nil) must return a no-op unsubscribe")
	}
	if spotlightHandleContinueActivity(spotlightItemActionType, `{}`) && globalApplication == nil {
		t.Fatal("package entry point must be false without an application")
	}
}

func TestSpotlightBridgesActivityContinuation(t *testing.T) {
	app := &App{}
	app.Activity = newActivityManager(app)
	app.Spotlight = newSpotlightManager(app)

	// Without an OnOpen handler the bridge is not installed, so the
	// ActivityManager reports the activity as unhandled.
	if app.Activity.continueActivity(UserActivity{Type: spotlightItemActionType, UserInfo: map[string]any{spotlightActivityIdentifierKey: "note-1"}}) {
		t.Fatal("expected false before OnOpen")
	}

	opened := make(chan string, 2)
	app.Spotlight.OnOpen(func(_ *Context, id, query string) { opened <- id + "|" + query })
	app.Spotlight.OnOpen(func(_ *Context, id, query string) {})
	app.Activity.handlerLock.RLock()
	bridges := len(app.Activity.continueFns)
	app.Activity.handlerLock.RUnlock()
	if bridges != 1 {
		t.Fatalf("expected exactly one continuation bridge, got %d", bridges)
	}

	if !app.Activity.continueActivity(UserActivity{Type: spotlightItemActionType, UserInfo: map[string]any{spotlightActivityIdentifierKey: "note-1"}}) {
		t.Fatal("item activity not handled through the ActivityManager")
	}
	select {
	case got := <-opened:
		if got != "note-1|" {
			t.Fatalf("unexpected call %q", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("handler was not called")
	}
	if app.Activity.continueActivity(UserActivity{Type: UserActivityTypeBrowsingWeb, WebpageURL: "https://example.com"}) {
		t.Fatal("browsing activity must not be claimed by Spotlight")
	}
}
