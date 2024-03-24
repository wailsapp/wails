package null

import "github.com/wailsapp/wails/v2/internal/frontend"

// OpenFileDialog does nothing
func (f *Frontend) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (result string, err error) {
	return "", nil
}

// OpenMultipleFilesDialog does nothing
func (f *Frontend) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	return []string{}, nil
}

// OpenDirectoryDialog does nothing
func (f *Frontend) OpenDirectoryDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	return "", nil
}

// SaveFileDialog does nothing
func (f *Frontend) SaveFileDialog(dialogOptions frontend.SaveDialogOptions) (string, error) {
	return "", nil
}

// MessageDialog does nothing
func (f *Frontend) MessageDialog(dialogOptions frontend.MessageDialogOptions) (string, error) {
	return "", nil
}
