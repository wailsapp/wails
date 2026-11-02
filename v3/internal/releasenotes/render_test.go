package releasenotes

import (
	"encoding/json"
	"os"
	"testing"
)

func TestRenderNoColour(t *testing.T) {
	data, err := os.ReadFile("testdata/no_colour.json")
	if err != nil {
		t.Fatal(err)
	}
	var want []string
	if err := json.Unmarshal(data, &want); err != nil {
		t.Fatal(err)
	}
	notes := []string{
		"# Release Notes for v3.0.0\n## Changes\n\n- Faster startup\n- A **bold** improvement\n\nCode: `hello()`\n",
		"# Release Notes for v3.0.0\n",
	}
	if len(want) != len(notes) {
		t.Fatalf("fixture count = %d, want %d", len(want), len(notes))
	}
	for i, note := range notes {
		if got := Render(note, true); got != want[i] {
			t.Errorf("case %d: got %q, want %q", i, got, want[i])
		}
	}
}
