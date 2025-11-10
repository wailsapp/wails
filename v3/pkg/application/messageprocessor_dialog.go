package application

import (
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

func (m *MessageProcessor) processDialogMethod(req *RuntimeRequest, window Window) (any, error) {
	args := req.Args.AsMap()

	switch req.Method {
	case DialogInfo, DialogWarning, DialogError, DialogQuestion:
		var options MessageDialogOptions
		err := req.Args.ToStruct(&options)
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

		resp := make(chan string, 1)
		for _, button := range options.Buttons {
			label := button.Label
			button.OnClick(func() {
				select {
				case resp <- label:
				default:
				}
			})
		}
		dialog.AddButtons(options.Buttons)
		dialog.Show()

		response := <-resp
		return response, nil

	case DialogOpenFile:
		var options OpenFileDialogOptions
		err := req.Args.ToStruct(&options)
		if err != nil {
			return nil, errs.WrapInvalidDialogCallErrorf(err, "error parsing dialog options")
		}
		var detached = args.Bool("Detached")
		if detached == nil || !*detached {
			options.Window = window
		}
		dialog := globalApplication.Dialog.OpenFileWithOptions(&options)

		if options.AllowsMultipleSelection {
			files, err := dialog.PromptForMultipleSelection()
			if err != nil {
				return nil, errs.WrapInvalidDialogCallErrorf(err, "Dialog.OpenFile failed: error getting selection")
			}

			return files, nil
		} else {
			file, err := dialog.PromptForSingleSelection()
			if err != nil {
				return nil, errs.WrapInvalidDialogCallErrorf(err, "Dialog.OpenFile failed, error getting selection")
			}
			return file, nil
		}

	case DialogSaveFile:
		var options SaveFileDialogOptions
		err := req.Args.ToStruct(&options)
		if err != nil {
			return nil, errs.WrapInvalidDialogCallErrorf(err, "error parsing dialog options")
		}
		var detached = args.Bool("Detached")
		if detached == nil || !*detached {
			options.Window = window
		}
		dialog := globalApplication.Dialog.SaveFileWithOptions(&options)

		file, err := dialog.PromptForSingleSelection()
		if err != nil {
			return nil, errs.WrapInvalidDialogCallErrorf(err, "Dialog.SaveFile failed: error getting selection")
		}
		return file, nil

	default:
		return nil, errs.NewInvalidDialogCallErrorf("unknown method: %d", req.Method)
	}
}
