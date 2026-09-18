package application

import (
	"reflect"
	"testing"
)

func TestParameter_IsType(t *testing.T) {
	param := &Parameter{
		Name:     "test",
		TypeName: "string",
	}

	if !param.IsType("string") {
		t.Error("IsType should return true for matching type")
	}

	if param.IsType("int") {
		t.Error("IsType should return false for non-matching type")
	}
}

func TestParameter_IsError(t *testing.T) {
	errorParam := &Parameter{
		Name:     "err",
		TypeName: "error",
	}

	if !errorParam.IsError() {
		t.Error("IsError should return true for error type")
	}

	stringParam := &Parameter{
		Name:     "s",
		TypeName: "string",
	}

	if stringParam.IsError() {
		t.Error("IsError should return false for non-error type")
	}
}

func TestNewParameter(t *testing.T) {
	stringType := reflect.TypeOf("")
	param := newParameter("myParam", stringType)

	if param.Name != "myParam" {
		t.Errorf("Name = %q, want %q", param.Name, "myParam")
	}

	if param.TypeName != "string" {
		t.Errorf("TypeName = %q, want %q", param.TypeName, "string")
	}

	if param.ReflectType != stringType {
		t.Error("ReflectType not set correctly")
	}
}

func TestCallError_Error(t *testing.T) {
	err := &CallError{
		Kind:    ReferenceError,
		Message: "test error",
	}

	if err.Error() != "test error" {
		t.Errorf("Error() = %q, want %q", err.Error(), "test error")
	}
}

func TestCallError_Kinds(t *testing.T) {
	tests := []struct {
		kind     ErrorKind
		expected string
	}{
		{ReferenceError, "ReferenceError"},
		{TypeError, "TypeError"},
		{RuntimeError, "RuntimeError"},
	}

	for _, tt := range tests {
		if string(tt.kind) != tt.expected {
			t.Errorf("ErrorKind = %q, want %q", string(tt.kind), tt.expected)
		}
	}
}

func TestCallError_WithCause(t *testing.T) {
	cause := map[string]string{"detail": "some detail"}
	err := &CallError{
		Kind:    RuntimeError,
		Message: "runtime error occurred",
		Cause:   cause,
	}

	if err.Error() != "runtime error occurred" {
		t.Error("Error() should return the message")
	}

	if err.Cause == nil {
		t.Error("Cause should be set")
	}
}

func TestCallOptions_Fields(t *testing.T) {
	opts := CallOptions{
		MethodID:   12345,
		MethodName: "TestService.Method",
	}

	if opts.MethodID != 12345 {
		t.Error("MethodID not set correctly")
	}

	if opts.MethodName != "TestService.Method" {
		t.Error("MethodName not set correctly")
	}
}
