package subsystem

import (
	"os"
	"sync"
	"testing"

	"github.com/matryer/is"
	"github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/messagedispatcher/message"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

func TestSingleTopic(t *testing.T) {

	is := is.New(t)

	var expected string = "I am a message!"
	var actual string

	var wg sync.WaitGroup

	// Create new bus
	myLogger := logger.New(os.Stdout)
	myLogger.SetLogLevel(logger.TRACE)
	bus := servicebus.New(myLogger)
	eventSubsystem, _ := NewEvent(bus, myLogger)
	eventSubsystem.Start()

	eventSubsystem.RegisterListener("test", func(data ...interface{}) {
		is.Equal(len(data), 1)
		actual = data[0].(string)
		wg.Done()
	})

	wg.Add(1)

	eventMessage := &message.EventMessage{
		Name: "test",
		Data: []interface{}{"I am a message!"},
	}

	bus.Start()
	bus.Publish("event:test:from:j", eventMessage)
	wg.Wait()
	bus.Stop()

	is.Equal(actual, expected)

}
