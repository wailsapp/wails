//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2ExecuteScriptResultVtbl struct {
	IUnknownVtbl
	GetSucceeded         ComProc
	GetResultAsJson      ComProc
	TryGetResultAsString ComProc
	GetException         ComProc
}

type ICoreWebView2ExecuteScriptResult struct {
	Vtbl *ICoreWebView2ExecuteScriptResultVtbl
}

func (i *ICoreWebView2ExecuteScriptResult) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2ExecuteScriptResult) GetSucceeded() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetSucceeded.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	value := _value != 0
	return value, nil
}

func (i *ICoreWebView2ExecuteScriptResult) GetResultAsJson() (string, error) {
	// Create *uint16 to hold result
	var _jsonResult *uint16

	hr, _, _ := i.Vtbl.GetResultAsJson.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_jsonResult)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	jsonResult := UTF16PtrToString(_jsonResult)
	CoTaskMemFree(unsafe.Pointer(_jsonResult))
	return jsonResult, nil
}

func (i *ICoreWebView2ExecuteScriptResult) TryGetResultAsString() (string, bool, error) {
	// Create *uint16 to hold result
	var _stringResult *uint16
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.TryGetResultAsString.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_stringResult)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", false, syscall.Errno(hr)
	}
	// Get result and cleanup
	stringResult := UTF16PtrToString(_stringResult)
	CoTaskMemFree(unsafe.Pointer(_stringResult))
	// Get result and cleanup
	value := _value != 0
	return stringResult, value, nil
}

func (i *ICoreWebView2ExecuteScriptResult) GetException() (*ICoreWebView2ScriptException, error) {

	var exception *ICoreWebView2ScriptException

	hr, _, _ := i.Vtbl.GetException.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&exception)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return exception, nil
}
