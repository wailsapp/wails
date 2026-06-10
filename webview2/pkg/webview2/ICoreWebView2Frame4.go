//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Frame4Vtbl struct {
	ICoreWebView2Frame3Vtbl
	PostSharedBufferToScript ComProc
}

type ICoreWebView2Frame4 struct {
	Vtbl *ICoreWebView2Frame4Vtbl
}

func (i *ICoreWebView2Frame4) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Frame4) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Frame4 queries the object for its ICoreWebView2Frame4 interface. The receiver
// is the root of ICoreWebView2Frame4's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Frame) GetICoreWebView2Frame4() (*ICoreWebView2Frame4, error) {
	var result *ICoreWebView2Frame4

	iidICoreWebView2Frame4 := NewGUID("{188782dc-92aa-4732-ab3c-fcc59f6f68b9}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Frame4)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Frame4) PostSharedBufferToScript(sharedBuffer *ICoreWebView2SharedBuffer, access COREWEBVIEW2_SHARED_BUFFER_ACCESS, additionalDataAsJson string) error {

	// Convert string 'additionalDataAsJson' to *uint16
	_additionalDataAsJson, err := UTF16PtrFromString(additionalDataAsJson)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PostSharedBufferToScript.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(sharedBuffer)),
		uintptr(access),
		uintptr(unsafe.Pointer(_additionalDataAsJson)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
