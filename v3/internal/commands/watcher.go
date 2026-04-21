package commands

import (
	"os"

	"github.com/atterpac/refresh/engine"
	"github.com/wailsapp/wails/v3/internal/signal"
	"gopkg.in/yaml.v3"
)

type WatcherOptions struct {
	Mode   string `description:"Whether standard dev mode or using delve" default:"dev"`
	Config string `description:"The config file including path" default:"."`
}

func Watcher(options *WatcherOptions) error {
	stopChan := make(chan struct{})

	// Parse the config file
	c, err := os.ReadFile(options.Config)
	if err != nil {
		return err
	}

	var usedConfig engine.Config

	switch options.Mode {
	case "dev":
		type devConfig struct {
			Config engine.Config `yaml:"dev_mode"`
		}
		var transitoryConfig devConfig
		err = yaml.Unmarshal(c, &transitoryConfig)
		if err != nil {
			return err
		}
		usedConfig = transitoryConfig.Config

	case "debug":
		type debugConfig struct {
			Config engine.Config `yaml:"debug_mode"`
		}
		var transitoryConfig debugConfig
		err = yaml.Unmarshal(c, &transitoryConfig)
		if err != nil {
			return err
		}
		usedConfig = transitoryConfig.Config
	}

	watcherEngine, err := engine.NewEngineFromConfig(usedConfig)
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
