//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2NavigationCompletedEventArgs2Vtbl struct {
	ICoreWebView2NavigationCompletedEventArgsVtbl
	GetHttpStatusCode ComProc
}

type ICoreWebView2NavigationCompletedEventArgs2 struct {
	Vtbl *ICoreWebView2NavigationCompletedEventArgs2Vtbl
}

func (i *ICoreWebView2NavigationCompletedEventArgs2) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2NavigationCompletedEventArgs2) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2NavigationCompletedEventArgs2 queries the object for its ICoreWebView2NavigationCompletedEventArgs2 interface. The receiver
// is the root of ICoreWebView2NavigationCompletedEventArgs2's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2NavigationCompletedEventArgs) GetICoreWebView2NavigationCompletedEventArgs2() (*ICoreWebView2NavigationCompletedEventArgs2, error) {
	var result *ICoreWebView2NavigationCompletedEventArgs2

	iidICoreWebView2NavigationCompletedEventArgs2 := NewGUID("{fdf8b738-ee1e-4db2-a329-8d7d7b74d792}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2NavigationCompletedEventArgs2)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2NavigationCompletedEventArgs2) GetHttpStatusCode() (int, error) {

	var value int

	hr, _, _ := i.Vtbl.GetHttpStatusCode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}
