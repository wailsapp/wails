//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_20Vtbl struct {
	ICoreWebView2_19Vtbl
	GetFrameId ComProc
}

type ICoreWebView2_20 struct {
	Vtbl *ICoreWebView2_20Vtbl
}

func (i *ICoreWebView2_20) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_20) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2_20 queries the object for its ICoreWebView2_20 interface. The receiver
// is the root of ICoreWebView2_20's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2) GetICoreWebView2_20() (*ICoreWebView2_20, error) {
	var result *ICoreWebView2_20

	iidICoreWebView2_20 := NewGUID("{b4bc1926-7305-11ee-b962-0242ac120002}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_20)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_20) GetFrameId() (uint32, error) {

	var value uint32

	hr, _, _ := i.Vtbl.GetFrameId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}
