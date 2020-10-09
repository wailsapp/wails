package subsystem

import (
	"strings"

	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Log is the Logging subsystem. It handles messages with topics starting
// with "log:"
type Log struct {
	logChannel  <-chan *servicebus.Message
	quitChannel <-chan *servicebus.Message
	running     bool

	// Logger!
	logger *logger.Logger
}

// NewLog creates a new log subsystem
func NewLog(bus *servicebus.ServiceBus, logger *logger.Logger) (*Log, error) {

	// Subscribe to log messages
	logChannel, err := bus.Subscribe("log")
	if err != nil {
		return nil, err
	}

	// Subscribe to quit messages
	quitChannel, err := bus.Subscribe("quit")
	if err != nil {
		return nil, err
	}

	result := &Log{
		logChannel:  logChannel,
		quitChannel: quitChannel,
		logger:      logger,
	}

	return result, nil
}

// Start the subsystem
func (l *Log) Start() error {

	l.running = true

	// Spin off a go routine
	go func() {
		for l.running {
			select {
			case <-l.quitChannel:
				l.running = false
				break
			case logMessage := <-l.logChannel:
				logType := strings.TrimPrefix(logMessage.Topic(), "log:")
				switch logType {
				case "trace":
					l.logger.Trace(logMessage.Data().(string))
				case "debug":
					l.logger.Debug(logMessage.Data().(string))
				case "info":
					l.logger.Info(logMessage.Data().(string))
				case "warning":
					l.logger.Warning(logMessage.Data().(string))
				case "error":
					l.logger.Error(logMessage.Data().(string))
				case "fatal":
					l.logger.Fatal(logMessage.Data().(string))
				default:
					l.logger.Error("unknown log message: %+v", logMessage)
				}
			}
		}
		l.logger.Trace("Logger Shutdown")
	}()

	return nil
}
