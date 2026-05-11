//go:build windows

package webview2

import (
	"unsafe"
)

type ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerVtbl struct {
	IUnknownVtbl
	Invoke ComProc
}

type ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandler struct {
	Vtbl *ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerVtbl
	impl ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerImpl
}

func (i *ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandler) AddRef() uintptr {
	refCounter, _, _ := i.Vtbl.AddRef.Call(uintptr(unsafe.Pointer(i)))
	return refCounter
}

func ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerIUnknownQueryInterface(this *ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandler, refiid, object uintptr) uintptr {
	return this.impl.QueryInterface(refiid, object)
}

func ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerIUnknownAddRef(this *ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandler) uintptr {
	return this.impl.AddRef()
}

func ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerIUnknownRelease(this *ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandler) uintptr {
	return this.impl.Release()
}

func ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerInvoke(this *ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandler, errorCode uintptr, result *ICoreWebView2PermissionSettingCollectionView) uintptr {
	return this.impl.GetNonDefaultPermissionSettingsCompleted(errorCode, result)
}

type ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerImpl interface {
	IUnknownImpl
	GetNonDefaultPermissionSettingsCompleted(errorCode uintptr, result *ICoreWebView2PermissionSettingCollectionView) uintptr
}

var ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerFn = ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerVtbl{
	IUnknownVtbl{
		NewComProc(ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerIUnknownQueryInterface),
		NewComProc(ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerIUnknownAddRef),
		NewComProc(ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerIUnknownRelease),
	},
	NewComProc(ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerInvoke),
}

func NewICoreWebView2GetNonDefaultPermissionSettingsCompletedHandler(impl ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerImpl) *ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandler {
	return &ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandler{
		Vtbl: &ICoreWebView2GetNonDefaultPermissionSettingsCompletedHandlerFn,
		impl: impl,
	}
}
