//go:build windows

package webview2

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type ICoreWebView2_8Vtbl struct {
	IUnknownVtbl
	AddIsMutedChanged                   ComProc
	RemoveIsMutedChanged                ComProc
	GetIsMuted                          ComProc
	PutIsMuted                          ComProc
	AddIsDocumentPlayingAudioChanged    ComProc
	RemoveIsDocumentPlayingAudioChanged ComProc
	GetIsDocumentPlayingAudio           ComProc
}

type ICoreWebView2_8 struct {
	Vtbl *ICoreWebView2_8Vtbl
}

func (i *ICoreWebView2_8) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func (i *ICoreWebView2) GetICoreWebView2_8() *ICoreWebView2_8 {
	var result *ICoreWebView2_8

	iidICoreWebView2_8 := NewGUID("{E9632730-6E1E-43AB-B7B8-7B2C9E62E094}")
	_, _, _ = i.Vtbl.QueryInterface.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(iidICoreWebView2_8)),
		uintptr(unsafe.Pointer(&result)))

	return result
}

func (i *ICoreWebView2_8) AddIsMutedChanged(eventHandler *ICoreWebView2IsMutedChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddIsMutedChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_8) RemoveIsMutedChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveIsMutedChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_8) GetIsMuted() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsMuted.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	value := _value != 0
	return value, nil
}

func (i *ICoreWebView2_8) PutIsMuted(value bool) error {

	hr, _, _ := i.Vtbl.PutIsMuted.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_8) AddIsDocumentPlayingAudioChanged(eventHandler *ICoreWebView2IsDocumentPlayingAudioChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddIsDocumentPlayingAudioChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2_8) RemoveIsDocumentPlayingAudioChanged(token EventRegistrationToken) error {

	hr, _, _ := i.Vtbl.RemoveIsDocumentPlayingAudioChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2_8) GetIsDocumentPlayingAudio() (bool, error) {
	// Create int32 to hold bool result
	var _value int32

	hr, _, _ := i.Vtbl.GetIsDocumentPlayingAudio.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return false, syscall.Errno(hr)
	}
	// Get result and cleanup
	value := _value != 0
	return value, nil
}
