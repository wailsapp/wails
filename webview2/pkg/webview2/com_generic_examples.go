//go:build windows

package webview2

// This file demonstrates the new generic COM call patterns
// Compare these with the existing implementations in ICoreWebView2_9.go

// OLD PATTERN (with unsafe pointers and manual error handling):
//
// func (i *ICoreWebView2_9) OpenDefaultDownloadDialog() error {
// 	hr, _, _ := i.Vtbl.OpenDefaultDownloadDialog.Call(
// 		uintptr(unsafe.Pointer(i)),
// 	)
// 	if windows.Handle(hr) != windows.S_OK {
// 		return syscall.Errno(hr)
// 	}
// 	return nil
// }

// NEW PATTERN (with generic InvokeVoid):
//
// func (i *ICoreWebView2_9) OpenDefaultDownloadDialog() error {
// 	return InvokeVoid(i.Vtbl.OpenDefaultDownloadDialog, uintptr(unsafe.Pointer(i)))
// }

// OLD PATTERN (getter for struct):
//
// func (i *ICoreWebView2_9) GetDefaultDownloadDialogMargin() (POINT, error) {
// 	var value POINT
// 	hr, _, _ := i.Vtbl.GetDefaultDownloadDialogMargin.Call(
// 		uintptr(unsafe.Pointer(i)),
// 		uintptr(unsafe.Pointer(&value)),
// 	)
// 	if windows.Handle(hr) != windows.S_OK {
// 		return POINT{}, syscall.Errno(hr)
// 	}
// 	return value, nil
// }

// NEW PATTERN (with generic InvokeValue):
//
// func (i *ICoreWebView2_9) GetDefaultDownloadDialogMargin() (POINT, error) {
// 	var value POINT
// 	return InvokeValue(i.Vtbl.GetDefaultDownloadDialogMargin, &value, uintptr(unsafe.Pointer(i)))
// }

// OLD PATTERN (getter for bool):
//
// func (i *ICoreWebView2_9) GetIsDefaultDownloadDialogOpen() (bool, error) {
// 	var _value int32
// 	hr, _, _ := i.Vtbl.GetIsDefaultDownloadDialogOpen.Call(
// 		uintptr(unsafe.Pointer(i)),
// 		uintptr(unsafe.Pointer(&_value)),
// 	)
// 	if windows.Handle(hr) != windows.S_OK {
// 		return false, syscall.Errno(hr)
// 	}
// 	value := _value != 0
// 	return value, nil
// }

// NEW PATTERN (with generic InvokeBool):
//
// func (i *ICoreWebView2_9) GetIsDefaultDownloadDialogOpen() (bool, error) {
// 	var _value int32
// 	return InvokeBool(i.Vtbl.GetIsDefaultDownloadDialogOpen, &_value, uintptr(unsafe.Pointer(i)))
// }

// OLD PATTERN (getter for interface):
//
// func (i *ICoreWebView2_13) GetProfile() (*ICoreWebView2Profile, error) {
// 	var value *ICoreWebView2Profile
// 	hr, _, _ := i.Vtbl.GetProfile.Call(
// 		uintptr(unsafe.Pointer(i)),
// 		uintptr(unsafe.Pointer(&value)),
// 	)
// 	if windows.Handle(hr) != windows.S_OK {
// 		return nil, syscall.Errno(hr)
// 	}
// 	return value, nil
// }

// NEW PATTERN (with generic InvokeInterface):
//
// func (i *ICoreWebView2_13) GetProfile() (*ICoreWebView2Profile, error) {
// 	return InvokeInterface(i.Vtbl.GetProfile, uintptr(unsafe.Pointer(i)))
// }

// OLD PATTERN (getter for string):
//
// func (i *ICoreWebView2_15) GetFaviconUri() (string, error) {
// 	var _value *uint16
// 	hr, _, _ := i.Vtbl.GetFaviconUri.Call(
// 		uintptr(unsafe.Pointer(i)),
// 		uintptr(unsafe.Pointer(_value)),
// 	)
// 	if windows.Handle(hr) != windows.S_OK {
// 		return "", syscall.Errno(hr)
// 	}
// 	value := UTF16PtrToString(_value)
// 	CoTaskMemFree(unsafe.Pointer(_value))
// 	return value, nil
// }

// NEW PATTERN (with generic InvokeString):
//
// func (i *ICoreWebView2_15) GetFaviconUri() (string, error) {
// 	return InvokeString(i.Vtbl.GetFaviconUri, uintptr(unsafe.Pointer(i)))
// }

// OLD PATTERN (event registration):
//
// func (i *ICoreWebView2_9) AddIsDefaultDownloadDialogOpenChanged(handler *ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandler) (EventRegistrationToken, error) {
// 	var token EventRegistrationToken
// 	hr, _, _ := i.Vtbl.AddIsDefaultDownloadDialogOpenChanged.Call(
// 		uintptr(unsafe.Pointer(i)),
// 		uintptr(unsafe.Pointer(handler)),
// 		uintptr(unsafe.Pointer(&token)),
// 	)
// 	if windows.Handle(hr) != windows.S_OK {
// 		return EventRegistrationToken{}, syscall.Errno(hr)
// 	}
// 	return token, nil
// }

// NEW PATTERN (with generic InvokeToken):
//
// func (i *ICoreWebView2_9) AddIsDefaultDownloadDialogOpenChanged(handler *ICoreWebView2IsDefaultDownloadDialogOpenChangedEventHandler) (EventRegistrationToken, error) {
// 	return InvokeToken(i.Vtbl.AddIsDefaultDownloadDialogOpenChanged, handler, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(handler)))
// }

// Example of using the new generic helpers in existing methods:
// This demonstrates how to migrate existing methods to use the new patterns
// The generator should be updated to output code in the new pattern

/*
MIGRATION GUIDE FOR THE GENERATOR:
====================================

When generating interface methods, use the following patterns:

1. For void methods (returns error only):
   OLD: if hr != S_OK { return error }
   NEW: return InvokeVoid(vtbl.Method, thisPtr)

2. For getters returning struct:
   OLD: var value T; if hr != S_OK { return T{}, error }; return value, nil
   NEW: var value T; return InvokeValue(vtbl.Method, &value, thisPtr)

3. For getters returning bool:
   OLD: var value int32; if hr != S_OK { return false, error }; return value != 0, nil
   NEW: var value int32; return InvokeBool(vtbl.Method, &value, thisPtr)

4. For getters returning interface:
   OLD: var value *T; if hr != S_OK { return nil, error }; return value, nil
   NEW: return InvokeInterface(vtbl.Method, thisPtr)

5. For getters returning string:
   OLD: var value *uint16; if hr != S_OK { return "", error }; str = UTF16PtrToString(value); CoTaskMemFree(value); return str, nil
   NEW: return InvokeString(vtbl.Method, thisPtr)

6. For event registration (returns token):
   OLD: var token EventRegistrationToken; if hr != S_OK { return EventRegistrationToken{}, error }; return token, nil
   NEW: return InvokeToken(vtbl.Method, &token, thisPtr, handlerPtr)

BENEFITS:
- Less boilerplate code
- Consistent error handling
- Automatic memory management for strings
- Better type safety with generics
- Improved IDE support with explicit return types
*/
