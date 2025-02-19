package startupshutdown

import (
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

func TestServiceStartupShutdown(t *testing.T) {
	var seq atomic.Int64

	services := []application.Service{
		svctest.Configure(&Service1{}, svctest.Config{Id: 0, T: t, Seq: &seq}),
		svctest.Configure(&Service2{}, svctest.Config{Id: 1, T: t, Seq: &seq}),
		svctest.Configure(&Service3{}, svctest.Config{Id: 2, T: t, Seq: &seq}),
		svctest.Configure(&Service4{}, svctest.Config{Id: 3, T: t, Seq: &seq}),
		svctest.Configure(&Service5{}, svctest.Config{Id: 4, T: t, Seq: &seq}),
		svctest.Configure(&Service6{}, svctest.Config{Id: 5, T: t, Seq: &seq}),
	}

	app := apptest.New(t, application.Options{
		Services: services[:3],
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
		if count := seq.Load(); count != int64(len(services)) {
			t.Errorf("Wrong startup call count: wanted %d, got %d", len(services), count)
		}
		seq.Store(0)
		app.Quit()
	})

	err := apptest.Run(t, app)
	if err != nil {
		t.Fatal(err)
	}

	if count := seq.Load(); count != int64(len(services)) {
		t.Errorf("Wrong shutdown call count: wanted %d, got %d", len(services), count)
	}

	bound := int64(len(services)) + 1
	validate(t, services[0], bound)
	validate(t, services[1], bound)
	validate(t, services[2], bound)
	validate(t, services[3], bound)
	validate(t, services[4], bound)
	validate(t, services[5], bound)
}

func validate(t *testing.T, svc application.Service, bound int64) {
	id := svc.Instance().(interface{ Id() int }).Id()
	startup := svc.Instance().(interface{ StartupSeq() int64 }).StartupSeq()
	shutdown := svc.Instance().(interface{ ShutdownSeq() int64 }).ShutdownSeq()

	if startup == 0 && shutdown == 0 {
		t.Errorf("Service #%d did not start nor shut down", id)
		return
	} else if startup == 0 {
		t.Errorf("Service #%d started, but did not shut down", id)
		return
	} else if shutdown == 0 {
		t.Errorf("Service #%d shut down, but did not start", id)
		return
	}

	if shutdown != bound-startup {
		t.Errorf("Wrong sequence numbers for service #%d: wanted either %d..%d or %d..%d, got %d..%d", id, startup, bound-startup, bound-shutdown, shutdown, startup, shutdown)
	}
}
