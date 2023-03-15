package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
)

func (m *MessageProcessor) dialogErrorCallback(window *WebviewWindow, message string, dialogID *string, err error) {
	errorMsg := fmt.Sprintf(message, err)
	m.Error(errorMsg)
	msg := "_wails.dialogErrorCallback('" + *dialogID + "', " + strconv.Quote(errorMsg) + ");"
	window.ExecJS(msg)
}

func (m *MessageProcessor) dialogCallback(window *WebviewWindow, dialogID *string, result string, isJSON bool) {
	msg := fmt.Sprintf("_wails.dialogCallback('%s', %s, %v);", *dialogID, strconv.Quote(result), isJSON)
	window.ExecJS(msg)
}

func (m *MessageProcessor) processDialogMethod(method string, rw http.ResponseWriter, r *http.Request, window *WebviewWindow, params QueryParams) {

	args, err := params.Args()
	if err != nil {
		m.httpError(rw, "Unable to parse arguments: %s", err)
		return
	}
	dialogID := args.String("dialog-id")
	if dialogID == nil {
		m.Error("dialog-id is required")
		return
	}
	switch method {
	case "Info", "Warning", "Error", "Question":
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
		case "Info":
			dialog = globalApplication.InfoDialog()
		case "Warning":
			dialog = globalApplication.WarningDialog()
		case "Error":
			dialog = globalApplication.ErrorDialog()
		case "Question":
			dialog = globalApplication.QuestionDialog()
		}
		// TODO: Add support for attaching Message dialogs to windows
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
	case "OpenFile":
		var options OpenFileDialogOptions
		err := params.ToStruct(&options)
		if err != nil {
			m.httpError(rw, "Error parsing dialog options: %s", err.Error())
			return
		}
		dialog := globalApplication.OpenFileDialogWithOptions(&options)

		go func() {
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
				}
			} else {
				file, err := dialog.PromptForSingleSelection()
				if err != nil {
					m.dialogErrorCallback(window, "Error getting selection: %s", dialogID, err)
					return
				}
				m.dialogCallback(window, dialogID, file, false)
			}
		}()
		m.ok(rw)
	case "SaveFile":
		var options SaveFileDialogOptions
		err := params.ToStruct(&options)
		if err != nil {
			m.httpError(rw, "Error parsing dialog options: %s", err.Error())
			return
		}
		dialog := globalApplication.SaveFileDialogWithOptions(&options)

		go func() {
			file, err := dialog.PromptForSingleSelection()
			if err != nil {
				m.dialogErrorCallback(window, "Error getting selection: %s", dialogID, err)
				return
			}
			m.dialogCallback(window, dialogID, file, false)
		}()
		m.ok(rw)

	default:
		m.httpError(rw, "Unknown dialog method: %s", method)
	}

}
