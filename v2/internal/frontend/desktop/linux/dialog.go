//go:build linux
// +build linux

package linux

import (
	"github.com/wailsapp/wails/v2/internal/frontend"
	"unsafe"
)

/*
#include <stdlib.h>
#include "gtk/gtk.h"
*/
import "C"

const (
	GTK_FILE_CHOOSER_ACTION_OPEN          C.GtkFileChooserAction = C.GTK_FILE_CHOOSER_ACTION_OPEN
	GTK_FILE_CHOOSER_ACTION_SAVE          C.GtkFileChooserAction = C.GTK_FILE_CHOOSER_ACTION_SAVE
	GTK_FILE_CHOOSER_ACTION_SELECT_FOLDER C.GtkFileChooserAction = C.GTK_FILE_CHOOSER_ACTION_SELECT_FOLDER
)

var openFileResults = make(chan []string)
var messageDialogResult = make(chan string)

func (f *Frontend) OpenFileDialog(dialogOptions frontend.OpenDialogOptions) (result string, err error) {
	f.mainWindow.OpenFileDialog(dialogOptions, 0, GTK_FILE_CHOOSER_ACTION_OPEN)
	results := <-openFileResults
	if len(results) == 1 {
		return results[0], nil
	}
	return "", nil
}

func (f *Frontend) OpenMultipleFilesDialog(dialogOptions frontend.OpenDialogOptions) ([]string, error) {
	f.mainWindow.OpenFileDialog(dialogOptions, 1, GTK_FILE_CHOOSER_ACTION_OPEN)
	result := <-openFileResults
	return result, nil
}

func (f *Frontend) OpenDirectoryDialog(dialogOptions frontend.OpenDialogOptions) (string, error) {
	f.mainWindow.OpenFileDialog(dialogOptions, 0, GTK_FILE_CHOOSER_ACTION_SELECT_FOLDER)
	result := <-openFileResults
	if len(result) == 1 {
		return result[0], nil
	}
	return "", nil
}

func (f *Frontend) SaveFileDialog(dialogOptions frontend.SaveDialogOptions) (string, error) {
	options := frontend.OpenDialogOptions{
		DefaultDirectory:     dialogOptions.DefaultDirectory,
		DefaultFilename:      dialogOptions.DefaultFilename,
		Title:                dialogOptions.Title,
		Filters:              dialogOptions.Filters,
		ShowHiddenFiles:      dialogOptions.ShowHiddenFiles,
		CanCreateDirectories: dialogOptions.CanCreateDirectories,
	}
	f.mainWindow.OpenFileDialog(options, 0, GTK_FILE_CHOOSER_ACTION_SAVE)
	results := <-openFileResults
	if len(results) == 1 {
		return results[0], nil
	}
	return "", nil
}

func (f *Frontend) MessageDialog(dialogOptions frontend.MessageDialogOptions) (string, error) {
	f.mainWindow.MessageDialog(dialogOptions)
	return <-messageDialogResult, nil
}

//export processOpenFileResult
func processOpenFileResult(carray **C.char) {
	// Create a Go slice from the C array
	var result []string
	goArray := (*[1024]*C.char)(unsafe.Pointer(carray))[:1024:1024]
	for _, s := range goArray {
		if s == nil {
			break
		}
		result = append(result, C.GoString(s))
	}
	openFileResults <- result
}

//export processMessageDialogResult
func processMessageDialogResult(result *C.char) {
	messageDialogResult <- C.GoString(result)
}
