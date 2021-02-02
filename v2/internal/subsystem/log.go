package subsystem

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/runtime"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// Log is the Logging subsystem. It handles messages with topics starting
// with "log:"
type Log struct {
	logChannel <-chan *servicebus.Message

	// quit flag
	shouldQuit bool

	// Logger!
	logger *logger.Logger

	// Loglevel store
	logLevelStore *runtime.Store

	// Context for shutdown
	ctx    context.Context
	cancel context.CancelFunc

	// internal waitgroup
	wg sync.WaitGroup
}

// NewLog creates a new log subsystem
func NewLog(bus *servicebus.ServiceBus, logger *logger.Logger, logLevelStore *runtime.Store) (*Log, error) {

	// Subscribe to log messages
	logChannel, err := bus.Subscribe("log")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	result := &Log{
		logChannel:    logChannel,
		logger:        logger,
		logLevelStore: logLevelStore,
		ctx:           ctx,
		cancel:        cancel,
	}

	return result, nil
}

// Start the subsystem
func (l *Log) Start() error {

	l.wg.Add(1)

	// Spin off a go routine
	go func() {
		defer l.logger.Trace("Logger Shutdown")

		for l.shouldQuit == false {
			select {
			case <-l.ctx.Done():
				l.wg.Done()
				return
			case logMessage := <-l.logChannel:
				logType := strings.TrimPrefix(logMessage.Topic(), "log:")
				switch logType {
				case "print":
					l.logger.Print(logMessage.Data().(string))
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
				case "setlevel":
					switch inLevel := logMessage.Data().(type) {
					case logger.LogLevel:
						l.logger.SetLogLevel(inLevel)
						l.logLevelStore.Set(inLevel)
					case string:
						uint64level, err := strconv.ParseUint(inLevel, 10, 8)
						if err != nil {
							l.logger.Error("Error parsing log level: %+v", inLevel)
							continue
						}
						level := logger.LogLevel(uint64level)
						l.logLevelStore.Set(level)
						l.logger.SetLogLevel(level)
					}

				default:
					l.logger.Error("unknown log message: %+v", logMessage)
				}
			}
		}
	}()

	return nil
}

func (l *Log) Close() {
	l.cancel()
	l.wg.Wait()
}
