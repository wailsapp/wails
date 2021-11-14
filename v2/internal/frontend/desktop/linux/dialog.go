//go:build linux
// +build linux

package linux

import "github.com/wailsapp/wails/v2/internal/frontend"

func (f *Frontend) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	panic("implement me")
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
