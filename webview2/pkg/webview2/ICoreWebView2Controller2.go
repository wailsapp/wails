//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Controller2Vtbl struct {
	ICoreWebView2ControllerVtbl
	GetDefaultBackgroundColor ComProc
	PutDefaultBackgroundColor ComProc
}

type ICoreWebView2Controller2 struct {
	Vtbl *ICoreWebView2Controller2Vtbl
}

func (i *ICoreWebView2Controller2) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Controller2) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Controller2 queries the object for its ICoreWebView2Controller2 interface. The receiver
// is the root of ICoreWebView2Controller2's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Controller) GetICoreWebView2Controller2() (*ICoreWebView2Controller2, error) {
	var result *ICoreWebView2Controller2

	iidICoreWebView2Controller2 := NewGUID("{c979903e-d4ca-4228-92eb-47ee3fa96eab}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Controller2)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Controller2) GetDefaultBackgroundColor() (COREWEBVIEW2_COLOR, error) {

	var value COREWEBVIEW2_COLOR

	hr, _, _ := i.Vtbl.GetDefaultBackgroundColor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return COREWEBVIEW2_COLOR{}, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2Controller2) PutDefaultBackgroundColor(value COREWEBVIEW2_COLOR) error {


	hr, _, _ := i.Vtbl.PutDefaultBackgroundColor.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(*(*uint32)(unsafe.Pointer(&value))),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
