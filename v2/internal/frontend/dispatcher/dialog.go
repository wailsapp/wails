package dispatcher

import (
	"encoding/json"
	"github.com/mooijtech/wails/v2/internal/frontend"
	"github.com/pkg/errors"
)

func (d *Dispatcher) processDialogMessage(message string, sender frontend.Frontend) (string, error) {
	if len(message) < 2 {
		return "", errors.New("Invalid Dialog Message: " + message)
	}

	switch message[1:] {
	case "OD":
		// OpenDirectoryDialog
		var dialogOptions frontend.OpenDialogOptions

		if err := json.Unmarshal([]byte(message[4:]), &dialogOptions); err != nil {
			return "", errors.WithStack(err)
		}

		return sender.OpenDirectoryDialog(dialogOptions)
	case "OMD":
		// OpenMultipleDirectoriesDialog
		var dialogOptions frontend.OpenDialogOptions

		if err := json.Unmarshal([]byte(message[5:]), &dialogOptions); err != nil {
			return "", errors.WithStack(err)
		}

		results, err := sender.OpenMultipleDirectoriesDialog(dialogOptions)

		if err != nil {
			return "", errors.WithStack(err)
		}

		resultsJSON, err := json.Marshal(results)

		if err != nil {
			return "", errors.WithStack(err)
		}

		return string(resultsJSON), nil
	case "OF":
		// OpenFileDialog
		var dialogOptions frontend.OpenDialogOptions

		if err := json.Unmarshal([]byte(message[4:]), &dialogOptions); err != nil {
			return "", errors.WithStack(err)
		}

		return sender.OpenFileDialog(dialogOptions)
	case "OMF":
		// OpenMultipleFilesDialog
		var dialogOptions frontend.OpenDialogOptions

		if err := json.Unmarshal([]byte(message[5:]), &dialogOptions); err != nil {
			return "", errors.WithStack(err)
		}

		results, err := sender.OpenMultipleFilesDialog(dialogOptions)

		if err != nil {
			return "", errors.WithStack(err)
		}

		resultsJSON, err := json.Marshal(results)

		if err != nil {
			return "", errors.WithStack(err)
		}

		return string(resultsJSON), nil
	case "SF":
		// SaveFileDialog
		var dialogOptions frontend.SaveDialogOptions

		if err := json.Unmarshal([]byte(message[4:]), &dialogOptions); err != nil {
			return "", errors.WithStack(err)
		}

		return sender.SaveFileDialog(dialogOptions)
	case "M":
		// MessageDialog
		var dialogOptions frontend.MessageDialogOptions

		if err := json.Unmarshal([]byte(message[3:]), &dialogOptions); err != nil {
			return "", errors.WithStack(err)
		}

		return sender.MessageDialog(dialogOptions)
	default:
		d.log.Error("unknown Dialog message: %s", message)
		return "", nil
	}
}
