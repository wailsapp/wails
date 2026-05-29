//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2_19Vtbl struct {
	IUnknownVtbl
	GetMemoryUsageTargetLevel ComProc
	PutMemoryUsageTargetLevel ComProc
}

type ICoreWebView2_19 struct {
	Vtbl *ICoreWebView2_19Vtbl
}

func (i *ICoreWebView2_19) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2_19) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2_19() (*ICoreWebView2_19, error) {
	var result *ICoreWebView2_19

	iidICoreWebView2_19 := NewGUID("{6921f954-79b0-437f-a997-c85811897c68}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_19)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2_19) GetMemoryUsageTargetLevel() (COREWEBVIEW2_MEMORY_USAGE_TARGET_LEVEL, error) {

	var value COREWEBVIEW2_MEMORY_USAGE_TARGET_LEVEL

	hr, _, _ := i.Vtbl.GetMemoryUsageTargetLevel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2_19) PutMemoryUsageTargetLevel(value COREWEBVIEW2_MEMORY_USAGE_TARGET_LEVEL) error {


	hr, _, _ := i.Vtbl.PutMemoryUsageTargetLevel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
