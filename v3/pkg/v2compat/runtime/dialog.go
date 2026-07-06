package runtime

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// DialogType mirrors the v2 runtime.DialogType type.
type DialogType string

// Dialog types, mirroring the v2 runtime constants.
const (
	InfoDialog     DialogType = "info"
	WarningDialog  DialogType = "warning"
	ErrorDialog    DialogType = "error"
	QuestionDialog DialogType = "question"
)

// FileFilter mirrors the v2 runtime.FileFilter type.
// v3 equivalent: application.FileFilter.
type FileFilter struct {
	DisplayName string // Filter information EG: "Image Files (*.jpg, *.png)"
	Pattern     string // semicolon separated list of extensions, EG: "*.jpg;*.png"
}

// OpenDialogOptions mirrors the v2 runtime.OpenDialogOptions type.
// v3 equivalent: application.OpenFileDialogOptions.
type OpenDialogOptions struct {
	DefaultDirectory           string
	DefaultFilename            string
	Title                      string
	Filters                    []FileFilter
	ShowHiddenFiles            bool
	CanCreateDirectories       bool
	ResolvesAliases            bool
	TreatPackagesAsDirectories bool
}

// SaveDialogOptions mirrors the v2 runtime.SaveDialogOptions type.
// v3 equivalent: application.SaveFileDialogOptions.
type SaveDialogOptions struct {
	DefaultDirectory           string
	DefaultFilename            string
	Title                      string
	Filters                    []FileFilter
	ShowHiddenFiles            bool
	CanCreateDirectories       bool
	TreatPackagesAsDirectories bool
}

// MessageDialogOptions mirrors the v2 runtime.MessageDialogOptions type.
// v3 equivalent: the application.MessageDialog builder API.
type MessageDialogOptions struct {
	Type          DialogType
	Title         string
	Message       string
	Buttons       []string
	DefaultButton string
	CancelButton  string
	Icon          []byte
}

// convertFilters maps v2 file filters onto their v3 equivalent.
func convertFilters(filters []FileFilter) []application.FileFilter {
	if len(filters) == 0 {
		return nil
	}
	result := make([]application.FileFilter, 0, len(filters))
	for _, filter := range filters {
		result = append(result, application.FileFilter{
			DisplayName: filter.DisplayName,
			Pattern:     filter.Pattern,
		})
	}
	return result
}

// OpenDirectoryDialog mirrors the v2 runtime.OpenDirectoryDialog function.
// v3 equivalent: app.Dialog.OpenFileWithOptions with CanChooseDirectories set.
func OpenDirectoryDialog(_ context.Context, dialogOptions OpenDialogOptions) (string, error) {
	a := app()
	if a == nil {
		return "", errNoApp
	}
	dialog := a.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
		CanChooseDirectories:            true,
		CanChooseFiles:                  false,
		CanCreateDirectories:            dialogOptions.CanCreateDirectories,
		ShowHiddenFiles:                 dialogOptions.ShowHiddenFiles,
		ResolvesAliases:                 dialogOptions.ResolvesAliases,
		TreatsFilePackagesAsDirectories: dialogOptions.TreatPackagesAsDirectories,
		Filters:                         convertFilters(dialogOptions.Filters),
		Title:                           dialogOptions.Title,
		Directory:                       dialogOptions.DefaultDirectory,
	})
	return dialog.PromptForSingleSelection()
}

// OpenFileDialog mirrors the v2 runtime.OpenFileDialog function. The v2
// DefaultFilename option is ignored as v3 open dialogs have no default filename.
// v3 equivalent: app.Dialog.OpenFileWithOptions.
func OpenFileDialog(_ context.Context, dialogOptions OpenDialogOptions) (string, error) {
	a := app()
	if a == nil {
		return "", errNoApp
	}
	dialog := a.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
		CanChooseFiles:                  true,
		CanCreateDirectories:            dialogOptions.CanCreateDirectories,
		ShowHiddenFiles:                 dialogOptions.ShowHiddenFiles,
		ResolvesAliases:                 dialogOptions.ResolvesAliases,
		TreatsFilePackagesAsDirectories: dialogOptions.TreatPackagesAsDirectories,
		Filters:                         convertFilters(dialogOptions.Filters),
		Title:                           dialogOptions.Title,
		Directory:                       dialogOptions.DefaultDirectory,
	})
	return dialog.PromptForSingleSelection()
}

