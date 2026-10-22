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
	Time     time.Time            `json:"time"`
	Ptr      *time.Time           `json:"ptr"`
	Optional *time.Time           `json:"optional,omitempty"`
	Slice    []time.Time          `json:"slice"`
	PtrSlice []*time.Time         `json:"ptrSlice"`
	Nested   [][]time.Time        `json:"nested"`
	Map      map[string]time.Time `json:"map"`
	Custom   time.Time            `json:"custom" ts_type:"Date"`
}

func newTimeFieldsTS() *TypeScriptify {
	ts := New()
	// Types in the same namespace are the ones the generator emits code for.
	ts.Namespace = "typescriptify"
	ts.KnownStructs = slicer.String()
	ts.Add(TimeFields{})
	return ts
}

func generateTimeFields(t *testing.T) string {
	t.Helper()
	code, err := newTimeFieldsTS().Convert(nil)
	if err != nil {
		t.Fatalf("Convert() error: %v", err)
	}
	return code
}

func TestTimeFieldsGeneratedAsString(t *testing.T) {
	code := generateTimeFields(t)

	// time.Time marshals to a JSON string (RFC 3339), so the generated type
	// must be string rather than any or a synthetic Time class.
	wantFields := []string{
		"    time: string;",
		"    ptr?: string;",
		"    optional?: string;",
		"    slice: string[];",
		"    ptrSlice: string[];",
		"    nested: string[][];",
		"    map: Record<string, string>;",
		"    custom: Date;", // ts_type must win over the built-in time handling
	}
	for _, want := range wantFields {
		if !strings.Contains(code, want) {
			t.Errorf("generated code does not contain %q\n%s", want, code)
		}
	}

	for _, field := range []string{"time", "ptr", "optional", "slice", "ptrSlice", "nested", "map"} {
		convert := `convertValues(source["` + field + `"]`
		if strings.Contains(code, convert) {
			t.Errorf("constructor calls %s for time field %q; a string needs no conversion\n%s", convert, field, code)
		}
	}

	if strings.Contains(code, "class Time {") {
		t.Errorf("generated code contains a synthetic Time class\n%s", code)
	}
}

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
