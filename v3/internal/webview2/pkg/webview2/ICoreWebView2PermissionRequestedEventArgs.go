//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2PermissionRequestedEventArgsVtbl struct {
	IUnknownVtbl
	GetUri             ComProc
	GetPermissionKind  ComProc
	GetIsUserInitiated ComProc
	GetState           ComProc
	PutState           ComProc
	GetDeferral        ComProc
}

type ICoreWebView2PermissionRequestedEventArgs struct {
	Vtbl *ICoreWebView2PermissionRequestedEventArgsVtbl
}

func (i *ICoreWebView2PermissionRequestedEventArgs) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2PermissionRequestedEventArgs) GetUri() (string, error) {
	// Create *uint16 to hold result
	var _uri *uint16

	hr, _, _ := i.Vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	uri := UTF16PtrToString(_uri)
	CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, nil
}

func (i *ICoreWebView2PermissionRequestedEventArgs) GetPermissionKind() (COREWEBVIEW2_PERMISSION_KIND, error) {

	var permissionKind COREWEBVIEW2_PERMISSION_KIND

	hr, _, _ := i.Vtbl.GetPermissionKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&permissionKind)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return permissionKind, nil
}

func (i *ICoreWebView2PermissionRequestedEventArgs) GetIsUserInitiated() (bool, error) {
	// Create int32 to hold bool result
	var _isUserInitiated int32

	hr, _, _ := i.Vtbl.GetIsUserInitiated.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isUserInitiated)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	isUserInitiated := _isUserInitiated != 0
	return isUserInitiated, nil
}

func (i *ICoreWebView2PermissionRequestedEventArgs) GetState() (COREWEBVIEW2_PERMISSION_STATE, error) {

	var state COREWEBVIEW2_PERMISSION_STATE

	hr, _, _ := i.Vtbl.GetState.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&state)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return state, nil
}

func (i *ICoreWebView2PermissionRequestedEventArgs) PutState(state COREWEBVIEW2_PERMISSION_STATE) error {

	hr, _, _ := i.Vtbl.PutState.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(state),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2PermissionRequestedEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {

	var deferral *ICoreWebView2Deferral

	hr, _, _ := i.Vtbl.GetDeferral.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&deferral)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return deferral, nil
}
