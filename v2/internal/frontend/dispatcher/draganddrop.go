package dispatcher

import (
	"errors"
	"strconv"
	"strings"
)

func (d *Dispatcher) processDragAndDropMessage(message string) (string, error) {
	switch message[1] {
	case 'D':
		msg := strings.SplitN(message[3:], ":", 3)
		if len(msg) != 3 {
			return "", errors.New("Invalid drag and drop Message: " + message)
		}

		x, err := strconv.Atoi(msg[0])
		if err != nil {
			return "", errors.New("Invalid x coordinate in drag and drop Message: " + message)
		}

		y, err := strconv.Atoi(msg[1])
		if err != nil {
			return "", errors.New("Invalid y coordinate in drag and drop Message: " + message)
		}

		paths := strings.Split(msg[2], "\n")
		if len(paths) < 1 {
			return "", errors.New("Invalid drag and drop Message: " + message)
		}

		d.events.Emit("wails:file-drop", x, y, paths)
	default:
		return "", errors.New("Invalid drag and drop Message: " + message)
	}

	return "", nil
}
