package application

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

type fakeActivity struct {
	published   map[uint64]UserActivity
	invalidated []uint64
	publishErr  error
}

func newFakeActivity() *fakeActivity {
	return &fakeActivity{published: map[uint64]UserActivity{}}
}

func (f *fakeActivity) publish(id uint64, activity UserActivity) error {
	if f.publishErr != nil {
		return f.publishErr
	}
	f.published[id] = activity
	return nil
}

func (f *fakeActivity) update(id uint64, userInfo map[string]any) error {
	activity, ok := f.published[id]
	if !ok {
		return ErrActivityInvalidated
	}
	activity.UserInfo = userInfo
	f.published[id] = activity
	return nil
}

func (f *fakeActivity) invalidate(id uint64) {
	delete(f.published, id)
	f.invalidated = append(f.invalidated, id)
}

func newTestActivityManager(impl activityImpl) *ActivityManager {
	return &ActivityManager{impl: impl, published: map[uint64]*PublishedActivity{}}
}

func TestUserActivityValidation(t *testing.T) {
	valid := UserActivity{
		Type:               "com.example.editing",
		Title:              "Editing notes",
		UserInfo:           map[string]any{"id": 42},
		WebpageURL:         "https://example.com/notes/42",
		EligibleForHandoff: true,
	}
	if err := valid.validate(); err != nil {
		t.Fatalf("valid activity rejected: %v", err)
	}
	cases := []struct {
		name   string
		mutate func(*UserActivity)
		want   string
	}{
		{"empty type", func(a *UserActivity) { a.Type = " " }, "Type is required"},
		{"type with space", func(a *UserActivity) { a.Type = "com.example. editing" }, "whitespace"},
		{"relative url", func(a *UserActivity) { a.WebpageURL = "/notes/42" }, "WebpageURL"},
		{"ftp url", func(a *UserActivity) { a.WebpageURL = "ftp://example.com/x" }, "WebpageURL"},
		{"browsing without url", func(a *UserActivity) { a.Type = UserActivityTypeBrowsingWeb; a.WebpageURL = "" }, "browsing activity"},
		{"handoff without title or url", func(a *UserActivity) { a.Title = ""; a.WebpageURL = "" }, "Handoff activity"},
		{"unserialisable userinfo", func(a *UserActivity) { a.UserInfo = map[string]any{"fn": func() {}} }, "UserInfo"},
	}
	for _, tc := range cases {
		activity := valid
		tc.mutate(&activity)
		err := activity.validate()
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%s: got %v, want error mentioning %q", tc.name, err, tc.want)
		}
	}
	// A Spotlight-only activity needs neither title nor URL.
	if err := (UserActivity{Type: "com.example.view", EligibleForSearch: true}).validate(); err != nil {
		t.Errorf("search-only activity rejected: %v", err)
	}
}

