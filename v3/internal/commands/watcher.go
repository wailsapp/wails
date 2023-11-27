package commands

import (
	"github.com/atterpac/refresh/engine"
	"strings"
)

type WatcherOptions struct {
	Path             string `description:"The path to watch" default:"."`
	PreExec          string `description:"The command to run before the main command"`
	Exec             string `description:"The command to run when a change is detected"`
	PostExec         string `description:"The command to run after the main command"`
	IgnoreFiles      string `description:"The files to ignore (comma separated)"`
	IgnoreDirs       string `description:"The directories to ignore (comma separated)"`
	IgnoreExtensions string `description:"The extensions to ignore (comma separated)"`
	Debounce         int    `description:"The debounce time in milliseconds" default:"1000"`
	PreWait          bool   `description:"Wait for the pre-exec command to finish before running the main command"`
}

func Watcher(options *WatcherOptions) error {

	ignore := engine.Ignore{
		File:      strings.Split(options.IgnoreFiles, ","),
		Dir:       strings.Split(options.IgnoreDirs, ","),
		Extension: strings.Split(options.IgnoreExtensions, ","),
	}
	config := engine.Config{
		RootPath:    options.Path,
		PreExec:     options.PreExec,
		ExecCommand: options.Exec,
		PostExec:    options.PostExec,
		Ignore:      ignore,
		LogLevel:    "info",
		Debounce:    options.Debounce,
		PreWait:     options.PreWait,
	}

	watch := engine.NewEngineFromConfig(config)

	watch.Start()
	<-make(chan struct{})
	return nil
}
