//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2BrowserProcessExitedEventArgsVtbl struct {
	IUnknownVtbl
	GetBrowserProcessExitKind ComProc
	GetBrowserProcessId ComProc
}

type ICoreWebView2BrowserProcessExitedEventArgs struct {
	Vtbl *ICoreWebView2BrowserProcessExitedEventArgsVtbl
}

func (i *ICoreWebView2BrowserProcessExitedEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2BrowserProcessExitedEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2BrowserProcessExitedEventArgs) GetBrowserProcessExitKind() (COREWEBVIEW2_BROWSER_PROCESS_EXIT_KIND, error) {

	var value COREWEBVIEW2_BROWSER_PROCESS_EXIT_KIND

	hr, _, _ := i.Vtbl.GetBrowserProcessExitKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2BrowserProcessExitedEventArgs) GetBrowserProcessId() (uint32, error) {

	var value uint32

	hr, _, _ := i.Vtbl.GetBrowserProcessId.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}
