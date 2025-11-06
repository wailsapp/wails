package application

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/errs"
)

const (
	DialogInfo     = 0
	DialogWarning  = 1
	DialogError    = 2
	DialogQuestion = 3
	DialogOpenFile = 4
	DialogSaveFile = 5
)

var dialogMethodNames = map[int]string{
	DialogInfo:     "Info",
	DialogWarning:  "Warning",
	DialogError:    "Error",
	DialogQuestion: "Question",
	DialogOpenFile: "OpenFile",
	DialogSaveFile: "SaveFile",
}

func (m *MessageProcessor) dialogErrorCallback(window Window, message string, dialogID *string, err error) {
	m.Error(message, "error", err)
	window.DialogError(*dialogID, err.Error())
}

func (m *MessageProcessor) dialogCallback(window Window, dialogID *string, result string, isJSON bool) {
	window.DialogResponse(*dialogID, result, isJSON)
}

func (m *MessageProcessor) processDialogMethod(req *RuntimeRequest, window Window) (any, error) {
	args, err := req.Params.Args()
	if err != nil {
		return nil, errs.WrapInvalidDialogCallErrorf(err, "unable to parse arguments")
	}

	dialogID := args.String("dialog-id")
	if dialogID == nil {
		return nil, errs.NewInvalidDialogCallErrorf("missing argument 'dialog-id'")
	}

	switch req.Method {
	case DialogInfo, DialogWarning, DialogError, DialogQuestion:
		var options MessageDialogOptions
		err := params.ToStruct(&options)
		if err != nil {
			return nil, errs.WrapInvalidDialogCallErrorf(err, "error parsing dialog options")
		}
		if len(options.Buttons) == 0 {
			switch runtime.GOOS {
			case "darwin":
				options.Buttons = []*Button{{Label: "OK", IsDefault: true}}
			}
		}
		var dialog *MessageDialog
		switch req.Method {
		case DialogInfo:
			dialog = InfoDialog()
		case DialogWarning:
			dialog = WarningDialog()
		case DialogError:
			dialog = ErrorDialog()
		case DialogQuestion:
			dialog = QuestionDialog()
		}
		var detached = args.Bool("Detached")
		if detached == nil || !*detached {
			dialog.AttachToWindow(window)
		}

		dialog.SetTitle(options.Title)
		dialog.SetMessage(options.Message)
		for _, button := range options.Buttons {
			label := button.Label
			button.OnClick(func() {
				m.dialogCallback(window, dialogID, label, false)
			})
		}
		dialog.AddButtons(options.Buttons)
		dialog.Show()
		return unit, nil

	case DialogOpenFile:
		var options OpenFileDialogOptions
		err := params.ToStruct(&options)
		if err != nil {
			return nil, errs.WrapInvalidDialogCallErrorf(err, "error parsing dialog options")
		}
		var detached = args.Bool("Detached")
		if detached == nil || !*detached {
			options.Window = window
		}
		dialog := globalApplication.Dialog.OpenFileWithOptions(&options)

		go func() {
			defer handlePanic()
			if options.AllowsMultipleSelection {
				files, err := dialog.PromptForMultipleSelection()
				if err != nil {
					m.dialogErrorCallback(window, "Dialog.OpenFile failed", dialogID, fmt.Errorf("error getting selection: %w", err))
					return
				} else {
					result, err := json.Marshal(files)
					if err != nil {
						m.dialogErrorCallback(window, "Dialog.OpenFile failed", dialogID, fmt.Errorf("error marshaling files: %w", err))
						return
					}
					m.dialogCallback(window, dialogID, string(result), true)
				}
			} else {
				file, err := dialog.PromptForSingleSelection()
				if err != nil {
					m.dialogErrorCallback(window, "Dialog.OpenFile failed", dialogID, fmt.Errorf("error getting selection: %w", err))
					return
				}
				m.dialogCallback(window, dialogID, file, false)
			}
		}()
		return unit, nil

	case DialogSaveFile:
		var options SaveFileDialogOptions
		err := params.ToStruct(&options)
		if err != nil {
			return nil, errs.WrapInvalidDialogCallErrorf(err, "error parsing dialog options")
		}
		var detached = args.Bool("Detached")
		if detached == nil || !*detached {
			options.Window = window
		}
		dialog := globalApplication.Dialog.SaveFileWithOptions(&options)

		go func() {
			defer handlePanic()
			file, err := dialog.PromptForSingleSelection()
			if err != nil {
				m.dialogErrorCallback(window, "Dialog.SaveFile failed", dialogID, fmt.Errorf("error getting selection: %w", err))
				return
			}
			m.dialogCallback(window, dialogID, file, false)
		}()
		return unit, nil

	default:
		return nil, errs.NewInvalidDialogCallErrorf("unknown method: %d", req.Method)
	}
}
