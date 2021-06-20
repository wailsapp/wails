// +build windows

package ffenestri

/*
#include "ffenestri.h"
*/
import "C"

import (
	"encoding/json"
	"github.com/harry1453/go-common-file-dialog/cfd"
	"golang.org/x/sys/windows"
	"log"
	"strconv"
	"syscall"

	"github.com/wailsapp/wails/v2/pkg/options/dialog"

	"github.com/wailsapp/wails/v2/internal/logger"
)

// Client is our implementation of messageDispatcher.Client
type Client struct {
	app    *Application
	logger logger.CustomLogger
}

func newClient(app *Application) *Client {
	return &Client{
		app:    app,
		logger: app.logger,
	}
}

// Quit the application
func (c *Client) Quit() {
	c.app.logger.Trace("Got shutdown message")
	C.Quit(c.app.app)
}

// NotifyEvent will pass on the event message to the frontend
func (c *Client) NotifyEvent(message string) {
	eventMessage := `window.wails._.Notify(` + strconv.Quote(message) + `);`
	c.app.logger.Trace("eventMessage = %+v", eventMessage)
	C.ExecJS(c.app.app, c.app.string2CString(eventMessage))
}

// CallResult contains the result of the call from JS
func (c *Client) CallResult(message string) {
	callbackMessage := `window.wails._.Callback(` + strconv.Quote(message) + `);`
	c.app.logger.Trace("callbackMessage = %+v", callbackMessage)
	C.ExecJS(c.app.app, c.app.string2CString(callbackMessage))
}

// WindowSetTitle sets the window title to the given string
func (c *Client) WindowSetTitle(title string) {
	C.SetTitle(c.app.app, c.app.string2CString(title))
}

// WindowFullscreen will set the window to be fullscreen
func (c *Client) WindowFullscreen() {
	C.Fullscreen(c.app.app)
}

// WindowUnFullscreen will unfullscreen the window
func (c *Client) WindowUnFullscreen() {
	C.UnFullscreen(c.app.app)
}

// WindowShow will show the window
func (c *Client) WindowShow() {
	C.Show(c.app.app)
}

// WindowHide will hide the window
func (c *Client) WindowHide() {
	C.Hide(c.app.app)
}

// WindowCenter will hide the window
func (c *Client) WindowCenter() {
	C.Center(c.app.app)
}

// WindowMaximise will maximise the window
func (c *Client) WindowMaximise() {
	C.Maximise(c.app.app)
}

// WindowMinimise will minimise the window
func (c *Client) WindowMinimise() {
	C.Minimise(c.app.app)
}

// WindowUnmaximise will unmaximise the window
func (c *Client) WindowUnmaximise() {
	C.Unmaximise(c.app.app)
}

// WindowUnminimise will unminimise the window
func (c *Client) WindowUnminimise() {
	C.Unminimise(c.app.app)
}

// WindowPosition will position the window to x,y on the
// monitor that the window is mostly on
func (c *Client) WindowPosition(x int, y int) {
	C.SetPosition(c.app.app, C.int(x), C.int(y))
}

// WindowSize will resize the window to the given
// width and height
func (c *Client) WindowSize(width int, height int) {
	C.SetSize(c.app.app, C.int(width), C.int(height))
}

// WindowSetMinSize sets the minimum window size
func (c *Client) WindowSetMinSize(width int, height int) {
	C.SetMinWindowSize(c.app.app, C.int(width), C.int(height))
}

// WindowSetMaxSize sets the maximum window size
func (c *Client) WindowSetMaxSize(width int, height int) {
	C.SetMaxWindowSize(c.app.app, C.int(width), C.int(height))
}

// WindowSetColour sets the window colour
func (c *Client) WindowSetColour(colour int) {
	r, g, b, a := intToColour(colour)
	C.SetColour(c.app.app, r, g, b, a)
}

func convertFilters(filters []dialog.FileFilter) []cfd.FileFilter {
	var result []cfd.FileFilter
	for _, filter := range filters {
		result = append(result, cfd.FileFilter(filter))
	}
	return result
}

// OpenFileDialog will open a dialog with the given title and filter
func (c *Client) OpenFileDialog(options *dialog.OpenDialog, callbackID string) {
	config := cfd.DialogConfig{
		Folder:      options.DefaultDirectory,
		FileFilters: convertFilters(options.Filters),
		FileName:    options.DefaultFilename,
	}
	thisdialog, err := cfd.NewOpenFileDialog(config)
	if err != nil {
		log.Fatal(err)
	}
	//thisdialog.SetParentWindowHandle(uintptr(C.GetWindowHandle(c.app.app)))
	defer func(thisdialog cfd.OpenFileDialog) {
		err := thisdialog.Release()
		if err != nil {
			log.Fatal(err)
		}
	}(thisdialog)
	result, err := thisdialog.ShowAndGetResult()
	if err != nil {
		log.Fatal(err)
	}

	resultJSON, err := json.Marshal([]string{result})
	if err != nil {
		log.Fatal(err)
	}
	dispatcher.DispatchMessage("DO" + callbackID + "|" + string(resultJSON))

}

