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

// Guards #6263 (promised in #2411): an OpenFile dialog that allows directories and multiple
// selection must offer multiple folders on Windows, as on macOS and Linux.
// It previously fell back to a single-folder picker.
func TestOpenDialogKindFor(t *testing.T) {
	cases := []struct {
		dirs, multi bool
		want        openDialogKind
	}{
		{false, false, openSingleFile},
		{false, true, openMultipleFiles},
		{true, false, openSingleFolder},
		{true, true, openMultipleFolders},
	}
	for _, c := range cases {
		if got := openDialogKindFor(c.dirs, c.multi); got != c.want {
			t.Errorf("openDialogKindFor(dirs=%v, multi=%v) = %v, want %v", c.dirs, c.multi, got, c.want)
		}
	}
}
