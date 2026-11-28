//go:build windows

package edge

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// ICoreWebView2NavigationStartingEventArgs 由 go-webview2 的生成器补齐。
// Wails 上游没有生成这个类型，而它正是「页面要跳到哪里」的唯一来源。
type _ICoreWebView2NavigationStartingEventArgsVtbl struct {
	_IUnknownVtbl
	GetUri             ComProc
	GetIsUserInitiated ComProc
	GetIsRedirected    ComProc
	GetRequestHeaders  ComProc
	GetCancel          ComProc
	PutCancel          ComProc
	GetNavigationId    ComProc
}

type ICoreWebView2NavigationStartingEventArgs struct {
	vtbl *_ICoreWebView2NavigationStartingEventArgsVtbl
}

func (i *ICoreWebView2NavigationStartingEventArgs) AddRef() uint32 {
	ret, _, _ := i.vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(ret)
}

func (i *ICoreWebView2NavigationStartingEventArgs) Release() uint32 {
	ret, _, _ := i.vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(ret)
}

// GetUri 返回即将导航到的地址。
func (i *ICoreWebView2NavigationStartingEventArgs) GetUri() (string, error) {
	var uri *uint16
	hr, _, _ := i.vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&uri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", windows.Errno(hr)
	}
	text := windows.UTF16PtrToString(uri)
	windows.CoTaskMemFree(unsafe.Pointer(uri))
	return text, nil
}

// GetIsUserInitiated 区分用户点击与脚本跳转。
func (i *ICoreWebView2NavigationStartingEventArgs) GetIsUserInitiated() (bool, error) {
	var result int32
	hr, _, _ := i.vtbl.GetIsUserInitiated.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&result)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return result != 0, nil
}

// GetIsRedirected 报告这次导航是否由服务端或客户端重定向产生。
func (i *ICoreWebView2NavigationStartingEventArgs) GetIsRedirected() (bool, error) {
	var result int32
	hr, _, _ := i.vtbl.GetIsRedirected.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&result)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return result != 0, nil
}

// GetCancel 报告宿主是否已经请求取消这次导航。
func (i *ICoreWebView2NavigationStartingEventArgs) GetCancel() (bool, error) {
	var result int32
	hr, _, _ := i.vtbl.GetCancel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&result)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, windows.Errno(hr)
	}
	return result != 0, nil
}

// PutCancel 设为 true 即取消这次导航。
func (i *ICoreWebView2NavigationStartingEventArgs) PutCancel(value bool) error {
	var flag int32
	if value {
		flag = 1
	}
	hr, _, _ := i.vtbl.PutCancel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(flag),
	)
	if windows.Handle(hr) != windows.S_OK {
		return windows.Errno(hr)
	}
	return nil
}
