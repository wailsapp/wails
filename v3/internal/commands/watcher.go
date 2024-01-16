package commands

import (
	"github.com/atterpac/refresh/engine"
)

type WatcherOptions struct {
	Config string `description:"The config file including path" default:"."`
}

func Watcher(options *WatcherOptions) error {
	engine.NewEngineFromTOML(options.Config).Start()
	<-make(chan struct{})
	return nil
}
