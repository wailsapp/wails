package wails

import homedir "github.com/mitchellh/go-homedir"

// RuntimeFileSystem exposes file system utilities to the runtime
type RuntimeFileSystem struct {
}

func newRuntimeFileSystem() *RuntimeFileSystem {
	return &RuntimeFileSystem{}
}

// HomeDir returns the user's home directory
func (r *RuntimeFileSystem) HomeDir() (string, error) {
	return homedir.Dir()
}
