package application

import (
	"encoding/json"
	"testing"
)

// Benchmark baseline allocation without pooling
func BenchmarkCallOptionsBaseline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opts := &CallOptions{
			MethodID:   12345,
			MethodName: "TestMethod",
			Args:       make([]json.RawMessage, 3),
		}
		opts.Args[0] = json.RawMessage(`{"key": "value1"}`)
		opts.Args[1] = json.RawMessage(`{"key": "value2"}`)
		opts.Args[2] = json.RawMessage(`{"key": "value3"}`)
		_ = opts // Prevent optimization
	}
}

// Benchmark pooled allocation
func BenchmarkCallOptionsPooled(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opts := GetCallOptions()
		opts.MethodID = 12345
		opts.MethodName = "TestMethod"
		opts.Args = append(opts.Args, json.RawMessage(`{"key": "value1"}`))
		opts.Args = append(opts.Args, json.RawMessage(`{"key": "value2"}`))
		opts.Args = append(opts.Args, json.RawMessage(`{"key": "value3"}`))
		PutCallOptions(opts)
	}
}

// Benchmark Args baseline allocation
func BenchmarkArgsBaseline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		args := &Args{
			data: make(map[string]any, 8),
		}
		args.data["call-id"] = "test-12345"
		args.data["method"] = "TestMethod"
		args.data["param1"] = "value1"
		args.data["param2"] = 42
		args.data["param3"] = true
		_ = args // Prevent optimization
	}
}

// Benchmark Args pooled allocation
func BenchmarkArgsPooled(b *testing.B) {
	for i := 0; i < b.N; i++ {
		args := GetArgs()
		args.data["call-id"] = "test-12345"
		args.data["method"] = "TestMethod"
		args.data["param1"] = "value1"
		args.data["param2"] = 42
		args.data["param3"] = true
		PutArgs(args)
	}
}

// Benchmark QueryParams baseline allocation
func BenchmarkQueryParamsBaseline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		qp := make(QueryParams, 4)
		qp["method"] = []string{"call"}
		qp["args"] = []string{`{"call-id":"test-12345","method":"TestMethod"}`}
		qp["window-id"] = []string{"123"}
		qp["type"] = []string{"binding"}
		_ = qp // Prevent optimization
	}
}

// Benchmark QueryParams pooled allocation
func BenchmarkQueryParamsPooled(b *testing.B) {
	for i := 0; i < b.N; i++ {
		qp := GetQueryParams()
		qp["method"] = []string{"call"}
		qp["args"] = []string{`{"call-id":"test-12345","method":"TestMethod"}`}
		qp["window-id"] = []string{"123"}
		qp["type"] = []string{"binding"}
		PutQueryParams(qp)
	}
}

// Benchmark Parameter slice baseline allocation
func BenchmarkParameterSliceBaseline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		params := make([]*Parameter, 0, 4)
		params = append(params, &Parameter{Name: "ctx", TypeName: "context.Context"})
		params = append(params, &Parameter{Name: "arg1", TypeName: "string"})
		params = append(params, &Parameter{Name: "arg2", TypeName: "int"})
		params = append(params, &Parameter{Name: "arg3", TypeName: "bool"})
		_ = params // Prevent optimization
	}
}

// Benchmark Parameter slice pooled allocation
func BenchmarkParameterSlicePooled(b *testing.B) {
	for i := 0; i < b.N; i++ {
		params := GetParameterSlice()
		params = append(params, &Parameter{Name: "ctx", TypeName: "context.Context"})
		params = append(params, &Parameter{Name: "arg1", TypeName: "string"})
		params = append(params, &Parameter{Name: "arg2", TypeName: "int"})
		params = append(params, &Parameter{Name: "arg3", TypeName: "bool"})
		PutParameterSlice(params)
	}
}

// Benchmark realistic workload - simulating method call processing
func BenchmarkMethodCallWorkloadBaseline(b *testing.B) {
	testQP := QueryParams{
		"method": []string{"call"},
		"args":   []string{`{"call-id":"test-12345","method":"TestMethod","params":["arg1",42,true]}`},
	}

	for i := 0; i < b.N; i++ {
		// Simulate Args creation
		args := &Args{data: make(map[string]any)}
		argData := testQP["args"]
		if len(argData) == 1 {
			json.Unmarshal([]byte(argData[0]), &args.data)
		}

		// Simulate CallOptions creation
		opts := &CallOptions{
			MethodID:   12345,
			MethodName: "TestMethod",
			Args:       make([]json.RawMessage, 3),
		}

		// Use the objects
		_ = args.String("call-id")
		_ = opts.MethodName
	}
}

