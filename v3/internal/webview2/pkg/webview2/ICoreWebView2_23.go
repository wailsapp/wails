//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_23Vtbl struct {
	ICoreWebView2_22Vtbl
	PostWebMessageAsJsonWithAdditionalObjects ComProc
}

type ICoreWebView2_23 struct {
	Vtbl *ICoreWebView2_23Vtbl
}

func (i *ICoreWebView2_23) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_23) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2_23 queries the object for its ICoreWebView2_23 interface. The receiver
// is the root of ICoreWebView2_23's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2) GetICoreWebView2_23() (*ICoreWebView2_23, error) {
	var result *ICoreWebView2_23

	iidICoreWebView2_23 := NewGUID("{508f0db5-90c4-5872-90a7-267a91377502}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_23)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_23) PostWebMessageAsJsonWithAdditionalObjects(webMessageAsJson string, additionalObjects *ICoreWebView2ObjectCollectionView) error {

	// Convert string 'webMessageAsJson' to *uint16
	_webMessageAsJson, err := UTF16PtrFromString(webMessageAsJson)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.PostWebMessageAsJsonWithAdditionalObjects.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_webMessageAsJson)),
		uintptr(unsafe.Pointer(additionalObjects)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
