package dispatcher

import (
	"github.com/pkg/errors"
	"github.com/wailsapp/wails/v2/internal/logger"
	pkgLogger "github.com/wailsapp/wails/v2/pkg/logger"
)

var logLevelMap = map[byte]logger.LogLevel{
	'1': pkgLogger.TRACE,
	'2': pkgLogger.DEBUG,
	'3': pkgLogger.INFO,
	'4': pkgLogger.WARNING,
	'5': pkgLogger.ERROR,
}

func (d *Dispatcher) processLogMessage(message string) (string, error) {
	if len(message) < 3 {
		return "", errors.New("Invalid Log Message: " + message)
	}

	messageText := message[2:]

	switch message[1] {
	case 'T':
		d.log.Trace(messageText)
	case 'P':
		d.log.Print(messageText)
	case 'D':
		d.log.Debug(messageText)
	case 'I':
		d.log.Info(messageText)
	case 'W':
		d.log.Warning(messageText)
	case 'E':
		d.log.Error(messageText)
	case 'F':
		d.log.Fatal(messageText)
	case 'S':
		loglevel, exists := logLevelMap[message[2]]
		if !exists {
			return "", errors.New("Invalid Set Log Level Message: " + message)
		}
		d.log.SetLogLevel(loglevel)
	default:
		return "", errors.New("Invalid Log Message: " + message)
	}
	return "", nil
}
