package runtime

import (
	"context"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// NewLifecycleService bridges the v2 application lifecycle hooks (OnStartup,
// OnDomReady, OnShutdown) onto a v3 service. The wails3 migrate command
// registers it as the last service of a migrated application.
func NewLifecycleService(onStartup, onDomReady, onShutdown func(context.Context)) application.Service {
	return application.NewService(&lifecycleService{
		onStartup:  onStartup,
		onDomReady: onDomReady,
		onShutdown: onShutdown,
	})
}

// lifecycleService implements application.ServiceStartup and
// application.ServiceShutdown to drive the v2 lifecycle hooks.
type lifecycleService struct {
	onStartup    func(context.Context)
	onDomReady   func(context.Context)
	onShutdown   func(context.Context)
	ctx          context.Context
	domReadyOnce sync.Once
}

// ServiceStartup calls the OnStartup hook and arranges for the OnDomReady
// hook to run exactly once when the first window signals that the runtime
// is ready.
func (s *lifecycleService) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	s.ctx = ctx
	if s.onStartup != nil {
		s.onStartup(ctx)
	}
	if s.onDomReady != nil {
		hook := func(window application.Window) {
			window.OnWindowEvent(events.Common.WindowRuntimeReady, func(*application.WindowEvent) {
				s.domReadyOnce.Do(func() {
					s.onDomReady(ctx)
				})
			})
		}
		if a := app(); a != nil {
			for _, window := range a.Window.GetAll() {
				hook(window)
			}
			a.Window.OnCreate(hook)
		}
	}
	return nil
}

// ServiceShutdown calls the OnShutdown hook.
func (s *lifecycleService) ServiceShutdown() error {
	if s.onShutdown != nil {
		ctx := s.ctx
		if ctx == nil {
			ctx = context.Background()
		}
		s.onShutdown(ctx)
	}
	return nil
}
