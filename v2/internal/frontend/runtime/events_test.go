package runtime_test

import (
	"fmt"
	"github.com/wailsapp/wails/v2/internal/frontend/runtime"
	"sync"
	"testing"
)
import "github.com/matryer/is"

type mockLogger struct {
	Log string
}

func (t *mockLogger) Trace(format string, args ...interface{}) {
	t.Log = fmt.Sprintf(format, args...)
}

func Test_EventsOn(t *testing.T) {
	i := is.New(t)
	l := &mockLogger{}
	manager := runtime.NewEvents(l)

	// Test On
	eventName := "test"
	counter := 0
	var wg sync.WaitGroup
	wg.Add(1)
	manager.On(eventName, func(args ...interface{}) {
		// This is called in a goroutine
		counter++
		wg.Done()
	})
	manager.Emit(eventName, "test payload")
	wg.Wait()
	i.Equal(1, counter)

}
