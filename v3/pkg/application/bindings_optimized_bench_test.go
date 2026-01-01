//go:build bench

package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"

	jsonv2 "github.com/go-json-experiment/json"
)

// This file contains optimized versions of BoundMethod.Call for benchmarking.
// These demonstrate potential optimizations that could be applied.

// Pools for reducing allocations
var (
	// Pool for []reflect.Value slices (sized for typical arg counts)
	callArgsPool = sync.Pool{
		New: func() any {
			// Pre-allocate for up to 8 args (covers vast majority of methods)
			return make([]reflect.Value, 0, 8)
		},
	}

	// Pool for []any slices
	anySlicePool = sync.Pool{
		New: func() any {
			return make([]any, 0, 4)
		},
	}

	// Pool for CallError structs
	callErrorPool = sync.Pool{
		New: func() any {
			return &CallError{}
		},
	}
)

// CallOptimized is an optimized version of BoundMethod.Call that uses sync.Pool
func (b *BoundMethod) CallOptimized(ctx context.Context, args []json.RawMessage) (result any, err error) {
	defer handlePanic(handlePanicOptions{skipEnd: 5})

	argCount := len(args)
	if b.needsContext {
		argCount++
	}

	if argCount != len(b.Inputs) {
		cerr := callErrorPool.Get().(*CallError)
		cerr.Kind = TypeError
		cerr.Message = fmt.Sprintf("%s expects %d arguments, got %d", b.FQN, len(b.Inputs), argCount)
		cerr.Cause = nil
		return nil, cerr
	}

	// Get callArgs from pool
	callArgs := callArgsPool.Get().([]reflect.Value)
	callArgs = callArgs[:0] // Reset length but keep capacity

	// Ensure capacity
	if cap(callArgs) < argCount {
		callArgs = make([]reflect.Value, 0, argCount)
	}
	callArgs = callArgs[:argCount]

	base := 0
	if b.needsContext {
		callArgs[0] = reflect.ValueOf(ctx)
		base++
	}

	// Iterate over given arguments
	for index, arg := range args {
		value := reflect.New(b.Inputs[base+index].ReflectType)
		err = json.Unmarshal(arg, value.Interface())
		if err != nil {
			// Return callArgs to pool before returning error
			callArgsPool.Put(callArgs[:0])

			cerr := callErrorPool.Get().(*CallError)
			cerr.Kind = TypeError
			cerr.Message = fmt.Sprintf("could not parse argument #%d: %s", index, err)
			cerr.Cause = json.RawMessage(b.marshalError(err))
			return nil, cerr
		}
		callArgs[base+index] = value.Elem()
	}

	// Do the call
	var callResults []reflect.Value
	if b.Method.Type().IsVariadic() {
		callResults = b.Method.CallSlice(callArgs)
	} else {
		callResults = b.Method.Call(callArgs)
	}

	// Return callArgs to pool
	callArgsPool.Put(callArgs[:0])

	// Get output slice from pool
	nonErrorOutputs := anySlicePool.Get().([]any)
	nonErrorOutputs = nonErrorOutputs[:0]
	defer func() {
		anySlicePool.Put(nonErrorOutputs[:0])
	}()

	var errorOutputs []error

	for _, field := range callResults {
		if field.Type() == errorType {
			if field.IsNil() {
				continue
			}
			if errorOutputs == nil {
				errorOutputs = make([]error, 0, len(callResults)-len(nonErrorOutputs))
				nonErrorOutputs = nil
			}
			errorOutputs = append(errorOutputs, field.Interface().(error))
		} else if nonErrorOutputs != nil {
			nonErrorOutputs = append(nonErrorOutputs, field.Interface())
		}
	}

	if len(errorOutputs) > 0 {
		info := make([]json.RawMessage, len(errorOutputs))
		for i, err := range errorOutputs {
			info[i] = b.marshalError(err)
		}

		cerr := &CallError{
			Kind:    RuntimeError,
			Message: errors.Join(errorOutputs...).Error(),
			Cause:   info,
		}
		if len(info) == 1 {
			cerr.Cause = info[0]
		}
		return nil, cerr
	}

	if len(nonErrorOutputs) == 1 {
		result = nonErrorOutputs[0]
	} else if len(nonErrorOutputs) > 1 {
		// Need to copy since we're returning the pooled slice
		resultSlice := make([]any, len(nonErrorOutputs))
		copy(resultSlice, nonErrorOutputs)
		result = resultSlice
	}

	return result, nil
}

// Benchmark comparing original vs optimized Call
func BenchmarkCallOriginal(b *testing.B) {
	_ = New(Options{})
	bindings := NewBindings(nil, nil)

	service := &benchService{}
	_ = bindings.Add(NewService(service))

	callOptions := &CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application.benchService.StringArg",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	args := []json.RawMessage{json.RawMessage(`"hello world"`)}
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = method.Call(ctx, args)
	}
}

func BenchmarkCallOptimized(b *testing.B) {
	_ = New(Options{})
	bindings := NewBindings(nil, nil)

	service := &benchService{}
	_ = bindings.Add(NewService(service))

	callOptions := &CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application.benchService.StringArg",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	args := []json.RawMessage{json.RawMessage(`"hello world"`)}
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = method.CallOptimized(ctx, args)
	}
}

// benchService for internal tests
type benchService struct{}

func (s *benchService) StringArg(str string) string {
	return str
}

func (s *benchService) MultipleArgs(s1 string, i int, b bool) (string, int, bool) {
	return s1, i, b
}

