package webview2loader

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetAvailableCoreWebView2BrowserVersionString get the browser version info including channel name.
func GetAvailableCoreWebView2BrowserVersionString(browserExecutableFolder string) (string, error) {
	if browserExecutableFolder != "" {
		clientPath, err := findEmbeddedClientDll(browserExecutableFolder)
		if err != nil {
			return "", err
		}

		return findEmbeddedBrowserVersion(clientPath)
	}

	return "", fmt.Errorf("not implemented yet for empty browserExecutableFolder ")
}

func findEmbeddedBrowserVersion(filename string) (string, error) {
	block, err := getFileVersionInfo(filename)
	if err != nil {
		return "", err
	}

	info, err := verQueryValueString(block, "\\StringFileInfo\\040904B0\\ProductVersion")
	if err != nil {
		return "", err
	}

	return info, nil
}

func findEmbeddedClientDll(embeddedEdgeSubFolder string) (outClientPath string, err error) {
	if !filepath.IsAbs(embeddedEdgeSubFolder) {
		exe, err := os.Executable()
		if err != nil {
			return "", err
		}

		embeddedEdgeSubFolder = filepath.Join(filepath.Dir(exe), embeddedEdgeSubFolder)
	}

	return findClientDllInFolder(embeddedEdgeSubFolder)
}

func findClientDllInFolder(folder string) (string, error) {
	arch := ""
	switch runtime.GOARCH {
	case "arm64":
		arch = "arm64"
	case "amd64":
		arch = "x64"
	case "386":
		arch = "x86"
	default:
		return "", fmt.Errorf("Unsupported architecture")
	}

	dllPath := filepath.Join(folder, "EBWebView", arch, "EmbeddedBrowserWebView.dll")
	if _, err := os.Stat(dllPath); err != nil {
		return "", err
	}
	return dllPath, nil
}
