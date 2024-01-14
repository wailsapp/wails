package dispatcher

import (
	"errors"
	"strings"
)

func (d *Dispatcher) processDragAndDropMessage(message string) (string, error) {
	switch message[1] {
	case 'D':
		sl := strings.Split(message[2:], "\n")
		if len(sl) < 1 {
			return "", errors.New("Invalid drag and drop drop Event Message: " + message)
		}
		d.events.Emit("wails.dnd.drop", sl)
	default:
		return "", errors.New("Invalid drag and drop drop Event Message: " + message)
	}

	return "", nil
}
