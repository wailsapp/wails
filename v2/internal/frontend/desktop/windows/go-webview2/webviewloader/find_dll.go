//go:build windows

package webviewloader

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/windows/registry"
)

var (
	errNoClientDLLFound = errors.New("no webview2 found")
)

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
		return "", mapFindErr(err)
	}
	return dllPath, nil
}

func mapFindErr(err error) error {
	if errors.Is(err, registry.ErrNotExist) {
		return errNoClientDLLFound
	}
	if errors.Is(err, os.ErrNotExist) {
		return errNoClientDLLFound
	}
	return err
}
