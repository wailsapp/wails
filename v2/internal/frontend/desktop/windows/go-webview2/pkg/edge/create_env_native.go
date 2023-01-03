//go:build windows && native_webview2loader

package edge

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/internal/frontend/desktop/windows/go-webview2/webviewloader"

	"golang.org/x/sys/windows"
)

func createCoreWebView2EnvironmentWithOptions(browserExecutableFolder, userDataFolder string, environmentCompletedHandle *iCoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, additionalBrowserArgs string) error {
	browserPathPtr, err := windows.UTF16PtrFromString(browserExecutableFolder)
	if err != nil {
		return fmt.Errorf("Error calling UTF16PtrFromString for %s: %v", browserExecutableFolder, err)
	}

	userPathPtr, err := windows.UTF16PtrFromString(userDataFolder)
	if err != nil {
		return fmt.Errorf("Error calling UTF16PtrFromString for %s: %v", userDataFolder, err)
	}

	hr, err := webviewloader.CreateCoreWebView2EnvironmentWithOptions(
		browserPathPtr,
		userPathPtr,
		uintptr(unsafe.Pointer(environmentCompletedHandle)),
		additionalBrowserArgs,
	)

	if hr != 0 {
		if err == nil || err == windows.ERROR_SUCCESS {
			err = syscall.Errno(hr)
		}
		return err
	}

	return nil
}
