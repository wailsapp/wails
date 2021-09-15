//go:build darwin

package darwin

import "C"
import (
	"github.com/wailsapp/wails/v2/internal/frontend"
)

// OpenDirectoryDialog prompts the user to select a directory
func (f *Frontend) OpenDirectoryDialog(options frontend.OpenDialogOptions) (string, error) {
	return "", nil
}

// OpenFileDialog prompts the user to select a file
func (f *Frontend) OpenFileDialog(options frontend.OpenDialogOptions) (string, error) {
	return "", nil
}

// OpenMultipleFilesDialog prompts the user to select a file
func (f *Frontend) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	return []string{}, nil
}

// SaveFileDialog prompts the user to select a file
func (f *Frontend) SaveFileDialog(dialogOptions frontend.SaveDialogOptions) (string, error) {
	return "", nil
}

// MessageDialog show a message dialog to the user
func (f *Frontend) MessageDialog(options frontend.MessageDialogOptions) (string, error) {
	return "", nil
}
