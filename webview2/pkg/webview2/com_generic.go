//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

// COMCallOptions provides optional configuration for COM calls
type COMCallOptions struct {
	// IgnoreError ignores HRESULT errors (useful for QueryInterface)
	IgnoreError bool
}

// CallVoid executes a COM method that returns only HRESULT
func CallVoid(hr uintptr) error {
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

// CallValue executes a COM method that returns a value and HRESULT
func CallValue[T any](hr uintptr, value T) (T, error) {
	if windows.Handle(hr) != windows.S_OK {
		var zero T
		return zero, syscall.Errno(hr)
	}
	return value, nil
}

// CallBool executes a COM method that returns an int32 bool and HRESULT
func CallBool(hr uintptr, value int32) (bool, error) {
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	return value != 0, nil
}

// CallString executes a COM method that returns a string pointer and HRESULT
// It handles COM memory cleanup (CoTaskMemFree) automatically
func CallString(hr uintptr, value *uint16) (string, error) {
	if windows.Handle(hr) != windows.S_OK {
		if value != nil {
			CoTaskMemFree(unsafe.Pointer(value))
		}
		return "", syscall.Errno(hr)
	}
	result := UTF16PtrToString(value)
	CoTaskMemFree(unsafe.Pointer(value))
	return result, nil
}

// CallInterface executes a COM method that returns an interface pointer and HRESULT
func CallInterface[T any](hr uintptr, value *T) (*T, error) {
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, nil
}

// CallToken executes a COM method that returns an EventRegistrationToken and HRESULT
func CallToken(hr uintptr, token EventRegistrationToken) (EventRegistrationToken, error) {
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

// COMMethod represents a COM vtable method that can be called generically
type COMMethod interface {
	Call(...uintptr) (uintptr, uintptr, uintptr)
}

// InvokeVoid invokes a COM method with void return (HRESULT only)
func InvokeVoid(method COMMethod, args ...uintptr) error {
	hr, _, _ := method.Call(args...)
	return CallVoid(hr)
}

// InvokeValue invokes a COM method that returns a value and HRESULT
// The result pointer should be the last argument passed to the method
func InvokeValue[T any](method COMMethod, result *T, args ...uintptr) (T, error) {
	fullArgs := append(args, uintptr(unsafe.Pointer(result)))
	hr, _, _ := method.Call(fullArgs...)
	return CallValue(hr, *result)
}

// InvokeBool invokes a COM method that returns a bool (int32) and HRESULT
func InvokeBool(method COMMethod, result *int32, args ...uintptr) (bool, error) {
	fullArgs := append(args, uintptr(unsafe.Pointer(result)))
	hr, _, _ := method.Call(fullArgs...)
	return CallBool(hr, *result)
}

// InvokeString invokes a COM method that returns a string pointer and HRESULT
// Handles COM memory cleanup automatically
func InvokeString(method COMMethod, result **uint16, args ...uintptr) (string, error) {
	fullArgs := append(args, uintptr(unsafe.Pointer(result)))
	hr, _, _ := method.Call(fullArgs...)
	return CallString(hr, *result)
}

// InvokeInterface invokes a COM method that returns an interface pointer and HRESULT
func InvokeInterface[T any](method COMMethod, result **T, args ...uintptr) (*T, error) {
	fullArgs := append(args, uintptr(unsafe.Pointer(result)))
	hr, _, _ := method.Call(fullArgs...)
	return CallInterface(hr, *result)
}

// InvokeToken invokes a COM method that returns an EventRegistrationToken and HRESULT
func InvokeToken(method COMMethod, result *EventRegistrationToken, args ...uintptr) (EventRegistrationToken, error) {
	fullArgs := append(args, uintptr(unsafe.Pointer(result)))
	hr, _, _ := method.Call(fullArgs...)
	return CallToken(hr, *result)
}

// QueryInterface performs a generic QueryInterface call
func QueryInterface[T any](obj interface{}, iid *GUID) (*T, error) {
	var result *T

	switch v := obj.(type) {
	case *IUnknown:
		hr, _, _ := v.Vtbl.QueryInterface.Call(
			uintptr(unsafe.Pointer(v)),
			uintptr(unsafe.Pointer(iid)),
			uintptr(unsafe.Pointer(&result)),
		)
		return CallInterface(hr, &result)
	default:
		return nil, syscall.EINVAL
	}
}
