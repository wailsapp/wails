//go:build windows

package webview2
import (
	"unsafe"
	"syscall"
	"golang.org/x/sys/windows"
)

type ICoreWebView2EnvironmentOptions7Vtbl struct {
	IUnknownVtbl
	GetChannelSearchKind ComProc
	PutChannelSearchKind ComProc
	GetReleaseChannels ComProc
	PutReleaseChannels ComProc
}

type ICoreWebView2EnvironmentOptions7 struct {
	Vtbl *ICoreWebView2EnvironmentOptions7Vtbl
}

func (i *ICoreWebView2EnvironmentOptions7) AddRef() uint32 {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}

func (i *ICoreWebView2EnvironmentOptions7) Release() uint32 {
	refCounter, _, _ := i.Vtbl.Release.Call(uintptr(unsafe.Pointer(i)))
	return uint32(refCounter)
}


func (i *ICoreWebView2EnvironmentOptions7) GetChannelSearchKind() (COREWEBVIEW2_CHANNEL_SEARCH_KIND, error) {

	var value COREWEBVIEW2_CHANNEL_SEARCH_KIND

	hr, _, err := i.Vtbl.GetChannelSearchKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2EnvironmentOptions7) PutChannelSearchKind(value COREWEBVIEW2_CHANNEL_SEARCH_KIND) error {


	hr, _, err := i.Vtbl.PutChannelSearchKind.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}

func (i *ICoreWebView2EnvironmentOptions7) GetReleaseChannels() (COREWEBVIEW2_RELEASE_CHANNELS, error) {

	var value COREWEBVIEW2_RELEASE_CHANNELS

	hr, _, err := i.Vtbl.GetReleaseChannels.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&value)),
	)
	if windows.Handle(hr) != windows.S_OK {
		return 0, syscall.Errno(hr)
	}
	return value, err
}

func (i *ICoreWebView2EnvironmentOptions7) PutReleaseChannels(value COREWEBVIEW2_RELEASE_CHANNELS) error {


	hr, _, err := i.Vtbl.PutReleaseChannels.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(value),
	)
	if windows.Handle(hr) != windows.S_OK {
		return syscall.Errno(hr)
	}
	return err
}
