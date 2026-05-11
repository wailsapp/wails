# Generic COM Call Helpers

This package provides type-safe generic helpers for COM method calls in the Webview2 binding. These helpers reduce boilerplate, improve type safety, and provide better IDE support.

## Overview

The traditional pattern for COM calls involves manual `unsafe.Pointer` handling and error checking:

```go
func (i *ICoreWebView2_9) GetDefaultDownloadDialogMargin() (POINT, error) {
	var value POINT
	hr, _, _ := i.Vtbl.GetDefaultDownloadDialogMargin.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return POINT{}, syscall.Errno(hr)
	}
	return value, nil
}
```

With generics, this becomes:

```go
func (i *ICoreWebView2_9) GetDefaultDownloadDialogMargin() (POINT, error) {
	var value POINT
	return InvokeValue(i.Vtbl.GetDefaultDownloadDialogMargin, &value, uintptr(unsafe.Pointer(i)))
}
```

## Available Helpers

### InvokeVoid

For methods that return only HRESULT (no output value):

```go
func (i *ICoreWebView2_9) OpenDefaultDownloadDialog() error {
	return InvokeVoid(i.Vtbl.OpenDefaultDownloadDialog, uintptr(unsafe.Pointer(i)))
}
```

### InvokeValue[T]

For getters that return a struct or primitive value:

```go
func (i *ICoreWebView2_9) GetDefaultDownloadDialogCornerAlignment() (COREWEBVIEW2_DEFAULT_DOWNLOAD_DIALOG_CORNER_ALIGNMENT, error) {
	var value COREWEBVIEW2_DEFAULT_DOWNLOAD_DIALOG_CORNER_ALIGNMENT
	return InvokeValue(i.Vtbl.GetDefaultDownloadDialogCornerAlignment, &value, uintptr(unsafe.Pointer(i)))
}
```

### InvokeBool

For getters that return a boolean value (stored as `int32` in COM):

```go
func (i *ICoreWebView2_9) GetIsDefaultDownloadDialogOpen() (bool, error) {
	var value int32
	return InvokeBool(i.Vtbl.GetIsDefaultDownloadDialogOpen, &value, uintptr(unsafe.Pointer(i)))
}
```

### InvokeString

For getters that return a string. **Automatically handles COM memory cleanup**:

```go
func (i *ICoreWebView2_15) GetFaviconUri() (string, error) {
	return InvokeString(i.Vtbl.GetFaviconUri, uintptr(unsafe.Pointer(i)))
}
```

**Important:** No need to call `CoTaskMemFree` manually - the helper does it for you.

### InvokeInterface[T]

For getters that return an interface pointer:

```go
func (i *ICoreWebView2_13) GetProfile() (*ICoreWebView2Profile, error) {
	return InvokeInterface(i.Vtbl.GetProfile, uintptr(unsafe.Pointer(i)))
}
```

### InvokeToken

For event registration methods that return an `EventRegistrationToken`:

```go
func (i *ICoreWebView2_9) AddIsDefaultDownloadDialogOpenChanged(handler *ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandler) (EventRegistrationToken, error) {
	return InvokeToken(i.Vtbl.AddIsDefaultDownloadDialogOpenChanged, handler, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(handler)))
}
```

## Benefits

### 1. Type Safety

The generic helpers ensure compile-time type safety. If you try to use the wrong type, the compiler will catch it:

```go
var value int32
// This won't compile if method returns a struct:
_, _ = InvokeValue(vtbl.GetStruct, &value, thisPtr)
```

### 2. Automatic Memory Management

String getters no longer require manual memory cleanup:

```go
// Old pattern - easy to forget cleanup
var _value *uint16
hr, _, _ := vtbl.GetString.Call(...)
if hr != S_OK {
	return "", error
}
str := UTF16PtrToString(_value)
CoTaskMemFree(unsafe.Pointer(_value)) // Easy to forget!
return str, nil

// New pattern - automatic cleanup
return InvokeString(vtbl.GetString, thisPtr) // CoTaskMemFree called automatically
```

### 3. Consistent Error Handling

All helpers use the same error handling pattern, reducing bugs from inconsistent HRESULT checking:

