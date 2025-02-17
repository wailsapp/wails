package shutdownerror

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
	Service1 struct{ svctest.Shutdowner }
	Service2 struct{ svctest.Shutdowner }
	Service3 struct{ svctest.Shutdowner }
	Service4 struct{ svctest.Shutdowner }
	Service5 struct{ svctest.Shutdowner }
	Service6 struct{ svctest.Shutdowner }
)

func TestServiceShutdownError(t *testing.T) {
	var seq atomic.Int64

	services := []application.Service{
		svctest.Configure(&Service1{}, svctest.Config{Id: 0, T: t, Seq: &seq}),
		svctest.Configure(&Service2{}, svctest.Config{Id: 1, T: t, Seq: &seq, ShutdownErr: true}),
		svctest.Configure(&Service3{}, svctest.Config{Id: 2, T: t, Seq: &seq}),
		svctest.Configure(&Service4{}, svctest.Config{Id: 3, T: t, Seq: &seq}),
		svctest.Configure(&Service5{}, svctest.Config{Id: 4, T: t, Seq: &seq, ShutdownErr: true}),
		svctest.Configure(&Service6{}, svctest.Config{Id: 5, T: t, Seq: &seq, ShutdownErr: true}),
	}

	expectedShutdownErrors := []int{5, 4, 1}
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
		app.Quit()
	})

	err := apptest.Run(t, app)
	if err != nil {
		t.Fatal(err)
	}

	if ec := errCount.Load(); ec != int64(len(expectedShutdownErrors)) {
		t.Errorf("Wrong shutdown error count: wanted %d, got %d", len(expectedShutdownErrors), ec)
	}

	if count := seq.Load(); count != int64(len(services)) {
		t.Errorf("Wrong shutdown call count: wanted %d, got %d", len(services), count)
	}

	validate(t, services[0], 5)
	validate(t, services[1], 4)
	validate(t, services[2], 2, 3)
	validate(t, services[3], 1)
	validate(t, services[4], 1)
	validate(t, services[5], 0)
}

func validate(t *testing.T, svc application.Service, prev ...int64) {
	id := svc.Instance().(interface{ Id() int }).Id()
	seq := svc.Instance().(interface{ ShutdownSeq() int64 }).ShutdownSeq()

	if seq == 0 {
		t.Errorf("Service #%d did not shut down", id)
		return
	}

	for _, p := range prev {
		if seq <= p {
			t.Errorf("Wrong shutdown sequence number for service #%d: wanted >%d, got %d", id, p, seq)
		}
	}
}
