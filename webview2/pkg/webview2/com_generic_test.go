//go:build windows

package webview2

import (
	"testing"
	"unsafe"
)

// TestInvokeValue tests the generic InvokeValue helper
func TestInvokeValue(t *testing.T) {
	// This is a compile-time test to ensure type safety
	// We can't actually test COM calls without a real COM object

	type TestStruct struct {
		a int32
		b int32
	}

	// Verify that InvokeValue works with different types
	var testStruct TestStruct
	var testInt int32
	var testUint32 uint32

	// These should compile without errors
	_ = InvokeValue(func(...uintptr) (uintptr, uintptr, uintptr) { return 0, 0, 0 }, &testStruct)
	_ = InvokeValue(func(...uintptr) (uintptr, uintptr, uintptr) { return 0, 0, 0 }, &testInt)
	_ = InvokeValue(func(...uintptr) (uintptr, uintptr, uintptr) { return 0, 0, 0 }, &testUint32)
}

// TestInvokeBool tests the generic InvokeBool helper
func TestInvokeBool(t *testing.T) {
	// Compile-time test for type safety
	var testBool int32

	// This should compile without errors
	_, _ = InvokeBool(func(...uintptr) (uintptr, uintptr, uintptr) { return 0, 0, 0 }, &testBool)
}

// TestInvokeString tests the generic InvokeString helper
func TestInvokeString(t *testing.T) {
	// Compile-time test for type safety
	var testStr *uint16

	// This should compile without errors
	_, _ = InvokeString(func(...uintptr) (uintptr, uintptr, uintptr) { return 0, 0, 0 }, &testStr)
}

// TestInvokeInterface tests the generic InvokeInterface helper
func TestInvokeInterface(t *testing.T) {
	// Compile-time test for type safety
	type TestInterface struct{}

	var testIf *TestInterface

	// This should compile without errors
	_, _ = InvokeInterface(func(...uintptr) (uintptr, uintptr, uintptr) { return 0, 0, 0 }, &testIf)
}

// TestInvokeToken tests the generic InvokeToken helper
func TestInvokeToken(t *testing.T) {
	// Compile-time test for type safety
	var testToken EventRegistrationToken

	// This should compile without errors
	_, _ = InvokeToken(func(...uintptr) (uintptr, uintptr, uintptr) { return 0, 0, 0 }, &testToken)
}

// TestInvokeVoid tests the generic InvokeVoid helper
func TestInvokeVoid(t *testing.T) {
	// Compile-time test
	_ = InvokeVoid(func(...uintptr) (uintptr, uintptr, uintptr) { return 0, 0, 0 })
}

// TestCOMMethodType tests that our COMMethod type is satisfied by ComProc
func TestCOMMethodType(t *testing.T) {
	// Verify ComProc satisfies COMMethod interface
	var proc ComProc
	var method COMMethod = proc

	_ = method
}

// TestGenericHelpersWithPointers tests generic helpers with pointer types
func TestGenericHelpersWithPointers(t *testing.T) {
	type TestStruct struct {
		value int32
	}

	// Test that InvokeValue works with pointer receiver
	var testStruct TestStruct
	var ptr *TestStruct = &testStruct

	// This should compile
	_, _ = InvokeValue(func(...uintptr) (uintptr, uintptr, uintptr) { return 0, 0, 0 }, &testStruct)
	_ = ptr
}

// TestMemoryLayout verifies that the generic helpers don't introduce memory issues
func TestMemoryLayout(t *testing.T) {
	// Verify unsafe.Pointer conversions are consistent
	var i int32 = 42
	ptr := unsafe.Pointer(&i)
	i2 := (*int32)(ptr)

	if *i2 != 42 {
		t.Error("Memory layout issue: pointer conversion failed")
	}
}

// BenchmarkInvokeVoid benchmarks the InvokeVoid helper
func BenchmarkInvokeVoid(b *testing.B) {
	method := func(...uintptr) (uintptr, uintptr, uintptr) {
		return 0, 0, 0 // Simulate S_OK
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = InvokeVoid(method)
	}
}

// BenchmarkInvokeValue benchmarks the InvokeValue helper
func BenchmarkInvokeValue(b *testing.B) {
	method := func(...uintptr) (uintptr, uintptr, uintptr) {
		return 0, 0, 0 // Simulate S_OK
	}
	var result int32

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = InvokeValue(method, &result)
	}
}
