//go:build windows
// +build windows

package cfd

import "github.com/go-ole/go-ole"

func initialize() {
	// Swallow error
	_ = ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_DISABLE_OLE1DDE)
}

// TODO doc
func NewOpenFileDialog(config DialogConfig) (OpenFileDialog, error) {
	initialize()

	openDialog, err := newIFileOpenDialog()
	if err != nil {
		return nil, err
	}
	err = config.apply(openDialog)
	if err != nil {
		return nil, err
	}
	return openDialog, nil
}

// TODO doc
func NewOpenMultipleFilesDialog(config DialogConfig) (OpenMultipleFilesDialog, error) {
	initialize()

	openDialog, err := newIFileOpenDialog()
	if err != nil {
		return nil, err
	}
	err = config.apply(openDialog)
	if err != nil {
		return nil, err
	}
	err = openDialog.setIsMultiselect(true)
	if err != nil {
		return nil, err
	}
	return openDialog, nil
}

// TODO doc
func NewSelectFolderDialog(config DialogConfig) (SelectFolderDialog, error) {
	initialize()

	openDialog, err := newIFileOpenDialog()
	if err != nil {
		return nil, err
	}
	err = config.apply(openDialog)
	if err != nil {
		return nil, err
	}
	err = openDialog.setPickFolders(true)
	if err != nil {
		return nil, err
	}
	return openDialog, nil
}

// TODO doc
func NewSaveFileDialog(config DialogConfig) (SaveFileDialog, error) {
	initialize()

	saveDialog, err := newIFileSaveDialog()
	if err != nil {
		return nil, err
	}
	err = config.apply(saveDialog)
	if err != nil {
		return nil, err
	}
	return saveDialog, nil
}
