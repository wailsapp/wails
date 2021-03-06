package subsystem

import (
	"context"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

// URL is the URL Handler subsystem. It handles messages with topics starting
// with "url:"
type URL struct {
	urlChannel <-chan *servicebus.Message

	// quit flag
	shouldQuit bool

	// Logger!
	logger *logger.Logger

	// Context for shutdown
	ctx    context.Context
	cancel context.CancelFunc

	// internal waitgroup
	wg sync.WaitGroup

	// Handlers
	handlers map[string]func(string)
}

// NewURL creates a new log subsystem
func NewURL(bus *servicebus.ServiceBus, logger *logger.Logger, handlers map[string]func(string)) (*URL, error) {

	// Subscribe to log messages
	urlChannel, err := bus.Subscribe("url")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	result := &URL{
		urlChannel: urlChannel,
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		handlers:   handlers,
	}

	return result, nil
}

// Start the subsystem
func (u *URL) Start() error {

	u.wg.Add(1)

	// Spin off a go routine
	go func() {
		defer u.logger.Trace("URL Shutdown")

		for u.shouldQuit == false {
			select {
			case <-u.ctx.Done():
				u.wg.Done()
				return
			case urlMessage := <-u.urlChannel:
				// Guard against nil messages
				if urlMessage == nil {
					continue
				}
				messageType := strings.TrimPrefix(urlMessage.Topic(), "url:")
				switch messageType {
				case "handler":
					url := urlMessage.Data().(string)
					splitURL := strings.Split(url, ":")
					protocol := splitURL[0]
					callback, ok := u.handlers[protocol]
					if ok {
						go callback(url)
					}
				default:
					u.logger.Error("unknown url message: %+v", urlMessage)
				}
			}
		}
	}()

	return nil
}

func (u *URL) Close() {
	u.cancel()
	u.wg.Wait()
}
