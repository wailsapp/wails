//go:build !windows
// +build !windows

package cfd

import "fmt"

var errUnsupported = fmt.Errorf("common file dialogs are only available on windows")

// TODO doc
func NewOpenFileDialog(config DialogConfig) (OpenFileDialog, error) {
	return nil, errUnsupported
}

// TODO doc
func NewOpenMultipleFilesDialog(config DialogConfig) (OpenMultipleFilesDialog, error) {
	return nil, errUnsupported
}

// TODO doc
func NewSelectFolderDialog(config DialogConfig) (SelectFolderDialog, error) {
	return nil, errUnsupported
}

// TODO doc
func NewSaveFileDialog(config DialogConfig) (SaveFileDialog, error) {
	return nil, errUnsupported
}
