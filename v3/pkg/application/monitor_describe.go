package application

import (
	"os"
	"runtime"
	"sort"
	"strings"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/internal/assetserver"
	"github.com/wailsapp/wails/v3/pkg/application/monitor"
)

// Window accessors (Position/Size/GetScreen/...) touch the native UI and MUST
// run on the GTK/Cocoa main thread. The monitor's describe requests arrive on a
// socket goroutine and can poll once per second, so we cannot block that
// goroutine on InvokeSync to fetch them: doing so holds the main loop hostage
// every tick (it stops answering the compositor's ping → "application not
// responding") and risks a hard deadlock if the native call ever waits on the
// caller. Instead we serve the most recent cached window state immediately and
// queue a single coalesced async refresh on the main thread for next time.
var (
	monitorWindowCache     atomic.Pointer[[]monitor.WindowInfo]
	monitorRefreshInFlight atomic.Bool
)

// scheduleWindowCacheRefresh queues exactly one main-thread refresh of the
// window cache. If a refresh is already pending it is a no-op, so a fast-polling
// consumer can never stack native work onto the main loop.
func scheduleWindowCacheRefresh() {
	if globalApplication == nil {
		return
	}
	if !monitorRefreshInFlight.CompareAndSwap(false, true) {
		return
	}
	InvokeAsync(func() {
		defer monitorRefreshInFlight.Store(false)
		windows := collectWindows()
		monitorWindowCache.Store(&windows)
	})
}

// collectSnapshot builds a point-in-time description of the app for the IPC
// monitor: windows, bound methods, backend event listeners, and metadata. It is
// registered with the monitor package (which cannot import this one) via
// monitor.SetDescribeFunc in the wiring.
func collectSnapshot() *monitor.Snapshot {
	defer func() { _ = recover() }()

	snap := &monitor.Snapshot{
		App: monitor.AppInfo{
			PID:       os.Getpid(),
			Platform:  runtime.GOOS,
			DebugMode: globalApplication.isDebugMode,
		},
	}
	if globalApplication != nil {
		snap.App.Name = globalApplication.options.Name
		snap.App.DevServer = assetserver.GetDevServerURL()
		// Serve the cached window state (collected on the main thread by a prior
		// refresh) without blocking, then queue the next refresh. Bindings and
		// listeners are plain Go maps and are safe to read off-thread.
		if cached := monitorWindowCache.Load(); cached != nil {
			snap.Windows = *cached
		}
		scheduleWindowCacheRefresh()
		snap.Bindings = collectBindings()
		snap.Listeners = collectListeners()
	}
	return snap
}

func collectWindows() []monitor.WindowInfo {
	if globalApplication == nil {
		return nil
	}
	windows := globalApplication.Window.GetAll()
	out := make([]monitor.WindowInfo, 0, len(windows))
	for _, w := range windows {
		if w == nil {
			continue
		}
		x, y := w.Position()
		width, height := w.Size()
		relX, relY := w.RelativePosition()
		wi := monitor.WindowInfo{
			ID:          w.ID(),
			Name:        w.Name(),
			Width:       width,
			Height:      height,
			X:           x,
			Y:           y,
			RelX:        relX,
			RelY:        relY,
			Focused:     w.IsFocused(),
			Fullscreen:  w.IsFullscreen(),
			Maximised:   w.IsMaximised(),
			Minimised:   w.IsMinimised(),
			Resizable:   w.Resizable(),
			IgnoreMouse: w.IsIgnoreMouseEvents(),
			Zoom:        w.GetZoom(),
		}
		if b := w.GetBorderSizes(); b != nil {
			wi.Border = monitor.Insets{Left: b.Left, Right: b.Right, Top: b.Top, Bottom: b.Bottom}
		}
		if scr, err := w.GetScreen(); err == nil && scr != nil {
			wi.Screen = &monitor.ScreenInfo{
				Name:        scr.Name,
				ID:          scr.ID,
				IsPrimary:   scr.IsPrimary,
				ScaleFactor: scr.ScaleFactor,
				Rotation:    scr.Rotation,
				BoundsX:     scr.Bounds.X,
				BoundsY:     scr.Bounds.Y,
				BoundsW:     scr.Bounds.Width,
				BoundsH:     scr.Bounds.Height,
				WorkW:       scr.WorkArea.Width,
				WorkH:       scr.WorkArea.Height,
			}
		}
		out = append(out, wi)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func collectBindings() []monitor.BindingInfo {
	if globalApplication == nil || globalApplication.bindings == nil {
		return nil
	}
	b := globalApplication.bindings
	out := make([]monitor.BindingInfo, 0, len(b.boundMethods))
	for _, m := range b.boundMethods {
		if m == nil {
			continue
		}
		out = append(out, monitor.BindingInfo{
			FQN:     m.FQN,
			Name:    m.Name,
			ID:      m.ID,
			Service: serviceFromFQN(m.FQN),
			Comment: strings.TrimSpace(m.Comments),
			Inputs:  toParams(m.Inputs),
			Outputs: toParams(m.Outputs),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].FQN < out[j].FQN })
	return out
}

func toParams(ps []*Parameter) []monitor.ParamInfo {
	if len(ps) == 0 {
		return nil
	}
	out := make([]monitor.ParamInfo, 0, len(ps))
	for _, p := range ps {
		if p == nil {
			continue
		}
		out = append(out, monitor.ParamInfo{Name: p.Name, Type: p.TypeName})
	}
	return out
}

// serviceFromFQN derives a service label from "pkg.Type.Method".
func serviceFromFQN(fqn string) string {
	if i := strings.LastIndex(fqn, "."); i > 0 {
		return fqn[:i]
	}
	return ""
}

func collectListeners() []monitor.EventListenerInfo {
	if globalApplication == nil || globalApplication.customEventProcessor == nil {
		return nil
	}
	ep := globalApplication.customEventProcessor
	ep.notifyLock.RLock()
	out := make([]monitor.EventListenerInfo, 0, len(ep.listeners))
	for name, ls := range ep.listeners {
		out = append(out, monitor.EventListenerInfo{Name: name, Count: len(ls)})
	}
	ep.notifyLock.RUnlock()
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
