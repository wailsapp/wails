//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_3Vtbl struct {
	IUnknownVtbl
	TrySuspend                          ComProc
	Resume                              ComProc
	GetIsSuspended                      ComProc
	SetVirtualHostNameToFolderMapping   ComProc
	ClearVirtualHostNameToFolderMapping ComProc
}

type ICoreWebView2_3 struct {
	Vtbl *ICoreWebView2_3Vtbl
}

func (i *ICoreWebView2_3) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_3() *ICoreWebView2_3 {
	var result *ICoreWebView2_3

	iidICoreWebView2_3 := NewGUID("{A0D6DF20-3B92-416D-AA0C-437A9C727857}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_3)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_3) TrySuspend(handler *ICoreWebView2TrySuspendCompletedHandler) error {

	hr, _, _ := i.Vtbl.TrySuspend.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_3) Resume() error {

	hr, _, _ := i.Vtbl.Resume.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_3) GetIsSuspended() (bool, error) {
	// Create int32 to hold bool result
	var _isSuspended int32

	hr, _, _ := i.Vtbl.GetIsSuspended.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_isSuspended)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	isSuspended := _isSuspended != 0
	return isSuspended, nil
}

func (i *ICoreWebView2_3) SetVirtualHostNameToFolderMapping(hostName string, folderPath string, accessKind COREWEBVIEW2_HOST_RESOURCE_ACCESS_KIND) error {

	// Convert string 'hostName' to *uint16
	_hostName, err := UTF16PtrFromString(hostName)
	if err != nil {
		return err
	}
	// Convert string 'folderPath' to *uint16
	_folderPath, err := UTF16PtrFromString(folderPath)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.SetVirtualHostNameToFolderMapping.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_hostName)),
		uintptr(unsafe.Pointer(_folderPath)),
		uintptr(accessKind),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_3) ClearVirtualHostNameToFolderMapping(hostName string) error {

	// Convert string 'hostName' to *uint16
	_hostName, err := UTF16PtrFromString(hostName)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.ClearVirtualHostNameToFolderMapping.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_hostName)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
