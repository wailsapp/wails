package commands

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/atterpac/refresh/engine"
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
	AppArgs string `name:"appargs" description:"Extra app arguments"`
	Config  string `description:"The config file including path" default:"."`
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

	ensureIgnored(&devconfig.Config.Ignore.File, "*_test.go")

	// If we have application arguments, append them to the wails task
	if options.AppArgs != "" {
		for idx, execAction := range devconfig.Config.ExecStruct {
			if execAction.Cmd == "wails3 task run" {
				clonedExec := execAction

				escapedAppArgs := options.AppArgs
				if runtime.GOOS == "windows" {
					escapedAppArgs = fmt.Sprintf("%q", options.AppArgs)
				} else {
					escapedAppArgs =
						"'" + strings.ReplaceAll(options.AppArgs, "'", `'"'"'`) + "'"
				}

				clonedExec.Cmd = fmt.Sprintf(
					"wails3 task run -appargs=%s",
					escapedAppArgs,
				)
				devconfig.Config.ExecStruct[idx] = clonedExec
				fmt.Printf("Executing %s\n", clonedExec.Cmd)
			}
		}
	}

	watcherEngine, err := engine.NewEngineFromConfig(devconfig.Config)
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
