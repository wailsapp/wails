package ffenestri

/*

#cgo linux CFLAGS: -DFFENESTRI_LINUX=1
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.0

#include <stdlib.h>
#include "ffenestri.h"

*/
import "C"

import (
	"strconv"

	"github.com/wailsapp/wails/v2/pkg/options/dialog"

	"github.com/wailsapp/wails/v2/internal/logger"
)

// Client is our implentation of messageDispatcher.Client
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

func (c *Client) WindowSetMinSize(width int, height int) {
	C.SetMinWindowSize(c.app.app, C.int(width), C.int(height))
}

func (c *Client) WindowSetMaxSize(width int, height int) {
	C.SetMaxWindowSize(c.app.app, C.int(width), C.int(height))
}

// WindowSetColour sets the window colour
func (c *Client) WindowSetColour(colour int) {
	r, g, b, a := intToColour(colour)
	C.SetColour(c.app.app, r, g, b, a)
}

// OpenDialog will open a dialog with the given title and filter
func (c *Client) OpenDialog(dialogOptions *dialog.OpenDialog, callbackID string) {
	C.OpenDialog(c.app.app,
		c.app.string2CString(callbackID),
		c.app.string2CString(dialogOptions.Title),
		c.app.string2CString(dialogOptions.Filters),
		c.app.string2CString(dialogOptions.DefaultFilename),
		c.app.string2CString(dialogOptions.DefaultDirectory),
		c.app.bool2Cint(dialogOptions.AllowFiles),
		c.app.bool2Cint(dialogOptions.AllowDirectories),
		c.app.bool2Cint(dialogOptions.AllowMultiple),
		c.app.bool2Cint(dialogOptions.ShowHiddenFiles),
		c.app.bool2Cint(dialogOptions.CanCreateDirectories),
		c.app.bool2Cint(dialogOptions.ResolvesAliases),
		c.app.bool2Cint(dialogOptions.TreatPackagesAsDirectories),
	)
}

// SaveDialog will open a dialog with the given title and filter
func (c *Client) SaveDialog(dialogOptions *dialog.SaveDialog, callbackID string) {
	C.SaveDialog(c.app.app,
		c.app.string2CString(callbackID),
		c.app.string2CString(dialogOptions.Title),
		c.app.string2CString(dialogOptions.Filters),
		c.app.string2CString(dialogOptions.DefaultFilename),
		c.app.string2CString(dialogOptions.DefaultDirectory),
		c.app.bool2Cint(dialogOptions.ShowHiddenFiles),
		c.app.bool2Cint(dialogOptions.CanCreateDirectories),
		c.app.bool2Cint(dialogOptions.TreatPackagesAsDirectories),
	)
}

// MessageDialog will open a message dialog with the given options
func (c *Client) MessageDialog(dialogOptions *dialog.MessageDialog, callbackID string) {

	// Sanity check button length
	if len(dialogOptions.Buttons) > 4 {
		c.app.logger.Error("Given %d message dialog buttons. Maximum is 4", len(dialogOptions.Buttons))
		return
	}

	// Process buttons
	buttons := []string{"", "", "", ""}
	for i, button := range dialogOptions.Buttons {
		buttons[i] = button
	}

	C.MessageDialog(c.app.app,
		c.app.string2CString(callbackID),
		c.app.string2CString(string(dialogOptions.Type)),
		c.app.string2CString(dialogOptions.Title),
		c.app.string2CString(dialogOptions.Message),
		c.app.string2CString(dialogOptions.Icon),
		c.app.string2CString(buttons[0]),
		c.app.string2CString(buttons[1]),
		c.app.string2CString(buttons[2]),
		c.app.string2CString(buttons[3]),
		c.app.string2CString(dialogOptions.DefaultButton),
		c.app.string2CString(dialogOptions.CancelButton))
}

func (c *Client) DarkModeEnabled(callbackID string) {
	C.DarkModeEnabled(c.app.app, c.app.string2CString(callbackID))
}

func (c *Client) SetApplicationMenu(applicationMenuJSON string) {
	C.SetApplicationMenu(c.app.app, c.app.string2CString(applicationMenuJSON))
}

func (c *Client) SetTrayMenu(trayMenuJSON string) {
	C.SetTrayMenu(c.app.app, c.app.string2CString(trayMenuJSON))
}

func (c *Client) UpdateTrayMenuLabel(JSON string) {
	C.UpdateTrayMenuLabel(c.app.app, c.app.string2CString(JSON))
}

func (c *Client) UpdateContextMenu(contextMenuJSON string) {
	C.UpdateContextMenu(c.app.app, c.app.string2CString(contextMenuJSON))
}

func (c *Client) DeleteTrayMenuByID(id string) {
	C.DeleteTrayMenuByID(c.app.app, c.app.string2CString(id))
}
