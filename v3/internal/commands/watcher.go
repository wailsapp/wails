package commands

import (
	"github.com/atterpac/refresh/engine"
	"github.com/wailsapp/wails/v3/internal/signal"
	"gopkg.in/yaml.v3"
	"os"
)

type WatcherOptions struct {
	Config string `description:"The config file including path" default:"."`
}

func Watcher(options *WatcherOptions) error {
	stopChan := make(chan struct{})

	// Parse the config file
	type devConfig struct {
		Config engine.Config `yaml:"dev_mode"`
	}

	var devconfig devConfig

	// Parse the config file
	c, err := os.ReadFile(options.Config)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(c, &devconfig)
	if err != nil {
		return err
	}

	watcherEngine, err := engine.NewEngineFromConfig(devconfig.Config)
	if err != nil {
		return err
	}
	signalHandler := signal.NewSignalHandler(func() {
		stopChan <- struct{}{}
	})
	signalHandler.ExitMessage = func(sig os.Signal) string {
		return ""
	}
	signalHandler.Start()
	err = watcherEngine.Start()
	if err != nil {
		return err
	}
	<-stopChan
	return nil
}
