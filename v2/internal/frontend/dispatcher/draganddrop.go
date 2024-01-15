package dispatcher

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

func (d *Dispatcher) processDragAndDropMessage(message string) (string, error) {
	switch message[1] {
	case 'H':
		log.Println(message)
		x, y, paths, err := formatDragData(message)
		if err != nil {
			return "", err
		}
		d.events.Emit("wails:file-drop-hover", x, y, paths)
	case 'D':
		x, y, paths, err := formatDragData(message)
		if err != nil {
			return "", err
		}
		d.events.Emit("wails:file-drop", x, y, paths)
	case 'C':
		d.events.Emit("wails:file-drop-cancelled")
	default:
		return "", errors.New("Invalid drag and drop drop Event Message: " + message)
	}

	return "", nil
}

func formatDragData(message string) (int, int, []string, error) {
	msg := strings.SplitN(message[3:], ":", 3)
	if len(msg) != 3 {
		return 0, 0, nil, errors.New("Invalid drag and drop drop Event Message: " + message)
	}
	paths := strings.Split(msg[2], "\n")
	if len(paths) < 1 {
		return 0, 0, nil, errors.New("Invalid drag and drop drop Event Message: " + message)
	}

	x, err := strconv.Atoi(msg[0])
	if err != nil {
		return 0, 0, nil, errors.New("Invalid x coordinate in drag and drop drop Event Message: " + message)
	}

	y, err := strconv.Atoi(msg[1])
	if err != nil {
		return 0, 0, nil, errors.New("Invalid y coordinate in drag and drop drop Event Message: " + message)
	}
	log.Println(x, y, paths)
	return x, y, paths, nil
}
