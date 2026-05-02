//go:build bench && goexperiment.jsonv2

package application_test

import (
	"encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"testing"
)

// Benchmark structures matching real Wails usage patterns

type SimpleArg struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type ComplexArg struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
	Nested   *NestedArg             `json:"nested,omitempty"`
}

type NestedArg struct {
	Value   float64 `json:"value"`
	Enabled bool    `json:"enabled"`
}

type CallResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Test data
var (
	simpleArgJSON  = []byte(`{"name":"test","value":42}`)
	complexArgJSON = []byte(`{
		"id": 12345,
		"name": "Test Complex Data",
		"tags": ["tag1", "tag2", "tag3", "tag4", "tag5"],
		"metadata": {"key1": "value1", "key2": 42, "key3": true},
		"nested": {"value": 3.14159, "enabled": true}
	}`)

	simpleResult = CallResult{
		Success: true,
		Data:    "hello world",
	}

	complexResult = CallResult{
		Success: true,
		Data: ComplexArg{
			ID:   12345,
			Name: "Result Data",
			Tags: []string{"a", "b", "c"},
			Metadata: map[string]interface{}{
				"processed": true,
				"count":     100,
			},
			Nested: &NestedArg{Value: 2.718, Enabled: true},
		},
	}
)

// === UNMARSHAL BENCHMARKS (argument parsing) ===

func BenchmarkJSONv1_Unmarshal_Simple(b *testing.B) {
	for b.Loop() {
		var arg SimpleArg
		_ = json.Unmarshal(simpleArgJSON, &arg)
	}
}

func BenchmarkJSONv2_Unmarshal_Simple(b *testing.B) {
	for b.Loop() {
		var arg SimpleArg
		_ = jsonv2.Unmarshal(simpleArgJSON, &arg)
	}
}

func BenchmarkJSONv1_Unmarshal_Complex(b *testing.B) {
	for b.Loop() {
		var arg ComplexArg
		_ = json.Unmarshal(complexArgJSON, &arg)
	}
}

func BenchmarkJSONv2_Unmarshal_Complex(b *testing.B) {
	for b.Loop() {
		var arg ComplexArg
		_ = jsonv2.Unmarshal(complexArgJSON, &arg)
	}
}

func BenchmarkJSONv1_Unmarshal_Interface(b *testing.B) {
	for b.Loop() {
		var arg interface{}
		_ = json.Unmarshal(complexArgJSON, &arg)
	}
}

func BenchmarkJSONv2_Unmarshal_Interface(b *testing.B) {
	for b.Loop() {
		var arg interface{}
		_ = jsonv2.Unmarshal(complexArgJSON, &arg)
	}
}

// === MARSHAL BENCHMARKS (result serialization) ===

func BenchmarkJSONv1_Marshal_Simple(b *testing.B) {
	for b.Loop() {
		_, _ = json.Marshal(simpleResult)
	}
}

func BenchmarkJSONv2_Marshal_Simple(b *testing.B) {
	for b.Loop() {
		_, _ = jsonv2.Marshal(simpleResult)
	}
}

func BenchmarkJSONv1_Marshal_Complex(b *testing.B) {
	for b.Loop() {
		_, _ = json.Marshal(complexResult)
	}
}

func BenchmarkJSONv2_Marshal_Complex(b *testing.B) {
	for b.Loop() {
		_, _ = jsonv2.Marshal(complexResult)
	}
}

// === RAW MESSAGE HANDLING (common in Wails bindings) ===

func BenchmarkJSONv1_RawMessage_Unmarshal(b *testing.B) {
	raw := json.RawMessage(complexArgJSON)
	for b.Loop() {
		var arg ComplexArg
		_ = json.Unmarshal(raw, &arg)
	}
}

func BenchmarkJSONv2_RawMessage_Unmarshal(b *testing.B) {
	raw := jsontext.Value(complexArgJSON)
	for b.Loop() {
		var arg ComplexArg
		_ = jsonv2.Unmarshal(raw, &arg)
	}
}

// === SLICE ARGUMENTS (common pattern) ===

var sliceArgJSON = []byte(`[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]`)
var largeSliceArgJSON = func() []byte {
	data, _ := json.Marshal(make([]int, 100))
	return data
}()

func BenchmarkJSONv1_Unmarshal_Slice(b *testing.B) {
	for b.Loop() {
		var arg []int
		_ = json.Unmarshal(sliceArgJSON, &arg)
	}
}

func BenchmarkJSONv2_Unmarshal_Slice(b *testing.B) {
	for b.Loop() {
		var arg []int
		_ = jsonv2.Unmarshal(sliceArgJSON, &arg)
	}
}

func BenchmarkJSONv1_Unmarshal_LargeSlice(b *testing.B) {
	for b.Loop() {
		var arg []int
		_ = json.Unmarshal(largeSliceArgJSON, &arg)
	}
}

func BenchmarkJSONv2_Unmarshal_LargeSlice(b *testing.B) {
	for b.Loop() {
		var arg []int
		_ = jsonv2.Unmarshal(largeSliceArgJSON, &arg)
	}
}

// === STRING ARGUMENT (most common) ===

var stringArgJSON = []byte(`"hello world this is a test string"`)

func BenchmarkJSONv1_Unmarshal_String(b *testing.B) {
	for b.Loop() {
		var arg string
		_ = json.Unmarshal(stringArgJSON, &arg)
	}
}

func BenchmarkJSONv2_Unmarshal_String(b *testing.B) {
	for b.Loop() {
		var arg string
		_ = jsonv2.Unmarshal(stringArgJSON, &arg)
	}
}

// === MULTIPLE ARGUMENTS (simulating method call) ===

var multiArgJSON = [][]byte{
	[]byte(`"arg1"`),
	[]byte(`42`),
	[]byte(`true`),
	[]byte(`{"key": "value"}`),
}

func BenchmarkJSONv1_Unmarshal_MultiArgs(b *testing.B) {
	for b.Loop() {
		var s string
		var i int
		var bl bool
		var m map[string]string
		_ = json.Unmarshal(multiArgJSON[0], &s)
		_ = json.Unmarshal(multiArgJSON[1], &i)
		_ = json.Unmarshal(multiArgJSON[2], &bl)
		_ = json.Unmarshal(multiArgJSON[3], &m)
	}
}

func BenchmarkJSONv2_Unmarshal_MultiArgs(b *testing.B) {
	for b.Loop() {
		var s string
		var i int
		var bl bool
		var m map[string]string
		_ = jsonv2.Unmarshal(multiArgJSON[0], &s)
		_ = jsonv2.Unmarshal(multiArgJSON[1], &i)
		_ = jsonv2.Unmarshal(multiArgJSON[2], &bl)
		_ = jsonv2.Unmarshal(multiArgJSON[3], &m)
	}
}

// === ERROR RESPONSE MARSHALING ===

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

var errorResp = ErrorResponse{
	Code:    500,
	Message: "Internal server error",
	Details: "Something went wrong while processing the request",
}

func BenchmarkJSONv1_Marshal_Error(b *testing.B) {
	for b.Loop() {
		_, _ = json.Marshal(errorResp)
	}
}

func BenchmarkJSONv2_Marshal_Error(b *testing.B) {
	for b.Loop() {
		_, _ = jsonv2.Marshal(errorResp)
	}
}
