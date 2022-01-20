//go:build linux
// +build linux

package linux

import (
	"github.com/wailsapp/wails/v2/internal/frontend"
)
import "C"

var openFileResults = make(chan string)

func (f *Frontend) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (result string, err error) {
	f.mainWindow.OpenFileDialog(dialogOptions)
	result = <-openFileResults
	return
}

func (f *Frontend) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	panic("implement me")
}

func (f *Frontend) OpenDirectoryDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	panic("implement me")
}

func (f *Frontend) SaveFileDialog(dialogOptions frontend.SaveDialogOptions) (string, error) {
	panic("implement me")
}

func (f *Frontend) MessageDialog(dialogOptions frontend.MessageDialogOptions) (string, error) {
	panic("implement me")
}

//export processOpenFileResult
func processOpenFileResult(result *C.char) {
	openFileResults <- C.GoString(result)
}
