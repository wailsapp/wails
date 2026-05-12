//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2DownloadStartingEventArgsVtbl struct {
	IUnknownVtbl
	GetDownloadOperation ComProc
	GetCancel ComProc
	PutCancel ComProc
	GetResultFilePath ComProc
	PutResultFilePath ComProc
	GetHandled ComProc
	PutHandled ComProc
	GetDeferral ComProc
}

type ICoreWebView2DownloadStartingEventArgs struct {
	Vtbl *ICoreWebView2DownloadStartingEventArgsVtbl
}

func (i *ICoreWebView2DownloadStartingEventArgs) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2DownloadStartingEventArgs) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2DownloadStartingEventArgs) GetDownloadOperation() (*ICoreWebView2DownloadOperation, error) {

	var downloadOperation *ICoreWebView2DownloadOperation

	hr, _, err := i.Vtbl.GetDownloadOperation.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&downloadOperation)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return downloadOperation, err
}

func (i *ICoreWebView2DownloadStartingEventArgs) GetCancel() (bool, error) {
	// Create int32 to hold bool result
	var _cancel int32

	hr, _, err := i.Vtbl.GetCancel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_cancel)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    cancel := _cancel != 0
	return cancel, err
}

func (i *ICoreWebView2DownloadStartingEventArgs) PutCancel(cancel bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _cancel int32
	if cancel {
		_cancel = 1
	}

	hr, _, err := i.Vtbl.PutCancel.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_cancel),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2DownloadStartingEventArgs) GetResultFilePath() (string, error) {
	// Create *uint16 to hold result
	var _resultFilePath *uint16


	hr, _, err := i.Vtbl.GetResultFilePath.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_resultFilePath)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	resultFilePath := UTF16PtrToString(_resultFilePath)
	CoTaskMemFree(unsafe.Pointer(_resultFilePath))
	return resultFilePath, err
}

func (i *ICoreWebView2DownloadStartingEventArgs) PutResultFilePath(resultFilePath string) error {

	// Convert string 'resultFilePath' to *uint16
	_resultFilePath, err := UTF16PtrFromString(resultFilePath)
	if err != nil {
		return err
	}

	hr, _, err := i.Vtbl.PutResultFilePath.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_resultFilePath)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2DownloadStartingEventArgs) GetHandled() (bool, error) {
	// Create int32 to hold bool result
	var _handled int32

	hr, _, err := i.Vtbl.GetHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_handled)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
    handled := _handled != 0
	return handled, err
}

func (i *ICoreWebView2DownloadStartingEventArgs) PutHandled(handled bool) error {

	// Convert Go bool to COM BOOL (int32)
	var _handled int32
	if handled {
		_handled = 1
	}

	hr, _, err := i.Vtbl.PutHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(_handled),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2DownloadStartingEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {

	var deferral *ICoreWebView2Deferral

	hr, _, err := i.Vtbl.GetDeferral.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&deferral)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return nil, syscall.Errno(hr)
	}
	return deferral, err
}