```go
// All helpers follow the same pattern:
if windows.Handle(hr) != windows.S_OK {
	return zeroValue, syscall.Errno(hr)
}
```

### 4. Reduced Boilerplate

Less repetitive code means fewer bugs and easier maintenance:

```go
// Old: 11 lines
func (i *ICoreWebView2_9) OpenDefaultDownloadDialog() error {
	hr, _, _ := i.Vtbl.OpenDefaultDownloadDialog.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

// New: 2 lines
func (i *ICoreWebView2_9) OpenDefaultDownloadDialog() error {
	return InvokeVoid(i.Vtbl.OpenDefaultDownloadDialog, uintptr(unsafe.Pointer(i)))
}
```

### 5. Better IDE Support

Generic functions provide better IDE autocompletion and type information:

- IDEs can infer return types from the generic parameter
- Type hints are explicit in the function signature
- Go-to-definition works better with fewer indirections

## Performance

The generic helpers have **zero runtime overhead** compared to the manual pattern. The compiler inlines the generic code, so the generated assembly is identical.

See `com_generic_test.go` for benchmarks that verify this.

## Migration Guide for the Generator

When updating the IDL-to-Go generator to use these helpers:

### Step 1: Determine the return type pattern

Check the IDL method signature:
- `void Method()` → Use `InvokeVoid`
- `HRESULT Method([out] T* result)` → Use `InvokeValue[T]`
- `HRESULT Method([out] BOOL* result)` → Use `InvokeBool`
- `HRESULT Method([out, retval] BSTR* result)` → Use `InvokeString`
- `HRESULT Method([out, retval] IInterface** result)` → Use `InvokeInterface[IInterface]`
- Event registration → Use `InvokeToken`

### Step 2: Generate the method signature

The method signature remains the same as before - only the body changes.

### Step 3: Generate the body using the appropriate helper

```go
// For InvokeVoid:
return InvokeVoid(i.Vtbl.MethodName, uintptr(unsafe.Pointer(i)))

// For InvokeValue:
var value ValueType
return InvokeValue(i.Vtbl.MethodName, &value, uintptr(unsafe.Pointer(i)))

// For InvokeBool:
var value int32
return InvokeBool(i.Vtbl.MethodName, &value, uintptr(unsafe.Pointer(i)))

// For InvokeString:
return InvokeString(i.Vtbl.MethodName, uintptr(unsafe.Pointer(i)))

// For InvokeInterface:
return InvokeInterface(i.Vtbl.MethodName, uintptr(unsafe.Pointer(i)))

// For InvokeToken:
return InvokeToken(i.Vtbl.MethodName, &token, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(handler)))
```

### Step 4: For setters (Put methods)

Setters don't need the value pointer argument pattern:

```go
// Old
func (i *ICoreWebView2_9) PutDefaultDownloadDialogCornerAlignment(value COREWEBVIEW2_DEFAULT_DOWNLOAD_DIALOG_CORNER_ALIGNMENT) error {
	hr, _, _ := i.Vtbl.PutDefaultDownloadDialogCornerAlignment.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

// New
func (i *ICoreWebView2_9) PutDefaultDownloadDialogCornerAlignment(value COREWEBVIEW2_DEFAULT_DOWNLOAD_DIALOG_CORNER_ALIGNMENT) error {
	return InvokeVoid(i.Vtbl.PutDefaultDownloadDialogCornerAlignment, uintptr(unsafe.Pointer(i)), uintptr(value))
}
```

## Examples

See `com_generic_examples.go` for comprehensive examples showing:
- Old vs. new patterns for each helper type
- Full method implementations
- Migration patterns for the generator

## Testing

Run tests to verify type safety:

```bash
cd webview2/pkg/webview2
go test -v -run "TestInvoke"
```

Run benchmarks to verify performance:

```bash
go test -bench="BenchmarkInvoke" -benchmem
```

## Requirements

- Go 1.18+ (for generics support)
- Windows (COM is Windows-specific)

## See Also

- `com_generic.go` - Generic helper implementations
- `com_generic_examples.go` - Usage examples and migration guide
- `com_generic_test.go` - Tests and benchmarks
- `WEBVIEW2_ARCHITECTURE.md` - Architecture decision record
