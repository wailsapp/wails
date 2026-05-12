//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2ProcessFailedEventArgs2Vtbl struct {
	IUnknownVtbl
	GetReason ComProc
	GetExitCode ComProc
	GetProcessDescription ComProc
	GetFrameInfosForFailedProcess ComProc
}

type ICoreWebView2ProcessFailedEventArgs2 struct {
	Vtbl *ICoreWebView2ProcessFailedEventArgs2Vtbl
}

func (i *ICoreWebView2ProcessFailedEventArgs2) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2ProcessFailedEventArgs2) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2) GetICoreWebView2ProcessFailedEventArgs2() (*ICoreWebView2ProcessFailedEventArgs2, error) {
	var result *ICoreWebView2ProcessFailedEventArgs2

	iidICoreWebView2ProcessFailedEventArgs2 := NewGUID("{4dab9422-46fa-4c3e-a5d2-41d2071d3680}")
	hr, _, _ := i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2ProcessFailedEventArgs2)),
		uintptr(unsafe.Pointer(&result)))
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return result, nil
}


func (i *ICoreWebView2ProcessFailedEventArgs2) GetReason() (COREWEBVIEW2_PROCESS_FAILED_REASON, error) {

	var reason COREWEBVIEW2_PROCESS_FAILED_REASON

	hr, _, err := i.Vtbl.GetReason.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&reason)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return reason, err
}

func (i *ICoreWebView2ProcessFailedEventArgs2) GetExitCode() (int, error) {

	var exitCode int

	hr, _, err := i.Vtbl.GetExitCode.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&exitCode)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return exitCode, err
}

func (i *ICoreWebView2ProcessFailedEventArgs2) GetProcessDescription() (string, error) {
	// Create *uint16 to hold result
	var _processDescription *uint16


	hr, _, err := i.Vtbl.GetProcessDescription.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_processDescription)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	processDescription := UTF16PtrToString(_processDescription)
	CoTaskMemFree(unsafe.Pointer(_processDescription))
	return processDescription, err
}

func (i *ICoreWebView2ProcessFailedEventArgs2) GetFrameInfosForFailedProcess() (*ICoreWebView2FrameInfoCollection, error) {

	var frames *ICoreWebView2FrameInfoCollection

	hr, _, err := i.Vtbl.GetFrameInfosForFailedProcess.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&frames)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return frames, err
}
