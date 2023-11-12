//go:build darwin
// +build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework WebKit
#import <Foundation/Foundation.h>
#import "Application.h"
#import "WailsContext.h"
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"unsafe"

	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/frontend"
)

// Obj-C dialog methods send the response to this channel
var (
	messageDialogResponse  = make(chan int)
	openFileDialogResponse = make(chan string)
	saveFileDialogResponse = make(chan string)
	dialogLock             sync.Mutex
)

// OpenDirectoryDialog prompts the user to select a directory
func (f *Frontend) OpenDirectoryDialog(options frontend.OpenDialogOptions) (string, error) {
	results, err := f.openDialog(&options, false, false, true)
	if err != nil {
		return "", err
	}
	var selected string
	if len(results) > 0 {
		selected = results[0]
	}
	return selected, nil
}

func (f *Frontend) openDialog(options *frontend.OpenDialogOptions, multiple bool, allowfiles bool, allowdirectories bool) ([]string, error) {
	dialogLock.Lock()
	defer dialogLock.Unlock()

	c := NewCalloc()
	defer c.Free()
	title := c.String(options.Title)
	defaultFilename := c.String(options.DefaultFilename)
	defaultDirectory := c.String(options.DefaultDirectory)
	allowDirectories := bool2Cint(allowdirectories)
	allowFiles := bool2Cint(allowfiles)
	canCreateDirectories := bool2Cint(options.CanCreateDirectories)
	treatPackagesAsDirectories := bool2Cint(options.TreatPackagesAsDirectories)
	resolveAliases := bool2Cint(options.ResolvesAliases)
	showHiddenFiles := bool2Cint(options.ShowHiddenFiles)
	allowMultipleFileSelection := bool2Cint(multiple)

	var filterStrings slicer.StringSlicer
	if options.Filters != nil {
		for _, filter := range options.Filters {
			thesePatterns := strings.Split(filter.Pattern, ";")
			for _, pattern := range thesePatterns {
				pattern = strings.TrimSpace(pattern)
				if pattern != "" {
					filterStrings.Add(pattern)
				}
			}
		}
		filterStrings.Deduplicate()
	}
	filters := filterStrings.Join(";")
	C.OpenFileDialog(f.mainWindow.context, title, defaultFilename, defaultDirectory, allowDirectories, allowFiles, canCreateDirectories, treatPackagesAsDirectories, resolveAliases, showHiddenFiles, allowMultipleFileSelection, c.String(filters))

	result := <-openFileDialogResponse

	var parsedResults []string
	err := json.Unmarshal([]byte(result), &parsedResults)

	return parsedResults, err
}

// OpenFileDialog prompts the user to select a file
func (f *Frontend) OpenFileDialog(options frontend.OpenDialogOptions) (string, error) {
	results, err := f.openDialog(&options, false, true, false)
	if err != nil {
		return "", err
	}
	var selected string
	if len(results) > 0 {
		selected = results[0]
	}
	return selected, nil
}

// OpenMultipleFilesDialog prompts the user to select a file
func (f *Frontend) OpenMultipleFilesDialog(options frontend.OpenDialogOptions) ([]string, error) {
	return f.openDialog(&options, true, true, false)
}

// SaveFileDialog prompts the user to select a file
func (f *Frontend) SaveFileDialog(options frontend.SaveDialogOptions) (string, error) {
	dialogLock.Lock()
	defer dialogLock.Unlock()

	c := NewCalloc()
	defer c.Free()
	title := c.String(options.Title)
	defaultFilename := c.String(options.DefaultFilename)
	defaultDirectory := c.String(options.DefaultDirectory)
	canCreateDirectories := bool2Cint(options.CanCreateDirectories)
	treatPackagesAsDirectories := bool2Cint(options.TreatPackagesAsDirectories)
	showHiddenFiles := bool2Cint(options.ShowHiddenFiles)

	var filterStrings slicer.StringSlicer
	if options.Filters != nil {
		for _, filter := range options.Filters {
			thesePatterns := strings.Split(filter.Pattern, ";")
			for _, pattern := range thesePatterns {
				pattern = strings.TrimSpace(pattern)
				if pattern != "" {
					filterStrings.Add(pattern)
				}
			}
		}
		filterStrings.Deduplicate()
	}
	filters := filterStrings.Join(";")
	C.SaveFileDialog(f.mainWindow.context, title, defaultFilename, defaultDirectory, canCreateDirectories, treatPackagesAsDirectories, showHiddenFiles, c.String(filters))

	result := <-saveFileDialogResponse

	return result, nil
}

// MessageDialog show a message dialog to the user
func (f *Frontend) MessageDialog(options frontend.MessageDialogOptions) (string, error) {
	dialogLock.Lock()
	defer dialogLock.Unlock()

	c := NewCalloc()
	defer c.Free()
	dialogType := c.String(string(options.Type))
	title := c.String(options.Title)
	message := c.String(options.Message)
	defaultButton := c.String(options.DefaultButton)
	cancelButton := c.String(options.CancelButton)
	const MaxButtons = 4
	var buttons [MaxButtons]*C.char
	for index, buttonText := range options.Buttons {
		if index == MaxButtons {
			return "", fmt.Errorf("max %d buttons supported (%d given)", MaxButtons, len(options.Buttons))
		}
		buttons[index] = c.String(buttonText)
	}

	var iconData unsafe.Pointer
	var iconDataLength C.int
	if options.Icon != nil {
		iconData = unsafe.Pointer(&options.Icon[0])
		iconDataLength = C.int(len(options.Icon))
	}

	C.MessageDialog(f.mainWindow.context, dialogType, title, message, buttons[0], buttons[1], buttons[2], buttons[3], defaultButton, cancelButton, iconData, iconDataLength)

	result := <-messageDialogResponse

	selectedC := buttons[result]
	var selected string
	if selectedC != nil {
		selected = options.Buttons[result]
	}
	return selected, nil
}

//export processMessageDialogResponse
func processMessageDialogResponse(selection int) {
	messageDialogResponse <- selection
}

//export processOpenFileDialogResponse
func processOpenFileDialogResponse(cselection *C.char) {
	selection := C.GoString(cselection)
	openFileDialogResponse <- selection
}

//export processSaveFileDialogResponse
func processSaveFileDialogResponse(cselection *C.char) {
	selection := C.GoString(cselection)
	saveFileDialogResponse <- selection
}
