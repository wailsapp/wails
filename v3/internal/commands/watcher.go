package commands

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"

	"github.com/atterpac/refresh/engine"
	"github.com/atterpac/refresh/process"
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

func ensurePrimaryExitPolicy(executes []process.Execute) {
	for index := range executes {
		execute := &executes[index]
		if execute.Type == process.Primary && execute.ExitPolicy == "" {
			execute.ExitPolicy = process.ExitPolicyShutdown
		}
	}
}

func isInterruptError(err error) bool {
	var exitErr *exec.ExitError
	return errors.As(err, &exitErr) && isInterruptProcessState(exitErr.ProcessState)
}

type WatcherOptions struct {
	Config string `description:"The config file including path" default:"."`
}

func Watcher(options *WatcherOptions) error {
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

	ensureIgnored(&devconfig.Config.Ignore.File, "*_test.go")
	ensurePrimaryExitPolicy(devconfig.Config.ExecStruct)

	watcherEngine, err := engine.NewEngineFromConfig(devconfig.Config)
	if err != nil {
		return err
	}
	if err := watcherEngine.Start(); err != nil {
		if isInterruptError(err) {
			slog.Warn("graceful exit requested", "signal", os.Interrupt)
			return nil
		}
		return err
	}
	return nil
}
