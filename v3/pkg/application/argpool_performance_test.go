package application

import (
	"encoding/json"
	"testing"
)

// Isolated struct allocation benchmarks
func BenchmarkCallOptionsStructOnly(b *testing.B) {
	// Pre-allocate JSON messages to isolate struct allocation
	arg1 := json.RawMessage(`{"key": "value1"}`)
	arg2 := json.RawMessage(`{"key": "value2"}`)

	b.Run("Baseline", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			opts := &CallOptions{
				MethodID:   12345,
				MethodName: "TestMethod",
			}
			// Simulate using pre-allocated slice capacity
			opts.Args = []json.RawMessage{arg1, arg2}
			_ = opts // Prevent optimization
		}
	})

	b.Run("Pooled", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			opts := GetCallOptions()
			opts.MethodID = 12345
			opts.MethodName = "TestMethod"
			opts.Args = []json.RawMessage{arg1, arg2}
			PutCallOptions(opts)
		}
	})
}

// Benchmark Args struct allocation only
func BenchmarkArgsStructOnly(b *testing.B) {
	b.Run("Baseline", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			args := &Args{
				data: map[string]any{
					"call-id": "test-12345",
					"method":  "TestMethod",
					"param1":  "value1",
					"param2":  42,
				},
			}
			_ = args // Prevent optimization
		}
	})

	b.Run("Pooled", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			args := GetArgs()
			args.data["call-id"] = "test-12345"
			args.data["method"] = "TestMethod"
			args.data["param1"] = "value1"
			args.data["param2"] = 42
			PutArgs(args)
		}
	})
}

// Benchmark contention for object pools
func BenchmarkPoolContentionRealWorld(b *testing.B) {
	testJSON := json.RawMessage(`{"key": "value"}`)

	b.Run("CallOptions-Contention", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				opts := GetCallOptions()
				opts.MethodID = 12345
				opts.MethodName = "TestMethod"
				opts.Args = append(opts.Args, testJSON)
				
				// Simulate some work
				_ = opts.MethodName
				_ = len(opts.Args)
				
				PutCallOptions(opts)
			}
		})
	})

	b.Run("Args-Contention", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				args := GetArgs()
				args.data["method"] = "TestMethod"
				args.data["id"] = 12345
				
				// Simulate some work
				_ = args.String("method")
				_ = args.Int("id")
				
				PutArgs(args)
			}
		})
	})
}

// Measure allocations in realistic workflow
func BenchmarkRealWorldWorkflow(b *testing.B) {
	queryParams := QueryParams{
		"method": []string{"call"},
		"args":   []string{`{"call-id":"test-12345","method":"TestService.TestMethod","params":[{"name":"test","value":42}]}`},
	}

	b.Run("Baseline-Workflow", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			// Traditional allocation workflow
			args := &Args{data: make(map[string]any)}
			argData := queryParams["args"]
			if len(argData) == 1 {
				json.Unmarshal([]byte(argData[0]), &args.data)
			}

			opts := &CallOptions{}
			queryParams.ToStruct(opts)

			// Simulate method lookup and execution
			callID := args.String("call-id")
			methodName := opts.MethodName
			
			// Use the values
			_ = callID
			_ = methodName
		}
	})

	b.Run("Pooled-Workflow", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			// Pooled allocation workflow
			args, _ := ArgsFromQueryParams(queryParams)
			opts := GetCallOptions()
			queryParams.ToStruct(opts)

			// Simulate method lookup and execution
			callID := args.String("call-id")
			methodName := opts.MethodName
			
			// Use the values
			_ = callID
			_ = methodName

			// Return to pools
			PutArgs(args)
			PutCallOptions(opts)
		}
	})
}

// High-frequency simulation test
func BenchmarkHighFrequencyMethodCalls(b *testing.B) {
	// Simulate high-frequency method calls like those in real applications
	methodCalls := []struct {
		id   uint32
		name string
		args []json.RawMessage
	}{
		{1001, "Service.GetUser", []json.RawMessage{json.RawMessage(`{"id": 123}`)}},
		{1002, "Service.UpdateData", []json.RawMessage{json.RawMessage(`{"data": "value"}`), json.RawMessage(`{"flag": true}`)}},
		{1003, "Service.ProcessEvent", []json.RawMessage{json.RawMessage(`{"event": "click"}`), json.RawMessage(`{"x": 100}`), json.RawMessage(`{"y": 200}`)}},
	}

	b.Run("Baseline-HighFreq", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			call := methodCalls[i%len(methodCalls)]
			opts := &CallOptions{
				MethodID:   call.id,
				MethodName: call.name,
				Args:       make([]json.RawMessage, len(call.args)),
			}
			copy(opts.Args, call.args)
			_ = opts
		}
	})

	b.Run("Pooled-HighFreq", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			call := methodCalls[i%len(methodCalls)]
			opts := GetCallOptions()
			opts.MethodID = call.id
			opts.MethodName = call.name
			opts.Args = append(opts.Args, call.args...)
			PutCallOptions(opts)
		}
	})
}

// Memory pressure test with GC
func BenchmarkMemoryPressure(b *testing.B) {
	b.Run("Baseline-MemPressure", func(b *testing.B) {
		b.ReportAllocs()
		var results []*CallOptions
		for i := 0; i < b.N; i++ {
			opts := &CallOptions{
				MethodID:   uint32(i),
				MethodName: "TestMethod",
				Args:       make([]json.RawMessage, 1),
			}
			opts.Args[0] = json.RawMessage(`{"id": ` + string(rune('0'+(i%10))) + `}`)
			results = append(results, opts)
			
			// Periodically clear to simulate GC pressure
			if len(results) > 1000 {
				results = results[:0]
			}
		}
		_ = results
	})

	b.Run("Pooled-MemPressure", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			opts := GetCallOptions()
			opts.MethodID = uint32(i)
			opts.MethodName = "TestMethod"
			opts.Args = append(opts.Args, json.RawMessage(`{"id": `+string(rune('0'+(i%10)))+`}`))
			
			// Simulate some processing time
			_ = opts.MethodName
			
			PutCallOptions(opts)
		}
	})
}