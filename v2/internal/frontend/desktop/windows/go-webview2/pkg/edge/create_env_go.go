//go:build exp_gowebview2loader

package edge

import (
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2/webviewloader"
)

func createCoreWebView2EnvironmentWithOptions(browserExecutableFolder, userDataFolder string, environmentCompletedHandle *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) error {
	e := &environmentCreatedHandler{environmentCompletedHandle}
	return webviewloader.CreateCoreWebView2EnvironmentWithOptions(
		e,
		webviewloader.WithBrowserExecutableFolder(browserExecutableFolder),
		webviewloader.WithUserDataFolder(userDataFolder),
	)
}

type environmentCreatedHandler struct {
	originalHandler *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler
}

func (r *environmentCreatedHandler) EnvironmentCompleted(errorCode webviewloader.HRESULT, createdEnvironment *webviewloader.ICoreWebView2Environment) webviewloader.HRESULT {
	env := (*ICoreWebView2Environment)(unsafe.Pointer(createdEnvironment))
	res := r.originalHandler.impl.EnvironmentCompleted(uintptr(errorCode), env)
	return webviewloader.HRESULT(res)
}
