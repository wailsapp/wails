//go:build bench

package application_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wailsapp/wails/v3/internal/hash"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// BenchmarkService provides methods with varying complexity for benchmarking
type BenchmarkService struct{}

func (s *BenchmarkService) NoArgs() {}

func (s *BenchmarkService) StringArg(str string) string {
	return str
}

func (s *BenchmarkService) IntArg(i int) int {
	return i
}

func (s *BenchmarkService) MultipleArgs(s1 string, i int, b bool) (string, int, bool) {
	return s1, i, b
}

func (s *BenchmarkService) StructArg(p BenchPerson) BenchPerson {
	return p
}

func (s *BenchmarkService) ComplexStruct(c ComplexData) ComplexData {
	return c
}

func (s *BenchmarkService) SliceArg(items []int) []int {
	return items
}

func (s *BenchmarkService) VariadicArg(items ...string) []string {
	return items
}

func (s *BenchmarkService) WithContext(ctx context.Context, s1 string) string {
	return s1
}

func (s *BenchmarkService) Method1()  {}
func (s *BenchmarkService) Method2()  {}
func (s *BenchmarkService) Method3()  {}
func (s *BenchmarkService) Method4()  {}
func (s *BenchmarkService) Method5()  {}
func (s *BenchmarkService) Method6()  {}
func (s *BenchmarkService) Method7()  {}
func (s *BenchmarkService) Method8()  {}
func (s *BenchmarkService) Method9()  {}
func (s *BenchmarkService) Method10() {}

type BenchPerson struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type ComplexData struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
	Nested   *NestedData            `json:"nested"`
}

type NestedData struct {
	Value   float64 `json:"value"`
	Enabled bool    `json:"enabled"`
}

// Helper to create JSON args
func benchArgs(jsonArgs ...string) []json.RawMessage {
	args := make([]json.RawMessage, len(jsonArgs))
	for i, j := range jsonArgs {
		args[i] = json.RawMessage(j)
	}
	return args
}

// BenchmarkMethodBinding measures the cost of registering services with varying method counts
func BenchmarkMethodBinding(b *testing.B) {
	// Initialize global application (required for bindings)
	_ = application.New(application.Options{})

	b.Run("SingleService", func(b *testing.B) {
		for b.Loop() {
			bindings := application.NewBindings(nil, nil)
			_ = bindings.Add(application.NewService(&BenchmarkService{}))
		}
	})

	b.Run("MultipleServices", func(b *testing.B) {
		for b.Loop() {
			bindings := application.NewBindings(nil, nil)
			_ = bindings.Add(application.NewService(&BenchmarkService{}))
			_ = bindings.Add(application.NewService(&BenchPerson{})) // Will fail but tests the path
		}
	})
}

// BenchmarkMethodLookupByID measures method lookup by ID performance
func BenchmarkMethodLookupByID(b *testing.B) {
	_ = application.New(application.Options{})
	bindings := application.NewBindings(nil, nil)
	_ = bindings.Add(application.NewService(&BenchmarkService{}))

	// Get a valid method ID
	callOptions := &application.CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.StringArg",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}
	methodID := method.ID

	b.ResetTimer()
	for b.Loop() {
		_ = bindings.GetByID(methodID)
	}
}

// BenchmarkMethodLookupByName measures method lookup by name performance
func BenchmarkMethodLookupByName(b *testing.B) {
	_ = application.New(application.Options{})
	bindings := application.NewBindings(nil, nil)
	_ = bindings.Add(application.NewService(&BenchmarkService{}))

	callOptions := &application.CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.StringArg",
	}

	b.ResetTimer()
	for b.Loop() {
		_ = bindings.Get(callOptions)
	}
}

// BenchmarkSimpleCall measures the cost of calling a method with a simple string argument
func BenchmarkSimpleCall(b *testing.B) {
	_ = application.New(application.Options{})
	bindings := application.NewBindings(nil, nil)
	_ = bindings.Add(application.NewService(&BenchmarkService{}))

	callOptions := &application.CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.StringArg",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	args := benchArgs(`"hello world"`)
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = method.Call(ctx, args)
	}
}

