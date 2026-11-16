package voidpackage

import (
	"bytes"
	"context"
	"slices"
	"testing"
)

func TestDependencyPatternsPreserveVersionConstraints(t *testing.T) {
	input := []string{"gtk4", "libwebkitgtk60", "lib-example", "gtk4>=4.0", "gtk4<5", "gtk4<=4.20", "gtk4>4", "foo-bar-1.0_2"}
	want := []string{"gtk4>=0", "libwebkitgtk60>=0", "lib-example>=0", "gtk4>=4.0", "gtk4<5", "gtk4<=4.20", "gtk4>4", "foo-bar-1.0_2"}
	original := slices.Clone(input)
	got, err := normalizePatterns(input)
	if err != nil || !slices.Equal(got, want) {
		t.Fatalf("patterns: %v, %v; want %v", got, err, want)
	}
	if !slices.Equal(input, original) {
		t.Fatal("normalization changed caller-owned dependency declarations")
	}
}

func TestMalformedDependencyCannotWriteAnArchive(t *testing.T) {
	for _, value := range []string{"", "gtk4>=", "gtk4<", "gtk4=1", "gtk4>>1", ">=1", "gtk4 >=4", "gtk4\n", "gtk4\x00"} {
		t.Run(value, func(t *testing.T) {
			options := fixture(t)
			options.Depends = []string{value}
			var output bytes.Buffer
			if err := Write(context.Background(), &output, options); err == nil || output.Len() != 0 {
				t.Fatalf("malformed pattern wrote %d archive bytes: %v", output.Len(), err)
			}
		})
	}
}

func TestConflictAndReplacementPatternsUseDependencySyntax(t *testing.T) {
	metadata := Metadata{Depends: []string{"gtk4"}, Conflicts: []string{"old-app"}, Replaces: []string{"old-app-1.0_1"}, Provides: []string{"virtual-app-1.0_1"}}
	got, err := normalizeDependencies(metadata)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got.Depends, []string{"gtk4>=0"}) || !slices.Equal(got.Conflicts, []string{"old-app>=0"}) || !slices.Equal(got.Replaces, metadata.Replaces) || !slices.Equal(got.Provides, metadata.Provides) {
		t.Fatalf("unexpected package relationships: %#v", got)
	}
	for _, invalid := range []Metadata{{Conflicts: []string{"old-app>="}}, {Replaces: []string{"old-app<"}}} {
		if _, err := normalizeDependencies(invalid); err == nil {
			t.Fatalf("accepted malformed package relationships: %#v", invalid)
		}
	}
}
