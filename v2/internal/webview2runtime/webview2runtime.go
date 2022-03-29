//go:build windows
// +build windows

package webview2runtime

import (
	_ "embed"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
)

// Info contains all the information about an installation of the webview2 runtime.
type Info struct {
	Location        string
	Name            string
	Version         string
	SilentUninstall string
}

// IsOlderThan returns true if the installed version is older than the given required version.
// Returns error if something goes wrong.
func (i *Info) IsOlderThan(requiredVersion string) (bool, error) {
	var mod = syscall.NewLazyDLL("WebView2Loader.dll")
	var CompareBrowserVersions = mod.NewProc("CompareBrowserVersions")
	v1, err := syscall.UTF16PtrFromString(i.Version)
	if err != nil {
		return false, err
	}
	v2, err := syscall.UTF16PtrFromString(requiredVersion)
	if err != nil {
		return false, err
	}
	var result int = 9
	_, _, err = CompareBrowserVersions.Call(uintptr(unsafe.Pointer(v1)), uintptr(unsafe.Pointer(v2)), uintptr(unsafe.Pointer(&result)))
	if result < -1 || result > 1 {
		return false, err
	}
	return result == -1, nil
}

func downloadBootstrapper() (string, error) {
	bootstrapperURL := `https://go.microsoft.com/fwlink/p/?LinkId=2124703`
	installer := filepath.Join(os.TempDir(), `MicrosoftEdgeWebview2Setup.exe`)

	// Download installer
	out, err := os.Create(installer)
	defer out.Close()
	if err != nil {
		return "", err
	}
	resp, err := http.Get(bootstrapperURL)
	defer resp.Body.Close()
	if err != nil {
		err = out.Close()
		return "", err
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return installer, nil
}

// InstallUsingEmbeddedBootstrapper will download the bootstrapper from Microsoft and run it to install
// the latest version of the runtime.
// Returns true if the installer ran successfully.
// Returns an error if something goes wrong
func InstallUsingEmbeddedBootstrapper() (bool, error) {
	installer, err := WriteInstaller(os.TempDir())
	if err != nil {
		return false, err
	}
	result, err := runInstaller(installer)
	if err != nil {
		return false, err
	}

	return result, os.Remove(installer)

}

// InstallUsingBootstrapper will extract the embedded bootstrapper from Microsoft and run it to install
// the latest version of the runtime.
// Returns true if the installer ran successfully.
// Returns an error if something goes wrong
func InstallUsingBootstrapper() (bool, error) {

	installer, err := downloadBootstrapper()
	if err != nil {
		return false, err
	}

	result, err := runInstaller(installer)
	if err != nil {
		return false, err
	}

	return result, os.Remove(installer)

}

func runInstaller(installer string) (bool, error) {
	// Credit: https://stackoverflow.com/a/10385867
	cmd := exec.Command(installer)
	if err := cmd.Start(); err != nil {
		return false, err
	}
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus() == 0, nil
			}
		}
	}
	return true, nil
}

// Confirm will prompt the user with a message and OK / CANCEL buttons.
// Returns true if OK is selected by the user.
// Returns an error if something went wrong.
func Confirm(caption string, title string) (bool, error) {
	var flags uint = 0x00000001 // MB_OKCANCEL
	result, err := MessageBox(caption, title, flags)
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// Error will an error message to the user.
// Returns an error if something went wrong.
func Error(caption string, title string) error {
	var flags uint = 0x00000010 // MB_ICONERROR
	_, err := MessageBox(caption, title, flags)
	return err
}

// MessageBox prompts the user with the given caption and title.
// Flags may be provided to customise the dialog.
// Returns an error if something went wrong.
func MessageBox(caption string, title string, flags uint) (int, error) {
	captionUTF16, err := syscall.UTF16PtrFromString(caption)
	if err != nil {
		return -1, err
	}
	titleUTF16, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return -1, err
	}
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		uintptr(0),
		uintptr(unsafe.Pointer(captionUTF16)),
		uintptr(unsafe.Pointer(titleUTF16)),
		uintptr(flags))

	return int(ret), nil
}

// OpenInstallerDownloadWebpage will open the browser on the WebView2 download page
func OpenInstallerDownloadWebpage() error {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", "https://developer.microsoft.com/en-us/microsoft-edge/webview2/")
	return cmd.Run()
}