// Benchmark realistic workload with pooling
func BenchmarkMethodCallWorkloadPooled(b *testing.B) {
	testQP := QueryParams{
		"method": []string{"call"},
		"args":   []string{`{"call-id":"test-12345","method":"TestMethod","params":["arg1",42,true]}`},
	}

	for i := 0; i < b.N; i++ {
		// Use pooled Args creation
		args, _ := ArgsFromQueryParams(testQP)

		// Use pooled CallOptions creation
		opts := GetCallOptions()
		opts.MethodID = 12345
		opts.MethodName = "TestMethod"

		// Use the objects
		_ = args.String("call-id")
		_ = opts.MethodName

		// Return to pools
		PutArgs(args)
		PutCallOptions(opts)
	}
}

// Benchmark contention scenarios with multiple goroutines
func BenchmarkCallOptionsContention(b *testing.B) {
	b.Run("Baseline", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				opts := &CallOptions{
					MethodID:   12345,
					MethodName: "TestMethod",
					Args:       make([]json.RawMessage, 2),
				}
				opts.Args[0] = json.RawMessage(`{"test": "value"}`)
				_ = opts
			}
		})
	})

	b.Run("Pooled", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				opts := GetCallOptions()
				opts.MethodID = 12345
				opts.MethodName = "TestMethod"
				opts.Args = append(opts.Args, json.RawMessage(`{"test": "value"}`))
				PutCallOptions(opts)
			}
		})
	})
}

// Test correctness of pooled implementations
func TestPooledImplementationsCorrectness(t *testing.T) {
	// Test CallOptions pooling
	opts1 := GetCallOptions()
	opts1.MethodID = 12345
	opts1.MethodName = "TestMethod"
	opts1.Args = append(opts1.Args, json.RawMessage(`{"test": "value"}`))

	if opts1.MethodID != 12345 || opts1.MethodName != "TestMethod" || len(opts1.Args) != 1 {
		t.Error("CallOptions pooling failed correctness test")
	}

	PutCallOptions(opts1)

	// Get another instance - should be reset
	opts2 := GetCallOptions()
	if opts2.MethodID != 0 || opts2.MethodName != "" || len(opts2.Args) != 0 {
		t.Error("CallOptions reset failed")
	}
	PutCallOptions(opts2)

	// Test Args pooling
	args1 := GetArgs()
	args1.data["test"] = "value"
	args1.data["number"] = 42

	if val := args1.String("test"); val == nil || *val != "value" {
		t.Error("Args pooling failed correctness test")
	}

	PutArgs(args1)

	// Get another instance - should be reset
	args2 := GetArgs()
	if len(args2.data) != 0 {
		t.Error("Args reset failed")
	}
	PutArgs(args2)

	// Test QueryParams pooling
	qp1 := GetQueryParams()
	qp1["test"] = []string{"value"}

	if vals := qp1["test"]; len(vals) != 1 || vals[0] != "value" {
		t.Error("QueryParams pooling failed correctness test")
	}

	PutQueryParams(qp1)

	// Get another instance - should be reset
	qp2 := GetQueryParams()
	if len(qp2) != 0 {
		t.Error("QueryParams reset failed")
	}
	PutQueryParams(qp2)
}

// Memory allocation benchmark to measure actual allocation reduction
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run("Baseline-Allocs", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			opts := &CallOptions{
				MethodID:   12345,
				MethodName: "TestMethod",
				Args:       make([]json.RawMessage, 2),
			}
			args := &Args{data: make(map[string]any, 4)}
			qp := make(QueryParams, 3)
			
			// Use objects to prevent dead code elimination
			opts.Args[0] = json.RawMessage(`{"key": "value"}`)
			args.data["test"] = "value"
			qp["test"] = []string{"value"}
		}
	})

	b.Run("Pooled-Allocs", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			opts := GetCallOptions()
			args := GetArgs()
			qp := GetQueryParams()
			
			// Use objects to prevent dead code elimination
			opts.Args = append(opts.Args, json.RawMessage(`{"key": "value"}`))
			args.data["test"] = "value"
			qp["test"] = []string{"value"}
			
			PutCallOptions(opts)
			PutArgs(args)
			PutQueryParams(qp)
		}
	})
}