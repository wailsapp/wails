package main

import (
	"strings"
	"testing"
)

// The runtime classes and selectors here exist on every macOS, so resolve must
// find them, and must report the names that do not exist.
func TestResolveDarwinClassesAndSelectors(t *testing.T) {
	r := &report{
		Classes: []*entry{
			{Name: "NSObject", Where: []string{"a.go:1"}},
			{Name: "Mine", Ours: true, Where: []string{"a.go:2"}},
			{Name: "NoSuchClassPuregocheck", Where: []string{"a.go:3"}},
			{Name: "NoSuchOptionalClassPuregocheck", Optional: true, Where: []string{"a.go:4"}},
		},
		Selectors: []*entry{
			{Name: "description", Where: []string{"a.go:5"}},
			{Name: "noSuchSelectorPuregocheck:", Where: []string{"a.go:6"}},
			{Name: "noSuchOptionalPuregocheck:", Optional: true, Where: []string{"a.go:7"}},
			{Name: "ownSelector:", Ours: true, Where: []string{"a.go:8"}},
		},
	}
	if err := resolve(r); err != nil {
		t.Fatal(err)
	}
	if len(r.Loaded) == 0 || !strings.HasPrefix(r.Loaded[0], "macOS ") {
		t.Errorf("Loaded = %v, want the macOS version first", r.Loaded)
	}
	missing := strings.Join(r.Missing, "|")
	for _, want := range []string{"class NoSuchClassPuregocheck", "selector noSuchSelectorPuregocheck:"} {
		if !strings.Contains(missing, want) {
			t.Errorf("Missing = %v, want %q", r.Missing, want)
		}
	}
	if strings.Contains(missing, "NSObject") || strings.Contains(missing, "description") || strings.Contains(missing, "Mine") || strings.Contains(missing, "ownSelector") {
		t.Errorf("resolved or backend-owned name reported missing: %v", r.Missing)
	}
	absent := strings.Join(r.Absent, "|")
	for _, want := range []string{"class NoSuchOptionalClassPuregocheck", "selector noSuchOptionalPuregocheck:"} {
		if !strings.Contains(absent, want) {
			t.Errorf("Absent = %v, want %q", r.Absent, want)
		}
	}
}
