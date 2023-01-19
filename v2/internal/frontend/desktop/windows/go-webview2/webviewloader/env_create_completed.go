//go:build windows && !native_webview2loader

package webviewloader

import (
	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2/pkg/combridge"
)

// HRESULT
//
// See https://docs.microsoft.com/en-us/windows/win32/seccrypto/common-hresult-values
type HRESULT int32

// ICoreWebView2Environment Represents the WebView2 Environment
//
// See https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environment
type ICoreWebView2Environment = combridge.IUnknownImpl

// ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler receives the WebView2Environment created using CreateCoreWebView2Environment.
type ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler interface {
	// EnvironmentCompleted is invoked to receive the created WebView2Environment
	//
	// See https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2createcorewebview2environmentcompletedhandler?#invoke
	EnvironmentCompleted(errorCode HRESULT, createdEnvironment *ICoreWebView2Environment) HRESULT
}

type iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler interface {
	combridge.IUnknown
	ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler
}

func init() {
	combridge.RegisterVTable[combridge.IUnknown, iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler](
		"{4e8a3389-c9d8-4bd2-b6b5-124fee6cc14d}",
		_iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke,
	)
}

func _iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerInvoke(this uintptr, errorCode HRESULT, env *combridge.IUnknownImpl) uintptr {
	res := combridge.Resolve[iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler](this).EnvironmentCompleted(errorCode, env)
	return uintptr(res)
}
