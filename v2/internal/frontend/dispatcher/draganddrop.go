package dispatcher

import (
	"errors"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"strconv"
	"strings"
)

func (d *Dispatcher) processDragAndDropMessage(message string, sender frontend.Frontend) (string, error) {
	switch message[1] {
	case 'M':
		sl := strings.Split(message[2:], ":")
		if len(sl) < 2 {
			return "", errors.New("Invalid drag and drop motion Event Message: " + message)
		}
		var x, y int
		var err error
		x, err = strconv.Atoi(sl[0])
		if err != nil {
			return "", errors.New("Invalid drag and drop motion Event Message: " + message)
		}
		y, err = strconv.Atoi(sl[1])
		if err != nil {
			return "", errors.New("Invalid drag and drop motion Event Message: " + message)
		}
		go d.events.Emit("wails.dnd.motion", x, y)
	case 'D':
		sl := strings.Split(message[2:], "\n")
		if len(sl) < 1 {
			return "", errors.New("Invalid drag and drop drop Event Message: " + message)
		}
		go d.events.Emit("wails.dnd.drop", sl)
	}

	return "", nil
}
