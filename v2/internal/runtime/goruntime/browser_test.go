package goruntime

import (
	"os"
	"testing"

	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/internal/servicebus"
)

func TestBrowserOpen(t *testing.T) {
	mylogger := logger.New(os.Stdout)
	myServiceBus := servicebus.New(mylogger)
	myRuntime := New(myServiceBus)
	myRuntime.Browser.Open("http://www.google.com")
}
