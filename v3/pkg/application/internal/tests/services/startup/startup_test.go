package startup

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
	Service1 struct{ svctest.Startupper }
	Service2 struct{ svctest.Startupper }
	Service3 struct{ svctest.Startupper }
	Service4 struct{ svctest.Startupper }
	Service5 struct{ svctest.Startupper }
	Service6 struct{ svctest.Startupper }
)

func TestServiceStartup(t *testing.T) {
	var seq atomic.Int64

	services := []application.Service{
		svctest.Configure(&Service1{}, svctest.Config{Id: 0, T: t, Seq: &seq}),
		svctest.Configure(&Service2{}, svctest.Config{Id: 1, T: t, Seq: &seq, Options: application.ServiceOptions{
			Name: "I am service 2",
		}}),
		svctest.Configure(&Service3{}, svctest.Config{Id: 2, T: t, Seq: &seq, Options: application.ServiceOptions{
			Route: "/mounted/here",
		}}),
		svctest.Configure(&Service4{}, svctest.Config{Id: 3, T: t, Seq: &seq}),
		svctest.Configure(&Service5{}, svctest.Config{Id: 4, Seq: &seq, Options: application.ServiceOptions{
			Name:  "I am service 5",
			Route: "/mounted/there",
		}}),
		svctest.Configure(&Service6{}, svctest.Config{Id: 5, T: t, Seq: &seq, Options: application.ServiceOptions{
			Name:  "I am service 6",
			Route: "/mounted/elsewhere",
		}}),
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
		app.Quit()
	})

	err := apptest.Run(t, app)
	if err != nil {
		t.Fatal(err)
	}

	if count := seq.Load(); count != int64(len(services)) {
		t.Errorf("Wrong startup call count: wanted %d, got %d", len(services), count)
	}

	validate(t, services[0], 0)
	validate(t, services[1], 1)
	validate(t, services[2], 2)
	validate(t, services[3], 3)
	validate(t, services[4], 3)
	validate(t, services[5], 4, 5)
}

func validate(t *testing.T, svc application.Service, prev ...int64) {
	id := svc.Instance().(interface{ Id() int }).Id()
	seq := svc.Instance().(interface{ StartupSeq() int64 }).StartupSeq()

	if seq == 0 {
		t.Errorf("Service #%d did not start up", id)
		return
	}

	for _, p := range prev {
		if seq <= p {
			t.Errorf("Wrong startup sequence number for service #%d: wanted >%d, got %d", id, p, seq)
		}
	}
}
