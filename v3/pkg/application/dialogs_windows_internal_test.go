//go:build windows

package application

import (
	"testing"

	"golang.org/x/sys/windows"
)

// Guards #4233: dialogs with an icon must keep their button configuration.
// The old implementation forced MB_OK alongside MB_USERICON, so a question
// dialog with an icon lost its Yes/No buttons.
func TestMessageDialogUserIconFlags(t *testing.T) {
	question := calculateMessageDialogFlags(MessageDialogOptions{DialogType: QuestionDialogType})
	got := messageDialogUserIconFlags(question)

	if got&windows.MB_YESNO != windows.MB_YESNO {
		t.Errorf("question dialog lost its Yes/No buttons: flags=%#x", got)
	}
	if got&windows.MB_USERICON == 0 {
		t.Errorf("user icon flag missing: flags=%#x", got)
	}

	info := calculateMessageDialogFlags(MessageDialogOptions{DialogType: InfoDialogType})
	got = messageDialogUserIconFlags(info)
	if got&windows.MB_ICONINFORMATION == windows.MB_ICONINFORMATION {
		t.Errorf("standard icon bits must be stripped when using MB_USERICON: flags=%#x", got)
	}
}
