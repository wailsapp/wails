package bridge

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/internal/messagedispatcher"

	"github.com/wailsapp/wails/v2/internal/logger"

	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/pkg/options/dialog"
)

type DialogClient struct {
	dispatcher *messagedispatcher.DispatchClient
	log        *logger.Logger
}

func (d *DialogClient) DeleteTrayMenuByID(id string) {
}

func NewDialogClient(log *logger.Logger) *DialogClient {
	return &DialogClient{
		log: log,
	}
}

func (d *DialogClient) Quit() {
}

func (d *DialogClient) NotifyEvent(message string) {
}

func (d *DialogClient) CallResult(message string) {
}

func (d *DialogClient) OpenDialog(dialogOptions *dialog.OpenDialog, callbackID string) {
}

func (d *DialogClient) SaveDialog(dialogOptions *dialog.SaveDialog, callbackID string) {
}

func (d *DialogClient) MessageDialog(dialogOptions *dialog.MessageDialog, callbackID string) {

	osa, err := exec.LookPath("osascript")
	if err != nil {
		d.log.Info("MessageDialog unavailable (osascript not found)")
		return
	}

	var btns slicer.StringSlicer
	defaultButton := ""
	cancelButton := ""
	for index, btn := range dialogOptions.Buttons {
		btns.Add(strconv.Quote(btn))
		if btn == dialogOptions.DefaultButton {
			defaultButton = fmt.Sprintf("default button %d", index+1)
		}
		if btn == dialogOptions.CancelButton {
			cancelButton = fmt.Sprintf("cancel button %d", index+1)
		}
	}
	buttons := "{" + btns.Join(",") + "}"
	script := fmt.Sprintf("display dialog \"%s\" buttons %s %s %s with title \"%s\"", dialogOptions.Message, buttons, defaultButton, cancelButton, dialogOptions.Title)
	go func() {
		out, err := exec.Command(osa, "-e", script).Output()
		if err != nil {
			// Assume user has pressed cancel button
			if dialogOptions.CancelButton != "" {
				d.dispatcher.DispatchMessage("DM" + callbackID + "|" + dialogOptions.CancelButton)
				return
			}
			d.log.Error("Dialog had bad exit code. If this was a Cancel button, add 'CancelButton' to the dialog.MessageDialog struct. Error: %s", err.Error())
			d.dispatcher.DispatchMessage("DM" + callbackID + "|error - check logs")
			return
		}

		buttonPressed := strings.TrimSpace(strings.TrimPrefix(string(out), "button returned:"))
		d.dispatcher.DispatchMessage("DM" + callbackID + "|" + buttonPressed)
	}()
}

func (d *DialogClient) WindowSetTitle(title string) {
}

func (d *DialogClient) WindowShow() {
}

func (d *DialogClient) WindowHide() {
}

func (d *DialogClient) WindowCenter() {
}

func (d *DialogClient) WindowMaximise() {
}

func (d *DialogClient) WindowUnmaximise() {
}

func (d *DialogClient) WindowMinimise() {
}

func (d *DialogClient) WindowUnminimise() {
}

func (d *DialogClient) WindowPosition(x int, y int) {
}

func (d *DialogClient) WindowSize(width int, height int) {
}

func (d *DialogClient) WindowSetMinSize(width int, height int) {
}

func (d *DialogClient) WindowSetMaxSize(width int, height int) {
}

func (d *DialogClient) WindowFullscreen() {
}

func (d *DialogClient) WindowUnFullscreen() {
}

func (d *DialogClient) WindowSetColour(colour int) {
}

func (d *DialogClient) DarkModeEnabled(callbackID string) {
}

func (d *DialogClient) SetApplicationMenu(menuJSON string) {
}

func (d *DialogClient) SetTrayMenu(trayMenuJSON string) {
}

func (d *DialogClient) UpdateTrayMenuLabel(trayMenuJSON string) {
}

func (d *DialogClient) UpdateContextMenu(contextMenuJSON string) {
}
