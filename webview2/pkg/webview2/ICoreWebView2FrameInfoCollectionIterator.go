//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2FrameInfoCollectionIteratorVtbl struct {
	IUnknownVtbl
	GetHasCurrent ComProc
	GetCurrent ComProc
	MoveNext ComProc
}

type ICoreWebView2FrameInfoCollectionIterator struct {
	Vtbl *ICoreWebView2FrameInfoCollectionIteratorVtbl
}

func (i *ICoreWebView2FrameInfoCollectionIterator) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2FrameInfoCollectionIterator) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2FrameInfoCollectionIterator) GetHasCurrent() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, err := i.Vtbl.GetHasCurrent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    value := _value != 0
	return value, err
}

func (i *ICoreWebView2FrameInfoCollectionIterator) GetCurrent() (*ICoreWebView2FrameInfo, error) {

	var value *ICoreWebView2FrameInfo

	hr, _, err := i.Vtbl.GetCurrent.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2FrameInfoCollectionIterator) MoveNext() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, err := i.Vtbl.MoveNext.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    value := _value != 0
	return value, err
}
