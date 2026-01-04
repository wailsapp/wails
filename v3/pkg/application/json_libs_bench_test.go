//go:build bench

// Disabled: goccy/go-json causes Windows panics. See PR #4859.

package application_test

/*
import (
	"encoding/json"
	"testing"

	"github.com/bytedance/sonic"
	gojson "github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
)

// Test structures matching real Wails binding patterns

type SimpleBindingArg struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type ComplexBindingArg struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
	Nested   *NestedBindingArg      `json:"nested,omitempty"`
}

type NestedBindingArg struct {
	Value   float64 `json:"value"`
	Enabled bool    `json:"enabled"`
}

// Test data simulating frontend calls
var (
	simpleJSON = []byte(`{"name":"test","value":42}`)

	complexJSON = []byte(`{"id":12345,"name":"Test Complex Data","tags":["tag1","tag2","tag3","tag4","tag5"],"metadata":{"key1":"value1","key2":42,"key3":true},"nested":{"value":3.14159,"enabled":true}}`)

	stringJSON = []byte(`"hello world this is a test string"`)

	multiArgsJSON = [][]byte{
		[]byte(`"arg1"`),
		[]byte(`42`),
		[]byte(`true`),
		[]byte(`{"key":"value"}`),
	}
)

// Configure jsoniter for maximum compatibility
var jsoniterStd = jsoniter.ConfigCompatibleWithStandardLibrary

// ============================================================================
// UNMARSHAL BENCHMARKS - This is the HOT PATH (bindings.go:289)
// ============================================================================

// --- Simple struct unmarshal ---

func BenchmarkUnmarshal_Simple_StdLib(b *testing.B) {
	for b.Loop() {
		var arg SimpleBindingArg
		_ = json.Unmarshal(simpleJSON, &arg)
	}
}

func BenchmarkUnmarshal_Simple_GoJSON(b *testing.B) {
	for b.Loop() {
		var arg SimpleBindingArg
		_ = gojson.Unmarshal(simpleJSON, &arg)
	}
}

func BenchmarkUnmarshal_Simple_JSONIter(b *testing.B) {
	for b.Loop() {
		var arg SimpleBindingArg
		_ = jsoniterStd.Unmarshal(simpleJSON, &arg)
	}
}

func BenchmarkUnmarshal_Simple_Sonic(b *testing.B) {
	for b.Loop() {
		var arg SimpleBindingArg
		_ = sonic.Unmarshal(simpleJSON, &arg)
	}
}

// --- Complex struct unmarshal ---

func BenchmarkUnmarshal_Complex_StdLib(b *testing.B) {
	for b.Loop() {
		var arg ComplexBindingArg
		_ = json.Unmarshal(complexJSON, &arg)
	}
}

func BenchmarkUnmarshal_Complex_GoJSON(b *testing.B) {
	for b.Loop() {
		var arg ComplexBindingArg
		_ = gojson.Unmarshal(complexJSON, &arg)
	}
}

func BenchmarkUnmarshal_Complex_JSONIter(b *testing.B) {
	for b.Loop() {
		var arg ComplexBindingArg
		_ = jsoniterStd.Unmarshal(complexJSON, &arg)
	}
}

func BenchmarkUnmarshal_Complex_Sonic(b *testing.B) {
	for b.Loop() {
		var arg ComplexBindingArg
		_ = sonic.Unmarshal(complexJSON, &arg)
	}
}

// --- String unmarshal (most common single arg) ---

func BenchmarkUnmarshal_String_StdLib(b *testing.B) {
	for b.Loop() {
		var arg string
		_ = json.Unmarshal(stringJSON, &arg)
	}
}

func BenchmarkUnmarshal_String_GoJSON(b *testing.B) {
	for b.Loop() {
		var arg string
		_ = gojson.Unmarshal(stringJSON, &arg)
	}
}

func BenchmarkUnmarshal_String_JSONIter(b *testing.B) {
	for b.Loop() {
		var arg string
		_ = jsoniterStd.Unmarshal(stringJSON, &arg)
	}
}

func BenchmarkUnmarshal_String_Sonic(b *testing.B) {
	for b.Loop() {
		var arg string
		_ = sonic.Unmarshal(stringJSON, &arg)
	}
}

// --- Interface{} unmarshal (dynamic typing) ---

func BenchmarkUnmarshal_Interface_StdLib(b *testing.B) {
	for b.Loop() {
		var arg interface{}
		_ = json.Unmarshal(complexJSON, &arg)
	}
}

func BenchmarkUnmarshal_Interface_GoJSON(b *testing.B) {
	for b.Loop() {
		var arg interface{}
		_ = gojson.Unmarshal(complexJSON, &arg)
	}
}

func BenchmarkUnmarshal_Interface_JSONIter(b *testing.B) {
	for b.Loop() {
		var arg interface{}
		_ = jsoniterStd.Unmarshal(complexJSON, &arg)
	}
}

func BenchmarkUnmarshal_Interface_Sonic(b *testing.B) {
	for b.Loop() {
		var arg interface{}
		_ = sonic.Unmarshal(complexJSON, &arg)
	}
}

// --- Multi-arg unmarshal (simulating typical method call) ---

func BenchmarkUnmarshal_MultiArgs_StdLib(b *testing.B) {
	for b.Loop() {
		var s string
		var i int
		var bl bool
		var m map[string]string
		_ = json.Unmarshal(multiArgsJSON[0], &s)
		_ = json.Unmarshal(multiArgsJSON[1], &i)
		_ = json.Unmarshal(multiArgsJSON[2], &bl)
		_ = json.Unmarshal(multiArgsJSON[3], &m)
	}
}

func BenchmarkUnmarshal_MultiArgs_GoJSON(b *testing.B) {
	for b.Loop() {
		var s string
		var i int
		var bl bool
		var m map[string]string
		_ = gojson.Unmarshal(multiArgsJSON[0], &s)
		_ = gojson.Unmarshal(multiArgsJSON[1], &i)
		_ = gojson.Unmarshal(multiArgsJSON[2], &bl)
		_ = gojson.Unmarshal(multiArgsJSON[3], &m)
	}
}

func BenchmarkUnmarshal_MultiArgs_JSONIter(b *testing.B) {
	for b.Loop() {
		var s string
		var i int
		var bl bool
		var m map[string]string
		_ = jsoniterStd.Unmarshal(multiArgsJSON[0], &s)
		_ = jsoniterStd.Unmarshal(multiArgsJSON[1], &i)
		_ = jsoniterStd.Unmarshal(multiArgsJSON[2], &bl)
		_ = jsoniterStd.Unmarshal(multiArgsJSON[3], &m)
	}
}

func BenchmarkUnmarshal_MultiArgs_Sonic(b *testing.B) {
	for b.Loop() {
		var s string
		var i int
		var bl bool
		var m map[string]string
		_ = sonic.Unmarshal(multiArgsJSON[0], &s)
		_ = sonic.Unmarshal(multiArgsJSON[1], &i)
		_ = sonic.Unmarshal(multiArgsJSON[2], &bl)
		_ = sonic.Unmarshal(multiArgsJSON[3], &m)
	}
}

// ============================================================================
// MARSHAL BENCHMARKS - Result serialization
// ============================================================================

type BindingResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var simpleResult = BindingResult{
	Success: true,
	Data:    "hello world",
}

var complexResult = BindingResult{
	Success: true,
	Data: ComplexBindingArg{
		ID:   12345,
		Name: "Result Data",
		Tags: []string{"a", "b", "c"},
		Metadata: map[string]interface{}{
			"processed": true,
			"count":     100,
		},
		Nested: &NestedBindingArg{Value: 2.718, Enabled: true},
	},
}

// --- Simple result marshal ---

func BenchmarkMarshal_Simple_StdLib(b *testing.B) {
	for b.Loop() {
		_, _ = json.Marshal(simpleResult)
	}
}

func BenchmarkMarshal_Simple_GoJSON(b *testing.B) {
	for b.Loop() {
		_, _ = gojson.Marshal(simpleResult)
	}
}

func BenchmarkMarshal_Simple_JSONIter(b *testing.B) {
	for b.Loop() {
		_, _ = jsoniterStd.Marshal(simpleResult)
	}
}

func BenchmarkMarshal_Simple_Sonic(b *testing.B) {
	for b.Loop() {
		_, _ = sonic.Marshal(simpleResult)
	}
}

// --- Complex result marshal ---

func BenchmarkMarshal_Complex_StdLib(b *testing.B) {
	for b.Loop() {
		_, _ = json.Marshal(complexResult)
	}
}

func BenchmarkMarshal_Complex_GoJSON(b *testing.B) {
	for b.Loop() {
		_, _ = gojson.Marshal(complexResult)
	}
}

func BenchmarkMarshal_Complex_JSONIter(b *testing.B) {
	for b.Loop() {
		_, _ = jsoniterStd.Marshal(complexResult)
	}
}

func BenchmarkMarshal_Complex_Sonic(b *testing.B) {
	for b.Loop() {
		_, _ = sonic.Marshal(complexResult)
	}
}
*/
