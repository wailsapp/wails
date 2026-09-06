//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_16Vtbl struct {
	ICoreWebView2_15Vtbl
	Print ComProc
	ShowPrintUI ComProc
	PrintToPdfStream ComProc
}

type ICoreWebView2_16 struct {
	Vtbl *ICoreWebView2_16Vtbl
}

func (i *ICoreWebView2_16) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_16) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2_16 queries the object for its ICoreWebView2_16 interface. The receiver
// is the root of ICoreWebView2_16's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2) GetICoreWebView2_16() (*ICoreWebView2_16, error) {
	var result *ICoreWebView2_16

	iidICoreWebView2_16 := NewGUID("{0EB34DC9-9F91-41E1-8639-95CD5943906B}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_16)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_16) Print(printSettings *ICoreWebView2PrintSettings, handler *ICoreWebView2PrintCompletedHandler) error {


	hr, _, _ := i.Vtbl.Print.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(printSettings)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_16) ShowPrintUI(printDialogKind COREWEBVIEW2_PRINT_DIALOG_KIND) error {


	hr, _, _ := i.Vtbl.ShowPrintUI.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(printDialogKind),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_16) PrintToPdfStream(printSettings *ICoreWebView2PrintSettings, handler *ICoreWebView2PrintToPdfStreamCompletedHandler) error {


	hr, _, _ := i.Vtbl.PrintToPdfStream.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(printSettings)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
