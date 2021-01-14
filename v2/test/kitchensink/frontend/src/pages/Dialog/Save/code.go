package main

import (
	"github.com/wailsapp/wails/v2/pkg/options/dialog"
	"io/ioutil"

	"github.com/wailsapp/wails/v2"
)

type Notepad struct {
	runtime *wails.Runtime
}

func (n *Notepad) WailsInit(runtime *wails.Runtime) error {
	n.runtime = runtime
	return nil
}

// SaveNotes attempts to save the given notes to disk.
// Returns false if the user cancelled the save, true on
// successful save.
func (n *Notepad) SaveNotes(notes string) (bool, error) {

	selectedFile := n.runtime.Dialog.Save(&dialog.SaveDialog{
		DefaultFilename: "notes.md",
		Filters:         "*.md",
	})

	// Check if the user pressed cancel
	if selectedFile == "" {
		// Cancelled
		return false, nil
	}

	// Save notes
	err := ioutil.WriteFile(selectedFile, []byte(notes), 0700)
	if err != nil {
		return false, err
	}

	return true, nil
}
