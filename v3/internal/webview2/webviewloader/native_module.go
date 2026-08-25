//go:build windows && native_webview2loader

package webviewloader

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/jchv/go-winloader"

	"golang.org/x/sys/windows"
)

func init() {
	preventEnvAndRegistryOverrides(nil, nil, "")
}

var (
	memOnce                                         sync.Once
	memModule                                       winloader.Module
	memCreate                                       winloader.Proc
	memCompareBrowserVersions                       winloader.Proc
	memGetAvailableCoreWebView2BrowserVersionString winloader.Proc
	memErr                                          error
)

const (
	// https://referencesource.microsoft.com/#system.web/Util/hresults.cs,20
	E_FILENOTFOUND = 0x80070002
)

// CompareBrowserVersions will compare the 2 given versions and return:
//
//	   Less than zero: v1 < v2
//	             zero: v1 == v2
//	Greater than zero: v1 > v2
func CompareBrowserVersions(v1 string, v2 string) (int, error) {
	_v1, err := windows.UTF16PtrFromString(v1)
	if err != nil {
		return 0, err
	}
	_v2, err := windows.UTF16PtrFromString(v2)
	if err != nil {
		return 0, err
	}

	err = loadFromMemory()
	if err != nil {
		return 0, err
	}

	var result int32
	_, _, err = memCompareBrowserVersions.Call(
		uint64(uintptr(unsafe.Pointer(_v1))),
		uint64(uintptr(unsafe.Pointer(_v2))),
		uint64(uintptr(unsafe.Pointer(&result))))

	if err != windows.ERROR_SUCCESS {
		return 0, err
	}
	return int(result), nil
}

// GetAvailableCoreWebView2BrowserVersionString returns version of the webview2 runtime.
// If path is empty, it will try to find installed webview2 is the system.
// If there is no version installed, a blank string is returned.
func GetAvailableCoreWebView2BrowserVersionString(path string) (string, error) {
	if path != "" {
		// The default implementation fails if CGO and a fixed browser path is used. It's caused by the go-winloader
		// which loads the native DLL from memory.
		// Use the new GoWebView2Loader in this case, in the future we will make GoWebView2Loader
		// feature-complete and remove the use of the native DLL and go-winloader.
		version, err := goGetAvailableCoreWebView2BrowserVersionString(path)
		if errors.Is(err, errNoClientDLLFound) {
			// WebView2 is not found
			return "", nil
		} else if err != nil {
			return "", err
		}

		return version, nil
	}

	err := loadFromMemory()
	if err != nil {
		return "", err
	}

	var browserPath *uint16 = nil
	if path != "" {
		browserPath, err = windows.UTF16PtrFromString(path)
		if err != nil {
			return "", fmt.Errorf("error calling UTF16PtrFromString for %s: %v", path, err)
		}
	}

	preventEnvAndRegistryOverrides(browserPath, nil, "")
	var result *uint16
	res, _, err := memGetAvailableCoreWebView2BrowserVersionString.Call(
		uint64(uintptr(unsafe.Pointer(browserPath))),
		uint64(uintptr(unsafe.Pointer(&result))))

	if res != 0 {
		if res == E_FILENOTFOUND {
			// WebView2 is not installed
			return "", nil
		}

		return "", fmt.Errorf("Unable to call GetAvailableCoreWebView2BrowserVersionString (%x): %w", res, err)
	}

	version := windows.UTF16PtrToString(result)
	windows.CoTaskMemFree(unsafe.Pointer(result))
	return version, nil
}

// CreateCoreWebView2EnvironmentWithOptions tries to load WebviewLoader2 and
// call the CreateCoreWebView2EnvironmentWithOptions routine.
func CreateCoreWebView2EnvironmentWithOptions(browserExecutableFolder, userDataFolder *uint16, environmentCompletedHandle uintptr, additionalBrowserArgs string) (uintptr, error) {
	err := loadFromMemory()
	if err != nil {
		return 0, err
	}

	preventEnvAndRegistryOverrides(browserExecutableFolder, userDataFolder, additionalBrowserArgs)
	res, _, _ := memCreate.Call(
		uint64(uintptr(unsafe.Pointer(browserExecutableFolder))),
		uint64(uintptr(unsafe.Pointer(userDataFolder))),
		0,
		uint64(environmentCompletedHandle),
	)
	return uintptr(res), nil
}

func loadFromMemory() error {
	var err error
	// DLL is not available natively. Try loading embedded copy.
	memOnce.Do(func() {
		memModule, memErr = winloader.LoadFromMemory(WebView2Loader)
		if memErr != nil {
			err = fmt.Errorf("Unable to load WebView2Loader.dll from memory: %w", memErr)
			return
		}
		memCreate = memModule.Proc("CreateCoreWebView2EnvironmentWithOptions")
		memCompareBrowserVersions = memModule.Proc("CompareBrowserVersions")
		memGetAvailableCoreWebView2BrowserVersionString = memModule.Proc("GetAvailableCoreWebView2BrowserVersionString")
	})
	return err
}

func preventEnvAndRegistryOverrides(browserFolder, userDataFolder *uint16, additionalBrowserArgs string) {
	// Setting these env variables to empty string also prevents registry overrides because webview2loader
	// checks for existence and not for empty value
	os.Setenv("WEBVIEW2_PIPE_FOR_SCRIPT_DEBUGGER", "")

	// Set these overrides to the values or empty to prevent registry and external env overrides
	os.Setenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", additionalBrowserArgs)
	os.Setenv("WEBVIEW2_RELEASE_CHANNEL_PREFERENCE", "0")
	os.Setenv("WEBVIEW2_BROWSER_EXECUTABLE_FOLDER", windows.UTF16PtrToString(browserFolder))
	os.Setenv("WEBVIEW2_USER_DATA_FOLDER", windows.UTF16PtrToString(userDataFolder))
}

func goGetAvailableCoreWebView2BrowserVersionString(browserExecutableFolder string) (string, error) {
	clientPath, err := findEmbeddedClientDll(browserExecutableFolder)
	if err != nil {
		return "", err
	}

	return findEmbeddedBrowserVersion(clientPath)
}
