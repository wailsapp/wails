package commands

import (
	"context"
	"os"
	ossignal "os/signal"
	"syscall"

	"github.com/atterpac/refresh/engine"
	"github.com/wailsapp/wails/v3/internal/monitor/tui"
	"github.com/wailsapp/wails/v3/internal/signal"
	"gopkg.in/yaml.v3"
)

func ensureIgnored(list *[]string, pattern string) {
	for _, item := range *list {
		if item == pattern {
			return
		}
	}
	*list = append(*list, pattern)
}

type WatcherOptions struct {
	Config string `description:"The config file including path" default:"."`
	TUI    bool   `description:"Run dev mode inside an interactive TUI"`
}

// devConfig mirrors the dev_mode section of build/config.yml.
type devConfig struct {
	Config engine.Config `yaml:"dev_mode"`
}

// parseDevConfig reads and unmarshals the dev_mode engine config, applying the
// shared defaults used by both the plain and TUI watchers.
func parseDevConfig(path string) (engine.Config, error) {
	var dc devConfig
	c, err := os.ReadFile(path)
	if err != nil {
		return engine.Config{}, err
	}
	if err := yaml.Unmarshal(c, &dc); err != nil {
		return engine.Config{}, err
	}
	ensureIgnored(&dc.Config.Ignore.File, "*_test.go")
	return dc.Config, nil
}

func Watcher(options *WatcherOptions) error {
	cfg, err := parseDevConfig(options.Config)
	if err != nil {
		return err
	}

	if options.TUI {
		return watchTUI(cfg)
	}

	return watchPlain(cfg)
}

// watchPlain is the original, unchanged dev watcher: the engine owns OS signals
// and prints process output straight to the terminal.
func watchPlain(cfg engine.Config) error {
	stopChan := make(chan struct{})

	watcherEngine, err := engine.NewEngineFromConfig(cfg)
	if err != nil {
		return err
	}

	// Setup cleanup function that stops the engine
	cleanup := func() {
		watcherEngine.Stop()
	}
	defer cleanup()

	// Signal handler needs to notify when to stop
	signalCleanup := func() {
		cleanup()
		stopChan <- struct{}{}
	}

	signalHandler := signal.NewSignalHandler(signalCleanup)
	signalHandler.ExitMessage = func(sig os.Signal) string {
		return ""
	}
	signalHandler.Start()

	// Start the engine
	err = watcherEngine.Start()
	if err != nil {
		return err
	}

	<-stopChan
	return nil
}

// watchTUI runs the engine via the embedding SDK and drives the dado TUI: the
// engine output and lifecycle events feed the processes view, and the TUI owns
// the terminal (so the engine installs no signal traps — Run(ctx) is used).
func watchTUI(cfg engine.Config) error {
	// Suppress the clir footer; the TUI owns the screen.
	DisableFooter = true

	// Opt the built app into exposing its IPC monitor socket so the bindings /
	// events / calls views can attach. Children inherit the dev-process env.
	_ = os.Setenv("WAILS_MONITOR", "1")

	// Silence engine log chatter; the TUI renders its own UI.
	cfg.LogLevel = "mute"

	store := tui.NewProcStore()
	cfg.Output = store.WriterFor
	cfg.OnProcessEvent = store.HandleEvent

	watcherEngine, err := engine.NewEngineFromConfig(cfg)
	if err != nil {
		return err
	}

	// Ctrl+C also tears everything down cleanly.
	ctx, stop := ossignal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Engine supervises in the background; cancelling ctx triggers a clean
	// shutdown of every process.
	engCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() { _ = watcherEngine.Run(engCtx) }()

	// Block on the TUI; on quit, cancel the engine context.
	return tui.RunDev(ctx, watcherEngine, store)
}