// OpenDirectoryDialog will open a dialog with the given title and filter
func (c *Client) OpenDirectoryDialog(dialogOptions *dialog.OpenDialog, callbackID string) {
	config := cfd.DialogConfig{
		Title:  dialogOptions.Title,
		Role:   "PickFolder",
		Folder: dialogOptions.DefaultDirectory,
	}
	thisDialog, err := cfd.NewSelectFolderDialog(config)
	if err != nil {
		log.Fatal()
	}
	//thisDialog.SetParentWindowHandle(uintptr(C.GetWindowHandle(c.app.app)))
	defer func(thisDialog cfd.SelectFolderDialog) {
		err := thisDialog.Release()
		if err != nil {
			log.Fatal(err)
		}
	}(thisDialog)
	result, err := thisDialog.ShowAndGetResult()
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	dispatcher.DispatchMessage("DD" + callbackID + "|" + string(resultJSON))
}

// OpenMultipleFilesDialog will open a dialog with the given title and filter
func (c *Client) OpenMultipleFilesDialog(dialogOptions *dialog.OpenDialog, callbackID string) {
	config := cfd.DialogConfig{
		Title:       dialogOptions.Title,
		Role:        "OpenMultipleFiles",
		FileFilters: convertFilters(dialogOptions.Filters),
		FileName:    dialogOptions.DefaultFilename,
		Folder:      dialogOptions.DefaultDirectory,
	}
	thisdialog, err := cfd.NewOpenMultipleFilesDialog(config)
	if err != nil {
		log.Fatal(err)
	}
	//thisdialog.SetParentWindowHandle(uintptr(C.GetWindowHandle(c.app.app)))
	defer func(thisdialog cfd.OpenMultipleFilesDialog) {
		err := thisdialog.Release()
		if err != nil {
			log.Fatal(err)
		}
	}(thisdialog)
	result, err := thisdialog.ShowAndGetResults()
	if err != nil {
		log.Fatal(err)
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	dispatcher.DispatchMessage("D*" + callbackID + "|" + string(resultJSON))
}

// SaveDialog will open a dialog with the given title and filter
func (c *Client) SaveDialog(dialogOptions *dialog.SaveDialog, callbackID string) {
	saveDialog, err := cfd.NewSaveFileDialog(cfd.DialogConfig{
		Title:       dialogOptions.Title,
		Role:        "SaveFile",
		FileFilters: convertFilters(dialogOptions.Filters),
		FileName:    dialogOptions.DefaultFilename,
		Folder:      dialogOptions.DefaultDirectory,
	})
	if err != nil {
		log.Fatal(err)
	}
	//saveDialog.SetParentWindowHandle(uintptr(C.GetWindowHandle(c.app.app)))
	if err := saveDialog.Show(); err != nil {
		log.Fatal(err)
	}
	result, err := saveDialog.GetResult()
	if err != nil {
		log.Fatal(err)
	}
	dispatcher.DispatchMessage("DS" + callbackID + "|" + result)
}

// MessageDialog will open a message dialog with the given options
func (c *Client) MessageDialog(options *dialog.MessageDialog, callbackID string) {

	title, err := syscall.UTF16PtrFromString(options.Title)
	if err != nil {
		log.Fatal(err)
	}
	message, err := syscall.UTF16PtrFromString(options.Message)
	if err != nil {
		log.Fatal(err)
	}
	var flags uint32
	switch options.Type {
	case dialog.InfoDialog:
		flags = windows.MB_OK | windows.MB_ICONINFORMATION
	case dialog.ErrorDialog:
		flags = windows.MB_ICONERROR | windows.MB_OK
	case dialog.QuestionDialog:
		flags = windows.MB_YESNO
	case dialog.WarningDialog:
		flags = windows.MB_OK | windows.MB_ICONWARNING
	}

	button, _ := windows.MessageBox(0, message, title, flags|windows.MB_SYSTEMMODAL)
	// This maps MessageBox return values to strings
	responses := []string{"", "Ok", "Cancel", "Abort", "Retry", "Ignore", "Yes", "No", "", "", "Try Again", "Continue"}
	result := "Error"
	if int(button) < len(responses) {
		result = responses[button]
	}
	dispatcher.DispatchMessage("DM" + callbackID + "|" + result)
}

// DarkModeEnabled sets the application to use dark mode
func (c *Client) DarkModeEnabled(callbackID string) {
	C.DarkModeEnabled(c.app.app, c.app.string2CString(callbackID))
}

// SetApplicationMenu sets the application menu
func (c *Client) SetApplicationMenu(applicationMenuJSON string) {
	C.SetApplicationMenu(c.app.app, c.app.string2CString(applicationMenuJSON))
}

// SetTrayMenu sets the tray menu
func (c *Client) SetTrayMenu(trayMenuJSON string) {
	C.SetTrayMenu(c.app.app, c.app.string2CString(trayMenuJSON))
}

// UpdateTrayMenuLabel updates a tray menu label
func (c *Client) UpdateTrayMenuLabel(JSON string) {
	C.UpdateTrayMenuLabel(c.app.app, c.app.string2CString(JSON))
}

// UpdateContextMenu will update the current context menu
func (c *Client) UpdateContextMenu(contextMenuJSON string) {
	C.UpdateContextMenu(c.app.app, c.app.string2CString(contextMenuJSON))
}

// DeleteTrayMenuByID will remove a tray menu based on the given id
func (c *Client) DeleteTrayMenuByID(id string) {
	C.DeleteTrayMenuByID(c.app.app, c.app.string2CString(id))
}
