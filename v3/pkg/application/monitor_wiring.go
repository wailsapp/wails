package application

import (
	"os"

	"github.com/wailsapp/wails/v3/pkg/application/monitor"
)

// Monitor enablement environment variables.
const (
	// envMonitorEnable opts the IPC monitor in (dev mode only).
	envMonitorEnable = "WAILS_MONITOR"
	// envMonitorSock both opts in AND overrides the socket path.
	envMonitorSock = "WAILS_MONITOR_SOCK"
	// envMonitorDangerousProd allows the monitor to run in a production build.
	// Value must match dangerousProdValue exactly.
	envMonitorDangerousProd = "WAILS_MONITOR_DANGEROUS_PRODUCTION"
	dangerousProdValue      = "i-understand-this-exposes-ipc"
)

// monitorStop is set when the monitor is running, so it can be torn down on
// application shutdown.
var monitorStop func()

// startMonitor wires up the in-process IPC monitor if (and only if) the app is
// in dev mode AND the developer has explicitly opted in, OR the production
// escape hatch is set.
//
// Dev mode is detected via the App.isDebugMode flag (set from the debug build
// tag), the existing Wails convention.
func (a *App) startMonitor() {
	defer func() { _ = recover() }()

	sockOverride := os.Getenv(envMonitorSock)
	optedIn := sockOverride != "" || os.Getenv(envMonitorEnable) == "1"
	dangerousProd := os.Getenv(envMonitorDangerousProd) == dangerousProdValue

	if !optedIn && !dangerousProd {
		return
	}

	if !a.isDebugMode {
		if !dangerousProd {
			// Opted in but not in dev mode and no escape hatch: do nothing.
			return
		}
		a.Logger.Warn(
			"================================================================\n" +
				"  WAILS IPC MONITOR ENABLED IN A PRODUCTION BUILD\n" +
				"  All IPC traffic (binding calls, args, results, events) is\n" +
				"  being exposed over a local unix socket.\n" +
				"  This MUST NEVER be enabled for shipped/distributed apps.\n" +
				"  Unset WAILS_MONITOR_DANGEROUS_PRODUCTION to disable.\n" +
				"================================================================")
	}

	appName := a.options.Name
	if appName == "" {
		appName = "wails-app"
	}
	pid := os.Getpid()

	sockPath := sockOverride
	if sockPath == "" {
		sockPath = monitor.DefaultSocketPath(appName, pid)
	}

	sink, err := monitor.Start(monitor.Config{SocketPath: sockPath})
	if err != nil {
		a.Logger.Error("Failed to start IPC monitor", "error", err, "sock", sockPath)
		return
	}

	// Register the on-demand snapshot collector (windows, bindings, listeners)
	// and prime the window cache so the first describe has data to serve.
	monitor.SetDescribeFunc(collectSnapshot)
	scheduleWindowCacheRefresh()

	cleanupDiscovery, err := monitor.WriteDiscovery(appName, pid, sockPath)
	if err != nil {
		a.Logger.Warn("IPC monitor started but discovery file could not be written", "error", err)
		cleanupDiscovery = func() {}
	}

	monitorStop = func() {
		monitor.SetDescribeFunc(nil)
		sink.Stop()
		cleanupDiscovery()
		// Drop cached window state and release the refresh latch so a later
		// restart starts clean.
		monitorWindowCache.Store(nil)
		monitorRefreshInFlight.Store(false)
	}

	a.Logger.Info("IPC monitor listening", "sock", sockPath)
}

// stopMonitor tears down the monitor if it is running.
func (a *App) stopMonitor() {
	defer func() { _ = recover() }()
	if monitorStop != nil {
		monitorStop()
		monitorStop = nil
	}
}
