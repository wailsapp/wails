package typescriptify

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/leaanthony/slicer"
)

// TimeFields covers the shapes a time.Time field can take in a bound struct.
type TimeFields struct {
	Time         time.Time               `json:"time"`
	Maybe        time.Time               `json:"maybe,omitempty"`
	Ptr          *time.Time              `json:"ptr"`
	Optional     *time.Time              `json:"optional,omitempty"`
	Slice        []time.Time             `json:"slice"`
	PtrSlice     []*time.Time            `json:"ptrSlice"`
	Nested       [][]time.Time           `json:"nested"`
	NestedPtr    [][]*time.Time          `json:"nestedPtr"`
	PtrNested    []*[]time.Time          `json:"ptrNested"`
	PtrNestedPtr []*[]*time.Time         `json:"ptrNestedPtr"`
	ArrayPtr     [2]*time.Time           `json:"arrayPtr"`
	Map          map[string]time.Time    `json:"map"`
	PtrMap       map[string]*time.Time   `json:"ptrMap"`
	MapSlice     map[string][]*time.Time `json:"mapSlice"`
	Custom       time.Time               `json:"custom" ts_type:"Date"`
}

// newTimeFieldsTS returns a generator configured for the TimeFields struct.
func newTimeFieldsTS() *TypeScriptify {
	ts := New()
	// Types in the same namespace are the ones the generator emits code for.
	ts.Namespace = "typescriptify"
	ts.KnownStructs = slicer.String()
	ts.Add(TimeFields{})
	return ts
}

// generateTimeFields returns the TypeScript generated for TimeFields.
func generateTimeFields(t *testing.T) string {
	t.Helper()
	code, err := newTimeFieldsTS().Convert(nil)
	if err != nil {
		t.Fatalf("Convert() error: %v", err)
	}
	return code
}

// TestTimeFieldsGeneratedAsString checks the generated type of every
// supported time.Time shape.
func TestTimeFieldsGeneratedAsString(t *testing.T) {
	code := generateTimeFields(t)

	// time.Time marshals to a JSON string (RFC 3339). A nil *time.Time at any
	// level marshals to null, so pointer shapes need an explicit null member.
	wantFields := []string{
		"    time: string;",
		"    maybe?: string;",
		"    ptr?: string | null;",
		"    optional?: string | null;",
		"    slice: string[];",
		"    ptrSlice: (string | null)[];",
		"    nested: string[][];",
		"    nestedPtr: (string | null)[][];",
		"    ptrNested: (string[] | null)[];",
		"    ptrNestedPtr: ((string | null)[] | null)[];",
		"    arrayPtr: (string | null)[];",
		"    map: Record<string, string>;",
		"    ptrMap: Record<string, string | null>;",
		"    mapSlice: Record<string, Array<string | null>>;",
		"    custom: Date;", // ts_type must win over the built-in time handling
	}
	for _, want := range wantFields {
		if !strings.Contains(code, want) {
			t.Errorf("generated code does not contain %q\n%s", want, code)
		}
	}

	for _, field := range []string{"time", "maybe", "ptr", "optional", "slice", "ptrSlice", "nested", "nestedPtr", "ptrNested", "ptrNestedPtr", "arrayPtr", "map", "ptrMap", "mapSlice"} {
		convert := `convertValues(source["` + field + `"]`
		if strings.Contains(code, convert) {
			t.Errorf("constructor calls %s for time field %q; a string needs no conversion\n%s", convert, field, code)
		}
	}

	if strings.Contains(code, "class Time {") {
		t.Errorf("generated code contains a synthetic Time class\n%s", code)
	}
}

// TestTimeFieldsDoNotReportNotFound checks that time.Time no longer triggers
// the generator's "Not found" diagnostic.
func TestTimeFieldsDoNotReportNotFound(t *testing.T) {
	// The "Not found: time.Time" diagnostic is printed with the builtin
	// println, which writes straight to the process stderr, so the check runs
	// in a child process where stderr can be captured portably.
	if os.Getenv("TYPESCRIPTIFY_TIME_CHILD") == "1" {
		if _, err := newTimeFieldsTS().Convert(nil); err != nil {
			t.Fatalf("Convert() error: %v", err)
		}
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestTimeFieldsDoNotReportNotFound$")
	cmd.Env = append(os.Environ(), "TYPESCRIPTIFY_TIME_CHILD=1")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("child test run failed: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	if strings.Contains(stderr.String(), "Not found: time.Time") {
		t.Errorf("conversion printed %q\nstderr:\n%s", "Not found: time.Time", stderr.String())
	}
}
