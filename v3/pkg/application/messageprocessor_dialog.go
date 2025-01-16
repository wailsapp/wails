package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
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
	errorMsg := fmt.Sprintf(message, err)
	m.Error(errorMsg)
	window.DialogError(*dialogID, errorMsg)
}

func (m *MessageProcessor) dialogCallback(window Window, dialogID *string, result string, isJSON bool) {
	window.DialogResponse(*dialogID, result, isJSON)
}

func (m *MessageProcessor) processDialogMethod(method int, rw http.ResponseWriter, r *http.Request, window Window, params QueryParams) {

	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Unable to parse arguments: %s", err.Error())
		return
	}
	dialogID := args.String("dialog-id")
	if dialogID == nil {
		m.Error("dialog-id is required")
		return
	}

	var methodName = "Dialog." + dialogMethodNames[method]

	switch method {
	case DialogInfo, DialogWarning, DialogError, DialogQuestion:
		var options MessageDialogOptions
		err := params.ToStruct(&options)
		if err != nil {
			m.dialogErrorCallback(window, "Error parsing dialog options: %s", dialogID, err)
			return
		}
		if len(options.Buttons) == 0 {
			switch runtime.GOOS {
			case "darwin":
				options.Buttons = []*Button{{Label: "OK", IsDefault: true}}
			}
		}
		var dialog *MessageDialog
		switch method {
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
		m.ok(rw)
		m.Info("Runtime Call:", "method", methodName, "options", options)

	case DialogOpenFile:
		var options OpenFileDialogOptions
		err := params.ToStruct(&options)
		if err != nil {
			m.httpError(rw, "Error parsing dialog options: %s", err.Error())
			return
		}
		var detached = args.Bool("Detached")
		if detached == nil || !*detached {
			options.Window = window.(*WebviewWindow)
		}
		dialog := OpenFileDialogWithOptions(&options)

		go func() {
			defer handlePanic()
			if options.AllowsMultipleSelection {
				files, err := dialog.PromptForMultipleSelection()
				if err != nil {
					m.dialogErrorCallback(window, "Error getting selection: %s", dialogID, err)
					return
				} else {
					result, err := json.Marshal(files)
					if err != nil {
						m.dialogErrorCallback(window, "Error marshalling files: %s", dialogID, err)
						return
					}
					m.dialogCallback(window, dialogID, string(result), true)
					m.Info("Runtime Call:", "method", methodName, "result", result)
				}
			} else {
				file, err := dialog.PromptForSingleSelection()
				if err != nil {
					m.dialogErrorCallback(window, "Error getting selection: %s", dialogID, err)
					return
				}
				m.dialogCallback(window, dialogID, file, false)
				m.Info("Runtime Call:", "method", methodName, "result", file)
			}
		}()
		m.ok(rw)
		m.Info("Runtime Call:", "method", methodName, "options", options)

	case DialogSaveFile:
		var options SaveFileDialogOptions
		err := params.ToStruct(&options)
		if err != nil {
			m.httpError(rw, "Error parsing dialog options: %s", err.Error())
			return
		}
		var detached = args.Bool("Detached")
		if detached == nil || !*detached {
			options.Window = window.(*WebviewWindow)
		}
		dialog := SaveFileDialogWithOptions(&options)

		go func() {
			defer handlePanic()
			file, err := dialog.PromptForSingleSelection()
			if err != nil {
				m.dialogErrorCallback(window, "Error getting selection: %s", dialogID, err)
				return
			}
			m.dialogCallback(window, dialogID, file, false)
			m.Info("Runtime Call:", "method", methodName, "result", file)
		}()
		m.ok(rw)
		m.Info("Runtime Call:", "method", methodName, "options", options)

	default:
		m.httpError(rw, "Unknown dialog method: %d", method)
	}

}