// OpenMultipleFilesDialog mirrors the v2 runtime.OpenMultipleFilesDialog
// function. The v2 DefaultFilename option is ignored as v3 open dialogs have
// no default filename.
// v3 equivalent: app.Dialog.OpenFileWithOptions with AllowsMultipleSelection set.
func OpenMultipleFilesDialog(_ context.Context, dialogOptions OpenDialogOptions) ([]string, error) {
	a := app()
	if a == nil {
		return nil, errNoApp
	}
	dialog := a.Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
		CanChooseFiles:                  true,
		AllowsMultipleSelection:         true,
		CanCreateDirectories:            dialogOptions.CanCreateDirectories,
		ShowHiddenFiles:                 dialogOptions.ShowHiddenFiles,
		ResolvesAliases:                 dialogOptions.ResolvesAliases,
		TreatsFilePackagesAsDirectories: dialogOptions.TreatPackagesAsDirectories,
		Filters:                         convertFilters(dialogOptions.Filters),
		Title:                           dialogOptions.Title,
		Directory:                       dialogOptions.DefaultDirectory,
	})
	return dialog.PromptForMultipleSelection()
}

// SaveFileDialog mirrors the v2 runtime.SaveFileDialog function.
// v3 equivalent: app.Dialog.SaveFileWithOptions.
func SaveFileDialog(_ context.Context, dialogOptions SaveDialogOptions) (string, error) {
	a := app()
	if a == nil {
		return "", errNoApp
	}
	dialog := a.Dialog.SaveFileWithOptions(&application.SaveFileDialogOptions{
		CanCreateDirectories:            dialogOptions.CanCreateDirectories,
		ShowHiddenFiles:                 dialogOptions.ShowHiddenFiles,
		TreatsFilePackagesAsDirectories: dialogOptions.TreatPackagesAsDirectories,
		Filters:                         convertFilters(dialogOptions.Filters),
		Title:                           dialogOptions.Title,
		Directory:                       dialogOptions.DefaultDirectory,
		Filename:                        dialogOptions.DefaultFilename,
	})
	return dialog.PromptForSingleSelection()
}

// MessageDialog mirrors the v2 runtime.MessageDialog function. It shows a
// message dialog and blocks until a button is clicked, returning the label of
// the clicked button.
//
// Unlike v2 there is no per-platform default-button fallback: dismissing the
// dialog without clicking a button (where the platform allows it) only
// delivers the cancel button's label if the platform wires the dismissal to
// the cancel button.
// v3 equivalent: the app.Dialog.Info/Question/Warning/Error builder API.
func MessageDialog(_ context.Context, dialogOptions MessageDialogOptions) (string, error) {
	a := app()
	if a == nil {
		return "", errNoApp
	}
	var dialog *application.MessageDialog
	switch dialogOptions.Type {
	case QuestionDialog:
		dialog = a.Dialog.Question()
	case WarningDialog:
		dialog = a.Dialog.Warning()
	case ErrorDialog:
		dialog = a.Dialog.Error()
	default:
		dialog = a.Dialog.Info()
	}
	dialog.SetTitle(dialogOptions.Title)
	dialog.SetMessage(dialogOptions.Message)
	if len(dialogOptions.Icon) > 0 {
		dialog.SetIcon(dialogOptions.Icon)
	}

	buttons := dialogOptions.Buttons
	if len(buttons) == 0 {
		buttons = []string{"Ok"}
	}

	result := make(chan string, 1)
	for _, label := range buttons {
		button := dialog.AddButton(label)
		button.OnClick(func() {
			select {
			case result <- label:
			default:
			}
		})
		if label == dialogOptions.DefaultButton {
			dialog.SetDefaultButton(button)
		}
		if label == dialogOptions.CancelButton {
			dialog.SetCancelButton(button)
		}
	}

	dialog.Show()
	return <-result, nil
}
