// +build windows

package runtime

import (
	"golang.org/x/sys/windows"
	"syscall"

	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/harry1453/go-common-file-dialog/cfdutil"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	dialogoptions "github.com/wailsapp/wails/v2/pkg/options/dialog"
)

// Dialog defines all Dialog related operations
type Dialog interface {
	OpenFile(dialogOptions *dialogoptions.OpenDialog) (string, error)
	OpenMultipleFiles(dialogOptions *dialogoptions.OpenDialog) ([]string, error)
	OpenDirectory(dialogOptions *dialogoptions.OpenDialog) (string, error)
	Save(dialogOptions *dialogoptions.SaveDialog) (string, error)
	Message(dialogOptions *dialogoptions.MessageDialog) (string, error)
}

// dialog exposes the Dialog interface
type dialog struct {
	bus *servicebus.ServiceBus
}

// newDialogs creates a new Dialogs struct
func newDialog(bus *servicebus.ServiceBus) Dialog {
	return &dialog{
		bus: bus,
	}
}

// processTitleAndFilter return the title and filter from the given params.
// title is the first string, filter is the second
func (r *dialog) processTitleAndFilter(params ...string) (string, string) {

	var title, filter string

	if len(params) > 0 {
		title = params[0]
	}

	if len(params) > 1 {
		filter = params[1]
	}

	return title, filter
}

func convertFilters(filters []dialogoptions.FileFilter) []cfd.FileFilter {
	var result []cfd.FileFilter
	for _, filter := range filters {
		result = append(result, cfd.FileFilter(filter))
	}
	return result
}

func pickMultipleFiles(options *dialogoptions.OpenDialog) ([]string, error) {

	results, err := cfdutil.ShowOpenMultipleFilesDialog(cfd.DialogConfig{
		Title:       options.Title,
		Role:        "OpenMultipleFiles",
		FileFilters: convertFilters(options.Filters),
		FileName:    options.DefaultFilename,
		Folder:      options.DefaultDirectory,
	})
	return results, err
}

func (r *dialog) OpenMultipleFiles(options *dialogoptions.OpenDialog) ([]string, error) {
	return pickMultipleFiles(options)
}

func (r *dialog) OpenDirectory(options *dialogoptions.OpenDialog) (string, error) {
	return cfdutil.ShowPickFolderDialog(cfd.DialogConfig{
		Title:  options.Title,
		Role:   "PickFolder",
		Folder: options.DefaultDirectory,
	})
}

func (r *dialog) OpenFile(options *dialogoptions.OpenDialog) (string, error) {
	result, err := cfdutil.ShowOpenFileDialog(cfd.DialogConfig{
		Folder:      options.DefaultDirectory,
		FileFilters: convertFilters(options.Filters),
		FileName:    options.DefaultFilename,
	})
	return result, err
}

// Save prompts the user to select a file
func (r *dialog) Save(options *dialogoptions.SaveDialog) (string, error) {

	result, err := cfdutil.ShowSaveFileDialog(cfd.DialogConfig{
		Title:       options.Title,
		Role:        "SaveFile",
		FileName:    options.DefaultFilename,
		FileFilters: convertFilters(options.Filters),
	})
	return result, err
}

// Message show a message to the user
func (r *dialog) Message(options *dialogoptions.MessageDialog) (string, error) {

	// TODO: error handling
	title, err := syscall.UTF16PtrFromString(options.Title)
	if err != nil {
		return "", err
	}
	message, err := syscall.UTF16PtrFromString(options.Message)
	if err != nil {
		return "", err
	}
	var flags uint32
	switch options.Type {
	case dialogoptions.InfoDialog:
		flags = windows.MB_OK | windows.MB_ICONINFORMATION
	case dialogoptions.ErrorDialog:
		flags = windows.MB_ICONERROR | windows.MB_OK
	case dialogoptions.QuestionDialog:
		flags = windows.MB_YESNO
	case dialogoptions.WarningDialog:
		flags = windows.MB_OK | windows.MB_ICONWARNING
	}

	result, _ := windows.MessageBox(0, message, title, flags|windows.MB_SYSTEMMODAL)
	if options.Type == dialogoptions.QuestionDialog {
		if result == 6 { // IDYES
			return "Yes", nil
		}
		if result == 7 { // IDNO
			return "No", nil
		}
	}
	return "", nil

}
