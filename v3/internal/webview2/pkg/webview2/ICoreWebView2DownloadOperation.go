//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2DownloadOperationVtbl struct {
	IUnknownVtbl
	AddBytesReceivedChanged       ComProc
	RemoveBytesReceivedChanged    ComProc
	AddEstimatedEndTimeChanged    ComProc
	RemoveEstimatedEndTimeChanged ComProc
	AddStateChanged               ComProc
	RemoveStateChanged            ComProc
	GetUri                        ComProc
	GetContentDisposition         ComProc
	GetMimeType                   ComProc
	GetTotalBytesToReceive        ComProc
	GetBytesReceived              ComProc
	GetEstimatedEndTime           ComProc
	GetResultFilePath             ComProc
	GetState                      ComProc
	GetInterruptReason            ComProc
	Cancel                        ComProc
	Pause                         ComProc
	Resume                        ComProc
	GetCanResume                  ComProc
}

type ICoreWebView2DownloadOperation struct {
	Vtbl *ICoreWebView2DownloadOperationVtbl
}

func (i *ICoreWebView2DownloadOperation) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2DownloadOperation) AddBytesReceivedChanged(eventHandler *ICoreWebView2BytesReceivedChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddBytesReceivedChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2DownloadOperation) RemoveBytesReceivedChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveBytesReceivedChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2DownloadOperation) AddEstimatedEndTimeChanged(eventHandler *ICoreWebView2EstimatedEndTimeChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddEstimatedEndTimeChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2DownloadOperation) RemoveEstimatedEndTimeChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveEstimatedEndTimeChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2DownloadOperation) AddStateChanged(eventHandler *ICoreWebView2StateChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddStateChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2DownloadOperation) RemoveStateChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveStateChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2DownloadOperation) GetUri() (string, error) {
	// Create *uint16 to hold result
	var _uri *uint16

	hr, _, _ := i.Vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_uri)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	uri := UTF16PtrToString(_uri)
	CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, nil
}

func (i *ICoreWebView2DownloadOperation) GetContentDisposition() (string, error) {
	// Create *uint16 to hold result
	var _contentDisposition *uint16

	hr, _, _ := i.Vtbl.GetContentDisposition.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_contentDisposition)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	contentDisposition := UTF16PtrToString(_contentDisposition)
	CoTaskMemFree(unsafe.Pointer(_contentDisposition))
	return contentDisposition, nil
}

func (i *ICoreWebView2DownloadOperation) GetMimeType() (string, error) {
	// Create *uint16 to hold result
	var _mimeType *uint16

	hr, _, _ := i.Vtbl.GetMimeType.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_mimeType)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	mimeType := UTF16PtrToString(_mimeType)
	CoTaskMemFree(unsafe.Pointer(_mimeType))
	return mimeType, nil
}

func (i *ICoreWebView2DownloadOperation) GetTotalBytesToReceive() (int64, error) {

	var totalBytesToReceive int64

	hr, _, _ := i.Vtbl.GetTotalBytesToReceive.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&totalBytesToReceive)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return totalBytesToReceive, nil
}

func (i *ICoreWebView2DownloadOperation) GetBytesReceived() (int64, error) {

	var bytesReceived int64

	hr, _, _ := i.Vtbl.GetBytesReceived.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&bytesReceived)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return bytesReceived, nil
}

func (i *ICoreWebView2DownloadOperation) GetEstimatedEndTime() (string, error) {
	// Create *uint16 to hold result
	var _estimatedEndTime *uint16

	hr, _, _ := i.Vtbl.GetEstimatedEndTime.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_estimatedEndTime)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	estimatedEndTime := UTF16PtrToString(_estimatedEndTime)
	CoTaskMemFree(unsafe.Pointer(_estimatedEndTime))
	return estimatedEndTime, nil
}

func (i *ICoreWebView2DownloadOperation) GetResultFilePath() (string, error) {
	// Create *uint16 to hold result
	var _resultFilePath *uint16

	hr, _, _ := i.Vtbl.GetResultFilePath.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(_resultFilePath)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return "", syscall.Errno(hr)
	}
	// Get result and cleanup
	resultFilePath := UTF16PtrToString(_resultFilePath)
	CoTaskMemFree(unsafe.Pointer(_resultFilePath))
	return resultFilePath, nil
}

func (i *ICoreWebView2DownloadOperation) GetState() (COREWEBVIEW2_DOWNLOAD_STATE, error) {

	var downloadState COREWEBVIEW2_DOWNLOAD_STATE

	hr, _, _ := i.Vtbl.GetState.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&downloadState)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return downloadState, nil
}

func (i *ICoreWebView2DownloadOperation) GetInterruptReason() (COREWEBVIEW2_DOWNLOAD_INTERRUPT_REASON, error) {

	var interruptReason COREWEBVIEW2_DOWNLOAD_INTERRUPT_REASON

	hr, _, _ := i.Vtbl.GetInterruptReason.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&interruptReason)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return interruptReason, nil
}

func (i *ICoreWebView2DownloadOperation) Cancel() error {

	hr, _, _ := i.Vtbl.Cancel.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2DownloadOperation) Pause() error {

	hr, _, _ := i.Vtbl.Pause.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2DownloadOperation) Resume() error {

	hr, _, _ := i.Vtbl.Resume.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2DownloadOperation) GetCanResume() (bool, error) {
	// Create int32 to hold bool result
	var _canResume int32

	hr, _, _ := i.Vtbl.GetCanResume.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_canResume)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	canResume := _canResume != 0
	return canResume, nil
}
