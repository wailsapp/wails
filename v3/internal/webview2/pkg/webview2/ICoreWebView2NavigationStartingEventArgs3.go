//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2NavigationStartingEventArgs3Vtbl struct {
	ICoreWebView2NavigationStartingEventArgs2Vtbl
	GetNavigationKind ComProc
}

type ICoreWebView2NavigationStartingEventArgs3 struct {
	Vtbl *ICoreWebView2NavigationStartingEventArgs3Vtbl
}

func (i *ICoreWebView2NavigationStartingEventArgs3) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2NavigationStartingEventArgs3) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2NavigationStartingEventArgs3 queries the object for its ICoreWebView2NavigationStartingEventArgs3 interface. The receiver
// is the root of ICoreWebView2NavigationStartingEventArgs3's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2NavigationStartingEventArgs) GetICoreWebView2NavigationStartingEventArgs3() (*ICoreWebView2NavigationStartingEventArgs3, error) {
	var result *ICoreWebView2NavigationStartingEventArgs3

	iidICoreWebView2NavigationStartingEventArgs3 := NewGUID("{ddffe494-4942-4bd2-ab73-35b8ff40e19f}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2NavigationStartingEventArgs3)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2NavigationStartingEventArgs3) GetNavigationKind() (COREWEBVIEW2_NAVIGATION_KIND, error) {

	var value COREWEBVIEW2_NAVIGATION_KIND

	hr, _, _ := i.Vtbl.GetNavigationKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}
