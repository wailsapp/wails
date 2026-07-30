package tests

import (
	"errors"
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
)

var appChan chan *application.App = make(chan *application.App, 1)
var errChan chan error = make(chan error, 1)
var endChan chan error = make(chan error, 1)

func init() { runtime.LockOSThread() }

func New(t *testing.T, options application.Options) *application.App {
	var app *application.App

	app = application.Get()
	if app != nil {
		return app
	}

	if options.Name == "" {
		options.Name = t.Name()
	}

	errorHandler := options.ErrorHandler
	var fatalErr error
	var fatalErrLock sync.Mutex
	options.ErrorHandler = func(err error) {
		if fatal := (*application.FatalError)(nil); errors.As(err, &fatal) {
			fatalErrLock.Lock()
			fatalErr = err
			fatalErrLock.Unlock()
		} else if errorHandler != nil {
			errorHandler(err)
		} else {
			app.Logger.Error(err.Error())
		}
	}

	postShutdown := options.PostShutdown
	options.PostShutdown = func() {
		if postShutdown != nil {
			postShutdown()
		}
		fatalErrLock.Lock()
		shutdownErr := fatalErr
		fatalErrLock.Unlock()
		endChan <- shutdownErr
		select {} // Block forever
	}

	return application.New(options)
}

func Run(t *testing.T, app *application.App) error {
	appChan <- app
	select {
	case err := <-errChan:
		return err
	case fatal := <-endChan:
		return fatal
	}
}

func Main(m *testing.M) {
	go func() {
		os.Exit(m.Run())
	}()

	errChan <- (<-appChan).Run()
	select {} // Block forever
}