// BenchmarkComplexCall measures the cost of calling a method with a complex struct argument
func BenchmarkComplexCall(b *testing.B) {
	_ = application.New(application.Options{})
	bindings := application.NewBindings(nil, nil)
	_ = bindings.Add(application.NewService(&BenchmarkService{}))

	callOptions := &application.CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.ComplexStruct",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	complexArg := `{
		"id": 12345,
		"name": "Test Complex Data",
		"tags": ["tag1", "tag2", "tag3", "tag4", "tag5"],
		"metadata": {"key1": "value1", "key2": 42, "key3": true},
		"nested": {"value": 3.14159, "enabled": true}
	}`
	args := benchArgs(complexArg)
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = method.Call(ctx, args)
	}
}

// BenchmarkVariadicCall measures the cost of calling a variadic method
func BenchmarkVariadicCall(b *testing.B) {
	_ = application.New(application.Options{})
	bindings := application.NewBindings(nil, nil)
	_ = bindings.Add(application.NewService(&BenchmarkService{}))

	callOptions := &application.CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.VariadicArg",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	args := benchArgs(`["one", "two", "three", "four", "five"]`)
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = method.Call(ctx, args)
	}
}

// BenchmarkCallWithContext measures the cost of calling a method that requires context
func BenchmarkCallWithContext(b *testing.B) {
	_ = application.New(application.Options{})
	bindings := application.NewBindings(nil, nil)
	_ = bindings.Add(application.NewService(&BenchmarkService{}))

	callOptions := &application.CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.WithContext",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	args := benchArgs(`"context test"`)
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = method.Call(ctx, args)
	}
}

// BenchmarkJSONMarshalResult measures JSON marshaling overhead for results
func BenchmarkJSONMarshalResult(b *testing.B) {
	person := BenchPerson{
		Name:    "John Doe",
		Age:     30,
		Email:   "john@example.com",
		Address: "123 Main St, City, Country",
	}

	b.Run("SimplePerson", func(b *testing.B) {
		for b.Loop() {
			_, _ = json.Marshal(person)
		}
	})

	complex := ComplexData{
		ID:   12345,
		Name: "Complex Test",
		Tags: []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
		Metadata: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
			"key3": true,
		},
		Nested: &NestedData{
			Value:   3.14159,
			Enabled: true,
		},
	}

	b.Run("ComplexData", func(b *testing.B) {
		for b.Loop() {
			_, _ = json.Marshal(complex)
		}
	})
}

// BenchmarkHashComputation measures the FNV hash computation used for method IDs
func BenchmarkHashComputation(b *testing.B) {
	testCases := []struct {
		name string
		fqn  string
	}{
		{"Short", "pkg.Service.Method"},
		{"Medium", "github.com/user/project/pkg.Service.Method"},
		{"Long", "github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.ComplexStruct"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				_ = hash.Fnv(tc.fqn)
			}
		})
	}
}

// BenchmarkJSONUnmarshal measures JSON unmarshaling overhead for arguments
func BenchmarkJSONUnmarshal(b *testing.B) {
	b.Run("String", func(b *testing.B) {
		data := []byte(`"hello world"`)
		for b.Loop() {
			var s string
			_ = json.Unmarshal(data, &s)
		}
	})

	b.Run("Int", func(b *testing.B) {
		data := []byte(`12345`)
		for b.Loop() {
			var i int
			_ = json.Unmarshal(data, &i)
		}
	})

	b.Run("Struct", func(b *testing.B) {
		data := []byte(`{"name":"John","age":30,"email":"john@example.com","address":"123 Main St"}`)
		for b.Loop() {
			var p BenchPerson
			_ = json.Unmarshal(data, &p)
		}
	})

	b.Run("ComplexStruct", func(b *testing.B) {
		data := []byte(`{"id":12345,"name":"Test","tags":["a","b","c"],"metadata":{"k":"v"},"nested":{"value":3.14,"enabled":true}}`)
		for b.Loop() {
			var c ComplexData
			_ = json.Unmarshal(data, &c)
		}
	})
}

// BenchmarkMethodLookupWithAliases measures method lookup with alias resolution
func BenchmarkMethodLookupWithAliases(b *testing.B) {
	_ = application.New(application.Options{})

	// Create aliases map
	aliases := make(map[uint32]uint32)
	for i := uint32(0); i < 100; i++ {
		aliases[i+1000] = i
	}

	bindings := application.NewBindings(nil, aliases)
	_ = bindings.Add(application.NewService(&BenchmarkService{}))

	callOptions := &application.CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.StringArg",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	b.Run("DirectLookup", func(b *testing.B) {
		id := method.ID
		for b.Loop() {
			_ = bindings.GetByID(id)
		}
	})

	b.Run("AliasLookup", func(b *testing.B) {
		// Add an alias for this method
		aliases[9999] = method.ID
		for b.Loop() {
			_ = bindings.GetByID(9999)
		}
	})
}

// BenchmarkReflectValueCall measures the overhead of reflect.Value.Call
func BenchmarkReflectValueCall(b *testing.B) {
	_ = application.New(application.Options{})
	bindings := application.NewBindings(nil, nil)
	_ = bindings.Add(application.NewService(&BenchmarkService{}))

	ctx := context.Background()

	b.Run("NoArgs", func(b *testing.B) {
		callOptions := &application.CallOptions{
			MethodName: "github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.NoArgs",
		}
		method := bindings.Get(callOptions)
		if method == nil {
			b.Fatal("method not found")
		}
		args := benchArgs()
		for b.Loop() {
			_, _ = method.Call(ctx, args)
		}
	})

	b.Run("MultipleArgs", func(b *testing.B) {
		callOptions := &application.CallOptions{
			MethodName: "github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.MultipleArgs",
		}
		method := bindings.Get(callOptions)
		if method == nil {
			b.Fatal("method not found")
		}
		args := benchArgs(`"test"`, `42`, `true`)
		for b.Loop() {
			_, _ = method.Call(ctx, args)
		}
	})
}

// BenchmarkBindingsScaling measures how bindings performance scales with service count
func BenchmarkBindingsScaling(b *testing.B) {
	_ = application.New(application.Options{})

	// We can only add one service of each type, so we test lookup scaling
	bindings := application.NewBindings(nil, nil)
	_ = bindings.Add(application.NewService(&BenchmarkService{}))

	// Generate method names for lookup
	methodNames := []string{
		"github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.NoArgs",
		"github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.StringArg",
		"github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.IntArg",
		"github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.MultipleArgs",
		"github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.StructArg",
	}

	b.Run("SequentialLookup", func(b *testing.B) {
		for b.Loop() {
			for _, name := range methodNames {
				_ = bindings.Get(&application.CallOptions{MethodName: name})
			}
		}
	})
}

// BenchmarkCallErrorPath measures the cost of error handling in method calls
func BenchmarkCallErrorPath(b *testing.B) {
	_ = application.New(application.Options{})
	bindings := application.NewBindings(nil, nil)
	_ = bindings.Add(application.NewService(&BenchmarkService{}))

	callOptions := &application.CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.StringArg",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	ctx := context.Background()

	b.Run("WrongArgCount", func(b *testing.B) {
		args := benchArgs() // No args when one is expected
		for b.Loop() {
			_, _ = method.Call(ctx, args)
		}
	})

	b.Run("WrongArgType", func(b *testing.B) {
		args := benchArgs(`123`) // Int when string is expected
		for b.Loop() {
			_, _ = method.Call(ctx, args)
		}
	})
}

// BenchmarkSliceArgSizes measures performance with varying slice sizes
func BenchmarkSliceArgSizes(b *testing.B) {
	_ = application.New(application.Options{})
	bindings := application.NewBindings(nil, nil)
	_ = bindings.Add(application.NewService(&BenchmarkService{}))

	callOptions := &application.CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application_test.BenchmarkService.SliceArg",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	ctx := context.Background()

	sizes := []int{1, 10, 100, 1000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
			// Build slice JSON
			slice := make([]int, size)
			for i := range slice {
				slice[i] = i
			}
			data, _ := json.Marshal(slice)
			args := []json.RawMessage{data}

			b.ResetTimer()
			for b.Loop() {
				_, _ = method.Call(ctx, args)
			}
		})
	}
}
