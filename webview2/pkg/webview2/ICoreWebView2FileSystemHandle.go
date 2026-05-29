//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2FileSystemHandleVtbl struct {
	IUnknownVtbl
	GetKind ComProc
	GetPath ComProc
	GetPermission ComProc
}

type ICoreWebView2FileSystemHandle struct {
	Vtbl *ICoreWebView2FileSystemHandleVtbl
}

func (i *ICoreWebView2FileSystemHandle) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FileSystemHandle) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2FileSystemHandle) GetKind() (COREWEBVIEW2_FILE_SYSTEM_HANDLE_KIND, error) {

	var value COREWEBVIEW2_FILE_SYSTEM_HANDLE_KIND

	hr, _, _ := i.Vtbl.GetKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2FileSystemHandle) GetPath() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, _ := i.Vtbl.GetPath.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}

func (i *ICoreWebView2FileSystemHandle) GetPermission() (COREWEBVIEW2_FILE_SYSTEM_HANDLE_PERMISSION, error) {

	var value COREWEBVIEW2_FILE_SYSTEM_HANDLE_PERMISSION

	hr, _, _ := i.Vtbl.GetPermission.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}
