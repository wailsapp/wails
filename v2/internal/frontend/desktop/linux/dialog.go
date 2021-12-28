//go:build linux
// +build linux

package linux

import (
	"github.com/wailsapp/wails/v2/internal/frontend"
)
import "C"

var openFileResults = make(chan string)

func (f *Frontend) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (result string, err error) {

	f.dispatch(func() {
		println("Before OpenFileDialog")
		f.mainWindow.OpenFileDialog(dialogOptions)
		println("After OpenFileDialog")
	})
	println("Waiting for result")
	result = <-openFileResults
	println("Got result")
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
