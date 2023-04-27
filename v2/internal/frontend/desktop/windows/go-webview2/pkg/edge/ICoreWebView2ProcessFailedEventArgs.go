//go:build windows

package edge

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2ProcessFailedEventArgsVtbl struct {
	_IUnknownVtbl
	GetProcessFailedKind ComProc
}

type ICoreWebView2ProcessFailedEventArgs struct {
	vtbl *_ICoreWebView2ProcessFailedEventArgsVtbl
}

func (i *ICoreWebView2ProcessFailedEventArgs) GetProcessFailedKind() (COREWEBVIEW2_PROCESS_FAILED_KIND, error) {
	kind := COREWEBVIEW2_PROCESS_FAILED_KIND(0xffffffff)
	hr, _, err := i.vtbl.GetProcessFailedKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&kind)),
	)

	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}

	if kind == 0xffffffff {
		if err == nil {
			err = fmt.Errorf("unknown error")
		}
		return 0, err
	}

	return kind, nil
}
