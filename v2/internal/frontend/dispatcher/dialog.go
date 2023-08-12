package dispatcher

import (
	"encoding/json"
	"github.com/mooijtech/wails/v2/internal/frontend"
	"github.com/pkg/errors"
)

func (d *Dispatcher) processDialogMessage(message string, sender frontend.Frontend) (any, error) {
	if len(message) < 2 {
		return "", errors.New("Invalid Dialog Message: " + message)
	}

	switch message[1:4] {
	case "OMF":
		// OpenMultipleFilesDialog
		var dialogOptions frontend.OpenDialogOptions

		if err := json.Unmarshal([]byte(message[5:]), &dialogOptions); err != nil {
			return "", errors.WithStack(err)
		}

		return sender.OpenMultipleFilesDialog(dialogOptions)
	case "OMD":
		// OpenMultipleDirectoriesDialog
		var dialogOptions frontend.OpenDialogOptions

		if err := json.Unmarshal([]byte(message[5:]), &dialogOptions); err != nil {
			return "", errors.WithStack(err)
		}

		return sender.OpenMultipleDirectoriesDialog(dialogOptions)
	}

	switch message[1:3] {
	case "OD":
		// OpenDirectoryDialog
		var dialogOptions frontend.OpenDialogOptions

		if err := json.Unmarshal([]byte(message[4:]), &dialogOptions); err != nil {
			return "", errors.WithStack(err)
		}

		return sender.OpenDirectoryDialog(dialogOptions)
	case "OF":
		// OpenFileDialog
		var dialogOptions frontend.OpenDialogOptions

		if err := json.Unmarshal([]byte(message[4:]), &dialogOptions); err != nil {
			return "", errors.WithStack(err)
		}

		return sender.OpenFileDialog(dialogOptions)
	case "SF":
		// SaveFileDialog
		var dialogOptions frontend.SaveDialogOptions

		if err := json.Unmarshal([]byte(message[4:]), &dialogOptions); err != nil {
			return "", errors.WithStack(err)
		}

		return sender.SaveFileDialog(dialogOptions)
	}

	switch message[2] {
	case 'M':
		// MessageDialog
		var dialogOptions frontend.MessageDialogOptions

		if err := json.Unmarshal([]byte(message[3:]), &dialogOptions); err != nil {
			return "", errors.WithStack(err)
		}

		return sender.MessageDialog(dialogOptions)
	}

	d.log.Error("unknown Dialog message: %s", message)

	return "", nil
}