func TestUserActivityUserInfoJSONRoundTrip(t *testing.T) {
	original := UserActivity{
		Type:                  "com.example.editing",
		Title:                 "Editing",
		UserInfo:              map[string]any{"id": float64(42), "tags": []any{"a", "b"}, "nested": map[string]any{"ok": true}, "none": nil},
		WebpageURL:            "https://example.com/42",
		EligibleForHandoff:    true,
		EligibleForSearch:     true,
		EligibleForPrediction: true,
		Keywords:              []string{"notes", "editing"},
	}
	encoded, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"type"`, `"title"`, `"userInfo"`, `"webpageURL"`, `"eligibleForHandoff"`, `"eligibleForSearch"`, `"eligibleForPrediction"`, `"keywords"`} {
		if !strings.Contains(string(encoded), key) {
			t.Errorf("encoded activity lacks %s: %s", key, encoded)
		}
	}
	var decoded UserActivity
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded, original) {
		t.Fatalf("round trip changed the activity:\n%#v\nwant\n%#v", decoded, original)
	}
}

func TestActivityPublishUpdateInvalidate(t *testing.T) {
	impl := newFakeActivity()
	m := newTestActivityManager(impl)

	if _, err := m.Publish(UserActivity{}); err == nil {
		t.Fatal("invalid activity published")
	}
	published, err := m.Publish(UserActivity{Type: "com.example.editing", Title: "Editing"})
	if err != nil {
		t.Fatalf("Publish: %v", err)
	}
	if got := published.Activity(); got.UserInfo == nil || got.Keywords == nil {
		t.Fatalf("Publish should normalise nil UserInfo and Keywords: %+v", got)
	}
	if len(impl.published) != 1 {
		t.Fatalf("impl published %v", impl.published)
	}

	if err := published.Update(map[string]any{"fn": func() {}}); err == nil {
		t.Fatal("unserialisable Update accepted")
	}
	if err := published.Update(map[string]any{"cursor": 7}); err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got := published.Activity().UserInfo["cursor"]; got != 7 {
		t.Fatalf("Update did not keep UserInfo: %v", published.Activity().UserInfo)
	}
	if got := impl.published[published.id].UserInfo["cursor"]; got != 7 {
		t.Fatalf("Update did not reach the platform: %v", impl.published[published.id].UserInfo)
	}

	published.Invalidate()
	published.Invalidate()
	if len(impl.invalidated) != 1 || len(impl.published) != 0 {
		t.Fatalf("Invalidate: invalidated=%v published=%v", impl.invalidated, impl.published)
	}
	if err := published.Update(map[string]any{}); !errors.Is(err, ErrActivityInvalidated) {
		t.Fatalf("Update after Invalidate: got %v", err)
	}
	m.lock.RLock()
	remaining := len(m.published)
	m.lock.RUnlock()
	if remaining != 0 {
		t.Fatalf("manager still tracks %d activities", remaining)
	}

	impl.publishErr = errors.New("boom")
	if _, err := m.Publish(UserActivity{Type: "com.example.editing", Title: "Editing"}); err == nil {
		t.Fatal("platform error swallowed")
	}

	unsupported := newTestActivityManager(nil)
	if _, err := unsupported.Publish(UserActivity{Type: "com.example.editing", Title: "Editing"}); !errors.Is(err, ErrActivityUnsupported) {
		t.Fatalf("unsupported platform: got %v", err)
	}
}

func TestActivityHandlers(t *testing.T) {
	m := newTestActivityManager(newFakeActivity())

	// Without handlers only browsing activities are accepted up front.
	if m.willContinue("com.example.editing") {
		t.Fatal("accepted an activity nobody handles")
	}
	if !m.willContinue(UserActivityTypeBrowsingWeb) {
		t.Fatal("universal links should be accepted by default")
	}

	var continued []UserActivity
	m.OnContinue(func(ctx *Context, activity UserActivity) bool {
		if ctx == nil {
			t.Error("nil context")
		}
		continued = append(continued, activity)
		return activity.Type == "com.example.editing"
	})
	if !m.willContinue("com.example.editing") {
		t.Fatal("an OnContinue handler should make willContinue accept")
	}
	m.OnWillContinue(func(activityType string) bool { return activityType == "com.example.special" })
	if m.willContinue("com.example.editing") || !m.willContinue("com.example.special") {
		t.Fatal("OnWillContinue handlers should decide once registered")
	}

	if !m.continueActivity(UserActivity{Type: "com.example.editing"}) {
		t.Fatal("handled activity reported as unhandled")
	}
	if m.continueActivity(UserActivity{Type: "com.example.other"}) {
		t.Fatal("unhandled activity reported as handled")
	}
	if len(continued) != 2 || continued[0].UserInfo == nil || continued[0].Keywords == nil {
		t.Fatalf("continue handlers saw %+v", continued)
	}

	var failedType string
	var failedErr error
	m.OnFailed(func(activityType string, err error) { failedType, failedErr = activityType, err })
	m.failed("com.example.editing", errors.New("transfer failed"))
	if failedType != "com.example.editing" || failedErr == nil || failedErr.Error() != "transfer failed" {
		t.Fatalf("OnFailed saw %q %v", failedType, failedErr)
	}

	var updated UserActivity
	m.OnUpdated(func(activity UserActivity) { updated = activity })
	m.updated(UserActivity{Type: "com.example.editing", Title: "Saved"})
	if updated.Title != "Saved" {
		t.Fatalf("OnUpdated saw %+v", updated)
	}

	// nil handlers are ignored rather than stored.
	m.OnContinue(nil)
	m.OnWillContinue(nil)
	m.OnFailed(nil)
	m.OnUpdated(nil)
	if len(m.continueFns) != 1 || len(m.willContinueFns) != 1 || len(m.failedFns) != 1 || len(m.updatedFns) != 1 {
		t.Fatal("nil handlers were stored")
	}
}
