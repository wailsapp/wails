//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Profile3Vtbl struct {
	ICoreWebView2Profile2Vtbl
	GetPreferredTrackingPreventionLevel ComProc
	PutPreferredTrackingPreventionLevel ComProc
}

type ICoreWebView2Profile3 struct {
	Vtbl *ICoreWebView2Profile3Vtbl
}

func (i *ICoreWebView2Profile3) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Profile3) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Profile3 queries the object for its ICoreWebView2Profile3 interface. The receiver
// is the root of ICoreWebView2Profile3's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Profile) GetICoreWebView2Profile3() (*ICoreWebView2Profile3, error) {
	var result *ICoreWebView2Profile3

	iidICoreWebView2Profile3 := NewGUID("{b188e659-5685-4e05-bdba-fc640e0f1992}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Profile3)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Profile3) GetPreferredTrackingPreventionLevel() (COREWEBVIEW2_TRACKING_PREVENTION_LEVEL, error) {

	var value COREWEBVIEW2_TRACKING_PREVENTION_LEVEL

	hr, _, _ := i.Vtbl.GetPreferredTrackingPreventionLevel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2Profile3) PutPreferredTrackingPreventionLevel(value COREWEBVIEW2_TRACKING_PREVENTION_LEVEL) error {


	hr, _, _ := i.Vtbl.PutPreferredTrackingPreventionLevel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
