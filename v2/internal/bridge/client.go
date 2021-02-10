package bridge

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/pkg/options/dialog"
	"golang.org/x/sync/semaphore"
)

type BridgeClient struct {
	session         *session
	dialogSemaphore *semaphore.Weighted
}

func (b BridgeClient) Quit() {
	b.session.log.Info("Quit unsupported in Bridge mode")
}

func (b BridgeClient) NotifyEvent(message string) {
	//b.session.sendMessage("n" + message)
	b.session.log.Info("NotifyEvent: %s", message)
	b.session.log.Info("NotifyEvent unsupported in Bridge mode")
}

func (b BridgeClient) CallResult(message string) {
	b.session.sendMessage("c" + message)
}

func (b BridgeClient) OpenDialog(dialogOptions *dialog.OpenDialog, callbackID string) {
	b.session.log.Info("OpenDialog unsupported in Bridge mode")
}

func (b BridgeClient) SaveDialog(dialogOptions *dialog.SaveDialog, callbackID string) {
	b.session.log.Info("SaveDialog unsupported in Bridge mode")
}

func (b BridgeClient) MessageDialog(dialogOptions *dialog.MessageDialog, callbackID string) {

	osa, err := exec.LookPath("osascript")
	if err != nil {
		b.session.log.Info("MessageDialog unavailable (osascript not found)")
		return
	}

	var btns slicer.StringSlicer
	defaultButton := ""
	for index, btn := range dialogOptions.Buttons {
		btns.Add(strconv.Quote(btn))
		if btn == dialogOptions.DefaultButton {
			defaultButton = fmt.Sprintf("default button %d", index+1)
		}
	}
	buttons := "{" + btns.Join(",") + "}"
	script := fmt.Sprintf("display dialog \"%s\" buttons %s %s with title \"%s\"", dialogOptions.Message, buttons, defaultButton, dialogOptions.Title)

	b.session.log.Info("OSASCRIPT: %s", script)
	go func() {
		out, err := exec.Command(osa, "-e", script).Output()
		if err != nil {
			b.session.log.Error(err.Error())
			return
		}

		b.session.log.Info(string(out))
		buttonPressed := strings.TrimSpace(strings.TrimPrefix(string(out), "button returned:"))
		b.session.client.DispatchMessage("DM" + callbackID + "|" + buttonPressed)
	}()
}

func (b BridgeClient) WindowSetTitle(title string) {
	b.session.log.Info("WindowSetTitle unsupported in Bridge mode")
}

func (b BridgeClient) WindowShow() {
	b.session.log.Info("WindowShow unsupported in Bridge mode")
}

func (b BridgeClient) WindowHide() {
	b.session.log.Info("WindowHide unsupported in Bridge mode")
}

func (b BridgeClient) WindowCenter() {
	b.session.log.Info("WindowCenter unsupported in Bridge mode")
}

func (b BridgeClient) WindowMaximise() {
	b.session.log.Info("WindowMaximise unsupported in Bridge mode")
}

func (b BridgeClient) WindowUnmaximise() {
	b.session.log.Info("WindowUnmaximise unsupported in Bridge mode")
}

func (b BridgeClient) WindowMinimise() {
	b.session.log.Info("WindowMinimise unsupported in Bridge mode")
}

func (b BridgeClient) WindowUnminimise() {
	b.session.log.Info("WindowUnminimise unsupported in Bridge mode")
}

func (b BridgeClient) WindowPosition(x int, y int) {
	b.session.log.Info("WindowPosition unsupported in Bridge mode")
}

func (b BridgeClient) WindowSize(width int, height int) {
	b.session.log.Info("WindowSize unsupported in Bridge mode")
}

func (b BridgeClient) WindowSetMinSize(width int, height int) {
	b.session.log.Info("WindowSetMinSize unsupported in Bridge mode")
}

func (b BridgeClient) WindowSetMaxSize(width int, height int) {
	b.session.log.Info("WindowSetMaxSize unsupported in Bridge mode")
}

func (b BridgeClient) WindowFullscreen() {
	b.session.log.Info("WindowFullscreen unsupported in Bridge mode")
}

func (b BridgeClient) WindowUnFullscreen() {
	b.session.log.Info("WindowUnFullscreen unsupported in Bridge mode")
}

func (b BridgeClient) WindowSetColour(colour int) {
	b.session.log.Info("WindowSetColour unsupported in Bridge mode")
}

func (b BridgeClient) DarkModeEnabled(callbackID string) {
	b.session.log.Info("DarkModeEnabled unsupported in Bridge mode")
}

func (b BridgeClient) SetApplicationMenu(menuJSON string) {
	b.session.log.Info("SetApplicationMenu unsupported in Bridge mode")
}

func (b BridgeClient) SetTrayMenu(trayMenuJSON string) {
	b.session.sendMessage("TS" + trayMenuJSON)
}

func (b BridgeClient) UpdateTrayMenuLabel(trayMenuJSON string) {
	b.session.sendMessage("TS" + trayMenuJSON)
}

func (b BridgeClient) UpdateContextMenu(contextMenuJSON string) {
	b.session.log.Info("UpdateContextMenu unsupported in Bridge mode")
}

func newBridgeClient(session *session) *BridgeClient {
	return &BridgeClient{
		session: session,
	}
}
