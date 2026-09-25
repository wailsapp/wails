//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Environment7Vtbl struct {
	ICoreWebView2Environment6Vtbl
	GetUserDataFolder ComProc
}

type ICoreWebView2Environment7 struct {
	Vtbl *ICoreWebView2Environment7Vtbl
}

func (i *ICoreWebView2Environment7) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Environment7) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Environment7 queries the object for its ICoreWebView2Environment7 interface. The receiver
// is the root of ICoreWebView2Environment7's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Environment) GetICoreWebView2Environment7() (*ICoreWebView2Environment7, error) {
	var result *ICoreWebView2Environment7

	iidICoreWebView2Environment7 := NewGUID("{43c22296-3bbd-43a4-9c00-5c0df6dd29a2}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Environment7)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Environment7) GetUserDataFolder() (string, error) {
	// Create *uint16 to hold result
	var _value *uint16


	hr, _, _ := i.Vtbl.GetUserDataFolder.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	value := UTF16PtrToString(_value)
	CoTaskMemFree(unsafe.Pointer(_value))
	return value, nil
}
