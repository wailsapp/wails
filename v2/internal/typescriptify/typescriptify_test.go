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
	CustomPtr    *time.Time              `json:"customPtr" ts_type:"Date"`
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
		"    custom: Date;",     // ts_type must win over the built-in time handling
		"    customPtr?: Date;", // including pointer fields, where ts_type supplies the type
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

// recursiveSlice is a legal recursive slice type used to check depth detection.
type recursiveSlice []recursiveSlice

// customByte marshals as a JSON number, so a slice of it is a JSON array.
type customByte byte

func (b *customByte) MarshalJSON() ([]byte, error) {
	return []byte("0"), nil
}

// namedBytes is a named byte slice; it still marshals to base64.
type namedBytes []byte

// textByte marshals as JSON text, so a slice of it is a JSON array.
type textByte byte

func (b *textByte) MarshalText() ([]byte, error) {
	return []byte("byte"), nil
}

// CustomArrayFields covers ts_type overrides on slice and array fields and
// keeps the non-array override, map and ts_transform paths in the same build.
type CustomArrayFields struct {
	Ints           []int           `json:"ints" ts_type:"Date"`
	Nested         [][]int         `json:"nested" ts_type:"Date"`
	Array          [2]int          `json:"array" ts_type:"Date"`
	Times          []time.Time     `json:"times" ts_type:"Date"`
	PtrTimes       []*time.Time    `json:"ptrTimes" ts_type:"Date"`
	PlainTimes     []*time.Time    `json:"plainTimes"`
	CustomMap      map[string]int  `json:"customMap" ts_type:"CustomType"`
	Bytes          []byte          `json:"bytes" ts_type:"{[key:string]:string}"`
	Union          []int           `json:"union" ts_type:"Date | null"`
	PtrInts        *[]int          `json:"ptrInts" ts_type:"Date"`
	NestedPtrTimes []*[]*time.Time `json:"nestedPtrTimes" ts_type:"Date"`
	CustomBytes    []customByte    `json:"customBytes" ts_type:"Date"`
	TextBytes      []textByte      `json:"textBytes" ts_type:"Date"`
	NestedBytes    [][]byte        `json:"nestedBytes" ts_type:"string"`
	DeepBytes      [][][]byte      `json:"deepBytes" ts_type:"string"`
	ByteArrays     [2][]byte       `json:"byteArrays" ts_type:"string"`
	PtrNestedBytes []*[]byte       `json:"ptrNestedBytes" ts_type:"string"`
	NestedCustom   [][]customByte  `json:"nestedCustomBytes" ts_type:"string"`
	NestedText     [][]textByte    `json:"nestedTextBytes" ts_type:"string"`
	NamedBytes     []namedBytes    `json:"namedBytes" ts_type:"string"`
	ByteArray      [2]byte         `json:"byteArray" ts_type:"string"`
	DeepPtrBytes   []*[][]byte     `json:"deepPtrBytes" ts_type:"string"`
	Recursive      recursiveSlice  `json:"recursive" ts_type:"Date"`
	Transformed    time.Time       `json:"transformed" ts_type:"Date" ts_transform:"new Date(__VALUE__)"`
}

// generateCustomArrayFields returns the TypeScript generated for
// CustomArrayFields.
func generateCustomArrayFields(t *testing.T) string {
	t.Helper()
	ts := New()
	ts.Namespace = "typescriptify"
	ts.KnownStructs = slicer.String()
	ts.Add(CustomArrayFields{})
	code, err := ts.Convert(nil)
	if err != nil {
		t.Fatalf("Convert() error: %v", err)
	}
	return code
}

// TestCustomArrayFieldsKeepDimensions checks that a ts_type tag on a slice or
// array supplies the element type while the array dimensions are preserved.
func TestCustomArrayFieldsKeepDimensions(t *testing.T) {
	code := generateCustomArrayFields(t)

	wantFields := []string{
		"    ints: Date[];",
		"    nested: Date[][];",
		"    array: Date[];",
		"    times: Date[];",
		"    ptrTimes: Date[];",
		"    plainTimes: (string | null)[];",
		"    customMap: CustomType;",
		"    bytes: {[key:string]:string};",
		"    union: (Date | null)[];",
		"    ptrInts?: Date[];",
		"    nestedPtrTimes: Date[][];",
		"    customBytes: Date[];",
		"    textBytes: Date[];",
		"    nestedBytes: string[];",
		"    deepBytes: string[][];",
		"    byteArrays: string[];",
		"    ptrNestedBytes: string[];",
		"    nestedCustomBytes: string[][];",
		"    nestedTextBytes: string[][];",
		"    namedBytes: string[];",
		"    byteArray: string[];",
		"    deepPtrBytes: string[][];",
		"    recursive: Date[];",
		"    transformed: Date;",
	}
	for _, want := range wantFields {
		if !strings.Contains(code, want) {
			t.Errorf("generated code does not contain %q\n%s", want, code)
		}
	}

	notWantFields := []string{
		"    ints: Date;",
		"    nested: Date;",
		"    times: string[];",
		"    ptrTimes: (string | null)[];",
		"    bytes: {[key:string]:string}[];",
		"    nestedBytes: string[][];",
		"    deepBytes: string[][][];",
		"    byteArrays: string[][];",
		"    ptrNestedBytes: string[][];",
		"    nestedCustomBytes: string[];",
		"    nestedTextBytes: string[];",
		"    deepPtrBytes: string[][][];",
	}
	for _, notWant := range notWantFields {
		if strings.Contains(code, notWant) {
			t.Errorf("generated code unexpectedly contains %q\n%s", notWant, code)
		}
	}

	if !strings.Contains(code, `this.transformed = new Date(source["transformed"]);`) {
		t.Errorf("ts_transform initializer was not generated\n%s", code)
	}
	if strings.Contains(code, "convertValues") {
		t.Errorf("generated code contains convertValues for simple array/custom fields\n%s", code)
	}
}

// ManageTypeArray covers a field option supplied through WithFieldOpts.
type ManageTypeArray struct {
	Values []int `json:"values"`
}

// TestManageTypeArrayKeepsDimensions checks that a WithFieldOpts ts_type on a
// slice keeps the array dimensions.
func TestManageTypeArrayKeepsDimensions(t *testing.T) {
	ts := New()
	ts.Namespace = "typescriptify"
	ts.KnownStructs = slicer.String()
	ts.Add(NewStruct(ManageTypeArray{}).WithFieldOpts([]int{}, TypeOptions{TSType: "Date"}))
	code, err := ts.Convert(nil)
	if err != nil {
		t.Fatalf("Convert() error: %v", err)
	}
	if !strings.Contains(code, "    values: Date[];") {
		t.Errorf("generated code does not contain %q\n%s", "    values: Date[];", code)
	}
}
