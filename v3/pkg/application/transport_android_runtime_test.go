package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"testing/synctest"
)

type androidCancellationProbe struct {
	check func(context.Context) bool
}

func (p *androidCancellationProbe) Check(ctx context.Context) bool {
	return p.check(ctx)
}

func TestAndroidBindingCancellation(t *testing.T) {
	probe := &androidCancellationProbe{
		check: func(ctx context.Context) bool { return ctx.Err() == context.Canceled },
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	app := &App{Logger: logger, bindings: NewBindings(nil, nil), windows: make(map[uint]Window)}
	app.Window = newWindowManager(app)
	app.Window.Add(NewWindow(WebviewWindowOptions{Name: "one"}))
	app.Window.Add(NewWindow(WebviewWindowOptions{Name: "two"}))
	previous := globalApplication
	globalApplication = app
	t.Cleanup(func() { globalApplication = previous })
	if err := app.bindings.Add(NewService(probe)); err != nil {
		t.Fatal(err)
	}
	m := NewMessageProcessor(logger)

	prepare := func(object int, client, window string) func() string {
		const methodName = "github.com/wailsapp/wails/v3/pkg/application.androidCancellationProbe.Check"
		payload := fmt.Sprintf(
			`{"object":%d,"method":0,"windowName":%q,"clientId":%q,"args":{"call-id":"same","methodName":%q,"args":[]}}`,
			object, window, client, methodName,
		)
		return m.prepareAndroidRuntimeCall(payload)
	}
	run := func(t *testing.T, worker func() string, want string) {
		t.Helper()
		var response struct {
			OK    bool
			Data  json.RawMessage
			Error string
		}
		raw := worker()
		if err := json.Unmarshal([]byte(raw), &response); err != nil || !response.OK {
			t.Fatalf("runtime call failed: %s (%v)", raw, err)
		}
		if got := string(response.Data); got != want {
			t.Fatalf("response data = %s, want %s", got, want)
		}
	}

	t.Run("cancel worker runs before binding worker", func(t *testing.T) {
		binding := prepare(callRequest, "session", "")
		// Hold the binding worker until the cancellation worker has returned.
		run(t, prepare(cancelCallRequest, "session", ""), "{}")
		run(t, binding, "true")
	})

	t.Run("cancel while binding runs and repeat cancellation", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			entered := make(chan struct{})
			release := make(chan struct{})
			original := probe.check
			probe.check = func(ctx context.Context) bool {
				close(entered)
				<-release
				return ctx.Err() == context.Canceled
			}
			defer func() { probe.check = original }()

			binding := prepare(callRequest, "session", "one")
			result := make(chan string, 1)
			go func() { result <- binding() }()
			<-entered
			run(t, prepare(cancelCallRequest, "session", "one"), "{}")
			run(t, prepare(cancelCallRequest, "session", "one"), "{}")
			close(release)
			run(t, func() string { return <-result }, "true")
		})
	})

	t.Run("late cancellation leaves no entry", func(t *testing.T) {
		run(t, prepare(callRequest, "session", "one"), "false")
		run(t, prepare(cancelCallRequest, "session", "one"), "{}")
		run(t, prepare(cancelCallRequest, "session", "one"), "{}")
		run(t, prepare(callRequest, "session", "one"), "false")
	})

	t.Run("same call ID across clients and windows", func(t *testing.T) {
		target := prepare(callRequest, "session", "one")
		otherClient := prepare(callRequest, "other", "one")
		otherWindow := prepare(callRequest, "session", "two")
		run(t, prepare(cancelCallRequest, "session", "one"), "{}")
		run(t, target, "true")
		run(t, otherClient, "false")
		run(t, otherWindow, "false")
	})

	if len(m.runningCalls) != 0 || len(m.windowCalls) != 0 {
		t.Fatal("binding registrations remained after completion")
	}
}
