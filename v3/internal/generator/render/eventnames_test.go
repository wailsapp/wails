package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/generator/collect"
)

// TestEventNamesEscapeComments guards comment terminators in discoverable names.
func TestEventNamesEscapeComments(t *testing.T) {
	for _, ts := range []bool{false, true} {
		renderer := NewRenderer(nil, &flags.GenerateBindingsOptions{TS: ts})
		events := &collect.EventMap{
			Imports: &collect.ImportMap{Self: "example"},
			Defs:    []*collect.EventInfo{{Name: "before*/after"}},
		}
		var output bytes.Buffer
		if err := renderer.EventNames(&output, events); err != nil {
			t.Fatal(err)
		}
		got := output.String()
		if !strings.Contains(got, `Name of the custom event "before* /after".`) {
			t.Errorf("TS=%v: unsafe comment:\n%s", ts, got)
		}
		if !strings.Contains(got, `export const BeforeAfter = "before*/after";`) {
			t.Errorf("TS=%v: event value changed:\n%s", ts, got)
		}
	}
}
