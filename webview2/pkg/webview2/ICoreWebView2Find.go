//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2FindVtbl struct {
	IUnknownVtbl
	GetActiveMatchIndex ComProc
	GetMatchCount ComProc
	AddActiveMatchIndexChanged ComProc
	RemoveActiveMatchIndexChanged ComProc
	AddMatchCountChanged ComProc
	RemoveMatchCountChanged ComProc
	Start ComProc
	FindNext ComProc
	FindPrevious ComProc
	Stop ComProc
}

type ICoreWebView2Find struct {
	Vtbl *ICoreWebView2FindVtbl
}

func (i *ICoreWebView2Find) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2Find) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2Find) GetActiveMatchIndex() (int32, error) {

	var value int32

	hr, _, _ := i.Vtbl.GetActiveMatchIndex.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2Find) GetMatchCount() (int32, error) {

	var value int32

	hr, _, _ := i.Vtbl.GetMatchCount.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, nil
}

func (i *ICoreWebView2Find) AddActiveMatchIndexChanged(eventHandler *ICoreWebView2FindActiveMatchIndexChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddActiveMatchIndexChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Find) RemoveActiveMatchIndexChanged(token EventRegistrationToken) error {


	hr, _, _ := i.Vtbl.RemoveActiveMatchIndexChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Find) AddMatchCountChanged(eventHandler *ICoreWebView2FindMatchCountChangedEventHandler) (EventRegistrationToken, error) {

	var token EventRegistrationToken

	hr, _, _ := i.Vtbl.AddMatchCountChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(eventHandler)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return EventRegistrationToken{}, syscall.Errno(hr)
	}
	return token, nil
}

func (i *ICoreWebView2Find) RemoveMatchCountChanged(token EventRegistrationToken) error {


	hr, _, _ := i.Vtbl.RemoveMatchCountChanged.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&token)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Find) Start(options *ICoreWebView2FindOptions, handler *ICoreWebView2FindStartCompletedHandler) error {


	hr, _, _ := i.Vtbl.Start.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(options)),
		uintptr(unsafe.Pointer(handler)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Find) FindNext() error {


	hr, _, _ := i.Vtbl.FindNext.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Find) FindPrevious() error {


	hr, _, _ := i.Vtbl.FindPrevious.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}

func (i *ICoreWebView2Find) Stop() error {


	hr, _, _ := i.Vtbl.Stop.Call(
		uintptr(unsafe.Pointer(i)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return nil
}
