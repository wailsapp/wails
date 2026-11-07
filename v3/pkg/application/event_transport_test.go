package application

import "testing"

func newEventTransportTestApp(t *testing.T, options Options) *App {
	t.Helper()
	previous := globalApplication
	globalApplication = nil
	a := New(options)
	t.Cleanup(func() {
		a.cancel()
		for _, shutdown := range a.shutdownTasks {
			shutdown()
		}
		globalApplication = previous
	})
	return a
}
