package startupshutdownerror

import (
	"errors"
	"slices"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/wailsapp/wails/v3/pkg/application"
	apptest "github.com/wailsapp/wails/v3/pkg/application/internal/tests"
	svctest "github.com/wailsapp/wails/v3/pkg/application/internal/tests/services"
	"github.com/wailsapp/wails/v3/pkg/events"
)

func TestMain(m *testing.M) {
	apptest.Main(m)
}

type (
	Service1 struct{ svctest.StartupShutdowner }
	Service2 struct{ svctest.StartupShutdowner }
	Service3 struct{ svctest.StartupShutdowner }
	Service4 struct{ svctest.StartupShutdowner }
	Service5 struct{ svctest.StartupShutdowner }
	Service6 struct{ svctest.StartupShutdowner }
)

func TestServiceStartupShutdownError(t *testing.T) {
	var seq atomic.Int64

	services := []application.Service{
		svctest.Configure(&Service1{}, svctest.Config{Id: 0, T: t, Seq: &seq}),
		svctest.Configure(&Service2{}, svctest.Config{Id: 1, T: t, Seq: &seq, ShutdownErr: true}),
		svctest.Configure(&Service3{}, svctest.Config{Id: 2, T: t, Seq: &seq}),
		svctest.Configure(&Service4{}, svctest.Config{Id: 3, T: t, Seq: &seq, StartupErr: true, ShutdownErr: true}),
		svctest.Configure(&Service5{}, svctest.Config{Id: 4, T: t, Seq: &seq, StartupErr: true, ShutdownErr: true}),
		svctest.Configure(&Service6{}, svctest.Config{Id: 5, T: t, Seq: &seq, StartupErr: true, ShutdownErr: true}),
	}

	expectedShutdownErrors := []int{1}
	var errCount atomic.Int64

	var app *application.App
	app = apptest.New(t, application.Options{
		Services: services[:3],
		ErrorHandler: func(err error) {
			var mock *svctest.Error
			if !errors.As(err, &mock) {
				app.Logger.Error(err.Error())
				return
			}

			i := int(errCount.Add(1) - 1)
			if i < len(expectedShutdownErrors) && mock.Id == expectedShutdownErrors[i] {
				return
			}

			cut := min(i, len(expectedShutdownErrors))
			if slices.Contains(expectedShutdownErrors[:cut], mock.Id) {
				t.Errorf("Late or duplicate shutdown error for service #%d", mock.Id)
			} else if slices.Contains(expectedShutdownErrors[cut:], mock.Id) {
				t.Errorf("Early shutdown error for service #%d", mock.Id)
			} else {
				t.Errorf("Unexpected shutdown error for service #%d", mock.Id)
			}
		},
	})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		app.RegisterService(services[3])
		wg.Done()
	}()
	go func() {
		app.RegisterService(services[4])
		wg.Done()
	}()
	wg.Wait()

	app.RegisterService(services[5])

	app.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
		t.Errorf("Application started")
		app.Quit()
	})

	var mock *svctest.Error

	err := apptest.Run(t, app)
	if err != nil {
		if !errors.As(err, &mock) {
			t.Fatal(err)
		}
	}

	if mock == nil {
		t.Fatal("Wanted error for service #3 or #4, got none")
	} else if mock.Id != 3 && mock.Id != 4 {
		t.Errorf("Wanted error for service #3 or #4, got #%d", mock.Id)
	}

	if ec := errCount.Load(); ec != int64(len(expectedShutdownErrors)) {
		t.Errorf("Wrong shutdown error count: wanted %d, got %d", len(expectedShutdownErrors), ec)
	}

	if count := seq.Load(); count != 4+3 {
		t.Errorf("Wrong startup+shutdown call count: wanted %d+%d, got %d", 4, 3, count)
	}

	validate(t, services[0], true, true)
	validate(t, services[1], true, true)
	validate(t, services[2], true, true)
	validate(t, services[3], mock.Id == 3, false)
	validate(t, services[4], mock.Id == 4, false)
	validate(t, services[5], false, false)
}

func validate(t *testing.T, svc application.Service, startup bool, shutdown bool) {
	id := svc.Instance().(interface{ Id() int }).Id()
	startupSeq := svc.Instance().(interface{ StartupSeq() int64 }).StartupSeq()
	shutdownSeq := svc.Instance().(interface{ ShutdownSeq() int64 }).ShutdownSeq()

	if startup != (startupSeq != 0) {
		if startupSeq == 0 {
			t.Errorf("Service #%d did not start up", id)
		} else {
			t.Errorf("Unexpected startup for service #%d at seq=%d", id, startupSeq)
		}
	}

	if shutdown != (shutdownSeq != 0) {
		if shutdownSeq == 0 {
			t.Errorf("Service #%d did not shut down", id)
		} else {
			t.Errorf("Unexpected shutdown for service #%d at seq=%d", id, shutdownSeq)
		}
	}
}
