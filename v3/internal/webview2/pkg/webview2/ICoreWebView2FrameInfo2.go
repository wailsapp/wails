//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2FrameInfo2Vtbl struct {
	IUnknownVtbl
	GetParentFrameInfo ComProc
	GetFrameId         ComProc
	GetFrameKind       ComProc
}

type ICoreWebView2FrameInfo2 struct {
	Vtbl *ICoreWebView2FrameInfo2Vtbl
}

func (i *ICoreWebView2FrameInfo2) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2FrameInfo2() *ICoreWebView2FrameInfo2 {
	var result *ICoreWebView2FrameInfo2

	iidICoreWebView2FrameInfo2 := NewGUID("{56f85cfa-72c4-11ee-b962-0242ac120002}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2FrameInfo2)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2FrameInfo2) GetParentFrameInfo() (*ICoreWebView2FrameInfo, error) {

	var frameInfo *ICoreWebView2FrameInfo

	hr, _, _ := i.Vtbl.GetParentFrameInfo.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&frameInfo)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return frameInfo, nil
}

func (i *ICoreWebView2FrameInfo2) GetFrameId() (uint32, error) {

	var id uint32

	hr, _, _ := i.Vtbl.GetFrameId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&id)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return id, nil
}

func (i *ICoreWebView2FrameInfo2) GetFrameKind() (COREWEBVIEW2_FRAME_KIND, error) {

	var kind COREWEBVIEW2_FRAME_KIND

	hr, _, _ := i.Vtbl.GetFrameKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&kind)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return kind, nil
}
