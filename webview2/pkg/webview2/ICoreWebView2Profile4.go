//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Profile4Vtbl struct {
	IUnknownVtbl
	SetPermissionState ComProc
	GetNonDefaultPermissionSettings ComProc
}

type ICoreWebView2Profile4 struct {
	Vtbl *ICoreWebView2Profile4Vtbl
}

func (i *ICoreWebView2Profile4) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Profile4) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2Profile4() (*ICoreWebView2Profile4, error) {
	var result *ICoreWebView2Profile4

	iidICoreWebView2Profile4 := NewGUID("{8f4ae680-192e-4ec8-833a-21cfadaef628}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Profile4)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Profile4) SetPermissionState(PermissionKind COREWEBVIEW2_PERMISSION_KIND, origin string, State COREWEBVIEW2_PERMISSION_STATE, handler *ICoreWebView2SetPermissionStateCompletedHandler) error {

	// Convert string 'origin' to *uint16
	_origin, err := UTF16PtrFromString(origin)
	if err != nil {
		return err
	}

	hr, _, err := i.Vtbl.SetPermissionState.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(PermissionKind),
		uintptr(unsafe.Pointer(_origin)),
		uintptr(State),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2Profile4) GetNonDefaultPermissionSettings(handler *ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandler) error {


	hr, _, err := i.Vtbl.GetNonDefaultPermissionSettings.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