func BenchmarkCallOriginal_MultiArgs(b *testing.B) {
	_ = New(Options{})
	bindings := NewBindings(nil, nil)

	service := &benchService{}
	_ = bindings.Add(NewService(service))

	callOptions := &CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application.benchService.MultipleArgs",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	args := []json.RawMessage{
		json.RawMessage(`"test"`),
		json.RawMessage(`42`),
		json.RawMessage(`true`),
	}
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = method.Call(ctx, args)
	}
}

func BenchmarkCallOptimized_MultiArgs(b *testing.B) {
	_ = New(Options{})
	bindings := NewBindings(nil, nil)

	service := &benchService{}
	_ = bindings.Add(NewService(service))

	callOptions := &CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application.benchService.MultipleArgs",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	args := []json.RawMessage{
		json.RawMessage(`"test"`),
		json.RawMessage(`42`),
		json.RawMessage(`true`),
	}
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = method.CallOptimized(ctx, args)
	}
}

// CallWithJSONv2 uses the new JSON v2 library for unmarshaling
func (b *BoundMethod) CallWithJSONv2(ctx context.Context, args []json.RawMessage) (result any, err error) {
	defer handlePanic(handlePanicOptions{skipEnd: 5})

	argCount := len(args)
	if b.needsContext {
		argCount++
	}

	if argCount != len(b.Inputs) {
		return nil, &CallError{
			Kind:    TypeError,
			Message: fmt.Sprintf("%s expects %d arguments, got %d", b.FQN, len(b.Inputs), argCount),
		}
	}

	// Convert inputs to values of appropriate type
	callArgs := make([]reflect.Value, argCount)
	base := 0

	if b.needsContext {
		callArgs[0] = reflect.ValueOf(ctx)
		base++
	}

	// Iterate over given arguments - use JSON v2 for unmarshaling
	for index, arg := range args {
		value := reflect.New(b.Inputs[base+index].ReflectType)
		err = jsonv2.Unmarshal(arg, value.Interface())
		if err != nil {
			return nil, &CallError{
				Kind:    TypeError,
				Message: fmt.Sprintf("could not parse argument #%d: %s", index, err),
				Cause:   json.RawMessage(b.marshalError(err)),
			}
		}
		callArgs[base+index] = value.Elem()
	}

	// Do the call
	var callResults []reflect.Value
	if b.Method.Type().IsVariadic() {
		callResults = b.Method.CallSlice(callArgs)
	} else {
		callResults = b.Method.Call(callArgs)
	}

	var nonErrorOutputs = make([]any, 0, len(callResults))
	var errorOutputs []error

	for _, field := range callResults {
		if field.Type() == errorType {
			if field.IsNil() {
				continue
			}
			if errorOutputs == nil {
				errorOutputs = make([]error, 0, len(callResults)-len(nonErrorOutputs))
				nonErrorOutputs = nil
			}
			errorOutputs = append(errorOutputs, field.Interface().(error))
		} else if nonErrorOutputs != nil {
			nonErrorOutputs = append(nonErrorOutputs, field.Interface())
		}
	}

	if len(errorOutputs) > 0 {
		info := make([]json.RawMessage, len(errorOutputs))
		for i, err := range errorOutputs {
			info[i] = b.marshalError(err)
		}

		cerr := &CallError{
			Kind:    RuntimeError,
			Message: errors.Join(errorOutputs...).Error(),
			Cause:   info,
		}
		if len(info) == 1 {
			cerr.Cause = info[0]
		}
		return nil, cerr
	}

	if len(nonErrorOutputs) == 1 {
		result = nonErrorOutputs[0]
	} else if len(nonErrorOutputs) > 1 {
		result = nonErrorOutputs
	}

	return result, nil
}

func BenchmarkCallJSONv2(b *testing.B) {
	_ = New(Options{})
	bindings := NewBindings(nil, nil)

	service := &benchService{}
	_ = bindings.Add(NewService(service))

	callOptions := &CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application.benchService.StringArg",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	args := []json.RawMessage{json.RawMessage(`"hello world"`)}
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = method.CallWithJSONv2(ctx, args)
	}
}

func BenchmarkCallJSONv2_MultiArgs(b *testing.B) {
	_ = New(Options{})
	bindings := NewBindings(nil, nil)

	service := &benchService{}
	_ = bindings.Add(NewService(service))

	callOptions := &CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application.benchService.MultipleArgs",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	args := []json.RawMessage{
		json.RawMessage(`"test"`),
		json.RawMessage(`42`),
		json.RawMessage(`true`),
	}
	ctx := context.Background()

	b.ResetTimer()
	for b.Loop() {
		_, _ = method.CallWithJSONv2(ctx, args)
	}
}

// Concurrent benchmark to test pool effectiveness under load
func BenchmarkCallOriginal_Concurrent(b *testing.B) {
	_ = New(Options{})
	bindings := NewBindings(nil, nil)

	service := &benchService{}
	_ = bindings.Add(NewService(service))

	callOptions := &CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application.benchService.StringArg",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	args := []json.RawMessage{json.RawMessage(`"hello world"`)}
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = method.Call(ctx, args)
		}
	})
}

func BenchmarkCallOptimized_Concurrent(b *testing.B) {
	_ = New(Options{})
	bindings := NewBindings(nil, nil)

	service := &benchService{}
	_ = bindings.Add(NewService(service))

	callOptions := &CallOptions{
		MethodName: "github.com/wailsapp/wails/v3/pkg/application.benchService.StringArg",
	}
	method := bindings.Get(callOptions)
	if method == nil {
		b.Fatal("method not found")
	}

	args := []json.RawMessage{json.RawMessage(`"hello world"`)}
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = method.CallOptimized(ctx, args)
		}
	})
}
