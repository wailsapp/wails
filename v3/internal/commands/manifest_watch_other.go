//go:build !windows

package commands

import "github.com/rjeczalik/notify"

func watchManifestDirectory(path string, events chan<- notify.EventInfo) error {
	return notify.Watch(path, events, notify.All)
}

func stopManifestWatchChannel(events chan<- notify.EventInfo) { notify.Stop(events) }
