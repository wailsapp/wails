package runtime

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/fs"
)

// FileFilter defines a filter for dialog boxes
type FileFilter = frontend.FileFilter

// OpenDialogOptions contains the options for the OpenDialogOptions runtime method
type OpenDialogOptions = frontend.OpenDialogOptions

// SaveDialogOptions contains the options for the SaveDialog runtime method
type SaveDialogOptions = frontend.SaveDialogOptions

type DialogType = frontend.DialogType

const (
	InfoDialog     = frontend.InfoDialog
	WarningDialog  = frontend.WarningDialog
	ErrorDialog    = frontend.ErrorDialog
	QuestionDialog = frontend.QuestionDialog
)

// MessageDialogOptions contains the options for the Message dialogs, EG Info, Warning, etc runtime methods
type MessageDialogOptions = frontend.MessageDialogOptions

// OpenDirectoryDialog prompts the user to select a directory
func OpenDirectoryDialog(ctx context.Context, dialogOptions OpenDialogOptions) (string, error) {
	appFrontend := getFrontend(ctx)
	if dialogOptions.DefaultDirectory != "" {
		if !fs.DirExists(dialogOptions.DefaultDirectory) {
			return "", fmt.Errorf("default directory '%s' does not exist", dialogOptions.DefaultDirectory)
		}
	}
	return appFrontend.OpenDirectoryDialog(dialogOptions)
}

// OpenFileDialog prompts the user to select a file
func OpenFileDialog(ctx context.Context, dialogOptions OpenDialogOptions) (string, error) {
	appFrontend := getFrontend(ctx)
	if dialogOptions.DefaultDirectory != "" {
		if !fs.DirExists(dialogOptions.DefaultDirectory) {
			return "", fmt.Errorf("default directory '%s' does not exist", dialogOptions.DefaultDirectory)
		}
	}
	return appFrontend.OpenFileDialog(dialogOptions)
}

// OpenMultipleFilesDialog prompts the user to select a file
func OpenMultipleFilesDialog(ctx context.Context, dialogOptions OpenDialogOptions) ([]string, error) {
	appFrontend := getFrontend(ctx)
	if dialogOptions.DefaultDirectory != "" {
		if !fs.DirExists(dialogOptions.DefaultDirectory) {
			return nil, fmt.Errorf("default directory '%s' does not exist", dialogOptions.DefaultDirectory)
		}
	}
	return appFrontend.OpenMultipleFilesDialog(dialogOptions)
}

// SaveFileDialog prompts the user to select a file
func SaveFileDialog(ctx context.Context, dialogOptions SaveDialogOptions) (string, error) {
	appFrontend := getFrontend(ctx)
	if dialogOptions.DefaultDirectory != "" {
		if !fs.DirExists(dialogOptions.DefaultDirectory) {
			return "", fmt.Errorf("default directory '%s' does not exist", dialogOptions.DefaultDirectory)
		}
	}
	return appFrontend.SaveFileDialog(dialogOptions)
}

// MessageDialog show a message dialog to the user
func MessageDialog(ctx context.Context, dialogOptions MessageDialogOptions) (string, error) {
	appFrontend := getFrontend(ctx)
	return appFrontend.MessageDialog(dialogOptions)
}
