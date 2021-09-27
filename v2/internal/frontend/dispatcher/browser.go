package dispatcher

import (
	"errors"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

// processBrowserMessage processing browser messages
func (d *Dispatcher) processBrowserMessage(message string, sender frontend.Frontend) (string, error) {
	if len(message) < 2 {
		return "", errors.New("Invalid Browser Message: " + message)
	}
	switch message[1] {
	case 'O':
		url := message[3:]
		go sender.BrowserOpenURL(url)
	default:
		d.log.Error("unknown Browser message: %s", message)
	}

	return "", nil
}
