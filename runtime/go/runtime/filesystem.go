package runtime

import homedir "github.com/mitchellh/go-homedir"

// FileSystem exposes file system utilities to the runtime
type FileSystem struct {}

// Creates a new FileSystem struct
func newFileSystem() *FileSystem {
	return &FileSystem{}
}

// HomeDir returns the user's home directory
func (r *FileSystem) HomeDir() (string, error) {
	return homedir.Dir()
}
