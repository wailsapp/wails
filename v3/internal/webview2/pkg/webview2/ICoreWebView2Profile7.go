//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2Profile7Vtbl struct {
	ICoreWebView2Profile6Vtbl
	AddBrowserExtension ComProc
	GetBrowserExtensions ComProc
}

type ICoreWebView2Profile7 struct {
	Vtbl *ICoreWebView2Profile7Vtbl
}

func (i *ICoreWebView2Profile7) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Profile7) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


// GetICoreWebView2Profile7 queries the object for its ICoreWebView2Profile7 interface. The receiver
// is the root of ICoreWebView2Profile7's inheritance chain — the object that actually
// implements it.
func (i *ICoreWebView2Profile) GetICoreWebView2Profile7() (*ICoreWebView2Profile7, error) {
	var result *ICoreWebView2Profile7

	iidICoreWebView2Profile7 := NewGUID("{7b4c7906-a1aa-4cb4-b723-db09f813d541}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2Profile7)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2Profile7) AddBrowserExtension(extensionFolderPath string, handler *ICoreWebView2ProfileAddBrowserExtensionCompletedHandler) error {

	// Convert string 'extensionFolderPath' to *uint16
	_extensionFolderPath, err := UTF16PtrFromString(extensionFolderPath)
	if err != nil {
		return err
	}

	hr, _, _ := i.Vtbl.AddBrowserExtension.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_extensionFolderPath)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Profile7) GetBrowserExtensions(handler *ICoreWebView2ProfileGetBrowserExtensionsCompletedHandler) error {


	hr, _, _ := i.Vtbl.GetBrowserExtensions.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
