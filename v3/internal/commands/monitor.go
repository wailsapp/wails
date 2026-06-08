package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wailsapp/wails/v3/internal/monitor"
	"github.com/wailsapp/wails/v3/internal/monitor/tui"
)

// MonitorOptions are the flags for `wails3 monitor`.
type MonitorOptions struct {
	Sock string `description:"Attach directly to this monitor unix socket path"`
	PID  int    `name:"pid" description:"Attach to the running app with this PID"`
	List bool   `name:"list" description:"List attachable running apps and exit"`
}

// Monitor attaches to a running Wails app's dev IPC monitor and shows a live
// TUI of binding calls and events crossing the IPC boundary.
//
// The target app must be running in dev mode with the monitor enabled, e.g.:
//
//	WAILS_MONITOR=1 wails3 dev
func Monitor(options *MonitorOptions) error {
	DisableFooter = true

	// Direct socket attach bypasses discovery.
	if options.Sock != "" {
		return runTUI(options.Sock, monitor.DiscoveryEntry{Name: "app", Sock: options.Sock})
	}

	entries, err := monitor.List()
	if err != nil {
		return fmt.Errorf("failed to scan for running apps: %w", err)
	}

	if options.List {
		if len(entries) == 0 {
			fmt.Println("No attachable Wails apps found. Start one with WAILS_MONITOR=1.")
			return nil
		}
		fmt.Printf("%-24s %-8s %s\n", "NAME", "PID", "SOCKET")
		for _, e := range entries {
			fmt.Printf("%-24s %-8d %s\n", e.Name, e.PID, e.Sock)
		}
		return nil
	}

	if len(entries) == 0 {
		return fmt.Errorf("no attachable Wails apps found.\n" +
			"Start your app in dev mode with the monitor enabled:\n" +
			"  WAILS_MONITOR=1 wails3 dev")
	}

	// Filter by PID if given.
	if options.PID != 0 {
		for _, e := range entries {
			if e.PID == options.PID {
				return runTUI(e.Sock, e)
			}
		}
		return fmt.Errorf("no running app with PID %d (try `wails3 monitor --list`)", options.PID)
	}

	// Exactly one → attach. Several → pick the first and tell the user how to
	// disambiguate (a full interactive picker can come later).
	if len(entries) == 1 {
		return runTUI(entries[0].Sock, entries[0])
	}

	fmt.Println("Multiple running apps found; attaching to the first.")
	fmt.Println("Use `wails3 monitor --pid <PID>` to choose. Available:")
	for _, e := range entries {
		fmt.Printf("  %-24s pid %d\n", e.Name, e.PID)
	}
	return runTUI(entries[0].Sock, entries[0])
}

func runTUI(sock string, entry monitor.DiscoveryEntry) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return tui.Run(ctx, entry)
}
