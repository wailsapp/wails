package application

import (
	"reflect"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/events"
)

func TestParsePosixLocale(t *testing.T) {
	cases := []struct {
		in   string
		want LocaleInfo
		ok   bool
	}{
		{"en_GB.UTF-8", LocaleInfo{Identifier: "en_GB", Language: "en", Region: "GB", Preferred: []string{"en-GB"}}, true},
		{"pt_BR", LocaleInfo{Identifier: "pt_BR", Language: "pt", Region: "BR", Preferred: []string{"pt-BR"}}, true},
		{"de", LocaleInfo{Identifier: "de", Language: "de", Preferred: []string{"de"}}, true},
		{"ca_ES@valencia", LocaleInfo{Identifier: "ca_ES", Language: "ca", Region: "ES", Preferred: []string{"ca-ES"}}, true},
		{"C", LocaleInfo{}, false},
		{"POSIX", LocaleInfo{}, false},
		{"C.UTF-8", LocaleInfo{}, false},
		{"", LocaleInfo{}, false},
	}
	for _, tc := range cases {
		got, ok := parsePosixLocale(tc.in)
		if ok != tc.ok || !reflect.DeepEqual(got, tc.want) {
			t.Errorf("parsePosixLocale(%q) = %+v, %v; want %+v, %v", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}

func TestLocaleFromEnvironmentPrecedence(t *testing.T) {
	env := map[string]string{"LANG": "fr_FR.UTF-8", "LC_MESSAGES": "es_ES", "LC_ALL": "C"}
	lookup := func(name string) string { return env[name] }
	// LC_ALL=C is skipped, LC_MESSAGES wins over LANG.
	if got := localeFromEnvironment(lookup); got.Identifier != "es_ES" {
		t.Fatalf("Identifier = %q, want es_ES", got.Identifier)
	}
	delete(env, "LC_MESSAGES")
	if got := localeFromEnvironment(lookup); got.Identifier != "fr_FR" {
		t.Fatalf("Identifier = %q, want fr_FR", got.Identifier)
	}
	if got := localeFromEnvironment(func(string) string { return "" }); got.Identifier != "" {
		t.Fatalf("Identifier with empty env = %q, want empty", got.Identifier)
	}
}

func TestEnvironmentEventsRegistered(t *testing.T) {
	cases := map[events.ApplicationEventType]string{
		events.Mac.ApplicationDidChangeAccessibilitySettings: "mac:ApplicationDidChangeAccessibilitySettings",
		events.Mac.ApplicationDidChangeKeyboardLayout:        "mac:ApplicationDidChangeKeyboardLayout",
		events.Mac.ApplicationDidChangeLocale:                "mac:ApplicationDidChangeLocale",
		events.Common.AccessibilitySettingsChanged:           "common:AccessibilitySettingsChanged",
	}
	for id, name := range cases {
		if id == 0 {
			t.Errorf("%s has no event id", name)
		}
		if got := events.JSEvent(uint(id)); got != name {
			t.Errorf("JSEvent(%d) = %q, want %q", id, got, name)
		}
		if !events.IsKnownEvent(name) {
			t.Errorf("%s is not a known event", name)
		}
	}
}
