//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Settings7Vtbl struct {
	IUnknownVtbl
	GetHiddenPdfToolbarItems ComProc
	PutHiddenPdfToolbarItems ComProc
}

type ICoreWebView2Settings7 struct {
	Vtbl *ICoreWebView2Settings7Vtbl
}

func (i *ICoreWebView2Settings7) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Settings7) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2Settings7() (*ICoreWebView2Settings7, error) {
	var result *ICoreWebView2Settings7

	iidICoreWebView2Settings7 := NewGUID("{488dc902-35ef-42d2-bc7d-94b65c4bc49c}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Settings7)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Settings7) GetHiddenPdfToolbarItems() (COREWEBVIEW2_PDF_TOOLBAR_ITEMS, error) {

	var value COREWEBVIEW2_PDF_TOOLBAR_ITEMS

	hr, _, _ := i.Vtbl.GetHiddenPdfToolbarItems.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2Settings7) PutHiddenPdfToolbarItems(value COREWEBVIEW2_PDF_TOOLBAR_ITEMS) error {


	hr, _, _ := i.Vtbl.PutHiddenPdfToolbarItems.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
