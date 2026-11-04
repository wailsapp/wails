package manifest

import (
	"testing"
	"unicode/utf8"
)

func FuzzManifestDecodeNeverPanics(f *testing.F) {
	f.Add([]byte(""))
	f.Add([]byte("version: 3"))
	f.Add([]byte("version: 3\nproject: {}"))
	f.Add([]byte("version: 3\nproject: &project {}"))
	f.Add([]byte("version: 3\nproject:\n  name: \u2603"))
	root := f.TempDir()
	f.Fuzz(func(t *testing.T, source []byte) {
		_, _ = decodeYAML(root, "fuzz.yaml", source, "")
	})
}

func FuzzManifestLiteralStringRoundTrip(f *testing.F) {
	f.Add("")
	f.Add("plain")
	f.Add("quote: \" and slash: \\")
	f.Add("line one\nline two")
	f.Add("snowman: ☃")
	f.Add("control: \x00\x01")
	f.Add("carriage\rreturn\tand tab")
	f.Add("template markers: ${value} and %{if true}")
	root := f.TempDir()
	f.Fuzz(func(t *testing.T, literal string) {
		if !utf8.ValidString(literal) {
			t.Skip()
		}
		document := NewDocument(Project{Name: "literal", ProductName: "Literal", Identifier: "com.example.literal", Version: "1.0.0", Description: literal})
		encoded, err := EncodeDocument(document)
		if err != nil {
			t.Fatal(err)
		}
		roundTripped, err := decodeYAML(root, "round-trip.yaml", encoded, "")
		if err != nil {
			t.Fatal(err)
		}
		if got := roundTripped.Config.Project.Description; got != literal {
			t.Fatalf("description = %q, want %q", got, literal)
		}
	})
}
