//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2ProcessExtendedInfoVtbl struct {
	IUnknownVtbl
	GetProcessInfo ComProc
	GetAssociatedFrameInfos ComProc
}

type ICoreWebView2ProcessExtendedInfo struct {
	Vtbl *ICoreWebView2ProcessExtendedInfoVtbl
}

func (i *ICoreWebView2ProcessExtendedInfo) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ProcessExtendedInfo) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2ProcessExtendedInfo) GetProcessInfo() (*ICoreWebView2ProcessInfo, error) {

	var processInfo *ICoreWebView2ProcessInfo

	hr, _, _ := i.Vtbl.GetProcessInfo.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&processInfo)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return processInfo, nil
}

func (i *ICoreWebView2ProcessExtendedInfo) GetAssociatedFrameInfos() (*ICoreWebView2FrameInfoCollection, error) {

	var frames *ICoreWebView2FrameInfoCollection

	hr, _, _ := i.Vtbl.GetAssociatedFrameInfos.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&frames)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return frames, nil
}
